package cbr

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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

// @API CBR POST /v3/{project_id}/vaults
// @API CBR POST /v3/{project_id}/vaults/{vault_id}/associatepolicy
// @API CBR POST /v3/{project_id}/vaults/{vault_id}/dissociatepolicy
// @API CBR GET /v3/{project_id}/vaults/{vault_id}
// @API CBR GET /v3/{project_id}/policies
// @API CBR PUT /v3/{project_id}/vaults/{vault_id}
// @API CBR POST /v3/{project_id}/vaults/{vault_id}/addresources
// @API CBR POST /v3/{project_id}/vaults/{vault_id}/removeresources
// @API CBR POST /v3/{project_id}/vault/{vault_id}/tags/action
// @API CBR DELETE /v3/{project_id}/vaults/{vault_id}
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

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the vault is located.",
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					"The cloud type of the vault.",
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the vault.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the vault.",
			},
			"protection_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The protection type of the vault.",
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
			"locked": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Locked status of the vault.",
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

func buildAssociateResourcesForServer(rType string, resources []interface{}) ([]map[string]interface{}, error) {
	// If no resource is set, send an empty slice to the CBR service.
	results := make([]map[string]interface{}, 0, len(resources))

	for _, val := range resources {
		serverId := utils.PathSearch("server_id", val, "").(string)
		if serverId == "" {
			// It means that the current resource object is an empty object: {}
			continue
		}

		if utils.PathSearch("includes", val, schema.NewSet(schema.HashString, nil)).(*schema.Set).Len() > 0 {
			return results, errors.New("server type vaults does not support 'includes'")
		}

		result := map[string]interface{}{
			"id":   utils.PathSearch("server_id", val, ""),
			"type": rType,
		}
		// The server vault only support excludes (blacklist).
		if excludes := utils.PathSearch("excludes", val, schema.NewSet(schema.HashString, nil)).(*schema.Set); excludes.Len() > 0 {
			result["extra_info"] = map[string]interface{}{
				"exclude_volumes": utils.ExpandToStringListBySet(excludes),
			}
		}
		results = append(results, result)
	}
	return results, nil
}

func buildAssociateResourcesForDisk(rType string, resources []interface{}) ([]map[string]interface{}, error) {
	if len(resources) > 1 {
		return nil, errors.New("the size of resources cannot grant than one for disk and turbo vault")
	} else if len(resources) == 0 || (len(resources) == 1 && resources[0] == nil) {
		// If no resource is set, send an empty slice to the CBR service.
		return make([]map[string]interface{}, 0), nil
	}

	if utils.PathSearch("excludes", resources[0], schema.NewSet(schema.HashString, nil)).(*schema.Set).Len() > 0 {
		return nil, errors.New("disk-type and turbo-type vaults does not support 'excludes'")
	}
	includes := utils.PathSearch("includes", resources[0], schema.NewSet(schema.HashString, nil)).(*schema.Set)
	if includes.Len() < 1 {
		return make([]map[string]interface{}, 0), nil
	}
	results := make([]map[string]interface{}, 0, includes.Len())
	for _, v := range includes.List() {
		results = append(results, map[string]interface{}{
			"id":   v,
			"type": rType,
		})
	}
	return results, nil
}

