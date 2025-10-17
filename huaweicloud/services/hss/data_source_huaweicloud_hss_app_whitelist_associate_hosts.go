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

// @API HSS GET /v5/{project_id}/app/host
func DataSourceAppWhitelistAssociateHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppWhitelistAssociateHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"learning_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"apply_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_value": {
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
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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
						"asset_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"learning_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apply_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"intercept": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAppWhitelistAssociateHostsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("policy_name"); ok {
		queryParams = fmt.Sprintf("%s&policy_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("learning_status"); ok {
		queryParams = fmt.Sprintf("%s&learning_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("apply_status"); ok {
		queryParams = fmt.Sprintf("%s&apply_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		queryParams = fmt.Sprintf("%s&os_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_id"); ok {
		queryParams = fmt.Sprintf("%s&policy_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		queryParams = fmt.Sprintf("%s&public_ip=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceAppWhitelistAssociateHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v5/{project_id}/app/host"
		offset   = 0
		result   = make([]interface{}, 0)
		totalNum float64
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildAppWhitelistAssociateHostsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving hosts associated with the process whitelist policy: %s", err)
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

		totalNum = utils.PathSearch("total_num", getRespBody, float64(0)).(float64)
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
		d.Set("data_list", flattenAppWhitelistAssociateHosts(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppWhitelistAssociateHosts(appHosts []interface{}) []interface{} {
	if len(appHosts) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(appHosts))
	for _, v := range appHosts {
		rst = append(rst, map[string]interface{}{
			"host_id":         utils.PathSearch("host_id", v, nil),
			"host_name":       utils.PathSearch("host_name", v, nil),
			"public_ip":       utils.PathSearch("public_ip", v, nil),
			"private_ip":      utils.PathSearch("private_ip", v, nil),
			"asset_value":     utils.PathSearch("asset_value", v, nil),
			"policy_name":     utils.PathSearch("policy_name", v, nil),
			"event_num":       utils.PathSearch("event_num", v, nil),
			"os_type":         utils.PathSearch("os_type", v, nil),
			"learning_status": utils.PathSearch("learning_status", v, nil),
			"apply_status":    utils.PathSearch("apply_status", v, nil),
			"intercept":       utils.PathSearch("intercept", v, nil),
			"policy_id":       utils.PathSearch("policy_id", v, nil),
			"policy_type":     utils.PathSearch("policy_type", v, nil),
		})
	}

	return rst
}
