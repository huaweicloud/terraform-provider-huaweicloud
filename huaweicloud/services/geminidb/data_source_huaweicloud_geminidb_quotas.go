package geminidb

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

// @API GeminiDB GET /v3/{project_id}/quotas
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datastore_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildQuotasQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("datastore_type"); ok {
		queryParams = fmt.Sprintf("%s&datastore_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("mode"); ok {
		queryParams = fmt.Sprintf("%s&mode=%v", queryParams, v)
	}
	if v, ok := d.GetOk("product_type"); ok {
		queryParams = fmt.Sprintf("%s&product_type=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/quotas"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildQuotasQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB database quotas: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("quotas", flattenQuotas(utils.PathSearch("quotas.resources", getRespBody,
			make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQuotas(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"type":  utils.PathSearch("type", v, nil),
			"quota": utils.PathSearch("quota", v, nil),
			"used":  utils.PathSearch("used", v, nil),
		})
	}

	return rst
}
