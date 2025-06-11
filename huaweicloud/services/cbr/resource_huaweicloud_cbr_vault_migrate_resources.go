package cbr

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableMigrationParams = []string{
	"vault_id",
	"destination_vault_id",
	"resource_ids",
}

// @API CBR POST /v3/{project_id}/vaults/{vault_id}/migrateresources
func ResourceVaultMigrateResources() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVaultMigrateResourcesCreate,
		ReadContext:   resourceVaultMigrateResourcesRead,
		UpdateContext: resourceVaultMigrateResourcesUpdate,
		DeleteContext: resourceVaultMigrateResourcesDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableMigrationParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level
region will be used.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the source vault ID from which resources will be migrated.`,
			},
			"destination_vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the destination vault ID where resources will be migrated to.`,
			},
			"resource_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the IDs of the resources to be migrated.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"migrated_resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Specifies the list of resources that have been successfully migrated.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func buildVaultMigrateResourcesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"destination_vault_id": d.Get("destination_vault_id"),
		"resource_ids":         d.Get("resource_ids"),
	}
}

func resourceVaultMigrateResourcesCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/vaults/{vault_id}/migrateresources"
		product = "cbr"
		vaultID = d.Get("vault_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{vault_id}", vaultID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildVaultMigrateResourcesBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error migrating CBR resources between vaults: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("migrated_resources", utils.PathSearch("migrated_resources", respBody, nil)); err != nil {
		log.Printf("[ERROR] error setting migrated_resources: %s", err)
	}

	d.SetId(vaultID)

	return nil
}

func resourceVaultMigrateResourcesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceVaultMigrateResourcesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceVaultMigrateResourcesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to migrate CBR resources between vaults. Deleting this 
resource will not change the current migration result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
