package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
)

func getFunction(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("fgs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	return fgs.GetFunctionMetadata(client, state.Primary.ID)
}

func TestAccFunction_basic(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceNameWithDash()

		withBase64Code   = "huaweicloud_fgs_function.with_base64_code"
		rcWithBase64Code = acceptance.InitResourceCheck(withBase64Code, &obj, getFunction)

		withTextCode   = "huaweicloud_fgs_function.with_text_code"
		rcWithTextCode = acceptance.InitResourceCheck(withTextCode, &obj, getFunction)

		withObsStorage   = "huaweicloud_fgs_function.with_obs_storage"
		rcWithObsStorage = acceptance.InitResourceCheck(withObsStorage, &obj, getFunction)

		withCustomImage   = "huaweicloud_fgs_function.with_custom_image"
		rcWithCustomImage = acceptance.InitResourceCheck(withCustomImage, &obj, getFunction)

		withDeprecatedParams   = "huaweicloud_fgs_function.with_deprecated_params"
		rcWithDeprecatedParams = acceptance.InitResourceCheck(withDeprecatedParams, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsAgency(t)
			acceptance.TestAccPreCheckFgsAppAgency(t)
			acceptance.TestAccPreCheckImageUrl(t)
			acceptance.TestAccPreCheckImageUrlUpdated(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithBase64Code.CheckResourceDestroy(),
			rcWithTextCode.CheckResourceDestroy(),
			rcWithObsStorage.CheckResourceDestroy(),
			rcWithCustomImage.CheckResourceDestroy(),
			rcWithDeprecatedParams.CheckResourceDestroy(),
		),

		Steps: []resource.TestStep{
			{
				Config: testAccFunction_basic_step1(name, "Python2.7"),
				Check: resource.ComposeTestCheckFunc(
					// Check the function which the code context is base64 encoded.
					rcWithBase64Code.CheckResourceExists(),
					resource.TestCheckResourceAttr(withBase64Code, "name", name+"-base64-code"),
					resource.TestCheckResourceAttr(withBase64Code, "memory_size", "128"),
					resource.TestCheckResourceAttr(withBase64Code, "runtime", "Python2.7"),
					resource.TestCheckResourceAttr(withBase64Code, "timeout", "3"),
					resource.TestCheckResourceAttr(withBase64Code, "app", "default"),
					resource.TestCheckResourceAttr(withBase64Code, "handler", "index.handler"),
					resource.TestCheckResourceAttr(withBase64Code, "code_type", "inline"),
					resource.TestCheckResourceAttr(withBase64Code, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(withBase64Code, "concurrency_num", "1"),
					resource.TestCheckResourceAttr(withBase64Code, "depend_list.#", "1"),
					resource.TestCheckResourceAttrPair(withBase64Code, "depend_list.0",
						"huaweicloud_fgs_dependency_version.test", "version_id"),
					resource.TestCheckResourceAttr(withBase64Code, "user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withBase64Code, "encrypted_user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withBase64Code, "tags.%", "2"),
					resource.TestCheckResourceAttr(withBase64Code, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(withBase64Code, "tags.key", "value"),
					resource.TestCheckResourceAttr(withBase64Code, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttr(withBase64Code, "enable_dynamic_memory", "true"),
					resource.TestCheckResourceAttr(withBase64Code, "is_stateful_function", "true"),
					resource.TestCheckResourceAttr(withBase64Code, "initializer_handler", "index.handler"),
					resource.TestCheckResourceAttr(withBase64Code, "initializer_timeout", "5"),
					resource.TestCheckResourceAttr(withBase64Code, "network_controller.#", "1"),
					resource.TestCheckResourceAttr(withBase64Code, "network_controller.0.trigger_access_vpcs.#", "2"),
					resource.TestCheckResourceAttr(withBase64Code, "network_controller.0.disable_public_network", "true"),
					resource.TestCheckResourceAttrSet(withBase64Code, "urn"),
					resource.TestCheckResourceAttrSet(withBase64Code, "version"),
					// Check the function which the code context is not base64 encoded.
					rcWithTextCode.CheckResourceExists(),
					resource.TestCheckResourceAttr(withTextCode, "name", name+"-text-code"),
					resource.TestCheckResourceAttr(withTextCode, "memory_size", "128"),
					resource.TestCheckResourceAttr(withTextCode, "runtime", "Python2.7"),
					resource.TestCheckResourceAttr(withTextCode, "timeout", "3"),
					resource.TestCheckResourceAttr(withTextCode, "app", "default"),
					resource.TestCheckResourceAttr(withTextCode, "handler", "index.handler"),
					resource.TestCheckResourceAttr(withTextCode, "code_type", "inline"),
					resource.TestCheckResourceAttr(withTextCode, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(withTextCode, "concurrency_num", "1"),
					resource.TestCheckResourceAttr(withTextCode, "depend_list.#", "1"),
					resource.TestCheckResourceAttrPair(withTextCode, "depend_list.0",
						"huaweicloud_fgs_dependency_version.test", "version_id"),
					resource.TestCheckResourceAttr(withTextCode, "user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withTextCode, "encrypted_user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withTextCode, "tags.%", "2"),
					resource.TestCheckResourceAttr(withTextCode, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(withTextCode, "tags.key", "value"),
					resource.TestCheckResourceAttr(withTextCode, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttr(withTextCode, "enable_dynamic_memory", "true"),
					resource.TestCheckResourceAttr(withTextCode, "is_stateful_function", "true"),
					resource.TestCheckResourceAttr(withTextCode, "network_controller.#", "1"),
					resource.TestCheckResourceAttr(withTextCode, "network_controller.0.trigger_access_vpcs.#", "2"),
					resource.TestCheckResourceAttr(withTextCode, "network_controller.0.disable_public_network", "true"),
					resource.TestCheckResourceAttrSet(withTextCode, "urn"),
					resource.TestCheckResourceAttrSet(withTextCode, "version"),
					// Check the function which the code file is storaged in the OBS bucket.
					rcWithObsStorage.CheckResourceExists(),
					resource.TestCheckResourceAttr(withObsStorage, "name", name+"-obs-storage"),
					resource.TestCheckResourceAttr(withObsStorage, "memory_size", "128"),
					resource.TestCheckResourceAttr(withObsStorage, "runtime", "Python2.7"),
					resource.TestCheckResourceAttr(withObsStorage, "timeout", "3"),
					resource.TestCheckResourceAttr(withObsStorage, "app", "default"),
					resource.TestCheckResourceAttr(withObsStorage, "handler", "index.handler"),
					resource.TestCheckResourceAttr(withObsStorage, "code_type", "obs"),
					resource.TestCheckResourceAttr(withObsStorage, "agency", acceptance.HW_FGS_AGENCY_NAME),
					resource.TestCheckResourceAttr(withObsStorage, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(withObsStorage, "concurrency_num", "1"),
					resource.TestCheckResourceAttr(withObsStorage, "depend_list.#", "1"),
					resource.TestCheckResourceAttrPair(withObsStorage, "depend_list.0",
						"huaweicloud_fgs_dependency_version.test", "version_id"),
					resource.TestCheckResourceAttr(withObsStorage, "user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withObsStorage, "encrypted_user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.%", "2"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.key", "value"),
					resource.TestCheckResourceAttr(withObsStorage, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttr(withObsStorage, "enable_dynamic_memory", "true"),
					resource.TestCheckResourceAttr(withObsStorage, "is_stateful_function", "true"),
					resource.TestCheckResourceAttr(withObsStorage, "network_controller.#", "1"),
					resource.TestCheckResourceAttr(withObsStorage, "network_controller.0.trigger_access_vpcs.#", "2"),
					resource.TestCheckResourceAttr(withObsStorage, "network_controller.0.disable_public_network", "true"),
					resource.TestCheckResourceAttrSet(withObsStorage, "urn"),
					resource.TestCheckResourceAttrSet(withObsStorage, "version"),
					// Check the function which is build via an SWR image.
					rcWithCustomImage.CheckResourceExists(),
					resource.TestCheckResourceAttr(withCustomImage, "name", name+"-custom-image"),
					resource.TestCheckResourceAttr(withCustomImage, "memory_size", "128"),
					resource.TestCheckResourceAttr(withCustomImage, "runtime", "Custom Image"),
					resource.TestCheckResourceAttr(withCustomImage, "timeout", "3"),
					resource.TestCheckResourceAttr(withCustomImage, "app", "default"),
					resource.TestCheckResourceAttr(withCustomImage, "handler", "-"),
					resource.TestCheckResourceAttr(withCustomImage, "code_type", "Custom-Image-Swr"),
					resource.TestCheckResourceAttr(withCustomImage, "agency", acceptance.HW_FGS_AGENCY_NAME),
					resource.TestCheckResourceAttr(withCustomImage, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrPair(withCustomImage, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(withCustomImage, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(withCustomImage, "enable_auth_in_header", "false"),
					resource.TestCheckResourceAttr(withCustomImage, "concurrency_num", "1"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.#", "1"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.url", acceptance.HW_BUILD_IMAGE_URL),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.command", ""),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.args", ""),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.working_dir", "/"),
					resource.TestCheckResourceAttr(withCustomImage, "user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withCustomImage, "encrypted_user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withCustomImage, "tags.%", "2"),
					resource.TestCheckResourceAttr(withCustomImage, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(withCustomImage, "tags.key", "value"),
					resource.TestCheckResourceAttr(withCustomImage, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttr(withCustomImage, "enable_dynamic_memory", "true"),
					resource.TestCheckResourceAttr(withCustomImage, "is_stateful_function", "true"),
					resource.TestCheckResourceAttr(withCustomImage, "network_controller.#", "1"),
					resource.TestCheckResourceAttr(withCustomImage, "network_controller.0.trigger_access_vpcs.#", "2"),
					resource.TestCheckResourceAttr(withCustomImage, "network_controller.0.disable_public_network", "true"),
					resource.TestCheckResourceAttrSet(withCustomImage, "urn"),
					resource.TestCheckResourceAttrSet(withCustomImage, "version"),
					rcWithDeprecatedParams.CheckResourceExists(),
					resource.TestCheckResourceAttr(withDeprecatedParams, "package", "default"),
					resource.TestCheckResourceAttr(withDeprecatedParams, "xrole", acceptance.HW_FGS_AGENCY_NAME),
				),
			},
			{
				ResourceName:      withDeprecatedParams,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
					"app", // Recommand parameter setting by default.
					"package",
					"agency", // Recommand parameter setting by default.
					"xrole",
				},
			},
			{
				Config: testAccFunction_basic_step2(name, "Python3.6"),
				Check: resource.ComposeTestCheckFunc(
					// Check the function which the code context is base64 encoded.
					rcWithBase64Code.CheckResourceExists(),
					resource.TestCheckResourceAttr(withBase64Code, "memory_size", "256"),
					resource.TestCheckResourceAttr(withBase64Code, "runtime", "Python3.6"),
					resource.TestCheckResourceAttr(withBase64Code, "timeout", "5"),
					resource.TestCheckResourceAttr(withBase64Code, "app", "default"),
					resource.TestCheckResourceAttr(withBase64Code, "handler", "index.handler"),
					resource.TestCheckResourceAttr(withBase64Code, "code_type", "inline"),
					resource.TestCheckResourceAttr(withBase64Code, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(withBase64Code, "concurrency_num", "500"),
					resource.TestCheckResourceAttr(withBase64Code, "depend_list.#", "1"),
					resource.TestCheckResourceAttrPair(withBase64Code, "depend_list.0",
						"huaweicloud_fgs_dependency_version.test", "version_id"),
					resource.TestCheckResourceAttr(withBase64Code, "user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withBase64Code, "encrypted_user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withBase64Code, "tags.%", "2"),
					resource.TestCheckResourceAttr(withBase64Code, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(withBase64Code, "tags.new_key", "value"),
					resource.TestCheckResourceAttr(withBase64Code, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttr(withBase64Code, "enable_dynamic_memory", "false"),
					resource.TestCheckResourceAttr(withBase64Code, "is_stateful_function", "false"),
					resource.TestCheckResourceAttr(withBase64Code, "network_controller.#", "1"),
					resource.TestCheckResourceAttr(withBase64Code, "network_controller.0.trigger_access_vpcs.#", "2"),
					resource.TestCheckResourceAttr(withBase64Code, "network_controller.0.disable_public_network", "false"),
					resource.TestCheckResourceAttrSet(withBase64Code, "urn"),
					resource.TestCheckResourceAttrSet(withBase64Code, "version"),
					// Check the function which the code context is not base64 encoded.
					rcWithTextCode.CheckResourceExists(),
					resource.TestCheckResourceAttr(withTextCode, "memory_size", "256"),
					resource.TestCheckResourceAttr(withTextCode, "runtime", "Python3.6"),
					resource.TestCheckResourceAttr(withTextCode, "timeout", "5"),
					resource.TestCheckResourceAttr(withTextCode, "app", "default"),
					resource.TestCheckResourceAttr(withTextCode, "handler", "index.handler"),
					resource.TestCheckResourceAttr(withTextCode, "code_type", "inline"),
					resource.TestCheckResourceAttr(withTextCode, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(withTextCode, "concurrency_num", "500"),
					resource.TestCheckResourceAttr(withTextCode, "depend_list.#", "1"),
					resource.TestCheckResourceAttrPair(withTextCode, "depend_list.0",
						"huaweicloud_fgs_dependency_version.test", "version_id"),
					resource.TestCheckResourceAttr(withTextCode, "user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withTextCode, "encrypted_user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withTextCode, "tags.%", "2"),
					resource.TestCheckResourceAttr(withTextCode, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(withTextCode, "tags.new_key", "value"),
					resource.TestCheckResourceAttr(withTextCode, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttr(withTextCode, "enable_dynamic_memory", "false"),
					resource.TestCheckResourceAttr(withTextCode, "is_stateful_function", "false"),
					resource.TestCheckResourceAttr(withTextCode, "network_controller.#", "1"),
					resource.TestCheckResourceAttr(withTextCode, "network_controller.0.trigger_access_vpcs.#", "2"),
					resource.TestCheckResourceAttr(withTextCode, "network_controller.0.disable_public_network", "false"),
					resource.TestCheckResourceAttrSet(withTextCode, "urn"),
					resource.TestCheckResourceAttrSet(withTextCode, "version"),
					// Check the function which the code file is storaged in the OBS bucket.
					rcWithObsStorage.CheckResourceExists(),
					resource.TestCheckResourceAttr(withObsStorage, "memory_size", "256"),
					resource.TestCheckResourceAttr(withObsStorage, "runtime", "Python3.6"),
					resource.TestCheckResourceAttr(withObsStorage, "timeout", "5"),
					resource.TestCheckResourceAttr(withObsStorage, "app", "default"),
					resource.TestCheckResourceAttr(withObsStorage, "handler", "index.handler"),
					resource.TestCheckResourceAttr(withObsStorage, "code_type", "obs"),
					resource.TestCheckResourceAttr(withObsStorage, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(withObsStorage, "concurrency_num", "500"),
					resource.TestCheckResourceAttr(withObsStorage, "depend_list.#", "1"),
					resource.TestCheckResourceAttrPair(withObsStorage, "depend_list.0",
						"huaweicloud_fgs_dependency_version.test", "version_id"),
					resource.TestCheckResourceAttr(withObsStorage, "user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withObsStorage, "encrypted_user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.%", "2"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.new_key", "value"),
					resource.TestCheckResourceAttr(withObsStorage, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttr(withObsStorage, "enable_dynamic_memory", "false"),
					resource.TestCheckResourceAttr(withObsStorage, "is_stateful_function", "false"),
					resource.TestCheckResourceAttr(withObsStorage, "network_controller.#", "1"),
					resource.TestCheckResourceAttr(withObsStorage, "network_controller.0.trigger_access_vpcs.#", "2"),
					resource.TestCheckResourceAttr(withObsStorage, "network_controller.0.disable_public_network", "false"),
					resource.TestCheckResourceAttrSet(withObsStorage, "urn"),
					resource.TestCheckResourceAttrSet(withObsStorage, "version"),
					// Check the function which is build via an SWR image.
					rcWithCustomImage.CheckResourceExists(),
					resource.TestCheckResourceAttr(withCustomImage, "memory_size", "256"),
					resource.TestCheckResourceAttr(withCustomImage, "runtime", "Custom Image"),
					resource.TestCheckResourceAttr(withCustomImage, "timeout", "5"),
					resource.TestCheckResourceAttr(withCustomImage, "app", "default"),
					resource.TestCheckResourceAttr(withCustomImage, "app_agency", acceptance.HW_FGS_APP_AGENCY_NAME),
					resource.TestCheckResourceAttr(withCustomImage, "handler", "-"),
					resource.TestCheckResourceAttr(withCustomImage, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttrPair(withCustomImage, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(withCustomImage, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(withCustomImage, "enable_auth_in_header", "true"),
					resource.TestCheckResourceAttr(withCustomImage, "concurrency_num", "500"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.#", "1"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.url", acceptance.HW_BUILD_IMAGE_URL_UPDATED),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.command", "/bin/sh"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.args", "-args,value"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.working_dir", "/"),
					resource.TestCheckResourceAttr(withCustomImage, "user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withCustomImage, "encrypted_user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withCustomImage, "tags.%", "2"),
					resource.TestCheckResourceAttr(withCustomImage, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(withCustomImage, "tags.new_key", "value"),
					resource.TestCheckResourceAttr(withCustomImage, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttr(withCustomImage, "enable_dynamic_memory", "false"),
					resource.TestCheckResourceAttr(withCustomImage, "is_stateful_function", "false"),
					resource.TestCheckResourceAttr(withCustomImage, "network_controller.#", "1"),
					resource.TestCheckResourceAttr(withCustomImage, "network_controller.0.trigger_access_vpcs.#", "2"),
					resource.TestCheckResourceAttr(withCustomImage, "network_controller.0.disable_public_network", "false"),
					resource.TestCheckResourceAttrSet(withCustomImage, "urn"),
					resource.TestCheckResourceAttrSet(withCustomImage, "version"),
				),
			},
			{
				ResourceName:      withBase64Code,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
					"tags",
					"encrypted_user_data",
				},
			},
			{
				ResourceName:      withTextCode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
					"tags",
					"encrypted_user_data",
				},
			},
			{
				ResourceName:      withObsStorage,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
					"tags",
					"encrypted_user_data",
				},
			},
			{
				ResourceName:      withCustomImage,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
					"tags",
					"encrypted_user_data",
				},
			},
		},
	})
}

const functionScriptVariableDefinition = `
variable "script_content" {
  type    = string
  default = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
}
`

func zipFileUploadResourcesConfig(name, runtime string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc" "test" {
  name = "%[2]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[2]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0), 1)
}

resource "huaweicloud_obs_bucket" "test" {
  bucket = "%[2]s"
  acl    = "private"

  provisioner "local-exec" {
    command = "echo '${var.script_content}' >> test.py\nzip -r test.zip test.py"
  }
  provisioner "local-exec" {
    command = "rm test.zip test.py"
    when    = destroy
  }
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket = huaweicloud_obs_bucket.test.bucket
  key    = "test.zip"
  source = abspath("./test.zip")
}

resource "huaweicloud_fgs_dependency_version" "test" {
  name    = "%[2]s"
  runtime = "%[3]s"
  link    = "%[4]s"
}
`, functionScriptVariableDefinition, name, runtime,
		acceptance.HW_FGS_DEPENDENCY_OBS_LINK)
}

func testAccFunction_basic_step1(name, runtime string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc" "trigger_access" {
  count = 3

  name = format("%[2]s-trigger-%%d", count.index)
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_fgs_function" "with_base64_code" {
  name                  = "%[2]s-base64-code"
  memory_size           = 128
  runtime               = "%[3]s"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  enable_dynamic_memory = true
  is_stateful_function  = true
  initializer_handler   = "index.handler"
  initializer_timeout   = 5

  user_data = jsonencode({
    "owner": "terraform"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform"
  })

  network_controller {
    disable_public_network = true

    dynamic "trigger_access_vpcs" {
      for_each = slice(huaweicloud_vpc.trigger_access[*].id, 0, 2)

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }

  # Update the runtime value will trigger the dependency version change.
  depend_list = [huaweicloud_fgs_dependency_version.test.version_id]
  
  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_fgs_function" "with_text_code" {
  name                  = "%[2]s-text-code"
  memory_size           = 128
  runtime               = "%[3]s"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = var.script_content
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  enable_dynamic_memory = true
  is_stateful_function  = true

  user_data = jsonencode({
    "owner": "terraform"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform"
  })

  network_controller {
    disable_public_network = true

    dynamic "trigger_access_vpcs" {
      for_each = slice(huaweicloud_vpc.trigger_access[*].id, 0, 2)

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }

  # Update the runtime value will trigger the dependency version change.
  depend_list = [huaweicloud_fgs_dependency_version.test.version_id]

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_fgs_function" "with_obs_storage" {
  name                  = "%[2]s-obs-storage"
  memory_size           = 128
  runtime               = "%[3]s"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "obs"
  code_url              = format("https://%%s/%%s", huaweicloud_obs_bucket.test.bucket_domain_name, huaweicloud_obs_bucket_object.test.key)
  agency                = "%[4]s"
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  enable_dynamic_memory = true
  is_stateful_function  = true

  user_data = jsonencode({
    "owner": "terraform"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform"
  })

  network_controller {
    disable_public_network = true

    dynamic "trigger_access_vpcs" {
      for_each = slice(huaweicloud_vpc.trigger_access[*].id, 0, 2)

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }

  # Update the runtime value will trigger the dependency version change.
  depend_list = [huaweicloud_fgs_dependency_version.test.version_id]

  tags = {
    foo = "bar"
    key = "value"
  }
}

# The dependencies are already packaged in the custom image, do not specifies the depend_list parameter.
resource "huaweicloud_fgs_function" "with_custom_image" {
  name                  = "%[2]s-custom-image"
  memory_size           = 128
  runtime               = "Custom Image"
  timeout               = 3
  app                   = "default"
  handler               = "-"
  agency                = "%[4]s"
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  enable_dynamic_memory = true
  is_stateful_function  = true
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id

  custom_image {
    url = "%[5]s"
  }

  user_data = jsonencode({
    "owner": "terraform"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform"
  })

  network_controller {
    disable_public_network = true

    dynamic "trigger_access_vpcs" {
      for_each = slice(huaweicloud_vpc.trigger_access[*].id, 0, 2)

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_fgs_function" "with_deprecated_params" {
  name        = "%[2]s-deprecated-params"
  memory_size = 128
  runtime     = "%[3]s"
  timeout     = 3
  package     = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  xrole       = "%[4]s"
}
`, zipFileUploadResourcesConfig(name, runtime), name, runtime,
		acceptance.HW_FGS_AGENCY_NAME,
		acceptance.HW_BUILD_IMAGE_URL)
}

func testAccFunction_basic_step2(name, runtime string) string {
	//nolint:revive
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc" "trigger_access" {
  count = 3

  name = format("%[2]s-trigger-%%d", count.index)
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_fgs_function" "with_base64_code" {
  name                  = "%[2]s-base64-code"
  memory_size           = 256
  runtime               = "%[3]s"
  timeout               = 5
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Updated by terraform script"
  functiongraph_version = "v2"
  concurrency_num       = 500

  user_data = jsonencode({
    "owner": "terraform",
    "usage": "acceptance test"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform",
    "usage": "acceptance test"
  })

  network_controller {
    dynamic "trigger_access_vpcs" {
      for_each = slice(huaweicloud_vpc.trigger_access[*].id, 1, 3)

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }

  # Update the runtime value will trigger the dependency version change.
  depend_list = [huaweicloud_fgs_dependency_version.test.version_id]
  
  tags = {
    foo     = "baar"
    new_key = "value"
  }
}

resource "huaweicloud_fgs_function" "with_text_code" {
  name                  = "%[2]s-text-code"
  memory_size           = 256
  runtime               = "%[3]s"
  timeout               = 5
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = var.script_content
  description           = "Updated by terraform script"
  functiongraph_version = "v2"
  concurrency_num       = 500

  user_data = jsonencode({
    "owner": "terraform",
    "usage": "acceptance test"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform",
    "usage": "acceptance test"
  })

  network_controller {
    dynamic "trigger_access_vpcs" {
      for_each = slice(huaweicloud_vpc.trigger_access[*].id, 1, 3)

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }

  # Update the runtime value will trigger the dependency version change.
  depend_list = [huaweicloud_fgs_dependency_version.test.version_id]

  tags = {
    foo     = "baar"
    new_key = "value"
  }
}

resource "huaweicloud_fgs_function" "with_obs_storage" {
  name                  = "%[2]s-obs-storage"
  memory_size           = 256
  runtime               = "%[3]s"
  timeout               = 5
  app                   = "default"
  handler               = "index.handler"
  code_type             = "obs"
  code_url              = format("https://%%s/%%s", huaweicloud_obs_bucket.test.bucket_domain_name, huaweicloud_obs_bucket_object.test.key)
  agency                = "%[4]s"
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
  description           = "Updated by terraform script"
  functiongraph_version = "v2"
  concurrency_num       = 500

  user_data = jsonencode({
    "owner": "terraform",
    "usage": "acceptance test"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform",
    "usage": "acceptance test"
  })

  network_controller {
    dynamic "trigger_access_vpcs" {
      for_each = slice(huaweicloud_vpc.trigger_access[*].id, 1, 3)

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }

  # Update the runtime value will trigger the dependency version change.
  depend_list = [huaweicloud_fgs_dependency_version.test.version_id]

  tags = {
    foo     = "baar"
    new_key = "value"
  }
}

# The dependencies are already packaged in the custom image, do not specifies the depend_list parameter.
resource "huaweicloud_fgs_function" "with_custom_image" {
  name                  = "%[2]s-custom-image"
  memory_size           = 256
  runtime               = "Custom Image"
  timeout               = 5
  app                   = "default"
  handler               = "-"
  agency                = "%[4]s"
  app_agency            = "%[5]s"
  description           = "Updated by terraform script"
  functiongraph_version = "v2"
  concurrency_num       = 500
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
  enable_auth_in_header = true

  custom_image {
    url         = "%[6]s"
    command     = "/bin/sh"
    args        = "-args,value"
    working_dir = "/"
  }

  user_data = jsonencode({
    "owner": "terraform",
    "usage": "acceptance test"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform",
    "usage": "acceptance test"
  })

  network_controller {
    dynamic "trigger_access_vpcs" {
      for_each = slice(huaweicloud_vpc.trigger_access[*].id, 1, 3)

      content {
        vpc_id = trigger_access_vpcs.value
      }
    }
  }

  tags = {
    foo     = "baar"
    new_key = "value"
  }
}
`, zipFileUploadResourcesConfig(name, runtime), name, runtime,
		acceptance.HW_FGS_AGENCY_NAME,
		acceptance.HW_FGS_APP_AGENCY_NAME,
		acceptance.HW_BUILD_IMAGE_URL_UPDATED)
}

func TestAccFunction_withEpsId(t *testing.T) {
	var (
		obj interface{}

		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_function.test"

		rc = acceptance.InitResourceCheck(resourceName, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_withEpsId(name, "0"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					// Default value is v2. Some regions support only v1, the default value is v1.
					resource.TestMatchResourceAttr(resourceName, "functiongraph_version", regexp.MustCompile(`v1|v2`)),
				),
			},
			{
				Config: testAccFunction_withEpsId(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_withEpsId(name, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name                  = "%[2]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  enterprise_project_id = "%[3]s"
  description           = "Created by terraform script"
}
`, functionScriptVariableDefinition, name, epsId)
}

func TestAccFunction_logConfig(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceName()

		createWithLtsParams      = "huaweicloud_fgs_function.create_with_lts_params"
		rcCreateWithLtsParams    = acceptance.InitResourceCheck(createWithLtsParams, &obj, getFunction)
		createWithoutLtsParams   = "huaweicloud_fgs_function.create_without_lts_params"
		rcCreateWithoutLtsParams = acceptance.InitResourceCheck(createWithoutLtsParams, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcCreateWithLtsParams.CheckResourceDestroy(),
			rcCreateWithoutLtsParams.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_logConfig_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcCreateWithLtsParams.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithLtsParams, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_group_id",
						"huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_group_name",
						"huaweicloud_lts_group.test.0", "group_name"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_stream_id",
						"huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_stream_name",
						"huaweicloud_lts_stream.test.0", "stream_name"),
					resource.TestCheckResourceAttr(createWithLtsParams, "lts_custom_tag.%", "2"),
					resource.TestCheckResourceAttr(createWithLtsParams, "lts_custom_tag.foo", "bar"),
					resource.TestCheckResourceAttr(createWithLtsParams, "lts_custom_tag.key", "value"),
					rcCreateWithoutLtsParams.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithoutLtsParams, "functiongraph_version", "v2"),
					// In some regions (such as 'cn-north-4'), the FunctionGraph service automatically binds the groups
					// and streams created by FunctionGraph to functions that do not have LTS set.
					resource.TestCheckResourceAttrSet(createWithoutLtsParams, "log_group_id"),
					resource.TestCheckNoResourceAttr(createWithoutLtsParams, "log_group_name"),
					resource.TestCheckResourceAttrSet(createWithoutLtsParams, "log_stream_id"),
					resource.TestCheckNoResourceAttr(createWithoutLtsParams, "log_stream_name"),
				),
			},
			{
				Config: testAccFunction_logConfig_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcCreateWithLtsParams.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithLtsParams, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_group_id",
						"huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_group_name",
						"huaweicloud_lts_group.test.1", "group_name"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_stream_id",
						"huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_stream_name",
						"huaweicloud_lts_stream.test.1", "stream_name"),
					resource.TestCheckResourceAttr(createWithLtsParams, "lts_custom_tag.%", "2"),
					resource.TestCheckResourceAttr(createWithLtsParams, "lts_custom_tag.foo", "baar"),
					resource.TestCheckResourceAttr(createWithLtsParams, "lts_custom_tag.new_key", "value"),
					rcCreateWithoutLtsParams.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithoutLtsParams, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttrPair(createWithoutLtsParams, "log_group_id",
						"huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(createWithoutLtsParams, "log_group_name",
						"huaweicloud_lts_group.test.0", "group_name"),
					resource.TestCheckResourceAttrPair(createWithoutLtsParams, "log_stream_id",
						"huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttrPair(createWithoutLtsParams, "log_stream_name",
						"huaweicloud_lts_stream.test.0", "stream_name"),
					resource.TestCheckResourceAttr(createWithoutLtsParams, "lts_custom_tag.%", "2"),
					resource.TestCheckResourceAttr(createWithoutLtsParams, "lts_custom_tag.foo", "bar"),
					resource.TestCheckResourceAttr(createWithoutLtsParams, "lts_custom_tag.key", "value"),
				),
			},
			{
				Config: testAccFunction_logConfig_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rcCreateWithLtsParams.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithLtsParams, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_group_id",
						"huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_group_name",
						"huaweicloud_lts_group.test.1", "group_name"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_stream_id",
						"huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttrPair(createWithLtsParams, "log_stream_name",
						"huaweicloud_lts_stream.test.1", "stream_name"),
					resource.TestCheckResourceAttr(createWithLtsParams, "lts_custom_tag.%", "0"),
				),
			},
		},
	})
}

func testAccFunction_logConfig_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  count = 2

  group_name  = format("%[2]s_%%d", count.index)
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  count = 2

  group_id    = huaweicloud_lts_group.test[count.index].id
  stream_name = format("%[2]s_%%d", count.index)
}
`, functionScriptVariableDefinition, name)
}

func testAccFunction_logConfig_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "create_with_lts_params" {
  name                  = "%[2]s_with_lts_params"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  agency                = "%[3]s"

  log_group_id    = huaweicloud_lts_group.test[0].id
  log_stream_id   = huaweicloud_lts_stream.test[0].id
  log_group_name  = huaweicloud_lts_group.test[0].group_name
  log_stream_name = huaweicloud_lts_stream.test[0].stream_name
  lts_custom_tag = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_fgs_function" "create_without_lts_params" {
  name                  = "%[2]s_without_lts_params"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  agency                = "%[3]s"
}
`, testAccFunction_logConfig_base(name), name, acceptance.HW_FGS_AGENCY_NAME)
}

func testAccFunction_logConfig_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "create_with_lts_params" {
  name                  = "%[2]s_with_lts_params"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  agency                = "%[3]s"

  log_group_id    = huaweicloud_lts_group.test[1].id
  log_stream_id   = huaweicloud_lts_stream.test[1].id
  log_group_name  = huaweicloud_lts_group.test[1].group_name
  log_stream_name = huaweicloud_lts_stream.test[1].stream_name
  lts_custom_tag  = {
    foo     = "baar"
    new_key = "value"
  }
}

resource "huaweicloud_fgs_function" "create_without_lts_params" {
  name                  = "%[2]s_without_lts_params"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  agency                = "%[3]s"

  log_group_id    = huaweicloud_lts_group.test[0].id
  log_stream_id   = huaweicloud_lts_stream.test[0].id
  log_group_name  = huaweicloud_lts_group.test[0].group_name
  log_stream_name = huaweicloud_lts_stream.test[0].stream_name
  lts_custom_tag  = {
    foo = "bar"
    key = "value"
  }
}
`, testAccFunction_logConfig_base(name), name, acceptance.HW_FGS_AGENCY_NAME)
}

func testAccFunction_logConfig_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "create_with_lts_params" {
  name                  = "%[2]s_with_lts_params"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
}
`, testAccFunction_logConfig_base(name), name, acceptance.HW_FGS_AGENCY_NAME)
}

func TestAccFunction_strategy(t *testing.T) {
	var (
		obj interface{}

		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_function.test"

		rc = acceptance.InitResourceCheck(resourceName, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_strategy_undefined(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "400"),
				),
			},
			{
				Config: testAccFunction_strategy_defined(name, 1000),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "1000"),
				),
			},
			{
				Config: testAccFunction_strategy_defined(name, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "0"),
				),
			},
			{
				Config: testAccFunction_strategy_defined(name, -1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "-1"),
				),
			},
			{
				Config: testAccFunction_strategy_undefined(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "max_instance_num", "-1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_strategy_undefined(name string) string {
	return fmt.Sprintf(`
%[1]s

# max_instance_num can only be configured in the ver.2 function.
resource "huaweicloud_fgs_function" "test" {
  name                  = "%[2]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
}
`, functionScriptVariableDefinition, name)
}

func testAccFunction_strategy_defined(name string, maxInstanceNum int) string {
	return fmt.Sprintf(`
%[1]s

# max_instance_num can only be configured in the ver.2 function.
resource "huaweicloud_fgs_function" "test" {
  name                  = "%[2]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
  max_instance_num      = %[3]d
}
`, functionScriptVariableDefinition, name, maxInstanceNum)
}

func TestAccFunction_versions(t *testing.T) {
	var (
		obj interface{}

		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_function.test"

		rc = acceptance.InitResourceCheck(resourceName, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_versions_step1(functionScriptContentDefinition("Hello, world!"), name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.#", "0"),
				),
			},
			{
				Config: testAccFunction_versions_step2(functionScriptContentDefinition("Hi, world!"), name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.#", "2"),
				),
			},
			{
				Config: testAccFunction_versions_step3(functionScriptContentDefinition("Hi, world!"), name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.#", "1"),
				),
			},
			{
				Config: testAccFunction_versions_step4(functionScriptContentDefinition("Yo, world!"), name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.#", "2"),
				),
			},
			{
				Config: testAccFunction_versions_step5(functionScriptContentDefinition("Goodbye, world!"), name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
		},
	})
}

func functionScriptContentDefinition(msg string) string {
	return fmt.Sprintf(`
variable "script_content" {
  type    = string
  default = <<EOT
def main():
    print("%[1]s")

if __name__ == "__main__":
    main()
EOT
}
`, msg)
}

func testAccFunction_versions_step1(funcScript, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name                  = "%[2]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"
}
`, funcScript, name)
}

// Before this configuration supplement, the func_code must be updated.
func testAccFunction_versions_step2(funcScript, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name                  = "%[2]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"

  versions {
    name = "latest"

    aliases {
      name        = "demo"
      description = "This is a description of the alias demo under the version latest."
    }
  }
  # The value of the parameter func_code must be modified before each custom version add.
  versions {
    name        = "v1.0"
    description = "This is a description of the version v1.0. (Prepare to update)"

    aliases {
      name        = "v1_0-alias"
      description = "This is a description of the alias v1_0-alias under the version v1.0."
    }
  }
}
`, funcScript, name)
}

// Delete the alias configuration and recreate the version v1.0.
func testAccFunction_versions_step3(funcScript, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name                  = "%[2]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"

  # The value of the parameter func_code must be modified before each custom version add.
  versions {
    name        = "v1.0"
    description = "This is a description of the version v1.0."

    aliases {
      name        = "v1_0-alias"
      description = "This is a description of the alias v1_0-alias under the version v1.0."
    }
  }
}
`, funcScript, name)
}

// Delete the alias configuration and recreate the version v1.0.
func testAccFunction_versions_step4(funcScript, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name                  = "%[2]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"

  # The value of the parameter func_code must be modified before each custom version add.
  versions {
    name        = "v1.0"
    description = "This is a description of the version v1.0."

    aliases {
      name        = "v1_0-alias"
      description = "This is a description of the alias v1_0-alias under the version v1.0."
    }
  }
  versions {
    name        = "v2.0"
    description = "This is a description of the version v2.0."

    aliases {
      name        = "v2_0-alias"
      description = "This is a description of the alias v2_0-alias under the version v2.0."

      additional_version_weights = jsonencode({
        "v1.0": 15
      })
    }
  }
}
`, funcScript, name)
}

// Before this configuration supplement, the func_code must be updated.
func testAccFunction_versions_step5(funcScript, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name                  = "%[2]s"
  memory_size           = 128
  runtime               = "Python2.7"
  timeout               = 3
  app                   = "default"
  handler               = "index.handler"
  code_type             = "inline"
  func_code             = base64encode(var.script_content)
  description           = "Created by terraform script"
  functiongraph_version = "v2"

  # The value of the parameter func_code must be modified before each custom version add.
  versions {
    name        = "v1.0"
    description = "This is a description of the version v1.0."

    aliases {
      name        = "v1_0-alias"
      description = "This is a description of the alias v1_0-alias under the version v1.0."
    }
  }
  versions {
    name        = "v2.0"
    description = "This is a description of the version v2.0."

    aliases {
      name        = "v2_0-alias"
      description = "This is a description of the alias v2_0-alias under the version v2.0."

      additional_version_weights = jsonencode({
        "v1.0": 15
      })
    }
  }
  versions {
    name        = "v3.0"
    description = "This is a description of the version v2.0."

    aliases {
      name        = "v3_0-alias"
      description = "This is a description of the alias v2_0-alias under the version v3.0."
      additional_version_strategy = jsonencode({
        "v2.0": {
          "combine_type": "or",
          "rules": [
            {
              "rule_type": "Header",
              "param": "version",
              "op": "=",
              "value": "v2_value"
            },
            {
              "rule_type": "Header",
              "param": "Owner",
              "op": "in",
              "value": "terraform,administrator"
            }
          ]
        }
      })
    }
  }
}
`, funcScript, name)
}

func TestAccFunction_network(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceName()

		createWithNetwork   = "huaweicloud_fgs_function.create_with_network"
		rcCreateWithNetwork = acceptance.InitResourceCheck(createWithNetwork, &obj, getFunction)

		createWithoutNetwork   = "huaweicloud_fgs_function.create_without_network"
		rcCreateWithoutNetwork = acceptance.InitResourceCheck(createWithoutNetwork, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcCreateWithNetwork.CheckResourceDestroy(),
			rcCreateWithoutNetwork.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_network_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcCreateWithNetwork.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithNetwork, "name", name+"_with_network"),
					resource.TestCheckResourceAttr(createWithNetwork, "agency", acceptance.HW_FGS_AGENCY_NAME),
					resource.TestCheckResourceAttrPair(createWithNetwork, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(createWithNetwork, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(createWithNetwork, "dns_list"),
					resource.TestCheckResourceAttrSet(createWithNetwork, "peering_cidr"),
					rcCreateWithoutNetwork.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithoutNetwork, "name", name+"_without_network"),
					resource.TestCheckResourceAttr(createWithoutNetwork, "agency", ""),
					resource.TestCheckResourceAttr(createWithoutNetwork, "vpc_id", ""),
					resource.TestCheckResourceAttr(createWithoutNetwork, "network_id", ""),
					resource.TestCheckResourceAttr(createWithoutNetwork, "dns_list", ""),
					resource.TestCheckResourceAttr(createWithoutNetwork, "peering_cidr", ""),
				),
			},
			{
				Config: testAccFunction_network_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcCreateWithNetwork.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(createWithNetwork, "dns_list"),
					resource.TestCheckResourceAttrSet(createWithNetwork, "peering_cidr"),
					rcCreateWithoutNetwork.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithoutNetwork, "agency", acceptance.HW_FGS_AGENCY_NAME),
					resource.TestCheckResourceAttrPair(createWithoutNetwork, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(createWithoutNetwork, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(createWithoutNetwork, "dns_list"),
					resource.TestCheckResourceAttrSet(createWithoutNetwork, "peering_cidr"),
				),
			},
			{
				ResourceName:      createWithNetwork,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
			{
				ResourceName:      createWithoutNetwork,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_network_base(name string) string {
	return fmt.Sprintf(`
%[1]s

variable "base_cidr" {
  default = "192.168.0.0/16"
}

resource "huaweicloud_vpc" "test" {
  name = "%[2]s"
  cidr = cidrsubnet(var.base_cidr, 4, 0)
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[2]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 2, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 2, 0), 1)
}

