package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSFSFileSystemV2_importBasic(t *testing.T) {
	resourceName := "huaweicloud_sfs_file_system_v2.sfs_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSFSFileSystemV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
