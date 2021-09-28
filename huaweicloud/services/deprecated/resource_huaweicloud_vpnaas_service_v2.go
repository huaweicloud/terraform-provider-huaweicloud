package deprecated

import (
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceVpnServiceV2() *schema.Resource {
	return &schema.Resource{
		Create:             resourceVpnServiceV2Create,
		Read:               resourceVpnServiceV2Read,
		Update:             resourceVpnServiceV2Update,
		Delete:             resourceVpnServiceV2Delete,
		DeprecationMessage: "VPN has been deprecated.",
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_v6_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_v4_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVpnServiceV2Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var createOpts services.CreateOptsBuilder

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts = VpnServiceCreateOpts{
		services.CreateOpts{
			Name:         d.Get("name").(string),
			Description:  d.Get("description").(string),
			AdminStateUp: &adminStateUp,
			TenantID:     d.Get("tenant_id").(string),
			SubnetID:     d.Get("subnet_id").(string),
			RouterID:     d.Get("router_id").(string),
		},
		MapValueSpecs(d),
	}

	logp.Printf("[DEBUG] Create service: %#v", createOpts)

	service, err := services.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NOT_CREATED"},
		Target:     []string{"PENDING_CREATE"},
		Refresh:    waitForServiceCreation(networkingClient, service.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()

	if err != nil {
		return err
	}

	logp.Printf("[DEBUG] Service created: %#v", service)

	d.SetId(service.ID)

	return resourceVpnServiceV2Read(d, meta)
}

func resourceVpnServiceV2Read(d *schema.ResourceData, meta interface{}) error {
	logp.Printf("[DEBUG] Retrieve information about service: %s", d.Id())

	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	service, err := services.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "service")
	}

	logp.Printf("[DEBUG] Read HuaweiCloud Service %s: %#v", d.Id(), service)

	d.Set("name", service.Name)
	d.Set("description", service.Description)
	d.Set("subnet_id", service.SubnetID)
	d.Set("admin_state_up", service.AdminStateUp)
	d.Set("tenant_id", service.TenantID)
	d.Set("router_id", service.RouterID)
	d.Set("status", service.Status)
	d.Set("external_v6_ip", service.ExternalV6IP)
	d.Set("external_v4_ip", service.ExternalV4IP)
	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceVpnServiceV2Update(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	opts := services.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &adminStateUp
		hasChange = true
	}

	var updateOpts services.UpdateOptsBuilder
	updateOpts = opts

	logp.Printf("[DEBUG] Updating service with id %s: %#v", d.Id(), updateOpts)

	if hasChange {
		service, err := services.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return err
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForServiceUpdate(networkingClient, service.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForState()

		if err != nil {
			return err
		}

		logp.Printf("[DEBUG] Updated service with id %s", d.Id())
	}

	return resourceVpnServiceV2Read(d, meta)
}

func resourceVpnServiceV2Delete(d *schema.ResourceData, meta interface{}) error {
	logp.Printf("[DEBUG] Destroy service: %s", d.Id())

	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	err = services.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForServiceDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()

	return err
}

func waitForServiceDeletion(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		serv, err := services.Get(networkingClient, id).Extract()
		logp.Printf("[DEBUG] Got service %s => %#v", id, serv)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Service %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmtp.Errorf("Unexpected error: %s", err)
		}

		logp.Printf("[DEBUG] Service %s deletion is pending", id)
		return serv, "DELETING", nil
	}
}

func waitForServiceCreation(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := services.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "NOT_CREATED", nil
		}
		return service, "PENDING_CREATE", nil
	}
}

func waitForServiceUpdate(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := services.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return service, "UPDATED", nil
	}
}
