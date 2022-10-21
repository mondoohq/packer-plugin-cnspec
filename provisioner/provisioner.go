package provisioner

//go:generate go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc mapstructure-to-hcl2 -type Config,SudoConfig
//go:generate packer-sdc struct-markdown

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/adapter"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/mitchellh/mapstructure"
	"go.mondoo.com/packer-plugin-mondoo/version"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	ctx                 interpolate.Context
	// The command to invoke mondoo. Defaults to `mondoo scan`.
	Command string
	// The alias by which the host should be known.
	// Defaults to `default`.
	HostAlias string `mapstructure:"host_alias"`
	// The `user` set for your communicator. Defaults to the `user` set
	// by packer.
	User string `mapstructure:"user"`
	// The port on which to attempt to listen for SSH
	//  connections. This value is a starting point. The provisioner will attempt
	//  listen for SSH connections on the first available of ten ports, starting at
	//  `local_port`. A system-chosen port is used when `local_port` is missing or
	//  empty.
	LocalPort uint `mapstructure:"local_port"`
	// The SSH key that will be used to run the SSH
	//  server on the host machine to forward commands to the target machine.
	//  packer connects to this server and will validate the identity of the
	//  server using the system known_hosts. The default behavior is to generate
	//  and use a onetime key.
	SSHHostKeyFile string `mapstructure:"ssh_host_key_file"`
	// The SSH public key of the packer `ssh_user`.
	// The default behavior is to generate and use a onetime key.
	SSHAuthorizedKeyFile string `mapstructure:"ssh_authorized_key_file"`
	// packer's SFTP proxy is not reliable on some unix/linux systems,
	// therefore we recommend to use scp as default for packer proxy
	UseSFTP bool `mapstructure:"use_sftp"`
	// Sets the log level to `DEBUG`
	Debug bool `mapstructure:"debug"`
	// The asset name passed to Mondoo Platform. Defaults to the hostname
	// of the instance.
	AssetName string `mapstructure:"asset_name"`
	// Array of environment variables for configuring Mondoo.
	MondooEnvVars []string `mapstructure:"mondoo_env_vars"`
	// Configure behavior whether packer should fail if `scan_threshold` is
	// not met. If `scan_threshold` configuration is omitted, the threshold
	// is set to `0` and builds will pass regardless of what score is
	// returned.
	// If `score_threshold` is set to a value, and `on_failure = "continue"`
	// builds will continue regardless of what score is returned.
	OnFailure string `mapstructure:"on_failure"`
	// Configure an optional map of labels for the asset data in Mondoo Platform.
	Labels map[string]string `mapstructure:"labels"`
	// Configure an optional map of `key/val` annotations for the asset data in
	// Mondoo Platform.
	Annotations map[string]string `mapstructure:"annotations"`
	// Configures incognito mode. Defaults to `true`. When set to false, scan results
	// will not be sent to the Mondoo platform.
	Incognito bool `mapstructure:"incognito"`
	// A list of policies to be executed (requires incognito mode).
	Policies []string `mapstructure:"policies"`
	// A path to local policy bundle file.
	PolicyBundle string `mapstructure:"policybundle"`
	// Run mondoo scan with `--sudo`. Defaults to none.
	Sudo *SudoConfig `mapstructure:"sudo"`
	// Configure WinRM user. Defaults to `user` set by the packer communicator.
	WinRMUser string `mapstructure:"winrm_user"`
	// Configure WinRM user password. Defaults to `password` set by the packer communicator.
	WinRMPassword string `mapstructure:"winrm_password"`
	// Use proxy to connect to host to scan. This configuration will fall-back to packer proxy
	// for cases where the provisioner cannot access the target directly
	// NOTE: we have seen cases with the vsphere builder
	UseProxy bool `mapstructure:"use_proxy"`
	// Set output format: summary, full, yaml, json, csv, compact, report, junit (default "compact")
	Output string `mapstructure:"output"`
	// An integer value to set the `score_threshold` of mondoo scans. Defaults to `0` which results in
	// a passing score regardless of what scan results are returned.
	ScoreThreshold int `mapstructure:"score_threshold"`
	// The path to the mondoo client config. Defaults to `$HOME/.config/mondoo/mondoo.yml`
	MondooConfigPath string `mapstructure:"mondoo_config_path"`
}

