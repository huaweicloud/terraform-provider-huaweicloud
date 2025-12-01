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

// @API HSS GET /v5/{project_id}/host-management/hosts-risk
func DataSourceHostsRisk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostsRiskRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_id_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"install_result_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detect_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vulnerability": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"baseline": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"intrusion": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildHostsRiskQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	for _, hostId := range d.Get("host_id_list").([]interface{}) {
		queryParams = fmt.Sprintf("%s&host_id_list=%v", queryParams, hostId)
	}

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceHostsRiskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/host-management/hosts-risk"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildHostsRiskQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS hosts risk: %s", err)
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
		d.Set("region", region),
		d.Set("total_num", utils.PathSearch("total_num", respBody, nil)),
		d.Set("data_list", flattenHostsRiskDataList(
			utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHostsRiskDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_id":             utils.PathSearch("host_id", v, nil),
			"agent_status":        utils.PathSearch("agent_status", v, nil),
			"install_result_code": utils.PathSearch("install_result_code", v, nil),
			"version":             utils.PathSearch("version", v, nil),
			"protect_status":      utils.PathSearch("protect_status", v, nil),
			"detect_result":       utils.PathSearch("detect_result", v, nil),
			"asset":               utils.PathSearch("asset", v, nil),
			"vulnerability":       utils.PathSearch("vulnerability", v, nil),
			"baseline":            utils.PathSearch("baseline", v, nil),
			"intrusion":           utils.PathSearch("intrusion", v, nil),
		})
	}

	return rst
}
