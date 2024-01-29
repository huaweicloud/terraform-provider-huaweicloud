package cbr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"
	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

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
	// VaultTypeTurbo is the object type of the Cloud desktop Backups.
	VaultTypeWorkspace = "workspace"
	// VaultTypeTurbo is the object type of the VM Ware Backups.
	VaultTypeVMware = "vmware"
	// VaultTypeTurbo is the object type of the Cloud files Backups.
	VaultTypeFile = "file"

	// ResourceTypeServer is the type of the Cloud Server resources to be backed up.
	ResourceTypeServer = "OS::Nova::Server"
	// ResourceTypeDisk is the type of the Cloud Disk resources to be backed up.
	ResourceTypeDisk = "OS::Cinder::Volume"
	// ResourceTypeTurbo is the type of the SFS Turbo resources to be backed up.
	ResourceTypeTurbo = "OS::Sfs::Turbo"
	// ResourceTypeWorkspace is the type of the Cloud desktop resources to be backed up.
	ResourceTypeWorkspace = "OS::Workspace::DesktopV2"
	// ResourceTypeNone is the type that used to mark no resource needs to be backed up.
	ResourceTypeNone = "No resource to backup"
)

var (
	resourceType = map[string]string{
		VaultTypeServer:    ResourceTypeServer,
		VaultTypeDisk:      ResourceTypeDisk,
		VaultTypeTurbo:     ResourceTypeTurbo,
		VaultTypeWorkspace: ResourceTypeWorkspace,
		VaultTypeVMware:    ResourceTypeNone,
		VaultTypeFile:      ResourceTypeNone,
	}
)

// @API CBR POST /v3/{project_id}/vaults/{id}/associatepolicy
// @API CBR POST /v3/{project_id}/vaults/{id}/dissociatepolicy
// @API CBR POST /v3/{project_id}/vaults/{id}/removeresources
// @API CBR DELETE /v3/{project_id}/vaults/{id}
// @API CBR GET /v3/{project_id}/vaults/{id}
// @API CBR PUT /v3/{project_id}/vaults/{id}
// @API CBR POST /v3/{project_id}/vaults
// @API CBR GET /v3/{project_id}/policies
// @API CBR POST /v3/{project_id}/vault/{id}/tags/action
// @API CBR POST /v3/{project_id}/vaults/{id}/addresources
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the vault is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The vault name.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The vault type.",
			},
			"protection_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The protection type.",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The capacity of the vault, in GB.",
			},
			"consistent_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "crash_consistent",
				Description: "The consistent level (specification) of the vault.",
			},
			"auto_expand": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable auto capacity expansion for the vault.",
			},
			"auto_bind": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether automatic association is supported.",
			},
			"bind_rules": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The rules for automatic association.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The enterprise project ID to which the vault belongs.",
			},
			"policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The policy ID.",
						},
						"destination_vault_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of destination vault to which the replication policy will associated.",
						},
					},
				},
				Description: "The policy details to associate with the CBR vault.",
			},
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the ECS instance to be backed up.",
						},
						"excludes": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The array of disk IDs which will be excluded in the backup.",
						},
						"includes": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The array of disk or SFS file systems which will be included in the backup.",
						},
					},
				},
				Description: "The array of one or more resources to attach to the CBR vault.",
			},
			"backup_name_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The backup name prefix.",
			},
			"is_multi_az": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Whether multiple availability zones are used for backing up.",
			},
			// Public parameters.
			"tags":          common.TagsSchema(),
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
			// Computed parameters.
			"allocated": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The allocated capacity, in GB.",
			},
			"used": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The used capacity, in GB.",
			},
			"spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The specification code.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The vault status.",
			},
			"storage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the bucket for the vault.",
			},
			// Deprecated arguments
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "schema:Deprecated; Using parameter 'policy' instead.",
			},
		},
	}
}

func buildAssociateResourcesForServer(rType string, resources []interface{}) ([]vaults.ResourceCreate, error) {
	// If no resource is set, send an empty slice to the CBR service.
	results := make([]vaults.ResourceCreate, len(resources))

	for i, val := range resources {
		res, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid object type of the parameter 'resources', want 'map[string]interface{}', "+
				"but got '%T'", val)
		}
		if includes, ok := res["includes"].(*schema.Set); ok && includes.Len() > 0 {
			return results, fmt.Errorf("server type vault does not support 'includes'")
		}

		result := vaults.ResourceCreate{
			ID:        res["server_id"].(string),
			Type:      rType,
			ExtraInfo: &vaults.ResourceExtraInfo{},
		}
		// The server vault only support excludes (blacklist).
		if excludes, ok := res["excludes"].(*schema.Set); ok && excludes.Len() > 0 {
			volumes := make([]string, excludes.Len())
			for i, v := range excludes.List() {
				volumes[i] = v.(string)
			}
			result.ExtraInfo.ExcludeVolumes = volumes
		}
		results[i] = result
	}
	return results, nil
}

