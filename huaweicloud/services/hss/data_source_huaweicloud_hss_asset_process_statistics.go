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

// @API HSS GET /v5/{project_id}/asset/process/statistics
func DataSourceAssetProcessStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetProcessStatisticsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the process path.",
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
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        assetProcessStatisticsSchema(),
				Description: "The process statistics list.",
			},
		},
	}
}

func assetProcessStatisticsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The process path.",
			},
			"num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of processes.",
			},
		},
	}
	return &sc
}

func buildAssetProcessStatisticsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	queryParams := "?limit=100"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("path"); ok {
		queryParams = fmt.Sprintf("%s&path=%v", queryParams, v)
	}
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	return queryParams
}

func flattenAssetProcessStatistics(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"path": utils.PathSearch("path", v, nil),
			"num":  utils.PathSearch("num", v, nil),
		})
	}
	return rst
}

func dataSourceAssetProcessStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	requestPath := client.Endpoint + "v5/{project_id}/asset/process/statistics"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAssetProcessStatisticsQueryParams(d, cfg)
	allProcessStatistics := make([]interface{}, 0)
	offset := 0

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS process statistics: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		processStatisticsResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(processStatisticsResp) == 0 {
			break
		}
		allProcessStatistics = append(allProcessStatistics, processStatisticsResp...)
		offset += len(processStatisticsResp)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data_list", flattenAssetProcessStatistics(allProcessStatistics)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
