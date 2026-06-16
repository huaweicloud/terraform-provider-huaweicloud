package css

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

// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/volume
func DataSourceCssClusterVolumeUsage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterVolumeUsageRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     clusterVolumeUsageDiskInfoSchema(),
			},
		},
	}
}

func clusterVolumeUsageDiskInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_used": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"percentage": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceClusterVolumeUsageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1.0/{project_id}/clusters/{cluster_id}/volume"
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving CSS cluster(%s) volume usage: %s", clusterId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("disk_info_list", flattenGetClusterVolumeUsageDiskInfoBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetClusterVolumeUsageDiskInfoBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("diskInfoList", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"group":         utils.PathSearch("group", v, nil),
			"role":          utils.PathSearch("role", v, nil),
			"disk_type":     utils.PathSearch("diskType", v, nil),
			"disk_capacity": utils.PathSearch("diskCapacity", v, nil),
			"disk_used":     utils.PathSearch("diskUsed", v, nil),
			"percentage":    utils.PathSearch("percentage", v, nil),
		})
	}
	return res
}
