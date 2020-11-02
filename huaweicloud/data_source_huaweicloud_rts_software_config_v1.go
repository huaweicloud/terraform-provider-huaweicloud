package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwareconfig"
)

func dataSourceRtsSoftwareConfigV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRtsSoftwareConfigV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"input_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			"output_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			"config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRtsSoftwareConfigV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud rts client: %s", err)
	}

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
