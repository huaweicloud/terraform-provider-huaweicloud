package cce

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{clusterid}/nodes
// @API ECS GET /v1/{project_id}/cloudservers/{id}/tags
func DataSourceNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodeRead,

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
			"hostname_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNodeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("unable to create CCE client : %s", err)
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
		return diag.Errorf("unable to retrieve Nodes: %s", err)
	}

	if len(refinedNodes) < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	if len(refinedNodes) > 1 {
		return diag.Errorf("your query returned more than one result." +
			" please try a more specific search criteria")
	}

	node := refinedNodes[0]

	log.Printf("[DEBUG] Retrieved Nodes using given filter %s: %+v", node.Metadata.Id, node)
	d.SetId(node.Metadata.Id)

	mErr := multierror.Append(nil,
		d.Set("node_id", node.Metadata.Id),
		d.Set("name", node.Metadata.Name),
		d.Set("flavor_id", node.Spec.Flavor),
		d.Set("availability_zone", node.Spec.Az),
		d.Set("os", node.Spec.Os),
		d.Set("billing_mode", node.Spec.BillingMode),
		d.Set("key_pair", node.Spec.Login.SshKey),
		d.Set("subnet_id", node.Spec.NodeNicSpec.PrimaryNic.SubnetId),
		d.Set("ecs_group_id", node.Spec.EcsGroupID),
		d.Set("server_id", node.Status.ServerID),
		d.Set("public_ip", node.Status.PublicIP),
		d.Set("private_ip", node.Status.PrivateIP),
		d.Set("status", node.Status.Phase),
		d.Set("region", config.GetRegion(d)),
		d.Set("hostname_config", flattenResourceNodeHostnameConfig(node.Spec.HostnameConfig)),
		d.Set("enterprise_project_id", node.Spec.ServerEnterpriseProjectID),
	)

	var volumes []map[string]interface{}
	for _, pairObject := range node.Spec.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_params"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	mErr = multierror.Append(mErr, d.Set("data_volumes", volumes))

	rootVolume := []map[string]interface{}{
		{
			"size":          node.Spec.RootVolume.Size,
			"volumetype":    node.Spec.RootVolume.VolumeType,
			"extend_params": node.Spec.RootVolume.ExtendParam,
		},
	}
	mErr = multierror.Append(mErr, d.Set("root_volume", rootVolume))

	// fetch tags from ECS instance
	computeClient, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	serverId := node.Status.ServerID

	if resourceTags, err := tags.Get(computeClient, "cloudservers", serverId).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagmap))
	} else {
		log.Printf("[WARN] Error fetching tags of CCE Node (%s): %s", serverId, err)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting node fields: %s", err)
	}

	return nil
}
