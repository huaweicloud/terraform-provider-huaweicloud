package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAsyncInvokeConfigurations_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_fgs_async_invoke_configurations.all"
		dc  = acceptance.InitDataSourceCheck(all)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAsyncInvokeConfigurations_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "configurations.#", "1"),
					// Check the attributes of the first configuration
					resource.TestCheckResourceAttrSet(all, "configurations.0.func_urn"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.max_async_event_age_in_seconds"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.max_async_retry_attempts"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.created_time"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.last_modified"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.enable_async_status_log"),
					// Check destination config structure
					resource.TestCheckResourceAttrSet(all, "configurations.0.destination_config.0.on_success.0.destination"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.destination_config.0.on_success.0.param"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.destination_config.0.on_failure.0.destination"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.destination_config.0.on_failure.0.param"),
					// Verify specific values
					resource.TestCheckResourceAttr(all, "configurations.0.max_async_event_age_in_seconds", "3500"),
					resource.TestCheckResourceAttr(all, "configurations.0.max_async_retry_attempts", "2"),
					resource.TestCheckResourceAttr(all, "configurations.0.enable_async_status_log", "true"),
					resource.TestCheckResourceAttr(all, "configurations.0.destination_config.0.on_success.0.destination", "OBS"),
					resource.TestCheckResourceAttr(all, "configurations.0.destination_config.0.on_failure.0.destination", "SMN"),
				),
			},
		},
	})
}

func testAccDataAsyncInvokeConfigurations_base(name string) string {
	return fmt.Sprintf(`
variable "script_content" {
  type    = string
  default = <<EOT
def main():
    print("Hello, World, It is good!")

if __name__ == "__main__":
    main()
EOT
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  memory_size = 128
  description = "Created by terraform script"
  handler     = "index.handler"
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  agency      = "%[2]s"
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_fgs_async_invoke_configuration" "test" {
  function_urn                   = huaweicloud_fgs_function.test.urn
  max_async_event_age_in_seconds = 3500
  max_async_retry_attempts       = 2
  enable_async_status_log        = true

  on_success {
    destination = "OBS"
    param       = jsonencode({
      bucket  = huaweicloud_obs_bucket.test.bucket
      prefix  = "/success"
      expires = 5
    })
  }

  on_failure {
    destination = "SMN"
    param       = jsonencode({
      topic_urn = huaweicloud_smn_topic.test.topic_urn
    })
  }
}
`, name, acceptance.HW_FGS_AGENCY_NAME)
}

func testAccDataAsyncInvokeConfigurations_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_async_invoke_configurations" "all" {
  depends_on = [huaweicloud_fgs_async_invoke_configuration.test]
  
  function_urn = huaweicloud_fgs_function.test.urn
}
`, testAccDataAsyncInvokeConfigurations_base(name))
}
