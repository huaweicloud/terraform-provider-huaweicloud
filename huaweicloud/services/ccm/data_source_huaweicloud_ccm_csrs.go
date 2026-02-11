package ccm

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SCM GET /v3/scm/csr
func DataSourceCcmCsrs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcmCsrsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the CSR name.`,
			},
			"private_key_algo": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the key algorithm type.`,
			},
			"csr_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The CSR list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CSR ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CSR name.`,
						},
						"csr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CSR content.`,
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name bound to the CSR.`,
						},
						"sans": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The additional domain name bound to the CSR.`,
						},
						"private_key_algo": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key algorithm.`,
						},
						"usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CSR usage.`,
						},
						"company_country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The country.`,
						},
						"company_province": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The province.`,
						},
						"company_city": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The city.`,
						},
						"company_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The company name.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The CSR creation time.`,
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The CSR update time.`,
						},
					},
				},
			},
		},
	}
}

func buildReadCsrsQueryParams(d *schema.ResourceData) string {
	rst := "?limit=50"

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		rst += fmt.Sprintf("&name=%v", v)
	}

	if v, ok := d.GetOk("private_key_algo"); ok && v.(string) != "" {
		rst += fmt.Sprintf("&private_key_algo=%v", v)
	}

	return rst
}

// The pagination parameter `offset` of this API has an issue.
// Temporarily, it's being handled by only querying the first page.
func dataSourceCcmCsrsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v3/scm/csr"
		product = "scm"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildReadCsrsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CCM SSL CSRs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	csrListRespArray := utils.PathSearch("csr_list", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("csr_list", flattenCcmCsrsResponse(csrListRespArray)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCcmCsrsResponse(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"csr":              utils.PathSearch("csr", v, nil),
			"domain_name":      utils.PathSearch("domain_name", v, nil),
			"sans":             utils.PathSearch("sans", v, nil),
			"private_key_algo": utils.PathSearch("private_key_algo", v, nil),
			"usage":            utils.PathSearch("usage", v, nil),
			"company_country":  utils.PathSearch("company_country", v, nil),
			"company_province": utils.PathSearch("company_province", v, nil),
			"company_city":     utils.PathSearch("company_city", v, nil),
			"company_name":     utils.PathSearch("company_name", v, nil),
			"create_time":      utils.PathSearch("create_time", v, nil),
			"update_time":      utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
