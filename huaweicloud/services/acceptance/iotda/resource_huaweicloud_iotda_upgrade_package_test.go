package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getUpgradePackageResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	getOpts := &model.ShowOtaPackageRequest{
		PackageId: state.Primary.ID,
	}

	resp, err := client.ShowOtaPackage(getOpts)
	if err != nil {
		return nil, fmt.Errorf("error querying IoTDA OTA upgrade package")
	}

	return resp, nil
}

func TestAccUpgradePackage_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_iotda_upgrade_package.test"
		rName        = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getUpgradePackageResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testUpgradePackage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "product_id", "huaweicloud_iotda_product.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "softwarePackage"),
					resource.TestCheckResourceAttr(resourceName, "version", "v1.0"),
					resource.TestCheckResourceAttr(resourceName, "file_location.0.obs_location.0.region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttrPair(resourceName, "file_location.0.obs_location.0.bucket_name",
						"huaweicloud_obs_bucket_object.test", "bucket"),
					resource.TestCheckResourceAttrPair(resourceName, "file_location.0.obs_location.0.object_key",
						"huaweicloud_obs_bucket_object.test", "key"),
					resource.TestCheckResourceAttr(resourceName, "support_source_versions.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", "description test"),
					resource.TestCheckResourceAttr(resourceName, "custom_info", "custom_info test"),
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

func TestAccUpgradePackage_derived(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_iotda_upgrade_package.test"
		rName        = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getUpgradePackageResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testUpgradePackage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "product_id", "huaweicloud_iotda_product.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "softwarePackage"),
					resource.TestCheckResourceAttr(resourceName, "version", "v1.0"),
					resource.TestCheckResourceAttr(resourceName, "file_location.0.obs_location.0.region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttrPair(resourceName, "file_location.0.obs_location.0.bucket_name",
						"huaweicloud_obs_bucket_object.test", "bucket"),
					resource.TestCheckResourceAttrPair(resourceName, "file_location.0.obs_location.0.object_key",
						"huaweicloud_obs_bucket_object.test", "key"),
					resource.TestCheckResourceAttr(resourceName, "support_source_versions.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", "description test"),
					resource.TestCheckResourceAttr(resourceName, "custom_info", "custom_info test"),
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

func testUpgradePackageWithObsBucket_base() string {
	randInt := acctest.RandInt()

	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket = "tf-acc-test-bucket-%d"
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket  = huaweicloud_obs_bucket.test.bucket
  key     = "test-key.img"
  content = "some_bucket_content"
}
`, randInt)
}

func testUpgradePackage_basic(name string) string {
	obsBucketBasic := testUpgradePackageWithObsBucket_base()
	productbasic := testProduct_basic(name)

	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_upgrade_package" "test" {
  space_id   = huaweicloud_iotda_space.test.id
  type       = "softwarePackage"
  product_id = huaweicloud_iotda_product.test.id
  version    = "v1.0"

  file_location {
    obs_location {
      region      = "%[3]s"
      bucket_name = huaweicloud_obs_bucket_object.test.bucket
      object_key  = huaweicloud_obs_bucket_object.test.key
    }
  }

  support_source_versions = [
    "v1.0",
    "v2.0",
    "v3.0",
  ]

  description = "description test"
  custom_info = "custom_info test"
}
`, productbasic, obsBucketBasic, acceptance.HW_REGION_NAME)
}
