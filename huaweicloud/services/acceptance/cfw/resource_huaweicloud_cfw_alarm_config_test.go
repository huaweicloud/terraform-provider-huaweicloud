package cfw

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
)

func getResourceAlarmConfigFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	product := "cfw"

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	alarmTypeStr := state.Primary.Attributes["alarm_type"]
	alarmType, err := strconv.Atoi(alarmTypeStr)
	if err != nil {
		return nil, fmt.Errorf("error converting alarm_type to int: %s", err)
	}

	return cfw.GetAlarmConfig(client, state.Primary.Attributes["fw_instance_id"], alarmType)
}

func TestAccResourceAlarmConfig_basic(t *testing.T) {
	var obj interface{}

	resourceName := "huaweicloud_cfw_alarm_config.test"
	randName := acceptance.RandomAccResourceName()
	baseConfig := testResourceAlarmConfig_base(randName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceAlarmConfigFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceAlarmConfig_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "frequency_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "alarm_time_period", "0"),
					resource.TestCheckResourceAttr(resourceName, "frequency_time", "10"),
					resource.TestCheckResourceAttr(resourceName, "severity", "LOW"),
					resource.TestCheckResourceAttrPair(resourceName, "topic_urn", "huaweicloud_smn_topic.t1", "topic_urn"),
				),
			},
			{
				Config: testResourceAlarmConfig_update(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "frequency_count", "6"),
					resource.TestCheckResourceAttr(resourceName, "alarm_time_period", "1"),
					resource.TestCheckResourceAttr(resourceName, "frequency_time", "5"),
					resource.TestCheckResourceAttr(resourceName, "severity", "MEDIUM,LOW"),
					resource.TestCheckResourceAttrPair(resourceName, "topic_urn", "huaweicloud_smn_topic.t2", "topic_urn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAlarmConfigImportState(resourceName),
			},
		},
	})
}

func testResourceAlarmConfig_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_alarm_config" "test" {
  fw_instance_id    = "%[2]s"
  alarm_type        = 0 
  frequency_count   = 5 
  alarm_time_period = 0 
  frequency_time    = 10
  severity          = "LOW"
  topic_urn         = huaweicloud_smn_topic.t1.topic_urn
}
`, baseConfig, acceptance.HW_CFW_INSTANCE_ID)
}

func testResourceAlarmConfig_update(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_alarm_config" "test" {
  fw_instance_id    = "%[2]s"
  alarm_type        = 0 
  frequency_count   = 6
  alarm_time_period = 1 
  frequency_time    = 5
  severity          = "MEDIUM,LOW"
  topic_urn         = huaweicloud_smn_topic.t2.topic_urn
}
`, baseConfig, acceptance.HW_CFW_INSTANCE_ID)
}

func testResourceAlarmConfig_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "t1" {
  name = "%[1]s-1"
}

resource "huaweicloud_smn_topic" "t2" {
  name = "%[1]s-2"
}
`, name)
}

func testAlarmConfigImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["fw_instance_id"] == "" {
			return "", fmt.Errorf("attribute (fw_instance_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["alarm_type"] == "" {
			return "", fmt.Errorf("attribute (alarm_type) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["fw_instance_id"] + "/" + rs.Primary.Attributes["alarm_type"], nil
	}
}
