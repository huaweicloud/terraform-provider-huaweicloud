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
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
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
			Create: schema.DefaultTimeout(30 * time.Minute),
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
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
				ForceNew:  true,
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
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
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
		Num:      "3",
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

	createOpts := instances.CreateGeminiDBOpts{
		Name:             d.Get("name").(string),
		Region:           GetRegion(d, config),
		AvailabilityZone: d.Get("availability_zone").(string),
		VpcId:            d.Get("vpc_id").(string),
		SubnetId:         d.Get("subnet_id").(string),
		SecurityGroupId:  d.Get("security_group_id").(string),
		Password:         d.Get("password").(string),
		Mode:             "Cluster",
		Flavor:           resourceGeminiDBFlavor(d),
		DataStore:        resourceGeminiDBDataStore(d),
		BackupStrategy:   resourceGeminiDBBackupStrategy(d),
	}
	/*
		backupStrategy := resourceGeminiDBBackupStrategy(d)
		if backupStrategy.StartTime != "" {
			createOpts.BackupStrategy = backupStrategy
		}
	*/
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating GeminiDB instance : %s", err)
	}

	d.SetId(instance.Id)
	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"normal"},
		Refresh:    GeminiDBInstanceStateRefreshFunc(client, instance.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      120 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		// the instance has created, but the status is unnormal
		deleteErr := resourceGeminiDBInstanceV3Delete(d, meta)
		if deleteErr != nil {
			log.Printf("[ERROR] Error deleting GeminiDB instance: %s", deleteErr)
		}
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
		return fmt.Errorf("Error fetching GeminiDB instance: deleted")
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

	nodesList := make([]map[string]interface{}, 0, 1)
	for _, group := range instance.Groups {
		for _, Node := range group.Nodes {
			node := map[string]interface{}{
				"id":         Node.Id,
				"name":       Node.Name,
				"status":     Node.Status,
				"private_ip": Node.PrivateIp,
			}
			nodesList = append(nodesList, node)
		}
	}
	d.Set("nodes", nodesList)

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
	result := instances.Delete(client, instanceId)
	if result.Err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"normal", "abnormal", "creating", "createfail", "enlargefail", "data_disk_full"},
		Target:     []string{"deleted"},
		Refresh:    GeminiDBInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
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
	tagErr := UpdateResourceTags(client, d, "instances")
	if tagErr != nil {
		return fmt.Errorf("Error updating tags of GeminiDB %q: %s", d.Id(), tagErr)
	}

	return resourceGeminiDBInstanceV3Read(d, meta)
}
