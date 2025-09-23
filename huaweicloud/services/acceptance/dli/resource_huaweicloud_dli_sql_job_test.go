package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v1/sqljob"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDliSQLJobResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Dli v1 client, err=%s", err)
	}
	return sqljob.Status(client, state.Primary.ID)
}

// check the DDL sql
func TestAccResourceDliSqlJob_basic(t *testing.T) {
	var sqlJobObj sqljob.SqlJobOpts
	resourceName := "huaweicloud_dli_sql_job.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&sqlJobObj,
		getDliSQLJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliSQLQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSQLJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlJobBaseResource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sql", fmt.Sprint("DESC ", name)),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "job_type", "DDL"),
					resource.TestCheckResourceAttr(resourceName, "queue_name", acceptance.HW_DLI_SQL_QUEUE_NAME),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rows", "schema"},
			},
		},
	})
}

func TestAccResourceDliSqlJob_async(t *testing.T) {
	var sqlJobObj sqljob.SqlJobOpts
	resourceName := "huaweicloud_dli_sql_job.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&sqlJobObj,
		getDliSQLJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSQLJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSQLJobResource_aync(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sql", fmt.Sprint("SELECT * FROM ", name)),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "queue_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "job_type", "QUERY"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rows", "schema", "conf", "duration", "status"},
			},
		},
	})
}

func testAccSqlJobBaseResource_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_sql_job" "test" {
  sql           = "DESC ${huaweicloud_dli_table.test.name}"
  database_name = huaweicloud_dli_database.test.name
  queue_name    = "%s"
}
`, testAccSQLJobBaseResource(name), acceptance.HW_DLI_SQL_QUEUE_NAME)
}

func testAccSQLJobBaseResource(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name        = "%s"
  description = "For terraform acc test"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%s"
  data_location = "DLI"

  columns {
    name        = "name"
    type        = "string"
    description = "person name"
  }

  columns {
    name        = "addrss"
    type        = "string"
    description = "home address"
  }
}
`, name, name)
}

func testAccSQLJobResource_aync(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_sql_job" "test" {
  sql           = "SELECT * FROM ${huaweicloud_dli_table.test.name}"
  database_name = huaweicloud_dli_database.test.name

  conf {
	spark_sql_max_records_per_file                = 1
    spark_sql_auto_broadcast_join_threshold       = -1
    spark_sql_shuffle_partitions                  = 4096
    spark_sql_dynamic_partition_overwrite_enabled = true
    spark_sql_files_max_partition_bytes           = 134216704
    spark_sql_bad_records_path                    = "test"
    dli_sql_sqlasync_enabled                      = true
    dli_sql_job_timeout                           = 300
  }
}
`, testAccSQLJobBaseResource(name))
}

func testAccCheckDliSQLJobDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Dli client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dli_sql_job" {
			continue
		}

		res, err := sqljob.Status(client, rs.Primary.ID)
		if err == nil && res != nil && (res.Status != sqljob.JobStatusCancelled &&
			res.Status != sqljob.JobStatusFinished && res.Status != sqljob.JobStatusFailed) {
			return fmt.Errorf("huaweicloud_dli_sql_job still exists:%s,%+v,%+v", rs.Primary.ID, err, res)
		}
	}

	return nil
}
