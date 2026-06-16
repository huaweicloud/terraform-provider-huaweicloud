package taurusdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/sql-filter/history-rules
func DataSourceSqlControlHistoryRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlControlHistoryRulesRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sql_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sql_filter_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlControlHistoryRulesSchema(),
			},
		},
	}
}

func sqlControlHistoryRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pattern": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_concurrency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expire_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"delete_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSqlControlHistoryRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/sql-filter/history-rules"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildGetSqlControlHistoryRulesQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB SQL control history rules: %s", err)
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
		d.Set("sql_filter_rules", flattenSqlControlHistoryRules(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSqlControlHistoryRulesQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&node_id=%v", res, d.Get("node_id").(string))
	if v, ok := d.GetOk("sql_type"); ok {
		res = fmt.Sprintf("%s&sql_type=%v", res, v.(string))
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenSqlControlHistoryRules(resp interface{}) []interface{} {
	curJson := utils.PathSearch("sql_filter_rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"instance_id":     utils.PathSearch("instance_id", v, nil),
			"node_id":         utils.PathSearch("node_id", v, nil),
			"pattern":         utils.PathSearch("pattern", v, nil),
			"sql_type":        utils.PathSearch("sql_type", v, nil),
			"max_concurrency": utils.PathSearch("max_concurrency", v, nil),
			"create_at":       utils.PathSearch("create_at", v, nil),
			"expire_at":       utils.PathSearch("expire_at", v, nil),
			"delete_at":       utils.PathSearch("delete_at", v, nil),
		})
	}
	return res
}
