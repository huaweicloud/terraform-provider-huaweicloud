package gaussdb

import (
	"context"
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
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/auditlog"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/backups"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/configurations"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/sqlfilter"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL GET /v3/{project_id}/instances
// @API GaussDBforMySQL GET /v3/{project_id}/configurations
// @API GaussDBforNoSQL GET /v3/{project_id}/dedicated-resources
// @API GaussDBforNoSQL POST /v3/{project_id}/instances
// @API GaussDBforMySQL POST /v3/{project_id}/instance/{instance_id}/audit-log/switch
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/sql-filter/switch
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/backups/policy/update
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/tags/action
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/name
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/password
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/action
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/nodes/enlarge
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/nodes/{nodeID}
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/volume/extend
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/backups/policy/update
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/proxy
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy/enlarge
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/proxy
// @API GaussDBforMySQL GET /v3/{project_id}/instance/{instance_id}/audit-log/switch-status
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/sql-filter/switch
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/tags
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}
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

		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("proxy_node_num") {
				mErr := multierror.Append(
					d.SetNewComputed("proxy_address"),
					d.SetNewComputed("proxy_port"),
				)
				return mErr.ErrorOrNil()
			}
			return nil
		},

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
				ForceNew: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.All(validation.IntBetween(40, 128000), validation.IntDivisibleBy(10)),
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
				ValidateFunc: validation.StringInSlice([]string{
					"single", "multi",
				}, true),
			},
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
							ValidateFunc: validation.StringInSlice([]string{
								"gaussdb-mysql",
							}, true),
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
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// only supported in some regions, so it's not shown in the doc
			"tags": common.TagsSchema(),

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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_write_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_user_name": {
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
				ValidateFunc: validation.IntBetween(1, 9),
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
		return diag.Errorf(
			"error waiting for instance (%s) to become ready: %s",
			id, err)
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
		if err = updateInstanceBackupStrategy(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("proxy_flavor"); ok {
		if err = enableInstanceProxy(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of Gaussdb mysql instance %s: %s", d.Id(), tagErr)
		}
	}

	return resourceGaussDBInstanceRead(ctx, d, meta)
}

func resourceGaussDBInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, err, "GaussDB instance")
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] retrieved instance %s: %#v", instanceID, instance)

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
	)

	if instance.ConfigurationId != "" {
		setConfigurationId(d, client, instance.ConfigurationId)
	}

	if instance.DedicatedResourceId != "" {
		setDedicatedResourceId(d, client, instance.DedicatedResourceId)
	}

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		d.Set("port", dbPort)
	}
	if len(instance.PrivateIps) > 0 {
		d.Set("private_write_ip", instance.PrivateIps[0])
	}

	// set data store
	setDatastore(d, instance.DataStore)
	// set nodes, read_replicas, volume_size, flavor
	setNodes(d, instance.Nodes)
	// set backup_strategy
	setBackupStrategy(d, instance.BackupStrategy)
	// set proxy
	setProxy(d, client, instanceID)
	// set audit log status
	setAuditLog(d, client, instanceID)
	// set sql filter status
	res, err := sqlfilter.Get(client, instanceID).Extract()
	if err != nil {
		log.Printf("[DEBUG] query instance %s sql filter status failed: %s", instanceID, err)
	} else {
		d.Set("sql_filter_enabled", res.SwitchStatus == "ON")
	}

	// save tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for Gaussdb mysql instance (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of Gaussdb mysql instance (%s): %s", d.Id(), err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func setConfigurationId(d *schema.ResourceData, client *golangsdk.ServiceClient, configurationId string) {
	configsList, err := configurations.List(client).Extract()
	if err != nil {
		log.Printf("unable to retrieve configurations: %s", err)
		return
	}
	for _, conf := range configsList {
		if conf.ID == configurationId {
			d.Set("configuration_name", conf.Name)
			break
		}
	}
}

func setDedicatedResourceId(d *schema.ResourceData, client *golangsdk.ServiceClient, dedicatedResourceId string) {
	pages, err := instances.ListDeh(client).AllPages()
	if err != nil {
		log.Printf("unable to retrieve dedicated resources: %s", err)
		return
	}
	allResources, err := instances.ExtractDehResources(pages)
	if err != nil {
		log.Printf("unable to extract dedicated resources: %s", err)
		return
	}
	for _, der := range allResources.Resources {
		if der.Id == dedicatedResourceId {
			d.Set("dedicated_resource_name", der.ResourceName)
			break
		}
	}
}

func setNodes(d *schema.ResourceData, nodes []instances.Nodes) {
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
	d.Set("nodes", nodesList)
	d.Set("read_replicas", slaveCount)
	d.Set("volume_size", volumeSize)
	if flavor != "" {
		log.Printf("[DEBUG] node flavor: %s", flavor)
		d.Set("flavor", flavor)
	}
}

func setDatastore(d *schema.ResourceData, datastore instances.DataStore) {
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
	d.Set("datastore", dbList)
}

func setBackupStrategy(d *schema.ResourceData, strategy instances.BackupStrategy) {
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": strategy.StartTime,
	}
	if days, err := strconv.Atoi(strategy.KeepDays); err == nil {
		backupStrategy["keep_days"] = days
	}
	backupStrategyList[0] = backupStrategy
	d.Set("backup_strategy", backupStrategyList)
}

func setProxy(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) {
	proxy, err := instances.GetProxy(client, instanceId).Extract()
	if err != nil {
		log.Printf("[DEBUG] instance %s proxy not enabled: %s", instanceId, err)
		return
	}
	d.Set("proxy_flavor", proxy.Flavor)
	d.Set("proxy_node_num", proxy.NodeNum)
	d.Set("proxy_address", proxy.Address)
	d.Set("proxy_port", proxy.Port)
}

func setAuditLog(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) {
	resp, err := auditlog.Get(client, instanceId)
	if err != nil {
		log.Printf("[DEBUG] query instance %s audit log status failed: %s", instanceId, err)
		return
	}
	var status bool
	if resp.SwitchStatus == "ON" {
		status = true
	}
	d.Set("audit_log_enabled", status)
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
		if err = updateInstanceBackupStrategy(client, d); err != nil {
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
		migrateOpts := enterpriseprojects.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "gaussdb",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := common.MigrateEnterpriseProject(ctx, cfg, d, migrateOpts); err != nil {
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
		return diag.Errorf(
			"error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
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

func updateInstanceBackupStrategy(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var updateOpts backups.UpdateOpts
	backupRaw := d.Get("backup_strategy").([]interface{})
	rawMap := backupRaw[0].(map[string]interface{})
	keepDays := rawMap["keep_days"].(int)
	updateOpts.KeepDays = &keepDays
	updateOpts.StartTime = rawMap["start_time"].(string)
	// Fixed to "1,2,3,4,5,6,7"
	updateOpts.Period = "1,2,3,4,5,6,7"
	log.Printf("[DEBUG] update backup_strategy: %#v", updateOpts)

	err := backups.Update(client, d.Id(), updateOpts).ExtractErr()
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
			OkCodes: []int{
				200,
			},
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
