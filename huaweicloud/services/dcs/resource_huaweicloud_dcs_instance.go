package dcs

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
)

// @API DCS GET /v2/available-zones
// @API DCS POST /v2/{project_id}/instances
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS PUT /v2/{project_id}/instance/{instance_id}/whitelist
// @API DCS GET /v2/{project_id}/instance/{instance_id}/whitelist
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/async-configs
// @API DCS GET /v2/{project_id}/jobs/{job_id}
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/bigkey/autoscan
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/hotkey/autoscan
// @API DCS GET /v2/{project_id}/instances/{instance_id}/bigkey/autoscan
// @API DCS GET /v2/{project_id}/instances/{instance_id}/hotkey/autoscan
// @API DCS GET /v2/{project_id}/instances/{instance_id}/configs
// @API DCS PUT /v2/{project_id}/instances/status
// @API DCS PUT /v2/{project_id}/instances/{instance_id}/ssl
// @API DCS GET /v2/{project_id}/instances/{instance_id}/ssl
// @API DCS GET /v2/{project_id}/instances/{instance_id}/tags
// @API DCS PUT /v2/{project_id}/instances/{instance_id}
// @API DCS POST /v2/{project_id}/instances/{instance_id}/password/reset
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
			"big_key_enable_auto_scan": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"big_key_schedule_at": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hot_key_enable_auto_scan": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"hot_key_schedule_at": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"big_key_updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hot_key_updated_at": {
				Type:     schema.TypeString,
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
		err = updateInstanceWhitelist(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// set parameters
	if v, ok := d.GetOk("parameters"); ok {
		parameters := v.(*schema.Set).List()
		err = updateParameters(ctx, d, client, schema.TimeoutUpdate, parameters)
		if err != nil {
			return diag.FromErr(err)
		}
		restart, err := checkDcsInstanceRestart(client, id, parameters)
		if err != nil {
			return diag.FromErr(err)
		}
		if restart {
			if err = restartDcsInstance(ctx, d, client, schema.TimeoutCreate); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if sslEnabled := d.Get("ssl_enable").(bool); sslEnabled {
		err = updateInstanceSsl(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("big_key_enable_auto_scan").(bool) || len(d.Get("big_key_schedule_at").([]interface{})) > 0 {
		err = updateBigKeyAutoScan(ctx, d, client, schema.TimeoutCreate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("hot_key_enable_auto_scan").(bool) || len(d.Get("hot_key_schedule_at").([]interface{})) > 0 {
		err = updateHotKeyAutoScan(ctx, d, client, schema.TimeoutCreate)
		if err != nil {
			return diag.FromErr(err)
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
		"rename_commands":        buildInstanceRenameCommandsBodyParams(d),
		"instance_backup_policy": buildInstanceBackupPolicyBodyParams(d),
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

func buildInstanceRenameCommandsBodyParams(d *schema.ResourceData) map[string]interface{} {
	renameCmds := d.Get("rename_commands").(map[string]interface{})
	if d.Get("engine") != "Redis" || len(renameCmds) == 0 {
		return nil
	}
	return renameCmds
}

func buildInstanceBackupPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
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

func updateParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string,
	parameters []interface{}) error {
	_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v2/{project_id}/instances/{instance_id}/async-configs",
		httpMethod:         "PUT",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   utils.RemoveNil(buildUpdateInstanceParametersBodyParams(parameters)),
		isRetry:            true,
		timeout:            timeout,
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) parameters: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateInstanceParametersBodyParams(parameters []interface{}) map[string]interface{} {
	params := make([]map[string]interface{}, 0, len(parameters))
	for _, parameter := range parameters {
		if v, ok := parameter.(map[string]interface{}); ok {
			params = append(params, map[string]interface{}{
				"param_id":    v["id"].(string),
				"param_name":  v["name"].(string),
				"param_value": v["value"].(string),
			})
		}
	}
	bodyParams := map[string]interface{}{
		"redis_config": params,
	}
	return bodyParams
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
	getRespBody, err := getInstanceField(client, getInstanceFieldParams{
		httpUrl:    "v2/{project_id}/instances/{instance_id}/configs",
		httpMethod: "GET",
		pathParams: map[string]string{"instance_id": instanceID},
	})
	if err != nil {
		log.Printf("[WARN] error fetching DCS instance(%s) parameters: %s", instanceID, err)
		return nil, nil, err
	}

	redisConfig := utils.PathSearch("redis_config", getRespBody, make([]interface{}, 0)).([]interface{})
	parametersMap := generateParametersMap(redisConfig)
	var params []map[string]interface{}
	restartParams := make([]string, 0)
	for _, parameter := range parameters {
		paramId := parameter.(map[string]interface{})["id"]
		if v, ok := parametersMap[paramId.(string)]; ok {
			name := utils.PathSearch("param_name", v, "").(string)
			params = append(params, map[string]interface{}{
				"id":    utils.PathSearch("param_id", v, nil),
				"name":  name,
				"value": utils.PathSearch("param_value", v, nil),
			})
			needRestart := utils.PathSearch("need_restart", v, false).(bool)
			if needRestart {
				restartParams = append(restartParams, name)
			}
		}
	}
	return params, restartParams, nil
}

func restartDcsInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string) error {
	_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:             "v2/{project_id}/instances/status",
		httpMethod:          "PUT",
		pathParams:          map[string]string{"instance_id": d.Id()},
		updateBodyParams:    utils.RemoveNil(buildRestartInstanceBodyParams(d)),
		isRetry:             true,
		timeout:             timeout,
		isWaitInstanceReady: true,
	})
	if err != nil {
		return fmt.Errorf("error restarting instance(%s): %s", d.Id(), err)
	}
	return nil
}

func buildRestartInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instances": []string{d.Id()},
		"action":    "restart",
	}
	return bodyParams
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
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", instance, make([]interface{}, 0)))),
	)

	mErr = multierror.Append(mErr, setDcsInstanceWhitelist(d, client)...)
	mErr = multierror.Append(mErr, setDcsInstanceBigKeyAutoScan(d, client)...)
	mErr = multierror.Append(mErr, setDcsInstanceHotKeyAutoScan(d, client)...)

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

func setDcsInstanceWhitelist(d *schema.ResourceData, client *golangsdk.ServiceClient) []error {
	getRespBody, err := getInstanceField(client, getInstanceFieldParams{
		httpUrl:    "v2/{project_id}/instance/{instance_id}/whitelist",
		httpMethod: "GET",
		pathParams: map[string]string{"instance_id": d.Id()},
	})
	if err != nil {
		log.Printf("[WARN] error fetching DCS instance(%s) whitelist: %s", d.Id(), err)
		return nil
	}

	var errs []error
	whitelists := flattenInstanceWhitelist(getRespBody)
	if len(whitelists) == 0 {
		// Set to the default value, otherwise it will prompt change after importing resources.
		errs = append(errs, d.Set("whitelist_enable", true))
	} else {
		errs = append(errs, d.Set("whitelist_enable", utils.PathSearch("enable_whitelist", getRespBody, nil)))
		errs = append(errs, d.Set("whitelists", flattenInstanceWhitelist(getRespBody)))
	}

	return errs
}

func setDcsInstanceBigKeyAutoScan(d *schema.ResourceData, client *golangsdk.ServiceClient) []error {
	getRespBody, err := getInstanceField(client, getInstanceFieldParams{
		httpUrl:    "v2/{project_id}/instances/{instance_id}/bigkey/autoscan",
		httpMethod: "GET",
		pathParams: map[string]string{"instance_id": d.Id()},
	})
	if err != nil {
		log.Printf("[WARN] error fetching DCS instance(%s) big key auto scan: %s", d.Id(), err)
		return nil
	}

	var errs []error
	errs = append(errs, d.Set("big_key_enable_auto_scan", utils.PathSearch("enable_auto_scan", getRespBody, nil)))
	errs = append(errs, d.Set("big_key_schedule_at", utils.PathSearch("schedule_at", getRespBody, nil)))
	errs = append(errs, d.Set("big_key_updated_at", utils.PathSearch("updated_at", getRespBody, nil)))

	return errs
}

func setDcsInstanceHotKeyAutoScan(d *schema.ResourceData, client *golangsdk.ServiceClient) []error {
	getRespBody, err := getInstanceField(client, getInstanceFieldParams{
		httpUrl:    "v2/{project_id}/instances/{instance_id}/hotkey/autoscan",
		httpMethod: "GET",
		pathParams: map[string]string{"instance_id": d.Id()},
	})
	if err != nil {
		log.Printf("[WARN] error fetching DCS instance(%s) hot key auto scan: %s", d.Id(), err)
		return nil
	}

	var errs []error
	errs = append(errs, d.Set("hot_key_enable_auto_scan", utils.PathSearch("enable_auto_scan", getRespBody, nil)))
	errs = append(errs, d.Set("hot_key_schedule_at", utils.PathSearch("schedule_at", getRespBody, nil)))
	errs = append(errs, d.Set("hot_key_updated_at", utils.PathSearch("updated_at", getRespBody, nil)))

	return errs
}

func flattenInstanceWhitelist(resp interface{}) []interface{} {
	curJson := utils.PathSearch("whitelist", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"group_name": utils.PathSearch("group_name", v, nil),
			"ip_address": utils.PathSearch("ip_list", v, nil),
		})
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

func generateParametersMap(redisConfig []interface{}) map[string]interface{} {
	parametersMap := make(map[string]interface{})
	for _, v := range redisConfig {
		parametersMap[utils.PathSearch("param_id", v, "").(string)] = v
	}
	return parametersMap
}

func resourceDcsInstancesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	// update basic params
	if d.HasChanges("port", "name", "description", "security_group_id", "backup_policy",
		"maintain_begin", "maintain_end", "rename_commands") {
		err = updateInstance(ctx, d, client)
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
		err = updateInstancePassword(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// resize instance
	if d.HasChanges("flavor", "capacity") {
		err = resizeDcsInstance(ctx, d, client, cfg)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = updateDcsTags(ctx, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("whitelists", "whitelist_enable") {
		err = updateInstanceWhitelist(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("parameters") {
		oRaw, nRaw := d.GetChange("parameters")
		changedParameters := nRaw.(*schema.Set).Difference(oRaw.(*schema.Set)).List()
		err = updateParameters(ctx, d, client, schema.TimeoutUpdate, changedParameters)
		if err != nil {
			return diag.FromErr(err)
		}
		// Sending parametersChanged to Read to warn users the instance needs a reboot.
		ctx = context.WithValue(ctx, ctxType("parametersChanged"), "true")
	}

	if d.HasChange("ssl_enable") {
		err = updateInstanceSsl(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("big_key_enable_auto_scan", "big_key_schedule_at") {
		err = updateBigKeyAutoScan(ctx, d, client, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("hot_key_enable_auto_scan", "hot_key_schedule_at") {
		err = updateHotKeyAutoScan(ctx, d, client, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "dcs",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", d.Id(), err)
		}
	}

	return resourceDcsInstancesRead(ctx, d, meta)
}

func updateInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v2/{project_id}/instances/{instance_id}",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: utils.RemoveNil(buildUpdateInstanceBodyParams(d)),
	})
	if err != nil {
		return fmt.Errorf("error updating instance: %s", err)
	}
	return nil
}

func buildUpdateInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                   d.Get("name"),
		"description":            d.Get("description"),
		"port":                   d.Get("port"),
		"rename_commands":        buildInstanceRenameCommandsBodyParams(d),
		"maintain_begin":         d.Get("maintain_begin"),
		"maintain_end":           d.Get("maintain_end"),
		"security_group_id":      utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		"instance_backup_policy": buildInstanceBackupPolicyBodyParams(d),
	}
	return bodyParams
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
		return fmt.Errorf("error while waiting for DCS instance(%s) port update completed: %#v", d.Id(), err)
	}
	return nil
}

func refreshDcsInstancePort(c *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getDcsInstanceByID(c, id)
		if err != nil {
			return nil, "ERROR", err
		}
		port := utils.PathSearch("port", instance, float64(0)).(float64)
		return instance, strconv.Itoa(int(port)), nil
	}
}

func updateInstancePassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v2/{project_id}/instances/{instance_id}/password/reset",
		httpMethod:       "POST",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateInstancePasswordBodyParams(d),
		isRetry:          true,
		timeout:          schema.TimeoutUpdate,
	})
	if err != nil {
		return fmt.Errorf("error updating instance password: %s", err)
	}
	return nil
}

func buildUpdateInstancePasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := make(map[string]interface{})
	password := d.Get("password")
	if password == "" {
		bodyParams["no_password_access"] = true
	} else {
		bodyParams["new_password"] = password
		bodyParams["no_password_access"] = false
	}
	return bodyParams
}

func resizeDcsInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, cfg *config.Config) error {
	oVal, nVal := d.GetChange("flavor")
	oldSpecCode := oVal.(string)
	newSpecCode := nVal.(string)
	bodyParams, err := buildResizeInstanceOpt(client, d, oldSpecCode, newSpecCode)
	if err != nil {
		return err
	}

	checkOrderExpression := ""
	var bssClient *golangsdk.ServiceClient
	if d.Get("charging_mode").(string) == chargeModePrePaid {
		checkOrderExpression = "order_id"
		c, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}
		bssClient = c
	}
	_, err = updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v2/{project_id}/instances/{instance_id}/resize",
		httpMethod:           "POST",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     bodyParams,
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkOrderExpression: checkOrderExpression,
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) flavor: %s", d.Id(), err)
	}

	// check the result of the change
	instance, err := getDcsInstanceByID(client, d.Id())
	if err != nil {
		return fmt.Errorf("error getting DCS instance: %s", err)
	}
	specCode := utils.PathSearch("spec_code", instance, "").(string)
	if specCode != newSpecCode {
		return fmt.Errorf("change flavor failed, after changed the DCS flavor still is: %s, expected: %s",
			specCode, newSpecCode)
	}
	return nil
}

