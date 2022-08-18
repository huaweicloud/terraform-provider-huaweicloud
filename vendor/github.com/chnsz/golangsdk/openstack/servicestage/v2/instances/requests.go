package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to deploy a specified component.
type CreateOpts struct {
	// Specified the component instance name.
	// The value can contain 2 to 63 characters, including lowercase letters, digits, and hyphens (-).
	// It must start with a lowercase letter and end with a lowercase letter or digit.
	Name string `json:"name" required:"true"`
	// Specified the environment ID.
	EnvId string `json:"environment_id" required:"true"`
	// Specified the number of instance replicas.
	Replica int `json:"replica" required:"true"`
	// Resource specifications, which can be obtained by using the API in Obtaining All Supported Flavors of Application
	// Resources. If you need to customize resource specifications, the format is as follows:
	//   CUSTOM-xxG:xxC-xxC:xxGi-xxGi
	// The meaning of each part is:
	//   xxG: storage capacity allocated to a component instance (reserved field). You can set it to a fixed number.
	//   xxC-xxC: the maximum and minimum number of CPU cores allocated to a component instance.
	//   xxGi-xxGi: the maximum and minimum memory allocated to a component instance.
	// For example, CUSTOM-10G:0.5C-0.25C:1.6Gi-0.8Gi indicates that the maximum number of CPU cores allocated to a
	// component instance is 0.5, the minimum number of CPU cores is 0.25, the maximum memory is 1.6 Gi, and the minimum
	// memory is 0.8 Gi.
	FlavorId string `json:"flavor_id" required:"true"`
	// Deployed resources.
	ReferResources []ReferResource `json:"refer_resources" required:"true"`
	// Application component version that meets version semantics. Example: 1.0.0.
	Version string `json:"version" required:"true"`
	// Artifact. key indicates the component name. In the Docker container scenario, key indicates the container name.
	// If the source parameters of a component specify the software package source, this parameter is optional, and the
	// software package source of the component is inherited by default. Otherwise, this parameter is mandatory.
	Artifacts map[string]Artifact `json:"artifacts,omitempty"`
	// Configuration parameters, such as environment variables, deployment configurations, and O&M monitoring.
	// By default, this parameter is left blank.
	Configuration Configuration `json:"configuration,omitempty"`
	// Description. The value can contain up to 128 characters.
	Description string `json:"description,omitempty"`
	// External network access.
	ExternalAccesses []ExternalAccess `json:"external_accesses,omitempty"`
}

