package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cts"
)

func getConfigurationResourceFunc(conf *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CTS client: %s", err)
	}

	return cts.GetConfiguration(client)
}

func TestAccConfiguration_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_cts_configuration.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getConfigurationResourceFunc)
	)
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfiguration_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_sync_global_trace", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_support_read_only", "true"),
					resource.TestCheckResourceAttr(resourceName, "support_read_only_services.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "support_read_only_services.0", "ECS"),
					resource.TestCheckResourceAttr(resourceName, "support_read_only_services.1", "EVS"),
				),
			},
			{
				Config: testAccConfiguration_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_sync_global_trace", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_support_read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "support_read_only_services.#", "0"),
				),
			},
			{
				Config: testAccConfiguration_basic_step3,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_sync_global_trace", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_support_read_only", "true"),
					resource.TestCheckResourceAttr(resourceName, "support_read_only_services.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "support_read_only_services.0", "VPC"),
				),
			},
		},
	})
}

const testAccConfiguration_basic_step1 = `
resource "huaweicloud_cts_configuration" "test" {
  is_sync_global_trace       = true
  is_support_read_only       = true
  support_read_only_services = ["ECS", "EVS"]
}
`

const testAccConfiguration_basic_step2 = `
resource "huaweicloud_cts_configuration" "test" {
  is_sync_global_trace = true
  is_support_read_only = false
}
`

const testAccConfiguration_basic_step3 = `
resource "huaweicloud_cts_configuration" "test" {
  is_sync_global_trace       = false
  is_support_read_only       = true
  support_read_only_services = ["VPC"]
}
`
