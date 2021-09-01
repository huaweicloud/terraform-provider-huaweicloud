package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceMemberV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceMemberV3Create,
		Read:   resourceMemberV3Read,
		Update: resourceMemberV3Update,
		Delete: resourceMemberV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			},

			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 {
						errors = append(errors, fmtp.Errorf(
							"Only numbers greater than 0 are supported values for 'weight'"))
					}
					return
				},
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMemberV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	createOpts := pools.CreateMemberOpts{
		Name:         d.Get("name").(string),
		Address:      d.Get("address").(string),
		ProtocolPort: d.Get("protocol_port").(int),
		Weight:       d.Get("weight").(int),
	}

	// Must omit if not set
	if v, ok := d.GetOk("subnet_id"); ok {
		createOpts.SubnetID = v.(string)
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	poolID := d.Get("pool_id").(string)
	member, err := pools.CreateMember(elbClient, poolID, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating member: %s", err)
	}

	d.SetId(member.ID)

	return resourceMemberV3Read(d, meta)
}

func resourceMemberV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	member, err := pools.GetMember(elbClient, d.Get("pool_id").(string), d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "member")
	}

	logp.Printf("[DEBUG] Retrieved member %s: %#v", d.Id(), member)

	d.Set("name", member.Name)
	d.Set("weight", member.Weight)
	d.Set("subnet_id", member.SubnetID)
	d.Set("address", member.Address)
	d.Set("protocol_port", member.ProtocolPort)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceMemberV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts pools.UpdateMemberOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("weight") {
		updateOpts.Weight = d.Get("weight").(int)
	}

	logp.Printf("[DEBUG] Updating member %s with options: %#v", d.Id(), updateOpts)
	poolID := d.Get("pool_id").(string)
	_, err = pools.UpdateMember(elbClient, poolID, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to update member %s: %s", d.Id(), err)
	}

	return resourceMemberV3Read(d, meta)
}

func resourceMemberV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	poolID := d.Get("pool_id").(string)
	err = pools.DeleteMember(elbClient, poolID, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Unable to delete member %s: %s", d.Id(), err)
	}
	return nil
}
