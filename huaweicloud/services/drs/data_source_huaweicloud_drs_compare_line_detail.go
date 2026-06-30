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

// @API DRS GET /v3/{project_id}/jobs/{job_id}/compare/{compare_job_id}/line-detail
func DataSourceDrsCompareLineDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsCompareLineDetailRead,

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
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_tb_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_line_compare_result_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tableLineCompareResultInfoSchema(),
			},
		},
	}
}

func tableLineCompareResultInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target_table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"difference_row_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"compare_line_config_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCompareLineDetailQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("db_name"); ok {
		queryParams = fmt.Sprintf("%s&db_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("target_db_name"); ok {
		queryParams = fmt.Sprintf("%s&target_db_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("query_tb_name"); ok {
		queryParams = fmt.Sprintf("%s&query_tb_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDrsCompareLineDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "drs"
		httpUrl      = "v3/{project_id}/jobs/{job_id}/compare/{compare_job_id}/line-detail"
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
		currentRequestPath := requestPath + buildCompareLineDetailQueryParams(d, limit, offset)
		resp, err := client.Request("GET", currentRequestPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS compare line detail: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		infosRaw := utils.PathSearch("table_line_compare_result_infos", respBody, make([]interface{}, 0)).([]interface{})
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
		d.Set("table_line_compare_result_infos", flattenTableLineCompareResultInfos(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTableLineCompareResultInfos(result []interface{}) []interface{} {
	if len(result) == 0 {
		return nil
	}

	infos := make([]interface{}, 0, len(result))
	for _, info := range result {
		infos = append(infos, map[string]interface{}{
			"source_table_name":          utils.PathSearch("source_table_name", info, nil),
			"source_row_num":             utils.PathSearch("source_row_num", info, nil),
			"target_table_name":          utils.PathSearch("target_table_name", info, nil),
			"target_row_num":             utils.PathSearch("target_row_num", info, nil),
			"difference_row_num":         utils.PathSearch("difference_row_num", info, nil),
			"status":                     utils.PathSearch("status", info, nil),
			"compare_line_config_filter": utils.PathSearch("compare_line_config_filter", info, nil),
		})
	}

	return infos
}
