package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/database-volume-summary
func DataSourceDataDiskSpaceUsage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataDiskSpaceUsageRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_disk_capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_disk_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"space_usage_growth_per_day": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"estimated_remaining_days": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cn_components": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDiskSpaceUsageCnComponentsSchema(),
			},
			"dn_components": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDiskSpaceUsageDnComponentsSchema(),
			},
		},
	}
}

func dataDiskSpaceUsageCnComponentsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataDiskSpaceUsageDnComponentsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDataDiskSpaceUsageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/database-volume-summary"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB data disk space usage: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("data_disk_capacity", utils.PathSearch("data_disk_capacity", getRespBody, nil)),
		d.Set("data_disk_usage", utils.PathSearch("data_disk_usage", getRespBody, nil)),
		d.Set("space_usage_growth_per_day", utils.PathSearch("space_usage_growth_per_day", getRespBody, nil)),
		d.Set("estimated_remaining_days", utils.PathSearch("estimated_remaining_days", getRespBody, nil)),
		d.Set("cn_components", flattenDataDiskSpaceUsageCnComponents(getRespBody)),
		d.Set("dn_components", flattenDataDiskSpaceUsageDnComponents(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataDiskSpaceUsageCnComponents(resp interface{}) []interface{} {
	curJson := utils.PathSearch("cn_components", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"node_id":      utils.PathSearch("node_id", v, nil),
			"component_id": utils.PathSearch("component_id", v, nil),
			"node_name":    utils.PathSearch("node_name", v, nil),
		})
	}
	return rst
}

func flattenDataDiskSpaceUsageDnComponents(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dn_components", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"node_id":      utils.PathSearch("node_id", v, nil),
			"component_id": utils.PathSearch("component_id", v, nil),
			"role":         utils.PathSearch("role", v, nil),
			"node_name":    utils.PathSearch("node_name", v, nil),
		})
	}
	return rst
}
