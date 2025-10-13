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

// @API HSS POST /v5/{project_id}/config
func DataSourceConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config_name_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"config_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildConfigsQueryParams(epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("%s?enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func buildConfigsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"config_name_list": utils.ExpandToStringList(d.Get("config_name_list").([]interface{})),
	}
}

func dataSourceConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/config"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildConfigsQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildConfigsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS configs: %s", err)
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
		d.Set("data_list", flattenConfigs(utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigs(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"config_name":  utils.PathSearch("config_name", v, nil),
			"config_value": utils.PathSearch("config_value", v, nil),
		})
	}

	return rst
}
