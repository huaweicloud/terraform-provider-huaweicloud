package rds

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/packlog/infos
func DataSourceMergedBinlogFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMergedBinlogFilesRead,

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
			"pack_log_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     mergedBinlogFilesSchema(),
			},
		},
	}
}

func mergedBinlogFilesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"size_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"query_end_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMergedBinlogFilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/packlog/infos"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""},
	)
	if err != nil {
		return diag.Errorf("error retrieving RDS merged binlog files: %s", err)
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

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("pack_log_infos", flattenGetMergedBinlogFilesBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetMergedBinlogFilesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("pack_log_infos", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"instance_id":      utils.PathSearch("instance_id", v, nil),
			"size":             utils.PathSearch("size", v, nil),
			"size_unit":        utils.PathSearch("size_unit", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"query_start_time": utils.PathSearch("query_start_time", v, nil),
			"query_end_time":   utils.PathSearch("query_end_time", v, nil),
			"file_name":        utils.PathSearch("file_name", v, nil),
		})
	}
	return res
}
