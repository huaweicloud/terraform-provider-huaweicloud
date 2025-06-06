package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsConfigurationCopy_basic(t *testing.T) {
	var obj interface{}
	sourceName := acceptance.RandomAccResourceName()
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_parametergroup_copy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getConfiguration,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsConfigurationCopy_basic(sourceName, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "8.0"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccRdsConfigurationCopy_update(sourceName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_id", "values"},
			},
		},
	})
}

func testAccRdsConfigurationCopy_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_parametergroup" "test" {
  name        = "%s"
  description = "description_1"

  datastore {
    type    = "mysql"
    version = "8.0"
  }
}
`, rName)
}

func testAccRdsConfigurationCopy_basic(sourceName, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_parametergroup_copy" "test" {
  config_id   = huaweicloud_rds_parametergroup.test.id
  name        = "%[2]s"
  description = "test description"

  values = {
    auto_increment_increment     = "2"
    binlog_rows_query_log_events = "ON"
  }
}
`, testAccRdsConfigurationCopy_base(sourceName), rName)
}

func testAccRdsConfigurationCopy_update(sourceName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_parametergroup_copy" "test" {
  config_id   = huaweicloud_rds_parametergroup.test.id
  name        = "%[2]s"
  description = ""

  values = {
    bulk_insert_buffer_size = "10"
    connect_timeout         = "10"
  }
}
`, testAccRdsConfigurationCopy_base(sourceName), updateName)
}
