package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/endpoints"
)

func ResourceVPCEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCEndpointCreate,
		Read:   resourceVPCEndpointRead,
		Update: resourceVPCEndpointUpdate,
		Delete: resourceVPCEndpointDelete,

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
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_dns": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"enable_whitelist": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"packet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceVPCEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcepClient, err := config.VPCEPClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	enableDNS := d.Get("enable_dns").(bool)
	enableACL := d.Get("enable_whitelist").(bool)
	createOpts := endpoints.CreateOpts{
		ServiceID:       d.Get("service_id").(string),
		VpcID:           d.Get("vpc_id").(string),
		SubnetID:        d.Get("network_id").(string),
		PortIP:          d.Get("ip_address").(string),
		EnableDNS:       &enableDNS,
		EnableWhitelist: &enableACL,
	}

	raw := d.Get("whitelist").(*schema.Set).List()
	if enableACL && len(raw) > 0 {
		whitelists := make([]string, len(raw))
		for i, v := range raw {
			whitelists[i] = v.(string)
		}
		createOpts.Whitelist = whitelists
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandResourceTags(tagRaw)
		createOpts.Tags = taglist
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	ep, err := endpoints.Create(vpcepClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud VPC endpoint: %s", err)
	}

	d.SetId(ep.ID)
	log.Printf("[INFO] Waiting for Huaweicloud VPC endpoint(%s) to become accepted", ep.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"accepted", "pendingAcceptance"},
		Refresh:    waitForVPCEndpointStatus(vpcepClient, ep.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for VPC endpoint(%s) to become accepted: %s",
			ep.ID, stateErr)
	}

	return resourceVPCEndpointRead(d, meta)
}

func resourceVPCEndpointRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcepClient, err := config.VPCEPClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	ep, err := endpoints.Get(vpcepClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Huaweicloud VPC endpoint: %s", err)
	}

	log.Printf("[DEBUG] retrieving Huaweicloud VPC endpoint: %#v", ep)
	d.Set("region", GetRegion(d, config))
	d.Set("status", ep.Status)
	d.Set("service_id", ep.ServiceID)
	d.Set("service_name", ep.ServiceName)
	d.Set("service_type", ep.ServiceType)
	d.Set("vpc_id", ep.VpcID)
	d.Set("network_id", ep.SubnetID)
	d.Set("ip_address", ep.IPAddr)
	d.Set("enable_dns", ep.EnableDNS)
	d.Set("enable_whitelist", ep.EnableWhitelist)
	d.Set("whitelist", ep.Whitelist)
	d.Set("packet_id", ep.MarkerID)

	if len(ep.DNSNames) > 0 {
		d.Set("private_domain_name", ep.DNSNames[0])
	} else {
		d.Set("private_domain_name", nil)
	}

	// fetch tags from endpoints.Endpoint
	tagmap := make(map[string]string)
	for _, val := range ep.Tags {
		tagmap[val.Key] = val.Value
	}
	d.Set("tags", tagmap)

	return nil
}

func resourceVPCEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcepClient, err := config.VPCEPClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	//update tags
	if d.HasChange("tags") {
		tagErr := UpdateResourceTags(vpcepClient, d, tagVPCEP, d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of VPC endpoint service %s: %s", d.Id(), tagErr)
		}
	}
	return resourceVPCEndpointRead(d, meta)
}

func resourceVPCEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcepClient, err := config.VPCEPClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	err = endpoints.Delete(vpcepClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting Huaweicloud VPC endpoint %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForVPCEndpointStatus(vpcepClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting Huaweicloud VPC endpoint %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func waitForVPCEndpointStatus(vpcepClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ep, err := endpoints.Get(vpcepClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud VPC endpoint %s", id)
				return ep, "deleted", nil
			}
			return ep, "error", err
		}

		return ep, ep.Status, nil
	}
}
