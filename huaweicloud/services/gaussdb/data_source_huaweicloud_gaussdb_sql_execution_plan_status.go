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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/sql/{node_id}/plans/query
func DataSourceSqlExecutionPlanStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlExecutionPlanStatusRead,

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
			"sql_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sql_plan_bind_state_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlExecutionPlanStatusSqlPlanBindStateListSchema(),
			},
		},
	}
}

func sqlExecutionPlanStatusSqlPlanBindStateListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"outline": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSqlExecutionPlanStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/sql/{node_id}/plans/query"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{node_id}", d.Get("node_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	offset := 0
	res := make([]interface{}, 0)
	for {
		listOpt.JSONBody = utils.RemoveNil(buildGetSqlExecutionPlanStatusBodyParams(d, offset))
		listResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB SQL execution plan status: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}
		sqlPlanBindStateList := flattenGetSqlExecutionPlanStatusSqlPlanBindStateListBody(listRespBody)
		if len(sqlPlanBindStateList) == 0 {
			break
		}
		res = append(res, sqlPlanBindStateList...)
		offset += 100
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("sql_plan_bind_state_list", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSqlExecutionPlanStatusBodyParams(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_id":    d.Get("sql_id"),
		"page_size": 100,
		"offset":    offset,
	}
	return bodyParams
}

func flattenGetSqlExecutionPlanStatusSqlPlanBindStateListBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("sql_plan_bind_state_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"outline":  utils.PathSearch("outline", v, nil),
			"cost":     utils.PathSearch("cost", v, nil),
			"status":   utils.PathSearch("status", v, nil),
			"sql_hash": utils.PathSearch("sql_hash", v, nil),
			"plan_id":  utils.PathSearch("plan_id", v, nil),
		})
	}
	return rst
}
