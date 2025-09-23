package cse

import (
	"context"
	"fmt"
	"log"
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

// @API CSE GET /v2/{project_id}/enginemgr/engines
func DataSourceMicroserviceEngines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMicroserviceEnginesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where dedicated microservice engines are located.`,
			},
			"engines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dedicated microservice engine.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the dedicated microservice engine.`,
						},
						"flavor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flavor name of the dedicated microservice engine.`,
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of availability zones.`,
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The network ID of the subnet to which the dedicated microservice engine belongs.`,
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The authentication method for the dedicated microservice engine.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the dedicated microservice engine.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID to which the dedicated microservice engine belongs.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the dedicated microservice engine.`,
						},
						"eip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The EIP ID to which the dedicated microservice engine assocated.`,
						},
						"extend_params": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The additional parameters for the dedicated microservice engine.`,
						},
						"service_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of the microservice resources.`,
						},
						"instance_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of the microservice instance resources.`,
						},
						"service_registry_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The internal access address.`,
									},
									"public": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The public access address.`,
									},
								},
							},
							Description: `The connection addresses of service center.`,
						},
						"config_center_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The internal access address.`,
									},
									"public": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The public access address.`,
									},
								},
							},
							Description: `The addresses of config center.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the dedicated microservice engine.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the dedicated microservice engine, in RFC3339 format.`,
						},
					},
				},
				Description: `All queried dedicated microservice engines.`,
			},
		},
	}
}

func queryMicroserviceEngines(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/enginemgr/engines?type=CSE&limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error querying microservice engines: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		engines := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		// When the number of remaining engines is 1, both offset 0 and 1 will return this engine object.
		// So it needs to be included in the total determination.
		if len(engines) < 1 || offset >= int(utils.PathSearch("total", respBody, float64(0)).(float64)) {
			break
		}
		result = append(result, engines...)
		offset += len(engines)
	}
	return result, nil
}

func flattenMicroserviceEngineServiceRegistryAddresses(engineDetail interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"private": utils.PathSearch("serviceEndpoint.serviceCenter.masterEntrypoint", engineDetail, nil),
			"public":  utils.PathSearch("publicServiceEndpoint.serviceCenter.masterEntrypoint", engineDetail, nil),
		},
	}
}

func flattenMicroserviceEngineConfigCenterAddresses(engineDetail interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"private": utils.PathSearch("serviceEndpoint.kie.masterEntrypoint", engineDetail, nil),
			"public":  utils.PathSearch("publicServiceEndpoint.kie.masterEntrypoint", engineDetail, nil),
		},
	}
}

func convertStrNumberToIntIgnoreErr(strNum string) int {
	result, err := strconv.Atoi(strNum)
	if err != nil {
		log.Printf("[ERROR] unable to convert object from type string to type int: %s", err)
		return 0
	}
	return result
}

func flattenAllMicroserviceEngines(engines []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(engines))

	for _, engineDetail := range engines {
		result = append(result, map[string]interface{}{
			"id":                         utils.PathSearch("id", engineDetail, nil),
			"name":                       utils.PathSearch("name", engineDetail, nil),
			"flavor":                     utils.PathSearch("flavor", engineDetail, nil),
			"availability_zones":         utils.PathSearch("reference.azList", engineDetail, nil),
			"network_id":                 utils.PathSearch("reference.networkId", engineDetail, nil),
			"auth_type":                  utils.PathSearch("authType", engineDetail, nil),
			"version":                    utils.PathSearch("version", engineDetail, nil),
			"enterprise_project_id":      utils.PathSearch("enterpriseProjectId", engineDetail, nil),
			"description":                utils.PathSearch("description", engineDetail, nil),
			"eip_id":                     utils.PathSearch("reference.publicIpId", engineDetail, nil),
			"extend_params":              utils.PathSearch("reference.inputs", engineDetail, nil),
			"service_limit":              convertStrNumberToIntIgnoreErr(utils.PathSearch("reference.serviceLimit", engineDetail, "").(string)),
			"instance_limit":             convertStrNumberToIntIgnoreErr(utils.PathSearch("reference.instanceLimit", engineDetail, "").(string)),
			"service_registry_addresses": flattenMicroserviceEngineServiceRegistryAddresses(engineDetail),
			"config_center_addresses":    flattenMicroserviceEngineConfigCenterAddresses(engineDetail),
			"status":                     utils.PathSearch("status", engineDetail, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("createTime",
				engineDetail, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceMicroserviceEnginesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	engines, err := queryMicroserviceEngines(client)
	if err != nil {
		return diag.Errorf("error querying CSE microservice engines: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("engines", flattenAllMicroserviceEngines(engines)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
