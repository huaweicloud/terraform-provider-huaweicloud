package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGEIPAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_global_eip_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGEIPResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGEIPAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "associate_instance.0.instance_type", "ECS"),
					resource.TestCheckResourceAttrPair(rName, "associate_instance.0.instance_id", "huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "gc_bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(rName, "gc_bandwidth.0.enterprise_project_id"),
					resource.TestCheckResourceAttr(rName, "gc_bandwidth.0.charge_mode", "95"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_reserve_gcb"},
			},
		},
	})
}

func testAccECS(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  key_pair           = huaweicloud_kps_keypair.test.name

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, name)
}

func testAccGEIPAssociate_basic(name string) string {
	return fmt.Sprintf(`
%s

%s

%s

resource "huaweicloud_global_eip_associate" "test" {
  depends_on = [huaweicloud_vpc_internet_gateway.test]

  global_eip_id = huaweicloud_global_eip.test.id

  associate_instance {
    region        = huaweicloud_compute_instance.test.region
    project_id    = "%s"
    instance_type = "ECS"
    instance_id   = huaweicloud_compute_instance.test.id
  }

  gc_bandwidth {
    name        = "%s"
    charge_mode = "95"
    size        = 100
  }
  
  is_reserve_gcb = false
}
`, testAccGEIP_basic(name), testAccIGW_basic(name), testAccECS(name), acceptance.HW_PROJECT_ID, name)
}

func TestAccGEIPAssociate_gcBandwidth(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_global_eip_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGEIPResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGEIPAssociate_gcBandwidth(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "associate_instance.0.instance_type", "ECS"),
					resource.TestCheckResourceAttrPair(rName, "associate_instance.0.instance_id", "huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "gc_bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(rName, "gc_bandwidth.0.enterprise_project_id"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_reserve_gcb"},
			},
		},
	})
}

func testAccGEIPAssociate_gcBandwidth(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

resource "huaweicloud_global_eip_associate" "test" {
  depends_on = [huaweicloud_vpc_internet_gateway.test]

  global_eip_id = huaweicloud_global_eip.test.id

  associate_instance {
    region        = huaweicloud_compute_instance.test.region
    project_id    = "%[5]s"
    instance_type = "ECS"
    instance_id   = huaweicloud_compute_instance.test.id
  }

  gc_bandwidth {
    id = huaweicloud_cc_global_connection_bandwidth.test.id
  }
  
  is_reserve_gcb = false
}
`, testAccGEIP_basic(name), testAccIGW_basic(name), testAccECS(name), testAccGEIPAssociate_gcBandwidth_base(name),
		acceptance.HW_PROJECT_ID)
}

func testAccGEIPAssociate_gcBandwidth_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_global_connection_bandwidth" "test" {
  name        = "%[1]s"
  type        = "Region"  
  bordercross = false
  charge_mode = "bwd"
  size        = 300
  description = "test"
  sla_level   = "Ag"
}
`, name)
}
