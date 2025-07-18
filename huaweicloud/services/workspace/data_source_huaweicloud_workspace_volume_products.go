package workspace

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

// @API Workspace GET /v2/{project_id}/volume/products
func DataSourceVolumeProducts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVolumeProductsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the volume products are located.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The availability zone where the volume products are located.`,
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of volume products.`,
			},
			"volume_products": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        volumeProductSchema(),
				Description: `The list of volume products that matched filter parameters.`,
			},
		},
	}
}

func volumeProductSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of volume product.`,
			},
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The volume type of volume product.`,
			},
			"volume_product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product type of volume product.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource type of volume product.`,
			},
			"cloud_service_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cloud service type of volume product.`,
			},
			"domain_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of domain IDs that support this volume.`,
			},
			"names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        volumeProductNameSchema(),
				Description: `The list of volume product name information.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the volume product.`,
			},
		},
	}
}

func volumeProductNameSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The language of volume product name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The volume product name in this language.`,
			},
		},
	}
}

func flattenVolumeProducts(products []interface{}) []interface{} {
	if len(products) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(products))
	for _, item := range products {
		result = append(result, map[string]interface{}{
			"resource_spec_code":  utils.PathSearch("resource_spec_code", item, nil),
			"volume_type":         utils.PathSearch("volume_type", item, nil),
			"volume_product_type": utils.PathSearch("volume_product_type", item, nil),
			"resource_type":       utils.PathSearch("resource_type", item, nil),
			"cloud_service_type":  utils.PathSearch("cloud_service_type", item, nil),
			"domain_ids":          utils.PathSearch("domain_ids", item, nil),
			"names":               flattenVolumeProductNames(utils.PathSearch("name", item, make([]interface{}, 0)).([]interface{})),
			"status":              utils.PathSearch("status", item, nil),
		})
	}
	return result
}

func flattenVolumeProductNames(names []interface{}) []interface{} {
	if len(names) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(names))
	for _, item := range names {
		result = append(result, map[string]interface{}{
			"language": utils.PathSearch("language", item, nil),
			"value":    utils.PathSearch("value", item, nil),
		})
	}

	return result
}

func buildVolumeProductsParams(d *schema.ResourceData) string {
	res := "?"

	if v, ok := d.GetOk("availability_zone"); ok {
		res = fmt.Sprintf("%savailability_zone=%v&", res, v)
	}
	if v, ok := d.GetOk("volume_type"); ok {
		res = fmt.Sprintf("%s&volume_type=%v", res, v)
	}

	return res
}

func queryVolumeProducts(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/volume/products"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildVolumeProductsParams(d)

	requestOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, requestOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("volumes", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceVolumeProductsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	volumeProducts, err := queryVolumeProducts(client, d)
	if err != nil {
		return diag.Errorf("error querying volume products: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("volume_products", flattenVolumeProducts(volumeProducts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
