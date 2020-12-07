package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cci/v1/networks"
)

func resourceCCINetworkV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCINetworkV1Create,
		Read:   resourceCCINetworkV1Read,
		Delete: resourceCCINetworkV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		//request and response parameters
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
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_id": {
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
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"available_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkAnnotationsV1(d *schema.ResourceData) map[string]string {
	m := map[string]string{
		"network.alpha.kubernetes.io/default-security-group": d.Get("security_group").(string),
		"network.alpha.kubernetes.io/domain_id":              d.Get("domain_id").(string),
		"network.alpha.kubernetes.io/project_id":             d.Get("project_id").(string),
	}
	return m
}

func resourceCCINetworkV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cciClient, err := config.CciV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Unable to create HuaweiCloud CCI client : %s", err)
	}

	createOpts := networks.CreateOpts{
		Kind:       "Network",
		ApiVersion: "networking.cci.io/v1beta1",
		Metadata: networks.CreateMetaData{
			Name:        d.Get("name").(string),
			Annotations: resourceNetworkAnnotationsV1(d),
		},
		Spec: networks.Spec{
			AttachedVPC:   d.Get("vpc_id").(string),
			NetworkType:   "underlay_neutron",
			NetworkID:     d.Get("network_id").(string),
			SubnetID:      d.Get("subnet_id").(string),
			AvailableZone: d.Get("available_zone").(string),
			Cidr:          d.Get("cidr").(string),
		},
	}

	ns := d.Get("namespace").(string)
	create, err := networks.Create(cciClient, ns, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI Network: %s", err)
	}

	log.Printf("[DEBUG] Waiting for HuaweiCloud CCI network (%s) to become available", create.Metadata.Name)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Initializing", "Pending"},
		Target:     []string{"Active"},
		Refresh:    waitForCCINetworkActive(cciClient, ns, create.Metadata.Name),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI network: %s", err)
	}
	d.SetId(create.Metadata.Name)

	return resourceCCINetworkV1Read(d, meta)
}

func resourceCCINetworkV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cciClient, err := config.CciV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	n, err := networks.Get(cciClient, ns, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HuaweiCloud CCI: %s", err)
	}

	d.Set("name", n.Metadata.Name)
	d.Set("vpc_id", n.Spec.AttachedVPC)
	d.Set("network_id", n.Spec.NetworkID)
	d.Set("subnet_id", n.Spec.SubnetID)
	d.Set("available_zone", n.Spec.AvailableZone)
	d.Set("cidr", n.Spec.Cidr)

	return nil
}

func resourceCCINetworkV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cciClient, err := config.CciV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI Client: %s", err)
	}

	ns := d.Get("namespace").(string)
	err = networks.Delete(cciClient, ns, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCI Network: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Terminating", "Active"},
		Target:     []string{"Deleted"},
		Refresh:    waitForCCINetworkDelete(cciClient, ns, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCI network: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCCINetworkActive(cciClient *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := networks.Get(cciClient, ns, name).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.State, nil
	}
}

func waitForCCINetworkDelete(cciClient *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud CCI network %s.\n", name)

		r, err := networks.Get(cciClient, ns, name).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud CCI network %s", name)
				return r, "Deleted", nil
			}
		}
		if r.Status.State == "Terminating" {
			return r, "Terminating", nil
		}
		log.Printf("[DEBUG] HuaweiCloud CCI network %s still available.\n", name)
		return r, "Active", nil
	}
}