func buildAssociateResourcesForDisk(rType string, resources []interface{}) ([]vaults.ResourceCreate, error) {
	if len(resources) > 1 {
		return nil, fmt.Errorf("the size of resources cannot grant than one for disk and turbo vault")
	} else if len(resources) == 0 {
		// If no resource is set, send an empty slice to the CBR service.
		return make([]vaults.ResourceCreate, 0), nil
	}

	res, ok := resources[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid object type of the parameter 'resources', want 'map[string]interface{}', "+
			"but got '%T'", resources[0])
	}
	if includes, ok := res["includes"].(*schema.Set); ok && includes.Len() > 0 {
		result := make([]vaults.ResourceCreate, includes.Len())
		for i, v := range includes.List() {
			result[i] = vaults.ResourceCreate{
				ID:   v.(string),
				Type: rType,
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("only includes can be set for disk type and turbo type vault")
}

func buildAssociateResources(vType string, resources *schema.Set) ([]vaults.ResourceCreate, error) {
	var result = make([]vaults.ResourceCreate, 0)
	var err error
	rType, ok := resourceType[vType]
	if !ok {
		return nil, fmt.Errorf("invalid resource type: %s", vType)
	}
	log.Printf("[DEBUG] The resource type is: %s", rType)
	switch rType {
	case ResourceTypeServer, ResourceTypeWorkspace:
		result, err = buildAssociateResourcesForServer(rType, resources.List())
	case ResourceTypeDisk, ResourceTypeTurbo:
		result, err = buildAssociateResourcesForDisk(rType, resources.List())
	case ResourceTypeNone:
		// Nothing to do.
	default:
		err = fmt.Errorf("invalid vault type: %s", vType)
	}
	return result, err
}

func buildDissociateResourcesForServer(resources []interface{}) ([]string, error) {
	result := make([]string, 0, len(resources))
	for _, val := range resources {
		res, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid object type of the parameter 'resources', want 'map[string]interface{}', "+
				"but got '%T'", resources[0])
		}
		// ID list of all servers attached in the specified vault.
		serverId, ok := res["server_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid type of the parameter 'server_id', want 'string', but got '%T'",
				res["server_id"])
		}
		result = append(result, serverId)
	}
	return result, nil
}

func buildDissociateResourcesForDisk(resources []interface{}) ([]string, error) {
	res, ok := resources[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type of the parameter 'resources', want 'map[string]interface{}', but got '%T'",
			resources[0])
	}
	includes, ok := res["includes"].(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("invalid type of the parameter 'resources.includes', want '*schema.Set', but got '%T'",
			res["includes"])
	}

	// All disks attached in the specified vault.
	return utils.ExpandToStringListBySet(includes), nil
}

func buildDissociateResources(vType string, resources *schema.Set) ([]string, error) {
	var result []string
	var err error
	rType, ok := resourceType[vType]
	if !ok {
		return nil, fmt.Errorf("invalid resource type: %s", vType)
	}
	log.Printf("[DEBUG] The resource type is %s", rType)
	switch rType {
	case ResourceTypeServer, ResourceTypeWorkspace:
		result, err = buildDissociateResourcesForServer(resources.List())
	case ResourceTypeDisk, ResourceTypeTurbo:
		return buildDissociateResourcesForDisk(resources.List())
	case ResourceTypeNone:
		// Nothing to do.
	default:
		err = fmt.Errorf("invalid vault type: %s", vType)
	}
	return result, err
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
		IsMultiAz:       d.Get("is_multi_az").(bool),
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

func buildVaultCreateOpts(cfg *config.Config, d *schema.ResourceData) (*vaults.CreateOpts, error) {
	res, ok := d.Get("resources").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("invalid type of the parameter 'resources', want '*schema.Set', but got '%T'",
			d.Get("resources"))
	}
	resources, err := buildAssociateResources(d.Get("type").(string), res)
	if err != nil {
		return nil, fmt.Errorf("error building the structure of associated resources: %s", err)
	}

	isAutoExpand, ok := d.GetOk("auto_expand")
	if ok && isPrePaid(d) {
		return nil, fmt.Errorf("the prepaid vault do not support the parameter 'auto_expand'")
	}

	result := vaults.CreateOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
		Resources:           resources,
		BackupPolicyID:      d.Get("policy_id").(string), // The deprecated parameter (can only bind backup policy).
		Billing:             buildBillingStructure(d),
		AutoExpand:          isAutoExpand.(bool),
		AutoBind:            d.Get("auto_bind").(bool),
		BackupNamePrefix:    d.Get("backup_name_prefix").(string),
	}

	bindRulesRaw, ok := d.Get("bind_rules").(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type of the parameter 'bind_rules', want 'map[string]interface{}', "+
			"but got '%T'", d.Get("bind_rules"))
	}
	binRulesList := utils.ExpandResourceTags(bindRulesRaw)
	if len(binRulesList) > 0 {
		bindRules := &vaults.VaultBindRules{
			Tags: binRulesList,
		}
		result.BindRules = bindRules
	}
	return &result, nil
}

func resourceVaultCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	opts, err := buildVaultCreateOpts(cfg, d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] The createOpts is: %+v", opts)
	result := vaults.Create(client, opts)

	if isPrePaid(d) {
		resp, err := result.ExtractOrder()
		if err != nil || len(resp.Orders) < 1 || resp.Orders[0].ID == "" {
			return diag.Errorf("unable to find any order information after creating CBR vault")
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		orderId := resp.Orders[0].ID
		if err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("the order is not completed while creating CBR vault: %v", err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	} else {
		vault, err := result.Extract()
		if err != nil {
			return diag.Errorf("error creating vaults: %s", err)
		}
		d.SetId(vault.ID)
	}

	// Bind backup(/replication) policy to the vault, not batch bind.
	if _, ok := d.GetOk("policy"); ok {
		err := updatePolicyBindings(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if err := utils.UpdateResourceTags(client, d, "vault", d.Id()); err != nil {
		return diag.Errorf("error setting tags of CBR vault: %s", err)
	}

	return resourceVaultRead(ctx, d, meta)
}

func parseVaultResourcesForServer(resources []vaults.ResourceResp) []map[string]interface{} {
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
			result["includes"] = includeVolumes
		}
		if len(res.ExtraInfo.ExcludeVolumes) > 0 {
			result["excludes"] = res.ExtraInfo.ExcludeVolumes
		}

		results[i] = result
	}
	return results
}

func parseVaultResourcesForDisk(resources []vaults.ResourceResp) []map[string]interface{} {
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

func flattenVaultResources(vType string, resources []vaults.ResourceResp) []map[string]interface{} {
	switch vType {
	case VaultTypeServer, VaultTypeWorkspace:
		return parseVaultResourcesForServer(resources)
	case VaultTypeDisk, VaultTypeTurbo:
		return parseVaultResourcesForDisk(resources)
	default:
		// Nothing to do for type file and type vmware.
	}
	return nil
}

func getPoliciesByVaultId(client *golangsdk.ServiceClient, vaultId string) ([]policies.Policy, error) {
	listOpts := policies.ListOpts{
		VaultID: vaultId,
	}
	allPages, err := policies.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error getting policy by vault ID (%s): %s", vaultId, err)
	}

	policyList, err := policies.ExtractPolicies(allPages)
	if err != nil {
		return nil, fmt.Errorf("error extracting policy list: %s", err)
	}

	return policyList, nil
}

// Convert Mega Bytes to Giga Bytes, the result is to two decimal places.
func getNumberInGB(megaBytes float64) float64 {
	denominator := float64(1024)
	result, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", megaBytes/denominator), 64)
	return result
}

func flattenPolicies(client *golangsdk.ServiceClient, vaultId string) []map[string]interface{} {
	policyList, err := getPoliciesByVaultId(client, vaultId)
	if err != nil {
		log.Printf("[ERROR] error querying CBR policies by vault ID (%s): %v", vaultId, err)
		return nil
	}
	if len(policyList) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, len(policyList))
	for i, val := range policyList {
		policy := map[string]interface{}{
			"id": val.ID,
		}
		if len(val.AssociatedVaults) < 1 {
			continue
		}
		for _, v := range val.AssociatedVaults {
			if v.VaultID == vaultId {
				policy["destination_vault_id"] = v.DestinationVaultID
				break
			}
		}

		result[i] = policy
	}
	return result
}

func parseVaultChargingMode(billing vaults.Billing) string {
	switch billing.ChargingMode {
	case "pre_paid":
		return "prePaid"
	case "post_paid":
		return "postPaid"
	}
	return ""
}

func resourceVaultRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	vaultId := d.Id()
	resp, err := vaults.Get(client, vaultId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CBR vault")
	}

	mErr := multierror.Append(
		// Required && Optional
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("type", resp.Billing.ObjectType),
		d.Set("protection_type", resp.Billing.ProtectType),
		d.Set("size", resp.Billing.Size),
		d.Set("consistent_level", resp.Billing.ConsistentLevel),
		d.Set("auto_expand", resp.AutoExpand),
		d.Set("auto_bind", resp.AutoBind),
		d.Set("enterprise_project_id", resp.EnterpriseProjectID),
		d.Set("backup_name_prefix", resp.BackupNamePrefix),
		d.Set("is_multi_az", resp.Billing.IsMultiAz),
		d.Set("tags", utils.TagsToMap(resp.Tags)),
		d.Set("bind_rules", utils.TagsToMap(resp.BindRules.Tags)),
		d.Set("policy", flattenPolicies(client, vaultId)),
		d.Set("resources", flattenVaultResources(resp.Billing.ObjectType, resp.Resources)),
		d.Set("charging_mode", parseVaultChargingMode(resp.Billing)),
		// Computed
		// The result of 'allocated' and 'used' is in MB, and now we need to use GB as the unit.
		d.Set("allocated", getNumberInGB(float64(resp.Billing.Allocated))),
		d.Set("used", getNumberInGB(float64(resp.Billing.Used))),
		d.Set("spec_code", resp.Billing.SpecCode),
		d.Set("status", resp.Billing.Status),
		d.Set("storage", resp.Billing.StorageUnit),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting vault resource fields: %s", err)
	}

	return nil
}

func waitForAllResourcesDissociated(ctx context.Context, client *golangsdk.ServiceClient, vaultId string,
	resourceIds []string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := vaults.Get(client, vaultId).Extract()
			if err != nil {
				return nil, "FAILED", fmt.Errorf("error getting vault by ID (%s)", vaultId)
			}
			for _, queryRes := range resp.Resources {
				for _, resId := range resourceIds {
					if queryRes.ID == resId {
						return resp, "PENDING", nil
					}
				}
			}
			return resp, "COMPLETED", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("timeout waiting for dissociate resources to complete: %s", err)
	}
	return nil
}

func waitForAllResourcesAssociated(ctx context.Context, client *golangsdk.ServiceClient, vaultId string,
	resources []vaults.ResourceCreate, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := vaults.Get(client, vaultId).Extract()
			if err != nil {
				return nil, "FAILED", fmt.Errorf("error getting vault by ID (%s)", vaultId)
			}
			for _, res := range resources {
				for _, queryRes := range resp.Resources {
					if res.ID != "" {
						if queryRes.ID == res.ID || len(queryRes.ExtraInfo.ExcludeVolumes) == len(res.ExtraInfo.ExcludeVolumes) {
							return resp, "COMPLETED", nil
						}
					} else if len(queryRes.ExtraInfo.IncludeVolumes) == len(res.ExtraInfo.IncludeVolumes) {
						return resp, "COMPLETED", nil
					}
				}
			}
			return resp, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("timeout waiting for associate resources to complete: %s", err)
	}
	return nil
}

func updateAssociatedResources(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		vaultId   = d.Id()
		vaultType = d.Get("type").(string)

		oldResources, newResources = d.GetChange("resources")
		addRaws                    = newResources.(*schema.Set).Difference(oldResources.(*schema.Set))
		delRaws                    = oldResources.(*schema.Set).Difference(newResources.(*schema.Set))
	)

	// Remove all resources bound to the vault.
	if delRaws.Len() > 0 {
		resources, err := buildDissociateResources(vaultType, delRaws)
		if err != nil {
			return fmt.Errorf("error building dissociate list of vault resources: %s", err)
		}
		opts := vaults.DissociateResourcesOpts{
			ResourceIDs: resources,
		}
		log.Printf("[DEBUG] The dissociate opts is: %#v", opts)
		_, err = vaults.DissociateResources(client, vaultId, opts).Extract()
		if err != nil {
			return fmt.Errorf("error dissociating resources: %s", err)
		}
		if waitForAllResourcesDissociated(ctx, client, vaultId, resources, d.Timeout(schema.TimeoutUpdate)) != nil {
			return err
		}
	}

	// Add resources to the specified vault.
	if addRaws.Len() > 0 {
		resources, err := buildAssociateResources(vaultType, addRaws)
		if err != nil {
			return fmt.Errorf("error building associate list of vault resources: %s", err)
		}
		opts := vaults.AssociateResourcesOpts{
			Resources: resources,
		}
		log.Printf("[DEBUG] The associate opts is: %#v", opts)
		_, err = vaults.AssociateResources(client, vaultId, opts).Extract()
		if err != nil {
			return fmt.Errorf("error associating resources: %s", err)
		}
		if waitForAllResourcesAssociated(ctx, client, vaultId, resources, d.Timeout(schema.TimeoutUpdate)) != nil {
			return err
		}
	}

	return nil
}

