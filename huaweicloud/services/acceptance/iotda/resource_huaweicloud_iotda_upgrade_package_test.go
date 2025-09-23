package iotda

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getUpgradePackageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/ota-upgrades/packages/{package_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{package_id}", state.Primary.ID)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IoTDA OTA upgrade package: %s", err)
	}

	return utils.FlattenResponse(getResp)
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
					resource.TestCheckResourceAttr(resourceName, "file_location.0.obs_location.0.sign",
						"0d6afb7e939f0936f40afdc759b5a354ea5427ec250a47e7b904ab1ea800a01d"),
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
	productBasic := testProduct_basic(name)

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
      sign        = "0d6afb7e939f0936f40afdc759b5a354ea5427ec250a47e7b904ab1ea800a01d"
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
`, productBasic, obsBucketBasic, acceptance.HW_REGION_NAME)
}
