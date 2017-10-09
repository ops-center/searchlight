package install

import (
	aci "github.com/appscode/searchlight/apis/monitoring"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Install registers the API group and adds types to a scheme
func Install(groupFactoryRegistry announced.APIGroupFactoryRegistry, registry *registered.APIRegistrationManager, scheme *runtime.Scheme) {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  aci.GroupName,
			VersionPreferenceOrder:     []string{api.SchemeGroupVersion.Version},
			ImportPrefix:               "gitub.com/appscode/searchlight/apis/monitoring",
			RootScopedKinds:            sets.NewString("CustomResourceDefinition"),
			AddInternalObjectsToScheme: aci.AddToScheme,
		},
		announced.VersionToSchemeFunc{
			api.SchemeGroupVersion.Version: api.AddToScheme,
		},
	).Announce(groupFactoryRegistry).RegisterAndEnable(registry, scheme); err != nil {
		panic(err)
	}
}
