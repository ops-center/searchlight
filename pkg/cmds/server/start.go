package server

import (
	"fmt"
	"io"
	"net"

	"github.com/appscode/go/log/golog"
	incidentsv1alpha1 "github.com/appscode/searchlight/apis/incidents/v1alpha1"
	"github.com/appscode/searchlight/pkg/operator"
	"github.com/appscode/searchlight/pkg/server"
	_ "github.com/go-openapi/loads"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
)

const defaultEtcdPathPrefix = "/registry/monitoring.appscode.com"

type SearchlightOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
	OperatorOptions    *OperatorOptions

	StdOut io.Writer
	StdErr io.Writer
}

func NewSearchlightOptions(out, errOut io.Writer) *SearchlightOptions {
	o := &SearchlightOptions{
		// TODO we will nil out the etcd storage options.  This requires a later level of k8s.io/apiserver
		RecommendedOptions: genericoptions.NewRecommendedOptions(defaultEtcdPathPrefix, server.Codecs.LegacyCodec(admissionv1beta1.SchemeGroupVersion)),
		OperatorOptions:    NewOperatorOptions(),
		StdOut:             out,
		StdErr:             errOut,
	}
	o.RecommendedOptions.Etcd = nil

	return o
}

func (o SearchlightOptions) AddFlags(fs *pflag.FlagSet) {
	o.RecommendedOptions.AddFlags(fs)
	o.OperatorOptions.AddFlags(fs)
}

func (o SearchlightOptions) Validate(args []string) error {
	return nil
}

func (o *SearchlightOptions) Complete(cmd *cobra.Command) error {
	o.OperatorOptions.verbosity = golog.ParseFlags(cmd.Flags()).Verbosity
	return nil
}

func (o SearchlightOptions) Config() (*server.SearchlightConfig, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewRecommendedConfig(server.Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig, server.Scheme); err != nil {
		return nil, err
	}
	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(incidentsv1alpha1.GetOpenAPIDefinitions, server.Scheme)
	serverConfig.OpenAPIConfig.Info.Title = "searchlight-server"
	serverConfig.OpenAPIConfig.Info.Version = incidentsv1alpha1.SchemeGroupVersion.Version
	serverConfig.OpenAPIConfig.IgnorePrefixes = []string{
		"/swaggerapi",
		"/apis/admission.monitoring.appscode.com/v1alpha1/admissionreviews",
	}
	serverConfig.SwaggerConfig = genericapiserver.DefaultSwaggerConfig()
	serverConfig.EnableMetrics = true

	controllerConfig := operator.NewOperatorConfig(serverConfig.ClientConfig)
	if err := o.OperatorOptions.ApplyTo(controllerConfig); err != nil {
		return nil, err
	}

	config := &server.SearchlightConfig{
		GenericConfig:  serverConfig,
		OperatorConfig: controllerConfig,
	}
	return config, nil
}

func (o SearchlightOptions) Run(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	s, err := config.Complete().New()
	if err != nil {
		return err
	}

	return s.Run(stopCh)
}
