package taurusdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/auditlog"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/autoscaling"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/backups"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/configurations"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/parameters"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/sqlfilter"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ctxType string

// @API GaussDBforMySQL GET /v3/{project_id}/instances
// @API GaussDBforMySQL GET /v3/{project_id}/configurations
// @API GaussDBforNoSQL GET /v3/{project_id}/dedicated-resources
// @API GaussDBforNoSQL POST /v3/{project_id}/instances
// @API GaussDBforMySQL POST /v3/{project_id}/instance/{instance_id}/audit-log/switch
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/sql-filter/switch
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/backups/policy/update
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/restart
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/ops-window
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/monitor-policy
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/internal-ip
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/port
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/dns
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/dns
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/ssl-option
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/alias
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/tags/action
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/name
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/password
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/action
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/nodes/enlarge
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/nodes/{nodeID}
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/volume/extend
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/proxy
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy/enlarge
// @API GaussDBforMySQL PUT /v3/{project_id}/configurations/{configuration_id}/apply
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/configurations
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/auto-scaling/policy
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/backups/encryption
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/slowlog/modify
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/proxy
// @API GaussDBforMySQL GET /v3/{project_id}/instance/{instance_id}/audit-log/switch-status
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/sql-filter/switch
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/monitor-policy
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/tags
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/configurations
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/auto-scaling/policy
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/backups/encryption
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/database-version
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/slowlog/query
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrat
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceGaussDBInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBInstanceCreate,
		UpdateContext: resourceGaussDBInstanceUpdate,
		ReadContext:   resourceGaussDBInstanceRead,
		DeleteContext: resourceGaussDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
				if d.HasChange("proxy_node_num") {
					mErr := multierror.Append(
						d.SetNewComputed("proxy_address"),
						d.SetNewComputed("proxy_port"),
					)
					return mErr.ErrorOrNil()
				}
				return nil
			},
			config.MergeDefaultTags(),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
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
				Optional: true,
				Computed: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dedicated_resource_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"table_name_case_sensitivity": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"read_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "UTC+08:00",
			},
			"availability_zone_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "single",
			},
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"private_write_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"private_dns_name_prefix": {
				Type:     schema.TypeString,
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
			"seconds_level_monitoring_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"seconds_level_monitoring_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"seconds_level_monitoring_enabled"},
			},
			"ssl_option": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"slow_log_show_original_switch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
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
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"sql_filter_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
				Optional: true,
				Computed: true,
			},
			"auto_scaling": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Required: true,
						},
						"scaling_strategy": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flavor_switch": {
										Type:     schema.TypeString,
										Required: true,
									},
									"read_only_switch": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"monitor_cycle": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"silence_cycle": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"enlarge_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"max_flavor": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"reduce_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"max_read_only_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"read_only_weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"silence_start_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_read_only_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"encryption_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encryption_type": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"encryption_status"},
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"encryption_status"},
			},
			// charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			// make ForceNew false here but do nothing in update method!
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
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
			"auto_pay": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"audit_log_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// only supported in some regions, so it's not shown in the doc
			"tags": common.TagsSchema(),

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_flag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"current_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_kernel_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_read_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Deprecated
			"proxy_flavor": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "use huaweicloud_gaussdb_mysql_proxy instead",
			},
			"proxy_node_num": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				Deprecated: "use huaweicloud_gaussdb_mysql_proxy instead",
			},
			"configuration_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Deprecated",
			},
			"proxy_address": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "use huaweicloud_gaussdb_mysql_proxy instead",
			},
			"proxy_port": {
				Type:       schema.TypeInt,
				Computed:   true,
				Deprecated: "use huaweicloud_gaussdb_mysql_proxy instead",
			},
		},
	}
}

func resourceGaussDBDataStore(d *schema.ResourceData) instances.DataStoreOpt {
	var db instances.DataStoreOpt

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
	} else {
		db.Type = "gaussdb-mysql"
		db.Version = "8.0"
	}
	return db
}

func GaussDBInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		if v.Id == "" {
			return v, "DELETED", nil
		}
		return v, v.Status, nil
	}
}

func resourceGaussDBInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GaussdbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if common.HasFilledOpt(d, "force_import") {
		log.Printf("[DEBUG] the Gaussdb mysql instance force_import is set, try to import it instead of creating")
		listOpts := instances.ListTaurusDBInstanceOpts{
			Name: d.Get("name").(string),
		}
		pages, err := instances.List(client, listOpts).AllPages()
		if err != nil {
			return diag.FromErr(err)
		}

		allInstances, err := instances.ExtractTaurusDBInstances(pages)
		if err != nil {
			return diag.Errorf("unable to retrieve instances: %s ", err)
		}
		if allInstances.TotalCount > 0 {
			instance := allInstances.Instances[0]
			log.Printf("[DEBUG] found existing mysql instance %s with name %s", instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceGaussDBInstanceRead(ctx, d, meta)
		}
	}

	createOpts := instances.CreateTaurusDBOpts{
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor").(string),
		Region:              cfg.GetRegion(d),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		DedicatedResourceId: d.Get("dedicated_resource_id").(string),
		TimeZone:            d.Get("time_zone").(string),
		SlaveCount:          d.Get("read_replicas").(int),
		Mode:                "Cluster",
		DataStore:           resourceGaussDBDataStore(d),
	}

	if d.Get("table_name_case_sensitivity").(bool) {
		lowerCaseTableNames := 0
		createOpts.LowerCaseTableNames = &lowerCaseTableNames
	}

	azMode := d.Get("availability_zone_mode").(string)
	createOpts.AZMode = azMode
	if azMode == "multi" {
		v, exist := d.GetOk("master_availability_zone")
		if !exist {
			return diag.Errorf("missing master_availability_zone in a multi availability zone mode")
		}
		createOpts.MasterAZ = v.(string)
	}

	if common.HasFilledOpt(d, "volume_size") {
		volume := &instances.VolumeOpt{
			Size: d.Get("volume_size").(int),
		}
		createOpts.Volume = volume
	}

	// configuration
	if d.Get("configuration_id") == "" && d.Get("configuration_name") != "" {
		configsList, err := configurations.List(client).Extract()
		if err != nil {
			return diag.Errorf("unable to retrieve configurations: %s", err)
		}
		confName := d.Get("configuration_name").(string)
		for _, conf := range configsList {
			if conf.Name == confName {
				createOpts.ConfigurationId = conf.ID
				break
			}
		}
		if createOpts.ConfigurationId == "" {
			return diag.Errorf("unable to find configuration named %s", confName)
		}
	}

	// dedicated resource
	if d.Get("dedicated_resource_id") == "" && d.Get("dedicated_resource_name") != "" {
		pages, err := instances.ListDeh(client).AllPages()
		if err != nil {
			return diag.Errorf("unable to retrieve dedicated resources: %s", err)
		}
		allResources, err := instances.ExtractDehResources(pages)
		if err != nil {
			return diag.Errorf("unable to extract dedicated resources: %s", err)
		}

		derName := d.Get("dedicated_resource_name").(string)
		for _, der := range allResources.Resources {
			if der.ResourceName == derName {
				createOpts.DedicatedResourceId = der.Id
				break
			}
		}
		if createOpts.DedicatedResourceId == "" {
			return diag.Errorf("unable to find dedicated resource named %s", derName)
		}
	}

	// PrePaid
	if d.Get("charging_mode") == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}

		chargeInfo := &instances.ChargeInfoOpt{
			ChargingMode: d.Get("charging_mode").(string),
			PeriodType:   d.Get("period_unit").(string),
			PeriodNum:    d.Get("period").(int),
			IsAutoRenew:  d.Get("auto_renew").(string),
			IsAutoPay:    common.GetAutoPay(d),
		}
		createOpts.ChargeInfo = chargeInfo
	}

	log.Printf("[DEBUG] create options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating GaussDB instance : %s", err)
	}

	id := instance.Instance.Id
	d.SetId(id)

	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BUILD", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      GaussDBInstanceStateRefreshFunc(client, id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        180 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to become ready: %s",
			id, err)
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second) // lintignore:R018

	// waiting for the instance to become ready again
	// as instance will become BACKING UP state after ACTIVE
	stateConf = &resource.StateChangeConf{
		Pending:      []string{"BUILD", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      GaussDBInstanceStateRefreshFunc(client, id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", id, err)
	}

	// audit-log switch
	if _, ok := d.GetOk("audit_log_enabled"); ok {
		err = switchAuditLog(ctx, client, d, schema.TimeoutCreate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// sql-filter switch
	if _, ok := d.GetOk("sql_filter_enabled"); ok {
		err = switchSQLFilter(ctx, client, d, schema.TimeoutCreate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("backup_strategy"); ok {
		if err = updateInstanceBackupStrategy(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("proxy_flavor"); ok {
		if err = enableInstanceProxy(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if parametersRaw := d.Get("parameters").(*schema.Set); parametersRaw.Len() > 0 {
		if err = initializeParameters(ctx, d, client, parametersRaw.List()); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("private_write_ip"); ok {
		if err = updatePrivateWriteIp(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("port"); ok {
		if err = updatePort(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("private_dns_name_prefix"); ok {
		if err = applyPrivateDNSName(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
		if err = updatePrivateDNSName(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("maintain_begin"); ok {
		if err = updateMaintainWindow(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("seconds_level_monitoring_enabled"); ok {
		if err = updatesSecondsLevelMonitoring(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("ssl_option"); ok && v == "false" {
		if err = updateSslOption(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("slow_log_show_original_switch"); ok {
		if err = updateSlowLogShowOriginalSwitch(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("description"); ok {
		if err = updateDescription(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("auto_scaling"); ok {
		if err = updateAutoScaling(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("encryption_status"); ok && v.(string) == "ON" {
		if err = updateEncryption(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", d.Id(), tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags for GaussDB MySQL instance %s: %s", d.Id(), tagErr)
		}
	}

	return resourceGaussDBInstanceRead(ctx, d, meta)
}

func buildGaussDBMySQLParameters(params []interface{}) parameters.UpdateParametersOpts {
	values := make(map[string]string)
	for _, v := range params {
		key := v.(map[string]interface{})["name"].(string)
		value := v.(map[string]interface{})["value"].(string)
		values[key] = value
	}
	return parameters.UpdateParametersOpts{ParameterValues: values}
}

func initializeParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, parametersRaw []interface{}) error {
	updateOpts := buildGaussDBMySQLParameters(parametersRaw)
	restartRequired, err := modifyParameters(ctx, client, d, schema.TimeoutCreate, &updateOpts)
	if err != nil {
		return err
	}

	if restartRequired {
		return restartGaussDBMySQLInstance(ctx, client, d, schema.TimeoutCreate)
	}
	return nil
}

func restartGaussDBMySQLInstance(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	opts := instances.RestartOpts{}
	// If parameters which requires restart changed, reboot the instance.
	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.Restart(client, d.Id(), opts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error restarting GaussDB MySQL instance (%s): %s", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func resourceGaussDBInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GaussdbV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}
	var mErr *multierror.Error

	instanceID := d.Id()
	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB MySQL instance")
	}
	if instance.Id == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB MySQL instance")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", instance.Name),
		d.Set("status", instance.Status),
		d.Set("mode", instance.Type),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("configuration_id", instance.ConfigurationId),
		d.Set("dedicated_resource_id", instance.DedicatedResourceId),
		d.Set("db_user_name", instance.DbUserName),
		d.Set("time_zone", instance.TimeZone),
		d.Set("availability_zone_mode", instance.AZMode),
		d.Set("master_availability_zone", instance.MasterAZ),
		d.Set("description", instance.Alias),
		d.Set("created_at", instance.Created),
		d.Set("updated_at", instance.Updated),
	)

	maintainWindow := strings.Split(instance.MaintenanceWindow, "-")
	if len(maintainWindow) == 2 {
		mErr = multierror.Append(mErr, d.Set("maintain_begin", maintainWindow[0]))
		mErr = multierror.Append(mErr, d.Set("maintain_end", maintainWindow[1]))
	}

	if instance.ConfigurationId != "" {
		mErr = multierror.Append(mErr, setConfigurationId(d, client, instance.ConfigurationId))
	}

	if instance.DedicatedResourceId != "" {
		mErr = multierror.Append(mErr, setDedicatedResourceId(d, client, instance.DedicatedResourceId))
	}

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		mErr = multierror.Append(mErr, d.Set("port", dbPort))
	}
	if len(instance.PrivateIps) > 0 {
		mErr = multierror.Append(mErr, d.Set("private_write_ip", instance.PrivateIps[0]))
	}
	if len(instance.PrivateDnsNames) > 0 {
		mErr = multierror.Append(mErr, d.Set("private_dns_name_prefix", strings.Split(instance.PrivateDnsNames[0], ".")[0]))
		mErr = multierror.Append(mErr, d.Set("private_dns_name", instance.PrivateDnsNames[0]))
	}

	// set data store
	mErr = multierror.Append(mErr, setDatastore(d, instance.DataStore))
	// set nodes, read_replicas, volume_size, flavor
	mErr = multierror.Append(mErr, setNodes(d, instance.Nodes)...)
	// set backup_strategy
	mErr = multierror.Append(mErr, setBackupStrategy(d, instance.BackupStrategy))
	// set proxy
	mErr = multierror.Append(mErr, setProxy(d, client, instanceID)...)
	// set audit log status
	mErr = multierror.Append(mErr, setAuditLog(d, client, instanceID))
	// set sql filter status
	mErr = multierror.Append(mErr, setSqlFilter(d, client, instanceID))
	// set seconds level monitoring
	mErr = multierror.Append(mErr, setSecondsLevelMonitoring(d, client, instanceID)...)
	// set auto scaling
	mErr = multierror.Append(mErr, setAutoScaling(d, client, instanceID))
	// set backup encryption
	mErr = multierror.Append(mErr, setEncryption(d, client, instanceID))
	// set version
	mErr = multierror.Append(mErr, setVersion(d, client, instanceID)...)
	// set slow log show original
	mErr = multierror.Append(mErr, setSlowLogShowOriginalSwitch(d, client, instanceID))

	// save tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] error fetching tags of GaussDB MySQL instance (%s): %s", d.Id(), err)
	}

	diagErr := setGaussDBMySQLParameters(ctx, d, client)
	resErr := append(diag.FromErr(mErr.ErrorOrNil()), diagErr...)

	return resErr
}

func setConfigurationId(d *schema.ResourceData, client *golangsdk.ServiceClient, configurationId string) error {
	configsList, err := configurations.List(client).Extract()
	if err != nil {
		log.Printf("[WARN] unable to retrieve configurations: %s", err)
		return nil
	}
	for _, conf := range configsList {
		if conf.ID == configurationId {
			return d.Set("configuration_name", conf.Name)
		}
	}
	return nil
}

func setDedicatedResourceId(d *schema.ResourceData, client *golangsdk.ServiceClient, dedicatedResourceId string) error {
	pages, err := instances.ListDeh(client).AllPages()
	if err != nil {
		log.Printf("[WARN] unable to retrieve dedicated resources: %s", err)
		return nil
	}
	allResources, err := instances.ExtractDehResources(pages)
	if err != nil {
		log.Printf("[WARN] unable to extract dedicated resources: %s", err)
		return nil
	}
	for _, der := range allResources.Resources {
		if der.Id == dedicatedResourceId {
			return d.Set("dedicated_resource_name", der.ResourceName)
		}
	}
	return nil
}

func setNodes(d *schema.ResourceData, nodes []instances.Nodes) []error {
	flavor := ""
	slaveCount := 0
	volumeSize := 0
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range nodes {
		node := map[string]interface{}{
			"id":                raw.Id,
			"name":              raw.Name,
			"status":            raw.Status,
			"type":              raw.Type,
			"availability_zone": raw.AvailabilityZone,
		}
		if len(raw.PrivateIps) > 0 {
			node["private_read_ip"] = raw.PrivateIps[0]
		}
		if raw.Volume.Size > 0 {
			volumeSize = raw.Volume.Size
		}
		nodesList = append(nodesList, node)
		if raw.Type == "slave" && (raw.Status == "ACTIVE" || raw.Status == "BACKING UP") {
			slaveCount++
		}
		if flavor == "" {
			flavor = raw.Flavor
		}
	}
	var errs []error
	errs = append(errs, d.Set("nodes", nodesList))
	errs = append(errs, d.Set("read_replicas", slaveCount))
	errs = append(errs, d.Set("volume_size", volumeSize))
	if flavor != "" {
		log.Printf("[DEBUG] node flavor: %s", flavor)
		errs = append(errs, d.Set("flavor", flavor))
	}
	return errs
}

func setDatastore(d *schema.ResourceData, datastore instances.DataStore) error {
	dbList := make([]map[string]interface{}, 1)
	db := map[string]interface{}{
		"version": datastore.Version,
	}
	// normalize engine
	engine := datastore.Type
	if engine == "GaussDB(for MySQL)" {
		engine = "gaussdb-mysql"
	}
	db["engine"] = engine
	dbList[0] = db
	return d.Set("datastore", dbList)
}

func setBackupStrategy(d *schema.ResourceData, strategy instances.BackupStrategy) error {
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": strategy.StartTime,
	}
	if days, err := strconv.Atoi(strategy.KeepDays); err == nil {
		backupStrategy["keep_days"] = days
	}
	backupStrategyList[0] = backupStrategy
	return d.Set("backup_strategy", backupStrategyList)
}

func setProxy(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) []error {
	proxy, err := instances.GetProxy(client, instanceId).Extract()
	if err != nil {
		log.Printf("[WARN] instance %s proxy not enabled: %s", instanceId, err)
		return nil
	}
	var errs []error
	errs = append(errs, d.Set("proxy_flavor", proxy.Flavor))
	errs = append(errs, d.Set("proxy_node_num", proxy.NodeNum))
	errs = append(errs, d.Set("proxy_address", proxy.Address))
	errs = append(errs, d.Set("proxy_port", proxy.Port))
	return errs
}

func setAuditLog(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) error {
	resp, err := auditlog.Get(client, instanceId)
	if err != nil {
		log.Printf("[WARN] query instance %s audit log status failed: %s", instanceId, err)
		return nil
	}
	var status bool
	if resp.SwitchStatus == "ON" {
		status = true
	}
	return d.Set("audit_log_enabled", status)
}

func setSqlFilter(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) error {
	resp, err := sqlfilter.Get(client, instanceId).Extract()
	if err != nil {
		log.Printf("[WARN] query instance %s sql filter status failed: %s", instanceId, err)
		return nil
	}
	var status bool
	if resp.SwitchStatus == "ON" {
		status = true
	}
	return d.Set("sql_filter_enabled", status)
}

func setSecondsLevelMonitoring(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) []error {
	resp, err := instances.GetSecondLevelMonitoring(client, instanceId).Extract()
	if err != nil {
		log.Printf("[WARN] query instance %s seconds level monitoring failed: %s", instanceId, err)
		return nil
	}
	var errs []error
	errs = append(errs, d.Set("seconds_level_monitoring_enabled", resp.MonitorSwitch))
	errs = append(errs, d.Set("seconds_level_monitoring_period", resp.Period))
	return errs
}

func setAutoScaling(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) error {
	resp, err := autoscaling.Get(client, instanceId).Extract()
	if err != nil {
		log.Printf("[WARN] query instance %s auto scaling failed: %s", instanceId, err)
		return nil
	}

	autoScaling := map[string]interface{}{
		"id":     resp.Id,
		"status": resp.Status,
		"scaling_strategy": []interface{}{
			map[string]interface{}{
				"flavor_switch":    resp.ScalingStrategy.FlavorSwitch,
				"read_only_switch": resp.ScalingStrategy.ReadOnlySwitch,
			},
		},
		"monitor_cycle":       resp.MonitorCycle,
		"silence_cycle":       resp.SilenceCycle,
		"enlarge_threshold":   resp.EnlargeThreshold,
		"max_flavor":          resp.MaxFavor,
		"reduce_enabled":      resp.ReduceEnabled,
		"min_flavor":          resp.MinFlavor,
		"silence_start_at":    resp.SilenceStartAt,
		"max_read_only_count": resp.MaxReadOnlyCount,
		"min_read_only_count": resp.MinReadOnlyCount,
		"read_only_weight":    resp.ReadOnlyWeight,
	}
	return d.Set("auto_scaling", []interface{}{autoScaling})
}

func setEncryption(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) error {
	resp, err := backups.GetEncryption(client, instanceId).Extract()
	if err != nil {
		log.Printf("[WARN] query instance %s backup encryption failed: %s", instanceId, err)
		return nil
	}
	return d.Set("encryption_status", strings.ToUpper(resp.EncryptionStatus))
}

func setVersion(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) []error {
	resp, err := instances.GetVersion(client, instanceId).Extract()
	if err != nil {
		log.Printf("[WARN] query instance %s version failed: %s", instanceId, err)
		return nil
	}
	var errs []error
	errs = append(errs, d.Set("upgrade_flag", resp.UpgradeFlag))
	errs = append(errs, d.Set("current_version", resp.Datastore.CurrentVersion))
	errs = append(errs, d.Set("current_kernel_version", resp.Datastore.CurrentKernelVersion))
	return errs
}

func setSlowLogShowOriginalSwitch(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) error {
	resp, err := instances.GetSlowLogShowOriginalSwitch(client, instanceId).Extract()
	if err != nil {
		log.Printf("[WARN] query instance %s slow log show original failed: %s", instanceId, err)
		return nil
	}
	slowLogShowOriginalSwitch, _ := strconv.ParseBool(resp.OpenSlowLogSwitch)
	return d.Set("slow_log_show_original_switch", slowLogShowOriginalSwitch)
}

func setGaussDBMySQLParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) diag.Diagnostics {
	parametersList, err := parameters.List(client, d.Id())
	if err != nil {
		return nil
	}

	var configurationRestart bool
	var params []map[string]interface{}
	rawParameterList := d.Get("parameters").(*schema.Set).List()
	rawParameterMap := make(map[string]bool)
	for _, rawParameter := range rawParameterList {
		rawParameterMap[rawParameter.(map[string]interface{})["name"].(string)] = true
	}
	for _, v := range parametersList {
		if v.RestartRequired {
			configurationRestart = true
		}
		if rawParameterMap[v.Name] {
			p := map[string]interface{}{
				"name":  v.Name,
				"value": v.Value,
			}
			params = append(params, p)
		}
	}

	var diagnostics diag.Diagnostics
	if len(params) > 0 {
		if err = d.Set("parameters", params); err != nil {
			log.Printf("error saving parameters to GaussDB MySQL instance (%s): %s", d.Id(), err)
		}
		if ctx.Value(ctxType("parametersChanged")) == "true" {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Parameters Changed",
				Detail:   "Parameters changed which needs reboot.",
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
func resourceGaussDBInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GaussdbV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s ", err)
	}
	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating bss V2 client: %s", err)
	}

	instanceId := d.Id()

	if d.HasChange("name") {
		if err = updateInstanceName(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password") {
		if err = updateInstancePassword(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("flavor") {
		if err = updateInstanceFlavor(ctx, client, bssClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("read_replicas") {
		if err = updateInstanceReadReplica(ctx, client, bssClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("volume_size") {
		if err = updateInstanceVolumeSize(ctx, client, bssClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backup_strategy") {
		if err = updateInstanceBackupStrategy(ctx, client, d, schema.TimeoutUpdate); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("proxy_flavor") {
		if err = updateInstanceProxyFlavor(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("proxy_node_num") {
		if err = updateInstanceProxyNodeNum(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("audit_log_enabled") {
		err = switchAuditLog(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("sql_filter_enabled") {
		err = switchSQLFilter(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("configuration_id") {
		ctx, err = updateConfiguration(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}

		// if parameters is set, it should be modified
		if params, ok := d.GetOk("parameters"); ok {
			updateOpts := buildGaussDBMySQLParameters(params.(*schema.Set).List())
			_, err = modifyParameters(ctx, client, d, schema.TimeoutUpdate, &updateOpts)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("parameters") && !d.HasChanges("configuration_id") {
		ctx, err = updateRdsParameters(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("private_write_ip") {
		err = updatePrivateWriteIp(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("port") {
		err = updatePort(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("security_group_id") {
		err = updateSecurityGroup(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("private_dns_name_prefix") {
		instance, err := instances.Get(client, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(instance.PrivateDnsNames) == 0 || len(instance.PrivateDnsNames[0]) == 0 {
			err = applyPrivateDNSName(ctx, client, d, schema.TimeoutUpdate)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		err = updatePrivateDNSName(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("maintain_begin", "maintain_end") {
		err = updateMaintainWindow(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("seconds_level_monitoring_enabled", "seconds_level_monitoring_period") {
		err = updatesSecondsLevelMonitoring(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("ssl_option") {
		err = updateSslOption(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("slow_log_show_original_switch") {
		err = updateSlowLogShowOriginalSwitch(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("description") {
		err = updateDescription(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_scaling") {
		err = updateAutoScaling(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("encryption_status", "encryption_type", "kms_key_id") {
		err = updateEncryption(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", instanceId)
		if tagErr != nil {
			return diag.Errorf("error updating tags of Gaussdb mysql instance %q: %s", instanceId, tagErr)
		}
	}

	if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", instanceId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "gaussdb",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGaussDBInstanceRead(ctx, d, meta)
}

func resourceGaussDBInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GaussdbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s ", err)
	}

	instanceId := d.Id()
	if d.Get("charging_mode") == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{instanceId}); err != nil {
			// try to delete the instance directly if unsubscribing failed
			res := instances.Delete(client, instanceId)
			if res.Err != nil {
				return common.CheckDeletedDiag(d, res.Err, "GaussDB instance")
			}
		}
	} else {
		result := instances.Delete(client, instanceId)
		if result.Err != nil {
			return common.CheckDeletedDiag(d, result.Err, "GaussDB instance")
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "BACKING UP", "FAILED"},
		Target:     []string{"DELETED"},
		Refresh:    GaussDBInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to be deleted: %s ", instanceId, err)
	}
	log.Printf("[DEBUG] successfully deleted instance %s", instanceId)
	return nil
}

func updateInstanceName(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	newName := d.Get("name").(string)
	updateNameOpts := instances.UpdateNameOpts{
		Name: newName,
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.UpdateName(client, d.Id(), updateNameOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating name for instance %s: %s ", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(schema.TimeoutUpdate))
}

func updateInstancePassword(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	newPass := d.Get("password").(string)
	updatePassOpts := instances.UpdatePassOpts{
		Password: newPass,
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.UpdatePass(client, d.Id(), updatePassOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating password for instance %s: %s ", d.Id(), err)
	}
	return nil
}

func updateInstanceFlavor(ctx context.Context, client, bssClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	newFlavor := d.Get("flavor").(string)
	resizeOpts := instances.ResizeOpts{
		Resize: instances.ResizeOpt{
			Spec: newFlavor,
		},
	}
	if d.Get("charging_mode") == "prePaid" {
		resizeOpts.IsAutoPay = common.GetAutoPay(d)
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.Resize(client, d.Id(), resizeOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating flavor for instance %s: %s ", d.Id(), err)
	}
	job := r.(*instances.JobResponse)

	// wait for job success
	if job.JobID != "" {
		if err = checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	// wait for order success
	if job.OrderID != "" {
		waitTime := int(d.Timeout(schema.TimeoutUpdate) / time.Second)
		if err = orders.WaitForOrderSuccess(bssClient, waitTime, job.OrderID); err != nil {
			return err
		}
		// check whether the order take effect
		instance, err := instances.Get(client, d.Id()).Extract()
		if err != nil {
			return err
		}
		currFlavor := ""
		for _, raw := range instance.Nodes {
			if currFlavor == "" {
				currFlavor = raw.Flavor
				break
			}
		}
		if currFlavor != newFlavor {
			return fmt.Errorf("error updating flavor for instance %s: order failed", d.Id())
		}
	}
	return nil
}

func updateInstanceReadReplica(ctx context.Context, client, bssClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldNum, newNum := d.GetChange("read_replicas")
	if newNum.(int) > oldNum.(int) {
		if err := createInstanceReadReplica(ctx, client, bssClient, d, newNum.(int), oldNum.(int)); err != nil {
			return err
		}
	}
	if newNum.(int) < oldNum.(int) {
		if err := deleteInstanceReadReplica(ctx, client, bssClient, d, newNum.(int), oldNum.(int)); err != nil {
			return err
		}
	}
	return nil
}

func createInstanceReadReplica(ctx context.Context, client, bssClient *golangsdk.ServiceClient, d *schema.ResourceData,
	newNum, oldNum int) error {
	expandSize := newNum - oldNum
	priorities := make([]int, 0)
	for i := 0; i < expandSize; i++ {
		priorities = append(priorities, 1)
	}
	createReplicaOpts := instances.CreateReplicaOpts{
		Priorities: priorities,
	}
	if d.Get("charging_mode") == "prePaid" {
		createReplicaOpts.IsAutoPay = common.GetAutoPay(d)
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.CreateReplica(client, d.Id(), createReplicaOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error creating read replicas for instance %s: %s ", d.Id(), err)
	}
	job := r.(*instances.JobResponse)

	// wait for job success
	if job.JobID != "" {
		jobList := strings.Split(job.JobID, ",")
		log.Printf("[DEBUG] create replica jobs: %#v", jobList)
		for i := 0; i < len(jobList); i++ {
			jobId := jobList[i]
			log.Printf("[DEBUG] waiting for job: %s", jobId)
			if err = checkGaussDBMySQLJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
	}
	// wait for order success
	if job.OrderID != "" {
		waitTime := int(d.Timeout(schema.TimeoutUpdate) / time.Second)
		if err = orders.WaitForOrderSuccess(bssClient, waitTime, job.OrderID); err != nil {
			return err
		}
		// check whether the order take effect
		instance, err := instances.Get(client, d.Id()).Extract()
		if err != nil {
			return err
		}
		slaveCount := 0
		for _, raw := range instance.Nodes {
			if raw.Type == "slave" && (raw.Status == "ACTIVE" || raw.Status == "BACKING UP") {
				slaveCount++
			}
		}
		if newNum != slaveCount {
			return fmt.Errorf("error updating read_replicas for instance %s: order failed", d.Id())
		}
	}
	return nil
}

func deleteInstanceReadReplica(ctx context.Context, client, bssClient *golangsdk.ServiceClient, d *schema.ResourceData,
	newNum, oldNum int) error {
	shrinkSize := oldNum - newNum
	slaveNodes := make([]string, 0)
	nodes := d.Get("nodes").([]interface{})
	for _, nodeRaw := range nodes {
		node := nodeRaw.(map[string]interface{})
		if node["type"].(string) == "slave" && node["status"] == "ACTIVE" {
			slaveNodes = append(slaveNodes, node["id"].(string))
		}
	}
	log.Printf("[DEBUG] Slave Nodes: %+v", slaveNodes)
	if len(slaveNodes) <= shrinkSize {
		return fmt.Errorf("error deleting read replicas for instance %s: Shrink Size is bigger than active slave nodes", d.Id())
	}
	for i := 0; i < shrinkSize; i++ {
		retryFunc := func() (interface{}, bool, error) {
			res, err := instances.DeleteReplica(client, d.Id(), slaveNodes[i]).ExtractJobResponse()
			retry, err := handleMultiOperationsError(err)
			return res, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error deleting read replicas for instance %s: %s ", d.Id(), err)
		}
		job := r.(*instances.JobResponse)
		// wait for job success
		if job.JobID != "" {
			if err = checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
		// wait for order success
		if job.OrderID != "" {
			waitTime := int(d.Timeout(schema.TimeoutUpdate) / time.Second)
			if err = orders.WaitForOrderSuccess(bssClient, waitTime, job.OrderID); err != nil {
				return err
			}
		}
	}
	// check whether the order take effect
	instance, err := instances.Get(client, d.Id()).Extract()
	if err != nil {
		return err
	}
	slaveCount := 0
	for _, raw := range instance.Nodes {
		if raw.Type == "slave" && (raw.Status == "ACTIVE" || raw.Status == "BACKING UP") {
			slaveCount++
		}
	}
	if newNum != slaveCount {
		return fmt.Errorf("error updating read_replicas for instance %s: order failed", d.Id())
	}
	return nil
}

func updateInstanceVolumeSize(ctx context.Context, client, bssClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	extendOpts := instances.ExtendVolumeOpts{
		Size:      d.Get("volume_size").(int),
		IsAutoPay: common.GetAutoPay(d),
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.ExtendVolume(client, d.Id(), extendOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error extending volume: %s", err)
	}
	job := r.(*instances.JobResponse)

	// wait for order success
	if job.OrderID != "" {
		waitTime := int(d.Timeout(schema.TimeoutUpdate) / time.Second)
		if err = orders.WaitForOrderSuccess(bssClient, waitTime, job.OrderID); err != nil {
			return err
		}
		// check whether the order take effect
		instance, err := instances.Get(client, d.Id()).Extract()
		if err != nil {
			return err
		}
		volumeSize := 0
		for _, raw := range instance.Nodes {
			if raw.Volume.Size > 0 {
				volumeSize = raw.Volume.Size
				break
			}
		}
		if volumeSize != d.Get("volume_size").(int) {
			return fmt.Errorf("error updating volume for instance %s: order failed", d.Id())
		}
	}
	return nil
}

func updateInstanceBackupStrategy(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	var updateOpts backups.UpdateOpts
	backupRaw := d.Get("backup_strategy").([]interface{})
	rawMap := backupRaw[0].(map[string]interface{})
	keepDays := rawMap["keep_days"].(int)
	updateOpts.KeepDays = &keepDays
	updateOpts.StartTime = rawMap["start_time"].(string)
	// Fixed to "1,2,3,4,5,6,7"
	updateOpts.Period = "1,2,3,4,5,6,7"
	log.Printf("[DEBUG] update backup_strategy: %#v", updateOpts)

	retryFunc := func() (interface{}, bool, error) {
		err := backups.Update(client, d.Id(), updateOpts).ExtractErr()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating backup_strategy: %s", err)
	}
	return nil
}

func updateInstanceProxyFlavor(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	if _, ok := d.GetOk("proxy_flavor"); ok {
		if err := enableInstanceProxy(ctx, client, d, schema.TimeoutUpdate); err != nil {
			return err
		}
	} else {
		if err := deleteInstanceProxy(ctx, client, d); err != nil {
			return err
		}
	}
	return nil
}

func enableInstanceProxy(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	proxyOpts := instances.ProxyOpts{
		Flavor:  d.Get("proxy_flavor").(string),
		NodeNum: d.Get("proxy_node_num").(int),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.EnableProxy(client, d.Id(), proxyOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error enabling proxy: %s", err)
	}
	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func deleteInstanceProxy(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.DeleteProxy(client, d.Id()).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error disabling proxy: %s", err)
	}
	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(schema.TimeoutUpdate))
}

func updateInstanceProxyNodeNum(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldNum, newNum := d.GetChange("proxy_node_num")
	if oldNum.(int) != 0 && newNum.(int) > oldNum.(int) && common.HasFilledOpt(d, "proxy_flavor") {
		enlargeSize := newNum.(int) - oldNum.(int)
		enlargeProxyOpts := instances.EnlargeProxyOpts{
			NodeNum: enlargeSize,
		}
		retryFunc := func() (interface{}, bool, error) {
			res, err := instances.EnlargeProxy(client, d.Id(), enlargeProxyOpts).ExtractJobResponse()
			retry, err := handleMultiOperationsError(err)
			return res, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error enlarging proxy: %s", err)
		}
		job := r.(*instances.JobResponse)
		if err = checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	if newNum.(int) < oldNum.(int) && !d.HasChange("proxy_flavor") {
		return fmt.Errorf("error updating proxy_node_num for instance %s: new num should be greater than old num", d.Id())
	}
	return nil
}

func switchAuditLog(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	var flag string
	if d.Get("audit_log_enabled").(bool) {
		flag = "ON"
	} else {
		flag = "OFF"
	}
	opts := auditlog.UpdateAuditlogOpts{
		SwitchStatus: flag,
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := auditlog.Update(client, d.Id(), opts)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("switch audit log to %q failed: %s", flag, err)
	}

	return nil
}

func switchSQLFilter(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	flag := "OFF"
	if d.Get("sql_filter_enabled").(bool) {
		flag = "ON"
	}
	opts := sqlfilter.UpdateSqlFilterOpts{
		SwitchStatus: flag,
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := sqlfilter.Update(client, d.Id(), opts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("switch SQL filter to %q failed: %s", flag, err)
	}
	job := r.(*sqlfilter.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(schema.TimeoutUpdate))
}

func updateConfiguration(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (context.Context, error) {
	opts := configurations.ApplyOpts{
		InstanceIds: []string{d.Id()},
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err := configurations.Apply(client, d.Get("configuration_id").(string), opts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return ctx, fmt.Errorf("error updating GausdDB MySQL instance configuration: %s ", err)
	}

	// wait 30 seconds for the instance apply configuration completed
	// lintignore:R018
	time.Sleep(30 * time.Second)

	// Sending configurationChanged to Read to warn users the instance needs a reboot.
	ctx = context.WithValue(ctx, ctxType("configurationChanged"), "true")

	return ctx, nil
}

func updateRdsParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (context.Context, error) {
	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	change := ns.Difference(os)
	if change.Len() > 0 {
		updateOpts := buildGaussDBMySQLParameters(change.List())
		restartRequired, err := modifyParameters(ctx, client, d, schema.TimeoutUpdate, &updateOpts)
		if err != nil {
			return ctx, err
		}
		if restartRequired {
			// Sending parametersChanged to Read to warn users the instance needs a reboot.
			ctx = context.WithValue(ctx, ctxType("parametersChanged"), "true")
		}
	}

	return ctx, nil
}

func modifyParameters(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string,
	parameterOpts *parameters.UpdateParametersOpts) (bool, error) {
	retryFunc := func() (interface{}, bool, error) {
		res, err := parameters.Update(client, d.Id(), *parameterOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return false, fmt.Errorf("error modifying parameters for GaussDB MySQL instance (%s): %s", d.Id(), err)
	}
	job := r.(*parameters.JobResponse)
	err = checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
	if err != nil {
		return false, err
	}
	return job.RestartRequired, nil
}

func updatePrivateWriteIp(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	updatePrivateWriteIpOpts := instances.UpdatePrivateIpOpts{
		InternalIp: d.Get("private_write_ip").(string),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.UpdatePrivateIp(client, d.Id(), updatePrivateWriteIpOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating private write IP for instance %s: %s ", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func updatePort(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	updatePortOpts := instances.UpdatePortOpts{
		Port: d.Get("port").(int),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.UpdatePort(client, d.Id(), updatePortOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating port for instance %s: %s ", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func updateSecurityGroup(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	updateSecurityGroupOpts := instances.UpdateSecurityGroupOpts{
		SecurityGroupId: d.Get("security_group_id").(string),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.UpdateSecurityGroup(client, d.Id(), updateSecurityGroupOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 2 * time.Second,
		PollInterval: 2 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating security group for instance %s: %s ", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func applyPrivateDNSName(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	opts := instances.ApplyPrivateDnsNameOpts{
		DnsType: "private",
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.ApplyPrivateDnsName(client, d.Id(), opts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 30 * time.Second,
		PollInterval: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error applying private DNS name for instance %s: %s ", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func updatePrivateDNSName(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	opts := instances.UpdatePrivateDnsNameOpts{
		DnsName: d.Get("private_dns_name_prefix").(string),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.UpdatePrivateDnsName(client, d.Id(), opts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 30 * time.Second,
		PollInterval: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating private DNS name for instance %s: %s ", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func updateMaintainWindow(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateMaintenanceWindowOpts := instances.UpdateMaintenanceWindowOpts{
		StartTime: d.Get("maintain_begin").(string),
		EndTime:   d.Get("maintain_end").(string),
	}

	_, err := instances.UpdateMaintenanceWindow(client, d.Id(), updateMaintenanceWindowOpts).ExtractUpdateMaintenanceWindowResponse()
	if err != nil {
		return fmt.Errorf("error updating maintenance window for instance %s: %s ", d.Id(), err)
	}

	return nil
}

func updatesSecondsLevelMonitoring(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout string) error {
	opts := instances.UpdateSecondLevelMonitoringOpts{
		MonitorSwitch: d.Get("seconds_level_monitoring_enabled").(bool),
		Period:        d.Get("seconds_level_monitoring_period").(int),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.UpdateSecondLevelMonitoring(client, d.Id(), opts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 30 * time.Second,
		PollInterval: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating seconds level monitoring for instance %s: %s ", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func updateSslOption(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	sslOption, _ := strconv.ParseBool(d.Get("ssl_option").(string))
	updateSslOptionOpts := instances.UpdateSslOptionOpts{
		SslOption: sslOption,
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.UpdateSslOption(client, d.Id(), updateSslOptionOpts).ExtractJobResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 30 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating ssl option for instance %s: %s ", d.Id(), err)
	}

	job := r.(*instances.JobResponse)
	return checkGaussDBMySQLJobFinish(ctx, client, job.JobID, d.Timeout(timeout))
}

func updateSlowLogShowOriginalSwitch(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	opts := instances.UpdateSlowLogShowOriginalSwitchOpts{
		OpenSlowLogSwitch: d.Get("slow_log_show_original_switch").(bool),
	}

	_, err := instances.UpdateSlowLogShowOriginalSwitch(client, d.Id(), opts).ExtractUpdateSlowLogShowOriginalSwitchResponse()
	if err != nil {
		return fmt.Errorf("error updating slow low show original switch for instance %s: %s ", d.Id(), err)
	}

	return nil
}

func updateDescription(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateAliasOpts := instances.UpdateAliasOpts{
		Alias: d.Get("description").(string),
	}

	_, err := instances.UpdateAlias(client, d.Id(), updateAliasOpts).ExtractUpdateAliasResponse()
	if err != nil {
		return fmt.Errorf("error updating description for instance %s: %s ", d.Id(), err)
	}

	return nil
}

func updateAutoScaling(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	rawParams := d.Get("auto_scaling").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	updateAutoScalingOpts := autoscaling.UpdateAutoScalingOpts{
		Status: raw["status"].(string),
		ScalingStrategy: &autoscaling.ScalingStrategy{
			FlavorSwitch:   raw["scaling_strategy"].([]interface{})[0].(map[string]interface{})["flavor_switch"].(string),
			ReadOnlySwitch: raw["scaling_strategy"].([]interface{})[0].(map[string]interface{})["read_only_switch"].(string),
		},
		MonitorCycle:     raw["monitor_cycle"].(int),
		SilenceCycle:     raw["silence_cycle"].(int),
		EnlargeThreshold: raw["enlarge_threshold"].(int),
		MaxFlavor:        raw["max_flavor"].(string),
		ReduceEnabled:    raw["reduce_enabled"].(bool),
		MaxReadOnlyCount: raw["max_read_only_count"].(int),
		ReadOnlyWeight:   raw["read_only_weight"].(int),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := autoscaling.Update(client, d.Id(), updateAutoScalingOpts).ExtractUpdateResponse()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 30 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating auto scaling for instance %s: %s ", d.Id(), err)
	}

	return nil
}

func updateEncryption(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateEncryptionOpts := backups.UpdateEncryptionOpts{
		EncryptionStatus: d.Get("encryption_status").(string),
		Type:             d.Get("encryption_type").(string),
		KmsKeyId:         d.Get("kms_key_id").(string),
	}

	_, err := backups.UpdateEncryption(client, d.Id(), updateEncryptionOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating backup encryption for instance %s: %s ", d.Id(), err)
	}

	return nil
}

func checkGaussDBMySQLJobFinish(ctx context.Context, client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending", "Running"},
		Target:       []string{"Completed"},
		Refresh:      gaussDBMysqlDatabaseStatusRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for GaussDB MySQL instance job (%s) to be completed: %s ", jobID, err)
	}
	return nil
}

func gaussDBMysqlDatabaseStatusRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getJobStatusHttpUrl = "v3/{project_id}/jobs?id={job_id}"
		)

		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", client.ProjectID)
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{job_id}", jobId)

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			return nil, "Failed", err
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("job.status", getJobStatusRespBody, "")
		return getJobStatusRespBody, status.(string), nil
	}
}
