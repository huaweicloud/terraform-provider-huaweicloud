package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAvailableIpNumberDataSource_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_available_ip_number.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		notExistSubnetName = "data.huaweicloud_workspace_available_ip_number.test_non_exist_subnet"
		dcNonExistSubnet   = acceptance.InitDataSourceCheck(notExistSubnetName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAvailableIpNumberDataSource_step1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "available_ip", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				Config: testAccAvailableIpNumberDataSource_step2,
				Check: resource.ComposeTestCheckFunc(
					dcNonExistSubnet.CheckResourceExists(),
					resource.TestCheckResourceAttr(notExistSubnetName, "available_ip", "0"),
				),
			},
		},
	})
}

const testAccAvailableIpNumberDataSource_step1 = `
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_available_ip_number" "test" {
  subnet_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
}
`

const testAccAvailableIpNumberDataSource_step2 = `
data "huaweicloud_workspace_available_ip_number" "test_non_exist_subnet" {
  subnet_id = "NOT_FOUND"
}
`
