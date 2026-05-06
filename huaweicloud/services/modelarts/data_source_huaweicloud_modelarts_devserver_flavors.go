package modelarts

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

// @API ModelArts GET /v1/{project_id}/dev-servers/resource-flavors
func DataSourceDevServerFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevServerFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the DevServer flavors are located.`,
			},
			"server_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The service type of the DevServer flavors.`,
			},
			"arch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The CPU architecture of the DevServer flavors.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The charging mode of the DevServer flavors.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of DevServer flavors that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IAAS flavor name.`,
						},
						"specification": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The specification configuration of the flavor.`,
						},
						"arch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CPU architecture of the flavor.`,
						},
						"server_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The server type of the flavor.`,
						},
						"sku_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SKU billing code of the flavor.`,
						},
						"charging_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The charging mode of the flavor.`,
						},
						"roce_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of NICs of the flavor.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of super node instances.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the server flavor.`,
						},
						"server_flavors": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The expandable flavors of the super node.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the server flavor.`,
									},
								},
							},
						},
						"flavor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The computing power card type of the flavor.`,
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The availability zones of the flavor.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the availability zone.`,
									},
									"is_sold_out": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the flavor is sold out in the availability zone.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildDevServerFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("server_type"); ok {
		res = fmt.Sprintf("%s&server_type=%v", res, v)
	}
	if v, ok := d.GetOk("arch"); ok {
		res = fmt.Sprintf("%s&arch=%v", res, v)
	}
	if v, ok := d.GetOk("charging_mode"); ok {
		res = fmt.Sprintf("%s&charging_mode=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func listDevServerFlavors(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/dev-servers/resource-flavors"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildDevServerFlavorsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
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

func flattenDevServerFlavorAvailabilityZones(availabilityZones []interface{}) []map[string]interface{} {
	if len(availabilityZones) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(availabilityZones))
	for _, az := range availabilityZones {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", az, nil),
			"is_sold_out": utils.PathSearch("is_sold_out", az, nil),
		})
	}

	return result
}

func flattenDevServerFlavorServerFlavors(serverFlavors []interface{}) []map[string]interface{} {
	if len(serverFlavors) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(serverFlavors))
	for _, serverFlavor := range serverFlavors {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", serverFlavor, nil),
		})
	}

	return result
}

func flattenDevServerFlavors(flavors []interface{}) []map[string]interface{} {
	if len(flavors) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(flavors))
	for _, flavor := range flavors {
		result = append(result, map[string]interface{}{
			"flavor":        utils.PathSearch("flavor", flavor, nil),
			"specification": utils.PathSearch("specification", flavor, nil),
			"arch":          utils.PathSearch("arch", flavor, nil),
			"server_type":   utils.PathSearch("server_type", flavor, nil),
			"sku_code":      utils.PathSearch("sku_code", flavor, nil),
			"charging_mode": utils.PathSearch("charging_mode", flavor, nil),
			"roce_num":      utils.PathSearch("roce_num", flavor, nil),
			"count":         utils.PathSearch("count", flavor, nil),
			"status":        utils.PathSearch("status", flavor, nil),
			"flavor_type":   utils.PathSearch("flavor_type", flavor, nil),
			"server_flavors": flattenDevServerFlavorServerFlavors(
				utils.PathSearch("server_flavors", flavor, make([]interface{}, 0)).([]interface{})),
			"availability_zones": flattenDevServerFlavorAvailabilityZones(
				utils.PathSearch("availability_zones", flavor, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceDevServerFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	flavors, err := listDevServerFlavors(client, d)
	if err != nil {
		return diag.Errorf("error querying DevServer flavors: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavors", flattenDevServerFlavors(flavors)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
