package huaweicloud

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCCENodeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodeV3Create,
		Read:   resourceCCENodeV3Read,
		Update: resourceCCENodeV3Update,
		Delete: resourceCCENodeV3Delete,
		Importer: &schema.ResourceImporter{
			State: resourceCCENodeV3Import,
		},

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
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
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
				Optional: true,
				Computed: true,
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
						"hw_passthrough": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"extend_param": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "use extend_params instead",
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
						"hw_passthrough": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"extend_param": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "use extend_params instead",
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
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
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"eip_ids", "iptype", "bandwidth_charge_mode", "bandwidth_size", "sharetype",
				},
			},
			"iptype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
			},
			"bandwidth_charge_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_ids", "eip_id"},
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"docker", "containerd",
				}, false),
			},
			"ecs_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ecs_performance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"public_key": {
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
			"labels": { //(k8s_tags)
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			//(node/ecs_tags)
			"tags": tagsSchema(),
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// charge info: charging_mode, period_unit, period, auto_renew
			"charging_mode": schemeChargingMode(nil),
			"period_unit":   schemaPeriodUnit(nil),
			"period":        schemaPeriod(nil),
			"auto_renew":    schemaAutoRenew(nil),

			"extend_param": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"keep_ecs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Deprecated
			"eip_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				ConflictsWith: []string{
					"eip_id", "iptype", "bandwidth_charge_mode", "bandwidth_size", "sharetype",
				},
				Deprecated: "use eip_id instead",
			},
			"billing_mode": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "use charging_mode instead",
			},
			"extend_param_charging_mode": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "use charging_mode instead",
			},
			"order_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "will be removed after v1.26.0",
			},
		},
	}
}

func resourceCCENodeAnnotationsV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceCCENodeK8sTags(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceCCENodeTags(d *schema.ResourceData) []tags.ResourceTag {
	tagRaw := d.Get("tags").(map[string]interface{})
	return utils.ExpandResourceTags(tagRaw)
}

func resourceCCERootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var root nodes.VolumeSpec
	volumeRaw := d.Get("root_volume").([]interface{})
	if len(volumeRaw) == 1 {
		rawMap := volumeRaw[0].(map[string]interface{})
		root.Size = rawMap["size"].(int)
		root.VolumeType = rawMap["volumetype"].(string)
		root.HwPassthrough = rawMap["hw_passthrough"].(bool)
		root.ExtendParam = rawMap["extend_params"].(map[string]interface{})
	}
	return root
}

func resourceCCEDataVolume(d *schema.ResourceData) []nodes.VolumeSpec {
	volumeRaw := d.Get("data_volumes").([]interface{})
	volumes := make([]nodes.VolumeSpec, len(volumeRaw))
	for i, raw := range volumeRaw {
		rawMap := raw.(map[string]interface{})
		volumes[i] = nodes.VolumeSpec{
			Size:          rawMap["size"].(int),
			VolumeType:    rawMap["volumetype"].(string),
			HwPassthrough: rawMap["hw_passthrough"].(bool),
			ExtendParam:   rawMap["extend_params"].(map[string]interface{}),
		}
	}
	return volumes
}

func resourceCCETaint(d *schema.ResourceData) []nodes.TaintSpec {
	taintRaw := d.Get("taints").([]interface{})
	taints := make([]nodes.TaintSpec, len(taintRaw))
	for i, raw := range taintRaw {
		rawMap := raw.(map[string]interface{})
		taints[i] = nodes.TaintSpec{
			Key:    rawMap["key"].(string),
			Value:  rawMap["value"].(string),
			Effect: rawMap["effect"].(string),
		}
	}
	return taints
}

func resourceCCEEipIDs(d *schema.ResourceData) []string {
	if v, ok := d.GetOk("eip_id"); ok {
		return []string{v.(string)}
	}
	rawID := d.Get("eip_ids").(*schema.Set)
	id := make([]string, rawID.Len())
	for i, raw := range rawID.List() {
		id[i] = raw.(string)
	}
	return id
}

