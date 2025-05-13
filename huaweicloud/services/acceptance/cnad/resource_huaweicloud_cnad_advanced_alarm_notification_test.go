package cnad

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cnad"
)

func getAlarmNotificationResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "aad"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CNAD client: %s", err)
	}

	return cnad.GetAlarmNotification(client)
}

func TestAccAlarmNotification_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cnad_advanced_alarm_notification.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAlarmNotificationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlarmNotification_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "topic_urn",
						"huaweicloud_smn_topic.topic_1", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "is_close_attack_source_flag"),
				),
			},
			{
				Config: testAlarmNotification_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "topic_urn",
						"huaweicloud_smn_topic.topic_2", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "is_close_attack_source_flag"),
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

func testAlarmNotification_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name = "%[1]s_1"
}

resource "huaweicloud_smn_topic" "topic_2" {
  name = "%[1]s_2"
}
`, name)
}

func testAlarmNotification_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cnad_advanced_alarm_notification" "test" {
  topic_urn = huaweicloud_smn_topic.topic_1.topic_urn
}
`, testAlarmNotification_base(name))
}

func testAlarmNotification_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cnad_advanced_alarm_notification" "test" {
  topic_urn = huaweicloud_smn_topic.topic_2.topic_urn
}
`, testAlarmNotification_base(name))
}
