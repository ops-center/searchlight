package options

import (
	"github.com/spf13/pflag"
)

type Config struct {
	Master           string
	KubeConfig       string
	IcingaSecretName string
}

func NewConfig() *Config {
	return &Config{
		Master:     "",
		KubeConfig: "",
		IcingaSecretName: "",
	}
}

func (s *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.KubeConfig, "kubeconfig", s.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	fs.StringVar(&s.IcingaSecretName, "icinga-secret", s.IcingaSecretName, "Icinga secret name")
}
