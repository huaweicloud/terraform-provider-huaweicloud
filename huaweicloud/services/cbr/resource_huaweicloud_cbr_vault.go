package cbr

import (
	"context"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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

func ResourceVault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVaultCreate,
		ReadContext:   resourceVaultRead,
		UpdateContext: resourceVaultUpdate,
		DeleteContext: resourceVaultDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
			"consistent_level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "crash_consistent",
				ValidateFunc: validation.StringInSlice([]string{
					"crash_consistent", "app_consistent",
				}, false),
			},
			"auto_expand": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_bind": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"bind_rules": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"tags":          common.TagsSchema(),
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenew(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
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
			return results, fmt.Errorf("server vault does not support includes")
		}
		results[i] = result
	}
	return results, nil
}

func buildAssociateResourcesForDisk(rType string, resources []interface{}) ([]vaults.ResourceCreate, error) {
	if len(resources) > 1 {
		return []vaults.ResourceCreate{},
			fmt.Errorf("the size of resources cannot grant than one for disk and turbo vault")
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
	return []vaults.ResourceCreate{}, fmt.Errorf("only includes can be set for disk and turbo vault")
}

func buildAssociateResources(vType string, resources *schema.Set) ([]vaults.ResourceCreate, error) {
	var result []vaults.ResourceCreate
	var err error
	rType, ok := ResourceType[vType]
	if !ok {
		return nil, fmt.Errorf("invalid resource type: %s", vType)
	}
	log.Printf("[DEBUG] The resource type is %s", rType)
	switch rType {
	case ResourceTypeServer:
		result, err = buildAssociateResourcesForServer(rType, resources.List())
	case ResourceTypeDisk, ResourceTypeTurbo:
		result, err = buildAssociateResourcesForDisk(rType, resources.List())
	default:
		err = fmt.Errorf("the vault type only support server, disk and turbo")
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
		return nil, fmt.Errorf("invalid resource type: %s", vType)
	}
	log.Printf("[DEBUG] The resource type is %s", rType)
	switch rType {
	case ResourceTypeServer:
		return buildDissociateResourcesForServer(rType, resources.List()), nil
	case ResourceTypeDisk, ResourceTypeTurbo:
		return buildDissociateResourcesForDisk(rType, resources.List()), nil
	default:
		return nil, fmt.Errorf("the vault type only support server, disk and turbo")
	}
}

func isPrePaid(d *schema.ResourceData) bool {
	return d.Get("charging_mode").(string) == "prePaid"
}

func buildBillingStructure(d *schema.ResourceData) *vaults.BillingCreate {
	billing := &vaults.BillingCreate{
		ObjectType:      d.Get("type").(string),
		ConsistentLevel: d.Get("consistent_level").(string),
		ProtectType:     d.Get("protection_type").(string),
		Size:            d.Get("size").(int),
	}

	if isPrePaid(d) {
		billing.ChargingMode = "pre_paid"
		billing.PeriodType = d.Get("period_unit").(string)
		billing.PeriodNum = d.Get("period").(int)
		billing.IsAutoRenew, _ = strconv.ParseBool(d.Get("auto_renew").(string))
		billing.IsAutoPay, _ = strconv.ParseBool(common.GetAutoPay(d))
	}

	return billing
}

func resourceVaultCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	resources, err := buildAssociateResources(d.Get("type").(string), d.Get("resources").(*schema.Set))
	if err != nil {
		return diag.Errorf("error building vault resources: %s", err)
	}

	ae, ok := d.GetOk("auto_expand")
	if ok && isPrePaid(d) {
		return diag.Errorf("the prepaid vault do not support the auto_expand parameter")
	}

	opts := vaults.CreateOpts{
		Name:                d.Get("name").(string),
		BackupPolicyID:      d.Get("policy_id").(string),
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
		Resources:           resources,
		Billing:             buildBillingStructure(d),
		AutoExpand:          ae.(bool),
		AutoBind:            d.Get("auto_bind").(bool),
	}

	bindRulesRaw := d.Get("bind_rules").(map[string]interface{})
	binRulesList := utils.ExpandResourceTags(bindRulesRaw)
	if len(binRulesList) > 0 {
		bindRules := &vaults.VaultBindRules{
			Tags: binRulesList,
		}
		opts.BindRules = bindRules
	}

	log.Printf("[DEBUG] The createOpts is: %+v", opts)
	result := vaults.Create(client, opts)

	if isPrePaid(d) {
		resp, err := result.ExtractOrder()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(resp.Orders) < 1 {
			return diag.Errorf("unable to find any order information after creating CBR vault")
		}
		if resp.Orders[0].ResourceId == "" {
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "Unsupported Region",
					Detail: fmt.Sprintf("Currently, we does not support prepaid creation completely in this region "+
						"(%s), because of the API response does not include vault ID. But the order has been created, "+
						"if you don't want it, you can unsubscribe in the console. Also you can manage it by import "+
						"operation using ID. You cannot create a new vault with the same configuration until you "+
						"unsubscribe.", region),
				},
			}
		}
		d.SetId(resp.Orders[0].ResourceId)
		err = common.WaitOrderComplete(ctx, d, config, resp.Orders[0].ID)
		if err != nil {
			return diag.Errorf("the order is not completed while creating CBR vault (%s): %v", d.Id(), err)
		}
	} else {
		vault, err := result.Extract()
		if err != nil {
			return diag.Errorf("error creating vaults: %s", err)
		}
		d.SetId(vault.ID)
	}

	if policy, ok := d.GetOk("policy_id"); ok {
		_, err := vaults.BindPolicy(client, d.Id(), vaults.BindPolicyOpts{PolicyID: policy.(string)}).Extract()
		if err != nil {
			return diag.Errorf("error binding policy to vault: %s", err)
		}
	}

	if err := utils.UpdateResourceTags(client, d, "vault", d.Id()); err != nil {
		return diag.Errorf("error setting tags of CBR vault: %s", err)
	}

	return resourceVaultRead(ctx, d, meta)
}

