package healthcheck

import (
	"github.com/huaweicloud/golangsdk"
)

type HealthCheck struct {
	HealthcheckInterval    int    `json:"healthcheck_interval"`
	ListenerID             string `json:"listener_id"`
	ID                     string `json:"id"`
	HealthcheckProtocol    string `json:"healthcheck_protocol"`
	UnhealthyThreshold     int    `json:"unhealthy_threshold"`
	UpdateTime             string `json:"update_time"`
	CreateTime             string `json:"create_time"`
	HealthcheckConnectPort int    `json:"healthcheck_connect_port"`
	HealthcheckTimeout     int    `json:"healthcheck_timeout"`
	HealthcheckUri         string `json:"healthcheck_uri"`
	HealthyThreshold       int    `json:"healthy_threshold"`
}

type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (*HealthCheck, error) {
	s := &HealthCheck{}
	return s, r.ExtractInto(s)
}

type CreateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	golangsdk.ErrResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}
