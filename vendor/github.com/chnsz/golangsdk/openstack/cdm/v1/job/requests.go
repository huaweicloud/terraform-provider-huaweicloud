package job

import "github.com/chnsz/golangsdk"

const (
	GenericJdbcConnector   = "generic-jdbc-connector"
	ObsConnector           = "obs-connector"
	ThirdpartyObsConnector = "thirdparty-obs-connector"
	HdfsConnector          = "hdfs-connector"
	HbaseConnector         = "hbase-connector"
	HiveConnector          = "hive-connector"
	SftpConnector          = "sftp-connector"
	FtpConnector           = "ftp-connector"
	MongodbConnector       = "mongodb-connector"
	KafkaConnector         = "kafka-connector"
	DisConnector           = "dis-connector"
	ElasticsearchConnector = "elasticsearch-connector"
	DliConnector           = "dli-connector"
	OpentsdbConnector      = "opentsdb-connector"
	DmsKafkaConnector      = "dms-kafka-connector"

	StatusBooting         = "BOOTING"
	StatusFailureOnSubmit = "FAILURE_ON_SUBMIT"
	StatusRunning         = "RUNNING"
	StatusSucceeded       = "SUCCEEDED"
	StatusFailed          = "FAILED"
	StatusUnknown         = "UNKNOWN"
	StatusNeverExecuted   = "NEVER_EXECUTED"
)

type JobCreateOpts struct {
	Jobs []Job `json:"jobs" required:"true"`
}

type Job struct {
	// Job type. The options are as follows:
	// NORMAL_JOB: table/file migration
	// BATCH_JOB: entire DB migration
	// SCENARIO_JOB: scenario migration
	JobType string `json:"job_type,omitempty"`
	// Job name, which contains 1 to 240 characters.
	Name string `json:"name,omitempty"`

	// Source link name
	FromLinkName string `json:"from-link-name,omitempty"`
	// Source link type
	FromConnectorName string `json:"from-connector-name,omitempty"`
	// Source link parameter configuration
	FromConfigValues JobConfigs `json:"from-config-values,omitempty"`

	// Destination link name
	ToLinkName string `json:"to-link-name,omitempty"`
	// Destination link type
	ToConnectorName string `json:"to-connector-name,omitempty"`
	// Destination link parameter configuration
	ToConfigValues JobConfigs `json:"to-config-values,omitempty"`

	// Job parameter configuration
	DriverConfigValues JobConfigs `json:"driver-config-values,omitempty"`

	CreationUser string `json:"creation-user,omitempty"`
	CreationDate *int   `json:"creation-date,omitempty"`
	UpdateDate   *int   `json:"update-date,omitempty"`
	UpdateUser   string `json:"update-user,omitempty"`
	// Status of a job. The options are as follows:
	// BOOTING: starting
	// RUNNING: running
	// SUCCEEDED: successful
	// FAILED: failed
	// NEW: not executed
	Status string `json:"status,omitempty"`
}

type JobConfigs struct {
	Configs []Configs `json:"configs,omitempty"`
}

type Configs struct {
	Inputs []Input `json:"inputs" required:"true"`
	Name   string  `json:"name" required:"true"`
}

type Input struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

type GetJobsOpts struct {
	// When job_name is all, this parameter is used for fuzzy job filtering.
	Filter string `q:"filter"`
	// Page number  Minimum: 1
	PageNo int `q:"page_no"`
	// Number of jobs on each page. The value ranges from 10 to 100.
	PageSize int `q:"page_size"`
	// Job type. The options are as follows:
	// jobType=NORMAL_JOB: table/file migration job
	// jobType=BATCH_JOB: entire DB migration job
	// jobType=SCENARIO_JOB: scenario migration job
	// If this parameter is not specified, only table/file migration jobs are queried by default.
	JobType string `q:"jobType"`
}

type ListJobSubmissionsOpts struct {
	JobName string `q:"jname"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, clusterId string, opts JobCreateOpts) (*CreateResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CreateResponse
	_, err = c.Post(createURL(c, clusterId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Update(c *golangsdk.ServiceClient, clusterId string, name string, opts JobCreateOpts) (*UpdateResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst UpdateResponse
	_, err = c.Put(updateURL(c, clusterId, name), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Delete(c *golangsdk.ServiceClient, clusterId string, name string) (*ErrorResponse, error) {
	var rst ErrorResponse
	_, err := c.DeleteWithResponse(deleteURL(c, clusterId, name), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Get(c *golangsdk.ServiceClient, clusterId string, jobName string, opts GetJobsOpts) (*JobsDetail, error) {
	url := getURL(c, clusterId, jobName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst JobsDetail
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Start(c *golangsdk.ServiceClient, clusterId string, jobName string) (*StartJobResponse, error) {
	var rst StartJobResponse
	_, err := c.Put(startURL(c, clusterId, jobName), nil, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Stop(c *golangsdk.ServiceClient, clusterId string, jobName string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Put(stopURL(c, clusterId, jobName), nil, nil, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r
}

func GetJobStatus(c *golangsdk.ServiceClient, clusterId string, jobName string) (*StatusResponse, error) {
	var rst StatusResponse
	_, err := c.Get(getStatusURL(c, clusterId, jobName), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func ListJobSubmissions(c *golangsdk.ServiceClient, clusterId string, opts ListJobSubmissionsOpts) (*ListSubmissionsRst,
	error) {
	url := ListJobSubmissionsURL(c, clusterId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst ListSubmissionsRst
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}
