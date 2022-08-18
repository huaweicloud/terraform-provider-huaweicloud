package huaweicloud

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCCENodeAttachV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCCENodeAttachV3Create,
		ReadContext:   resourceCCENodeV3Read,
		UpdateContext: resourceCCENodeAttachV3Update,
		DeleteContext: resourceCCENodeAttachV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
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
			},
			"key_pair": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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
			},
			"nic_multi_queue": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nic_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"preinstall": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				StateFunc: utils.DecodeHashAndHexEncode,
			},
			"postinstall": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				StateFunc: utils.DecodeHashAndHexEncode,
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
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					}},
			},
			//(k8s_tags)
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			//(node/ecs_tags)
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
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
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

func resourceCCENodeAttachV3ServerConfig(d *schema.ResourceData) *nodes.ServerConfig {
	if hasFilledOpt(d, "tags") || hasFilledOpt(d, "image_id") {
		serverConfig := nodes.ServerConfig{
			UserTags: resourceCCENodeTags(d),
		}

		if v, ok := d.GetOk("image_id"); ok {
			rootVolume := nodes.RootVolume{
				ImageID: v.(string),
			}
			serverConfig.RootVolume = &rootVolume
		}
		return &serverConfig
	}
	return nil
}

func resourceCCENodeAttachV3VolumeConfig(d *schema.ResourceData) *nodes.VolumeConfig {
	if v, ok := d.GetOk("lmv_config"); ok {
		volumeConfig := nodes.VolumeConfig{
			LvmConfig: v.(string),
		}
		return &volumeConfig
	}
	return nil
}

func resourceCCENodeAttachV3RuntimeConfig(d *schema.ResourceData) *nodes.RuntimeConfig {
	if v, ok := d.GetOk("docker_base_size"); ok {
		runtimeConfig := nodes.RuntimeConfig{
			DockerBaseSize: v.(int),
		}
		return &runtimeConfig
	}
	return nil
}

func resourceCCENodeAttachV3K8sOptions(d *schema.ResourceData) *nodes.K8sOptions {
	if hasFilledOpt(d, "labels") || hasFilledOpt(d, "taints") || hasFilledOpt(d, "max_pods") ||
		hasFilledOpt(d, "nic_multi_queue") || hasFilledOpt(d, "nic_threshold") {
		k8sOptions := nodes.K8sOptions{
			Labels:        resourceCCENodeK8sTags(d),
			Taints:        resourceCCETaint(d),
			MaxPods:       d.Get("max_pods").(int),
			NicMultiQueue: d.Get("nic_multi_queue").(string),
			NicThreshold:  d.Get("nic_threshold").(string),
		}
		return &k8sOptions
	}

	return nil
}

func resourceCCENodeAttachV3Lifecycle(d *schema.ResourceData) *nodes.Lifecycle {
	if hasFilledOpt(d, "preinstall") || hasFilledOpt(d, "postinstall") {
		lifecycle := nodes.Lifecycle{
			Preinstall:  d.Get("preinstall").(string),
			PostInstall: d.Get("postinstall").(string),
		}
		return &lifecycle
	}
	return nil
}

