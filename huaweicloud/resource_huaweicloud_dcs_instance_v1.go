package huaweicloud

import (
	"strconv"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dcs/v1/instances"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/whitelists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceDcsInstanceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDcsInstancesV1Create,
		Read:   resourceDcsInstancesV1Read,
		Update: resourceDcsInstancesV1Update,
		Delete: resourceDcsInstancesV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Redis", "Memcached",
				}, true),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"whitelists"},
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ForceNew:  true,
			},
			"maintain_begin": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"maintain_end"},
				Computed:     true,
			},
			"maintain_end": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"maintain_begin"},
				Computed:     true,
			},
			"backup_policy": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"backup_type", "begin_at", "period_type", "backup_at", "save_days"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"save_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"begin_at": {
							Type:     schema.TypeString,
							Required: true,
						},
						"period_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"backup_at": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"whitelist_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"whitelists": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 4,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"tags": tagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"internal_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Deprecated
			"save_days": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use `backup_policy` instead",
			},
			"backup_type": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use `backup_policy` instead",
			},
			"begin_at": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_type", "backup_at", "save_days", "backup_type"},
				Deprecated:   "Please use `backup_policy` instead",
			},
			"period_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"begin_at", "backup_at", "save_days", "backup_type"},
				Deprecated:   "Please use `backup_policy` instead",
			},
			"backup_at": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_type", "begin_at", "save_days", "backup_type"},
				Deprecated:   "Please use `backup_policy` instead",
				Elem:         &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func formatAts(src []interface{}) []int {
	res := make([]int, len(src))
	for i, at := range src {
		res[i] = at.(int)
	}
	return res
}

func getInstanceBackupPolicy(d *schema.ResourceData) *instances.InstanceBackupPolicy {
	var instanceBackupPolicy *instances.InstanceBackupPolicy
	if _, ok := d.GetOk("backup_policy"); !ok { // deprecated branch
		if v, ok := d.GetOk("backup_at"); ok {
			backupAts := v.([]interface{})
			return &instances.InstanceBackupPolicy{
				SaveDays:   d.Get("save_days").(int),
				BackupType: d.Get("backup_type").(string),
				PeriodicalBackupPlan: instances.PeriodicalBackupPlan{
					BeginAt:    d.Get("begin_at").(string),
					PeriodType: d.Get("period_type").(string),
					BackupAt:   formatAts(backupAts),
				},
			}
		}
		// neither backup_policy nor backup_at is specified
		return nil
	}

	backupPolicyList := d.Get("backup_policy").([]interface{})
	if len(backupPolicyList) == 0 {
		return nil
	}
	backupPolicy := backupPolicyList[0].(map[string]interface{})
	backupAts := backupPolicy["backup_at"].([]interface{})
	instanceBackupPolicy = &instances.InstanceBackupPolicy{
		SaveDays:   backupPolicy["save_days"].(int),
		BackupType: backupPolicy["backup_type"].(string),
		PeriodicalBackupPlan: instances.PeriodicalBackupPlan{
			BeginAt:    backupPolicy["begin_at"].(string),
			PeriodType: backupPolicy["period_type"].(string),
			BackupAt:   formatAts(backupAts),
		},
	}

	return instanceBackupPolicy
}

func resourceDcsInstancesCheck(d *schema.ResourceData) error {
	engineVersion := d.Get("engine_version").(string)
	secGroupID := d.Get("security_group_id").(string)

	// check for Redis 4.0 and 5.0
	if engineVersion == "4.0" || engineVersion == "5.0" {
		if secGroupID != "" {
			return fmtp.Errorf("security_group_id is not supported for Redis 4.0 and 5.0. please configure the whitelists alternatively")
		}
	} else {
		// check for Memcached and Redis 3.0
		if secGroupID == "" {
			return fmtp.Errorf("security_group_id is mandatory for this DCS instance")
		}
	}

	return nil
}

