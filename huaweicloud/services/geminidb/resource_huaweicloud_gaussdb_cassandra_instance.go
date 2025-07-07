package geminidb

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/geminidb/v3/backups"
	"github.com/chnsz/golangsdk/openstack/geminidb/v3/configurations"
	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type defaultValues struct {
	Mode      string
	dbType    string
	dbVersion string
	logName   string
}

// @API GaussDBforNoSQL GET /v3/{project_id}/instances
// @API GaussDBforNoSQL GET /v3/{project_id}/dedicated-resources
// @API GaussDBforNoSQL POST /v3/{project_id}/instances
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/tags/action
// @API GaussDBforNoSQL GET /v3/{project_id}/instances/{instance_id}/tags
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/name
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/password
// @API GaussDBforNoSQL PUT /v3/{project_id}/configurations/{config_id}/apply
// @API GaussDBforNoSQL GET /v3/{project_id}/configurations/{config_id}
// @API GaussDBforNoSQL GET /v3/{project_id}/instances/{instance_id}/configurations
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/extend-volume
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/enlarge-node
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/reduce-node
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/resize
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/security-group
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/backups/policy
// @API GaussDBforNoSQL DELETE /v3/{project_id}/instances/{instance_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceGeminiDBInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBCassandraInstanceCreate,
		ReadContext:   resourceGeminiDBInstanceV3Read,
		UpdateContext: resourceGaussDBCassandraInstanceUpdate,
		DeleteContext: resourceGeminiDBInstanceV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"volume_size": {
				Type:     schema.TypeInt,
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
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
								"cassandra", "GeminiDB-Cassandra",
							}, true),
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == "GeminiDB-Cassandra"
							},
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"rocksDB",
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
							Computed: true,
						},
					},
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_port": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_reduce": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			// charge info: charging_mode, period_unit, period, auto_renew
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

			"tags": common.TagsSchema(),
		},
	}
}

func resourceGaussDBCassandraInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaults := defaultValues{
		Mode:      "Cluster",
		dbType:    "cassandra",
		dbVersion: "3.11",
		logName:   "cassandra",
	}
	return resourceGeminiDBInstanceV3Create(ctx, d, meta, defaults)
}

func resourceGeminiDBDataStore(d *schema.ResourceData, defaults defaultValues) instances.DataStore {
	var db instances.DataStore

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
		db.StorageEngine = datastore["storage_engine"].(string)
	} else {
		db.Type = defaults.dbType
		db.Version = defaults.dbVersion
		db.StorageEngine = "rocksDB"
	}
	return db
}

func resourceGeminiDBBackupStrategy(d *schema.ResourceData) *instances.BackupStrategyOpt {
	if _, ok := d.GetOk("backup_strategy"); ok {
		opt := &instances.BackupStrategyOpt{
			StartTime: d.Get("backup_strategy.0.start_time").(string),
		}
		// The default value of keepdays is 7, but empty value of keepdays will be converted to 0.
		if v, ok := d.GetOk("backup_strategy.0.keep_days"); ok {
			opt.KeepDays = strconv.Itoa(v.(int))
		}
		return opt
	}
	return nil
}

func resourceGeminiDBFlavor(d *schema.ResourceData) []instances.FlavorOpt {
	var flavorList []instances.FlavorOpt
	flavor := instances.FlavorOpt{
		Num:      strconv.Itoa(d.Get("node_num").(int)),
		Size:     d.Get("volume_size").(int),
		Storage:  "ULTRAHIGH",
		SpecCode: d.Get("flavor").(string),
	}
	flavorList = append(flavorList, flavor)
	return flavorList
}

func GeminiDBInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.GetInstanceByID(client, instanceID)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "deleted", nil
			}
			return nil, "", err
		}

		return instance, instance.Status, nil
	}
}

func resourceGeminiDBInstanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}, defaults defaultValues) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GeminiDBV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if common.HasFilledOpt(d, "force_import") {
		log.Printf("[DEBUG] Gaussdb %s instance force_import is set, try to import it instead of creating", defaults.logName)
		listOpts := instances.ListGeminiDBInstanceOpts{
			Name: d.Get("name").(string),
		}
		pages, err := instances.List(client, listOpts).AllPages()
		if err != nil {
			return diag.FromErr(err)
		}

		allInstances, err := instances.ExtractGeminiDBInstances(pages)
		if err != nil {
			return diag.Errorf("unable to retrieve instances: %s ", err)
		}
		if allInstances.TotalCount > 0 {
			instance := allInstances.Instances[0]
			log.Printf("[DEBUG] found existing %s instance %s with name %s", defaults.logName, instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceGeminiDBInstanceV3Read(ctx, d, meta)
		}
	}

	createOpts := instances.CreateGeminiDBOpts{
		Name:                d.Get("name").(string),
		Region:              cfg.GetRegion(d),
		AvailabilityZone:    d.Get("availability_zone").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		DedicatedResourceId: d.Get("dedicated_resource_id").(string),
		Mode:                defaults.Mode,
		Flavor:              resourceGeminiDBFlavor(d),
		DataStore:           resourceGeminiDBDataStore(d, defaults),
		BackupStrategy:      resourceGeminiDBBackupStrategy(d),
	}
	if ssl := d.Get("ssl").(bool); ssl {
		createOpts.Ssl = "1"
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
			IsAutoPay:    "true",
			IsAutoRenew:  d.Get("auto_renew").(string),
		}
		createOpts.ChargeInfo = chargeInfo
	}
	log.Printf("[DEBUG] create options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating GeminiDB instance : %s", err)
	}

	d.SetId(instance.Id)
	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"creating"},
		Target:       []string{"normal"},
		Refresh:      GeminiDBInstanceStateRefreshFunc(client, instance.Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to become ready: %s",
			instance.Id, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of GeminiDB %s: %s", d.Id(), tagErr)
		}
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second) // lintignore:R018

	return resourceGeminiDBInstanceV3Read(ctx, d, meta)
}

func resourceGeminiDBInstanceV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GeminiDBV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.GetInstanceByID(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "GeminiDB")
	}
	if instance.Id == "" {
		d.SetId("")
		log.Printf("[WARN] failed to fetch GeminiDB instance: deleted")
		return nil
	}

	log.Printf("[DEBUG] retrieved instance %s: %#v", instanceID, instance)

	mErr := multierror.Append(
		d.Set("name", instance.Name),
		d.Set("region", instance.Region),
		d.Set("status", instance.Status),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("dedicated_resource_id", instance.DedicatedResourceId),
		d.Set("mode", instance.Mode),
		d.Set("db_user_name", instance.DbUserName),
		d.Set("lb_ip_address", instance.LbIpAddress),
		d.Set("lb_port", instance.LbPort),
	)

	if instance.DedicatedResourceId != "" {
		pages, err := instances.ListDeh(client).AllPages()
		if err != nil {
			log.Printf("[DEBUG] unable to retrieve dedicated resources: %s", err)
		} else {
			allResources, err := instances.ExtractDehResources(pages)
			if err != nil {
				log.Printf("[DEBUG] unable to extract dedicated resources: %s", err)
			} else {
				for _, der := range allResources.Resources {
					if der.Id == instance.DedicatedResourceId {
						mErr = multierror.Append(mErr, d.Set("dedicated_resource_name", der.ResourceName))
						break
					}
				}
			}
		}
	}

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		mErr = multierror.Append(mErr, d.Set("port", dbPort))
	}

	dbList := make([]map[string]interface{}, 0, 1)
	db := map[string]interface{}{
		"engine":         instance.DataStore.Type,
		"version":        instance.DataStore.Version,
		"storage_engine": instance.Engine,
	}
	dbList = append(dbList, db)
	mErr = multierror.Append(mErr, d.Set("datastore", dbList))

	specCode := ""
	wrongFlavor := "Inconsistent Flavor"
	ipsList := []string{}
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, group := range instance.Groups {
		for _, Node := range group.Nodes {
			node := map[string]interface{}{
				"id":             Node.Id,
				"name":           Node.Name,
				"status":         Node.Status,
				"private_ip":     Node.PrivateIp,
				"support_reduce": Node.SupportReduce,
			}
			if specCode == "" {
				specCode = Node.SpecCode
			} else if specCode != Node.SpecCode && specCode != wrongFlavor {
				specCode = wrongFlavor
			}
			nodesList = append(nodesList, node)
			// Only return Node private ips which doesn't support reduce
			if !Node.SupportReduce {
				ipsList = append(ipsList, Node.PrivateIp)
			}
		}
		if volSize, err := strconv.Atoi(group.Volume.Size); err == nil {
			mErr = multierror.Append(mErr, d.Set("volume_size", volSize))
		}
		if specCode != "" {
			log.Printf("[DEBUG] node specCode: %s", specCode)
			mErr = multierror.Append(mErr, d.Set("flavor", specCode))
		}
	}
	mErr = multierror.Append(
		mErr,
		d.Set("nodes", nodesList),
		d.Set("private_ips", ipsList),
		d.Set("node_num", len(nodesList)),
	)

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	mErr = multierror.Append(mErr, d.Set("backup_strategy", backupStrategyList))

	// save geminidb tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for geminidb (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of geminidb (%s): %s", d.Id(), err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGeminiDBInstanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GeminiDBV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	instanceId := d.Id()
	if d.Get("charging_mode") == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{instanceId}); err != nil {
			// Try to delete resource directly when unsubscrbing failed
			res := instances.Delete(client, instanceId)
			if res.Err != nil {
				return diag.FromErr(res.Err)
			}
		}
	} else {
		result := instances.Delete(client, instanceId)
		if result.Err != nil {
			return diag.FromErr(result.Err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"normal", "abnormal", "creating", "createfail", "enlargefail", "data_disk_full"},
		Target:       []string{"deleted"},
		Refresh:      GeminiDBInstanceStateRefreshFunc(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        15 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}

	return nil
}

func resourceGaussDBCassandraInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaults := defaultValues{
		Mode:      "Cluster",
		dbType:    "cassandra",
		dbVersion: "3.11",
		logName:   "cassandra",
	}
	return resourceGeminiDBInstanceV3Update(ctx, d, meta, defaults)
}

