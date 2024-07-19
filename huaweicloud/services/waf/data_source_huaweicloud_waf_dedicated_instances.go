package waf

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	instances "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API WAF GET /v1/{project_id}/premium-waf/instance/{instance_id}
// @API WAF GET /v1/{project_id}/premium-waf/instance
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
			"enterprise_project_id": {
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
						"ecs_flavor": {
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
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	instanceId, hasId := d.GetOk("id")
	epsId := d.Get("enterprise_project_id").(string)
	var items []instances.DedicatedInstance
	if hasId {
		instance, err := instances.GetWithEpsId(client, instanceId.(string), epsId)
		if err != nil {
			return diag.Errorf("Your query returned no results. " +
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
			InstanceName:        d.Get("name").(string),
			EnterpriseProjectId: epsId,
		}

		rst, err := instances.ListInstance(client, opts)
		if err != nil {
			return diag.Errorf("error retrieving WAF dedicated instances %s", err)
		}
		items = rst.Items
	}

	if len(items) == 0 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	ids := make([]string, 0, len(items))
	insts := make([]map[string]interface{}, 0, len(items))

	for _, r := range items {
		eng := map[string]interface{}{
			"id":               r.Id,
			"name":             r.InstanceName,
			"available_zone":   r.Zone,
			"cpu_architecture": r.Arch,
			"ecs_flavor":       r.CupFlavor,
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
		insts = append(insts, eng)
		ids = append(ids, r.Id)
	}

	if !hasId {
		d.SetId(hashcode.Strings(ids))
	}
	mErr := multierror.Append(nil,
		d.Set("instances", insts),
		d.Set("region", conf.GetRegion(d)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting WAF dedicated instance fields: %s", mErr)
	}

	return nil
}
