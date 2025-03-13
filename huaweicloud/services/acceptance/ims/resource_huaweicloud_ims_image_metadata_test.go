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

func getImageMetadataResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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
		return nil, fmt.Errorf("error retrieving IMS image metadata: %s", err)
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

func TestAccImageMetadata_basic(t *testing.T) {
	var (
		image        interface{}
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_ims_image_metadata.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImageMetadataResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccImageMetadata_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "__os_version", "Ubuntu 20.04 server 64bit"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "private"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protected", "false"),
					resource.TestCheckResourceAttr(resourceName, "container_format", "bare"),
					resource.TestCheckResourceAttr(resourceName, "disk_format", "vhd"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(resourceName, "min_disk", "80"),
					resource.TestCheckResourceAttrSet(resourceName, "file"),
					resource.TestCheckResourceAttrSet(resourceName, "self"),
					resource.TestCheckResourceAttrSet(resourceName, "schema"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccImageMetadata_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ims_image_metadata" "test" {
  __os_version     = "Ubuntu 20.04 server 64bit"
  visibility       = "private"
  name             = "%s"
  protected        = false
  container_format = "bare"
  disk_format      = "vhd"

  tags = [
    "test=testvalue",
    "image=imagevalue"
  ]

  min_ram  = 1024
  min_disk = 80
}
`, name)
}
