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
// @API DEW POST /v1.0/{project_id}/kms/describe-key
// @API DEW POST /v1.0/{project_id}/kms/update-key-alias
// @API DEW POST /v1.0/{project_id}/kms/update-key-description
// @API DEW POST /v1.0/{project_id}/kms/schedule-key-deletion
// @API DEW POST /v1.0/{project_id}/kms/enable-key
// @API DEW POST /v1.0/{project_id}/kms/disable-key
// @API DEW POST /v1.0/{project_id}/kms/enable-key-rotation
// @API DEW POST /v1.0/{project_id}/kms/disable-key-rotation
// @API DEW POST /v1.0/{project_id}/kms/get-key-rotation-status
// @API DEW POST /v1.0/{project_id}/kms/update-key-rotation-interval
// @API DEW POST /v1.0/{project_id}/kms/{key_id}/tags/action
// @API DEW GET /v1.0/{project_id}/kms/{key_id}/tags
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
			},
			"tags": common.TagsSchema(),
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
		return fmt.Errorf("invalid argument: 'rotation_interval' is only valid when the KMS key rotation is enabled")
	}
	return nil
}

func ResourceKmsKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	keyClient, err := cfg.KmsKeyV1Client(region)
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
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] The request body information for creating the KMS key: %#v", createOpts)
	resp, err := keys.Create(keyClient, createOpts).ExtractKeyInfo()
	if err != nil {
		return diag.Errorf("error creating KMS key: %s", err)
	}

	if resp.KeyID == "" {
		return diag.Errorf("error creating KMS key: ID is not found in API response")
	}

	d.SetId(resp.KeyID)

	// Wait for the key to become enabled.
	log.Printf("[DEBUG] Waiting for KMS key (%s) to become enabled", resp.KeyID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{WaitingForEnableState, DisabledState},
		Target:       []string{EnabledState, PendingImportState},
		Refresh:      keyV1StateRefreshFunc(keyClient, resp.KeyID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 3 * time.Second,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for KMS key (%s) to become ready: %s", resp.KeyID, err)
	}

	keyInfo := result.(*keys.Key)
	if !d.Get("is_enabled").(bool) && keyInfo.KeyState == EnabledState {
		key, err := keys.DisableKey(keyClient, resp.KeyID).ExtractKeyInfo()
		if err != nil {
			return diag.Errorf("error disabling KMS key: %s", err)
		}

		if key.KeyState != DisabledState {
			return diag.Errorf("error disabling KMS key, the key state is: %s", key.KeyState)
		}
	}

	if keyTags, ok := d.GetOk("tags"); ok {
		tagList := utils.ExpandResourceTags(keyTags.(map[string]interface{}))
		err = tags.Create(keyClient, "kms", d.Id(), tagList).ExtractErr()
		if err != nil {
			return diag.Errorf("error creating tags to the KMS key: %s", err)
		}
	}

	// enable rotation and change interval if necessary
	// Only kms key support rotation
	if _, ok := d.GetOk("rotation_enabled"); ok && isKmsKey(d) {
		rotationOpts := &rotation.RotationOpts{
			KeyID: resp.KeyID,
		}
		err := rotation.Enable(keyClient, rotationOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error enabling KMS key rotation: %s", err)
		}

		if v, ok := d.GetOk("rotation_interval"); ok {
			intervalOpts := &rotation.IntervalOpts{
				KeyID:    resp.KeyID,
				Interval: v.(int),
			}
			err := rotation.Update(keyClient, intervalOpts).ExtractErr()
			if err != nil {
				return diag.Errorf("error updating KMS key rotation interval: %s", err)
			}
		}
	}

	return ResourceKmsKeyRead(ctx, d, meta)
}

func ResourceKmsKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	keyClient, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS key client: %s", err)
	}

	key, err := keys.Get(keyClient, d.Id()).ExtractKeyInfo()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving KMS key")
	}

	if key.KeyState == PendingDeletionState {
		log.Printf("[WARN] Please remove the KMS key (%s), because it's already gone", d.Id())
		d.SetId("")
		return nil
	}
	d.SetId(key.KeyID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("key_id", key.KeyID),
		d.Set("domain_id", key.DomainID),
		d.Set("key_alias", key.KeyAlias),
		d.Set("key_description", key.KeyDescription),
		d.Set("key_algorithm", key.KeySpec),
		d.Set("creation_date", key.CreationDate),
		d.Set("scheduled_deletion_date", key.ScheduledDeletionDate),
		d.Set("default_key_flag", key.DefaultKeyFlag),
		d.Set("expiration_time", key.ExpirationTime),
		d.Set("enterprise_project_id", key.EnterpriseProjectID),
		d.Set("origin", key.Origin),
		d.Set("key_usage", key.KeyUsage),
		d.Set("key_state", key.KeyState),
		d.Set("keystore_id", key.KeyStoreID),
		utils.SetResourceTagsToState(d, keyClient, "kms", d.Id()),
	)

	if key.KeyState == EnabledState || key.KeyState == DisabledState {
		mErr = multierror.Append(mErr,
			d.Set("is_enabled", key.KeyState == EnabledState),
		)
	}

	// Set KMS key rotation
	rotationOpts := &rotation.RotationOpts{
		KeyID: key.KeyID,
	}

	resp, err := rotation.Get(keyClient, rotationOpts).Extract()
	if err == nil {
		mErr = multierror.Append(mErr,
			d.Set("rotation_enabled", resp.Enabled),
			d.Set("rotation_interval", resp.Interval),
			d.Set("rotation_number", resp.NumberOfRotations),
		)
	} else {
		log.Printf("[WARN] Error retrieving KMS key rotation information: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceKmsKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	keyClient, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS key client: %s", err)
	}

	if err := resourceKmsKeyValidation(d); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("key_alias") {
		updateAliasOpts := keys.UpdateAliasOpts{
			KeyID:    d.Id(),
			KeyAlias: d.Get("key_alias").(string),
		}
		_, err = keys.UpdateAlias(keyClient, updateAliasOpts).ExtractKeyInfo()
		if err != nil {
			return diag.Errorf("error updating KMS key: %s", err)
		}
	}

	if d.HasChange("key_description") {
		updateDesOpts := keys.UpdateDesOpts{
			KeyID:          d.Id(),
			KeyDescription: d.Get("key_description").(string),
		}
		_, err = keys.UpdateDes(keyClient, updateDesOpts).ExtractKeyInfo()
		if err != nil {
			return diag.Errorf("error updating KMS key: %s", err)
		}
	}

	keyState := d.Get("key_state").(string)
	if d.HasChange("is_enabled") {
		if err := updateKeyState(d, keyClient, d.Id(), keyState); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err := utils.UpdateResourceTags(keyClient, d, "kms", d.Id()); err != nil {
			return diag.Errorf("error updating tags of KMS key (%s): %s", d.Id(), err)
		}
	}

	if isKmsKey(d) {
		if err := updateRotation(d, keyClient, d.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "kms",
			RegionId:     region,
			ProjectId:    keyClient.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
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
			return fmt.Errorf("error updating KMS key rotation information: %s", rotationErr)
		}
	}

	if rotationEnabled && d.HasChange("rotation_interval") {
		intervalOpts := &rotation.IntervalOpts{
			KeyID:    keyID,
			Interval: d.Get("rotation_interval").(int),
		}
		err := rotation.Update(client, intervalOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("error updating KMS key rotation interval: %s", err)
		}
	}

	return nil
}

func ResourceKmsKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	keyClient, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS key client: %s", err)
	}

	resp, err := keys.Get(keyClient, d.Id()).ExtractKeyInfo()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving KMS key")
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
	if resp.KeyState != PendingDeletionState {
		resp, err = keys.Delete(keyClient, deleteOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		if resp.KeyState != PendingDeletionState {
			return diag.Errorf("error deleting KMS key")
		}
	}

	log.Printf("[DEBUG] The KMS Key (%s) deactivated", d.Id())
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
