package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ces"
)

func getResourceNotificationMaskFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("ces", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CES client: %s", err)
	}

	return ces.GetNotificationMask(client, state.Primary.Attributes["relation_type"],
		state.Primary.Attributes["relation_id"], state.Primary.ID,
	)
}

func TestAccResourceNotificationMask_resource(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_ces_notification_mask.test"
	name := acceptance.RandomAccResourceNameWithDash()
	baseConfig := testResourceNotificationMask_base(name)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceNotificationMaskFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceNotificationMask_resource(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "mask_name", name),
					resource.TestCheckResourceAttr(rName, "relation_type", "RESOURCE"),
					resource.TestCheckResourceAttr(rName, "mask_type", "FOREVER_TIME"),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
				),
			},
			{
				Config: testResourceNotificationMask_resource_update(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "mask_name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "relation_type", "RESOURCE"),
					resource.TestCheckResourceAttr(rName, "mask_type", "FOREVER_TIME"),
					resource.TestCheckResourceAttr(rName, "resources.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testNotificationMaskImportState(rName),
			},
		},
	})
}

func TestAccResourceNotificationMask_policy(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_ces_notification_mask.test"
	name := acceptance.RandomAccResourceNameWithDash()
	baseConfig := testResourceNotificationMask_base(name)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceNotificationMaskFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCesAlarmPolicies(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceNotificationMask_policy(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "mask_name", name),
					resource.TestCheckResourceAttr(rName, "relation_type", "RESOURCE_POLICY_NOTIFICATION"),
					resource.TestCheckResourceAttr(rName, "mask_type", "FOREVER_TIME"),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
					resource.TestCheckResourceAttr(rName, "relation_ids.#", "2"),
				),
			},
			{
				Config: testResourceNotificationMask_policy_update(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "mask_name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "relation_type", "RESOURCE_POLICY_NOTIFICATION"),
					resource.TestCheckResourceAttr(rName, "resources.#", "1"),
					resource.TestCheckResourceAttr(rName, "relation_ids.#", "1"),
					resource.TestCheckResourceAttr(rName, "mask_type", "START_END_TIME"),
					resource.TestCheckResourceAttr(rName, "start_date", "2025-03-25"),
					resource.TestCheckResourceAttr(rName, "end_date", "2027-03-26"),
					resource.TestCheckResourceAttr(rName, "start_time", "12:00:00"),
					resource.TestCheckResourceAttr(rName, "end_time", "20:00:00"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testNotificationMaskImportState(rName),
			},
		},
	})
}

func TestAccResourceNotificationMask_alarmRule(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_ces_notification_mask.test"
	name := acceptance.RandomAccResourceNameWithDash()
	baseConfig := testResourceNotificationMask_base(name)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceNotificationMaskFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceNotificationMask_alarmRule(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "relation_type", "ALARM_RULE"),
					resource.TestCheckResourceAttr(rName, "mask_type", "FOREVER_TIME"),
					resource.TestCheckResourceAttrPair(rName, "relation_id", "huaweicloud_ces_alarmrule.test", "id"),
				),
			},
			{
				Config: testResourceNotificationMask_alarmRule_update(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "relation_type", "ALARM_RULE"),
					resource.TestCheckResourceAttr(rName, "mask_type", "START_END_TIME"),
					resource.TestCheckResourceAttr(rName, "start_date", "2025-03-25"),
					resource.TestCheckResourceAttr(rName, "end_date", "2027-03-26"),
					resource.TestCheckResourceAttr(rName, "start_time", "12:00:00"),
					resource.TestCheckResourceAttr(rName, "end_time", "20:00:00"),
					resource.TestCheckResourceAttrPair(rName, "relation_id", "huaweicloud_ces_alarmrule.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testNotificationMaskImportStateForAlarmRule(rName),
			},
		},
	})
}

