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

var geminiDbSessionsCloseParams = []string{
	"instance_id",
}

// @API GeminiDB DELETE /v3/{project_id}/instances/{instance_id}/sessions
func ResourceGeminiDBSessionsClose() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBSessionsCloseCreate,
		UpdateContext: resourceGeminiDBSessionsCloseUpdate,
		ReadContext:   resourceGeminiDBSessionsCloseRead,
		DeleteContext: resourceGeminiDBSessionsCloseDelete,

		CustomizeDiff: config.FlexibleForceNew(geminiDbSessionsCloseParams),

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

func resourceGeminiDBSessionsCloseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	httpUrl := "v3/{project_id}/instances/{instance_id}/sessions"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error closing sessions for instance %s: %s", instanceID, err)
	}

	d.SetId(instanceID)

	return resourceGeminiDBSessionsCloseRead(ctx, d, meta)
}

func resourceGeminiDBSessionsCloseRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBSessionsCloseUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBSessionsCloseDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting Geminidb closing sessions for instance resource is not supported. The Geminidb closing sessions for instance " +
		"resource is only removed from the state, the Geminidb nodes on an Instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
