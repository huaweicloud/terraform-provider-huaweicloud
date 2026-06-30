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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/ddl
func DataSourceDrsDdls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsDdlsRead,

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
			"start_seq_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_seq_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ddls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"seqno": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"checkpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ddl_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ddl_text": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exe_result": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"record_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"clean_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDrsDdlsQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := ""

	if v, ok := d.GetOk("start_seq_no"); ok {
		queryParams += fmt.Sprintf("&start_seq_no=%s", v.(string))
	}
	if v, ok := d.GetOk("end_seq_no"); ok {
		queryParams += fmt.Sprintf("&end_seq_no=%s", v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams += fmt.Sprintf("&status=%s", v.(string))
	}
	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceDrsDdlsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/ddl"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	jobId := d.Get("job_id").(string)
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithQuery := requestPath + buildDrsDdlsQueryParams(d, offset)
		resp, err := client.Request("GET", requestPathWithQuery, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS DDLs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		ddlsResp := utils.PathSearch("ddl_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(ddlsResp) == 0 {
			break
		}

		result = append(result, ddlsResp...)
		offset += len(ddlsResp)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ddls", flattenDrsDdls(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDrsDdls(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"seqno":         utils.PathSearch("seqno", item, nil),
			"checkpoint":    utils.PathSearch("checkpoint", item, nil),
			"status":        utils.PathSearch("status", item, nil),
			"ddl_timestamp": utils.PathSearch("ddl_timestamp", item, nil),
			"ddl_text":      utils.PathSearch("ddl_text", item, nil),
			"exe_result":    utils.PathSearch("exe_result", item, nil),
			"record_time":   utils.PathSearch("record_time", item, nil),
			"clean_time":    utils.PathSearch("clean_time", item, nil),
		})
	}

	return result
}
