package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
)

func resourceRdsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Delete: resourceInstanceDelete,
		Update: resourceInstanceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_rds_instance resource instead",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								return ValidateStringList(v, k, []string{"PostgreSQL", "SQLServer", "MySQL"})
							},
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					}},
			},

			"flavorref": {
				Type:     schema.TypeString,
				Required: true,
			},

			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					}},
			},

			"availabilityzone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"nics": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnetid": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					}},
			},

			"securitygroup": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					}},
			},

			"dbport": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"backupstrategy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"starttime": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"keepdays": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					}},
			},

			"dbrtpd": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ha": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"replicationmode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								return ValidateStringList(v, k, []string{"async", "sync", "semisync"})
							},
						},
					}},
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceInstanceDataStore(d *schema.ResourceData) instances.DataStoreOps {
	var dataStore instances.DataStoreOps
	datastoreRaw := d.Get("datastore").([]interface{})
	log.Printf("[DEBUG] datastoreRaw: %+v", datastoreRaw)
	if len(datastoreRaw) == 1 {
		dataStore.Type = datastoreRaw[0].(map[string]interface{})["type"].(string)
		dataStore.Version = datastoreRaw[0].(map[string]interface{})["version"].(string)
	}
	log.Printf("[DEBUG] datastore: %+v", dataStore)
	return dataStore
}

func resourceInstanceVolume(d *schema.ResourceData) instances.VolumeOps {
	var volume instances.VolumeOps
	volumeRaw := d.Get("volume").([]interface{})
	log.Printf("[DEBUG] volumeRaw: %+v", volumeRaw)
	if len(volumeRaw) == 1 {
		volume.Type = volumeRaw[0].(map[string]interface{})["type"].(string)
		volume.Size = volumeRaw[0].(map[string]interface{})["size"].(int)
	}
	log.Printf("[DEBUG] volume: %+v", volume)
	return volume
}

func resourceInstanceNics(d *schema.ResourceData) instances.NicsOps {
	var nics instances.NicsOps
	nicsRaw := d.Get("nics").([]interface{})
	log.Printf("[DEBUG] nicsRaw: %+v", nicsRaw)
	if len(nicsRaw) == 1 {
		nics.SubnetId = nicsRaw[0].(map[string]interface{})["subnetid"].(string)
	}
	log.Printf("[DEBUG] nics: %+v", nics)
	return nics
}

func resourceInstanceSecurityGroup(d *schema.ResourceData) instances.SecurityGroupOps {
	var securityGroup instances.SecurityGroupOps
	SecurityGroupRaw := d.Get("securitygroup").([]interface{})
	log.Printf("[DEBUG] SecurityGroupOpsRaw: %+v", SecurityGroupRaw)
	if len(SecurityGroupRaw) == 1 {
		securityGroup.Id = SecurityGroupRaw[0].(map[string]interface{})["id"].(string)
	}
	log.Printf("[DEBUG] securityGroup: %+v", securityGroup)
	return securityGroup

}

func resourceInstanceBackupStrategy(d *schema.ResourceData) instances.BackupStrategyOps {
	var backupStrategy instances.BackupStrategyOps
	backupStrategyRaw := d.Get("backupstrategy").([]interface{})
	log.Printf("[DEBUG] backupStrategyRaw: %+v", backupStrategyRaw)
	if len(backupStrategyRaw) == 1 {
		backupStrategy.StartTime = backupStrategyRaw[0].(map[string]interface{})["starttime"].(string)
		backupStrategy.KeepDays = backupStrategyRaw[0].(map[string]interface{})["keepdays"].(int)
	} else {
		backupStrategy.StartTime = "00:00:00"
		backupStrategy.KeepDays = 0
	}
	log.Printf("[DEBUG] backupStrategy: %+v", backupStrategy)
	return backupStrategy
}

