package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPluginBatchApisAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	return apig.GetLocalBoundApiIdsForPlugin(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["plugin_id"], state.Primary.Attributes["env_id"],
		utils.ParseStateAttributeToListWithSeparator(state.Primary.Attributes["api_ids_origin"], ","))
}

func TestAccPluginBatchApisAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		resourceNamePart1 = "huaweicloud_apig_plugin_batch_apis_associate.part1"
		resourceNamePart2 = "huaweicloud_apig_plugin_batch_apis_associate.part2"
		rcPart1           = acceptance.InitResourceCheck(resourceNamePart1, &obj, getPluginBatchApisAssociateFunc)
		rcPart2           = acceptance.InitResourceCheck(resourceNamePart2, &obj, getPluginBatchApisAssociateFunc)
		baseConfig        = testAccPluginBatchApisAssociate_base()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcPart1.CheckResourceDestroy(),
			rcPart2.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPluginBatchApisAssociate_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNamePart1, "api_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceNamePart1, "api_ids_origin.#", "2"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNamePart2, "api_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceNamePart2, "api_ids_origin.#", "1"),
				),
			},
			{
				Config: testAccPluginBatchApisAssociate_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					// After resources refreshed, the api_ids will be overridden as all APIs under the same
					// environment are bound.
					resource.TestCheckResourceAttr(resourceNamePart1, "api_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceNamePart1, "api_ids_origin.#", "2"),
					rcPart2.CheckResourceExists(),
					// After resources refreshed, the api_ids will be overridden as all APIs under the same
					// environment are bound.
					resource.TestCheckResourceAttr(resourceNamePart2, "api_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceNamePart2, "api_ids_origin.#", "1"),
				),
			},
			{
				Config: testAccPluginBatchApisAssociate_basic_step3(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					// When multiple resources are used to manage the same function, api_ids will store the results
					// modified by other resources, resulting in api_ids displaying all binding results except for the
					// first change.
					resource.TestMatchResourceAttr(resourceNamePart1, "api_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(resourceNamePart1, "api_ids_origin.#", "1"),
					rcPart2.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceNamePart2, "api_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(resourceNamePart2, "api_ids_origin.#", "2"),
				),
			},
			{
				ResourceName:      resourceNamePart1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"api_ids",
				},
			},
			{
				// After resource part1 is imported, then api_ids will be overridden as all APIs under the same
				// environment are bound.
				Config: testAccPluginBatchApisAssociate_basic_step3(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNamePart1, "api_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceNamePart1, "api_ids_origin.#", "1"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNamePart2, "api_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceNamePart2, "api_ids_origin.#", "2"),
				),
			},
		},
	})
}

func testAccPluginBatchApisAssociate_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[3]s"
  }
}

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[4]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_group" "test" {
  name        = "%[2]s"
  instance_id = local.instance_id
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  port             = 8000
  balance_strategy = 2
  member_type      = "ecs"
  type             = 2

  health_check {
    protocol           = "HTTPS"
    threshold_normal   = 10  # maximum value
    threshold_abnormal = 10  # maximum value
    interval           = 300 # maximum value
    timeout            = 30  # maximum value
    path               = "/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  member {
    id   = huaweicloud_compute_instance.test.id
    name = huaweicloud_compute_instance.test.name
  }
}

resource "huaweicloud_apig_api" "test" {
  count = 4

  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s_${count.index}"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/${count.index}"
  security_authentication = "APP"
  matching                = "Exact"

  web {
    path             = "/getUserAge/${count.index}"
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
  }
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}

resource "huaweicloud_apig_api_publishment" "test" {
  count = 4

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test[count.index].id
  env_id      = huaweicloud_apig_environment.test.id
}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
  type        = "cors"
  content     = jsonencode(
    {
      allow_origin      = "*"
      allow_methods     = "GET,PUT,DELETE,HEAD,PATCH"
      allow_headers     = "Content-Type,Accept,Cache-Control"
      expose_headers    = "X-Request-Id,X-Apig-Latency"
      max_age           = 12700
      allow_credentials = true
    }
  )
}
`, common.TestBaseComputeResources(name), name,
		acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID,
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccPluginBatchApisAssociate_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_plugin_batch_apis_associate" "part1" {
  depends_on = [huaweicloud_apig_api_publishment.test]

  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_ids     = slice(huaweicloud_apig_api.test[*].id, 0, 2)
}

resource "huaweicloud_apig_plugin_batch_apis_associate" "part2" {
  depends_on = [huaweicloud_apig_api_publishment.test]

  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_ids     = slice(huaweicloud_apig_api.test[*].id, 3, 4)
}
`, baseConfig)
}

func testAccPluginBatchApisAssociate_basic_step2(baseConfig string) string {
	// Refresh the api_ids for all plugin associate resources.
	return testAccPluginBatchApisAssociate_basic_step1(baseConfig)
}

func testAccPluginBatchApisAssociate_basic_step3(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_plugin_batch_apis_associate" "part1" {
  depends_on = [huaweicloud_apig_api_publishment.test]

  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_ids     = slice(huaweicloud_apig_api.test[*].id, 0, 1)
}

resource "huaweicloud_apig_plugin_batch_apis_associate" "part2" {
  depends_on = [huaweicloud_apig_api_publishment.test]

  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_ids     = slice(huaweicloud_apig_api.test[*].id, 2, 4)
}
`, baseConfig)
}
