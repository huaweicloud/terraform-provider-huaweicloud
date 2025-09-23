package waf

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	instances "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API WAF GET /v1/{project_id}/premium-waf/instance
func DataSourceWafDedicatedInstances() *schema.Resource {
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

						// Deprecated; Reasons for abandonment are as follows:
						// `group_id`: Legacy fields are no longer supported.
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `schema: Deprecated;`,
						},
					},
				},
			},
		},
	}
}

func filterAndFlattenInstances(instanceResp *instances.DedicatedInstanceList, d *schema.ResourceData) []interface{} {
	if instanceResp == nil {
		return nil
	}

	var rst []interface{}
	for _, v := range instanceResp.Items {
		// In order to be compatible with historical logic, the following code is retained.
		if instanceID, hasID := d.GetOk("id"); hasID && instanceID != v.Id {
			continue
		}

		ins := map[string]interface{}{
			"id":               v.Id,
			"name":             v.InstanceName,
			"available_zone":   v.Zone,
			"cpu_architecture": v.Arch,
			"ecs_flavor":       v.CupFlavor,
			"vpc_id":           v.VpcId,
			"subnet_id":        v.SubnetId,
			"security_group":   v.SecurityGroupIds,
			"server_id":        v.ServerId,
			"service_ip":       v.ServiceIp,
			"run_status":       v.RunStatus,
			"access_status":    v.AccessStatus,
			"upgradable":       v.Upgradable,
			"group_id":         v.PoolId,
		}
		rst = append(rst, ins)
	}
	return rst
}

// Keep the historical code logic.
// If an ID is configured, the value of the ID will be used first as the ID of the datasource.
func generateDatasourceID(d *schema.ResourceData) (string, error) {
	instanceID, hasID := d.GetOk("id")
	if hasID {
		return instanceID.(string), nil
	}

	return uuid.GenerateUUID()
}

func dataSourceWafDedicatedInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WafDedicatedV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	opts := instances.ListInstanceOpts{
		InstanceName:        d.Get("name").(string),
		EnterpriseProjectId: conf.GetEnterpriseProjectID(d),
	}

	instanceResp, err := instances.ListInstance(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving WAF dedicated instances %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("instances", filterAndFlattenInstances(instanceResp, d)),
		d.Set("region", conf.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting WAF dedicated instance fields: %s", err)
	}

	dataSourceId, err := generateDatasourceID(d)
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	return nil
}
