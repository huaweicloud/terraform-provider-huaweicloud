package dew

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/kms/v1/keys"
	"github.com/chnsz/golangsdk/openstack/kms/v1/rotation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	WaitingForEnableState = "1"
	EnabledState          = "2"
	DisabledState         = "3"
	PendingDeletionState  = "4"
	PendingImportState    = "5"
)

// @API DEW POST /v1.0/{project_id}/kms/create-key
// @API DEW POST /v1.0/{project_id}/kms/disable-key
// @API DEW POST /v1.0/{project_id}/{resourceType}/{id}/tags/action
// @API DEW POST /v1.0/{project_id}/kms/enable-key-rotation
// @API DEW POST /v1.0/{project_id}/kms/update-key-rotation-interval
// @API DEW POST /v1.0/{project_id}/kms/describe-key
// @API DEW GET /v1.0/{project_id}/{resourceType}/{id}/tags
// @API DEW POST /v1.0/{project_id}/kms/get-key-rotation-status
// @API DEW POST /v1.0/{project_id}/kms/update-key-alias
// @API DEW POST /v1.0/{project_id}/kms/update-key-description
// @API DEW POST /v1.0/{project_id}/kms/enable-key
// @API DEW POST /v1.0/{project_id}/kms/disable-key-rotation
// @API DEW POST /v1.0/{project_id}/kms/schedule-key-deletion
// @API DEW POST /v1.0/{project_id}/kms/{id}/tags/action
func ResourceKmsKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceKmsKeyCreate,
		ReadContext:   ResourceKmsKeyRead,
		UpdateContext: ResourceKmsKeyUpdate,
		DeleteContext: ResourceKmsKeyDelete,
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
			"key_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"pending_days": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rotation_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rotation_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"rotation_enabled"},
				ValidateFunc: validation.IntBetween(30, 365),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_usage": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"keystore_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_deletion_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_key_flag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotation_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKmsKeyValidation(d *schema.ResourceData) error {
	_, rotationEnabled := d.GetOk("rotation_enabled")
	_, hasInterval := d.GetOk("rotation_interval")

	if !rotationEnabled && hasInterval {
		return fmt.Errorf("invalid arguments: rotation_interval is only valid when rotation is enabled")
	}
	return nil
}

func ResourceKmsKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	kmsKeyV1Client, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	if err := resourceKmsKeyValidation(d); err != nil {
		return diag.FromErr(err)
	}

	createOpts := &keys.CreateOpts{
		KeyAlias:            d.Get("key_alias").(string),
		KeyDescription:      d.Get("key_description").(string),
		KeySpec:             d.Get("key_algorithm").(string),
		KeyUsage:            d.Get("key_usage").(string),
		Origin:              d.Get("origin").(string),
		KeyStoreID:          d.Get("keystore_id").(string),
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, cfg),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := keys.Create(kmsKeyV1Client, createOpts).ExtractKeyInfo()
	if err != nil {
		return diag.Errorf("error creating KMS key: %s", err)
	}

	// Store the key ID
	d.SetId(v.KeyID)

	// Wait for the key to become enabled.
	log.Printf("[DEBUG] Waiting for KMS key (%s) to become enabled", v.KeyID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{WaitingForEnableState, DisabledState},
		Target:       []string{EnabledState, PendingImportState},
		Refresh:      keyV1StateRefreshFunc(kmsKeyV1Client, v.KeyID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 3 * time.Second,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for KMS key (%s) to become ready: %s", v.KeyID, err)
	}
	keyInfo := result.(*keys.Key)
	if !d.Get("is_enabled").(bool) && keyInfo.KeyState == EnabledState {
		key, err := keys.DisableKey(kmsKeyV1Client, v.KeyID).ExtractKeyInfo()
		if err != nil {
			return diag.Errorf("error disabling KMS key: %s", err)
		}

		if key.KeyState != DisabledState {
			return diag.Errorf("error disabling KMS key, the key state is: %s", key.KeyState)
		}
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		tagErr := tags.Create(kmsKeyV1Client, "kms", v.KeyID, taglist).ExtractErr()
		if tagErr != nil {
			return diag.Errorf("error creating tags of KMS key(%s): %s", v.KeyID, tagErr)
		}
	}

	// enable rotation and change interval if necessary
	// Only kms key support rotation
	if _, ok := d.GetOk("rotation_enabled"); ok && isKmsKey(d) {
		rotationOpts := &rotation.RotationOpts{
			KeyID: v.KeyID,
		}
		err := rotation.Enable(kmsKeyV1Client, rotationOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("failed to enable KMS key rotation: %s", err)
		}

		if i, ok := d.GetOk("rotation_interval"); ok {
			intervalOpts := &rotation.IntervalOpts{
				KeyID:    v.KeyID,
				Interval: i.(int),
			}
			err := rotation.Update(kmsKeyV1Client, intervalOpts).ExtractErr()
			if err != nil {
				return diag.Errorf("failed to change KMS key rotation interval: %s", err)
			}
		}
	}

	return ResourceKmsKeyRead(ctx, d, meta)
}

func ResourceKmsKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	kmsKeyV1Client, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS key client: %s", err)
	}

	v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "KMS key")
	}

	log.Printf("[DEBUG] Kms key %s: %+v", d.Id(), v)
	if v.KeyState == PendingDeletionState {
		log.Printf("[WARN] removing KMS key %s because it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(v.KeyID)
	mErr := multierror.Append(nil,
		d.Set("key_id", v.KeyID),
		d.Set("domain_id", v.DomainID),
		d.Set("key_alias", v.KeyAlias),
		d.Set("region", region),
		d.Set("key_description", v.KeyDescription),
		d.Set("key_algorithm", v.KeySpec),
		d.Set("creation_date", v.CreationDate),
		d.Set("scheduled_deletion_date", v.ScheduledDeletionDate),
		d.Set("default_key_flag", v.DefaultKeyFlag),
		d.Set("expiration_time", v.ExpirationTime),
		d.Set("enterprise_project_id", v.EnterpriseProjectID),
		d.Set("origin", v.Origin),
		d.Set("key_usage", v.KeyUsage),
		d.Set("key_state", v.KeyState),
		d.Set("keystore_id", v.KeyStoreID),
		utils.SetResourceTagsToState(d, kmsKeyV1Client, "kms", d.Id()),
	)

	if v.KeyState == EnabledState || v.KeyState == DisabledState {
		mErr = multierror.Append(mErr,
			d.Set("is_enabled", v.KeyState == EnabledState),
		)
	}

	// Set KMS rotation
	rotationOpts := &rotation.RotationOpts{
		KeyID: v.KeyID,
	}
	r, err := rotation.Get(kmsKeyV1Client, rotationOpts).Extract()
	if err == nil {
		mErr = multierror.Append(mErr,
			d.Set("rotation_enabled", r.Enabled),
			d.Set("rotation_interval", r.Interval),
			d.Set("rotation_number", r.NumberOfRotations),
		)
	} else {
		log.Printf("[WARN] error fetching details about KMS key rotation: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceKmsKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	kmsKeyV1Client, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS key client: %s", err)
	}

	if err := resourceKmsKeyValidation(d); err != nil {
		return diag.FromErr(err)
	}

	keyID := d.Id()
	if d.HasChange("key_alias") {
		updateAliasOpts := keys.UpdateAliasOpts{
			KeyID:    keyID,
			KeyAlias: d.Get("key_alias").(string),
		}
		_, err = keys.UpdateAlias(kmsKeyV1Client, updateAliasOpts).ExtractKeyInfo()
		if err != nil {
			return diag.Errorf("error updating KMS key: %s", err)
		}
	}

	if d.HasChange("key_description") {
		updateDesOpts := keys.UpdateDesOpts{
			KeyID:          keyID,
			KeyDescription: d.Get("key_description").(string),
		}
		_, err = keys.UpdateDes(kmsKeyV1Client, updateDesOpts).ExtractKeyInfo()
		if err != nil {
			return diag.Errorf("error updating KMS key: %s", err)
		}
	}

	keyState := d.Get("key_state").(string)
	if d.HasChange("is_enabled") {
		err := updateKeyState(d, kmsKeyV1Client, keyID, keyState)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(kmsKeyV1Client, d, "kms", keyID)
		if tagErr != nil {
			return diag.Errorf("error updating tags of kms: %s, err: %s", keyID, err)
		}
	}

	if isKmsKey(d) {
		err = updateRotation(d, kmsKeyV1Client, keyID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourceKmsKeyRead(ctx, d, meta)
}

func updateKeyState(d *schema.ResourceData, client *golangsdk.ServiceClient, keyID, keyState string) error {
	if d.Get("is_enabled").(bool) && keyState == DisabledState {
		key, err := keys.EnableKey(client, keyID).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("error enabling key: %s", err)
		}
		if key.KeyState != EnabledState {
			return fmt.Errorf("error enabling key, the key state is: %s", key.KeyState)
		}
	}

	if !d.Get("is_enabled").(bool) && keyState == EnabledState {
		key, err := keys.DisableKey(client, keyID).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("error disabling key: %s", err)
		}
		if key.KeyState != DisabledState {
			return fmt.Errorf("error disabling key, the key state is: %s", key.KeyState)
		}
	}

	return nil
}

func updateRotation(d *schema.ResourceData, client *golangsdk.ServiceClient, keyID string) error {
	rotationEnabled := d.Get("rotation_enabled").(bool)

	if d.HasChange("rotation_enabled") {
		var rotationErr error
		rotationOpts := &rotation.RotationOpts{
			KeyID: keyID,
		}
		if rotationEnabled {
			rotationErr = rotation.Enable(client, rotationOpts).ExtractErr()
		} else {
			rotationErr = rotation.Disable(client, rotationOpts).ExtractErr()
		}

		if rotationErr != nil {
			return fmt.Errorf("failed to update key rotation status: %s", rotationErr)
		}
	}

	if rotationEnabled && d.HasChange("rotation_interval") {
		intervalOpts := &rotation.IntervalOpts{
			KeyID:    keyID,
			Interval: d.Get("rotation_interval").(int),
		}
		err := rotation.Update(client, intervalOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("failed to change key rotation interval: %s", err)
		}
	}

	return nil
}

func ResourceKmsKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	kmsKeyV1Client, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS key client: %s", err)
	}

	v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "failed to retrieve key")
	}

	deleteOpts := &keys.DeleteOpts{
		KeyID:       d.Id(),
		PendingDays: "7",
	}
	if v, ok := d.GetOk("pending_days"); ok {
		deleteOpts.PendingDays = v.(string)
	}

	// It's possible that this key was used as a boot device and is currently
	// in a pending deletion state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.KeyState != PendingDeletionState {
		v, err = keys.Delete(kmsKeyV1Client, deleteOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		if v.KeyState != PendingDeletionState {
			return diag.Errorf("failed to delete KMS key")
		}
	}

	log.Printf("[DEBUG] KMS Key %s deactivated", d.Id())
	return nil
}

func keyV1StateRefreshFunc(client *golangsdk.ServiceClient, keyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := keys.Get(client, keyID).ExtractKeyInfo()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, PendingDeletionState, nil
			}
			return nil, "", err
		}
		return v, v.KeyState, nil
	}
}

func isKmsKey(d *schema.ResourceData) bool {
	if v, ok := d.GetOk("origin"); ok && v.(string) == "external" {
		return false
	}

	return true
}
