package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
)

func getResourceBackupDownloadPolicyFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dds", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS client: %s", err)
	}

	return dds.GetBackupDownloadPolicyInfo(client)
}

func TestAccBackupDownloadPolicy_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_dds_backup_download_policy.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceBackupDownloadPolicyFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupDownloadPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "action", "open"),
				),
			},
			{
				Config: testAccBackupDownloadPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "action", "close"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBackupDownloadPolicy_basic() string {
	return (`
resource "huaweicloud_dds_backup_download_policy" "test"{
  action = "open"
}
`)
}

func testAccBackupDownloadPolicy_update() string {
	return (`
resource "huaweicloud_dds_backup_download_policy" "test"{
  action = "close"
}
`)
}