func buildResizeInstanceOpt(client *golangsdk.ServiceClient, d *schema.ResourceData, oldSpecCode,
	newSpecCode string) (map[string]interface{}, error) {
	if oldSpecCode == newSpecCode {
		return nil, errors.New("the param flavor is invalid")
	}
	bodyParams := map[string]interface{}{
		"spec_code":    newSpecCode,
		"new_capacity": d.Get("capacity"),
	}
	if d.Get("charging_mode").(string) == chargeModePrePaid {
		bodyParams["bss_param"] = map[string]interface{}{
			"is_auto_pay": true,
		}
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
	bodyParams["change_type"] = changeType
	if changeType == "createReplication" {
		azCodes, err := getAzCode(d, client)
		if err != nil {
			return nil, err
		}
		bodyParams["available_zones"] = azCodes
	}
	if changeType == "deleteReplication" {
		newCacheMode := utils.PathSearch("cache_mode", oldFlavor, "").(string)
		if newCacheMode == "ha" {
			bodyParams["node_list"] = utils.ExpandToStringList(d.Get("deleted_nodes").([]interface{}))
		} else if newCacheMode == "cluster" {
			azCodes, err := getAzCode(d, client)
			if err != nil {
				return nil, err
			}
			bodyParams["reserved_ip"] = utils.ExpandToStringList(d.Get("reserved_ips").([]interface{}))
			bodyParams["available_zones"] = azCodes
		}
	}
	return bodyParams, nil
}

func getFlavorChangeType(oldFlavor, newFlavor interface{}) string {
	oldCacheMode := utils.PathSearch("cache_mode", oldFlavor, "").(string)
	newCacheMode := utils.PathSearch("cache_mode", newFlavor, "").(string)
	// if the cache mode is different, it indicates the type has been changed
	if oldCacheMode != newCacheMode {
		return "instanceType"
	}

	oldReplicaCount := utils.PathSearch("replica_count", oldFlavor, float64(0)).(float64)
	newReplicaCount := utils.PathSearch("replica_count", newFlavor, float64(0)).(float64)
	// indicates the replica count increase, should add replica
	if oldReplicaCount < newReplicaCount {
		return "createReplication"
	}
	// indicates the replica count decrease, should delete replica
	if oldReplicaCount > newReplicaCount {
		return "deleteReplication"
	}
	// indicates only the capacity been changed
	return ""
}

func getFlavorBySpecCode(client *golangsdk.ServiceClient, specCode string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/flavors?spec_code={spec_code}"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{spec_code}", specCode)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting DCS flavor by specCode %s: %s", specCode, err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	flavor := utils.PathSearch("flavors[0]", getRespBody, nil)
	if flavor == nil {
		return nil, fmt.Errorf("the result queried by specCode(%s) is empty", specCode)
	}
	return flavor, nil
}

