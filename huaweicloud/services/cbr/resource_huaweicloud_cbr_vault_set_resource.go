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

var nonUpdatableVaultSetResourceParams = []string{
	"vault_id",
	"resource_ids",
	"action",
}

// @API CBR PUT /v3/{project_id}/vaults/{vault_id}/set-resources
func ResourceVaultSetResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVaultSetResourceCreate,
		ReadContext:   resourceVaultSetResourceRead,
		UpdateContext: resourceVaultSetResourceUpdate,
		DeleteContext: resourceVaultSetResourceDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableVaultSetResourceParams),

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
				Description: `Specifies ID of the CBR vault to configure resource backup settings.`,
			},
			"resource_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of resource IDs for which to configure backup settings.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the action to configure backup settings.`,
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

func buildVaultResourceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"resource_ids": d.Get("resource_ids"),
		"action":       d.Get("action"),
	}
}

func resourceVaultSetResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v3/{project_id}/vaults/{vault_id}/set-resources"
		region  = cfg.GetRegion(d)
		vaultId = d.Get("vault_id").(string)
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{vault_id}", vaultId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildVaultResourceBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error configuring resource backup settings: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceVaultSetResourceRead(ctx, d, meta)
}

func resourceVaultSetResourceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceVaultSetResourceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceVaultSetResourceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to configure resource backup settings for a CBR vault.
Deleting this resource will not change the backup settings of the resources, but will only remove the resource
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