func buildAssociateResources(vType string, resources *schema.Set) ([]map[string]interface{}, error) {
	var result = make([]map[string]interface{}, 0)
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

func isPrePaid(d *schema.ResourceData) bool {
	return d.Get("charging_mode").(string) == "prePaid"
}

func buildBillingStructure(d *schema.ResourceData) map[string]interface{} {
	billing := map[string]interface{}{
		"cloud_type":       utils.ValueIgnoreEmpty(d.Get("cloud_type").(string)),
		"object_type":      d.Get("type").(string),
		"consistent_level": d.Get("consistent_level").(string),
		"protect_type":     d.Get("protection_type").(string),
		"size":             d.Get("size").(int),
		"is_multi_az":      d.Get("is_multi_az").(bool),
	}

	if isPrePaid(d) {
		billing["charging_mode"] = "pre_paid"
		billing["period_type"] = d.Get("period_unit").(string)
		billing["period_num"] = d.Get("period").(int)
		billing["is_auto_renew"], _ = strconv.ParseBool(d.Get("auto_renew").(string))
		billing["is_auto_pay"], _ = strconv.ParseBool(common.GetAutoPay(d))
	}

	return utils.RemoveNil(billing)
}

func buildBindRules(rules map[string]interface{}) []map[string]interface{} {
	// return an empty array instead of null value that used to reset the bind rules
	result := make([]map[string]interface{}, 0, len(rules))

	for k, v := range rules {
		tag := map[string]interface{}{
			"key":   k,
			"value": v,
		}
		result = append(result, tag)
	}

	return result
}

func buildVaultCreateOpts(cfg *config.Config, d *schema.ResourceData) (map[string]interface{}, error) {
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
		return nil, errors.New("the prepaid vault do not support the parameter 'auto_expand'")
	}

	result := map[string]interface{}{
		"name":                  d.Get("name").(string),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		// If no resources are bound when creating, enter an empty list.
		"resources":          resources,
		"backup_policy_id":   utils.ValueIgnoreEmpty(d.Get("policy_id")), // The deprecated parameter (can only bind backup policy).
		"billing":            buildBillingStructure(d),
		"auto_expand":        isAutoExpand.(bool),
		"auto_bind":          d.Get("auto_bind").(bool),
		"backup_name_prefix": utils.ValueIgnoreEmpty(d.Get("backup_name_prefix")),
		"locked":             d.Get("locked").(bool),
	}

	bindRulesRaw, ok := d.Get("bind_rules").(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type of the parameter 'bind_rules', want 'map[string]interface{}', "+
			"but got '%T'", d.Get("bind_rules"))
	}
	if len(bindRulesRaw) > 0 {
		result["bind_rules"] = map[string]interface{}{
			"tags": buildBindRules(bindRulesRaw),
		}
	}

	// Except resources field, remove all keys in which the related values are empty.
	result = utils.RemoveNil(result)
	result["resources"] = resources

	return result, nil
}

func resourceVaultCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/vaults"
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	body, err := buildVaultCreateOpts(cfg, d)
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"vault": body,
		},
	}
	createResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating CBR vault: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if isPrePaid(d) {
		orders := utils.PathSearch("orders[*].orderId", createRespBody, make([]interface{}, 0)).([]interface{})
		if len(orders) < 1 {
			return diag.Errorf("unable to find any order information after creating CBR vault")
		}

		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		orderId := fmt.Sprintf("%v", orders[0])
		if err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("the order is not completed while creating CBR vault: %v", err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	} else {
		d.SetId(utils.PathSearch("vault.id", createRespBody, "").(string))
	}

	vaultId := d.Id()
	// Bind backup(/replication) policy to the vault, not batch bind.
	if policies, ok := d.GetOk("policy"); ok {
		err := updatePoliciesBinding(client, vaultId, schema.NewSet(schema.HashString, nil), policies)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if err := utils.UpdateResourceTags(client, d, "vault", d.Id()); err != nil {
		return diag.Errorf("error setting tags of CBR vault: %s", err)
	}

	return resourceVaultRead(ctx, d, meta)
}

func parseVaultResourcesForServer(resources []interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(resources))
	for _, res := range resources {
		results = append(results, map[string]interface{}{
			"server_id": utils.PathSearch("id", res, ""),
			"includes":  utils.PathSearch("extra_info.include_volumes[*].id", res, make([]interface{}, 0)),
			"excludes":  utils.PathSearch("extra_info.exclude_volumes", res, make([]interface{}, 0)),
		})
	}
	return results
}

func parseVaultResourcesForDisk(resources []interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"includes": utils.PathSearch("[*].id", resources, make([]interface{}, 0)),
		},
	}
}

