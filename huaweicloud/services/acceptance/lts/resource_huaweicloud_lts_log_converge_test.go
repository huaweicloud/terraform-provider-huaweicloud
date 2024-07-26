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
		name  = acceptance.RandomAccResourceName()
		rc    = acceptance.InitResourceCheck(rName, &obj, getLogConvergeResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSLogConvergeBaseConfig(t)
			acceptance.TestAccPreCheckLTSLogConvergeMappingConfig(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogConverge_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "organization_id", acceptance.HW_LTS_LOG_CONVERGE_ORGANIZATION_ID),
					resource.TestCheckResourceAttr(rName, "management_account_id", acceptance.HW_LTS_LOG_CONVERGE_MANAGEMENT_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "member_account_id", acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.source_log_group_id",
						acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_GROUP_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.target_log_group_name", name),
					resource.TestCheckResourceAttrSet(rName, "log_mapping_config.0.target_log_group_id"),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.#", "0"),
				),
			},
			{
				Config: testLogConverge_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "organization_id", acceptance.HW_LTS_LOG_CONVERGE_ORGANIZATION_ID),
					resource.TestCheckResourceAttr(rName, "management_account_id", acceptance.HW_LTS_LOG_CONVERGE_MANAGEMENT_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "member_account_id", acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.source_log_group_id",
						acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_GROUP_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.target_log_group_name", name),
					resource.TestCheckResourceAttrSet(rName, "log_mapping_config.0.target_log_group_id"),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.#", "1"),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.0.source_log_stream_id",
						acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_STREAM_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.0.target_log_stream_name", name),
					resource.TestCheckResourceAttrSet(rName, "log_mapping_config.0.log_stream_config.0.target_log_stream_id"),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.0.log_stream_config.0.target_log_stream_ttl", "90"),
				),
			},
			{
				Config: testLogConverge_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "organization_id", acceptance.HW_LTS_LOG_CONVERGE_ORGANIZATION_ID),
					resource.TestCheckResourceAttr(rName, "management_account_id", acceptance.HW_LTS_LOG_CONVERGE_MANAGEMENT_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "member_account_id", acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "log_mapping_config.#", "2"),
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

func testLogConverge_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_log_converge" "test" {
  organization_id       = "%[1]s"
  management_account_id = "%[2]s"
  member_account_id     = "%[3]s"

  log_mapping_config {
    source_log_group_id = "%[4]s"
    # Do not use the existing log group, and automatically create it through the LTS service
    target_log_group_name = "%[5]s"
  }
}
`, acceptance.HW_LTS_LOG_CONVERGE_ORGANIZATION_ID,
		acceptance.HW_LTS_LOG_CONVERGE_MANAGEMENT_ACCOUNT_ID,
		acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID,
		acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_GROUP_ID,
		name,
	)
}

func testLogConverge_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_log_converge" "test" {
  organization_id       = "%[1]s"
  management_account_id = "%[2]s"
  member_account_id     = "%[3]s"

  log_mapping_config {
    source_log_group_id = "%[4]s"
    # Do not use the existing log group, and automatically create it through the LTS service
    target_log_group_name = "%[5]s"

    log_stream_config {
      source_log_stream_id = "%[6]s"
      # Do not use the existing log stream, and automatically create it through the LTS service
      target_log_stream_name = "%[5]s"
      target_log_stream_ttl  = 90
    }
  }
}
`, acceptance.HW_LTS_LOG_CONVERGE_ORGANIZATION_ID,
		acceptance.HW_LTS_LOG_CONVERGE_MANAGEMENT_ACCOUNT_ID,
		acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID,
		acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_GROUP_ID,
		name,
		acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_STREAM_ID,
	)
}

func testLogConverge_basic_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[5]s_manual"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[5]s_manual"
  ttl_in_days = 60
}

resource "huaweicloud_lts_log_converge" "test" {
  organization_id       = "%[1]s"
  management_account_id = "%[2]s"
  member_account_id     = "%[3]s"

  log_mapping_config {
    source_log_group_id = "%[4]s"
    # Do not use the existing log group, and automatically create it through the LTS service
    target_log_group_name = "%[5]s"

    log_stream_config {
      source_log_stream_id = "%[6]s"
      # Do not use the existing log stream, and automatically create it through the LTS service
      target_log_stream_name = "%[5]s"
      target_log_stream_ttl  = 150
    }
  }
  log_mapping_config {
    source_log_group_id = "%[4]s"
    # Use the existing log group
    target_log_group_name = huaweicloud_lts_group.test.group_name
    target_log_group_id   = huaweicloud_lts_group.test.id

    log_stream_config {
      source_log_stream_id   = "%[6]s"
      target_log_stream_name = huaweicloud_lts_stream.test.stream_name
      # Use the existing log stream
      target_log_stream_id  = huaweicloud_lts_stream.test.id
      target_log_stream_ttl = 90
    }
  }
}
`, acceptance.HW_LTS_LOG_CONVERGE_ORGANIZATION_ID,
		acceptance.HW_LTS_LOG_CONVERGE_MANAGEMENT_ACCOUNT_ID,
		acceptance.HW_LTS_LOG_CONVERGE_MEMBER_ACCOUNT_ID,
		acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_GROUP_ID,
		name,
		acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_STREAM_ID,
	)
}
