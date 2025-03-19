package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceElbAsynchronousTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_asynchronous_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbAsynchronousTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.job_type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.job_type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.entities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.entities.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.entities.0.resource_type"),
					resource.TestCheckOutput("job_id_filter_is_useful", "true"),
					resource.TestCheckOutput("job_type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("error_code_filter_is_useful", "true"),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),
					resource.TestCheckOutput("begin_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceElbAsynchronousTasks_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name               = "%[2]s"
  vpc_id             = huaweicloud_vpc.test.id
  ipv4_subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  waf_failure_action = "discard"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]
}

resource "huaweicloud_elb_loadbalancer_copy" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  name            = "%[2]s"
}
`, common.TestBaseNetwork(rName), rName)
}

func testDataSourceElbAsynchronousTasks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_elb_asynchronous_tasks" "test" {}

locals {
  job_id = data.huaweicloud_elb_asynchronous_tasks.test.jobs[0].job_id
}

data "huaweicloud_elb_asynchronous_tasks" "job_id_filter" {
  depends_on = [
    huaweicloud_elb_loadbalancer_copy.test,
    data.huaweicloud_elb_asynchronous_tasks.test
  ]

  job_id = data.huaweicloud_elb_asynchronous_tasks.test.jobs[0].job_id
}

output "job_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_asynchronous_tasks.job_id_filter.jobs) > 0 && alltrue(
  [for v in data.huaweicloud_elb_asynchronous_tasks.job_id_filter.jobs[*].job_id : v == local.job_id]
  )
}

locals {
  job_type = "CLONE"
}

data "huaweicloud_elb_asynchronous_tasks" "job_type_filter" {
  depends_on = [
    huaweicloud_elb_loadbalancer_copy.test,
  ]

  job_type = "CLONE"
}

output "job_type_filter_is_useful" {
  value = length(data.huaweicloud_elb_asynchronous_tasks.job_type_filter.jobs) > 0 && alltrue(
  [for v in data.huaweicloud_elb_asynchronous_tasks.job_type_filter.jobs[*].job_type : v == local.job_type]
  )
}

locals {
  status = "COMPLETE"
}

data "huaweicloud_elb_asynchronous_tasks" "status_filter" {
  depends_on = [
    huaweicloud_elb_loadbalancer_copy.test,
  ]

  status = "COMPLETE"
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_elb_asynchronous_tasks.status_filter.jobs) > 0 && alltrue(
  [for v in data.huaweicloud_elb_asynchronous_tasks.status_filter.jobs[*].status : v == local.status]
  )
}

locals {
  error_code = ""
}

data "huaweicloud_elb_asynchronous_tasks" "error_code_filter" {
  depends_on = [
    huaweicloud_elb_loadbalancer_copy.test,
  ]

  error_code = ""
}

output "error_code_filter_is_useful" {
  value = length(data.huaweicloud_elb_asynchronous_tasks.error_code_filter.jobs) > 0 && alltrue(
  [for v in data.huaweicloud_elb_asynchronous_tasks.error_code_filter.jobs[*].error_code : v == local.error_code]
  )
}

locals {
  resource_id = huaweicloud_elb_loadbalancer.test.id
}

data "huaweicloud_elb_asynchronous_tasks" "resource_id_filter" {
  depends_on = [
    huaweicloud_elb_loadbalancer_copy.test,
  ]

  resource_id = huaweicloud_elb_loadbalancer.test.id
}

output "resource_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_asynchronous_tasks.resource_id_filter.jobs) > 0 && alltrue(
  [for v in data.huaweicloud_elb_asynchronous_tasks.resource_id_filter.jobs[*].resource_id : v == local.resource_id]
  )
}

locals {
  begin_time = data.huaweicloud_elb_asynchronous_tasks.test.jobs[0].begin_time
}

data "huaweicloud_elb_asynchronous_tasks" "begin_time_filter" {
  depends_on = [
    huaweicloud_elb_loadbalancer_copy.test,
    data.huaweicloud_elb_asynchronous_tasks.test
  ]

  begin_time = data.huaweicloud_elb_asynchronous_tasks.test.jobs[0].begin_time
}

output "begin_time_filter_is_useful" {
  value = length(data.huaweicloud_elb_asynchronous_tasks.begin_time_filter.jobs) > 0
}
`, testDataSourceElbAsynchronousTasks_base(name))
}
