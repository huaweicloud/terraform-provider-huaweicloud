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

			// Optional parameters.
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version used to filter the microservice engine flavors.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID to which the microservice engine flavors belong.`,
			},

			// Attributes.
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
				Description: `The list of microservice engine flavors that matched filter parameters.`,
			},
		},
	}
}

func buildMicroserviceEngineFlavorsQueryParams(d *schema.ResourceData) string {
	result := ""

	if v, ok := d.GetOk("version"); ok {
		result = fmt.Sprintf("%s&specType=%s", result, v)
	}

	if result != "" {
		return "?" + result[1:]
	}
	return result
}

func listMicroserviceEngineFlavors(client *golangsdk.ServiceClient, d *schema.ResourceData,
	enterpriseProjectId string) ([]interface{}, error) {
	var (
		httpUrl  = "v2/{project_id}/enginemgr/flavors"
		listPath = client.Endpoint + httpUrl
	)

	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildMicroserviceEngineFlavorsQueryParams(d)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if enterpriseProjectId != "" {
		listOpts.MoreHeaders = buildRequestMoreHeaders(enterpriseProjectId)
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
	if len(flavors) < 1 {
		return nil
	}

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
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		enterpriseProjectId = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	flavors, err := listMicroserviceEngineFlavors(client, d, enterpriseProjectId)
	if err != nil {
		return diag.Errorf("error querying microservice engine flavors: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("enterprise_project_id", enterpriseProjectId),
		d.Set("flavors", flattenMicroserviceEngineFlavors(flavors)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