type SudoConfig struct {
	Active bool `mapstructure:"active"`
}

func validateFileConfig(name string, config string, req bool) error {
	if req {
		if name == "" {
			return fmt.Errorf("%s must be specified", config)
		}
	}
	info, err := os.Stat(name)
	if err != nil {
		return fmt.Errorf("%s: %s is invalid: %s", config, name, err)
	} else if info.IsDir() {
		return fmt.Errorf("%s: %s must point to a file", config, name)
	}
	return nil
}

type Provisioner struct {
	config             Config
	buildInfo          BuildInfo
	adapter            *adapter.Adapter
	adapterPrivKeyFile string
	done               chan struct{}
}

func (p *Provisioner) Prepare(raws ...interface{}) error {
	p.done = make(chan struct{})

	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	if p.config.Command == "" {
		p.config.Command = "mondoo"
	}

	var errs *packer.MultiError
	if len(p.config.SSHAuthorizedKeyFile) > 0 {
		err = validateFileConfig(p.config.SSHAuthorizedKeyFile, "ssh_authorized_key_file", true)
		if err != nil {
			log.Println(p.config.SSHAuthorizedKeyFile, "does not exist")
			errs = packer.MultiErrorAppend(errs, err)
		}
	}

	// ensure that we disable ssh auth, since the packer proxy only allows one auth mechanism
	p.config.MondooEnvVars = append(p.config.MondooEnvVars, "SSH_AUTH_SOCK=")

	if !p.config.UseSFTP {
		p.config.MondooEnvVars = append(p.config.MondooEnvVars, "MONDOO_SSH_SCP=on")
	}

	if p.config.Debug {
		p.config.MondooEnvVars = append(p.config.MondooEnvVars, "DEBUG=1")
	}

	if p.config.LocalPort > 65535 {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("local_port: %d must be a valid port", p.config.LocalPort))
	}

	if p.config.User == "" {
		usr, err := user.Current()
		if err != nil {
			errs = packer.MultiErrorAppend(errs, err)
		} else {
			p.config.User = usr.Username
		}
	}
	if p.config.User == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("user: could not determine current user from environment"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *Provisioner) Provision(ctx context.Context, ui packer.Ui, comm packer.Communicator, generatedData map[string]interface{}) error {
	ui.Say("Running Mondoo packer provisioner (Version: " + version.Version + ", Build: " + version.Build + ")")

	err := mapstructure.Decode(generatedData, &p.buildInfo)
	if err != nil {
		ui.Error(fmt.Sprintf("could not decode packer meta information:\n%v", generatedData))
		return err
	}

	if p.config.Debug {
		data, err := json.Marshal(p.buildInfo)
		if err != nil {
			return err
		}
		ui.Say(string(data))
	}

	// configure ssh proxy
	if p.config.UseProxy {
		ui.Say("configure ssh proxy")
		k, err := newUserKey(p.config.SSHAuthorizedKeyFile)
		if err != nil {
			return err
		}
		p.adapterPrivKeyFile = k.privKeyFile

		hostSigner, err := newSigner(p.config.SSHHostKeyFile)
		if err != nil {
			return err
		}

		// Remove the private key file when we're done with scanning
		if len(k.privKeyFile) > 0 {
			defer os.Remove(k.privKeyFile)
		}

		keyChecker := ssh.CertChecker{
			UserKeyFallback: func(conn ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
				if user := conn.User(); user != p.config.User {
					return nil, fmt.Errorf("authentication failed: %s is not a valid user", user)
				}

				if !bytes.Equal(k.Marshal(), pubKey.Marshal()) {
					return nil, errors.New("authentication failed: unauthorized key")
				}

				return nil, nil
			},
		}

		config := &ssh.ServerConfig{
			AuthLogCallback: func(conn ssh.ConnMetadata, method string, err error) {
				log.Printf("ssh proxy authentication attempt from %s to %s as %s using %s", conn.RemoteAddr(), conn.LocalAddr(), conn.User(), method)
			},
			PublicKeyCallback: keyChecker.Authenticate,
		}

		config.AddHostKey(hostSigner)

		localListener, err := func() (net.Listener, error) {

			port := p.config.LocalPort
			tries := 1
			if port != 0 {
				tries = 10
			}
			for i := 0; i < tries; i++ {
				l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
				port++
				if err != nil {
					ui.Say(err.Error())
					continue
				}
				_, portStr, err := net.SplitHostPort(l.Addr().String())
				if err != nil {
					ui.Say(err.Error())
					continue
				}
				portUint64, err := strconv.ParseUint(portStr, 10, 0)
				if err != nil {
					ui.Say(err.Error())
					continue
				}
				p.config.LocalPort = uint(portUint64)
				return l, nil
			}
			return nil, errors.New("error setting up SSH proxy connection")
		}()

		if err != nil {
			return err
		}

		// initialize ssh adapter
		p.adapter = adapter.NewAdapter(p.done, localListener, config, "sftp -e", ui, comm)
		defer func() {
			ui.Say("shutting down the SSH proxy")
			close(p.done)
			p.adapter.Shutdown()
		}()
		go p.adapter.Serve()
	}

	ui = &packer.SafeUi{
		Sem: make(chan int, 1),
		Ui:  ui,
	}

	// run mondoo policies
	err = p.executeMondoo(ctx, ui, comm)
	if err != nil {
		ui.Error(err.Error())
	}

	// NOTE: if we got an error but user set the continue option, we do not error the execution
	if err != nil && p.config.OnFailure != "continue" {
		return err
	}

	return nil
}

