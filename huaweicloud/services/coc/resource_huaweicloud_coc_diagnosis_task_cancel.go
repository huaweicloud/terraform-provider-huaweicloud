package coc

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var diagnosisTaskCancelNonUpdatableParams = []string{"task_id"}

// @API COC POST /v1/diagnosis/tasks/{task_id}/cancel
func ResourceDiagnosisTaskCancel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiagnosisTaskCancelCreate,
		ReadContext:   resourceDiagnosisTaskCancelRead,
		UpdateContext: resourceDiagnosisTaskCancelUpdate,
		DeleteContext: resourceDiagnosisTaskCancelDelete,

		CustomizeDiff: config.FlexibleForceNew(diagnosisTaskCancelNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDiagnosisTaskCancelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/diagnosis/tasks/{task_id}/cancel"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	taskID := d.Get("task_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error canceling the COC diagnosis task (%s): %s", taskID, err)
	}

	d.SetId(taskID)

	return resourceDiagnosisTaskCancelRead(ctx, d, meta)
}

func resourceDiagnosisTaskCancelRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDiagnosisTaskCancelUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDiagnosisTaskCancelDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting diagnosis task cancel resource is not supported. The diagnosis task cancel resource is only" +
		" removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
