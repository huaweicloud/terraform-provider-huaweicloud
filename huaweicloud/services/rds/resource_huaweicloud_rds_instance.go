package rds

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/rds/v3/backups"
	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ctxType string

// ResourceRdsInstance is the impl for huaweicloud_rds_instance resource
// @API RDS POST /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/jobs
// @API RDS GET /v3/{project_id}/instances
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/alias
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/ssl
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/ops-window
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/failover/strategy
// @API RDS POST /v3/{project_id}/instances/{id}/tags/action
// @API RDS PUT /v3.1/{project_id}/configurations/{config_id}/apply
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/configurations
// @API RDS PUT /v3.1/{project_id}/instances/{instance_id}/configurations
// @API RDS POST /v3/{project_id}/instances/{instance_id}/action
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/disk-auto-expansion
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/backups/policy
// @API RDS POST /v3/{project_id}/instances/{instance_id}/action/shutdown
// @API RDS POST /v3/{project_id}/instances/{instance_id}/action/startup
// @API RDS GET /v3/{project_id}/instances/{instance_id}/disk-auto-expansion
// @API RDS GET /v3/{project_id}/instances/{instance_id}/backups/policy
// @API RDS GET /v3/{project_id}/instances/{instance_id}/configurations
// @API RDS GET /v3/{project_id}/instances/{instance_id}/binlog/clear-policy
// @API RDS GET /v3/{project_id}/instances/{instance_id}/msdtc/hosts
// @API RDS GET /v3/{project_id}/instances/{instance_id}/tde-status
// @API RDS GET /v3/{project_id}/instances/{instance_id}/second-level-monitor
// @API RDS GET /v3/{project_id}/instances/{instance_id}/db-auto-upgrade
// @API RDS GET /v3/{project_id}/instances/{instance_id}/storage-used-space
// @API RDS GET /v3/{project_id}/instances/{instance_id}/replication/status
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/name
// @API RDS POST /v3/{project_id}/instances/{instance_id}/migrateslave
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/failover/mode
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/collations
// @API RDS POST /v3/{project_id}/instances/{instance_id}/msdtc/host
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/msdtc/host
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/tde
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/readonly-status
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/modify-dns
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/second-level-monitor
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/slowlog-sensitization/{status}
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/db-auto-upgrade
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/port
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/ip
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/security-group
// @API RDS POST /v3/{project_id}/instances/{instance_id}/password
// @API RDS POST /v3/{project_id}/instances/{instance_id}/to-period
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/binlog/clear-policy
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrat
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceRdsInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsInstanceCreate,
		ReadContext:   resourceRdsInstanceRead,
		UpdateContext: resourceRdsInstanceUpdate,
		DeleteContext: resourceRdsInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(30 * time.Minute),
			Update:  schema.DefaultTimeout(30 * time.Minute),
			Delete:  schema.DefaultTimeout(30 * time.Minute),
			Default: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: utils.SuppressCaseDiffs(),
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk_encryption_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"limit_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"volume.0.trigger_threshold"},
						},
						"trigger_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"volume.0.limit_size"},
						},
					},
				},
			},
			"restore": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"period"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"backup_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"database_name": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"keep_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "schema: Required",
						},
						"period": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"lower_case_table_names": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateIP,
			},
			"private_dns_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ha_replication_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"power_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ON", "OFF", "REBOOT",
				}, false),
			},
			"param_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"switch_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"binlog_retention_hours": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"msdtc_hosts": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tde_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"read_write_permissions": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"readonly", "readwrite",
				}, false),
			},
			"rotate_day": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"tde_enabled"},
			},
			"secret_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"tde_enabled"},
			},
			"secret_name": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"tde_enabled"},
			},
			"secret_version": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"tde_enabled"},
			},
			"seconds_level_monitoring_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"seconds_level_monitoring_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"seconds_level_monitoring_enabled"},
			},
			"minor_version_auto_upgrade_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"slow_log_show_original_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dss_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"maintain_begin"},
			},
			"is_flexus": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_dns_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"storage_used_space": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"replication_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// charging_mode,  period_unit and period only support changing post-paid to pre-paid billing mode.
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":   common.SchemaAutoPay(nil),
		},
	}
}

func buildRdsInstanceDBPort(d *schema.ResourceData) string {
	if v, ok := d.GetOk("db.0.port"); ok {
		return strconv.Itoa(v.(int))
	}
	return ""
}

func isMySQLDatabase(d *schema.ResourceData) bool {
	dbType := d.Get("db.0.type").(string)
	// Database type is not case sensitive.
	return strings.ToLower(dbType) == "mysql"
}

func isSQLServerDatabase(d *schema.ResourceData) bool {
	dbType := d.Get("db.0.type").(string)
	// Database type is not case sensitive.
	return strings.ToLower(dbType) == "sqlserver"
}

func resourceRdsInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if d.Get("ssl_enable").(bool) && !isMySQLDatabase(d) {
		return diag.Errorf("only MySQL database support SSL enable and disable")
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestBody := buildCreateInstanceBodyParams(d, cfg, region)
	requestBody["password"] = utils.ValueIgnoreEmpty(d.Get("db.0.password"))
	createOpt.JSONBody = utils.RemoveNil(requestBody)

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating RDS instance: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := utils.PathSearch("instance.id", createRespBody, "").(string)
	if instanceID == "" {
		return diag.Errorf("error creating RDS instance: ID is not found in API response")
	}
	d.SetId(instanceID)

	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	// wait for order success
	if orderId != "" {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId != "" {
		if err = checkRDSInstanceJobFinish(client, jobId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error creating instance (%s): %s", instanceID, err)
		}
	}
	// for prePaid charge mode
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE", "BACKING UP"},
		Refresh:      rdsInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RDS instance (%s) creation completed: %s", instanceID, err)
	}

	if err = updateRdsInstanceDescription(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceSSLConfig(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceMaintainWindow(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("switch_strategy"); ok && v.(string) != "reliability" {
		if err = updateRdsInstanceSwitchStrategy(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if err = updateBinlogRetentionHours(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateMsdtcHosts(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateTde(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("read_write_permissions"); ok && v.(string) == "readonly" {
		if err = updateReadWritePermissions(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if err = updatePrivateDNSNamePrefix(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateSecondLevelMonitoring(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("slow_log_show_original_status"); ok && v.(string) == "on" {
		if err = updateSlowLogShowOriginalStatus(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("minor_version_auto_upgrade_enabled").(bool) {
		if err = updateAutoUpgradeSwitchOption(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", instanceID, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of RDS instance (%s): %s", instanceID, tagErr)
		}
	}

	// Set Parameters
	if parameters := d.Get("parameters").(*schema.Set); parameters.Len() > 0 {
		clientV31, err := cfg.RdsV31Client(region)
		if err != nil {
			return diag.Errorf("error creating RDS V3.1 client: %s", err)
		}
		if err = initializeParameters(ctx, d, client, clientV31, instanceID, parameters); err != nil {
			return diag.FromErr(err)
		}
	}

	if size := d.Get("volume.0.limit_size").(int); size > 0 {
		if err = enableVolumeAutoExpand(ctx, d, client, size); err != nil {
			return diag.FromErr(err)
		}
	}

	if err = updateRdsInstanceBackupStrategy(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if action, ok := d.GetOk("power_action"); ok && action == "OFF" {
		err = updatePowerAction(ctx, d, client, action.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRdsInstanceRead(ctx, d, meta)
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, cfg *config.Config, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"flavor_ref":            d.Get("flavor"),
		"vpc_id":                d.Get("vpc_id"),
		"subnet_id":             d.Get("subnet_id"),
		"security_group_id":     d.Get("security_group_id"),
		"configuration_id":      utils.ValueIgnoreEmpty(d.Get("param_group_id")),
		"time_zone":             utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"data_vip":              utils.ValueIgnoreEmpty(d.Get("fixed_ip")),
		"disk_encryption_id":    utils.ValueIgnoreEmpty(d.Get("volume.0.disk_encryption_id")),
		"collation":             utils.ValueIgnoreEmpty(d.Get("collation")),
		"port":                  utils.ValueIgnoreEmpty(buildRdsInstanceDBPort(d)),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"region":                region,
		"availability_zone":     buildRdsInstanceAvailabilityZone(d),
		"datastore":             buildCreateInstanceDatastoreBodyParams(d),
		"volume":                buildCreateInstanceVolumeBodyParams(d),
		"ha":                    buildCreateInstanceHaBodyParams(d),
		"unchangeable_param":    buildCreateInstanceUnchangeableParamBodyParams(d),
		"restore_point":         buildCreateInstanceRestoreBodyParams(d),
		"dsspool_id":            utils.ValueIgnoreEmpty(d.Get("dss_pool_id")),
		"charge_info":           buildCreateInstanceChargeInfoBodyParams(d),
		"is_flexus":             utils.ValueIgnoreEmpty(d.Get("is_flexus")),
	}
	return bodyParams
}

func buildCreateInstanceDatastoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	dbRaw := d.Get("db").([]interface{})
	if len(dbRaw) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"type":    dbRaw[0].(map[string]interface{})["type"],
		"version": dbRaw[0].(map[string]interface{})["version"],
	}
	return bodyParams
}

func buildCreateInstanceVolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	volumeRaw := d.Get("volume").([]interface{})
	if len(volumeRaw) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"type": volumeRaw[0].(map[string]interface{})["type"],
		"size": volumeRaw[0].(map[string]interface{})["size"],
	}
	return bodyParams
}

func buildCreateInstanceHaBodyParams(d *schema.ResourceData) map[string]interface{} {
	v, ok := d.GetOk("ha_replication_mode")
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"mode":             "ha",
		"replication_mode": v.(string),
	}
	return bodyParams
}

func buildCreateInstanceUnchangeableParamBodyParams(d *schema.ResourceData) map[string]interface{} {
	v, ok := d.GetOk("lower_case_table_names")
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"lower_case_table_names": v.(string),
	}
	return bodyParams
}

func buildCreateInstanceRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	restoreRaw := d.Get("restore").([]interface{})
	if len(restoreRaw) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"type":          "backup",
		"instance_id":   restoreRaw[0].(map[string]interface{})["instance_id"],
		"backup_id":     utils.ValueIgnoreEmpty(restoreRaw[0].(map[string]interface{})["backup_id"]),
		"database_name": utils.ValueIgnoreEmpty(restoreRaw[0].(map[string]interface{})["database_name"]),
	}
	return bodyParams
}

func buildCreateInstanceChargeInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	if d.Get("charging_mode") != "prePaid" {
		return nil
	}

	bodyParams := map[string]interface{}{
		"charge_mode": d.Get("charging_mode"),
		"period_type": d.Get("period_unit"),
		"period_num":  d.Get("period"),
	}
	if d.Get("auto_pay").(string) != "false" {
		bodyParams["is_auto_pay"] = true
	}
	if d.Get("auto_renew").(string) != "false" {
		bodyParams["is_auto_renew"] = true
	}
	return bodyParams
}

func resourceRdsInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Id()
	instance, err := GetRdsInstanceByID(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting RDS instance")
	}

	mErr := multierror.Append(nil,
		d.Set("region", utils.PathSearch("region", instance, nil)),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("description", utils.PathSearch("alias", instance, nil)),
		d.Set("status", utils.PathSearch("status", instance, nil)),
		d.Set("created", utils.PathSearch("created", instance, nil)),
		d.Set("ha_replication_mode", utils.PathSearch("ha.replication_mode", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", instance, nil)),
		d.Set("flavor", utils.PathSearch("flavor_ref", instance, nil)),
		d.Set("time_zone", utils.PathSearch("time_zone", instance, nil)),
		d.Set("collation", utils.PathSearch("collation", instance, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", instance, nil)),
		d.Set("switch_strategy", utils.PathSearch("switch_strategy", instance, nil)),
		d.Set("charging_mode", utils.PathSearch("charge_info.charge_mode", instance, nil)),
		d.Set("ssl_enable", utils.PathSearch("enable_ssl", instance, nil)),
		d.Set("private_dns_names", utils.PathSearch("private_dns_names", instance, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", instance, make([]interface{}, 0)))),
		d.Set("public_ips", utils.PathSearch("public_ips", instance, nil)),
		d.Set("private_ips", utils.PathSearch("private_ips", instance, nil)),
		d.Set("fixed_ip", utils.PathSearch("private_ips[0]", instance, nil)),
		d.Set("volume", flattenInstanceVolume(client, instance, instanceID)),
		d.Set("db", flattenInstanceDb(d, instance)),
		d.Set("nodes", flattenInstanceNodes(instance)),
	)
	if v := utils.PathSearch("private_dns_names", instance, make([]interface{}, 0)).([]interface{}); len(v) > 0 {
		privateDNSNamePrefix := strings.Split(v[0].(string), ".")[0]
		mErr = multierror.Append(mErr, d.Set("private_dns_name_prefix", privateDNSNamePrefix))
	}
	if v := utils.PathSearch("maintenance_window", instance, "").(string); v != "" {
		maintainWindow := strings.Split(v, "-")
		mErr = multierror.Append(mErr, d.Set("maintain_begin", maintainWindow[0]))
		mErr = multierror.Append(mErr, d.Set("maintain_end", maintainWindow[1]))
	}
	status := utils.PathSearch("status", instance, "").(string)
	if status != "SHUTDOWN" {
		mErr = multierror.Append(mErr, d.Set("backup_strategy", flattenInstanceBackupStrategy(client, instance, instanceID)))
	}

	if isMySQLDatabase(d) {
		binlogRetentionHours, err := instances.GetBinlogRetentionHours(client, instanceID).Extract()
		if err != nil {
			log.Printf("[WARN] error getting RDS binlog retention hours: %s", err)
		} else {
			mErr = multierror.Append(mErr, d.Set("binlog_retention_hours", binlogRetentionHours.BinlogRetentionHours))
		}
	}

	if isSQLServerDatabase(d) && status != "SHUTDOWN" {
		msdtcHosts, err := instances.GetMsdtcHosts(client, instanceID)
		if err != nil {
			log.Printf("[WARN] error getting RDS msdtc hosts: %s", err)
		} else {
			hosts := make([]map[string]interface{}, 0, len(msdtcHosts))
			for _, msdtcHost := range msdtcHosts {
				hosts = append(hosts, map[string]interface{}{
					"id":        msdtcHost.Id,
					"ip":        msdtcHost.Host,
					"host_name": msdtcHost.HostName,
				})
			}
			mErr = multierror.Append(mErr, d.Set("msdtc_hosts", hosts))
		}
	}

	if isSQLServerDatabase(d) {
		tdeStatus, err := instances.GetTdeStatus(client, instanceID).Extract()
		if err != nil {
			log.Printf("[WARN] error getting TDE of the instance: %s", err)
		} else {
			tdeEnabled := false
			if tdeStatus.TdeStatus == "open" {
				tdeEnabled = true
			}
			mErr = multierror.Append(mErr, d.Set("tde_enabled", tdeEnabled))
		}
	}

	if isMySQLDatabase(d) {
		secondsLevelMonitoring, err := instances.GetSecondLevelMonitoring(client, instanceID).Extract()
		if err != nil {
			log.Printf("[WARN] fetching RDS seconds level monitoring failed: %s", err)
		} else {
			mErr = multierror.Append(mErr, d.Set("seconds_level_monitoring_enabled", secondsLevelMonitoring.SwitchOption))
			mErr = multierror.Append(mErr, d.Set("seconds_level_monitoring_interval", secondsLevelMonitoring.Interval))
		}
	}

	mErr = multierror.Append(mErr, setAutoUpgradeSwitchOption(d, client))
	mErr = multierror.Append(mErr, setStorageUsedSpace(d, client))
	mErr = multierror.Append(mErr, setReplicationStatus(d, client))

	diagErr := setRdsInstanceParameters(ctx, d, client, instanceID)
	resErr := append(diag.FromErr(mErr.ErrorOrNil()), diagErr...)

	return resErr
}

func flattenInstanceVolume(client *golangsdk.ServiceClient, instance interface{}, instanceID string) []interface{} {
	volume := map[string]interface{}{
		"type":               utils.PathSearch("volume.type", instance, nil),
		"size":               utils.PathSearch("volume.size", instance, nil),
		"disk_encryption_id": utils.PathSearch("disk_encryption_id", instance, nil),
	}

	// Only MySQL engines are supported.
	resp, err := instances.GetAutoExpand(client, instanceID)
	if err != nil {
		log.Printf("[ERROR] error query automatic expansion configuration of the instance storage: %s", err)
	}
	if resp.SwitchOption {
		volume["limit_size"] = resp.LimitSize
		volume["trigger_threshold"] = resp.TriggerThreshold
	}

	return []interface{}{volume}
}

func flattenInstanceDb(d *schema.ResourceData, instance interface{}) []interface{} {
	database := map[string]interface{}{
		"type":      utils.PathSearch("datastore.type", instance, nil),
		"version":   utils.PathSearch("datastore.version", instance, nil),
		"port":      utils.PathSearch("port", instance, nil),
		"user_name": utils.PathSearch("db_user_name", instance, nil),
	}
	if len(d.Get("db").([]interface{})) > 0 {
		database["password"] = d.Get("db.0.password")
	}
	return []interface{}{database}
}

func flattenInstanceBackupStrategy(client *golangsdk.ServiceClient, instance interface{}, instanceID string) []interface{} {
	backupStrategy, err := backups.Get(client, instanceID).Extract()
	if err != nil {
		log.Printf("[ERROR] error query backup strategy of the instance storage: %s", err)
	}

	backup := map[string]interface{}{
		"start_time": utils.PathSearch("backup_strategy.start_time", instance, nil),
		"keep_days":  utils.PathSearch("backup_strategy.keep_days", instance, nil),
	}
	if backupStrategy != nil {
		backup["period"] = backupStrategy.Period
	}

	return []interface{}{backup}
}

func flattenInstanceNodes(instance interface{}) []interface{} {
	nodesJson := utils.PathSearch("nodes", instance, make([]interface{}, 0))
	nodeArray := nodesJson.([]interface{})
	if len(nodeArray) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(nodeArray))
	for _, v := range nodeArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"role":              utils.PathSearch("role", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
		})
	}
	return rst
}

func setRdsInstanceParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) diag.Diagnostics {
	// Set Parameters
	configs, err := instances.GetConfigurations(client, instanceID).Extract()
	if err != nil {
		log.Printf("[WARN] error fetching parameters of instance (%s): %s", instanceID, err)
		return nil
	}

	var configurationRestart bool
	var paramRestart []string
	var params []map[string]interface{}
	rawParameterList := d.Get("parameters").(*schema.Set).List()
	for _, v := range configs.Parameters {
		if v.Restart {
			configurationRestart = true
		}
		for _, parameter := range rawParameterList {
			name := parameter.(map[string]interface{})["name"]
			if v.Name == name {
				p := map[string]interface{}{
					"name":  v.Name,
					"value": v.Value,
				}
				params = append(params, p)
				if v.Restart {
					paramRestart = append(paramRestart, v.Name)
				}
				break
			}
		}
	}

	var diagnostics diag.Diagnostics
	if len(params) > 0 {
		if err = d.Set("parameters", params); err != nil {
			log.Printf("error saving parameters to RDS instance (%s): %s", instanceID, err)
		}
		if len(paramRestart) > 0 && ctx.Value(ctxType("parametersChanged")) == "true" {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Parameters Changed",
				Detail:   fmt.Sprintf("Parameters %s changed which needs reboot.", paramRestart),
			})
		}
	}
	if configurationRestart && ctx.Value(ctxType("configurationChanged")) == "true" {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Configuration Changed",
			Detail:   "Configuration changed which needs reboot.",
		})
	}
	if len(diagnostics) > 0 {
		return diagnostics
	}
	return nil
}

func setAutoUpgradeSwitchOption(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	getRespBody, err := getInstanceField(client, getInstanceFieldParams{
		httpUrl:    "v3/{project_id}/instances/{instance_id}/db-auto-upgrade",
		httpMethod: "GET",
		pathParams: map[string]string{"instance_id": d.Id()},
	})
	if err != nil {
		log.Printf("[WARN] error retrieving RDS instance(%s) auto upgrade switch option: %s", d.Id(), err)
		return nil
	}
	return d.Set("minor_version_auto_upgrade_enabled", utils.PathSearch("switch_option", getRespBody, nil))
}

func setStorageUsedSpace(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	getRespBody, err := getInstanceField(client, getInstanceFieldParams{
		httpUrl:    "v3/{project_id}/instances/{instance_id}/storage-used-space",
		httpMethod: "GET",
		pathParams: map[string]string{"instance_id": d.Id()},
	})
	if err != nil {
		log.Printf("[WARN] error retrieving get RDS instance(%s) storage used space: %s", d.Id(), err)
		return nil
	}
	return d.Set("storage_used_space", flattenInstanceResponseBodyStorageUsedSpace(getRespBody))
}

func flattenInstanceResponseBodyStorageUsedSpace(resp interface{}) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"node_id": utils.PathSearch("node_id", resp, nil),
			"used":    utils.PathSearch("used", resp, nil),
		},
	}
	return rst
}