func resourceCCENodeAttachV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE Node client: %s", err)
	}

	// wait for the cce cluster to become available
	clusterID := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:       []string{"Available"},
		Refresh:      waitForClusterAvailable(nodeClient, clusterID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateCluster.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error waiting for HuaweiCloud CCE cluster to be Available: %s", err)
	}

	addOpts := nodes.AddOpts{
		Kind:       "List",
		ApiVersion: "v3",
	}

	addNode := nodes.AddNode{
		ServerID: d.Get("server_id").(string),
		Spec: nodes.AddNodeSpec{
			Os:            d.Get("os").(string),
			Name:          d.Get("name").(string),
			ServerConfig:  resourceCCENodeAttachV3ServerConfig(d),
			VolumeConfig:  resourceCCENodeAttachV3VolumeConfig(d),
			RuntimeConfig: resourceCCENodeAttachV3RuntimeConfig(d),
			K8sOptions:    resourceCCENodeAttachV3K8sOptions(d),
			Lifecycle:     resourceCCENodeAttachV3Lifecycle(d),
		},
	}

	addOpts.NodeList = append(addOpts.NodeList, addNode)

	logp.Printf("[DEBUG] Add node Options: %#v", addOpts)
	// Add loginSpec here so it wouldn't go in the above log entry
	var loginSpec nodes.LoginSpec
	if hasFilledOpt(d, "key_pair") {
		loginSpec = nodes.LoginSpec{
			SshKey: d.Get("key_pair").(string),
		}
	} else {
		password, err := utils.TryPasswordEncrypt(d.Get("password").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: password,
			},
		}
	}
	addOpts.NodeList[0].Spec.Login = loginSpec

	s, err := nodes.Add(nodeClient, clusterID, addOpts).ExtractAddNode()
	if err != nil {
		return fmtp.DiagErrorf("Error adding HuaweiCloud Node: %s", err)
	}

	nodeID, err := getResourceIDFromJob(ctx, nodeClient, s.JobID, "CreateNode", "InstallNode",
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(nodeID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Build", "Installing"},
		Target:       []string{"Active"},
		Refresh:      waitForCceNodeActive(nodeClient, clusterID, nodeID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error adding HuaweiCloud CCE Node: %s", err)
	}

	return resourceCCENodeV3Read(ctx, d, meta)
}

func resourceCCENodeAttachV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	if d.HasChanges("os", "key_pair", "password") {
		nodeClient, err := config.CceV3Client(GetRegion(d, config))
		if err != nil {
			return fmtp.DiagErrorf("Error creating HuaweiCloud CCE client: %s", err)
		}
		clusterID := d.Get("cluster_id").(string)
		resetOpts := nodes.ResetOpts{
			Kind:       "List",
			ApiVersion: "v3",
		}

		resetNode := nodes.ResetNode{
			NodeID: d.Id(),
			Spec: nodes.AddNodeSpec{
				Os:            d.Get("os").(string),
				Name:          d.Get("name").(string),
				ServerConfig:  resourceCCENodeAttachV3ServerConfig(d),
				VolumeConfig:  resourceCCENodeAttachV3VolumeConfig(d),
				RuntimeConfig: resourceCCENodeAttachV3RuntimeConfig(d),
				K8sOptions:    resourceCCENodeAttachV3K8sOptions(d),
				Lifecycle:     resourceCCENodeAttachV3Lifecycle(d),
			},
		}

		resetOpts.NodeList = append(resetOpts.NodeList, resetNode)

		logp.Printf("[DEBUG] Reset node Options: %#v", resetOpts)
		// Add loginSpec here so it wouldn't go in the above log entry
		var loginSpec nodes.LoginSpec
		if hasFilledOpt(d, "key_pair") {
			loginSpec = nodes.LoginSpec{
				SshKey: d.Get("key_pair").(string),
			}
		} else {
			password, err := utils.TryPasswordEncrypt(d.Get("password").(string))
			if err != nil {
				return diag.FromErr(err)
			}
			loginSpec = nodes.LoginSpec{
				UserPassword: nodes.UserPassword{
					Username: "root",
					Password: password,
				},
			}
		}
		resetOpts.NodeList[0].Spec.Login = loginSpec

		s, err := nodes.Reset(nodeClient, clusterID, resetOpts).ExtractAddNode()
		if err != nil {
			return fmtp.DiagErrorf("Error resetting HuaweiCloud Node: %s", err)
		}

		nodeID, err := getResourceIDFromJob(ctx, nodeClient, s.JobID, "CreateNode", "InstallNode",
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(nodeID)

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Build", "Installing"},
			Target:       []string{"Active"},
			Refresh:      waitForCceNodeActive(nodeClient, clusterID, nodeID),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        20 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmtp.DiagErrorf("Error resetting HuaweiCloud CCE Node: %s", err)
		}

		return resourceCCENodeV3Read(ctx, d, config)

	} else {
		return resourceCCENodeV3Update(ctx, d, config)
	}
}

func resourceCCENodeAttachV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)

	var removeOpts nodes.RemoveOpts
	var loginSpec nodes.LoginSpec

	if hasFilledOpt(d, "key_pair") {
		loginSpec = nodes.LoginSpec{
			SshKey: d.Get("key_pair").(string),
		}
	} else {
		password, err := utils.TryPasswordEncrypt(d.Get("password").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: password,
			},
		}
	}
	removeOpts.Spec.Login = loginSpec

	nodeItem := nodes.NodeItem{
		Uid: d.Id(),
	}
	removeOpts.Spec.Nodes = append(removeOpts.Spec.Nodes, nodeItem)

	err = nodes.Remove(nodeClient, clusterID, removeOpts).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error removing HuaweiCloud CCE node: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCceNodeDelete(nodeClient, clusterID, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud CCE Node: %s", err)
	}

	d.SetId("")
	return nil
}
