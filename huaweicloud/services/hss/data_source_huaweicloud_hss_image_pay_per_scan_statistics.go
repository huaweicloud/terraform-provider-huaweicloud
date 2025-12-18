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

// @API HSS GET /v5/{project_id}/image/pay-per-scan/statistics
func DataSourceImagePayPerScanStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagePayPerScanStatisticsRead,

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
			// Attributes
			"repository_scan_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cicd_scan_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"free_quota_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"already_given": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"image_free_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"high_risk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceImageTypeRiskInfoSchema(),
			},
			"has_risk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceImageTypeRiskInfoSchema(),
			},
			"total": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceImageTypeRiskInfoSchema(),
			},
			"unscan": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceImageTypeRiskInfoSchema(),
			},
		},
	}
}

func dataSourceImageTypeRiskInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"local": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"registriy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cicd": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceImagePayPerScanStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/image/pay-per-scan/statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath = fmt.Sprintf("%s?enterprise_project_id=%v", requestPath, epsId)
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS image pay per scan statistics: %s", err)
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

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("repository_scan_num", utils.PathSearch("repository_scan_num", respBody, nil)),
		d.Set("cicd_scan_num", utils.PathSearch("cicd_scan_num", respBody, nil)),
		d.Set("free_quota_num", utils.PathSearch("free_quota_num", respBody, nil)),
		d.Set("already_given", utils.PathSearch("already_given", respBody, nil)),
		d.Set("image_free_quota", utils.PathSearch("image_free_quota", respBody, nil)),
		d.Set("high_risk",
			flattenDataSourceImageTypeRiskInfo(utils.PathSearch("high_risk", respBody, nil))),
		d.Set("has_risk",
			flattenDataSourceImageTypeRiskInfo(utils.PathSearch("has_risk", respBody, nil))),
		d.Set("total",
			flattenDataSourceImageTypeRiskInfo(utils.PathSearch("total", respBody, nil))),
		d.Set("unscan",
			flattenDataSourceImageTypeRiskInfo(utils.PathSearch("unscan", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataSourceImageTypeRiskInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"local":     utils.PathSearch("local", resp, nil),
			"registriy": utils.PathSearch("registriy", resp, nil),
			"cicd":      utils.PathSearch("cicd", resp, nil),
		},
	}
}