func setReplicationStatus(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	getRespBody, err := getInstanceField(client, getInstanceFieldParams{
		httpUrl:    "v3/{project_id}/instances/{instance_id}/replication/status",
		httpMethod: "GET",
		pathParams: map[string]string{"instance_id": d.Id()},
	})
	if err != nil {
		log.Printf("[WARN] error retrieving RDS instance(%s) replication status: %s", d.Id(), err)
		return nil
	}
	return d.Set("replication_status", utils.PathSearch("replication_status", getRespBody, nil))
}

func resourceRdsInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.RdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}
	clientV31, err := cfg.RdsV31Client(region)
	if err != nil {
		return diag.Errorf("error creating RDS V3.1 client: %s", err)
	}

	instanceID := d.Id()

	// if power_action is changed from OFF to ON, the instance should be start first
	powerAction := d.Get("power_action").(string)
	if d.HasChanges("power_action") && powerAction == "ON" {
		if err = updatePowerAction(ctx, d, client, powerAction); err != nil {
			return diag.FromErr(err)
		}
	}

	if err = updateRdsInstanceName(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceDescription(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateAvailabilityZone(ctx, cfg, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceFlavor(ctx, d, cfg, client, true); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceVolumeSize(ctx, d, cfg, client); err != nil {
		return diag.FromErr(err)
	}
	if err = updateRdsInstanceBackupStrategy(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceMaintainWindow(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceReplicationMode(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceSwitchStrategy(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceCollation(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceDBPort(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceFixedIp(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceSecurityGroup(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceSSLConfig(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsRootPassword(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", instanceID)
		if tagErr != nil {
			return diag.Errorf("error updating tags of RDS instance (%s): %s", instanceID, tagErr)
		}
	}

	if d.HasChange("charging_mode") {
		if d.Get("charging_mode").(string) == "postPaid" {
			return diag.Errorf("error updating the charging mode of the RDS instance (%s): %s", d.Id(),
				"only support changing post-paid instance to pre-paid")
		}
		if err = updateBillingModeToPeriod(ctx, d, cfg, client); err != nil {
			return diag.FromErr(err)
		}
	} else if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceID); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", instanceID, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceID,
			ResourceType: "rds",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err = cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if ctx, err = updateConfiguration(ctx, d, client, clientV31, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if ctx, err = updateRdsParameters(ctx, d, client, clientV31, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateVolumeAutoExpand(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateBinlogRetentionHours(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateMsdtcHosts(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateTde(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateReadWritePermissions(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updatePrivateDNSNamePrefix(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateSecondLevelMonitoring(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateSlowLogShowOriginalStatus(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateAutoUpgradeSwitchOption(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	// if power_action is changed from ON to OFF/REBOOT, the instance should be close/restart at the end
	if d.HasChanges("power_action") && (powerAction == "OFF" || powerAction == "REBOOT") {
		if err = updatePowerAction(ctx, d, client, powerAction); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRdsInstanceRead(ctx, d, meta)
}

func resourceRdsInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating rds client: %s ", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Deleting Instance %s", id)
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		resourceIds := []string{id}
		// the image of SQL server is come from cloud market, when creating an SQL server instance resource, two order
		// will be created, one is instance order, the other is market image order, so it is needed to unsubscribe the
		// two order when unsubscribe the instance
		if strings.ToLower(d.Get("db.0.type").(string)) == "sqlserver" {
			resourceIds = append(resourceIds, fmt.Sprintf("%s%s", id, ".marketimage"))
		}
		retryFunc := func() (interface{}, bool, error) {
			err = common.UnsubscribePrePaidResource(d, config, resourceIds)
			retry, err := handleDeletionError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, id),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error unsubscribe RDS instance: %s", err)
		}
	} else {
		retryFunc := func() (interface{}, bool, error) {
			result := instances.Delete(client, id)
			retry, err := handleDeletionError(result.Err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, id),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE", "BACKING UP"},
		Target:       []string{"DELETED"},
		Refresh:      rdsInstanceStateRefreshFunc(client, id),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        15 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for rds instance (%s) to be deleted: %s ",
			id, err)
	}

	log.Printf("[DEBUG] Successfully deleted RDS instance %s", id)
	return nil
}

func GetRdsInstanceByID(client *golangsdk.ServiceClient, instanceID string) (interface{}, error) {
	instance, err := getRdsInstanceByIdAndFlexus(client, instanceID, false)
	if err != nil {
		return nil, err
	}
	if instance != nil {
		return instance, nil
	}

	// if rds instance is nil, then get flexus instance
	instance, err = getRdsInstanceByIdAndFlexus(client, instanceID, true)
	if err != nil {
		return nil, err
	}
	if instance != nil {
		return instance, nil
	}
	return nil, golangsdk.ErrDefault404{}
}

func getRdsInstanceByIdAndFlexus(client *golangsdk.ServiceClient, instanceID string, isFlexus bool) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances?id={instance_id}"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	if isFlexus {
		getPath = fmt.Sprintf("%s&group_type=flexus", getPath)
	}

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

	return utils.PathSearch("instances[0]", getRespBody, nil), nil
}

func buildRdsInstanceAvailabilityZone(d *schema.ResourceData) string {
	azList := utils.ExpandToStringList(d.Get("availability_zone").([]interface{}))
	return strings.Join(azList, ",")
}

func buildRdsInstanceParameters(params *schema.Set) instances.ModifyConfigurationOpts {
	var configOpts instances.ModifyConfigurationOpts

	values := make(map[string]string)
	for _, v := range params.List() {
		key := v.(map[string]interface{})["name"].(string)
		value := v.(map[string]interface{})["value"].(string)
		values[key] = value
	}
	configOpts.Values = values
	return configOpts
}

func initializeParameters(ctx context.Context, d *schema.ResourceData, client, clientV31 *golangsdk.ServiceClient,
	instanceID string, parametersRaw *schema.Set) error {
	configOpts := buildRdsInstanceParameters(parametersRaw)
	err := modifyParameters(ctx, d, client, clientV31, instanceID, &configOpts)
	if err != nil {
		return err
	}

	// Check if we need to restart
	restart, err := checkRdsInstanceRestart(client, instanceID, parametersRaw.List())
	if err != nil {
		return err
	}

	if restart {
		return restartRdsInstance(ctx, d.Timeout(schema.TimeoutCreate), client, d)
	}
	return nil
}

func checkRdsInstanceRestart(client *golangsdk.ServiceClient, instanceID string, parameters []interface{}) (bool, error) {
	configs, err := instances.GetConfigurations(client, instanceID).Extract()
	if err != nil {
		return false, fmt.Errorf("error fetching the instance parameters (%s): %s", instanceID, err)
	}

	for _, parameter := range parameters {
		name := parameter.(map[string]interface{})["name"]
		for _, v := range configs.Parameters {
			if v.Name == name && v.Restart {
				return true, nil
			}
		}
	}
	return false, nil
}

func restartRdsInstance(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	// If parameters which requires restart changed, reboot the instance.
	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.RebootInstance(client, d.Id()).Extract()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error rebooting for RDS instance (%s): %s", d.Id(), err)
	}

	// wait for the instance state to be 'ACTIVE'.
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Refresh:      rdsInstanceStateRefreshFunc(client, d.Id()),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) become active status: %s", d.Id(), err)
	}
	return nil
}

func updateRdsInstanceName(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("name") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/name",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateRdsInstanceNameBodyParams(d),
	})
	if err != nil {
		return fmt.Errorf("error renaming RDS instance (%s): %s", d.Id(), err)
	}
	return nil
}

func buildUpdateRdsInstanceNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name").(string),
	}
	return bodyParams
}

func updateRdsInstanceDescription(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("description") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/alias",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: utils.RemoveNil(buildUpdateRdsInstanceDescriptionBodyParams(d)),
	})
	if err != nil {
		return fmt.Errorf("error modify RDS instance (%s) description: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateRdsInstanceDescriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alias": utils.ValueIgnoreEmpty(d.Get("description").(string)),
	}
	return bodyParams
}

func updateAvailabilityZone(ctx context.Context, cfg *config.Config, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("availability_zone") {
		return nil
	}

	availabilityZone := d.Get("availability_zone").([]interface{})
	if len(availabilityZone) == 1 {
		return errors.New("master node does not support modifying availability zone")
	}

	oldRaws, newRaws := d.GetChange("availability_zone")
	oldAzList := oldRaws.([]interface{})
	newAzList := newRaws.([]interface{})
	if oldAzList[0].(string) != newAzList[0].(string) {
		return errors.New("master node does not support modifying availability zone")
	}

	if len(oldAzList) == 1 && len(newAzList) == 2 {
		return changeSingleToPrimaryStandby(ctx, cfg, d, client, newAzList[1].(string))
	}

	instance, err := GetRdsInstanceByID(client, d.Id())
	if err != nil {
		return fmt.Errorf("error getting RDS instance: %s", err)
	}

	slaveNodeId := utils.PathSearch("nodes[?role=='slave']|[0].id", instance, "").(string)
	err = migrateStandbyNode(ctx, d, client, slaveNodeId, newAzList[1].(string))
	if err != nil {
		return err
	}

	return nil
}

func changeSingleToPrimaryStandby(ctx context.Context, cfg *config.Config, d *schema.ResourceData, client *golangsdk.ServiceClient,
	azCode string) error {
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	err = updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v3/{project_id}/instances/{instance_id}/action",
		httpMethod:           "POST",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     utils.RemoveNil(buildChangeSingleToPrimaryStandbyBodyParams(d, azCode)),
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkJobExpression:   "job_id",
		checkOrderExpression: "order_id",
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error changing instance from Single to Primary/Standby: %s", err)
	}
	return nil
}

func buildChangeSingleToPrimaryStandbyBodyParams(d *schema.ResourceData, azCode string) map[string]interface{} {
	params := map[string]interface{}{
		"az_code_new_node": azCode,
		"dsspool_id":       utils.ValueIgnoreEmpty(d.Get("dss_pool_id").(string)),
	}
	if d.Get("charging_mode").(string) == "prePaid" {
		params["is_auto_pay"] = true
	}
	bodyParams := map[string]interface{}{
		"single_to_ha": params,
	}
	return bodyParams
}

func migrateStandbyNode(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, slaveNodeId,
	azCode string) error {
	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/migrateslave",
		httpMethod:         "POST",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   buildMigrateStandbySlaveBodyParams(slaveNodeId, azCode),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "workflowId",
	})
	if err != nil {
		return fmt.Errorf("error migrating slave node(%s): %s", slaveNodeId, err)
	}
	return nil
}

func buildMigrateStandbySlaveBodyParams(slaveNodeId, azCode string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"nodeId": slaveNodeId,
		"azCode": azCode,
	}
	return bodyParams
}

func updateRdsInstanceFlavor(ctx context.Context, d *schema.ResourceData, cfg *config.Config,
	client *golangsdk.ServiceClient, isSupportAutoPay bool) error {
	if !d.HasChange("flavor") {
		return nil
	}

	instance, err := GetRdsInstanceByID(client, d.Id())
	if err != nil {
		return fmt.Errorf("error getting RDS instance: %s", err)
	}

	// if the instance is changed from single to primary/standby, the flavor is not needed to update
	flavor := utils.PathSearch("flavor_ref", instance, "").(string)
	if v := d.Get("flavor").(string); v == flavor {
		return nil
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	err = updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v3/{project_id}/instances/{instance_id}/action",
		httpMethod:           "POST",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     utils.RemoveNil(buildUpdateInstanceFlavorBodyParams(d, isSupportAutoPay)),
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkJobExpression:   "job_id",
		checkOrderExpression: "order_id",
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error updating instance flavor: %s", err)
	}
	return nil
}

func buildUpdateInstanceFlavorBodyParams(d *schema.ResourceData, isSupportAutoPay bool) map[string]interface{} {
	params := map[string]interface{}{
		"spec_code":   d.Get("flavor").(string),
		"is_auto_pay": true,
	}
	if isSupportAutoPay && d.Get("auto_pay").(string) == "false" {
		params["is_auto_pay"] = false
	}
	bodyParams := map[string]interface{}{
		"resize_flavor": params,
	}
	return bodyParams
}

func updateRdsInstanceVolumeSize(ctx context.Context, d *schema.ResourceData, cfg *config.Config,
	client *golangsdk.ServiceClient) error {
	if !d.HasChange("volume.0.size") {
		return nil
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	err = updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v3/{project_id}/instances/{instance_id}/action",
		httpMethod:           "POST",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     utils.RemoveNil(buildUpdateInstanceVolumeBodyParams(d)),
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkJobExpression:   "job_id",
		checkOrderExpression: "order_id",
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error updating instance volume size: %s", err)
	}
	return nil
}

func buildUpdateInstanceVolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"size":        d.Get("volume").([]interface{})[0].(map[string]interface{})["size"],
		"is_auto_pay": true,
	}
	bodyParams := map[string]interface{}{
		"enlarge_volume": params,
	}
	return bodyParams
}

func updateRdsInstanceBackupStrategy(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("backup_strategy") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/backups/policy",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateInstanceBackupStrategyBodyParams(d),
	})
	if err != nil {
		return fmt.Errorf("error updating instance backup strategy: %s", err)
	}

	return nil
}

func buildUpdateInstanceBackupStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	period := d.Get("backup_strategy.0.period").(string)
	if period == "" {
		period = "1,2,3,4,5,6,7"
	}
	params := map[string]interface{}{
		"keep_days":  d.Get("backup_strategy.0.keep_days"),
		"start_time": d.Get("backup_strategy.0.start_time"),
		"period":     period,
	}
	bodyParams := map[string]interface{}{
		"backup_policy": params,
	}
	return bodyParams
}

func updateRdsInstanceDBPort(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("db.0.port") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:             "v3/{project_id}/instances/{instance_id}/port",
		httpMethod:          "PUT",
		pathParams:          map[string]string{"instance_id": d.Id()},
		updateBodyParams:    buildUpdateInstancePortBodyParams(d),
		isRetry:             true,
		timeout:             schema.TimeoutUpdate,
		checkJobExpression:  "workflowId",
		isWaitInstanceReady: true,
	})
	if err != nil {
		return fmt.Errorf("error updating instance database port: %s", err)
	}

	return nil
}

func buildUpdateInstancePortBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"port": d.Get("db.0.port"),
	}
	return bodyParams
}

func updateRdsInstanceFixedIp(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("fixed_ip") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/ip",
		httpMethod:         "PUT",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   buildUpdateInstanceFixedIpBodyParams(d),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "workflowId",
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) fixed IP: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateInstanceFixedIpBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_ip": d.Get("fixed_ip"),
	}
	return bodyParams
}

