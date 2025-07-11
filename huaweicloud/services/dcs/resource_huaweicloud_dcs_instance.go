package dcs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/availablezones"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/flavors"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/instances"
	dcsTags "github.com/chnsz/golangsdk/openstack/dcs/v2/tags"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/whitelists"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ctxType string

const (
	chargeModePostPaid = "postPaid"
	chargeModePrePaid  = "prePaid"
)

var (
	chargingMode = map[int]string{
		0: chargeModePostPaid,
		1: chargeModePrePaid,
	}

	operateErrorCode = map[string]bool{
		// current state not support
		"DCS.4026": true,
		// instance status is not running
		"DCS.4049": true,
		// backup
		"DCS.4096": true,
		// restore
		"DCS.4097": true,
		// restart
		"DCS.4111": true,
		// resize
		"DCS.4113": true,
		// change config
		"DCS.4114": true,
		// change password
		"DCS.4115": true,
		// upgrade
		"DCS.4116": true,
		// rollback
		"DCS.4117": true,
		// create
		"DCS.4118": true,
		// freeze
		"DCS.4120": true,
		// creating/restarting
		"DCS.4975": true,
	}
)

// @API DCS GET /v2/available-zones
// @API DCS POST /v2/{project_id}/instances
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS PUT /v2/{project_id}/instance/{instance_id}/whitelist
// @API DCS GET /v2/{project_id}/instance/{instance_id}/whitelist
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/configs
// @API DCS GET /v2/{project_id}/instances/{instance_id}/configs
// @API DCS PUT /v2/{project_id}/instances/status
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/ssl
// @API DCS GET /v2/{project_id}/instances/{instance_id}/ssl
// @API DCS GET /v2/{project_id}/instances/{instance_id}/tags
// @API DCS PUT /v2/{project_id}/instances/{instance_id}
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/password
// @API DCS POST /v2/{project_id}/instances/{instance_id}/resize
// @API DCS POST /v3/{project_id}/instances/{instance_id}/tags/action
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrat
// @API DCS DELETE /v2/{project_id}/instances/{instance_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceDcsInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsInstancesCreate,
		ReadContext:   resourceDcsInstancesRead,
		UpdateContext: resourceDcsInstancesUpdate,
		DeleteContext: resourceDcsInstancesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "schema: Required",
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "schema: Required",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"whitelists"},
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"whitelist_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"whitelists": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 4,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"maintain_begin": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"maintain_end"},
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_policy": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"backup_type", "begin_at", "period_type", "backup_at", "save_days"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"save_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"begin_at": {
							Type:     schema.TypeString,
							Required: true,
						},
						"period_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "weekly",
						},
						"backup_at": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"rename_commands": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
			"tags":          common.TagsSchema(),
			"deleted_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 1,
			},
			"reserved_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"launched_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_info": {
				Type:     schema.TypeList,
				Elem:     bandwidthSchema(),
				Computed: true,
			},
			"cache_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replica_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"readonly_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transparent_client_ip_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"product_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sharding_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Deprecated
			"product_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"product_id", "flavor"},
				Deprecated:   "Deprecated, please use `flavor` instead",
			},
			"available_zones": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"available_zones", "availability_zones"},
				Deprecated:   "Deprecated, please use `availability_zones` instead",
			},
			"enterprise_project_name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "Deprecated, this is a non-public attribute.",
			},
			"internal_version": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated, please us `engine_version` instead.",
			},
			"ip": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated, please us `private_ip` instead.",
			},
			"user_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated",
			},
			"user_name": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated",
			},
			"save_days": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Deprecated, please use `backup_policy` instead",
			},
			"backup_type": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Deprecated, please use `backup_policy` instead",
			},
			"begin_at": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_type", "backup_at", "save_days", "backup_type"},
				Deprecated:   "Deprecated, please use `backup_policy` instead",
			},
			"period_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"begin_at", "backup_at", "save_days", "backup_type"},
				Deprecated:   "Please use `backup_policy` instead",
			},
			"backup_at": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_type", "begin_at", "save_days", "backup_type"},
				Deprecated:   "Deprecated, please use `backup_policy` instead",
				Elem:         &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func bandwidthSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expand_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expand_effect_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expand_interval_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_expand_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"next_expand_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_running": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildBackupPolicyParams(d *schema.ResourceData) *instances.InstanceBackupPolicyOpts {
	if _, ok := d.GetOk("backup_policy"); !ok { // deprecated branch
		if v, ok := d.GetOk("backup_at"); ok {
			backupAts := v.([]interface{})
			return &instances.InstanceBackupPolicyOpts{
				SaveDays:   d.Get("save_days").(int),
				BackupType: d.Get("backup_type").(string),
				PeriodicalBackupPlan: instances.BackupPlan{
					BeginAt:    d.Get("begin_at").(string),
					PeriodType: d.Get("period_type").(string),
					BackupAt:   utils.ExpandToIntList(backupAts),
				},
			}
		}
		// neither backup_policy nor backup_at is specified
		return nil
	}

	backupPolicyList := d.Get("backup_policy").([]interface{})
	if len(backupPolicyList) == 0 {
		return nil
	}
	backupPolicy := backupPolicyList[0].(map[string]interface{})
	backupType := backupPolicy["backup_type"].(string)
	if len(backupType) == 0 || backupType == "manual" {
		return nil
	}
	// build backup policy options
	backupAt := utils.ExpandToIntList(backupPolicy["backup_at"].([]interface{}))
	backupPlan := instances.BackupPlan{
		BackupAt:   backupAt,
		PeriodType: backupPolicy["period_type"].(string),
		BeginAt:    backupPolicy["begin_at"].(string),
	}
	backupPolicyOpts := &instances.InstanceBackupPolicyOpts{
		BackupType:           backupPolicy["backup_type"].(string),
		SaveDays:             backupPolicy["save_days"].(int),
		PeriodicalBackupPlan: backupPlan,
	}
	return backupPolicyOpts
}

