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

func getImageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "ims"
		httpUrl = "v2/cloudimages"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath += fmt.Sprintf("?id=%s", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS image: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	image := utils.PathSearch("images[0]", getRespBody, nil)
	// If the list API return empty, return `404` error code.
	if image == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return image, nil
}

func TestAccImageRegistration_basic(t *testing.T) {
	var (
		image        interface{}
		resourceName = "huaweicloud_ims_image_registration.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test requires ensuring that there is an image file in the OBS bucket.
			acceptance.TestAccPreCheckImsImageUrl(t)
			acceptance.TestAccPreCheckIMSImageMetadataID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccImageRegistration_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_IMS_IMAGE_METADATA_ID),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(resourceName, "visibility"),
					resource.TestCheckResourceAttrSet(resourceName, "protected"),
					resource.TestCheckResourceAttrSet(resourceName, "container_format"),
					resource.TestCheckResourceAttrSet(resourceName, "min_ram"),
					resource.TestCheckResourceAttrSet(resourceName, "disk_format"),
					resource.TestCheckResourceAttrSet(resourceName, "min_disk"),
					resource.TestCheckResourceAttrSet(resourceName, "__os_version"),
					resource.TestCheckResourceAttrSet(resourceName, "file"),
					resource.TestCheckResourceAttrSet(resourceName, "self"),
					resource.TestCheckResourceAttrSet(resourceName, "schema"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "__os_bit"),
					resource.TestCheckResourceAttrSet(resourceName, "__isregistered"),
					resource.TestCheckResourceAttrSet(resourceName, "__platform"),
					resource.TestCheckResourceAttrSet(resourceName, "__os_type"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_env_type"),
					resource.TestCheckResourceAttrSet(resourceName, "__image_source_type"),
					resource.TestCheckResourceAttrSet(resourceName, "__imagetype"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "__image_size"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_url"},
			},
		},
	})
}

func testAccImageRegistration_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ims_image_registration" "test" {
  image_id  = "%[1]s"
  image_url = "%[2]s"
}
`, acceptance.HW_IMS_IMAGE_METADATA_ID, acceptance.HW_IMS_IMAGE_URL)
}
