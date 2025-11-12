package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePublications_basic(t *testing.T) {
	rName := "data.huaweicloud_rds_publications.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsSubscriberInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePublications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "publications.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.id"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.status"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.publication_name"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.publication_database"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.subscription_count"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.subscription_options.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.subscription_options.0.independent_agent"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.subscription_options.0.snapshot_always_available"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.subscription_options.0.replicate_ddl"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.subscription_options.0.allow_initialize_from_backup"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.id"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.job_schedule_type"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.one_time_occurrence.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.frequency.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.frequency.0.freq_type"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.frequency.0.freq_interval"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.job_schedule.0.frequency.0.freq_interval_monthly"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.job_schedule.0.frequency.0.freq_relative_interval_monthly"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.daily_frequency.#"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.job_schedule.0.daily_frequency.0.freq_subday_type"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.job_schedule.0.daily_frequency.0.active_start_time"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.job_schedule.0.daily_frequency.0.active_end_time"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.job_schedule.0.daily_frequency.0.freq_subday_interval"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.job_schedule.0.daily_frequency.0.freq_interval_unit"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.duration.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.duration.0.active_start_date"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.job_schedule.0.duration.0.active_end_date"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.is_select_all_table"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.extend_tables.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.table_name"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.schema"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.columns.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.primary_key.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.filter_statement"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.filter.#"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.filter.0.relation"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.filter.0.filters"),
					resource.TestCheckResourceAttrSet(rName, "publications.0.tables.0.article_properties.#"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.tables.0.article_properties.0.destination_object_name"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.tables.0.article_properties.0.destination_object_owner"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.tables.0.article_properties.0.insert_delivery_format"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.tables.0.article_properties.0.insert_stored_procedure"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.tables.0.article_properties.0.update_delivery_format"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.tables.0.article_properties.0.update_stored_procedure"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.tables.0.article_properties.0.delete_delivery_format"),
					resource.TestCheckResourceAttrSet(rName,
						"publications.0.tables.0.article_properties.0.delete_stored_procedure"),
					resource.TestCheckOutput("publication_name_filter_is_useful", "true"),
					resource.TestCheckOutput("publication_db_name_filter_is_useful", "true"),
					resource.TestCheckOutput("subscriber_instance_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePublications_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_publications" "test" {
  instance_id = "%[1]s"
}

locals {
  publication_name = data.huaweicloud_rds_publications.test.publications[0].publication_name
}
data "huaweicloud_rds_publications" "publication_name_filter" {
  instance_id      = "%[1]s"
  publication_name = data.huaweicloud_rds_publications.test.publications[0].publication_name
}
output "publication_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_publications.publication_name_filter.publications) > 0 && alltrue(
    [for v in data.huaweicloud_rds_publications.publication_name_filter.publications[*].publication_name :
    v == local.publication_name]
  )
}

locals {
  publication_db_name = data.huaweicloud_rds_publications.test.publications[0].publication_database
}
data "huaweicloud_rds_publications" "publication_db_name_filter" {
  instance_id         = "%[1]s"
  publication_db_name = data.huaweicloud_rds_publications.test.publications[0].publication_database
}
output "publication_db_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_publications.publication_db_name_filter.publications) > 0 && alltrue(
    [for v in data.huaweicloud_rds_publications.publication_db_name_filter.publications[*].publication_database :
    v == local.publication_db_name]
  )
}

data "huaweicloud_rds_publications" "subscriber_instance_id_filter" {
  instance_id            = "%[1]s"
  subscriber_instance_id = "%[2]s"
}
output "subscriber_instance_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_publications.subscriber_instance_id_filter.publications) > 0 && alltrue(
  [for v in data.huaweicloud_rds_publications.subscriber_instance_id_filter.publications[*].subscription_count : v > 0]
  )
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_SUBSCRIBER_INSTANCE_ID)
}
