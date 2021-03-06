// +build !ignore_autogenerated

/*
Copyright 2019 The Searchlight Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAlert) DeepCopyInto(out *ClusterAlert) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAlert.
func (in *ClusterAlert) DeepCopy() *ClusterAlert {
	if in == nil {
		return nil
	}
	out := new(ClusterAlert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterAlert) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAlertList) DeepCopyInto(out *ClusterAlertList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterAlert, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAlertList.
func (in *ClusterAlertList) DeepCopy() *ClusterAlertList {
	if in == nil {
		return nil
	}
	out := new(ClusterAlertList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterAlertList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterAlertSpec) DeepCopyInto(out *ClusterAlertSpec) {
	*out = *in
	out.CheckInterval = in.CheckInterval
	out.AlertInterval = in.AlertInterval
	if in.Receivers != nil {
		in, out := &in.Receivers, &out.Receivers
		*out = make([]Receiver, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Vars != nil {
		in, out := &in.Vars, &out.Vars
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterAlertSpec.
func (in *ClusterAlertSpec) DeepCopy() *ClusterAlertSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterAlertSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Incident) DeepCopyInto(out *Incident) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Incident.
func (in *Incident) DeepCopy() *Incident {
	if in == nil {
		return nil
	}
	out := new(Incident)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Incident) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IncidentList) DeepCopyInto(out *IncidentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Incident, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IncidentList.
func (in *IncidentList) DeepCopy() *IncidentList {
	if in == nil {
		return nil
	}
	out := new(IncidentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IncidentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IncidentNotification) DeepCopyInto(out *IncidentNotification) {
	*out = *in
	if in.Author != nil {
		in, out := &in.Author, &out.Author
		*out = new(string)
		**out = **in
	}
	if in.Comment != nil {
		in, out := &in.Comment, &out.Comment
		*out = new(string)
		**out = **in
	}
	in.FirstTimestamp.DeepCopyInto(&out.FirstTimestamp)
	in.LastTimestamp.DeepCopyInto(&out.LastTimestamp)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IncidentNotification.
func (in *IncidentNotification) DeepCopy() *IncidentNotification {
	if in == nil {
		return nil
	}
	out := new(IncidentNotification)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IncidentStatus) DeepCopyInto(out *IncidentStatus) {
	*out = *in
	if in.Notifications != nil {
		in, out := &in.Notifications, &out.Notifications
		*out = make([]IncidentNotification, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IncidentStatus.
func (in *IncidentStatus) DeepCopy() *IncidentStatus {
	if in == nil {
		return nil
	}
	out := new(IncidentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAlert) DeepCopyInto(out *NodeAlert) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAlert.
func (in *NodeAlert) DeepCopy() *NodeAlert {
	if in == nil {
		return nil
	}
	out := new(NodeAlert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeAlert) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAlertList) DeepCopyInto(out *NodeAlertList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeAlert, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAlertList.
func (in *NodeAlertList) DeepCopy() *NodeAlertList {
	if in == nil {
		return nil
	}
	out := new(NodeAlertList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeAlertList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAlertSpec) DeepCopyInto(out *NodeAlertSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeName != nil {
		in, out := &in.NodeName, &out.NodeName
		*out = new(string)
		**out = **in
	}
	out.CheckInterval = in.CheckInterval
	out.AlertInterval = in.AlertInterval
	if in.Receivers != nil {
		in, out := &in.Receivers, &out.Receivers
		*out = make([]Receiver, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Vars != nil {
		in, out := &in.Vars, &out.Vars
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAlertSpec.
func (in *NodeAlertSpec) DeepCopy() *NodeAlertSpec {
	if in == nil {
		return nil
	}
	out := new(NodeAlertSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PluginArguments) DeepCopyInto(out *PluginArguments) {
	*out = *in
	if in.Vars != nil {
		in, out := &in.Vars, &out.Vars
		*out = new(PluginVars)
		(*in).DeepCopyInto(*out)
	}
	if in.Host != nil {
		in, out := &in.Host, &out.Host
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PluginArguments.
func (in *PluginArguments) DeepCopy() *PluginArguments {
	if in == nil {
		return nil
	}
	out := new(PluginArguments)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PluginVarField) DeepCopyInto(out *PluginVarField) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PluginVarField.
func (in *PluginVarField) DeepCopy() *PluginVarField {
	if in == nil {
		return nil
	}
	out := new(PluginVarField)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PluginVars) DeepCopyInto(out *PluginVars) {
	*out = *in
	if in.Fields != nil {
		in, out := &in.Fields, &out.Fields
		*out = make(map[string]PluginVarField, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Required != nil {
		in, out := &in.Required, &out.Required
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PluginVars.
func (in *PluginVars) DeepCopy() *PluginVars {
	if in == nil {
		return nil
	}
	out := new(PluginVars)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodAlert) DeepCopyInto(out *PodAlert) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodAlert.
func (in *PodAlert) DeepCopy() *PodAlert {
	if in == nil {
		return nil
	}
	out := new(PodAlert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodAlert) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodAlertList) DeepCopyInto(out *PodAlertList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PodAlert, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodAlertList.
func (in *PodAlertList) DeepCopy() *PodAlertList {
	if in == nil {
		return nil
	}
	out := new(PodAlertList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodAlertList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodAlertSpec) DeepCopyInto(out *PodAlertSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.PodName != nil {
		in, out := &in.PodName, &out.PodName
		*out = new(string)
		**out = **in
	}
	out.CheckInterval = in.CheckInterval
	out.AlertInterval = in.AlertInterval
	if in.Receivers != nil {
		in, out := &in.Receivers, &out.Receivers
		*out = make([]Receiver, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Vars != nil {
		in, out := &in.Vars, &out.Vars
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodAlertSpec.
func (in *PodAlertSpec) DeepCopy() *PodAlertSpec {
	if in == nil {
		return nil
	}
	out := new(PodAlertSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Receiver) DeepCopyInto(out *Receiver) {
	*out = *in
	if in.To != nil {
		in, out := &in.To, &out.To
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Receiver.
func (in *Receiver) DeepCopy() *Receiver {
	if in == nil {
		return nil
	}
	out := new(Receiver)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchlightPlugin) DeepCopyInto(out *SearchlightPlugin) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchlightPlugin.
func (in *SearchlightPlugin) DeepCopy() *SearchlightPlugin {
	if in == nil {
		return nil
	}
	out := new(SearchlightPlugin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchlightPlugin) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchlightPluginList) DeepCopyInto(out *SearchlightPluginList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SearchlightPlugin, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchlightPluginList.
func (in *SearchlightPluginList) DeepCopy() *SearchlightPluginList {
	if in == nil {
		return nil
	}
	out := new(SearchlightPluginList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SearchlightPluginList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SearchlightPluginSpec) DeepCopyInto(out *SearchlightPluginSpec) {
	*out = *in
	if in.Webhook != nil {
		in, out := &in.Webhook, &out.Webhook
		*out = new(WebhookServiceSpec)
		**out = **in
	}
	if in.AlertKinds != nil {
		in, out := &in.AlertKinds, &out.AlertKinds
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Arguments.DeepCopyInto(&out.Arguments)
	if in.States != nil {
		in, out := &in.States, &out.States
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SearchlightPluginSpec.
func (in *SearchlightPluginSpec) DeepCopy() *SearchlightPluginSpec {
	if in == nil {
		return nil
	}
	out := new(SearchlightPluginSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookServiceSpec) DeepCopyInto(out *WebhookServiceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookServiceSpec.
func (in *WebhookServiceSpec) DeepCopy() *WebhookServiceSpec {
	if in == nil {
		return nil
	}
	out := new(WebhookServiceSpec)
	in.DeepCopyInto(out)
	return out
}
