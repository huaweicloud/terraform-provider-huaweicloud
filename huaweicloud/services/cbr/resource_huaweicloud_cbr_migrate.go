package cbr

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var nonUpdatableMigrateParams = []string{
	"all_regions",
	"reservation",
}

// @API CBR POST /v3/migrates
func ResourceMigrate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrateCreate,
		ReadContext:   resourceMigrateRead,
		UpdateContext: resourceMigrateUpdate,
		DeleteContext: resourceMigrateDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableMigrateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			"all_regions": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether to trigger migration in other regions.`,
			},
			"reservation": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: `Specifies the default expansion ratio of the vault.`,
			},
		},
	}
}

func buildMigrateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"all_regions": d.Get("all_regions"),
		"reservation": d.Get("reservation"),
	}
	return bodyParams
}

func resourceMigrateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/migrates"
		product = "cbr"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildMigrateBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error migrating CBR resources: %s", err)
	}

	// Generate a UUID for the resource ID
	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	return resourceMigrateRead(ctx, d, meta)
}

func resourceMigrateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceMigrateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceMigrateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to migrate CBR resources. Deleting this 
resource will not change the current migration result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
