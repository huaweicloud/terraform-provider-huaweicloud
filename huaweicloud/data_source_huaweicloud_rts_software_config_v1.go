package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwareconfig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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
	config := meta.(*config.Config)
	orchestrationClient, err := config.OrchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud rts client: %s", err)
	}

	n, err := softwareconfig.Get(orchestrationClient, d.Get("id").(string)).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud RTS Config: %s", err)
	}

	logp.Printf("[INFO] Retrieved RTS Software Config using given id %s", n.Id)
	d.SetId(n.Id)

	d.Set("name", n.Name)
	d.Set("group", n.Group)
	d.Set("region", GetRegion(d, config))
	d.Set("config", n.Config)
	d.Set("options", n.Options)

	if err := d.Set("input_values", n.Inputs); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving inputs to state for HuaweiCloud RTS Software Config (%s): %s", d.Id(), err)
	}
	if err := d.Set("output_values", n.Outputs); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving outputs to state for HuaweiCloud RTS Software Config (%s): %s", d.Id(), err)
	}

	return nil
}
