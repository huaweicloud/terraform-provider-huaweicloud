package huaweicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/dds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceDdsInstanceV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceDdsInstanceV3Create,
		Read:   resourceDdsInstanceV3Read,
		Update: resourceDdsInstanceV3Update,
		Delete: resourceDdsInstanceV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
							ValidateFunc: validation.StringInSlice([]string{
								"DDS-Community", "DDS-Enhanced",
							}, true),
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"4.0", "3.4", "3.2",
							}, true),
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"wiredTiger", "rocksDB",
							}, true),
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"disk_encryption_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Sharding", "ReplicaSet", "Single",
				}, true),
			},
			"flavor": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"mongos", "shard", "config", "replica", "single",
							}, true),
						},
						"num": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"storage": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ULTRAHIGH",
							}, true),
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"spec_code": {
							Type:     schema.TypeString,
							Required: true,
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
							Required: true,
						},
					},
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"db_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
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
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
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
		},
	}
}

func resourceDdsDataStore(d *schema.ResourceData) instances.DataStore {
	var dataStore instances.DataStore
	datastoreRaw := d.Get("datastore").([]interface{})
	logp.Printf("[DEBUG] datastoreRaw: %+v", datastoreRaw)
	if len(datastoreRaw) == 1 {
		dataStore.Type = datastoreRaw[0].(map[string]interface{})["type"].(string)
		dataStore.Version = datastoreRaw[0].(map[string]interface{})["version"].(string)
		dataStore.StorageEngine = datastoreRaw[0].(map[string]interface{})["storage_engine"].(string)
	}
	logp.Printf("[DEBUG] datastore: %+v", dataStore)
	return dataStore
}

func resourceDdsFlavors(d *schema.ResourceData) []instances.Flavor {
	var flavors []instances.Flavor
	flavorRaw := d.Get("flavor").([]interface{})
	logp.Printf("[DEBUG] flavorRaw: %+v", flavorRaw)
	for i := range flavorRaw {
		flavor := flavorRaw[i].(map[string]interface{})
		flavorReq := instances.Flavor{
			Type:     flavor["type"].(string),
			Num:      flavor["num"].(int),
			Storage:  flavor["storage"].(string),
			Size:     flavor["size"].(int),
			SpecCode: flavor["spec_code"].(string),
		}
		flavors = append(flavors, flavorReq)
	}
	logp.Printf("[DEBUG] flavors: %+v", flavors)
	return flavors
}

func resourceDdsBackupStrategy(d *schema.ResourceData) instances.BackupStrategy {
	var backupStrategy instances.BackupStrategy
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	logp.Printf("[DEBUG] backupStrategyRaw: %+v", backupStrategyRaw)
	startTime := "00:00-01:00"
	keepDays := 7
	if len(backupStrategyRaw) == 1 {
		startTime = backupStrategyRaw[0].(map[string]interface{})["start_time"].(string)
		keepDays = backupStrategyRaw[0].(map[string]interface{})["keep_days"].(int)
	}
	backupStrategy.StartTime = startTime
	backupStrategy.KeepDays = &keepDays
	logp.Printf("[DEBUG] backupStrategy: %+v", backupStrategy)
	return backupStrategy
}

func DdsInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := instances.ListInstanceOpts{
			Id: instanceID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return nil, "", err
		}
		instancesList, err := instances.ExtractInstances(allPages)
		if err != nil {
			return nil, "", err
		}

		if instancesList.TotalCount == 0 {
			var instance instances.InstanceResponse
			return instance, "deleted", nil
		}
		insts := instancesList.Instances

		status := insts[0].Status
		// wait for updating
		if status == "normal" && len(insts[0].Actions) > 0 {
			status = "updating"
		}
		return insts[0], status, nil
	}
}

