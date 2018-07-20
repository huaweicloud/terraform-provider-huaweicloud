package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
)

func TestAccSFSFileSystemV2_basic(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSFSFileSystemV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists("huaweicloud_sfs_file_system_v2.sfs_1", &share),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "name", "sfs-test1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "share_proto", "NFS"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "status", "available"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "size", "1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "access_level", "rw"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "access_to", OS_VPC_ID),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "access_type", "cert"),
				),
			},
			resource.TestStep{
				Config: testAccSFSFileSystemV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists("huaweicloud_sfs_file_system_v2.sfs_1", &share),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "name", "sfs-test2"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "share_proto", "NFS"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "status", "available"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "size", "2"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "access_level", "rw"),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "access_to", OS_VPC_ID),
					resource.TestCheckResourceAttr(
						"huaweicloud_sfs_file_system_v2.sfs_1", "access_type", "cert"),
				),
			},
		},
	})
}

func TestAccSFSFileSystemV2_timeout(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSFSFileSystemV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists("huaweicloud_sfs_file_system_v2.sfs_1", &share),
				),
			},
		},
	})
}

func testAccCheckSFSFileSystemV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sfsClient, err := config.sfsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_sfs_file_system_v2" {
			continue
		}

		_, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Share File still exists")
		}
	}

	return nil
}

func testAccCheckSFSFileSystemV2Exists(n string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.sfsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud sfs client: %s", err)
		}

		found, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("share file not found")
		}

		*share = *found

		return nil
	}
}

var testAccSFSFileSystemV2_basic = fmt.Sprintf(`
resource "huaweicloud_sfs_file_system_v2" "sfs_1" {
	share_proto = "NFS"
	size=1
	name="sfs-test1"
  	availability_zone="%s"
	access_to="%s"
  	access_type="cert"
  	access_level="rw"
	description="sfs_c2c_test-file"
}
`, OS_AVAILABILITY_ZONE, OS_VPC_ID)

var testAccSFSFileSystemV2_update = fmt.Sprintf(`
resource "huaweicloud_sfs_file_system_v2" "sfs_1" {
	share_proto = "NFS"
	size=2
	name="sfs-test2"
  	availability_zone="%s"
	access_to="%s"
  	access_type="cert"
  	access_level="rw"
	description="sfs_c2c_test-file"
}
`, OS_AVAILABILITY_ZONE, OS_VPC_ID)

var testAccSFSFileSystemV2_timeout = fmt.Sprintf(`
resource "huaweicloud_sfs_file_system_v2" "sfs_1" {
	share_proto = "NFS"
	size=1
	name="sfs-test1"
  	availability_zone="%s"
	access_to="%s"
  	access_type="cert"
  	access_level="rw"
	description="sfs_c2c_test-file"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}`, OS_AVAILABILITY_ZONE, OS_VPC_ID)
