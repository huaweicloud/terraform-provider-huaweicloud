package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getImsImageCopyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getImageCopy: query IMS image copy
	var (
		getImageCopyHttpUrl = "v2/cloudimages"
		getImageCopyProduct = "ims"
	)

	var getImageCopyClient *golangsdk.ServiceClient
	var err error

	targetRegion := state.Primary.Attributes["target_region"]
	if targetRegion == "" {
		getImageCopyClient, err = cfg.NewServiceClient(getImageCopyProduct, region)
		if err != nil {
			return nil, fmt.Errorf("error creating IMS Client: %s", err)
		}
	} else {
		getImageCopyClient, err = cfg.NewServiceClient(getImageCopyProduct, targetRegion)
		if err != nil {
			return nil, fmt.Errorf("error creating IMS Client: %s", err)
		}
	}

	getImageCopyPath := getImageCopyClient.Endpoint + getImageCopyHttpUrl

	getImageCopyQueryParams := buildGetImageCopyQueryParams(state)
	getImageCopyPath += getImageCopyQueryParams

	getImageCopyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getImageCopyResp, err := getImageCopyClient.Request("GET", getImageCopyPath, &getImageCopyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ImsImageCopy: %s", err)
	}
	getImageCopyRespBody, err := utils.FlattenResponse(getImageCopyResp)
	if err != nil {
		return nil, err
	}
	images := utils.PathSearch("images", getImageCopyRespBody, nil).([]interface{})
	if len(images) == 0 {
		return nil, fmt.Errorf("error get copy image")
	}
	return images[0], nil
}

func buildGetImageCopyQueryParams(state *terraform.ResourceState) string {
	res := ""
	res = fmt.Sprintf("%s&id=%v", res, state.Primary.ID)

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func TestAccImsImageCopy_basic(t *testing.T) {
	var obj interface{}

	sourceImageName := acceptance.RandomAccResourceName()
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_images_image_copy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImsImageCopyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImsImageCopy_basic(sourceImageName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "it's a test"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
				),
			},
			{
				Config: testImsImageCopy_update(sourceImageName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "it's a test"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
		},
	})
}

func TestAccImsImageCrossRegionCopy_basic(t *testing.T) {
	var obj interface{}

	sourceImageName := acceptance.RandomAccResourceName()
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_images_image_copy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImsImageCopyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImsImageCrossRegionCopy_basic(sourceImageName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "it's a test"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
				),
			},
			{
				Config: testImsImageCrossRegionCopy_update(sourceImageName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "it's a test"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
		},
	})
}

func testImsImageCopy_basic(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_images_image_copy" "test" {
 source_image_id = huaweicloud_images_image.test.id
 name            = "%s"
 description     = "it's a test"

 tags = {
    key1 = "value1"
    key2 = "value2"
 }
}
`, testAccImsImage_basic(baseImageName), copyImageName)
}

func testImsImageCopy_update(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_images_image_copy" "test" {
 source_image_id = huaweicloud_images_image.test.id
 name            = "%s"
 description     = "it's a test"

 tags = {
    key1 = "value1_update"
    key3 = "value3"
    key4 = "value4"
 }
}
`, testAccImsImage_basic(baseImageName), copyImageName)
}

func testImsImageCrossRegionCopy_basic(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_images_image_copy" "test" {
 source_image_id = huaweicloud_images_image.test.id
 name            = "%s"
 description     = "it's a test"
 target_region   = "cn-north-9"
 agency_name     = "ims_admin_agency"

 tags = {
    key1 = "value1"
    key2 = "value2"
 }
}
`, testAccImsImage_basic(baseImageName), copyImageName)
}

func testImsImageCrossRegionCopy_update(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_images_image_copy" "test" {
 source_image_id = huaweicloud_images_image.test.id
 name            = "%s"
 description     = "it's a test"
 target_region   = "cn-north-9"
 agency_name     = "ims_admin_agency"

 tags = {
    key1 = "value1_update"
    key3 = "value3"
    key4 = "value4"
 }
}
`, testAccImsImage_basic(baseImageName), copyImageName)
}