func flattenVaultResources(vType string, resources []interface{}) []map[string]interface{} {
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

func getPoliciesByVaultId(client *golangsdk.ServiceClient, vaultId string) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/policies?vault_id={vault_id}"
	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{vault_id}", vaultId)

	qeuryOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", queryPath, &qeuryOpts)
	if err != nil {
		return nil, fmt.Errorf("error querying policies from the vault (%s): %s", vaultId, err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("policies", respBody, make([]interface{}, 0)).([]interface{}), nil
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
	results := make([]map[string]interface{}, 0, len(policyList))
	for _, val := range policyList {
		policy := map[string]interface{}{
			"id": utils.PathSearch("id", val, ""),
		}
		if destVaultId := utils.PathSearch(fmt.Sprintf("associated_vaults[?vault_id=='%s'].destination_vault_id|[0]",
			vaultId), val, "").(string); destVaultId != "" {
			policy["destination_vault_id"] = destVaultId
		}
		results = append(results, policy)
	}
	return results
}

func parseVaultChargingMode(chargingMode string) string {
	switch chargingMode {
	case "pre_paid":
		return "prePaid"
	case "post_paid":
		return "postPaid"
	}
	return ""
}

func GetVaultById(client *golangsdk.ServiceClient, vaultId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/vaults/{vault_id}"
	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{vault_id}", vaultId)

	qeuryOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", queryPath, &qeuryOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("vault", respBody, nil), nil
}

func resourceVaultRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		vaultId = d.Id()
	)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	respBody, err := GetVaultById(client, vaultId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying policies from the vault (%s)", vaultId))
	}

	objectType := utils.PathSearch("billing.object_type", respBody, "").(string)
	mErr := multierror.Append(
		// Required && Optional
		d.Set("region", region),
		d.Set("cloud_type", utils.PathSearch("billing.cloud_type", respBody, nil)),
		d.Set("type", objectType),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("protection_type", utils.PathSearch("billing.protect_type", respBody, nil)),
		d.Set("size", utils.PathSearch("billing.size", respBody, nil)),
		d.Set("consistent_level", utils.PathSearch("billing.consistent_level", respBody, nil)),
		d.Set("auto_expand", utils.PathSearch("auto_expand", respBody, nil)),
		d.Set("auto_bind", utils.PathSearch("auto_bind", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("backup_name_prefix", utils.PathSearch("backup_name_prefix", respBody, nil)),
		d.Set("locked", utils.PathSearch("locked", respBody, nil)),
		d.Set("is_multi_az", utils.PathSearch("billing.is_multi_az", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", respBody, nil))),
		d.Set("bind_rules", utils.FlattenTagsToMap(utils.PathSearch("bind_rules.tags", respBody, nil))),
		d.Set("policy", flattenPolicies(client, vaultId)),
		d.Set("resources", flattenVaultResources(objectType,
			utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("charging_mode", parseVaultChargingMode(utils.PathSearch("billing.charging_mode", respBody, "").(string))),
		// Computed
		// The result of 'allocated' and 'used' is in MB, and now we need to use GB as the unit.
		d.Set("allocated", getNumberInGB(utils.PathSearch("billing.allocated", respBody, float64(0)).(float64))),
		d.Set("used", getNumberInGB(utils.PathSearch("billing.used", respBody, float64(0)).(float64))),
		d.Set("spec_code", utils.PathSearch("billing.spec_code", respBody, nil)),
		d.Set("status", utils.PathSearch("billing.status", respBody, nil)),
		d.Set("storage", utils.PathSearch("billing.storage_unit", respBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting vault resource fields: %s", err)
	}

	return nil
}

func buildDissociateResources(vType string, resources *schema.Set) ([]interface{}, error) {
	rType, ok := resourceType[vType]
	if !ok {
		return nil, fmt.Errorf("invalid resource type: %s", vType)
	}
	log.Printf("[DEBUG] The resource type is %s", rType)
	switch rType {
	case ResourceTypeServer, ResourceTypeWorkspace:
		return utils.PathSearch("[*].server_id", resources.List(), make([]interface{}, 0)).([]interface{}), nil
	case ResourceTypeDisk, ResourceTypeTurbo:
		return utils.PathSearch("[*].includes|[0]", resources.List(), schema.NewSet(schema.HashString, nil)).(*schema.Set).List(), nil
	case ResourceTypeNone:
		// Nothing to do.
	default:
		return nil, fmt.Errorf("invalid vault type: %s", vType)
	}
	return nil, nil
}

func waitForAllResourcesDissociated(ctx context.Context, client *golangsdk.ServiceClient, vaultId string,
	resourceIds []interface{}, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetVaultById(client, vaultId)
			if err != nil {
				return nil, "FAILED", fmt.Errorf("error getting vault by ID (%s): %s", vaultId, err)
			}
			if utils.IsSliceContainsAnyAnotherSliceElement(utils.ExpandToStringList(utils.PathSearch("resources[*].id", respBody,
				make([]interface{}, 0)).([]interface{})), utils.ExpandToStringList(resourceIds), false, true) {
				return respBody, "PENDING", nil
			}
			return respBody, "COMPLETED", nil
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
	resources []interface{}, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetVaultById(client, vaultId)
			if err != nil {
				return nil, "FAILED", fmt.Errorf("error getting vault by ID (%s): %s", vaultId, err)
			}
			for _, res := range resources {
				serverId := utils.PathSearch("server_id", res, "").(string)
				if serverId != "" && utils.PathSearch(fmt.Sprintf("length([?id=='%s'])", serverId), respBody, 0).(int) > 0 &&
					!utils.StrSliceContainsAnother(
						utils.ExpandToStringList(utils.PathSearch(fmt.Sprintf("resources[?id=='%s'].extra_info.exclude_volumes|[0]", serverId),
							respBody,
							make([]interface{}, 0),
						).([]interface{})),
						utils.ExpandToStringListBySet(utils.PathSearch("excludes", res, schema.NewSet(schema.HashString, nil)).(*schema.Set))) {
					return respBody, "PENDING", nil
				} else if len(utils.PathSearch(
					fmt.Sprintf("resources[?id=='%s'].extra_info.include_volumes[*].id|[0]", serverId),
					respBody,
					make([]interface{}, 0),
				).([]interface{})) != utils.PathSearch("includes", res, schema.NewSet(schema.HashString, nil)).(*schema.Set).Len() {
					return respBody, "PENDING", nil
				}
			}
			return respBody, "COMPLETED", nil
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
	if delRaws.Len() > 0 && delRaws.List()[0] != nil {
		httpUrl := "v3/{project_id}/vaults/{vault_id}/removeresources"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{vault_id}", vaultId)
		resources, err := buildDissociateResources(vaultType, delRaws)
		if err != nil {
			return fmt.Errorf("error building dissociate list of vault resources: %s", err)
		}
		updateOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"resource_ids": resources,
			},
		}

		_, err = client.Request("POST", updatePath, &updateOpts)
		if err != nil {
			return fmt.Errorf("error updating CBR vault (%s): %s", vaultId, err)
		}
		if waitForAllResourcesDissociated(ctx, client, vaultId, resources, d.Timeout(schema.TimeoutUpdate)) != nil {
			return err
		}
	}

	// Add resources to the specified vault.
	if addRaws.Len() > 0 && addRaws.List()[0] != nil {
		httpUrl := "v3/{project_id}/vaults/{vault_id}/addresources"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{vault_id}", vaultId)
		resources, err := buildAssociateResources(vaultType, addRaws)
		if err != nil {
			return fmt.Errorf("error building associate list of vault resources: %s", err)
		}
		if len(resources) > 0 {
			updateOpts := golangsdk.RequestOpts{
				KeepResponseBody: true,
				JSONBody: map[string]interface{}{
					"resources": resources,
				},
			}

			_, err = client.Request("POST", updatePath, &updateOpts)
			if err != nil {
				return fmt.Errorf("error updating CBR vault (%s): %s", vaultId, err)
			}
			if waitForAllResourcesAssociated(ctx, client, vaultId, addRaws.List(), d.Timeout(schema.TimeoutUpdate)) != nil {
				return err
			}
		}
	}

	return nil
}

func unbindPolicyFromVault(client *golangsdk.ServiceClient, vaultId, policyId string) error {
	httpUrl := "v3/{project_id}/vaults/{vault_id}/dissociatepolicy"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vault_id}", vaultId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"policy_id": policyId,
		},
	}
	_, err := client.Request("POST", createPath, &createOpts)
	return err
}

func bindPolicyToVault(client *golangsdk.ServiceClient, vaultId, destVaultId, policyId string) error {
	httpUrl := "v3/{project_id}/vaults/{vault_id}/associatepolicy"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vault_id}", vaultId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"policy_id":            policyId,
			"destination_vault_id": utils.ValueIgnoreEmpty(destVaultId),
		}),
	}
	_, err := client.Request("POST", createPath, &createOpts)
	return err
}

