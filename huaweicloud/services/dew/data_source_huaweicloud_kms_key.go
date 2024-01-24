package dew

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/kms/v1/keys"
	"github.com/chnsz/golangsdk/openstack/kms/v1/rotation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1.0/{project_id}/kms/list-keys
// @API DEW POST /v1.0/{project_id}/kms/get-key-rotation-status
// @API DEW GET /v1.0/{project_id}/kms/{key_id}/tags
func DataSourceKmsKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKmsKeyRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_alias": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					EnabledState, DisabledState, PendingDeletionState,
				}, false),
			},
			"default_key_flag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rotation_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rotation_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rotation_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceKmsKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	kmsKeyV1Client, err := cfg.KmsKeyV1Client(region)
	if err != nil {
		return diag.Errorf("error creating kms key client: %s", err)
	}

	isListKey := true
	nextMarker := ""
	allKeys := []keys.Key{}
	for isListKey {
		req := &keys.ListOpts{
			KeyState: d.Get("key_state").(string),
			Limit:    "",
			Marker:   nextMarker,
		}

		v, err := keys.List(kmsKeyV1Client, req).ExtractListKey()
		if err != nil {
			return diag.FromErr(err)
		}

		isListKey = v.Truncated == "true"
		nextMarker = v.NextMarker
		allKeys = append(allKeys, v.KeyDetails...)
	}

	filter := map[string]interface{}{
		"KeyDescription":      d.Get("key_description"),
		"KeyID":               d.Get("key_id"),
		"KeyAlias":            d.Get("key_alias"),
		"DefaultKeyFlag":      d.Get("default_key_flag"),
		"DomainID":            d.Get("domain_id"),
		"EnterpriseProjectID": d.Get("enterprise_project_id"),
	}
	rst, err := utils.FilterSliceWithField(allKeys, filter)
	if err != nil {
		return diag.Errorf("erroring filting kms keey list: %s", err)
	}

	if len(rst) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(rst) > 1 {
		return diag.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	key := rst[0].(keys.Key)
	log.Printf("[DEBUG] Kms key : %+v", key)

	d.SetId(key.KeyID)
	d.Set("key_id", key.KeyID)
	d.Set("domain_id", key.DomainID)
	d.Set("key_alias", key.KeyAlias)
	d.Set("region", region)
	d.Set("key_description", key.KeyDescription)
	d.Set("creation_date", key.CreationDate)
	d.Set("scheduled_deletion_date", key.ScheduledDeletionDate)
	d.Set("key_state", key.KeyState)
	d.Set("default_key_flag", key.DefaultKeyFlag)
	d.Set("expiration_time", key.ExpirationTime)
	d.Set("enterprise_project_id", key.EnterpriseProjectID)

	if resourceTags, err := tags.Get(kmsKeyV1Client, "kms", key.KeyID).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for kms key(%s): %s", key.KeyID, err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of kms key(%s): %s", key.KeyID, err)
	}

	// Set KMS rotation
	rotationOpts := &rotation.RotationOpts{
		KeyID: key.KeyID,
	}
	r, err := rotation.Get(kmsKeyV1Client, rotationOpts).Extract()
	if err == nil {
		d.Set("rotation_enabled", r.Enabled)
		if r.Enabled {
			d.Set("rotation_interval", r.Interval)
			d.Set("rotation_number", r.NumberOfRotations)
		}
	} else {
		log.Printf("[WARN] error fetching details about key rotation: %s", err)
	}

	return nil
}