func testResourceNotificationMask_resource(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "RESOURCE"
  mask_name     = "%[2]s"
  
  resources {
    namespace = "SYS.OBS"

    dimensions {
      name  = "bucket_name"
      value = huaweicloud_obs_bucket.bucket[0].bucket
    }
  }

  mask_type = "FOREVER_TIME"

  depends_on = [huaweicloud_ces_alarmrule.test]
}
`, baseConfig, name)
}

func testResourceNotificationMask_resource_update(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "RESOURCE"
  mask_name     = "%[2]s-update"
  
  resources {
    namespace = "SYS.OBS"

    dimensions {
      name  = "bucket_name"
      value = huaweicloud_obs_bucket.bucket[0].bucket
    }
  }

  resources {
    namespace = "SYS.OBS"

    dimensions {
      name  = "bucket_name"
      value = huaweicloud_obs_bucket.bucket[1].bucket
    }
  }

  mask_type = "FOREVER_TIME"

  depends_on = [huaweicloud_ces_alarmrule.test]
}
`, baseConfig, name)
}

func testResourceNotificationMask_policy(baseConfig, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "RESOURCE_POLICY_NOTIFICATION"
  mask_name     = "%[2]s"
  relation_ids  = [
    "%[3]s",
    "%[4]s",
  ]
  
  resources {
    namespace = "SYS.OBS"

    dimensions {
      name  = "bucket_name"
      value = "*"
    }
  }

  mask_type = "FOREVER_TIME"
}
`, baseConfig, name, acceptance.HW_CES_ALARM_POLICY_1, acceptance.HW_CES_ALARM_POLICY_2)
}

func testResourceNotificationMask_policy_update(baseConfig, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "RESOURCE_POLICY_NOTIFICATION"
  mask_name     = "%[2]s-update"
  relation_ids  = ["%[3]s"]
  
  resources {
    namespace = "SYS.OBS"

    dimensions {
      name  = "bucket_name"
      value = "*"
    }
  }

  mask_type  = "START_END_TIME"
  start_date = "2025-03-25"
  end_date   = "2027-03-26"
  start_time = "12:00:00"
  end_time   = "20:00:00"
}
`, baseConfig, name, acceptance.HW_CES_ALARM_POLICY_1)
}

func testResourceNotificationMask_alarmRule(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "ALARM_RULE"
  relation_id   = huaweicloud_ces_alarmrule.test.id
  mask_type     = "FOREVER_TIME"

  depends_on = [huaweicloud_ces_alarmrule.test]
}
`, baseConfig)
}

func testResourceNotificationMask_alarmRule_update(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_notification_mask" "test" {
  relation_type = "ALARM_RULE"
  relation_id   = huaweicloud_ces_alarmrule.test.id

  mask_type  = "START_END_TIME"
  start_date = "2025-03-25"
  end_date   = "2027-03-26"
  start_time = "12:00:00"
  end_time   = "20:00:00"

  depends_on = [huaweicloud_ces_alarmrule.test]
}
`, baseConfig)
}

func testResourceNotificationMask_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "smn-%[1]s"
  display_name = "The display name of smn topic"
}

resource "huaweicloud_obs_bucket" "bucket" {
  count         = 2
  bucket        = "%[1]s-${count.index}"
  storage_class = "STANDARD"
  acl           = "private"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%[1]s"
  alarm_action_enabled = true
  alarm_type           = "ALL_INSTANCE"

  metric {
    namespace = "SYS.OBS"
  }

  resources {
    dimensions {
      name = "bucket_name"
    }
  }

  condition  {
    period              = 1
    filter              = "average"
    comparison_operator = ">"
    value               = 300
    unit                = "B/s"
    count               = 1
    suppress_duration   = 86400
    metric_name         = "request_count_monitor_2XX"
    alarm_level         = 1
  }

  condition  {
    period              = 1
    filter              = "average"
    comparison_operator = ">"
    value               = 50
    unit                = "B/s"
    count               = 1
    suppress_duration   = 86400
    metric_name         = "request_count_4xx"
    alarm_level         = 3
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }

  depends_on = [huaweicloud_obs_bucket.bucket]
}
`, name)
}

func testNotificationMaskImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["relation_type"] == "" {
			return "", fmt.Errorf("Attribute (relation_type) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("Attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["relation_type"] + "/" + rs.Primary.ID, nil
	}
}

func testNotificationMaskImportStateForAlarmRule(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["relation_type"] == "" {
			return "", fmt.Errorf("Attribute (relation_type) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["relation_id"] == "" {
			return "", fmt.Errorf("Attribute (relation_id) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["relation_type"] + "/" + rs.Primary.Attributes["relation_id"], nil
	}
}