func updatePoliciesBinding(client *golangsdk.ServiceClient, vaultId string, oPolicies, nPolicies interface{}) error {
	var (
		rmRaw  = oPolicies.(*schema.Set).Difference(nPolicies.(*schema.Set))
		newRaw = nPolicies.(*schema.Set).Difference(oPolicies.(*schema.Set))
	)
	for _, policy := range rmRaw.List() {
		err := unbindPolicyFromVault(client, vaultId, utils.PathSearch("id", policy, "").(string))
		if err != nil {
			return fmt.Errorf("error unbinding policy from vault (%s): %w", vaultId, err)
		}
	}
	for _, policy := range newRaw.List() {
		// Although the BindPolicy method can override the old policy binding, it is difficult for us to know what type
		// of policy is in the old configuration. Overwriting rashly will only cause problems in unbinding.
		err := bindPolicyToVault(client, vaultId, utils.PathSearch("destination_vault_id", policy, "").(string),
			utils.PathSearch("id", policy, "").(string))
		if err != nil {
			return fmt.Errorf("error binding policy to vault (%s): %w", vaultId, err)
		}
	}
	return nil
}

func updateBasicParameters(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		requestBody = make(map[string]interface{})
		billing     = make(map[string]interface{})
		httpUrl     = "v3/{project_id}/vaults/{vault_id}"
		vaultId     = d.Id()
	)
	if d.HasChange("name") {
		requestBody["name"] = d.Get("name").(string)
	}

	if d.HasChange("consistent_level") {
		billing["consistent_level"] = d.Get("consistent_level").(string)
	}

	if d.HasChanges("size", "auto_expand", "auto_bind") {
		if isPrePaid(d) {
			return errors.New("cannot update 'size', 'auto_expand' or 'auto_bind' if the vault is prepaid mode")
		}
		requestBody["auto_expand"] = d.Get("auto_expand").(bool)
		requestBody["auto_bind"] = d.Get("auto_bind").(bool)
		billing["size"] = d.Get("size").(int)
	}

	if d.HasChange("locked") {
		requestBody["locked"] = d.Get("locked").(bool)
	}

	if d.HasChanges("bind_rules") {
		bindRulesRaw, ok := d.Get("bind_rules").(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid type of the parameter 'bind_rules', want 'map[string]interface{}', "+
				"but got '%T'", d.Get("bind_rules"))
		}
		bindRules := map[string]interface{}{
			"tags": buildBindRules(bindRulesRaw),
		}
		requestBody["bind_rules"] = bindRules
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{vault_id}", vaultId)

	if len(billing) > 0 {
		requestBody["billing"] = billing
	}
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"vault": requestBody,
		},
	}
	_, err := client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return fmt.Errorf("error updating CBR vault (%s): %s", vaultId, err)
	}
	return nil
}

func resourceVaultUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		vaultId = d.Id()
	)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	if d.HasChanges("name", "consistent_level", "size", "auto_expand", "auto_bind", "bind_rules", "locked") {
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
		oPolicies, nPolicies := d.GetChange("policy")
		if err := updatePoliciesBinding(client, vaultId, oPolicies, nPolicies); err != nil {
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
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   vaultId,
			ResourceType: "vault",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVaultRead(ctx, d, meta)
}

func resourceVaultDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/vaults/{vault_id}"
		vaultId = d.Id()
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	if isPrePaid(d) {
		err = common.UnsubscribePrePaidResource(d, cfg, []string{vaultId})
		if err != nil {
			return diag.Errorf("error unsubscribing vault (%s): %s", vaultId, err)
		}
	} else {
		deletePath := client.Endpoint + httpUrl
		deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
		deletePath = strings.ReplaceAll(deletePath, "{vault_id}", vaultId)

		deleteOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("DELETE", deletePath, &deleteOpts)
		if err != nil {
			return diag.Errorf("error deleting CBR vault (%s): %s", vaultId, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"available", "deleting"},
		Target:       []string{"deleted"},
		Refresh:      vaultStateRefreshFunc(client, vaultId),
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
		resp, err := GetVaultById(client, vaultId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "RESOURCE_NOT_FOUND", "deleted", nil
			}
			return resp, "available", err
		}
		return resp, utils.PathSearch("billing.status", resp, "STATUS_NOT_FOUND").(string), nil
	}
}
