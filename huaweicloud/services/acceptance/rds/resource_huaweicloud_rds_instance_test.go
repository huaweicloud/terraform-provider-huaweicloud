package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

func TestAccRdsInstance_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "09:00"),
					resource.TestCheckResourceAttr(resourceName, "private_dns_name_prefix", "terraformTest"),
					resource.TestCheckResourceAttr(resourceName, "minor_version_auto_upgrade_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "private_dns_names.0"),
				),
			},
			{
				Config: testAccRdsInstance_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_updated"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.230"),
					resource.TestCheckResourceAttr(resourceName, "private_ips.0", "192.168.0.230"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8636"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "15:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "17:00"),
					resource.TestCheckResourceAttr(resourceName, "private_dns_name_prefix", "terraformTestUpdate"),
					resource.TestCheckResourceAttr(resourceName, "minor_version_auto_upgrade_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "db.0.password"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
					"status",
					"availability_zone",
					"slow_log_show_original_status",
				},
			},
		},
	})
}

func TestAccRdsInstance_ha(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_ha(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "async"),
					resource.TestCheckResourceAttr(resourceName, "switch_strategy", "availability"),
				),
			},
			{
				Config: testAccRdsInstance_ha_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "switch_strategy", "reliability"),
				),
			},
		},
	})
}

func TestAccRdsInstance_mysql(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_mysql_step1(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "400"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "15"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "back_log"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "2000"),
					resource.TestCheckResourceAttr(resourceName, "binlog_retention_hours", "12"),
					resource.TestCheckResourceAttr(resourceName, "seconds_level_monitoring_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "seconds_level_monitoring_interval", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_used_space.#"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_used_space.0.node_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_used_space.0.used"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_status"),
				),
			},
			{
				Config: testAccRdsInstance_mysql_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.1.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "500"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "20"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3308"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "connect_timeout"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "14"),
					resource.TestCheckResourceAttr(resourceName, "binlog_retention_hours", "0"),
					resource.TestCheckResourceAttr(resourceName, "seconds_level_monitoring_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "seconds_level_monitoring_interval", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "db.0.password"),
				),
			},
			{
				Config: testAccRdsInstance_mysql_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "binlog_retention_hours", "6"),
					resource.TestCheckResourceAttr(resourceName, "seconds_level_monitoring_enabled", "false"),
				),
			},
		},
	})
}

func TestAccRdsInstance_mysql_power_action(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_mysql_power_action(name, "OFF", []string{"SHUTDOWN"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckOutput("instance_status_contains", "true"),
				),
			},
			{
				Config: testAccRdsInstance_mysql_power_action(name, "ON", []string{"ACTIVE", "BACKING UP"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
				),
			},
			{
				Config: testAccRdsInstance_mysql_power_action(name, "ON", []string{"ACTIVE", "BACKING UP"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckOutput("instance_status_contains", "true"),
				),
			},
			{
				Config: testAccRdsInstance_mysql_power_action(name, "REBOOT", []string{"ACTIVE", "BACKING UP"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
				),
			},
			{
				Config: testAccRdsInstance_mysql_power_action(name, "REBOOT", []string{"ACTIVE", "BACKING UP"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckOutput("instance_status_contains", "true"),
				),
			},
		},
	})
}

func TestAccRdsInstance_sqlserver(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_sqlserver(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "collation", "Chinese_PRC_CI_AS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
					resource.TestCheckResourceAttr(resourceName, "tde_enabled", "true"),
				),
			},
			{
				Config: testAccRdsInstance_sqlserver_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "collation", "Chinese_PRC_CI_AI"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
				),
			},
		},
	})
}

