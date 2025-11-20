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

// @API HSS GET /v5/{project_id}/webtamper/rasp/protect-history
func DataSourceWebtamperRaspProtectHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebtamperRaspProtectHistoryRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_level": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"threat_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_level": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attacked_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildWebtamperRaspProtectHistoryQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?start_time=%v&end_time=%v&limit=200", d.Get("start_time"), d.Get("start_time"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("alarm_level"); ok {
		queryParams = fmt.Sprintf("%s&alarm_level=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceWebtamperRaspProtectHistoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/webtamper/rasp/protect-history"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWebtamperRaspProtectHistoryQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving dynamic WTP events: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenWebtamperRaspProtectHistory(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWebtamperRaspProtectHistory(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_ip":      utils.PathSearch("host_ip", v, nil),
			"host_name":    utils.PathSearch("host_name", v, nil),
			"alarm_time":   utils.PathSearch("alarm_time", v, nil),
			"threat_type":  utils.PathSearch("threat_type", v, nil),
			"alarm_level":  utils.PathSearch("alarm_level", v, nil),
			"source_ip":    utils.PathSearch("source_ip", v, nil),
			"attacked_url": utils.PathSearch("attacked_url", v, nil),
		})
	}

	return rst
}