func updateRdsInstanceSecurityGroup(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("security_group_id") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/security-group",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateInstanceSecurityGroupBodyParams(d),
		isRetry:          true,
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) security group: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateInstanceSecurityGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"security_group_id": d.Get("security_group_id"),
	}
	return bodyParams
}

func updateRdsInstanceSSLConfig(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("ssl_enable") {
		return nil
	}
	if !isMySQLDatabase(d) {
		return fmt.Errorf("only MySQL database support SSL enable and disable")
	}
	return configRdsInstanceSSL(ctx, d, client)
}

func updateBillingModeToPeriod(ctx context.Context, d *schema.ResourceData, cfg *config.Config, client *golangsdk.ServiceClient) error {
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}

	err = updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v3/{project_id}/instances/{instance_id}/to-period",
		httpMethod:           "POST",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     utils.RemoveNil(buildUpdateBillingModeToPeriodBodyParams(d)),
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkOrderExpression: "order_id",
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) billing mode from post-paid to pre-paid: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateBillingModeToPeriodBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"period_type":     strings.ToUpper(d.Get("period_unit").(string)),
		"period_num":      d.Get("period").(int),
		"auto_pay_policy": "YES",
	}
	if d.Get("auto_renew").(string) == "true" {
		bodyParams["auto_renew_policy"] = "YES"
	}
	return bodyParams
}

