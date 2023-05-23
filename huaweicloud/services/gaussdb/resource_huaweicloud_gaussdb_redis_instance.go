package gaussdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceGaussRedisInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussRedisInstanceV3Create,
		ReadContext:   resourceGaussRedisInstanceV3Read,
		UpdateContext: resourceGaussRedisInstanceV3Update,
		DeleteContext: resourceGaussRedisInstanceV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(2, 12),
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"force_import": {
				Type:     schema.TypeBool,
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
							ValidateFunc: validation.StringInSlice([]string{
								"redis",
							}, true),
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
							ValidateFunc: validation.StringInSlice([]string{
								"5.0",
							}, true),
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_reduce": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"private_ip": {
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

			"tags": common.TagsSchema(),
		},
	}
}

func resourceGaussRedisDataStore(d *schema.ResourceData) instances.DataStore {
	var db instances.DataStore

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
		db.StorageEngine = datastore["storage_engine"].(string)
	} else {
		db.Type = "redis"
		db.Version = "5.0"
		db.StorageEngine = "rocksDB"
	}
	return db
}

func resourceGaussRedisBackupStrategy(d *schema.ResourceData) *instances.BackupStrategyOpt {
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

func resourceGaussRedisFlavor(d *schema.ResourceData) []instances.FlavorOpt {
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

func GaussRedisInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.GetInstanceByID(client, instanceID)

		if err != nil {
			return nil, "", err
		}
		if instance.Id == "" {
			return instance, "deleted", nil
		}

		return instance, instance.Status, nil
	}
}

func resourceGaussRedisInstanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GeminiDBV3Client(region)
	if err != nil {
		return diag.Errorf("error creating HuaweiCloud GaussDB for Redis client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if common.HasFilledOpt(d, "force_import") {
		logp.Printf("[DEBUG] Gaussdb Redis instance force_import is set, try to import it instead of creating")
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
			logp.Printf("[DEBUG] Found existing redis instance %s with name %s", instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceGaussRedisInstanceV3Read(ctx, d, meta)
		}
	}

	createOpts := instances.CreateGeminiDBOpts{
		Name:                d.Get("name").(string),
		Region:              region,
		AvailabilityZone:    d.Get("availability_zone").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		Mode:                "Cluster",
		Flavor:              resourceGaussRedisFlavor(d),
		DataStore:           resourceGaussRedisDataStore(d),
		BackupStrategy:      resourceGaussRedisBackupStrategy(d),
	}

	// PrePaid
	if d.Get("charging_mode") == "prePaid" {
		if err = common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}

		chargeInfo := &instances.ChargeInfoOpt{
			ChargingMode: d.Get("charging_mode").(string),
			PeriodType:   d.Get("period_unit").(string),
			PeriodNum:    d.Get("period").(int),
			IsAutoPay:    common.GetAutoPay(d),
			IsAutoRenew:  d.Get("auto_renew").(string),
		}
		createOpts.ChargeInfo = chargeInfo
	}
	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating GeminiDB instance : %s", err)
	}

	var delayTime time.Duration = 120
	// 1. wait for order success
	if instance.OrderId != "" {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutCreate)/time.Second),
			instance.OrderId); err != nil {
			return diag.FromErr(err)
		}
		delayTime = 10
	}

	d.SetId(instance.Id)
	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"creating", "BACKUP"},
		Target:       []string{"normal"},
		Refresh:      GaussRedisInstanceStateRefreshFunc(client, instance.Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        delayTime * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", instance.Id, err)
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
	time.Sleep(360 * time.Second) //lintignore:R018

	return resourceGaussRedisInstanceV3Read(ctx, d, meta)
}

