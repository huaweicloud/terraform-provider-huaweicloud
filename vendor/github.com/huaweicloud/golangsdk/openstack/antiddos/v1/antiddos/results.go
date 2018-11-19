package antiddos

import (
	"encoding/json"
	"time"

	"github.com/huaweicloud/golangsdk"
)

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*CreateResponse, error) {
	var response CreateResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type CreateResponse struct {
	// Internal error code
	ErrorCode string `json:"error_code,"`

	// Internal error description
	ErrorDescription string `json:"error_description,"`

	// ID of a task. This ID can be used to query the status of the task. This field is reserved for use in task auditing later. It is temporarily unused.
	TaskId string `json:"task_id,"`
}

type DailyReportResult struct {
	commonResult
}

func (r DailyReportResult) Extract() ([]Data, error) {
	var s DailyReportResponse
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Data, nil
}

type DailyReportResponse struct {
	// Traffic in the last 24 hours
	Data []Data `json:"data"`
}
type Data struct {
	// Start time
	PeriodStart int `json:"period_start,"`

	// Inbound traffic (bit/s)
	BpsIn int `json:"bps_in,"`

	// Attack traffic (bit/s)
	BpsAttack int `json:"bps_attack,"`

	// Total traffic
	TotalBps int `json:"total_bps,"`

	// Inbound packet rate (number of packets per second)
	PpsIn int `json:"pps_in,"`

	// Attack packet rate (number of packets per second)
	PpsAttack int `json:"pps_attack,"`

	// Total packet rate
	TotalPps int `json:"total_pps,"`
}

type DeleteResult struct {
	commonResult
}

func (r DeleteResult) Extract() (*DeleteResponse, error) {
	var response DeleteResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type DeleteResponse struct {
	// Internal error code
	ErrorCode string `json:"error_code,"`

	// Internal error description
	ErrorDescription string `json:"error_description,"`

	// ID of a task. This ID can be used to query the status of the task. This field is reserved for use in task auditing later. It is temporarily unused.
	TaskId string `json:"task_id,"`
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*GetResponse, error) {
	var response GetResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type GetResponse struct {
	// Whether L7 defense has been enabled
	EnableL7 bool `json:"enable_L7,"`

	// Position ID of traffic. The value ranges from 1 to 9.
	TrafficPosId int `json:"traffic_pos_id,"`

	// Position ID of number of HTTP requests. The value ranges from 1 to 15.
	HttpRequestPosId int `json:"http_request_pos_id,"`

	// Position ID of access limit during cleaning. The value ranges from 1 to 8.
	CleaningAccessPosId int `json:"cleaning_access_pos_id,"`

	// Application type ID. Possible values: 0 1
	AppTypeId int `json:"app_type_id,"`
}

type GetStatusResult struct {
	commonResult
}

func (r GetStatusResult) Extract() (*GetStatusResponse, error) {
	var response GetStatusResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type GetStatusResponse struct {
	// Defense status
	Status string `json:"status,"`
}

type GetTaskResult struct {
	commonResult
}

func (r GetTaskResult) Extract() (*GetTaskResponse, error) {
	var response GetTaskResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type GetTaskResponse struct {
	// Status of a task, which can be one of the following: success, failed, waiting, running, preprocess, ready
	TaskStatus string `json:"task_status,"`

	// Additional information about a task
	TaskMsg string `json:"task_msg,"`
}

type ListConfigsResult struct {
	commonResult
}

func (r ListConfigsResult) Extract() (*ListConfigsResponse, error) {
	var response ListConfigsResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type ListConfigsResponse struct {
	// List of traffic limits
	TrafficLimitedList []struct {
		// Position ID of traffic
		TrafficPosId int `json:"traffic_pos_id,"`

		// Threshold of traffic per second (Mbit/s)
		TrafficPerSecond int `json:"traffic_per_second,"`

		// Threshold of number of packets per second
		PacketPerSecond int `json:"packet_per_second,"`
	} `json:"traffic_limited_list,"`

	// List of HTTP limits
	HttpLimitedList []struct {
		// Position ID of number of HTTP requests
		HttpRequestPosId int `json:"http_request_pos_id,"`

		// Threshold of number of HTTP requests per second
		HttpPacketPerSecond int `json:"http_packet_per_second,"`
	} `json:"http_limited_list,"`

	// List of limits of numbers of connections
	ConnectionLimitedList []struct {
		// Position ID of access limit during cleaning
		CleaningAccessPosId int `json:"cleaning_access_pos_id,"`

		// Position ID of access limit during cleaning
		NewConnectionLimited int `json:"new_connection_limited,"`

		// Position ID of access limit during cleaning
		TotalConnectionLimited int `json:"total_connection_limited,"`
	} `json:"connection_limited_list,"`
}

type ListLogsResult struct {
	commonResult
}

func (r ListLogsResult) Extract() ([]Logs, error) {
	var s ListLogsResponse
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Logs, nil
}

type ListLogsResponse struct {
	// Total number of EIPs
	Total int `json:"total,"`

	// List of events
	Logs []Logs `json:"logs,"`
}

type Logs struct {
	// Start time
	StartTime int `json:"start_time,"`

	// End time
	EndTime int `json:"end_time,"`

	// Defense status, the possible value of which is one of the following: 1: indicates that traffic cleaning is underway. 2: indicates that traffic is discarded.
	Status int `json:"status,"`

	// Traffic at the triggering point.
	TriggerBps int `json:"trigger_bps,"`

	// Packet rate at the triggering point
	TriggerPps int `json:"trigger_pps,"`

	// HTTP request rate at the triggering point
	TriggerHttpPps int `json:"trigger_http_pps,"`
}
type ListStatusResult struct {
	commonResult
}

// Extract is a function that accepts a ListStatusOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func (r commonResult) Extract() ([]DdosStatus, error) {
	var s ListStatusResponse
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.DdosStatus, nil
}

type ListStatusResponse struct {
	// Total number of EIPs
	Total int `json:"total,"`

	// List of defense statuses
	DdosStatus []DdosStatus `json:"ddosStatus,"`
}

type DdosStatus struct {
	// Floating IP address
	FloatingIpAddress string `json:"floating_ip_address,"`

	// ID of an EIP
	FloatingIpId string `json:"floating_ip_id,"`

	// EIP type.
	NetworkType string `json:"network_type,"`

	// Defense status
	Status string `json:"status,"`
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*UpdateResponse, error) {
	var response UpdateResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type UpdateResponse struct {
	// Internal error code
	ErrorCode string `json:"error_code,"`

	// Internal error description
	ErrorDescription string `json:"error_description,"`

	// ID of a task. This ID can be used to query the status of the task. This field is reserved for use in task auditing later. It is temporarily unused.
	TaskId string `json:"task_id,"`
}

type WeeklyReportResult struct {
	commonResult
}

func (r WeeklyReportResult) Extract() (*WeeklyReportResponse, error) {
	var response WeeklyReportResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type WeekData struct {
	// Number of DDoS attacks intercepted
	DdosInterceptTimes int `json:"ddos_intercept_times,"`

	// Number of DDoS blackholes
	DdosBlackholeTimes int `json:"ddos_blackhole_times,"`

	// Maximum attack traffic
	MaxAttackBps int `json:"max_attack_bps,"`

	// Maximum number of attack connections
	MaxAttackConns int `json:"max_attack_conns,"`

	// Start date
	PeriodStartDate time.Time `json:"period_start_date,"`
}

type WeeklyReportResponse struct {
	// Number of DDoS attacks intercepted in a week
	DdosInterceptTimes int `json:"ddos_intercept_times,"`

	// Number of DDoS attacks intercepted in a week
	Weekdata []WeekData `json:"-"`

	// Top 10 attacked IP addresses
	Top10 []struct {
		// EIP
		FloatingIpAddress string `json:"floating_ip_address,"`

		// Number of DDoS attacks intercepted, including cleaning operations and blackholes
		Times int `json:"times,"`
	} `json:"top10,"`
}

func (r *WeeklyReportResponse) UnmarshalJSON(b []byte) error {
	type tmp WeeklyReportResponse
	var s struct {
		tmp
		Weekdata []struct {
			// Number of DDoS attacks intercepted
			DdosInterceptTimes int `json:"ddos_intercept_times,"`

			// Number of DDoS blackholes
			DdosBlackholeTimes int `json:"ddos_blackhole_times,"`

			// Maximum attack traffic
			MaxAttackBps int `json:"max_attack_bps,"`

			// Maximum number of attack connections
			MaxAttackConns int `json:"max_attack_conns,"`

			// Start date
			PeriodStartDate int64 `json:"period_start_date,"`
		} `json:"weekdata,"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = WeeklyReportResponse(s.tmp)
	r.Weekdata = make([]WeekData, len(s.Weekdata))

	for idx, val := range s.Weekdata {
		r.Weekdata[idx] = WeekData{
			DdosInterceptTimes: val.DdosBlackholeTimes,
			DdosBlackholeTimes: val.DdosBlackholeTimes,
			MaxAttackBps:       val.MaxAttackBps,
			MaxAttackConns:     val.MaxAttackConns,
			PeriodStartDate:    time.Unix(val.PeriodStartDate/1000, 0).UTC(),
		}
	}

	return nil
}
