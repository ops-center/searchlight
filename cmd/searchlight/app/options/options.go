package options

import (
	"github.com/spf13/pflag"
)

type Config struct {
	Master     string
	KubeConfig string
}

func NewConfig() *Config {
	return &Config{
		Master:     "127.0.0.1:8080",
		KubeConfig: "",
	}
}

func (s *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Master, "master", "127.0.0.1:8080", "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.KubeConfig, "kubeconfig", "", "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
}
