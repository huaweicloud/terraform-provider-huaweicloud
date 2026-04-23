package rfs

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

// @API RFS GET /v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}/prices
func DataSourceExecutionPlanPrices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExecutionPlanPricesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_plan_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"currency": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"module_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"unsupported_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_price": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"charge_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sale_price": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"discount": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"original_price": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"period_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"period_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"best_discount_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"best_discount_price": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"official_website_discount_price": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildExecutionPlanPricesQueryParams(d *schema.ResourceData) string {
	rst := ""

	if v, ok := d.GetOk("stack_id"); ok {
		rst += fmt.Sprintf("&stack_id=%s", v.(string))
	}

	if v, ok := d.GetOk("execution_plan_id"); ok {
		rst += fmt.Sprintf("&execution_plan_id=%s", v.(string))
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceExecutionPlanPricesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}/prices"
		stackName = d.Get("stack_name").(string)
		planName  = d.Get("execution_plan_name").(string)
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestPath = strings.ReplaceAll(requestPath, "{execution_plan_name}", planName)
	requestPath += buildExecutionPlanPricesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving RFS execution plan prices: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(reqUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("currency", utils.PathSearch("currency", respBody, nil)),
		d.Set("items", flattenExecutionPlanPricesItems(
			utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenExecutionPlanPricesItems(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(items))
	for _, v := range items {
		rst = append(rst, map[string]interface{}{
			"resource_type":       utils.PathSearch("resource_type", v, nil),
			"resource_name":       utils.PathSearch("resource_name", v, nil),
			"index":               utils.PathSearch("index", v, nil),
			"module_address":      utils.PathSearch("module_address", v, nil),
			"supported":           utils.PathSearch("supported", v, nil),
			"unsupported_message": utils.PathSearch("unsupported_message", v, nil),
			"resource_price": flattenResourcePrice(
				utils.PathSearch("resource_price", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenResourcePrice(resourcePriceResp []interface{}) []interface{} {
	if len(resourcePriceResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resourcePriceResp))
	for _, v := range resourcePriceResp {
		rst = append(rst, map[string]interface{}{
			"charge_mode":         utils.PathSearch("charge_mode", v, nil),
			"sale_price":          utils.PathSearch("sale_price", v, nil),
			"discount":            utils.PathSearch("discount", v, nil),
			"original_price":      utils.PathSearch("original_price", v, nil),
			"period_type":         utils.PathSearch("period_type", v, nil),
			"period_count":        utils.PathSearch("period_count", v, nil),
			"best_discount_type":  utils.PathSearch("best_discount_type", v, nil),
			"best_discount_price": utils.PathSearch("best_discount_price", v, nil),
			"official_website_discount_price": utils.PathSearch(
				"official_website_discount_price", v, nil),
		})
	}

	return rst
}
