package waf

import (
	"context"
	"strings"

	instances "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_instances"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceWafDedicatedInstancesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafDedicatedInstanceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"run_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"access_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"upgradable": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWafDedicatedInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.WafDedicatedV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client: %s", err)
	}

	instanceId, hasId := d.GetOk("id")
	var items []instances.DedicatedInstance
	if hasId {
		instance, err := instances.GetInstance(client, instanceId.(string))
		if err != nil {
			return fmtp.DiagErrorf("Your query returned no results. " +
				"Please change your search criteria and try again.")
		}
		d.SetId(instanceId.(string))
		items = append(items, *instance)

		if n, ok := d.GetOk("name"); ok {
			// If the instance name does not match name form schema, then clear the items
			if !strings.Contains(strings.ToLower(instance.InstanceName), strings.ToLower(n.(string))) {
				items = []instances.DedicatedInstance{}
			}
		}
	} else {
		// If the instance id is not set, or the name value is not set, the query list can be used.
		opts := instances.ListInstanceOpts{
			InstanceName: d.Get("name").(string),
		}

		rst, err := instances.ListInstance(client, opts)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "Error obtain WAF dedicated instance information.")
		}
		items = rst.Items
	}

	if len(items) == 0 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	ids := make([]string, 0, len(items))
	instances := make([]map[string]interface{}, 0, len(items))

	for _, r := range items {
		eng := map[string]interface{}{
			"id":               r.Id,
			"name":             r.InstanceName,
			"available_zone":   r.Zone,
			"cpu_architecture": r.Arch,
			"cpu_flavor":       r.CupFlavor,
			"vpc_id":           r.VpcId,
			"subnet_id":        r.SubnetId,
			"security_group":   r.SecurityGroupIds,
			"server_id":        r.ServerId,
			"service_ip":       r.ServiceIp,
			"run_status":       r.RunStatus,
			"access_status":    r.AccessStatus,
			"upgradable":       r.Upgradable,
			"group_id":         r.PoolId,
		}
		instances = append(instances, eng)
		ids = append(ids, r.Id)
	}

	if !hasId {
		d.SetId(hashcode.Strings(ids))
	}
	mErr := multierror.Append(nil,
		d.Set("instances", instances),
		d.Set("region", config.GetRegion(d)),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting WAF dedicated instance fields: %s", mErr)
	}

	return nil
}
