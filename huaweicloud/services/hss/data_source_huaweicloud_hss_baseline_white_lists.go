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

// @API HSS GET /v5/{project_id}/baseline/whitelists
func DataSourceBaselineWhiteLists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineWhiteListsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standard": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBaselineWhiteListsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("check_rule_name"); ok {
		queryParams = fmt.Sprintf("%s&check_rule_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		queryParams = fmt.Sprintf("%s&os_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("rule_type"); ok {
		queryParams = fmt.Sprintf("%s&rule_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("tag"); ok {
		queryParams = fmt.Sprintf("%s&tag=%v", queryParams, v)
	}
	if v, ok := d.GetOk("description"); ok {
		queryParams = fmt.Sprintf("%s&description=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceBaselineWhiteListsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		result  = make([]interface{}, 0)
		offset  = 0
		httpUrl = "v5/{project_id}/baseline/whitelists"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBaselineWhiteListsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving baseline white lists: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data_list", flattenBaselineWhiteListsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBaselineWhiteListsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"rule_type":       utils.PathSearch("rule_type", v, nil),
			"os_type":         utils.PathSearch("os_type", v, nil),
			"index_version":   utils.PathSearch("index_version", v, nil),
			"check_type":      utils.PathSearch("check_type", v, nil),
			"standard":        utils.PathSearch("standard", v, nil),
			"tag":             utils.PathSearch("tag", v, nil),
			"check_rule_name": utils.PathSearch("check_rule_name", v, nil),
			"description":     utils.PathSearch("description", v, nil),
		})
	}

	return rst
}
