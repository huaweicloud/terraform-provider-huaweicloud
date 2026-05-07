package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_cluster_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceClusterLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "type"),
					resource.TestCheckResourceAttrSet(dataSource, "completed"),
					resource.TestCheckResourceAttrSet(dataSource, "log_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "log_list.0.date"),
					resource.TestCheckResourceAttrSet(dataSource, "log_list.0.content"),
					resource.TestCheckResourceAttr(dataSource, "instance_log", ""),
					resource.TestCheckOutput("level_filter_is_useful", "true"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
					resource.TestCheckOutput("limit_filter_is_useful", "true"),
					resource.TestCheckOutput("keyword_filter_is_useful", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceClusterLogs_logstash_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_cluster_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSLogStashClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceClusterLogs_logstash_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_log"),
					resource.TestCheckResourceAttr(dataSource, "log_list.#", "0"),
					resource.TestCheckResourceAttr(dataSource, "type", ""),
					resource.TestCheckOutput("keyword_filter_is_useful", "true"),
					resource.TestCheckOutput("limit_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceClusterLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_clusters" "test" {
  cluster_id = "%s"
}

locals {
  css_cluster_id    = data.huaweicloud_css_clusters.test.clusters.0.id
  css_instance_name = data.huaweicloud_css_clusters.test.clusters.0.instances[0].name
}

data "huaweicloud_css_cluster_logs" "test" {
  cluster_id    = local.css_cluster_id
  instance_name = local.css_instance_name
  log_type      = "instance"
  limit          = 10000
}

data "huaweicloud_css_cluster_logs" "level_filter" {
  cluster_id    = local.css_cluster_id
  instance_name = local.css_instance_name
  log_type      = "instance"
  limit          = 10000
  level         = "WARN"

  depends_on = [data.huaweicloud_css_cluster_logs.test]
}

output "level_filter_is_useful" {
  value = alltrue([for log in data.huaweicloud_css_cluster_logs.level_filter.log_list : log.level == "WARN"])
}

data "huaweicloud_css_cluster_logs" "limit_filter" {
  cluster_id    = local.css_cluster_id
  instance_name = local.css_instance_name
  log_type      = "instance"
  limit         = 10
  level         = "WARN"

  depends_on = [data.huaweicloud_css_cluster_logs.level_filter]
}

output "limit_filter_is_useful" {
  value = length(data.huaweicloud_css_cluster_logs.limit_filter.log_list) == 10
}

data "huaweicloud_css_cluster_logs" "keyword_filter" {
  cluster_id    = local.css_cluster_id
  instance_name = local.css_instance_name
  log_type      = "instance"
  limit         = 10
  level         = "WARN"
  keyword       = "Obs"

  depends_on = [data.huaweicloud_css_cluster_logs.limit_filter]
}

output "keyword_filter_is_useful" {
  value = alltrue([for log in data.huaweicloud_css_cluster_logs.keyword_filter.log_list : strcontains(log.content, "Obs")])
}

data "huaweicloud_css_cluster_logs" "time_filter" {
  cluster_id    = local.css_cluster_id
  instance_name = local.css_instance_name
  log_type      = "instance"
  limit         = 10
  level         = "WARN"
  keyword       = "Obs"
  time_index    = "2000-01-01T00:00:00,000"

  depends_on = [data.huaweicloud_css_cluster_logs.keyword_filter]
}

output "time_filter_is_useful" {
  value = length(data.huaweicloud_css_cluster_logs.time_filter.log_list) == 0
}
`, acceptance.HW_CSS_CLUSTER_ID)
}

func testDataSourceClusterLogs_logstash_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_clusters" "test" {
  cluster_id = "%s"
}

locals {
  logstash_cluster_id    = data.huaweicloud_css_clusters.test.clusters.0.id
  logstash_instance_name = data.huaweicloud_css_clusters.test.clusters.0.instances[0].name
}

data "huaweicloud_css_cluster_logs" "test" {
  cluster_id    = local.logstash_cluster_id
  instance_name = local.logstash_instance_name
  log_type      = "instance"
}

data "huaweicloud_css_cluster_logs" "keyword_filter" {
  cluster_id    = local.logstash_cluster_id
  instance_name = local.logstash_instance_name
  log_type      = "instance"
  keyword       = "logstash"

  depends_on = [data.huaweicloud_css_cluster_logs.test]
}

locals {
  keyword_filter_log_lines = [for line in split("\n", data.huaweicloud_css_cluster_logs.keyword_filter.instance_log) : line if length(line) > 0]
}

output "keyword_filter_is_useful" {
  value = alltrue([for line in local.keyword_filter_log_lines : strcontains(line, "logstash")])
}

data "huaweicloud_css_cluster_logs" "limit_filter" {
  cluster_id    = local.logstash_cluster_id
  instance_name = local.logstash_instance_name
  log_type      = "instance"
  keyword       = "logstash"
  limit         = 1

  depends_on = [data.huaweicloud_css_cluster_logs.keyword_filter]
}

locals {
  limit_filter_log_lines = [for line in split("\n", data.huaweicloud_css_cluster_logs.limit_filter.instance_log) : line if length(line) > 0]
}

output "limit_filter_is_useful" {
  value = length(local.limit_filter_log_lines) == 1
}
`, acceptance.HW_CSS_LOGSTASH_CLUSTER_ID)
}
