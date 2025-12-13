package workspace

import (
	"context"
	"fmt"
	"strconv"
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

// @API Workspace GET /v2/{project_id}/desktop-pools/{pool_id}/desktops
func DataSourceDesktopPoolAssociatedDesktops() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopPoolAssociatedDesktopsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the associated desktops are located.`,
			},

			// Required parameters.
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the desktop pool to which the associated desktops belong.`,
			},

			// Attributes.
			"desktops": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of associated desktops.`,
				Elem:        desktopPoolAssociatedDesktopSchema(),
			},
		},
	}
}

func desktopPoolAssociatedDesktopSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"desktop_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the desktop.`,
			},
			"computer_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the desktop.`,
			},
			"os_host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OS host name of the desktop.`,
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of IP addresses of the desktop.`,
			},
			"ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IPv4 address of the desktop.`,
			},
			"ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IPv6 address of the desktop.`,
			},
			"desktop_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the desktop.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the desktop.`,
			},
			"in_maintenance_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the desktop is in maintenance mode.`,
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the desktop.`,
			},
			"login_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The login status of the desktop.`,
			},
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product ID of the desktop.`,
			},
			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolAssociatedDesktopVolumeSchema(),
				Description: `The root volume information of the desktop.`,
			},
			"data_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolAssociatedDesktopVolumeSchema(),
				Description: `The list of data volumes of the desktop.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone of the desktop.`,
			},
			"site_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The site type of the desktop.`,
			},
			"site_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The site name of the desktop.`,
			},
			"product": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolAssociatedDesktopProductSchema(),
				Description: `The product information of the desktop.`,
			},
			"os_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OS version of the desktop.`,
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SID of the desktop.`,
			},
			"tags": common.TagsComputedSchema(`The tags of the desktop.`),
			"is_support_internet": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the desktop supports internet access.`,
			},
			"is_attaching_eip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the desktop is attaching an EIP.`,
			},
			"attach_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attach state of the desktop.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID of the desktop.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet ID of the desktop.`,
			},
			"bill_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The billing resource ID of the desktop.`,
			},
		},
	}
}

func desktopPoolAssociatedDesktopVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the volume.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The size of the volume in GB.`,
			},
			"device": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The device name of the volume.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the volume.`,
			},
			"volume_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The volume ID.`,
			},
			"bill_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The billing resource ID of the volume.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the volume.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The display name of the volume.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource specification code of the volume.`,
			},
		},
	}
}

func desktopPoolAssociatedDesktopProductSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product ID.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor ID.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product type.`,
			},
			"cpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CPU specification.`,
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The memory specification.`,
			},
			"descriptions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product description.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The charging mode.`,
			},
			"architecture": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The architecture of the product.`,
			},
			"is_gpu": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the product is GPU type.`,
			},
			"package_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The package type of the product.`,
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
			"contain_data_disk": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the product contains data disk.`,
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
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the product.`,
			},
		},
	}
}

func listDesktopPoolAssociatedDesktops(client *golangsdk.ServiceClient, poolId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktop-pools/{pool_id}/desktops?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{pool_id}", poolId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		desktops := utils.PathSearch("pool_desktops", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, desktops...)
		if len(desktops) < limit {
			break
		}
		offset += len(desktops)
	}

	return result, nil
}

func flattenDesktopPoolAssociatedDesktopRootVolume(volume interface{}) []map[string]interface{} {
	if volume == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":               utils.PathSearch("type", volume, nil),
			"size":               utils.PathSearch("size", volume, nil),
			"device":             utils.PathSearch("device", volume, nil),
			"id":                 utils.PathSearch("id", volume, nil),
			"volume_id":          utils.PathSearch("volume_id", volume, nil),
			"bill_resource_id":   utils.PathSearch("bill_resource_id", volume, nil),
			"create_time":        utils.PathSearch("create_time", volume, nil),
			"display_name":       utils.PathSearch("display_name", volume, nil),
			"resource_spec_code": utils.PathSearch("resource_spec_code", volume, nil),
		},
	}
}

func flattenDesktopPoolAssociatedDesktopVolumes(volumes []interface{}) []map[string]interface{} {
	if len(volumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(volumes))
	for _, volume := range volumes {
		result = append(result, map[string]interface{}{
			"type":               utils.PathSearch("type", volume, nil),
			"size":               utils.PathSearch("size", volume, nil),
			"device":             utils.PathSearch("device", volume, nil),
			"id":                 utils.PathSearch("id", volume, nil),
			"volume_id":          utils.PathSearch("volume_id", volume, nil),
			"bill_resource_id":   utils.PathSearch("bill_resource_id", volume, nil),
			"create_time":        utils.PathSearch("create_time", volume, nil),
			"display_name":       utils.PathSearch("display_name", volume, nil),
			"resource_spec_code": utils.PathSearch("resource_spec_code", volume, nil),
		})
	}

	return result
}

