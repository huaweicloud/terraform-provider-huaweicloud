package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/sfs/v2/shares"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccSFSFileSystemV2_basic(t *testing.T) {
	var share shares.Share
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_sfs_file_system.sfs_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "cert"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccSFSFileSystemV2_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
		},
	})
}

func TestAccSFSFileSystemV2_withEpsId(t *testing.T) {
	var share shares.Share
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_sfs_file_system.sfs_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_epsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccSFSFileSystemV2_withoutRule(t *testing.T) {
	var share shares.Share
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_sfs_file_system.sfs_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_withoutRule(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "unavailable"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
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

func testAccCheckSFSFileSystemV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	sfsClient, err := config.SfsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_sfs_file_system" {
			continue
		}

		_, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Share File still exists")
		}
	}

	return nil
}

func testAccCheckSFSFileSystemV2Exists(n string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		sfsClient, err := config.SfsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud sfs client: %s", err)
		}

		found, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("share file not found")
		}

		*share = *found

		return nil
	}
}

func testAccSFSFileSystemV2_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_file_system" "sfs_1" {
  share_proto  = "NFS"
  size         = 10
  name         = "%s"
  description  = "sfs_c2c_test-file"
  access_to    = huaweicloud_vpc.test.id
  access_level = "rw"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, rName)
}

func testAccSFSFileSystemV2_epsId(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_file_system" "sfs_1" {
  share_proto  = "NFS"
  size         = 10
  name         = "%s"
  description  = "sfs_c2c_test-file"
  access_to    = huaweicloud_vpc.test.id
  access_type  = "cert"
  access_level = "rw"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  enterprise_project_id = "%s"
}
`, rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccSFSFileSystemV2_update(rName, updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_file_system" "sfs_1" {
  share_proto = "NFS"
  size        = 20
  name        = "%s"
  description = "sfs_c2c_test-file"
  access_to   = huaweicloud_vpc.test.id
  access_type = "cert"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]

  tags = {
    foo   = "bar"
    owner = "terraform_update"
  }
}
`, rName, updateName)
}

func testAccSFSFileSystemV2_withoutRule(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_sfs_file_system" "sfs_1" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "sfs_c2c_test-file"
}
`, rName)
}