resource "huaweicloud_vpc" "source" {
  count = 3

  name = format("%[2]s_peering_source_%%d", count.index)
  cidr = cidrsubnet(cidrsubnet(var.base_cidr, 4, 1), 2, count.index)
}

resource "huaweicloud_vpc" "target" {
  count = 3

  name = format("%[2]s_peering_target_%%d", count.index)
  cidr = cidrsubnet(cidrsubnet(var.base_cidr, 4, 2), 2, count.index)
}

resource "huaweicloud_vpc_peering_connection" "test" {
  count = 3

  name        = format("%[2]s_peering_connection_%%d", count.index)
  vpc_id      = huaweicloud_vpc.source[count.index].id
  peer_vpc_id = huaweicloud_vpc.target[count.index].id
}

resource "huaweicloud_dns_zone" "test" {
  count = 3

  zone_type = "private"
  name      = format("functiondebug.example%%d.com.", count.index)

  router {
    router_id = huaweicloud_vpc.test.id
  }
}
`, functionScriptVariableDefinition, name)
}

func testAccFunction_network_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "create_with_network" {
  name        = "%[2]s_with_network"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"

  # VPC access and DNS configuration
  agency     = "%[3]s" # Allow VPC and DNS permissions for FunctionGraph service
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
  dns_list   = jsonencode(
    [for v in slice(huaweicloud_dns_zone.test[*], 0, 2) : tomap({id=v.id, domain_name=v.name})]
  )
  peering_cidr = join(";", slice(huaweicloud_vpc.target[*].cidr, 0, 2))
}

resource "huaweicloud_fgs_function" "create_without_network" {
  name        = "%[2]s_without_network"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"
}
`, testAccFunction_network_base(name), name, acceptance.HW_FGS_AGENCY_NAME)
}

