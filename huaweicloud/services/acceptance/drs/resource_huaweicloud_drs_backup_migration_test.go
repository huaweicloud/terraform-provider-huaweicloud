package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/drs"
)

func getBackupMigrationResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("drs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DRS client: %s", err)
	}

	return drs.QueryMigrationDetail(client, state.Primary.ID)
}

func TestAccResourceBackupMigration_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_drs_backup_migration.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBackupMigrationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsBackupMigrationInstanceId(t)
			acceptance.TestAccPreCheckDrsBackupMigrationBucketName(t)
			acceptance.TestAccPreCheckDrsBackupMigrationBucketPath(t)
			acceptance.TestAccPreCheckDrsBackupMigrationFileName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupMigration_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "base_info.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "base_info.0.engine_type", "sqlserver"),
					resource.TestCheckResourceAttr(resourceName, "base_info.0.description", "test backup migration"),
					resource.TestCheckResourceAttr(resourceName, "target_db_info.0.target_instance_id",
						acceptance.HW_DRS_BACKUP_MIGRATION_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "backup_info.0.file_source", "OBS"),
					resource.TestCheckResourceAttr(resourceName, "backup_info.0.bucket_name",
						acceptance.HW_DRS_BACKUP_MIGRATION_BUCKET_NAME),
					resource.TestCheckResourceAttr(resourceName, "backup_info.0.files.0.name",
						acceptance.HW_DRS_BACKUP_MIGRATION_FILE_NAME),
					resource.TestCheckResourceAttr(resourceName, "backup_info.0.files.0.obs_path",
						acceptance.HW_DRS_BACKUP_MIGRATION_BUCKET_PATH),
					resource.TestCheckResourceAttr(resourceName, "options.0.is_last_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.is_precheck", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.recovery_mode", "full"),
					resource.TestCheckResourceAttr(resourceName, "options.0.is_cover", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.is_default_restore", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.is_delete_backup_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.db_names.0", "db-test-name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: testAccBackupMigration_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "base_info.0.name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "base_info.0.description", "test backup migration update"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_info.0.files"},
			},
		},
	})
}

func testAccBackupMigration_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_backup_migration" "test" {
  base_info {
    name        = "%[1]s"
    engine_type = "sqlserver"
    description = "test backup migration"
  }

  target_db_info {
    target_instance_id = "%[2]s"
  }

  backup_info {
    file_source = "OBS"
    bucket_name = "%[3]s"

    files {
      name     = "%[4]s"
      obs_path = "%[5]s"
    }
  }

  options {
    is_cover              = "false"
    is_default_restore    = "true"
    is_delete_backup_file = "false"
    is_last_backup        = true
    is_precheck           = true
    recovery_mode         = "full"
    db_names              = ["db-test-name"]
  }
}
`, name, acceptance.HW_DRS_BACKUP_MIGRATION_INSTANCE_ID, acceptance.HW_DRS_BACKUP_MIGRATION_BUCKET_NAME,
		acceptance.HW_DRS_BACKUP_MIGRATION_FILE_NAME, acceptance.HW_DRS_BACKUP_MIGRATION_BUCKET_PATH)
}

func testAccBackupMigration_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_backup_migration" "test" {
  base_info {
    name        = "%[1]s-update"
    engine_type = "sqlserver"
    description = "test backup migration update"
  }

  target_db_info {
    target_instance_id = "%[2]s"
  }

  backup_info {
    file_source = "OBS"
    bucket_name = "%[3]s"

    files {
      name     = "%[4]s"
      obs_path = "%[5]s"
    }
  }

  options {
    is_cover              = "false"
    is_default_restore    = "true"
    is_delete_backup_file = "false"
    is_last_backup        = true
    is_precheck           = true
    recovery_mode         = "full"
    db_names              = ["db-test-name"]
  }
}
`, name, acceptance.HW_DRS_BACKUP_MIGRATION_INSTANCE_ID, acceptance.HW_DRS_BACKUP_MIGRATION_BUCKET_NAME,
		acceptance.HW_DRS_BACKUP_MIGRATION_FILE_NAME, acceptance.HW_DRS_BACKUP_MIGRATION_BUCKET_PATH)
}
