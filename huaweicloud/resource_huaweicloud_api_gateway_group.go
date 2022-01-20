package huaweicloud

import (
	"time"

	"github.com/chnsz/golangsdk/openstack/apigw/shared/v1/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceAPIGatewayGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAPIGatewayGroupCreate,
		Read:   resourceAPIGatewayGroupRead,
		Update: resourceAPIGatewayGroupUpdate,
		Delete: resourceAPIGatewayGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAPIGatewayGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	apigwClient, err := config.ApiGatewayV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud api gateway client: %s", err)
	}

	createOpts := &groups.CreateOpts{
		Name:   d.Get("name").(string),
		Remark: d.Get("description").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := groups.Create(apigwClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud api group: %s", err)
	}

	// Store the ID now
	d.SetId(v.ID)

	return resourceAPIGatewayGroupRead(d, meta)
}

func resourceAPIGatewayGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	apigwClient, err := config.ApiGatewayV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud api gateway client: %s", err)
	}

	v, err := groups.Get(apigwClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "API GateWay group")
	}

	logp.Printf("[DEBUG] Retrieved api group %s: %+v", d.Id(), v)

	d.Set("name", v.Name)
	d.Set("description", v.Remark)
	d.Set("status", v.Status)

	return nil
}

func resourceAPIGatewayGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	apigwClient, err := config.ApiGatewayV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud api gateway client: %s", err)
	}

	updateOpts := groups.UpdateOpts{
		Name:   d.Get("name").(string),
		Remark: d.Get("description").(string),
	}

	_, err = groups.Update(apigwClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud api group: %s", err)
	}

	return resourceAPIGatewayGroupRead(d, meta)
}

func resourceAPIGatewayGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	apigwClient, err := config.ApiGatewayV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud api gateway client: %s", err)
	}

	if err := groups.Delete(apigwClient, d.Id()).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "api groups")
	}

	d.SetId("")
	return nil
}
