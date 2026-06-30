package drs

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

// @API DRS GET /v3/{project_id}/jobs/{job_id}/object/compare/{compare_type}
func DataSourceObjectCompareDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectCompareDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"compare_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"compare_job_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compare_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_db_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_db_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildObjectCompareDetailQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := "?limit=1000"

	if v, ok := d.GetOk("compare_job_id"); ok {
		queryParams += fmt.Sprintf("&compare_job_id=%s", v.(string))
	}
	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	return queryParams
}

func dataSourceObjectCompareDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/{job_id}/object/compare/{compare_type}"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	jobId := d.Get("job_id").(string)
	compareType := d.Get("compare_type").(string)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestPath = strings.ReplaceAll(requestPath, "{compare_type}", compareType)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithQuery := requestPath + buildObjectCompareDetailQueryParams(d, offset)
		resp, err := client.Request("GET", requestPathWithQuery, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS object compare detail: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		compareDetails := utils.PathSearch("compare_detail", respBody, make([]interface{}, 0)).([]interface{})
		if len(compareDetails) == 0 {
			break
		}

		result = append(result, compareDetails...)
		offset += len(compareDetails)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("compare_details", flattenObjectCompareDetails(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenObjectCompareDetails(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"source_db_name":  utils.PathSearch("source_db_name", item, nil),
			"target_db_name":  utils.PathSearch("target_db_name", item, nil),
			"source_db_value": utils.PathSearch("source_db_value", item, nil),
			"target_db_value": utils.PathSearch("target_db_value", item, nil),
			"status":          utils.PathSearch("status", item, nil),
			"error_message":   utils.PathSearch("error_message", item, nil),
		})
	}

	return result
}
