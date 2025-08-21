package deprecated

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/provider"
	"github.com/chnsz/golangsdk/openstack/networking/v2/networks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceNetworkingNetworkV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingNetworkV2Create,
		Read:   resourceNetworkingNetworkV2Read,
		Update: resourceNetworkingNetworkV2Update,
		Delete: resourceNetworkingNetworkV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use VPC resources instead",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"shared": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"segments": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"physical_network": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"segmentation_id": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
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

func resourceNetworkingNetworkV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	createOpts := NetworkCreateOpts{
		networks.CreateOpts{
			Name:     d.Get("name").(string),
			TenantID: d.Get("tenant_id").(string),
		},
		MapValueSpecs(d),
	}

	asuRaw := d.Get("admin_state_up").(string)
	if asuRaw != "" {
		asu, err := strconv.ParseBool(asuRaw)
		if err != nil {
			return fmtp.Errorf("admin_state_up, if provided, must be either 'true' or 'false'")
		}
		createOpts.AdminStateUp = &asu
	}

	sharedRaw := d.Get("shared").(string)
	if sharedRaw != "" {
		shared, err := strconv.ParseBool(sharedRaw)
		if err != nil {
			return fmtp.Errorf("shared, if provided, must be either 'true' or 'false': %v", err)
		}
		createOpts.Shared = &shared
	}

	segments := resourceNetworkingNetworkV2Segments(d)

	var n *networks.Network
	if len(segments) > 0 {
		providerCreateOpts := provider.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			Segments:          segments,
		}
		logp.Printf("[DEBUG] Create Options: %#v", providerCreateOpts)
		n, err = networks.Create(networkingClient, providerCreateOpts).Extract()
	} else {
		logp.Printf("[DEBUG] Create Options: %#v", createOpts)
		n, err = networks.Create(networkingClient, createOpts).Extract()
	}

	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Neutron network: %s", err)
	}

	logp.Printf("[INFO] Network ID: %s", n.ID)

	logp.Printf("[DEBUG] Waiting for Network (%s) to become available", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNetworkActive(networkingClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(n.ID)

	return resourceNetworkingNetworkV2Read(d, meta)
}

func resourceNetworkingNetworkV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	n, err := networks.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "network")
	}

	logp.Printf("[DEBUG] Retrieved Network %s: %+v", d.Id(), n)

	d.Set("name", n.Name)
	d.Set("admin_state_up", strconv.FormatBool(n.AdminStateUp))
	d.Set("shared", strconv.FormatBool(n.Shared))
	d.Set("tenant_id", n.TenantID)
	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceNetworkingNetworkV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var updateOpts networks.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("admin_state_up") {
		asuRaw := d.Get("admin_state_up").(string)
		if asuRaw != "" {
			asu, err := strconv.ParseBool(asuRaw)
			if err != nil {
				return fmtp.Errorf("admin_state_up, if provided, must be either 'true' or 'false'")
			}
			updateOpts.AdminStateUp = &asu
		}
	}
	if d.HasChange("shared") {
		sharedRaw := d.Get("shared").(string)
		if sharedRaw != "" {
			shared, err := strconv.ParseBool(sharedRaw)
			if err != nil {
				return fmtp.Errorf("shared, if provided, must be either 'true' or 'false': %v", err)
			}
			updateOpts.Shared = &shared
		}
	}

	logp.Printf("[DEBUG] Updating Network %s with options: %+v", d.Id(), updateOpts)

	_, err = networks.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud Neutron Network: %s", err)
	}

	return resourceNetworkingNetworkV2Read(d, meta)
}

func resourceNetworkingNetworkV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNetworkDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Neutron Network: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceNetworkingNetworkV2Segments(d *schema.ResourceData) (providerSegments []provider.Segment) {
	segments := d.Get("segments").([]interface{})
	for _, v := range segments {
		var segment provider.Segment
		segmentMap := v.(map[string]interface{})

		if v, ok := segmentMap["physical_network"].(string); ok {
			segment.PhysicalNetwork = v
		}

		if v, ok := segmentMap["network_type"].(string); ok {
			segment.NetworkType = v
		}

		if v, ok := segmentMap["segmentation_id"].(int); ok {
			segment.SegmentationID = v
		}

		providerSegments = append(providerSegments, segment)
	}
	return
}

func waitForNetworkActive(networkingClient *golangsdk.ServiceClient, networkId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := networks.Get(networkingClient, networkId).Extract()
		if err != nil {
			return nil, "", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Neutron Network: %+v", n)
		if n.Status == "DOWN" || n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, n.Status, nil
	}
}

func waitForNetworkDelete(networkingClient *golangsdk.ServiceClient, networkId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Network %s.\n", networkId)

		n, err := networks.Get(networkingClient, networkId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Network %s", networkId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = networks.Delete(networkingClient, networkId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Network %s", networkId)
				return n, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault409); ok {
				return n, "ACTIVE", nil
			}
			return n, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Network %s still active.\n", networkId)
		return n, "ACTIVE", nil
	}
}