func resourceGaussRedisInstanceV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GeminiDBV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating HuaweiCloud GaussRedis client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.GetInstanceByID(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "GaussRedis")
	}
	if instance.Id == "" {
		d.SetId("")
		logp.Printf("[WARN] failed to fetch GausssDB for Redis instance: deleted")
		return nil
	}

	logp.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("name", instance.Name)
	d.Set("region", instance.Region)
	d.Set("status", instance.Status)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("mode", instance.Mode)
	d.Set("db_user_name", instance.DbUserName)
	d.Set("lb_ip_address", instance.LbIpAddress)
	d.Set("lb_port", instance.LbPort)

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		d.Set("port", dbPort)
	}

	dbList := make([]map[string]interface{}, 0, 1)
	db := map[string]interface{}{
		"engine":         instance.DataStore.Type,
		"version":        instance.DataStore.Version,
		"storage_engine": instance.Engine,
	}
	dbList = append(dbList, db)
	d.Set("datastore", dbList)

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
			d.Set("volume_size", volSize)
		}
		if specCode != "" {
			logp.Printf("[DEBUG] Node SpecCode: %s", specCode)
			d.Set("flavor", specCode)
		}
	}
	d.Set("nodes", nodesList)
	d.Set("private_ips", ipsList)
	d.Set("node_num", len(nodesList))

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	d.Set("backup_strategy", backupStrategyList)

	// save geminidb tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err = d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for geminidb (%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] Error fetching tags of geminidb (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceGaussRedisInstanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GeminiDBV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating HuaweiCloud GaussRedis client: %s ", err)
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
		return diag.Errorf("Error waiting for instance (%s) to be deleted: %s ", instanceId, err)
	}
	logp.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	d.SetId("")
	return nil
}

