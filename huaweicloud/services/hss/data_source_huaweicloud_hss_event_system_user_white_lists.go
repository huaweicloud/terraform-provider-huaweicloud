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

// @API HSS GET /v5/{project_id}/event/white-list/userlist
func DataSourceEventSystemUserWhiteLists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventSystemUserWhiteListsRead,
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
			"system_user_name": {
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
			"remain_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"limit_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enterprise_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_user_name_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remarks": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildEventSystemUserWhiteListsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=100"

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
	if v, ok := d.GetOk("system_user_name"); ok {
		queryParams = fmt.Sprintf("%s&system_user_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceEventSystemUserWhiteListsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "hss"
		epsId     = cfg.GetEnterpriseProjectID(d)
		result    = make([]interface{}, 0)
		offset    = 0
		totalNum  float64
		remainNum float64
		limitNum  float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/event/white-list/userlist"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEventSystemUserWhiteListsQueryParams(d, epsId)
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS event system user white lists: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		remainNum = utils.PathSearch("remain_num", respBody, float64(0)).(float64)
		limitNum = utils.PathSearch("limit_num", respBody, float64(0)).(float64)

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
		d.Set("remain_num", remainNum),
		d.Set("limit_num", limitNum),
		d.Set("data_list", flattenEventSystemUserWhiteListsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEventSystemUserWhiteListsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"enterprise_project_name": utils.PathSearch("enterprise_project_name", v, nil),
			"host_id":                 utils.PathSearch("host_id", v, nil),
			"host_name":               utils.PathSearch("host_name", v, nil),
			"private_ip":              utils.PathSearch("private_ip", v, nil),
			"public_ip":               utils.PathSearch("public_ip", v, nil),
			"system_user_name_list":   utils.ExpandToStringList(utils.PathSearch("system_user_name_list", v, make([]interface{}, 0)).([]interface{})),
			"update_time":             utils.PathSearch("update_time", v, nil),
			"remarks":                 utils.PathSearch("remarks", v, nil),
		})
	}

	return rst
}
