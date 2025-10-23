package dew

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

// @API DEW GET /v1/{project_id}/dew/cpcs/images
func DataSourceCpcsImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCpcsImagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the image name.`,
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the service type of the image.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort attribute.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort direction.`,
			},
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the images.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image ID.`,
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image name.`,
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The service type of the image.`,
						},
						"arch_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The system architecture of the image.`,
						},
						"specification_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The specification ID.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the image.`,
						},
						"version_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version type.`,
						},
						"trust_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain whitelist.`,
						},
						"vendor_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The vendor name.`,
						},
						"vendor_image_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The vendor image version.`,
						},
						"ccsp_version_need": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The required platform version.`,
						},
						"hw_image_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Huawei image version.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description.`,
						},
					},
				},
			},
		},
	}
}

func buildDataSourceCpcsImagesQueryParams(d *schema.ResourceData, pageNum int) string {
	rst := fmt.Sprintf("?page_num=%d", pageNum)

	if v, ok := d.GetOk("image_name"); ok {
		rst += fmt.Sprintf("&image_name=%v", v)
	}

	if v, ok := d.GetOk("service_type"); ok {
		rst += fmt.Sprintf("&service_type=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	return rst
}

// The first page of page_num is `1`. Because of the quality issues of the API, we need to add a quantitative comparison
// of the value of `total_num`.
func dataSourceCpcsImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/dew/cpcs/images"
		product   = "kms"
		pageNum   = 1
		allImages = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithPageNum := requestPath + buildDataSourceCpcsImagesQueryParams(d, pageNum)
		resp, err := client.Request("GET", requestPathWithPageNum, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DEW CPCS images: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.Errorf("error flattening DEW CPCS images response: %s", err)
		}

		results := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
		if len(results) == 0 {
			break
		}
		allImages = append(allImages, results...)

		totalNum := int(utils.PathSearch("total_num", respBody, float64(0)).(float64))
		if len(allImages) >= totalNum {
			break
		}

		pageNum++
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(generateId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("images", flattenCpcsImagesResponseBody(allImages)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCpcsImagesResponseBody(images []interface{}) []interface{} {
	if len(images) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(images))
	for _, image := range images {
		result = append(result, map[string]interface{}{
			"image_id":             utils.PathSearch("image_id", image, nil),
			"image_name":           utils.PathSearch("image_name", image, nil),
			"service_type":         utils.PathSearch("service_type", image, nil),
			"arch_type":            utils.PathSearch("arch_type", image, nil),
			"specification_id":     utils.PathSearch("specification_id", image, nil),
			"create_time":          utils.PathSearch("create_time", image, nil),
			"version_type":         utils.PathSearch("version_type", image, nil),
			"trust_domain":         utils.PathSearch("trust_domain", image, nil),
			"vendor_name":          utils.PathSearch("vendor_name", image, nil),
			"vendor_image_version": utils.PathSearch("vendor_image_version", image, nil),
			"ccsp_version_need":    utils.PathSearch("ccsp_version_need", image, nil),
			"hw_image_version":     utils.PathSearch("hw_image_version", image, nil),
			"description":          utils.PathSearch("description", image, nil),
		})
	}

	return result
}
