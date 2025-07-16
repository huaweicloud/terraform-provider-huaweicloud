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

// @API HSS GET /v5/{project_id}/event/white-list/login
func DataSourceEventLoginWhiteLists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventLoginWhiteListsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_user_name": {
				Type:     schema.TypeString,
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remarks": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildEventLoginWhiteListsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("login_ip"); ok {
		queryParams = fmt.Sprintf("%s&login_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("login_user_name"); ok {
		queryParams = fmt.Sprintf("%s&login_user_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceEventLoginWhiteListsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	requestPath := client.Endpoint + "v5/{project_id}/event/white-list/login"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEventLoginWhiteListsQueryParams(d, epsId)
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS event login white lists: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)

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
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenEventLoginWhiteListsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEventLoginWhiteListsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"private_ip":              utils.PathSearch("private_ip", v, nil),
			"login_ip":                utils.PathSearch("login_ip", v, nil),
			"login_user_name":         utils.PathSearch("login_user_name", v, nil),
			"update_time":             utils.PathSearch("update_time", v, nil),
			"remarks":                 utils.PathSearch("remarks", v, nil),
			"enterprise_project_name": utils.PathSearch("enterprise_project_name", v, nil),
		})
	}

	return rst
}
