package cbr

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var replicateBackupNonUpdatableParams = []string{
	"backup_id",
	"replicate.*.destination_project_id",
	"replicate.*.destination_region",
	"replicate.*.destination_vault_id",
	"replicate.*.name",
	"replicate.*.description",
	"replicate.*.enable_acceleration",
}

// @API CBR POST /v3/{project_id}/backups/{backup_id}/replicate
func ResourceReplicateBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReplicateBackupCreate,
		ReadContext:   resourceReplicateBackupRead,
		UpdateContext: resourceReplicateBackupUpdate,
		DeleteContext: resourceReplicateBackupDelete,

		CustomizeDiff: config.FlexibleForceNew(replicateBackupNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replicate": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_project_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"destination_region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"destination_vault_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable_acceleration": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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

func buildReplicateBackupBodyParams(d *schema.ResourceData) map[string]interface{} {
	replicateRaw := d.Get("replicate").([]interface{})
	if len(replicateRaw) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"replicate": map[string]interface{}{},
	}

	replicate, ok := replicateRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	bodyParams["replicate"] = map[string]interface{}{
		"destination_project_id": replicate["destination_project_id"],
		"destination_region":     replicate["destination_region"],
		"destination_vault_id":   replicate["destination_vault_id"],
		"name":                   utils.ValueIgnoreEmpty(replicate["name"]),
		"description":            utils.ValueIgnoreEmpty(replicate["description"]),
		"enable_acceleration":    utils.ValueIgnoreEmpty(replicate["enable_acceleration"]),
	}

	return bodyParams
}

func resourceReplicateBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v3/{project_id}/backups/{backup_id}/replicate"
		product  = "cbr"
		backupId = d.Get("backup_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{backup_id}", backupId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildReplicateBackupBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error replicating CBR backup: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return resourceReplicateBackupRead(ctx, d, meta)
}

func resourceReplicateBackupRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceReplicateBackupUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceReplicateBackupDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to replicate backup.
Deleting this resource will not reset the replicated backup, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
