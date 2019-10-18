package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func dataSourceCCEClusterV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCCEClusterV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"highway_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_network_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					}},
			},
		},
	}
}

func dataSourceCCEClusterV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Unable to create HuaweiCloud CCE client : %s", err)
	}

	listOpts := clusters.ListOpts{
		ID:    d.Get("id").(string),
		Name:  d.Get("name").(string),
		Type:  d.Get("cluster_type").(string),
		Phase: d.Get("status").(string),
		VpcID: d.Get("vpc_id").(string),
	}

	refinedClusters, err := clusters.List(cceClient, listOpts)
	log.Printf("[DEBUG] Value of allClusters: %#v", refinedClusters)
	if err != nil {
		return fmt.Errorf("Unable to retrieve clusters: %s", err)
	}

	if len(refinedClusters) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedClusters) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Cluster := refinedClusters[0]

	log.Printf("[DEBUG] Retrieved Clusters using given filter %s: %+v", Cluster.Metadata.Id, Cluster)
	var v []map[string]interface{}
	for _, endpoint := range Cluster.Status.Endpoints {

		mapping := map[string]interface{}{
			"url":  endpoint.Url,
			"type": endpoint.Type,
		}
		v = append(v, mapping)
	}

	d.SetId(Cluster.Metadata.Id)

	d.Set("name", Cluster.Metadata.Name)
	d.Set("flavor_id", Cluster.Spec.Flavor)
	d.Set("description", Cluster.Spec.Description)
	d.Set("cluster_version", Cluster.Spec.Version)
	d.Set("cluster_type", Cluster.Spec.Type)
	d.Set("billing_mode", Cluster.Spec.BillingMode)
	d.Set("vpc_id", Cluster.Spec.HostNetwork.VpcId)
	d.Set("subnet_id", Cluster.Spec.HostNetwork.SubnetId)
	d.Set("highway_subnet_id", Cluster.Spec.HostNetwork.HighwaySubnet)
	d.Set("container_network_cidr", Cluster.Spec.ContainerNetwork.Cidr)
	d.Set("container_network_type", Cluster.Spec.ContainerNetwork.Mode)
	d.Set("status", Cluster.Status.Phase)
	if err := d.Set("endpoints", v); err != nil {
		return err
	}
	d.Set("region", GetRegion(d, config))

	return nil
}