func buildDcsTagsParams(tagsMap map[string]interface{}) []dcsTags.ResourceTag {
	tagArr := make([]dcsTags.ResourceTag, 0, len(tagsMap))
	for k, v := range tagsMap {
		tag := dcsTags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		tagArr = append(tagArr, tag)
	}
	return tagArr
}

func buildWhiteListParams(d *schema.ResourceData) whitelists.WhitelistOpts {
	enable := d.Get("whitelist_enable").(bool)
	groupList := d.Get("whitelists").(*schema.Set).List()

	groups := make([]whitelists.WhitelistGroupOpts, len(groupList))
	for i, v := range groupList {
		item := v.(map[string]interface{})
		groups[i] = whitelists.WhitelistGroupOpts{
			GroupName: item["group_name"].(string),
			IPList:    utils.ExpandToStringList(item["ip_address"].([]interface{})),
		}
	}

	whitelistOpts := whitelists.WhitelistOpts{
		Enable: &enable,
		Groups: groups,
	}
	return whitelistOpts
}

func buildSslParam(enable bool) instances.SslOpts {
	sslOpts := instances.SslOpts{
		Enable: &enable,
	}
	return sslOpts
}

func waitForWhiteListCompleted(ctx context.Context, c *golangsdk.ServiceClient, d *schema.ResourceData) error {
	enable := d.Get("whitelist_enable").(bool)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{strconv.FormatBool(!enable)},
		Target:       []string{strconv.FormatBool(enable)},
		Refresh:      refreshForWhiteListEnableStatus(c, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshForWhiteListEnableStatus(c *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := whitelists.Get(c, id).Extract()
		if err != nil {
			return nil, "Error", err
		}
		return r, strconv.FormatBool(r.Enable), nil
	}
}

func resourceDcsInstancesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// azCodes
	var azCodes []string
	availabilityZones, ok := d.GetOk("availability_zones")
	if ok {
		azCodes = utils.ExpandToStringList(availabilityZones.([]interface{}))
	} else {
		azCodes, err = getAvailableZoneCodeByID(client, d.Get("available_zones").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	requestBody := buildCreateInstanceBodyParams(d, azCodes, cfg)
	requestBody["password"] = utils.ValueIgnoreEmpty(d.Get("password"))
	createOpt.JSONBody = utils.RemoveNil(requestBody)

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS instance: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("instances[0].instance_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating DCS instance: ID is not found in API response")
	}
	d.SetId(id)

	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	if orderId != "" {
		err = waitForOrderComplete(ctx, d, cfg, region, orderId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = waitForDcsInstanceRunning(ctx, client, id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	// create whitelist when the function is enabled and configured
	enabled := d.Get("whitelist_enable").(bool)
	if enabled && d.Get("whitelists").(*schema.Set).Len() > 0 {
		whitelistOpts := buildWhiteListParams(d)
		log.Printf("[DEBUG] Create whitelist options: %#v", whitelistOpts)

		err = whitelists.Put(client, id, whitelistOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error creating whitelist for DCS instance (%s): %s", id, err)
		}
		// wait for whitelist created
		err = waitForWhiteListCompleted(ctx, client, d)
		if err != nil {
			return diag.Errorf("Error while waiting to create DCS whitelist: %s", err)
		}
	}

	// set parameters
	if v, ok := d.GetOk("parameters"); ok {
		parameters := v.(*schema.Set).List()
		err = updateParameters(ctx, d.Timeout(schema.TimeoutCreate), client, id, parameters)
		if err != nil {
			return diag.FromErr(err)
		}
		restart, err := checkDcsInstanceRestart(client, id, parameters)
		if err != nil {
			return diag.FromErr(err)
		}
		if restart {
			if err = restartDcsInstance(ctx, d.Timeout(schema.TimeoutCreate), client, id); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if sslEnabled := d.Get("ssl_enable").(bool); sslEnabled {
		sslOpts := buildSslParam(sslEnabled)
		_, err := instances.UpdateSsl(client, id, sslOpts)
		if err != nil {
			return diag.Errorf("error updating SSL for the instance (%s): %s", id, err)
		}

		err = waitForSslCompleted(ctx, client, d)
		if err != nil {
			return diag.Errorf("error waiting for updating SSL to complete: %s", err)
		}
	}

	return resourceDcsInstancesRead(ctx, d, meta)
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, azCodes []string, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                   d.Get("name"),
		"engine":                 d.Get("engine"),
		"engine_version":         utils.ValueIgnoreEmpty(d.Get("engine_version")),
		"capacity":               d.Get("capacity"),
		"instance_num":           1,
		"az_codes":               azCodes,
		"port":                   utils.ValueIgnoreEmpty(d.Get("port")),
		"vpc_id":                 d.Get("vpc_id"),
		"subnet_id":              d.Get("subnet_id"),
		"security_group_id":      utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		"enterprise_project_id":  utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"description":            utils.ValueIgnoreEmpty(d.Get("description")),
		"private_ip":             utils.ValueIgnoreEmpty(d.Get("private_ip")),
		"maintain_begin":         utils.ValueIgnoreEmpty(d.Get("maintain_begin")),
		"maintain_end":           utils.ValueIgnoreEmpty(d.Get("maintain_end")),
		"access_user":            utils.ValueIgnoreEmpty(d.Get("access_user")),
		"template_id":            utils.ValueIgnoreEmpty(d.Get("template_id")),
		"bss_param":              buildCreateInstanceBssParamBodyParams(d),
		"rename_commands":        buildCreateInstanceRenameCommandsBodyParams(d),
		"instance_backup_policy": buildCreateInstanceBackupPolicyBodyParams(d),
		"tags":                   utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}
	// resourceSpecCode
	resourceSpecCode := d.Get("flavor").(string)
	if pid, ok := d.GetOk("product_id"); ok {
		productID := pid.(string)
		resourceSpecCode = productID[0 : len(productID)-2]
	}
	bodyParams["spec_code"] = resourceSpecCode

	// noPasswordAccess
	if d.Get("access_user").(string) == "" && d.Get("password").(string) == "" {
		bodyParams["no_password_access"] = true
	}

	return bodyParams
}

func buildCreateInstanceBssParamBodyParams(d *schema.ResourceData) map[string]interface{} {
	if d.Get("charging_mode") != "prePaid" {
		return nil
	}

	bodyParams := map[string]interface{}{
		"charging_mode": d.Get("charging_mode"),
		"period_type":   d.Get("period_unit"),
		"period_num":    d.Get("period"),
		"is_auto_renew": d.Get("auto_renew"),
	}
	if d.Get("auto_pay").(string) != "false" {
		bodyParams["is_auto_pay"] = true
	}
	return bodyParams
}

func buildCreateInstanceRenameCommandsBodyParams(d *schema.ResourceData) map[string]interface{} {
	renameCmds := d.Get("rename_commands").(map[string]interface{})
	if d.Get("engine") != "Redis" || len(renameCmds) == 0 {
		return nil
	}
	return renameCmds
}

func buildCreateInstanceBackupPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	backupPolicyList, backupPolicyListOk := d.GetOk("backup_policy")
	_, backupAtOk := d.GetOk("backup_at")
	if !backupPolicyListOk && !backupAtOk {
		return nil
	}

	if backupAtOk {
		bodyParams := map[string]interface{}{
			"backup_type": d.Get("backup_type").(string),
			"save_days":   utils.ValueIgnoreEmpty(d.Get("save_days")),
			"periodical_backup_plan": map[string]interface{}{
				"begin_at":    d.Get("begin_at"),
				"period_type": d.Get("period_type"),
				"backup_at":   utils.ExpandToIntList(d.Get("backup_at").([]interface{})),
			},
		}
		return bodyParams
	}

	if len(backupPolicyList.([]interface{})) == 0 {
		return nil
	}

	backupPolicy := backupPolicyList.([]interface{})[0].(map[string]interface{})
	backupType := backupPolicy["backup_type"].(string)
	if len(backupType) == 0 || backupType == "manual" {
		return nil
	}

	bodyParams := map[string]interface{}{
		"backup_type":            backupPolicy["backup_type"],
		"save_days":              backupPolicy["save_days"],
		"periodical_backup_plan": buildCreateInstanceBackupPolicyPlanBodyParams(backupPolicy),
	}
	return bodyParams
}

func buildCreateInstanceBackupPolicyPlanBodyParams(backupPolicy map[string]interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"begin_at":    backupPolicy["begin_at"].(string),
		"period_type": backupPolicy["period_type"].(string),
		"backup_at":   utils.ExpandToIntList(backupPolicy["backup_at"].([]interface{})),
	}
	return bodyParams
}

func updateParameters(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient, instanceID string,
	parameters []interface{}) error {
	parameterOpts := buildUpdateParametersOpt(parameters)
	retryFunc := func() (interface{}, bool, error) {
		log.Printf("[DEBUG] Update DCS instance parameters params: %#v", parameterOpts)
		_, err := instances.ModifyConfiguration(client, instanceID, parameterOpts)
		retry, err := handleOperationError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error modifying parameters for DCS instance (%s): %s", instanceID, err)
	}
	return nil
}

func buildUpdateParametersOpt(parameters []interface{}) instances.ModifyRedisConfigOpts {
	parameterOpts := make([]instances.RedisConfigOpt, 0, len(parameters))
	for _, parameter := range parameters {
		if v, ok := parameter.(map[string]interface{}); ok {
			parameterOpts = append(parameterOpts, instances.RedisConfigOpt{
				ParamId:    v["id"].(string),
				ParamName:  v["name"].(string),
				ParamValue: v["value"].(string),
			})
		}
	}
	return instances.ModifyRedisConfigOpts{RedisConfig: parameterOpts}
}

func checkDcsInstanceRestart(client *golangsdk.ServiceClient, instanceID string, parameters []interface{}) (bool, error) {
	_, needStartParams, err := getParameters(client, instanceID, parameters)
	if err != nil {
		return false, err
	}
	if len(needStartParams) > 0 {
		return true, nil
	}
	return false, nil
}

func getParameters(client *golangsdk.ServiceClient, instanceID string, parameters []interface{}) ([]map[string]interface{},
	[]string, error) {
	configParameters, err := instances.GetConfigurations(client, instanceID)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching the DCS instance parameters (%s): %s", instanceID, err)
	}
	parametersMap := generateParametersMap(configParameters)
	var params []map[string]interface{}
	restartParams := make([]string, 0)
	for _, parameter := range parameters {
		paramId := parameter.(map[string]interface{})["id"]
		if v, ok := parametersMap[paramId.(string)]; ok {
			params = append(params, map[string]interface{}{
				"id":    v.ParamId,
				"name":  v.ParamName,
				"value": v.ParamValue,
			})
			if v.NeedRestart {
				restartParams = append(restartParams, v.ParamName)
			}
		}
	}
	return params, restartParams, nil
}

func restartDcsInstance(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient, instanceID string) error {
	restartOpts := instances.RestartOrFlushInstanceOpts{
		Instances: []string{instanceID},
		Action:    "restart",
	}
	retryFunc := func() (interface{}, bool, error) {
		log.Printf("[DEBUG] Restart DCS instance params: %#v", restartOpts)
		_, err := instances.RestartOrFlushInstance(client, restartOpts)
		retry, err := handleOperationError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error waiting for DCS instance (%s) become running status: %s", instanceID, err)
	}
	return nil
}

func createRenameCommandsOpt(renameCmds map[string]interface{}) instances.RedisCommand {
	renameCommands := instances.RedisCommand{}
	if v, ok := renameCmds["command"]; ok {
		renameCommands.Command = v.(string)
	}
	if v, ok := renameCmds["keys"]; ok {
		renameCommands.Keys = v.(string)
	}
	if v, ok := renameCmds["flushdb"]; ok {
		renameCommands.Flushdb = v.(string)
	}
	if v, ok := renameCmds["flushall"]; ok {
		renameCommands.Flushdb = v.(string)
	}
	if v, ok := renameCmds["hgetall"]; ok {
		renameCommands.Hgetall = v.(string)
	}
	return renameCommands
}

func waitForOrderComplete(ctx context.Context, d *schema.ResourceData, cfg *config.Config, region, orderId string) error {
	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[DEBUG] error the order is not completed while "+
			"creating DCS instance. %s : %v", d.Id(), err)
	}
	_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	return err
}

