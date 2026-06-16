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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/auto-sql-limiting/log
func DataSourceTaurusDBSqlAutoThrottlingRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlAutoThrottlingRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the instance ID.",
			},
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the node ID. The node role must be the primary node.",
			},
			"logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of auto throttling records.",
				Elem:        sqlAutoThrottlingRecordsSchema(),
			},
		},
	}
}

func sqlAutoThrottlingRecordsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The node ID.",
			},
			"pattern": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SQL throttling rule.",
			},
			"sql_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The throttling type.",
			},
			"max_concurrency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of concurrent requests.",
			},
			"create_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timestamp when throttling starts.",
			},
			"expire_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timestamp when throttling expires.",
			},
		},
	}
}

func dataSourceSqlAutoThrottlingRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/auto-sql-limiting/log"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{node_id}", d.Get("node_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB SQL auto throttling records: %s", err)
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
		d.Set("logs", flattenSqlAutoThrottlingRecordsLogs(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSqlAutoThrottlingRecordsLogs(resp interface{}) []interface{} {
	curJson := utils.PathSearch("logs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"node_id":         utils.PathSearch("node_id", v, nil),
			"pattern":         utils.PathSearch("pattern", v, nil),
			"sql_type":        utils.PathSearch("sql_type", v, nil),
			"max_concurrency": utils.PathSearch("max_concurrency", v, nil),
			"create_at":       utils.PathSearch("create_at", v, nil),
			"expire_at":       utils.PathSearch("expire_at", v, nil),
		})
	}
	return res
}
