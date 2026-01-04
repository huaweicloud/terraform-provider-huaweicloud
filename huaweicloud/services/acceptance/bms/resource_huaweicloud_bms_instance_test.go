package bms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bms/v1/baremetalservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceInstance(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("server.status", getRespBody, "").(string)
	if status == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccBmsInstance_basic(t *testing.T) {
	var instance baremetalservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_bms_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "user_data"),
					resource.TestCheckResourceAttr(resourceName, "nics.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "nics.0.subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "nics.0.ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "nics.0.mac_address"),
					resource.TestCheckResourceAttrSet(resourceName, "nics.0.port_id"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "test111"),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value"),
				),
			},
			{
				Config: testAccBmsInstance_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "nics.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "test222"),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1"),
				),
			},
			{
				Config: testAccBmsInstance_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "nics.#", "1"),
				),
			},
		},
	})
}

func TestAccBmsInstance_password_basic(t *testing.T) {
	var instance baremetalservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_bms_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_password_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "admin_pass", "Test@123"),
				),
			},
		},
	})
}

func TestAccBmsInstance_userdata(t *testing.T) {
	var instance baremetalservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_bms_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_userdata(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "user_data", "IyEvYmluL2Jhc2ggCmVjaG8gJ3Jvb3Q6VGVzdEAxMjMnIHwgY2hwYXNzd2Q="),
				),
			},
		},
	})
}

func testAccCheckBmsInstanceDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	bmsClient, err := cfg.BmsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating bms client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_bms_instance" {
			continue
		}

		server, err := baremetalservers.Get(bmsClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func TestAccBmsInstance_updateWithEpsId(t *testing.T) {
	var instance baremetalservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_bms_instance.test"
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_withEpsId(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccBmsInstance_withEpsId(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func TestAccBmsInstance_power_action(t *testing.T) {
	var instance interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_bms_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getResourceInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_power_action(rName, "OFF"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				Config: testAccBmsInstance_power_action(rName, "ON"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccBmsInstance_power_action(rName, "REBOOT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccBmsInstance_power_action(rName, "OFF"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
		},
	})
}

func testAccCheckBmsInstanceExists(n string, instance *baremetalservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		bmsClient, err := cfg.BmsV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating bms client: %s", err)
		}

		found, err := baremetalservers.Get(bmsClient, rs.Primary.ID).Extract()
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

func testAccBmsInstance_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_bms_flavors" "test" {
  cpu_arch          = "x86_64"
  vcpus             = 56
  memory            = 192
  availability_zone = try(element(data.huaweicloud_availability_zones.test.names, 0), "")
}

data "huaweicloud_images_images" "test" {
  name_regex = "x86"
  os         = "CentOS"
  image_type = "Ironic"
}

locals {
  x86_images = [for v in data.huaweicloud_images_images.test.images: v.id if v.container_format == "bare"]
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%s"
}

resource "huaweicloud_vpc_subnet" "test1" {
  name       = "%[2]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.5.0/24"
  gateway_ip = "192.168.5.99"
}`, common.TestBaseNetwork(rName), rName)
}

// Both `key_pair` and `user_data` are specified, `user_data` only injects user data with plain text.
func testAccBmsInstance_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = try(local.x86_images[0], "")

  user_data = <<EOF
#!/bin/bash 
sudo mkdir /example
EOF

  name                  = "%[2]s"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  system_disk_type = "GPSSD"
  system_disk_size = 150

  tags = {
    foo = "bar"
    key = "value"
  }

  agency_name = "test111"

  metadata = {
    foo = "bar"
    key = "value"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "false"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccBmsInstance_password_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  admin_pass        = "Test@123"
  image_id          = try(local.x86_images[0], "")

  name                  = "%[2]s"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  system_disk_type = "GPSSD"
  system_disk_size = 150

  tags = {
    foo = "bar"
    key = "value"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// The `user_data` field is specified for a Linux BMS that is created using an image with Cloud-Init installed,
// the `admin_pass` field becomes invalid and `user_data` will be injected into the BMS as a password.
// The `user_data` field is the result of base64 encoding of the following command:
// #!/bin/bash
// echo 'root:Test@123' | chpasswd
func testAccBmsInstance_userdata(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  user_data         = "IyEvYmluL2Jhc2ggCmVjaG8gJ3Jvb3Q6VGVzdEAxMjMnIHwgY2hwYXNzd2Q="
  image_id          = try(local.x86_images[0], "")

  name                  = "%[2]s"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  system_disk_type = "GPSSD"
  system_disk_size = 150

  tags = {
    foo = "bar"
    key = "value"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccBmsInstance_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = try(local.x86_images[0], "")

  name                  = "%[2]s_update"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  nics {
    subnet_id = huaweicloud_vpc_subnet.test1.id
  }

  system_disk_type = "GPSSD"
  system_disk_size = 150

  tags = {
    tag1 = "value1"
    tag2 = "value2"
  }

  agency_name = "test222"

  metadata = {
    foo1 = "bar1"
    key1 = "value1"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccBmsInstance_update2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = try(local.x86_images[0], "")

  name                  = "%[2]s_update"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  tags = {
    tag1 = "value1"
    tag2 = "value2"
  }

  system_disk_type = "GPSSD"
  system_disk_size = 150

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccBmsInstance_withEpsId(rName, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = "physical.s4.3xlarge"
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = try(local.x86_images[0], "")

  name                  = "%[2]s"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  tags = {
    tag1 = "value1"
    tag2 = "value2"
  }

  system_disk_type = "GPSSD"
  system_disk_size = 150

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, epsId)
}

func testAccBmsInstance_power_action(rName, powerAction string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = try(local.x86_images[0], "")
  name              = "%[2]s"
  user_id           = "%[3]s"
  power_action      = "%[4]s"
  system_disk_type  = "GPSSD"
  system_disk_size  = 150

  user_data = <<EOF
#!/bin/bash 
sudo mkdir /example
EOF

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  metadata = {
    foo1 = "bar1"
    key1 = "value1"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "false"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, powerAction)
}
