package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPublication(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/publications"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}

	publication := utils.PathSearch(fmt.Sprintf("publications[?id=='%s']|[0]", state.Primary.ID), listRespBody, nil)
	if publication == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/instances/{instance_id}/replication/publications",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the RDS publication (%s) does not exist", state.Primary.ID)),
			},
		}
	}

	return publication, nil
}

func TestAccPublication_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_publication.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPublication,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublication_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "publication_name", rName),
					resource.TestCheckResourceAttr(resourceName, "publication_database", "test_db"),
					resource.TestCheckResourceAttr(resourceName,
						"subscription_options.0.independent_agent", "false"),
					resource.TestCheckResourceAttr(resourceName,
						"subscription_options.0.snapshot_always_available", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscription_options.0.replicate_ddl", "false"),
					resource.TestCheckResourceAttr(resourceName,
						"subscription_options.0.allow_initialize_from_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "job_schedule.0.job_schedule_type", "one_time"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.one_time_occurrence.0.active_start_date", "2026-04-06"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.one_time_occurrence.0.active_start_time", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "extend_tables.0", "test_table_2"),
					resource.TestCheckResourceAttr(resourceName, "is_select_all_table", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subscription_count"),
				),
			},
			{
				Config: testAccPublication_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName,
						"subscription_options.0.independent_agent", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"subscription_options.0.snapshot_always_available", "false"),
					resource.TestCheckResourceAttr(resourceName, "subscription_options.0.replicate_ddl", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"subscription_options.0.allow_initialize_from_backup", "false"),
					resource.TestCheckResourceAttr(resourceName, "job_schedule.0.job_schedule_type", "one_time"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.one_time_occurrence.0.active_start_date", "2026-05-06"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.one_time_occurrence.0.active_start_time", "07:20:30"),
				),
			},
			{
				Config: testAccPublication_with_recurring_daily(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.job_schedule_type", "recurring"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_type", "daily"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_interval", "2"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.daily_frequency.0.freq_subday_type", "once"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.daily_frequency.0.active_start_time", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.duration.0.active_start_date", "2020-01-01"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.duration.0.active_end_date", "2099-12-31"),
				),
			},
			{
				Config: testAccPublication_with_recurring_weekly(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.job_schedule_type", "recurring"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_type", "weekly"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_interval_weekly.#", "3"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.daily_frequency.0.freq_subday_type", "multiple"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.daily_frequency.0.active_start_time", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.daily_frequency.0.active_end_time", "20:00:00"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.daily_frequency.0.freq_subday_interval", "10"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.daily_frequency.0.freq_interval_unit", "minute"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.duration.0.active_start_date", "2010-07-15"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.duration.0.active_end_date", "2089-10-20"),
				),
			},
			{
				Config: testAccPublication_with_recurring_monthly_day(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.job_schedule_type", "recurring"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_type", "monthly_day"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_interval_day_monthly", "15"),
				),
			},
			{
				Config: testAccPublication_with_recurring_monthly_week(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.job_schedule_type", "recurring"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_type", "monthly_week"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_interval_monthly", "weekday"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_relative_interval_monthly", "second"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPublicationResourceImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"is_create_snapshot_immediately",
				},
			},
		},
	})
}

func TestAccPublication_with_tables(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_publication.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPublication,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublication_with_tables_filter(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.job_schedule_type", "recurring"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_type", "monthly_week"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_interval_monthly", "weekday"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_relative_interval_monthly", "second"),
					resource.TestCheckResourceAttr(resourceName, "is_select_all_table", "false"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.table_name", "test_table_1"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.schema", "test_schema_1"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.columns.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.primary_key.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.filter.0.column", "id"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.filter.0.condition", "="),
					resource.TestCheckResourceAttr(resourceName, "tables.0.filter.0.value", "123"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.destination_object_name", "test_table_1"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.destination_object_owner", "test_schema_1"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.insert_delivery_format", "insert"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.insert_stored_procedure", "sp_MSins_test_schema_1test_table_1"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.update_delivery_format", "update"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.update_stored_procedure", "sp_MSupd_test_schema_1test_table_1"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.delete_delivery_format", "do_not_delete"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.delete_stored_procedure", "sp_MSdel_test_schema_1test_table_1"),
				),
			},
			{
				Config: testAccPublication_with_tables_filter_filters(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.job_schedule_type", "recurring"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_type", "monthly_week"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_interval_monthly", "weekday"),
					resource.TestCheckResourceAttr(resourceName,
						"job_schedule.0.frequency.0.freq_relative_interval_monthly", "second"),
					resource.TestCheckResourceAttr(resourceName, "is_select_all_table", "false"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.table_name", "test_table_1"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.schema", "test_schema_1"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.columns.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.primary_key.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.filter.0.relation", "AND"),
					resource.TestCheckResourceAttr(resourceName, "tables.0.filter.0.filters.#", "2"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.destination_object_name", "test_table_1"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.destination_object_owner", "test_schema_1"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.insert_delivery_format", "call_procedure"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.insert_stored_procedure", "sp_MSins_test_schema_1test_table_1"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.update_delivery_format", "scall_procedure"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.update_stored_procedure", "sp_MSupd_test_schema_1test_table_1"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.delete_delivery_format", "call_procedure"),
					resource.TestCheckResourceAttr(resourceName,
						"tables.0.article_properties.0.delete_stored_procedure", "sp_MSdel_test_schema_1test_table_1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPublicationResourceImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"is_create_snapshot_immediately",
				},
			},
		},
	})
}

