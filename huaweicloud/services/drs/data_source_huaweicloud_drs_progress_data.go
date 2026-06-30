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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/progress-data/{type}
func DataSourceDrsProgressData() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsProgressDataRead,

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
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flow_compare_data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildProgressDataQueryParams(offset int) string {
	queryParams := "?limit=1000"
	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	return queryParams
}

func dataSourceDrsProgressDataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "drs"
		httpUrl    = "v5/{project_id}/jobs/{job_id}/progress-data/{type}"
		offset     = 0
		result     = make([]interface{}, 0)
		createTime string
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{type}", d.Get("type").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithQuery := requestPath + buildProgressDataQueryParams(offset)
		resp, err := client.Request("GET", requestPathWithQuery, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS progress data: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		createTime = utils.PathSearch("create_time", respBody, "").(string)

		flowCompareData := utils.PathSearch("flow_compare_data", respBody, make([]interface{}, 0)).([]interface{})
		if len(flowCompareData) == 0 {
			break
		}

		result = append(result, flowCompareData...)
		offset += len(flowCompareData)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("create_time", createTime),
		d.Set("flow_compare_data", flattenFlowCompareData(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFlowCompareData(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"src_db":   utils.PathSearch("src_db", item, nil),
			"src_tb":   utils.PathSearch("src_tb", item, nil),
			"dst_db":   utils.PathSearch("dst_db", item, nil),
			"dst_tb":   utils.PathSearch("dst_tb", item, nil),
			"progress": utils.PathSearch("progress", item, nil),
		})
	}

	return result
}
