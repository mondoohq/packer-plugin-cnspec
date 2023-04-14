package provisioner

//go:generate go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc mapstructure-to-hcl2 -type Config,SudoConfig
//go:generate packer-sdc struct-markdown

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/adapter"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	config_loader "go.mondoo.com/cnquery/cli/config"
	"go.mondoo.com/cnquery/logger"
	"go.mondoo.com/cnquery/motor/asset"
	inventory "go.mondoo.com/cnquery/motor/inventory/v1"
	"go.mondoo.com/cnquery/motor/providers"
	"go.mondoo.com/cnquery/motor/vault"
	"go.mondoo.com/cnquery/upstream"
	cnspec_config "go.mondoo.com/cnspec/apps/cnspec/cmd/config"
	"go.mondoo.com/cnspec/cli/reporter"
	"go.mondoo.com/cnspec/policy"
	"go.mondoo.com/cnspec/policy/scan"
	"go.mondoo.com/packer-plugin-cnspec/version"
	"go.mondoo.com/ranger-rpc"
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

	if p.config.LocalPort > 65535 {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("local_port: %d must be a valid port", p.config.LocalPort))
	}

	if p.config.User == "" {
		userName, err := user.Current()
		if err != nil {
			errs = packer.MultiErrorAppend(errs, err)
		} else {
			p.config.User = userName.Username
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
	ui.Say("Running cnspec packer provisioner by Mondoo (Version: " + version.Version + ", Build: " + version.Build + ")")

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
		ui.Message("build info: " + string(data))
	}

	// configure ssh proxy
	if p.config.UseProxy {
		ui.Message("configure ssh proxy")
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
				if userName := conn.User(); userName != p.config.User {
					return nil, fmt.Errorf("authentication failed: %s is not a valid user", userName)
				}

				if !bytes.Equal(k.Marshal(), pubKey.Marshal()) {
					return nil, errors.New("authentication failed: unauthorized key")
				}

				return nil, nil
			},
		}

		proxyConfig := &ssh.ServerConfig{
			AuthLogCallback: func(conn ssh.ConnMetadata, method string, err error) {
				log.Printf("ssh proxy authentication attempt from %s to %s as %s using %s", conn.RemoteAddr(), conn.LocalAddr(), conn.User(), method)
			},
			PublicKeyCallback: keyChecker.Authenticate,
		}

		proxyConfig.AddHostKey(hostSigner)

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
		p.adapter = adapter.NewAdapter(p.done, localListener, proxyConfig, "sftp -e", ui, comm)
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

	// run policies
	err = p.executeCnspec(ui, comm)
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

