package main

// Version is set via ldflags
var Version string

// Build version is set via ldflags
var Build string

// Build date is set via ldflags
var Date string

type VulnOpts struct {
	Assets         []*Asset        `json:"assets,omitempty" mapstructure:"assets"`
	Report         *VulnOptsReport `json:"report,omitempty" mapstructure:"report"`
	Collector      string          `json:"collector,omitempty" mapstructure:"collector"`
	Async          bool            `json:"async,omitempty" mapstructure:"async"`
	IdDetector     string          `json:"id-detector,omitempty" mapstructure:"id-detector"`
	Incognito      bool            `json:"incognito,omitempty" mapstructure:"incognito"`
	Insecure       bool            `json:"insecure,omitempty" mapstructure:"insecure"`
	Policies       []string        `json:"policies,omitempty" mapstructure:"policies"`
	Sudo           VulnOptsSudo    `json:"sudo,omitempty" mapstructure:"sudo"`
	Output         string          `json:"output,omitempty" mapstructure:"output"`
	ScoreThreshold int             `json:"score_threshold,omitempty" mapstructure:"score_threshold"`
}

type VulnOptsSudo struct {
	Active bool `json:"active,omitempty" mapstructure:"active"`
}

type Asset struct {
	Name         string            `json:"name" mapstructure:"name"`
	Mrn          string            `json:"assetmrn,omitempty" mapstructure:"assetmrn"`
	Connection   string            `json:"connection,omitempty" mapstructure:"connection"`
	IdentityFile string            `json:"identityfile,omitempty" mapstructure:"identityfile"`
	Password     string            `json:"password,omitempty" mapstructure:"password"`
	Annotations  map[string]string `json:"annotations,omitempty" mapstructure:"annotations"`
	Labels       map[string]string `json:"labels,omitempty" mapstructure:"labels"`
}

type VulnOptsReport struct {
	Format string `json:"format,omitempty"`
}