func testAccFunction_network_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "create_with_network" {
  name        = "%[2]s_with_network"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"

  # VPC access and DNS configuration
  agency     = "%[3]s" # Allow VPC and DNS permissions for FunctionGraph service
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
  dns_list   = jsonencode(
    [for v in slice(huaweicloud_dns_zone.test[*], 1, 3) : tomap({id=v.id, domain_name=v.name})]
  )
  peering_cidr = join(";", slice(huaweicloud_vpc.target[*].cidr, 1, 3))
}

resource "huaweicloud_fgs_function" "create_without_network" {
  name        = "%[2]s_without_network"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"

  # VPC access and DNS configuration
  agency     = "%[3]s" # Allow VPC and DNS permissions for FunctionGraph service
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
  dns_list   = jsonencode(
    [for v in slice(huaweicloud_dns_zone.test[*], 0, 2) : tomap({id=v.id, domain_name=v.name})]
  )
  peering_cidr = join(";", slice(huaweicloud_vpc.target[*].cidr, 0, 2))
}
`, testAccFunction_network_base(name), name, acceptance.HW_FGS_AGENCY_NAME)
}

func TestAccFunction_reservedInstance(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceName()

		withVersion   = "huaweicloud_fgs_function.with_version"
		rcWithVersion = acceptance.InitResourceCheck(withVersion, &obj, getFunction)

		withAlias   = "huaweicloud_fgs_function.with_alias"
		rcWithAlias = acceptance.InitResourceCheck(withAlias, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithVersion.CheckResourceDestroy(),
			rcWithAlias.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_reservedInstance_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithVersion.CheckResourceExists(),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.#", "1"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.qualifier_name", "latest"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.qualifier_type", "version"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.count", "1"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.idle_mode", "true"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.tactics_config.0.cron_configs.#", "1"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.tactics_config.0.cron_configs.0.count", "2"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.tactics_config.0.cron_configs.0.cron", "0 */10 * * * ?"),
					resource.TestCheckResourceAttrSet(withVersion, "reserved_instances.0.tactics_config.0.cron_configs.0.start_time"),
					resource.TestCheckResourceAttrSet(withVersion, "reserved_instances.0.tactics_config.0.cron_configs.0.expired_time"),
					resource.TestCheckResourceAttrSet(withVersion, "reserved_instances.0.tactics_config.0.cron_configs.0.name"),
					rcWithAlias.CheckResourceExists(),
					resource.TestCheckResourceAttr(withAlias, "versions.#", "1"),
					resource.TestCheckResourceAttr(withAlias, "versions.0.name", "latest"),
					resource.TestCheckResourceAttr(withAlias, "versions.0.aliases.#", "1"),
					resource.TestCheckResourceAttr(withAlias, "versions.0.aliases.0.name", name),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.#", "1"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.qualifier_name", name),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.qualifier_type", "alias"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.count", "1"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.idle_mode", "true"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.tactics_config.0.cron_configs.#", "1"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.tactics_config.0.cron_configs.0.count", "2"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.tactics_config.0.cron_configs.0.cron", "0 */10 * * * ?"),
					resource.TestCheckResourceAttrSet(withAlias, "reserved_instances.0.tactics_config.0.cron_configs.0.start_time"),
					resource.TestCheckResourceAttrSet(withAlias, "reserved_instances.0.tactics_config.0.cron_configs.0.expired_time"),
					resource.TestCheckResourceAttrSet(withAlias, "reserved_instances.0.tactics_config.0.cron_configs.0.name"),
				),
			},
			{
				Config: testAccFunction_reservedInstance_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithVersion.CheckResourceExists(),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.count", "2"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.idle_mode", "false"),
					resource.TestCheckResourceAttr(withVersion, "reserved_instances.0.tactics_config.#", "0"),
					rcWithAlias.CheckResourceExists(),
					resource.TestCheckResourceAttr(withAlias, "versions.#", "1"),
					resource.TestCheckResourceAttr(withAlias, "versions.0.name", "latest"),
					resource.TestCheckResourceAttr(withAlias, "versions.0.aliases.#", "1"),
					resource.TestCheckResourceAttr(withAlias, "versions.0.aliases.0.name", name+"_new"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.qualifier_name", name+"_new"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.qualifier_type", "alias"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.count", "2"),
					resource.TestCheckResourceAttr(withAlias, "reserved_instances.0.idle_mode", "false"),
				),
			},
			{
				ResourceName:      withVersion,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
			{
				ResourceName:      withAlias,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_reservedInstance_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# Using the current time as the start time.
resource "time_static" "test" {}

# Using the current time one day later as the expiration time.
resource "time_offset" "test" {
  offset_days = 1
}

resource "huaweicloud_fgs_function" "with_version" {
  name        = "%[2]s_with_version"
  memory_size = 128
  runtime     = "Node.js16.17"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"

  reserved_instances {
    qualifier_type = "version"
    qualifier_name = "latest"
    count          = 1
    idle_mode      = true

    tactics_config {
      cron_configs {
        name         = "scheme-waekcy"
        cron         = "0 */10 * * * ?"
        start_time   = time_static.test.unix
        expired_time = time_offset.test.unix
        count        = 2
      }
    }
  }
}

resource "huaweicloud_fgs_function" "with_alias" {
  name        = "%[2]s_with_alias"
  memory_size = 128
  runtime     = "Node.js16.17"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"

  versions {
    name = "latest"
	
    aliases {
      name = "%[2]s"
    }
  }

  reserved_instances {
    qualifier_type = "alias"
    qualifier_name = "%[2]s"
    count          = 1
    idle_mode      = true

    tactics_config {
      cron_configs {
        name         = "scheme-waekcy"
        cron         = "0 */10 * * * ?"
        start_time   = time_static.test.unix
        expired_time = time_offset.test.unix
        count        = 2
      }
    }
  }
}
`, functionScriptVariableDefinition, name)
}

