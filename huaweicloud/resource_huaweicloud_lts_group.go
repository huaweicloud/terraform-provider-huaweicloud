package huaweicloud

import (
	"github.com/chnsz/golangsdk/openstack/lts/huawei/loggroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceLTSGroupV2() *schema.Resource {
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
	config := meta.(*config.Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	createOpts := &loggroups.CreateOpts{
		LogGroupName: d.Get("group_name").(string),
		TTL:          d.Get("ttl_in_days").(int),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	groupCreate, err := loggroups.Create(client, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating log group: %s", err)
	}

	d.SetId(groupCreate.ID)
	return resourceGroupV2Read(d, meta)
}

func resourceGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	groups, err := loggroups.List(client).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting HuaweiCloud log group list: %s", err)
	}
	for _, group := range groups.LogGroups {
		if group.ID == d.Id() {
			d.SetId(group.ID)
			d.Set("group_name", group.Name)
			d.Set("ttl_in_days", group.TTLinDays)
			return nil
		}
	}

	return fmtp.Errorf("Error HuaweiCloud log group %s: No Found", d.Id())
}

func resourceGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	updateOpts := &loggroups.UpdateOpts{
		TTL: d.Get("ttl_in_days").(int),
	}

	logp.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = loggroups.Update(client, updateOpts, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error update log group: %s", err)
	}

	return resourceGroupV2Read(d, meta)
}

func resourceGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	err = loggroups.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting log group")
	}

	d.SetId("")
	return nil
}
