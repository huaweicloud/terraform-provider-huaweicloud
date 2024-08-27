package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sfs/v2/shares"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getSfsAccessRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SfsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS client: %s", err)
	}

	resourceID := state.Primary.ID
	sfsID := state.Primary.Attributes["sfs_id"]
	rules, err := shares.ListAccessRights(client, sfsID).ExtractAccessRights()
	if err != nil {
		return nil, err
	}

	for _, item := range rules {
		if item.ID == resourceID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("the sfs access rule %s does not exist", resourceID)
}

func TestAccSFSAccessRuleV2_basic(t *testing.T) {
	var rule shares.AccessRight
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_access_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&rule,
		getSfsAccessRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configAccSFSAccessRuleV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				Config: configAccSFSAccessRuleV2_ipAuth(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
		},
	})
}

func configAccSFSAccessRuleV2_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "test" {
  name                  = "vpc-default"
  enterprise_project_id = "0"
}

resource "huaweicloud_sfs_file_system" "test" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "sfs file system created by terraform testacc"
}

resource "huaweicloud_sfs_access_rule" "test" {
  sfs_id    = huaweicloud_sfs_file_system.test.id
  access_to = data.huaweicloud_vpc.test.id
}`, rName)
}

func configAccSFSAccessRuleV2_ipAuth(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "test" {
  name                  = "vpc-default"
  enterprise_project_id = "0"
}

resource "huaweicloud_sfs_file_system" "test" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "sfs file system created by terraform testacc"
}

resource "huaweicloud_sfs_access_rule" "test" {
  sfs_id    = huaweicloud_sfs_file_system.test.id
  access_to = join("#", [data.huaweicloud_vpc.test.id, "192.168.10.0/24", "0", "no_all_squash,no_root_squash"])
}`, rName)
}