func resourceInstanceHa(d *schema.ResourceData) instances.HaOps {
	var ha instances.HaOps
	haRaw := d.Get("ha").([]interface{})
	log.Printf("[DEBUG] haRaw: %+v", haRaw)
	if len(haRaw) == 1 {
		ha.Enable = haRaw[0].(map[string]interface{})["enable"].(bool)
		if ha.Enable == true {
			ha.ReplicationMode = haRaw[0].(map[string]interface{})["replicationmode"].(string)
		}
	} else {
		ha.Enable = false
	}
	log.Printf("[DEBUG] ha: %+v", ha)
	return ha
}

func InstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}

		return instance, instance.Status, nil
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud rds client: %s ", err)
	}

	createOpts := instances.CreateOps{
		Name:             d.Get("name").(string),
		DataStore:        resourceInstanceDataStore(d),
		FlavorRef:        d.Get("flavorref").(string),
		Volume:           resourceInstanceVolume(d),
		Region:           GetRegion(d, config),
		AvailabilityZone: d.Get("availabilityzone").(string),
		Vpc:              d.Get("vpc").(string),
		Nics:             resourceInstanceNics(d),
		SecurityGroup:    resourceInstanceSecurityGroup(d),
		DbPort:           d.Get("dbport").(string),
		BackupStrategy:   resourceInstanceBackupStrategy(d),
		DbRtPd:           d.Get("dbrtpd").(string),
		Ha:               resourceInstanceHa(d),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error getting instance from result: %s ", err)
	}
	log.Printf("[DEBUG] Create : instance %s: %#v", instance.ID, instance)

	d.SetId(instance.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    InstanceStateRefreshFunc(client, instance.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			instance.ID, err)
	}

	if instance.ID != "" {
		return resourceInstanceRead(d, meta)
	}
	return fmt.Errorf("Unexpected conversion error in resourceInstanceCreate. ")
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud rds client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "instance")
	}

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("hostname", instance.HostName)
	d.Set("type", instance.Type)
	d.Set("region", instance.Region)
	d.Set("availabilityzone", instance.AvailabilityZone)
	d.Set("vpc", instance.Vpc)
	d.Set("status", instance.Status)

	nicsList := make([]map[string]interface{}, 0, 1)
	nics := map[string]interface{}{
		"subnetid": instance.Nics.SubnetId,
	}
	nicsList = append(nicsList, nics)
	log.Printf("[DEBUG] nicsList: %+v", nicsList)
	if err := d.Set("nics", nicsList); err != nil {
		return fmt.Errorf("[DEBUG] Error saving nics to Rds instance (%s): %s", d.Id(), err)
	}

	securitygroupList := make([]map[string]interface{}, 0, 1)
	securitygroup := map[string]interface{}{
		"id": instance.SecurityGroup.Id,
	}
	securitygroupList = append(securitygroupList, securitygroup)
	log.Printf("[DEBUG] securitygroupList: %+v", securitygroupList)
	if err := d.Set("securitygroup", securitygroupList); err != nil {
		return fmt.Errorf("[DEBUG] Error saving securitygroup to Rds instance (%s): %s", d.Id(), err)
	}

	d.Set("flavorref", instance.Flavor.Id)

	volumeList := make([]map[string]interface{}, 0, 1)
	volume := map[string]interface{}{
		"type": instance.Volume.Type,
		"size": instance.Volume.Size,
	}
	volumeList = append(volumeList, volume)
	if err := d.Set("volume", volumeList); err != nil {
		return fmt.Errorf(
			"[DEBUG] Error saving volume to Rds instance (%s): %s", d.Id(), err)
	}

	d.Set("dbport", instance.DbPort)

	datastoreList := make([]map[string]interface{}, 0, 1)
	datastore := map[string]interface{}{
		"type":    instance.DataStore.Type,
		"version": instance.DataStore.Version,
	}
	datastoreList = append(datastoreList, datastore)
	if err := d.Set("datastore", datastoreList); err != nil {
		return fmt.Errorf(
			"[DEBUG] Error saving datastore to Rds instance (%s): %s", d.Id(), err)
	}

	d.Set("updated", instance.Updated)
	d.Set("created", instance.Created)
	return nil
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud rds client: %s ", err)
	}

	log.Printf("[DEBUG] Deleting Instance %s", d.Id())

	id := d.Id()
	result := instances.Delete(client, id)
	if result.Err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    InstanceStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			id, err)
	}

	log.Printf("[DEBUG] Successfully deleted instance %s", id)
	return nil
}

func InstanceStateUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID string, size int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}
		log.Printf("[DEBUG] Updating instance.Volume : %+v", instance.Volume)
		if instance.Volume.Size == size {
			return instance, "UPDATED", nil
		}

		return instance, instance.Status, nil
	}
}

func InstanceStateFlavorUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID string, flavorID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}

		return instance, instance.Status, nil
	}
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error Updating HuaweiCloud rds client: %s ", err)
	}

	log.Printf("[DEBUG] Updating instances %s", d.Id())
	id := d.Id()

	if d.HasChange("volume") {
		var updateOpts instances.UpdateOps
		volume := make(map[string]interface{})
		volumeRaw := d.Get("volume").([]interface{})
		log.Printf("[DEBUG] volumeRaw: %+v", volumeRaw)
		if len(volumeRaw) == 1 {
			volume["size"] = volumeRaw[0].(map[string]interface{})["size"].(int)
		}
		log.Printf("[DEBUG] volume: %+v", volume)
		updateOpts.Volume = volume
		_, err = instances.UpdateVolumeSize(client, updateOpts, id).Extract()
		if err != nil {
			return fmt.Errorf("Error updating instance volume from result: %s ", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"ACTIVE"},
			Target:     []string{"UPDATED"},
			Refresh:    InstanceStateUpdateRefreshFunc(client, id, updateOpts.Volume["size"].(int)),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) volume to be Updated: %s ",
				id, err)
		}
		log.Printf("[DEBUG] Successfully updated instance %s volume: %+v", id, volume)
	}

	if d.HasChange("flavorref") {
		var updateFlavorOpts instances.UpdateFlavorOps

		log.Printf("[DEBUG] Update flavorref: %s", d.Get("flavorref").(string))

		updateFlavorOpts.FlavorRef = d.Get("flavorref").(string)
		_, err = instances.UpdateFlavorRef(client, updateFlavorOpts, id).Extract()
		if err != nil {
			return fmt.Errorf("Error updating instance Flavor from result: %s ", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"MODIFYING"},
			Target:     []string{"ACTIVE"},
			Refresh:    InstanceStateFlavorUpdateRefreshFunc(client, id, d.Get("flavorref").(string)),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) flavor to be Updated: %s ",
				id, err)
		}
		log.Printf("[DEBUG] Successfully updated instance %s flavor: %s", id, d.Get("flavorref").(string))
	}

	if d.HasChange("backupstrategy") {
		var updatepolicyOpts instances.UpdatePolicyOps
		backupstrategyRaw := d.Get("backupstrategy").([]interface{})
		log.Printf("[DEBUG] backupstrategyRaw: %+v", backupstrategyRaw)
		if len(backupstrategyRaw) == 1 {
			updatepolicyOpts.StartTime = backupstrategyRaw[0].(map[string]interface{})["starttime"].(string)
			updatepolicyOpts.KeepDays = backupstrategyRaw[0].(map[string]interface{})["keepdays"].(int)
		}
		log.Printf("[DEBUG] updatepolicyOpts: %+v", updatepolicyOpts)
		_, err = instances.UpdatePolicy(client, updatepolicyOpts, id).Extract()
		if err != nil {
			return fmt.Errorf("Error updating instance policy from result: %s ", err)
		}

		log.Printf("[DEBUG] Successfully updated instance %s policy: %+v", id, updatepolicyOpts)
	}

	log.Printf("[DEBUG] Successfully updated instance %s", id)
	d.SetId(id)
	return resourceInstanceRead(d, meta)
}
