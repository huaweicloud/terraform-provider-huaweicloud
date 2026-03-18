package cfw

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

// @API CFW GET /v1/{project_id}/report-profile
func DataSourceReportProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReportProfilesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Compared to the API documentation, the `data` structure hierarchy has been simplified, eliminating the
			// `records` hierarchy, and the `total` has been extracted to the same level as the `data`.
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"report_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildReportProfilesQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?fw_instance_id=%s&limit=%d&offset=%d",
		d.Get("fw_instance_id").(string), limit, offset)
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%s", queryParams, v.(string))
	}

	return queryParams
}

func dataSourceReportProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/report-profile"
		recordsResult = make([]interface{}, 0)
		limit         = 1024
		offset        = 0
		total         float64
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		pathWithParams := requestPath + buildReportProfilesQueryParams(d, limit, offset)
		resp, err := client.Request("GET", pathWithParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CFW report profiles: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		total = utils.PathSearch("data.total", respBody, float64(0)).(float64)
		records := utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{})
		recordsResult = append(recordsResult, records...)
		if len(records) < limit {
			break
		}

		offset += len(records)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenReportProfilesData(recordsResult)),
		d.Set("total", total),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReportProfilesData(records []interface{}) []interface{} {
	if len(records) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(records))
	for _, v := range records {
		result = append(result, map[string]interface{}{
			"profile_id": utils.PathSearch("profile_id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"category":   utils.PathSearch("category", v, nil),
			"status":     utils.PathSearch("status", v, nil),
			"report_id":  utils.PathSearch("report_id", v, nil),
			"last_time":  utils.PathSearch("last_time", v, nil),
		})
	}

	return result
}
