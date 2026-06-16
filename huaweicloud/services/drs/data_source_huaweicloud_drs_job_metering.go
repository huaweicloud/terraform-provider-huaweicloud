package drs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/jobs/{job_id}/metering
func DataSourceDrsJobMetering() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsJobMeteringRead,

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
			"product_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_size_measure_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"usage_factor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"usage_value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"usage_measure_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDrsJobMeteringRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		jobId   = d.Get("job_id").(string)
		httpUrl = "v5/{project_id}/jobs/{job_id}/metering"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS job metering: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("product_info_list", flattenProductInfoList(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenProductInfoList(respBody interface{}) []interface{} {
	resultsRaw := utils.PathSearch("product_info_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(resultsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resultsRaw))
	for _, item := range resultsRaw {
		result = append(result, map[string]interface{}{
			"id":                       utils.PathSearch("id", item, nil),
			"cloud_service_type":       utils.PathSearch("cloud_service_type", item, nil),
			"resource_type":            utils.PathSearch("resource_type", item, nil),
			"resource_spec_code":       utils.PathSearch("resource_spec_code", item, nil),
			"resource_size":            utils.PathSearch("resource_size", item, nil),
			"resource_size_measure_id": utils.PathSearch("resource_size_measure_id", item, nil),
			"usage_factor":             utils.PathSearch("usage_factor", item, nil),
			"usage_value":              utils.PathSearch("usage_value", item, nil),
			"usage_measure_id":         utils.PathSearch("usage_measure_id", item, nil),
		})
	}
	return result
}
