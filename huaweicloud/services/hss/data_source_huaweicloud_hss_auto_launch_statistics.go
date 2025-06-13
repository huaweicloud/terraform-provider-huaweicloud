package hss

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/asset/auto-launch/statistics
func DataSourceAutoLaunchStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoLaunchStatisticsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the auto launch name.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the auto launch type.",
			},
			"data_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        resourceAutoLaunchSchema(),
				Description: "The auto launch statistics list.",
			},
		},
	}
}

func resourceAutoLaunchSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The auto launch name.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The auto launch type.",
			},
			"num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of hosts that have this auto launch item.",
			},
		},
	}
	return &sc
}

func buildAutoLaunchStatisticsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	queryParams := "?limit=200"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("name"); ok {
		// There are special characters in name, which need to be escaped
		encodedName := url.QueryEscape(v.(string))
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, encodedName)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	return queryParams
}

func flattenAutoLaunchStatistics(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"type": utils.PathSearch("type", v, nil),
			"num":  utils.PathSearch("num", v, nil),
		})
	}
	return rst
}

func dataSourceAutoLaunchStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/asset/auto-launch/statistics"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAutoLaunchStatisticsQueryParams(d, cfg)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS auto launch statistics: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	statisticsResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data_list", flattenAutoLaunchStatistics(statisticsResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
