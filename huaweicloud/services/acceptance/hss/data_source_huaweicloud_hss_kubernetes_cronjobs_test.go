package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `data_list`.
func TestAccDataSourceKubernetesCronJobs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_kubernetes_cronjobs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubernetesCronJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_update_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cronjob_info_list.#"),
				),
			},
		},
	})
}

func testAccDataSourceKubernetesCronJobs_basic() string {
	return `
data "huaweicloud_hss_kubernetes_cronjobs" "test" {
  enterprise_project_id = "0"
}
`
}
