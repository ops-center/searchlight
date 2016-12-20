package env

import (
	"errors"
	"os"
	"strings"
)

type Environment string

const (
	Dev    Environment = "dev"
	QA     Environment = "qa"
	Prod   Environment = "prod"
	Onebox Environment = "onebox"
	BoxDev Environment = "boxdev"
)

const (
	Key           = "APPSCODE_ENV"
	ProdApiServer = "https://api.appscode.com:3443"
	QAApiServer   = "https://api.appscode.info:3443"
)

func (e Environment) IsPublic() bool {
	switch e {
	case Prod, Onebox:
		return true
	default:
		return false
	}
}

func (e Environment) IsHosted() bool {
	switch e {
	case Dev, QA, Prod:
		return true
	default:
		return false
	}
}

func (e Environment) DebugEnabled() bool {
	switch e {
	case Dev, QA, BoxDev:
		return true
	default:
		return false
	}
}

func (e Environment) DevMode() bool {
	return e == Dev || e == BoxDev
}

func (e Environment) String() string {
	return string(e)
}

func (e *Environment) MarshalJSON() ([]byte, error) {
	return []byte(`"` + *e + `"`), nil
}

func (e *Environment) UnmarshalJSON(data []byte) error {
	if e == nil {
		return errors.New("jsontypes.ArrayOrInt: UnmarshalJSON on nil pointer")
	}
	*e = FromString(string(data[1 : len(data)-1]))
	return nil
}

func FromHost() Environment {
	return FromString(strings.ToLower(strings.TrimSpace(os.Getenv(Key))))
}

func FromString(e string) Environment {
	switch e {
	case "prod":
		return Prod
	case "onebox":
		return Onebox
	case "qa":
		return QA
	case "boxdev":
		return BoxDev
	default:
		return Dev
	}
}
