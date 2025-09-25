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

// @API SWR GET /v2/{project_id}/instances/{instance_id}/signature/policies
func DataSourceSwrEnterpriseImageSignaturePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseImageSignaturePoliciesRead,

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
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the image signature policies.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the policy ID.`,
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the namespace name.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the policy name.`,
						},
						"signature_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the signature method.`,
						},
						"signature_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the signature algorithm.`,
						},
						"signature_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the signature key.`,
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
						"scope_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the scope rules`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
													Elem:        dataSourceSchemaSwrEnterpriseImageSignaturePolicyRuleSelector(),
												},
											},
										},
									},
									"repo_scope_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the repository select method.`,
									},
									"tag_selectors": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Indicates the repository version selector.`,
										Elem:        dataSourceSchemaSwrEnterpriseImageSignaturePolicyRuleSelector(),
									},
								},
							},
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the policy is enabled.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description of policy.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator`,
						},
						"namespace_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the namespace ID`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last update time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSchemaSwrEnterpriseImageSignaturePolicyRuleSelector() *schema.Resource {
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

func dataSourceSwrEnterpriseImageSignaturePoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}
	listHttpUrl := "v2/{project_id}/instances/{instance_id}/signature/policies?limit=100"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listPath + fmt.Sprintf("&offset=%v", offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error querying SWR image signature policies: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.Errorf("error flattening SWR image signature policies response: %s", err)
		}

		policies := utils.PathSearch("policies", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(policies) == 0 {
			break
		}
		for _, policy := range policies {
			results = append(results, map[string]interface{}{
				"id":                  utils.PathSearch("id", policy, nil),
				"namespace_name":      utils.PathSearch("namespace_name", policy, nil),
				"name":                utils.PathSearch("name", policy, nil),
				"enabled":             utils.PathSearch("enabled", policy, nil),
				"description":         utils.PathSearch("description", policy, nil),
				"signature_method":    utils.PathSearch("signature_method", policy, nil),
				"signature_algorithm": utils.PathSearch("signature_algorithm", policy, nil),
				"signature_key":       utils.PathSearch("signature_key", policy, nil),
				"namespace_id":        utils.PathSearch("namespace_id", policy, nil),
				"creator":             utils.PathSearch("creator", policy, nil),
				"created_at":          utils.PathSearch("created_at", policy, nil),
				"updated_at":          utils.PathSearch("updated_at", policy, nil),
				"trigger":             flattenSwrEnterpriseImageSignaturePolicyTrigger(policy),
				"scope_rules":         flattenSwrEnterpriseImageSignaturePolicyScopeRules(policy),
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
