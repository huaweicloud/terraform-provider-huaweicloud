package clusters

import (
	"encoding/json"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

type ListCluster struct {
	// API type, fixed value Cluster
	Kind string `json:"kind"`
	//API version, fixed value v3
	ApiVersion string `json:"apiVersion"`
	//all Clusters
	Clusters []Clusters `json:"items"`
}

type Clusters struct {
	// API type, fixed value Cluster
	Kind string `json:"kind" required:"true"`
	//API version, fixed value v3
	ApiVersion string `json:"apiversion" required:"true"`
	//Metadata of a Cluster
	Metadata MetaData `json:"metadata" required:"true"`
	//specifications of a Cluster
	Spec Spec `json:"spec" required:"true"`
	//status of a Cluster
	Status Status `json:"status"`
}

// Metadata required to create a cluster
type MetaData struct {
	//Cluster unique name
	Name string `json:"name"`
	//Cluster unique Id
	Id string `json:"uid"`
	// Cluster tag, key/value pair format
	Labels map[string]string `json:"labels,omitempty"`
	//Cluster annotation, key/value pair format
	Annotations map[string]string `json:"annotations,omitempty"`
	// Cluster alias
	Alias string `json:"alias"`
	// Cluster timezone
	Timezone string `json:"timezone"`
}

// Specifications to create a cluster
type Spec struct {
	//Cluster Type: VirtualMachine, BareMetal, or Windows
	Type string `json:"type" required:"true"`
	// Cluster specifications
	Flavor string `json:"flavor" required:"true"`
	// For the cluster version, please fill in v1.7.3-r10 or v1.9.2-r1. Currently only Kubernetes 1.7 and 1.9 clusters are supported.
	Version string `json:"version,omitempty"`
	//Cluster description
	Description string `json:"description,omitempty"`
	//Public IP ID
	PublicIP string `json:"publicip_id,omitempty"`
	// Node network parameters
	HostNetwork HostNetworkSpec `json:"hostNetwork" required:"true"`
	//Container network parameters
	ContainerNetwork ContainerNetworkSpec `json:"containerNetwork" required:"true"`
	//ENI network parameters
	EniNetwork *EniNetworkSpec `json:"eniNetwork,omitempty"`
	// Enable Distributed Cluster Management
	EnableDistMgt bool `json:"enableDistMgt,omitempty"`
	//Authentication parameters
	Authentication AuthenticationSpec `json:"authentication,omitempty"`
	// Charging mode of the cluster, which is 0 (on demand)
	BillingMode int `json:"billingMode,omitempty"`
	//Extended parameter for a cluster
	ExtendParam map[string]interface{} `json:"extendParam,omitempty"`
	//Advanced configuration of master node
	Masters []MasterSpec `json:"masters,omitempty"`
	//Range of kubernetes clusterIp
	KubernetesSvcIPRange string `json:"kubernetesSvcIpRange,omitempty"`
	// Service network, use this to replace KubernetesSvcIPRange
	ServiceNetwork *ServiceNetwork `json:"serviceNetwork,omitempty"`
	//Custom san list for certificates
	CustomSan []string `json:"customSan,omitempty"`
	// Tags of cluster, key value pair format
	ClusterTags []tags.ResourceTag `json:"clusterTags,omitempty"`
	// configurationsOverride
	ConfigurationsOverride []PackageConfiguration `json:"configurationsOverride,omitempty"`
	// Whether to enable IPv6
	IPv6Enable bool `json:"ipv6enable,omitempty"`
	// K8s proxy mode
	KubeProxyMode string `json:"kubeProxyMode,omitempty"`
	// Whether to enable Istio
	SupportIstio bool `json:"supportIstio,omitempty"`
	// The category, the value can be CCE and CCE
	Category string `json:"category,omitempty"`
	// The Encrytion Config
	EncryptionConfig *EncryptionConfig `json:"encryptionConfig,omitempty"`
}

type ServiceNetwork struct {
	IPv4Cidr string `json:"IPv4CIDR,omitempty"`
}

type PackageConfiguration struct {
	Name           string        `json:"name,omitempty"`
	Configurations []interface{} `json:"configurations,omitempty"`
}

type EncryptionConfig struct {
	Mode     string `json:"mode,omitempty"`
	KmsKeyID string `json:"kmsKeyID,omitempty"`
}

// Node network parameters
type HostNetworkSpec struct {
	//The ID of the VPC used to create the node
	VpcId string `json:"vpc" required:"true"`
	//The ID of the subnet used to create the node
	SubnetId string `json:"subnet" required:"true"`
	// The ID of the high speed network used to create bare metal nodes.
	// This parameter is required when creating a bare metal cluster.
	HighwaySubnet string `json:"highwaySubnet,omitempty"`
	//The ID of the Security Group used to create the node
	SecurityGroup string `json:"SecurityGroup,omitempty"`
}

// Container network parameters
type ContainerNetworkSpec struct {
	//Container network type: overlay_l2 , underlay_ipvlan or vpc-router
	Mode string `json:"mode" required:"true"`
	//Container network segment: 172.16.0.0/16 ~ 172.31.0.0/16. If there is a network segment conflict, it will be automatically reselected.
	Cidr string `json:"cidr,omitempty"`
	// List of container CIDR blocks. In clusters of v1.21 and later, the cidrs field is used.
	// When the cluster network type is vpc-router, you can add multiple container CIDR blocks.
	// In versions earlier than v1.21, if the cidrs field is used, the first CIDR element in the array is used as the container CIDR block.
	Cidrs []CidrSpec `json:"cidrs,omitempty"`
}

type CidrSpec struct {
	// Container network segment. Recommended: 10.0.0.0/12-19, 172.16.0.0/16-19, and 192.168.0.0/16-19
	Cidr string `json:"cidr" required:"true"`
}

type EniNetworkSpec struct {
	//Eni network subnet id, will be deprecated in the future
	SubnetId string `json:"eniSubnetId,omitempty"`
	//Eni network cidr, will be deprecated in the future
	Cidr string `json:"eniSubnetCIDR,omitempty"`
	// Eni network subnet IDs
	Subnets []EniSubnetSpec `json:"subnets" required:"true"`
}

type EniSubnetSpec struct {
	SubnetID string `json:"subnetID" required:"true"`
}

// Authentication parameters
type AuthenticationSpec struct {
	//Authentication mode: rbac , x509 or authenticating_proxy
	Mode                string            `json:"mode" required:"true"`
	AuthenticatingProxy map[string]string `json:"authenticatingProxy" required:"true"`
}

type MasterSpec struct {
	// AZ of master node
	MasterAZ string `json:"availabilityZone,omitempty"`
}

type Status struct {
	//The state of the cluster
	Phase string `json:"phase"`
	//The ID of the Job that is operating asynchronously in the cluster
	JobID string `json:"jobID"`
	//Reasons for the cluster to become current
	Reason string `json:"reason"`
	//The status of each component in the cluster
	Conditions Conditions `json:"conditions"`
	//Kube-apiserver access address in the cluster
	Endpoints []Endpoints `json:"-"`
}

type Conditions struct {
	//The type of component
	Type string `json:"type"`
	//The state of the component
	Status string `json:"status"`
	//The reason that the component becomes current
	Reason string `json:"reason"`
}

type Endpoints struct {
	//The address accessed within the user's subnet - Huawei
	Url string `json:"url"`
	//Public network access address - Huawei
	Type string `json:"type"`
	//Internal network address - OTC
	Internal string `json:"internal"`
	//External network address - OTC
	External string `json:"external"`
	//Endpoint of the cluster to be accessed through API Gateway - OTC
	ExternalOTC string `json:"external_otc"`
}

type Certificate struct {
	//API type, fixed value Config
	Kind string `json:"kind"`
	//API version, fixed value v1
	ApiVersion string `json:"apiVersion"`
	//Cluster list
	Clusters []CertClusters `json:"clusters"`
	//User list
	Users []CertUsers `json:"users"`
	//Context list
	Contexts []CertContexts `json:"contexts"`
	//The current context
	CurrentContext string `json:"current-context"`
}

type CertClusters struct {
	//Cluster name
	Name string `json:"name"`
	//Cluster information
	Cluster CertCluster `json:"cluster"`
}

type CertCluster struct {
	//Server IP address
	Server string `json:"server"`
	//Certificate data
	CertAuthorityData string `json:"certificate-authority-data"`
	//whether skip tls verify
	InsecureSkipTLSVerify bool `json:"insecure-skip-tls-verify"`
}

type CertUsers struct {
	//User name
	Name string `json:"name"`
	//Cluster information
	User CertUser `json:"user"`
}

type CertUser struct {
	//Client certificate
	ClientCertData string `json:"client-certificate-data"`
	//Client key data
	ClientKeyData string `json:"client-key-data"`
}

type CertContexts struct {
	//Context name
	Name string `json:"name"`
	//Context information
	Context CertContext `json:"context"`
}

type CertContext struct {
	//Cluster name
	Cluster string `json:"cluster"`
	//User name
	User string `json:"user"`
}

// UnmarshalJSON helps to unmarshal Status fields into needed values.
// OTC and Huawei have different data types and child fields for `endpoints` field in Cluster Status.
// This function handles the unmarshal for both
func (r *Status) UnmarshalJSON(b []byte) error {
	type tmp Status
	var s struct {
		tmp
		Endpoints []Endpoints `json:"endpoints"`
	}

	err := json.Unmarshal(b, &s)

	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError: //check if type error occurred (handles the different endpoint structure for huawei and otc)
			var s struct {
				tmp
				Endpoints Endpoints `json:"endpoints"`
			}
			err := json.Unmarshal(b, &s)
			if err != nil {
				return err
			}
			*r = Status(s.tmp)
			r.Endpoints = []Endpoints{{Internal: s.Endpoints.Internal,
				External:    s.Endpoints.External,
				ExternalOTC: s.Endpoints.ExternalOTC}}
			return nil
		default:
			return err
		}
	}

	*r = Status(s.tmp)
	r.Endpoints = s.Endpoints

	return err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a cluster.
func (r commonResult) Extract() (*Clusters, error) {
	var s Clusters
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractCluster is a function that accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func (r commonResult) ExtractClusters() ([]Clusters, error) {
	var s ListCluster
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s.Clusters, nil

}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Cluster.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Cluster.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Cluster.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// ListResult represents the result of a list operation. Call its ExtractCluster
// method to interpret it as a Cluster.
type ListResult struct {
	commonResult
}

type GetCertResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a cluster.
func (r GetCertResult) Extract() (*Certificate, error) {
	var s Certificate
	err := r.ExtractInto(&s)
	return &s, err
}

// UpdateIpResult represents the result of an update operation. Call its Extract
// method to interpret it as a Cluster.
type UpdateIpResult struct {
	golangsdk.ErrResult
}

type OperationResult struct {
	golangsdk.ErrResult
}
