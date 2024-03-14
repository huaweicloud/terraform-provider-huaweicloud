package dew

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/kms/v1/keys"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DEW POST /v1.0/{project_id}/kms/delete-imported-key-material
// @API DEW POST /v1.0/{project_id}/kms/import-key-material
// @API DEW POST /v1.0/{project_id}/kms/describe-key
func ResourceKmsKeyMaterial() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceKmsKeyMaterialCreate,
		ReadContext:   ResourceKmsKeyMaterialRead,
		DeleteContext: ResourceKmsKeyMaterialDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"import_token": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encrypted_key_material": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encrypted_privatekey": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"key_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceKmsKeyMaterialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	kmsKeyV1Client, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	importMaterialOpts := &keys.ImportMaterialOpts{
		KeyID:                d.Get("key_id").(string),
		ImportToken:          d.Get("import_token").(string),
		EncryptedKeyMaterial: d.Get("encrypted_key_material").(string),
		EncryptedPrivatekey:  d.Get("encrypted_privatekey").(string),
		ExpirationTime:       d.Get("expiration_time").(string),
	}

	v, err := keys.ImportKeyMaterial(kmsKeyV1Client, importMaterialOpts).Extract()
	if err != nil || v == nil {
		return diag.Errorf("error importing KMS key material: %s", err)
	}

	d.SetId(d.Get("key_id").(string))

	return ResourceKmsKeyMaterialRead(ctx, d, meta)
}

func ResourceKmsKeyMaterialRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	kmsKeyV1Client, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS key client: %s", err)
	}

	v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the KMS key")
	}
	if v.KeyState == PendingDeletionState || v.KeyState == PendingImportState {
		return common.CheckDeletedDiag(d, err,
			"The KMS key is pending deletion or the key material is pending import")
	}

	expirationTime := flatternExpirationTime(v.ExpirationTime)

	d.SetId(v.KeyID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("key_id", v.KeyID),
		d.Set("key_state", v.KeyState),
		d.Set("expiration_time", expirationTime),
		d.Set("key_usage", v.KeyUsage),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flatternExpirationTime(expTimeStr string) string {
	if expTimeStr == "" {
		return ""
	}
	expTime, _ := strconv.ParseInt(expTimeStr, 10, 64)
	return strconv.FormatInt(expTime/1000, 10)
}

func ResourceKmsKeyMaterialDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	kmsKeyV1Client, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS key client: %s", err)
	}

	v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the KMS key")
	}
	if v.KeyState == PendingDeletionState || v.KeyState == PendingImportState {
		return common.CheckDeletedDiag(d, err, "The KMS key is pending deletion or the key material is deleted")
	}

	deleteMaterialOpts := &keys.DeleteKeyMaterialOpts{
		KeyID: d.Id(),
	}

	// The key material of the asymmetric key does not support deletion.
	// Deleting the key material of an asymmetric key will return {"error":{"error_msg":"xx","error_code":"KMS.2702"}}
	// The key material of the symmetric key support deletion.
	_, err = keys.DeleteKeyMaterial(kmsKeyV1Client, deleteMaterialOpts).Extract()
	if _, ok := err.(golangsdk.ErrDefault400); ok {
		errCode, errMessage := parseDeleteResponseError(err)
		if errCode == "KMS.2702" {
			log.Printf("[WARN] failed to delete key material, errCode : %s, errMsg: %s", errCode, errMessage)
			errorMessage := "The asymmetric key material can't be deleted. The project is only removed from the state," +
				" but it remains in the cloud."
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  errorMessage,
				},
			}
		}
	}
	if err != nil {
		diag.Errorf("error deleting key material (%s): %s", d.Id(), err)
	}

	return nil
}

func parseDeleteResponseError(err error) (errorCode, errorMsg string) {
	var response interface{}
	if jsonErr := json.Unmarshal(err.(golangsdk.ErrDefault400).Body, &response); jsonErr == nil {
		errorCode, parseErr := jmespath.Search("error.error_code", response)
		if parseErr != nil {
			log.Printf("[WARN] failed to parse error_code from response body: %s", parseErr)
		}
		errMsg, parseErr := jmespath.Search("error.error_msg", response)
		if parseErr != nil {
			log.Printf("[WARN] failed to parse error_msg from response body: %s", parseErr)
		}
		return errorCode.(string), errMsg.(string)
	}
	log.Printf("[WARN] failed to parse KMS error message from response body: %s", err)
	return "", ""
}
