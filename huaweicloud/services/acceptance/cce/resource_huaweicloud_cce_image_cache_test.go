package cce

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getImageCacheFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getImageCacheHttpUrl = "v5/imagecaches/{image_cache_id}"
		getImageCacheProduct = "cce"
	)
	getImageCacheClient, err := cfg.NewServiceClient(getImageCacheProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE Client: %s", err)
	}

	getImageCachePath := getImageCacheClient.Endpoint + getImageCacheHttpUrl
	getImageCachePath = strings.ReplaceAll(getImageCachePath, "{image_cache_id}", state.Primary.ID)

	getImageCacheOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getImageCacheResp, err := getImageCacheClient.Request("GET", getImageCachePath, &getImageCacheOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE image cache: %s", err)
	}

	return utils.FlattenResponse(getImageCacheResp)
}

func TestAccImageCache_basic(t *testing.T) {
	var (
		imageCache   interface{}
		resourceName = "huaweicloud_cce_image_cache.test"
		rName        = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&imageCache,
			getImageCacheFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccImageCache_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "images.0",
						fmt.Sprintf("swr.%s.myhuaweicloud.com/xpanse/kafka:latest", acceptance.HW_REGION_NAME)),
					resource.TestCheckResourceAttrPair(resourceName, "building_config.0.cluster",
						"huaweicloud_cce_autopilot_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "building_config.0.image_pull_secrets.0",
						"default:default-secret"),
					resource.TestCheckResourceAttr(resourceName, "image_cache_size", "20"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "7"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
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

func testAccImageCache_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_autopilot_cluster" "test" {
  name        = "%[2]s"
  flavor      = "cce.autopilot.cluster"
  description = "created by terraform"

  host_network {
    vpc    = huaweicloud_vpc.test.id
    subnet = huaweicloud_vpc_subnet.test.id
  }

  container_network {
    mode = "eni"
  }

  eni_network {
    subnets {
      subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
    }
  }

  tags = {
    "foo" = "bar"
    "key" = "value"
  }
}
`, common.TestVpc(name), name)
}

func testAccImageCache_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_image_cache" "test" {
  name   = "%[2]s"
  images = ["swr.%[3]s.myhuaweicloud.com/xpanse/kafka:latest"]

  building_config {
    cluster            = huaweicloud_cce_autopilot_cluster.test.id
    image_pull_secrets = ["default:default-secret"]
  }

  image_cache_size = 20
  retention_days   = 7
}
`, testAccImageCache_base(name), name, acceptance.HW_REGION_NAME)
}
