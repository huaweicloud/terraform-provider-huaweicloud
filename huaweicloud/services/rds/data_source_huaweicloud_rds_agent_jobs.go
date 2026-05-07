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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/db-jobs
func DataSourceAgentJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAgentJobsRead,

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
			"job_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     agentJobSchema(),
			},
		},
	}
}

func agentJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"run_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_failure_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failure_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"agent_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAgentJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/db-jobs"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildGetAgentJobsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""},
	)
	if err != nil {
		return diag.Errorf("error retrieving RDS agent jobs: %s", err)
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
		d.Set("jobs", flattenGetAgentJobsBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAgentJobsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("job_type"); ok {
		res = fmt.Sprintf("%s&job_type=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetAgentJobsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("jobs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"job_id":            utils.PathSearch("job_id", v, nil),
			"job_name":          utils.PathSearch("job_name", v, nil),
			"is_enabled":        utils.PathSearch("is_enabled", v, nil),
			"run_time":          utils.PathSearch("run_time", v, nil),
			"run_status":        utils.PathSearch("run_status", v, nil),
			"last_failure_time": utils.PathSearch("last_failure_time", v, nil),
			"failure_count":     utils.PathSearch("failure_count", v, nil),
			"agent_type":        utils.PathSearch("agent_type", v, nil),
			"profile_id":        utils.PathSearch("profile_id", v, nil),
			"profile_name":      utils.PathSearch("profile_name", v, nil),
		})
	}
	return res
}