func resourceCCEExtendParam(d *schema.ResourceData) map[string]interface{} {
	extendParam := make(map[string]interface{})
	if v, ok := d.GetOk("extend_param"); ok {
		for key, val := range v.(map[string]interface{}) {
			extendParam[key] = val.(string)
		}
		if v, ok := extendParam["periodNum"]; ok {
			periodNum, err := strconv.Atoi(v.(string))
			if err != nil {
				logp.Printf("[WARNING] PeriodNum %s invalid, Type conversion error: %s", v.(string), err)
			}
			extendParam["periodNum"] = periodNum
		}
	}

	// assemble the charge info
	var isPrePaid bool
	var billingMode int

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		isPrePaid = true
	}
	if v, ok := d.GetOk("billing_mode"); ok {
		billingMode = v.(int)
	}
	if isPrePaid || billingMode == 2 {
		extendParam["chargingMode"] = 2
		extendParam["isAutoPay"] = "true"
		extendParam["isAutoRenew"] = "false"
	}

	if v, ok := d.GetOk("period_unit"); ok {
		extendParam["periodType"] = v.(string)
	}
	if v, ok := d.GetOk("period"); ok {
		extendParam["periodNum"] = v.(int)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		extendParam["isAutoRenew"] = v.(string)
	}

	if v, ok := d.GetOk("ecs_performance_type"); ok {
		extendParam["ecs:performancetype"] = v.(string)
	}
	if v, ok := d.GetOk("max_pods"); ok {
		extendParam["maxPods"] = v.(int)
	}
	if v, ok := d.GetOk("order_id"); ok {
		extendParam["orderID"] = v.(string)
	}
	if v, ok := d.GetOk("product_id"); ok {
		extendParam["productID"] = v.(string)
	}
	if v, ok := d.GetOk("public_key"); ok {
		extendParam["publicKey"] = v.(string)
	}
	if v, ok := d.GetOk("preinstall"); ok {
		extendParam["alpha.cce/preInstall"] = installScriptEncode(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		extendParam["alpha.cce/postInstall"] = installScriptEncode(v.(string))
	}

	return extendParam
}

func resourceCCENodeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE Node client: %s", err)
	}

	// validation
	billingMode := 0
	if d.Get("charging_mode").(string) == "prePaid" || d.Get("billing_mode").(int) == 2 {
		billingMode = 2
		if err := validatePrePaidChargeInfo(d); err != nil {
			return err
		}
	}
	// eipCount must be specified when bandwidth_size parameters was set
	eipCount := 0
	if _, ok := d.GetOk("bandwidth_size"); ok {
		eipCount = 1
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

	createOpts := nodes.CreateOpts{
		Kind:       "Node",
		ApiVersion: "v3",
		Metadata: nodes.CreateMetaData{
			Name:        d.Get("name").(string),
			Annotations: resourceCCENodeAnnotationsV2(d),
		},
		Spec: nodes.Spec{
			Flavor:      d.Get("flavor_id").(string),
			Az:          d.Get("availability_zone").(string),
			Os:          d.Get("os").(string),
			RootVolume:  resourceCCERootVolume(d),
			DataVolumes: resourceCCEDataVolume(d),
			PublicIP: nodes.PublicIPSpec{
				Ids:   resourceCCEEipIDs(d),
				Count: eipCount,
				Eip: nodes.EipSpec{
					IpType: d.Get("iptype").(string),
					Bandwidth: nodes.BandwidthOpts{
						ChargeMode: d.Get("bandwidth_charge_mode").(string),
						Size:       d.Get("bandwidth_size").(int),
						ShareType:  d.Get("sharetype").(string),
					},
				},
			},
			BillingMode: billingMode,
			Count:       1,
			NodeNicSpec: nodes.NodeNicSpec{
				PrimaryNic: nodes.PrimaryNic{
					SubnetId: d.Get("subnet_id").(string),
				},
			},
			EcsGroupID:  d.Get("ecs_group_id").(string),
			ExtendParam: resourceCCEExtendParam(d),
			Taints:      resourceCCETaint(d),
			K8sTags:     resourceCCENodeK8sTags(d),
			UserTags:    resourceCCENodeTags(d),
		},
	}

	if v, ok := d.GetOk("fixed_ip"); ok {
		createOpts.Spec.NodeNicSpec.PrimaryNic.FixedIps = []string{v.(string)}
	}
	if v, ok := d.GetOk("runtime"); ok {
		createOpts.Spec.RunTime = &nodes.RunTimeSpec{
			Name: v.(string),
		}
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
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
	createOpts.Spec.Login = loginSpec

	s, err := nodes.Create(nodeClient, clusterid, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveCreate(nodeClient, createOpts, clusterid, 403)
			if err == "fail" {
				return fmtp.Errorf("Error creating HuaweiCloud Node")
			}
			s = retryNode
		} else {
			return fmtp.Errorf("Error creating HuaweiCloud Node: %s", err)
		}
	}

	nodeID, err := getResourceIDFromJob(nodeClient, s.Status.JobID, "CreateNode", "CreateNodeVM")
	if err != nil {
		return err
	}
	d.SetId(nodeID)

	logp.Printf("[DEBUG] Waiting for CCE Node (%s) to become available", s.Metadata.Name)
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
		return fmtp.Errorf("Error creating HuaweiCloud CCE Node: %s", err)
	}

	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE Node client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	s, err := nodes.Get(nodeClient, clusterid, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud Node: %s", err)
	}

	d.Set("region", GetRegion(d, config))
	d.Set("name", s.Metadata.Name)
	d.Set("flavor_id", s.Spec.Flavor)
	d.Set("availability_zone", s.Spec.Az)
	d.Set("os", s.Spec.Os)
	d.Set("key_pair", s.Spec.Login.SshKey)
	d.Set("subnet_id", s.Spec.NodeNicSpec.PrimaryNic.SubnetId)
	d.Set("ecs_group_id", s.Spec.EcsGroupID)
	if s.Spec.BillingMode != 0 {
		d.Set("charging_mode", "prePaid")
	}
	if s.Spec.RunTime != nil {
		d.Set("runtime", s.Spec.RunTime.Name)
	}

	var volumes []map[string]interface{}
	for _, pairObject := range s.Spec.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["hw_passthrough"] = pairObject.HwPassthrough
		volume["extend_params"] = pairObject.ExtendParam
		volume["extend_param"] = ""
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving dataVolumes to state for HuaweiCloud Node (%s): %s", d.Id(), err)
	}

	rootVolume := []map[string]interface{}{
		{
			"size":           s.Spec.RootVolume.Size,
			"volumetype":     s.Spec.RootVolume.VolumeType,
			"hw_passthrough": s.Spec.RootVolume.HwPassthrough,
			"extend_params":  s.Spec.RootVolume.ExtendParam,
			"extend_param":   "",
		},
	}
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving root Volume to state for HuaweiCloud Node (%s): %s", d.Id(), err)
	}

	// set computed attributes
	serverId := s.Status.ServerID
	d.Set("server_id", serverId)
	d.Set("private_ip", s.Status.PrivateIP)
	d.Set("public_ip", s.Status.PublicIP)
	d.Set("status", s.Status.Phase)

	// fetch tags from ECS instance
	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	if resourceTags, err := tags.Get(computeClient, "cloudservers", serverId).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmtp.Errorf("Error saving tags to state for CCE Node (%s): %s", serverId, err)
		}
	} else {
		logp.Printf("[WARN] Error fetching tags of CCE Node (%s): %s", serverId, err)
	}

	return nil
}

func resourceCCENodeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	if d.HasChange("name") {
		var updateOpts nodes.UpdateOpts
		updateOpts.Metadata.Name = d.Get("name").(string)

		clusterid := d.Get("cluster_id").(string)
		_, err = nodes.Update(nodeClient, clusterid, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud cce node: %s", err)
		}
	}

	//update tags
	if d.HasChange("tags") {
		computeClient, err := config.ComputeV1Client(GetRegion(d, config))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		serverId := d.Get("server_id").(string)
		tagErr := utils.UpdateResourceTags(computeClient, d, "cloudservers", serverId)
		if tagErr != nil {
			return fmtp.Errorf("Error updating tags of cce node %s: %s", d.Id(), tagErr)
		}
	}

	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	clusterid := d.Get("cluster_id").(string)
	// remove node without deleting ecs
	if d.Get("keep_ecs").(bool) {
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
	} else {
		// for prePaid node, firstly, we should unsubscribe the ecs server, and then delete it
		if d.Get("charging_mode").(string) == "prePaid" || d.Get("billing_mode").(int) == 2 {
			serverID := d.Get("server_id").(string)
			publicIP := d.Get("public_ip").(string)

			resourceIDs := make([]string, 0, 2)
			computeClient, err := config.ComputeV2Client(GetRegion(d, config))
			if err != nil {
				return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
			}

			// check whether the ecs server of the perPaid exists before unsubscribe it
			// because resource could not be found cannot be unsubscribed
			if serverID != "" {
				server, err := servers.Get(computeClient, serverID).Extract()

				if err != nil {
					if _, ok := err.(golangsdk.ErrDefault404); !ok {
						return fmtp.Errorf("Error retrieving HuaweiCloud ecs intance: %s", err)
					}
				} else {
					if server.Status != "DELETED" && server.Status != "SOFT_DELETED" {
						resourceIDs = append(resourceIDs, serverID)
					}
				}
			}

			// unsubscribe the eip if necessary
			if _, ok := d.GetOk("iptype"); ok && publicIP != "" {
				eipClient, err := config.NetworkingV1Client(GetRegion(d, config))
				if err != nil {
					return fmtp.Errorf("Error creating networking client: %s", err)
				}

				if eipID, err := getEipIDbyAddress(eipClient, publicIP); err == nil {
					resourceIDs = append(resourceIDs, eipID)
				} else {
					logp.Printf("[WARN] Error fetching EIP ID of CCE Node (%s): %s", d.Id(), err)
				}
			}

			if len(resourceIDs) > 0 {
				if err := UnsubscribePrePaidResource(d, config, resourceIDs); err != nil {
					return fmtp.Errorf("Error unsubscribing HuaweiCloud CCE node: %s", err)
				}
			}

			// wait for the ecs server of the prePaid node to be deleted
			pending := []string{"ACTIVE", "SHUTOFF"}
			target := []string{"DELETED", "SOFT_DELETED"}
			deleteTimeout := d.Timeout(schema.TimeoutDelete)
			if err := waitForServerTargetState(computeClient, serverID, pending, target, deleteTimeout); err != nil {
				return fmtp.Errorf("State waiting timeout: %s", err)
			}

			nodes.Delete(nodeClient, clusterid, d.Id())
		} else {
			err = nodes.Delete(nodeClient, clusterid, d.Id()).ExtractErr()
			if err != nil {
				return fmtp.Errorf("Error deleting HuaweiCloud CCE node: %s", err)
			}
		}
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

func getResourceIDFromJob(client *golangsdk.ServiceClient, jobID, jobType, subJobType string) (string, error) {
	// prePaid: waiting for the job to become running
	stateJob := &resource.StateChangeConf{
		Pending:      []string{"Initializing"},
		Target:       []string{"Running", "Success"},
		Refresh:      waitForJobStatus(client, jobID),
		Timeout:      5 * time.Minute,
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}

	v, err := stateJob.WaitForState()
	if err != nil {
		return "", fmtp.Errorf("Error waiting for job (%s) to become running: %s", jobID, err)
	}

	job := v.(*nodes.Job)
	if len(job.Spec.SubJobs) == 0 {
		return "", fmtp.Errorf("Error fetching sub jobs from %s", jobID)
	}

	var subJobID string
	var refreshJob bool
	for _, s := range job.Spec.SubJobs {
		// postPaid: should get details of sub job ID
		if s.Spec.Type == jobType {
			subJobID = s.Metadata.ID
			refreshJob = true
			break
		}
	}

	if refreshJob {
		job, err = nodes.GetJobDetails(client, subJobID).ExtractJob()
		if err != nil {
			return "", fmtp.Errorf("Error fetching sub Job %s: %s", subJobID, err)
		}
	}

	var nodeid string
	for _, s := range job.Spec.SubJobs {
		if s.Spec.Type == subJobType {
			nodeid = s.Spec.ResourceID
			break
		}
	}
	if nodeid == "" {
		return "", fmtp.Errorf("Error fetching %s Job resource id", subJobType)
	}
	return nodeid, nil
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
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud CCE Node %s", nodeId)

		r, err := nodes.Get(cceClient, clusterId, nodeId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud CCE Node %s", nodeId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		return r, r.Status.Phase, nil
	}
}

func waitForClusterAvailable(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[INFO] Waiting for CCE Cluster %s to be available", clusterId)
		n, err := clusters.Get(cceClient, clusterId).Extract()

		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func waitForJobStatus(cceClient *golangsdk.ServiceClient, jobID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		job, err := nodes.GetJobDetails(cceClient, jobID).ExtractJob()
		if err != nil {
			return nil, "", err
		}

		return job, job.Status.Phase, nil
	}
}

func recursiveCreate(cceClient *golangsdk.ServiceClient, opts nodes.CreateOptsBuilder, ClusterID string, errCode int) (*nodes.Nodes, string) {
	if errCode == 403 {
		stateCluster := &resource.StateChangeConf{
			Target:       []string{"Available"},
			Refresh:      waitForClusterAvailable(cceClient, ClusterID),
			Timeout:      15 * time.Minute,
			Delay:        20 * time.Second,
			PollInterval: 10 * time.Second,
		}
		_, stateErr := stateCluster.WaitForState()
		if stateErr != nil {
			logp.Printf("[INFO] Cluster Unavailable %s.\n", stateErr)
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

func installScriptHashSum(script string) string {
	// Check whether the preinstall/postinstall is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(script)
	if base64DecodeError != nil {
		v = []byte(script)
	}

	hash := sha1.Sum(v)
	return hex.EncodeToString(hash[:])
}

func installScriptEncode(script string) string {
	if _, err := base64.StdEncoding.DecodeString(script); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(script))
	}
	return script
}

func getEipIDbyAddress(client *golangsdk.ServiceClient, address string) (string, error) {
	listOpts := &eips.ListOpts{
		PublicIp: address,
	}
	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return "", fmtp.Errorf("Unable to retrieve eips: %s ", err)
	}
	if len(allEips) != 1 {
		return "", fmtp.Errorf("queried none or more results")
	}

	return allEips[0].ID, nil
}

func resourceCCENodeV3Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("Invalid format specified for CCE Node. Format must be <cluster id>/<node id>")
		return nil, err
	}

	clusterID := parts[0]
	nodeID := parts[1]

	d.SetId(nodeID)
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}
