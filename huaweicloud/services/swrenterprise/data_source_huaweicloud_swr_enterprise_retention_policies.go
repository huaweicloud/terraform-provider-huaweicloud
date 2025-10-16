package swrenterprise

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/{project_id}/instances/{instance_id}/retention/policies
func DataSourceSwrEnterpriseRetentionPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseRetentionPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the namespace ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the trigger name.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the policies.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the policy ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the policy name.`,
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the namespace name.`,
						},
						"namespace_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the namespace ID`,
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the algorithm of policy.`,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the policy is enabled.`,
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the retention rules.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the priority.`,
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the policy action.`,
									},
									"template": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the template type.`,
									},
									// `params` type is Map<String, Object>
									"params": {
										Type:        schema.TypeMap,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Indicates the params.`,
									},
									"repo_scope_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the repo scope mode.`,
									},
									"tag_selectors": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Indicates the repository version selector.`,
										Elem:        dataSourceSchemaSwrEnterpriseRetentionPolicyRuleSelector(),
									},
									"scope_selectors": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Indicates the repository selectors.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `Indicates the repository selector key.`,
												},
												"value": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: `Indicates the repository selector value.`,
													Elem:        dataSourceSchemaSwrEnterpriseRetentionPolicyRuleSelector(),
												},
											},
										},
									},
									"disabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Indicates whether the policy rule is disabled.`,
									},
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the retention policy rule ID.`,
									},
								},
							},
						},
						"trigger": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the trigger config.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the trigger type.`,
									},
									"trigger_settings": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Indicates the trigger settings.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cron": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `Indicates the scheduled setting.`,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceSchemaSwrEnterpriseRetentionPolicyRuleSelector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the matching type.`,
			},
			"decoration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the selector matching type.`,
			},
			"pattern": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pattern.`,
			},
			"extras": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the extra infos.`,
			},
		},
	}
}

func dataSourceSwrEnterpriseRetentionPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}
	listRetentionPoliciesHttpUrl := "v2/{project_id}/instances/{instance_id}/retention/policies"
	listRetentionPoliciesPath := client.Endpoint + listRetentionPoliciesHttpUrl
	listRetentionPoliciesPath = strings.ReplaceAll(listRetentionPoliciesPath, "{project_id}", client.ProjectID)
	listRetentionPoliciesPath = strings.ReplaceAll(listRetentionPoliciesPath, "{instance_id}", d.Get("instance_id").(string))
	listRetentionPoliciesPath += buildSwrEnterpriseRetentionPoliciesQueryParams(d)
	listRetentionPoliciesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listRetentionPoliciesPath + fmt.Sprintf("&offset=%v", offset)
		listRetentionPoliciesResp, err := client.Request("GET", currentPath, &listRetentionPoliciesOpt)
		if err != nil {
			return diag.Errorf("error querying SWR retention policies: %s", err)
		}
		listRetentionPoliciesRespBody, err := utils.FlattenResponse(listRetentionPoliciesResp)
		if err != nil {
			return diag.Errorf("error flattening SWR retention policies response: %s", err)
		}

		policies := utils.PathSearch("retentions", listRetentionPoliciesRespBody, make([]interface{}, 0)).([]interface{})
		if len(policies) == 0 {
			break
		}
		for _, policy := range policies {
			results = append(results, map[string]interface{}{
				"id":             utils.PathSearch("id", policy, nil),
				"name":           utils.PathSearch("name", policy, nil),
				"enabled":        utils.PathSearch("enabled", policy, nil),
				"namespace_id":   utils.PathSearch("namespace_id", policy, nil),
				"namespace_name": utils.PathSearch("namespace_name", policy, nil),
				"algorithm":      utils.PathSearch("algorithm", policy, nil),
				"trigger":        flattenSwrEnterpriseRetentionPolicyTrigger(policy),
				"rules":          flattenSwrEnterpriseRetentionPolicyRetentionRules(policy),
			})
		}

		// offset must be the multiple of limit
		offset += 100
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSwrEnterpriseRetentionPoliciesQueryParams(d *schema.ResourceData) string {
	res := "?limit=100"

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("namespace_id"); ok {
		res = fmt.Sprintf("%s&namespace_id=%v", res, v)
	}

	return res
}
