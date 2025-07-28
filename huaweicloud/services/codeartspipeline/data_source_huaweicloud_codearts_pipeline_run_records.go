package codeartspipeline

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/list
func DataSourceCodeArtsPipelineRunRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineRunRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"pipeline_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline ID.`,
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of status.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting attribute.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting sequence.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the pipeline record list.`,
				Elem:        schemaPipelineRunRecords(),
			},
		},
	}
}

func schemaPipelineRunRecords() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pipeline_run_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline run ID.`,
			},
			"executor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the executor ID.`,
			},
			"executor_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the executor name.`,
			},
			"stage_status_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the stage information list.`,
				Elem:        schemaPipelineRunRecordsStageStatusList(),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of pipeline run.`,
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the trigger type.`,
			},
			"run_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the pipeline running sequence number.`,
			},
			"build_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the build parameters.`,
				Elem:        schemaPipelineRunDetailSourcesParamBuildParams(),
			},
			"artifact_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the artifacts after running a pipeline.`,
				Elem:        schemaPipelineRunRecordsArtifactParams(),
			},
			"start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the start time.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the end time.`,
			},
			"detail_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the address of the details page.`,
			},
			"modify_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the address of the editing page.`,
			},
		},
	}
}

func schemaPipelineRunRecordsStageStatusList() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the stage ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the stage name.`,
			},
			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the serial number.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the stage status.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the end time.`,
			},
		},
	}
}

func schemaPipelineRunRecordsArtifactParams() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"package_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the package name.`,
			},
			"branch_filter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the branch filter.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the package version.`,
			},
			"organization": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the docker organization.`,
			},
		},
	}
}

func dataSourceCodeArtsPipelineRunRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/list"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", d.Get("pipeline_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildCodeArtsPipelineRunRecordsQueryParams(d, offset))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving pipeline run records: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		if err := checkResponseError(getRespBody, ""); err != nil {
			return diag.Errorf("error retrieving pipeline run records: %s", err)
		}

		records := utils.PathSearch("pipeline_runs", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}
		for _, record := range records {
			rst = append(rst, flattenPipelineRunRecords(record))
		}

		offset += len(records)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCodeArtsPipelineRunRecordsQueryParams(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sort_dir":   utils.ValueIgnoreEmpty(d.Get("sort_dir")),
		"sort_key":   utils.ValueIgnoreEmpty(d.Get("sort_key")),
		"start_time": utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":   utils.ValueIgnoreEmpty(d.Get("end_time")),
		"status":     utils.ValueIgnoreEmpty(d.Get("status")),
		"limit":      100,
		"offset":     offset,
	}

	return bodyParams
}

func flattenPipelineRunRecords(resp interface{}) map[string]interface{} {
	return map[string]interface{}{
		"pipeline_run_id":   utils.PathSearch("pipeline_run_id", resp, nil),
		"executor_id":       utils.PathSearch("executor_id", resp, nil),
		"executor_name":     utils.PathSearch("executor_name", resp, nil),
		"stage_status_list": flattenPipelineRunRecordsStageStatusList(resp),
		"status":            utils.PathSearch("status", resp, nil),
		"trigger_type":      utils.PathSearch("trigger_type", resp, nil),
		"run_number":        utils.PathSearch("run_number", resp, nil),
		"build_params":      flattenPipelineRunBuildParams(resp),
		"artifact_params":   flattenPipelineRunRecordsArtifactParams(resp),
		"start_time":        utils.PathSearch("start_time", resp, nil),
		"end_time":          utils.PathSearch("end_time", resp, nil),
		"detail_url":        utils.PathSearch("detail_url", resp, nil),
		"modify_url":        utils.PathSearch("modify_url", resp, nil),
	}
}

func flattenPipelineRunRecordsStageStatusList(resp interface{}) []interface{} {
	stages := utils.PathSearch("stage_status_list", resp, make([]interface{}, 0)).([]interface{})
	if len(stages) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(stages))
	for _, s := range stages {
		stage := s.(map[string]interface{})
		m := map[string]interface{}{
			"id":         utils.PathSearch("id", stage, nil),
			"name":       utils.PathSearch("name", stage, nil),
			"sequence":   utils.PathSearch("sequence", stage, nil),
			"status":     utils.PathSearch("status", stage, nil),
			"start_time": utils.PathSearch("start_time", stage, nil),
			"end_time":   utils.PathSearch("end_time", stage, nil),
		}
		result = append(result, m)
	}

	return result
}

func flattenPipelineRunRecordsArtifactParams(resp interface{}) []map[string]interface{} {
	params := utils.PathSearch("artifact_params", resp, nil)
	if params == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"package_name":  utils.PathSearch("package_name", params, nil),
			"branch_filter": utils.PathSearch("branch_filter", params, nil),
			"version":       utils.PathSearch("version", params, nil),
			"organization":  utils.PathSearch("organization", params, nil),
		},
	}
}