func updateConfiguration(ctx context.Context, d *schema.ResourceData, client, clientV31 *golangsdk.ServiceClient,
	instanceID string) (context.Context, error) {
	if !d.HasChange("param_group_id") {
		return ctx, nil
	}
	if _, ok := d.GetOk("param_group_id"); !ok {
		return ctx, nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3.1/{project_id}/configurations/{config_id}/apply",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"config_id": d.Get("param_group_id").(string)},
		updateBodyParams: utils.RemoveNil(buildUpdateInstanceConfigurationBodyParams(d)),
		isRetry:          true,
		timeout:          schema.TimeoutUpdate,
	})
	if err != nil {
		return ctx, fmt.Errorf("error updating instance(%s) configuration: %s ", d.Id(), err)
	}

	// wait 30 seconds for the instance to enter the modified status, or the modification has been completed
	// lintignore:R018
	time.Sleep(30 * time.Second)

	// if parameters is set, it should be modified
	if parameters, ok := d.GetOk("parameters"); ok {
		parametersOpts := buildRdsInstanceParameters(parameters.(*schema.Set))
		err = modifyParameters(ctx, d, client, clientV31, instanceID, &parametersOpts)
		if err != nil {
			return ctx, err
		}
	}

	// Sending configurationChanged to Read to warn users the instance needs a reboot.
	ctx = context.WithValue(ctx, ctxType("configurationChanged"), "true")

	return ctx, nil
}

func buildUpdateInstanceConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_ids": []string{d.Id()},
	}
	return bodyParams
}

func updateRdsParameters(ctx context.Context, d *schema.ResourceData, client, clientV31 *golangsdk.ServiceClient,
	instanceID string) (context.Context, error) {
	if !d.HasChange("parameters") {
		return ctx, nil
	}
	values := make(map[string]string)

	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	change := ns.Difference(os).List()
	if len(change) > 0 {
		for _, v := range change {
			key := v.(map[string]interface{})["name"].(string)
			value := v.(map[string]interface{})["value"].(string)
			values[key] = value
		}

		configOpts := instances.ModifyConfigurationOpts{
			Values: values,
		}
		err := modifyParameters(ctx, d, client, clientV31, instanceID, &configOpts)
		if err != nil {
			return ctx, nil
		}
	}

	// Sending parametersChanged to Read to warn users the instance needs a reboot.
	ctx = context.WithValue(ctx, ctxType("parametersChanged"), "true")

	return ctx, nil
}

