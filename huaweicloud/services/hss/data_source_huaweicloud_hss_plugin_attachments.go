package hss

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

// @API HSS GET /v5/{project_id}/plugins/attachments
func DataSourcePluginAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluginAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"plugin_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host_status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"agent_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_arch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourcePluginAttachmentsDataListSchema(),
			},
		},
	}
}

func dataSourcePluginAttachmentsDataListSchema() *schema.Resource {
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
			"host_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
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
			"agent_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_arch": {
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
			"plugin_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plugin_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"install_progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remaining_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildPluginAttachmentsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?plugin_code=%s&limit=200", d.Get("plugin_code").(string))
	hostIdList := d.Get("host_ids").([]interface{})
	hostStatusList := d.Get("host_status").([]interface{})

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	for _, key := range []string{
		"plugin_version",
		"plugin_status",
		"host_name",
		"agent_status",
		"os_type",
		"os_arch",
		"host_type",
	} {
		if v, ok := d.GetOk(key); ok {
			queryParams = fmt.Sprintf("%s&%s=%v", queryParams, key, v)
		}
	}

	if len(hostIdList) > 0 {
		for _, v := range hostIdList {
			queryParams = fmt.Sprintf("%s&host_ids=%v", queryParams, v)
		}
	}
	if len(hostStatusList) > 0 {
		for _, v := range hostStatusList {
			queryParams = fmt.Sprintf("%s&host_status=%v", queryParams, v)
		}
	}

	return queryParams
}

func dataSourcePluginAttachmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/plugins/attachments"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildPluginAttachmentsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS plugin attachments: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		dataListResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataListResp) == 0 {
			break
		}

		result = append(result, dataListResp...)
		offset += len(dataListResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenPluginAttachmentsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPluginAttachmentsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_id":          utils.PathSearch("host_id", v, nil),
			"host_name":        utils.PathSearch("host_name", v, nil),
			"host_type":        utils.PathSearch("host_type", v, nil),
			"private_ip":       utils.PathSearch("private_ip", v, nil),
			"public_ip":        utils.PathSearch("public_ip", v, nil),
			"host_status":      utils.PathSearch("host_status", v, nil),
			"agent_status":     utils.PathSearch("agent_status", v, nil),
			"agent_version":    utils.PathSearch("agent_version", v, nil),
			"asset_value":      utils.PathSearch("asset_value", v, nil),
			"os_type":          utils.PathSearch("os_type", v, nil),
			"os_arch":          utils.PathSearch("os_arch", v, nil),
			"os_name":          utils.PathSearch("os_name", v, nil),
			"os_version":       utils.PathSearch("os_version", v, nil),
			"plugin_status":    utils.PathSearch("plugin_status", v, nil),
			"plugin_version":   utils.PathSearch("plugin_version", v, nil),
			"status_detail":    utils.PathSearch("status_detail", v, nil),
			"install_progress": utils.PathSearch("install_progress", v, nil),
			"remaining_time":   utils.PathSearch("remaining_time", v, nil),
			"protect_status":   utils.PathSearch("protect_status", v, nil),
		})
	}

	return rst
}
