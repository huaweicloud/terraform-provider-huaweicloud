package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sms/v3/tasks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getMigrationTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.SmsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	return tasks.Get(client, state.Primary.ID)
}

func TestAccMigrationTask_basic(t *testing.T) {
	var migration tasks.MigrateTask
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sms_task.migration"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&migration,
		getMigrationTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMigrationTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "state", "READY"),
					resource.TestCheckResourceAttr(resourceName, "use_public_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "speed_limit.0.start", "00:00"),
					resource.TestCheckResourceAttr(resourceName, "speed_limit.0.end", "23:59"),
					resource.TestCheckResourceAttr(resourceName, "configurations.0.config_status", "testA"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.#"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.0.size"),
					resource.TestCheckResourceAttrPair(resourceName, "vm_template_id",
						"huaweicloud_sms_server_template.test", "id"),
				),
			},
			{
				Config: testAccMigrationTask_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "state", "READY"),
					resource.TestCheckResourceAttr(resourceName, "use_public_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "speed_limit.0.start", "01:00"),
					resource.TestCheckResourceAttr(resourceName, "speed_limit.0.end", "22:59"),
					resource.TestCheckResourceAttr(resourceName, "configurations.0.config_status", "testB"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.#"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.0.size"),
					resource.TestCheckResourceAttrPair(resourceName, "vm_template_id",
						"huaweicloud_sms_server_template.test", "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"use_public_ip", "syncing", "action"},
			},
		},
	})
}

func testAccMigrationTask_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_sms_task" "migration" {
  type             = "MIGRATE_FILE"
  os_type          = "LINUX"
  source_server_id = data.huaweicloud_sms_source_servers.source.servers[0].id
  vm_template_id   = huaweicloud_sms_server_template.test.id

  speed_limit {
    start                = "00:00"
    end                  = "23:59"
    speed                = 1000
    over_speed_threshold = 50.0
  }

  configurations {
    config_key    = "LINUX_CPU_LIMIT"
    config_value  = "60"
    config_status = "testA"
  }
}
`, testAccMigrationTask_base(name))
}

func testAccMigrationTask_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_sms_task" "migration" {
  type             = "MIGRATE_FILE"
  os_type          = "LINUX"
  source_server_id = data.huaweicloud_sms_source_servers.source.servers[0].id
  vm_template_id   = huaweicloud_sms_server_template.test.id

  speed_limit {
    start                = "01:00"
    end                  = "22:59"
    speed                = 100
    over_speed_threshold = 60.0
  }

  configurations {
    config_key    = "LINUX_CPU_LIMIT"
    config_value  = "50"
    config_status = "testB"
  }
}
`, testAccMigrationTask_base(name))
}

func testAccMigrationTask_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_sms_source_servers" "source" {
  name = "%s"
}

resource "huaweicloud_sms_server_template" "test" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`, acceptance.HW_SMS_SOURCE_SERVER, name)
}
