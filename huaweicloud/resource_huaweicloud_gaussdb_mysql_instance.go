package huaweicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/taurusdb/v3/instances"
)

func resourceGaussDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceGaussDBInstanceCreate,
		Update: resourceGaussDBInstanceUpdate,
		Read:   resourceGaussDBInstanceRead,
		Delete: resourceGaussDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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
				Optional: true,
				ForceNew: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"read_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  1,
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

func resourceGaussDBBackupStrategy(d *schema.ResourceData) *instances.BackupStrategyOpt {
	var backupOpt instances.BackupStrategyOpt

	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 1 {
		strategy := backupStrategyRaw[0].(map[string]interface{})
		backupOpt.StartTime = strategy["start_time"].(string)
		backupOpt.KeepDays = strconv.Itoa(strategy["keep_days"].(int))
	} else {
		// set defautl backup strategy
		backupOpt.StartTime = "00:00-01:00"
		backupOpt.KeepDays = "7"
	}

	return &backupOpt
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

func resourceGaussDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.initServiceClient("gaussdb", GetRegion(d, config), "mysql/v3")
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	createOpts := instances.CreateTaurusDBOpts{
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor").(string),
		Password:            d.Get("password").(string),
		Region:              GetRegion(d, config),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
		TimeZone:            d.Get("time_zone").(string),
		SlaveCount:          d.Get("read_replicas").(int),
		Mode:                "Cluster",
		DataStore:           resourceGaussDBDataStore(d),
		BackupStrategy:      resourceGaussDBBackupStrategy(d),
	}
	azMode := d.Get("availability_zone_mode").(string)
	createOpts.AZMode = azMode
	if azMode == "multi" {
		v, exist := d.GetOk("master_availability_zone")
		if !exist {
			return fmt.Errorf("missing master_availability_zone in a multi availability zone mode")
		}
		createOpts.MasterAZ = v.(string)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating GaussDB instance : %s", err)
	}

	id := instance.Instance.Id
	d.SetId(id)

	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "BACKING UP"},
		Target:     []string{"ACTIVE"},
		Refresh:    GaussDBInstanceStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      180 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		// the instance has created, but the status is unnormal
		deleteErr := resourceGaussDBInstanceDelete(d, meta)
		if deleteErr != nil {
			log.Printf("[ERROR] Error deleting GaussDB instance: %s", deleteErr)
		}

		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			id, err)
	}

	return resourceGaussDBInstanceRead(d, meta)
}

func resourceGaussDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.initServiceClient("gaussdb", region, "mysql/v3")
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "GaussDB instance")
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("region", region)
	d.Set("name", instance.Name)
	d.Set("status", instance.Status)
	d.Set("mode", instance.Type)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("configuration_id", instance.ConfigurationId)
	d.Set("db_user_name", instance.DbUserName)
	d.Set("time_zone", instance.TimeZone)
	d.Set("availability_zone_mode", instance.AZMode)
	d.Set("master_availability_zone", instance.MasterAZ)

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		d.Set("port", dbPort)
	}
	if len(instance.PrivateIps) > 0 {
		d.Set("private_write_ip", instance.PrivateIps[0])
	}

	// set data store
	dbList := make([]map[string]interface{}, 1)
	db := map[string]interface{}{
		"version": instance.DataStore.Version,
	}
	// normalize engine
	engine := instance.DataStore.Type
	if engine == "GaussDB(for MySQL)" {
		engine = "gaussdb-mysql"
	}
	db["engine"] = engine
	dbList[0] = db
	d.Set("datastore", dbList)

	// set nodes
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range instance.Nodes {
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
		nodesList = append(nodesList, node)
	}
	d.Set("nodes", nodesList)

	// set backup_strategy
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
	}
	if days, err := strconv.Atoi(instance.BackupStrategy.KeepDays); err == nil {
		backupStrategy["keep_days"] = days
	}
	backupStrategyList[0] = backupStrategy
	d.Set("backup_strategy", backupStrategyList)

	return nil
}

func resourceGaussDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.initServiceClient("gaussdb", GetRegion(d, config), "mysql/v3")
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}
	instanceId := d.Id()

	if d.HasChange("read_replicas") {
		old, newnum := d.GetChange("read_replicas")
		if newnum.(int) > old.(int) {
			expand_size := newnum.(int) - old.(int)
			priorities := []int{}
			for i := 0; i < expand_size; i++ {
				priorities = append(priorities, 1)
			}
			createReplicaOpts := instances.CreateReplicaOpts{
				Priorities: priorities,
			}
			log.Printf("[DEBUG] Create Replica Options: %+v", createReplicaOpts)

			n, err := instances.CreateReplica(client, instanceId, createReplicaOpts).ExtractJobResponse()
			if err != nil {
				d.Set("read_replicas", old.(int))
				return fmt.Errorf("Error creating read replicas for instance %s: %s ", instanceId, err)
			}

			job_list := strings.Split(n.JobID, ",")
			log.Printf("[DEBUG] Create Replica Jobs: %#v", job_list)
			for i := 0; i < len(job_list); i++ {
				job_id := job_list[i]
				log.Printf("[DEBUG] Waiting for job: %s", job_id)
				if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), job_id); err != nil {
					return err
				}
			}
		}
		if newnum.(int) < old.(int) {
			shrink_size := old.(int) - newnum.(int)

			slave_nodes := []string{}
			nodes := d.Get("nodes").([]interface{})
			for i := range nodes {
				node := nodes[i].(map[string]interface{})
				if node["type"].(string) == "slave" && node["status"] == "ACTIVE" {
					slave_nodes = append(slave_nodes, node["id"].(string))
				}
			}
			log.Printf("[DEBUG] Slave Nodes: %+v", slave_nodes)
			if len(slave_nodes) <= shrink_size {
				d.Set("read_replicas", old.(int))
				return fmt.Errorf("Error deleting read replicas for instance %s: Shrink Size is bigger than slave nodes", instanceId)
			}
			for i := 0; i < shrink_size; i++ {
				n, err := instances.DeleteReplica(client, instanceId, slave_nodes[i]).ExtractJobResponse()
				read_replicas := old.(int) - i
				if err != nil {
					d.Set("read_replicas", read_replicas)
					return fmt.Errorf("Error creating read replica %s for instance %s: %s ", slave_nodes[i], instanceId, err)
				}

				if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
					d.Set("read_replicas", read_replicas)
					return err
				}
				log.Printf("[DEBUG] Deleted Read Replica: %s", slave_nodes[i])
			}
		}
	}

	return resourceGaussDBInstanceRead(d, meta)
}

func resourceGaussDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.initServiceClient("gaussdb", GetRegion(d, config), "mysql/v3")
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	instanceId := d.Id()
	result := instances.Delete(client, instanceId)
	if result.Err != nil {
		return CheckDeleted(d, result.Err, "GaussDB instance")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "BACKING UP", "FAILED"},
		Target:     []string{"DELETED"},
		Refresh:    GaussDBInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	return nil
}
