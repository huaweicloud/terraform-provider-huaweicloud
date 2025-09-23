package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sms"
)

func getMigrationProjectResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("sms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	return sms.GetMigrationProject(client, state.Primary.ID)
}

func TestAccResourceMigrationProject_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sms_migration_project.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getMigrationProjectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMigrationProject_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "region", "cn-north-9"),
					resource.TestCheckResourceAttr(resourceName, "use_public_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "exist_server", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "MIGRATE_BLOCK"),
					resource.TestCheckResourceAttr(resourceName, "syncing", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "is_default"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "start_target_server"),
					resource.TestCheckResourceAttrSet(resourceName, "speed_limit"),
					resource.TestCheckResourceAttrSet(resourceName, "use_public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "exist_server"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project"),
					resource.TestCheckResourceAttrSet(resourceName, "syncing"),
					resource.TestCheckResourceAttrSet(resourceName, "start_network_check"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMigrationProject_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "region", "cn-north-4"),
					resource.TestCheckResourceAttr(resourceName, "use_public_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "exist_server", "false"),
					resource.TestCheckResourceAttr(resourceName, "type", "MIGRATE_FILE"),
					resource.TestCheckResourceAttr(resourceName, "syncing", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "is_default"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "start_target_server"),
					resource.TestCheckResourceAttrSet(resourceName, "speed_limit"),
					resource.TestCheckResourceAttrSet(resourceName, "use_public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "exist_server"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project"),
					resource.TestCheckResourceAttrSet(resourceName, "syncing"),
					resource.TestCheckResourceAttrSet(resourceName, "start_network_check"),
				),
			},
		},
	})
}

func testAccMigrationProject_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_sms_migration_project" "test" {
  name          = "%[1]s"
  region        = "cn-north-9"
  use_public_ip = true
  exist_server  = true
  type          = "MIGRATE_BLOCK"
  syncing       = true
}
`, name)
}

func testAccMigrationProject_updated(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_sms_migration_project" "test" {
  name          = "%[1]s"
  region        = "cn-north-4"
  use_public_ip = false
  exist_server  = false
  type          = "MIGRATE_FILE"
  syncing       = false
}
`, name)
}
