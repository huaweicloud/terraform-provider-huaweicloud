package dew

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	scopeDomainValue = "domain"
	scopeDomainLabel = "account"
)

// @API DEW POST /v3/{project_id}/keypairs
// @API DEW PUT /v3/{project_id}/keypairs/{keypair_name}
// @API DEW GET /v3/{project_id}/keypairs/{keypair_name}
// @API DEW DELETE /v3/{project_id}/keypairs/{keypair_name}
// @API DEW POST /v3/{project_id}/keypairs/private-key/import
// @API DEW DELETE /v3/{project_id}/keypairs/{keypair_name}/private-key
func ResourceKeypair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeypairCreate,
		UpdateContext: resourceKeypairUpdate,
		DeleteContext: resourceKeypairDelete,
		ReadContext:   resourceKeypairRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encryption_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kms_key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"key_file"},
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_managed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildKeypairScopeParams(d *schema.ResourceData) string {
	scope := d.Get("scope").(string)
	if scope == scopeDomainLabel {
		return scopeDomainValue
	}

	return scope
}

func buildKeypairKeyProtectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	// The API requires that `encryption_type` must be configured before other fields can be configured.
	encryptionType, ok := d.GetOk("encryption_type")
	if !ok {
		return nil
	}

	encryptionBodyParams := map[string]interface{}{
		"type":         encryptionType,
		"kms_key_name": utils.ValueIgnoreEmpty(d.Get("kms_key_name")),
		"kms_key_id":   utils.ValueIgnoreEmpty(d.Get("kms_key_id")),
	}

	return map[string]interface{}{
		"private_key": utils.ValueIgnoreEmpty(d.Get("private_key")),
		"encryption":  encryptionBodyParams,
	}
}

func buildKeypairBodyParams(d *schema.ResourceData) map[string]interface{} {
	keypairParams := map[string]interface{}{
		"name":           d.Get("name"),
		"type":           "ssh",
		"public_key":     utils.ValueIgnoreEmpty(d.Get("public_key")),
		"user_id":        utils.ValueIgnoreEmpty(d.Get("user_id")),
		"scope":          utils.ValueIgnoreEmpty(buildKeypairScopeParams(d)),
		"key_protection": buildKeypairKeyProtectionBodyParams(d),
	}

	return map[string]interface{}{
		"keypair": keypairParams,
	}
}

func preCheckCreateArguments(d *schema.ResourceData) error {
	if v, ok := d.GetOk("encryption_type"); ok && v.(string) == "kms" {
		kmsKeyID := d.Get("kms_key_id").(string)
		kmsKeyName := d.Get("kms_key_name").(string)

		if kmsKeyID == "" && kmsKeyName == "" {
			return fmt.Errorf("'kms_key_name' or 'kms_key_id' is mandatory when the 'encryption_type' value is 'kms'")
		}
	}

	return nil
}

func buildUpdateKeypairDescriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	descriptionBodyParam := map[string]interface{}{
		"description": d.Get("description"),
	}

	return map[string]interface{}{
		"keypair": descriptionBodyParam,
	}
}

func updateKeypairDescription(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/{project_id}/keypairs/{keypair_name}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{keypair_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateKeypairDescriptionBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating keypair description: %s", err)
	}

	return nil
}

func writePrivateToKeyFile(respBody interface{}, d *schema.ResourceData, filePath string) error {
	privateKey := utils.PathSearch("keypair.private_key", respBody, "").(string)
	if privateKey == "" {
		log.Printf("[DEBUG] unable to write private key to key file: Private key is empty in API response")
		return nil
	}

	if err := utils.WriteToPemFile(filePath, privateKey); err != nil {
		return fmt.Errorf("unable to generate private key: %s", err)
	}

	if err := d.Set("key_file", filePath); err != nil {
		log.Printf("[DEBUG] error setting key_file attribute to local state: %s", err)
	}

	return nil
}

func resourceKeypairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/keypairs"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	if err := preCheckCreateArguments(d); err != nil {
		return diag.FromErr(err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildKeypairBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating KPS keypair: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	keyPairName := utils.PathSearch("keypair.name", respBody, "").(string)
	if keyPairName == "" {
		return diag.Errorf("error creating KPS key pair: keypair name is not found in API response")
	}
	d.SetId(keyPairName)

	// update description
	if _, ok := d.GetOk("description"); ok {
		if err := updateKeypairDescription(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// write private key to local. only when it is not import public_key and the key_file is not empty
	if filePath, ok := d.GetOk("key_file"); ok {
		if err := writePrivateToKeyFile(respBody, d, filePath.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceKeypairRead(ctx, d, meta)
}

func flattenScopeAttribute(scope string) string {
	if scope == scopeDomainValue {
		scope = scopeDomainLabel
	}

	return scope
}

func resourceKeypairRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/keypairs/{keypair_name}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{keypair_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving KPS keypair")
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	createTime := utils.PathSearch("keypair.create_time", respBody, float64(0)).(float64)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("keypair.name", respBody, nil)),
		d.Set("scope", flattenScopeAttribute(utils.PathSearch("keypair.scope", respBody, "").(string))),
		d.Set("public_key", utils.PathSearch("keypair.public_key", respBody, nil)),
		d.Set("description", utils.PathSearch("keypair.description", respBody, nil)),
		d.Set("user_id", utils.PathSearch("keypair.user_id", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampUTC(int64(createTime)/1000)),
		d.Set("fingerprint", utils.PathSearch("keypair.fingerprint", respBody, nil)),
		d.Set("is_managed", utils.PathSearch("keypair.is_key_protection", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func clearKeypairPrivateKey(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/{project_id}/keypairs/{keypair_name}/private-key"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{keypair_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		var errDefault404 golangsdk.ErrDefault404
		if errors.As(err, &errDefault404) {
			log.Printf("[DEBUG] The KPS keypair private key has already been cleared")
			return nil
		}
		return fmt.Errorf("error clearing KPS keypair private key: %s", err)
	}

	return nil
}

func buildImportPrivateKeyKeyProtectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	encryptionBodyParams := map[string]interface{}{
		"type":         d.Get("encryption_type"),
		"kms_key_name": utils.ValueIgnoreEmpty(d.Get("kms_key_name")),
		"kms_key_id":   utils.ValueIgnoreEmpty(d.Get("kms_key_id")),
	}

	return map[string]interface{}{
		"private_key": d.Get("private_key"),
		"encryption":  encryptionBodyParams,
	}
}

func buildImportPrivateKeyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"name":           d.Get("name"),
		"user_id":        utils.ValueIgnoreEmpty(d.Get("user_id")),
		"key_protection": buildImportPrivateKeyKeyProtectionBodyParams(d),
	}

	return map[string]interface{}{
		"keypair": bodyParam,
	}
}

func importKeypairPrivateKey(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/{project_id}/keypairs/private-key/import"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildImportPrivateKeyBodyParams(d)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error importing KPS keypair private key: %s", err)
	}

	return nil
}

func preCheckUpdatePrivateKeyArguments(d *schema.ResourceData) error {
	encryptionType := d.Get("encryption_type").(string)
	if encryptionType == "" {
		return fmt.Errorf("'encryption_type' is mandatory when importing private key")
	}

	if encryptionType == "kms" {
		kmsKeyID := d.Get("kms_key_id").(string)
		kmsKeyName := d.Get("kms_key_name").(string)

		if kmsKeyID == "" && kmsKeyName == "" {
			return fmt.Errorf("'kms_key_name' or 'kms_key_id' is mandatory when the 'encryption_type' value is 'kms'")
		}
	}

	return nil
}

func updateKeypairPrivateKey(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	privateKey := d.Get("private_key").(string)
	if privateKey == "" {
		if err := clearKeypairPrivateKey(client, d); err != nil {
			return err
		}
	}

	if privateKey != "" {
		if err := preCheckUpdatePrivateKeyArguments(d); err != nil {
			return err
		}

		if err := importKeypairPrivateKey(client, d); err != nil {
			return err
		}
	}

	return nil
}

func preCheckUpdateEncryptionAndUserIDArguments(d *schema.ResourceData) error {
	privateKey := d.Get("private_key").(string)
	if privateKey == "" {
		return fmt.Errorf("'private_key' is mandatory when updating encryption values or user ID")
	}

	encryptionType := d.Get("encryption_type").(string)
	if encryptionType == "" {
		return fmt.Errorf("'encryption_type' is mandatory when updating encryption values or user ID")
	}

	if encryptionType == "kms" {
		kmsKeyID := d.Get("kms_key_id").(string)
		kmsKeyName := d.Get("kms_key_name").(string)

		if kmsKeyID == "" && kmsKeyName == "" {
			return fmt.Errorf("'kms_key_name' or 'kms_key_id' is mandatory when the 'encryption_type' value is 'kms'")
		}
	}

	return nil
}

func updateKeypairEncryptionAndUserID(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	if err := preCheckUpdateEncryptionAndUserIDArguments(d); err != nil {
		return err
	}

	// Edit 'encryption_type', 'kms_key_id', 'kms_key_name', or 'user_id' by importing private key API
	return importKeypairPrivateKey(client, d)
}

func resourceKeypairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	if d.HasChange("description") {
		if err := updateKeypairDescription(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("private_key") {
		if err := updateKeypairPrivateKey(client, d); err != nil {
			return diag.Errorf("error updating KPS keypair private key: %s", err)
		}
	}

	// The API parameter restriction logic for importing private key is relatively complex.
	// For readability and maintainability, the logic for editing 'private_key' and other fields is separated.
	if d.HasChanges("encryption_type", "kms_key_name", "kms_key_id", "user_id") {
		if err := updateKeypairEncryptionAndUserID(client, d); err != nil {
			return diag.Errorf("error updating KPS keypair encryption type: %s", err)
		}
	}

	return resourceKeypairRead(ctx, d, meta)
}

func resourceKeypairDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/keypairs/{keypair_name}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{keypair_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting KPS keypair")
	}

	return nil
}
