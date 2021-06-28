package huaweicloud

import (
	"math"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cbr/v3/policies"
	"github.com/huaweicloud/golangsdk/openstack/cbr/v3/vaults"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

var resourceType map[string]string = map[string]string{
	"server": "OS::Nova::Server",
	"disk":   "OS::Cinder::Volume",
	"turbo":  "OS::Sfs::Turbo",
}

func resourceCBRVaultV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCBRVaultV3Create,
		Read:   resourceCBRVaultV3Read,
		Update: resourceCBRVaultV3Update,
		Delete: resourceCBRVaultV3Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//If the validation content has changed, please update the resource type map.
				ValidateFunc: validation.StringInSlice([]string{
					"server", "disk", "turbo",
				}, false),
			},
			"consistent_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"crash_consistent", "app_consistent",
				}, false),
			},
			"protection_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"backup", "replication",
				}, false),
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 10485760),
			},
			"auto_expand": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"exclude_volumes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"include_volumes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"allocated": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceCBRVaultV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR v3 client: %s", err)
	}

	opts := vaults.CreateOpts{
		Name:                d.Get("name").(string),
		AutoExpand:          d.Get("auto_expand").(bool),
		BackupPolicyID:      d.Get("policy_id").(string),
		EnterpriseProjectID: GetEnterpriseProjectID(d, config),
		Resources:           buildCBRVaultResources(d),
		Billing:             buildCBRVaultBilling(d),
	}

	vault, err := vaults.Create(client, opts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating vaults: %s", err)
	}
	d.SetId(vault.ID)

	if policy, ok := d.GetOk("policy_id"); ok {
		_, err := vaults.BindPolicy(client, d.Id(), vaults.BindPolicyOpts{PolicyID: policy.(string)}).Extract()
		if err != nil {
			return fmtp.Errorf("Error binding policy to vault: %s", err)
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		tagRaw := v.(map[string]interface{})
		taglist := utils.ExpandResourceTags(tagRaw)
		tagErr := tags.Create(client, "vault", d.Id(), taglist).ExtractErr()
		if tagErr != nil {
			return fmtp.Errorf("Error setting tags of CBR vault: %s", tagErr)
		}
	}

	return resourceCBRVaultV3Read(d, meta)
}

func resourceCBRVaultV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR v3 client: %s", err)
	}

	vault, err := vaults.Get(client, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting vault details: %s", err)
	}

	resourceList := make([]map[string]interface{}, len(vault.Resources))
	for i, v := range vault.Resources {
		resource := map[string]interface{}{
			"id":             v.ID,
			"name":           v.Name,
			"protect_status": v.ProtectStatus,
			"size":           v.Size,
			"backup_size":    v.BackupSize,
			"backup_count":   v.BackupCount,
		}
		excludeVolumes := make([]string, len(v.ExtraInfo.ExcludeVolumes))
		for i, v := range v.ExtraInfo.ExcludeVolumes {
			excludeVolumes[i] = v
		}
		resource["exclude_volumes"] = excludeVolumes
		includeVolumes := make([]string, len(v.ExtraInfo.IncludeVolumes))
		for i, v := range v.ExtraInfo.IncludeVolumes {
			includeVolumes[i] = v.ID
		}
		resource["include_volumes"] = includeVolumes

		resourceList[i] = resource
	}

	mErr := multierror.Append(
		//required && optional
		d.Set("name", vault.Name),
		d.Set("consistent_level", vault.Billing.ConsistentLevel),
		d.Set("type", vault.Billing.ObjectType),
		d.Set("protection_type", vault.Billing.ProtectType),
		d.Set("resources", resourceList),
		d.Set("size", vault.Billing.Size),
		d.Set("auto_expand", vault.AutoExpand),
		d.Set("enterprise_project_id", vault.EnterpriseProjectID),
		d.Set("tags", utils.TagsToMap(vault.Tags)),
		//computed
		//The result of 'allocated' and 'used' is in MB, and now we need to use GB as the unit.
		d.Set("allocated", getNumberInGB(float64(vault.Billing.Allocated))),
		d.Set("used", getNumberInGB(float64(vault.Billing.Used))),
		d.Set("spec_code", vault.Billing.SpecCode),
		d.Set("status", vault.Billing.Status),
		d.Set("storage", vault.Billing.StorageUnit),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting vault fields: %s", err)
	}
	listOpts := policies.ListOpts{
		VaultID: d.Id(),
	}
	allPages, err := policies.List(client, listOpts).AllPages()
	if err != nil {
		return fmtp.Errorf("Error getting policy by ID (%s): %s", d.Id(), err)
	}
	policyList, err := policies.ExtractPolicies(allPages)
	if len(policyList) == 1 {
		d.Set("policy_id", policyList[0].ID)
	}

	return nil
}

