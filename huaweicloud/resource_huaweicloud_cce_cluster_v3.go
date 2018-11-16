package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func resourceCCEClusterV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCEClusterV3Create,
		Read:   resourceCCEClusterV3Read,
		Update: resourceCCEClusterV3Update,
		Delete: resourceCCEClusterV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		//request and response parameters
		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"billing_mode": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"highway_subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"extend_param": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"container_network_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_network_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceClusterLabelsV3(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceClusterAnnotationsV3(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceClusterExtendParamV3(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("extend_param").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceCCEClusterV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Unable to create HuaweiCloud CCE client : %s", err)
	}

	createOpts := clusters.CreateOpts{
		Kind:       "Cluster",
		ApiVersion: "v3",
		Metadata: clusters.CreateMetaData{Name: d.Get("name").(string),
			Labels:      resourceClusterLabelsV3(d),
			Annotations: resourceClusterAnnotationsV3(d)},
		Spec: clusters.Spec{
			Type:        d.Get("cluster_type").(string),
			Flavor:      d.Get("flavor").(string),
			Version:     d.Get("cluster_version").(string),
			Description: d.Get("description").(string),
			HostNetwork: clusters.HostNetworkSpec{VpcId: d.Get("vpc_id").(string),
				SubnetId:      d.Get("subnet_id").(string),
				HighwaySubnet: d.Get("highway_subnet_id").(string)},
			ContainerNetwork: clusters.ContainerNetworkSpec{Mode: d.Get("container_network_type").(string),
				Cidr: d.Get("container_network_cidr").(string)},
			BillingMode: d.Get("billing_mode").(int),
			ExtendParam: resourceClusterExtendParamV3(d),
		},
	}

	create, err := clusters.Create(cceClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for HuaweiCloud CCE cluster (%s) to become available", create.Metadata.Id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating"},
		Target:     []string{"Available"},
		Refresh:    waitForCCEClusterActive(cceClient, create.Metadata.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	d.SetId(create.Metadata.Id)

	return resourceCCEClusterV3Read(d, meta)

}

func resourceCCEClusterV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	n, err := clusters.Get(cceClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HuaweiCloud CCE: %s", err)
	}

	d.Set("id", n.Metadata.Id)
	d.Set("name", n.Metadata.Name)
	d.Set("status", n.Status.Phase)
	d.Set("flavor", n.Spec.Flavor)
	d.Set("cluster_version", n.Spec.Version)
	d.Set("cluster_type", n.Spec.Type)
	d.Set("description", n.Spec.Description)
	d.Set("billing_mode", n.Spec.BillingMode)
	d.Set("vpc_id", n.Spec.HostNetwork.VpcId)
	d.Set("subnet_id", n.Spec.HostNetwork.SubnetId)
	d.Set("highway_subnet_id", n.Spec.HostNetwork.HighwaySubnet)
	d.Set("container_network_type", n.Spec.ContainerNetwork.Mode)
	d.Set("container_network_cidr", n.Spec.ContainerNetwork.Cidr)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCCEClusterV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Client: %s", err)
	}

	var updateOpts clusters.UpdateOpts

	if d.HasChange("description") {
		updateOpts.Spec.Description = d.Get("description").(string)
	}
	_, err = clusters.Update(cceClient, d.Id(), updateOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud CCE: %s", err)
	}

	return resourceCCEClusterV3Read(d, meta)
}

func resourceCCEClusterV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Client: %s", err)
	}
	err = clusters.Delete(cceClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE Cluster: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deleting", "Available", "Unavailable"},
		Target:     []string{"Deleted"},
		Refresh:    waitForCCEClusterDelete(cceClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE cluster: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCCEClusterActive(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := clusters.Get(cceClient, clusterId).Extract()
		if err != nil {
			return nil, "", err
		}
		if n.Status.Phase != "Creating" {
			return n, "Creating", nil
		}

		return n, n.Status.Phase, nil
	}
}

func waitForCCEClusterDelete(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud CCE cluster %s.\n", clusterId)

		r, err := clusters.Get(cceClient, clusterId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud CCE cluster %s", clusterId)
				return r, "Deleted", nil
			}
		}
		if r.Status.Phase == "Deleting" {
			return r, "Deleting", nil
		}
		log.Printf("[DEBUG] HuaweiCloud CCE cluster %s still available.\n", clusterId)
		return r, "Available", nil
	}
}
