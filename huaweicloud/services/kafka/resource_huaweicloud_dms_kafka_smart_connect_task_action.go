package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	// pauseTaskURL is the url to pause the task.
	pauseTaskURL = "v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}/pause"
	// resumeTaskURL is the url to resume the task.
	resumeTaskURL = "v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}/resume"
	// restartOrStartTaskURL is the url to restart or start the task.
	restartOrStartTaskURL = "v2/{project_id}/kafka/instances/{instance_id}/connector/tasks/{task_id}/restart"
)

var (
	actionURL = map[string]string{
		"pause":   pauseTaskURL,
		"resume":  resumeTaskURL,
		"restart": restartOrStartTaskURL,
		"start":   restartOrStartTaskURL,
	}
)

// @API Kafka PUT /v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}/pause
// @API Kafka PUT /v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}/resume
// @API Kafka PUT /v2/{project_id}/kafka/instances/{instance_id}/connector/tasks/{task_id}/restart
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}
func ResourceDmsKafkaSmartConnectTaskAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaSmartConnectTaskActionCreate,
		ReadContext:   resourceDmsKafkaSmartConnectTaskActionRead,
		UpdateContext: resourceDmsKafkaSmartConnectTaskActionUpdate,
		DeleteContext: resourceDmsKafkaSmartConnectTaskActionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
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
				Description: `Specifies the kafka instance ID.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the smart connect task ID.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the source type of the smart connect task.`,
				ValidateFunc: validation.StringInSlice([]string{
					"pause", "resume", "start", "restart",
				}, false),
			},
		},
	}
}

func resourceDmsKafkaSmartConnectTaskActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	taskID := d.Get("task_id").(string)
	action := d.Get("action").(string)

	err = createSmartConnectTaskAction(ctx, client, d.Timeout(schema.TimeoutCreate), instanceID, taskID, action)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID + "/" + taskID)

	return resourceDmsKafkaSmartConnectTaskActionRead(ctx, d, meta)
}

func createSmartConnectTaskAction(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceID, taskID, action string) error {
	createActionHttpUrl := actionURL[action]
	createActionPath := client.Endpoint + createActionHttpUrl
	createActionPath = strings.ReplaceAll(createActionPath, "{project_id}", client.ProjectID)
	createActionPath = strings.ReplaceAll(createActionPath, "{instance_id}", instanceID)
	createActionPath = strings.ReplaceAll(createActionPath, "{task_id}", taskID)
	createActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	if _, err := client.Request("PUT", createActionPath, &createActionOpt); err != nil {
		return fmt.Errorf("error creating DMS kafka smart connect task: %v", err)
	}

	var pending []string
	var target []string

	switch action {
	case "start":
		pending = []string{"WAITING"}
		target = []string{"RUNNING"}
	case "pause":
		pending = []string{"RUNNING"}
		target = []string{"PAUSED"}
	case "resume":
		pending = []string{"PAUSED"}
		target = []string{"RUNNING"}
	case "restart":
		pending = []string{"PAUSED", "RESTARTING"}
		target = []string{"RUNNING"}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      pending,
		Target:       target,
		Refresh:      kafkav2SmartConnectTaskStateRefreshFunc(client, instanceID, taskID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for the smart connect task action(%s) to be done: %s", action, err)
	}

	return nil
}

func resourceDmsKafkaSmartConnectTaskActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	taskID := d.Get("task_id").(string)
	action := d.Get("action").(string)

	if d.HasChange("action") {
		err = createSmartConnectTaskAction(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceID, taskID, action)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDmsKafkaSmartConnectTaskActionRead(ctx, d, meta)
}

func resourceDmsKafkaSmartConnectTaskActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsKafkaSmartConnectTaskActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting action resource is not supported. The action resource is only removed from the state," +
		" the task remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
