package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerNetworkInfo_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_network_info.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires the preparation of a CCE cluster under the default enterprise project.
			acceptance.TestAccPreCheckHSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContainerNetworkInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnet"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_group"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cidrs"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_proxy_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_support_egress"),
				),
			},
		},
	})
}

func testAccDataSourceContainerNetworkInfo_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_container_network_info" "test" {
  cluster_id            = "%s"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_CLUSTER_ID)
}