func TestAccRdsInstance_sqlserver_msdtc_hosts(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_sqlserver_msdtcHosts(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "msdtc_hosts.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "msdtc_hosts.0.ip",
						"huaweicloud_compute_instance.ecs_1", "access_ip_v4"),
					resource.TestCheckResourceAttr(resourceName, "msdtc_hosts.0.host_name", "msdtc-host-name-1"),
					resource.TestCheckResourceAttrSet(resourceName, "msdtc_hosts.0.id"),
				),
			},
			{
				Config: testAccRdsInstance_sqlserver_msdtcHosts_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "msdtc_hosts.0.ip",
						"huaweicloud_compute_instance.ecs_2", "access_ip_v4"),
					resource.TestCheckResourceAttr(resourceName, "msdtc_hosts.0.host_name", "msdtc-host-name-2"),
				),
			},
		},
	})
}

func TestAccRdsInstance_mariadb(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_mariadb(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
					resource.TestCheckResourceAttrSet(resourceName, "db.0.password"),
				),
			},
		},
	})
}

func TestAccRdsInstance_prePaid(t *testing.T) {
	var (
		instance instances.RdsInstanceResponse

		resourceType = "huaweicloud_rds_instance"
		resourceName = "huaweicloud_rds_instance.test"
		name         = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_prePaid(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttrPair(resourceName, "enterprise_project_id",
						"huaweicloud_enterprise_project.test.0", "id"),
				),
			},
			{
				Config: testAccRdsInstance_prePaid_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
					resource.TestCheckResourceAttrPair(resourceName, "enterprise_project_id",
						"huaweicloud_enterprise_project.test.1", "id"),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_mysql(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test_backup"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_mysql(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "400"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "15"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_sqlserver(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test_backup"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_sqlserver(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_pg(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test_backup"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_pg(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8732"),
				),
			},
		},
	})
}

func TestAccRdsInstance_change_billing_mode_to_prepaid(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_change_billing_mode_to_prepaid(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
				),
			},
			{
				Config: testAccRdsInstance_change_billing_mode_to_prepaid_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
				),
			},
		},
	})
}

func testAccCheckRdsInstanceDestroy(rsType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating rds client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != rsType {
				continue
			}

			id := rs.Primary.ID
			instance, err := rds.GetRdsInstanceByID(client, id)
			if err != nil {
				return err
			}
			if instance.Id != "" {
				return fmt.Errorf("%s (%s) still exists", rsType, id)
			}
		}
		return nil
	}
}

func testAccCheckRdsInstanceExists(name string, instance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating rds client: %s", err)
		}

		found, err := rds.GetRdsInstanceByID(client, id)
		if err != nil {
			return fmt.Errorf("error checking %s exist, err=%s", name, err)
		}
		if found.Id == "" {
			return fmt.Errorf("resource %s does not exist", name)
		}

		instance = found
		return nil
	}
}