func resourceCBRVaultV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR v3 client: %s", err)
	}

	if d.HasChanges("name", "size", "auto_expand") {
		opts := vaults.UpdateOpts{}
		opts.Name = d.Get("name").(string)
		opts.Billing = &vaults.BillingUpdate{
			Size: d.Get("size").(int),
		}
		ae := d.Get("auto_expand").(bool)
		opts.AutoExpand = &ae

		_, err := vaults.Update(client, d.Id(), opts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating the vault: %s", err)
		}
	}

	if d.HasChange("resources") {
		if err := updateResources(d, client); err != nil {
			return err
		}
	}
	if d.HasChange("policy_id") {
		if err := updatePolicy(d, client); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		if err = utils.UpdateResourceTags(client, d, "vault", d.Id()); err != nil {
			return fmtp.Errorf("Failed to update tags: %s", err)
		}
	}

	return resourceCBRVaultV3Read(d, meta)
}

func resourceCBRVaultV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR v3 client: %s", err)
	}

	if err := vaults.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Error deleting CBR v3 vault: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIncludeVolume(volumes []interface{}) []vaults.ResourceExtraInfoIncludeVolumes {
	includeVolumes := make([]vaults.ResourceExtraInfoIncludeVolumes, len(volumes))
	for i, v := range volumes {
		includeVolumes[i] = vaults.ResourceExtraInfoIncludeVolumes{
			ID: v.(string),
		}
	}
	return includeVolumes
}

func resourceExcludeVolume(volumes []interface{}) []string {
	includeVolumes := make([]string, len(volumes))
	for i, v := range volumes {
		includeVolumes[i] = v.(string)
	}
	return includeVolumes
}

func buildCBRVaultResources(d *schema.ResourceData) []vaults.ResourceCreate {
	res := make([]vaults.ResourceCreate, 0)
	vaultType := d.Get("type").(string)
	for _, v := range d.Get("resources").(*schema.Set).List() {
		resourceID := v.(map[string]interface{})["id"].(string)
		resourceType := resourceType[vaultType]
		includes := resourceIncludeVolume(v.(map[string]interface{})["include_volumes"].([]interface{}))
		excludes := resourceExcludeVolume(v.(map[string]interface{})["exclude_volumes"].([]interface{}))
		extra := vaults.ResourceExtraInfo{
			IncludeVolumes: includes,
			ExcludeVolumes: excludes,
		}
		res = append(res, vaults.ResourceCreate{
			ID:        resourceID,
			Type:      resourceType,
			ExtraInfo: &extra,
		})
	}
	return res
}

func buildCBRVaultBilling(d *schema.ResourceData) *vaults.BillingCreate {
	billing := &vaults.BillingCreate{
		ObjectType:      d.Get("type").(string),
		ConsistentLevel: d.Get("consistent_level").(string),
		ProtectType:     d.Get("protection_type").(string),
		Size:            d.Get("size").(int),
	}

	return billing
}

//Take out and save all the IDs nested in the list of resources map as a new string array.
func getResourceIDs(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		resMap := v.(map[string]interface{})
		result[i] = resMap["id"].(string)
	}
	return result
}

//Take out and save all the IDs and resource type nested in the list of resources map as a new ResourceCreate array.
func getResources(slice []interface{}, resType string) []vaults.ResourceCreate {
	result := make([]vaults.ResourceCreate, len(slice))
	for i, v := range slice {
		resMap := v.(map[string]interface{})
		result[i] = vaults.ResourceCreate{
			ID:   resMap["id"].(string),
			Type: resType,
		}
	}
	return result
}

func updateResources(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldRaws, newRaws := d.GetChange("resources")
	oldRawsSet := oldRaws.(*schema.Set)
	newRawsSet := newRaws.(*schema.Set)
	addRaws := newRawsSet.Difference(oldRawsSet)
	removeRaws := oldRawsSet.Difference(oldRawsSet)

	vaultType := d.Get("type").(string)
	if removeRaws.Len() != 0 {
		_, err := vaults.DissociateResources(client, d.Id(), vaults.DissociateResourcesOpts{
			ResourceIDs: getResourceIDs(removeRaws.List()),
		}).Extract()
		if err != nil {
			return fmtp.Errorf("Error unbinding resources: %s", err)
		}
	}

	if addRaws.Len() != 0 {
		_, err := vaults.AssociateResources(client, d.Id(), vaults.AssociateResourcesOpts{
			Resources: getResources(addRaws.List(), resourceType[vaultType]),
		}).Extract()
		if err != nil {
			return fmtp.Errorf("Error binding resources: %s", err)
		}
	}

	return nil
}

func updatePolicy(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldP, newP := d.GetChange("policy_id")
	if newP != "" {
		_, err := vaults.BindPolicy(client, d.Id(), vaults.BindPolicyOpts{
			PolicyID: newP.(string),
		}).Extract()
		if err != nil {
			return fmtp.Errorf("Error binding policy to vault: %s", err)
		}
	} else {
		_, err := vaults.UnbindPolicy(client, d.Id(), vaults.BindPolicyOpts{
			PolicyID: oldP.(string),
		}).Extract()
		if err != nil {
			return fmtp.Errorf("Error unbinding policy from vault: %s", err)
		}
	}
	return nil
}

//Convert Mega Bytes to Giga Bytes, the result is to two decimal places
func getNumberInGB(megaBytes float64) float64 {
	denominator := float64(1024)
	return math.Trunc(float64(megaBytes) / denominator * 1e2 * 1e-2)
}
