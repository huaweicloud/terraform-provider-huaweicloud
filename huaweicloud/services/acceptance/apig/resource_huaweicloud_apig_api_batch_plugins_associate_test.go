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

func getApiBatchPluginsAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient(acceptance.HW_REGION_NAME, "apig")
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetLocalBoundPluginIdsForApi(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["api_id"], state.Primary.Attributes["env_id"],
		utils.ParseStateAttributeToListWithSeparator(state.Primary.Attributes["plugin_ids_origin"], ","))
}

func TestAccApiBatchPluginsAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		resourceNamePart1 = "huaweicloud_apig_api_batch_plugins_associate.part1"
		resourceNamePart2 = "huaweicloud_apig_api_batch_plugins_associate.part2"
		rcPart1           = acceptance.InitResourceCheck(resourceNamePart1, &obj, getApiBatchPluginsAssociateFunc)
		rcPart2           = acceptance.InitResourceCheck(resourceNamePart2, &obj, getApiBatchPluginsAssociateFunc)
		baseConfig        = testAccApiBatchPluginsAssociate_base()
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
				Config: testAccApiBatchPluginsAssociate_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNamePart1, "plugin_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceNamePart1, "plugin_ids_origin.#", "2"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNamePart2, "plugin_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceNamePart2, "plugin_ids_origin.#", "1"),
				),
			},
			{
				Config: testAccApiBatchPluginsAssociate_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					// After resources refreshed, the plugin_ids will be overridden as all plugins under the same
					// environment are bound.
					resource.TestCheckResourceAttr(resourceNamePart1, "plugin_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceNamePart1, "plugin_ids_origin.#", "2"),
					rcPart2.CheckResourceExists(),
					// After resources refreshed, the plugin_ids will be overridden as all plugins under the same
					// environment are bound.
					resource.TestCheckResourceAttr(resourceNamePart2, "plugin_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceNamePart2, "plugin_ids_origin.#", "1"),
				),
			},
			{
				Config: testAccApiBatchPluginsAssociate_basic_step3(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					// When multiple resources are used to manage the same function, plugin_ids will store the results
					// modified by other resources, resulting in plugin_ids displaying all binding results except for the
					// first change.
					resource.TestMatchResourceAttr(resourceNamePart1, "plugin_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(resourceNamePart1, "plugin_ids_origin.#", "1"),
					rcPart2.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceNamePart2, "plugin_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(resourceNamePart2, "plugin_ids_origin.#", "2"),
				),
			},
			{
				ResourceName:      resourceNamePart1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"plugin_ids",
				},
			},
			{
				// After resource part1 is imported, then plugin_ids will be overridden as all plugins under the same
				// environment are bound.
				Config: testAccApiBatchPluginsAssociate_basic_step3(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNamePart1, "plugin_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceNamePart1, "plugin_ids_origin.#", "1"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNamePart2, "plugin_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceNamePart2, "plugin_ids_origin.#", "2"),
				),
			},
		},
	})
}

func testAccApiBatchPluginsAssociate_base() string {
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
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "PATCH"
  request_path            = "/user_info"
  security_authentication = "APP"
  matching                = "Exact"

  web {
    path             = "/getUserAge"
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
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
}

resource "huaweicloud_apig_application" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}

resource "huaweicloud_apig_plugin" "http_response" {
  instance_id = local.instance_id
  name        = "%[2]s_http_response"
  type        = "set_resp_headers"
  content     = jsonencode(
    {
      response_headers = [{
        name       = "X-Custom-Pwd"
        value      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
        value_type = "custom_value"
        action     = "delete"
      },
      {
        name       = "X-Custom-Log-PATH"
        value      = "/tmp/debug.log"
        value_type = "custom_value"
        action     = "add"
      },
      {
        name       = "Sys-Param-updated"
        value      = "$context.cacheStatus"
        value_type = "system_parameter"
        action     = "append"
      }]
    }
  )
}

resource "huaweicloud_apig_plugin" "cors" {
  instance_id = local.instance_id
  name        = "%[2]s_cors"
  type        = "cors"
  content     = jsonencode(
    {
      allow_origin      = "*.terraform.test.com"
      allow_methods     = "POST,PATCH"
      allow_headers     = "Content-Type,Accept,Accept-Ranges"
      expose_headers    = "X-Request-Id,X-Apig-Auth-Type"
      max_age           = 800
      allow_credentials = false
    }
  )
}

resource "huaweicloud_apig_plugin" "third_auth" {
  instance_id = local.instance_id
  name        = "%[2]s_third_auth"
  type        = "third_auth"
  content     = jsonencode({
    "auth_downgrade_enabled": true,
    "auth_request": {
      "method": "GET",
      "path": "/test",
      "protocol": "HTTPS",
      "timeout": 5000,
      "url_domain": "",
      "vpc_channel_enabled": true,
      "vpc_channel_info": {
        "vpc_id": huaweicloud_apig_channel.test.id,
        "vpc_proxy_host": "test"
      }
    },
    "carry_body": {
      "enabled": true,
      "max_body_size": 1000
    }
    "carry_path_enabled": true,
    "carry_resp_headers": [],
    "custom_forbid_limit": 100,
    "identities": {
      "headers": [
        {
          "name": "X-Custom-Token"
        }
      ],
      "queries": null
    },
    "match_auth": null,
    "parameters": null,
    "rules": null,
    "return_resp_body_enabled": false,
    "rule_enabled":  false,
    "rule_type": "allow",
    "simple_auth_mode_enabled": true
  })
}

resource "huaweicloud_apig_plugin" "rate_limit" {
  instance_id = local.instance_id
  name        = "%[2]s_rate_limit"
  type        = "rate_limit"
  content     = jsonencode(
    {
      "scope": "basic",
      "default_time_unit": "minute",
      "default_interval": 2,
      "api_limit": 35,
      "app_limit": 15,
      "user_limit": 25,
      "ip_limit": 30,
      "algorithm": "haht",
      "specials": [
        {
          "type": "app",
          "policies": [
            {
              "key": huaweicloud_apig_application.test.id,
              "limit": 15
            }
          ]
        },
      ],
      "parameters": [
        {
          "type": "system",
          "name": "serverName",
          "value": "serverName"
        },
      ],
      "rules": [
        {
          "rule_name": "rule-0003",
          "match_regex": "[\"serverName\",\"==\",\"terraform\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 15
        },
      ]
    }
  )
}
`, common.TestBaseComputeResources(name), name,
		acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID,
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccApiBatchPluginsAssociate_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api_batch_plugins_associate" "part1" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  plugin_ids  = [
    huaweicloud_apig_plugin.http_response.id,
    huaweicloud_apig_plugin.third_auth.id,
  ]

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

resource "huaweicloud_apig_api_batch_plugins_associate" "part2" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  plugin_ids  = [
    huaweicloud_apig_plugin.cors.id,
  ]

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}
`, baseConfig)
}

func testAccApiBatchPluginsAssociate_basic_step2(baseConfig string) string {
	// Refresh the plugin_ids for all api batch plugins associate resources.
	return testAccApiBatchPluginsAssociate_basic_step1(baseConfig)
}

func testAccApiBatchPluginsAssociate_basic_step3(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api_batch_plugins_associate" "part1" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  plugin_ids  = [
    huaweicloud_apig_plugin.http_response.id,
  ]

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

resource "huaweicloud_apig_api_batch_plugins_associate" "part2" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  plugin_ids  = [
    huaweicloud_apig_plugin.cors.id,
    huaweicloud_apig_plugin.rate_limit.id,
  ]

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}
`, baseConfig)
}