func resourceGaussRedisInstanceV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GeminiDBV3Client(region)
	if err != nil {
		return diag.Errorf("error creating Huaweicloud GaussRedis client: %s", err)
	}
	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating HuaweiCloud bss V2 client: %s", err)
	}
	// update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of GaussDB for Redis %q: %s", d.Id(), tagErr)
		}
	}

	if d.HasChange("name") {
		updateNameOpts := instances.UpdateNameOpts{
			Name: d.Get("name").(string),
		}

		err = instances.UpdateName(client, d.Id(), updateNameOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating name for huaweicloud_gaussdb_redis_instance %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("password") {
		updatePassOpts := instances.UpdatePassOpts{
			Password: d.Get("password").(string),
		}

		err = instances.UpdatePass(client, d.Id(), updatePassOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating password for huaweicloud_gaussdb_redis_instance %s: %s",
				d.Id(), err)
		}
	}

	if d.HasChange("volume_size") {
		err = gaussRedisInstanceUpdateVolumeSize(ctx, d, client, bssClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("flavor") {
		err = gaussRedisInstanceUpdateFlavor(ctx, d, client, bssClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("security_group_id") {
		err = gaussRedisInstanceUpdateSecurityGroup(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("node_num") {
		err = gaussRedisInstanceUpdateNodeNum(ctx, d, client, bssClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", d.Id(), err)
		}
	}

	return resourceGaussRedisInstanceV3Read(ctx, d, meta)
}

func gaussRedisInstanceUpdateVolumeSize(ctx context.Context, d *schema.ResourceData,
	client, bssClient *golangsdk.ServiceClient) error {
	extendOpts := instances.ExtendVolumeOpts{
		Size: d.Get("volume_size").(int),
	}
	if d.Get("charging_mode") == "prePaid" {
		extendOpts.IsAutoPay = common.GetAutoPay(d)
	}

	var res *instances.ExtendResponse
	var err error
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		res, err = instances.ExtendVolume(client, d.Id(), extendOpts).Extract()
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
		return fmt.Errorf("error extending huaweicloud_gaussdb_redis_instance %s size: %s", d.Id(), err)
	}
	// 1. wait for order success
	if res.OrderId != "" {
		if err = orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second),
			res.OrderId); err != nil {
			return err
		}
	}

	// 2. wait instance status
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RESIZE_VOLUME", "PERIOD_RESOURCE_SPEC_CHG"},
		Target:  []string{"available"},
		Refresh: GaussRedisInstanceUpdateRefreshFunc(client, d.Id(),
			map[string]bool{"RESIZE_VOLUME": true, "PERIOD_RESOURCE_SPEC_CHG": true}),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf(
			"error waiting for huaweicloud_gaussdb_redis_instance %s to become ready: %s", d.Id(), err)
	}

	// 3. check whether the order take effect
	if res.OrderId != "" {
		instance, err := instances.GetInstanceByID(client, d.Id())
		if err != nil {
			return err
		}
		volumeSize := 0
		for _, group := range instance.Groups {
			if volSize, err := strconv.Atoi(group.Volume.Size); err == nil {
				volumeSize = volSize
				break
			}
		}
		if volumeSize != d.Get("volume_size").(int) {
			return fmt.Errorf("error extending volume for instance %s: order failed", d.Id())
		}
	}
	return nil
}

func gaussRedisInstanceUpdateSecurityGroup(ctx context.Context, d *schema.ResourceData,
	client *golangsdk.ServiceClient) error {
	updateSgOpts := instances.UpdateSgOpts{
		SecurityGroupID: d.Get("security_group_id").(string),
	}

	var err error
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		res := instances.UpdateSg(client, d.Id(), updateSgOpts)
		isRetry, err := handleOperationError(res.Err)
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf(
			"error updating security group for huaweicloud_gaussdb_redis_instance %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"MODIFY_SECURITYGROUP"},
		Target:       []string{"available"},
		Refresh:      GaussRedisInstanceUpdateRefreshFunc(client, d.Id(), map[string]bool{"MODIFY_SECURITYGROUP": true}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf(
			"error waiting for huaweicloud_gaussdb_redis_instance %s to become ready: %s", d.Id(), err)
	}
	return nil
}

// lintignore:R018
func gaussRedisInstanceUpdateNodeNum(ctx context.Context, d *schema.ResourceData,
	client, bssClient *golangsdk.ServiceClient) error {
	oldNum, newNum := d.GetChange("node_num")
	if newNum.(int) > oldNum.(int) {
		// Enlarge Nodes
		return gaussRedisInstanceEnlargeNodeNum(ctx, d, client, bssClient, newNum.(int), oldNum.(int))
	}
	// Reduce Nodes
	return gaussRedisInstanceReduceNodeNum(ctx, d, client, bssClient, newNum.(int), oldNum.(int))
}

func gaussRedisInstanceEnlargeNodeNum(ctx context.Context, d *schema.ResourceData,
	client, bssClient *golangsdk.ServiceClient, newNum, oldNum int) error {
	expandSize := newNum - oldNum
	enlargeNodeOpts := instances.EnlargeNodeOpts{
		Num: expandSize,
	}
	if d.Get("charging_mode") == "prePaid" {
		enlargeNodeOpts.IsAutoPay = common.GetAutoPay(d)
	}
	logp.Printf("[DEBUG] Enlarge Node Options: %+v", enlargeNodeOpts)

	var res *instances.ExtendResponse
	var err error
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		res, err = instances.EnlargeNode(client, d.Id(), enlargeNodeOpts).Extract()
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
		return fmt.Errorf("error enlarging huaweicloud_redis_cassandra_instance %s node size: %s",
			d.Id(), err)
	}
	// 1. wait for order success
	if res.OrderId != "" {
		if err = orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second),
			res.OrderId); err != nil {
			return err
		}
	}

	// 2. wait instance status
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"GROWING"},
		Target:       []string{"available"},
		Refresh:      GaussRedisInstanceUpdateRefreshFunc(client, d.Id(), map[string]bool{"GROWING": true}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        15 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf(
			"error waiting for huaweicloud_gaussdb_redis_instance %s to become ready: %s", d.Id(), err)
	}

	// 3. check whether the order take effect
	if res.OrderId != "" {
		instance, err := instances.GetInstanceByID(client, d.Id())
		if err != nil {
			return err
		}
		nodeNum := 0
		for _, group := range instance.Groups {
			nodeNum += len(group.Nodes)
		}
		if nodeNum != newNum {
			return fmt.Errorf("error enlarging node for instance %s: order failed", d.Id())
		}
	}
	return nil
}