func flattenDesktopPoolAssociatedDesktopProduct(product interface{}) []map[string]interface{} {
	if product == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"product_id":          utils.PathSearch("product_id", product, nil),
			"flavor_id":           utils.PathSearch("flavor_id", product, nil),
			"type":                utils.PathSearch("type", product, nil),
			"cpu":                 utils.PathSearch("cpu", product, nil),
			"memory":              utils.PathSearch("memory", product, nil),
			"descriptions":        utils.PathSearch("descriptions", product, nil),
			"charge_mode":         utils.PathSearch("charge_mode", product, nil),
			"architecture":        utils.PathSearch("architecture", product, nil),
			"is_gpu":              utils.PathSearch("is_gpu", product, nil),
			"package_type":        utils.PathSearch("package_type", product, nil),
			"system_disk_type":    utils.PathSearch("system_disk_type", product, nil),
			"system_disk_size":    utils.PathSearch("system_disk_size", product, nil),
			"contain_data_disk":   utils.PathSearch("contain_data_disk", product, nil),
			"resource_type":       utils.PathSearch("resource_type", product, nil),
			"cloud_service_type":  utils.PathSearch("cloud_service_type", product, nil),
			"volume_product_type": utils.PathSearch("volume_product_type", product, nil),
			"status":              utils.PathSearch("status", product, nil),
		},
	}
}

func flattenDesktopPoolAssociatedDesktops(desktops []interface{}) []map[string]interface{} {
	if len(desktops) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(desktops))
	for _, desktop := range desktops {
		result = append(result, map[string]interface{}{
			"desktop_id":          utils.PathSearch("desktop_id", desktop, nil),
			"computer_name":       utils.PathSearch("computer_name", desktop, nil),
			"os_host_name":        utils.PathSearch("os_host_name", desktop, nil),
			"ip_addresses":        utils.PathSearch("ip_addresses", desktop, make([]interface{}, 0)),
			"ipv4":                utils.PathSearch("ipv4", desktop, nil),
			"ipv6":                utils.PathSearch("ipv6", desktop, nil),
			"desktop_type":        utils.PathSearch("desktop_type", desktop, nil),
			"status":              utils.PathSearch("status", desktop, nil),
			"in_maintenance_mode": utils.PathSearch("in_maintenance_mode", desktop, nil),
			"created":             utils.PathSearch("created", desktop, nil),
			"login_status":        utils.PathSearch("login_status", desktop, nil),
			"product_id":          utils.PathSearch("product_id", desktop, nil),
			"root_volume":         flattenDesktopPoolAssociatedDesktopRootVolume(utils.PathSearch("root_volume", desktop, nil)),
			"data_volumes": flattenDesktopPoolAssociatedDesktopVolumes(utils.PathSearch("data_volumes",
				desktop, make([]interface{}, 0)).([]interface{})),
			"availability_zone":     utils.PathSearch("availability_zone", desktop, nil),
			"site_type":             utils.PathSearch("site_type", desktop, nil),
			"site_name":             utils.PathSearch("site_name", desktop, nil),
			"product":               flattenDesktopPoolAssociatedDesktopProduct(utils.PathSearch("product", desktop, nil)),
			"os_version":            utils.PathSearch("os_version", desktop, nil),
			"sid":                   utils.PathSearch("sid", desktop, nil),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", desktop, make([]interface{}, 0)).([]interface{})),
			"is_support_internet":   utils.PathSearch("is_support_internet", desktop, nil),
			"is_attaching_eip":      utils.PathSearch("is_attaching_eip", desktop, nil),
			"attach_state":          utils.PathSearch("attach_state", desktop, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", desktop, nil),
			"subnet_id":             utils.PathSearch("subnet_id", desktop, nil),
			"bill_resource_id":      utils.PathSearch("bill_resource_id", desktop, nil),
		})
	}

	return result
}

func dataSourceDesktopPoolAssociatedDesktopsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	desktops, err := listDesktopPoolAssociatedDesktops(client, d.Get("pool_id").(string))
	if err != nil {
		return diag.Errorf("error querying the associated desktops under the Workspace desktop pool: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("desktops", flattenDesktopPoolAssociatedDesktops(desktops)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