func getDcsInstanceWhitelist(d *schema.ResourceData) whitelists.WhitelistOpts {
	groupsRaw := d.Get("whitelists").(*schema.Set).List()
	whitelitGroups := make([]whitelists.WhitelistGroupOpts, len(groupsRaw))
	for i, v := range groupsRaw {
		groups := v.(map[string]interface{})

		ipRaw := groups["ip_address"].([]interface{})
		ipList := make([]string, len(ipRaw))
		for j, ip := range ipRaw {
			ipList[j] = ip.(string)
		}

		whitelitGroups[i] = whitelists.WhitelistGroupOpts{
			GroupName: groups["group_name"].(string),
			IPList:    ipList,
		}
	}

	enable := d.Get("whitelist_enable").(bool)
	if len(groupsRaw) == 0 {
		enable = false
	}

	return whitelists.WhitelistOpts{
		Enable: &enable,
		Groups: whitelitGroups,
	}
}

func flattenDcsInstanceWhitelist(object *whitelists.Whitelist) interface{} {
	whilteList := make([]map[string]interface{}, len(object.Groups))
	for i, group := range object.Groups {
		whilteList[i] = map[string]interface{}{
			"group_name": group.GroupName,
			"ip_address": group.IPList,
		}
	}
	return whilteList
}

func resourceDcsInstancesV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dcsV1Client, err := config.DcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dcs instance v1 client: %s", err)
	}

	if err := resourceDcsInstancesCheck(d); err != nil {
		return err
	}

	noPasswordAccess := "true"
	if d.Get("access_user").(string) != "" || d.Get("password").(string) != "" {
		noPasswordAccess = "false"
	}
	createOpts := &instances.CreateOps{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Engine:                d.Get("engine").(string),
		EngineVersion:         d.Get("engine_version").(string),
		Capacity:              d.Get("capacity").(float64),
		NoPasswordAccess:      noPasswordAccess,
		AccessUser:            d.Get("access_user").(string),
		VPCID:                 d.Get("vpc_id").(string),
		SecurityGroupID:       d.Get("security_group_id").(string),
		SubnetID:              d.Get("subnet_id").(string),
		AvailableZones:        getAllAvailableZones(d),
		ProductID:             d.Get("product_id").(string),
		InstanceBackupPolicy:  getInstanceBackupPolicy(d),
		MaintainBegin:         d.Get("maintain_begin").(string),
		MaintainEnd:           d.Get("maintain_end").(string),
		EnterpriseProjectID:   GetEnterpriseProjectID(d, config),
		EnterpriseProjectName: d.Get("enterprise_project_name").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	v, err := instances.Create(dcsV1Client, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud instance: %s", err)
	}
	logp.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"RUNNING"},
		Refresh:    DcsInstancesV1StateRefreshFunc(dcsV1Client, v.InstanceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	// set whitelist
	dcsV2Client, err := config.DcsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dcs instance v2 client: %s", err)
	}
	whitelistOpts := getDcsInstanceWhitelist(d)
	logp.Printf("[DEBUG] Create whitelist options: %#v", whitelistOpts)

	if *whitelistOpts.Enable {
		err = whitelists.Put(dcsV2Client, d.Id(), whitelistOpts).ExtractErr()
		if err != nil {
			return fmtp.Errorf("Error creating whitelist for instance (%s): %s", d.Id(), err)
		}
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(dcsV2Client, "dcs", v.InstanceID, taglist).ExtractErr(); tagErr != nil {
			return fmtp.Errorf("Error setting tags of DCS instance %s: %s", v.InstanceID, tagErr)
		}
	}

	return resourceDcsInstancesV1Read(d, meta)
}

