package huaweicloud

import (
	"strings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodepools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCCENodePoolV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCceNodePoolsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"initial_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"current_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": { //(k8s_tags)
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"extend_params": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
			},
			"max_pods": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"extend_param": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scall_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"min_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scale_down_cooldown_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceCceNodePoolsV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("unable to create HuaweiCloud CCE client : %s", err)
	}

	listOpts := nodepools.ListOpts{
		Uid:   d.Get("node_pool_id").(string),
		Name:  d.Get("name").(string),
		Phase: d.Get("status").(string),
	}

	refinedNodePools, err := nodepools.List(cceClient, d.Get("cluster_id").(string), listOpts)

	if err != nil {
		return fmtp.Errorf("unable to retrieve Node Pools: %s", err)
	}

	if len(refinedNodePools) < 1 {
		return fmtp.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(refinedNodePools) > 1 {
		return fmtp.Errorf("your query returned more than one result, please try a more specific search criteria")
	}

	NodePool := refinedNodePools[0]

	logp.Printf("[DEBUG] Retrieved Node Pools using given filter %s: %+v", NodePool.Metadata.Id, NodePool)
	d.SetId(NodePool.Metadata.Id)

	d.Set("node_pool_id", NodePool.Metadata.Id)
	d.Set("name", NodePool.Metadata.Name)
	d.Set("type", NodePool.Spec.Type)
	d.Set("flavor_id", NodePool.Spec.NodeTemplate.Flavor)
	d.Set("availability_zone", NodePool.Spec.NodeTemplate.Az)
	d.Set("os", NodePool.Spec.NodeTemplate.Os)
	d.Set("key_pair", NodePool.Spec.NodeTemplate.Login.SshKey)
	d.Set("scall_enable", NodePool.Spec.Autoscaling.Enable)
	d.Set("initial_node_count", NodePool.Spec.InitialNodeCount)
	d.Set("current_node_count", NodePool.Status.CurrentNode)
	d.Set("min_node_count", NodePool.Spec.Autoscaling.MinNodeCount)
	d.Set("max_node_count", NodePool.Spec.Autoscaling.MaxNodeCount)
	d.Set("scale_down_cooldown_time", NodePool.Spec.Autoscaling.ScaleDownCooldownTime)
	d.Set("priority", NodePool.Spec.Autoscaling.Priority)
	d.Set("subnet_id", NodePool.Spec.NodeTemplate.NodeNicSpec.PrimaryNic.SubnetId)
	d.Set("status", NodePool.Status.Phase)
	d.Set("region", GetRegion(d, config))

	// set extend_param
	var extendParam = NodePool.Spec.NodeTemplate.ExtendParam
	d.Set("max_pods", extendParam["maxPods"])
	delete(extendParam, "maxPods")
	if len(extendParam) > 0 {
		d.Set("extend_param", extendParam)
	}

	labels := map[string]string{}
	for key, val := range NodePool.Spec.NodeTemplate.K8sTags {
		if strings.Contains(key, "cce.cloud.com") {
			continue
		}
		labels[key] = val
	}
	d.Set("labels", labels)

	volumes := make([]map[string]interface{}, 0, len(NodePool.Spec.NodeTemplate.DataVolumes))
	for _, pairObject := range NodePool.Spec.NodeTemplate.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_params"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving dataVolumes to state for HuaweiCloud Node Pool (%s): %s", d.Id(), err)
	}

	rootVolume := []map[string]interface{}{
		{
			"size":          NodePool.Spec.NodeTemplate.RootVolume.Size,
			"volumetype":    NodePool.Spec.NodeTemplate.RootVolume.VolumeType,
			"extend_params": NodePool.Spec.NodeTemplate.RootVolume.ExtendParam,
		},
	}
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving root Volume to state for HuaweiCloud Node Pool (%s): %s", d.Id(), err)
	}

	tagmap := utils.TagsToMap(NodePool.Spec.NodeTemplate.UserTags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmtp.Errorf("error saving tags to state for CCE Node Pool(%s): %s", d.Id(), err)
	}

	return nil
}
