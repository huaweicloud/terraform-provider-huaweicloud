package iotda

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

// @API IoTDA GET /v5/iot/{project_id}/products
func DataSourceProducts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"space_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"device_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"space_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manufacturer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"industry": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildProductsQueryParams(d *schema.ResourceData) string {
	rst := ""
	if v, ok := d.GetOk("space_id"); ok {
		rst += fmt.Sprintf("&app_id=%v", v)
	}

	if v, ok := d.GetOk("product_name"); ok {
		rst += fmt.Sprintf("&product_name=%v", v)
	}

	return rst
}

func dataSourceProductsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v5/iot/{project_id}/products?limit=50"
		product     = "iotda"
		allProducts []interface{}
		offset      = 0
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildProductsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving IoTDA products: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		products := utils.PathSearch("products", respBody, make([]interface{}, 0)).([]interface{})
		if len(products) == 0 {
			break
		}

		allProducts = append(allProducts, products...)
		offset += len(products)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	targetProducts := filterListProducts(allProducts, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("products", flattenProducts(targetProducts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListProducts(products []interface{}, d *schema.ResourceData) []interface{} {
	if len(products) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(products))
	for _, v := range products {
		if productID, ok := d.GetOk("product_id"); ok &&
			fmt.Sprint(productID) != utils.PathSearch("product_id", v, "").(string) {
			continue
		}

		if spaceName, ok := d.GetOk("space_name"); ok &&
			fmt.Sprint(spaceName) != utils.PathSearch("app_name", v, "").(string) {
			continue
		}

		if deviceType, ok := d.GetOk("device_type"); ok &&
			fmt.Sprint(deviceType) != utils.PathSearch("device_type", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenProducts(products []interface{}) []interface{} {
	if len(products) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(products))
	for _, v := range products {
		rst = append(rst, map[string]interface{}{
			"space_id":          utils.PathSearch("app_id", v, nil),
			"space_name":        utils.PathSearch("app_name", v, nil),
			"id":                utils.PathSearch("product_id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"device_type":       utils.PathSearch("device_type", v, nil),
			"protocol_type":     utils.PathSearch("protocol_type", v, nil),
			"data_type":         utils.PathSearch("data_format", v, nil),
			"manufacturer_name": utils.PathSearch("manufacturer_name", v, nil),
			"industry":          utils.PathSearch("industry", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"created_at":        utils.PathSearch("create_time", v, nil),
		})
	}

	return rst
}
