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

// @API DRS GET /v3/{project_id}/jobs/{job_id}/compare/{compare_job_id}/content-overview
func DataSourceDrsCompareContentOverview() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsCompareContentOverviewRead,

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
			"compare_job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_compare_result_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     contentCompareResultOverviewInfoSchema(),
			},
		},
	}
}

func contentCompareResultOverviewInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_db": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_db": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"compare_end_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_inconsistent_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"uncomparable_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDrsCompareContentOverviewRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "drs"
		httpUrl      = "v3/{project_id}/jobs/{job_id}/compare/{compare_job_id}/content-overview"
		jobId        = d.Get("job_id").(string)
		compareJobId = d.Get("compare_job_id").(string)
		result       = make([]interface{}, 0)
		limit        = 1000
		offset       = 0
		mErr         *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestPath = strings.ReplaceAll(requestPath, "{compare_job_id}", compareJobId)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		currentRequestPath := fmt.Sprintf("%s?limit=%d&offset=%d", requestPath, limit, offset)
		resp, err := client.Request("GET", currentRequestPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS compare content overview: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		infosRaw := utils.PathSearch("content_compare_result_infos", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, infosRaw...)

		if len(infosRaw) < limit {
			break
		}

		offset += len(infosRaw)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("content_compare_result_infos", flattenContentCompareResultInfos(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContentCompareResultInfos(result []interface{}) []interface{} {
	if len(result) == 0 {
		return nil
	}

	infos := make([]interface{}, 0, len(result))
	for _, info := range result {
		infos = append(infos, map[string]interface{}{
			"status":                utils.PathSearch("status", info, nil),
			"source_db":             utils.PathSearch("source_db", info, nil),
			"target_db":             utils.PathSearch("target_db", info, nil),
			"compare_num":           utils.PathSearch("compare_num", info, nil),
			"compare_end_num":       utils.PathSearch("compare_end_num", info, nil),
			"data_inconsistent_num": utils.PathSearch("data_inconsistent_num", info, nil),
			"uncomparable_num":      utils.PathSearch("uncomparable_num", info, nil),
		})
	}

	return infos
}
