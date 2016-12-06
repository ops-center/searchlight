package client

import (
	"appscode/pkg/clients/kube"
	"encoding/json"
	"io"
	"net/url"
	"reflect"
	"strings"

	"github.com/appscode/log"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/runtime/serializer"
)

// TODO(@sadlil): Find a better way to replace ExtendedCodec to encode and decode objects.
// Follow the guide to replace it with api.Codec and api.ParameterCodecs.
var ExtendedCodec = &extendedCodec{}

func ContentConfig() restclient.ContentConfig {
	streamingInfo, _ := api.Codecs.StreamingSerializerForMediaType("application/json;stream=watch", nil)
	return restclient.ContentConfig{
		AcceptContentTypes: runtime.ContentTypeJSON,
		ContentType:        runtime.ContentTypeJSON,
		GroupVersion: &unversioned.GroupVersion{
			Group:   "appscode.com",
			Version: "v1",
		},
		NegotiatedSerializer: serializer.NegotiatedSerializerWrapper(runtime.SerializerInfo{Serializer: ExtendedCodec}, streamingInfo),
	}
}

type extendedCodec struct{}

func (*extendedCodec) Decode(data []byte, gvk *unversioned.GroupVersionKind, obj runtime.Object) (runtime.Object, *unversioned.GroupVersionKind, error) {
	if obj == nil {
		metadata := &unversioned.TypeMeta{}
		err := json.Unmarshal(data, metadata)
		if err != nil {
			return obj, gvk, err
		}
		log.V(7).Infoln("Detected metadata type for nil object, got", metadata.APIVersion, metadata.Kind)
		obj, err = setDefaultType(metadata)
		if err != nil {
			log.Errorln("faild to create type", err)
		}
	}
	err := json.Unmarshal(data, obj)
	if err != nil {
		return obj, gvk, err
	}
	return obj, gvk, nil
}

func (*extendedCodec) Encode(obj runtime.Object, w io.Writer) error {
	setDefaultVersionKind(obj)
	return json.NewEncoder(w).Encode(obj)
}

// DecodeParameters converts the provided url.Values into an object of type From with the kind of into, and then
// converts that object to into (if necessary). Returns an error if the operation cannot be completed.
func (*extendedCodec) DecodeParameters(parameters url.Values, from unversioned.GroupVersion, into runtime.Object) error {
	if len(parameters) == 0 {
		return nil
	}
	_, okDelete := into.(*api.DeleteOptions)
	if _, okList := into.(*api.ListOptions); okList || okDelete {
		from = unversioned.GroupVersion{Version: "v1"}
	}
	return runtime.NewParameterCodec(api.Scheme).DecodeParameters(parameters, from, into)
}

// EncodeParameters converts the provided object into the to version, then converts that object to url.Values.
// Returns an error if conversion is not possible.
func (c *extendedCodec) EncodeParameters(obj runtime.Object, to unversioned.GroupVersion) (url.Values, error) {
	result := url.Values{}
	if obj == nil {
		return result, nil
	}
	_, okDelete := obj.(*api.DeleteOptions)
	if _, okList := obj.(*api.ListOptions); okList || okDelete {
		to = unversioned.GroupVersion{Version: "v1"}
	}
	return runtime.NewParameterCodec(api.Scheme).EncodeParameters(obj, to)
}

func setDefaultVersionKind(obj runtime.Object) {
	// Check the values can are In type Extended Ingress
	defaultGVK := unversioned.GroupVersionKind{
		Group:   kube.V1beta1SchemeGroupVersion.Group,
		Version: kube.V1beta1SchemeGroupVersion.Version,
	}

	fullyQualifiedKind := reflect.ValueOf(obj).Type().String()
	lastIndexOfDot := strings.LastIndex(fullyQualifiedKind, ".")
	if lastIndexOfDot > 0 {
		defaultGVK.Kind = fullyQualifiedKind[lastIndexOfDot+1:]
	}

	obj.GetObjectKind().SetGroupVersionKind(defaultGVK)
}

func setDefaultType(metadata *unversioned.TypeMeta) (runtime.Object, error) {
	return api.Scheme.New(metadata.GroupVersionKind())
}
