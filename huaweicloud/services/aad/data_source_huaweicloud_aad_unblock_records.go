package aad

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v1/unblockservice/{domain_id}/unblock-record
func DataSourceUnblockRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUnblockRecordsRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specified the account ID of IAM user.",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specified the start time of unblock record.",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specified the end time of unblock record.",
			},
			"unblock_record": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        unblockRecord(),
				Description: `The unblock record.`,
			},
		},
	}
}

func unblockRecord() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address.",
			},
			"executor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The executor.",
			},
			"block_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The block id.`,
			},
			"blocking_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The blocking time.",
			},
			"unblocking_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unblocking time.",
			},
			"unblock_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unblock type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unblock status.",
			},
			"sort_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The sort time.",
			},
		},
	}
	return &sc
}

func dataSourceUnblockRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)

	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	httpUrl := "v1/unblockservice/{domain_id}/unblock-record"
	httpUrl = strings.ReplaceAll(httpUrl, "{domain_id}", d.Get("domain_id").(string))
	httpUrl += fmt.Sprintf("?start_time=%d&end_time=%d", d.Get("start_time").(int), d.Get("end_time").(int))
	listPath := client.Endpoint + httpUrl
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving AAD unblock records: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("unblock_record", flattenUnblockRecord(utils.PathSearch("unblock_record", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenUnblockRecord(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"ip":              utils.PathSearch("ip", v, nil),
			"executor":        utils.PathSearch("executor", v, nil),
			"block_id":        utils.PathSearch("block_id", v, nil),
			"blocking_time":   utils.PathSearch("blocking_time", v, nil),
			"unblocking_time": utils.PathSearch("unblocking_time", v, nil),
			"unblock_type":    utils.PathSearch("unblock_type", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"sort_time":       utils.PathSearch("sort_time", v, nil),
		})
	}
	return rst
}
