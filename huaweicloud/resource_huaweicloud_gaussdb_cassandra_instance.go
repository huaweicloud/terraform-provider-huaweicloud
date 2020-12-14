package huaweicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/bss/v2/orders"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/backups"
	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/configurations"
	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/instances"
)

func resourceGeminiDBInstanceV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeminiDBInstanceV3Create,
		Read:   resourceGeminiDBInstanceV3Read,
		Update: resourceGeminiDBInstanceV3Update,
		Delete: resourceGeminiDBInstanceV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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
				ValidateFunc: validation.IntBetween(3, 200),
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
								"GeminiDB-Cassandra",
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
								"3.11",
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
				ForceNew: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, true),
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, true),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"region": {
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
			"tags": tagsSchema(),
		},
	}
}

func resourceGeminiDBDataStore(d *schema.ResourceData) instances.DataStore {
	var db instances.DataStore

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
		db.StorageEngine = datastore["storage_engine"].(string)
	} else {
		db.Type = "GeminiDB-Cassandra"
		db.Version = "3.11"
		db.StorageEngine = "rocksDB"
	}
	return db
}

func resourceGeminiDBBackupStrategy(d *schema.ResourceData) *instances.BackupStrategyOpt {
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 1 {
		strategy := backupStrategyRaw[0].(map[string]interface{})
		return &instances.BackupStrategyOpt{
			StartTime: strategy["start_time"].(string),
			KeepDays:  strconv.Itoa(strategy["keep_days"].(int)),
		}
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
			return nil, "", err
		}
		if instance.Id == "" {
			return instance, "deleted", nil
		}

		return instance, instance.Status, nil
	}
}

func resourceGeminiDBInstanceV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.GeminiDBV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if hasFilledOpt(d, "force_import") {
		log.Printf("[DEBUG] Gaussdb cassandra instance force_import is set, try to import it instead of creating")
		listOpts := instances.ListGeminiDBInstanceOpts{
			Name: d.Get("name").(string),
		}
		pages, err := instances.List(client, listOpts).AllPages()
		if err != nil {
			return err
		}

		allInstances, err := instances.ExtractGeminiDBInstances(pages)
		if err != nil {
			return fmt.Errorf("Unable to retrieve instances: %s ", err)
		}
		if allInstances.TotalCount > 0 {
			instance := allInstances.Instances[0]
			log.Printf("[DEBUG] Found existing cassandra instance %s with name %s", instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceGeminiDBInstanceV3Read(d, meta)
		}
	}

	createOpts := instances.CreateGeminiDBOpts{
		Name:                d.Get("name").(string),
		Region:              GetRegion(d, config),
		AvailabilityZone:    d.Get("availability_zone").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
		Password:            d.Get("password").(string),
		Mode:                "Cluster",
		Flavor:              resourceGeminiDBFlavor(d),
		DataStore:           resourceGeminiDBDataStore(d),
		BackupStrategy:      resourceGeminiDBBackupStrategy(d),
	}
	if ssl := d.Get("ssl").(bool); ssl {
		createOpts.Ssl = "1"
	}

	// PrePaid
	if d.Get("charging_mode") == "prePaid" {
		chargeInfo := &instances.ChargeInfoOpt{
			ChargingMode: d.Get("charging_mode").(string),
			PeriodType:   d.Get("period_unit").(string),
			PeriodNum:    d.Get("period").(int),
			IsAutoPay:    "true",
			IsAutoRenew:  d.Get("auto_renew").(string),
		}
		createOpts.ChargeInfo = chargeInfo
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating GeminiDB instance : %s", err)
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

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			instance.Id, err)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of GeminiDB %s: %s", d.Id(), tagErr)
		}
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second)

	return resourceGeminiDBInstanceV3Read(d, meta)
}