func updateDcsTags(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("tags")
	oldTags := oldRaw.(map[string]interface{})
	newTags := newRaw.(map[string]interface{})
	// remove old tags
	if len(oldTags) > 0 {
		_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
			httpUrl:          "v2/{project_id}/dcs/{instance_id}/tags/action",
			httpMethod:       "POST",
			pathParams:       map[string]string{"instance_id": d.Id()},
			updateBodyParams: buildDcsTagsParams("delete", oldTags),
		})
		if err != nil {
			return fmt.Errorf("error updating instance(%s) tags: %s", d.Id(), err)
		}
	}

	if len(newTags) > 0 {
		_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
			httpUrl:          "v2/{project_id}/dcs/{instance_id}/tags/action",
			httpMethod:       "POST",
			pathParams:       map[string]string{"instance_id": d.Id()},
			updateBodyParams: buildDcsTagsParams("create", newTags),
		})
		if err != nil {
			return fmt.Errorf("error updating instance(%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func buildDcsTagsParams(action string, rawTags map[string]interface{}) map[string]interface{} {
	if len(rawTags) == 0 {
		return nil
	}
	tags := make([]map[string]interface{}, 0)
	for key, value := range rawTags {
		tags = append(tags, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}
	bodyPrams := map[string]interface{}{
		"action": action,
		"tags":   tags,
	}
	return bodyPrams
}

func updateInstanceWhitelist(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v2/{project_id}/instance/{instance_id}/whitelist",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: utils.RemoveNil(buildUpdateInstanceWhitelistBodyParams(d)),
	})
	if err != nil {
		return fmt.Errorf("error updating instance whitelist: %s", err)
	}
	return nil
}

