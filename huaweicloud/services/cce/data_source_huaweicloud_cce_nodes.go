package cce

import (
	"context"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCCENodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceNodesRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"key_pair": {
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
						"billing_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCceNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud CCE client : %s", err)
	}

	listOpts := nodes.ListOpts{
		Uid:   d.Get("node_id").(string),
		Name:  d.Get("name").(string),
		Phase: d.Get("status").(string),
	}

	refinedNodes, err := nodes.List(cceClient, d.Get("cluster_id").(string), listOpts)

	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve Nodes: %s", err)
	}

	ids := make([]string, 0, len(refinedNodes))
	nodesToSet := make([]map[string]interface{}, 0, len(refinedNodes))

	for _, v := range refinedNodes {
		logp.Printf("[DEBUG] Retrieved Nodes using given filter %s: %+v", v.Metadata.Id, v)
		ids = append(ids, v.Metadata.Id)
		node := map[string]interface{}{
			"id":                v.Metadata.Id,
			"name":              v.Metadata.Name,
			"flavor_id":         v.Spec.Flavor,
			"availability_zone": v.Spec.Az,
			"os":                v.Spec.Os,
			"billing_mode":      v.Spec.BillingMode,
			"key_pair":          v.Spec.Login.SshKey,
			"subnet_id":         v.Spec.NodeNicSpec.PrimaryNic.SubnetId,
			"ecs_group_id":      v.Spec.EcsGroupID,
			"server_id":         v.Status.ServerID,
			"public_ip":         v.Status.PublicIP,
			"private_ip":        v.Status.PrivateIP,
			"status":            v.Status.Phase,
		}

		var volumes []map[string]interface{}
		for _, pairObject := range v.Spec.DataVolumes {
			volume := make(map[string]interface{})
			volume["size"] = pairObject.Size
			volume["volumetype"] = pairObject.VolumeType
			volume["extend_params"] = pairObject.ExtendParam
			volumes = append(volumes, volume)
		}
		node["data_volumes"] = volumes

		rootVolume := []map[string]interface{}{
			{
				"size":          v.Spec.RootVolume.Size,
				"volumetype":    v.Spec.RootVolume.VolumeType,
				"extend_params": v.Spec.RootVolume.ExtendParam,
			},
		}
		node["root_volume"] = rootVolume

		// fetch tags from ECS instance
		computeClient, err := config.ComputeV1Client(config.GetRegion(d))
		if err != nil {
			return fmtp.DiagErrorf("Error creating HuaweiCloud compute client: %s", err)
		}

		serverId := v.Status.ServerID

		if resourceTags, err := tags.Get(computeClient, "cloudservers", serverId).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			node["tags"] = tagmap
		} else {
			logp.Printf("[WARN] Error fetching tags of CCE Node (%s): %s", serverId, err)
		}

		nodesToSet = append(nodesToSet, node)
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("nodes", nodesToSet),
		d.Set("ids", ids),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting cce nodes fields: %s", err)
	}

	return nil
}
