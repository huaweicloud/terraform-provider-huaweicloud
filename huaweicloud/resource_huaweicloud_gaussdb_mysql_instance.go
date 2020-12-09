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
	"github.com/huaweicloud/golangsdk/openstack/bss/v2/orders"
	"github.com/huaweicloud/golangsdk/openstack/taurusdb/v3/backups"
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
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"read_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
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
	client, err := config.gaussdbV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if hasFilledOpt(d, "force_import") {
		log.Printf("[DEBUG] Gaussdb mysql instance force_import is set, try to import it instead of creating")
		listOpts := instances.ListTaurusDBInstanceOpts{
			Name: d.Get("name").(string),
		}
		pages, err := instances.List(client, listOpts).AllPages()
		if err != nil {
			return err
		}

		allInstances, err := instances.ExtractTaurusDBInstances(pages)
		if err != nil {
			return fmt.Errorf("Unable to retrieve instances: %s ", err)
		}
		if allInstances.TotalCount > 0 {
			instance := allInstances.Instances[0]
			log.Printf("[DEBUG] Found existing mysql instance %s with name %s", instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceGaussDBInstanceRead(d, meta)
		}
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
		return fmt.Errorf("Error creating GaussDB instance : %s", err)
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

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			id, err)
	}

	if hasFilledOpt(d, "backup_strategy") {
		var updateOpts backups.UpdateOpts
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keep_days := rawMap["keep_days"].(int)
		updateOpts.KeepDays = &keep_days
		updateOpts.StartTime = rawMap["start_time"].(string)
		// Fixed to "1,2,3,4,5,6,7"
		updateOpts.Period = "1,2,3,4,5,6,7"
		log.Printf("[DEBUG] Update backup_strategy: %#v", updateOpts)

		err = backups.Update(client, id, updateOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error updating backup_strategy: %s", err)
		}
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second)

	return resourceGaussDBInstanceRead(d, meta)
}

func resourceGaussDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.gaussdbV3Client(region)
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
	flavor := ""
	slave_count := 0
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
		if raw.Type == "slave" && raw.Status == "ACTIVE" {
			slave_count += 1
		}
		if flavor == "" {
			flavor = raw.Flavor
		}
	}
	d.Set("nodes", nodesList)
	d.Set("read_replicas", slave_count)
	if flavor != "" {
		log.Printf("[DEBUG] Node Flavor: %s", flavor)
		d.Set("flavor", flavor)
	}

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
	client, err := config.gaussdbV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}
	instanceId := d.Id()

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		updateNameOpts := instances.UpdateNameOpts{
			Name: newName,
		}
		log.Printf("[DEBUG] Update Name Options: %+v", updateNameOpts)

		n, err := instances.UpdateName(client, instanceId, updateNameOpts).ExtractJobResponse()
		if err != nil {
			return fmt.Errorf("Error updating name for instance %s: %s ", instanceId, err)
		}

		if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.JobID); err != nil {
			return err
		}
		log.Printf("[DEBUG] Updated Name to %s for instance %s", newName, instanceId)
	}

	if d.HasChange("password") {
		newPass := d.Get("password").(string)
		updatePassOpts := instances.UpdatePassOpts{
			Password: newPass,
		}
		log.Printf("[DEBUG] Update Password Options: %+v", updatePassOpts)

		_, err := instances.UpdatePass(client, instanceId, updatePassOpts).ExtractJobResponse()
		if err != nil {
			return fmt.Errorf("Error updating password for instance %s: %s ", instanceId, err)
		}
		log.Printf("[DEBUG] Updated Password for instance %s", instanceId)
	}

	if d.HasChange("flavor") {
		newFlavor := d.Get("flavor").(string)
		resizeOpts := instances.ResizeOpts{
			Spec: newFlavor,
		}
		log.Printf("[DEBUG] Update Flavor Options: %+v", resizeOpts)

		n, err := instances.Resize(client, instanceId, resizeOpts).ExtractJobResponse()
		if err != nil {
			return fmt.Errorf("Error updating flavor for instance %s: %s ", instanceId, err)
		}

		if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.JobID); err != nil {
			return err
		}
		log.Printf("[DEBUG] Updated Flavor for instance %s", instanceId)
	}

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
				return fmt.Errorf("Error creating read replicas for instance %s: %s ", instanceId, err)
			}

			job_list := strings.Split(n.JobID, ",")
			log.Printf("[DEBUG] Create Replica Jobs: %#v", job_list)
			for i := 0; i < len(job_list); i++ {
				job_id := job_list[i]
				log.Printf("[DEBUG] Waiting for job: %s", job_id)
				if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), job_id); err != nil {
					return err
				}
			}
		}
		if newnum.(int) < old.(int) {
			shrink_size := old.(int) - newnum.(int)

			slave_nodes := []string{}
			nodes := d.Get("nodes").([]interface{})
			for _, nodeRaw := range nodes {
				node := nodeRaw.(map[string]interface{})
				if node["type"].(string) == "slave" && node["status"] == "ACTIVE" {
					slave_nodes = append(slave_nodes, node["id"].(string))
				}
			}
			log.Printf("[DEBUG] Slave Nodes: %+v", slave_nodes)
			if len(slave_nodes) <= shrink_size {
				return fmt.Errorf("Error deleting read replicas for instance %s: Shrink Size is bigger than active slave nodes", instanceId)
			}
			for i := 0; i < shrink_size; i++ {
				n, err := instances.DeleteReplica(client, instanceId, slave_nodes[i]).ExtractJobResponse()
				if err != nil {
					return fmt.Errorf("Error creating read replica %s for instance %s: %s ", slave_nodes[i], instanceId, err)
				}

				if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.JobID); err != nil {
					return err
				}
				log.Printf("[DEBUG] Deleted Read Replica: %s", slave_nodes[i])
			}
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

	return resourceGaussDBInstanceRead(d, meta)
}

func resourceGaussDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.gaussdbV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
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
			return CheckDeleted(d, result.Err, "GaussDB instance")
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

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	return nil
}
