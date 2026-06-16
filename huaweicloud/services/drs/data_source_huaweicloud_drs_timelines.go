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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/timelines
func DataSourceDrsTimelines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsTimelinesRead,

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
			"timelines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildTimelinesQueryParams(offset int) string {
	defaultLimit := 1000
	if offset == 0 {
		return fmt.Sprintf("?limit=%d", defaultLimit)
	}

	return fmt.Sprintf("?limit=%d&offset=%d", defaultLimit, offset)
}

func dataSourceDrsTimelinesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/timelines"
		result  = make([]interface{}, 0)
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{job_id}", d.Get("job_id").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	for {
		currentListPath := listPath + buildTimelinesQueryParams(offset)

		listResp, err := client.Request("GET", currentListPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS timelines: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		timelines := utils.PathSearch("timelines", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(timelines) == 0 {
			break
		}

		result = append(result, timelines...)
		offset += len(timelines)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("timelines", flattenTimelines(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTimelines(timelinesResp []interface{}) []interface{} {
	if len(timelinesResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(timelinesResp))
	for _, v := range timelinesResp {
		rst = append(rst, map[string]interface{}{
			"name":           utils.PathSearch("name", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"operation_time": utils.PathSearch("operation_time", v, nil),
			"user_name":      utils.PathSearch("user_name", v, nil),
		})
	}
	return rst
}
