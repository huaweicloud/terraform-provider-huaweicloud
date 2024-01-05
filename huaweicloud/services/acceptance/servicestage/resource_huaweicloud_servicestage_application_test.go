package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/applications"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAppResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ServiceStageV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage v2 client: %s", err)
	}
	return applications.Get(c, state.Primary.ID)
}

func TestAccApplication_basic(t *testing.T) {
	var (
		app          applications.Application
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_servicestage_application.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&app,
		getAppResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "environment.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "environment.0.id",
						"huaweicloud_servicestage_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "environment.0.variable.#", "3"),
				),
			},
			{
				Config: testAccApplication_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform test"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "environment.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "environment.0.id",
						"huaweicloud_servicestage_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "environment.0.variable.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "environment.0.variable.0.name", "owner"),
					resource.TestCheckResourceAttr(resourceName, "environment.0.variable.0.value", "terraform"),
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

func testAccApplication_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

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

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"

  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}

resource "huaweicloud_networking_secgroup" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  key_pair           = huaweicloud_kps_keypair.test.name
  security_group_ids = [huaweicloud_networking_secgroup.test.id]

  enterprise_project_id = "%[2]s"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_servicestage_environment" "test" {
  name   = "%[1]s"
  vpc_id = huaweicloud_vpc.test.id

  basic_resources {
    type = "ecs"
	id   = huaweicloud_compute_instance.test.id
  }
}`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApplication_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_servicestage_application" "test" {
  name        = "%s"
  description = "Created by terraform test"

  enterprise_project_id = "%s"

  environment {
    id = huaweicloud_servicestage_environment.test.id

    variable {
      name  = "_underscore-.001"
      value = "special characters: ~!@#$%%&^*()-_=+{[]}\\|;'<.?/,"
    }
    variable {
      name  = "-hyphen_.002"
      value = "abcdefghijklmnopqrstuvwxyz"
    }
    variable {
      name  = "letter-_.003"
      value = "1234567890"
    }
  }
}
`, testAccApplication_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApplication_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_servicestage_application" "test" {
  name        = "%s-update"
  description = "Updated by terraform test"

  enterprise_project_id = "%s"

  environment {
    id = huaweicloud_servicestage_environment.test.id

    variable {
      name  = "owner"
      value = "terraform"
    }
  }
}
`, testAccApplication_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