func testAccPublication_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_publication" "test" {
  instance_id                    = "%[1]s"
  publication_name               = "%[2]s"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "false"
    snapshot_always_available    = "true"
    replicate_ddl                = "false"
    allow_initialize_from_backup = "true"
  }

  job_schedule {
    job_schedule_type = "one_time"

    one_time_occurrence {
      active_start_date = "2026-04-06"
      active_start_time = "06:00:00"
    }
  }

  extend_tables = ["test_table_2"]

  is_select_all_table = "true"
}
`, acceptance.HW_RDS_INSTANCE_ID, rName)
}

func testAccPublication_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_publication" "test" {
  instance_id                    = "%[1]s"
  publication_name               = "%[2]s"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "true"
    snapshot_always_available    = "false"
    replicate_ddl                = "true"
    allow_initialize_from_backup = "false"
  }

  job_schedule{
    job_schedule_type = "one_time"

    one_time_occurrence {
      active_start_date = "2026-05-06"
      active_start_time = "07:20:30"
    }
  }

  is_select_all_table = "true"
}
`, acceptance.HW_RDS_INSTANCE_ID, rName)
}

func testAccPublication_with_recurring_daily(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_publication" "test" {
  instance_id                    = "%[1]s"
  publication_name               = "%[2]s"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "true"
    snapshot_always_available    = "false"
    replicate_ddl                = "true"
    allow_initialize_from_backup = "false"
  }

  job_schedule{
    job_schedule_type = "recurring"

    one_time_occurrence {
      active_start_date = "2026-05-06"
      active_start_time = "07:20:30"
    }

    frequency {
      freq_type     = "daily"
      freq_interval = 2
    }

    daily_frequency {
      freq_subday_type  = "once"
      active_start_time = "02:00:00"
    }

    duration {
      active_start_date = "2020-01-01"
      active_end_date   = "2099-12-31"
    }
  }

  is_select_all_table = "true"
}
`, acceptance.HW_RDS_INSTANCE_ID, rName)
}

func testAccPublication_with_recurring_weekly(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_publication" "test" {
  instance_id                    = "%[1]s"
  publication_name               = "%[2]s"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "true"
    snapshot_always_available    = "false"
    replicate_ddl                = "true"
    allow_initialize_from_backup = "false"
  }

  job_schedule{
    job_schedule_type = "recurring"

    one_time_occurrence {
      active_start_date = "2026-05-06"
      active_start_time = "07:20:30"
    }

    frequency {
      freq_type            = "weekly"
      freq_interval        = 2
      freq_interval_weekly = ["Monday", "Tuesday", "Saturday"]
    }

    daily_frequency {
      freq_subday_type     = "multiple"
      active_start_time    = "02:00:00"
      active_end_time      = "20:00:00"
      freq_subday_interval = "10"
      freq_interval_unit   = "minute"
    }

    duration {
      active_start_date = "2010-07-15"
      active_end_date   = "2089-10-20"
    }
  }

  is_select_all_table = "true"
}
`, acceptance.HW_RDS_INSTANCE_ID, rName)
}

func testAccPublication_with_recurring_monthly_day(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_publication" "test" {
  instance_id                    = "%[1]s"
  publication_name               = "%[2]s"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "true"
    snapshot_always_available    = "false"
    replicate_ddl                = "true"
    allow_initialize_from_backup = "false"
  }

  job_schedule{
    job_schedule_type = "recurring"

    one_time_occurrence {
      active_start_date = "2026-05-06"
      active_start_time = "07:20:30"
    }

    frequency {
      freq_type                 = "monthly_day"
      freq_interval             = 2
      freq_interval_weekly      = ["Monday", "Tuesday", "Saturday"]
      freq_interval_day_monthly = 15
    }

    daily_frequency {
      freq_subday_type     = "multiple"
      active_start_time    = "02:00:00"
      active_end_time      = "20:00:00"
      freq_subday_interval = "10"
      freq_interval_unit   = "minute"
    }

    duration {
      active_start_date = "2010-07-15"
      active_end_date   = "2089-10-20"
    }
  }

  is_select_all_table = "true"
}
`, acceptance.HW_RDS_INSTANCE_ID, rName)
}

func testAccPublication_with_recurring_monthly_week(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_publication" "test" {
  instance_id                    = "%[1]s"
  publication_name               = "%[2]s"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "true"
    snapshot_always_available    = "false"
    replicate_ddl                = "true"
    allow_initialize_from_backup = "false"
  }

  job_schedule{
    job_schedule_type = "recurring"

    one_time_occurrence {
      active_start_date = "2026-05-06"
      active_start_time = "07:20:30"
    }

    frequency {
      freq_type                      = "monthly_week"
      freq_interval                  = 2
      freq_interval_monthly          = "weekday"
      freq_relative_interval_monthly = "second"
      freq_interval_weekly           = ["Monday", "Tuesday", "Saturday"]
      freq_interval_day_monthly      = 15
    }

    daily_frequency {
      freq_subday_type     = "multiple"
      active_start_time    = "02:00:00"
      active_end_time      = "20:00:00"
      freq_subday_interval = "10"
      freq_interval_unit   = "minute"
    }

    duration {
      active_start_date = "2010-07-15"
      active_end_date   = "2089-10-20"
    }
  }

  is_select_all_table = "true"
}
`, acceptance.HW_RDS_INSTANCE_ID, rName)
}

