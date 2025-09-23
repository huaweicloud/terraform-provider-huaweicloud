package autoscaling

import "github.com/chnsz/golangsdk"

type UpdateResponse struct {
	InstanceId   string       `json:"instance_id"`
	InstanceName string       `json:"instance_name"`
	SwitchStatus SwitchStatus `json:"switch_status"`
}

type SwitchStatus struct {
	ScalingSwitch  string `json:"scaling_switch"`
	FlavorSwitch   string `json:"flavor_switch"`
	ReadOnlySwitch string `json:"read_only_switch"`
}

type UpdateResult struct {
	golangsdk.Result
}

func (r UpdateResult) ExtractUpdateResponse() (*UpdateResponse, error) {
	var updateResponse UpdateResponse
	err := r.ExtractInto(&updateResponse)
	return &updateResponse, err
}

type AutoScaling struct {
	Id               string          `json:"id"`
	InstanceId       string          `json:"instance_id"`
	InstanceName     string          `json:"instance_name"`
	Status           string          `json:"status"`
	ScalingStrategy  ScalingStrategy `json:"scaling_strategy"`
	MonitorCycle     int             `json:"monitor_cycle"`
	SilenceCycle     int             `json:"silence_cycle"`
	EnlargeThreshold int             `json:"enlarge_threshold"`
	MaxFavor         string          `json:"max_flavor"`
	ReduceEnabled    bool            `json:"reduce_enabled"`
	MinFlavor        string          `json:"min_flavor"`
	SilenceStartAt   string          `json:"silence_start_at"`
	MaxReadOnlyCount int             `json:"max_read_only_count"`
	MinReadOnlyCount int             `json:"min_read_only_count"`
	ReadOnlyWeight   int             `json:"read_only_weight"`
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*AutoScaling, error) {
	var autoScaling AutoScaling
	err := r.ExtractInto(&autoScaling)
	return &autoScaling, err
}
