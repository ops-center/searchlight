package plugin

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/appscode/searchlight/pkg/controller/host"
)

const (
	CRITICAL int = 2
	OK       int = 0
)

func GetKubeObjectInfo(hostname string) (objectType string, objectName string, namespace string) {
	parts := strings.Split(hostname, "@")
	if len(parts) != 2 {
		Fatalln(errors.New("Invalid icinga host.name"))

	}
	name := parts[0]
	namespace = parts[1]

	objectType = ""
	objectName = ""
	if name != host.CheckCommandPodExists && name != host.CheckCommandPodStatus {
		parts = strings.Split(name, "|")
		if len(parts) == 1 {
			objectType = host.TypePods
			objectName = parts[0]
		} else if len(parts) == 2 {
			objectType = parts[0]
			objectName = parts[1]
		} else {
			Fatalln(errors.New("Invalid icinga host.name"))
		}
	}
	return
}

func Fatalln(i interface{}) {
	if i != nil {
		fmt.Println(i)
		os.Exit(1)
	}
}
