package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var batchImportKeypairNonUpdatableParams = []string{
	"keypairs",
	"keypairs.*.name",
	"keypairs.*.type",
	"keypairs.*.public_key",
	"keypairs.*.scope",
	"keypairs.*.user_id",
	"keypairs.*.key_protection",
	"keypairs.*.key_protection.*.private_key",
	"keypairs.*.key_protection.*.encryption",
	"keypairs.*.key_protection.*.encryption.*.type",
	"keypairs.*.key_protection.*.encryption.*.kms_key_name",
	"keypairs.*.key_protection.*.encryption.*.kms_key_id",
}

// @API DEW POST /v3/{project_id}/keypairs/batch-import
func ResourceKpsBatchImportKeypair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKpsBatchImportKeypairCreate,
		ReadContext:   resourceKpsBatchImportKeypairRead,
		UpdateContext: resourceKpsBatchImportKeypairUpdate,
		DeleteContext: resourceKpsBatchImportKeypairDelete,

		CustomizeDiff: config.FlexibleForceNew(batchImportKeypairNonUpdatableParams),

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
				Description: `Specifies the list of keypairs to import.`,
				Elem:        batchImportKpsKeypairsSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"succeeded_keypairs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of successfully imported keypairs.`,
				Elem:        batchImportKpsKeypairsSucceededKeypairsSchema(),
			},
			"failed_keypairs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of failed imported keypairs.`,
				Elem:        batchImportKpsKeypairsFailedKeypairsSchema(),
			},
		},
	}
}

func batchImportKpsKeypairsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the SSH keypair name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the SSH keypair type.`,
			},
			"public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the public key of the keypair.`,
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the scope of the keypair.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the user ID of the keypair.`,
			},
			"key_protection": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Specifies the key protection configuration.`,
				Elem:        keypairsKeyProtectionSchema(),
			},
		},
	}
}

func keypairsKeyProtectionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Specifies the private key of the keypair.`,
			},
			"encryption": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the encryption configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the encryption type.`,
						},
						"kms_key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the KMS key name.`,
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the KMS key ID.`,
						},
					},
				},
			},
		},
	}
}

func batchImportKpsKeypairsFailedKeypairsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"keypair_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SSH keypair name.`,
			},
			"failed_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The failed reason.`,
			},
		},
	}
}

func batchImportKpsKeypairsSucceededKeypairsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SSH keypair name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SSH keypair type.`,
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public key of the keypair.`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The private key of the keypair.`,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fingerprint of the keypair.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID of the keypair.`,
			},
		},
	}
}

func buildImportKeypairEncryptionBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"type":         rawMap["type"],
		"kms_key_name": utils.ValueIgnoreEmpty(rawMap["kms_key_name"]),
		"kms_key_id":   utils.ValueIgnoreEmpty(rawMap["kms_key_id"]),
	}
}

func buildImportKeypairKeyProtectionBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"private_key": utils.ValueIgnoreEmpty(rawMap["private_key"]),
		"encryption":  buildImportKeypairEncryptionBodyParams(rawMap["encryption"].([]interface{})),
	}
}

func buildImportKeypairsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawArray := d.Get("keypairs").([]interface{})
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		element := map[string]interface{}{
			"name":           rawMap["name"],
			"type":           utils.ValueIgnoreEmpty(rawMap["type"]),
			"public_key":     utils.ValueIgnoreEmpty(rawMap["public_key"]),
			"scope":          utils.ValueIgnoreEmpty(rawMap["scope"]),
			"user_id":        utils.ValueIgnoreEmpty(rawMap["user_id"]),
			"key_protection": buildImportKeypairKeyProtectionBodyParams(rawMap["key_protection"].([]interface{})),
		}

		rst = append(rst, map[string]interface{}{
			"keypair": utils.RemoveNil(element),
		})
	}

	return rst
}

func flattenBatchImportSucceededKeypairs(respArray []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("keypair.keypair.name", v, nil),
			"type":        utils.PathSearch("keypair.keypair.type", v, nil),
			"public_key":  utils.PathSearch("keypair.keypair.public_key", v, nil),
			"private_key": utils.PathSearch("keypair.keypair.private_key", v, nil),
			"fingerprint": utils.PathSearch("keypair.keypair.fingerprint", v, nil),
			"user_id":     utils.PathSearch("keypair.keypair.user_id", v, nil),
		})
	}
	return rst
}

func flattenBatchImportFailedKeypairs(respArray []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"keypair_name":  utils.PathSearch("keypair_name", v, nil),
			"failed_reason": utils.PathSearch("failed_reason", v, nil),
		})
	}
	return rst
}

func resourceKpsBatchImportKeypairCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/keypairs/batch-import"
		product = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW KPS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildImportKeypairsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch importing DEW KPS keypairs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error flattening response body: %s", err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(generateId)

	succeededKeypairs := utils.PathSearch("succeeded_keypairs", respBody, make([]interface{}, 0)).([]interface{})
	failedKeypairs := utils.PathSearch("failed_keypairs", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("succeeded_keypairs", flattenBatchImportSucceededKeypairs(succeededKeypairs)),
		d.Set("failed_keypairs", flattenBatchImportFailedKeypairs(failedKeypairs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKpsBatchImportKeypairRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKpsBatchImportKeypairUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKpsBatchImportKeypairDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to batch import keypairs. Deleting this resource
	will not recover the imported keypairs, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
