package dataarts

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

// @API DataArtsStudio POST /v3/{project_id}/metadata/tasks/search
func DataSourceCatalogMetadataTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCatalogMetadataTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the metadata tasks are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the metadata tasks belongs.`,
			},

			// Optional parameters.
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user name which the metadata tasks are created.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the metadata task.",
			},
			"data_source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The data source type of the metadata tasks.",
			},
			"data_connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The data connection id of the metadata tasks.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The start time of the metadata tasks.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The end time of the metadata tasks.",
			},
			"directory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The directory ID of the metadata tasks.",
			},

			// Attributes.
			"metadata_tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataResourcesElem(),
				Description: `The list of metadata tasks that matched filter parameters.`,
			},
		},
	}
}

func dataResourcesElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the metadata task, in UUID format.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the metadata task.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the metadata task.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID which the metadata task is created.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the metadata task.",
			},
			"dir_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The directory ID of the metadata task.",
			},
			"schedule_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataMetadataTasksScheduleConfigElemSchema(),
				Description: "The dispatch information of the metadata task.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the metadata task.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user name which the metadata task is created.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The directory path of the metadata task.",
			},
			"last_run_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last run time of the metadata task.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start time of the metadata task.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The end time of the metadata task.",
			},
			"next_run_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The next run time of the metadata task.",
			},
			"duty_person": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The duty person of the metadata task.",
			},
			"data_source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data source type of the metadata task.",
			},
			"task_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The config information of the metadata task, in JSON format.",
			},
		},
	}
	return &sc
}

func dataMetadataTasksScheduleConfigElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cron_expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cron expression of the schedule task.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The end time of the schedule task.",
			},
			"max_time_out": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The max time out of the schedule task.",
			},
			"interval": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The interval time of the schedule task.",
			},
			"schedule_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The schedule type of the schedule task.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start time of the schedule task.",
			},
			"enabled": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether to enable the schedule task.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The job ID of the schedule task.",
			},
		},
	}
	return &sc
}

func buildCatalogMetadataTasksQueryParams(d *schema.ResourceData, limit, offset int) map[string]interface{} {
	return map[string]interface{}{
		"user_name":          d.Get("user_name"),
		"name":               d.Get("name"),
		"data_source_type":   d.Get("data_source_type"),
		"data_connection_id": d.Get("data_connection_id"),
		"start_time":         d.Get("start_time"),
		"end_time":           d.Get("end_time"),
		"directory_id":       d.Get("directory_id"),
		"limit":              limit,
		"offset":             offset,
	}
}

func listCatalogMetadataTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/metadata/tasks/search"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      bulidCatalogMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listOpt.JSONBody = utils.RemoveNil(buildCatalogMetadataTasksQueryParams(d, limit, offset))
		requestResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		metadataTasks := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, metadataTasks...)

		if len(metadataTasks) < limit {
			break
		}
		offset += len(metadataTasks)
	}

	return result, nil
}

func flattenCatalogMetadataTasks(metadataTasks []interface{}) []map[string]interface{} {
	if len(metadataTasks) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(metadataTasks))
	for _, metadataTask := range metadataTasks {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", metadataTask, nil),
			"name":        utils.PathSearch("name", metadataTask, nil),
			"description": utils.PathSearch("description", metadataTask, nil),
			"user_id":     utils.PathSearch("user_id", metadataTask, nil),
			"dir_id":      utils.PathSearch("dir_id", metadataTask, nil),
			"schedule_config": flattenScheduleConfig(utils.PathSearch("schedule_config", metadataTask,
				make(map[string]interface{})).(map[string]interface{})),
			"data_source_type": utils.PathSearch("data_source_type", metadataTask, nil),
			"task_config":      utils.JsonToString(utils.PathSearch("task_config", metadataTask, nil)),
			"create_time":      utils.PathSearch("create_time", metadataTask, nil),
			"update_time":      utils.PathSearch("update_time", metadataTask, nil),
			"user_name":        utils.PathSearch("user_name", metadataTask, nil),
			"path":             utils.PathSearch("path", metadataTask, nil),
			"last_run_time":    utils.PathSearch("last_run_time", metadataTask, nil),
			"start_time":       utils.PathSearch("start_time", metadataTask, nil),
			"end_time":         utils.PathSearch("end_time", metadataTask, nil),
			"next_run_time":    utils.PathSearch("next_run_time", metadataTask, nil),
			"duty_person":      utils.PathSearch("duty_person", metadataTask, nil),
		})
	}
	return result
}

func dataSourceCatalogMetadataTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	metadataTasks, err := listCatalogMetadataTasks(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Catalog metadata tasks: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("metadata_tasks", flattenCatalogMetadataTasks(metadataTasks)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
