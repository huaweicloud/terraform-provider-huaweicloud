package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/vpcs"
)

func ResourceIecVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceIecVpcV1Create,
		Read:   resourceIecVpcV1Read,
		Update: resourceIecVpcV1Update,
		Delete: resourceIecVpcV1Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SYSTEM",
			},
			"subnet_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceIecVpcV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecV1Client, err := config.IECV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	createOpts := vpcs.CreateOpts{
		Name: d.Get("name").(string),
		Cidr: d.Get("cidr").(string),
		Mode: d.Get("mode").(string),
	}

	n, err := vpcs.Create(iecV1Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC VPC: %s", err)
	}

	log.Printf("[INFO] IEC VPC ID: %s", n.ID)
	d.SetId(n.ID)

	return resourceIecVpcV1Read(d, meta)
}

func resourceIecVpcV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecV1Client, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	n, err := vpcs.Get(iecV1Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving Huaweicloud IEC VPC")
	}

	d.Set("name", n.Name)
	d.Set("cidr", n.Cidr)
	d.Set("mode", n.Mode)
	d.Set("subnet_num", n.SubnetNum)

	return nil
}

func resourceIecVpcV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecV1Client, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	var updateOpts vpcs.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("cidr") {
		updateOpts.Cidr = d.Get("cidr").(string)
	}

	_, err = vpcs.Update(iecV1Client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Huaweicloud IEC VPC: %s", err)
	}

	return resourceIecVpcV1Read(d, meta)
}

func resourceIecVpcV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecV1Client, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	err = vpcs.Delete(iecV1Client, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving Huaweicloud IEC VPC")
	}

	d.SetId("")
	return nil
}
