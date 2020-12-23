package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodepools"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
)

func ResourceCCENodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodePoolCreate,
		Read:   resourceCCENodePoolRead,
		Update: resourceCCENodePoolUpdate,
		Delete: resourceCCENodePoolDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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
			"initial_node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"labels": { //(k8s_tags)
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "random",
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_pair": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					}},
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"preinstall": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return installScriptHashSum(v.(string))
					default:
						return ""
					}
				},
			},
			"postinstall": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return installScriptHashSum(v.(string))
					default:
						return ""
					}
				},
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"scall_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"min_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scale_down_cooldown_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCCENodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Node Pool client: %s", err)
	}

	var loginSpec nodes.LoginSpec
	if hasFilledOpt(d, "key_pair") {
		loginSpec = nodes.LoginSpec{SshKey: d.Get("key_pair").(string)}
	} else if hasFilledOpt(d, "password") {
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: d.Get("password").(string),
			},
		}
	}

	var base64PreInstall, base64PostInstall string
	if v, ok := d.GetOk("preinstall"); ok {
		base64PreInstall = installScriptEncode(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		base64PostInstall = installScriptEncode(v.(string))
	}

	initialNodeCount := d.Get("initial_node_count").(int)

	createOpts := nodepools.CreateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.CreateMetaData{
			Name: d.Get("name").(string),
		},
		Spec: nodepools.CreateSpec{
			Type: d.Get("type").(string),
			NodeTemplate: nodes.Spec{
				Flavor:      d.Get("flavor_id").(string),
				Az:          d.Get("availability_zone").(string),
				Os:          d.Get("os").(string),
				Login:       loginSpec,
				RootVolume:  resourceCCERootVolume(d),
				DataVolumes: resourceCCEDataVolume(d),
				K8sTags:     resourceCCENodeK8sTags(d),
				BillingMode: 0,
				Count:       1,
				NodeNicSpec: nodes.NodeNicSpec{
					PrimaryNic: nodes.PrimaryNic{
						SubnetId: d.Get("subnet_id").(string),
					},
				},
				ExtendParam: nodes.ExtendParam{
					PreInstall:  base64PreInstall,
					PostInstall: base64PostInstall,
				},
				Taints: resourceCCETaint(d),
			},
			Autoscaling: nodepools.AutoscalingSpec{
				Enable:                d.Get("scall_enable").(bool),
				MinNodeCount:          d.Get("min_node_count").(int),
				MaxNodeCount:          d.Get("max_node_count").(int),
				ScaleDownCooldownTime: d.Get("scale_down_cooldown_time").(int),
				Priority:              d.Get("priority").(int),
			},
			InitialNodeCount: &initialNodeCount,
		},
	}

	clusterid := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:     []string{"Available"},
		Refresh:    waitForClusterAvailable(nodePoolClient, clusterid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err = stateCluster.WaitForState()

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	s, err := nodepools.Create(nodePoolClient, clusterid, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveNodePoolCreate(nodePoolClient, createOpts, clusterid, 403)
			if err == "fail" {
				return fmt.Errorf("Error creating HuaweiCloud Node Pool")
			}
			s = retryNode
		} else {
			return fmt.Errorf("Error creating HuaweiCloud Node Pool: %s", err)
		}
	}

	if len(s.Metadata.Id) == 0 {
		return fmt.Errorf("Error fetching CreateNodePool id")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Synchronizing"},
		Target:       []string{""},
		Refresh:      waitForCceNodePoolActive(nodePoolClient, clusterid, s.Metadata.Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Node Pool: %s", err)
	}

	log.Printf("[DEBUG] Create node pool: %v", s)

	d.SetId(s.Metadata.Id)
	return resourceCCENodePoolRead(d, meta)
}

func resourceCCENodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Node Pool client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	s, err := nodepools.Get(nodePoolClient, clusterid, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HuaweiCloud Node Pool: %s", err)
	}

	d.Set("name", s.Metadata.Name)
	d.Set("flavor_id", s.Spec.NodeTemplate.Flavor)
	d.Set("availability_zone", s.Spec.NodeTemplate.Az)
	d.Set("os", s.Spec.NodeTemplate.Os)
	d.Set("billing_mode", s.Spec.NodeTemplate.BillingMode)
	d.Set("key_pair", s.Spec.NodeTemplate.Login.SshKey)
	d.Set("initial_node_count", s.Spec.InitialNodeCount)
	d.Set("scall_enable", s.Spec.Autoscaling.Enable)
	d.Set("min_node_count", s.Spec.Autoscaling.MinNodeCount)
	d.Set("max_node_count", s.Spec.Autoscaling.MaxNodeCount)
	d.Set("scale_down_cooldown_time", s.Spec.Autoscaling.ScaleDownCooldownTime)
	d.Set("priority", s.Spec.Autoscaling.Priority)

	labels := map[string]string{}
	for key, val := range s.Spec.NodeTemplate.K8sTags {
		if strings.Contains(key, "cce.cloud.com") {
			continue
		}
		labels[key] = val
	}
	d.Set("labels", labels)

	var volumes []map[string]interface{}
	for _, pairObject := range s.Spec.NodeTemplate.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_param"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving dataVolumes to state for HuaweiCloud Node Pool (%s): %s", d.Id(), err)
	}

	rootVolume := []map[string]interface{}{
		{
			"size":         s.Spec.NodeTemplate.RootVolume.Size,
			"volumetype":   s.Spec.NodeTemplate.RootVolume.VolumeType,
			"extend_param": s.Spec.NodeTemplate.RootVolume.ExtendParam,
		},
	}
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmt.Errorf("[DEBUG] Error saving root Volume to state for HuaweiCloud Node Pool (%s): %s", d.Id(), err)
	}

	d.Set("status", s.Status.Phase)

	return nil
}

func resourceCCENodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	initialNodeCount := d.Get("initial_node_count").(int)

	updateOpts := nodepools.UpdateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.UpdateMetaData{
			Name: d.Get("name").(string),
		},
		Spec: nodepools.UpdateSpec{
			InitialNodeCount: &initialNodeCount,
			Autoscaling: nodepools.AutoscalingSpec{
				Enable:                d.Get("scall_enable").(bool),
				MinNodeCount:          d.Get("min_node_count").(int),
				MaxNodeCount:          d.Get("max_node_count").(int),
				ScaleDownCooldownTime: d.Get("scale_down_cooldown_time").(int),
				Priority:              d.Get("priority").(int),
			},
			Type: d.Get("type").(string),
		},
	}

	clusterid := d.Get("cluster_id").(string)
	_, err = nodepools.Update(nodePoolClient, clusterid, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud Node Node Pool: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Synchronizing"},
		Target:     []string{""},
		Refresh:    waitForCceNodePoolActive(nodePoolClient, clusterid, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Node Pool: %s", err)
	}

	return resourceCCENodePoolRead(d, meta)
}

func resourceCCENodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	err = nodepools.Delete(nodePoolClient, clusterid, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE Node Pool: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCceNodePoolDelete(nodePoolClient, clusterid, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE Node Pool: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceNodePoolActive(cceClient *golangsdk.ServiceClient, clusterId, nodePoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodepools.Get(cceClient, clusterId, nodePoolId).Extract()
		if err != nil {
			return nil, "", err
		}
		return n, n.Status.Phase, nil
	}
}

func waitForCceNodePoolDelete(cceClient *golangsdk.ServiceClient, clusterId, nodePoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud CCE Node Pool %s.\n", nodePoolId)

		r, err := nodepools.Get(cceClient, clusterId, nodePoolId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud CCE Node Pool %s", nodePoolId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		log.Printf("[DEBUG] HuaweiCloud CCE Node Pool %s still available.\n", nodePoolId)
		return r, r.Status.Phase, nil
	}
}

func recursiveNodePoolCreate(cceClient *golangsdk.ServiceClient, opts nodepools.CreateOptsBuilder, ClusterID string, errCode int) (*nodepools.NodePool, string) {
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
		s, err := nodepools.Create(cceClient, ClusterID, opts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				return recursiveNodePoolCreate(cceClient, opts, ClusterID, 403)
			} else {
				return s, "fail"
			}
		} else {
			return s, "success"
		}
	}
	return nil, "fail"
}
