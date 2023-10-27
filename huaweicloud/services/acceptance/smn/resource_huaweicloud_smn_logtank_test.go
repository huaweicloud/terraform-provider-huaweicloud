package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/smn/v2/logtank"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/smn"
)

func getResourceSMNLogtankFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SmnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating smn v2 client: %s", err)
	}

	logtankGets, err := logtank.List(client, state.Primary.ID).Extract()
	if err != nil {
		return nil, fmt.Errorf("error list logtanks: %s", err)
	}

	logtankGet := smn.GetLogtankById(logtankGets, state.Primary.Attributes["logtank_id"])
	if logtankGet == nil {
		return nil, fmt.Errorf("the logtank does not exist")
	}
	return logtankGet, nil
}

func TestAccSMNV2Logtank_basic(t *testing.T) {
	var (
		logtankGet          logtank.LogtankGet
		rName               = acceptance.RandomAccResourceNameWithDash()
		logtankResourceName = fmt.Sprintf("huaweicloud_smn_logtank.%s", rName)
	)
	rc := acceptance.InitResourceCheck(
		logtankResourceName,
		&logtankGet,
		getResourceSMNLogtankFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2LogtankConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(logtankResourceName, "topic_urn",
						"huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttrPair(logtankResourceName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(logtankResourceName, "log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
				),
			},
			{
				Config: testAccSMNV2LogtankConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(logtankResourceName, "topic_urn",
						"huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttrPair(logtankResourceName, "log_group_id",
						"huaweicloud_lts_group.test_update", "id"),
					resource.TestCheckResourceAttrPair(logtankResourceName, "log_stream_id",
						"huaweicloud_lts_stream.test_update", "id"),
				),
			},
			{
				ResourceName:      logtankResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSMNV2ImportStateIdFunc(logtankResourceName),
			},
		},
	})
}

func testAccSMNV2ImportStateIdFunc(logtankResourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		resourceLogtank, ok := s.RootModule().Resources[logtankResourceName]
		if !ok {
			return "", fmt.Errorf("smn logtank not found: %s", resourceLogtank)
		}
		topicURN := resourceLogtank.Primary.ID
		logtankID := resourceLogtank.Primary.Attributes["logtank_id"]
		if len(topicURN) == 0 || len(logtankID) == 0 {
			return "", fmt.Errorf("resource not found: %s/%s", topicURN, logtankID)
		}
		return fmt.Sprintf("%s/%s", topicURN, logtankID), nil
	}
}

func testAccSMNV2LogtankConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_smn_logtank" "%[1]s" {
  topic_urn      = huaweicloud_smn_topic.test.topic_urn
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
}
`, rName)
}

func testAccSMNV2LogtankConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_lts_group" "test_update" {
  group_name  = "%[1]s_update"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test_update" {
  group_id    = huaweicloud_lts_group.test_update.id
  stream_name = "%[1]s_update"
}

resource "huaweicloud_smn_logtank" "%[1]s" {
  topic_urn     = huaweicloud_smn_topic.test.topic_urn
  log_group_id  = huaweicloud_lts_group.test_update.id
  log_stream_id = huaweicloud_lts_stream.test_update.id
}
`, rName)
}
