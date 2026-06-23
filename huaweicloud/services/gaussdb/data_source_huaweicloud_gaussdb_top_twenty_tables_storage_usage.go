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

// @API GaussDB GET /v3.1/{project_id}/instances/{instance_id}/top-table-volume
func DataSourceTopTwentyTablesStorageUsage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTopTwentyTablesStorageUsageRead,

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
			"job_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topTableVolumeSchema(),
			},
		},
	}
}

func topTableVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_part_type": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_hash_cluster_key": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tuples": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"average_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"skew_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"skew_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"skew_stddev": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTopTwentyTablesStorageUsageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3.1/{project_id}/instances/{instance_id}/top-table-volume"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildTopTwentyTablesStorageUsageQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB top twenty tables storage usage: %s", err)
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
		d.Set("job_id", utils.PathSearch("job_id", getRespBody, nil)),
		d.Set("node_id", utils.PathSearch("node_id", getRespBody, nil)),
		d.Set("state", utils.PathSearch("state", getRespBody, nil)),
		d.Set("table_volumes", flattenTopTableVolumeBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildTopTwentyTablesStorageUsageQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("job_id"); ok {
		res += "&job_id=" + v.(string)
	}
	if v, ok := d.GetOk("node_id"); ok {
		res += "&node_id=" + v.(string)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenTopTableVolumeBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("table_volumes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"table_name":          utils.PathSearch("table_name", v, nil),
			"table_owner":         utils.PathSearch("table_owner", v, nil),
			"database_name":       utils.PathSearch("database_name", v, nil),
			"schema_name":         utils.PathSearch("schema_name", v, nil),
			"is_part_type":        utils.PathSearch("is_part_type", v, nil),
			"is_hash_cluster_key": utils.PathSearch("is_hash_cluster_key", v, nil),
			"tuples":              utils.PathSearch("tuples", v, nil),
			"create_time":         utils.PathSearch("create_time", v, nil),
			"update_time":         utils.PathSearch("update_time", v, nil),
			"average_size":        utils.PathSearch("average_size", v, nil),
			"max_ratio":           utils.PathSearch("max_ratio", v, nil),
			"min_ratio":           utils.PathSearch("min_ratio", v, nil),
			"table_size":          utils.PathSearch("table_size", v, nil),
			"skew_size":           utils.PathSearch("skew_size", v, nil),
			"skew_ratio":          utils.PathSearch("skew_ratio", v, nil),
			"skew_stddev":         utils.PathSearch("skew_stddev", v, nil),
		})
	}
	return rst
}
