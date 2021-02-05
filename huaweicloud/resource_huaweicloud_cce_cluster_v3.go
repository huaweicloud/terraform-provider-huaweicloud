package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func ResourceCCEClusterV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCEClusterV3Create,
		Read:   resourceCCEClusterV3Read,
		Update: resourceCCEClusterV3Update,
		Delete: resourceCCEClusterV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "VirtualMachine",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"highway_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"extend_param": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"container_network_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_network_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"authentication_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "rbac",
			},
			"authenticating_proxy_ca": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"multi_az": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"masters"},
			},
			"masters": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				MaxItems:      3,
				ConflictsWith: []string{"multi_az"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"eip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIP,
			},
			"service_network_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"kube_proxy_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_config_raw": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_authority_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"certificate_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_certificate_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
func resourceClusterExtendParamV3(d *schema.ResourceData, config *Config) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("extend_param").(map[string]interface{}) {
		m[key] = val.(string)
	}
	if multi_az, ok := d.GetOk("multi_az"); ok && multi_az == true {
		m["clusterAZ"] = "multi_az"
	}
	if kube_proxy_mode, ok := d.GetOk("kube_proxy_mode"); ok {
		m["kubeProxyMode"] = kube_proxy_mode.(string)
	}
	if eip, ok := d.GetOk("eip"); ok {
		m["clusterExternalIP"] = eip.(string)
	}

	epsID := GetEnterpriseProjectID(d, config)

	if epsID != "" {
		m["enterpriseProjectId"] = epsID
	}
	return m
}

func resourceClusterMastersV3(d *schema.ResourceData) ([]clusters.MasterSpec, error) {
	if v, ok := d.GetOk("masters"); ok {
		flavorId := d.Get("flavor_id").(string)
		mastersRaw := v.([]interface{})
		if strings.Contains(flavorId, "s1") && len(mastersRaw) != 1 {
			return nil, fmt.Errorf("Error creating HuaweiCloud Cluster: "+
				"single-master cluster need 1 az for master node, but got %d", len(mastersRaw))
		}
		if strings.Contains(flavorId, "s2") && len(mastersRaw) != 3 {
			return nil, fmt.Errorf("Error creating HuaweiCloud Cluster: "+
				"high-availability cluster need 3 az for master nodes, but got %d", len(mastersRaw))
		}
		masters := make([]clusters.MasterSpec, len(mastersRaw))
		for i, raw := range mastersRaw {
			rawMap := raw.(map[string]interface{})
			masters[i] = clusters.MasterSpec{
				MasterAZ: rawMap["availability_zone"].(string),
			}
		}
		return masters, nil
	}

	return nil, nil
}

func resourceCCEClusterV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.CceV3Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Unable to create HuaweiCloud CCE client : %s", err)
	}

	authenticating_proxy := make(map[string]string)
	if hasFilledOpt(d, "authenticating_proxy_ca") {
		authenticating_proxy["ca"] = d.Get("authenticating_proxy_ca").(string)
	}
	createOpts := clusters.CreateOpts{
		Kind:       "Cluster",
		ApiVersion: "v3",
		Metadata: clusters.CreateMetaData{Name: d.Get("name").(string),
			Labels:      resourceClusterLabelsV3(d),
			Annotations: resourceClusterAnnotationsV3(d)},
		Spec: clusters.Spec{
			Type:        d.Get("cluster_type").(string),
			Flavor:      d.Get("flavor_id").(string),
			Version:     d.Get("cluster_version").(string),
			Description: d.Get("description").(string),
			HostNetwork: clusters.HostNetworkSpec{
				VpcId:         d.Get("vpc_id").(string),
				SubnetId:      d.Get("subnet_id").(string),
				HighwaySubnet: d.Get("highway_subnet_id").(string),
			},
			ContainerNetwork: clusters.ContainerNetworkSpec{
				Mode: d.Get("container_network_type").(string),
				Cidr: d.Get("container_network_cidr").(string),
			},
			Authentication: clusters.AuthenticationSpec{
				Mode:                d.Get("authentication_mode").(string),
				AuthenticatingProxy: authenticating_proxy,
			},
			BillingMode:          d.Get("billing_mode").(int),
			ExtendParam:          resourceClusterExtendParamV3(d, config),
			KubernetesSvcIPRange: d.Get("service_network_cidr").(string),
		},
	}

	masters, err := resourceClusterMastersV3(d)
	if err != nil {
		return err
	}
	createOpts.Spec.Masters = masters

	create, err := clusters.Create(cceClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for HuaweiCloud CCE cluster (%s) to become available", create.Metadata.Id)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Creating"},
		Target:       []string{"Available"},
		Refresh:      waitForCCEClusterActive(cceClient, create.Metadata.Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        150 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE cluster: %s", err)
	}
	d.SetId(create.Metadata.Id)

	return resourceCCEClusterV3Read(d, meta)

}