func gaussRedisInstanceReduceNodeNum(ctx context.Context, d *schema.ResourceData,
	client, bssClient *golangsdk.ServiceClient, newNum, oldNum int) error {
	shrinkSize := oldNum - newNum
	reduceNodeOpts := instances.ReduceNodeOpts{
		Num: 1,
	}
	logp.Printf("[DEBUG] Reduce Node Options: %+v", reduceNodeOpts)

	for i := 0; i < shrinkSize; i++ {
		var res *instances.ExtendResponse
		var err error
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			res, err = instances.ReduceNode(client, d.Id(), reduceNodeOpts).Extract()
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
			return fmt.Errorf(
				"error shrinking huaweicloud_gaussdb_redis_instance %s node size: %s", d.Id(), err)
		}

		// wait for order success
		if res.OrderId != "" {
			if err = orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second),
				res.OrderId); err != nil {
				return err
			}
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"REDUCING"},
			Target:       []string{"available"},
			Refresh:      GaussRedisInstanceUpdateRefreshFunc(client, d.Id(), map[string]bool{"REDUCING": true}),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        15 * time.Second,
			PollInterval: 20 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf(
				"error waiting for huaweicloud_gaussdb_redis_instance %s to become ready: %s", d.Id(), err)
		}

		// 3. check whether the order take effect
		if res.OrderId != "" {
			instance, err := instances.GetInstanceByID(client, d.Id())
			if err != nil {
				return err
			}
			nodeNum := 0
			for _, group := range instance.Groups {
				nodeNum += len(group.Nodes)
			}
			if nodeNum != newNum {
				return fmt.Errorf("error enlarging node for instance %s: order failed", d.Id())
			}
		}
	}
	return nil
}

func GaussRedisInstanceUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID string,
	states map[string]bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.GetInstanceByID(client, instanceID)

		if err != nil {
			return nil, "", err
		}
		if instance.Id == "" {
			return instance, "deleted", nil
		}
		for _, action := range instance.Actions {
			if _, ok := states[action]; ok {
				return instance, action, nil
			}
		}

		return instance, "available", nil
	}
}

func gaussRedisInstanceUpdateFlavor(ctx context.Context, d *schema.ResourceData,
	client, bssClient *golangsdk.ServiceClient) error {
	instance, err := instances.GetInstanceByID(client, d.Id())
	if err != nil {
		return fmtp.Errorf("error fetching huaweicloud_gaussdb_redis_instance %s: %s", d.Id(), err)
	}

	// Do resize action
	resizeOpts := instances.ResizeOpts{
		Resize: instances.ResizeOpt{
			InstanceID: d.Id(),
			SpecCode:   d.Get("flavor").(string),
		},
	}
	if d.Get("charging_mode") == "prePaid" {
		resizeOpts.IsAutoPay = common.GetAutoPay(d)
	}

	var res *instances.ExtendResponse
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		res, err = instances.Resize(client, d.Id(), resizeOpts).Extract()
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
		return fmtp.Errorf("Error resizing huaweicloud_gaussdb_redis_instance %s: %s", d.Id(), err)
	}
	// 1. wait for order success
	if res.OrderId != "" {
		if err = orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second),
			res.OrderId); err != nil {
			return err
		}
	}

	// 2. wait for instance status.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RESIZE_FLAVOR"},
		Target:       []string{"available"},
		Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, d.Id(), "RESIZE_FLAVOR"),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf(
			"error waiting for huaweicloud_gaussdb_redis_instance %s to become ready: %s", d.Id(), err)
	}

	// 3. check whether the order take effect
	if res.OrderId == "" {
		return nil
	}

	instance, err = instances.GetInstanceByID(client, d.Id())
	if err != nil {
		return err
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
		return fmtp.Errorf("Error updating flavor for instance %s: order failed", d.Id())
	}
	return nil
}

func handleOperationError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault403); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}
		if errorCode == "DBS.200019" {
			return true, err
		}
	}
	return false, err
}
