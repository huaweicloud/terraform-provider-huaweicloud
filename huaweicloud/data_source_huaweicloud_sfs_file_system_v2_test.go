package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSFSFileSystemV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_sfs_file_system.share_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
				),
			},
		},
	})
}

func testAccCheckSFSFileSystemV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find share file data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("share file data source ID not set ")
		}

		return nil
	}
}

func testAccSFSFileSystemV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "vpc_default" {
  name = "vpc-default"
}

data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_file_system" "sfs_1" {
  share_proto  = "NFS"
  size         = 10
  name         = "%s"
  description  = "sfs_c2c_test-file"
  access_to    = data.huaweicloud_vpc.vpc_default.id
  access_level = "rw"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
}

data "huaweicloud_sfs_file_system" "share_1" {
  id = huaweicloud_sfs_file_system.sfs_1.id
}
`, rName)
}
