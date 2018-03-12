package backendecs

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
)

type Listener struct {
	ID string `json:"id"`
}

type Backend struct {
	ServerAddress string     `json:"server_address"`
	ID            string     `json:"id"`
	Address       string     `json:"address"`
	Status        string     `json:"status"`
	HealthStatus  string     `json:"health_status"`
	UpdateTime    string     `json:"update_time"`
	CreateTime    string     `json:"create_time"`
	ServerName    string     `json:"server_name"`
	ServerID      string     `json:"server_id"`
	Listeners     []Listener `json:"listeners"`
}

// GetResult represents the result of a get operation.
type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*Backend, error) {
	var s []Backend
	err := r.ExtractInto(&s)
	if err == nil {
		if len(s) != 1 {
			return nil, fmt.Errorf("get %d backends: %#v", len(s), s)
		}
		return &(s[0]), nil
	}
	return nil, err
}
