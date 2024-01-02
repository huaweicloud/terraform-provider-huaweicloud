package alarmreminding

import (
	"github.com/chnsz/golangsdk"
)

func GetWarnAlert(client *golangsdk.ServiceClient) (r WarnAlertResult) {
	url := GetWarnAlertURL(client)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

type UpdateOps struct {
	// The topic urn.
	TopicUrn string `json:"topic_urn"`

	// The display name of alert config.
	DisplayName string `json:"display_name"`

	// Alarm configuration.
	WarnConfig *WarnConfig `json:"warn_config"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type WarnConfig struct {
	// Whether to enable alarm config
	EnableAntiDDoS *bool `json:"antiDDoS,omitempty"`

	// Deprecated
	// Brute force cracking (system logins, FTP, and DB)
	BruceForce *bool `json:"bruce_force,omitempty"`

	// Deprecated
	// Alarms about remote logins
	RemoteLogin *bool `json:"remote_login,omitempty"`

	// Deprecated
	// Weak passwords (system and database)
	WeakPassword *bool `json:"weak_password,omitempty"`

	// Deprecated
	// Overly high rights of a database process
	HighPrivilege *bool `json:"high_privilege,omitempty"`

	// Deprecated
	// Web page backdoors
	BackDoors *bool `json:"back_doors,omitempty"`

	// Deprecated
	// Reserved
	Waf *bool `json:"waf,omitempty"`

	// Deprecated
	// Possible values: 0: indicates that alarms are sent once a day. 1: indicates that alarms are sent once every half hour. This parameter is mandatory for the Host Intrusion Detection (HID) service.
	SendFrequency *int `json:"send_frequency,omitempty"`
}

type UpdateOpsBuilder interface {
	ToUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOps) ToUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateWarnAlert(client *golangsdk.ServiceClient, opts UpdateOpsBuilder) (r UpdateResult) {
	reuestBody, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(UpdateWarnAlertURL(client), reuestBody, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}
