package das

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

// @API DAS GET /v3/{project_id}/instances/{instance_id}/transaction
func DataSourceHistoryTransactions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHistoryTransactionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the DAS history transactions are located.",
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the database instance.",
			},
			"datastore_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database type.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The start time of the query range, in RFC3339 format.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The end time of the query range, in RFC3339 format.",
			},

			// Optional parameters.
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The sort order.",
			},
			"order_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The field used for sorting.",
			},
			"last_sec_min": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The minimum duration of the transaction, in seconds.",
			},
			"last_sec_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum duration of the transaction, in seconds.",
			},

			// Attributes.
			"transactions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of history transactions that matched the filter parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_sec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The transaction duration, in seconds.",
						},
						"wait_locks": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of wait locks.",
						},
						"hold_locks": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of hold locks.",
						},
						"occurrence_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The occurrence time, in RFC3339 format.",
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The transaction content.",
						},
						"collect_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collect time, in RFC3339 format.",
						},
					},
				},
			},
		},
	}
}

func dataSourceHistoryTransactionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	transactions, err := listHistoryTransactions(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS history transactions: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("transactions", flattenHistoryTransactions(transactions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listHistoryTransactions(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/transaction"
		perPage = 50
		curPage = 1
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildHistoryTransactionsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithPage := fmt.Sprintf("%s&page_num=%d&page_size=%d", listPath, curPage, perPage)

		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		transactionList := utils.PathSearch("transaction_info_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, transactionList...)
		if len(transactionList) < perPage {
			break
		}
		curPage++
	}

	return result, nil
}

func buildHistoryTransactionsQueryParams(d *schema.ResourceData) string {
	startTime := utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string))
	endTime := utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string))

	res := fmt.Sprintf("?datastore_type=%v&start_at=%d&end_at=%d",
		d.Get("datastore_type").(string), startTime, endTime)

	if v, ok := d.GetOk("order_by"); ok {
		res = fmt.Sprintf("%s&order_by=%v", res, v)
	}
	if v, ok := d.GetOk("order_field"); ok {
		res = fmt.Sprintf("%s&order=%v", res, v)
	}
	if v, ok := d.GetOk("last_sec_min"); ok {
		res = fmt.Sprintf("%s&last_sec_min=%v", res, v)
	}
	if v, ok := d.GetOk("last_sec_max"); ok {
		res = fmt.Sprintf("%s&last_sec_max=%v", res, v)
	}

	return res
}

func flattenHistoryTransactions(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"last_sec":   utils.PathSearch("last_sec", item, nil),
			"wait_locks": utils.PathSearch("wait_locks", item, nil),
			"hold_locks": utils.PathSearch("hold_locks", item, nil),
			"detail":     utils.PathSearch("detail", item, nil),
			"occurrence_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("occurrence_time", item, float64(0)).(float64))/1000, false),
			"collect_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("collect_time", item, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
