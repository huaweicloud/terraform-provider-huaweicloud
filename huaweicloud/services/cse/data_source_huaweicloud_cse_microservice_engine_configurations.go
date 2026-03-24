package cse

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSE POST /v4/token
// @API CSE GET /v1/{project_id}/kie/kv
func DataSourceMicroserviceEngineConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMicroserviceEngineConfigurationsRead,

		Schema: map[string]*schema.Schema{
			// Special parameters.
			// These parameters are used to specify the address that used to request the access token and access the
			// microservice engine.
			"auth_address": {
				Type:     schema.TypeString,
				Optional: true,
				Description: utils.SchemaDesc(
					`The address that used to request the access token.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The address that used to access the engine and query configurations.`,
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

			// Optional parameters.
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID to which the microservice engine configurations belong.`,
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
							Description: `The configuration key (item name).`,
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
				Description: `The list of configurations that matched filter parameters.`,
			},
		},
	}
}

func listMicroserviceEngineConfigurations(client *golangsdk.ServiceClient, authInfo MicroserviceEngineAuthInfo) ([]interface{}, error) {
	var (
		httpUrl  = "v1/{project_id}/kie/kv"
		listPath = client.Endpoint + httpUrl
	)

	// The project ID of the microservice engine configurations is the fixed value "default".
	// No region parameter needs to be defined because this data source does not use IAM authentication.
	listPath = strings.ReplaceAll(listPath, "{project_id}", microserviceEngineConfigurationDefaultProjectId)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(authInfo.EnterpriseProjectId),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(authInfo.AuthAddress, authInfo.AdminUser, authInfo.AdminPass)
	if err != nil {
		return nil, err
	}
	// If the microservice instance has RBAC authentication enabled, the Authorization header will use a special token
	// provided by the CSE service to replace the original IAM authentication information (AKSK authentication) in the
	// request header.
	if token != "" {
		listOpts.MoreHeaders["Authorization"] = token
	}

	requestResp, err := client.Request("GET", listPath, &listOpts)
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
	if len(configurations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(configurations))
	for _, configuration := range configurations {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", configuration, nil),
			"key":             utils.PathSearch("key", configuration, nil),
			"value_type":      utils.PathSearch("value_type", configuration, nil),
			"value":           utils.PathSearch("value", configuration, nil),
			"status":          utils.PathSearch("status", configuration, nil),
			"tags":            utils.PathSearch("labels", configuration, nil),
			"created_at":      formatKieKvTime(utils.PathSearch("create_time", configuration, nil)),
			"updated_at":      formatKieKvTime(utils.PathSearch("update_time", configuration, nil)),
			"create_revision": utils.PathSearch("create_revision", configuration, nil),
			"update_revision": utils.PathSearch("update_revision", configuration, nil),
		})
	}

	return result
}

func dataSourceMicroserviceEngineConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
		// Querying microservice engine configurations requires building a client based on the
		// microservice engine's connection address, which does not use IAM authentication.
		client                     = common.NewCustomClient(true, d.Get("connect_address").(string))
		microserviceEngineAuthInfo = MicroserviceEngineAuthInfo{
			AuthAddress:         getAuthAddress(d),
			AdminUser:           d.Get("admin_user").(string),
			AdminPass:           d.Get("admin_pass").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}
	)

	configurations, err := listMicroserviceEngineConfigurations(client, microserviceEngineAuthInfo)
	if err != nil {
		return diag.Errorf("error querying microservice engine configurations: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("enterprise_project_id", microserviceEngineAuthInfo.EnterpriseProjectId),
		d.Set("configurations", flattenMicroserviceEngineConfigurations(configurations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
