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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}/attached-apis
func DataSourceOrchestrationRuleAssociatedApis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrchestrationRuleAssociatedApisRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the orchestration rule belongs.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the orchestration rule.",
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the API associated with the orchestration rule.",
			},
			"api_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the API associated with the orchestration rule, fuzzy matching is supported.",
			},
			"apis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API.`,
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the API.`,
						},
						"req_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request address of the API.`,
						},
						"req_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request method of the API.`,
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The security authentication mode of the API request.`,
						},
						"match_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The matching mode of the API.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API group to which the API belongs.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the API group to which the API belongs.`,
						},
						"attached_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the orchestration rule is associated with the API, in RFC3339 format.`,
						},
					},
				},
				Description: `All associated APIs that match the filter parameters.`,
			},
		},
	}
}

func buildOrchestrationRuleAssociatedApisQueryParams(d *schema.ResourceData) string {
	res := ""
	if apiId, ok := d.GetOk("api_id"); ok {
		res = fmt.Sprintf("%s&api_id=%v", res, apiId)
	}

	if apiName, ok := d.GetOk("api_name"); ok {
		res = fmt.Sprintf("%s&api_name=%v", res, apiName)
	}
	return res
}

func listOrchestrationRuleAssociatedApis(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl             = "v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}/attached-apis"
		instanceId          = d.Get("instance_id").(string)
		orchestrationRuleId = d.Get("rule_id").(string)
		limit               = 500
		offset              = 0
		result              = make([]interface{}, 0)
	)

	listPath := client.Endpoint + fmt.Sprintf("%s?limit=%d", httpUrl, limit)
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{orchestration_id}", orchestrationRuleId)
	listPath += buildOrchestrationRuleAssociatedApisQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error query the APIs associated with the orchestration rule (%s) under specified dedicated instance (%s): %s",
				orchestrationRuleId, instanceId, err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		rules := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, rules...)
		if len(rules) < limit {
			break
		}
		offset += len(rules)
	}
	return result, nil
}

func dataSourceOrchestrationRuleAssociatedApisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	rules, err := listOrchestrationRuleAssociatedApis(client, d)
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
		d.Set("apis", flattenOrchestrationRuleAssociatedApis(rules)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOrchestrationRuleAssociatedApis(apis []interface{}) []interface{} {
	if len(apis) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(apis))
	for _, rule := range apis {
		result = append(result, map[string]interface{}{
			"api_id":     utils.PathSearch("api_id", rule, nil),
			"api_name":   utils.PathSearch("api_name", rule, nil),
			"req_uri":    utils.PathSearch("req_uri", rule, nil),
			"req_method": utils.PathSearch("req_method", rule, nil),
			"auth_type":  utils.PathSearch("auth_type", rule, nil),
			"match_mode": utils.PathSearch("match_mode", rule, nil),
			"group_id":   utils.PathSearch("group_id", rule, nil),
			"group_name": utils.PathSearch("group_name", rule, nil),
			"attached_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("attached_time", rule, "").(string))/1000, false),
		})
	}
	return result
}
