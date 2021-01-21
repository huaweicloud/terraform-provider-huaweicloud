package huaweicloud

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/kms/v1/keys"
)

func dataSourceKmsKeyV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKmsKeyV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateKmsKeyStatus,
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
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceKmsKeyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsRegion := GetRegion(d, config)
	kmsKeyV1Client, err := config.kmsKeyV1Client(kmsRegion)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	is_list_key := true
	next_marker := ""
	allKeys := []keys.Key{}
	for is_list_key {
		req := &keys.ListOpts{
			KeyState: d.Get("key_state").(string),
			Limit:    "",
			Marker:   next_marker,
		}

		v, err := keys.List(kmsKeyV1Client, req).ExtractListKey()
		if err != nil {
			return err
		}

		is_list_key = v.Truncated == "true"
		next_marker = v.NextMarker
		allKeys = append(allKeys, v.KeyDetails...)
	}

	keyProperties := map[string]string{}
	if v, ok := d.GetOk("key_description"); ok {
		keyProperties["KeyDescription"] = v.(string)
	}
	if v, ok := d.GetOk("key_id"); ok {
		keyProperties["KeyID"] = v.(string)
	}
	if v, ok := d.GetOk("key_alias"); ok {
		keyProperties["KeyAlias"] = v.(string)
	}
	if v, ok := d.GetOk("default_key_flag"); ok {
		keyProperties["DefaultKeyFlag"] = v.(string)
	}
	if v, ok := d.GetOk("domain_id"); ok {
		keyProperties["DomainID"] = v.(string)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		keyProperties["EnterpriseProjectID"] = v.(string)
	}

	if len(allKeys) > 1 && len(keyProperties) > 0 {
		var filteredKeys []keys.Key
		for _, key := range allKeys {
			match := true
			for searchKey, searchValue := range keyProperties {
				r := reflect.ValueOf(&key)
				f := reflect.Indirect(r).FieldByName(searchKey)
				if !f.IsValid() {
					match = false
					break
				}

				keyValue := f.String()
				if searchValue != keyValue {
					match = false
					break
				}
			}

			if match {
				filteredKeys = append(filteredKeys, key)
			}
		}
		allKeys = filteredKeys
	}

	if len(allKeys) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allKeys) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	key := allKeys[0]
	log.Printf("[DEBUG] Kms key : %+v", key)

	d.SetId(key.KeyID)
	d.Set("key_id", key.KeyID)
	d.Set("domain_id", key.DomainID)
	d.Set("key_alias", key.KeyAlias)
	d.Set("region", kmsRegion)
	d.Set("key_description", key.KeyDescription)
	d.Set("creation_date", key.CreationDate)
	d.Set("scheduled_deletion_date", key.ScheduledDeletionDate)
	d.Set("key_state", key.KeyState)
	d.Set("default_key_flag", key.DefaultKeyFlag)
	d.Set("expiration_time", key.ExpirationTime)
	d.Set("enterprise_project_id", key.EnterpriseProjectID)

	if resourceTags, err := tags.Get(kmsKeyV1Client, "kms", key.KeyID).Extract(); err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmt.Errorf("Error saving tags to state for kms key(%s): %s", key.KeyID, err)
		}
	} else {
		log.Printf("[WARN] Error fetching tags of kms key(%s): %s", key.KeyID, err)
	}

	return nil
}
