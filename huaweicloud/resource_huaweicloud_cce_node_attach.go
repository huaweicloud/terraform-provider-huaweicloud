package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCCENodeAttachV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodeAttachV3Create,
		Read:   resourceCCENodeV3Read,
		Update: resourceCCENodeV3Update,
		Delete: resourceCCENodeAttachV3Delete,

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
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  50,
			},
			"lvm_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"docker_base_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  10,
			},
			"nic_multi_queue": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "[{\"queue\":4}]",
			},
			"nic_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "0.3:0.6",
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"tags": tagsSchema(),

			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hw_passthrough": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hw_passthrough": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
			},
			"runtime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ecs_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCCENodeAttachV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE Node client: %s", err)
	}

	// wait for the cce cluster to become available
	clusterid := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:       []string{"Available"},
		Refresh:      waitForClusterAvailable(nodeClient, clusterid),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateCluster.WaitForState()

	addOpts := nodes.AddOpts{
		Kind:       "Node",
		ApiVersion: "v3",
	}

	addNode := nodes.AddNode{
		ServerID: d.Get("server_id").(string),
		Spec: nodes.AddNodeSpec{
			Os: d.Get("os").(string),
		},
	}

	if v, ok := d.GetOk("lmv_config"); ok {
		volumeConfig := nodes.VolumeConfig{
			LvmConfig: v.(string),
		}
		addNode.Spec.VolumeConfig = &volumeConfig
	}
	if v, ok := d.GetOk("docker_base_size"); ok {
		runtimeConfig := nodes.RuntimeConfig{
			DockerBaseSize: v.(int),
		}
		addNode.Spec.RuntimeConfig = &runtimeConfig
	}

	k8sOptions := nodes.K8sOptions{
		MaxPods:       d.Get("max_pods").(int),
		NicMultiQueue: d.Get("nic_multi_queue").(string),
		NicThreshold:  d.Get("nic_threshold").(string),
	}
	if (k8sOptions != nodes.K8sOptions{}) {
		addNode.Spec.K8sOptions = &k8sOptions
	}

	if v, ok := d.GetOk("image_id"); ok {
		extendParam := map[string]interface{}{
			"alpha.cce/NodeImageID": v.(int),
		}
		addNode.Spec.ExtendParam = extendParam
	}

	if hasFilledOpt(d, "preinstall") || hasFilledOpt(d, "postinstall") {
		lifecycle := nodes.Lifecycle{
			Preinstall:  d.Get("preinstall").(string),
			PostInstall: d.Get("postinstall").(string),
		}
		addNode.Spec.Lifecycle = &lifecycle
	}

	addOpts.NodeList = append(addOpts.NodeList, addNode)

	logp.Printf("[DEBUG] Add node Options: %#v", addOpts)
	// Add loginSpec here so it wouldn't go in the above log entry
	var loginSpec nodes.LoginSpec
	if hasFilledOpt(d, "key_pair") {
		loginSpec = nodes.LoginSpec{
			SshKey: d.Get("key_pair").(string),
		}
	} else if hasFilledOpt(d, "password") {
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: d.Get("password").(string),
			},
		}
	}
	addOpts.NodeList[0].Spec.Login = loginSpec

	s, err := nodes.Add(nodeClient, clusterid, addOpts).ExtractAddNode()

	if err != nil {
		return fmtp.Errorf("Error adding HuaweiCloud Node: %s", err)
	}

	nodeID, err := getResourceIDFromJob(nodeClient, s.JobID, "CreateNode", "InstallNode")
	if err != nil {
		return err
	}
	d.SetId(nodeID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Build", "Installing"},
		Target:       []string{"Active"},
		Refresh:      waitForCceNodeActive(nodeClient, clusterid, nodeID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error adding HuaweiCloud CCE Node: %s", err)
	}

	return resourceCCENodeV3Update(d, meta)
}

func resourceCCENodeAttachV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	clusterid := d.Get("cluster_id").(string)

	var removeOpts nodes.RemoveOpts

	var loginSpec nodes.LoginSpec
	if hasFilledOpt(d, "key_pair") {
		loginSpec = nodes.LoginSpec{
			SshKey: d.Get("key_pair").(string),
		}
	} else if hasFilledOpt(d, "password") {
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: d.Get("password").(string),
			},
		}
	}
	removeOpts.Spec.Login = loginSpec

	nodeItem := nodes.NodeItem{
		Uid: d.Id(),
	}
	removeOpts.Spec.Nodes = append(removeOpts.Spec.Nodes, nodeItem)

	err = nodes.Remove(nodeClient, clusterid, removeOpts).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error removing HuaweiCloud CCE node: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCceNodeDelete(nodeClient, clusterid, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud CCE Node: %s", err)
	}

	d.SetId("")
	return nil
}
