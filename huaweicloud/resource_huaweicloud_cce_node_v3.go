package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
)

func resourceCCENodeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodeV3Create,
		Read:   resourceCCENodeV3Read,
		Update: resourceCCENodeV3Update,
		Delete: resourceCCENodeV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"az": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"root_volume": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumetype": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"data_volumes": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumetype": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"eip_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"eip_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"iptype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"bandwidth_charge_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"sharetype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"bandwidth_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"billing_mode": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"extend_param_charging_mode": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ecs_performance_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"order_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"product_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"max_pods": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceCCENodeLabelsV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceCCENodeAnnotationsV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceCCEDataVolume(d *schema.ResourceData) []nodes.VolumeSpec {
	volumeRaw := d.Get("data_volumes").(*schema.Set).List()
	volumes := make([]nodes.VolumeSpec, len(volumeRaw))
	for i, raw := range volumeRaw {
		rawMap := raw.(map[string]interface{})
		volumes[i] = nodes.VolumeSpec{
			Size:        rawMap["size"].(int),
			VolumeType:  rawMap["volumetype"].(string),
			ExtendParam: rawMap["extend_param"].(string),
		}
	}
	return volumes
}
func resourceCCERootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var nics nodes.VolumeSpec
	nicsRaw := d.Get("root_volume").([]interface{})
	if len(nicsRaw) == 1 {
		nics.Size = nicsRaw[0].(map[string]interface{})["size"].(int)
		nics.VolumeType = nicsRaw[0].(map[string]interface{})["volumetype"].(string)
		nics.ExtendParam = nicsRaw[0].(map[string]interface{})["extend_param"].(string)
	}
	return nics
}
func resourceCCEEipIDs(d *schema.ResourceData) []string {
	rawID := d.Get("eip_ids").(*schema.Set)
	id := make([]string, rawID.Len())
	for i, raw := range rawID.List() {
		id[i] = raw.(string)
	}
	return id
}
func resourceCCENodeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Node client: %s", err)
	}

	createOpts := nodes.CreateOpts{
		Kind:       "Node",
		ApiVersion: "v3",
		Metadata: nodes.CreateMetaData{
			Name:        d.Get("name").(string),
			Labels:      resourceCCENodeLabelsV2(d),
			Annotations: resourceCCENodeAnnotationsV2(d),
		},
		Spec: nodes.Spec{
			Flavor:      d.Get("flavor").(string),
			Az:          d.Get("az").(string),
			Login:       nodes.LoginSpec{SshKey: d.Get("sshkey").(string)},
			RootVolume:  resourceCCERootVolume(d),
			DataVolumes: resourceCCEDataVolume(d),
			PublicIP: nodes.PublicIPSpec{
				Ids:   resourceCCEEipIDs(d),
				Count: d.Get("eip_count").(int),
				Eip: nodes.EipSpec{
					IpType: d.Get("iptype").(string),
					Bandwidth: nodes.BandwidthOpts{
						ChargeMode: d.Get("bandwidth_charge_mode").(string),
						Size:       d.Get("bandwidth_size").(int),
						ShareType:  d.Get("sharetype").(string),
					},
				},
			},
			BillingMode: d.Get("billing_mode").(int),
			Count:       1,
			ExtendParam: nodes.ExtendParam{
				ChargingMode:       d.Get("extend_param_charging_mode").(int),
				EcsPerformanceType: d.Get("ecs_performance_type").(string),
				MaxPods:            d.Get("max_pods").(int),
				OrderID:            d.Get("order_id").(string),
				ProductID:          d.Get("product_id").(string),
				PublicKey:          d.Get("public_key").(string),
			},
		},
	}

	clusterid := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:     []string{"Available"},
		Refresh:    waitForClusterAvailable(nodeClient, clusterid),
		Timeout:    15 * time.Minute,
		Delay:      15 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateCluster.WaitForState()

	s, err := nodes.Create(nodeClient, clusterid, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveCreate(nodeClient, createOpts, clusterid, 403)
			if err == "fail" {
				return fmt.Errorf("Error creating HuaweiCloud Node")
			}
			s = retryNode
		} else {
			return fmt.Errorf("Error creating HuaweiCloud Node: %s", err)
		}
	}

	job, err := nodes.GetJobDetails(nodeClient, s.Status.JobID).ExtractJob()
	if err != nil {
		return fmt.Errorf("Error fetching HuaweiCloud Job Details: %s", err)
	}
	jobResorceId := job.Spec.SubJobs[0].Metadata.ID

	subjob, err := nodes.GetJobDetails(nodeClient, jobResorceId).ExtractJob()
	if err != nil {
		return fmt.Errorf("Error fetching HuaweiCloud Job Details: %s", err)
	}
	nodeid := subjob.Spec.SubJobs[0].Spec.ResourceID

	log.Printf("[DEBUG] Waiting for CCE Node (%s) to become available", s.Metadata.Name)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Build", "Installing"},
		Target:     []string{"Active", "Abnormal"},
		Refresh:    waitForCceNodeActive(nodeClient, clusterid, nodeid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()

	node, err := nodes.Get(nodeClient, clusterid, nodeid).Extract()
	d.SetId(node.Metadata.Id)
	d.Set("iptype", s.Spec.PublicIP.Eip.IpType)
	d.Set("bandwidth_charge_mode", s.Spec.PublicIP.Eip.Bandwidth.ChargeMode)
	d.Set("bandwidth_size", s.Spec.PublicIP.Eip.Bandwidth.Size)
	d.Set("sharetype", s.Spec.PublicIP.Eip.Bandwidth.ShareType)
	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Node client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	s, err := nodes.Get(nodeClient, clusterid, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HuaweiCloud Node: %s", err)
	}

	d.Set("name", s.Metadata.Name)
	d.Set("labels", s.Metadata.Labels)
	d.Set("annotations", s.Metadata.Annotations)
	d.Set("flavor", s.Spec.Flavor)
	d.Set("az", s.Spec.Az)
	d.Set("billing_mode", s.Spec.BillingMode)
	d.Set("node_count", s.Spec.Count)
	d.Set("extend_param_charging_mode", s.Spec.ExtendParam.ChargingMode)
	d.Set("ecs:performance_type", s.Spec.ExtendParam.PublicKey)
	d.Set("order_id", s.Spec.ExtendParam.OrderID)
	d.Set("product_id", s.Spec.ExtendParam.ProductID)
	d.Set("max_pods", s.Spec.ExtendParam.MaxPods)
	d.Set("ecs_performance_type", s.Spec.ExtendParam.EcsPerformanceType)
	d.Set("sshkey", s.Spec.Login.SshKey)
	var volumes []map[string]interface{}
	for _, pairObject := range s.Spec.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_param"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving dataVolumes to state for HuaweiCloud Node (%s): %s", d.Id(), err)
	}

	rootVolume := []map[string]interface{}{
		{
			"size":         s.Spec.RootVolume.Size,
			"volumetype":   s.Spec.RootVolume.VolumeType,
			"extend_param": s.Spec.RootVolume.ExtendParam,
		},
	}
	d.Set("root_volume", rootVolume)
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmt.Errorf("[DEBUG] Error saving root Volume to state for HuaweiCloud Node (%s): %s", d.Id(), err)
	}

	d.Set("eip_ids", s.Spec.PublicIP.Ids)
	d.Set("eip_count", s.Spec.PublicIP.Count)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCCENodeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	var updateOpts nodes.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Metadata.Name = d.Get("name").(string)
	}

	clusterid := d.Get("cluster_id").(string)
	_, err = nodes.Update(nodeClient, clusterid, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud  Node: %s", err)
	}

	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	err = nodes.Delete(nodeClient, clusterid, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE Cluster: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"Deleted"},
		Refresh:    waitForCceNodeDelete(nodeClient, clusterid, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE Node: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceNodeActive(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodes.Get(cceClient, clusterId, nodeId).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func waitForCceNodeDelete(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud CCE Node %s.\n", nodeId)

		r, err := nodes.Get(cceClient, clusterId, nodeId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud CCE Node %s", nodeId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		log.Printf("[DEBUG] HuaweiCloud CCE Node %s still available.\n", nodeId)
		return r, r.Status.Phase, nil
	}
}

func waitForClusterAvailable(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Waiting for HuaweiCloud Cluster to be available %s.\n", clusterId)
		n, err := clusters.Get(cceClient, clusterId).Extract()

		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func recursiveCreate(cceClient *golangsdk.ServiceClient, opts nodes.CreateOptsBuilder, ClusterID string, errCode int) (*nodes.Nodes, string) {
	if errCode == 403 {
		stateCluster := &resource.StateChangeConf{
			Target:     []string{"Available"},
			Refresh:    waitForClusterAvailable(cceClient, ClusterID),
			Timeout:    15 * time.Minute,
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, stateErr := stateCluster.WaitForState()
		if stateErr != nil {
			log.Printf("[INFO] Cluster Unavailable %s.\n", stateErr)
		}
		s, err := nodes.Create(cceClient, ClusterID, opts).Extract()
		if err != nil {
			//if err.(golangsdk.ErrUnexpectedResponseCode).Actual == 403 {
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				return recursiveCreate(cceClient, opts, ClusterID, 403)
			} else {
				return s, "fail"
			}
		} else {
			return s, "success"
		}
	}
	return nil, "fail"
}
