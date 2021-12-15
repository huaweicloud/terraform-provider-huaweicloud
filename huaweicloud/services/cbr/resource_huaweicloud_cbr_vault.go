package cbr

import (
	"math"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	// VaultTypeServer is the object type of the Cloud Server Backups.
	VaultTypeServer = "server"
	// VaultTypeDisk is the object type of the Cloud Disk Backups.
	VaultTypeDisk = "disk"
	// VaultTypeTurbo is the object type of the SFS Turbo Backups.
	VaultTypeTurbo = "turbo"

	// ResourceTypeServer is the type of the Cloud Server resources to be backed up.
	ResourceTypeServer = "OS::Nova::Server"
	// ResourceTypeDisk is the type of the Cloud Disk resources to be backed up.
	ResourceTypeDisk = "OS::Cinder::Volume"
	// ResourceTypeTurbo is the type of the SFS Turbo resources to be backed up.
	ResourceTypeTurbo = "OS::Sfs::Turbo"
)

var ResourceType map[string]string = map[string]string{
	VaultTypeServer: ResourceTypeServer,
	VaultTypeDisk:   ResourceTypeDisk,
	VaultTypeTurbo:  ResourceTypeTurbo,
}

func ResourceCBRVaultV3() *schema.Resource {
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
					VaultTypeServer, VaultTypeDisk, VaultTypeTurbo,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"excludes": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"includes": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"allocated": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"used": {
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
		},
	}
}

func buildAssociateResourcesForServer(rType string, resources []interface{}) ([]vaults.ResourceCreate, error) {
	if len(resources) == 0 {
		return []vaults.ResourceCreate{}, nil
	}
	results := make([]vaults.ResourceCreate, len(resources))
	for i, val := range resources {
		res := val.(map[string]interface{})
		result := vaults.ResourceCreate{
			ID:        res["server_id"].(string),
			Type:      rType,
			ExtraInfo: &vaults.ResourceExtraInfo{},
		}
		// The server vault only support excludes (blacklist).
		if res["excludes"].(*schema.Set).Len() > 0 {
			volumes := make([]string, res["excludes"].(*schema.Set).Len())
			for i, v := range res["excludes"].(*schema.Set).List() {
				volumes[i] = v.(string)
			}
			result.ExtraInfo.ExcludeVolumes = volumes
		}
		if res["includes"].(*schema.Set).Len() > 0 {
			return results, fmtp.Errorf("Server vault does not support includes.")
		}
		results[i] = result
	}
	return results, nil
}

func buildAssociateResourcesForDisk(rType string, resources []interface{}) ([]vaults.ResourceCreate, error) {
	if len(resources) > 1 {
		return []vaults.ResourceCreate{},
			fmtp.Errorf("The size of resources cannot grant than one for disk and turbo vault.")
	} else if len(resources) == 0 {
		return []vaults.ResourceCreate{}, nil
	}
	res := resources[0].(map[string]interface{})
	if res["includes"].(*schema.Set).Len() > 0 {
		result := make([]vaults.ResourceCreate, res["includes"].(*schema.Set).Len())
		for i, v := range res["includes"].(*schema.Set).List() {
			result[i] = vaults.ResourceCreate{
				ID:   v.(string),
				Type: rType,
			}
		}
		return result, nil
	}
	return []vaults.ResourceCreate{}, fmtp.Errorf("Only includes can be set for disk and turbo vault.")
}

func buildAssociateResources(vType string, resources *schema.Set) ([]vaults.ResourceCreate, error) {
	var result []vaults.ResourceCreate
	var err error
	rType, ok := ResourceType[vType]
	if !ok {
		fmtp.Errorf("Invalid resource type: %s", vType)
	}
	logp.Printf("[DEBUG] The resource type is %s", rType)
	switch rType {
	case ResourceTypeServer:
		result, err = buildAssociateResourcesForServer(rType, resources.List())
	case ResourceTypeDisk, ResourceTypeTurbo:
		result, err = buildAssociateResourcesForDisk(rType, resources.List())
	default:
		err = fmtp.Errorf("The vault type only support server, disk and turbo.")
	}
	return result, err
}