func resourceDcsInstancesV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	dcsV1Client, err := config.DcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dcs instance v1 client: %s", err)
	}
	v, err := instances.Get(dcsV1Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "DCS instance")
	}

	logp.Printf("[DEBUG] Dcs instance %s: %+v", d.Id(), v)

	d.SetId(v.InstanceID)
	d.Set("name", v.Name)
	d.Set("engine", v.Engine)
	d.Set("engine_version", v.EngineVersion)
	d.Set("used_memory", v.UsedMemory)
	d.Set("max_memory", v.MaxMemory)
	d.Set("ip", v.IP)
	d.Set("port", v.Port)
	d.Set("status", v.Status)
	d.Set("description", v.Description)
	d.Set("resource_spec_code", v.ResourceSpecCode)
	d.Set("internal_version", v.InternalVersion)
	d.Set("vpc_id", v.VPCID)
	d.Set("vpc_name", v.VPCName)
	d.Set("created_at", v.CreatedAt)
	d.Set("product_id", v.ProductID)
	d.Set("security_group_id", v.SecurityGroupID)
	d.Set("security_group_name", v.SecurityGroupName)
	d.Set("subnet_id", v.SubnetID)
	d.Set("subnet_name", v.SubnetName)
	d.Set("user_id", v.UserID)
	d.Set("user_name", v.UserName)
	d.Set("order_id", v.OrderID)
	d.Set("maintain_begin", v.MaintainBegin)
	d.Set("maintain_end", v.MaintainEnd)
	d.Set("access_user", v.AccessUser)
	d.Set("enterprise_project_id", v.EnterpriseProjectID)
	d.Set("enterprise_project_name", v.EnterpriseProjectName)

	// set capacity by Capacity and CapacityMinor
	var capacity float64 = float64(v.Capacity)
	if v.CapacityMinor != "" {
		if minor, err := strconv.ParseFloat(v.CapacityMinor, 64); err == nil {
			capacity += minor
		}
	}
	d.Set("capacity", capacity)

	dcsV2Client, err := config.DcsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dcs instance v2 client: %s", err)
	}
	object, err := whitelists.Get(dcsV2Client, d.Id()).Extract()

	enable := object.Enable
	// change enable to true when none whitelist groups exists
	if len(object.Groups) == 0 {
		enable = true
	}
	d.Set("whitelist_enable", enable)
	err = d.Set("whitelists", flattenDcsInstanceWhitelist(object))
	if err != nil {
		return fmtp.Errorf("Error setting whitelists for DCS instance, err: %s", err)
	}

	// set tags
	if resourceTags, err := tags.Get(dcsV2Client, "instances", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmtp.Errorf("[DEBUG] Error saving tag to state for DCS instance (%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] fetching tags of DCS instance failed: %s", err)
	}

	return nil
}

func resourceDcsInstancesV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	if err := resourceDcsInstancesCheck(d); err != nil {
		return err
	}

	//lintignore:R019
	if d.HasChanges("name", "description", "security_group_id", "maintain_begin", "maintain_end", "backup_policy") {
		dcsV1Client, err := config.DcsV1Client(GetRegion(d, config))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud dcs instance v1 client: %s", err)
		}

		description := d.Get("description").(string)
		updateOpts := instances.UpdateOpts{
			Name:                 d.Get("name").(string),
			Description:          &description,
			MaintainBegin:        d.Get("maintain_begin").(string),
			MaintainEnd:          d.Get("maintain_end").(string),
			SecurityGroupID:      d.Get("security_group_id").(string),
			InstanceBackupPolicy: getInstanceBackupPolicy(d),
		}

		err = instances.Update(dcsV1Client, d.Id(), updateOpts).Err
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud Dcs Instance: %s", err)
		}
	}

	if d.HasChanges("whitelists", "tags") {
		dcsV2Client, err := config.DcsV2Client(GetRegion(d, config))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud dcs instance v2 client: %s", err)
		}

		if d.HasChange("whitelists") {
			whitelistOpts := getDcsInstanceWhitelist(d)
			logp.Printf("[DEBUG] update whitelist options: %#v", whitelistOpts)

			err = whitelists.Put(dcsV2Client, d.Id(), whitelistOpts).ExtractErr()
			if err != nil {
				return fmtp.Errorf("Error updating whitelist for instance (%s): %s", d.Id(), err)
			}
		}

		// update tags
		tagErr := utils.UpdateResourceTags(dcsV2Client, d, "dcs", d.Id())
		if tagErr != nil {
			return fmtp.Errorf("Error updating tags of DCS instance:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceDcsInstancesV1Read(d, meta)
}

func resourceDcsInstancesV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dcsV1Client, err := config.DcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dcs instance v1 client: %s", err)
	}

	_, err = instances.Get(dcsV1Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "instance")
	}

	err = instances.Delete(dcsV1Client, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud instance: %s", err)
	}

	// Wait for the instance to delete before moving on.
	logp.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING", "RUNNING"},
		Target:     []string{"DELETED"},
		Refresh:    DcsInstancesV1StateRefreshFunc(dcsV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	logp.Printf("[DEBUG] Dcs instance %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func DcsInstancesV1StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		return v, v.Status, nil
	}
}
