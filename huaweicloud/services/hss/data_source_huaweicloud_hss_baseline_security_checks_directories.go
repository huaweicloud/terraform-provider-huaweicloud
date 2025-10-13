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

// @API HSS GET /v5/{project_id}/baseline/security-checks/directory
func DataSourceBaselineSecurityChecksDirectories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineSecurityChecksDirectoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_os": {
				Type:     schema.TypeString,
				Required: true,
			},
			"select_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"day_of_week": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"hour": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"minute": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"random_offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"baseline_directory_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standard": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"pwd_directory_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// The `checked` field in the API documentation is of type boolean,
						// but it actually returns type string.
						"checked": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBaselineSecurityChecksDirectoriesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?support_os=%v", d.Get("support_os"))
	queryParams = fmt.Sprintf("%s&select_type=%v", queryParams, d.Get("select_type"))

	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceBaselineSecurityChecksDirectoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/baseline/security-checks/directory"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBaselineSecurityChecksDirectoriesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS baseline security checks directories: %s", err)
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
		d.Set("task_condition", flattenTaskCondition(respBody)),
		d.Set("baseline_directory_list", flattenBaselineDirectoryList(respBody)),
		d.Set("pwd_directory_list", flattenPwdDirectoryList(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTaskCondition(resp interface{}) []interface{} {
	taskCondition := utils.PathSearch("task_condition", resp, nil)
	if taskCondition == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":          utils.PathSearch("type", taskCondition, nil),
			"day_of_week":   utils.ExpandToStringList(utils.PathSearch("day_of_week", taskCondition, make([]interface{}, 0)).([]interface{})),
			"hour":          utils.PathSearch("hour", taskCondition, nil),
			"minute":        utils.PathSearch("minute", taskCondition, nil),
			"random_offset": utils.PathSearch("random_offset", taskCondition, nil),
		},
	}
}

func flattenBaselineDirectoryList(resp interface{}) []interface{} {
	baselineDirList := utils.PathSearch("baseline_directory_list", resp, make([]interface{}, 0)).([]interface{})
	if len(baselineDirList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(baselineDirList))
	for _, v := range baselineDirList {
		rst = append(rst, map[string]interface{}{
			"type":      utils.PathSearch("type", v, nil),
			"standard":  utils.PathSearch("standard", v, nil),
			"data_list": flattenBaselineDirectoryListDataList(utils.PathSearch("data_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenBaselineDirectoryListDataList(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		rst = append(rst, map[string]interface{}{
			"name":   utils.PathSearch("name", v, nil),
			"enable": utils.PathSearch("enable", v, nil),
		})
	}

	return rst
}

func flattenPwdDirectoryList(resp interface{}) []interface{} {
	pwdDirList := utils.PathSearch("pwd_directory_list", resp, make([]interface{}, 0)).([]interface{})
	if len(pwdDirList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(pwdDirList))
	for _, v := range pwdDirList {
		rst = append(rst, map[string]interface{}{
			"tag":     utils.PathSearch("tag", v, nil),
			"sub_tag": utils.PathSearch("sub_tag", v, nil),
			"checked": utils.PathSearch("checked", v, nil),
			"key":     utils.PathSearch("key", v, nil),
		})
	}

	return rst
}
