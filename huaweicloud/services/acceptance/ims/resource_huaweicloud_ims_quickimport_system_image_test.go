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

func getQuickImportSystemImageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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
		return nil, fmt.Errorf("error retrieving IMS quick import system image: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	image := utils.PathSearch("images[0]", getRespBody, nil)
	// If the list API return empty, then return `404` error code.
	if image == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return image, nil
}

func TestAccQuickImportSystemImage_basic(t *testing.T) {
	var (
		image        interface{}
		name         = acceptance.RandomAccResourceName()
		updateName   = name + "-update"
		resourceName = "huaweicloud_ims_quickimport_system_image.test"
		defaultEpsId = "0"
		migrateEpsId = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getQuickImportSystemImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need setting a non default enterprise project ID.
			acceptance.TestAccPreCheckEpsID(t)
			// This test requires ensuring that there is an image file in the OBS bucket.
			acceptance.TestAccPreCheckImsImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccQuickImportSystemImage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "os_version", "Ubuntu 20.04 server 64bit"),
					resource.TestCheckResourceAttr(resourceName, "image_url", acceptance.HW_IMS_IMAGE_URL),
					resource.TestCheckResourceAttr(resourceName, "min_disk", "100"),
					resource.TestCheckResourceAttr(resourceName, "description", "description test"),
					resource.TestCheckResourceAttr(resourceName, "hw_firmware_type", "uefi"),
					resource.TestCheckResourceAttr(resourceName, "type", "ECS"),
					resource.TestCheckResourceAttr(resourceName, "architecture", "x86"),
					resource.TestCheckResourceAttr(resourceName, "license_type", "byol"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
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
				Config: testAccQuickImportSystemImage_update(updateName, migrateEpsId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "hw_firmware_type", "bios"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", migrateEpsId),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"image_url", "type", "architecture", "license_type",
				},
			},
		},
	})
}

func testAccQuickImportSystemImage_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ims_quickimport_system_image" "test" {
  name                  = "%[1]s"
  os_version            = "Ubuntu 20.04 server 64bit"
  image_url             = "%[2]s"
  min_disk              = 100
  description           = "description test"
  hw_firmware_type      = "uefi"
  type                  = "ECS"
  architecture          = "x86"
  license_type          = "byol"
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_IMS_IMAGE_URL)
}

func testAccQuickImportSystemImage_update(updateName, migrateEpsId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ims_quickimport_system_image" "test" {
  name                  = "%[1]s"
  os_version            = "Ubuntu 20.04 server 64bit"
  image_url             = "%[2]s"
  min_disk              = 100
  hw_firmware_type      = "bios"
  type                  = "ECS"
  architecture          = "x86"
  license_type          = "byol"
  enterprise_project_id = "%[3]s"

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, updateName, acceptance.HW_IMS_IMAGE_URL, migrateEpsId)
}
