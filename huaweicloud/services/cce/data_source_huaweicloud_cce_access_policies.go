package cce

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/access-policies
func DataSourceCCEAccessPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEAccessPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_policy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type AccessPoliciesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCCEAccessPoliciesDSWrapper(d *schema.ResourceData, meta interface{}) *AccessPoliciesDSWrapper {
	return &AccessPoliciesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCCEAccessPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCCEAccessPoliciesDSWrapper(d, meta)
	accessPoliciesRst, err := wrapper.getAccessPolicies()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.accessPoliciesToSchema(accessPoliciesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /api/v3/access-policies
func (w *AccessPoliciesDSWrapper) getAccessPolicies() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/api/v3/access-policies"
	params := map[string]any{
		"cluster_id": w.Get("cluster_id"),
	}
	params = utils.RemoveNil(params)

	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *AccessPoliciesDSWrapper) accessPoliciesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("access_policy_list", schemas.SliceToList(body.Get("accessPolicyList"),
			func(accesspolicy gjson.Result) any {
				return map[string]any{
					"name":        accesspolicy.Get("name").Value(),
					"policy_id":   accesspolicy.Get("policyId").Value(),
					"policy_type": accesspolicy.Get("policyType").Value(),
					"create_time": accesspolicy.Get("createTime").Value(),
					"update_time": accesspolicy.Get("updateTime").Value(),
					"clusters":    schemas.SliceToStrList(accesspolicy.Get("clusters")),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
