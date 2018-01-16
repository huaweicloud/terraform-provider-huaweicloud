package huaweicloud

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/kms/v3/keys"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

const WaitingForEnableState = "1"
const EnabledState = "2"
const DisabledState = "3"
const PendingDeletionState = "4"

func resourceKmsKeyV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceKmsKeyV3Create,
		Read:   resourceKmsKeyV3Read,
		Update: resourceKmsKeyV3Update,
		Delete: resourceKmsKeyV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_alias": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"key_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"realm": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_usage": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sequence": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"creation_date": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scheduled_deletion_date": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"default_key_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"expiration_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"origin": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"key_rotation_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"pending_days": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceKmsKeyV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsKeyV3Client, err := config.kmsKeyV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	createOpts := &keys.CreateOpts{
		KeyAlias:       d.Get("key_alias").(string),
		KeyDescription: d.Get("key_description").(string),
		Realm:          d.Get("realm").(string),
		KeyPolicy:      d.Get("key_policy").(string),
		KeyUsage:       d.Get("key_usage").(string),
		KeyType:        d.Get("key_type").(string),
		Sequence:       d.Get("sequence").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := keys.Create(kmsKeyV3Client, createOpts).ExtractKeyInfo()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud key: %s", err)
	}
	log.Printf("[INFO] Key ID: %s", v.KeyID)

	// Wait for the key to become enabled.
	log.Printf("[DEBUG] Waiting for leu (%s) to become enabled", v.KeyID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{WaitingForEnableState, DisabledState},
		Target:     []string{EnabledState},
		Refresh:    KeyV3StateRefreshFunc(kmsKeyV3Client, v.KeyID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for key (%s) to become ready: %s",
			v.KeyID, err)
	}

	// Store the key ID now
	d.SetId(v.KeyID)
	d.Set("key_id", v.KeyID)

	return resourceKmsKeyV3Read(d, meta)
}

func resourceKmsKeyV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	kmsKeyV3Client, err := config.kmsKeyV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}
	getOpts := &keys.ListOpts{
		KeyID: d.Id(),
	}
	v, err := keys.Get(kmsKeyV3Client, getOpts).ExtractKeyInfo()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Kms key %s: %+v", d.Id(), v)
	if v.KeyState == PendingDeletionState {
		log.Printf("[WARN] Removing KMS key %s because it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(v.KeyID)
	d.Set("key_id", v.KeyID)
	d.Set("domain_id", v.DomainID)
	d.Set("key_alias", v.KeyAlias)
	d.Set("realm", v.Realm)
	d.Set("key_description", v.KeyDescription)
	d.Set("creation_date", v.CreationDate)
	d.Set("scheduled_deletion_date", v.ScheduledDeletionDate)
	d.Set("key_state", v.KeyState)
	d.Set("default_key_flag", v.DefaultKeyFlag)
	d.Set("key_type", v.KeyType)
	d.Set("expiration_time", v.ExpirationTime)
	d.Set("origin", v.Origin)
	d.Set("key_rotation_enabled", v.KeyRotationEnabled)

	return nil
}

func resourceKmsKeyV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsKeyV3Client, err := config.kmsKeyV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	if d.HasChange("key_alias") {
		updateAliasOpts := keys.UpdateAliasOpts{
			KeyID:    d.Id(),
			KeyAlias: d.Get("key_alias").(string),
		}
		_, err = keys.UpdateAlias(kmsKeyV3Client, updateAliasOpts).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud key: %s", err)
		}
	}

	if d.HasChange("key_description") {
		updateDesOpts := keys.UpdateDesOpts{
			KeyID:          d.Id(),
			KeyDescription: d.Get("key_description").(string),
		}
		_, err = keys.UpdateDes(kmsKeyV3Client, updateDesOpts).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud key: %s", err)
		}
	}

	return resourceKmsKeyV3Read(d, meta)
}

func resourceKmsKeyV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsKeyV3Client, err := config.kmsKeyV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	getOpts := &keys.ListOpts{
		KeyID: d.Id(),
	}
	v, err := keys.Get(kmsKeyV3Client, getOpts).ExtractKeyInfo()
	if err != nil {
		return CheckDeleted(d, err, "key")
	}

	deleteOpts := &keys.DeleteOpts{
		KeyID: d.Id(),
	}
	if v, ok := d.GetOk("pending_days"); ok {
		deleteOpts.PendingDays = v.(string)
	}
	if v, ok := d.GetOk("sequence"); ok {
		deleteOpts.Sequence = v.(string)
	}

	// It's possible that this key was used as a boot device and is currently
	// in a pending deletion state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.KeyState != PendingDeletionState {
		v, err = keys.Delete(kmsKeyV3Client, deleteOpts).Extract()
		if err != nil {
			return err
		}

		if v.KeyState != PendingDeletionState {
			return fmt.Errorf("failed to delete key")
		}
	}

	log.Printf("[DEBUG] KMS Key %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func KeyV3StateRefreshFunc(client *gophercloud.ServiceClient, keyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := keys.Get(client, &keys.ListOpts{KeyID: keyID}).ExtractKeyInfo()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return v, PendingDeletionState, nil
			}
			return nil, "", err
		}

		return v, v.KeyState, nil
	}
}
