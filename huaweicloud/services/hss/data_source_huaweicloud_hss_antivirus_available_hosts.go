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

// @API HSS GET /v5/{project_id}/antivirus/host-management/hosts
func DataSourceAntivirusAvailableHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAntivirusAvailableHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scan_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_type": {
				Type:     schema.TypeString,
				Required: true,
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
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_start_time": {
				Type:     schema.TypeInt,
				Optional: true,
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
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAntivirusAvailableHostsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	queryParams = fmt.Sprintf("%s&scan_type=%v", queryParams, d.Get("scan_type"))
	queryParams = fmt.Sprintf("%s&start_type=%v", queryParams, d.Get("start_type"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
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
	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_id"); ok {
		queryParams = fmt.Sprintf("%s&policy_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("next_start_time"); ok {
		queryParams = fmt.Sprintf("%s&next_start_time=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceAntivirusAvailableHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/antivirus/host-management/hosts"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAntivirusAvailableHostsQueryParams(d, epsId)
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS antivirus available hosts: %s", err)
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

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		if int(totalNum) == len(result) {
			break
		}

		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenAntivirusAvailableHostsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAntivirusAvailableHostsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_id":    utils.PathSearch("host_id", v, nil),
			"host_name":  utils.PathSearch("host_name", v, nil),
			"public_ip":  utils.PathSearch("public_ip", v, nil),
			"private_ip": utils.PathSearch("private_ip", v, nil),
			"agent_id":   utils.PathSearch("agent_id", v, nil),
			"os_type":    utils.PathSearch("os_type", v, nil),
			"group_id":   utils.PathSearch("group_id", v, nil),
		})
	}

	return rst
}
