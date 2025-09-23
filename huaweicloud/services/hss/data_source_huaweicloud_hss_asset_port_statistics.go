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

// @API HSS GET /v5/{project_id}/asset/port/statistics
func DataSourceAssetPortStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetPortStatisticsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the port number.",
			},
			"port_string": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the port string.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the port type.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the port status.",
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sort key.",
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sort direction.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type. The default value is host.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the resource belongs.",
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port number.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port type.",
						},
						"num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of ports.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port status.",
						},
					},
				},
				Description: "The port statistics list.",
			},
		},
	}
}

func buildAssetPortStatisticsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	queryParams := "?limit=100"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("port"); ok {
		queryParams = fmt.Sprintf("%s&port=%v", queryParams, v)
	}
	if v, ok := d.GetOk("port_string"); ok {
		queryParams = fmt.Sprintf("%s&port_string=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	return queryParams
}

func flattenAssetPortStatistics(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"port":   utils.PathSearch("port", v, nil),
			"type":   utils.PathSearch("type", v, nil),
			"num":    utils.PathSearch("num", v, nil),
			"status": utils.PathSearch("status", v, nil),
		})
	}
	return rst
}

func dataSourceAssetPortStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	requestPath := client.Endpoint + "v5/{project_id}/asset/port/statistics"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAssetPortStatisticsQueryParams(d, cfg)
	allPortStatistics := make([]interface{}, 0)
	offset := 0

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS asset port statistics: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		portStatisticsResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(portStatisticsResp) == 0 {
			break
		}
		allPortStatistics = append(allPortStatistics, portStatisticsResp...)
		offset += len(portStatisticsResp)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data_list", flattenAssetPortStatistics(allPortStatistics)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