func testAccPublication_with_tables_filter(rName string) string {
	return fmt.Sprintf(`

resource "huaweicloud_rds_publication" "test" {
  instance_id                    = "%[1]s"
  publication_name               = "%[2]s"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "true"
    snapshot_always_available    = "false"
    replicate_ddl                = "true"
    allow_initialize_from_backup = "false"
  }

  job_schedule{
    job_schedule_type = "recurring"

    frequency {
      freq_type                      = "monthly_week"
      freq_interval                  = 2
      freq_interval_monthly          = "weekday"
      freq_relative_interval_monthly = "second"
    }

    daily_frequency {
      freq_subday_type     = "multiple"
      active_start_time    = "02:00:00"
      active_end_time      = "20:00:00"
      freq_subday_interval = "10"
      freq_interval_unit   = "minute"
    }

    duration {
      active_start_date = "2010-07-15"
      active_end_date   = "2089-10-20"
    }
  }

  is_select_all_table = "false"

  tables {
    table_name       = "test_table_1"
    schema           = "test_schema_1"
    columns          = ["id", "name", "address"]
    primary_key      = ["id"]
    filter_statement = "WHERE id = '123'"

    filter {
      column    = "id"
      condition = "="
      value     = "123"
    }

    article_properties {
      destination_object_name  = "test_table_1"
      destination_object_owner = "test_schema_1"
      insert_delivery_format   = "insert"
      insert_stored_procedure  = "sp_MSins_test_schema_1test_table_1"
      update_delivery_format   = "update"
      update_stored_procedure  = "sp_MSupd_test_schema_1test_table_1"
      delete_delivery_format   = "do_not_delete"
      delete_stored_procedure  = "sp_MSdel_test_schema_1test_table_1"
    }
  }
}
`, acceptance.HW_RDS_INSTANCE_ID, rName)
}

func testAccPublication_with_tables_filter_filters(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_publication" "test" {
  instance_id                    = "%[1]s"
  publication_name               = "%[2]s"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "true"
    snapshot_always_available    = "false"
    replicate_ddl                = "true"
    allow_initialize_from_backup = "false"
  }

  job_schedule{
    job_schedule_type = "recurring"

    frequency {
      freq_type                      = "monthly_week"
      freq_interval                  = 2
      freq_interval_monthly          = "weekday"
      freq_relative_interval_monthly = "second"
    }

    daily_frequency {
      freq_subday_type     = "multiple"
      active_start_time    = "02:00:00"
      active_end_time      = "20:00:00"
      freq_subday_interval = "3"
      freq_interval_unit   = "minute"
    }

    duration {
      active_start_date = "2010-07-15"
      active_end_date   = "2089-10-20"
    }
  }

  is_select_all_table = "false"

  tables {
    table_name       = "test_table_1"
    schema           = "test_schema_1"
    columns          = ["id", "name", "address"]
    primary_key      = ["id"]
    filter_statement = "WHERE id = '111' AND (id = '222' AND name = '123')"

    filter {
      relation = "AND"
      filters  = [
        jsonencode({
          "column": "id",
          "condition": "=",
          "value": "111"
        }),
        jsonencode({
          "relation": "AND",
          "filters": [
            {
              "column": "id",
              "condition": "=",
              "value": "222"
            },
            {
              "column": "name",
              "condition": "=",
              "value": "123"
            }
          ]
        })
      ]
    }

    article_properties {
      destination_object_name  = "test_table_1"
      destination_object_owner = "test_schema_1"
      insert_delivery_format   = "call_procedure"
      insert_stored_procedure  = "sp_MSins_test_schema_1test_table_1"
      update_delivery_format   = "scall_procedure"
      update_stored_procedure  = "sp_MSupd_test_schema_1test_table_1"
      delete_delivery_format   = "call_procedure"
      delete_stored_procedure  = "sp_MSdel_test_schema_1test_table_1"
    }
  }
}
`, acceptance.HW_RDS_INSTANCE_ID, rName)
}

func testPublicationResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}
