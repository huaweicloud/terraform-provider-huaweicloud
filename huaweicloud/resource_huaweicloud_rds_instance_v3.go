package huaweicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/backups"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func resourceRdsInstanceV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdsInstanceV3Create,
		Read:   resourceRdsInstanceV3Read,
		Update: resourceRdsInstanceV3Update,
		Delete: resourceRdsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"db": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
							ForceNew:  true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
					},
				},
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidateIP,
			},

			"ha_replication_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"param_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),

			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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

			"public_ips": {
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

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRdsInstanceV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	client, err := config.RdsV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud RDS client: %s", err)
	}

	createRdsOpts := instances.CreateRdsOpts{
		Name:                d.Get("name").(string),
		FlavorRef:           d.Get("flavor").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("param_group_id").(string),
		TimeZone:            d.Get("time_zone").(string),
		FixedIp:             d.Get("fixed_ip").(string),
		Password:            d.Get("db.0.password").(string),
		DiskEncryptionId:    d.Get("volume.0.disk_encryption_id").(string),
		Port:                strconv.Itoa(d.Get("db.0.port").(int)),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
		Region:              region,
		AvailabilityZone:    resourceRdsInstanceAvailabilityZone(d),
		Datastore:           resourceRdsInstanceDatastore(d),
		Volume:              resourceRdsInstanceVolume(d),
		BackupStrategy:      resourceRdsInstanceBackupStrategy(d),
		Ha:                  resourceRdsInstanceHaReplicationMode(d),
	}

	res, err := instances.Create(client, createRdsOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud RDS instance: %s", err)
	}
	d.SetId(res.Instance.Id)
	instanceID := d.Id()

	if err := checkRDSInstanceJobFinish(client, res.JobId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("Error creating instance (%s): %s", instanceID, err)
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", instanceID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of RDS instance (%s): %s", instanceID, tagErr)
		}
	}

	return resourceRdsInstanceV3Read(d, meta)
}

func resourceRdsInstanceV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud RDS client: %s", err)
	}

	instanceID := d.Id()
	instance, err := getRdsInstanceByID(client, instanceID)
	if err != nil {
		return fmt.Errorf("Error getting huaweicloud RDS instance: %s", err)
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] Retrieved RDS instance (%s): %#v", instanceID, instance)

	d.Set("region", instance.Region)
	d.Set("name", instance.Name)
	d.Set("status", instance.Status)
	d.Set("port", instance.Port)
	d.Set("type", instance.Type)
	d.Set("created", instance.Created)
	d.Set("ha_replication_mode", instance.Ha.ReplicationMode)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("flavor", instance.FlavorRef)
	d.Set("disk_encryption_id", instance.DiskEncryptionId)
	d.Set("time_zone", instance.TimeZone)
	d.Set("enterprise_project_id", instance.EnterpriseProjectId)

	publicIps := make([]interface{}, len(instance.PublicIps))
	for i, v := range instance.PublicIps {
		publicIps[i] = v
	}
	d.Set("public_ips", publicIps)

	privateIps := make([]string, len(instance.PrivateIps))
	for i, v := range instance.PrivateIps {
		privateIps[i] = v
	}
	d.Set("private_ips", privateIps)
	d.Set("fixed_ip", privateIps[0])

	volume := make([]map[string]interface{}, 1)
	volume[0] = map[string]interface{}{
		"type":               instance.Volume.Type,
		"size":               instance.Volume.Size,
		"disk_encryption_id": instance.DiskEncryptionId,
	}
	if err := d.Set("volume", volume); err != nil {
		return fmt.Errorf("[DEBUG] Error saving volume to RDS instance (%s): %s", instanceID, err)
	}

	dbList := make([]map[string]interface{}, 1)
	database := map[string]interface{}{
		"type":      instance.DataStore.Type,
		"version":   instance.DataStore.Version,
		"port":      instance.Port,
		"user_name": instance.DbUserName,
	}
	if len(d.Get("db").([]interface{})) > 0 {
		database["password"] = d.Get("db.0.password")
	}
	dbList[0] = database
	if err := d.Set("db", dbList); err != nil {
		return fmt.Errorf("[DEBUG] Error saving data base to RDS instance (%s): %s", instanceID, err)
	}

	backup := make([]map[string]interface{}, 1)
	backup[0] = map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
	}
	if err := d.Set("backup_strategy", backup); err != nil {
		return fmt.Errorf("[DEBUG] Error saving backup strategy to RDS instance (%s): %s", instanceID, err)
	}

	nodes := make([]map[string]interface{}, len(instance.Nodes))
	for i, v := range instance.Nodes {
		nodes[i] = map[string]interface{}{
			"id":                v.Id,
			"name":              v.Name,
			"role":              v.Role,
			"status":            v.Status,
			"availability_zone": v.AvailabilityZone,
		}
	}
	if err := d.Set("nodes", nodes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving nodes to RDS instance (%s): %s", instanceID, err)
	}

	d.Set("tags", utils.TagsToMap(instance.Tags))

	az1 := instance.Nodes[0].AvailabilityZone
	if strings.HasSuffix(d.Get("flavor").(string), ".ha") {
		if len(instance.Nodes) < 2 {
			return fmt.Errorf("[DEBUG] Error saving availability zone to RDS instance (%s): "+
				"HA mode must have two availability zone", instanceID)
		}
		az2 := instance.Nodes[1].AvailabilityZone
		if instance.Nodes[1].Role == "master" {
			d.Set("availability_zone", []string{az2, az1})
		} else {
			d.Set("availability_zone", []string{az1, az2})
		}
	} else {
		d.Set("availability_zone", []string{az1})
	}

	return nil
}

