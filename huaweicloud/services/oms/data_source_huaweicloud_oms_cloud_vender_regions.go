package oms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OMS GET /v2/{project_id}/objectstorage/data-center
func DataSourceCloudVenderRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudVenderRegionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cloud_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
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

func dataSourceCloudVenderRegionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/objectstorage/data-center"
	)

	client, err := cfg.NewServiceClient("oms", region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the cloud vendor supported regions: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("region_info", flattenCloudVenderRegions(
			utils.PathSearch("[*]", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCloudVenderRegions(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		rst = append(rst, map[string]interface{}{
			"service_name": utils.PathSearch("service_name", v, nil),
			"region_list": flattenObjectstorageRegionList(
				utils.PathSearch("region_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenObjectstorageRegionList(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		rst = append(rst, map[string]interface{}{
			"cloud_type":  utils.PathSearch("cloud_type", v, nil),
			"value":       utils.PathSearch("value", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}

	return rst
}
