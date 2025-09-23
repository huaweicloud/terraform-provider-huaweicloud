// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CPH
// ---------------------------------------------------------------

package cph

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CPH GET /v1/{project_id}/cloud-phone/phone-images
func DataSourcePhoneImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourcePhoneImagesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_public": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  `The image type.`,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
			"image_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The label of image.`,
			},
			"images": {
				Type:        schema.TypeList,
				Elem:        phoneImagesImagesSchema(),
				Computed:    true,
				Description: `The list of images detail.`,
			},
		},
	}
}

func phoneImagesImagesSchema() *schema.Resource {
	sc := schema.Resource{
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
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The os type of the image.`,
			},
			"os_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The os name of the image.`,
			},
			"is_public": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The image type.`,
			},
			"image_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The label of the image.`,
			},
		},
	}
	return &sc
}

func resourcePhoneImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listImages: Query the list of CPH phone images
	var (
		listImagesHttpUrl = "v1/{project_id}/cloud-phone/phone-images"
		listImagesProduct = "cph"
	)
	listImagesClient, err := cfg.NewServiceClient(listImagesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	listImagesPath := listImagesClient.Endpoint + listImagesHttpUrl
	listImagesPath = strings.ReplaceAll(listImagesPath, "{project_id}", listImagesClient.ProjectID)

	listImagesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listImagesResp, err := listImagesClient.Request("GET", listImagesPath, &listImagesOpt)
	if err != nil {
		return diag.Errorf("error retrieving CPH phone images: %s", err)
	}

	listImagesRespBody, err := utils.FlattenResponse(listImagesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("images", filterListPhoneImagesImages(
			flattenListPhoneImagesImages(listImagesRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListPhoneImagesImages(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("phone_images", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("image_name", v, nil),
			"os_type":     utils.PathSearch("os_type", v, nil),
			"os_name":     utils.PathSearch("os_name", v, nil),
			"id":          utils.PathSearch("image_id", v, nil),
			"is_public":   utils.PathSearch("is_public", v, nil),
			"image_label": utils.PathSearch("image_label", v, nil),
		})
	}
	return rst
}

func filterListPhoneImagesImages(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("is_public"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("is_public", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("image_label"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("image_label", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