func resourceDdsInstanceV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DDS client: %s ", err)
	}

	createOpts := instances.CreateOpts{
		Name:                d.Get("name").(string),
		DataStore:           resourceDdsDataStore(d),
		Region:              GetRegion(d, config),
		AvailabilityZone:    d.Get("availability_zone").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		DiskEncryptionId:    d.Get("disk_encryption_id").(string),
		Mode:                d.Get("mode").(string),
		Flavor:              resourceDdsFlavors(d),
		BackupStrategy:      resourceDdsBackupStrategy(d),
		EnterpriseProjectID: GetEnterpriseProjectID(d, config),
	}
	if d.Get("ssl").(bool) {
		createOpts.Ssl = "1"
	} else {
		createOpts.Ssl = "0"
	}
	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting instance from result: %s ", err)
	}
	logp.Printf("[DEBUG] Create : instance %s: %#v", instance.Id, instance)

	d.SetId(instance.Id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating"},
		Target:     []string{"normal"},
		Refresh:    DdsInstanceStateRefreshFunc(client, instance.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      120 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			instance.Id, err)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", instance.Id, taglist).ExtractErr(); tagErr != nil {
			return fmtp.Errorf("Error setting tags of DDS instance %s: %s", instance.Id, tagErr)
		}
	}

	return resourceDdsInstanceV3Read(d, meta)
}

func resourceDdsInstanceV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DDS client: %s", err)
	}

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return fmtp.Errorf("Error fetching DDS instance: %s", err)
	}
	instances, err := instances.ExtractInstances(allPages)
	if err != nil {
		return fmtp.Errorf("Error extracting DDS instance: %s", err)
	}
	if instances.TotalCount == 0 {
		logp.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		d.SetId("")
		return nil
	}
	insts := instances.Instances
	instance := insts[0]

	logp.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("region", instance.Region)
	d.Set("name", instance.Name)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("disk_encryption_id", instance.DiskEncryptionId)
	d.Set("mode", instance.Mode)
	d.Set("db_username", instance.DbUserName)
	d.Set("status", instance.Status)
	d.Set("port", instance.Port)
	d.Set("enterprise_project_id", instance.EnterpriseProjectID)

	sslEnable := true
	if instance.Ssl == 0 {
		sslEnable = false
	}
	d.Set("ssl", sslEnable)

	datastoreList := make([]map[string]interface{}, 0, 1)
	datastore := map[string]interface{}{
		"type":           instance.DataStore.Type,
		"version":        instance.DataStore.Version,
		"storage_engine": instance.Engine,
	}
	datastoreList = append(datastoreList, datastore)
	d.Set("datastore", datastoreList)

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	d.Set("backup_strategy", backupStrategyList)

	// save nodes attribute
	err = d.Set("nodes", flattenDdsInstanceV3Nodes(instance))
	if err != nil {
		return fmtp.Errorf("Error setting nodes of DDS instance, err: %s", err)
	}

	// save tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmtp.Errorf("Error saving tags to state for DDS instance (%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] Error fetching tags of DDS instance (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceDdsInstanceV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DDS client: %s ", err)
	}

	var opts []instances.UpdateOpt
	if d.HasChange("name") {
		opt := instances.UpdateOpt{
			Param:  "new_instance_name",
			Value:  d.Get("name").(string),
			Action: "modify-name",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("password") {
		opt := instances.UpdateOpt{
			Param:  "user_pwd",
			Value:  d.Get("password").(string),
			Action: "reset-password",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("ssl") {
		opt := instances.UpdateOpt{
			Param:  "ssl_option",
			Action: "switch-ssl",
			Method: "post",
		}
		if d.Get("ssl").(bool) {
			opt.Value = "1"
		} else {
			opt.Value = "0"
		}
		opts = append(opts, opt)
	}

	if d.HasChange("security_group_id") {
		opt := instances.UpdateOpt{
			Param:  "security_group_id",
			Value:  d.Get("security_group_id").(string),
			Action: "modify-security-group",
			Method: "post",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("backup_strategy") {
		backupStrategy := resourceDdsBackupStrategy(d)
		backupStrategy.Period = "1,2,3,4,5,6,7"
		opt := instances.UpdateOpt{
			Param:  "backup_policy",
			Value:  backupStrategy,
			Action: "backups/policy",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	r := instances.Update(client, d.Id(), opts)
	if r.Err != nil {
		return fmtp.Errorf("Error updating instance from result: %s ", r.Err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"normal"},
		Refresh:    DdsInstanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			d.Id(), err)
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", d.Id())
		if tagErr != nil {
			return fmtp.Errorf("Error updating tags of DDS instance:%s, err:%s", d.Id(), tagErr)
		}
	}

	// update flavor
	if d.HasChange("flavor") {
		for i := range d.Get("flavor").([]interface{}) {
			numIndex := fmt.Sprintf("flavor.%d.num", i)
			volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
			specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
			if d.HasChange(numIndex) {
				err := flavorNumUpdate(client, d, i)
				if err != nil {
					return err
				}
			}
			if d.HasChange(volumeSizeIndex) {
				err := flavorSizeUpdate(client, d, i)
				if err != nil {
					return err
				}
			}
			if d.HasChange(specCodeIndex) {
				err := flavorSpecCodeUpdate(client, d, i)
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceDdsInstanceV3Read(d, meta)
}

func resourceDdsInstanceV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DDS client: %s ", err)
	}

	instanceId := d.Id()
	result := instances.Delete(client, instanceId)
	if result.Err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"normal", "abnormal", "frozen", "createfail", "enlargefail", "data_disk_full"},
		Target:     []string{"deleted"},
		Refresh:    DdsInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}
	logp.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	return nil
}

func flattenDdsInstanceV3Nodes(dds instances.InstanceResponse) interface{} {
	nodesList := make([]map[string]interface{}, 0)
	for _, group := range dds.Groups {
		groupType := group.Type
		for _, Node := range group.Nodes {
			node := map[string]interface{}{
				"type":       groupType,
				"id":         Node.Id,
				"name":       Node.Name,
				"role":       Node.Role,
				"status":     Node.Status,
				"private_ip": Node.PrivateIP,
				"public_ip":  Node.PublicIP,
			}
			nodesList = append(nodesList, node)
		}
	}
	return nodesList
}

func getDdsInstanceV3ShardGroupID(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]string, error) {
	groupIDs := make([]string, 0)

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return groupIDs, fmtp.Errorf("Error fetching DDS instance: %s", err)
	}
	instances, err := instances.ExtractInstances(allPages)
	if err != nil {
		return groupIDs, fmtp.Errorf("Error extracting DDS instance: %s", err)
	}
	if instances.TotalCount == 0 {
		logp.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		return groupIDs, nil
	}
	insts := instances.Instances
	instance := insts[0]

	logp.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	for _, group := range instance.Groups {
		if group.Type == "shard" {
			groupIDs = append(groupIDs, group.Id)
		}
	}

	return groupIDs, nil

}

func getDdsInstanceV3MongosNodeID(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]string, error) {
	nodeIDs := make([]string, 0)

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return nodeIDs, fmtp.Errorf("Error fetching DDS instance: %s", err)
	}
	instances, err := instances.ExtractInstances(allPages)
	if err != nil {
		return nodeIDs, fmtp.Errorf("Error extracting DDS instance: %s", err)
	}
	if instances.TotalCount == 0 {
		logp.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		return nodeIDs, nil
	}
	insts := instances.Instances
	instance := insts[0]

	logp.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	for _, group := range instance.Groups {
		if group.Type == "mongos" {
			for _, node := range group.Nodes {
				nodeIDs = append(nodeIDs, node.Id)
			}
		}
	}

	return nodeIDs, nil

}

func flavorUpdate(client *golangsdk.ServiceClient, d *schema.ResourceData, opts []instances.UpdateOpt) error {
	r := instances.Update(client, d.Id(), opts)
	if r.Err != nil {
		return fmtp.Errorf("Error updating instance from result: %s ", r.Err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"normal"},
		Refresh:    DdsInstanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			d.Id(), err)
	}

	return nil
}

