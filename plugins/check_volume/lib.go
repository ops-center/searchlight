package check_volume

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/appscode/searchlight/pkg/config"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	kApi "k8s.io/kubernetes/pkg/api"
)

const (
	awsElasticBlockStorePluginName = "kubernetes.io~aws-ebs"
	azureDataDiskPluginName        = "kubernetes.io~azure-disk"
	azureFilePluginName            = "kubernetes.io~azure-file"
	cephfsPluginName               = "kubernetes.io~cephfs"
	cinderVolumePluginName         = "kubernetes.io~cinder"
	configMapPluginName            = "kubernetes.io~configmap"
	downwardAPIPluginName          = "kubernetes.io~downward-api"
	emptyDirPluginName             = "kubernetes.io~empty-dir"
	fcPluginName                   = "kubernetes.io~fc"
	flockerPluginName              = "kubernetes.io~flocker"
	gcePersistentDiskPluginName    = "kubernetes.io~gce-pd"
	gitRepoPluginName              = "kubernetes.io~git-repo"
	glusterfsPluginName            = "kubernetes.io~glusterfs"
	hostPathPluginName             = "kubernetes.io~host-path"
	iscsiPluginName                = "kubernetes.io~iscsi"
	nfsPluginName                  = "kubernetes.io~nfs"
	quobytePluginName              = "kubernetes.io~quobyte"
	rbdPluginName                  = "kubernetes.io~rbd"
	secretPluginName               = "kubernetes.io~secret"
	vsphereVolumePluginName        = "kubernetes.io~vsphere-volume"
)

func getVolumePluginName(volumeSource *kApi.VolumeSource) string {
	if volumeSource.AWSElasticBlockStore != nil {
		return awsElasticBlockStorePluginName
	} else if volumeSource.AzureDisk != nil {
		return azureDataDiskPluginName
	} else if volumeSource.AzureFile != nil {
		return azureFilePluginName
	} else if volumeSource.CephFS != nil {
		return cephfsPluginName
	} else if volumeSource.Cinder != nil {
		return cinderVolumePluginName
	} else if volumeSource.ConfigMap != nil {
		return configMapPluginName
	} else if volumeSource.DownwardAPI != nil {
		return downwardAPIPluginName
	} else if volumeSource.EmptyDir != nil {
		return emptyDirPluginName
	} else if volumeSource.FC != nil {
		return fcPluginName
	} else if volumeSource.Flocker != nil {
		return flockerPluginName
	} else if volumeSource.GCEPersistentDisk != nil {
		return gcePersistentDiskPluginName
	} else if volumeSource.GitRepo != nil {
		return gitRepoPluginName
	} else if volumeSource.Glusterfs != nil {
		return glusterfsPluginName
	} else if volumeSource.HostPath != nil {
		return hostPathPluginName
	} else if volumeSource.ISCSI != nil {
		return iscsiPluginName
	} else if volumeSource.NFS != nil {
		return nfsPluginName
	} else if volumeSource.Quobyte != nil {
		return quobytePluginName
	} else if volumeSource.RBD != nil {
		return rbdPluginName
	} else if volumeSource.Secret != nil {
		return secretPluginName
	} else if volumeSource.VsphereVolume != nil {
		return vsphereVolumePluginName
	}
	return ""
}

func getPersistentVolumePluginName(volumeSource *kApi.PersistentVolumeSource) string {
	if volumeSource.AWSElasticBlockStore != nil {
		return awsElasticBlockStorePluginName
	} else if volumeSource.AzureDisk != nil {
		return azureDataDiskPluginName
	} else if volumeSource.AzureFile != nil {
		return azureFilePluginName
	} else if volumeSource.CephFS != nil {
		return cephfsPluginName
	} else if volumeSource.Cinder != nil {
		return cinderVolumePluginName
	} else if volumeSource.FC != nil {
		return fcPluginName
	} else if volumeSource.Flocker != nil {
		return flockerPluginName
	} else if volumeSource.GCEPersistentDisk != nil {
		return gcePersistentDiskPluginName
	} else if volumeSource.Glusterfs != nil {
		return glusterfsPluginName
	} else if volumeSource.HostPath != nil {
		return hostPathPluginName
	} else if volumeSource.ISCSI != nil {
		return iscsiPluginName
	} else if volumeSource.NFS != nil {
		return nfsPluginName
	} else if volumeSource.Quobyte != nil {
		return quobytePluginName
	} else if volumeSource.RBD != nil {
		return rbdPluginName
	} else if volumeSource.VsphereVolume != nil {
		return vsphereVolumePluginName
	}
	return ""
}

