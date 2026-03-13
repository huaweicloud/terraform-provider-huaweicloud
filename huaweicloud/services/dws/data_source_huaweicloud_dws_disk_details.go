package dws

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v1.0/{project_id}/dms/disk
func DataSourceDiskDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDiskDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the disk details are located.`,
			},

			// Optional parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The cluster ID.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The instance ID.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The instance name.`,
			},

			// Attributes.
			"disk_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        diskDetailSchema(),
				Description: `The list of disk details that matched filter parameters.`,
			},
		},
	}
}

func diskDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance name.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance ID.`,
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The host name.`,
			},
			"disk_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk name.`,
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk type.`,
			},
			"total": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The total disk capacity in GB.`,
			},
			"used": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The used disk capacity in GB.`,
			},
			"available": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The available disk capacity in GB.`,
			},
			"used_percentage": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The disk usage percentage.`,
			},
			"await": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The I/O wait time in milliseconds.`,
			},
			"svctm": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The I/O service time in milliseconds.`,
			},
			"util": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The I/O usage percentage.`,
			},
			"read_rate": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The disk read rate in KB/s.`,
			},
			"write_rate": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The disk write rate in KB/s.`,
			},
		},
	}
}

func buildDiskDetailsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("cluster_id"); ok {
		res = fmt.Sprintf("%s&cluster_id=%v", res, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		res = fmt.Sprintf("%s&instance_name=%v", res, v)
	}

	return res
}

func listDiskDetails(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1.0/{project_id}/dms/disk?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildDiskDetailsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPathWithLimit + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		diskDetails, ok := respBody.([]interface{})
		if !ok {
			return nil, fmt.Errorf("error retrieving disk details: unexpected response type %T", respBody)
		}

		result = append(result, diskDetails...)
		if len(diskDetails) < limit {
			break
		}
		offset += len(diskDetails)
	}

	return result, nil
}

func flattenDiskDetails(all []interface{}) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(all))
	for _, item := range all {
		result = append(result, map[string]interface{}{
			"instance_name":   utils.PathSearch("instance_name", item, nil),
			"instance_id":     utils.PathSearch("instance_id", item, nil),
			"host_name":       utils.PathSearch("host_name", item, nil),
			"disk_name":       utils.PathSearch("disk_name", item, nil),
			"disk_type":       utils.PathSearch("disk_type", item, nil),
			"total":           utils.PathSearch("total", item, nil),
			"used":            utils.PathSearch("used", item, nil),
			"available":       utils.PathSearch("available", item, nil),
			"used_percentage": utils.PathSearch("used_percentage", item, nil),
			"await":           utils.PathSearch("await", item, nil),
			"svctm":           utils.PathSearch("svctm", item, nil),
			"util":            utils.PathSearch("util", item, nil),
			"read_rate":       utils.PathSearch("read_rate", item, nil),
			"write_rate":      utils.PathSearch("write_rate", item, nil),
		})
	}

	return result
}

func dataSourceDiskDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	diskDetails, err := listDiskDetails(client, d)
	if err != nil {
		return diag.Errorf("error retrieving disk details: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("disk_details", flattenDiskDetails(diskDetails)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