func testAccFunction_reservedInstance_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "with_version" {
  name        = "%[2]s_with_version"
  memory_size = 128
  runtime     = "Node.js16.17"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"
	  
  reserved_instances {
    qualifier_type = "version"
    qualifier_name = "latest"
    count          = 2
    idle_mode      = false
  }
}

resource "huaweicloud_fgs_function" "with_alias" {
  name        = "%[2]s_with_alias"
  memory_size = 128
  runtime     = "Node.js16.17"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"

  versions {
    name = "latest"
	
    aliases {
      name = "%[2]s_new"
    }
  }

  reserved_instances {
    qualifier_type = "alias"
    qualifier_name = "%[2]s_new"
    count          = 2
    idle_mode      = false
  }
}
`, functionScriptVariableDefinition, name)
}

func TestAccFunction_gpuMemory(t *testing.T) {
	var (
		obj interface{}

		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_function.test"

		rc = acceptance.InitResourceCheck(resourceName, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Only the cn-east-3 region supports this function, and you need to submit a service ticket to enable the function.
			acceptance.TestAccPreCheckFgsGpuType(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_gpuMemory_undefined(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "gpu_memory", "0"),
					resource.TestCheckResourceAttr(resourceName, "gpu_type", ""),
				),
			},
			{
				Config: testAccFunction_gpuMemory_defined(name, 1024),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "gpu_memory", "1024"),
					resource.TestCheckResourceAttr(resourceName, "gpu_type", acceptance.HW_FGS_GPU_TYPE),
				),
			},
			{
				Config: testAccFunction_gpuMemory_defined(name, 2048),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "gpu_memory", "2048"),
					resource.TestCheckResourceAttr(resourceName, "gpu_type", acceptance.HW_FGS_GPU_TYPE),
				),
			},
			{
				Config: testAccFunction_gpuMemory_undefined(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "gpu_memory", "0"),
					resource.TestCheckResourceAttr(resourceName, "gpu_type", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_gpuMemory_undefined(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name        = "%[2]s"
  memory_size = 128
  runtime     = "Custom"
  timeout     = 3
  app         = "default"
  handler     = "bootstrap"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"
}
`, functionScriptVariableDefinition, name)
}

