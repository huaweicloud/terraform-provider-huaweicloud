package dli

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v1/tables"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
)

func getDliTableResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Dli v1 client, err=%s", err)
	}
	databaseName, tableName := dli.ParseTableInfoFromId(state.Primary.ID)
	return tables.Get(client, databaseName, tableName)
}

// check the dli table
func TestAccResourceDliTable_basic(t *testing.T) {
	var tableObj tables.CreateTableOpts
	resourceName := "huaweicloud_dli_table.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&tableObj,
		getDliTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliTableResource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "data_location", tables.TableTypeDLI),
					// The description returned by the API may contain extra characters "}" at the end.
					resource.TestMatchResourceAttr(resourceName, "description", regexp.MustCompile(`^dli table test?}$`)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"conf", "schema", "rows", "job_mode"},
			},
		},
	})
}

func testAccDliTableResource_basic(name string) string {
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
    name = "name"
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

func TestAccResourceDliTable_OBS(t *testing.T) {
	var tableObj tables.CreateTableOpts
	resourceName := "huaweicloud_dli_table.test"
	name := acceptance.RandomAccResourceName()
	obsBucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&tableObj,
		getDliTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliTableResource_OBS(name, obsBucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "data_location", tables.TableTypeOBS),
					// The description returned by the API may contain extra characters "}" at the end.
					resource.TestMatchResourceAttr(resourceName, "description", regexp.MustCompile(`^dli table test?}$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDliTableResource_OBS(name string, obsBucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}


resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "user/data/user.csv"
  content      = "Jason,Tokyo"
  content_type = "text/plain"
}

resource "huaweicloud_dli_database" "test" {
  name        = "%s"
  description = "For terraform acc test"
}

resource "huaweicloud_dli_table" "test" {
  database_name   = huaweicloud_dli_database.test.name
  name            = "%s"
  data_location   = "OBS"
  description     = "dli table test"
  data_format     = "csv"
  bucket_location = "obs://${huaweicloud_obs_bucket_object.test.bucket}/user/data"

  columns {
    name        = "name"
    type        = "string"
    description = "person name"
  }

  columns {
    name         = "addrss"
    type         = "string"
    description  = "home address"
    is_partition = true
  }

}
`, obsBucketName, name, name)
}