func waitForDcsInstanceRunning(ctx context.Context, c *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"PENDING"},
		Target:                    []string{"RUNNING"},
		Refresh:                   refreshDcsInstanceState(c, id),
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting instance(%s) to ready: %s", id, err)
	}
	return nil
}

func refreshDcsInstanceState(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getDcsInstanceByID(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "", "DELETED", nil
			}
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status", instance, "").(string)

		failStatus := []string{"CREATEFAILED", "ERROR", "FROZEN"}
		if utils.StrSliceContains(failStatus, status) {
			return instance, "ERROR", fmt.Errorf("unexpect status: %s", status)
		}
		if status == "RUNNING" {
			return instance, status, nil
		}
		return instance, "PENDING", nil
	}
}

func resourceDcsInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instance, err := getDcsInstanceByID(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting DCS instance")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("engine", utils.PathSearch("engine", instance, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", instance, nil)),
		d.Set("capacity", utils.PathSearch("capacity", instance, nil)),
		d.Set("flavor", utils.PathSearch("spec_code", instance, nil)),
		d.Set("availability_zones", utils.PathSearch("az_codes", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", instance, nil)),
		d.Set("vpc_name", utils.PathSearch("vpc_name", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", instance, nil)),
		d.Set("subnet_name", utils.PathSearch("subnet_name", instance, nil)),
		d.Set("subnet_cidr", utils.PathSearch("subnet_cidr", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", instance, nil)),
		d.Set("security_group_name", utils.PathSearch("security_group_name", instance, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", instance, nil)),
		d.Set("description", utils.PathSearch("description", instance, nil)),
		d.Set("private_ip", utils.PathSearch("ip", instance, nil)),
		d.Set("ip", utils.PathSearch("ip", instance, nil)),
		d.Set("maintain_begin", utils.PathSearch("maintain_begin", instance, nil)),
		d.Set("maintain_end", utils.PathSearch("maintain_end", instance, nil)),
		d.Set("charging_mode", chargingMode[int(utils.PathSearch("charging_mode", instance, float64(0)).(float64))]),
		d.Set("port", utils.PathSearch("port", instance, nil)),
		d.Set("status", utils.PathSearch("status", instance, nil)),
		d.Set("used_memory", utils.PathSearch("used_memory", instance, nil)),
		d.Set("max_memory", utils.PathSearch("max_memory", instance, nil)),
		d.Set("domain_name", utils.PathSearch("domain_name", instance, nil)),
		d.Set("user_id", utils.PathSearch("user_id", instance, nil)),
		d.Set("user_name", utils.PathSearch("user_name", instance, nil)),
		d.Set("access_user", utils.PathSearch("access_user", instance, nil)),
		d.Set("ssl_enable", utils.PathSearch("enable_ssl", instance, nil)),
		d.Set("created_at", utils.PathSearch("created_at", instance, nil)),
		d.Set("launched_at", utils.PathSearch("launched_at", instance, nil)),
		d.Set("cache_mode", utils.PathSearch("cache_mode", instance, nil)),
		d.Set("cpu_type", utils.PathSearch("cpu_type", instance, nil)),
		d.Set("readonly_domain_name", utils.PathSearch("readonly_domain_name", instance, nil)),
		d.Set("replica_count", utils.PathSearch("replica_count", instance, nil)),
		d.Set("transparent_client_ip_enable", utils.PathSearch("transparent_client_ip_enable", instance, nil)),
		d.Set("bandwidth_info", flattenInstanceBandWidth(instance)),
		d.Set("product_type", utils.PathSearch("product_type", instance, nil)),
		d.Set("sharding_count", utils.PathSearch("sharding_count", instance, nil)),
		d.Set("backup_policy", flattenInstanceBackupPolicy(instance)),
	)

	// set tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagMap); err != nil {
			return diag.Errorf("[DEBUG] error saving tag to state for DCS instance (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] fetching tags of DCS instance failed: %s", err)
	}

	// set white list
	// some regions (cn-south-1) will fail to call the API due to the cloud reason
	// ignore the error temporarily.
	wList, err := whitelists.Get(client, d.Id()).Extract()
	if err != nil || wList == nil || len(wList.Groups) == 0 {
		log.Printf("error fetching whitelists for DCS instance, error: %s", err)
		// Set to the default value, otherwise it will prompt change after importing resources.
		mErr = multierror.Append(
			mErr,
			d.Set("whitelist_enable", true),
		)
		return diag.FromErr(mErr.ErrorOrNil())
	}

	log.Printf("[DEBUG] Find DCS instance white list : %#v", wList.Groups)
	whiteList := make([]map[string]interface{}, len(wList.Groups))
	for i, group := range wList.Groups {
		whiteList[i] = map[string]interface{}{
			"group_name": group.GroupName,
			"ip_address": group.IPList,
		}
	}
	mErr = multierror.Append(
		mErr,
		d.Set("whitelists", whiteList),
		d.Set("whitelist_enable", wList.Enable),
	)

	diagErr := setDcsInstanceParameters(ctx, d, client, d.Id())
	return append(diagErr, diag.FromErr(mErr.ErrorOrNil())...)
}

func getDcsInstanceByID(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenInstanceBandWidth(instance interface{}) []interface{} {
	bandwidth := utils.PathSearch("bandwidth_info", instance, nil)
	if bandwidth == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"bandwidth": utils.PathSearch("bandwidth", bandwidth, nil),
			"begin_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("begin_time",
				bandwidth, float64(0)).(float64))/1000, false),
			"current_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("current_time",
				bandwidth, float64(0)).(float64))/1000, false),
			"end_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("end_time",
				bandwidth, float64(0)).(float64))/1000, false),
			"expand_count":         utils.PathSearch("expand_count", bandwidth, nil),
			"expand_effect_time":   utils.PathSearch("expand_effect_time", bandwidth, nil),
			"expand_interval_time": utils.PathSearch("expand_interval_time", bandwidth, nil),
			"max_expand_count":     utils.PathSearch("max_expand_count", bandwidth, nil),
			"next_expand_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("next_expand_time",
				bandwidth, float64(0)).(float64))/1000, false),
			"task_running": utils.PathSearch("task_running", bandwidth, nil),
		},
	}
	return rst
}