func modifyParameters(ctx context.Context, d *schema.ResourceData, client, clientV31 *golangsdk.ServiceClient,
	instanceID string, configOpts *instances.ModifyConfigurationOpts) error {
	modifyApiClient := &clientV31
	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.ModifyConfiguration(*modifyApiClient, instanceID, *configOpts).Extract()
		// if the api is not exists, the v3 client should be used
		if apiNotExists := handleApiNotExistsError(err); apiNotExists {
			modifyApiClient = &client
			_, err = instances.ModifyConfiguration(*modifyApiClient, instanceID, *configOpts).Extract()
		}
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil && !handleTimeoutError(err) {
		return fmt.Errorf("error modifying parameters for RDS instance (%s): %s", instanceID, err)
	}

	return checkParameterUpdateCompleted(ctx, d, client, instanceID, d.Timeout(schema.TimeoutUpdate))
}

func checkParameterUpdateCompleted(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      rdsInstanceParamRefreshFunc(client, d, instanceID),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) parameter to be updated: %s ", instanceID, err)
	}
	return nil
}

func rdsInstanceParamRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		configs, err := instances.GetConfigurations(client, instanceID).Extract()
		if err != nil {
			return nil, "ERROR", err
		}
		for _, parameter := range d.Get("parameters").(*schema.Set).List() {
			name := parameter.(map[string]interface{})["name"]
			value := parameter.(map[string]interface{})["value"]
			for _, v := range configs.Parameters {
				if v.Name == name && v.Value != value {
					return configs, "PENDING", nil
				}
			}
		}
		return configs, "SUCCESS", nil
	}
}

func updateVolumeAutoExpand(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("volume.0.limit_size", "volume.0.trigger_threshold") {
		return nil
	}

	limitSize := d.Get("volume.0.limit_size").(int)
	if limitSize > 0 {
		if err := enableVolumeAutoExpand(ctx, d, client, limitSize); err != nil {
			return err
		}
	} else {
		if err := disableVolumeAutoExpand(ctx, schema.TimeoutUpdate, client, d); err != nil {
			return err
		}
	}
	return nil
}

func updateBinlogRetentionHours(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("binlog_retention_hours") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/binlog/clear-policy",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: utils.RemoveNil(buildUpdateBinlogRetentionHoursBodyParams(d)),
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) binlog retention hours: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateBinlogRetentionHoursBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"binlog_retention_hours": d.Get("binlog_retention_hours"),
	}
	return bodyParams
}

func updateMsdtcHosts(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("msdtc_hosts") {
		return nil
	}
	oldRaws, newRaws := d.GetChange("msdtc_hosts")
	addHosts := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	deleteHosts := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))

	if deleteHosts.Len() > 0 {
		err := doUpdateMsdtcHosts(ctx, d, client, "DELETE", deleteHosts.List())
		if err != nil {
			return err
		}
	}
	if addHosts.Len() > 0 {
		err := doUpdateMsdtcHosts(ctx, d, client, "POST", addHosts.List())
		if err != nil {
			return err
		}
	}

	return nil
}

func doUpdateMsdtcHosts(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, method string,
	hostsRaw []interface{}) error {
	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/msdtc/host",
		httpMethod:         method,
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   utils.RemoveNil(buildMsdtcHostsBodyParams(hostsRaw)),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) MSDTC hosts: %s", d.Id(), err)
	}

	return nil
}

func buildMsdtcHostsBodyParams(hostsRaw []interface{}) map[string]interface{} {
	parameters := make([]map[string]interface{}, len(hostsRaw))
	for i, v := range hostsRaw {
		raw := v.(map[string]interface{})
		parameters[i] = map[string]interface{}{
			"ip":        raw["ip"],
			"host_name": raw["host_name"],
		}
	}
	bodyParams := map[string]interface{}{
		"hosts": parameters,
	}
	return bodyParams
}

func updateTde(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("tde_enabled") {
		return nil
	}

	if !d.Get("tde_enabled").(bool) {
		return fmt.Errorf("TDE cannot be disabled after being enabled")
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/tde",
		httpMethod:         "PUT",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   utils.RemoveNil(buildUpdateTdeBodyParams(d)),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) TDE: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateTdeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"rotate_day":     d.Get("rotate_day"),
		"secret_id":      d.Get("secret_id"),
		"secret_name":    d.Get("secret_name"),
		"secret_version": d.Get("secret_version"),
	}
	return bodyParams
}

func updateReadWritePermissions(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("read_write_permissions") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/readonly-status",
		httpMethod:         "PUT",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   buildUpdateReadWritePermissionsBodyParams(d),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) read write permissions: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateReadWritePermissionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	readonly := false
	if d.Get("read_write_permissions") == "readonly" {
		readonly = true
	}
	bodyParams := map[string]interface{}{
		"readonly": readonly,
	}
	return bodyParams
}

func updateSecondLevelMonitoring(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("seconds_level_monitoring_enabled", "seconds_level_monitoring_interval") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/second-level-monitor",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateSecondLevelMonitoringBodyParams(d),
		isRetry:          true,
		timeout:          schema.TimeoutUpdate,
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) seconds level monitoring: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateSecondLevelMonitoringBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch_option": d.Get("seconds_level_monitoring_enabled"),
		"interval":      d.Get("seconds_level_monitoring_interval"),
	}
	return bodyParams
}

func updatePrivateDNSNamePrefix(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("private_dns_name_prefix") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/modify-dns",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdatePrivateDNSNamePrefixBodyParams(d),
		isRetry:          true,
		timeout:          schema.TimeoutUpdate,
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) private DNS name prefix: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      rdsInstancePrivateDNSNameRefreshFunc(client, d.Id(), d.Get("private_dns_name_prefix").(string)),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        1 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) updating instance private DNS name prefix "+
			"completed: %s", d.Id(), err)
	}
	return nil
}

func buildUpdatePrivateDNSNamePrefixBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dns_name": d.Get("private_dns_name_prefix"),
	}
	return bodyParams
}

func rdsInstancePrivateDNSNameRefreshFunc(client *golangsdk.ServiceClient, instanceID,
	privateDNSNamePrefix string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := GetRdsInstanceByID(client, instanceID)
		if err != nil {
			return nil, "ERROR", err
		}
		instanceId := utils.PathSearch("id", instance, "").(string)
		if instanceId == "" {
			return instance, "DELETED", fmt.Errorf("the instance(%s) has been deleted", instanceID)
		}
		privateDNSNames := utils.PathSearch("private_dns_names", instance, make([]interface{}, 0)).([]interface{})
		if len(privateDNSNames) == 0 {
			return instance, "ERROR", fmt.Errorf("error getting private DNS names of the instance(%s)", instanceID)
		}
		prefix := strings.Split(privateDNSNames[0].(string), ".")[0]
		if privateDNSNamePrefix != prefix {
			return instance, "PENDING", nil
		}
		return instance, "COMPLETED", nil
	}
}

