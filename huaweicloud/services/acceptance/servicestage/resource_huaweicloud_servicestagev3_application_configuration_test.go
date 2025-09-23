package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
)

func getV3ApplicationConfigurationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("servicestage", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}

	return servicestage.GetV3ApplicationConfiguration(client, state.Primary.Attributes["environment_id"],
		state.Primary.Attributes["application_id"])
}

func TestAccV3ApplicationConfiguration_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_servicestagev3_application_configuration.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3ApplicationConfigurationFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3ApplicationConfiguration_basic(basicV3ApplicationConfigurationEnvs, name, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id",
						"huaweicloud_servicestagev3_environment.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "application_id",
						"huaweicloud_servicestagev3_application.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.env.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.assign_strategy", "false"),
				),
			},
			{
				Config: testAccV3ApplicationConfiguration_basic(updateV3ApplicationConfigurationEnvs, name, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id",
						"huaweicloud_servicestagev3_environment.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "application_id",
						"huaweicloud_servicestagev3_application.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.env.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.assign_strategy", "true"),
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

const basicV3ApplicationConfigurationEnvs = `
variable "environment_variables" {
  type = list(object({
    name  = string
    value = string
  }))

  default = [
    {
      name  = "ENV_NAME_AND_START_WITH_A_UPPERCASE_ENGLISH_LETTER",
      value = "ENV_VALUE_AND_START_WITH_A_UPPERCASE_ENGLISH_LETTER",
    },
    {
      name  = "env_name_and_start_with_a_lowercase_english_letter",
      value = "env_value_and_start_with_a_lowercase_english_letter",
    },
    {
      name  = "-env_name_and_start_with_a_hyphen",
      value = "-env_value_and_start_with_a_hyphen",
    },
    {
      name  = "_env_name_and_start_with_a_underscore",
      value = "_env_value_and_start_with_a_underscore",
    },
    {
      name  = "env_name_and_end_with_a_dot_.",
      value = "._env_value_and_start_with_a_dot",
    },
    {
      name  = "env_name",
      value = "!@#$%^&*()+<>,/|\\~_env_value_and_start_with_some_special_characters",
    },
  ]
}`

const updateV3ApplicationConfigurationEnvs = `
variable "environment_variables" {
  type = list(object({
    name  = string
    value = string
  }))

  default = [
    {
      name  = "ENV_NAME_AND_START_WITH_A_UPPERCASE_ENGLISH_LETTER",
      value = "ENV_VALUE_AND_START_WITH_A_UPPERCASE_ENGLISH_LETTER",
    },
    {
      name  = "new_env_name_and_start_with_a_lowercase_english_letter",
      value = "env_value_and_start_with_a_lowercase_english_letter",
    },
    {
      name  = "-env_name_and_start_with_a_hyphen",
      value = "-new_env_value_and_start_with_a_hyphen",
    },
    {
      name  = "_new_env_name_and_start_with_a_underscore",
      value = "_new_env_value_and_start_with_a_underscore",
    },
    {
      name  = "env_name_and_end_with_a_digit_0",
      value = "0_env_value_and_start_with_a_digit",
    },
  ]
}`

func testAccV3ApplicationConfiguration_basic(envConfig, name string, assignStrategyEnabled bool) string {
	return fmt.Sprintf(`
%[1]s

variable "assign_strategy_enabled" {
  type    = bool
  default = %[2]v
}

resource "huaweicloud_vpc" "test" {
  name = "%[3]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_servicestagev3_application" "test" {
  name = "%[3]s"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}

resource "huaweicloud_servicestagev3_environment" "test" {
  name   = "%[3]s"
  vpc_id = huaweicloud_vpc.test.id
}

resource "huaweicloud_servicestagev3_application_configuration" "test" {
  application_id = huaweicloud_servicestagev3_application.test.id
  environment_id = huaweicloud_servicestagev3_environment.test.id
  
  configuration {
    dynamic "env" {
      for_each = var.environment_variables

      content {
        name  = env.value["name"]
        value = env.value["value"]
      }
    }
    assign_strategy = var.assign_strategy_enabled ? var.assign_strategy_enabled : null
  }
}
`, envConfig, assignStrategyEnabled, name)
}