func flattenInstanceBackupPolicy(instance interface{}) []interface{} {
	backupPolicy := utils.PathSearch("instance_backup_policy.policy", instance, nil)
	if backupPolicy == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"backup_type": utils.PathSearch("backup_type", backupPolicy, nil),
			"save_days":   utils.PathSearch("save_days", backupPolicy, nil),
			"begin_at":    utils.PathSearch("periodical_backup_plan.begin_at", backupPolicy, nil),
			"period_type": utils.PathSearch("periodical_backup_plan.period_type", backupPolicy, nil),
			"backup_at":   utils.PathSearch("periodical_backup_plan.backup_at", backupPolicy, nil),
		},
	}
	return rst
}

func setDcsInstanceParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) diag.Diagnostics {
	params, needStartParams, err := getParameters(client, instanceID, d.Get("parameters").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(params) > 0 {
		if err = d.Set("parameters", params); err != nil {
			return diag.FromErr(err)
		}
		if len(needStartParams) > 0 && ctx.Value(ctxType("parametersChanged")) == "true" {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Parameters Changed",
					Detail:   fmt.Sprintf("Parameters %s changed which needs reboot.", needStartParams),
				},
			}
		}
	}
	return nil
}

func generateParametersMap(configurations *instances.Configuration) map[string]instances.RedisConfig {
	parametersMap := make(map[string]instances.RedisConfig)
	for _, redisConfig := range configurations.RedisConfig {
		parametersMap[redisConfig.ParamId] = redisConfig
	}
	return parametersMap
}

func resourceDcsInstancesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Id()
	client, err := cfg.DcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DCS Client(v2): %s", err)
	}

	// update basic params
	if d.HasChanges("port", "name", "description", "security_group_id", "backup_policy",
		"maintain_begin", "maintain_end", "rename_commands") {
		desc := d.Get("description").(string)
		securityGroupID := d.Get("security_group_id").(string)
		renameCommandsOpt := createRenameCommandsOpt(d.Get("rename_commands").(map[string]interface{}))
		opts := instances.ModifyInstanceOpt{
			Name:            d.Get("name").(string),
			Port:            d.Get("port").(int),
			Description:     &desc,
			MaintainBegin:   d.Get("maintain_begin").(string),
			MaintainEnd:     d.Get("maintain_end").(string),
			SecurityGroupId: &securityGroupID,
			BackupPolicy:    buildBackupPolicyParams(d),
			RenameCommands:  &renameCommandsOpt,
		}
		log.Printf("[DEBUG] Update DCS instance options : %#v", opts)

		_, err = instances.Update(client, instanceId, opts)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.HasChange("port") {
			// Modifying the port is asynchronous and needs to wait for completion.
			err = waitForPortUpdated(ctx, client, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("password") {
		oldVal, newVal := d.GetChange("password")
		opts := instances.UpdatePasswordOpts{
			OldPassword: oldVal.(string),
			NewPassword: newVal.(string),
		}
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err = instances.UpdatePassword(client, instanceId, opts)
			isRetry, err := handleOperationError(err)
			if isRetry {
				return resource.RetryableError(err)
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// resize instance
	err = resizeDcsInstance(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// update tags
	if d.HasChange("tags") {
		oldVal, newVal := d.GetChange("tags")
		err = updateDcsTags(client, instanceId, oldVal.(map[string]interface{}), newVal.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update whitelist
	if d.HasChanges("whitelists", "whitelist_enable") {
		whitelistOpts := buildWhiteListParams(d)
		log.Printf("[DEBUG] Update DCS instance whitelist options: %#v", whitelistOpts)

		err = whitelists.Put(client, instanceId, whitelistOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating whitelist for instance (%s): %s", instanceId, err)
		}

		// wait for whitelist updated
		err = waitForWhiteListCompleted(ctx, client, d)
		if err != nil {
			return diag.Errorf("error while waiting to create DCS whitelist: %s", err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", instanceId, err)
		}
	}

	if d.HasChange("parameters") {
		oRaw, nRaw := d.GetChange("parameters")
		changedParameters := nRaw.(*schema.Set).Difference(oRaw.(*schema.Set)).List()
		err = updateParameters(ctx, d.Timeout(schema.TimeoutUpdate), client, instanceId, changedParameters)
		if err != nil {
			return diag.FromErr(err)
		}
		// Sending parametersChanged to Read to warn users the instance needs a reboot.
		ctx = context.WithValue(ctx, ctxType("parametersChanged"), "true")
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "dcs",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	// update SSL
	if d.HasChange("ssl_enable") {
		sslOpts := buildSslParam(d.Get("ssl_enable").(bool))
		_, err = instances.UpdateSsl(client, instanceId, sslOpts)
		if err != nil {
			return diag.Errorf("error updating SSL for the instance (%s): %s", instanceId, err)
		}

		// wait for SSL updated
		err = waitForSslCompleted(ctx, client, d)
		if err != nil {
			return diag.Errorf("error waiting for updating SSL to complete: %s", err)
		}
	}

	return resourceDcsInstancesRead(ctx, d, meta)
}

func waitForPortUpdated(ctx context.Context, c *golangsdk.ServiceClient, d *schema.ResourceData) error {
	op, np := d.GetChange("port")
	stateConf := &resource.StateChangeConf{
		Pending:      []string{strconv.Itoa(op.(int))},
		Target:       []string{strconv.Itoa(np.(int))},
		Refresh:      refreshDcsInstancePort(c, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("[DEBUG] error while waiting to create/resize/delete DCS instance. %s : %#v",
			d.Id(), err)
	}
	return nil
}

func refreshDcsInstancePort(c *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.Get(c, id)
		if err != nil {
			return nil, "Error", err
		}
		return r, strconv.Itoa(r.Port), nil
	}
}

func updateDcsTags(c *golangsdk.ServiceClient, id string, oldVal, newVal map[string]interface{}) error {
	// remove old tags
	if len(oldVal) > 0 {
		tagList := buildDcsTagsParams(oldVal)
		err := dcsTags.Delete(c, id, tagList)
		if err != nil {
			return err
		}
	}

	// add new tags
	if len(newVal) > 0 {
		tagList := buildDcsTagsParams(newVal)
		err := dcsTags.Create(c, id, tagList)
		if err != nil {
			return err
		}
	}
	return nil
}

func resizeDcsInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DcsV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating DCS Client(v2): %s", err)
	}

	if d.HasChanges("flavor", "capacity") {
		oVal, nVal := d.GetChange("flavor")
		oldSpecCode := oVal.(string)
		newSpecCode := nVal.(string)
		opts, err := buildResizeInstanceOpt(client, d, oldSpecCode, newSpecCode)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Resize DCS instance options : %#v", *opts)

		var res *instances.ResizeResponse
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			res, err = instances.ResizeInstance(client, d.Id(), *opts)
			isRetry, err := handleOperationError(err)
			if isRetry {
				return resource.RetryableError(err)
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error resize DCS instance: %s", err)
		}

		if d.Get("charging_mode").(string) == chargeModePrePaid {
			// wait for order pay
			bssClient, err := cfg.BssV2Client(region)
			if err != nil {
				return fmt.Errorf("error creating BSS v2 client: %s", err)
			}
			err = common.WaitOrderComplete(ctx, bssClient, res.OrderId, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
		}

		// wait for dcs instance change
		err = waitForDcsInstanceRunning(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

		// check the result of the change
		instance, err := instances.Get(client, d.Id())
		if err != nil {
			return common.CheckDeleted(d, err, "DCS instance")
		}
		if instance.SpecCode != d.Get("flavor").(string) {
			return fmt.Errorf("change flavor failed, after changed the DCS flavor still is: %s, expected: %s",
				instance.SpecCode, newSpecCode)
		}
	}
	return nil
}

func buildResizeInstanceOpt(client *golangsdk.ServiceClient, d *schema.ResourceData, oldSpecCode,
	newSpecCode string) (*instances.ResizeInstanceOpts, error) {
	opts := instances.ResizeInstanceOpts{
		SpecCode:    newSpecCode,
		NewCapacity: d.Get("capacity").(float64),
	}
	if d.Get("charging_mode").(string) == chargeModePrePaid {
		opts.BssParam = instances.DcsBssParamOpts{
			IsAutoPay: "true",
		}
	}
	if oldSpecCode == newSpecCode {
		return nil, fmt.Errorf("the param flavor is invalid")
	}
	oldFlavor, err := getFlavorBySpecCode(client, oldSpecCode)
	if err != nil {
		return nil, err
	}
	newFlavor, err := getFlavorBySpecCode(client, newSpecCode)
	if err != nil {
		return nil, err
	}
	changeType := getFlavorChangeType(oldFlavor, newFlavor)
	opts.ChangeType = changeType
	if changeType == "createReplication" {
		azCodes, err := getAzCode(d, client)
		if err != nil {
			return nil, err
		}
		opts.AvailableZones = azCodes
	}
	if changeType == "deleteReplication" {
		if newFlavor.CacheMode == "ha" {
			opts.NodeList = utils.ExpandToStringList(d.Get("deleted_nodes").([]interface{}))
		} else if newFlavor.CacheMode == "cluster" {
			azCodes, err := getAzCode(d, client)
			if err != nil {
				return nil, err
			}
			opts.ReservedIp = utils.ExpandToStringList(d.Get("reserved_ips").([]interface{}))
			opts.AvailableZones = azCodes
		}
	}
	return &opts, nil
}

func getFlavorChangeType(oldFlavor, newFlavor *flavors.Flavor) string {
	// if the cache mode is different, it indicates the type has been changed
	if oldFlavor.CacheMode != newFlavor.CacheMode {
		return "instanceType"
	}
	// indicates the replica count increase, should add replica
	if oldFlavor.ReplicaCount < newFlavor.ReplicaCount {
		return "createReplication"
	}
	// indicates the replica count decrease, should delete replica
	if oldFlavor.ReplicaCount > newFlavor.ReplicaCount {
		return "deleteReplication"
	}
	// indicates only the capacity been changed
	return ""
}

func getFlavorBySpecCode(client *golangsdk.ServiceClient, specCode string) (*flavors.Flavor, error) {
	list, err := flavors.List(client, &flavors.ListOpts{SpecCode: specCode}).Extract()
	if err != nil {
		return nil, fmt.Errorf("error getting dcs flavors list by specCode %s: %s", specCode, err)
	}
	if len(list) < 1 {
		return nil, fmt.Errorf("the result queried by specCode(%s) is empty", specCode)
	}
	return &list[0], nil
}

func handleOperationError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}
		// CBC.99003651: Another operation is being performed.
		if operateErrorCode[errorCode.(string)] || errorCode == "CBC.99003651" {
			return true, err
		}
	}
	return false, err
}

func resourceDcsInstancesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == chargeModePrePaid {
		retryFunc := func() (interface{}, bool, error) {
			err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
			retry, err := handleOperationError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     refreshDcsInstanceState(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error unsubscribe DCS instance: %s", err)
		}
	} else {
		err = deleteDcsInstance(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Waiting to delete success
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING", "PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      refreshDcsInstanceState(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting instance(%s) to delete: %s", d.Id(), err)
	}

	return nil
}

func deleteDcsInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	httpUrl := "v2/{project_id}/instances/{instance_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleOperationError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return err
	}
	return nil
}

func getAzCode(d *schema.ResourceData, client *golangsdk.ServiceClient) ([]string, error) {
	var azCodes []string
	availabilityZones, ok := d.GetOk("availability_zones")
	if ok {
		azCodes = utils.ExpandToStringList(availabilityZones.([]interface{}))
	} else {
		availableZonesCodes, err := getAvailableZoneCodeByID(client, d.Get("available_zones").([]interface{}))
		if err != nil {
			return nil, err
		}
		azCodes = availableZonesCodes
	}
	return azCodes, nil
}

func getAvailableZoneCodeByID(client *golangsdk.ServiceClient, azIds []interface{}) ([]string, error) {
	azCodes := make([]string, 0, len(azIds))
	if len(azIds) == 0 {
		return azCodes, fmt.Errorf("availability_zones are required")
	}

	list, err := availablezones.List(client)
	if err != nil {
		return azCodes, err
	}

	mapping := make(map[string]string)
	for _, v := range list.AvailableZones {
		mapping[v.ID] = v.Code
	}

	for _, id := range azIds {
		azID := id.(string)
		if _, ok := mapping[azID]; !ok {
			return azCodes, fmt.Errorf("invalid available zone code: %s", azID)
		}
		azCodes = append(azCodes, mapping[azID])
	}

	return azCodes, nil
}

func waitForSslCompleted(ctx context.Context, c *golangsdk.ServiceClient, d *schema.ResourceData) error {
	enable := d.Get("ssl_enable").(bool)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{strconv.FormatBool(!enable)},
		Target:       []string{strconv.FormatBool(enable)},
		Refresh:      updateSslStatusRefreshFunc(c, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func updateSslStatusRefreshFunc(c *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.GetSsl(c, id)
		if err != nil {
			return nil, "Error", err
		}
		return r, strconv.FormatBool(r.Enable), nil
	}
}
