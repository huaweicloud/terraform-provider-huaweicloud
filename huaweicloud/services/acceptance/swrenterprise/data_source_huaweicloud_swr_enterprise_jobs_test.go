package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseJobs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_jobs.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSwrEnterpriseJobs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.updated_at"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSwrEnterpriseJobs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_jobs" "test" {
  depends_on = [huaweicloud_swr_enterprise_instance.test]
}

data "huaweicloud_swr_enterprise_jobs" "filter_by_status" {
  status = try(data.huaweicloud_swr_enterprise_jobs.test.jobs[0].status, "")
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_jobs.filter_by_status.jobs) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_jobs.filter_by_status.jobs[*].status : 
	  v == data.huaweicloud_swr_enterprise_jobs.test.jobs[0].status]
  )
}
`, testAccSwrEnterpriseInstance_update(name))
}