const (
	hostFactPort = 56977
)

type request struct {
	host      string
	name      string
	warning   float64
	critical  float64
	node_stat bool
}

type usageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

func getUsage(hostIP, path string) (*usageStat, error) {
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/du", hostIP, hostFactPort))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("p", path)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	usages := make([]*usageStat, 1)
	if err = json.Unmarshal(respData, &usages); err != nil {
		return nil, err
	}

	return usages[0], nil
}

func checkResult(field string, warning, critical, result float64) {
	if result >= critical {
		fmt.Fprintln(os.Stdout, util.State[2], fmt.Sprintf("%v used more than %v%%", field, critical))
		os.Exit(2)
	}
	if result >= warning {
		fmt.Fprintln(os.Stdout, util.State[1], fmt.Sprintf("%v used more than %v%%", field, warning))
		os.Exit(1)
	}
}

func checkDiskStat(req *request, nodeIP, path string) {
	usage, err := getUsage(nodeIP, path)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	warning := req.warning
	critical := req.critical
	checkResult("Disk", warning, critical, usage.UsedPercent)
	checkResult("Inodes", warning, critical, usage.InodesUsedPercent)

	fmt.Fprintln(os.Stdout, util.State[0], "(Disk & Inodes)")
	os.Exit(0)
}

func checkNodeDiskStat(req *request) {
	host := req.host
	parts := strings.Split(host, "@")
	if len(parts) != 2 {
		fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
		os.Exit(3)
	}

	kubeClient, err := config.GetKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	node_name := parts[0]
	node, err := kubeClient.Nodes().Get(node_name)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	if node == nil {
		fmt.Fprintln(os.Stdout, util.State[3], "Node not found")
		os.Exit(3)
	}

	hostIP := ""
	for _, address := range node.Status.Addresses {
		if address.Type == kApi.NodeInternalIP {
			hostIP = address.Address
		}
	}

	if hostIP == "" {
		fmt.Fprintln(os.Stdout, util.State[3], "Node InternalIP not found")
		os.Exit(3)
	}
	checkDiskStat(req, hostIP, "/")
}
func checkPodVolumeStat(req *request) {
	host := req.host
	name := req.name
	parts := strings.Split(host, "@")
	if len(parts) != 2 {
		fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
		os.Exit(3)
	}

	kubeClient, err := config.GetKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	pod_name := parts[0]
	namespace := parts[1]
	pod, err := kubeClient.Pods(namespace).Get(pod_name)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	var volumeSourcePluginName = ""
	var volumeSourceName = ""
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == name {
			if volume.PersistentVolumeClaim != nil {
				claim, err := kubeClient.PersistentVolumeClaims(namespace).Get(volume.PersistentVolumeClaim.ClaimName)
				if err != nil {
					fmt.Fprintln(os.Stdout, util.State[3], err)
					os.Exit(3)

				}
				volume, err := kubeClient.PersistentVolumes().Get(claim.Spec.VolumeName)
				if err != nil {
					fmt.Fprintln(os.Stdout, util.State[3], err)
					os.Exit(3)
				}
				volumeSourcePluginName = getPersistentVolumePluginName(&volume.Spec.PersistentVolumeSource)
				volumeSourceName = volume.Name

			} else {
				volumeSourcePluginName = getVolumePluginName(&volume.VolumeSource)
				volumeSourceName = volume.Name
			}
			break
		}
	}

	if volumeSourcePluginName == "" {
		fmt.Fprintln(os.Stdout, util.State[3], errors.New("Invalid volume source"))
		os.Exit(3)
	}

	path := fmt.Sprintf("/var/lib/kubelet/pods/%v/volumes/%v/%v", pod.UID, volumeSourcePluginName, volumeSourceName)
	checkDiskStat(req, pod.Status.HostIP, path)
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "volume",
		Short:   "Check kubernetes volume",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			util.EnsureFlagsSet(cmd, "host")
			if req.node_stat {
				checkNodeDiskStat(&req)
			} else {
				util.EnsureFlagsSet(cmd, "name")
				checkPodVolumeStat(&req)
			}
		},
	}

	c.Flags().BoolVar(&req.node_stat, "node_stat", false, "Checking Node disk size")
	c.Flags().StringVarP(&req.host, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.name, "name", "N", "", "Volume name")
	c.Flags().Float64VarP(&req.warning, "warning", "w", 75.0, "Warning level value (usage percentage)")
	c.Flags().Float64VarP(&req.critical, "critical", "c", 90.0, "Critical level value (usage percentage)")
	return c
}
