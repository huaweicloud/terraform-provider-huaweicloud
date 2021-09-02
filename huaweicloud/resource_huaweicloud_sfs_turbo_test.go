package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/sfs_turbo/v1/shares"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccSFSTurbo_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	turboName := fmt.Sprintf("sfs-turbo-acc-%s", randSuffix)
	resourceName := "huaweicloud_sfs_turbo.sfs-turbo1"
	var turbo shares.Turbo

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSTurboDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboExists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", turboName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "enhanced", "false"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSFSTurbo_update(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboExists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "size", "600"),
					resource.TestCheckResourceAttr(resourceName, "status", "221"),
				),
			},
		},
	})
}

func TestAccSFSTurbo_crypt(t *testing.T) {
	randSuffix := acctest.RandString(5)
	turboName := fmt.Sprintf("sfs-turbo-acc-%s", randSuffix)
	resourceName := "huaweicloud_sfs_turbo.sfs-turbo1"
	var turbo shares.Turbo

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSTurboDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_crypt(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboExists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", turboName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "enhanced", "false"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
					resource.TestCheckResourceAttrSet(resourceName, "crypt_key_id"),
				),
			},
		},
	})
}

func TestAccSFSTurbo_withEpsId(t *testing.T) {
	randSuffix := acctest.RandString(5)
	turboName := fmt.Sprintf("sfs-turbo-acc-%s", randSuffix)
	resourceName := "huaweicloud_sfs_turbo.sfs-turbo1"
	var turbo shares.Turbo

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSTurboDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_withEpsId(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboExists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", turboName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckSFSTurboDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	sfsClient, err := config.SfsV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud sfs turbo client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_sfs_turbo" {
			continue
		}

		_, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("SFS Turbo still exists")
		}
	}

	return nil
}

func testAccCheckSFSTurboExists(n string, share *shares.Turbo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		sfsClient, err := config.SfsV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud sfs turbo client: %s", err)
		}

		found, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("sfs turbo not found")
		}

		*share = *found
		return nil
	}
}

func testAccNetworkPreConditions(suffix string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_v1" "test" {
  name = "tf-acc-vpc-%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet_v1" "test" {
  name       = "tf-acc-subnet-%s"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc_v1.test.id
}

resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name        = "tf-acc-sg-%s"
  description = "terraform security group for sfs turbo acceptance test"
}
`, suffix, suffix, suffix)
}

func testAccSFSTurbo_basic(suffix string) string {
	return fmt.Sprintf(`
%s
data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_turbo" "sfs-turbo1" {
  name        = "sfs-turbo-acc-%s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = huaweicloud_vpc_v1.test.id
  subnet_id   = huaweicloud_vpc_subnet_v1.test.id
  security_group_id = huaweicloud_networking_secgroup_v2.secgroup.id
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
}
`, testAccNetworkPreConditions(suffix), suffix)
}

func testAccSFSTurbo_update(suffix string) string {
	return fmt.Sprintf(`
%s
data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_turbo" "sfs-turbo1" {
  name        = "sfs-turbo-acc-%s"
  size        = 600
  share_proto = "NFS"
  vpc_id      = huaweicloud_vpc_v1.test.id
  subnet_id   = huaweicloud_vpc_subnet_v1.test.id
  security_group_id = huaweicloud_networking_secgroup_v2.secgroup.id
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
}
`, testAccNetworkPreConditions(suffix), suffix)
}

func testAccSFSTurbo_crypt(suffix string) string {
	return fmt.Sprintf(`
%s
data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_kms_key" "key_1" {
  key_alias = "kms-acc-%s"
  pending_days = "7"
}

resource "huaweicloud_sfs_turbo" "sfs-turbo1" {
  name        = "sfs-turbo-acc-%s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = huaweicloud_vpc_v1.test.id
  subnet_id   = huaweicloud_vpc_subnet_v1.test.id
  security_group_id = huaweicloud_networking_secgroup_v2.secgroup.id
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  crypt_key_id      = huaweicloud_kms_key.key_1.id
}
`, testAccNetworkPreConditions(suffix), suffix, suffix)
}

func testAccSFSTurbo_withEpsId(suffix string) string {
	return fmt.Sprintf(`
%s
data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_turbo" "sfs-turbo1" {
  name                   = "sfs-turbo-acc-%s"
  size                   = 500
  share_proto            = "NFS"
  vpc_id                 = huaweicloud_vpc_v1.test.id
  subnet_id              = huaweicloud_vpc_subnet_v1.test.id
  security_group_id      = huaweicloud_networking_secgroup_v2.secgroup.id
  availability_zone      = data.huaweicloud_availability_zones.myaz.names[0]
  enterprise_project_id  = "%s"
}
`, testAccNetworkPreConditions(suffix), suffix, HW_ENTERPRISE_PROJECT_ID_TEST)
}
