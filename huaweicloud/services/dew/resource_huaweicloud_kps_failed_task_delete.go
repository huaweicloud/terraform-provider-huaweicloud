package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var failedTaskDeleteNonUpdatableParams = []string{
	"task_id",
}

// @API DEW DELETE /v3/{project_id}/failed-tasks/{task_id}
func ResourceKpsFailedTaskDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKpsFailedTaskDeleteCreate,
		ReadContext:   resourceKpsFailedTaskDeleteRead,
		UpdateContext: resourceKpsFailedTaskDeleteUpdate,
		DeleteContext: resourceKpsFailedTaskDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(failedTaskDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the task ID.`,
			},
			"enable_force_new": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceKpsFailedTaskDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/failed-tasks/{task_id}"
		product = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KPS client: %s", err)
	}

	taskID := d.Get("task_id").(string)
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{task_id}", taskID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting failed task: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

func resourceKpsFailedTaskDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKpsFailedTaskDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKpsFailedTaskDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to delete a specific failed task.
Deleting this resource will not recover the deleted task, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
