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
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"access_scope": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"namespaces": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
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
			func(accessPolicy gjson.Result) any {
				return map[string]any{
					"kind":        accessPolicy.Get("kind").Value(),
					"api_version": accessPolicy.Get("apiVersion").Value(),
					"name":        accessPolicy.Get("name").Value(),
					"policy_id":   accessPolicy.Get("policyId").Value(),
					"clusters":    schemas.SliceToStrList(accessPolicy.Get("clusters")),
					"access_scope": schemas.ObjectToList(accessPolicy.Get("accessScope"),
						func(accessScope gjson.Result) any {
							return map[string]any{
								"namespaces": schemas.SliceToStrList(accessScope.Get("namespaces")),
							}
						},
					),
					"policy_type": accessPolicy.Get("policyType").Value(),
					"principal": schemas.ObjectToList(accessPolicy.Get("principal"),
						func(principal gjson.Result) any {
							return map[string]any{
								"type": principal.Get("type").Value(),
								"ids":  schemas.SliceToStrList(principal.Get("ids")),
							}
						},
					),
					"create_time": accessPolicy.Get("createTime").Value(),
					"update_time": accessPolicy.Get("updateTime").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