func (p *Provisioner) executeCnspec(ui packer.Ui, comm packer.Communicator) error {

	assetConfig := &providers.Config{
		Backend: providers.ProviderType_UNKNOWN,
		Options: map[string]string{},
	}

	if p.config.Sudo != nil && p.config.Sudo.Active {
		ui.Message("activated sudo")
		assetConfig.Sudo = &providers.Sudo{
			Active: p.config.Sudo.Active,
		}
	}

	if p.buildInfo.ConnType == "" || p.buildInfo.ConnType == "ssh" {
		ui.Message("detected packer build via ssh")
		assetConfig.Backend = providers.ProviderType_SSH
		assetConfig.Host = p.buildInfo.Host
		assetConfig.Port = int32(p.buildInfo.Port)
		assetConfig.Insecure = true // we do not check the hostkey for the packer build
		assetConfig.Credentials = []*vault.Credential{}

		if !p.config.UseSFTP {
			assetConfig.Options["ssh_scp"] = "on"
		}

		// use proxy
		if p.config.UseProxy {
			ui.Message("use packer's ssh proxy")
			// overwrite host since we go via the proxy now
			assetConfig.Host = "127.0.0.1"
			assetConfig.Port = int32(p.config.LocalPort)

			// NOTE: packer proxy only allows one auth mechanism
			cred, err := vault.NewPrivateKeyCredentialFromPath(p.config.User, p.adapterPrivKeyFile, "")
			if err != nil {
				return errors.Wrap(err, "could not gather private key file for proxy from: "+p.adapterPrivKeyFile)
			}
			assetConfig.Credentials = append(assetConfig.Credentials, cred)
		} else if len(p.buildInfo.SSHPrivateKey) > 0 {
			cred := vault.NewPrivateKeyCredential(p.buildInfo.User, []byte(p.buildInfo.SSHPrivateKey), "")
			assetConfig.Credentials = append(assetConfig.Credentials, cred)
		} else {
			// fallback to password auth
			cred := vault.NewPasswordCredential(p.buildInfo.User, p.buildInfo.Password)
			assetConfig.Credentials = append(assetConfig.Credentials, cred)
		}
	} else if p.buildInfo.ConnType == "winrm" {
		ui.Message("detected packer build via winrm")
		assetConfig.Backend = providers.ProviderType_WINRM
		assetConfig.Host = p.buildInfo.Host
		assetConfig.Port = int32(p.buildInfo.Port)
		assetConfig.Insecure = true // we do not check the hostkey for the packer build
		cred := vault.NewPasswordCredential(p.buildInfo.User, p.buildInfo.Password)
		assetConfig.Credentials = append(assetConfig.Credentials, cred)
	} else if p.buildInfo.ConnType == "docker" {
		ui.Message("detected packer container image build")
		assetConfig.Backend = providers.ProviderType_DOCKER
		// buildInfo.ID containers the docker container image id
		assetConfig.Host = fmt.Sprintf("%s", p.buildInfo.ID)
	} else {
		ui.Message("detected packer build via unknown connection type: " + p.buildInfo.ConnType)
		return errors.New("unsupported connection type: " + p.buildInfo.ConnType)
	}

	var policyBundle *policy.Bundle
	var policyFilters []string

	if p.config.PolicyBundle != "" {
		ui.Message("load policy bundle from: " + p.config.PolicyBundle)
		var err error
		policyBundle, err = policy.BundleFromPaths(p.config.PolicyBundle)
		if err != nil {
			return errors.Wrap(err, "could not load policy bundle from "+p.config.PolicyBundle)
		}
		p.config.Incognito = true
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
	conf := inventory.New(inventory.WithAssets(&asset.Asset{
		Name:        p.config.AssetName,
		Connections: []*providers.Config{assetConfig},
		Annotations: p.config.Annotations,
		Labels: map[string]string{
			"packer.io/buildname": p.config.PackerBuildName,
			"packer.io/buildtype": p.config.PackerBuilderType,
		},
	}))

	scanJob := &scan.Job{
		Inventory:     conf,
		Bundle:        policyBundle,
		PolicyFilters: policyFilters,
		ReportType:    scan.ReportType_FULL,
	}

	debugLogBuffer := &bytes.Buffer{}
	logger.SetWriter(debugLogBuffer)
	if p.config.Debug {
		data, _ := json.Marshal(scanJob)
		ui.Message(fmt.Sprintf("cnspec job configuration: %v", string(data)))

		// configure stderr logger
		logger.Set("debug")

		// log output for debug/error logs
		defer func() {
			ui.Message(debugLogBuffer.String())
		}()
	}

	// base 64 config env setting has always precedence
	viper.SetConfigType("yaml")
	if value := os.Getenv("MONDOO_CONFIG_BASE64"); len(value) > 0 {
		ui.Message("load config from detected MONDOO_CONFIG_BASE64")
		decodedData, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			wErr := errors.Wrap(err, "cannot parse config from MONDOO_CONFIG_BASE64")
			ui.Error(wErr.Error())
			return wErr
		}
		err = viper.ReadConfig(bytes.NewBuffer(decodedData))
		if err != nil {
			return errors.Wrap(err, "could not read config from MONDOO_CONFIG_BASE64")
		}
	} else {
		// load first config we find in the following order:
		// 1. MONDOO_CONFIG_PATH env variable
		// 2. MondooConfigPath from config
		// 3. If no MondooConfigPath was set: home directory & system directory
		paths := []string{}

		if path := os.Getenv("MONDOO_CONFIG_PATH"); len(path) > 0 {
			paths = append(paths, path)
		}

		if p.config.MondooConfigPath != "" {
			paths = append(paths, p.config.MondooConfigPath)
		} else {
			config_loader.AppFs = afero.NewOsFs() // TODO fix in config_loader package, this should not be here
			homeConfig, exists, err := config_loader.HomePath()
			if err == nil && exists {
				paths = append(paths, homeConfig)
			}

			if path, ok := config_loader.SystemPath(); ok {
				paths = append(paths, path)
			}
		}

		foundConfig := false
		for i := range paths {
			path := paths[i]
			if path == "" {
				continue
			}

			_, err := os.Stat(path)
			if err != nil {
				continue
			}

			ui.Message("load config from detected " + path)
			data, err := os.ReadFile(path)
			if err != nil {
				wErr := errors.Wrap(err, "cannot parse config from "+path)
				ui.Error(wErr.Error())
				return wErr
			}
			err = viper.ReadConfig(bytes.NewBuffer(data))
			if err != nil {
				wErr := errors.Wrap(err, "cannot parse config from "+path)
				ui.Error(wErr.Error())
				return wErr
			}

			foundConfig = true
			break
		}

		if !foundConfig {
			ui.Message("no configuration provided")
			p.config.Incognito = true
		}
	}

	var result *scan.ScanResult
	var err error
	if p.config.Incognito {
		ui.Message("scan packer build in incognito mode")
		scanService := scan.NewLocalScanner()
		result, err = scanService.RunIncognito(context.Background(), scanJob)
		if err != nil {
			return err
		}
	} else {
		cfg, err := cnspec_config.ReadConfig()
		if err != nil {
			wErr := errors.Wrap(err, "could not parse cnspec configuration")
			ui.Error(wErr.Error())
			return wErr
		}

		var scannerOpts []scan.ScannerOption
		serviceAccount := cfg.GetServiceCredential()
		if serviceAccount != nil {
			ui.Message("using service account credentials")
			scannerOpts = append(scannerOpts, scan.WithUpstream(cfg.UpstreamApiEndpoint(), cfg.GetParentMrn(), ranger.DefaultHttpClient()))
			certAuth, err := upstream.NewServiceAccountRangerPlugin(serviceAccount)
			if err != nil {
				ui.Error("could not create service account plugin: " + err.Error())
				return err
			}
			plugins := []ranger.ClientPlugin{certAuth}
			scannerOpts = append(scannerOpts, scan.WithPlugins(plugins))
		}

		ui.Message("scan packer build")
		scanService := scan.NewLocalScanner(scannerOpts...)
		result, err = scanService.Run(context.Background(), scanJob)
		if err != nil {
			ui.Error("scan failed: " + err.Error())
			return err
		}
	}

	ui.Message("scan completed successfully")

	// render terminal output
	output := p.config.Output
	r, err := reporter.New(output)
	if err != nil {
		return err
	}
	r.IsIncognito = p.config.Incognito

	fullReport := result.GetFull()
	if fullReport == nil {
		rErr := errors.New("could not gather the full report")
		ui.Error(rErr.Error())
		return rErr
	}
	buf := &bytes.Buffer{}
	err = r.Print(fullReport, buf)
	if err != nil {
		return err
	}
	ui.Message(buf.String())

	// default is to pass all controls
	scoreThreshold := 100
	if p.config.OnFailure == "continue" {
		// ignore the result of the scan
		scoreThreshold = 0
	} else if p.config.ScoreThreshold != 0 {
		// user overwrite the default score threshold
		scoreThreshold = p.config.ScoreThreshold
	}

	if fullReport.GetWorstScore() < uint32(scoreThreshold) {
		return fmt.Errorf("scan has completed with %d score, does not pass score threshold %d", fullReport.GetWorstScore(), scoreThreshold)
	}

	return nil
}
