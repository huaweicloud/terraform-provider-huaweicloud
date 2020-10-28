package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/instances"
)

func resourceDmsInstancesV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDmsInstancesV1Create,
		Read:   resourceDmsInstancesV1Read,
		Update: resourceDmsInstancesV1Update,
		Delete: resourceDmsInstancesV1Delete,
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
					"rabbitmq", "kafka",
				}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"partition_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
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
			"connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDmsInstancesV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dmsV1Client, err := config.dmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
	}

	ssl_enable := false
	if d.Get("access_user").(string) != "" || d.Get("password").(string) != "" {
		ssl_enable = true
	}
	createOpts := &instances.CreateOps{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Engine:          d.Get("engine").(string),
		EngineVersion:   d.Get("engine_version").(string),
		StorageSpace:    d.Get("storage_space").(int),
		Password:        d.Get("password").(string),
		AccessUser:      d.Get("access_user").(string),
		VPCID:           d.Get("vpc_id").(string),
		SecurityGroupID: d.Get("security_group_id").(string),
		SubnetID:        d.Get("subnet_id").(string),
		AvailableZones:  getAllAvailableZones(d),
		ProductID:       d.Get("product_id").(string),
		MaintainBegin:   d.Get("maintain_begin").(string),
		MaintainEnd:     d.Get("maintain_end").(string),
		PartitionNum:    d.Get("partition_num").(int),
		Specification:   d.Get("specification").(string),
		StorageSpecCode: d.Get("storage_spec_code").(string),
		SslEnable:       ssl_enable,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := instances.Create(dmsV1Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud instance: %s", err)
	}
	log.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"RUNNING"},
		Refresh:    DmsInstancesV1StateRefreshFunc(dmsV1Client, v.InstanceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		dmsV2Client, err := config.dmsV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud dms instance v2 client: %s", err)
		}

		taglist := expandResourceTags(tagRaw)
		engine := d.Get("engine").(string)
		if tagErr := tags.Create(dmsV2Client, engine, v.InstanceID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of dms instance %s: %s", v.InstanceID, tagErr)
		}
	}

	return resourceDmsInstancesV1Read(d, meta)
}

func resourceDmsInstancesV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	dmsV1Client, err := config.dmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
	}
	v, err := instances.Get(dmsV1Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "DMS instance")
	}

	log.Printf("[DEBUG] Dms instance %s: %+v", d.Id(), v)

	d.SetId(v.InstanceID)
	d.Set("name", v.Name)
	d.Set("engine", v.Engine)
	d.Set("engine_version", v.EngineVersion)
	d.Set("specification", v.Specification)
	d.Set("used_storage_space", v.UsedStorageSpace)
	d.Set("connect_address", v.ConnectAddress)
	d.Set("port", v.Port)
	d.Set("status", v.Status)
	d.Set("description", v.Description)
	d.Set("instance_id", v.InstanceID)
	d.Set("resource_spec_code", v.ResourceSpecCode)
	d.Set("type", v.Type)
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

	// set tags
	dmsV2Client, err := config.dmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud dms instance v2 client: %s", err)
	}

	engine := d.Get("engine").(string)
	resourceTags, err := tags.Get(dmsV2Client, engine, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching tags of dms instance: %s", err)
	}
	tagmap := tagsToMap(resourceTags.Tags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("[DEBUG] Error saving tag to state for dms instance (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceDmsInstancesV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChanges("name", "description", "maintain_begin", "maintain_end", "security_group_id") {
		dmsV1Client, err := config.dmsV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud dms instance client: %s", err)
		}

		var updateOpts instances.UpdateOpts
		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}
		if d.HasChange("description") {
			description := d.Get("description").(string)
			updateOpts.Description = &description
		}
		if d.HasChange("maintain_begin") {
			maintain_begin := d.Get("maintain_begin").(string)
			updateOpts.MaintainBegin = maintain_begin
		}
		if d.HasChange("maintain_end") {
			maintain_end := d.Get("maintain_end").(string)
			updateOpts.MaintainEnd = maintain_end
		}
		if d.HasChange("security_group_id") {
			security_group_id := d.Get("security_group_id").(string)
			updateOpts.SecurityGroupID = security_group_id
		}

		err = instances.Update(dmsV1Client, d.Id(), updateOpts).Err
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud Dms Instance: %s", err)
		}
	}

	if d.HasChange("tags") {
		dmsV2Client, err := config.dmsV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud dms instance v2 client: %s", err)
		}
		// update tags
		engine := d.Get("engine").(string)
		tagErr := UpdateResourceTags(dmsV2Client, d, engine, d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of dms instance:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceDmsInstancesV1Read(d, meta)
}

func resourceDmsInstancesV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dmsV1Client, err := config.dmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
	}

	_, err = instances.Get(dmsV1Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "instance")
	}

	err = instances.Delete(dmsV1Client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud instance: %s", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING", "RUNNING"},
		Target:     []string{"DELETED"},
		Refresh:    DmsInstancesV1StateRefreshFunc(dmsV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	log.Printf("[DEBUG] Dms instance %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func DmsInstancesV1StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
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
