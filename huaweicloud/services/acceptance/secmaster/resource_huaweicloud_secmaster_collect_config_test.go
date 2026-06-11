package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceCollectConfigFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := secmaster.GetCollectConfigInfo(client, state.Primary.Attributes["region_id"], state.Primary.ID)
	if err != nil {
		return nil, err
	}

	// If the config is disabled (enable=inactive), treat it as deleted
	enableStatus := utils.PathSearch("config.enable", respBody, "").(string)
	if enableStatus == "inactive" {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccResourceCollectConfig_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_secmaster_collect_config.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceCollectConfigFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterDataspaceId(t)
			acceptance.TestAccPreCheckSecMasterDataspaceName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCollectConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "workspace_id"),
					resource.TestCheckResourceAttrSet(rName, "dataspace_id"),
					resource.TestCheckResourceAttrSet(rName, "dataspace_name"),
					resource.TestCheckResourceAttrSet(rName, "region_id"),
					resource.TestCheckResourceAttr(rName, "config.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "config.0.csvc_display"),
					resource.TestCheckResourceAttrSet(rName, "config.0.source_id"),
					resource.TestCheckResourceAttrSet(rName, "config.0.shards"),
					resource.TestCheckResourceAttrSet(rName, "config.0.ttl"),
				),
			},
			{
				Config: testAccCollectConfig_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "config.#", "1"),
					resource.TestCheckResourceAttr(rName, "config.0.alert", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"lts_config",
				},
				ImportStateIdFunc: testAccCollectConfigImportStateFunc(rName),
			},
		},
	})
}

func testAccCollectConfig_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_collect_config" "test" {
  workspace_id   = "%[1]s"
  dataspace_id   = "%[2]s"
  dataspace_name = "%[3]s"
  region_id      = "%[4]s"

  config {
    source_id      = "%[5]s"
    alert          = true
    ttl            = 7
    shards         = 1
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
`, acceptance.HW_SECMASTER_WORKSPACE_ID,
		acceptance.HW_SECMASTER_DATASPACE_ID,
		acceptance.HW_SECMASTER_DATASPACE_NAME,
		acceptance.HW_REGION_NAME,
		acceptance.HW_SECMASTER_SOURCE_ID)
}

func testAccCollectConfig_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_collect_config" "test" {
  workspace_id   = "%[1]s"
  dataspace_id   = "%[2]s"
  dataspace_name = "%[3]s"
  region_id      = "%[4]s"

  config {
    source_id      = "%[5]s"
    alert          = false
    ttl            = 7
    shards         = 1
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
`, acceptance.HW_SECMASTER_WORKSPACE_ID,
		acceptance.HW_SECMASTER_DATASPACE_ID,
		acceptance.HW_SECMASTER_DATASPACE_NAME,
		acceptance.HW_REGION_NAME,
		acceptance.HW_SECMASTER_SOURCE_ID)
}

func testAccCollectConfigImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var regionId, sourceId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		regionId = rs.Primary.Attributes["region_id"]
		sourceId = rs.Primary.ID

		if regionId == "" || sourceId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<region_id>/<id>', but got '%s/%s'",
				regionId, sourceId)
		}

		return fmt.Sprintf("%s/%s", regionId, sourceId), nil
	}
}