func makeVaultResourcesForServer(resources []vaults.ResourceResp) []map[string]interface{} {
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

// MakeVaultResourcesForDisk is a method for constructing a map list based on the resources response of the server.
func MakeVaultResourcesForDisk(resources []vaults.ResourceResp) []map[string]interface{} {
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

func makeVaultResources(vType string, resources []vaults.ResourceResp) []map[string]interface{} {
	var result []map[string]interface{}
	switch vType {
	case VaultTypeServer:
		result = makeVaultResourcesForServer(resources)
	case VaultTypeDisk, VaultTypeTurbo:
		result = MakeVaultResourcesForDisk(resources)
	}
	return result
}

func setResources(d *schema.ResourceData, vType string, resources []vaults.ResourceResp) error {
	result := makeVaultResources(vType, resources)
	if len(result) != 0 {
		return d.Set("resources", result)
	}
	return nil
}

func getPolicyByVaultId(client *golangsdk.ServiceClient, vaultId string) (string, error) {
	listOpts := policies.ListOpts{
		VaultID: vaultId,
	}
	allPages, err := policies.List(client, listOpts).AllPages()
	if err != nil {
		return "", fmt.Errorf("error getting policy by vault ID (%s): %s", vaultId, err)
	}

	policyList, err := policies.ExtractPolicies(allPages)
	if err != nil {
		return "", fmt.Errorf("error extracting vault list: %s", err)
	}

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

func setPolicyId(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	policyId, err := getPolicyByVaultId(client, d.Id())
	if err != nil {
		return err
	}
	return d.Set("policy_id", policyId)
}

func setCbrVaultCharging(d *schema.ResourceData, billing vaults.Billing) error {
	switch billing.ChargingMode {
	case "pre_paid":
		return d.Set("charging_mode", "prePaid")
	case "post_paid":
		return d.Set("charging_mode", "postPaid")
	}
	return nil
}

func resourceVaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	resp, err := vaults.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting vault details")
	}

	mErr := multierror.Append(
		// Required && Optional
		d.Set("name", resp.Name),
		d.Set("type", resp.Billing.ObjectType),
		d.Set("protection_type", resp.Billing.ProtectType),
		d.Set("size", resp.Billing.Size),
		d.Set("consistent_level", resp.Billing.ConsistentLevel),
		d.Set("auto_expand", resp.AutoExpand),
		d.Set("auto_bind", resp.AutoBind),
		d.Set("enterprise_project_id", resp.EnterpriseProjectID),
		d.Set("tags", utils.TagsToMap(resp.Tags)),
		d.Set("bind_rules", utils.TagsToMap(resp.BindRules.Tags)),
		setResources(d, resp.Billing.ObjectType, resp.Resources),
		setPolicyId(d, client),
		setCbrVaultCharging(d, resp.Billing),
		// Computed
		// The result of 'allocated' and 'used' is in MB, and now we need to use GB as the unit.
		d.Set("allocated", getNumberInGB(float64(resp.Billing.Allocated))),
		d.Set("used", getNumberInGB(float64(resp.Billing.Used))),
		d.Set("spec_code", resp.Billing.SpecCode),
		d.Set("status", resp.Billing.Status),
		d.Set("storage", resp.Billing.StorageUnit),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting vault fields: %s", err)
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
			return fmt.Errorf("error building dissociate list of vault resources: %s", err)
		}
		dissociateOpt := vaults.DissociateResourcesOpts{
			ResourceIDs: resources,
		}
		log.Printf("[DEBUG] The dissociate opt is: %+v", dissociateOpt)
		_, err = vaults.DissociateResources(client, d.Id(), dissociateOpt).Extract()
		if err != nil {
			return fmt.Errorf("error dissociating resources: %s", err)
		}
	}

	// Add resources to the specified vault.
	if addRaws.Len() > 0 {
		resources, err := buildAssociateResources(d.Get("type").(string), addRaws)
		if err != nil {
			return fmt.Errorf("error building associate list of vault resources: %s", err)
		}
		associateOpt := vaults.AssociateResourcesOpts{
			Resources: resources,
		}
		log.Printf("[DEBUG] The associate opt is: %+v", associateOpt)
		_, err = vaults.AssociateResources(client, d.Id(), associateOpt).Extract()
		if err != nil {
			return fmt.Errorf("error binding resources: %s", err)
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
			return fmt.Errorf("error binding policy to vault: %s", err)
		}
	} else {
		_, err := vaults.UnbindPolicy(client, d.Id(), vaults.BindPolicyOpts{
			PolicyID: oldP.(string),
		}).Extract()
		if err != nil {
			return fmt.Errorf("error unbinding policy from vault: %s", err)
		}
	}
	return nil
}

func resourceVaultUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	opts := vaults.UpdateOpts{}
	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}

	if d.HasChanges("size", "auto_expand", "auto_bind") {
		if isPrePaid(d) {
			return diag.Errorf("cannot update size or auto_expand if the vault is prepaid mode")
		}
		ae := d.Get("auto_expand").(bool)
		ab := d.Get("auto_bind").(bool)
		opts.AutoExpand = &ae
		opts.AutoBind = &ab
		opts.Billing = &vaults.BillingUpdate{
			Size: d.Get("size").(int),
		}
	}

	if d.HasChanges("bind_rules") {
		bindRulesRaw := d.Get("bind_rules").(map[string]interface{})
		binRulesList := utils.ExpandResourceTags(bindRulesRaw)
		bindRules := &vaults.VaultBindRules{
			Tags: binRulesList,
		}
		opts.BindRules = bindRules
	}

	if !reflect.DeepEqual(opts, vaults.UpdateOpts{}) {
		_, err := vaults.Update(client, d.Id(), opts).Extract()
		if err != nil {
			return diag.Errorf("error updating the vault: %s", err)
		}
	}

	if d.HasChange("resources") {
		if err := updateResources(d, client); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("policy_id") {
		if err := updatePolicy(d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err = utils.UpdateResourceTags(client, d, "vault", d.Id()); err != nil {
			return diag.Errorf("failed to update tags: %s", err)
		}
	}

	return resourceVaultRead(ctx, d, meta)
}

func resourceVaultDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	if isPrePaid(d) {
		err = common.UnsubscribePrePaidResource(d, config, []string{d.Id()})
		if err != nil {
			return diag.Errorf("error unsubscribing vault (%s): %s", d.Id(), err)
		}
	} else {
		if err := vaults.Delete(client, d.Id()).ExtractErr(); err != nil {
			return diag.Errorf("error deleting CBR v3 vault: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    vaultStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("timeout waiting for vault deletion to complete: %s", err)
	}
	d.SetId("")

	return nil
}

func vaultStateRefreshFunc(client *golangsdk.ServiceClient, vaultId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := vaults.Get(client, vaultId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resp, "deleted", nil
			}
			return resp, "available", err
		}
		return resp, resp.Billing.Status, nil
	}
}
