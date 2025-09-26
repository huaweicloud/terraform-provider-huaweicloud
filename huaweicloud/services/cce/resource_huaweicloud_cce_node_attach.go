package cce

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/add
// @API CCE GET /api/v3/projects/{project_id}/jobs/{job_id}
// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/{node_id}
// @API ECS GET /v1/{project_id}/cloudservers/{id}/tags
// @API ECS GET /v1/{project_id}/cloudservers/{server_id}
// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/reset
// @API CCE PUT /api/v3/projects/{project_id}/clusters/{clusterid}/nodes/{node_id}
// @API ECS PUT /v1/{project_id}/cloudservers/{server_id}/os-reset-password
// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/tags/action
// @API KMS POST /v3/{project_id}/keypairs/associate
// @API KMS POST /v3/{project_id}/keypairs/disassociate
// @API CCE PUT /api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/operation/remove

var nodeAttachSchema = map[string]*schema.Schema{
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
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	},
	"private_key": {
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	},
	"max_pods": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"lvm_config": {
		Type:     schema.TypeString,
		Optional: true,
		ConflictsWith: []string{
			"storage",
		},
	},
	"docker_base_size": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"runtime": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"nic_multi_queue": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "schema: Internal",
	},
	"nic_threshold": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "schema: Internal",
	},
	"image_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "schema: Internal",
	},
	"system_disk_kms_key_id": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"preinstall": {
		Type:      schema.TypeString,
		Optional:  true,
		StateFunc: utils.DecodeHashAndHexEncode,
	},
	"postinstall": {
		Type:      schema.TypeString,
		Optional:  true,
		StateFunc: utils.DecodeHashAndHexEncode,
	},
	"initialized_conditions": {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"storage": resourceNodeStorageSchema(),
	"taints": {
		Type:     schema.TypeList,
		Optional: true,
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
	// (k8s_tags)
	"labels": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	// (node/ecs_tags)
	"tags": common.TagsSchema(),
	"hostname_config": {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		},
	},

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
				"extend_params": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"kms_key_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"dss_pool_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"iops": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"throughput": {
					Type:     schema.TypeInt,
					Computed: true,
				},

				// Internal parameters
				"hw_passthrough": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "schema: Internal",
				},

				// Deprecated parameters
				"extend_param": {
					Type:       schema.TypeString,
					Computed:   true,
					Deprecated: "use extend_params instead",
				},
			},
		},
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
				"extend_params": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"kms_key_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"dss_pool_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"iops": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"throughput": {
					Type:     schema.TypeInt,
					Computed: true,
				},

				// Internal parameters
				"hw_passthrough": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "schema: Internal",
				},

				// Deprecated parameters
				"extend_param": {
					Type:       schema.TypeString,
					Computed:   true,
					Deprecated: "use extend_params instead",
				},
			},
		},
	},
	"ecs_group_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"subnet_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"extension_nics": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"subnet_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
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
	"charging_mode": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enterprise_project_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ResourceNodeAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeAttachCreate,
		ReadContext:   resourceNodeAttachRead,
		UpdateContext: resourceNodeAttachUpdate,
		DeleteContext: resourceNodeAttachDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: nodeAttachSchema,
	}
}

func resourceNodeAttachServerConfig(d *schema.ResourceData) *nodes.ServerConfig {
	var res nodes.ServerConfig
	if common.HasFilledOpt(d, "tags") {
		res.UserTags = buildResourceNodeTags(d)
	}

	if common.HasFilledOpt(d, "image_id") || common.HasFilledOpt(d, "system_disk_kms_key_id") {
		rootVolume := nodes.RootVolume{
			ImageID: d.Get("image_id").(string),
			CmkID:   d.Get("system_disk_kms_key_id").(string),
		}
		res.RootVolume = &rootVolume
	}

	return &res
}

func resourceNodeAttachVolumeConfig(d *schema.ResourceData) *nodes.VolumeConfig {
	// only one of lvm_config and storage can be specified
	if v, ok := d.GetOk("lvm_config"); ok {
		volumeConfig := nodes.VolumeConfig{
			LvmConfig: v.(string),
		}
		return &volumeConfig
	}

	if _, ok := d.GetOk("storage"); ok {
		volumeConfig := nodes.VolumeConfig{
			Storage: buildResourceNodeStorage(d),
		}
		return &volumeConfig
	}
	return nil
}

func resourceNodeAttachRuntimeConfig(d *schema.ResourceData) *nodes.RuntimeConfig {
	var res nodes.RuntimeConfig

	if v, ok := d.GetOk("docker_base_size"); ok {
		res.DockerBaseSize = v.(int)
	}

	if v, ok := d.GetOk("runtime"); ok {
		res.Runtime = &nodes.RunTimeSpec{
			Name: v.(string),
		}
	}

	return &res
}

