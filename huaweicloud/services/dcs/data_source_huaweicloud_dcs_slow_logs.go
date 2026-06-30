package dcs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances/{instance_id}/slowlog
func DataSourceDcsSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsSlowLogsRead,

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
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slowlogs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"command": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shard_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDcsSlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	httpUrl := "v2/{project_id}/instances/{instance_id}/slowlog"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceID)
	listPath += buildGetSlowLogsQueryParams(d)

	listResp, err := pagination.ListAllItems(client, "offset", listPath, &pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DCS slow logs: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
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
		d.Set("slowlogs", flattenDcsSlowLogs(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSlowLogsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?start_time=%v", d.Get("start_time"))
	res = fmt.Sprintf("%s&end_time=%v", res, d.Get("end_time"))

	if v, ok := d.GetOk("sort_key"); ok {
		res = fmt.Sprintf("%s&sort_key=%v", res, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		res = fmt.Sprintf("%s&sort_dir=%v", res, v)
	}

	return res
}

func flattenDcsSlowLogs(resp interface{}) []interface{} {
	curJson := utils.PathSearch("slowlogs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"command":     utils.PathSearch("command", v, nil),
			"start_time":  utils.PathSearch("start_time", v, nil),
			"duration":    utils.PathSearch("duration", v, nil),
			"shard_name":  utils.PathSearch("shard_name", v, nil),
			"database_id": utils.PathSearch("database_id", v, nil),
			"username":    utils.PathSearch("username", v, nil),
			"node_role":   utils.PathSearch("node_role", v, nil),
			"client_ip":   utils.PathSearch("client_ip", v, nil),
		})
	}
	return res
}
