package vpcep

import (
	"context"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceVPCEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCEndpointCreate,
		ReadContext:   resourceVPCEndpointRead,
		UpdateContext: resourceVPCEndpointUpdate,
		DeleteContext: resourceVPCEndpointDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"tags": common.TagsSchema(),
		},
	}
}

func resourceVPCEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
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
		taglist := utils.ExpandResourceTags(tagRaw)
		createOpts.Tags = taglist
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	ep, err := endpoints.Create(vpcepClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC endpoint: %s", err)
	}

	d.SetId(ep.ID)
	log.Printf("[INFO] Waiting for VPC endpoint(%s) to become accepted", ep.ID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"creating"},
		Target:       []string{"accepted", "pendingAcceptance"},
		Refresh:      waitForVPCEndpointStatus(vpcepClient, ep.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf("error waiting for VPC endpoint(%s) to become accepted: %s", ep.ID, stateErr)
	}

	return resourceVPCEndpointRead(ctx, d, meta)
}

func resourceVPCEndpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	ep, err := endpoints.Get(vpcepClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC endpoint")
	}

	log.Printf("[DEBUG] retrieving VPC endpoint: %#v", ep)
	d.Set("region", region)
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

func resourceVPCEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	//update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(vpcepClient, d, tagVPCEP, d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPC endpoint service %s: %s", d.Id(), tagErr)
		}
	}
	return resourceVPCEndpointRead(ctx, d, meta)
}

func resourceVPCEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	err = endpoints.Delete(vpcepClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting VPC endpoint %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"deleting"},
		Target:       []string{"deleted"},
		Refresh:      waitForVPCEndpointStatus(vpcepClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting VPC endpoint %s: %s", d.Id(), err)
	}

	return nil
}

func waitForVPCEndpointStatus(vpcepClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ep, err := endpoints.Get(vpcepClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted VPC endpoint %s", id)
				return ep, "deleted", nil
			}
			return ep, "error", err
		}

		return ep, ep.Status, nil
	}
}
