package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSubscriptions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_subscriptions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRdsSubscriptions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.publication_id"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.publication_name"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.is_cloud"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.subscription_database"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.subscription_type"),

					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.local_subscription.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.local_subscription.0.publication_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.local_subscription.0.publication_instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.job_schedule.#"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.job_schedule.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.job_schedule.0.job_schedule_type"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.job_schedule.0.one_time_occurrence.#"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.job_schedule.0.frequency.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.frequency.0.freq_type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.frequency.0.freq_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.job_schedule.0.daily_frequency.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.daily_frequency.0.freq_subday_type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.daily_frequency.0.active_start_time"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.daily_frequency.0.active_end_time"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.daily_frequency.0.freq_subday_interval"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.daily_frequency.0.freq_interval_unit"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.job_schedule.0.duration.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.duration.0.active_start_date"),
					resource.TestCheckResourceAttrSet(dataSource,
						"subscriptions.0.job_schedule.0.duration.0.active_end_date"),

					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_subscriptions.publication_id_filter",
						"subscriptions.0.publication_subscription.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_subscriptions.publication_id_filter",
						"subscriptions.0.publication_subscription.0.subscription_instance_name"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_subscriptions.publication_id_filter",
						"subscriptions.0.publication_subscription.0.subscription_instance_ip"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_subscriptions.publication_id_filter",
						"subscriptions.0.publication_subscription.0.subscription_instance_id"),

					resource.TestCheckOutput("publication_id_filter_is_useful", "true"),
					resource.TestCheckOutput("is_cloud_filter_is_useful", "true"),
					resource.TestCheckOutput("publication_name_filter_is_useful", "true"),
					resource.TestCheckOutput("subscription_db_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRdsSubscriptions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_publications" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_subscriptions" "test" {
  instance_id = "%[2]s"
}

locals {
  publication_id = data.huaweicloud_rds_publications.test.publications[0].id
}
data "huaweicloud_rds_subscriptions" "publication_id_filter" {
  instance_id    = "%[1]s"
  publication_id = data.huaweicloud_rds_publications.test.publications[0].id
}
output "publication_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_subscriptions.publication_id_filter.subscriptions) > 0 && alltrue(
    [for v in data.huaweicloud_rds_subscriptions.publication_id_filter.subscriptions[*].publication_id :
    v == local.publication_id]
  )
}

locals {
  is_cloud = "true"
}
data "huaweicloud_rds_subscriptions" "is_cloud_filter" {
  instance_id = "%[2]s"
  is_cloud    = "true"
}
output "is_cloud_filter_is_useful" {
  value = length(data.huaweicloud_rds_subscriptions.is_cloud_filter.subscriptions) > 0 && alltrue(
    [for v in data.huaweicloud_rds_subscriptions.is_cloud_filter.subscriptions[*].is_cloud : v]
  )
}

locals {
  publication_name = data.huaweicloud_rds_publications.test.publications[0].publication_name
}
data "huaweicloud_rds_subscriptions" "publication_name_filter" {
  instance_id      = "%[2]s"
  publication_name = data.huaweicloud_rds_publications.test.publications[0].publication_name
}
output "publication_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_subscriptions.publication_name_filter.subscriptions) > 0 && alltrue(
    [for v in data.huaweicloud_rds_subscriptions.publication_name_filter.subscriptions[*].publication_name :
    v == local.publication_name]
  )
}

locals {
  subscription_db_name = data.huaweicloud_rds_subscriptions.test.subscriptions[0].subscription_database
}
data "huaweicloud_rds_subscriptions" "subscription_db_name_filter" {
  instance_id          = "%[2]s"
  subscription_db_name = data.huaweicloud_rds_subscriptions.test.subscriptions[0].subscription_database
}
output "subscription_db_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_subscriptions.subscription_db_name_filter.subscriptions) > 0 && alltrue(
    [for v in data.huaweicloud_rds_subscriptions.subscription_db_name_filter.subscriptions[*].subscription_database :
    v == local.subscription_db_name]
  )
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_SUBSCRIBER_INSTANCE_ID)
}
