package geminidb

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

var geminiDbScheduledTaskCancelParams = []string{
	"job_id",
}

// @API GeminiDB PUT /v3/{project_id}/instances/disaster-recovery/settings
func ResourceGeminiDBScheduledTaskCancel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBScheduledTaskCancelCreate,
		UpdateContext: resourceGeminiDBScheduledTaskCancelUpdate,
		ReadContext:   resourceGeminiDBScheduledTaskCancelRead,
		DeleteContext: resourceGeminiDBScheduledTaskCancelDelete,

		CustomizeDiff: config.FlexibleForceNew(geminiDbScheduledTaskCancelParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
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

func resourceGeminiDBScheduledTaskCancelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	jobID := d.Get("job_id").(string)

	httpUrl := "v3/{project_id}/scheduled-jobs/{job_id}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{job_id}", jobID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error canceling scheduled task: %s", err)
	}

	d.SetId(jobID)

	return resourceGeminiDBScheduledTaskCancelRead(ctx, d, meta)
}

func resourceGeminiDBScheduledTaskCancelUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBScheduledTaskCancelRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBScheduledTaskCancelDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting Geminidb scheduled task cancel resource is not supported. The Geminidb scheduled " +
		"task cancel resource is only removed from the state, the Geminidb scheduled task remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
