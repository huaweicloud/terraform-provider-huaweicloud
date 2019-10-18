package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCSBSBackupV1_importBasic(t *testing.T) {
	resourceName := "huaweicloud_csbs_backup_v1.csbs"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}
