package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/custom/rule/config
func DataSourceCustomRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     customRulesDataListSchema(),
			},
		},
	}
}

func customRulesDataListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_block": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hash_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_all_host": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCustomRulesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("rule_id"); ok {
		queryParams = fmt.Sprintf("%s&rule_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("rule_name"); ok {
		queryParams = fmt.Sprintf("%s&rule_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceCustomRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		result  = make([]interface{}, 0)
		offset  = 0
		httpUrl = "v5/{project_id}/custom/rule/config"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildCustomRulesQueryParams(d)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS custom rules: %s", err)
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

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenCustomRulesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCustomRulesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"rule_id":     utils.PathSearch("rule_id", v, nil),
			"host_num":    utils.PathSearch("host_num", v, nil),
			"rule_name":   utils.PathSearch("rule_name", v, nil),
			"rule_status": utils.PathSearch("rule_status", v, nil),
			"rule_type":   utils.PathSearch("rule_type", v, nil),
			"auto_block":  utils.PathSearch("auto_block", v, nil),
			"hash_type":   utils.PathSearch("hash_type", v, nil),
			"is_all_host": utils.PathSearch("is_all_host", v, nil),
			"create_time": utils.PathSearch("create_time", v, nil),
			"update_time": utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