func resourceNodeAttachK8sOptions(d *schema.ResourceData) *nodes.K8sOptions {
	if common.HasFilledOpt(d, "labels") || common.HasFilledOpt(d, "taints") || common.HasFilledOpt(d, "max_pods") ||
		common.HasFilledOpt(d, "nic_multi_queue") || common.HasFilledOpt(d, "nic_threshold") {
		k8sOptions := nodes.K8sOptions{
			Labels:        buildResourceNodeK8sTags(d),
			Taints:        buildResourceNodeTaint(d),
			MaxPods:       d.Get("max_pods").(int),
			NicMultiQueue: d.Get("nic_multi_queue").(string),
			NicThreshold:  d.Get("nic_threshold").(string),
		}
		return &k8sOptions
	}

	return nil
}

func resourceNodeAttachLifecycle(d *schema.ResourceData) *nodes.Lifecycle {
	if common.HasFilledOpt(d, "preinstall") || common.HasFilledOpt(d, "postinstall") {
		lifecycle := nodes.Lifecycle{
			Preinstall:  d.Get("preinstall").(string),
			PostInstall: d.Get("postinstall").(string),
		}
		return &lifecycle
	}
	return nil
}

func buildNodeAttachCreateOpts(d *schema.ResourceData) (*nodes.AddOpts, error) {
	result := nodes.AddOpts{
		Kind:       "List",
		ApiVersion: "v3",
		NodeList: []nodes.AddNode{
			{
				ServerID: d.Get("server_id").(string),
				Spec: nodes.AddNodeSpec{
					Os:                    d.Get("os").(string),
					Name:                  d.Get("name").(string),
					ServerConfig:          resourceNodeAttachServerConfig(d),
					VolumeConfig:          resourceNodeAttachVolumeConfig(d),
					RuntimeConfig:         resourceNodeAttachRuntimeConfig(d),
					K8sOptions:            resourceNodeAttachK8sOptions(d),
					Lifecycle:             resourceNodeAttachLifecycle(d),
					InitializedConditions: utils.ExpandToStringList(d.Get("initialized_conditions").([]interface{})),
					HostnameConfig:        buildResourceNodeHostnameConfig(d),
				},
			},
		},
	}

	log.Printf("[DEBUG] Add node Options: %#v", result)
	// Add loginSpec here so it wouldn't go in the above log entry
	loginSpec, err := buildResourceNodeLoginSpec(d)
	if err != nil {
		diag.FromErr(err)
	}
	result.NodeList[0].Spec.Login = loginSpec
	return &result, nil
}

func resourceNodeAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cceClient, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	// Wait for the cce cluster to become available
	clusterID := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(cceClient, clusterID, []string{"Available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateCluster.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE cluster to become available: %s", err)
	}

	addOpts, err := buildNodeAttachCreateOpts(d)
	if err != nil {
		return diag.Errorf("error creating AddOpts structure of 'Add' method for CCE node attach: %s", err)
	}
	resp, err := nodes.Add(cceClient, clusterID, addOpts).ExtractAddNode()
	if err != nil {
		return diag.Errorf("error adding node to the cluster (%s): %s", clusterID, err)
	}

	nodeID, err := getResourceIDFromJob(ctx, cceClient, resp.JobID, "CreateNode", "InstallNode",
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(nodeID)

	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "Build" and "Installing".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      nodeStateRefreshFunc(cceClient, clusterID, nodeID, []string{"Active"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE node attach to the cluster: %s", err)
	}

	return resourceNodeRead(ctx, d, meta)
}

func buildNodeAttachUpdateOpts(d *schema.ResourceData) (*nodes.ResetOpts, error) {
	result := nodes.ResetOpts{
		Kind:       "List",
		ApiVersion: "v3",
		NodeList: []nodes.ResetNode{
			{
				NodeID: d.Id(),
				Spec: nodes.AddNodeSpec{
					Os:                    d.Get("os").(string),
					Name:                  d.Get("name").(string),
					ServerConfig:          resourceNodeAttachServerConfig(d),
					VolumeConfig:          resourceNodeAttachVolumeConfig(d),
					RuntimeConfig:         resourceNodeAttachRuntimeConfig(d),
					K8sOptions:            resourceNodeAttachK8sOptions(d),
					Lifecycle:             resourceNodeAttachLifecycle(d),
					InitializedConditions: utils.ExpandToStringList(d.Get("initialized_conditions").([]interface{})),
				},
			},
		},
	}

	log.Printf("[DEBUG] Add node Options: %#v", result)
	// Add loginSpec here so it wouldn't go in the above log entry
	loginSpec, err := buildResourceNodeLoginSpec(d)
	if err != nil {
		diag.FromErr(err)
	}
	result.NodeList[0].Spec.Login = loginSpec
	return &result, nil
}

func resourceNodeAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	nodeClient, err := cfg.CceV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CCE Node client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	s, err := nodes.Get(nodeClient, clusterid, d.Id()).Extract()

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE Node")
	}

	// The following parameters are not returned:
	// password, private_key, fixed_ip, extension_nics, eip_id, iptype, bandwidth_charge_mode, bandwidth_size,
	// sharetype, extend_params, dedicated_host_id, initialized_conditions, labels, taints, period_unit, period, auto_renew, auto_pay
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", s.Metadata.Name),
		d.Set("flavor_id", s.Spec.Flavor),
		d.Set("availability_zone", s.Spec.Az),
		d.Set("os", s.Spec.Os),
		d.Set("key_pair", s.Spec.Login.SshKey),
		d.Set("subnet_id", s.Spec.NodeNicSpec.PrimaryNic.SubnetId),
		d.Set("ecs_group_id", s.Spec.EcsGroupID),
		d.Set("server_id", s.Status.ServerID),
		d.Set("private_ip", s.Status.PrivateIP),
		d.Set("public_ip", s.Status.PublicIP),
		d.Set("status", s.Status.Phase),
		d.Set("root_volume", flattenResourceNodeRootVolume(d, s.Spec.RootVolume)),
		d.Set("data_volumes", flattenResourceNodeDataVolume(d, s.Spec.DataVolumes)),
		d.Set("initialized_conditions", s.Spec.InitializedConditions),
		d.Set("hostname_config", flattenResourceNodeHostnameConfig(s.Spec.HostnameConfig)),
		d.Set("enterprise_project_id", s.Spec.ServerEnterpriseProjectID),
		d.Set("tags", utils.TagsToMap(s.Spec.UserTags)),
		d.Set("storage", flattenResourceNodeStorage(s.Spec.Storage)),
		d.Set("extension_nics", flattenExtensionNics(s.Spec.NodeNicSpec.ExtNics)),
	)

	if s.Spec.BillingMode != 0 {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	}
	if s.Spec.RunTime != nil {
		mErr = multierror.Append(mErr, d.Set("runtime", s.Spec.RunTime.Name))
	}

	computeClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	serverId := s.Status.ServerID
	// fetch key_pair from ECS instance
	if server, err := cloudservers.Get(computeClient, serverId).Extract(); err == nil {
		mErr = multierror.Append(mErr, d.Set("key_pair", server.KeyName))
	} else {
		log.Printf("[WARN] Error fetching ECS instance (%s): %s", serverId, err)
	}

	// fetch tags from ECS instance
	if resourceTags, err := tags.Get(computeClient, "cloudservers", serverId).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagmap))
	} else {
		log.Printf("[WARN] Error fetching tags of ECS instance (%s): %s", serverId, err)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCE Node fields: %s", err)
	}
	return nil
}

func resourceNodeAttachUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	if d.HasChanges("name", "tags", "key_pair", "password") {
		return resourceNodeUpdate(ctx, d, cfg)
	}

	cceClient, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}
	clusterID := d.Get("cluster_id").(string)

	resetOpts, err := buildNodeAttachUpdateOpts(d)
	if err != nil {
		return diag.Errorf("error creating ResetOpts structure of 'Reset' method for CCE node attach: %s", err)
	}
	resp, err := nodes.Reset(cceClient, clusterID, resetOpts).ExtractAddNode()
	if err != nil {
		return diag.Errorf("error resetting node: %s", err)
	}

	nodeID, err := getResourceIDFromJob(ctx, cceClient, resp.JobID, "CreateNode", "InstallNode",
		d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(nodeID)

	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "Build" and "Installing".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      nodeStateRefreshFunc(cceClient, clusterID, nodeID, []string{"Active"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE Node reset complete: %s", err)
	}

	return resourceNodeRead(ctx, d, cfg)
}

func resourceNodeAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cceClient, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)

	var removeOpts nodes.RemoveOpts
	var loginSpec nodes.LoginSpec

	loginSpec, err = buildResourceNodeLoginSpec(d)
	if err != nil {
		diag.FromErr(err)
	}
	removeOpts.Spec.Login = loginSpec

	nodeItem := nodes.NodeItem{
		Uid: d.Id(),
	}
	removeOpts.Spec.Nodes = append(removeOpts.Spec.Nodes, nodeItem)

	err = nodes.Remove(cceClient, clusterID, removeOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("error removing CCE node: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "Deleting".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      nodeStateRefreshFunc(cceClient, clusterID, d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting CCE Node: %s", err)
	}
	return nil
}
