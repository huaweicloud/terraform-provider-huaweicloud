package flinkjob

import (
	"github.com/chnsz/golangsdk"
)

// StreamGraphOpts is the structure that represents used to generate a static stream graph or simplified stream graph.
type StreamGraphOpts struct {
	// The job ID of the flink job.
	JobId string `json:"-" required:"true"`
	// Stream SQL statement.
	SqlBody string `json:"sql_body" required:"true"`
	// The total number of CUs.
	CuNumber *int `json:"cu_number,omitempty"`
	// The number of CUs of the management unit.
	ManagerCuNumber *int `json:"manager_cu_number,omitempty"`
	// The number of parallel jobs
	ParallelNumber *int `json:"parallel_number,omitempty"`
	// The number of CUs in a taskManager.
	TmCus *int `json:"tm_cus,omitempty"`
	// The number of slots in a taskManager.
	TmSlotNum *int `json:"tm_slot_num,omitempty"`
	// The operator configurations.
	OperatorConfig string `json:"operator_config,omitempty"`
	// Whether to estimate static resources.
	StaticEstimator *bool `json:"static_estimator,omitempty"`
	// Job type. Only flink_opensource_sql_job job is supported.
	JobType string `json:"job_type,omitempty"`
	// Stream graph type. The valid values aer as follows:
	// + simple_graph: Simplified stream graph.
	// + job_graph: Static stream graph.
	GraphType string `json:"graph_type,omitempty"`
	// Traffic or hit ratio of each operator, which is a string in JSON format.
	StaticEstimatorConfig string `json:"static_estimator_config,omitempty"`
	// Flink version. Currently, only 1.10 and 1.12 are supported.
	FlinkVersion string `json:"flink_version,omitempty"`
}

// CreateFlinkSqlJobGraph is a method to generate a stream graph for a Flink SQL job using given parameters.
func CreateFlinkSqlJobGraph(c *golangsdk.ServiceClient, opts StreamGraphOpts) (*streamGraphResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r streamGraphResp
	_, err = c.Post(streamGraphURL(c, opts.JobId), b, &r, &golangsdk.RequestOpts{})
	return &r, err
}
