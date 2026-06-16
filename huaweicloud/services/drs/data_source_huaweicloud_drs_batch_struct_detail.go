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

// @API DRS POST /v3/{project_id}/jobs/{type}/batch-struct-detail
func DataSourceDrsBatchStructDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsBatchStructDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cur_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resultsStructSchema(),
			},
		},
	}
}

func resultsStructSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"struct_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     structDetailSchema(),
			},
		},
	}
}

func structDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"total_record": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     listItemSchema(),
			},
		},
	}
}

func listItemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"src_db": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_tb": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_db": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_tb": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildBatchStructDetailBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"jobs": d.Get("jobs").([]interface{}),
	}

	if v, ok := d.GetOk("cur_page"); ok {
		bodyParams["cur_page"] = v.(int)
	}
	if v, ok := d.GetOk("per_page"); ok {
		bodyParams["per_page"] = v.(int)
	}

	return bodyParams
}

func dataSourceDrsBatchStructDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/{type}/batch-struct-detail"
		jobType = d.Get("type").(string)
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{type}", jobType)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildBatchStructDetailBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS batch struct detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	results := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("results", flattenResults(results)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, v := range results {
		rst = append(rst, map[string]interface{}{
			"job_id":        utils.PathSearch("job_id", v, nil),
			"error_code":    utils.PathSearch("error_code", v, nil),
			"error_message": utils.PathSearch("error_message", v, nil),
			"struct_detail": flattenStructDetail(utils.PathSearch("struct_detail", v, nil)),
		})
	}
	return rst
}

func flattenStructDetail(detail interface{}) []interface{} {
	if detail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"total_record": utils.PathSearch("total_record", detail, nil),
			"create_time":  utils.PathSearch("create_time", detail, nil),
			"list":         flattenListItem(utils.PathSearch("list", detail, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenListItem(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(items))
	for _, item := range items {
		rst = append(rst, map[string]interface{}{
			"progress": utils.PathSearch("progress", item, nil),
			"src_db":   utils.PathSearch("src_db", item, nil),
			"src_tb":   utils.PathSearch("src_tb", item, nil),
			"dst_db":   utils.PathSearch("dst_db", item, nil),
			"dst_tb":   utils.PathSearch("dst_tb", item, nil),
		})
	}
	return rst
}
