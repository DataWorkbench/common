package helm

// ***************************************************************
// K8s Type Zone
// ***************************************************************

// firstAtAll, create docker registry secret by kubectl:
// ***************************************************************
// ImageConfig
// ***************************************************************
// kubectl create secret docker-registry my-docker-registry-secret
//                                       --docker-server=<your-registry-server>
//                                       --docker-username=<your-name>
//                                       --docker-password=<your-pword>
//                                       --docker-email=<your-email>
type Image struct {
	// TODO: check pull secrets
	Registry    string   `json:"registry,omitempty" yaml:"registry,omitempty"`
	PullSecrets []string `json:"pullSecrets,omitempty" yaml:"pullSecrets,omitempty"`
	PullPolicy  string   `json:"pullPolicy,omitempty" yaml:"pullPolicy,omitempty"`

	Tag string `json:"tag,omitempty" yaml:"-"`
}

// update from other if the field is ""
func (i *Image) Update(source *Image) {
	if source == nil {
		return
	}
	if i.Registry == "" && source.Registry != "" {
		i.Registry = source.Registry
	}
	if len(i.PullSecrets) == 0 && len(source.PullSecrets) > 0 {
		i.PullSecrets = source.PullSecrets
	}
	if i.PullPolicy == "" && source.PullPolicy != "" {
		i.PullPolicy = source.PullPolicy
	}
}

type Resource struct {
	Cpu    string `json:"cpu,omitempty"    yaml:"cpu,omitempty"`
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type Resources struct {
	Limits   Resource `json:"limits,omitempty"   yaml:"limits,omitempty"`
	Requests Resource `json:"requests,omitempty" yaml:"requests,omitempty"`
}

// k8s workload(a deployment / statefulset / .. in a Chart) configurations
type Workload struct {
	Replicas       int8   `json:"replicas,omitempty"       yaml:"replicas,omitempty"`
	UpdateStrategy string `json:"updateStrategy,omitempty" yaml:"updateStrategy,omitempty"`
	TimeoutSecond  int    `json:"-"                        yaml:"timeoutSecond,omitempty"`

	Resources  *Resources  `json:"resources,omitempty"  yaml:"resources,omitempty"`
	Persistent *Persistent `json:"persistent,omitempty" yaml:"persistent,omitempty"`
}

type LocalPv struct {
	Nodes []string `json:"nodes" yaml:"nodes"`
	Home  string   `json:"home"  yaml:"-"`
}

type Persistent struct {
	// for local pv, eg: 10Gi
	Size     string   `json:"size,omitempty"    yaml:"size"`
	HostPath string   `json:"hostPath"          yaml:"-"`
	LocalPv  *LocalPv `json:"localPv,omitempty" yaml:"localPv"`
}

func (p *Persistent) UpdateLocalPv(localPvHome string, nodes []string) {
	// TODO: check if localPv exist and start with localPvHome
	p.LocalPv.Home = localPvHome

	if len(p.LocalPv.Nodes) == 0 {
		p.LocalPv.Nodes = nodes
	}
}
