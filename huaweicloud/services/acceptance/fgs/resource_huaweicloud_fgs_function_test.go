package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
					resource.TestCheckResourceAttr(withObsStorage, "functiongraph_version", "v2"),
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
					resource.TestCheckResourceAttr(withObsStorage, "functiongraph_version", "v2"),
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
					resource.TestCheckResourceAttr(withCustomImage, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrPair(withCustomImage, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(withCustomImage, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(withCustomImage, "concurrency_num", "1"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.#", "1"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.url", acceptance.HW_BUILD_IMAGE_URL),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.command", ""),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.args", ""),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.working_dir", "/"),
					resource.TestCheckResourceAttr(withObsStorage, "user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withObsStorage, "encrypted_user_data", "{\"owner\":\"terraform\"}"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.%", "2"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.key", "value"),
					resource.TestCheckResourceAttr(withObsStorage, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttrSet(withObsStorage, "urn"),
					resource.TestCheckResourceAttrSet(withObsStorage, "version"),
				),
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
					resource.TestCheckResourceAttr(withObsStorage, "functiongraph_version", "v2"),
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
					resource.TestCheckResourceAttr(withObsStorage, "functiongraph_version", "v2"),
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
					resource.TestCheckResourceAttr(withObsStorage, "concurrency_num", "500"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.#", "1"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.url", acceptance.HW_BUILD_IMAGE_URL_UPDATED),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.command", "/bin/sh"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.args", "-args,value"),
					resource.TestCheckResourceAttr(withCustomImage, "custom_image.0.working_dir", "/"),
					resource.TestCheckResourceAttr(withObsStorage, "user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withObsStorage, "encrypted_user_data", "{\"owner\":\"terraform\",\"usage\":\"acceptance test\"}"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.%", "2"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(withObsStorage, "tags.new_key", "value"),
					resource.TestCheckResourceAttr(withObsStorage, "functiongraph_version", "v2"),
					resource.TestCheckResourceAttrSet(withObsStorage, "urn"),
					resource.TestCheckResourceAttrSet(withObsStorage, "version"),
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

  user_data = jsonencode({
    "owner": "terraform"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform"
  })

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

  user_data = jsonencode({
    "owner": "terraform"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform"
  })

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

  user_data = jsonencode({
    "owner": "terraform"
  })
  encrypted_user_data = jsonencode({
    "owner": "terraform"
  })

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

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, zipFileUploadResourcesConfig(name, runtime), name, runtime,
		acceptance.HW_FGS_AGENCY_NAME,
		acceptance.HW_BUILD_IMAGE_URL)
}

func testAccFunction_basic_step2(name, runtime string) string {
	//nolint:revive
	return fmt.Sprintf(`
%[1]s

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
				Config: testAccFunction_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					// Default value is v2. Some regions support only v1, the default value is v1.
					resource.TestMatchResourceAttr(resourceName, "functiongraph_version", regexp.MustCompile(`v1|v2`)),
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

func testAccFunction_withEpsId(name string) string {
	//nolint:revive
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
`, functionScriptVariableDefinition,
		name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccFunction_logConfig(t *testing.T) {
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
				Config: testAccFunction_logConfig_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "functiongraph_version", "v1"),
					resource.TestCheckResourceAttrSet(resourceName, "log_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "log_stream_id"),
				),
			},
			{
				Config: testAccFunction_logConfig_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "functiongraph_version", "v1"),
					resource.TestCheckResourceAttrSet(resourceName, "log_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "log_stream_id"),
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
  functiongraph_version = "v1"

  log_group_id    = huaweicloud_lts_group.test[0].id
  log_stream_id   = huaweicloud_lts_stream.test[0].id
  log_group_name  = huaweicloud_lts_group.test[0].group_name
  log_stream_name = huaweicloud_lts_stream.test[0].stream_name
}
`, testAccFunction_logConfig_base(name), name)
}

func testAccFunction_logConfig_step2(name string) string {
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
  functiongraph_version = "v1"

  log_group_id    = huaweicloud_lts_group.test[1].id
  log_stream_id   = huaweicloud_lts_stream.test[1].id
  log_group_name  = huaweicloud_lts_group.test[1].group_name
  log_stream_name = huaweicloud_lts_stream.test[1].stream_name
}
`, testAccFunction_logConfig_base(name), name)
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
				Config: testAccFunction_versions_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.#", "0"),
				),
			},
			{
				Config: testAccFunction_versions_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.#", "1"),
				),
			},
			{
				Config: testAccFunction_versions_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "versions.#", "2"),
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

func testAccFunction_versions_step1(name string) string {
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
`, functionScriptVariableDefinition, name)
}

func testAccFunction_versions_step2(name string) string {
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
    name = "%[2]s"

    aliases {
      name        = "custom_alias"
      description = "This is a description of the custom alias"
    }
  }
}
`, functionScriptVariableDefinition, name)
}

func testAccFunction_versions_step3(name string) string {
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
      name = "demo"
    }
  }
  versions {
    name = "%[2]s"

    aliases {
      name = "custom_alias_update"
    }
  }
}
`, functionScriptVariableDefinition, name)
}

func TestAccFunction_domain(t *testing.T) {
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
				Config: testAccFunction_domain_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccFunction_domain_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

func testAccFunction_domain_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_dns_zone" "test" {
  count = 3

  zone_type = "private"
  name      = format("functiondebug.example%%d.com.", count.index)

  router {
    router_id = huaweicloud_vpc.test.id
  }
}
`, functionScriptVariableDefinition, common.TestVpc(name))
}

func testAccFunction_domain_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name        = "%[2]s"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"

  # VPC access and DNS configuration
  agency     = "function_all_trust" # Allow VPC and DNS permissions for FunctionGraph service
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
  dns_list   = jsonencode(
    [for v in slice(huaweicloud_dns_zone.test[*], 0, 2) : tomap({id=v.id, domain_name=v.name})]
  )
}
`, testAccFunction_domain_base(name), name)
}

func testAccFunction_domain_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test" {
  name        = "%[2]s"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  description = "Created by terraform script"

  # VPC access and DNS configuration
  agency     = "function_all_trust" # Allow VPC and DNS permissions for FunctionGraph service
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
  dns_list   = jsonencode(
    [for v in slice(huaweicloud_dns_zone.test[*], 1, 3) : tomap({id=v.id, domain_name=v.name})]
  )
}
`, testAccFunction_domain_base(name), name)
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