func resourceRdsInstanceV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud RDS Client: %s", err)
	}
	instanceID := d.Id()
	// Since the instance will throw an exception when making an API interface call in 'BACKING UP' state,
	// wait for the instance state to be updated to 'ACTIVE' before calling the interface.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BACKING UP", "ACTIVE"},
		Target:     []string{"ACTIVE"},
		Refresh:    rdsInstanceStateRefreshFunc(client, instanceID),
		Timeout:    d.Timeout(schema.TimeoutDefault),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for RDS instance (%s) become active state: %s", instanceID, err)
	}

	if err := updateRdsInstanceName(d, client, instanceID); err != nil {
		return fmt.Errorf("[ERROR] %s", err)
	}

	if err := updateRdsInstanceFlavor(d, client, instanceID); err != nil {
		return fmt.Errorf("[ERROR] %s", err)
	}

	if err := updateRdsInstanceVolume(d, client, instanceID); err != nil {
		return fmt.Errorf("[ERROR] %s", err)
	}

	if err := updateRdsInstanceBackpStrategy(d, client, instanceID); err != nil {
		return fmt.Errorf("[ERROR] %s", err)
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", instanceID)
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of RDS instance (%s): %s", instanceID, tagErr)
		}
	}

	return resourceRdsInstanceV3Read(d, meta)
}

func resourceRdsInstanceAvailabilityZone(d *schema.ResourceData) string {
	azList := make([]string, len(d.Get("availability_zone").([]interface{})))
	for i, az := range d.Get("availability_zone").([]interface{}) {
		azList[i] = az.(string)
	}
	return strings.Join(azList, ",")
}

func resourceRdsInstanceDatastore(d *schema.ResourceData) *instances.Datastore {
	var database *instances.Datastore
	dbRaw := d.Get("db").([]interface{})

	if len(dbRaw) == 1 {
		database = new(instances.Datastore)
		database.Type = dbRaw[0].(map[string]interface{})["type"].(string)
		database.Version = dbRaw[0].(map[string]interface{})["version"].(string)
	}
	return database
}

func resourceRdsInstanceVolume(d *schema.ResourceData) *instances.Volume {
	var volume *instances.Volume
	volumeRaw := d.Get("volume").([]interface{})

	if len(volumeRaw) == 1 {
		volume = new(instances.Volume)
		volume.Type = volumeRaw[0].(map[string]interface{})["type"].(string)
		volume.Size = volumeRaw[0].(map[string]interface{})["size"].(int)
	}
	return volume
}

func resourceRdsInstanceBackupStrategy(d *schema.ResourceData) *instances.BackupStrategy {
	var backupStrategy *instances.BackupStrategy
	backupRaw := d.Get("backup_strategy").([]interface{})

	if len(backupRaw) == 1 {
		backupStrategy = new(instances.BackupStrategy)
		backupStrategy.StartTime = backupRaw[0].(map[string]interface{})["start_time"].(string)
		backupStrategy.KeepDays = backupRaw[0].(map[string]interface{})["keep_days"].(int)
	}
	return backupStrategy
}

func resourceRdsInstanceHaReplicationMode(d *schema.ResourceData) *instances.Ha {
	var ha *instances.Ha
	if hasFilledOpt(d, "ha_replication_mode") {
		ha.Mode = "ha"
		ha.ReplicationMode = d.Get("ha_replication_mode").(string)
	}
	return ha
}

func updateRdsInstanceBackpStrategy(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string) error {
	if d.HasChange("backup_strategy") {
		var updateOpts backups.UpdateOpts

		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keep_days := rawMap["keep_days"].(int)
		updateOpts.KeepDays = &keep_days
		updateOpts.StartTime = rawMap["start_time"].(string)
		// TODO(zhenguo): Make Period configured by users
		updateOpts.Period = "1,2,3,4,5,6,7"
		log.Printf("[DEBUG] updateOpts: %#v", updateOpts)
		err := backups.Update(client, instanceID, updateOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud RDS instance (%s): %s", instanceID, err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"BACKING UP"},
			Target:     []string{"ACTIVE"},
			Refresh:    rdsInstanceStateRefreshFunc(client, instanceID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("Error waiting for RDS instance (%s) backup to be updated: %s ", instanceID, err)
		}
	}
	return nil
}

func updateRdsInstanceName(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string) error {
	if d.HasChange("name") {
		var renameRdsInstanceOpts instances.RenameRdsInstanceOpts

		renameRdsInstanceOpts.Name = d.Get("name").(string)
		r := instances.Rename(client, renameRdsInstanceOpts, instanceID)
		if r.Result.Err != nil {
			return fmt.Errorf("Error renaming HuaweiCloud RDS instance (%s): %s", instanceID, r.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"MODIFYING", "ACTIVE"},
			Target:     []string{"ACTIVE"},
			Refresh:    rdsInstanceStateRefreshFunc(client, instanceID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("Error waiting for RDS instance (%s) flavor to be updated: %s ", instanceID, err)
		}
	}
	return nil
}

func checkRDSInstanceJobFinish(client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed", "Failed"},
		Refresh:      rdsInstanceJobRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for RDS instance (%s) job to be completed: %s ", jobID, err)
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
