package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/db-jobs/{job_id}/histories
func DataSourceAgentJobHistories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAgentJobHistoriesRead,

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
				Required: true,
			},
			"run_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"histories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     agentJobHistoriesSchema(),
			},
		},
	}
}

func agentJobHistoriesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"history_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_duration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAgentJobHistoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/db-jobs/{job_id}/histories"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{job_id}", d.Get("job_id").(string))
	listPath += buildGetAgentJobHistoriesQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""},
	)
	if err != nil {
		return diag.Errorf("error retrieving RDS agent job histories: %s", err)
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

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("histories", flattenGetAgentJobHistoriesBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAgentJobHistoriesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("run_status"); ok {
		res = fmt.Sprintf("%s&run_status=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetAgentJobHistoriesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("histories", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"history_id":   utils.PathSearch("history_id", v, nil),
			"run_status":   utils.PathSearch("run_status", v, nil),
			"run_time":     utils.PathSearch("run_time", v, nil),
			"run_duration": utils.PathSearch("run_duration", v, nil),
			"message":      utils.PathSearch("message", v, nil),
		})
	}
	return res
}
