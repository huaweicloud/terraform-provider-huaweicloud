package cc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCcBandwidthPackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcBandwidthPackagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the bandwidth package ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the bandwidth package name.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the enterprise project that the bandwidth package belongs to.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the bandwidth package status.`,
			},
			"billing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the billing mode of the bandwidth package.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the resource that the bandwidth package is bound to.`,
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies Bandwidth range specified for the bandwidth package.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the tags of the bandwidth package.`,
			},
			"bandwidth_packages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Bandwidth package list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The bandwidth package ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The bandwidth package name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The bandwidth package description.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the account that the bandwidth package belongs to.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project that the bandwidth package belongs to.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Project ID of the bandwidth package.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Time when the resource was created.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Time when the resource was updated.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource that the bandwidth package is bound to.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Type of the resource that the bandwidth package is bound to.`,
						},
						"local_area_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of a local access point.`,
						},
						"remote_area_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of a remote access point.`,
						},
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specification code of the bandwidth package.`,
						},
						"billing_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Billing mode of the bandwidth package.`,
						},
						"tags": common.TagsComputedSchema(),
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Status of the bandwidth package.`,
						},
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Order ID of the bandwidth package.`,
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Product ID of the bandwidth package.`,
						},
						"charge_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Billing option.`,
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Bandwidth range specified for the bandwidth package.`,
						},
						"interflow_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Interflow mode of the bandwidth package.`,
						},
					},
				},
			},
		},
	}
}

// @API CC GET /v3/{domain_id}/ccaas/bandwidth-packages
// @API CC POST /v3/{domain_id}/ccaas/bandwidth-packages/filter
func dataSourceCcBandwidthPackagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)

	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	result, err := getBandwidthPackages(client, cfg, d)
	if err != nil {
		return diag.Errorf("error retrieving bandwidth packages: %s", err)
	}

	if tags, ok := d.GetOk("tags"); ok {
		resourceIDs, err := filterBandwidthPackageByTags(tags.(map[string]interface{}), client, cfg.DomainID)
		if err != nil {
			return diag.Errorf("error filtering bandwidth packages by tags: %s", err)
		}

		result = filter(result, resourceIDs)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("bandwidth_packages", flattenListBandwidthPackageResponseBody(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getBandwidthPackages(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v3/{domain_id}/ccaas/bandwidth-packages"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)

	params := buildBandwidthPackagesQueryParams(d, cfg)
	path += params

	resp, err := pagination.ListAllItems(
		client,
		"marker",
		path,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return nil, err
	}

	curJson := utils.PathSearch("bandwidth_packages", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))

	if bandwidthFilter, ok := d.GetOk("bandwidth"); ok {
		for _, item := range curArray {
			itemMap := item.(map[string]interface{})
			if int(itemMap["bandwidth"].(float64)) == bandwidthFilter {
				rst = append(rst, item)
			}
		}
	} else {
		rst = append(rst, curArray...)
	}

	return rst, nil
}

func flattenListBandwidthPackageResponseBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"domain_id":             utils.PathSearch("domain_id", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"project_id":            utils.PathSearch("project_id", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"updated_at":            utils.PathSearch("updated_at", v, nil),
			"resource_id":           utils.PathSearch("resource_id", v, nil),
			"resource_type":         utils.PathSearch("resource_type", v, nil),
			"local_area_id":         utils.PathSearch("local_area_id", v, nil),
			"remote_area_id":        utils.PathSearch("remote_area_id", v, nil),
			"spec_code":             utils.PathSearch("spec_code", v, nil),
			"billing_mode":          utils.PathSearch("billing_mode", v, nil),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"status":                utils.PathSearch("status", v, nil),
			"order_id":              utils.PathSearch("order_id", v, nil),
			"product_id":            utils.PathSearch("product_id", v, nil),
			"charge_mode":           utils.PathSearch("charge_mode", v, nil),
			"bandwidth":             utils.PathSearch("bandwidth", v, nil),
			"interflow_mode":        utils.PathSearch("interflow_mode", v, nil),
		})
	}
	return rst
}

func filterBandwidthPackageByTags(tags map[string]interface{}, client *golangsdk.ServiceClient, domainID string) ([]string, error) {
	var resourceIDs []string
	httpUrl := "v3/{domain_id}/ccaas/bandwidth-packages/filter"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", domainID)

	if len(tags) < 1 {
		return nil, nil
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildFilterBandwidthPackageByTagsOpts(tags),
	}

	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	curJson := utils.PathSearch("bandwidth_packages[*].id", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	for _, item := range curArray {
		resourceIDs = append(resourceIDs, item.(string))
	}

	return resourceIDs, nil
}

func buildFilterBandwidthPackageByTagsOpts(tagmap map[string]interface{}) map[string]interface{} {
	taglist := make([]interface{}, 0, len(tagmap))

	for k, v := range tagmap {
		taglist = append(taglist, map[string]interface{}{
			"key":    k,
			"values": []string{v.(string)},
		})
	}

	return map[string]interface{}{
		"tags": taglist,
	}
}

func buildBandwidthPackagesQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("billing_mode"); ok {
		res = fmt.Sprintf("%s&billing_mode=%v", res, v)
	}
	if v, ok := d.GetOk("resource_id"); ok {
		res = fmt.Sprintf("%s&resource_id=%v", res, v)
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
