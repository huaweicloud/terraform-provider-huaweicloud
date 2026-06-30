package hss

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

// @API HSS GET /v5/{project_id}/common/hosts
func DataSourceCommonHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCommonHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"feature_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_name_upper": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_name_lower": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     commonHostsDataListSchema(),
			},
		},
	}
}

func commonHostsDataListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCommonHostsQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		queryParams = fmt.Sprintf("%s&public_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("feature_name"); ok {
		queryParams = fmt.Sprintf("%s&feature_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}
	if v, ok := d.GetOk("agent_status"); ok {
		queryParams = fmt.Sprintf("%s&agent_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		queryParams = fmt.Sprintf("%s&cluster_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		queryParams = fmt.Sprintf("%s&cluster_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("version_name_upper"); ok {
		queryParams = fmt.Sprintf("%s&version_name_upper=%v", queryParams, v)
	}
	if v, ok := d.GetOk("version_name_lower"); ok {
		queryParams = fmt.Sprintf("%s&version_name_lower=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceCommonHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		result  = make([]interface{}, 0)
		offset  = 0
		httpUrl = "v5/{project_id}/common/hosts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildCommonHostsQueryParams(d)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS common hosts: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenCommonHostsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCommonHostsDataList(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"host_id":      utils.PathSearch("host_id", v, nil),
			"host_name":    utils.PathSearch("host_name", v, nil),
			"public_ip":    utils.PathSearch("public_ip", v, nil),
			"private_ip":   utils.PathSearch("private_ip", v, nil),
			"agent_id":     utils.PathSearch("agent_id", v, nil),
			"os_type":      utils.PathSearch("os_type", v, nil),
			"host_status":  utils.PathSearch("host_status", v, nil),
			"agent_status": utils.PathSearch("agent_status", v, nil),
			"os_name":      utils.PathSearch("os_name", v, nil),
			"os_version":   utils.PathSearch("os_version", v, nil),
			"asset_value":  utils.PathSearch("asset_value", v, nil),
			"cluster_id":   utils.PathSearch("cluster_id", v, nil),
			"cluster_name": utils.PathSearch("cluster_name", v, nil),
			"group_id":     utils.PathSearch("group_id", v, nil),
			"group_name":   utils.PathSearch("group_name", v, nil),
		})
	}

	return rst
}
