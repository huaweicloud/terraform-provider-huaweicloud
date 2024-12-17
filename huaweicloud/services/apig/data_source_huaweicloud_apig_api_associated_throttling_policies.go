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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/binded-throttles
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}/throttle-specials
func DataSourceApiAssociatedThrottlingPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiAssociatedThrottlingPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the throttling policies belong.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API bound to the throttling policy.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the throttling policy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the throttling policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the throttling policy.`,
			},
			"env_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the environment where the API is published.`,
			},
			"period_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The time unit for limiting the number of API calls.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All throttling policies that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the throttling policy.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the throttling policy.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the throttling policy.`,
						},
						"period_unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time unit for limiting the number of API calls.`,
						},
						"period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The period of time for limiting the number of API calls.`,
						},
						"max_api_requests": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of times an API can be accessed within a specified period.`,
						},
						"max_app_requests": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of times the API can be accessed by an app within the same period.`,
						},
						"max_ip_requests": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of times the API can be accessed by an IP address within the same period.`,
						},
						"max_user_requests": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of times the API can be accessed by a user within the same period.`,
						},
						"env_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the environment where the API is published.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the throttling policy.`,
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
						"bind_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The bind ID.`,
						},
						"bind_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time that the throttling policy is bound to the API, in RFC3339 format.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the throttling policy, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func dataSpecialThrottleSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_api_requests": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of times an API can be accessed within a specified period.`,
			},
			"throttling_object_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The object ID which the special user/application throttling policy belongs.`,
			},
			"throttling_object_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The object name which the special user/application throttling policy belongs.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the special user/application throttling policy.`,
			},
		},
	}
}

func buildListApiAssociatedThrottlingPoliciesParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("policy_id"); ok {
		res = fmt.Sprintf("%s&throttle_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&throttle_name=%v", res, v)
	}
	return res
}

func queryApiAssociatedThrottlingPolicies(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/binded-throttles?api_id={api_id}"
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{api_id}", apiId)

	queryParams := buildListApiAssociatedThrottlingPoliciesParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&limit=100&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving associated throttling policies (bound to the API: %s) under "+
				"specified dedicated instance (%s): %s", apiId, instanceId, err)
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

func dataSourceApiAssociatedThrottlingPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	policies, err := queryApiAssociatedThrottlingPolicies(client, d)
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
		d.Set("policies", filterAssociatedThrottlingPolicies(flattenAssociatedThrottlingPolicies(client, instanceId, policies), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterAssociatedThrottlingPolicies(all []interface{}, d *schema.ResourceData) []interface{} {
	result := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("env_name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("env_name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("period_unit"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("period_unit", v, nil)) {
			continue
		}

		result = append(result, v)
	}
	return result
}

func querySpecialThrottlingPolicies(client *golangsdk.ServiceClient, instanceId, policyId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/throttles/{throttle_id}/throttle-specials"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{throttle_id}", policyId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s?limit=100&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving special throttling policies: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		policies := utils.PathSearch("throttle_specials", respBody, make([]interface{}, 0)).([]interface{})
		if len(policies) < 1 {
			break
		}
		result = append(result, policies...)
		offset += len(policies)
	}
	return result, nil
}

func flattenSpecialThrottlingPolicies(specPolicies []interface{}) (userThrottles, appThrottles []interface{}) {
	// The special throttling policies contain IAM user throttles and app throttles.
	// According to the rules of append method, the maximum memory is expanded to 32,
	// and the average waste of memory is less than the waste caused by directly setting it to 30.
	for _, specPolicy := range specPolicies {
		objectType := utils.PathSearch("object_type", specPolicy, "").(string)
		throttle := map[string]interface{}{
			"max_api_requests":       utils.PathSearch("call_limits", specPolicy, ""),
			"throttling_object_id":   utils.PathSearch("object_id", specPolicy, ""),
			"throttling_object_name": utils.PathSearch("object_name", specPolicy, ""),
			"id":                     utils.PathSearch("id", specPolicy, ""),
		}
		switch objectType {
		case string(PolicyTypeApplication):
			appThrottles = append(appThrottles, throttle)
		case string(PolicyTypeUser):
			userThrottles = append(userThrottles, throttle)
		default:
			log.Printf("invalid object type, want 'APP' or 'USER', but got '%v'", objectType)
		}
	}
	return
}

func flattenAssociatedThrottlingPolicies(client *golangsdk.ServiceClient, instanceId string, policies []interface{}) []interface{} {
	if len(policies) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(policies))
	for _, policy := range policies {
		policyId := utils.PathSearch("id", policy, "").(string)
		policyDetail := map[string]interface{}{
			"id":                policyId,
			"name":              utils.PathSearch("name", policy, nil),
			"type":              *analyseThrottlingPolicyType(int(utils.PathSearch("type", policy, float64(0)).(float64))),
			"period_unit":       utils.PathSearch("time_unit", policy, nil),
			"period":            utils.PathSearch("time_interval", policy, nil),
			"max_api_requests":  utils.PathSearch("api_call_limits", policy, nil),
			"max_app_requests":  utils.PathSearch("app_call_limits", policy, nil),
			"max_ip_requests":   utils.PathSearch("ip_call_limits", policy, nil),
			"max_user_requests": utils.PathSearch("user_call_limits", policy, nil),
			"env_name":          utils.PathSearch("env_name", policy, nil),
			"description":       utils.PathSearch("remark", policy, nil),
			"bind_id":           utils.PathSearch("bind_id", policy, nil),
			"bind_time":         utils.PathSearch("bind_time", policy, nil),   // Already in RFC3339 format.
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
