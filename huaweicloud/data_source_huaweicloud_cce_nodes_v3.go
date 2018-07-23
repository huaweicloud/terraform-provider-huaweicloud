package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
	"log"
)

//Creates schema for data source
func dataSourceCceNodesV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCceNodesV3Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"az": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"share_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"extend_param": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_volumes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"billing_mode": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"spec_extend_param": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceCceNodesV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Unable to create HuaweiCloud CCE client : %s", err)
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
		return fmt.Errorf("Unable to retrieve Nodes: %s", err)
	}

	if len(refinedNodes) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNodes) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Node := refinedNodes[0]

	var v []map[string]interface{}
	for _, volume := range Node.Spec.DataVolumes {

		mapping := map[string]interface{}{
			"disk_size":   volume.Size,
			"volume_type": volume.VolumeType,
		}
		v = append(v, mapping)
	}

	pids := Node.Spec.PublicIP.Ids
	PublicIDs := make([]string, len(pids))
	for i, val := range pids {
		PublicIDs[i] = val
	}
	log.Printf("[DEBUG] Retrieved Nodes using given filter %s: %+v", Node.Metadata.Id, Node)
	d.SetId(Node.Metadata.Id)
	d.Set("node_id", Node.Metadata.Id)
	d.Set("name", Node.Metadata.Name)
	d.Set("flavor", Node.Spec.Flavor)
	d.Set("az", Node.Spec.Az)
	d.Set("billing_mode", Node.Spec.BillingMode)
	d.Set("status", Node.Status.Phase)
	d.Set("data_volumes", v)
	d.Set("disk_size", Node.Spec.RootVolume.Size)
	d.Set("volume_type", Node.Spec.RootVolume.VolumeType)
	d.Set("extend_param", Node.Spec.RootVolume.ExtendParam)
	d.Set("sshkey", Node.Spec.Login.SshKey)
	d.Set("charge_mode", Node.Spec.PublicIP.Eip.Bandwidth.ChargeMode)
	d.Set("bandwidth_size", Node.Spec.PublicIP.Eip.Bandwidth.Size)
	d.Set("share_type", Node.Spec.PublicIP.Eip.Bandwidth.ShareType)
	d.Set("ip_type", Node.Spec.PublicIP.Eip.IpType)
	d.Set("server_id", Node.Status.ServerID)
	d.Set("public_ip", Node.Status.PublicIP)
	d.Set("private_ip", Node.Status.PrivateIP)
	d.Set("spec_extend_param", Node.Spec.ExtendParam)
	d.Set("eip_count", Node.Spec.PublicIP.Count)
	d.Set("eip_ids", PublicIDs)

	return nil
}
