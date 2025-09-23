package codeartsbuild

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsBuild GET /v1/record/{build_project_id}/records
func DataSourceCodeArtsBuildTaskRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsBuildTaskRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"build_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the build project ID.`,
			},
			"triggers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of triggers to search.`,
			},
			"branches": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of branches to search.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of tags to search.`,
			},
			"from_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start date for the query, format: yyyy-MM-dd HH:mm:ss.`,
			},
			"to_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end date for the query, format: yyyy-MM-dd HH:mm:ss.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the build record list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the unique identifier of the build record.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the status of the build record.`,
						},
						"status_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the status code of the build record.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time of the build record.`,
						},
						"schedule_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the scheduled time of the build record.`,
						},
						"queued_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the queued time of the build record.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the start time of the build record.`,
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the finish time of the build record.`,
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the duration of the build record.`,
						},
						"build_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the build duration of the build record.`,
						},
						"pending_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the pending duration of the build record.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the project ID of the build record.`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the display name of the build record.`,
						},
						"trigger_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the trigger name of the build record.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the group name of the build record.`,
						},
						"execution_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the execution ID of the build record.`,
						},
						"parameters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the parameters of the build record.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the parameter name.`,
									},
									"secret": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Indicates whether the parameter is secret.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the parameter value.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the parameter type.`,
									},
								},
							},
						},
						"repository": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the repository of the build record.`,
						},
						"branch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the branch of the build record.`,
						},
						"revision": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the revision (commitId) of the build record.`,
						},
						"build_yml_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the build yaml path of the build record.`,
						},
						"build_yml_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the build yaml URL of the build record.`,
						},
						"daily_build_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the daily build number of the build record.`,
						},
						"build_record_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the parameters of the build record.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rerun": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Indicates whether the record is rerun.`,
									},
									"trigger_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the trigger type.`,
									},
									"record_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the record type.`,
									},
									"is_rerun": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Indicates the whether the record is rerun.`,
									},
								},
							},
						},
						"trigger_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the trigger type of the build record.`,
						},
						"scm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the SCM type of the build record.`,
						},
						"scm_web_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the SCM web URL of the build record.`,
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user ID of the build record.`,
						},
						"build_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the build number of the build record.`,
						},
						"daily_build_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the daily build number of the build record.`,
						},
						"dev_cloud_build_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the build type of the build record.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsBuildTaskRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_build", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	buildProjectId := d.Get("build_project_id").(string)
	getHttpUrl := "v1/record/{build_project_id}/records"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{build_project_id}", buildProjectId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPath += fmt.Sprintf("?limit=%v", 100)
	getPath += buildCodeArtsBuildTaskRecordsQueryParams(d)
	pageIndex := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&page=%d", pageIndex)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving build task records: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		records := utils.PathSearch("result.data", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}
		for _, record := range records {
			rst = append(rst, map[string]interface{}{
				"id":                   utils.PathSearch("id", record, nil),
				"status":               utils.PathSearch("status", record, nil),
				"status_code":          utils.PathSearch("status_code", record, nil),
				"create_time":          utils.PathSearch("create_time", record, nil),
				"schedule_time":        utils.PathSearch("schedule_time", record, nil),
				"queued_time":          utils.PathSearch("queued_time", record, nil),
				"start_time":           utils.PathSearch("start_time", record, nil),
				"finish_time":          utils.PathSearch("finish_time", record, nil),
				"duration":             utils.PathSearch("duration", record, nil),
				"build_duration":       utils.PathSearch("build_duration", record, nil),
				"pending_duration":     utils.PathSearch("pending_duration", record, nil),
				"project_id":           utils.PathSearch("project_id", record, nil),
				"display_name":         utils.PathSearch("display_name", record, nil),
				"trigger_name":         utils.PathSearch("trigger_name", record, nil),
				"group_name":           utils.PathSearch("group_name", record, nil),
				"execution_id":         utils.PathSearch("execution_id", record, nil),
				"repository":           utils.PathSearch("repository", record, nil),
				"branch":               utils.PathSearch("branch", record, nil),
				"revision":             utils.PathSearch("revision", record, nil),
				"build_yml_path":       utils.PathSearch("build_yml_path", record, nil),
				"build_yml_url":        utils.PathSearch("build_yml_url", record, nil),
				"daily_build_number":   utils.PathSearch("daily_build_number", record, nil),
				"trigger_type":         utils.PathSearch("trigger_type", record, nil),
				"scm_type":             utils.PathSearch("scm_type", record, nil),
				"scm_web_url":          utils.PathSearch("scm_web_url", record, nil),
				"user_id":              utils.PathSearch("user_id", record, nil),
				"build_no":             utils.PathSearch("build_no", record, nil),
				"daily_build_no":       utils.PathSearch("daily_build_no", record, nil),
				"dev_cloud_build_type": utils.PathSearch("dev_cloud_build_type", record, nil),
				"parameters":           flattenBuildTaskRecordsParameters(record),
				"build_record_type":    flattenBuildTaskRecordsBuildRecordType(record),
			})
		}

		pageIndex++
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

func buildCodeArtsBuildTaskRecordsQueryParams(d *schema.ResourceData) string {
	params := ""

	if v, ok := d.GetOk("triggers"); ok {
		triggers := v.([]interface{})
		for _, trigger := range triggers {
			params += fmt.Sprintf("&triggers=%v", trigger)
		}
	}
	if v, ok := d.GetOk("branches"); ok {
		branches := v.([]interface{})
		for _, branch := range branches {
			params += fmt.Sprintf("&branches=%v", branch)
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := v.([]interface{})
		for _, tag := range tags {
			params += fmt.Sprintf("&tags=%v", tag)
		}
	}
	if v, ok := d.GetOk("from_date"); ok {
		params += fmt.Sprintf("&from_date=%v", v)
	}
	if v, ok := d.GetOk("to_date"); ok {
		params += fmt.Sprintf("&to_date=%v", v)
	}

	return params
}

func flattenBuildTaskRecordsParameters(resp interface{}) []interface{} {
	parametersList := utils.PathSearch("parameters", resp, make([]interface{}, 0)).([]interface{})
	if len(parametersList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(parametersList))
	for _, p := range parametersList {
		parameter := p.(map[string]interface{})
		parameterMap := map[string]interface{}{
			"name":   utils.PathSearch("name", parameter, nil),
			"secret": utils.PathSearch("secret", parameter, nil),
			"value":  utils.PathSearch("value", parameter, nil),
			"type":   utils.PathSearch("type", parameter, nil),
		}
		result = append(result, parameterMap)
	}

	return result
}

func flattenBuildTaskRecordsBuildRecordType(resp interface{}) []map[string]interface{} {
	buildRecordType := utils.PathSearch("build_record_type", resp, nil)
	if buildRecordType == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"rerun":        utils.PathSearch("rerun", buildRecordType, nil),
			"trigger_type": utils.PathSearch("trigger_type", buildRecordType, nil),
			"record_type":  utils.PathSearch("record_type", buildRecordType, nil),
			"is_rerun":     utils.PathSearch("is_rerun", buildRecordType, nil),
		},
	}
}
