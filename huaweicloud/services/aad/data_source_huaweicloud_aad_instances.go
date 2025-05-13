package aad

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET v1/aad/instances
func DataSourceAADInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAADInstancesRead,

		Schema: map[string]*schema.Schema{
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        instanceInfo(),
				Description: `The list of the AAD instances.`,
			},
		},
	}
}

func instanceInfo() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AAD instance ID.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the AAD instance.",
			},
			"ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        instanceIpInfo(),
				Description: `The list of the AAD instance IPs.`,
			},
			"expire_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The expiration time of the AAD instance.",
			},
			"service_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The service bandwidth of the AAD instance.",
			},
			"instance_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The AAD instance status.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise project ID of the AAD instance.",
			},
			"overseas_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The AAD instance type, `0`-mainland China, `1`-overseas.",
			},
		},
	}
	return &sc
}

func instanceIpInfo() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP ID of the AAD instance.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP of the AAD instance.",
			},
			"basic_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The basic bandwidth of the AAD instance.",
			},
			"elastic_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The elastic bandwidth of the AAD instance.",
			},
			"ip_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The IP status of the AAD instance.",
			},
		},
	}
	return &sc
}

func dataSourceAADInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)

	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	httpUrl := "v1/aad/instances"
	listPath := client.Endpoint + httpUrl
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving AAD instances:%s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("items", flattenGetAADInstancesResponseBodyItems(utils.PathSearch("items", respBody, nil))),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error retrieving AAD instances: %s", mErr)
	}

	return nil
}

func flattenGetAADInstancesResponseBodyItems(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := resp.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"instance_id":           utils.PathSearch("instance_id", v, nil),
			"instance_name":         utils.PathSearch("instance_name", v, nil),
			"ips":                   flattenGetAADInstancesResponseBodyIps(utils.PathSearch("ips", v, nil)),
			"expire_time":           utils.PathSearch("expire_time", v, nil),
			"service_bandwidth":     utils.PathSearch("service_bandwidth", v, nil),
			"instance_status":       utils.PathSearch("instance_status", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"overseas_type":         utils.PathSearch("overseas_type", v, nil),
		})
	}
	return rst
}

func flattenGetAADInstancesResponseBodyIps(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := resp.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"ip_id":             utils.PathSearch("ip_id", v, nil),
			"ip":                utils.PathSearch("ip", v, nil),
			"basic_bandwidth":   utils.PathSearch("basic_bandwidth", v, nil),
			"elastic_bandwidth": utils.PathSearch("elastic_bandwidth", v, nil),
			"ip_status":         utils.PathSearch("ip_status", v, nil),
		})
	}
	return rst
}
