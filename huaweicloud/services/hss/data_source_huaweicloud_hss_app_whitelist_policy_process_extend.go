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

// @API HSS GET /v5/{project_id}/app/{policy_id}/process-extend
func DataSourceAppWhitelistPolicyProcessExtend() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppWhitelistPolicyProcessExtendRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
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
						"process_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cmdline": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAppWhitelistPolicyProcessExtendQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?host_id=%s", d.Get("host_id").(string))

	if epsId != "" {
		queryParams += fmt.Sprintf("&enterprise_project_id=%s", epsId)
	}

	return queryParams
}

func dataSourceAppWhitelistPolicyProcessExtendRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		epsId    = cfg.GetEnterpriseProjectID(d)
		product  = "hss"
		httpUrl  = "v5/{project_id}/app/{policy_id}/process-extend"
		policyId = d.Get("policy_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", policyId)
	requestPath += buildAppWhitelistPolicyProcessExtendQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS app whitelist policy process extend: %s", err)
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
		d.Set("data_list", flattenAppWhitelistPolicyProcessExtendDataList(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppWhitelistPolicyProcessExtendDataList(resp interface{}) []interface{} {
	dataList := utils.PathSearch("data_list", resp, make([]interface{}, 0)).([]interface{})
	if len(dataList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		rst = append(rst, map[string]interface{}{
			"process_name": utils.PathSearch("process_name", v, nil),
			"process_path": utils.PathSearch("process_path", v, nil),
			"process_hash": utils.PathSearch("process_hash", v, nil),
			"container_id": utils.PathSearch("container_id", v, nil),
			"cmdline":      utils.PathSearch("cmdline", v, nil),
			"file_size":    utils.PathSearch("file_size", v, nil),
		})
	}

	return rst
}
