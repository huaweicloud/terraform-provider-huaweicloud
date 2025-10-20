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

// @API SWR GET /v2/{project_id}/instances/{instance_id}/immutabletagrules
func DataSourceSwrEnterpriseImmutableTagRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseImmutableTagRulesRead,

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
				Description: `Specifies the enterprise instance namespace ID.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the immutable tag rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the immutable tag rule ID`,
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
						"tag_selectors": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the repository version selector.`,
							Elem:        dataSourceSchemaSwrEnterpriseImmutableTagRuleRuleSelector(),
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
										Elem:        dataSourceSchemaSwrEnterpriseImmutableTagRuleRuleSelector(),
									},
								},
							},
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the priority.`,
						},
						"disabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the policy rule is disabled.`,
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
					},
				},
			},
		},
	}
}

func dataSourceSchemaSwrEnterpriseImmutableTagRuleRuleSelector() *schema.Resource {
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

func dataSourceSwrEnterpriseImmutableTagRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	listHttpUrl := "v2/{project_id}/instances/{instance_id}/immutabletagrules"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildSwrEnterpriseImmutableTagRulesQueryParams(d)
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
			return diag.Errorf("error querying SWR immutable tag rules: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.Errorf("error flattening SWR immutable tag rules response: %s", err)
		}

		rules := utils.PathSearch("immutable_rules", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(rules) == 0 {
			break
		}
		for _, rule := range rules {
			results = append(results, map[string]interface{}{
				"id":             utils.PathSearch("id", rule, nil),
				"namespace_name": utils.PathSearch("namespace_name", rule, nil),
				"priority":       utils.PathSearch("priority", rule, nil),
				"disabled":       utils.PathSearch("disabled", rule, nil),
				"action":         utils.PathSearch("action", rule, nil),
				"template":       utils.PathSearch("template", rule, nil),
				"namespace_id":   utils.PathSearch("namespace_id", rule, nil),
				"tag_selectors": flattenSwrEnterpriseImmutableTagRuleScopeRulesRuleSelector(
					utils.PathSearch("tag_selectors", rule, make([]interface{}, 0))),
				"scope_selectors": flattenSwrEnterpriseImmutableTagRuleScopeRulesScopeSelectors(
					utils.PathSearch("scope_selectors", rule, nil)),
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
		d.Set("rules", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSwrEnterpriseImmutableTagRulesQueryParams(d *schema.ResourceData) string {
	res := "?limit=100"

	if v, ok := d.GetOk("namespace_id"); ok {
		res = fmt.Sprintf("%s&namespace_id=%v", res, v)
	}

	return res
}
