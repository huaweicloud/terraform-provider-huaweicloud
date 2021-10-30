package batches

import "github.com/chnsz/golangsdk"

// CreateOpts is a struct which will be used to submit a spark job.
type CreateOpts struct {
	// Name of the package that is of the JAR or pyFile type and has been uploaded to the DLI resource management
	// system. You can also specify an OBS path, for example, obs://Bucket name/Package name.
	File string `json:"file" required:"true"`
	// Queue name. Set this parameter to the name of the created DLI queue.
	// NOTE: This parameter is compatible with the cluster_name parameter. That is, if cluster_name is used to specify a
	//       queue, the queue is still valid.
	// You are advised to use the queue parameter. The queue and cluster_name parameters cannot coexist.
	Queue string `json:"queue" required:"true"`
	// Java/Spark main class of the batch processing job.
	ClassName *string `json:"class_name,omitempty"`
	// Queue name. Set this parameter to the created DLI queue name.
	// NOTE: You are advised to use the queue parameter. The queue and cluster_name parameters cannot coexist.
	ClusterName string `json:"cluster_name,omitempty"`
	// Input parameters of the main class, that is, application parameters.
	Arguments []string `json:"args,omitempty"`
	// Compute resource type. Currently, resource types A, B, and C are available.
	// If this parameter is not specified, the minimum configuration (type A) is used.
	Specification string `json:"sc_type,omitempty"`
	// Name of the package that is of the JAR type and has been uploaded to the DLI resource management system.
	// You can also specify an OBS path, for example, obs://Bucket name/Package name.
	Jars []string `json:"jars,omitempty"`
	// Name of the package that is of the PyFile type and has been uploaded to the DLI resource management system.
	// You can also specify an OBS path, for example, obs://Bucket name/Package name.
	PythonFiles []string `json:"python_files,omitempty"`
	// Name of the package that is of the file type and has been uploaded to the DLI resource management system.
	// You can also specify an OBS path, for example, obs://Bucket name/Package name.
	Files []string `json:"files,omitempty"`
	// Name of the dependent system resource module. You can view the module name using the API related to Querying
	// Resource Packages in a Group. DLI provides dependencies for executing datasource jobs.
	// The following table lists the dependency modules corresponding to different services.
	//   CloudTable/MRS HBase: sys.datasource.hbase
	//   CloudTable/MRS OpenTSDB: sys.datasource.opentsdb
	//   RDS MySQL: sys.datasource.rds
	//   RDS Postgre: preset
	//   DWS: preset
	//   CSS: sys.datasource.css
	Modules []string `json:"modules,omitempty"`
	// JSON object list, including the name and type of the JSON package that has been uploaded to the queue.
	Resources []Resource `json:"resources,omitempty"`
	// JSON object list, including the package group resource. For details about the format, see the request example.
	// If the type of the name in resources is not verified, the package with the name exists in the group.
	Groups []Group `json:"groups,omitempty"`
	// Batch configuration item. For details, see Spark Configuration.
	Configurations map[string]interface{} `json:"conf,omitempty"`
	// Batch processing task name. The value contains a maximum of 128 characters.
	Name string `json:"name,omitempty"`
	// Driver memory of the Spark application, for example, 2 GB and 2048 MB. This configuration item replaces the
	// default parameter in sc_type. The unit must be provided. Otherwise, the startup fails.
	DriverMemory string `json:"driver_memory,omitempty"`
	// Number of CPU cores of the Spark application driver.
	// This configuration item replaces the default parameter in sc_type.
	DriverCores int `json:"driver_cores,omitempty"`
	// Executor memory of the Spark application, for example, 2 GB and 2048 MB. This configuration item replaces the
	// default parameter in sc_type. The unit must be provided. Otherwise, the startup fails.
	ExecutorMemory string `json:"executor_memory,omitempty"`
	// Number of CPU cores of each Executor in the Spark application.
	// This configuration item replaces the default parameter in sc_type.
	ExecutorCores int `json:"executor_cores,omitempty"`
	// Number of Executors in a Spark application. This configuration item replaces the default parameter in sc_type.
	NumExecutors int `json:"num_executors,omitempty"`
	// OBS bucket for storing the Spark jobs. Set this parameter when you need to save jobs.
	ObsBucket string `json:"obs_bucket,omitempty"`
	// Whether to enable the retry function.
	// If enabled, Spark jobs will be automatically retried after an exception occurs. The default value is false.
	AutoRecovery bool `json:"auto_recovery,omitempty"`
	// Maximum retry times. The maximum value is 100, and the default value is 20.
	MaxRetryTimes int `json:"max_retry_times,omitempty"`
	// Job feature. Type of the Spark image used by a job.
	// basic: indicates that the basic Spark image provided by DLI is used.
	// custom: indicates that the user-defined Spark image is used.
	// ai: indicates that the AI image provided by DLI is used.
	Feature string `json:"feature,omitempty"`
	// Version of the Spark component used by a job. Set this parameter when feature is set to basic or ai.
	// If this parameter is not set, the default Spark component version 2.3.2 is used.
	SparkVersion string `json:"spark_version,omitempty"`
	// Custom image. The format is Organization name/Image name:Image version.
	// This parameter is valid only when feature is set to custom.
	// You can use this parameter with the feature parameter to specify a user-defined Spark image for job running.
	Image string `json:"image,omitempty"`
	// To access metadata, set this parameter to DLI.
	CatalogName string `json:"catalog_name,omitempty"`
}

// Group is an object which will be build up a package group.
type Group struct {
	// User group name.
	Name string `json:"name,omitempty"`
	// User group resource.
	Resources []Resource `json:"resources,omitempty"`
}

// Resource is an object which specified the user group resource.
type Resource struct {
	// Resource name. You can also specify an OBS path, for example, obs://Bucket name/Package name.
	Name string `json:"name,omitempty"`
	// Resource type.
	Type string `json:"type,omitempty"`
}

// Create is a method to submit a Spark job with given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, nil)
	if err == nil {
		var r CreateResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Get is a method to obtain the specified Spark job with job ID.
func Get(c *golangsdk.ServiceClient, jobId string) (*CreateResp, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, jobId), &rst.Body, nil)
	if err == nil {
		var r CreateResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// GetState is a method to obtain the state of specified Spark job with job ID.
func GetState(c *golangsdk.ServiceClient, jobId string) (*StateResp, error) {
	var rst golangsdk.Result
	_, err := c.Get(stateURL(c, jobId), &rst.Body, nil)
	if err == nil {
		var r StateResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Delete is a method to cancel the unfinished spark job.
func Delete(c *golangsdk.ServiceClient, jobId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, jobId), nil)
	return &r
}
