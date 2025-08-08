package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Each workspace has default endpoint service data, so you can run test cases directly.
func TestAccDataSourceSecmasterVpcEndpointServices_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_vpc_endpoint_services.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterVpcEndpointServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.deprecated"),
				),
			},
		},
	})
}

func testDataSourceSecmasterVpcEndpointServices_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_vpc_endpoint_services" "test" {
  workspace_id = "%[1]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
