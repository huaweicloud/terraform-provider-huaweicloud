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
	getImageCopyClient, err := cfg.NewServiceClient(getImageCopyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS Client: %s", err)
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
				Config: testImsImageCopy_basic(name, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				Config: testImsImageCopy_basic(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_id"},
			},
		},
	})
}

func testImsImageCopy_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_images_image" "test" {
  name        = "%[1]s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testImsImageCopy_basic(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`


resource "huaweicloud_images_image_copy" "test" {
 //image_id = huaweicloud_images_image.test.id
 image_id = "b95678d3-9627-43f7-9f41-ef6778dde3f9"
 name     = "%s"
 target_region = "cn-north-9"
 agency_name = "ims_admin_agency"
}
`, copyImageName)
}
