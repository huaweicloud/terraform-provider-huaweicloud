package taurusdb

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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/nodes
func DataSourceTaurusDBHtapStarrocksNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksNodesRead,

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
			"node_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     htapStarrocksNodesSchema(),
			},
		},
	}
}

func htapStarrocksNodesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBHtapStarrocksNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/nodes"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB HTAP StarRocks nodes: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("node_list", flattenHtapStarrocksNodesBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHtapStarrocksNodesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("node_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"node_id":   utils.PathSearch("node_id", v, nil),
			"node_name": utils.PathSearch("node_name", v, nil),
			"role":      utils.PathSearch("role", v, nil),
			"status":    utils.PathSearch("status", v, nil),
		})
	}
	return res
}
