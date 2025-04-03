package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceMonitors_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_monitors.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMonitors_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "monitors.#"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.domain_name"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.id"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.interval"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.status_code"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.http_method"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.max_retries"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.port"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.pool_id"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.timeout"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.protocol"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.url_path"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "monitors.0.updated_at"),
					resource.TestCheckOutput("domain_name_filter_is_useful", "true"),
					resource.TestCheckOutput("monitor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("pool_id_filter_is_useful", "true"),
					resource.TestCheckOutput("interval_filter_is_useful", "true"),
					resource.TestCheckOutput("max_retries_filter_is_useful", "true"),
					resource.TestCheckOutput("timeout_filter_is_useful", "true"),
					resource.TestCheckOutput("status_code_filter_is_useful", "true"),
					resource.TestCheckOutput("http_method_filter_is_useful", "true"),
					resource.TestCheckOutput("url_path_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceMonitors_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_monitors" "test" {
  depends_on = [huaweicloud_elb_monitor.monitor_1]

  monitor_id = huaweicloud_elb_monitor.monitor_1.id
}

locals {
  domain_name = data.huaweicloud_elb_monitors.test.monitors[0].domain_name
}
data "huaweicloud_elb_monitors" "domain_name_filter" {
  domain_name = local.domain_name
}
output "domain_name_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.domain_name_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.domain_name_filter.monitors[*].domain_name :v == local.domain_name]
  )
}

locals {
  monitor_id = data.huaweicloud_elb_monitors.test.monitors[0].id
}
data "huaweicloud_elb_monitors" "monitor_id_filter" {
  monitor_id = local.monitor_id
}
output "monitor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.monitor_id_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.monitor_id_filter.monitors[*].id :v == local.monitor_id]
  )
}

locals {
  pool_id = data.huaweicloud_elb_monitors.test.monitors[0].pool_id
}
data "huaweicloud_elb_monitors" "pool_id_filter" {
  pool_id = local.pool_id
}
output "pool_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.pool_id_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.pool_id_filter.monitors[*].pool_id :v == local.pool_id]
  )
}

locals {
  interval = data.huaweicloud_elb_monitors.test.monitors[0].interval
}
data "huaweicloud_elb_monitors" "interval_filter" {
  interval = local.interval
}
output "interval_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.interval_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.interval_filter.monitors[*].interval :v == local.interval]
  )
}

locals {
  max_retries = data.huaweicloud_elb_monitors.test.monitors[0].max_retries
}
data "huaweicloud_elb_monitors" "max_retries_filter" {
  max_retries = local.max_retries
}
output "max_retries_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.max_retries_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.max_retries_filter.monitors[*].max_retries :v == local.max_retries]
  )
}

locals {
  timeout = data.huaweicloud_elb_monitors.test.monitors[0].timeout
}
data "huaweicloud_elb_monitors" "timeout_filter" {
  timeout = local.timeout
}
output "timeout_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.timeout_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.timeout_filter.monitors[*].timeout :v == local.timeout]
  )
}

locals {
  status_code = data.huaweicloud_elb_monitors.test.monitors[0].status_code
}
data "huaweicloud_elb_monitors" "status_code_filter" {
  status_code = local.status_code
}
output "status_code_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.status_code_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.status_code_filter.monitors[*].status_code :v == local.status_code]
  )
}

locals {
  http_method = data.huaweicloud_elb_monitors.test.monitors[0].http_method
}
data "huaweicloud_elb_monitors" "http_method_filter" {
  http_method =  data.huaweicloud_elb_monitors.test.monitors[0].http_method
}
output "http_method_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.http_method_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.http_method_filter.monitors[*].http_method :v == local.http_method]
  )
}

locals {
  url_path = data.huaweicloud_elb_monitors.test.monitors[0].url_path
}
data "huaweicloud_elb_monitors" "url_path_filter" {
  url_path = local.url_path
}
output "url_path_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.url_path_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.url_path_filter.monitors[*].url_path :v == local.url_path]
  )
}

locals {
  protocol = data.huaweicloud_elb_monitors.test.monitors[0].protocol
}
data "huaweicloud_elb_monitors" "protocol_filter" {
  protocol = local.protocol
}
output "protocol_filter_is_useful" {
  value = length(data.huaweicloud_elb_monitors.protocol_filter.monitors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_monitors.protocol_filter.monitors[*].protocol :v == local.protocol]
  )
}

`, testAccElbV3MonitorConfig_basic(name))
}
