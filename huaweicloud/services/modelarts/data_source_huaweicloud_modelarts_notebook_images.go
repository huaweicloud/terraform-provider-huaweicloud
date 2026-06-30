package modelarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/images
func DataSourceNotebookImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotebookImagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the images are located.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the image to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "BUILD_IN",
				Description: `The type of the image to be queried.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the namespace to which images belong.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID to which images belong.`,
			},
			"cpu_arch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The CPU architecture of the image to be queried.`,
			},

			// Deprecated parameters.
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
				Description: utils.SchemaDesc(
					`The name of the organization (namespace) to which images belong`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},

			// Attributes.
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of images that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the image.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the image.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the namespace to which the image belongs.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The workspace ID to which the image belongs.`,
						},
						"resource_categories": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The resource categories of the image.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"dev_services": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The dev services of the image.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The service type of the image.`,
						},
						"show_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the image to be displayed.`,
						},
						"swr_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage path of the image.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the image.`,
						},
						"description": {
							Type:        schema.TypeString,
							Description: `The description of the image.`,
							Computed:    true,
						},
						"cpu_arch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CPU architecture of the image.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the image.`,
						},
					},
				},
			},
		},
	}
}

func buildListNotebookImagesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("organization"); ok {
		res = fmt.Sprintf("%s&namespace=%v", res, v)
	}

	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		res = fmt.Sprintf("%s&workspace_id=%v", res, v)
	}

	return res
}

func listNotebookImages(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/images?limit={limit}"
		limit   = 200
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildListNotebookImagesQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		images := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, images...)
		if len(images) < limit {
			break
		}
		offset += len(images)
	}

	return result, nil
}

func manualFilterNotebookImages(images []interface{}, architecture string) []interface{} {
	if architecture != "" {
		return utils.PathSearch(fmt.Sprintf("[?arch=='%s']", architecture), images, make([]interface{}, 0)).([]interface{})
	}
	return images
}

func flattenNotebookImages(images []interface{}) []interface{} {
	if len(images) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(images))
	for _, image := range images {
		result = append(result, map[string]interface{}{
			"id":                  utils.PathSearch("id", image, nil),
			"name":                utils.PathSearch("name", image, nil),
			"namespace":           utils.PathSearch("namespace", image, nil),
			"workspace_id":        utils.PathSearch("workspace_id", image, nil),
			"resource_categories": utils.PathSearch("resource_categories", image, nil),
			"dev_services":        utils.PathSearch("dev_services", image, nil),
			"service_type":        utils.PathSearch("service_type", image, nil),
			"show_name":           utils.PathSearch("show_name", image, nil),
			"swr_path":            utils.PathSearch("swr_path", image, nil),
			"type":                utils.PathSearch("type", image, nil),
			"description":         utils.PathSearch("description", image, nil),
			"cpu_arch":            utils.PathSearch("arch", image, nil),
			"status":              utils.PathSearch("status", image, nil),
		})
	}

	return result
}

func dataSourceNotebookImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	images, err := listNotebookImages(client, d)
	if err != nil {
		return diag.Errorf("error listing ModelArts notebook images: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("images", flattenNotebookImages(manualFilterNotebookImages(images, d.Get("cpu_arch").(string)))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
