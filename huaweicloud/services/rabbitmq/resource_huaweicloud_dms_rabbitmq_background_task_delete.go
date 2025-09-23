package rabbitmq

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API RabbitMQ DELETE /v2/{project_id}/instances/{instance_id}/tasks/{task_id}
func ResourceDmsRabbitMQBackgroundTaskDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRabbitMQBackgroundTaskDeleteCreate,
		ReadContext:   resourceRabbitMQBackgroundTaskDeleteRead,
		DeleteContext: resourceRabbitMQBackgroundTaskDeleteDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRabbitMQBackgroundTaskDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	taskID := d.Get("task_id").(string)

	httpUrl := "v2/{project_id}/instances/{instance_id}/tasks/{task_id}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting RabbitMQ background task(%s): %s", taskID, err)
	}

	d.SetId(taskID)

	return nil
}

func resourceRabbitMQBackgroundTaskDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRabbitMQBackgroundTaskDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RabbitMQ background task delete resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
