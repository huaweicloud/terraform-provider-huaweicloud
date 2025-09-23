package kafka

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v2/{project_id}/connectors/{connector_id}/sink-tasks
func DataSourceDmsKafkaSmartConnectTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsKafkaSmartConnectTasksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connector_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the connector ID of the kafka instance.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the smart connect task.`,
			},
			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the smart connect task.`,
			},
			"destination_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the destination type of the smart connect task.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the smart connect task.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Elem:        tasksSchema(),
				Computed:    true,
				Description: `Indicates the list of the smart connect tasks.`,
			},
		},
	}
}

func tasksSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the smart connect task.`,
			},
			"task_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the smart connect task.`,
			},
			"destination_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the destination type of the smart connect task.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the smart connect task.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the smart connect task.`,
			},
			"topics": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the topic names separated by commas or the topic regular expression of the smart connect task.`,
			},
		},
	}
	return &sc
}

func resourceDmsKafkaSmartConnectTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getKafkaSmartConnectTasks: query DMS kafka smart connect tasks
	var (
		getKafkaSmartConnectTasksHttpUrl = "v2/{project_id}/connectors/{connector_id}/sink-tasks"
		getKafkaSmartConnectTasksProduct = "dms"
	)
	getKafkaSmartConnectTasksClient, err := cfg.NewServiceClient(getKafkaSmartConnectTasksProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	connectorID := d.Get("connector_id").(string)
	getKafkaSmartConnectTasksPath := getKafkaSmartConnectTasksClient.Endpoint + getKafkaSmartConnectTasksHttpUrl
	getKafkaSmartConnectTasksPath = strings.ReplaceAll(getKafkaSmartConnectTasksPath, "{project_id}",
		getKafkaSmartConnectTasksClient.ProjectID)
	getKafkaSmartConnectTasksPath = strings.ReplaceAll(getKafkaSmartConnectTasksPath, "{connector_id}", connectorID)

	getKafkaSmartConnectTasksOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getKafkaSmartConnectTasksResp, err := getKafkaSmartConnectTasksClient.Request("GET", getKafkaSmartConnectTasksPath,
		&getKafkaSmartConnectTasksOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS kafka smart connect tasks")
	}

	getKafkaSmartConnectTasksRespBody, respBodyerr := utils.FlattenResponse(getKafkaSmartConnectTasksResp)
	if respBodyerr != nil {
		return diag.FromErr(respBodyerr)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("tasks", flattenTasks(filterTasks(d, getKafkaSmartConnectTasksRespBody))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTasks(taskRespBody []interface{}) []interface{} {
	if taskRespBody == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(taskRespBody))
	for _, v := range taskRespBody {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("task_id", v, nil),
			"task_name":        utils.PathSearch("task_name", v, nil),
			"destination_type": utils.PathSearch("destination_type", v, nil),
			"created_at":       utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", v, float64(0)).(float64)/1000), false),
			"status":           utils.PathSearch("status", v, nil),
			"topics":           utils.PathSearch("topics", v, nil),
		})
	}
	return rst
}

func filterTasks(d *schema.ResourceData, resp interface{}) []interface{} {
	taskJson := utils.PathSearch("tasks", resp, make([]interface{}, 0))
	taskArray := taskJson.([]interface{})
	if len(taskArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(taskArray))

	rawTaskId, rawTaskIdOK := d.GetOk("task_id")
	rawTaskName, rawTaskNameOK := d.GetOk("task_name")
	rawDestinationType, rawDestinationTypeOK := d.GetOk("destination_type")
	rawStatus, rawStatusOK := d.GetOk("status")

	for _, task := range taskArray {
		taskID := utils.PathSearch("task_id", task, nil)
		taskName := utils.PathSearch("task_name", task, nil)
		destinationType := utils.PathSearch("destination_type", task, nil)
		status := utils.PathSearch("status", task, nil)

		if rawTaskIdOK && rawTaskId != taskID {
			continue
		}
		if rawTaskNameOK && rawTaskName != taskName {
			continue
		}
		if rawDestinationTypeOK && rawDestinationType != destinationType {
			continue
		}
		if rawStatusOK && rawStatus != status {
			continue
		}
		result = append(result, task)
	}

	return result
}
