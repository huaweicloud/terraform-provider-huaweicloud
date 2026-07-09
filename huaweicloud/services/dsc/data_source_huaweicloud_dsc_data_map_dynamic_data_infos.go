package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// Currently, this resource only supports attribute fields in the API documentation that have actual return values.
// Other attribute fields that do not actually return or are inconsistent with the API documentation are temporarily
// not supported.

// @API DSC GET /v2/{project_id}/data-map/data-infos/dbss-list
func DataSourceDataMapDynamicDataInfos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataMapDynamicDataInfosRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_dbss_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapVpcDynamicGroupSchema(),
			},
		},
	}
}

func dataMapVpcDynamicGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dbss": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapDbssInfoSchema(),
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataMapDbssInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dbss_instance_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapDbssInstanceInfoSchema(),
			},
			"dbss_rds_database_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapRdsDatabaseSchema(),
			},
		},
	}
}

func dataMapDbssInstanceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataMapRdsDatabaseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"configured": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDataMapDynamicDataInfosQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("vpc_id"); ok {
		return fmt.Sprintf("?vpc_id=%v", v)
	}

	return ""
}

func dataSourceDataMapDynamicDataInfosRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/data-map/data-infos/dbss-list"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDataMapDynamicDataInfosQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC data map dynamic data infos: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_dbss_list", flattenDataMapVpcDynamicGroupList(
			utils.PathSearch("vpc_dbss_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataMapVpcDynamicGroupList(vpcDbssList []interface{}) []interface{} {
	if len(vpcDbssList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(vpcDbssList))
	for _, v := range vpcDbssList {
		rst = append(rst, map[string]interface{}{
			"dbss": flattenDataMapDbssInfoList(
				utils.PathSearch("dbss", v, make([]interface{}, 0)).([]interface{})),
			"total":    utils.PathSearch("total", v, nil),
			"vpc_id":   utils.PathSearch("vpc_id", v, nil),
			"vpc_name": utils.PathSearch("vpc_name", v, nil),
		})
	}

	return rst
}

func flattenDataMapDbssInfoList(dbssList []interface{}) []interface{} {
	if len(dbssList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dbssList))
	for _, v := range dbssList {
		rst = append(rst, map[string]interface{}{
			"dbss_instance_info_list": flattenDataMapDbssInstanceInfoList(
				utils.PathSearch("dbss_instance_info_list", v, make([]interface{}, 0)).([]interface{})),
			"dbss_rds_database_list": flattenDataMapRdsDatabaseList(
				utils.PathSearch("dbss_rds_database_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDataMapDbssInstanceInfoList(instanceInfoList []interface{}) []interface{} {
	if len(instanceInfoList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(instanceInfoList))
	for _, v := range instanceInfoList {
		rst = append(rst, map[string]interface{}{
			"instance_id":   utils.PathSearch("instance_id", v, nil),
			"instance_name": utils.PathSearch("instance_name", v, nil),
		})
	}

	return rst
}

func flattenDataMapRdsDatabaseList(rdsDatabaseList []interface{}) []interface{} {
	if len(rdsDatabaseList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rdsDatabaseList))
	for _, v := range rdsDatabaseList {
		rst = append(rst, map[string]interface{}{
			"configured": utils.PathSearch("configured", v, nil),
			"db_name":    utils.PathSearch("db_name", v, nil),
			"id":         utils.PathSearch("id", v, nil),
			"type":       utils.PathSearch("type", v, nil),
		})
	}

	return rst
}
