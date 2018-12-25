package antiddos

import (
	"reflect"
	"strconv"
	"time"

	"github.com/huaweicloud/golangsdk"
)

type CreateOpts struct {
	// Whether to enable L7 defense
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

type CreateOptsBuilder interface {
	ToCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, floatingIpId string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(CreateURL(client, floatingIpId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func DailyReport(client *golangsdk.ServiceClient, floatingIpId string) (r DailyReportResult) {
	url := DailyReportURL(client, floatingIpId)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, floatingIpId string) (r DeleteResult) {
	url := DeleteURL(client, floatingIpId)
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		JSONResponse: &r.Body,
		OkCodes:      []int{200},
	})
	return
}

func Get(client *golangsdk.ServiceClient, floatingIpId string) (r GetResult) {
	url := GetURL(client, floatingIpId)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func GetStatus(client *golangsdk.ServiceClient, floatingIpId string) (r GetStatusResult) {
	url := GetStatusURL(client, floatingIpId)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type GetTaskOpts struct {
	// Task ID (nonnegative integer) character string
	TaskId string `q:"task_id"`
}

type GetTaskOptsBuilder interface {
	ToGetTaskQuery() (string, error)
}

func (opts GetTaskOpts) ToGetTaskQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

func GetTask(client *golangsdk.ServiceClient, opts GetTaskOptsBuilder) (r GetTaskResult) {
	url := GetTaskURL(client)
	if opts != nil {
		query, err := opts.ToGetTaskQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func ListConfigs(client *golangsdk.ServiceClient) (r ListConfigsResult) {
	url := ListConfigsURL(client)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type ListLogsOpts struct {
	// Limit of number of returned results or the maximum number of returned results of a query. The value ranges from 1 to 100, and this parameter is used together with the offset parameter. If neither limit nor offset is used, query results of all ECSs are returned.
	Limit int `q:"limit"`

	// Offset. This parameter is valid only when used together with the limit parameter.
	Offset int `q:"offset"`

	// Possible values: desc: indicates that query results are given and sorted by time in descending order. asc: indicates that query results are given and sorted by time in ascending order.The default value is desc.
	SortDir string `q:"sort_dir"`
}

type ListLogsOptsBuilder interface {
	ToListLogsQuery() (string, error)
}

func (opts ListLogsOpts) ToListLogsQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

func ListLogs(client *golangsdk.ServiceClient, floatingIpId string, opts ListLogsOptsBuilder) (r ListLogsResult) {
	url := ListLogsURL(client, floatingIpId)
	if opts != nil {
		query, err := opts.ToListLogsQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type ListStatusOpts struct {
	// ID of an EIP
	FloatingIpId string

	// If this parameter is not used, the defense statuses of all ECSs are displayed in the Neutron-queried order by default.
	Status string `q:"status"`

	// Limit of number of returned results
	Limit int `q:"limit"`

	// Offset
	Offset int `q:"offset"`

	// IP address. Both IPv4 and IPv6 addresses are supported. For example, if you enter ?ip=192.168, the defense status of EIPs corresponding to 192.168.111.1 and 10.192.168.8 is returned.
	Ip string `q:"ip"`
}

// ListStatus returns collection of DdosStatus. It accepts a ListStatusOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListStatus(client *golangsdk.ServiceClient, opts ListStatusOpts) ([]DdosStatus, error) {
	var r ListStatusResult

	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := ListStatusURL(client) + q.String()

	_, r.Err = client.Get(u, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	allStatus, err := r.Extract()
	if err != nil {
		return nil, err
	}

	return FilterDdosStatus(allStatus, opts)
}

func FilterDdosStatus(ddosStatus []DdosStatus, opts ListStatusOpts) ([]DdosStatus, error) {

	var refinedDdosStatus []DdosStatus
	var matched bool
	m := map[string]interface{}{}

	if opts.FloatingIpId != "" {
		m["FloatingIpId"] = opts.FloatingIpId
	}

	if len(m) > 0 && len(ddosStatus) > 0 {
		for _, ddosStatus := range ddosStatus {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&ddosStatus, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedDdosStatus = append(refinedDdosStatus, ddosStatus)
			}
		}
	} else {
		refinedDdosStatus = ddosStatus
	}
	return refinedDdosStatus, nil
}

func getStructField(v *DdosStatus, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

type UpdateOpts struct {
	// Whether to enable L7 defense
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

type UpdateOptsBuilder interface {
	ToUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(client *golangsdk.ServiceClient, floatingIpId string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(UpdateURL(client, floatingIpId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type WeeklyReportOpts struct {
	// Start date of a seven-day period
	PeriodStartDate time.Time `q:""`
	//PeriodStartDate string `q:"period_start_date"`
}

type WeeklyReportOptsBuilder interface {
	ToWeeklyReportQuery() (string, error)
}

func (opts WeeklyReportOpts) ToWeeklyReportQuery() (string, error) {
	return "?period_start_date=" + strconv.FormatInt(time.Time(opts.PeriodStartDate).Unix()*1000, 10), nil //q.String(), err
}

func WeeklyReport(client *golangsdk.ServiceClient, opts WeeklyReportOptsBuilder) (r WeeklyReportResult) {
	url := WeeklyReportURL(client)
	if opts != nil {
		query, err := opts.ToWeeklyReportQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
