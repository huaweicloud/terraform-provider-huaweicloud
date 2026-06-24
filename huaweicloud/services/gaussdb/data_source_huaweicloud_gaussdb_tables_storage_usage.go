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

// @API GaussDB POST /v3.1/{project_id}/instances/{instance_id}/table-volume
func DataSourceTablesStorageUsage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTablesStorageUsageRead,

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
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema_names": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tablesStorageUsageSchema(),
			},
		},
	}
}

func tablesStorageUsageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"schema_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_name": {
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

func dataSourceTablesStorageUsageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3.1/{project_id}/instances/{instance_id}/table-volume"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	offset := 0
	maxLimit := 100
	res := make([]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildGetTablesStorageUsageBodyParams(d, offset, maxLimit))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB tables storage usage: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		tablesStorageUsage := flattenTablesStorageUsageBody(getRespBody)
		if len(tablesStorageUsage) == 0 {
			break
		}
		res = append(res, tablesStorageUsage...)
		offset += maxLimit
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("table_volumes", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetTablesStorageUsageBodyParams(d *schema.ResourceData, offset, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"database_name": d.Get("database_name"),
		"schema_names":  d.Get("schema_names"),
		"table_name":    utils.ValueIgnoreEmpty(d.Get("table_name")),
		"user_name":     utils.ValueIgnoreEmpty(d.Get("user_name")),
		"sort_key":      utils.ValueIgnoreEmpty(d.Get("sort_key")),
		"sort_order":    utils.ValueIgnoreEmpty(d.Get("sort_order")),
		"offset":        offset,
		"limit":         limit,
	}

	return bodyParams
}

func flattenTablesStorageUsageBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("table_volumes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"table_size":          utils.PathSearch("table_size", v, nil),
			"id":                  utils.PathSearch("id", v, nil),
			"table_name":          utils.PathSearch("table_name", v, nil),
			"table_owner":         utils.PathSearch("table_owner", v, nil),
			"schema_name":         utils.PathSearch("schema_name", v, nil),
			"database_name":       utils.PathSearch("database_name", v, nil),
			"is_part_type":        utils.PathSearch("is_part_type", v, nil),
			"is_hash_cluster_key": utils.PathSearch("is_hash_cluster_key", v, nil),
			"tuples":              utils.PathSearch("tuples", v, nil),
			"create_time":         utils.PathSearch("create_time", v, nil),
			"update_time":         utils.PathSearch("update_time", v, nil),
			"average_size":        utils.PathSearch("average_size", v, nil),
			"max_ratio":           utils.PathSearch("max_ratio", v, nil),
			"min_ratio":           utils.PathSearch("min_ratio", v, nil),
			"skew_size":           utils.PathSearch("skew_size", v, nil),
			"skew_ratio":          utils.PathSearch("skew_ratio", v, nil),
			"skew_stddev":         utils.PathSearch("skew_stddev", v, nil),
		})
	}
	return rst
}
