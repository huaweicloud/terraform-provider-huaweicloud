package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/lts/huawei/loggroups"
)

func resourceLTSGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupV2Create,
		Read:   resourceGroupV2Read,
		Update: resourceGroupV2Update,
		Delete: resourceGroupV2Delete,
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
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.ltsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	createOpts := &loggroups.CreateOpts{
		LogGroupName: d.Get("group_name").(string),
		TTL:          d.Get("ttl_in_days").(int),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	groupCreate, err := loggroups.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating log group: %s", err)
	}

	d.SetId(groupCreate.ID)
	return resourceGroupV2Read(d, meta)
}

func resourceGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.ltsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	groups, err := loggroups.List(client).Extract()
	if err != nil {
		return fmt.Errorf("Error getting HuaweiCloud log group list: %s", err)
	}
	for _, group := range groups.LogGroups {
		if group.ID == d.Id() {
			d.SetId(group.ID)
			d.Set("group_name", group.Name)
			d.Set("ttl_in_days", group.TTLinDays)
			return nil
		}
	}

	return fmt.Errorf("Error HuaweiCloud log group %s: No Found", d.Id())
}

func resourceGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.ltsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	updateOpts := &loggroups.UpdateOpts{
		TTL: d.Get("ttl_in_days").(int),
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = loggroups.Update(client, updateOpts, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error update log group: %s", err)
	}

	return resourceGroupV2Read(d, meta)
}

func resourceGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.ltsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	err = loggroups.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting log group")
	}

	d.SetId("")
	return nil
}
