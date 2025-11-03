package aad

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v2/aad/instances
func DataSourceV2Instances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2InstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_access_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pp_support": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pp_enable": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"overseas_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"isp": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildV2InstancesQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("instance_access_type"); ok {
		return fmt.Sprintf("?instance_access_type=%v", v)
	}

	return ""
}

func dataSourceV2InstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/instances"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildV2InstancesQueryParams(d)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving AAD v2 instances: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("items", flattenV2InstancesItems(utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV2InstancesItems(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"instance_id":           utils.PathSearch("instance_id", v, nil),
			"instance_name":         utils.PathSearch("instance_name", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"instance_access_type":  utils.PathSearch("instance_access_type", v, nil),
			"pp_support":            utils.PathSearch("pp_support", v, nil),
			"pp_enable":             utils.PathSearch("pp_enable", v, nil),
			"overseas_type":         utils.PathSearch("overseas_type", v, nil),
			"vips":                  flattenV2InstancesVips(utils.PathSearch("vips", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenV2InstancesVips(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"ip":  utils.PathSearch("ip", v, nil),
			"isp": utils.PathSearch("isp", v, nil),
		})
	}

	return rst
}