func buildUpdateInstanceWhitelistBodyParams(d *schema.ResourceData) map[string]interface{} {
	groupList := d.Get("whitelists").(*schema.Set)
	params := make([]map[string]interface{}, 0, groupList.Len())
	for _, v := range groupList.List() {
		group := v.(map[string]interface{})
		params = append(params, map[string]interface{}{
			"group_name": group["group_name"],
			"ip_list":    utils.ExpandToStringList(group["ip_address"].([]interface{})),
		})
	}
	bodyParams := map[string]interface{}{
		"enable_whitelist": d.Get("whitelist_enable"),
		"whitelist":        params,
	}
	return bodyParams
}

func updateInstanceSsl(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v2/{project_id}/instances/{instance_id}/ssl",
		httpMethod:         "PUT",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   utils.RemoveNil(buildUpdateInstanceSslBodyParams(d)),
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error updating instance SSL: %s", err)
	}
	return nil
}

func buildUpdateInstanceSslBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enabled": d.Get("ssl_enable"),
	}
	return bodyParams
}

func updateBigKeyAutoScan(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string) error {
	_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v2/{project_id}/instances/{instance_id}/bigkey/autoscan",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateBigKeyAutoScanBodyParams(d),
		isRetry:          true,
		timeout:          timeout,
	})
	if err != nil {
		return fmt.Errorf("error updating instance big key auto scan: %s", err)
	}
	return nil
}

func buildUpdateBigKeyAutoScanBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enable_auto_scan": d.Get("big_key_enable_auto_scan"),
		"schedule_at":      utils.ExpandToStringList(d.Get("big_key_schedule_at").([]interface{})),
	}
	return bodyParams
}

func updateHotKeyAutoScan(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string) error {
	_, err := updateDcsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v2/{project_id}/instances/{instance_id}/hotkey/autoscan",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: utils.RemoveNil(buildUpdateHotKeyAutoScanBodyParams(d)),
		isRetry:          true,
		timeout:          timeout,
	})
	if err != nil {
		return fmt.Errorf("error updating instance hot key auto scan: %s", err)
	}
	return nil
}

func buildUpdateHotKeyAutoScanBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enable_auto_scan": d.Get("hot_key_enable_auto_scan"),
		"schedule_at":      utils.ExpandToStringList(d.Get("hot_key_schedule_at").([]interface{})),
	}
	return bodyParams
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

	var (
		httpUrl = "v2/available-zones"
	)
	getPath := client.Endpoint + httpUrl

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	availableZones := utils.PathSearch("available_zones", getRespBody, make([]interface{}, 0)).([]interface{})
	mapping := make(map[string]string)
	for _, v := range availableZones {
		id := utils.PathSearch("id", v, "").(string)
		code := utils.PathSearch("code", v, "").(string)
		mapping[id] = code
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
