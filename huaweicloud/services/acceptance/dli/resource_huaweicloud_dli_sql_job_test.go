package dli

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v1/sqljob"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDliSqlJobResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating Dli v1 client, err=%s", err)
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
		getDliSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSqlJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlJobBaseResource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sql", fmt.Sprint("DESC ", name)),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "job_type", "DDL"),
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

func TestAccResourceDliSqlJob_query(t *testing.T) {
	var sqlJobObj sqljob.SqlJobOpts
	resourceName := "huaweicloud_dli_sql_job.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&sqlJobObj,
		getDliSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSqlJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlJobBaseResource_query(name),
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
		getDliSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSqlJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlJobResource_aync(name),
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
}
`, testAccSqlJobBaseResource(name))
}

func testAccSqlJobBaseResource(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name        = "%s"
  description = "For terraform acc test"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%s"
  data_location = "DLI"
  description   = "dli table test"

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

func testAccSqlJobBaseResource_query(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_sql_job" "test" {
  sql           = "SELECT * FROM ${huaweicloud_dli_table.test.name}"
  database_name = huaweicloud_dli_database.test.name

}
`, testAccSqlJobBaseResource(name))
}

func testAccSqlJobResource_aync(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_sql_job" "test" {
  sql           = "SELECT * FROM ${huaweicloud_dli_table.test.name}"
  database_name = huaweicloud_dli_database.test.name

  conf {
    dli_sql_sqlasync_enabled = true
  }
}
`, testAccSqlJobBaseResource(name))
}

func testAccCheckDliSqlJobDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("error creating Dli client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dli_sql_job" {
			continue
		}

		res, err := sqljob.Status(client, rs.Primary.ID)
		if err == nil && res != nil && (res.Status != sqljob.JobStatusCancelled &&
			res.Status != sqljob.JobStatusFinished && res.Status != sqljob.JobStatusFailed) {
			return fmtp.Errorf("huaweicloud_dli_sql_job still exists:%s,%+v,%+v", rs.Primary.ID, err, res)
		}
	}

	return nil
}
