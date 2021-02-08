package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/datastores"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/flavors"
	instanceV1 "github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
)

func resourceRdsReadReplicaInstance() *schema.Resource {

	return &schema.Resource{

		Create: resourceRdsReadReplicaInstanceCreate,
		Read:   resourceRdsReadReplicaInstanceRead,
		Update: resourceRdsReadReplicaInstanceUpdate,
		Delete: resourceRdsReadReplicaInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				Required: true,
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"primary_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},

			"volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk_encryption_id": {
							Type:     schema.TypeString,
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

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
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

			"db": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceRdsReadReplicaInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud rds client: %s ", err)
	}

	createOpts := instances.CreateReplicaOpts{
		Name:                d.Get("name").(string),
		ReplicaOfId:         d.Get("primary_instance_id").(string),
		FlavorRef:           d.Get("flavor").(string),
		Region:              GetRegion(d, config),
		AvailabilityZone:    d.Get("availability_zone").(string),
		Volume:              resourceReplicaInstanceVolume(d),
		DiskEncryptionId:    resourceDiskEncryptionID(d),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
	}
	log.Printf("[DEBUG] Create replica instance Options: %#v", createOpts)

	resp, err := instances.CreateReplica(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating replica instance: %s ", err)
	}

	instance := resp.Instance
	d.SetId(instance.Id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "RESTORING"},
		Target:     []string{"ACTIVE"},
		Refresh:    rdsInstanceStateRefreshFunc(client, instance.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for read replica instance (%s) to become ready: %s ", instance.Id, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := expandResourceTags(tagRaw)
		err := tags.Create(client, "instances", instance.Id, tagList).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error setting tags of Rds read replica instance %s: %s", instance.Id, err)
		}
	}

	return resourceRdsReadReplicaInstanceRead(d, meta)
}

func resourceRdsReadReplicaInstanceRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud rds client: %s", err)
	}

	instanceID := d.Id()
	instance, err := getRdsInstanceByID(client, instanceID)
	if err != nil {
		return err
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved rds read replica instance %s: %#v", instanceID, instance)

	az := readAvailabilityZone(instance)
	d.Set("name", instance.Name)
	d.Set("flavor", instance.FlavorRef)
	d.Set("region", instance.Region)
	d.Set("availability_zone", az)
	d.Set("private_ips", instance.PrivateIps)
	d.Set("public_ips", instance.PublicIps)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("type", instance.Type)
	d.Set("status", instance.Status)
	d.Set("enterprise_project_id", instance.EnterpriseProjectId)

	primaryInstanceID, err := readPrimaryInstanceID(instance)
	if err != nil {
		return err
	}
	d.Set("primary_instance_id", primaryInstanceID)

	// set tags
	tagsMap := make(map[string]interface{})
	for _, tag := range instance.Tags {
		tagsMap[tag.Key] = tag.Value
	}
	d.Set("tags", tagsMap)

	// save volume
	volumeList := make([]map[string]interface{}, 0, 1)
	volume := map[string]interface{}{
		"type":               instance.Volume.Type,
		"size":               instance.Volume.Size,
		"disk_encryption_id": instance.DiskEncryptionId,
	}
	volumeList = append(volumeList, volume)
	if err := d.Set("volume", volumeList); err != nil {
		return fmt.Errorf("[DEBUG] Error saving volume to RDS read replica instance (%s): %s", d.Id(), err)
	}

	// save database
	dbList := make([]map[string]interface{}, 0, 1)
	database := map[string]interface{}{
		"type":      instance.DataStore.Type,
		"version":   instance.DataStore.Version,
		"port":      instance.Port,
		"user_name": instance.DbUserName,
	}
	dbList = append(dbList, database)
	if err := d.Set("db", dbList); err != nil {
		return fmt.Errorf("[DEBUG] Error saving data base to RDS read replica instance (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceRdsReadReplicaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	clientV3, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud rds v3 client: %s ", err)
	}
	clientV1, err := config.RdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud rds v1 client: %s ", err)
	}

	// Fetching node id
	instance, err := getRdsInstanceByID(clientV3, d.Id())
	if err != nil {
		return err
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}
	nodeID := instance.Nodes[0].Id
	log.Printf("[DEBUG] primary instance node id: %s", nodeID)

	// update flavor
	if d.HasChange("flavor") {
		newFlavor := d.Get("flavor")

		// Fetch flavor id
		db := d.Get("db").([]interface{})
		datastoreName := db[0].(map[string]interface{})["type"].(string)
		datastoreVersion := db[0].(map[string]interface{})["version"].(string)
		datastoresList, err := datastores.List(clientV1, datastoreName).Extract()
		if err != nil {
			return fmt.Errorf("Unable to retrieve datastores: %s ", err)
		}
		if len(datastoresList) < 1 {
			return fmt.Errorf("Returned no datastore result. ")
		}
		var datastoreID string
		for _, datastore := range datastoresList {
			if strings.HasPrefix(datastore.Name, datastoreVersion) {
				datastoreID = datastore.ID
				break
			}
		}
		if datastoreID == "" {
			return fmt.Errorf("Returned no datastore ID. ")
		}
		log.Printf("[DEBUG] Received datastore Id: %s", datastoreID)
		flavorsList, err := flavors.List(clientV1, datastoreID, GetRegion(d, config)).Extract()
		if err != nil {
			return fmt.Errorf("Unable to retrieve flavors: %s", err)
		}
		if len(flavorsList) < 1 {
			return fmt.Errorf("Returned no flavor result. ")
		}
		var rdsFlavor flavors.Flavor
		for _, flavor := range flavorsList {
			if flavor.SpecCode == newFlavor.(string) {
				rdsFlavor = flavor
				break
			}
		}

		log.Printf("[DEBUG] Update flavor: %s", newFlavor.(string))
		var updateFlavorOpts instanceV1.UpdateFlavorOps
		updateFlavorOpts.FlavorRef = rdsFlavor.ID
		_, err = instanceV1.UpdateFlavorRef(clientV1, updateFlavorOpts, nodeID).Extract()
		if err != nil {
			return fmt.Errorf("Error updating replica instance Flavor from result: %s ", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"MODIFYING"},
			Target:     []string{"ACTIVE"},
			Refresh:    replicaInstanceStateFlavorUpdateRefreshFunc(clientV1, nodeID, d.Get("flavor").(string)),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for replica instance (%s) flavor to be Updated: %s ", d.Id(), err)
		}
		log.Printf("[DEBUG] Successfully updated replica instance %s flavor: %s", d.Id(), d.Get("flavor").(string))

	}

	// update volume
	if d.HasChange("volume") {
		newVolume := d.Get("volume")
		var updateOpts instanceV1.UpdateOps
		volume := make(map[string]interface{})
		volumeRaw := newVolume.([]interface{})
		if len(volumeRaw) == 1 {
			if m, ok := volumeRaw[0].(map[string]interface{}); ok {
				volume["size"] = m["size"].(int)
			}
		}

		updateOpts.Volume = volume
		_, err := instanceV1.UpdateVolumeSize(clientV1, updateOpts, nodeID).Extract()
		if err != nil {
			return fmt.Errorf("Error updating read replica volume from result: %s ", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"MODIFYING"},
			Target:     []string{"UPDATED"},
			Refresh:    replicaInstanceStateUpdateRefreshFunc(clientV1, nodeID, updateOpts.Volume["size"].(int)),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for read replica (%s) volume to be Updated: %s ", d.Id(), err)
		}
		log.Printf("[DEBUG] Successfully updated read replica %s volume: %+v", d.Id(), volume)
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := UpdateResourceTags(clientV3, d, "instances", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of RDS read replica instance: %s, err: %s", d.Id(), tagErr)
		}
	}

	return resourceRdsReadReplicaInstanceRead(d, meta)
}

func resourceRdsReadReplicaInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud rds client: %s ", err)
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
		Refresh:    rdsInstanceStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for read replica instance (%s) to be deleted: %s ",
			id, err)
	}

	log.Printf("[DEBUG] Successfully deleted rds read replica instance %s", id)
	return nil
}

func getRdsInstanceByID(client *golangsdk.ServiceClient, instanceID string) (*instances.RdsInstanceResponse, error) {
	listOpts := instances.ListRdsInstanceOpts{
		Id: instanceID,
	}
	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("An error occured while querying rds read replica instance %s: %s", instanceID, err)
	}

	resp, err := instances.ExtractRdsInstances(pages)
	if err != nil {
		return nil, err
	}

	instanceList := resp.Instances
	if len(instanceList) == 0 {
		// return an empty rds instance
		log.Printf("[WARN] can not find the specified rds read replica instance %s", instanceID)
		instance := new(instances.RdsInstanceResponse)
		return instance, nil
	}

	if len(instanceList) > 1 {
		return nil, fmt.Errorf("retrieving more than one rds read replica instance by %s", instanceID)
	}
	if instanceList[0].Id != instanceID {
		return nil, fmt.Errorf("the id of rds read replica instance was expected %s, but got %s",
			instanceID, instanceList[0].Id)
	}

	return &instanceList[0], nil
}

func readAvailabilityZone(resp *instances.RdsInstanceResponse) string {
	node := resp.Nodes[0]
	return node.AvailabilityZone
}

func readPrimaryInstanceID(resp *instances.RdsInstanceResponse) (string, error) {
	relatedInst := resp.RelatedInstance
	for _, relate := range relatedInst {
		if relate.Type == "replica_of" {
			return relate.Id, nil
		}
	}
	return "", fmt.Errorf("Error when get primary instance id for replica %s", resp.Id)
}

func resourceDiskEncryptionID(d *schema.ResourceData) string {
	var encryptionID string
	volumeRaw := d.Get("volume").([]interface{})

	if len(volumeRaw) == 1 {
		encryptionID = volumeRaw[0].(map[string]interface{})["disk_encryption_id"].(string)
	}

	return encryptionID
}

func resourceReplicaInstanceVolume(d *schema.ResourceData) *instances.Volume {
	var volume *instances.Volume
	volumeRaw := d.Get("volume").([]interface{})

	if len(volumeRaw) == 1 {
		volume = new(instances.Volume)
		volume.Type = volumeRaw[0].(map[string]interface{})["type"].(string)
		volume.Size = volumeRaw[0].(map[string]interface{})["size"].(int)
		// the size is optional and invalid for replica, but it's required in sdk
		// so just set 100 if not specified
		if volume.Size == 0 {
			volume.Size = 100
		}
	}
	log.Printf("[DEBUG] volume: %+v", volume)
	return volume
}

func rdsInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getRdsInstanceByID(client, instanceID)
		if err != nil {
			return nil, "FOUND ERROR", err
		}
		if instance.Id == "" {
			return instance, "DELETED", nil
		}

		return instance, instance.Status, nil
	}
}

func replicaInstanceStateUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID string, size int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instanceV1.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}
		log.Printf("[DEBUG] Updating read replica instance.Volume : %+v", instance.Volume)
		if instance.Volume.Size == size {
			return instance, "UPDATED", nil
		}

		return instance, instance.Status, nil
	}
}

func replicaInstanceStateFlavorUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID string, flavorID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instanceV1.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}

		return instance, instance.Status, nil
	}
}
