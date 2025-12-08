package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAvailableIpNumber_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_available_ip_number.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAvailableIpNumber_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "available_ip", "0"),
				),
			},
			{
				Config: testAccDataAvailableIpNumber_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "available_ip", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

const testAccDataAvailableIpNumber_basic_step1 = `
data "huaweicloud_workspace_available_ip_number" "test" {
  subnet_id = "NOT_FOUND"
}
`

const testAccDataAvailableIpNumber_basic_step2 = `
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_available_ip_number" "test" {
  subnet_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
}
`