// Cancel just exists when provision is cancelled
func (p *Provisioner) Cancel() {
	if p.done != nil {
		close(p.done)
	}
	if p.adapter != nil {
		p.adapter.Shutdown()
	}
	os.Exit(0)
}

func (p *Provisioner) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

func (p *Provisioner) executeMondoo(ctx context.Context, ui packer.Ui, comm packer.Communicator) error {
	var envvars []string

	if len(p.config.MondooEnvVars) > 0 {
		envvars = append(envvars, p.config.MondooEnvVars...)
	}

	// Always available Packer provided env vars
	p.config.MondooEnvVars = append(p.config.MondooEnvVars, fmt.Sprintf("PACKER_BUILD_NAME=%s", p.config.PackerBuildName))
	p.config.MondooEnvVars = append(p.config.MondooEnvVars, fmt.Sprintf("PACKER_BUILDER_TYPE=%s", p.config.PackerBuilderType))

	cmdargs := []string{"scan"}

	connType := "local"
	var endpoint string
	var user string
	var password string
	var privKeyFile string

	if p.buildInfo.ConnType == "" || p.buildInfo.ConnType == "ssh" {
		connType = "ssh"
		endpoint = fmt.Sprintf("%s:%d", p.buildInfo.Host, p.buildInfo.Port)
		user = p.buildInfo.User
		password = p.buildInfo.Password
		// if we get a private key, cache that key locally
		if len(p.buildInfo.SSHPrivateKey) > 0 {
			tmpfile, err := ioutil.TempFile("", "packer")
			if err != nil {
				return err
			}
			// clean up ssh key after scan
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(p.buildInfo.SSHPrivateKey)); err != nil {
				return err
			}
			if err := tmpfile.Close(); err != nil {
				return err
			}
			privKeyFile = tmpfile.Name()
		}

		// use proxy
		if p.config.UseProxy {
			ui.Say("use packer's ssh proxy")
			endpoint = fmt.Sprintf("%s:%d", "127.0.0.1", p.config.LocalPort)
			privKeyFile = p.adapterPrivKeyFile
			user = p.config.User
			password = ""
		}
	} else if p.buildInfo.ConnType == "winrm" {
		connType = "winrm"
		endpoint = fmt.Sprintf("%s:%d", p.buildInfo.Host, p.buildInfo.Port)
		user = p.buildInfo.User
		password = p.buildInfo.Password
	} else if p.buildInfo.ConnType == "docker" {
		connType = "docker"
		endpoint = fmt.Sprintf("%s", p.buildInfo.ID)
		ui.Say(endpoint)
	} else {
		return errors.New("unsupported connection type: " + p.buildInfo.ConnType)
	}
	// mondoo scan local or mondoo scan ssh ec2-user@3.219.56.31 or mondoo scan winrm ec2-user@3.219.56.31
	cmdargs = append(cmdargs, connType)
	if connType == "ssh" || connType == "winrm" {
		cmdargs = append(cmdargs, fmt.Sprintf("%s@%s", user, endpoint))
	}

	if connType == "docker" {
		cmdargs = append(cmdargs, endpoint)
	}

	if p.config.PolicyBundle != "" {
		cmdargs = append(cmdargs, "--policy-bundle "+p.config.PolicyBundle)
	}

	if p.config.Output != "" {
		cmdargs = append(cmdargs, []string{"--output", p.config.Output}...)
	}

	if password != "" {
		cmdargs = append(cmdargs, []string{"--password", password}...)
	}

	if privKeyFile != "" {
		cmdargs = append(cmdargs, []string{"--identity-file", privKeyFile}...)
	}

	if p.config.OnFailure == "continue" {
		// ignore the result of the scan
		cmdargs = append(cmdargs, []string{"--score-threshold", strconv.Itoa(0)}...)
	} else if p.config.ScoreThreshold != 0 {
		// user overwrite the default score threshold
		cmdargs = append(cmdargs, []string{"--score-threshold", strconv.Itoa(p.config.ScoreThreshold)}...)
	} else {
		// expects all controls to pass
		cmdargs = append(cmdargs, []string{"--score-threshold", strconv.Itoa(100)}...)
	}

	if p.config.MondooConfigPath != "" {
		cmdargs = append(cmdargs, []string{"--config", p.config.MondooConfigPath}...)
	}

	// If annotations are not specified, this will error out so make sure to init the map.
	if p.config.Annotations == nil {
		p.config.Annotations = map[string]string{}
	}

	// labels are deprecated, therefore we merge them with annotations
	for k := range p.config.Labels {
		p.config.Annotations[k] = p.config.Labels[k]
	}
	// build configuration
	connection := fmt.Sprintf("%s://%s", connType, endpoint)
	if user != "" {
		connection = fmt.Sprintf("%s://%s@%s", connType, user, endpoint)
	}

	conf := &VulnOpts{
		Assets: []*Asset{
			{
				Name:         p.config.AssetName,
				Connection:   connection,
				IdentityFile: privKeyFile,
				Password:     password,
				Annotations:  p.config.Annotations,
			},
		},
		Insecure: true, // we do not check the hostkey for the packer build
	}

	if p.config.Sudo != nil && p.config.Sudo.Active {
		ui.Say("activated sudo")
		conf.Sudo = VulnOptsSudo{
			Active: p.config.Sudo.Active,
		}
	}

	// pass incognito to mondoo scan
	conf.Incognito = p.config.Incognito

	// pass policies into mondoo config
	conf.Policies = p.config.Policies

	// prep config for mondoo executable
	mondooScanConf, err := json.Marshal(conf)

	if err != nil {
		return err
	}

	if p.config.Debug {
		ui.Say(fmt.Sprintf("mondoo configuration: %v", string(mondooScanConf)))
	}

	cmd := exec.Command(p.config.Command, cmdargs...)

	cmd.Env = os.Environ()
	if len(envvars) > 0 {
		cmd.Env = append(cmd.Env, envvars...)
	}
	cmd.Env = append(cmd.Env, "CI=true")
	cmd.Env = append(cmd.Env, "PACKER_PIPELINE=true")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdin.Write(mondooScanConf)
	stdin.Close()

	wg := sync.WaitGroup{}
	repeat := func(r io.ReadCloser) {
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadString('\n')
			if line != "" {
				line = strings.TrimRightFunc(line, unicode.IsSpace)
				ui.Message(line)
			}
			if err != nil {
				if err == io.EOF {
					break
				} else {
					ui.Error(err.Error())
					break
				}
			}
		}
		wg.Done()
	}
	wg.Add(2)
	go repeat(stdout)
	go repeat(stderr)

	ui.Say(fmt.Sprintf("Executing Mondoo: %s", cmd.Args))
	if err := cmd.Start(); err != nil {
		return err
	}
	wg.Wait()
	err = cmd.Wait()

	if err != nil {
		return fmt.Errorf("non-zero exit status: %s", err)
	}

	return nil
}
