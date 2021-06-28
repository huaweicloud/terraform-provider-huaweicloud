package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccApigInstanceV2_basic(t *testing.T) {
	var resourceName = "huaweicloud_apig_instance.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var instance instances.Instance

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
				),
			},
			{
				Config: testAccApigInstance_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
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

func TestAccApigInstanceV2_egress(t *testing.T) {
	var resourceName = "huaweicloud_apig_instance.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var instance instances.Instance

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
				),
			},
			{
				Config: testAccApigInstance_egress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "egress_address"),
				),
			},
			{
				Config: testAccApigInstance_egressUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "egress_address"),
				),
			},
			{
				Config: testAccApigInstance_basic(rName), // Unbind egress nat
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
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

func TestAccApigInstanceV2_ingress(t *testing.T) {
	var resourceName = "huaweicloud_apig_instance.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var instance instances.Instance

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
				),
			},
			{
				Config: testAccApigInstance_ingress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttrSet(resourceName, "eip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ingress_address"),
				),
			},
			{
				Config: testAccApigInstance_ingressUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttrSet(resourceName, "eip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ingress_address"),
				),
			},
			{
				Config: testAccApigInstance_basic(rName), // Unbind ingress eip
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
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

func testAccCheckApigInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_instance" {
			continue
		}
		_, err := instances.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("APIG v2 instance (%s) is still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckApigInstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Resource %s not found", n)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}

		found, err := instances.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*instance = *found
		return nil
	}
}

func testAccApigInstance_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
  cidr       = "192.168.0.0/24"
}
`, rName, rName)
}

func testAccApigInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%s"
  maintain_begin        = "14:00:00"
  description           = "created by acc test"

  available_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApigInstance_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "update" {
  name = "%s-update"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%s-update"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.update.id
  enterprise_project_id = "%s"
  maintain_begin        = "18:00:00"
  description           = "updated by acc test"

  available_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApigInstance_egress(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%s"
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  bandwidth_size        = 3

  available_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApigInstance_egressUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%s"
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  bandwidth_size        = 5

  available_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApigInstance_ingress(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 3
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%s"
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  eip_id                = huaweicloud_vpc_eip.test.id

  available_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApigInstance_ingressUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "update" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s-update"
    size        = 4
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%s"
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  eip_id                = huaweicloud_vpc_eip.update.id

  available_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}