func resourceGeminiDBInstanceV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}, defaults defaultValues) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GeminiDBV3Client(region)
	if err != nil {
		return diag.Errorf("error creating Vpc: %s", err)
	}
	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating bss V2 client: %s", err)
	}
	// update tags
	instanceId := d.Id()
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", instanceId)
		if tagErr != nil {
			return diag.Errorf("error updating tags of GeminiDB %q: %s", instanceId, tagErr)
		}
	}

	if d.HasChange("name") {
		updateNameOpts := instances.UpdateNameOpts{
			Name: d.Get("name").(string),
		}

		err := instances.UpdateName(client, instanceId, updateNameOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating name for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}
	}

	if d.HasChange("password") {
		updatePassOpts := instances.UpdatePassOpts{
			Password: d.Get("password").(string),
		}

		err := instances.UpdatePass(client, instanceId, updatePassOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating password for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}
	}

	if d.HasChange("configuration_id") {
		instanceIds := []string{d.Id()}
		applyOpts := configurations.ApplyOpts{
			InstanceIds: instanceIds,
		}

		configId := d.Get("configuration_id").(string)
		ret, err := configurations.Apply(client, configId, applyOpts).Extract()
		if err != nil || !ret.Success {
			return diag.Errorf("error updating configuration_id for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"SET_CONFIGURATION"},
			Target:     []string{"available"},
			Refresh:    GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "SET_CONFIGURATION"),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
		}

		// Compare the target configuration and the instance configuration
		cfg, err := configurations.Get(client, configId).Extract()
		if err != nil {
			return diag.Errorf("error fetching configuration %s: %s", configId, err)
		}
		configParams := cfg.Parameters
		log.Printf("[DEBUG] configuration parameters %#v", configParams)

		instanceConfig, err := configurations.GetInstanceConfig(client, instanceId).Extract()
		if err != nil {
			return diag.Errorf("error fetching instance configuration for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}
		instanceConfigParams := instanceConfig.Parameters
		log.Printf("[DEBUG] instance configuration parameters %#v", instanceConfigParams)

		if len(configParams) != len(instanceConfigParams) {
			return diag.Errorf("error updating configuration for instance: %s", instanceId)
		}
		for i := range configParams {
			if !configParams[i].ReadOnly && configParams[i] != instanceConfigParams[i] {
				return diag.Errorf("error updating configuration for instance: %s", instanceId)
			}
		}
	}

	if d.HasChange("volume_size") {
		extendOpts := instances.ExtendVolumeOpts{
			Size: d.Get("volume_size").(int),
		}
		if d.Get("charging_mode") == "prePaid" {
			extendOpts.IsAutoPay = "true"
		}

		n, err := instances.ExtendVolume(client, instanceId, extendOpts).Extract()
		if err != nil {
			return diag.Errorf("error extending gaussdb_%s_instance %s size: %s", defaults.logName, instanceId, err)
		}
		// 1. wait for order success
		if n.OrderId != "" {
			if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.OrderId); err != nil {
				return diag.FromErr(err)
			}
		}

		// 2. wait instance status
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"RESIZE_VOLUME"},
			Target:     []string{"available"},
			Refresh:    GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "RESIZE_VOLUME"),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
		}

		// 3. check whether the order take effect
		if n.OrderId != "" {
			instance, err := instances.GetInstanceByID(client, instanceId)
			if err != nil {
				return diag.FromErr(err)
			}
			volumeSize := 0
			for _, group := range instance.Groups {
				if volSize, err := strconv.Atoi(group.Volume.Size); err == nil {
					volumeSize = volSize
					break
				}
			}
			if volumeSize != d.Get("volume_size").(int) {
				return diag.Errorf("error extending volume for instance %s: order failed", instanceId)
			}
		}
	}

	if d.HasChange("node_num") {
		old, newnum := d.GetChange("node_num")
		if newnum.(int) > old.(int) {
			// Enlarge Nodes
			expandSize := newnum.(int) - old.(int)
			enlargeNodeOpts := instances.EnlargeNodeOpts{
				Num: expandSize,
			}
			if d.Get("charging_mode") == "prePaid" {
				enlargeNodeOpts.IsAutoPay = "true"
			}
			log.Printf("[DEBUG] enlarge node options: %+v", enlargeNodeOpts)

			n, err := instances.EnlargeNode(client, instanceId, enlargeNodeOpts).Extract()
			if err != nil {
				return diag.Errorf("error enlarging gaussdb_%s_instance %s node size: %s", defaults.logName, instanceId, err)
			}
			// 1. wait for order success
			if n.OrderId != "" {
				if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.OrderId); err != nil {
					return diag.FromErr(err)
				}
			}

			// 2. wait instance status
			stateConf := &resource.StateChangeConf{
				Pending:      []string{"GROWING"},
				Target:       []string{"available"},
				Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "GROWING"),
				Timeout:      d.Timeout(schema.TimeoutUpdate),
				Delay:        15 * time.Second,
				PollInterval: 20 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf(
					"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
			}

			// 3. check whether the order take effect
			if n.OrderId != "" {
				instance, err := instances.GetInstanceByID(client, instanceId)
				if err != nil {
					return diag.FromErr(err)
				}
				nodeNum := 0
				for _, group := range instance.Groups {
					nodeNum += len(group.Nodes)
				}
				if nodeNum != newnum.(int) {
					return diag.Errorf("error enlarging node for instance %s: order failed", instanceId)
				}
			}
		}
		if newnum.(int) < old.(int) {
			if defaults.dbType == "influxdb" {
				return diag.Errorf("shrinking gaussdb %s instance node size is not allowed", defaults.logName)
			}
			// Reduce Nodes
			shrinkSize := old.(int) - newnum.(int)
			// the API accepts maxinum num of 10
			reduceNum := 10
			loopSize := shrinkSize / reduceNum
			lastNum := shrinkSize % reduceNum
			if lastNum > 0 {
				loopSize++
			}

			for i := 0; i < loopSize; i++ {
				if lastNum > 0 && (i == loopSize-1) {
					reduceNum = lastNum
				}
				reduceNodeOpts := instances.ReduceNodeOpts{
					Num: reduceNum,
				}
				log.Printf("[DEBUG] reduce node options: %+v", reduceNodeOpts)

				n, err := instances.ReduceNode(client, instanceId, reduceNodeOpts).Extract()
				if err != nil {
					return diag.Errorf("error shrinking gaussdb %s instance %s node size: %s", defaults.logName, instanceId, err)
				}

				// 1. wait for order success
				if n.OrderId != "" {
					if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.OrderId); err != nil {
						return diag.FromErr(err)
					}
				}

				// 2. wait instance status
				stateConf := &resource.StateChangeConf{
					Pending:      []string{"REDUCING"},
					Target:       []string{"available"},
					Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "REDUCING"),
					Timeout:      d.Timeout(schema.TimeoutUpdate),
					Delay:        15 * time.Second,
					PollInterval: 20 * time.Second,
				}

				_, err = stateConf.WaitForStateContext(ctx)
				if err != nil {
					return diag.Errorf(
						"error waiting for gaussdb %s instance %s to become ready: %s", defaults.logName, instanceId, err)
				}
			}
		}
	}

	if d.HasChange("flavor") {
		instance, err := instances.GetInstanceByID(client, instanceId)
		if err != nil {
			return diag.Errorf(
				"error fetching gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}

		specCode := ""
		for _, action := range instance.Actions {
			if action == "RESIZE_FLAVOR" {
				// Wait here if the instance already in RESIZE_FLAVOR state
				stateConf := &resource.StateChangeConf{
					Pending:      []string{"RESIZE_FLAVOR"},
					Target:       []string{"available"},
					Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "RESIZE_FLAVOR"),
					Timeout:      d.Timeout(schema.TimeoutUpdate),
					PollInterval: 20 * time.Second,
				}

				_, err = stateConf.WaitForStateContext(ctx)
				if err != nil {
					return diag.Errorf(
						"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
				}

				instance, err := instances.GetInstanceByID(client, instanceId)
				if err != nil {
					return diag.Errorf(
						"error fetching gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
				}

				// Fetch node flavor
				wrongFlavor := "Inconsistent Flavor"
				for _, group := range instance.Groups {
					for _, Node := range group.Nodes {
						if specCode == "" {
							specCode = Node.SpecCode
						} else if specCode != Node.SpecCode && specCode != wrongFlavor {
							specCode = wrongFlavor
						}
					}
				}
				break
			}
		}

		flavor := d.Get("flavor").(string)
		if specCode != flavor {
			log.Printf("[DEBUG] inconsistent node specCode: %s, flavor: %s", specCode, flavor)
			// Do resize action
			resizeOpts := instances.ResizeOpts{
				Resize: instances.ResizeOpt{
					InstanceID: d.Id(),
					SpecCode:   d.Get("flavor").(string),
				},
			}
			if d.Get("charging_mode") == "prePaid" {
				resizeOpts.IsAutoPay = "true"
			}

			n, err := instances.Resize(client, instanceId, resizeOpts).Extract()
			if err != nil {
				return diag.Errorf("error resizing gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
			}
			// 1. wait for order success
			if n.OrderId != "" {
				if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.OrderId); err != nil {
					return diag.FromErr(err)
				}
			}

			// 2. wait for instance status.
			stateConf := &resource.StateChangeConf{
				Pending:      []string{"RESIZE_FLAVOR"},
				Target:       []string{"available"},
				Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "RESIZE_FLAVOR"),
				Timeout:      d.Timeout(schema.TimeoutUpdate),
				PollInterval: 20 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf(
					"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
			}

			// 3. check whether the order take effect
			if n.OrderId != "" {
				instance, err := instances.GetInstanceByID(client, instanceId)
				if err != nil {
					return diag.FromErr(err)
				}
				currFlavor := ""
				for _, group := range instance.Groups {
					for _, Node := range group.Nodes {
						if currFlavor == "" {
							currFlavor = Node.SpecCode
							break
						}
					}
				}
				if currFlavor != d.Get("flavor").(string) {
					return diag.Errorf("error updating flavor for instance %s: order failed", instanceId)
				}
			}
		}
	}

	if d.HasChange("security_group_id") {
		updateSgOpts := instances.UpdateSgOpts{
			SecurityGroupID: d.Get("security_group_id").(string),
		}

		result := instances.UpdateSg(client, instanceId, updateSgOpts)
		if result.Err != nil {
			return diag.Errorf("error updating security group for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, result.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"MODIFY_SECURITYGROUP"},
			Target:       []string{"available"},
			Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "MODIFY_SECURITYGROUP"),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			PollInterval: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
		}
	}

	if d.HasChange("backup_strategy") {
		var updateOpts backups.UpdateOpts
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keepDays := rawMap["keep_days"].(int)
		updateOpts.KeepDays = &keepDays
		updateOpts.StartTime = rawMap["start_time"].(string)
		// Fixed to "1,2,3,4,5,6,7"
		updateOpts.Period = "1,2,3,4,5,6,7"
		log.Printf("[DEBUG] update backup_strategy: %#v", updateOpts)

		err = backups.Update(client, instanceId, updateOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating backup_strategy: %s", err)
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

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "nosql",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGeminiDBInstanceV3Read(ctx, d, meta)
}

func GeminiDBInstanceUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID, state string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.GetInstanceByID(client, instanceID)

		if err != nil {
			return nil, "", err
		}
		if instance.Id == "" {
			return instance, "deleted", nil
		}
		for _, action := range instance.Actions {
			if state == "REDUCING" {
				if action == "REDUCING" || action == "PERIOD_RESOURCE_DELETE" {
					return instance, state, nil
				}
			} else {
				if action == state {
					return instance, state, nil
				}
			}
		}

		return instance, "available", nil
	}
}
