package connection_monitors

import (
	"github.com/chnsz/golangsdk"
)

type CommonResult struct {
	golangsdk.Result
}

type ConnectionMonitor struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	VpnConnectionId string `json:"vpn_connection_id"`
	Type            string `json:"type"`
	SourceIp        string `json:"source_ip"`
	DestinationIp   string `json:"destination_ip"`
	ProtoType       string `json:"proto_type"`
}

type ListResp struct {
	ConnectionMonitors []ConnectionMonitor `json:"connection_monitors"`
	RequestId          string              `json:"request_id"`
}

func (r CommonResult) Extract() ([]ConnectionMonitor, error) {
	var s *ListResp
	err := r.ExtractInto(&s)
	return s.ConnectionMonitors, err
}
