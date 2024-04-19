package lts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
)

func getLogConvergeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("lts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	return lts.GetLogConvergeConfigsById(client, state.Primary.Attributes["member_account_id"])
}

func TestAccLogConverge_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_lts_log_converge.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getLogConvergeResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSLogConvergeBaseConfig(t)
			acceptance.TestAccPreCheckLTSLogConvergeGroupConfig(t)
			acceptance.TestAccPreCheckLTSLogConvergeStreamConfig(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogConverge_basic(90),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "organization_id", acceptance.HW_LTS_LOG_CONVERGE_ORGANIZATION_ID),
					resource.TestCheckResourceAttr(rName, "management_account_id", acceptance.HW_LTS_LOG_CONVERGE_MANAGEMENT_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "member_account_id", acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.source_log_group_id",
						acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_GROUP_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.target_log_group_name",
						acceptance.HW_LTS_LOG_CONVERGE_TARGET_LOG_GROUP_NAME),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.target_log_group_id",
						acceptance.HW_LTS_LOG_CONVERGE_TARGET_LOG_GROUP_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.0.source_log_stream_id",
						acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_STREAM_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.0.target_log_stream_name",
						acceptance.HW_LTS_LOG_CONVERGE_TARGET_LOG_STREAM_NAME),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.0.target_log_stream_id",
						acceptance.HW_LTS_LOG_CONVERGE_TARGET_LOG_STREAM_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.0.target_log_stream_ttl",
						"90"),
				),
			},
			{
				Config: testLogConverge_basic(150),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.0.target_log_stream_ttl",
						"150"),
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

func testLogConverge_basic(ttl int) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_log_converge" "test" {
  organization_id       = "%[1]s"
  management_account_id = "%[2]s"
  member_account_id     = "%[3]s"

  log_mapping_config {
    source_log_group_id   = "%[4]s"
    target_log_group_name = "%[5]s"
    target_log_group_id   = "%[6]s"

    log_stream_config {
      source_log_stream_id   = "%[7]s"
      target_log_stream_name = "%[8]s"
      target_log_stream_id   = "%[9]s"
      target_log_stream_ttl  = %[10]d
    }
  }
}
`, acceptance.HW_LTS_LOG_CONVERGE_ORGANIZATION_ID,
		acceptance.HW_LTS_LOG_CONVERGE_MANAGEMENT_ACCOUNT_ID,
		acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID,
		acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_GROUP_ID,
		acceptance.HW_LTS_LOG_CONVERGE_TARGET_LOG_GROUP_NAME,
		acceptance.HW_LTS_LOG_CONVERGE_TARGET_LOG_GROUP_ID,
		acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_STREAM_ID,
		acceptance.HW_LTS_LOG_CONVERGE_TARGET_LOG_STREAM_NAME,
		acceptance.HW_LTS_LOG_CONVERGE_TARGET_LOG_STREAM_ID,
		ttl,
	)
}
