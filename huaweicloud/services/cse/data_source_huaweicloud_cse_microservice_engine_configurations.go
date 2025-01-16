package cse

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSE GET /v1/{project_id}/kie/kv
func DataSourceMicroserviceEngineConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMicroserviceEngineConfigurationsRead,

		Schema: map[string]*schema.Schema{
			// Authentication and request parameters.
			"auth_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The address that used to request the access token.`,
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The address that used to send requests and manage configuration.`,
			},
			"admin_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user name that used to pass the RBAC control.`,
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"admin_user"},
				Description:  `The user password that used to pass the RBAC control.`,
			},
			// Attributes.
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the microservice engine configuration.`,
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration key (item name).",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the configuration value.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The configuration value.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The configuration status.`,
						},
						"tags": common.TagsComputedSchema(
							`The key/value pairs to associate with the configuration.`,
						),
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the configuration, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the configuration, in RFC3339 format.`,
						},
						"create_revision": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The create version of the configuration.`,
						},
						"update_revision": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The update version of the configuration.`,
						},
					},
				},
				Description: `All queried configurations of the dedicated microservice engine.`,
			},
		},
	}
}

func queryMicroserviceEngineConfigurations(d *schema.ResourceData) ([]interface{}, error) {
	var (
		client    = common.NewCustomClient(true, d.Get("connect_address").(string), "v1", "default")
		httpUrl   = "kie/kv"
		queryPath = client.ResourceBase + httpUrl
	)

	token, err := GetAuthorizationToken(d.Get("auth_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return nil, err
	}

	queryOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if token != "" {
		queryOpts.MoreHeaders = map[string]string{
			"Authorization": token,
		}
	}

	requestResp, err := client.Request("GET", queryPath, &queryOpts)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenMicroserviceEngineConfigurations(configurations []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(configurations))

	for _, configuration := range configurations {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", configuration, nil),
			"key":             utils.PathSearch("key", configuration, nil),
			"value_type":      utils.PathSearch("value_type", configuration, nil),
			"value":           utils.PathSearch("value", configuration, nil),
			"status":          utils.PathSearch("status", configuration, nil),
			"tags":            utils.PathSearch("labels", configuration, nil),
			"created_at":      utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", configuration, float64(0)).(float64)), false),
			"updated_at":      utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", configuration, float64(0)).(float64)), false),
			"create_revision": utils.PathSearch("create_revision", configuration, nil),
			"update_revision": utils.PathSearch("update_revision", configuration, nil),
		})
	}

	return result
}

func dataSourceMicroserviceEngineConfigurationsRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	configurations, err := queryMicroserviceEngineConfigurations(d)
	if err != nil {
		return diag.Errorf("error querying microservice engine configurations: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("configurations", flattenMicroserviceEngineConfigurations(configurations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
