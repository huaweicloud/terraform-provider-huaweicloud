package huaweicloud

import (
	"context"

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

func DataSourceCCENodeV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceNodesV3Read,

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
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
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
	}
}

func dataSourceCceNodesV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud CCE client : %s", err)
	}

	listOpts := nodes.ListOpts{
		Uid:   d.Get("node_id").(string),
		Name:  d.Get("name").(string),
		Phase: d.Get("status").(string),
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("node_id"); ok {
		listOpts.Uid = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Phase = v.(string)
	}

	refinedNodes, err := nodes.List(cceClient, d.Get("cluster_id").(string), listOpts)

	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve Nodes: %s", err)
	}

	if len(refinedNodes) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNodes) > 1 {
		return fmtp.DiagErrorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Node := refinedNodes[0]

	logp.Printf("[DEBUG] Retrieved Nodes using given filter %s: %+v", Node.Metadata.Id, Node)
	d.SetId(Node.Metadata.Id)

	mErr := multierror.Append(nil,
		d.Set("node_id", Node.Metadata.Id),
		d.Set("name", Node.Metadata.Name),
		d.Set("flavor_id", Node.Spec.Flavor),
		d.Set("availability_zone", Node.Spec.Az),
		d.Set("os", Node.Spec.Os),
		d.Set("billing_mode", Node.Spec.BillingMode),
		d.Set("key_pair", Node.Spec.Login.SshKey),
		d.Set("subnet_id", Node.Spec.NodeNicSpec.PrimaryNic.SubnetId),
		d.Set("ecs_group_id", Node.Spec.EcsGroupID),
		d.Set("server_id", Node.Status.ServerID),
		d.Set("public_ip", Node.Status.PublicIP),
		d.Set("private_ip", Node.Status.PrivateIP),
		d.Set("status", Node.Status.Phase),
		d.Set("region", GetRegion(d, config)),
	)

	var volumes []map[string]interface{}
	for _, pairObject := range Node.Spec.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_params"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	mErr = multierror.Append(mErr, d.Set("data_volumes", volumes))

	rootVolume := []map[string]interface{}{
		{
			"size":          Node.Spec.RootVolume.Size,
			"volumetype":    Node.Spec.RootVolume.VolumeType,
			"extend_params": Node.Spec.RootVolume.ExtendParam,
		},
	}
	mErr = multierror.Append(mErr, d.Set("root_volume", rootVolume))

	// fetch tags from ECS instance
	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud compute client: %s", err)
	}

	serverId := Node.Status.ServerID

	if resourceTags, err := tags.Get(computeClient, "cloudservers", serverId).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagmap))
	} else {
		logp.Printf("[WARN] Error fetching tags of CCE Node (%s): %s", serverId, err)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting node fields: %s", err)
	}

	return nil
}
