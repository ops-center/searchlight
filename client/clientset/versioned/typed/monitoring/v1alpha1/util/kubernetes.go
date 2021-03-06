package util

import (
	"errors"
	"fmt"
	"reflect"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kmodules.xyz/client-go/meta"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return api.SchemeGroupVersion.WithKind(meta.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
	}

	switch u := v.(type) {
	case *api.ClusterAlert:
		u.APIVersion = api.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *api.NodeAlert:
		u.APIVersion = api.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *api.PodAlert:
		u.APIVersion = api.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
