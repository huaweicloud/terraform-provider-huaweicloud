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

// @API SWR GET /v2/{project_id}/instances/{instance_id}/webhook/policies
func DataSourceSwrEnterpriseTriggers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseTriggersRead,

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
			"order_column": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the order column.`,
			},
			"order_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the order type.`,
			},
			"triggers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the triggers.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the trigger ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the trigger name.`,
						},
						"targets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the target params.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the trigger type.`,
									},
									"address_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the trigger address type.`,
									},
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the trigger address.`,
									},
									"auth_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the auth header.`,
									},
									"skip_cert_verify": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Indicates whether to skip the verification of the certificate.`,
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
													Elem:        dataSourceSchemaSwrEnterpriseTriggerRuleSelector(),
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
										Elem:        dataSourceSchemaSwrEnterpriseTriggerRuleSelector(),
									},
								},
							},
						},
						"event_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the event types of trigger.`,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the trigger is enabled.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description of trigger.`,
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
						"namespace_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the namespace name.`,
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

func dataSourceSchemaSwrEnterpriseTriggerRuleSelector() *schema.Resource {
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
		},
	}
}

func dataSourceSwrEnterpriseTriggersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}
	listTriggersHttpUrl := "v2/{project_id}/instances/{instance_id}/webhook/policies"
	listTriggersPath := client.Endpoint + listTriggersHttpUrl
	listTriggersPath = strings.ReplaceAll(listTriggersPath, "{project_id}", client.ProjectID)
	listTriggersPath = strings.ReplaceAll(listTriggersPath, "{instance_id}", d.Get("instance_id").(string))
	listTriggersPath += buildSwrEnterpriseTriggersQueryParams(d)
	listTriggersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listTriggersPath + fmt.Sprintf("&offset=%v", offset)
		listTriggersResp, err := client.Request("GET", currentPath, &listTriggersOpt)
		if err != nil {
			return diag.Errorf("error querying SWR triggers: %s", err)
		}
		listTriggersRespBody, err := utils.FlattenResponse(listTriggersResp)
		if err != nil {
			return diag.Errorf("error flattening SWR triggers response: %s", err)
		}

		triggers := utils.PathSearch("policies", listTriggersRespBody, make([]interface{}, 0)).([]interface{})
		if len(triggers) == 0 {
			break
		}
		for _, trigger := range triggers {
			results = append(results, map[string]interface{}{
				"id":             utils.PathSearch("id", trigger, nil),
				"name":           utils.PathSearch("name", trigger, nil),
				"enabled":        utils.PathSearch("enabled", trigger, nil),
				"description":    utils.PathSearch("description", trigger, nil),
				"event_types":    utils.PathSearch("event_types", trigger, nil),
				"namespace_id":   utils.PathSearch("namespace_id", trigger, nil),
				"namespace_name": utils.PathSearch("namespace_name", trigger, nil),
				"creator":        utils.PathSearch("creator", trigger, nil),
				"created_at":     utils.PathSearch("created_at", trigger, nil),
				"updated_at":     utils.PathSearch("updated_at", trigger, nil),
				"targets":        flattenSwrEnterpriseTriggerTargets(trigger),
				"scope_rules":    flattenSwrEnterpriseTriggerScopeRules(trigger),
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
		d.Set("triggers", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSwrEnterpriseTriggersQueryParams(d *schema.ResourceData) string {
	res := "?limit=100"

	if v, ok := d.GetOk("order_column"); ok {
		res = fmt.Sprintf("%s&order_column=%v", res, v)
	}
	if v, ok := d.GetOk("order_type"); ok {
		res = fmt.Sprintf("%s&order_type=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("namespace_id"); ok {
		res = fmt.Sprintf("%s&namespace_id=%v", res, v)
	}

	return res
}
