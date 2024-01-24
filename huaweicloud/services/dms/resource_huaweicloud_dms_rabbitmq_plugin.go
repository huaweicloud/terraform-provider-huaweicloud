package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RabbitMQ PUT /v2/{project_id}/instances/{instance_id}/rabbitmq/plugins
// @API RabbitMQ GET /v2/{project_id}/instances/{instance_id}/rabbitmq/plugins
// @API RabbitMQ GET /v2/{project_id}/instances/{instance_id}/tasks
func ResourceDmsRabbitmqPlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRabbitmqPluginCreate,
		ReadContext:   resourceDmsRabbitmqPluginRead,
		DeleteContext: resourceDmsRabbitmqPluginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(50 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the RabbitMQ instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the plugin.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the plugin is enabled.`,
			},
			"running": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the plugin is running.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version of the plugin.`,
			},
		},
	}
}

func buildRabbitmqPluginBody(d *schema.ResourceData, enable bool) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"plugins": d.Get("name"),
		"enable":  enable,
	}
	return bodyParams
}

func resourceDmsRabbitmqPluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createRabbitmqPluginHttpUrl = "v2/{project_id}/instances/{instance_id}/rabbitmq/plugins"
		createRabbitmqPluginProduct = "dmsv2"
	)

	createRabbitmqPluginClient, err := cfg.NewServiceClient(createRabbitmqPluginProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createRabbitmqPluginPath := createRabbitmqPluginClient.Endpoint + createRabbitmqPluginHttpUrl
	createRabbitmqPluginPath = strings.ReplaceAll(createRabbitmqPluginPath, "{project_id}", createRabbitmqPluginClient.ProjectID)
	createRabbitmqPluginPath = strings.ReplaceAll(createRabbitmqPluginPath, "{instance_id}", instanceID)

	jsonBody := utils.RemoveNil(buildRabbitmqPluginBody(d, true))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBody,
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err = createRabbitmqPluginClient.Request("PUT", createRabbitmqPluginPath, &reqOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rabbitmqInstanceStateRefreshFunc(createRabbitmqPluginClient, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error enabling the plugin: %s", err)
	}

	name := d.Get("name").(string)
	id := fmt.Sprintf("%s/%s", instanceID, name)
	d.SetId(id)

	// The RabbitMQ enabling plugin is done if the status of its task is SUCCESS.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      rabbitmqPluginTaskRefreshFunc(createRabbitmqPluginClient, instanceID, name, "rabbitmqPluginModify"),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the RabbitMQ (%s) enabling plugin (%s) to be done: %s", instanceID, name, err)
	}

	return resourceDmsRabbitmqPluginRead(ctx, d, cfg)
}

func resourceDmsRabbitmqPluginRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getRabbitmqPluginHttpUrl = "v2/{project_id}/instances/{instance_id}/rabbitmq/plugins"
		getRabbitmqPluginProduct = "dms"
	)

	getRabbitmqPluginClient, err := cfg.NewServiceClient(getRabbitmqPluginProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceID := parts[0]
	name := parts[1]
	getRabbitmqPluginPath := getRabbitmqPluginClient.Endpoint + getRabbitmqPluginHttpUrl
	getRabbitmqPluginPath = strings.ReplaceAll(getRabbitmqPluginPath, "{project_id}", getRabbitmqPluginClient.ProjectID)
	getRabbitmqPluginPath = strings.ReplaceAll(getRabbitmqPluginPath, "{instance_id}", instanceID)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getRabbitmqPluginResp, err := getRabbitmqPluginClient.Request("GET", getRabbitmqPluginPath, &reqOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the plugin")
	}

	getRabbitmqPluginRespBody, err := utils.FlattenResponse(getRabbitmqPluginResp)
	if err != nil {
		return diag.FromErr(err)
	}

	plugin := utils.PathSearch(fmt.Sprintf("plugins|[?name=='%s']|[0]", name), getRabbitmqPluginRespBody, nil)
	enable := utils.PathSearch("enable", plugin, false).(bool)
	if !enable {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("name", name),
		d.Set("enable", enable),
		d.Set("version", utils.PathSearch("version", plugin, nil)),
		d.Set("running", utils.PathSearch("running", plugin, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRabbitmqPluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteRabbitmqPluginHttpUrl = "v2/{project_id}/instances/{instance_id}/rabbitmq/plugins"
		deleteRabbitmqPluginProduct = "dmsv2"
	)

	deleteRabbitmqPluginClient, err := cfg.NewServiceClient(deleteRabbitmqPluginProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}
	instanceID := d.Get("instance_id").(string)
	deleteRabbitmqPluginPath := deleteRabbitmqPluginClient.Endpoint + deleteRabbitmqPluginHttpUrl
	deleteRabbitmqPluginPath = strings.ReplaceAll(deleteRabbitmqPluginPath, "{project_id}", deleteRabbitmqPluginClient.ProjectID)
	deleteRabbitmqPluginPath = strings.ReplaceAll(deleteRabbitmqPluginPath, "{instance_id}", instanceID)
	jsonBody := utils.RemoveNil(buildRabbitmqPluginBody(d, false))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBody,
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err = deleteRabbitmqPluginClient.Request("PUT", deleteRabbitmqPluginPath, &reqOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rabbitmqInstanceStateRefreshFunc(deleteRabbitmqPluginClient, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error disabling the plugin: %s", err)
	}

	// The RabbitMQ disabling plugin is done if the status of its task is SUCCESS.
	name := d.Get("name").(string)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      rabbitmqPluginTaskRefreshFunc(deleteRabbitmqPluginClient, instanceID, name, "rabbitmqPluginModify"),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        1 * time.Second,
		PollInterval: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the RabbitMQ (%s) disabling plugin (%s) to be done: %s", instanceID, d.Get("name"), err)
	}

	return resourceDmsRabbitmqPluginRead(ctx, d, cfg)
}

func rabbitmqPluginTaskRefreshFunc(client *golangsdk.ServiceClient, instanceID, pluginName, taskName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// getRabbitmqPluginTask: get RabbitMQ plugin task
		getRabbitmqPluginTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/tasks"
		getRabbitmqPluginTaskPath := client.Endpoint + getRabbitmqPluginTaskHttpUrl
		getRabbitmqPluginTaskPath = strings.ReplaceAll(getRabbitmqPluginTaskPath, "{project_id}",
			client.ProjectID)
		getRabbitmqPluginTaskPath = strings.ReplaceAll(getRabbitmqPluginTaskPath, "{instance_id}", instanceID)

		reqOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getRabbitmqPluginTaskResp, err := client.Request("GET", getRabbitmqPluginTaskPath, &reqOpt)

		if err != nil {
			return nil, "QUERY ERROR", err
		}

		getRabbitmqPluginTaskRespBody, err := utils.FlattenResponse(getRabbitmqPluginTaskResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		task := getRabbitmqPluginTask(taskName, pluginName, getRabbitmqPluginTaskRespBody)
		if task == nil {
			return nil, "NOT FOUND", nil
		}
		status := utils.PathSearch("status", task, "").(string)
		return task, status, nil
	}
}

// resp is a struct including a list of tasks, the task example as below:
// {
// "name": "rabbitmqPluginModify",
// "params": "{\"enable\":true,\"plugins\":\"rabbitmq_delayed_message_exchange\"}\n",
// "status": "SUCCESS",
// ...
// }

func getRabbitmqPluginTask(taskName string, pluginName string, resp interface{}) interface{} {
	pluginTasks := utils.PathSearch(fmt.Sprintf("tasks|[?name=='%s']", taskName), resp, make([]interface{}, 0)).([]interface{})

	for _, task := range pluginTasks {
		params := utils.PathSearch("params", task, nil).(string)
		paramsData := []byte(params)
		var paramsJons interface{}
		err := json.Unmarshal(paramsData, &paramsJons)
		if err != nil {
			log.Printf("[WARN] error parsing params from the plugin task %#v", task)
			return nil
		}
		plugins := utils.PathSearch("plugins", paramsJons, nil)
		if pluginName == plugins {
			return task
		}
	}

	return nil
}
