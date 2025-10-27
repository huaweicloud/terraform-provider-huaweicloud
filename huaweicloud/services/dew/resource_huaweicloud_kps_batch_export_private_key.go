package dew

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var batchExportPrivateKeyNonUpdatableParams = []string{
	"keypairs",
	"keypairs.*.name",
	"keypairs.*.type",
	"keypairs.*.scope",
	"keypairs.*.public_key",
	"keypairs.*.fingerprint",
	"keypairs.*.is_key_protection",
	"keypairs.*.frozen_state",
	"export_file_name",
}

// @API DEW POST /v3/{project_id}/keypairs/private-key/batch-export
func ResourceKpsBatchExportPrivateKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKpsBatchExportPrivateKeyCreate,
		ReadContext:   resourceKpsBatchExportPrivateKeyRead,
		UpdateContext: resourceKpsBatchExportPrivateKeyUpdate,
		DeleteContext: resourceKpsBatchExportPrivateKeyDelete,

		CustomizeDiff: config.FlexibleForceNew(batchExportPrivateKeyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			"keypairs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the list of keypairs to export.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the SSH keypair name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the SSH keypair type.`,
						},
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the scope of the keypair.`,
						},
						"public_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the public key of the keypair.`,
						},
						"fingerprint": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the fingerprint of the keypair.`,
						},
						"is_key_protection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether the key is protected.`,
						},
						"frozen_state": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the frozen state of the keypair.`,
						},
					},
				},
			},
			"export_file_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the directory file for storing exported key pairs, and requires the file ending in .zip`,
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

func buildBatchExportPrivateKeyRequestBody(d *schema.ResourceData) []map[string]interface{} {
	keypairs := d.Get("keypairs").([]interface{})
	rstArray := make([]map[string]interface{}, 0, len(keypairs))
	for _, v := range keypairs {
		keypairMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rstArray = append(rstArray, map[string]interface{}{
			"keypair": map[string]interface{}{
				"name":              utils.ValueIgnoreEmpty(keypairMap["name"]),
				"type":              utils.ValueIgnoreEmpty(keypairMap["type"]),
				"scope":             utils.ValueIgnoreEmpty(keypairMap["scope"]),
				"public_key":        utils.ValueIgnoreEmpty(keypairMap["public_key"]),
				"fingerprint":       utils.ValueIgnoreEmpty(keypairMap["fingerprint"]),
				"is_key_protection": utils.ValueIgnoreEmpty(keypairMap["is_key_protection"]),
				"frozen_state":      utils.ValueIgnoreEmpty(keypairMap["frozen_state"]),
			},
		})
	}
	return rstArray
}

func resourceKpsBatchExportPrivateKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		httpUrl        = "v3/{project_id}/keypairs/private-key/batch-export"
		product        = "kms"
		exportFileName = d.Get("export_file_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW KPS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildBatchExportPrivateKeyRequestBody(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch exporting DEW KPS private keys: %s", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error reading response body: %s", err)
	}

	if err := os.WriteFile(exportFileName, bodyBytes, 0600); err != nil {
		return diag.Errorf("failed to write zip file to (%s): %s", exportFileName, err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceKpsBatchExportPrivateKeyRead(ctx, d, meta)
}

func resourceKpsBatchExportPrivateKeyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKpsBatchExportPrivateKeyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKpsBatchExportPrivateKeyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to batch export private keys. Deleting this resource
	will not recover the exported keys, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
