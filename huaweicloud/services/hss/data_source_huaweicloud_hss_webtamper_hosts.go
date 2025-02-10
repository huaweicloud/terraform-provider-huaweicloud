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

// @API HSS GET /v5/{project_id}/webtamper/hosts
func DataSourceWebTamperHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebTamperHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rasp_protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_bit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rasp_protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"anti_tampering_times": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"detect_tampering_times": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildWebTamperHostsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=20"
	queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		queryParams = fmt.Sprintf("%s&public_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_name"); ok {
		queryParams = fmt.Sprintf("%s&group_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		queryParams = fmt.Sprintf("%s&os_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_status"); ok {
		queryParams = fmt.Sprintf("%s&protect_status=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceWebTamperHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
		httpUrl = "v5/{project_id}/webtamper/hosts"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildWebTamperHostsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS web tamper hosts: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		hostsResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(hostsResp) == 0 {
			break
		}

		result = append(result, hostsResp...)

		totalNum := utils.PathSearch("total_num", getRespBody, float64(0)).(float64)
		if int(totalNum) == len(result) {
			break
		}

		offset += len(hostsResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("hosts", flattenWebTamperHosts(filterWebTamperHosts(result, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterWebTamperHosts(hostsResp []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(hostsResp))
	for _, v := range hostsResp {
		if raspProtectStatus, ok := d.GetOk("rasp_protect_status"); ok &&
			fmt.Sprint(raspProtectStatus) != utils.PathSearch("rasp_protect_status", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenWebTamperHosts(hostsResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(hostsResp))
	for _, v := range hostsResp {
		rst = append(rst, map[string]interface{}{
			"id":                     utils.PathSearch("host_id", v, nil),
			"name":                   utils.PathSearch("host_name", v, nil),
			"public_ip":              utils.PathSearch("public_ip", v, nil),
			"private_ip":             utils.PathSearch("private_ip", v, nil),
			"group_name":             utils.PathSearch("group_name", v, nil),
			"os_bit":                 utils.PathSearch("os_bit", v, nil),
			"os_type":                utils.PathSearch("os_type", v, nil),
			"protect_status":         utils.PathSearch("protect_status", v, nil),
			"rasp_protect_status":    utils.PathSearch("rasp_protect_status", v, nil),
			"anti_tampering_times":   utils.PathSearch("anti_tampering_times", v, nil),
			"detect_tampering_times": utils.PathSearch("detect_tampering_times", v, nil),
		})
	}

	return rst
}