func flavorNumUpdate(client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType != "mongos" && groupType != "shard" {
		return fmtp.Errorf("Error updating instance: %s does not support adding nodes", groupType)
	}
	specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
	volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
	volumeSize := d.Get(volumeSizeIndex).(int)
	numIndex := fmt.Sprintf("flavor.%d.num", i)
	oldNumRaw, newNumRaw := d.GetChange(numIndex)
	oldNum := oldNumRaw.(int)
	newNum := newNumRaw.(int)
	if newNum < oldNum {
		return fmtp.Errorf("Error updating instance: the new num(%d) must be greater than the old num(%d)", newNum, oldNum)
	}

	var numUpdateOpts []instances.UpdateOpt

	if groupType == "mongos" {
		opt := instances.UpdateOpt{
			Param: "",
			Value: instances.UpdateNodeNumOpts{
				Type:     groupType,
				SpecCode: d.Get(specCodeIndex).(string),
				Num:      newNum - oldNum,
			},
			Action: "enlarge",
			Method: "post",
		}
		numUpdateOpts = append(numUpdateOpts, opt)
	} else {
		volume := instances.UpdateVolumeOpts{
			Size: &volumeSize,
		}

		opt := instances.UpdateOpt{
			Param: "",
			Value: instances.UpdateNodeNumOpts{
				Type:     groupType,
				SpecCode: d.Get(specCodeIndex).(string),
				Num:      newNum - oldNum,
				Volume:   &volume,
			},
			Action: "enlarge",
			Method: "post",
		}
		numUpdateOpts = append(numUpdateOpts, opt)
	}
	err := flavorUpdate(client, d, numUpdateOpts)
	if err != nil {
		return err
	}
	return nil
}

