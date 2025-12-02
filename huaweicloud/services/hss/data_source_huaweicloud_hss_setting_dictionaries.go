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

// @API HSS GET /v5/{project_id}/setting/dictionaries
func DataSourceSettingDictionaries() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSettingDictionariesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scene": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func buildSettingDictionariesQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?group_code=%v&limit=200", d.Get("group_code"))

	if v, ok := d.GetOk("scene"); ok {
		queryParams = fmt.Sprintf("%s&scene=%v", queryParams, v)
	}

	if v, ok := d.GetOk("code"); ok {
		queryParams = fmt.Sprintf("%s&code=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceSettingDictionariesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/setting/dictionaries"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildSettingDictionariesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving dictionaries: %s", err)
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
		d.Set("data_list", flattenSettingDictionaries(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSettingDictionaries(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		rst = append(rst, map[string]interface{}{
			"code":   utils.PathSearch("code", v, nil),
			"value":  utils.PathSearch("value", v, nil),
			"values": utils.PathSearch("values", v, nil),
		})
	}

	return rst
}