func updateSlowLogShowOriginalStatus(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("slow_log_show_original_status") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:    "v3/{project_id}/instances/{instance_id}/slowlog-sensitization/{status}",
		httpMethod: "PUT",
		pathParams: map[string]string{
			"instance_id": d.Id(),
			"status":      d.Get("slow_log_show_original_status").(string),
		},
		isRetry: true,
		timeout: schema.TimeoutUpdate,
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) slow log show original status: %s", d.Id(), err)
	}
	return nil
}

func updateAutoUpgradeSwitchOption(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("minor_version_auto_upgrade_enabled") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/db-auto-upgrade",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildAutoUpgradeSwitchOptionBodyParams(d),
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) auto upgrade switch option: %s", d.Id(), err)
	}

	return nil
}

func buildAutoUpgradeSwitchOptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch_option": d.Get("minor_version_auto_upgrade_enabled"),
	}
	return bodyParams
}

func updatePowerAction(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, powerAction string) error {
	switch powerAction {
	case "ON":
		if err := startupInstance(ctx, d, client); err != nil {
			return err
		}
	case "OFF":
		if err := shutdownInstance(ctx, d, client); err != nil {
			return err
		}
	case "REBOOT":
		if err := rebootInstance(ctx, d, client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("the value of power_action(%s) is error, it should be in [ON, OFF, BEBOOT]", powerAction)
	}

	return nil
}

func startupInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/action/startup",
		httpMethod:         "POST",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   make(map[string]interface{}),
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error starting instance (%s): %s", d.Id(), err)
	}
	return nil
}

func shutdownInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/action/shutdown",
		httpMethod:         "POST",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   make(map[string]interface{}),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error stopping instance (%s): %s", d.Id(), err)
	}
	return nil
}

func rebootInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/action",
		httpMethod:         "POST",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   buildRebootInstanceBodyParams(),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error rebooting instance (%s): %s", d.Id(), err)
	}
	return nil
}

func buildRebootInstanceBodyParams() map[string]interface{} {
	bodyParams := map[string]interface{}{
		"restart": make(map[string]interface{}),
	}
	return bodyParams
}

func enableVolumeAutoExpand(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, limitSize int) error {
	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/disk-auto-expansion",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: utils.RemoveNil(buildEnableVolumeAutoExpandBodyParams(d, limitSize)),
		isRetry:          true,
		timeout:          schema.TimeoutUpdate,
	})
	if err != nil {
		return fmt.Errorf("error enabling automatic expansion of instance(%s) storage: %s ", d.Id(), err)
	}

	return nil
}

func buildEnableVolumeAutoExpandBodyParams(d *schema.ResourceData, limitSize int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch_option":     true,
		"limit_size":        limitSize,
		"trigger_threshold": d.Get("volume.0.trigger_threshold"),
	}
	return bodyParams
}

func disableVolumeAutoExpand(ctx context.Context, timeout string, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/disk-auto-expansion",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildDisableVolumeAutoExpandBodyParams(),
		isRetry:          true,
		timeout:          timeout,
	})
	if err != nil {
		return fmt.Errorf("error disabling automatic expansion of instance(%s) storage: %s ", d.Id(), err)
	}

	return nil
}

func buildDisableVolumeAutoExpandBodyParams() map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch_option": false,
	}
	return bodyParams
}

func configRdsInstanceSSL(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/ssl",
		httpMethod:         "PUT",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   buildUpdateInstanceSSLBodyParams(d),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "job_id",
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) SSL configuration: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateInstanceSSLBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ssl_option": d.Get("ssl_enable"),
	}
	return bodyParams
}

func checkRDSInstanceJobFinish(client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      rdsInstanceJobRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for RDS instance job (%s) to be completed: %s ", jobID, err)
	}
	return nil
}

func rdsInstanceJobRefreshFunc(client *golangsdk.ServiceClient, jobID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		jobOpts := instances.RDSJobOpts{
			JobID: jobID,
		}
		jobList, err := instances.GetRDSJob(client, jobOpts).Extract()
		if err != nil {
			return nil, "FOUND ERROR", err
		}

		return jobList.Job, jobList.Job.Status, nil
	}
}

func rdsInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := GetRdsInstanceByID(client, instanceID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "", "DELETED", nil
			}
			return nil, "FOUND ERROR", err
		}
		status := utils.PathSearch("status", instance, "").(string)
		if status == "FAILED" {
			return instance, status, fmt.Errorf("the instance status is: %s", status)
		}
		return instance, status, nil
	}
}

func updateRdsRootPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("db.0.password") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/password",
		httpMethod:       "POST",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateInstancePasswordBodyParams(d),
		isRetry:          true,
		timeout:          schema.TimeoutUpdate,
	})
	if err != nil {
		return fmt.Errorf("error updating instance(%s) root password: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateInstancePasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_user_pwd": d.Get("db.0.password"),
	}
	return bodyParams
}

func updateRdsInstanceMaintainWindow(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("maintain_begin", "maintain_end") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/ops-window",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateInstanceMaintainWindowBodyParams(d),
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) maintain window: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateInstanceMaintainWindowBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time": d.Get("maintain_begin").(string),
		"end_time":   d.Get("maintain_end").(string),
	}
	return bodyParams
}

func updateRdsInstanceReplicationMode(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("ha_replication_mode") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/failover/mode",
		httpMethod:         "PUT",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   buildUpdateInstanceReplicationModeBodyParams(d),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "workflowId",
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) replication mode: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateInstanceReplicationModeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"mode": d.Get("ha_replication_mode"),
	}
	return bodyParams
}

func updateRdsInstanceSwitchStrategy(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("switch_strategy") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/failover/strategy",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateInstanceSwitchStrategyBodyParams(d),
		isRetry:          true,
		timeout:          schema.TimeoutUpdate,
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) switch strategy: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateInstanceSwitchStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"repairStrategy": d.Get("switch_strategy"),
	}
	return bodyParams
}

func updateRdsInstanceCollation(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("collation") {
		return nil
	}

	err := updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:            "v3/{project_id}/instances/{instance_id}/collations",
		httpMethod:         "PUT",
		pathParams:         map[string]string{"instance_id": d.Id()},
		updateBodyParams:   buildUpdateInstanceCollationBodyParams(d),
		isRetry:            true,
		timeout:            schema.TimeoutUpdate,
		checkJobExpression: "job_id||jobId",
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) collation: %s", d.Id(), err)
	}

	return nil
}

func buildUpdateInstanceCollationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"collation": d.Get("collation"),
	}
	return bodyParams
}

func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + m["value"].(string))
}
