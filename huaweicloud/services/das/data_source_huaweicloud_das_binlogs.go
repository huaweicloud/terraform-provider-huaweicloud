package das

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

// @API DAS POST /v3/{project_id}/connections/{connection_id}/binlog-parse/list-file
func DataSourceBinlogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBinlogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the DAS binlogs are located.",
			},

			// Required parameters.
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database user ID.",
			},
			"binlog_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The binlog file type.",
			},

			// Optional parameters.
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The start time of the query range, in RFC3339 format.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The end time of the query range, in RFC3339 format.",
			},

			// Attributes.
			"binlogs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of binlog files that matched the filter parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file name.",
						},
						"backup_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID that has already been backed up.",
						},
						"file_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The file size.",
						},
						"task_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The archive log parse information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The task ID.",
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The task creation time, in RFC3339 format.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The task modification time, in RFC3339 format.",
									},
									"project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tenant ID of the task.",
									},
									"project_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tenant name of the task.",
									},
									"user_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user ID of the task.",
									},
									"user_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user name of the task.",
									},
									"connection_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The connection ID of the task.",
									},
									"binlog_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The binlog type of the task.",
									},
									"file_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The binlog file name of the task.",
									},
									"backup_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The backup file ID of the task.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The status of the task.",
									},
									"err_msg": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The error message.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBinlogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	binlogs, err := listBinlogs(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS binlogs: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("binlogs", flattenBinlogs(binlogs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listBinlogs(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/connections/{connection_id}/binlog-parse/list-file"
		perPage = 100
		curPage = 1
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{connection_id}", d.Get("user_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			// Use this header key-value pair to make the API response formatted as underscores
			// eg: `fileName` to `file_name`
			"X-Source-Service": "das",
		},
	}

	for {
		requestBody := buildBinlogsRequestBody(d, curPage, perPage)
		opt.JSONBody = requestBody

		requestResp, err := client.Request("POST", listPath, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		fileList := utils.PathSearch("file_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, fileList...)

		if len(fileList) < perPage {
			break
		}
		curPage++
	}

	return result, nil
}

func buildBinlogsRequestBody(d *schema.ResourceData, curPage, perPage int) map[string]interface{} {
	body := map[string]interface{}{
		"binlog_type": d.Get("binlog_type").(string),
		"cur_page":    curPage,
		"per_page":    perPage,
	}

	if v, ok := d.GetOk("start_time"); ok {
		body["start_time"] = utils.ConvertTimeStrToNanoTimestamp(v.(string)) / 1000
	}
	if v, ok := d.GetOk("end_time"); ok {
		body["end_time"] = utils.ConvertTimeStrToNanoTimestamp(v.(string)) / 1000
	}

	return body
}

func flattenBinlogs(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"file_name": utils.PathSearch("file_name", item, nil),
			"file_size": utils.PathSearch("file_size", item, nil),
			"backup_id": utils.PathSearch("backup_id", item, nil),
			"task_info": flattenBinlogTaskInfo(
				utils.PathSearch("task_info", item, make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return result
}

func flattenBinlogTaskInfo(taskInfo map[string]interface{}) []map[string]interface{} {
	if len(taskInfo) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, 1)
	result = append(result, map[string]interface{}{
		"id":            utils.PathSearch("id", taskInfo, nil),
		"project_id":    utils.PathSearch("tenant_id", taskInfo, nil),
		"project_name":  utils.PathSearch("tenant_name", taskInfo, nil),
		"user_id":       utils.PathSearch("user_id", taskInfo, nil),
		"user_name":     utils.PathSearch("user_name", taskInfo, nil),
		"connection_id": utils.PathSearch("connection_id", taskInfo, nil),
		"binlog_type":   utils.PathSearch("binlog_type", taskInfo, nil),
		"file_name":     utils.PathSearch("file_name", taskInfo, nil),
		"backup_id":     utils.PathSearch("backup_id", taskInfo, nil),
		"status":        utils.PathSearch("status", taskInfo, nil),
		"err_msg":       utils.PathSearch("err_msg", taskInfo, nil),
		"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("gmt_create",
			taskInfo, float64(0)).(float64))/1000, false),
		"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("gmt_modified",
			taskInfo, float64(0)).(float64))/1000, false),
	})

	return result
}
