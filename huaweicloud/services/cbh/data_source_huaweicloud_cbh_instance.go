package cbh

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
			"bastion_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the bastion.`,
			},
			"bastion_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the current version of the instance image`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        CbhInstancesInstanceSchema(),
				Computed:    true,
				Description: `Indicates the list of CBH instance.`,
			},
		},
	}
}

func CbhInstancesInstanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"publicip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the elastic IP.`,
			},
			"exp_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the expire time of the instance.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the start time of the instance.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the end time of the instance.`,
			},
			"release_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the release time of the instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance name.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the server id of the instance.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the private ip of the instance.`,
			},
			"task_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task status of the instance.`,
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
			"update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates whether the instance image can be upgraded.`,
			},
			"instance_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the instance.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the resource.`,
			},
			"period": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the duration of tenant purchase.`,
			},
			"bastion_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the bastion.`,
			},
			"alter_permit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates whether the front-end displays the capacity expansion button.`,
			},
			"bastion_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the current version of the instance image.`,
			},
			"new_bastion_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest version of the instance image.`,
			},
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the bastion.`,
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates whether auto renew is enabled.`,
			},
		},
	}
	return &sc
}

func resourceCbhInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getCbhInstances: Query the List of CBH instances
	var (
		getCbhInstancesHttpUrl = "v1/{project_id}/cbs/instance/list"
		getCbhInstancesProduct = "cbh"
	)
	getCbhInstancesClient, err := config.NewServiceClient(getCbhInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CbhInstances Client: %s", err)
	}

	getCbhInstancesPath := getCbhInstancesClient.Endpoint + getCbhInstancesHttpUrl
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", getCbhInstancesClient.ProjectID)

	getCbhInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCbhInstancesResp, err := getCbhInstancesClient.Request("GET", getCbhInstancesPath, &getCbhInstancesOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CbhInstances")
	}

	getCbhInstancesRespBody, err := utils.FlattenResponse(getCbhInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	instances := utils.PathSearch("instance", getCbhInstancesRespBody, make([]interface{}, 0)).([]interface{})

	name := d.Get("name").(string)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	securityGroupId := d.Get("security_group_id").(string)
	flavorId := d.Get("flavor_id").(string)
	bastionType := d.Get("bastion_type").(string)
	bastionVersion := d.Get("bastion_version").(string)

	res := make([]interface{}, 0)
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if len(name) > 0 && instance["name"].(string) != name {
			continue
		}
		if len(vpcId) > 0 && instance["vpc_id"].(string) != vpcId {
			continue
		}
		if len(subnetId) > 0 && instance["subnet_id"].(string) != subnetId {
			continue
		}
		if len(securityGroupId) > 0 && instance["security_group_id"].(string) != securityGroupId {
			continue
		}
		if len(flavorId) > 0 && instance["flavor_id"].(string) != flavorId {
			continue
		}
		if len(bastionType) > 0 && instance["bastion_type"].(string) != bastionType {
			continue
		}
		if len(bastionVersion) > 0 && instance["bastion_version"].(string) != bastionVersion {
			continue
		}
		res = append(res, v)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instances", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
