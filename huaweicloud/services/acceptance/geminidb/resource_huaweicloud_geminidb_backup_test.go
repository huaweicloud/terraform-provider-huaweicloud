package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/geminidb"
)

func TestAccGeminiDBBackup_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_geminidb_backup.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBBackup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test backup"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
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

func TestAccGeminiDBBackup_withDatabaseTables(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_geminidb_backup.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBBackup_withDatabaseTables(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test backup with database tables"),
					resource.TestCheckResourceAttr(resourceName, "database_tables.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_tables"},
			},
		},
	})
}

func getBackupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geminidb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	backup, err := geminidb.GetBackup(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}
	if backup == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return backup, nil
}

func testAccGeminiDBBackup_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_backup" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "test backup"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, rName)
}

func testAccGeminiDBBackup_withDatabaseTables(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_backup" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "test backup with database tables"

  database_tables {
    database_name = "%[3]s"
    table_names   = ["%[4]s"]
  }
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, rName, acceptance.HW_GEMINIDB_DATABASE_NAME, acceptance.HW_GEMINIDB_TABLE_NAME)
}