func testAccFunction_gpuMemory_defined(name string, memory int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name        = "%[2]s"
  memory_size = 128
  runtime     = "Custom"
  timeout     = 3
  app         = "default"
  handler     = "bootstrap"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"
  gpu_memory  = %[3]d
  gpu_type    = "%[4]s"
}
`, functionScriptVariableDefinition, name, memory, acceptance.HW_FGS_GPU_TYPE)
}

func TestAccFunction_java(t *testing.T) {
	var (
		obj interface{}

		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_function.test"

		rc = acceptance.InitResourceCheck(resourceName, &obj, getFunction)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_java_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memory_size", "128"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Java11"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "app", "default"),
					resource.TestCheckResourceAttr(resourceName, "handler", "com.huawei.demo.TriggerTests.apigTest"),
					resource.TestCheckResourceAttr(resourceName, "code_type", "zip"),
					resource.TestCheckResourceAttr(resourceName, "code_filename", "java-demo.zip"),
					resource.TestCheckResourceAttr(resourceName, "enable_class_isolation", "true"),
					resource.TestCheckResourceAttr(resourceName, "ephemeral_storage", "512"),
					resource.TestCheckResourceAttr(resourceName, "heartbeat_handler", "com.huawei.demo.TriggerTests.heartBeat"),
					resource.TestCheckResourceAttr(resourceName, "restore_hook_handler", "com.huawei.demo.TriggerTests.restoreHook"),
					resource.TestCheckResourceAttr(resourceName, "restore_hook_timeout", "10"),
				),
			},
			{
				Config: testAccFunction_java_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enable_auth_in_header", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_class_isolation", "false"),
					resource.TestCheckResourceAttr(resourceName, "ephemeral_storage", "10240"),
					resource.TestCheckResourceAttr(resourceName, "heartbeat_handler", "com.huawei.demo.TriggerTests.heartBeat"),
					resource.TestCheckResourceAttr(resourceName, "restore_hook_handler", "com.huawei.demo.TriggerTests.restoreHook"),
					resource.TestCheckResourceAttr(resourceName, "restore_hook_timeout", "20"),
				),
			},
			{
				Config: testAccFunction_java_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enable_auth_in_header", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_class_isolation", "false"),
					resource.TestCheckResourceAttr(resourceName, "ephemeral_storage", "10240"),
					resource.TestCheckResourceAttr(resourceName, "heartbeat_handler", ""),
					resource.TestCheckResourceAttr(resourceName, "restore_hook_handler", ""),
					resource.TestCheckResourceAttr(resourceName, "restore_hook_timeout", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"func_code",
				},
			},
		},
	})
}

func testAccFunction_java_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name          = "%[1]s"
  memory_size   = 128
  runtime       = "Java11"
  timeout       = 15
  app           = "default"
  handler       = "com.huawei.demo.TriggerTests.apigTest"
  code_type     = "zip"
  code_filename = "java-demo.zip"
  agency        = "%[2]s"

  enable_class_isolation = true
  ephemeral_storage      = 512
  heartbeat_handler      = "com.huawei.demo.TriggerTests.heartBeat"
  restore_hook_handler   = "com.huawei.demo.TriggerTests.restoreHook"
  restore_hook_timeout   = 10
}
`, name, acceptance.HW_FGS_AGENCY_NAME)
}

func testAccFunction_java_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name          = "%[1]s"
  memory_size   = 128
  runtime       = "Java11"
  timeout       = 15
  app           = "default"
  handler       = "com.huawei.demo.TriggerTests.apigTest"
  code_type     = "zip"
  code_filename = "java-demo.zip"
  agency        = "%[2]s"

  enable_class_isolation = false
  ephemeral_storage      = 10240
  heartbeat_handler      = "com.huawei.demo.TriggerTests.heartBeat"
  restore_hook_handler   = "com.huawei.demo.TriggerTests.restoreHook"
  restore_hook_timeout   = 20
}
`, name, acceptance.HW_FGS_AGENCY_NAME)
}

func testAccFunction_java_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name          = "%[1]s"
  memory_size   = 128
  runtime       = "Java11"
  timeout       = 15
  app           = "default"
  handler       = "com.huawei.demo.TriggerTests.apigTest"
  code_type     = "zip"
  code_filename = "java-demo.zip"
  agency        = "%[2]s"
}
`, name, acceptance.HW_FGS_AGENCY_NAME)
}
