package ims

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCbrWholeImage_basic(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = acceptance.RandomAccResourceName()
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_ims_cbr_whole_image.test"
		defaultEpsId = "0"
		migrateEpsId = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImsImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need setting a non default enterprise project ID.
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCbrWholeImage_basic(rName, 2048, 4096),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform description test"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "2048"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttrSet(resourceName, "visibility"),
					resource.TestCheckResourceAttrSet(resourceName, "os_version"),
					resource.TestCheckResourceAttrSet(resourceName, "min_disk"),
					resource.TestCheckResourceAttrSet(resourceName, "disk_format"),
					resource.TestCheckResourceAttrSet(resourceName, "data_origin"),
					resource.TestMatchResourceAttr(resourceName, "active_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrPair(resourceName, "backup_id", "huaweicloud_cbr_checkpoint.test", "backups.0.id"),
				),
			},
			{
				Config: testAccCbrWholeImage_update1(rName, rNameUpdate, migrateEpsId, 1024, 2048),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "2048"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", migrateEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				Config: testAccCbrWholeImage_update2(rName, rNameUpdate, defaultEpsId, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttr(resourceName, "is_delete_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"is_delete_backup",
				},
			},
		},
	})
}

func testAccCbrWholeImage_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200

  resources {
    server_id = huaweicloud_compute_instance.test.id
  }
}

resource "huaweicloud_cbr_checkpoint" "test" {
  vault_id = huaweicloud_cbr_vault.test.id
  name     = "%[2]s"

  backups {
    type        = "OS::Nova::Server"
    resource_id = huaweicloud_compute_instance.test.id
  }
}
`, testAccEcsSystemImage_base(rName), rName)
}

func testAccCbrWholeImage_basic(rName string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_cbr_whole_image" "test" {
  name        = "%[2]s"
  backup_id   = try(tolist(huaweicloud_cbr_checkpoint.test.backups)[0].id, "")
  description = "terraform description test"
  min_ram     = %[3]d
  max_ram     = %[4]d

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCbrWholeImage_base(rName), rName, minRAM, maxRAM)
}

func testAccCbrWholeImage_update1(rName, rNameUpdate, migrateEpsId string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_cbr_whole_image" "test" {
  name                  = "%[2]s"
  backup_id             = try(tolist(huaweicloud_cbr_checkpoint.test.backups)[0].id, "")
  enterprise_project_id = "%[3]s"
  min_ram               = %[4]d
  max_ram               = %[5]d

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccCbrWholeImage_base(rName), rNameUpdate, migrateEpsId, minRAM, maxRAM)
}

func testAccCbrWholeImage_update2(rName, rNameUpdate, defaultEpsId string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_cbr_whole_image" "test" {
  name                  = "%[2]s"
  backup_id             = try(tolist(huaweicloud_cbr_checkpoint.test.backups)[0].id, "")
  enterprise_project_id = "%[3]s"
  min_ram               = %[4]d
  max_ram               = %[5]d
  is_delete_backup      = true

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccCbrWholeImage_base(rName), rNameUpdate, defaultEpsId, minRAM, maxRAM)
}
