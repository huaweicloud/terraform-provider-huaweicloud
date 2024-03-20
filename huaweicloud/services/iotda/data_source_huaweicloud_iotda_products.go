package iotda

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

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

func dataSourceProductsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allProducts []model.ProductSummary
		limit       = int32(50)
		offset      int32
	)

	for {
		listOpts := model.ListProductsRequest{
			AppId:       utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			ProductName: utils.StringIgnoreEmpty(d.Get("product_name").(string)),
			Limit:       utils.Int32(limit),
			Offset:      &offset,
		}

		listResp, listErr := client.ListProducts(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA products: %s", listErr)
		}

		if len(*listResp.Products) == 0 {
			break
		}
		allProducts = append(allProducts, *listResp.Products...)
		offset += limit
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

func filterListProducts(products []model.ProductSummary, d *schema.ResourceData) []model.ProductSummary {
	if len(products) == 0 {
		return nil
	}

	rst := make([]model.ProductSummary, 0, len(products))
	for _, v := range products {
		if productID, ok := d.GetOk("product_id"); ok &&
			fmt.Sprint(productID) != utils.StringValue(v.ProductId) {
			continue
		}

		if spaceName, ok := d.GetOk("space_name"); ok &&
			fmt.Sprint(spaceName) != utils.StringValue(v.AppName) {
			continue
		}

		if deviceType, ok := d.GetOk("device_type"); ok &&
			fmt.Sprint(deviceType) != utils.StringValue(v.DeviceType) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenProducts(products []model.ProductSummary) []interface{} {
	if len(products) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(products))
	for _, v := range products {
		rst = append(rst, map[string]interface{}{
			"space_id":          v.AppId,
			"space_name":        v.AppName,
			"id":                v.ProductId,
			"name":              v.Name,
			"device_type":       v.DeviceType,
			"protocol_type":     v.ProtocolType,
			"data_type":         v.DataFormat,
			"manufacturer_name": v.ManufacturerName,
			"industry":          v.Industry,
			"description":       v.Description,
			"created_at":        v.CreateTime,
		})
	}

	return rst
}
