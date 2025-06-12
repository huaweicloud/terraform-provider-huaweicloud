package cbr

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableBatchUpdateVaultParams = []string{
	"smn_notify",
	"threshold",
}

// @API CBR PUT /v3/{project_id}/vaults/batch-update
func ResourceBatchUpdateVault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchUpdateVaultCreate,
		ReadContext:   resourceBatchUpdateVaultRead,
		UpdateContext: resourceBatchUpdateVaultUpdate,
		DeleteContext: resourceBatchUpdateVaultDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableBatchUpdateVaultParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Description: `Specifies the region in which to execute the request. 
					If omitted, the provider-level region will be used.`,
			},
			"smn_notify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable SMN notification for the vault.`,
			},
			"threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the threshold of the vault capacity in GB.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"updated_vaults_id": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of vault IDs that have been updated.`,
			},
		},
	}
}

func buildBatchUpdateVaultBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"vault": map[string]interface{}{
			"smn_notify": d.Get("smn_notify"),
			"threshold":  utils.ValueIgnoreEmpty(d.Get("threshold")),
		},
	}
}

func resourceBatchUpdateVaultCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v3/{project_id}/vaults/batch-update"
		region  = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildBatchUpdateVaultBodyParams(d)),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error executing batch update CBR vault operation: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	updatedVaults := utils.PathSearch("updated_vaults_id", respBody, nil)
	if err = d.Set("updated_vaults_id", updatedVaults); err != nil {
		log.Printf("[DEBUG] error setting updated_vaults_id: %s", err)
	}

	// Generate a UUID for the resource ID
	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	return nil
}

func resourceBatchUpdateVaultRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchUpdateVaultUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchUpdateVaultDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to batch update CBR vaults. 
Deleting this resource will not change the current vault configuration, but will only remove the resource 
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