func testAccRdsInstance_base() string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}
`)
}

func testAccRdsInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name                               = "%[2]s"
  description                        = "test_description"
  flavor                             = "rds.pg.n1.large.2"
  availability_zone                  = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id                  = data.huaweicloud_networking_secgroup.test.id
  subnet_id                          = data.huaweicloud_vpc_subnet.test.id
  vpc_id                             = data.huaweicloud_vpc.test.id
  time_zone                          = "UTC+08:00"
  maintain_begin                     = "06:00"
  maintain_end                       = "09:00"
  private_dns_name_prefix            = "terraformTest"
  slow_log_show_original_status      = "on"
  minor_version_auto_upgrade_enabled = true

  db {
    type     = "PostgreSQL"
    version  = "12"
    port     = 8634
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstance_base(), name)
}

// name, volume.size, backup_strategy, flavor, tags and password will be updated
func testAccRdsInstance_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test_update" {
  name                 = "%s"
  delete_default_rules = true
}

resource "huaweicloud_rds_instance" "test" {
  name                               = "%[2]s-update"
  flavor                             = "rds.pg.n1.large.2"
  availability_zone                  = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id                  = huaweicloud_networking_secgroup.test_update.id
  subnet_id                          = data.huaweicloud_vpc_subnet.test.id
  vpc_id                             = data.huaweicloud_vpc.test.id
  time_zone                          = "UTC+08:00"
  fixed_ip                           = "192.168.0.230"
  maintain_begin                     = "15:00"
  maintain_end                       = "17:00"
  private_dns_name_prefix            = "terraformTestUpdate"
  slow_log_show_original_status      = "off"
  minor_version_auto_upgrade_enabled = false

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8636
  }
  volume {
    type = "CLOUDSSD"
    size = 100
  }
  backup_strategy {
    start_time = "09:00-10:00"
    keep_days  = 2
  }

  tags = {
    key1 = "value"
    foo  = "bar_updated"
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_ha(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name                = "%[2]s"
  flavor              = "rds.pg.n1.large.2.ha"
  security_group_id   = data.huaweicloud_networking_secgroup.test.id
  subnet_id           = data.huaweicloud_vpc_subnet.test.id
  vpc_id              = data.huaweicloud_vpc.test.id
  time_zone           = "UTC+08:00"
  ha_replication_mode = "async"
  switch_strategy     = "availability"
  availability_zone   = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
  ]

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8634
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_ha_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name                = "%[2]s"
  flavor              = "rds.pg.n1.large.2.ha"
  security_group_id   = data.huaweicloud_networking_secgroup.test.id
  subnet_id           = data.huaweicloud_vpc_subnet.test.id
  vpc_id              = data.huaweicloud_vpc.test.id
  time_zone           = "UTC+08:00"
  ha_replication_mode = "sync"
  switch_strategy     = "reliability"
  availability_zone   = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[2],
  ]

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8634
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstance_base(), name)
}

// if the instance flavor has been changed, then a temp instance will be kept for 12 hours,
// the binding relationship between instance and security group or subnet cannot be unbound
// when deleting the instance in this period time, so we cannot create a new vpc, subnet and
// security group in the test case, otherwise, they cannot be deleted when destroy the resource
func testAccRdsInstance_mysql_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name                   = "%[2]s"
  flavor                 = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id      = data.huaweicloud_networking_secgroup.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  vpc_id                 = data.huaweicloud_vpc.test.id
  availability_zone      = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable             = true  
  binlog_retention_hours = "12"
  read_write_permissions = "readonly"

  seconds_level_monitoring_enabled  = true
  seconds_level_monitoring_interval = 1

  db {
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
    period     = 1
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 400
    trigger_threshold = 15
  }

  parameters {
    name  = "back_log"
    value = "2000"
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_mysql_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name                   = "%[3]s"
  flavor                 = data.huaweicloud_rds_flavors.test.flavors[1].name
  security_group_id      = data.huaweicloud_networking_secgroup.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  vpc_id                 = data.huaweicloud_vpc.test.id
  availability_zone      = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable             = false
  param_group_id         = huaweicloud_rds_parametergroup.pg_1.id
  binlog_retention_hours = "0"
  read_write_permissions = "readwrite"

  seconds_level_monitoring_enabled  = true
  seconds_level_monitoring_interval = 5

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3308
  }

  backup_strategy {
    start_time = "18:15-19:15"
    keep_days  = 5
    period     = 3
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 500
    trigger_threshold = 20
  }

  parameters {
    name  = "connect_timeout"
    value = "14"
  }
}
`, testAccRdsInstance_base(), testAccRdsConfig_basic(name), name)
}

func testAccRdsInstance_mysql_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name                   = "%[3]s"
  flavor                 = data.huaweicloud_rds_flavors.test.flavors[1].name
  security_group_id      = data.huaweicloud_networking_secgroup.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  vpc_id                 = data.huaweicloud_vpc.test.id
  availability_zone      = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable             = false
  param_group_id         = huaweicloud_rds_parametergroup.pg_1.id
  binlog_retention_hours = "6"
  read_write_permissions = "readwrite"

  seconds_level_monitoring_enabled = false

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3308
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  parameters {
    name  = "connect_timeout"
    value = "14"
  }
}
`, testAccRdsInstance_base(), testAccRdsConfig_basic(name), name)
}

func testAccRdsInstance_mysql_power_action(name, action string, status []string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  power_action      = "%[3]s"

  db {
    type    = "MySQL"
    version = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

output "instance_status_contains" {
  value = contains(split(",", "%[4]s"), huaweicloud_rds_instance.test.status)
}
`, testAccRdsInstance_base(), name, action, strings.Join(status, ","))
}

