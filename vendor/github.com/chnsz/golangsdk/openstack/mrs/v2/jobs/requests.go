package jobs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// CreateOpts is a structure representing information of the job creation.
type CreateOpts struct {
	// Type of a job, and the valid values are as follows:
	//   MapReduce
	//   SparkSubmit
	//   HiveScript
	//   HiveSql
	//   DistCp, importing and exporting data
	//   SparkScript
	//   SparkSql
	//   Flink
	// NOTE:
	//   Spark, Hive, and Flink jobs can be added to only clusters that include Spark, Hive, and Flink components.
	JobType string `json:"job_type" required:"true"`
	// Job name. It contains 1 to 64 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// NOTE:
	// Identical job names are allowed but not recommended.
	JobName string `json:"job_name" required:"true"`
	// Key parameter for program execution.
	// The parameter is specified by the function of the user's program.
	// MRS is only responsible for loading the parameter.
	// The parameter contains a maximum of 4,096 characters, excluding special characters such as ;|&>'<$,
	// and can be left blank.
	// NOTE:
	//   If you enter a parameter with sensitive information (such as the login password), the parameter may be exposed
	//   in the job details display and log printing. Exercise caution when performing this operation.
	//   For MRS 1.9.2 or later, a file path on OBS can start with obs://. To use this format to submit HiveScript or
	//   HiveSQL jobs, choose Components > Hive > Service Configuration on the cluster details page, set Type to All,
	//   and search for core.site.customized.configs. Add the endpoint configuration item (fs.obs.endpoint) of OBS and
	//   enter the endpoint corresponding to OBS in Value. Obtain the value from Regions and Endpoints.
	//   For MRS 3.0.2 or later, a file path on OBS can start with obs://. To use this format to submit HiveScript or
	//   HiveSQL jobs, choose Components > Hive > Service Configuration on Manager. Switch Basic to All, and search for
	//   core.site.customized.configs. Add the endpoint configuration item (fs.obs.endpoint) of OBS and enter the
	//   endpoint corresponding to OBS in Value. Obtain the value from Regions and Endpoints.
	Arguments []string `json:"arguments,omitempty"`
	// Program system parameter.
	// The parameter contains a maximum of 2,048 characters, excluding special characters such as ><|'`&!\, and can be
	// left blank.
	Properties map[string]string `json:"properties,omitempty"`
}

// CreateOptsBuilder is an interface which to support request body build of the job creation.
type CreateOptsBuilder interface {
	ToJobCreateMap() (map[string]interface{}, error)
}

// ToJobCreateMap is a method which to build a request body by the CreateOpts.
func (opts CreateOpts) ToJobCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method to create a new mapreduce job.
func Create(client *golangsdk.ServiceClient, clusterId string, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToJobCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, clusterId), reqBody, &r.Body, nil)
	return
}

// Get is a method to get an existing mapreduce job by cluster ID and job ID.
func Get(client *golangsdk.ServiceClient, clsuterId, jobId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, clsuterId, jobId), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
		OkCodes:     []int{200, 202},
	})
	return
}

// ListOpts is a structure representing information of the job updation.
type ListOpts struct {
	// Job name. It contains 1 to 64 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
	JobName string `q:"job_name"`
	// Type of a job, and the valid values are as follows:
	//   MapReduce
	//   SparkSubmit
	//   HiveScript
	//   HiveSql
	//   DistCp, importing and exporting data
	//   SparkScript
	//   SparkSql
	//   Flink
	JobType string `q:"job_type"`
	// Execution status of a job.
	//   FAILED: indicates that the job fails to be executed.
	//   KILLED: indicates that the job is terminated.
	//   New: indicates that the job is created.
	//   NEW_SAVING: indicates that the job has been created and is being saved.
	//   SUBMITTED: indicates that the job is submitted.
	//   ACCEPTED: indicates that the job is accepted.
	//   RUNNING: indicates that the job is running.
	//   FINISHED: indicates that the job is completed.
	JobState string `q:"job_state"`
	// Execution result of a job.
	//   FAILED: indicates that the job fails to be executed.
	//   KILLED: indicates that the job is manually terminated during execution.
	//   UNDEFINED: indicates that the job is being executed.
	//   SUCCEEDED: indicates that the job has been successfully executed.
	JobResult string `q:"job_result"`
	// Number of records displayed on each page in the returned result. The default value is 10.
	Limit int `q:"limit"`
	// Offset.
	// The default offset from which the job list starts to be queried is 1.
	Offset int `q:"offset"`
	// Ranking mode of returned results. The default value is desc.
	//   asc: indicates that the returned results are ranked in ascending order.
	//   desc: indicates that the returned results are ranked in descending order.
	SortBy string `q:"sort_by"`
	// UTC timestamp after which a job is submitted, in milliseconds. For example, 1562032041362.
	SubmittedTimeBegin int `q:"submitted_time_begin"`
	// UTC timestamp before which a job is submitted, in milliseconds. For example, 1562032041362.
	SubmittedTimeEnd int `q:"submitted_time_end"`
}

// ListOptsBuilder is an interface which to support request query build of the job list operation.
type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

// ToListQuery is a method which to build a request query by the ListOpts.
func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more mapreduce jobs according to the query parameters.
func List(client *golangsdk.ServiceClient, clusterId string, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, clusterId)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return JobPage{pagination.SinglePageBase(r)}
	})
}

// DeleteOpts is a structure representing information of the job delete operation.
type DeleteOpts struct {
	JobIds []string `json:"job_id_list,omitempty"`
}

// DeleteOptsBuilder is an interface which to support request body build of the job delete operation.
type DeleteOptsBuilder interface {
	ToJobDeleteMap() (map[string]interface{}, error)
}

// ToJobDeleteMap is a method which to build a request body by the DeleteOpts.
func (opts DeleteOpts) ToJobDeleteMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Delete is a method to delete an existing mapreduce job.
func Delete(client *golangsdk.ServiceClient, clusterId string, opts DeleteOptsBuilder) (r DeleteResult) {
	reqBody, err := opts.ToJobDeleteMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(deleteURL(client, clusterId), reqBody, nil, nil)
	return
}
