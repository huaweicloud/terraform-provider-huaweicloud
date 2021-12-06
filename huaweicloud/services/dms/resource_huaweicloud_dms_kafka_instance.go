package dms

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceDmsKafkaInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaInstanceCreate,
		ReadContext:   resourceDmsKafkaInstanceRead,
		UpdateContext: resourceDmsKafkaInstanceUpdate,
		DeleteContext: resourceDmsKafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"100MB", "300MB", "900MB", "1200MB",
				}, false),
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"password",
				},
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ForceNew:  true,
				RequiredWith: []string{
					"access_user",
				},
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_ip_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"produce_reject", "time_base",
				}, false),
			},
			"dumping": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_auto_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partition_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_public_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
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
			"manegement_connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsKafkaPublicIpIDs(d *schema.ResourceData) (string, error) {
	bandwidth := d.Get("bandwidth").(string)
	publicIpIDsRaw := d.Get("public_ip_ids").([]interface{})

	IdNumMap := map[string]int{
		"100MB":  3,
		"300MB":  3,
		"600MB":  4,
		"1200MB": 8,
	}

	if IdNumMap[bandwidth] != len(publicIpIDsRaw) {
		return "", fmtp.Errorf("error creating HuaweiCloud DMS kafka instance: "+
			"%d public ip IDs needed when bandwidth is set to %s, but got %d",
			IdNumMap[bandwidth], bandwidth, len(publicIpIDsRaw))
	}

	publicIpIDs := make([]string, 0, len(publicIpIDsRaw))
	for _, v := range publicIpIDsRaw {
		publicIpIDs = append(publicIpIDs, v.(string))
	}
	return strings.Join(publicIpIDs, ","), nil
}

func resourceDmsKafkaInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud dms instance client: %s", err)
	}

	partitionNumMap := map[string]int{
		"100MB":  300,
		"300MB":  900,
		"600MB":  1800,
		"1200MB": 1800,
	}

	ssl_enable := false
	if d.Get("access_user").(string) != "" && d.Get("password").(string) != "" {
		ssl_enable = true
	}

	rawZones := d.Get("available_zones").([]interface{})
	utils.ExpandToStringList(rawZones)
	zones := utils.ExpandToStringList(rawZones)

	createOpts := &instances.CreateOps{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Engine:              "kafka",
		EngineVersion:       d.Get("engine_version").(string),
		Specification:       d.Get("bandwidth").(string),
		StorageSpace:        d.Get("storage_space").(int),
		PartitionNum:        partitionNumMap[d.Get("bandwidth").(string)],
		AccessUser:          d.Get("access_user").(string),
		VPCID:               d.Get("vpc_id").(string),
		SecurityGroupID:     d.Get("security_group_id").(string),
		SubnetID:            d.Get("network_id").(string),
		AvailableZones:      zones,
		ProductID:           d.Get("product_id").(string),
		KafkaManagerUser:    d.Get("manager_user").(string),
		MaintainBegin:       d.Get("maintain_begin").(string),
		MaintainEnd:         d.Get("maintain_end").(string),
		SslEnable:           ssl_enable,
		RetentionPolicy:     d.Get("retention_policy").(string),
		ConnectorEnalbe:     d.Get("dumping").(bool),
		EnableAutoTopic:     d.Get("enable_auto_topic").(bool),
		StorageSpecCode:     d.Get("storage_spec_code").(string),
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, config),
	}

	if _, ok := d.GetOk("public_ip_ids"); ok {
		publicIpIDs, err := resourceDmsKafkaPublicIpIDs(d)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.EnablePublicIP = true
		createOpts.PublicIpID = publicIpIDs
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		createOpts.Tags = taglist
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)
	createOpts.KafkaManagerPassword = d.Get("manager_password").(string)

	v, err := instances.Create(dmsV2Client, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud dms kafka instance: %s", err)
	}
	logp.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      DmsKafkaInstanceStateRefreshFunc(dmsV2Client, v.InstanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        300 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.DiagErrorf(
			"error waiting for instance (%s) to become ready: %s",
			v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	return resourceDmsKafkaInstanceRead(ctx, d, meta)
}

func resourceDmsKafkaInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)

	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud dms instance client: %s", err)
	}
	v, err := instances.Get(dmsV2Client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS instance")
	}

	logp.Printf("[DEBUG] Dms kafka instance %s: %+v", d.Id(), v)

	d.SetId(v.InstanceID)
	d.Set("region", config.GetRegion(d))
	d.Set("name", v.Name)
	d.Set("description", v.Description)
	d.Set("engine", v.Engine)
	d.Set("engine_version", v.EngineVersion)
	d.Set("bandwidth", v.Specification)
	// storage_space indicates total_storage_space while creating
	// set value of total_storage_space to storage_space to keep consistent
	d.Set("storage_space", v.TotalStorageSpace)

	partitionNum, _ := strconv.ParseInt(v.PartitionNum, 10, 64)
	d.Set("partition_num", partitionNum)

	d.Set("vpc_id", v.VPCID)
	d.Set("security_group_id", v.SecurityGroupID)
	d.Set("network_id", v.SubnetID)
	d.Set("available_zones", v.AvailableZones)
	d.Set("product_id", v.ProductID)
	d.Set("manager_user", v.KafkaManagerUser)
	d.Set("maintain_begin", v.MaintainBegin)
	d.Set("maintain_end", v.MaintainEnd)
	d.Set("enable_public_ip", v.EnablePublicIP)
	d.Set("ssl_enable", v.SslEnable)
	d.Set("retention_policy", v.RetentionPolicy)
	d.Set("dumping", v.ConnectorEnalbe)
	d.Set("enable_auto_topic", v.EnableAutoTopic)
	d.Set("storage_spec_code", v.StorageSpecCode)
	d.Set("enterprise_project_id", v.EnterpriseProjectID)
	d.Set("used_storage_space", v.UsedStorageSpace)
	d.Set("connect_address", v.ConnectAddress)
	d.Set("port", v.Port)
	d.Set("status", v.Status)
	d.Set("resource_spec_code", v.ResourceSpecCode)
	d.Set("user_id", v.UserID)
	d.Set("user_name", v.UserName)
	d.Set("manegement_connect_address", v.ManagementConnectAddress)
	d.Set("type", v.Type)
	d.Set("access_user", v.AccessUser)

	// set tags
	engine := "kafka"
	if resourceTags, err := tags.Get(dmsV2Client, engine, d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmtp.DiagErrorf("error saving tags to state for dms kafka instance (%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] error fetching tags of dms kafka instance (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceDmsKafkaInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud dms instance client: %s", err)
	}

	//lintignore:R019
	if d.HasChanges("name", "description", "maintain_begin", "maintain_end",
		"security_group_id", "retention_policy", "enterprise_project_id") {
		description := d.Get("description").(string)
		updateOpts := instances.UpdateOpts{
			Description:         &description,
			MaintainBegin:       d.Get("maintain_begin").(string),
			MaintainEnd:         d.Get("maintain_end").(string),
			SecurityGroupID:     d.Get("security_group_id").(string),
			RetentionPolicy:     d.Get("retention_policy").(string),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		}

		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}

		err = instances.Update(dmsV2Client, d.Id(), updateOpts).Err
		if err != nil {
			return fmtp.DiagErrorf("error updating HuaweiCloud Dms kafka Instance: %s", err)
		}
	}

	if d.HasChange("tags") {
		// update tags
		engine := "kafka"
		tagErr := utils.UpdateResourceTags(dmsV2Client, d, engine, d.Id())
		if tagErr != nil {
			return fmtp.DiagErrorf("error updating tags of dms kafka instance:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceDmsKafkaInstanceRead(ctx, d, meta)
}

func resourceDmsKafkaInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud dms instance client: %s", err)
	}

	err = instances.Delete(dmsV2Client, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("error deleting HuaweiCloud instance: %s", err)
	}

	// Wait for the instance to delete before moving on.
	logp.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING"},
		Target:       []string{"DELETED"},
		Refresh:      DmsKafkaInstanceStateRefreshFunc(dmsV2Client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        90 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.DiagErrorf(
			"error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	logp.Printf("[DEBUG] Dms instance %s deactivated", d.Id())
	d.SetId("")
	return nil
}

func DmsKafkaInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
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