func updatePolicyBindings(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		vaultId        = d.Id()
		oldVal, newVal = d.GetChange("policy")
		rmRaw          = oldVal.(*schema.Set).Difference(newVal.(*schema.Set))
		newRaw         = newVal.(*schema.Set).Difference(oldVal.(*schema.Set))
	)
	for _, policy := range rmRaw.List() {
		pm := policy.(map[string]interface{})
		_, err := vaults.UnbindPolicy(client, vaultId, vaults.BindPolicyOpts{
			PolicyID: pm["id"].(string),
		}).Extract()
		if err != nil {
			return fmt.Errorf("error unbinding policy from vault (%s): %w", vaultId, err)
		}
	}
	for _, policy := range newRaw.List() {
		pm := policy.(map[string]interface{})
		// Although the BindPolicy method can override the old policy binding, it is difficult for us to know what type
		// of policy is in the old configuration. Overwriting rashly will only cause problems in unbinding.
		_, err := vaults.BindPolicy(client, vaultId, vaults.BindPolicyOpts{
			DestinationVaultId: pm["destination_vault_id"].(string),
			PolicyID:           pm["id"].(string),
		}).Extract()
		if err != nil {
			return fmt.Errorf("error binding policy to vault (%s): %w", vaultId, err)
		}
	}
	return nil
}

func updateBasicParameters(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		billing = vaults.BillingUpdate{}
		opts    = vaults.UpdateOpts{
			Billing: &billing,
		}
	)
	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}
	if d.HasChange("consistent_level") {
		billing.ConsistentLevel = d.Get("consistent_level").(string)
	}

	if d.HasChanges("size", "auto_expand", "auto_bind") {
		if isPrePaid(d) {
			return fmt.Errorf("cannot update 'size', 'auto_expand' or 'auto_bind' if the vault is prepaid mode")
		}
		opts.AutoExpand = utils.Bool(d.Get("auto_expand").(bool))
		opts.AutoBind = utils.Bool(d.Get("auto_bind").(bool))
		billing.Size = d.Get("size").(int)
	}

	if d.HasChanges("bind_rules") {
		bindRulesRaw, ok := d.Get("bind_rules").(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid type of the parameter 'bind_rules', want 'map[string]interface{}', "+
				"but got '%T'", d.Get("bind_rules"))
		}
		binRulesList := utils.ExpandResourceTags(bindRulesRaw)
		bindRules := &vaults.VaultBindRules{
			Tags: binRulesList,
		}
		opts.BindRules = bindRules
	}
	_, err := vaults.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmt.Errorf("error updating the vault: %s", err)
	}
	return nil
}

func resourceVaultUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	vaultId := d.Id()
	if d.HasChanges("name", "consistent_level", "size", "auto_expand", "auto_bind", "bind_rules") {
		if err = updateBasicParameters(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("resources") {
		if err := updateAssociatedResources(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("policy") {
		if err := updatePolicyBindings(d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err = utils.UpdateResourceTags(client, d, "vault", vaultId); err != nil {
			return diag.Errorf("failed to update tags: %s", err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), vaultId); err != nil {
			return diag.Errorf("error updating the auto-renew of the vault (%s): %s", vaultId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := enterpriseprojects.MigrateResourceOpts{
			ResourceId:   vaultId,
			ResourceType: "vault",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := common.MigrateEnterpriseProject(ctx, cfg, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVaultRead(ctx, d, meta)
}

func resourceVaultDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	if isPrePaid(d) {
		err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
		if err != nil {
			return diag.Errorf("error unsubscribing vault (%s): %s", d.Id(), err)
		}
	} else {
		if err := vaults.Delete(client, d.Id()).ExtractErr(); err != nil {
			return diag.Errorf("error deleting CBR v3 vault: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"available", "deleting"},
		Target:       []string{"deleted"},
		Refresh:      vaultStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 20 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("timeout waiting for vault deletion to complete: %s", err)
	}

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
