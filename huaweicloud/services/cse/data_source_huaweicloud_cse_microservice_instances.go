package cse

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSE POST /v4/token
// @API CSE GET /v4/{project_id}/registry/microservices/{service_id}/instances
func DataSourceMicroserviceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMicroserviceInstancesRead,

		Schema: map[string]*schema.Schema{
			// Special parameters.
			// These parameters are used to specify the address that used to request the access token and access the
			// microservice engine.
			"auth_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The address that used to request the access token.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The address that used to access engine and query microservice instances.`,
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

			// Required parameters.
			"microservice_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the microservice to which the microservice instances belong.`,
			},

			// Attributes.
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the microservice instance.`,
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The host name of the microservice instance.`,
						},
						"endpoints": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of the access addresses of the microservice instance.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the microservice instance.`,
						},
						"properties": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The extended attributes of the microservice instance, in key/value format.`,
						},
						"health_check": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The heartbeat mode of the health check.`,
									},
									"interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The heartbeat interval of the health check, in seconds.`,
									},
									"max_retries": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The maximum retry number of the health check.`,
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The port of the health check.`,
									},
								},
							},
							Description: `The health check configuration of the microservice instance.`,
						},
						"data_center": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the data center.`,
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The custom region name of the data center.`,
									},
									"availability_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The custom availability zone of the data center.`,
									},
								},
							},
							Description: `The data center configuration of the microservice instance.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the microservice instance.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the microservice instance, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the microservice instance, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of the microservice instances.`,
			},
		},
	}
}

func listMicroserviceInstances(client *golangsdk.ServiceClient, authAddress, adminUser, adminPass,
	microserviceId string) ([]interface{}, error) {
	httpUrl := "v4/{project_id}/registry/microservices/{service_id}/instances"

	listPath := client.Endpoint + httpUrl
	// The project ID of the microservice instance is the fixed value "default".
	// No region parameter needs to be defined because this data source does not use IAM authentication.
	listPath = strings.ReplaceAll(listPath, "{project_id}", microserviceInstanceProjectId)
	listPath = strings.ReplaceAll(listPath, "{service_id}", microserviceId)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
	}

	// When a user configures both the `admin_user` and `admin_pass` fields, it indicates that the microservice engine
	// has enabled RBAC authentication. Subsequent requests will require the token information obtained via the
	// `POST /v4/token` interface.
	token, err := GetAuthorizationToken(authAddress, adminUser, adminPass)
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
	return utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{}), nil
}

// parseMicroserviceInstanceTimestamp converts the timestamp from API response (string, float64, or int) to RFC3339 format.
func parseMicroserviceInstanceTimestamp(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		if v == "" {
			return ""
		}
		r, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("[ERROR] unable to convert the string (%s) to int: %v", v, err)
			return ""
		}
		return utils.FormatTimeStampRFC3339(int64(r), false)
	case float64:
		return utils.FormatTimeStampRFC3339(int64(v), false)
	case int:
		return utils.FormatTimeStampRFC3339(int64(v), false)
	default:
		return ""
	}
}

func flattenMicroserviceInstances(instances []interface{}) []map[string]interface{} {
	if len(instances) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(instances))
	for _, instance := range instances {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("instanceId", instance, nil),
			"host_name":    utils.PathSearch("hostName", instance, nil),
			"endpoints":    utils.PathSearch("endpoints", instance, nil),
			"version":      utils.PathSearch("version", instance, nil),
			"properties":   buildInstanceProperties(utils.PathSearch("properties", instance, make(map[string]interface{})).(map[string]interface{})),
			"health_check": flattenHealthCheck(utils.PathSearch("healthCheck", instance, nil)),
			"data_center":  flattenDataCenter(utils.PathSearch("dataCenterInfo", instance, nil)),
			"status":       utils.PathSearch("status", instance, nil),
			"created_at":   parseMicroserviceInstanceTimestamp(utils.PathSearch("timestamp", instance, nil)),
			"updated_at":   parseMicroserviceInstanceTimestamp(utils.PathSearch("modTimestamp", instance, nil)),
		})
	}
	return result
}

func dataSourceMicroserviceInstancesRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var (
		// Querying microservice instances in the microservice engine requires building a client based on the
		// microservice engine's connection address, which does not use IAM authentication.
		client         = common.NewCustomClient(true, d.Get("connect_address").(string))
		authAddress    = getAuthAddress(d)
		adminUser      = d.Get("admin_user").(string)
		adminPass      = d.Get("admin_pass").(string)
		microserviceId = d.Get("microservice_id").(string)
	)

	instances, err := listMicroserviceInstances(client, authAddress, adminUser, adminPass, microserviceId)
	if err != nil {
		return diag.Errorf("error querying microservice instances: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("instances", flattenMicroserviceInstances(instances)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