// ReferResource is an object that specifies the deployed basic and optional resources.
type ReferResource struct {
	// Resource ID.
	// Note: If type is set to ecs, the value of this parameter must be Default.
	ID string `json:"id" required:"true"`
	// Resource type.
	// Basic resources: Cloud Container Engine (CCE), Auto Scaling (AS), and Elastic Cloud Server (ECS).
	// Optional resources: Relational Database Service (RDS), Distributed Cache Service (DCS),
	// Elastic Load Balance (ELB), and other services.
	Type string `json:"type" required:"true"`
	// Application alias, which is provided only in DCS scenario. Value: distributed_session, distributed_cache, or
	// distributed_session, distributed_cache. Default value: distributed_session, distributed_cache.
	ReferAlias string `json:"refer_alias,omitempty"`
	// Reference resource parameter.
	// NOTICE:
	// When type is set to cce, this parameter is mandatory. You need to specify the namespace of the cluster where the
	// component is to be deployed. Example: {"namespace": "default"}.
	// When type is set to ecs, this parameter is mandatory. You need to specify the hosts where the component is to be
	// deployed. Example: {"hosts":["04d9f887-9860-4029-91d1-7d3102903a69", "04d9f887-9860-4029-91d1-7d3102903a70"]}}.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// Artifact is an object that specifies the image storage of the software package.
type Artifact struct {
	// Storage mode. Value: swr, devcloud, or obs.
	Storage string `json:"storage" required:"true"`
	// Type. Value: package (VM-based deployment) or image (container-based deployment).
	Type string `json:"type" required:"true"`
	// Software package or image address.
	// For a component deployed on a VM, this parameter is the software package address.
	// For a component deployed based on a container, this parameter is the image address or component name:v${index}.
	// The latter indicates that the component source code or the image automatically built using the software package
	// will be used.
	URL string `json:"url" required:"true"`
	// Authentication mode. Value: iam or none. Default value: iam.
	Auth string `json:"auth,omitempty"`
	// Version number.
	Version string `json:"version,omitempty"`
	// Property information.
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Configuration is an object that specifies the build configuration for a component instance.
type Configuration struct {
	// Environment variable.
	EnvVariables []Variable `json:"env,omitempty"`
	// Data storage configuration.
	Storages []Storage `json:"storage,omitempty"`
	// Upgrade policy.
	Strategy *Strategy `json:"strategy,omitempty"`
	// Lifecycle.
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
	// Policy list of log collection.
	LogCollectionPolicies []LogCollectionPolicy `json:"logs,omitempty"`
	// Scheduling policy.
	Scheduler *Scheduler `json:"scheduler,omitempty"`
	// Health check.
	Probe *Probe `json:"probes,omitempty"`
}

// Configuration is an object that specifies the environment variable for the component instance.
type Variable struct {
	// Environment variable name.
	// The value contains 1 to 64 characters, including letters, digits, underscores (_), hyphens (-), and dots (.),
	// and cannot start with a digit.
	Name string `json:"name" required:"true"`
	// Environment variable value.
	Value string `json:"value" required:"true"`
}

// Storage is an object that specifies the data storage.
type Storage struct {
	// Storage type. Value:
	// HostPath: host path mounting.
	// EmptyDir: temporary directory mounting.
	// ConfigMap: configuration item mounting.
	// Secret: secret volume mounting.
	// PersistentVolumeClaim: cloud storage mounting.
	Type string `json:"type" required:"true"`
	// Storage parameter.
	Parameters *StorageParams `json:"parameters" required:"true"`
	// Directory mounted to the container.
	Mounts []Mount `json:"mounts" required:"true"`
}

// StorageParams is an extend object that specifies the storage path and name.
type StorageParams struct {
	// Host path. This parameter is applicable to the HostPath storage type.
	Path string `json:"path,omitempty"`
	// Name of a configuration item or secret. This parameter is applicable to the ConfigMap and Secret storage type.
	Name string `json:"name,omitempty"`
	// PVC name. This parameter is applicable to the PersistentVolumeClaim storage type.
	ClaimName string `json:"claimName,omitempty"`
}

// Mount is an object that specifies the directory mounted to the container.
type Mount struct {
	// Specifies the mounted disk path.
	Path string `json:"path" required:"true"`
	// Specifies the mounted disk permission is read-only or read-write.
	ReadOnly *bool `json:"readOnly" required:"true"`
	// Specifies the subpath of the mounted disk.
	SubPath string `json:"subPath,omitempty"`
}

// Strategy is an object that specifies the upgrade type, including in-place upgrade and rolling upgrade.
type Strategy struct {
	// Upgrade policy. Value: Recreate or RollingUpdate (default).
	// The former indicates in-place upgrade while the latter indicates rolling upgrade.
	Upgrade string `json:"upgrade,omitempty"`
}

// Lifecycle is an object that specifies the lifecycle of the component deployment.
type Lifecycle struct {
	// Startup command.
	Entrypoint *Entrypoint `json:"entrypoint,omitempty"`
	// Post-start processing.
	PostStart *Process `json:"post-start,omitempty"`
	// Pre-stop processing.
	PreStop *Process `json:"pre-stop,omitempty"`
}

// Entrypoint is an object that specifies the commands when launching up the deployment.
type Entrypoint struct {
	// Command that can be executed.
	Commands []string `json:"command" required:"true"`
	// Running parameters.
	Args []string `json:"args" required:"true"`
}

// Process is an object that specifies the post-processing or stop pre-processing.
type Process struct {
	// Process type. The value is command or http.
	// The command is to execute the command line, and http is to send an http request.
	Type string `json:"type" required:"true"`
	// Start post-processing or stop pre-processing parameters.
	Parameters *ProcessParams `json:"parameters" required:"true"`
}

// ProcessParams is an object that specifies the arguments of the post-processing or stop pre-processing.
type ProcessParams struct {
	// Command parameters, such as ["sleep", "1"]. Applies to command type.
	Commands []string `json:"command,omitempty"`
	// The port number. Applies to http type.
	Port int `json:"port,omitempty"`
	// Request URL. Applies to http type.
	Path string `json:"path,omitempty"`
	// Defaults to the IP address of the POD instance. You can also specify it yourself. Applies to http type.
	Host string `json:"host,omitempty"`
}

// LogCollectionPolicy is an object that specifies the policy of the log collection.
type LogCollectionPolicy struct {
	// Container mounting path.
	LogPath string `json:"logPath" required:"ture"`
	// Aging period.
	AgingPeriod string `json:"rotate" required:"ture"`
	// The extended host path, the valid values are as follows:
	//	None
	//	PodUID
	//	PodName
	//	PodUID/ContainerName
	//	PodName/ContainerName
	// If omited, means container mounting.
	HostExtendPath string `json:"hostExtendPath,omitempty"`
	// Host mounting path.
	HostPath string `json:"hostPath,omitempty"`
}

// Scheduler is an object that specifies the scheduling policy.
type Scheduler struct {
	// Affinity.
	Affinity *Affinity `json:"affinity,omitempty"`
	// Anti-affinity.
	AntiAffinity *Affinity `json:"anti-affinity,omitempty"`
}

// Affinity is an object that specifies the configuration details of the affinity or anti-affinity.
type Affinity struct {
	// AZ list.
	AvailabilityZones []string `json:"az,omitempty"`
	// Node private IP address list.
	Nodes []string `json:"node,omitempty"`
	// List of component instance names.
	Applications []string `json:"application,omitempty"`
}

// Probe is an object that specifies which probe members we have.
type Probe struct {
	// Component liveness probe.
	LivenessProbe *ProbeDetail `json:"livenessProbe,omitempty"`
	// Component service probe.
	ReadinessProbe *ProbeDetail `json:"readinessProbe,omitempty"`
}

// ProbeDetail is an object that specifies the configuration details of the liveness probe and service probe.
type ProbeDetail struct {
	// Value: http, tcp, or command.
	// The check methods are HTTP request check, TCP port check, and command execution check, respectively.
	Type string `json:"type" required:"true"`
	// Parameters.
	// If type is set to http, the object is HttpObject.
	// If type is set to tcp, the object is CommandObject.
	// If type is set to command, the object is TcpObject.
	Parameters map[string]interface{} `json:"parameters" required:"true"`
	// Interval between the startup and detection.
	Delay int `json:"delay,omitempty"`
	// Detection timeout interval.
	Timeout int `json:"timeout,omitempty"`
}

// HttpObject is an object that specifies the check parameters of the HTTP probe.
type HttpObject struct {
	// Value: HTTP or HTTPS.
	Scheme string `json:"scheme" required:"true"`
	// Port number.
	Port int `json:"port" required:"true"`
	// Request path.
	Path string `json:"path" required:"true"`
	// Pod IP address (default). You can specify an IP address.
	Host string `json:"host,omitempty"`
}

// CommandObject is an object that specifies the check parameters of the command probe.
type CommandObject struct {
	// Command list.
	Command []string `json:"command" required:"true"`
}

// TcpObject is an object that specifies the check parameters of the TCP probe.
type TcpObject struct {
	// Port number.
	Port int `json:"port" required:"true"`
}

// ExternalAccess is an object that specifies the configuration of the external IP access.
type ExternalAccess struct {
	// Protocol. Value: http or https.
	Protocol string `json:"protocol,omitempty"`
	// Access address.
	Address string `json:"address,omitempty"`
	// Port number.
	ForwardPort int `json:"forward_port,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a new instance under ServiceStage application using create option.
// Environment is a collection of infrestructures, covering computing, storage and networks, used for application
// deployment and running.
func Create(c *golangsdk.ServiceClient, appId string, componentId string, opts CreateOpts) (*JobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r JobResp
	_, err = c.Post(rootURL(c, appId, componentId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain the details of a specified component instance (deployment) using its ID.
func Get(c *golangsdk.ServiceClient, appId, componentId, instanceId string) (*Instance, error) {
	var r Instance
	_, err := c.Get(resourceURL(c, appId, componentId, instanceId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// Value range: 0–100.
	// Default value: 1000, indicating that a maximum of 1000 records can be queried and all records are displayed on
	// the same page.
	Limit int `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
	// Sorting field. By default, query results are sorted by creation time.
	// The following enumerated values are supported: create_time, name, and update_time.
	OrderBy string `q:"order_by"`
	// Descending or ascending order. Default value: desc.
	Order string `q:"order"`
}

// List is a method to query the list of the component instances (deployment) using given opts.
func List(c *golangsdk.ServiceClient, appId, componentId string, opts ListOpts) ([]Instance, error) {
	url := rootURL(c, appId, componentId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := InstancePage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractInstances(pages)
}

// UpdateOpts is the structure required by the Update method to update the configuration of the component instance.
type UpdateOpts struct {
	// Application component version that meets version semantics. Example: 1.0.0.
	Version string `json:"version" required:"true"`
	// Resource specifications, which can be obtained by using the API in Obtaining All Supported Flavors of Application
	// Resources. If you need to customize resource specifications, the format is as follows:
	//   CUSTOM-xxG:xxC-xxC:xxGi-xxGi
	// The meaning of each part is:
	//   xxG: storage capacity allocated to a component instance (reserved field). You can set it to a fixed number.
	//   xxC-xxC: the maximum and minimum number of CPU cores allocated to a component instance.
	//   xxGi-xxGi: the maximum and minimum memory allocated to a component instance.
	// For example, CUSTOM-10G:0.5C-0.25C:1.6Gi-0.8Gi indicates that the maximum number of CPU cores allocated to a
	// component instance is 0.5, the minimum number of CPU cores is 0.25, the maximum memory is 1.6 Gi, and the minimum
	// memory is 0.8 Gi.
	FlavorId string `json:"flavor_id,omitempty"`
	// Artifact. key indicates the component name. In the Docker container scenario, key indicates the container name.
	// If the source parameters of a component specify the software package source, this parameter is optional, and the
	// software package source of the component is inherited by default. Otherwise, this parameter is mandatory.
	Artifacts map[string]Artifact `json:"artifacts,omitempty"`
	// Configuration parameters, such as environment variables, deployment configurations, and O&M monitoring.
	// By default, this parameter is left blank.
	Configuration Configuration `json:"configuration,omitempty"`
	// Description. The value can contain up to 128 characters.
	Description *string `json:"description,omitempty"`
	// External network access.
	ExternalAccesses []ExternalAccess `json:"external_accesses,omitempty"`
	// Deployed resources.
	ReferResources []ReferResource `json:"refer_resources" required:"true"`
}

// Update is a method to update the current dependency configuration.
func Update(c *golangsdk.ServiceClient, appId, componentId, instanceId string, opts UpdateOpts) (*JobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r JobResp
	_, err = c.Put(resourceURL(c, appId, componentId, instanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to remove an existing instance.
func Delete(c *golangsdk.ServiceClient, appId, componentId, instanceId string) (*JobResp, error) {
	var r JobResp
	_, err := c.Delete(resourceURL(c, appId, componentId, instanceId), &golangsdk.RequestOpts{
		JSONResponse: &r,
		MoreHeaders:  requestOpts.MoreHeaders,
	})
	return &r, err
}

// UpdateOpts is the structure required by the DoAction method to change component status or upgrade it.
type ActionOpts struct {
	// Specified the actions, the valid actions are: start, stop, restart, scale and rollback.
	Action string `json:"action" required:"true"`
	// Specified the action parameters, required if action is scale or rollback.
	Parameters ActionParams `json:"parameters,omitempty"`
}

// ActionParams is an object that specifies the operate parameters.
type ActionParams struct {
	// Specified the number of replica, required if action is scale.
	Replica string `json:"replica,omitempty"`
	// Specified the list of ECS instance to be deployed during VM scaling, required if action is scale.
	Host []string `json:"hosts,omitempty"`
	// Specified the verison number, required if action is rollback.
	Version string `json:"version,omitempty"`
}

// DoAction is a method to change component status or upgrade the instance.
func DoAction(c *golangsdk.ServiceClient, appId string, componentId, instanceId string, opts CreateOpts) (*Job, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst Job
	_, err = c.Post(rootURL(c, appId, componentId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &rst, err
}