func flavorSizeUpdate(client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
	oldSizeRaw, newSizeRaw := d.GetChange(volumeSizeIndex)
	oldSize := oldSizeRaw.(int)
	newSize := newSizeRaw.(int)
	if newSize < oldSize {
		return fmtp.Errorf("Error updating instance: the new size(%d) must be greater than the old size(%d)", newSize, oldSize)
	}
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType != "replica" && groupType != "single" && groupType != "shard" {
		return fmtp.Errorf("Error updating instance: %s does not support scaling up storage space", groupType)
	}

	if groupType == "shard" {
		groupIDs, err := getDdsInstanceV3ShardGroupID(client, d)
		if err != nil {
			return err
		}

		for _, groupID := range groupIDs {
			var sizeUpdateOpts []instances.UpdateOpt
			opt := instances.UpdateOpt{
				Param: "volume",
				Value: instances.UpdateVolumeOpts{
					GroupID: groupID,
					Size:    &newSize,
				},
				Action: "enlarge-volume",
				Method: "post",
			}
			sizeUpdateOpts = append(sizeUpdateOpts, opt)
			err := flavorUpdate(client, d, sizeUpdateOpts)
			if err != nil {
				return err
			}
		}
	} else {
		var sizeUpdateOpts []instances.UpdateOpt
		opt := instances.UpdateOpt{
			Param: "volume",
			Value: instances.UpdateVolumeOpts{
				Size: &newSize,
			},
			Action: "enlarge-volume",
			Method: "post",
		}
		sizeUpdateOpts = append(sizeUpdateOpts, opt)
		err := flavorUpdate(client, d, sizeUpdateOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func flavorSpecCodeUpdate(client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType == "config" {
		return fmtp.Errorf("Error updating instance: %s does not support updating spec_code", groupType)
	}
	if groupType == "mongos" {
		nodeIDs, err := getDdsInstanceV3MongosNodeID(client, d)
		if err != nil {
			return err
		}
		for _, ID := range nodeIDs {
			var specUpdateOpts []instances.UpdateOpt
			opt := instances.UpdateOpt{
				Param: "resize",
				Value: instances.UpdateSpecOpts{
					TargetType:     "mongos",
					TargetID:       ID,
					TargetSpecCode: d.Get(specCodeIndex).(string),
				},
				Action: "resize",
				Method: "post",
			}
			specUpdateOpts = append(specUpdateOpts, opt)
			err := flavorUpdate(client, d, specUpdateOpts)
			if err != nil {
				return err
			}
		}
	} else if groupType == "shard" {
		groupIDs, err := getDdsInstanceV3ShardGroupID(client, d)
		if err != nil {
			return err
		}

		for _, ID := range groupIDs {
			var specUpdateOpts []instances.UpdateOpt
			opt := instances.UpdateOpt{
				Param: "resize",
				Value: instances.UpdateSpecOpts{
					TargetType:     "shard",
					TargetID:       ID,
					TargetSpecCode: d.Get(specCodeIndex).(string),
				},
				Action: "resize",
				Method: "post",
			}
			specUpdateOpts = append(specUpdateOpts, opt)
			err := flavorUpdate(client, d, specUpdateOpts)
			if err != nil {
				return err
			}
		}
	} else {
		var specUpdateOpts []instances.UpdateOpt
		opt := instances.UpdateOpt{
			Param: "resize",
			Value: instances.UpdateSpecOpts{
				TargetID:       d.Id(),
				TargetSpecCode: d.Get(specCodeIndex).(string),
			},
			Action: "resize",
			Method: "post",
		}
		specUpdateOpts = append(specUpdateOpts, opt)
		err := flavorUpdate(client, d, specUpdateOpts)
		if err != nil {
			return err
		}
	}
	return nil
}