func resourceCCEClusterV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.CceV3Client(GetRegion(d, config))
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

	d.Set("name", n.Metadata.Name)
	d.Set("status", n.Status.Phase)
	d.Set("flavor_id", n.Spec.Flavor)
	d.Set("cluster_version", n.Spec.Version)
	d.Set("cluster_type", n.Spec.Type)
	d.Set("description", n.Spec.Description)
	d.Set("billing_mode", n.Spec.BillingMode)
	d.Set("vpc_id", n.Spec.HostNetwork.VpcId)
	d.Set("subnet_id", n.Spec.HostNetwork.SubnetId)
	d.Set("highway_subnet_id", n.Spec.HostNetwork.HighwaySubnet)
	d.Set("container_network_type", n.Spec.ContainerNetwork.Mode)
	d.Set("container_network_cidr", n.Spec.ContainerNetwork.Cidr)
	d.Set("authentication_mode", n.Spec.Authentication.Mode)
	d.Set("security_group_id", n.Spec.HostNetwork.SecurityGroup)
	d.Set("region", GetRegion(d, config))
	d.Set("enterprise_project_id", n.Spec.ExtendParam["enterpriseProjectId"])
	d.Set("service_network_cidr", n.Spec.KubernetesSvcIPRange)

	r := clusters.GetCert(cceClient, d.Id())

	kubeConfigRaw, err := jsonMarshal(r.Body)

	if err != nil {
		log.Printf("Error marshaling r.Body: %s", err)
	}

	d.Set("kube_config_raw", string(kubeConfigRaw))

	cert, err := r.Extract()

	if err != nil {
		log.Printf("Error retrieving HuaweiCloud CCE cluster cert: %s", err)
	}

	//Set Certificate Clusters
	var clusterList []map[string]interface{}
	for _, clusterObj := range cert.Clusters {
		clusterCert := make(map[string]interface{})
		clusterCert["name"] = clusterObj.Name
		clusterCert["server"] = clusterObj.Cluster.Server
		clusterCert["certificate_authority_data"] = clusterObj.Cluster.CertAuthorityData
		clusterList = append(clusterList, clusterCert)
	}
	d.Set("certificate_clusters", clusterList)

	//Set Certificate Users
	var userList []map[string]interface{}
	for _, userObj := range cert.Users {
		userCert := make(map[string]interface{})
		userCert["name"] = userObj.Name
		userCert["client_certificate_data"] = userObj.User.ClientCertData
		userCert["client_key_data"] = userObj.User.ClientKeyData
		userList = append(userList, userCert)
	}
	d.Set("certificate_users", userList)

	// Set masters
	var masterList []map[string]interface{}
	for _, masterObj := range n.Spec.Masters {
		master := make(map[string]interface{})
		master["availability_zone"] = masterObj.MasterAZ
		masterList = append(masterList, master)
	}
	d.Set("masters", masterList)

	return nil
}

func resourceCCEClusterV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.CceV3Client(GetRegion(d, config))
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
	cceClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Client: %s", err)
	}
	err = clusters.Delete(cceClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE Cluster: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting", "Available", "Unavailable"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCCEClusterDelete(cceClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
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
