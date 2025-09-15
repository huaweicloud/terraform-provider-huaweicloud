package swrenterprise

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/{project_id}/instances
func DataSourceSwrEnterpriseInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance status.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID of the instance.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the instances.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the instance.`,
						},
						"spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the specification of the instance.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the VPC ID .`,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the subnet ID.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the enterprise project ID.`,
						},
						"obs_bucket_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the OBS bucket name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance version.`,
						},
						"charge_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the charge mode of instance.`,
						},
						"access_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the access address of instance.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last update time.`,
						},
						"expires_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the expired time.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance status.`,
						},
						"user_def_obs": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user specifies the OBS bucket.`,
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the VPC name.`,
						},
						"vpc_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the range of available subnets for the VPC.`,
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the subnet name.`,
						},
						"subnet_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the range of available subnets for the subnet.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSwrEnterpriseInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}
	listInstancesHttpUrl := "v2/{project_id}/instances"
	listInstancesPath := client.Endpoint + listInstancesHttpUrl
	listInstancesPath = strings.ReplaceAll(listInstancesPath, "{project_id}", client.ProjectID)
	listInstancesPath += buildSwrEnterpriseInstancesQueryParams(d)
	listInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listInstancesPath + fmt.Sprintf("&offset=%v", offset)
		listInstancesResp, err := client.Request("GET", currentPath, &listInstancesOpt)
		if err != nil {
			return diag.Errorf("error querying SWR instances: %s", err)
		}
		listInstancesRespBody, err := utils.FlattenResponse(listInstancesResp)
		if err != nil {
			return diag.Errorf("error flattening SWR instances response: %s", err)
		}

		instances := utils.PathSearch("instances", listInstancesRespBody, make([]interface{}, 0)).([]interface{})
		if len(instances) == 0 {
			break
		}
		for _, instance := range instances {
			results = append(results, map[string]interface{}{
				"id":                    utils.PathSearch("id", instance, nil),
				"name":                  utils.PathSearch("name", instance, nil),
				"spec":                  utils.PathSearch("spec", instance, nil),
				"vpc_id":                utils.PathSearch("vpc_id", instance, nil),
				"subnet_id":             utils.PathSearch("subnet_id", instance, nil),
				"enterprise_project_id": utils.PathSearch("enterprise_project_id", instance, nil),
				"obs_bucket_name":       utils.PathSearch("obs_bucket_name", instance, nil),
				"description":           utils.PathSearch("description", instance, nil),
				"version":               utils.PathSearch("version", instance, nil),
				"charge_mode":           utils.PathSearch("charge_mode", instance, nil),
				"access_address":        utils.PathSearch("access_address", instance, nil),
				"created_at":            utils.PathSearch("created_at", instance, nil),
				"updated_at":            utils.PathSearch("updated_at", instance, nil),
				"expires_at":            utils.PathSearch("expires_at", instance, nil),
				"status":                utils.PathSearch("status", instance, nil),
				"user_def_obs":          utils.PathSearch("user_def_obs", instance, nil),
				"vpc_name":              utils.PathSearch("vpc_name", instance, nil),
				"vpc_cidr":              utils.PathSearch("vpc_cidr", instance, nil),
				"subnet_name":           utils.PathSearch("subnet_name", instance, nil),
				"subnet_cidr":           utils.PathSearch("subnet_cidr", instance, nil),
			})
		}

		// offset must be the multiple of limit
		offset += 100
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSwrEnterpriseInstancesQueryParams(d *schema.ResourceData) string {
	res := "?limit=100"

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}

	return res
}
