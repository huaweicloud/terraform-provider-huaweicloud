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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis
func DataSourceApiBasicConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiBasicConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the dedicated instance to which the APIs belong.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the API.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the API.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the API group to which the APIs belong.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the API.`,
			},
			"request_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the request method of the API.",
			},
			"request_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the request address of the API.",
			},
			"request_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the request protocol of the API.",
			},
			"security_authentication": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the security authentication mode of the API request.",
			},
			"vpc_channel_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the VPC channel.",
			},
			"precise_search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the parameter name for exact matching.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the environment where the API is published.`,
			},
			"env_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the environment where the API is published.`,
			},
			"backend_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the backend type of the API.`,
			},
			"configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All API configurations that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the API.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the API.`,
						},
						"request_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The request method of the API.",
						},
						"request_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The request address of the API.",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The request protocol of the API.",
						},
						"security_authentication": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security authentication mode of the API request.",
						},
						"simple_authentication": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the authentication of the application code is enabled.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the authorizer to which the API request used.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of group corresponding to the API.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of group corresponding to the API.`,
						},
						"group_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of group corresponding to the API.`,
						},
						"env_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the environment where the API is published.`,
						},
						"env_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the environment where the API is published.`,
						},
						"publish_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of publish corresponding to the API.`,
						},
						"backend_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The backend type of the API.`,
						},
						"cors": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether CORS is supported.",
						},
						"matching": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The matching mode of the API.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the API.",
						},
						"tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of tags configuration.",
						},
						// The format is `yyyy-MM-ddTHH:mm:ss{timezone}`, e.g. `2006-01-02 15:04:05+08:00`.
						"registered_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The registered time of the API, in RFC3339 format.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the API, in RFC3339 format.",
						},
						"published_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The published time of the API, in RFC3339 format.",
						},
					},
				},
			},
		},
	}
}

func buildListApisParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("api_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		res = fmt.Sprintf("%s&group_id=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, buildApiType(v.(string)))
	}
	if v, ok := d.GetOk("request_method"); ok {
		res = fmt.Sprintf("%s&req_method=%v", res, v)
	}
	if v, ok := d.GetOk("request_path"); ok {
		res = fmt.Sprintf("%s&req_uri=%v", res, v)
	}
	if v, ok := d.GetOk("request_protocol"); ok {
		res = fmt.Sprintf("%s&req_protocol=%v", res, v)
	}
	if v, ok := d.GetOk("security_authentication"); ok {
		res = fmt.Sprintf("%s&auth_type=%v", res, v)
	}
	if v, ok := d.GetOk("vpc_channel_name"); ok {
		res = fmt.Sprintf("%s&vpc_channel_name=%v", res, v)
	}
	if v, ok := d.GetOk("precise_search"); ok {
		res = fmt.Sprintf("%s&precise_search=%v", res, v)
	}
	if v, ok := d.GetOk("env_id"); ok {
		res = fmt.Sprintf("%s&env_id=%v", res, v)
	}
	return res
}

func GetApiBasicConfigurations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apis?limit=500"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	queryParams := buildListApisParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving apis under specified dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		apiBasicConfigurations := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		if len(apiBasicConfigurations) < 1 {
			break
		}
		result = append(result, apiBasicConfigurations...)
		offset += len(apiBasicConfigurations)
	}
	return result, nil
}

func dataSourceApiBasicConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	apiBasicConfigurations, err := GetApiBasicConfigurations(client, d)
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
		d.Set("configurations", filterApiBasicConfigurations(flattenApiBasicConfigurations(apiBasicConfigurations), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterApiBasicConfigurations(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("env_name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("env_name", v, nil)) {
			continue
		}

		if param, ok := d.GetOk("backend_type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("backend_type", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenApiBasicConfigurations(configurations []interface{}) []interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(configurations))
	for _, conf := range configurations {
		result = append(result, map[string]interface{}{
			"id":                      utils.PathSearch("id", conf, nil),
			"name":                    utils.PathSearch("name", conf, nil),
			"type":                    analyseApiType(int(utils.PathSearch("type", conf, float64(0)).(float64))),
			"request_method":          utils.PathSearch("req_method", conf, nil),
			"request_path":            utils.PathSearch("req_uri", conf, nil),
			"request_protocol":        utils.PathSearch("req_protocol", conf, nil),
			"security_authentication": utils.PathSearch("auth_type", conf, nil),
			"simple_authentication":   flattenSimpleAuth(utils.PathSearch("auth_opt.app_code_auth_type", conf, "").(string)),
			"authorizer_id":           utils.PathSearch("authorizer_id", conf, nil),
			"group_id":                utils.PathSearch("group_id", conf, nil),
			"group_name":              utils.PathSearch("group_name", conf, nil),
			"group_version":           utils.PathSearch("group_version", conf, nil),
			"env_id":                  utils.PathSearch("run_env_id", conf, nil),
			"env_name":                utils.PathSearch("run_env_name", conf, nil),
			"publish_id":              utils.PathSearch("publish_id", conf, nil),
			"backend_type":            utils.PathSearch("backend_type", conf, nil),
			"cors":                    utils.PathSearch("cors", conf, nil),
			"matching":                analyseApiMatchMode(utils.PathSearch("match_mode", conf, "").(string)),
			"description":             utils.PathSearch("remark", conf, nil),
			"tags":                    utils.PathSearch("tags", conf, nil),
			"registered_at":           flattenTimeToRFC3339(utils.PathSearch("register_time", conf, "").(string)),
			"updated_at":              flattenTimeToRFC3339(utils.PathSearch("update_time", conf, "").(string)),
			"published_at":            flattenPulishTime(utils.PathSearch("publish_time", conf, "").(string)),
		})
	}
	return result
}

func flattenSimpleAuth(authType string) bool {
	return authType == string(AppCodeAuthTypeEnable)
}

// Formats the time according to the local computer's time.
func flattenTimeToRFC3339(timeStr string) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(timeStr)/1000, false)
}

func flattenPulishTime(utcTime string) string {
	pulishTime := utils.ConvertTimeStrToNanoTimestamp(utcTime, "2006-01-02 15:04:05")
	return utils.FormatTimeStampRFC3339(pulishTime/1000, false)
}
