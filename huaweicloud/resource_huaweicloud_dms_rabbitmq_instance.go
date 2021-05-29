package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/dms/v2/rabbitmq/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func resourceDmsRabbitmqInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceDmsRabbitmqInstanceCreate,
		Read:   resourceDmsRabbitmqInstanceRead,
		Update: resourceDmsRabbitmqInstanceUpdate,
		Delete: resourceDmsRabbitmqInstanceDelete,
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
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "3.7.17",
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
			"access_user": {
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
			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_public_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
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
			"connect_address": {
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

func resourceDmsRabbitmqInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud dms instance client: %s", err)
	}

	createOpts := &instances.CreateOps{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Engine:              "rabbitmq",
		EngineVersion:       d.Get("engine_version").(string),
		StorageSpace:        d.Get("storage_space").(int),
		AccessUser:          d.Get("access_user").(string),
		VPCID:               d.Get("vpc_id").(string),
		SecurityGroupID:     d.Get("security_group_id").(string),
		SubnetID:            d.Get("network_id").(string),
		AvailableZones:      getAllAvailableZones(d),
		ProductID:           d.Get("product_id").(string),
		MaintainBegin:       d.Get("maintain_begin").(string),
		MaintainEnd:         d.Get("maintain_end").(string),
		SslEnable:           d.Get("ssl_enable").(bool),
		StorageSpecCode:     d.Get("storage_spec_code").(string),
		EnterpriseProjectID: GetEnterpriseProjectID(d, config),
	}

	if v, ok := d.GetOk("public_ip_id"); ok {
		createOpts.EnablePublicIP = true
		createOpts.PublicIpID = v.(string)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		createOpts.Tags = taglist
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	v, err := instances.Create(dmsV2Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud dms rabbitmq instance: %s", err)
	}
	log.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      DmsRabbitmqInstanceStateRefreshFunc(dmsV2Client, v.InstanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        500 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"error waiting for instance (%s) to become ready: %s",
			v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	return resourceDmsRabbitmqInstanceRead(d, meta)
}

func resourceDmsRabbitmqInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	dmsV2Client, err := config.DmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud dms instance client: %s", err)
	}
	v, err := instances.Get(dmsV2Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "DMS instance")
	}

	log.Printf("[DEBUG] Dms rabbitmq instance %s: %+v", d.Id(), v)

	d.SetId(v.InstanceID)
	d.Set("region", GetRegion(d, config))
	d.Set("name", v.Name)
	d.Set("description", v.Description)
	d.Set("engine", v.Engine)
	d.Set("engine_version", v.EngineVersion)
	d.Set("specification", v.Specification)
	// storage_space indicates total_storage_space while creating
	// set value of total_storage_space to storage_space to keep consistent
	d.Set("storage_space", v.TotalStorageSpace)

	d.Set("vpc_id", v.VPCID)
	d.Set("security_group_id", v.SecurityGroupID)
	d.Set("network_id", v.SubnetID)
	d.Set("available_zones", v.AvailableZones)
	d.Set("product_id", v.ProductID)
	d.Set("maintain_begin", v.MaintainBegin)
	d.Set("maintain_end", v.MaintainEnd)
	d.Set("enable_public_ip", v.EnablePublicIP)
	d.Set("public_ip_id", v.PublicIPID)
	d.Set("ssl_enable", v.SslEnable)
	d.Set("storage_spec_code", v.StorageSpecCode)
	d.Set("enterprise_project_id", v.EnterpriseProjectID)
	d.Set("used_storage_space", v.UsedStorageSpace)
	d.Set("connect_address", v.ConnectAddress)
	d.Set("manegement_connect_address", v.ManagementConnectAddress)
	d.Set("port", v.Port)
	d.Set("status", v.Status)
	d.Set("resource_spec_code", v.ResourceSpecCode)
	d.Set("user_id", v.UserID)
	d.Set("user_name", v.UserName)
	d.Set("type", v.Type)
	d.Set("access_user", v.AccessUser)

	// set tags
	engine := "rabbitmq"
	if resourceTags, err := tags.Get(dmsV2Client, engine, d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmt.Errorf("error saving tags to state for dms rabbitmq instance (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of dms rabbitmq instance (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceDmsRabbitmqInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud dms instance client: %s", err)
	}

	//lintignore:R019
	if d.HasChanges("name", "description", "maintain_begin", "maintain_end",
		"security_group_id", "public_ip_id", "enterprise_project_id") {
		description := d.Get("description").(string)
		updateOpts := instances.UpdateOpts{
			Description:         &description,
			MaintainBegin:       d.Get("maintain_begin").(string),
			MaintainEnd:         d.Get("maintain_end").(string),
			SecurityGroupID:     d.Get("security_group_id").(string),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		}

		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}

		if d.HasChange("public_ip_id") {
			if v, ok := d.GetOk("public_ip_id"); ok {
				enablePublicIP := true
				updateOpts.EnablePublicIP = &enablePublicIP
				updateOpts.PublicIpID = v.(string)
			} else {
				enablePublicIP := false
				updateOpts.EnablePublicIP = &enablePublicIP
			}
		}

		err = instances.Update(dmsV2Client, d.Id(), updateOpts).Err
		if err != nil {
			return fmt.Errorf("error updating HuaweiCloud Dms rabbitmq Instance: %s", err)
		}
	}

	if d.HasChange("tags") {
		// update tags
		engine := "rabbitmq"
		tagErr := utils.UpdateResourceTags(dmsV2Client, d, engine, d.Id())
		if tagErr != nil {
			return fmt.Errorf("error updating tags of dms rabbitmq instance:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceDmsRabbitmqInstanceRead(d, meta)
}

func resourceDmsRabbitmqInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud dms instance client: %s", err)
	}

	err = instances.Delete(dmsV2Client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting HuaweiCloud instance: %s", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING"},
		Target:       []string{"DELETED"},
		Refresh:      DmsRabbitmqInstanceStateRefreshFunc(dmsV2Client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        90 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	log.Printf("[DEBUG] Dms instance %s deactivated", d.Id())
	d.SetId("")
	return nil
}

func DmsRabbitmqInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
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
