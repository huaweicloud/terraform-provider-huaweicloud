package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
)

func getAsyncLogConfigurationResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	fgsClient, err := cfg.NewServiceClient("fgs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}
	ltsClient, err := cfg.NewServiceClient("lts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}
	return fgs.GetAsyncLogConfigurationAndCheck(fgsClient, ltsClient)
}

func TestAccAsyncLogConfiguration_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_fgs_async_log_configuration.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getAsyncLogConfigurationResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAsyncLogConfiguration_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "force_delete", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"max_retries",
				},
			},
			{
				Config: testAccAsyncLogConfiguration_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "force_delete", "true"),
				),
			},
		},
	})
}

const testAccAsyncLogConfiguration_basic_step1 string = `
resource "huaweicloud_fgs_async_log_configuration" "test" {
  force_delete = false

  lifecycle {
    ignore_changes = [
      force_delete,
    ]
  }
}
`

const testAccAsyncLogConfiguration_basic_step2 string = `
resource "huaweicloud_fgs_async_log_configuration" "test" {
  force_delete = true
}
`
