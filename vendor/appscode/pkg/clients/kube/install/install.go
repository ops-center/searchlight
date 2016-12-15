package install

import (
	"appscode/pkg/clients/kube"

	"k8s.io/kubernetes/pkg/apimachinery/announced"
	"k8s.io/kubernetes/pkg/util/sets"
)

func init() {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  kube.GroupName,
			VersionPreferenceOrder:     []string{kube.V1beta1SchemeGroupVersion.Version},
			ImportPrefix:               "appscode/pkg/clients/kube",
			RootScopedKinds:            sets.NewString("PodSecurityPolicy", "ThirdPartyResource"),
			AddInternalObjectsToScheme: kube.AddToScheme,
		},
		announced.VersionToSchemeFunc{
			kube.V1beta1SchemeGroupVersion.Version: kube.V1betaAddToScheme,
		},
	).Announce().RegisterAndEnable(); err != nil {
		panic(err)
	}
}
