package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/dev-servers/images
func DataSourceDevServerImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevServerImagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the DevServer images are located.`,
			},

			// Optional parameters.
			"server_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the server to be queried.`,
			},
			"resource_flavor_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the resource flavor to be queried.`,
			},

			// Attributes.
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of DevServer images that matched filter parameters.`,
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
						"server_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the server.`,
						},
						"arch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The architecture of the image.`,
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

func buildDevServerImagesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("server_type"); ok {
		res = fmt.Sprintf("%s&server_type=%v", res, v)
	}

	if v, ok := d.GetOk("resource_flavor_name"); ok {
		res = fmt.Sprintf("%s&resource_flavor_name=%v", res, v)
	}

	if len(res) > 0 {
		res = "?" + res[1:]
	}

	return res
}

func listDevServerImages(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/dev-servers/images"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildDevServerImagesQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenDevServerImages(images []interface{}) []interface{} {
	if len(images) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(images))
	for _, image := range images {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("image_id", image, nil),
			"name":        utils.PathSearch("name", image, nil),
			"server_type": utils.PathSearch("server_type", image, nil),
			"arch":        utils.PathSearch("arch", image, nil),
			"status":      utils.PathSearch("status", image, nil),
		})
	}

	return result
}

func dataSourceDevServerImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	images, err := listDevServerImages(client, d)
	if err != nil {
		return diag.Errorf("error listing ModelArts DevServer images: %s", err)
	}

	randomUUID, err := uuid.NewUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("images", flattenDevServerImages(images)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
