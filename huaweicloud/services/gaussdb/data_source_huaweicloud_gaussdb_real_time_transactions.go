package gaussdb

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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/transactions/list
func DataSourceGaussdbRealTimeTransactions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbRealTimeTransactionsRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transaction_query_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     gaussdbTransactionQueryInfoSchema(),
			},
			"rows": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussdbRealTimeTransactionSchema(),
			},
		},
	}
}

func gaussdbTransactionQueryInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"xlog_quantity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datnames": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"usenames": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_addrs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func gaussdbRealTimeTransactionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sessionid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"query_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"datname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_addr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"usename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backend_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"xact_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_change": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"waiting": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_sql_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"top_xid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_xid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exec_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"xlog_quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussdbRealTimeTransactionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/transactions/list"
	postPath := client.Endpoint + httpUrl
	postPath = strings.ReplaceAll(postPath, "{project_id}", client.ProjectID)
	postPath = strings.ReplaceAll(postPath, "{instance_id}", d.Get("instance_id").(string))

	postOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	limit := 100
	offset := 0
	res := make([]interface{}, 0)

	for {
		postOpt.JSONBody = utils.RemoveNil(buildGetGaussdbRealTimeTransactionsBodyParams(d, offset, limit))
		postResp, err := client.Request("POST", postPath, &postOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB real time transactions: %s", err)
		}

		postRespBody, err := utils.FlattenResponse(postResp)
		if err != nil {
			return diag.FromErr(err)
		}

		transactions := flattenGetGaussdbRealTimeTransactionsBody(postRespBody)
		if len(transactions) == 0 {
			break
		}
		res = append(res, transactions...)
		offset += limit
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("rows", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetGaussdbRealTimeTransactionsBodyParams(d *schema.ResourceData, offset, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id":      d.Get("node_id"),
		"component_id": d.Get("component_id"),
		"offset":       offset,
		"limit":        limit,
	}

	if v, ok := d.GetOk("transaction_query_info"); ok {
		queryInfo := v.([]interface{})[0].(map[string]interface{})
		transactionQueryInfo := map[string]interface{}{}
		if val, ok := queryInfo["exec_time"]; ok && val.(string) != "" {
			transactionQueryInfo["exec_time"] = val
		}
		if val, ok := queryInfo["xlog_quantity"]; ok && val.(string) != "" {
			transactionQueryInfo["xlog_quantity"] = val
		}
		if val, ok := queryInfo["datnames"]; ok {
			transactionQueryInfo["datnames"] = val
		}
		if val, ok := queryInfo["usenames"]; ok {
			transactionQueryInfo["usenames"] = val
		}
		if val, ok := queryInfo["client_addrs"]; ok {
			transactionQueryInfo["client_addrs"] = val
		}
		if len(transactionQueryInfo) > 0 {
			bodyParams["transaction_query_info"] = transactionQueryInfo
		}
	}

	return bodyParams
}

func flattenGetGaussdbRealTimeTransactionsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("rows", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"sessionid":        utils.PathSearch("sessionid", v, nil),
			"pid":              utils.PathSearch("pid", v, nil),
			"query_id":         utils.PathSearch("query_id", v, nil),
			"datname":          utils.PathSearch("datname", v, nil),
			"client_addr":      utils.PathSearch("client_addr", v, nil),
			"client_port":      utils.PathSearch("client_port", v, nil),
			"usename":          utils.PathSearch("usename", v, nil),
			"query":            utils.PathSearch("query", v, nil),
			"backend_start":    utils.PathSearch("backend_start", v, nil),
			"xact_start":       utils.PathSearch("xact_start", v, nil),
			"application_name": utils.PathSearch("application_name", v, nil),
			"state":            utils.PathSearch("state", v, nil),
			"state_change":     utils.PathSearch("state_change", v, nil),
			"query_start":      utils.PathSearch("query_start", v, nil),
			"waiting":          utils.PathSearch("waiting", v, nil),
			"unique_sql_id":    utils.PathSearch("unique_sql_id", v, nil),
			"top_xid":          utils.PathSearch("top_xid", v, nil),
			"current_xid":      utils.PathSearch("current_xid", v, nil),
			"exec_time":        utils.PathSearch("exec_time", v, nil),
			"xlog_quantity":    utils.PathSearch("xlog_quantity", v, nil),
		})
	}
	return res
}
