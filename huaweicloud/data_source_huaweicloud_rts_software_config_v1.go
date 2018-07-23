package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwareconfig"
)

func dataSourceRtsSoftwareConfigV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRtsSoftwareConfigV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"input_values": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"output_values": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"config": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"options": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRtsSoftwareConfigV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))

	n, err := softwareconfig.Get(orchestrationClient, d.Get("id").(string)).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HuaweiCloud RTS Config: %s", err)
	}

	log.Printf("[INFO] Retrieved RTS Software Config using given id %s", n.Id)
	d.SetId(n.Id)

	d.Set("name", n.Name)
	d.Set("group", n.Group)
	d.Set("region", GetRegion(d, config))
	d.Set("config", n.Config)
	d.Set("options", n.Options)

	if err := d.Set("input_values", n.Inputs); err != nil {
		return fmt.Errorf("[DEBUG] Error saving inputs to state for HuaweiCloud RTS Software Config (%s): %s", d.Id(), err)
	}
	if err := d.Set("output_values", n.Outputs); err != nil {
		return fmt.Errorf("[DEBUG] Error saving outputs to state for HuaweiCloud RTS Software Config (%s): %s", d.Id(), err)
	}

	return nil
}
