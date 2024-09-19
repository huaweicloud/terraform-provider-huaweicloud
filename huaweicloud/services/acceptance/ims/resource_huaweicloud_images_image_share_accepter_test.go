package ims

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getImsImageShareAccepterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region          = acceptance.HW_REGION_NAME
		getImageHttpUrl = "v2/cloudimages"
		getImageProduct = "ims"
	)
	getImageClient, err := cfg.NewServiceClient(getImageProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS Client: %s", err)
	}

	getImagePath := getImageClient.Endpoint + getImageHttpUrl
	getImagePath = strings.ReplaceAll(getImagePath, "{project_id}", getImageClient.ProjectID)

	imageId := state.Primary.Attributes["image_id"]
	getImageQueryParams := buildGetImageQueryParams(imageId)
	getImagePath += getImageQueryParams

	getImageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getImageResp, err := getImageClient.Request("GET", getImagePath, &getImageOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS shared images: %s", err)
	}

	getImageRespBody, err := utils.FlattenResponse(getImageResp)
	if err != nil {
		return nil, err
	}

	images := utils.PathSearch("images", getImageRespBody, nil)
	if images == nil || len(images.([]interface{})) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return images, nil
}

func TestAccImsImageShareAccepter_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_images_image_share_accepter.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImsImageShareAccepterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckSourceImage(t)
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImsImageShareAccepter_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "image_id", acceptance.HW_IMAGE_SHARE_SOURCE_IMAGE_ID),
				),
			},
		},
	})
}

func testImsImageShareAccepter_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_images_image_share_accepter" "test" {
 image_id = "%s"
}
`, acceptance.HW_IMAGE_SHARE_SOURCE_IMAGE_ID)
}
