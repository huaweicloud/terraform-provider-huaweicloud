package cbh

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CBH GET /v1/{project_id}/cbs/instance/list
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
				Description: `Indicates the private ip of the instance.`,
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCbhInstances: Query the List of CBH instances
	var (
		getCbhInstancesProduct = "cbh"
	)
	getCbhInstancesClient, err := cfg.NewServiceClient(getCbhInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CbhInstances Client: %s", err)
	}

	instances, err := getInstanceList(getCbhInstancesClient)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	securityGroupId := d.Get("security_group_id").(string)
	flavorId := d.Get("flavor_id").(string)
	version := d.Get("version").(string)

	res := make([]interface{}, 0)
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if len(name) > 0 && instance["name"].(string) != name {
			continue
		}
		if len(vpcId) > 0 && instance["vpcId"].(string) != vpcId {
			continue
		}
		if len(subnetId) > 0 && instance["subnetId"].(string) != subnetId {
			continue
		}
		if len(securityGroupId) > 0 && instance["securityGroupId"].(string) != securityGroupId {
			continue
		}
		if len(flavorId) > 0 && instance["specification"].(string) != flavorId {
			continue
		}
		if len(version) > 0 && instance["bastionVersion"].(string) != version {
			continue
		}
		publicIpId := instance["publicId"]
		var publicIp string
		if publicIpId != nil && strings.TrimSpace(publicIpId.(string)) != "" {
			publicIp, err = getPublicAddressById(d, cfg, strings.TrimSpace(publicIpId.(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		res = append(res, map[string]interface{}{
			"id":                instance["instanceId"],
			"public_ip_id":      publicIpId,
			"public_ip":         publicIp,
			"name":              instance["name"],
			"private_ip":        instance["privateIp"],
			"status":            instance["status"],
			"vpc_id":            instance["vpcId"],
			"subnet_id":         instance["subnetId"],
			"security_group_id": instance["securityGroupId"],
			"flavor_id":         instance["specification"],
			"availability_zone": instance["zone"],
			"version":           instance["bastionVersion"],
		})
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instances", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
