package apig

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttles
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}/throttle-specials
func DataSourceThrottlingPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceThrottlingPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the dedicated instance to which the throttling policies belong.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the throttling policy.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the throttling policy.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the throttling policy.",
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the throttling policy.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the throttling policy.",
						},
						"period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The period of time for limiting the number of API calls.",
						},
						"period_unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time unit for limiting the number of API calls.",
						},
						"max_api_requests": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of times an API can be accessed within a specified period.",
						},
						"max_app_requests": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of times the API can be accessed by an app within the same period.",
						},
						"max_ip_requests": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of times the API can be accessed by an IP address within the same period.",
						},
						"max_user_requests": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of times the API can be accessed by a user within the same period.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the throttling policy.",
						},
						"user_throttles": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataSpecialThrottleSchemaResource(),
							Description: "The array of one or more special throttling policies for IAM user limit.",
						},
						"app_throttles": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataSpecialThrottleSchemaResource(),
							Description: "The array of one or more special throttling policies for APP limit.",
						},
						"bind_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of APIs bound to the throttling policy.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of throttling policy.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the throttling policy, in RFC3339 format.",
						},
					},
				},
			},
		},
	}
}

func buildListThrottlingPoliciesParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("policy_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	return res
}

func getThrottlingPolicies(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/throttles?limit=500"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	queryParams := buildListThrottlingPoliciesParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving throttling policies under specified dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		policies := utils.PathSearch("throttles", respBody, make([]interface{}, 0)).([]interface{})
		if len(policies) < 1 {
			break
		}
		result = append(result, policies...)
		offset += len(policies)
	}
	return result, nil
}

func dataSourceThrottlingPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	policies, err := getThrottlingPolicies(client, d)
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
		d.Set("policies", filterThrottlingPolicies(flattenThrottlingPolicies(client, instanceId, policies), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterThrottlingPolicies(all []interface{}, d *schema.ResourceData) []interface{} {
	result := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}

		result = append(result, v)
	}
	return result
}

func flattenThrottlingPolicies(client *golangsdk.ServiceClient, instanceId string, policies []interface{}) []interface{} {
	if len(policies) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(policies))
	for _, policy := range policies {
		policyId := utils.PathSearch("id", policy, "").(string)
		policyDetail := map[string]interface{}{
			"id":                policyId,
			"name":              utils.PathSearch("name", policy, nil),
			"period":            utils.PathSearch("time_interval", policy, nil),
			"period_unit":       utils.PathSearch("time_unit", policy, nil),
			"max_api_requests":  utils.PathSearch("api_call_limits", policy, nil),
			"max_app_requests":  utils.PathSearch("app_call_limits", policy, nil),
			"max_ip_requests":   utils.PathSearch("ip_call_limits", policy, 0),
			"max_user_requests": utils.PathSearch("user_call_limits", policy, 0),
			"type":              *analyseThrottlingPolicyType(int(utils.PathSearch("type", policy, float64(0)).(float64))),
			"description":       utils.PathSearch("remark", policy, nil),
			"bind_num":          utils.PathSearch("bind_num", policy, nil),
			"created_at":        utils.PathSearch("create_time", policy, nil), // Already in RFC3339 format.
		}

		if int(utils.PathSearch("is_inclu_special_throttle", policy, float64(0)).(float64)) == includeSpecialThrottle {
			// Get related special throttling policies.
			specResp, err := querySpecialThrottlingPolicies(client, instanceId, policyId)
			if err != nil {
				log.Printf("[ERROR] unable to find the special throttles from policy: %s", err)
			} else {
				userThrottles, appThrottles := flattenSpecialThrottlingPolicies(specResp)
				policyDetail["user_throttles"] = userThrottles
				policyDetail["app_throttles"] = appThrottles
			}
		}
		result = append(result, policyDetail)
	}
	return result
}
