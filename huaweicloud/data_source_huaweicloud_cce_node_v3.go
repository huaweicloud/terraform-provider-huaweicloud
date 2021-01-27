package huaweicloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
)

func DataSourceCCENodeV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCceNodesV3Read,

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
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extend_param": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
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
			"eip_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"spec_extend_param": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceCceNodesV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.CceV3Client(GetRegion(d, config))
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
	d.Set("flavor_id", Node.Spec.Flavor)
	d.Set("availability_zone", Node.Spec.Az)
	d.Set("billing_mode", Node.Spec.BillingMode)
	d.Set("status", Node.Status.Phase)
	d.Set("data_volumes", v)
	d.Set("disk_size", Node.Spec.RootVolume.Size)
	d.Set("volume_type", Node.Spec.RootVolume.VolumeType)
	d.Set("extend_param", Node.Spec.RootVolume.ExtendParam)
	d.Set("key_pair", Node.Spec.Login.SshKey)
	d.Set("charge_mode", Node.Spec.PublicIP.Eip.Bandwidth.ChargeMode)
	d.Set("bandwidth_size", Node.Spec.PublicIP.Eip.Bandwidth.Size)
	d.Set("share_type", Node.Spec.PublicIP.Eip.Bandwidth.ShareType)
	d.Set("ip_type", Node.Spec.PublicIP.Eip.IpType)
	d.Set("server_id", Node.Status.ServerID)
	d.Set("public_ip", Node.Status.PublicIP)
	d.Set("private_ip", Node.Status.PrivateIP)
	if byte, err := json.Marshal(Node.Spec.ExtendParam); err == nil {
		if err = d.Set("spec_extend_param", string(byte)); err != nil {
			return fmt.Errorf("Saving spec extend param ERROR: %s", err)
		}
	} else {
		return fmt.Errorf("Spec extend param translate ERROR: %s", err)
	}
	d.Set("eip_count", Node.Spec.PublicIP.Count)
	d.Set("eip_ids", PublicIDs)

	return nil
}
