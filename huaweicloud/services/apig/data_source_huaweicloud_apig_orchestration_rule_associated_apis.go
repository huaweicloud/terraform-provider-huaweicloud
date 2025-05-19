package apig

import (
	"context"
	"fmt"
	"strconv"
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the orchestration rule is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the orchestration rule belongs.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the orchestration rule to be queried.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The associated API ID under the orchestration rule.`,
			},
			"apis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the associated API.`,
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the associated API.`,
						},
						"req_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request method of the associated API.`,
						},
						"req_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request URI of the associated API.`,
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The auth type of the associated API.`,
						},
						"match_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The match mode of the associated API.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The group ID to which the associated API belongs.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The group name to which the associated API belongs.`,
						},
						"attached_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attached time of the associated API.`,
						},
					},
				},
				Description: `All API list that match the filter parameters under a specified orchestration rule.`,
			},
		},
	}
}

func buildOrchestrationRuleAssociatedApisQueryParams(d *schema.ResourceData) string {
	res := ""
	if apiId, ok := d.GetOk("api_id"); ok {
		res = fmt.Sprintf("%s&api_id=%v", res, apiId)
	}
	return res
}

func listOrchestrationRuleAssociatedApis(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}/attached-apis?limit={limit}"
		instanceId = d.Get("instance_id").(string)
		ruleId     = d.Get("rule_id").(string)
		limit      = 500
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{orchestration_id}", ruleId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildOrchestrationRuleAssociatedApisQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error query API list under a specified orchestration rule (%s): %s", ruleId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		assoiciatedApis := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, assoiciatedApis...)
		pageSize := int(utils.PathSearch("size", respBody, float64(0)).(float64))
		if pageSize < limit {
			break
		}
		offset += pageSize
	}

	return result, nil
}

func flattenOrchestrationRuleAssoiciatedApis(assoiciatedApis []interface{}) []interface{} {
	if len(assoiciatedApis) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(assoiciatedApis))
	for _, assoiciatedApi := range assoiciatedApis {
		result = append(result, map[string]interface{}{
			"api_id":        utils.PathSearch("api_id", assoiciatedApi, nil),
			"api_name":      utils.PathSearch("api_name", assoiciatedApi, nil),
			"req_method":    utils.PathSearch("req_method", assoiciatedApi, nil),
			"req_uri":       utils.PathSearch("req_uri", assoiciatedApi, nil),
			"auth_type":     utils.PathSearch("auth_type", assoiciatedApi, nil),
			"match_mode":    utils.PathSearch("match_mode", assoiciatedApi, nil),
			"group_id":      utils.PathSearch("group_id", assoiciatedApi, nil),
			"group_name":    utils.PathSearch("group_name", assoiciatedApi, nil),
			"attached_time": utils.PathSearch("attached_time", assoiciatedApi, nil),
		})
	}
	return result
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

	assoiciatedApis, err := listOrchestrationRuleAssociatedApis(client, d)
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
		d.Set("apis", flattenOrchestrationRuleAssoiciatedApis(assoiciatedApis)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
