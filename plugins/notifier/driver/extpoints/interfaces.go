package extpoints

import api "github.com/appscode/api/kubernetes/v1beta1"

type Driver interface {
	Notify(*api.IncidentNotifyRequest) error
}
