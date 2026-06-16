package gaussdb

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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/hba-info/history
func DataSourceGaussdbClientAuthConfigHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbClientAuthConfigHistoryRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hba_histories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     clientAuthConfigHistorySchema(),
			},
		},
	}
}

func clientAuthConfigHistorySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"before_confs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hbaConfSchema(),
			},
			"after_confs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     hbaConfSchema(),
			},
		},
	}
}

func hbaConfSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"method": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussdbClientAuthConfigHistoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/hba-info/history"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	currentTotal := 0
	results := make([]interface{}, 0)
	for {
		currentPath := listPath + buildGetClientAuthConfigHistoryQueryParams(d, currentTotal)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB client auth config history: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.Errorf("error flattening GaussDB client auth config history response: %s", err)
		}

		hbaHistories := utils.PathSearch("hba_histories", listRespBody, make([]interface{}, 0)).([]interface{})
		results = append(results, hbaHistories...)

		totalCount := utils.PathSearch("total_count", listRespBody, float64(0)).(float64)
		currentTotal += len(hbaHistories)
		if currentTotal >= int(totalCount) {
			break
		}
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("hba_histories", flattenGetClientAuthConfigHistoryBody(results)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetClientAuthConfigHistoryQueryParams(d *schema.ResourceData, offset int) string {
	res := fmt.Sprintf("?limit=100&offset=%d", offset)

	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}

	return res
}

func flattenGetClientAuthConfigHistoryBody(resp []interface{}) []interface{} {
	res := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		res = append(res, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"status":       utils.PathSearch("status", v, nil),
			"time":         utils.PathSearch("time", v, nil),
			"fail_reason":  utils.PathSearch("fail_reason", v, nil),
			"before_confs": flattenGetHbaConfsBody(utils.PathSearch("before_confs", v, make([]interface{}, 0)).([]interface{})),
			"after_confs":  flattenGetHbaConfsBody(utils.PathSearch("after_confs", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return res
}

func flattenGetHbaConfsBody(resp []interface{}) []interface{} {
	res := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		res = append(res, map[string]interface{}{
			"type":     utils.PathSearch("type", v, nil),
			"database": utils.PathSearch("database", v, nil),
			"user":     utils.PathSearch("user", v, nil),
			"address":  utils.PathSearch("address", v, nil),
			"method":   utils.PathSearch("method", v, nil),
		})
	}
	return res
}
