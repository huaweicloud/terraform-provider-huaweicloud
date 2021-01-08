package servers

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
)

// GetResult contains the server result
type GetResult struct {
	golangsdk.Result
}

type ServerDetail struct {
	Server *Server `json:"server"`
}

// Server IEC服务实例的对外呈现结构
type Server struct {
	cloudservers.CloudServer
	CommonFiled
}

//CommonFiled is common filed
type CommonFiled struct {
	ServerID      string             `json:"origin_server_id,omitempty"`
	EdgeCloudID   string             `json:"edgecloud_id,omitempty"`
	EdgeCloudName string             `json:"edgecloud_name,omitempty"`
	Location      common.GeoLocation `json:"geolocation"`
	Operator      common.Operator    `json:"operator"`
	DomainID      string             `json:"domain_id"`
}

// ExtractServerDetail 输出边缘实例详情
func (r GetResult) ExtractServerDetail() (ServerDetail, error) {
	var serverDetail ServerDetail
	err := r.ExtractInto(&serverDetail)
	return serverDetail, err
}

// UpdateResult contains the update result
type UpdateResult struct {
	golangsdk.Result
}

// ExtractUpdateToServer extract CloudServer struct from UpdateResult
func (r UpdateResult) ExtractUpdateToServer() (interface{}, error) {
	var updateServer Server
	err := r.ExtractInto(&updateServer)
	return &updateServer, err
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type JobExecResult struct {
	commonResult
}

//执行创建image异步接口时返回的jobid结构
type Job struct {
	// job id of create image
	Id string `json:"job_id"`
}

func (r commonResult) ExtractJob() (Job, error) {
	var j Job
	err := r.ExtractInto(&j)
	return j, err
}

//Server struct is a serverIds structure returned when you create a server
type ServerIDs struct {
	IDs []string `json:"server_ids"`
}

//ExtractServer is used to extract server struct in response
func (r commonResult) ExtractServer() (ServerIDs, error) {
	var s ServerIDs
	err := r.ExtractInto(&s)
	return s, err
}

//CreateCloudServerResponse is a structure for creating server return values
type CreateCloudServerResponse struct {
	Job       Job
	ServerIDs ServerIDs
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