func buildDissociateResourcesForServer(rType string, resources []interface{}) []string {
	result := make([]string, len(resources))
	for i, val := range resources {
		res := val.(map[string]interface{})
		// ID list of all servers attached in the specified vault.
		result[i] = res["server_id"].(string)
	}
	return result
}

func buildDissociateResourcesForDisk(rType string, resources []interface{}) []string {
	rMap := resources[0].(map[string]interface{})

	// All disks attached in the specified vault.
	return utils.ExpandToStringList(rMap["includes"].(*schema.Set).List())
}

func buildDissociateResources(vType string, resources *schema.Set) ([]string, error) {
	rType, ok := ResourceType[vType]
	if !ok {
		return nil, fmtp.Errorf("Invalid resource type: %s", vType)
	}
	logp.Printf("[DEBUG] The resource type is %s", rType)
	switch rType {
	case ResourceTypeServer:
		return buildDissociateResourcesForServer(rType, resources.List()), nil
	case ResourceTypeDisk, ResourceTypeTurbo:
		return buildDissociateResourcesForDisk(rType, resources.List()), nil
	default:
		return nil, fmtp.Errorf("The vault type only support server, disk and turbo.")
	}
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

func resourceCBRVaultV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR v3 client: %s", err)
	}

	resources, err := buildAssociateResources(d.Get("type").(string), d.Get("resources").(*schema.Set))
	if err != nil {
		return fmtp.Errorf("Error building vault resources: %s", err)
	}
	opts := vaults.CreateOpts{
		Name:                d.Get("name").(string),
		AutoExpand:          d.Get("auto_expand").(bool),
		BackupPolicyID:      d.Get("policy_id").(string),
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
		Resources:           resources,
		Billing:             buildCBRVaultBilling(d),
	}

	logp.Printf("[DEBUG] The createOpts is: %+v", opts)
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

	if err := utils.UpdateResourceTags(client, d, "vault", d.Id()); err != nil {
		return fmtp.Errorf("Error setting tags of CBR vault: %s", err)
	}

	return resourceCBRVaultV3Read(d, meta)
}

func makeCbrVaultResourcesForServer(resources []vaults.ResourceResp) []map[string]interface{} {
	results := make([]map[string]interface{}, len(resources))
	for i, res := range resources {
		result := map[string]interface{}{
			"server_id": res.ID,
		}
		if len(res.ExtraInfo.IncludeVolumes) > 0 {
			includeVolumes := make([]string, len(res.ExtraInfo.IncludeVolumes))
			for i, v := range res.ExtraInfo.IncludeVolumes {
				includeVolumes[i] = v.ID
			}
			result["includes"] = res.ExtraInfo.ExcludeVolumes
		}
		if len(res.ExtraInfo.ExcludeVolumes) > 0 {
			result["excludes"] = res.ExtraInfo.ExcludeVolumes
		}

		results[i] = result
	}
	return results
}

// MakeCbrVaultResourcesForDisk is a method for constructing a map list based on the resources response of the server.
func MakeCbrVaultResourcesForDisk(resources []vaults.ResourceResp) []map[string]interface{} {
	includeVolumes := make([]string, len(resources))
	for i, res := range resources {
		includeVolumes[i] = res.ID
	}
	return []map[string]interface{}{
		{
			"includes": includeVolumes,
		},
	}
}

func makeCbrVaultResources(vType string, resources []vaults.ResourceResp) []map[string]interface{} {
	var result []map[string]interface{}
	switch vType {
	case VaultTypeServer:
		result = makeCbrVaultResourcesForServer(resources)
	case VaultTypeDisk, VaultTypeTurbo:
		result = MakeCbrVaultResourcesForDisk(resources)
	}
	return result
}

func setCbrResources(d *schema.ResourceData, vType string, resources []vaults.ResourceResp) error {
	result := makeCbrVaultResources(vType, resources)
	if len(result) != 0 {
		return d.Set("resources", result)
	}
	return nil
}

