package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCollectConfig_basic(t *testing.T) {
	var (
		rName = "huaweicloud_secmaster_collect_config.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterDataspaceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "workspace_id"),
					resource.TestCheckResourceAttrSet(rName, "dataspace_id"),
					resource.TestCheckResourceAttrSet(rName, "dataspace_name"),
					resource.TestCheckResourceAttrSet(rName, "region_id"),
					resource.TestCheckResourceAttr(rName, "config.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "config.0.csvc_display"),
					resource.TestCheckResourceAttrSet(rName, "config.0.source_id"),
					resource.TestCheckResourceAttrSet(rName, "config.0.shards"),
					resource.TestCheckResourceAttrSet(rName, "config.0.ttl"),
					resource.TestCheckResourceAttrSet(rName, "config.0.alert"),
				),
			},
		},
	})
}

func testAccCollectConfig_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_collect_config" "test" {
  workspace_id   = "%[1]s"
  dataspace_id   = "%[2]s"
  dataspace_name = "isap-cloudlogs-%[1]s"
  region_id      = "%[3]s"

  config {
    source_id      = 1201
    alert          = true
    ttl            = 7
    shards         = 1
    enable         = 1
    csvc_display   = "数据库审计服务 DBSS"
    source_display = "数据库审计服务告警"
    csvc           = "dbss"
    source_name    = "dbss-alarm"
  }

  lts_config {
    config_name = "test_name"
    description = "test_description"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATASPACE_ID, acceptance.HW_REGION_NAME)
}
