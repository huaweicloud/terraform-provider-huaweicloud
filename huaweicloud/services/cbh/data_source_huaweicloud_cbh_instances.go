package cbh

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBH GET /v2/{project_id}/cbs/instance/list
func DataSourceCbhInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCbhInstancesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance name.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of a VPC.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of a subnet.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of a security group.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the specification of the instance.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the current version of the instance image`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        instancesInstanceSchema(),
				Computed:    true,
				Description: `Indicates the list of CBH instance.`,
			},
		},
	}
}

func instancesInstanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the instance.`,
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the elastic IP.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the elastic IP address.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance name.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the private IP address of the instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the instance.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a VPC.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a subnet.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a security group.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the specification of the instance.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the availability zone name.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the current version of the instance image.`,
			},
		},
	}
	return &sc
}

func resourceCbhInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                    = meta.(*config.Config)
		region                 = cfg.GetRegion(d)
		mErr                   *multierror.Error
		getCbhInstancesProduct = "cbh"
	)
	client, err := cfg.NewServiceClient(getCbhInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instances", flattenCBHResponseInstances(filterCBHResponseInstances(d, instances))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCBHResponseInstances(instances []interface{}) []interface{} {
	res := make([]interface{}, 0, len(instances))
	for _, v := range instances {
		// When EIP is not configured, the query interface field will return a space string.
		publicIpId := strings.TrimSpace(utils.PathSearch("network.public_id", v, "").(string))
		res = append(res, map[string]interface{}{
			"id":                utils.PathSearch("server_id", v, nil),
			"public_ip_id":      publicIpId,
			"public_ip":         utils.PathSearch("network.public_ip", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"private_ip":        utils.PathSearch("network.private_ip", v, nil),
			"status":            utils.PathSearch("status_info.status", v, nil),
			"vpc_id":            utils.PathSearch("network.vpc_id", v, nil),
			"subnet_id":         utils.PathSearch("network.subnet_id", v, nil),
			"security_group_id": utils.PathSearch("network.security_group_id", v, nil),
			"flavor_id":         utils.PathSearch("resource_info.specification", v, nil),
			"availability_zone": utils.PathSearch("az_info.zone", v, nil),
			"version":           utils.PathSearch("bastion_version", v, nil),
		})
	}
	return res
}

func filterCBHResponseInstances(d *schema.ResourceData, instances []interface{}) []interface{} {
	name := d.Get("name").(string)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	securityGroupId := d.Get("security_group_id").(string)
	flavorId := d.Get("flavor_id").(string)
	version := d.Get("version").(string)
	res := make([]interface{}, 0, len(instances))
	for _, v := range instances {
		nameResp := utils.PathSearch("name", v, "").(string)
		if len(name) > 0 && name != nameResp {
			continue
		}
		vpcIdResp := utils.PathSearch("network.vpc_id", v, "").(string)
		if len(vpcId) > 0 && vpcId != vpcIdResp {
			continue
		}
		subnetIdResp := utils.PathSearch("network.subnet_id", v, "").(string)
		if len(subnetId) > 0 && subnetId != subnetIdResp {
			continue
		}
		securityGroupIdResp := utils.PathSearch("network.security_group_id", v, "").(string)
		if len(securityGroupId) > 0 && securityGroupId != securityGroupIdResp {
			continue
		}
		flavorIdResp := utils.PathSearch("resource_info.specification", v, "").(string)
		if len(flavorId) > 0 && flavorId != flavorIdResp {
			continue
		}
		versionResp := utils.PathSearch("bastion_version", v, "").(string)
		if len(version) > 0 && version != versionResp {
			continue
		}

		res = append(res, v)
	}
	return res
}