func testAccRdsInstance_sqlserver(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2017_EE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AS"
  tde_enabled       = true

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2017_EE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_sqlserver_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2017_EE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AI"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2017_EE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_sqlserver_msdtcHosts_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_sqlserver_msdtcHosts(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "ecs_1" {
  name               = "%[2]s_ecs_1"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_rds_instance" "test" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AS"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  msdtc_hosts {
    ip        = huaweicloud_compute_instance.ecs_1.access_ip_v4
    host_name = "msdtc-host-name-1"
  }
}
`, testAccRdsInstance_sqlserver_msdtcHosts_base(name), name)
}

func testAccRdsInstance_sqlserver_msdtcHosts_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "ecs_1" {
  name               = "%[2]s_ecs_1"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_compute_instance" "ecs_2" {
  name               = "%[2]s_ecs_2"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_rds_instance" "test" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AS"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  msdtc_hosts {
    ip        = huaweicloud_compute_instance.ecs_2.access_ip_v4
    host_name = "msdtc-host-name-2"
  }
}
`, testAccRdsInstance_sqlserver_msdtcHosts_base(name), name)
}

func testAccRdsInstance_mariadb(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MariaDB"
  db_version    = "10.5"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    password = "Huangwei!120521"
    type     = "MariaDB"
    version  = "10.5"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_prePaid(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_enterprise_project" "test" {
  count = 2

  name        = "%[2]s_${count.index}"
  description = "terraform test"
}

resource "huaweicloud_rds_instance" "test" {
  depends_on             = [huaweicloud_networking_secgroup_rule.ingress]
  vpc_id                 = data.huaweicloud_vpc.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test.id
  lower_case_table_names = "1"
  enterprise_project_id  = huaweicloud_enterprise_project.test[0].id
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  name      = "%[2]s"
  flavor    = data.huaweicloud_rds_flavors.test.flavors[0].name
  collation = "Chinese_PRC_CI_AS"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8638
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_prePaid_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 8
}

resource "huaweicloud_enterprise_project" "test" {
  count = 2

  name        = "%[2]s_${count.index}"
  description = "terraform test"
}

resource "huaweicloud_rds_instance" "test" {
  depends_on             = [huaweicloud_networking_secgroup_rule.ingress]
  vpc_id                 = data.huaweicloud_vpc.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test.id
  lower_case_table_names = "1"
  enterprise_project_id  = huaweicloud_enterprise_project.test[1].id
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  name      = "%[2]s"
  flavor    = data.huaweicloud_rds_flavors.test.flavors[0].name
  collation = "Chinese_PRC_CI_AS"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8638
  }

  volume {
    type = "CLOUDSSD"
    size = 60
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_restore_mysql(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable        = true  

  restore {
    instance_id = huaweicloud_rds_backup.test.instance_id
    backup_id   = huaweicloud_rds_backup.test.id
  }

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
    period     = 1
  }

  volume {
    type              = "CLOUDSSD"
    size              = 50
    limit_size        = 400
    trigger_threshold = 15
  }
}
`, testBackup_mysql_basic(name), name)
}

func testAccRdsInstance_restore_sqlserver(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test_backup" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  restore {
    instance_id = huaweicloud_rds_backup.test.instance_id
    backup_id   = huaweicloud_rds_backup.test.id
  }

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testBackup_sqlserver_basic(name), name)
}

func testAccRdsInstance_restore_pg(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  restore {
    instance_id = huaweicloud_rds_backup.test.instance_id
    backup_id   = huaweicloud_rds_backup.test.id
  }

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8732
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testBackup_pg_basic(name), name, pwd)
}

func testAccRdsInstance_change_billing_mode_to_prepaid(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "14"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 8
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8632
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_change_billing_mode_to_prepaid_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "14"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 8
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8632
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}
`, testAccRdsInstance_base(), name)
}
