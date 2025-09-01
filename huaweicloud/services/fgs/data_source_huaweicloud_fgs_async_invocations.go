package fgs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/async-invocations
func DataSourceAsyncInvocations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAsyncInvocationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the async invocations are located.`,
			},

			// Required parameters.
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The function URN to which the async invocations belong.`,
			},

			// Optional parameters.
			"request_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The specified request ID of async invocation to be queried.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of async invocations to be queried.`,
			},
			"query_begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The begin time to query async invocations, in RFC3339 format (UTC time).`,
			},
			"query_end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time to query async invocations, in RFC3339 format (UTC time).`,
			},

			// Attributes.
			"invocations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of async invocations that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request ID of the async invocation.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the async invocation.`,
						},
						"error_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error code of the async invocation.`,
						},
						"error_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error message of the async invocation.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the async invocation, in RFC3339 format (UTC time).`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time of the async invocation, in RFC3339 format (UTC time).`,
						},
					},
				},
			},
		},
	}
}

func buildAsyncInvocationsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("request_id"); ok {
		res = fmt.Sprintf("%s&request_id=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("query_begin_time"); ok {
		res = fmt.Sprintf("%s&query_begin_time=%v", res,
			utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(v.(string))/1000, true, "2006-01-02T15:04:05"))
	}
	if v, ok := d.GetOk("query_end_time"); ok {
		res = fmt.Sprintf("%s&query_end_time=%v", res,
			utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(v.(string))/1000, true, "2006-01-02T15:04:05"))
	}

	return res
}

func listAsyncInvocations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/functions/{function_urn}/async-invocations?limit={limit}"
		limit   = 100
		marker  = "0"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{function_urn}", d.Get("function_urn").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildAsyncInvocationsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithMarker := fmt.Sprintf("%s&marker=%v", listPath, marker)
		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		invocations := utils.PathSearch("invocations", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, invocations...)
		if len(invocations) < limit {
			break
		}
		// In this API, marker has the same meaning as offset.
		nextMarker := utils.PathSearch("next_marker", respBody, "0").(string)
		if nextMarker == marker || nextMarker == "0" {
			// Make sure the next marker value is correct, not the previous marker or zero (in the last page).
			break
		}
		marker = nextMarker
	}

	return result, nil
}

func flattenAsyncInvocations(invocations []interface{}) []map[string]interface{} {
	if len(invocations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(invocations))
	for _, invocation := range invocations {
		result = append(result, map[string]interface{}{
			"request_id":    utils.PathSearch("request_id", invocation, nil),
			"status":        utils.PathSearch("status", invocation, nil),
			"error_message": utils.PathSearch("error_message", invocation, nil),
			"error_code":    utils.PathSearch("error_code", invocation, nil),
			"start_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("start_time",
				invocation, nil).(string), "2006-01-02T15:04:05")/1000, true, "2006-01-02T15:04:05Z"),
			"end_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("end_time",
				invocation, nil).(string), "2006-01-02T15:04:05")/1000, true, "2006-01-02T15:04:05Z"),
		})
	}

	return result
}

func dataSourceAsyncInvocationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	invocations, err := listAsyncInvocations(client, d)
	if err != nil {
		return diag.Errorf("error querying FunctionGraph async invocations: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("invocations", flattenAsyncInvocations(invocations)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
