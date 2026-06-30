package secmaster

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

// @API SecMaster GET /v1/{project_id}/subscriptions/orders
func DataSourceSubscriptionOrders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubscriptionOrdersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"page": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csb_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ecs_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     subscriptionOrdersResourceSchema(),
			},
			"subscription_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     subscriptionOrdersSubscriptionSchema(),
			},
		},
	}
}

func subscriptionOrdersResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cloud_service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"to_period": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_reset_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota_reset_cycle_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota_reset_cycle": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"amount": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"original_amount": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"measure_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     subscriptionOrdersTagListSchema(),
			},
			"usages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     subscriptionOrdersUsageSchema(),
			},
		},
	}
}

func subscriptionOrdersTagListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func subscriptionOrdersUsageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_percent": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"quota": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"free": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func subscriptionOrdersSubscriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildSubscriptionOrdersQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("page"); ok {
		return fmt.Sprintf("?page=%s", v.(string))
	}

	return ""
}

func dataSourceSubscriptionOrdersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/subscriptions/orders"
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildSubscriptionOrdersQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// There is an issue with the `limit` and `offset` paging parameters, which are currently not supported.
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster subscription orders: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("csb_version", utils.PathSearch("csb_version", respBody, nil)),
		d.Set("ecs_count", utils.PathSearch("ecs_count", respBody, nil)),
		d.Set("resources", flattenSubscriptionOrdersResources(
			utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("subscription_count", utils.PathSearch("subscription_count", respBody, nil)),
		d.Set("subscriptions", flattenSubscriptionOrdersSubscriptions(
			utils.PathSearch("subscriptions", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSubscriptionOrdersResources(resources []interface{}) []interface{} {
	if len(resources) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		rst = append(rst, map[string]interface{}{
			"resource_id":            utils.PathSearch("resource_id", v, nil),
			"resource_type_name":     utils.PathSearch("resource_type_name", v, nil),
			"resource_size":          utils.PathSearch("resource_size", v, nil),
			"cloud_service":          utils.PathSearch("cloud_service", v, nil),
			"resource_type":          utils.PathSearch("resource_type", v, nil),
			"resource_spec_code":     utils.PathSearch("resource_spec_code", v, nil),
			"to_period":              utils.PathSearch("to_period", v, nil),
			"create_time":            utils.PathSearch("create_time", v, nil),
			"update_time":            utils.PathSearch("update_time", v, nil),
			"expire_time":            utils.PathSearch("expire_time", v, nil),
			"resource_status":        utils.PathSearch("resource_status", v, nil),
			"order_id":               utils.PathSearch("order_id", v, nil),
			"charging_mode":          utils.PathSearch("charging_mode", v, nil),
			"quota_reset_mode":       utils.PathSearch("quota_reset_mode", v, nil),
			"quota_reset_cycle_type": utils.PathSearch("quota_reset_cycle_type", v, nil),
			"quota_reset_cycle":      utils.PathSearch("quota_reset_cycle", v, nil),
			"amount":                 utils.PathSearch("amount", v, nil),
			"original_amount":        utils.PathSearch("original_amount", v, nil),
			"measure_name":           utils.PathSearch("measure_name", v, nil),
			"tag_list": flattenSubscriptionOrdersTagList(
				utils.PathSearch("tag_list", v, make([]interface{}, 0)).([]interface{})),
			"usages": flattenSubscriptionOrdersUsages(
				utils.PathSearch("usages", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenSubscriptionOrdersTagList(tagList []interface{}) []interface{} {
	if len(tagList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tagList))
	for _, v := range tagList {
		rst = append(rst, map[string]interface{}{
			"key":         utils.PathSearch("key", v, nil),
			"value":       utils.PathSearch("value", v, nil),
			"create_time": utils.PathSearch("create_time", v, nil),
			"update_time": utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}

func flattenSubscriptionOrdersUsages(usages []interface{}) []interface{} {
	if len(usages) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(usages))
	for _, v := range usages {
		rst = append(rst, map[string]interface{}{
			"unit":                      utils.PathSearch("unit", v, nil),
			"resource_type_name":        utils.PathSearch("resource_type_name", v, nil),
			"source_resource_spec_code": utils.PathSearch("source_resource_spec_code", v, nil),
			"resource_spec_code":        utils.PathSearch("resource_spec_code", v, nil),
			"source_type":               utils.PathSearch("source_type", v, nil),
			"used_percent":              utils.PathSearch("used_percent", v, nil),
			"quota":                     utils.PathSearch("quota", v, nil),
			"used":                      utils.PathSearch("used", v, nil),
			"free":                      utils.PathSearch("free", v, nil),
		})
	}

	return rst
}

func flattenSubscriptionOrdersSubscriptions(subscriptions []interface{}) []interface{} {
	if len(subscriptions) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(subscriptions))
	for _, v := range subscriptions {
		rst = append(rst, map[string]interface{}{
			"owner":            utils.PathSearch("owner", v, nil),
			"endpoint":         utils.PathSearch("endpoint", v, nil),
			"protocol":         utils.PathSearch("protocol", v, nil),
			"subscription_urn": utils.PathSearch("subscription_urn", v, nil),
			"topic_urn":        utils.PathSearch("topic_urn", v, nil),
			"status":           utils.PathSearch("status", v, nil),
		})
	}

	return rst
}
