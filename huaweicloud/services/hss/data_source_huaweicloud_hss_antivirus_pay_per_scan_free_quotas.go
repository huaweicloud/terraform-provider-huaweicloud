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

// @API HSS GET /v5/{project_id}/antivirus/pay-per-scan/free-quotas
func DataSourceAntivirusPayPerScanFreeQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAntivirusPayPerScanFreeQuotasRead,

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
			"free_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"occupied_free_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"antivirus_already_given": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"antivirus_free_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"report_already_given": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"report_free_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildAntivirusPayPerScanFreeQuotasQueryParams(epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return queryParams
}

func dataSourceAntivirusPayPerScanFreeQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/antivirus/pay-per-scan/free-quotas"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAntivirusPayPerScanFreeQuotasQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS antivirus pay-per-scan free quotas: %s", err)
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
		d.Set("free_quota", utils.PathSearch("free_quota", respBody, nil)),
		d.Set("occupied_free_quota", utils.PathSearch("occupied_free_quota", respBody, nil)),
		d.Set("antivirus_already_given", utils.PathSearch("antivirus_already_given", respBody, nil)),
		d.Set("antivirus_free_quota", utils.PathSearch("antivirus_free_quota", respBody, nil)),
		d.Set("report_already_given", utils.PathSearch("report_already_given", respBody, nil)),
		d.Set("report_free_quota", utils.PathSearch("report_free_quota", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
