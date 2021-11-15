package vaults

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	Billing             *BillingCreate     `json:"billing" required:"true"`
	Name                string             `json:"name" required:"true"`
	Resources           []ResourceCreate   `json:"resources" required:"true"`
	AutoBind            bool               `json:"auto_bind,omitempty"`
	AutoExpand          bool               `json:"auto_expand,omitempty"`
	BackupPolicyID      string             `json:"backup_policy_id,omitempty"`
	BindRules           *VaultBindRules    `json:"bind_rules,omitempty"`
	Description         string             `json:"description,omitempty"`
	EnterpriseProjectID string             `json:"enterprise_project_id,omitempty"`
	Tags                []tags.ResourceTag `json:"tags,omitempty"`
}

type BillingCreate struct {
	ConsistentLevel string                  `json:"consistent_level" required:"true"`
	ObjectType      string                  `json:"object_type" required:"true"`
	ProtectType     string                  `json:"protect_type" required:"true"`
	Size            int                     `json:"size" required:"true"`
	ChargingMode    string                  `json:"charging_mode,omitempty"`
	CloudType       string                  `json:"cloud_type,omitempty"`
	ConsoleURL      string                  `json:"console_url,omitempty"`
	ExtraInfo       *BillingCreateExtraInfo `json:"extra_info,omitempty"`
	PeriodNum       int                     `json:"period_num,omitempty"`
	PeriodType      string                  `json:"period_type,omitempty"`
	IsAutoRenew     bool                    `json:"is_auto_renew,omitempty"`
	IsAutoPay       bool                    `json:"is_auto_pay,omitempty"`
}

type BillingCreateExtraInfo struct {
	CombinedOrderECSNum int    `json:"combined_order_ecs_num,omitempty"`
	CombinedOrderID     string `json:"combined_order_id,omitempty"`
}

type ResourceCreate struct {
	ID        string             `json:"id" required:"true"`
	Type      string             `json:"type" required:"true"`
	Name      string             `json:"name,omitempty"`
	ExtraInfo *ResourceExtraInfo `json:"extra_info,omitempty"`
}

type ResourceExtraInfo struct {
	ExcludeVolumes []string                          `json:"exclude_volumes,omitempty"`
	IncludeVolumes []ResourceExtraInfoIncludeVolumes `json:"include_volumes,omitempty"`
}

type ResourceExtraInfoIncludeVolumes struct {
	ID        string `json:"id" required:"true"`
	OSVersion string `json:"os_version,omitempty"`
}

type VaultBindRules struct {
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

type CreateOptsBuilder interface {
	ToVaultCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToVaultCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "vault")
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToVaultCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, err = client.Post(rootURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	r.Err = err
	return
}

func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

type UpdateOpts struct {
	Billing    *BillingUpdate  `json:"billing,omitempty"`
	Name       string          `json:"name,omitempty"`
	AutoBind   *bool           `json:"auto_bind,omitempty"`
	BindRules  *VaultBindRules `json:"bind_rules,omitempty"`
	AutoExpand *bool           `json:"auto_expand,omitempty"`
}

type BillingUpdate struct {
	Size int `json:"size,omitempty"`
}

type UpdateOptsBuilder interface {
	ToVaultUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToVaultUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "vault")
}

func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToVaultUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type ListOpts struct {
	CloudType           string `q:"cloud_type"`
	EnterpriseProjectID string `q:"enterprise_project_id"`
	ID                  string `q:"id"`
	Limit               int    `q:"limit"`
	Name                string `q:"name"`
	ObjectType          string `q:"object_type"`
	Offset              int    `q:"offset"`
	PolicyID            string `q:"policy_id"`
	ProtectType         string `q:"protect_type"`
	ResourceIDs         string `q:"resource_ids"`
	Status              string `q:"status"`
}

func (opts ListOpts) ToPolicyListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

type ListOptsBuilder interface {
	ToPolicyListQuery() (string, error)
}

//List is a method to obtain the specified CBR vaults according to the vault ID, vault name and so on.
//This method can also obtain all the CBR vaults through the default parameter settings.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToPolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return VaultPage{pagination.SinglePageBase(r)}
	})
}

type BindPolicyOpts struct {
	PolicyID string `json:"policy_id" required:"true"`
}

func (opts BindPolicyOpts) ToBindPolicyMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

type BindPolicyOptsBuilder interface {
	ToBindPolicyMap() (map[string]interface{}, error)
}

func BindPolicy(client *golangsdk.ServiceClient, vaultID string, opts BindPolicyOptsBuilder) (r BindPolicyResult) {
	reqBody, err := opts.ToBindPolicyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(bindPolicyURL(client, vaultID), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func UnbindPolicy(client *golangsdk.ServiceClient, vaultID string, opts BindPolicyOptsBuilder) (r UnbindPolicyResult) {
	reqBody, err := opts.ToBindPolicyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(unbindPolicyURL(client, vaultID), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type AssociateResourcesOpts struct {
	Resources []ResourceCreate `json:"resources" required:"true"`
}

func (opts AssociateResourcesOpts) ToAssociateResourcesMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

type AssociateResourcesOptsBuilder interface {
	ToAssociateResourcesMap() (map[string]interface{}, error)
}

func AssociateResources(client *golangsdk.ServiceClient, vaultID string, opts AssociateResourcesOptsBuilder) (r AssociateResourcesResult) {
	reqBody, err := opts.ToAssociateResourcesMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(addResourcesURL(client, vaultID), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type DissociateResourcesOpts struct {
	ResourceIDs []string `json:"resource_ids" required:"true"`
}

func (opts DissociateResourcesOpts) ToDissociateResourcesMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

type DissociateResourcesOptsBuilder interface {
	ToDissociateResourcesMap() (map[string]interface{}, error)
}

func DissociateResources(client *golangsdk.ServiceClient, vaultID string, opts DissociateResourcesOptsBuilder) (r DissociateResourcesResult) {
	reqBody, err := opts.ToDissociateResourcesMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(removeResourcesURL(client, vaultID), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
