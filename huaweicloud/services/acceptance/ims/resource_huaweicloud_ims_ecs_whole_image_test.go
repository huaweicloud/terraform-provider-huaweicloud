package ims

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEcsWholeImage_basic(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = acceptance.RandomAccResourceName()
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_ims_ecs_whole_image.test"
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
				Config: testAccEcsWholeImage_basic(rName, defaultEpsId, 2048, 4096),
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
					resource.TestCheckResourceAttrSet(resourceName, "backup_id"),
					resource.TestCheckResourceAttrSet(resourceName, "min_disk"),
					resource.TestCheckResourceAttrSet(resourceName, "disk_format"),
					resource.TestCheckResourceAttrSet(resourceName, "data_origin"),
					resource.TestCheckResourceAttrSet(resourceName, "os_version"),
					resource.TestMatchResourceAttr(resourceName, "active_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "vault_id", "huaweicloud_cbr_vault.test", "id"),
				),
			},
			{
				Config: testAccEcsWholeImage_update1(rName, rNameUpdate, migrateEpsId, 1024, 2048),
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
				Config: testAccEcsWholeImage_update2(rName, rNameUpdate, defaultEpsId, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
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

func testAccEcsWholeImage_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200
}
`, testAccEcsSystemImage_base(rName), rName)
}

func testAccEcsWholeImage_basic(rName, defaultEpsId string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_ecs_whole_image" "test" {
  name                  = "%[2]s"
  instance_id           = huaweicloud_compute_instance.test.id
  vault_id              = huaweicloud_cbr_vault.test.id
  description           = "terraform description test"
  enterprise_project_id = "%[3]s"
  min_ram               = %[4]d
  max_ram               = %[5]d

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEcsWholeImage_base(rName), rName, defaultEpsId, minRAM, maxRAM)
}

func testAccEcsWholeImage_update1(rName, rNameUpdate, migrateEpsId string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_ecs_whole_image" "test" {
  name                  = "%[2]s"
  instance_id           = huaweicloud_compute_instance.test.id
  vault_id              = huaweicloud_cbr_vault.test.id
  enterprise_project_id = "%[3]s"
  min_ram               = %[4]d
  max_ram               = %[5]d

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccEcsWholeImage_base(rName), rNameUpdate, migrateEpsId, minRAM, maxRAM)
}

func testAccEcsWholeImage_update2(rName, rNameUpdate, defaultEpsId string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_ecs_whole_image" "test" {
  name                  = "%[2]s"
  instance_id           = huaweicloud_compute_instance.test.id
  vault_id              = huaweicloud_cbr_vault.test.id
  enterprise_project_id = "%[3]s"
  min_ram               = %[4]d
  max_ram               = %[5]d

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccEcsWholeImage_base(rName), rNameUpdate, defaultEpsId, minRAM, maxRAM)
}

// When using fully automated script testing to set `is_delete_backup` to true, it can cause CBR vault resource to be
// deleted very slowly, resulting in a timeout error.
// So here we use environment variables to inject the `is_delete_backup` parameter to complete the testing.
func TestAccEcsWholeImage_withDeleteBackup(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_ims_ecs_whole_image.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImsImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need setting a CBR vault ID.
			acceptance.TestAccPreCheckImsVaultId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEcsWholeImage_withDeleteBackup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "is_delete_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "vault_id", acceptance.HW_IMS_VAULT_ID),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_compute_instance.test", "id"),
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

func testAccEcsWholeImage_withDeleteBackup_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_ecs_whole_image" "test" {
  name             = "%[2]s"
  instance_id      = huaweicloud_compute_instance.test.id
  vault_id         = "%[3]s"
  is_delete_backup = true
}
`, testAccEcsSystemImage_base(rName), rName, acceptance.HW_IMS_VAULT_ID)
}
