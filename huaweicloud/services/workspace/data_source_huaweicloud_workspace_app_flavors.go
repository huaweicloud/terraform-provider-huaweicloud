package workspace

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

// @API Workspace GET /v1/{project_id}/product
func DataSourceAppFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The product ID used to filter the app flavor list.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The flavor ID used to filter the app flavor list.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The availability zone used to filter the app flavor list.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The operating system type used to filter the app flavor list.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The charge mode used to filter the app flavor list.`,
			},
			"architecture": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The architecture type used to filter the app flavor list.`,
			},

			// attributes
			"flavors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of app flavors that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product ID of app flavors that matched filter parameters.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flavor ID of app flavors that matched filter parameters.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flavor type.`,
						},
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flavor architecture of app flavors that matched filter parameters.`,
						},
						"cpu": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CPU core count.`,
						},
						"memory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The memory size in MB.`,
						},
						"is_gpu": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the flavor is GPU type.`,
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The system disk type.`,
						},
						"system_disk_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The system disk size.`,
						},
						"descriptions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flavor description.`,
						},
						"charge_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The charge mode of app flavors that matched filter parameters.`,
						},
						"contain_data_disk": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the flavor includes data disk.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type.`,
						},
						"cloud_service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cloud service type.`,
						},
						"volume_product_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The volume product type.`,
						},
						"sessions": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of sessions supported by the flavor.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flavor status.`,
						},
						"cond_operation_az": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flavor status in availability zones.`,
						},
						"domain_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The domain IDs that the flavor belongs to.`,
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The package type.`,
						},
					},
				},
			},
		},
	}
}

func buildAppFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("product_id"); ok {
		res = fmt.Sprintf("%s&product_id=%v", res, v)
	}
	if v, ok := d.GetOk("flavor_id"); ok {
		res = fmt.Sprintf("%s&flavor_id=%v", res, v)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		res = fmt.Sprintf("%s&availability_zone=%v", res, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		res = fmt.Sprintf("%s&os_type=%v", res, v)
	}
	if v, ok := d.GetOk("charge_mode"); ok {
		res = fmt.Sprintf("%s&charge_mode=%v", res, v)
	}
	if v, ok := d.GetOk("architecture"); ok {
		res = fmt.Sprintf("%s&architecture=%v", res, v)
	}

	return res
}

func getAppFlavors(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/product"
	queryParams := buildAppFlavorsQueryParams(d)
	if queryParams != "" {
		httpUrl = httpUrl + "?" + queryParams[1:]
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("products", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceAppFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace app client: %s", err)
	}

	flavors, err := getAppFlavors(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace app flavors: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavors", flattenAppFlavors(flavors)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppFlavors(flavors []interface{}) []map[string]interface{} {
	if len(flavors) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(flavors))
	for _, flavor := range flavors {
		result = append(result, map[string]interface{}{
			"product_id":          utils.PathSearch("product_id", flavor, nil),
			"id":                  utils.PathSearch("flavor_id", flavor, nil),
			"type":                utils.PathSearch("type", flavor, nil),
			"architecture":        utils.PathSearch("architecture", flavor, nil),
			"cpu":                 utils.PathSearch("cpu", flavor, nil),
			"memory":              utils.PathSearch("memory", flavor, nil),
			"is_gpu":              utils.PathSearch("is_gpu", flavor, false),
			"system_disk_type":    utils.PathSearch("system_disk_type", flavor, nil),
			"system_disk_size":    utils.PathSearch("system_disk_size", flavor, nil),
			"descriptions":        utils.PathSearch("descriptions", flavor, nil),
			"charge_mode":         utils.PathSearch("charge_mode", flavor, nil),
			"contain_data_disk":   utils.PathSearch("contain_data_disk", flavor, false),
			"resource_type":       utils.PathSearch("resource_type", flavor, nil),
			"cloud_service_type":  utils.PathSearch("cloud_service_type", flavor, nil),
			"volume_product_type": utils.PathSearch("volume_product_type", flavor, nil),
			"sessions":            utils.PathSearch("sessions", flavor, 0),
			"status":              utils.PathSearch("status", flavor, nil),
			"cond_operation_az":   utils.PathSearch("cond_operation_az", flavor, nil),
			"domain_ids":          utils.PathSearch("domain_ids", flavor, nil),
			"package_type":        utils.PathSearch("package_type", flavor, nil),
		})
	}

	return result
}
