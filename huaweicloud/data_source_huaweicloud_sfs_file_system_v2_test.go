package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSFSFileSystemV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2DataSourceID("data.huaweicloud_sfs_file_system_v2.shares"),
					resource.TestCheckResourceAttr("data.huaweicloud_sfs_file_system_v2.shares", "name", rName),
					resource.TestCheckResourceAttr("data.huaweicloud_sfs_file_system_v2.shares", "status", "available"),
					resource.TestCheckResourceAttr("data.huaweicloud_sfs_file_system_v2.shares", "size", "1"),
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
  enterprise_project_id = "0"
}

data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_file_system_v2" "sfs_1" {
	share_proto = "NFS"
	size=1
	name="%s"
	availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
	access_to = data.huaweicloud_vpc.vpc_default.id
  	access_type="cert"
  	access_level="rw"
	description="sfs_c2c_test-file"
}
data "huaweicloud_sfs_file_system_v2" "shares" {
  id = huaweicloud_sfs_file_system_v2.sfs_1.id
}
`, rName)
}
