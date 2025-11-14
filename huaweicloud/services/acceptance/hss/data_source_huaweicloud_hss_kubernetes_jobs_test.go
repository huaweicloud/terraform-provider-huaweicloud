package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKubernetesJobs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_kubernetes_jobs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubernetesJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_update_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_info_list.#"),
				),
			},
		},
	})
}

func testAccDataSourceKubernetesJobs_basic() string {
	return `
data "huaweicloud_hss_kubernetes_jobs" "test" {
  enterprise_project_id = "0"
}
`
}
