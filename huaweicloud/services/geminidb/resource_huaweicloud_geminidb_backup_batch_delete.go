package geminidb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var backupBatchDeleteNonUpdatableParams = []string{
	"backup_ids",
}

// @API GeminiDB DELETE /v3/{project_id}/instances/backups
func ResourceBackupBatchDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupBatchDeleteCreate,
		ReadContext:   resourceBackupBatchDeleteRead,
		UpdateContext: resourceBackupBatchDeleteUpdate,
		DeleteContext: resourceBackupBatchDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(backupBatchDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"backup_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func buildBackupBatchDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"backup_ids": utils.ExpandToStringList(d.Get("backup_ids").([]interface{})),
	}

	return bodyParams
}

func resourceBackupBatchDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/backups"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
		JSONBody:         buildBackupBatchDeleteBodyParams(d),
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting manual backups: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	return nil
}

func resourceBackupBatchDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBackupBatchDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBackupBatchDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GeminiDB backup batch delete resource is not supported. The resource is only removed from the " +
		"state, the resource remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
