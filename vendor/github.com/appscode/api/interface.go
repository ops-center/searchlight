package api

import (
	"github.com/appscode/api/dtypes"
	"github.com/xeipuuv/gojsonschema"
)

type Request interface {
	IsRequest()
	IsValid() (*gojsonschema.Result, error)

	Reset()
	String() string
	ProtoMessage()
}

type Response interface {
	Reset()
	String() string
	ProtoMessage()
	GetStatus() *dtypes.Status
	SetStatus(*dtypes.Status)
}