func getCbrPolicyByVaultId(client *golangsdk.ServiceClient, vaultId string) (string, error) {
	listOpts := policies.ListOpts{
		VaultID: vaultId,
	}
	allPages, err := policies.List(client, listOpts).AllPages()
	if err != nil {
		return "", fmtp.Errorf("Error getting policy by vault ID (%s): %s", vaultId, err)
	}
	policyList, err := policies.ExtractPolicies(allPages)
	if len(policyList) >= 1 {
		return policyList[0].ID, nil
	}
	return "", nil
}

// Convert Mega Bytes to Giga Bytes, the result is to two decimal places.
func getNumberInGB(megaBytes float64) float64 {
	denominator := float64(1024)
	return math.Trunc(float64(megaBytes) / denominator * 1e2 * 1e-2)
}

func setCbrPolicyId(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	policyId, err := getCbrPolicyByVaultId(client, d.Id())
	if err != nil {
		return err
	}
	return d.Set("policy_id", policyId)
}

func resourceCBRVaultV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR v3 client: %s", err)
	}

	resp, err := vaults.Get(client, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting vault details: %s", err)
	}

	mErr := multierror.Append(
		// Required && Optional
		d.Set("name", resp.Name),
		d.Set("consistent_level", resp.Billing.ConsistentLevel),
		d.Set("type", resp.Billing.ObjectType),
		d.Set("protection_type", resp.Billing.ProtectType),
		d.Set("size", resp.Billing.Size),
		d.Set("auto_expand", resp.AutoExpand),
		d.Set("enterprise_project_id", resp.EnterpriseProjectID),
		d.Set("tags", utils.TagsToMap(resp.Tags)),
		setCbrResources(d, resp.Billing.ObjectType, resp.Resources),
		setCbrPolicyId(d, client),
		// Computed
		// The result of 'allocated' and 'used' is in MB, and now we need to use GB as the unit.
		d.Set("allocated", getNumberInGB(float64(resp.Billing.Allocated))),
		d.Set("used", getNumberInGB(float64(resp.Billing.Used))),
		d.Set("spec_code", resp.Billing.SpecCode),
		d.Set("status", resp.Billing.Status),
		d.Set("storage", resp.Billing.StorageUnit),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting vault fields: %s", err)
	}

	return nil
}

func updateResources(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldResources, newResources := d.GetChange("resources")
	addRaws := newResources.(*schema.Set).Difference(oldResources.(*schema.Set))
	delRaws := oldResources.(*schema.Set).Difference(newResources.(*schema.Set))

	// Remove all resources bound to the vault.
	if delRaws.Len() > 0 {
		resources, err := buildDissociateResources(d.Get("type").(string), delRaws)
		if err != nil {
			return fmtp.Errorf("Error building dissociate list of vault resources: %s", err)
		}
		dissociateOpt := vaults.DissociateResourcesOpts{
			ResourceIDs: resources,
		}
		logp.Printf("[DEBUG] The dissociate opt is: %+v", dissociateOpt)
		_, err = vaults.DissociateResources(client, d.Id(), dissociateOpt).Extract()
		if err != nil {
			return fmtp.Errorf("Error dissociating resources: %s", err)
		}
	}

	// Add resources to the specified vault.
	if addRaws.Len() > 0 {
		resources, err := buildAssociateResources(d.Get("type").(string), addRaws)
		if err != nil {
			return fmtp.Errorf("Error building associate list of vault resources: %s", err)
		}
		associateOpt := vaults.AssociateResourcesOpts{
			Resources: resources,
		}
		logp.Printf("[DEBUG] The associate opt is: %+v", associateOpt)
		_, err = vaults.AssociateResources(client, d.Id(), associateOpt).Extract()
		if err != nil {
			return fmtp.Errorf("Error binding resources: %s", err)
		}
	}

	return nil
}

func updatePolicy(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldP, newP := d.GetChange("policy_id")
	if newP != "" {
		// The BindPolicy method can overwrite the old policy binding.
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

func resourceCBRVaultV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
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
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR v3 client: %s", err)
	}

	if err := vaults.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Error deleting CBR v3 vault: %s", err)
	}

	d.SetId("")

	return nil
}
