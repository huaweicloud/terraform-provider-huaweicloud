package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHostsRisk_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_hosts_risk.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHostsRisk_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.agent_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.detect_result"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.asset"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.vulnerability"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.baseline"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.intrusion"),
				),
			},
		},
	})
}

func testDataSourceHostsRisk_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_hosts_risk" "test" {
  host_id_list          = ["%s"]
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
