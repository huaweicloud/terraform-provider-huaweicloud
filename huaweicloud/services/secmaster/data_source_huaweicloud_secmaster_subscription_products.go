package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/subscriptions/products
func DataSourceSubscriptionProducts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubscriptionProductsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"basic": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildSubscriptionProductsCodeObjectSchema(),
			},
			"standard": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildSubscriptionProductsCodeObjectSchema(),
			},
			"professional": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildSubscriptionProductsCodeObjectSchema(),
			},
			"large_screen": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildSubscriptionProductsCodeObjectSchema(),
			},
			"log_collection": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildSubscriptionProductsCodeObjectSchema(),
			},
			"log_retention": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildSubscriptionProductsCodeObjectSchema(),
			},
			"log_analysis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildSubscriptionProductsCodeObjectSchema(),
			},
			"soar": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildSubscriptionProductsCodeObjectSchema(),
			},
		},
	}
}

func buildSubscriptionProductsCodeObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cloud_service_type": {
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
			"resource_size_measure_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"usage_factor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage_measure_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSubscriptionProductsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/subscriptions/products"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"X-Language": "en-us"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster subscription products: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("basic", flattenSubscriptionProductsCodeObject(utils.PathSearch("basic", respBody, nil))),
		d.Set("standard", flattenSubscriptionProductsCodeObject(utils.PathSearch("standard", respBody, nil))),
		d.Set("professional", flattenSubscriptionProductsCodeObject(utils.PathSearch("professional", respBody, nil))),
		d.Set("large_screen", flattenSubscriptionProductsCodeObject(utils.PathSearch("large_screen", respBody, nil))),
		d.Set("log_collection", flattenSubscriptionProductsCodeObject(utils.PathSearch("log_collection", respBody, nil))),
		d.Set("log_retention", flattenSubscriptionProductsCodeObject(utils.PathSearch("log_retention", respBody, nil))),
		d.Set("log_analysis", flattenSubscriptionProductsCodeObject(utils.PathSearch("log_analysis", respBody, nil))),
		d.Set("soar", flattenSubscriptionProductsCodeObject(utils.PathSearch("soar", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSubscriptionProductsCodeObject(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cloud_service_type":       utils.PathSearch("cloud_service_type", respBody, nil),
			"resource_type":            utils.PathSearch("resource_type", respBody, nil),
			"resource_spec_code":       utils.PathSearch("resource_spec_code", respBody, nil),
			"resource_size_measure_id": utils.PathSearch("resource_size_measure_id", respBody, nil),
			"usage_factor":             utils.PathSearch("usage_factor", respBody, nil),
			"usage_measure_id":         utils.PathSearch("usage_measure_id", respBody, nil),
			"region_id":                utils.PathSearch("region_id", respBody, nil),
		},
	}
}
