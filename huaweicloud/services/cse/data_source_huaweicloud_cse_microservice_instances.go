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

// @API CSE GET /v4/{project_id}/registry/microservices/{service_id}/instances
func DataSourceMicroserviceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMicroserviceInstancesRead,

		Schema: map[string]*schema.Schema{
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
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"admin_pass"},
				Description:  `The user name that used to pass the RBAC control.`,
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"admin_user"},
				Description:  `The user password that used to pass the RBAC control.`,
			},
			"microservice_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated microservice to which the instances belong.`,
			},
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
							Description: `The host name ID of the microservice instance.`,
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
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The custom region name of the data center.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the data center.`,
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

func queryMicroserviceInstances(client *golangsdk.ServiceClient, token, microserviceId string) ([]interface{}, error) {
	var (
		httpUrl = "registry/microservices/{service_id}/instances"
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)

	if token != "" {
		listOpt.MoreHeaders = map[string]string{
			"Authorization": token,
		}
	}

	listPath := client.ResourceBase + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{service_id}", microserviceId)

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenMicroserviceInstances(instances []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(instances))

	for _, instance := range instances {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("instanceId", instance, nil),
			"host_name":    utils.PathSearch("hostName", instance, nil),
			"endpoints":    utils.PathSearch("endpoints", instance, nil),
			"version":      utils.PathSearch("version", instance, nil),
			"properties":   utils.PathSearch("properties", instance, nil),
			"health_check": flattenMicroserviceInstancesHealthCheck(utils.PathSearch("healthCheck", instance, nil)),
			"data_center":  flattenMicroserviceInstancesDataCenter(utils.PathSearch("dataCenterInfo", instance, nil)),
			"status":       utils.PathSearch("status", instance, nil),
			"created_at":   parseMicroserviceInstancesTime(utils.PathSearch("timestamp", instance, "").(string)),
			"updated_at":   parseMicroserviceInstancesTime(utils.PathSearch("modTimestamp", instance, "").(string)),
		})
	}

	return result
}

func flattenMicroserviceInstancesHealthCheck(healthCheck interface{}) []map[string]interface{} {
	if healthCheck == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mode":        utils.PathSearch("mode", healthCheck, nil),
			"interval":    utils.PathSearch("interval", healthCheck, nil),
			"max_retries": utils.PathSearch("times", healthCheck, nil),
			"port":        utils.PathSearch("port", healthCheck, nil),
		},
	}
}

func flattenMicroserviceInstancesDataCenter(dataCenter interface{}) []map[string]interface{} {
	if dataCenter == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"region":            utils.PathSearch("region", dataCenter, nil),
			"name":              utils.PathSearch("name", dataCenter, nil),
			"availability_zone": utils.PathSearch("availableZone", dataCenter, nil),
		},
	}
}

func parseMicroserviceInstancesTime(timeStampStr string) string {
	r, err := strconv.Atoi(timeStampStr)
	if err != nil {
		log.Printf("[ERROR] unable to convert the string (%s) to int", timeStampStr)
		return ""
	}

	return utils.FormatTimeStampRFC3339(int64(r), false)
}

func dataSourceMicroserviceInstancesRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	token, err := GetAuthorizationToken(d.Get("auth_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v4", "default")
	instances, err := queryMicroserviceInstances(client, token, d.Get("microservice_id").(string))
	if err != nil {
		return diag.Errorf("error querying CSE microservice instances: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("instances", flattenMicroserviceInstances(instances)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
