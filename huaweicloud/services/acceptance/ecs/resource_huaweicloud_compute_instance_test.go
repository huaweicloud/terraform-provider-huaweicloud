package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeInstance_basic(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(resourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "volume_attached.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.0.port"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "network.0.source_dest_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "stop_before_destroy", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_eip_on_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "50"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "test111"),
					resource.TestCheckResourceAttr(resourceName, "agent_list", "hss"),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "auto_terminate_time", "2025-10-10T11:11:00Z"),
				),
			},
			{
				Config: testAccComputeInstance_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "60"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "test222"),
					resource.TestCheckResourceAttr(resourceName, "agent_list", "ces"),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "network.0.source_dest_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_terminate_time", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy", "delete_eip_on_termination", "data_disks", "metadata", "user_data",
				},
			},
		},
	})
}

func TestAccComputeInstance_prePaid(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "delete_eip_on_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_time"),
				),
			},
			{
				Config: testAccComputeInstance_prePaidUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
		},
	})
}

func TestAccComputeInstance_spot(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_spot(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "spot"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy", "delete_eip_on_termination",
					"spot_maximum_price", "spot_duration", "spot_duration_count",
				},
			},
		},
	})
}

func TestAccComputeInstance_powerAction(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_powerAction(rName, "OFF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "ON"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "ON"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "REBOOT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "REBOOT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "FORCE-REBOOT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "FORCE-REBOOT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "FORCE-OFF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "FORCE-OFF"),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"delete_eip_on_termination",
					"power_action",
				},
			},
		},
	})
}

func TestAccComputeInstance_disk_encryption(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_disk_encryption(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_attached.1.kms_key_id",
						"huaweicloud_kms_key.test", "id"),
				),
			},
		},
	})
}

func TestAccComputeInstance_diskIopsThroughput(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_diskIopsThroughput(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "system_disk_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_throughput", "125"),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "data_disks.0.throughput", "200"),
					resource.TestCheckResourceAttr(resourceName, "data_disks.1.type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "data_disks.1.iops", "5000"),
					resource.TestCheckResourceAttr(resourceName, "data_disks.1.throughput", "300"),
				),
			},
		},
	})
}

func TestAccComputeInstance_withEPS(t *testing.T) {
	var instance cloudservers.CloudServer

	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_withEPS(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccComputeInstance_withEPS(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy", "delete_eip_on_termination",
				},
			},
		},
	})
}

func TestAccComputeInstance_changeNetwork(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_changeNetwork(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "network.0.uuid",
						"huaweicloud_vpc_subnet.test1", "id"),
				),
			},
			{
				Config: testAccComputeInstance_changeNetworkUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test update"),
					resource.TestCheckResourceAttrPair(resourceName, "network.0.uuid",
						"huaweicloud_vpc_subnet.test2", "id"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	computeClient, err := cfg.ComputeV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_instance" {
			continue
		}

		server, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeInstanceExists(n string, instance *cloudservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		computeClient, err := cfg.ComputeV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		found, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found")
		}

		*instance = *found

		return nil
	}
}

const testAccCompute_data = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name                = "%s"
  description         = "terraform test"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.huaweicloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agency_name         = "test111"
  agent_list          = "hss"
  auto_terminate_time = "2025-10-10T11:11:00Z"

  user_data = <<EOF
#! /bin/bash
echo user_test > /home/user.txt
EOF

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }

  metadata = {
    foo = "bar"
    key = "value"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name                = "%s-update"
  description         = "terraform test update"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.huaweicloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agency_name         = "test222"
  agent_list          = "ces"
  auto_terminate_time = ""

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = true
  }

  system_disk_type = "SAS"
  system_disk_size = 60

  data_disks {
    type = "SAS"
    size = "10"
  }

  metadata = {
    foo  = "bar2"
    key2 = "value2"
  }

  tags = {
    foo  = "bar2"
    key2 = "value2"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_prePaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  eip_type = "5_bgp"
  bandwidth {
    share_type  = "PER"
    size        = 5
    charge_mode = "bandwidth"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_prePaidUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  eip_type = "5_bgp"
  bandwidth {
    share_type  = "PER"
    size        = 5
    charge_mode = "bandwidth"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_spot(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  charging_mode      = "spot"
  spot_duration      = 2

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_powerAction(rName, powerAction string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  power_action       = "%s"

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, powerAction)
}

func testAccComputeInstance_disk_encryption(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_kms_key" "test" {
  key_alias       = "%s"
  pending_days    = "7"
  key_description = "first test key"
  is_enabled      = true
}

resource "huaweicloud_compute_instance" "test" {
  name                = "%s"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.huaweicloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agent_list          = "hss"

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
    kms_key_id = huaweicloud_kms_key.test.id
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeInstance_withEPS(rName, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name                  = "%s"
  description           = "terraform test"
  image_id              = data.huaweicloud_images_image.test.id
  flavor_id             = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids    = [data.huaweicloud_networking_secgroup.test.id]
  enterprise_project_id = "%s"
  system_disk_type      = "SAS"
  system_disk_size      = 40

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName, epsID)
}

func testAccComputeInstance_diskIopsThroughput(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name                = "%s"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.huaweicloud_networking_secgroup.test.id]
  stop_before_destroy = true

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type       = "GPSSD2"
  system_disk_size       = 50
  system_disk_iops       = 3000
  system_disk_throughput = 125

  data_disks {
    type       = "GPSSD2"
    size       = "50"
    iops       = 4000
    throughput = 200
  }
  data_disks {
    type       = "GPSSD2"
    size       = "50"
    iops       = 5000
    throughput = 300
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_changeNetwork(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc" "test" {
  name = "%[2]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test1" {
  name       = "%[2]s-1"
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  description        = "terraform test"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test1.id
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_changeNetworkUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc" "test" {
  name = "%[2]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test1" {
  name       = "%[2]s-1"
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_subnet" "test2" {
  name       = "%[2]s-2"
  cidr       = "192.168.2.0/24"
  gateway_ip = "192.168.2.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s-update"
  description        = "terraform test update"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test2.id
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }
}
`, testAccCompute_data, rName)
}
