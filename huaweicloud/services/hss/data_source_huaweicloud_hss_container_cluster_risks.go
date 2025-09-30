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

// @API HSS GET /v5/{project_id}/container/cluster/risks
func DataSourceContainerClusterRisks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerClusterRisksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"risk_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"risk_status": {
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
			"risk_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"risk_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"risk_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"risk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"risk_name": {
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
						"risk_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"risk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remediation": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerClusterRisksQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("risk_type"); ok {
		queryParams = fmt.Sprintf("%s&cicd_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("risk_status"); ok {
		queryParams = fmt.Sprintf("%s&risk_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		queryParams = fmt.Sprintf("%s&cluster_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		queryParams = fmt.Sprintf("%s&cluster_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("risk_name"); ok {
		queryParams = fmt.Sprintf("%s&risk_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("risk_level"); ok {
		queryParams = fmt.Sprintf("%s&risk_level=%v", queryParams, v)
	}
	if v, ok := d.GetOk("risk_category"); ok {
		queryParams = fmt.Sprintf("%s&risk_category=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceContainerClusterRisksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/container/cluster/risks"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerClusterRisksQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS container cluster risks: %s", err)
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
		d.Set("data_list", flattenContainerClusterRisks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerClusterRisks(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"risk_id":        utils.PathSearch("risk_id", v, nil),
			"risk_name":      utils.PathSearch("risk_name", v, nil),
			"cluster_id":     utils.PathSearch("cluster_id", v, nil),
			"cluster_name":   utils.PathSearch("cluster_name", v, nil),
			"risk_level":     utils.PathSearch("risk_level", v, nil),
			"risk_category":  utils.PathSearch("risk_category", v, nil),
			"risk_num":       utils.PathSearch("risk_num", v, nil),
			"last_scan_time": utils.PathSearch("last_scan_time", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"remediation":    utils.PathSearch("remediation", v, nil),
		})
	}

	return rst
}
