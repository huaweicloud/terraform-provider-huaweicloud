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

// Due to limitations in the testing environment, `file_id` could not be obtained, resulting in unsuccessful API calls
// and the resource not actually executing successfully.

// @API HSS GET /v5/{project_id}/container/iac/file/risks
func DataSourceContainerIacFileRisks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerIacFileRisksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"file_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
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
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"risk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"risk_name": {
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
						"build_instruction": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerIacFileRisksQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?limit=200&file_id=%v", d.Get("file_id"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
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

	return queryParams
}

func dataSourceContainerIacFileRisksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/container/iac/file/risks"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerIacFileRisksQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		getResp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS container IAC file risks: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
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
		d.Set("data_list", flattenContainerIacFileRisks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerIacFileRisks(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"risk_id":           utils.PathSearch("risk_id", v, nil),
			"rule_id":           utils.PathSearch("rule_id", v, nil),
			"risk_name":         utils.PathSearch("risk_name", v, nil),
			"risk_level":        utils.PathSearch("risk_level", v, nil),
			"risk_category":     utils.PathSearch("risk_category", v, nil),
			"risk_num":          utils.PathSearch("risk_num", v, nil),
			"last_scan_time":    utils.PathSearch("last_scan_time", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"remediation":       utils.PathSearch("remediation", v, nil),
			"build_instruction": utils.PathSearch("build_instruction", v, nil),
		})
	}

	return rst
}