func resourceGeminiDBInstanceV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.GeminiDBV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.GetInstanceByID(client, instanceID)
	if err != nil {
		return CheckDeleted(d, err, "GeminiDB")
	}
	if instance.Id == "" {
		d.SetId("")
		log.Printf("[WARN] failed to fetch GeminiDB instance: deleted")
		return nil
	}

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("name", instance.Name)
	d.Set("region", instance.Region)
	d.Set("status", instance.Status)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("mode", instance.Mode)
	d.Set("db_user_name", instance.DbUserName)

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
			log.Printf("[DEBUG] Node SpecCode: %s", specCode)
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

	//save geminidb tags
	resourceTags, err := tags.Get(client, "instances", d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching HuaweiCloud geminidb tags: %s", err)
	}

	tagmap := tagsToMap(resourceTags.Tags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags for HuaweiCloud geminidb (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceGeminiDBInstanceV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.GeminiDBV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s ", err)
	}

	instanceId := d.Id()
	if d.Get("charging_mode") == "prePaid" {
		bssV2Client, err := config.BssV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud bss V2 client: %s", err)
		}

		resourceIds := []string{instanceId}
		unsubscribeOpts := orders.UnsubscribeOpts{
			ResourceIds:     resourceIds,
			UnsubscribeType: 1,
		}
		_, err = orders.Unsubscribe(bssV2Client, unsubscribeOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error unsubscribe HuaweiCloud GaussDB instance: %s", err)
		}
	} else {
		result := instances.Delete(client, instanceId)
		if result.Err != nil {
			return result.Err
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

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	d.SetId("")
	return nil
}

func resourceGeminiDBInstanceV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.GeminiDBV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc: %s", err)
	}
	//update tags
	if d.HasChange("tags") {
		tagErr := UpdateResourceTags(client, d, "instances", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of GeminiDB %q: %s", d.Id(), tagErr)
		}
	}

	if d.HasChange("name") {
		updateNameOpts := instances.UpdateNameOpts{
			Name: d.Get("name").(string),
		}

		err := instances.UpdateName(client, d.Id(), updateNameOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error updating name for huaweicloud_gaussdb_cassandra_instance %s: %s", d.Id(), err)
		}

	}

	if d.HasChange("password") {
		updatePassOpts := instances.UpdatePassOpts{
			Password: d.Get("password").(string),
		}

		err := instances.UpdatePass(client, d.Id(), updatePassOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error updating password for huaweicloud_gaussdb_cassandra_instance %s: %s", d.Id(), err)
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
			return fmt.Errorf("Error updating configuration_id for huaweicloud_gaussdb_cassandra_instance %s: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"SET_CONFIGURATION"},
			Target:     []string{"available"},
			Refresh:    GeminiDBInstanceUpdateRefreshFunc(client, d.Id(), "SET_CONFIGURATION"),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 10 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for huaweicloud_gaussdb_cassandra_instance %s to become ready: %s", d.Id(), err)
		}

		// Compare the target configuration and the instance configuration
		config, err := configurations.Get(client, configId).Extract()
		if err != nil {
			return fmt.Errorf("Error fetching configuration %s: %s", configId, err)
		}
		configParams := config.Parameters
		log.Printf("[DEBUG] Configuration Parameters %#v", configParams)

		instanceConfig, err := configurations.GetInstanceConfig(client, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("Error fetching instance configuration for huaweicloud_gaussdb_cassandra_instance %s: %s", d.Id(), err)
		}
		instanceConfigParams := instanceConfig.Parameters
		log.Printf("[DEBUG] Instance Configuration Parameters %#v", instanceConfigParams)

		if len(configParams) != len(instanceConfigParams) {
			return fmt.Errorf("Error updating configuration for instance: %s", d.Id())
		}
		for i, _ := range configParams {
			if !configParams[i].ReadOnly && configParams[i] != instanceConfigParams[i] {
				return fmt.Errorf("Error updating configuration for instance: %s", d.Id())
			}
		}
	}

	if d.HasChange("volume_size") {
		extendOpts := instances.ExtendVolumeOpts{
			Size: d.Get("volume_size").(int),
		}

		result := instances.ExtendVolume(client, d.Id(), extendOpts)
		if result.Err != nil {
			return fmt.Errorf("Error extending huaweicloud_gaussdb_cassandra_instance %s size: %s", d.Id(), result.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"RESIZE_VOLUME"},
			Target:     []string{"available"},
			Refresh:    GeminiDBInstanceUpdateRefreshFunc(client, d.Id(), "RESIZE_VOLUME"),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 10 * time.Second,
		}

		_, err := stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for huaweicloud_gaussdb_cassandra_instance %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("node_num") {
		old, newnum := d.GetChange("node_num")
		if newnum.(int) > old.(int) {
			//Enlarge Nodes
			expand_size := newnum.(int) - old.(int)
			enlargeNodeOpts := instances.EnlargeNodeOpts{
				Num: expand_size,
			}
			log.Printf("[DEBUG] Enlarge Node Options: %+v", enlargeNodeOpts)

			result := instances.EnlargeNode(client, d.Id(), enlargeNodeOpts)
			if result.Err != nil {
				return fmt.Errorf("Error enlarging huaweicloud_gaussdb_cassandra_instance %s node size: %s", d.Id(), result.Err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:      []string{"GROWING"},
				Target:       []string{"available"},
				Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, d.Id(), "GROWING"),
				Timeout:      d.Timeout(schema.TimeoutUpdate),
				Delay:        15 * time.Second,
				PollInterval: 20 * time.Second,
			}

			_, err := stateConf.WaitForState()
			if err != nil {
				return fmt.Errorf(
					"Error waiting for huaweicloud_gaussdb_cassandra_instance %s to become ready: %s", d.Id(), err)
			}
		}
		if newnum.(int) < old.(int) {
			//Reduce Nodes
			shrink_size := old.(int) - newnum.(int)
			reduceNodeOpts := instances.ReduceNodeOpts{
				Num: 1,
			}
			log.Printf("[DEBUG] Reduce Node Options: %+v", reduceNodeOpts)

			for i := 0; i < shrink_size; i++ {
				result := instances.ReduceNode(client, d.Id(), reduceNodeOpts)
				if result.Err != nil {
					return fmt.Errorf("Error shrinking huaweicloud_gaussdb_cassandra_instance %s node size: %s", d.Id(), result.Err)
				}

				stateConf := &resource.StateChangeConf{
					Pending:      []string{"REDUCING"},
					Target:       []string{"available"},
					Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, d.Id(), "REDUCING"),
					Timeout:      d.Timeout(schema.TimeoutUpdate),
					Delay:        15 * time.Second,
					PollInterval: 20 * time.Second,
				}

				_, err := stateConf.WaitForState()
				if err != nil {
					return fmt.Errorf(
						"Error waiting for huaweicloud_gaussdb_cassandra_instance %s to become ready: %s", d.Id(), err)
				}
			}
		}
	}

	if d.HasChange("flavor") {
		instance, err := instances.GetInstanceByID(client, d.Id())
		if err != nil {
			return fmt.Errorf(
				"Error fetching huaweicloud_gaussdb_cassandra_instance %s: %s", d.Id(), err)
		}

		specCode := ""
		for _, action := range instance.Actions {
			if action == "RESIZE_FLAVOR" {
				// Wait here if the instance already in RESIZE_FLAVOR state
				stateConf := &resource.StateChangeConf{
					Pending:      []string{"RESIZE_FLAVOR"},
					Target:       []string{"available"},
					Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, d.Id(), "RESIZE_FLAVOR"),
					Timeout:      d.Timeout(schema.TimeoutUpdate),
					PollInterval: 20 * time.Second,
				}

				_, err = stateConf.WaitForState()
				if err != nil {
					return fmt.Errorf(
						"Error waiting for huaweicloud_gaussdb_cassandra_instance %s to become ready: %s", d.Id(), err)
				}

				instance, err := instances.GetInstanceByID(client, d.Id())
				if err != nil {
					return fmt.Errorf(
						"Error fetching huaweicloud_gaussdb_cassandra_instance %s: %s", d.Id(), err)
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
			log.Printf("[DEBUG] Inconsistent Node SpecCode: %s, Flavor: %s", specCode, flavor)
			// Do resize action
			resizeOpts := instances.ResizeOpts{
				InstanceID: d.Id(),
				SpecCode:   d.Get("flavor").(string),
			}

			result := instances.Resize(client, d.Id(), resizeOpts)
			if result.Err != nil {
				return fmt.Errorf("Error resizing huaweicloud_gaussdb_cassandra_instance %s: %s", d.Id(), result.Err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:      []string{"RESIZE_FLAVOR"},
				Target:       []string{"available"},
				Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, d.Id(), "RESIZE_FLAVOR"),
				Timeout:      d.Timeout(schema.TimeoutUpdate),
				PollInterval: 20 * time.Second,
			}

			_, err := stateConf.WaitForState()
			if err != nil {
				return fmt.Errorf(
					"Error waiting for huaweicloud_gaussdb_cassandra_instance %s to become ready: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("security_group_id") {
		updateSgOpts := instances.UpdateSgOpts{
			SecurityGroupID: d.Get("security_group_id").(string),
		}

		result := instances.UpdateSg(client, d.Id(), updateSgOpts)
		if result.Err != nil {
			return fmt.Errorf("Error updating security group for huaweicloud_gaussdb_cassandra_instance %s: %s", d.Id(), result.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"MODIFY_SECURITYGROUP"},
			Target:       []string{"available"},
			Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, d.Id(), "MODIFY_SECURITYGROUP"),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			PollInterval: 3 * time.Second,
		}

		_, err := stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for huaweicloud_gaussdb_cassandra_instance %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("backup_strategy") {
		var updateOpts backups.UpdateOpts
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keep_days := rawMap["keep_days"].(int)
		updateOpts.KeepDays = &keep_days
		updateOpts.StartTime = rawMap["start_time"].(string)
		// Fixed to "1,2,3,4,5,6,7"
		updateOpts.Period = "1,2,3,4,5,6,7"
		log.Printf("[DEBUG] Update backup_strategy: %#v", updateOpts)

		err = backups.Update(client, d.Id(), updateOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error updating backup_strategy: %s", err)
		}
	}

	return resourceGeminiDBInstanceV3Read(d, meta)
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
			if action == state {
				return instance, state, nil
			}
		}

		return instance, "available", nil
	}
}
