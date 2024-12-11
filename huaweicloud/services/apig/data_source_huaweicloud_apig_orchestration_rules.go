package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/orchestrations
func DataSourceOrchestrationRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrchestrationRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the orchestration rules are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the orchestration rules belong.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the orchestration rule to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the orchestration rule to be queried, fuzzy matching is supported.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All orchestration rules that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the orchestration rule.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the orchestration rule.`,
						},
						"strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the orchestration rule.`,
						},
						"is_preprocessing": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether rule is a preprocessing rule.`,
						},
						"mapped_param": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The parameter configuration after orchestration, in JSON format.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the orchestration rule, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the orchestration rule, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func buildOrchestrationRulesQueryParams(d *schema.ResourceData) string {
	res := ""
	if signId, ok := d.GetOk("rule_id"); ok {
		res = fmt.Sprintf("%s&orchestration_id=%v", res, signId)
	}
	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&orchestration_name=%v", res, name)
	}
	return res
}

func listOrchestrationRules(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/orchestrations?limit=100"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath += buildOrchestrationRulesQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error query orchestration rules under specified dedicated instance (%s): %s",
				instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		rules := utils.PathSearch("orchestrations", respBody, make([]interface{}, 0)).([]interface{})
		if len(rules) < 1 {
			break
		}
		result = append(result, rules...)
		offset += len(rules)
	}
	return result, nil
}

func dataSourceOrchestrationRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	rules, err := listOrchestrationRules(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenOrchestrationRules(rules)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOrchestrationRules(rules []interface{}) []interface{} {
	if len(rules) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("orchestration_id", rule, nil),
			"name":             utils.PathSearch("orchestration_name", rule, nil),
			"strategy":         utils.PathSearch("orchestration_strategy", rule, nil),
			"is_preprocessing": utils.PathSearch("is_preprocessing", rule, nil),
			"mapped_param": utils.JsonToString(utils.PathSearch("orchestration_mapped_param",
				rule, make(map[string]interface{}))),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("orchestration_create_time", rule, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("orchestration_update_time", rule, "").(string))/1000, false),
		})
	}
	return result
}
