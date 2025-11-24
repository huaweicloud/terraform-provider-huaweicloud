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

// @API HSS GET /v5/{project_id}/ransomware/optional/vaults
func DataSourceRansomwareBackupVaults() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRansomwareBackupVaultsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vault_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_id": {
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
						"vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vault_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vault_allocated": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vault_charging_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_policy_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resources_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildRansomwareBackupVaultsQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("vault_name"); ok {
		queryParams = fmt.Sprintf("%s&vault_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("vault_id"); ok {
		queryParams = fmt.Sprintf("%s&vault_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceRansomwareBackupVaultsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/ransomware/optional/vaults"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildRansomwareBackupVaultsQueryParams(d)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS ransomware backup vaults: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenRansomwareBackupVaultsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRansomwareBackupVaultsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"vault_id":              utils.PathSearch("vault_id", v, nil),
			"vault_name":            utils.PathSearch("vault_name", v, nil),
			"vault_size":            utils.PathSearch("vault_size", v, nil),
			"vault_used":            utils.PathSearch("vault_used", v, nil),
			"vault_allocated":       utils.PathSearch("vault_allocated", v, nil),
			"vault_charging_mode":   utils.PathSearch("vault_charging_mode", v, nil),
			"vault_status":          utils.PathSearch("vault_status", v, nil),
			"backup_policy_id":      utils.PathSearch("backup_policy_id", v, nil),
			"backup_policy_name":    utils.PathSearch("backup_policy_name", v, nil),
			"backup_policy_enabled": utils.PathSearch("backup_policy_enabled", v, nil),
			"resources_num":         utils.PathSearch("resources_num", v, nil),
		})
	}

	return rst
}
