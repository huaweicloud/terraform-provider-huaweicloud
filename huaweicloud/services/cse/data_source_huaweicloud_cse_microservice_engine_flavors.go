package cse

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

// @API CSE GET /v2/{project_id}/enginemgr/flavors
func DataSourceMicroserviceEngineFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMicroserviceEngineFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the microservice engine flavors are located.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version that used to filter the microservice engine flavors.`,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the microservice engine flavor.`,
						},
						"spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"available_cpu_memory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The CPU and memory combinations that the flavor is allowed.`,
									},
									"linear": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Whether the microservice engine flavor is a linear flavor.`,
									},
									"available_prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The flavor name prefix of the available node.`,
									},
								},
							},
							Description: `The specification detail of the microservice engine flavor.`,
						},
					},
				},
				Description: `All microservice engine flavors that match the filter parameters.`,
			},
		},
	}
}

func queryMicroserviceEngineFlavors(client *golangsdk.ServiceClient, specType string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/enginemgr/flavors"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	if specType != "" {
		listPath = fmt.Sprintf("%s?specType=%s", listPath, specType)
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenMicroserviceEngineFlavorSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"available_cpu_memory": utils.PathSearch("availableCpuMemory", spec, nil),
			"linear":               utils.PathSearch("linear", spec, nil),
			"available_prefix":     utils.PathSearch("availablePrefix", spec, nil),
		},
	}
}

func flattenMicroserviceEngineFlavors(flavors []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(flavors))

	for _, flavor := range flavors {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("flavor", flavor, nil),
			"spec": flattenMicroserviceEngineFlavorSpec(utils.PathSearch("spec", flavor, nil)),
		})
	}

	return result
}

func dataSourceMicroserviceEngineFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	flavors, err := queryMicroserviceEngineFlavors(client, d.Get("version").(string))
	if err != nil {
		return diag.Errorf("error querying microservice engine flavors: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavors", flattenMicroserviceEngineFlavors(flavors)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
