package bandwidths

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/structs"
)

type UpdateOpts struct {
	Bandwidth   Bandwidth    `json:"bandwidth" required:"true"`
	ExtendParam *ExtendParam `json:"extendParam,omitempty"`
}
type Bandwidth struct {
	Name string `json:"name,omitempty"`
	Size int    `json:"size,omitempty"`
}
type ExtendParam struct {
	IsAutoPay string `json:"is_auto_pay,omitempty"`
}

func (opts UpdateOpts) ToBandWidthUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

type CreateOptsBuilder interface {
	ToBandWidthCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	Name                string `json:"name" required:"true"`
	Size                *int   `json:"size" required:"true"`
	ChargeMode          string `json:"charge_mode,omitempty"`
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	PublicBorderGroup   string `json:"public_border_group,omitempty"`
	BandwidthType       string `json:"bandwidth_type,omitempty"`
}

func (opts CreateOpts) ToBandWidthCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "bandwidth")
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBandWidthCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(PostURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

type BatchCreateOptsBuilder interface {
	ToBandWidthBatchCreateMap() (map[string]interface{}, error)
}

type BatchCreateOpts struct {
	Name  string `json:"name" required:"true"`
	Size  *int   `json:"size" required:"true"`
	Count *int   `json:"count" required:"true"`
}

func (opts BatchCreateOpts) ToBandWidthBatchCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "bandwidth")
}

type BandWidthInsertOptsBuilder interface {
	ToBandWidthInsertMap() (map[string]interface{}, error)
}

type BandWidthRemoveOptsBuilder interface {
	ToBandWidthBatchRemoveMap() (map[string]interface{}, error)
}

type BandWidthInsertOpts struct {
	PublicipInfo []PublicIpInfoID `json:"publicip_info" required:"true"`
}

func (opts BandWidthInsertOpts) ToBandWidthInsertMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "bandwidth")
}

type BandWidthRemoveOpts struct {
	ChargeMode   string           `json:"charge_mode" required:"true"`
	Size         *int             `json:"size" required:"true"`
	PublicipInfo []PublicIpInfoID `json:"publicip_info" required:"true"`
}

func (opts BandWidthRemoveOpts) ToBandWidthBatchRemoveMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "bandwidth")
}

type PublicIpInfoID struct {
	PublicIPID   string `json:"publicip_id" required:"true"`
	PublicIPType string `json:"publicip_type,omitempty"`
}

func Insert(client *golangsdk.ServiceClient, bandwidthID string, opts BandWidthInsertOptsBuilder) (r CreateResult) {
	b, err := opts.ToBandWidthInsertMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(InsertURL(client, bandwidthID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

func Remove(client *golangsdk.ServiceClient, bandwidthID string, opts BandWidthRemoveOptsBuilder) (r DeleteResult) {
	b, err := opts.ToBandWidthBatchRemoveMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(RemoveURL(client, bandwidthID), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 204},
	})
	return
}

func BatchCreate(client *golangsdk.ServiceClient, opts BatchCreateOptsBuilder) (r BatchCreateResult) {
	b, err := opts.ToBandWidthBatchCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(BatchPostURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, bandwidthID string) (r DeleteResult) {
	url := DeleteURL(client, bandwidthID)
	_, r.Err = client.Delete(url, nil)
	return
}

func Update(c *golangsdk.ServiceClient, bandwidthID string, opts UpdateOpts) (r UpdateResult) {
	body, err := opts.ToBandWidthUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(UpdateURL(c, bandwidthID), body, &r.Body, nil)
	return
}

type ChangeToPeriodOpts struct {
	BandwidthIDs []string           `json:"bandwidth_ids" required:"true"`
	ExtendParam  structs.ChargeInfo `json:"extendParam" required:"true"`
}

// ChangeToPeriod is is used to change the bandwidth to prePaid billing mode
func ChangeToPeriod(c *golangsdk.ServiceClient, opts ChangeToPeriodOpts) (r ChangeResult) {
	body, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(changeToPeriodURL(c), body, &r.Body, nil)
	return
}
