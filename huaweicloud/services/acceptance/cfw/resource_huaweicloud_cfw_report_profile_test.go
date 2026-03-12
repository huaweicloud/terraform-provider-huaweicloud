package cfw

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
)

func getReportProfileResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "cfw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	return cfw.GetReportProfileDetail(client, state.Primary.ID, state.Primary.Attributes["fw_instance_id"])
}

func TestAccReportProfile_daily(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_cfw_report_profile.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getReportProfileResourceFunc,
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
				Config: testReportProfile_daily(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "daily"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "send_period", "3"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
					resource.TestCheckResourceAttr(rName, "subscription_type", "0"),
					resource.TestCheckResourceAttrPair(rName, "topic_urn", "huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "topic_name"),
				),
			},
			{
				Config: testReportProfile_daily_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "daily"),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "send_period", "4"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "subscription_type", "1"),
					resource.TestCheckResourceAttrPair(rName, "topic_urn", "huaweicloud_smn_topic.test_another", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "topic_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testReportProfileImportState(rName),
				ImportStateVerifyIgnore: []string{
					"description",
				},
			},
		},
	})
}

func testReportProfile_daily(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "%[1]s"
}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = "%[2]s"
  category          = "daily"
  name              = "%[1]s"
  topic_urn         = huaweicloud_smn_topic.test.topic_urn
  send_period       = "3"
  subscription_type = "0"
  status            = "0"
  description       = "test description"
}
`, name, acceptance.HW_CFW_INSTANCE_ID)
}

func testReportProfile_daily_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test_another" {
  name         = "%[1]s_another"
  display_name = "%[1]s_another"
}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = "%[2]s"
  category          = "daily"
  name              = "%[1]s_update"
  topic_urn         = huaweicloud_smn_topic.test_another.topic_urn
  send_period       = "4"
  subscription_type = "1"
  status            = "1"
  description       = "test description update"
}
`, name, acceptance.HW_CFW_INSTANCE_ID)
}

func TestAccReportProfile_weekly(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_cfw_report_profile.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getReportProfileResourceFunc,
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
				Config: testReportProfile_weekly(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "weekly"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "send_period", "1"),
					resource.TestCheckResourceAttr(rName, "send_week_day", "1"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "subscription_type", "0"),
					resource.TestCheckResourceAttrPair(rName, "topic_urn", "huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "topic_name"),
				),
			},
			{
				Config: testReportProfile_weekly_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "weekly"),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "send_period", "3"),
					resource.TestCheckResourceAttr(rName, "send_week_day", "2"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
					resource.TestCheckResourceAttr(rName, "subscription_type", "1"),
					resource.TestCheckResourceAttrPair(rName, "topic_urn", "huaweicloud_smn_topic.test_another", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "topic_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testReportProfileImportState(rName),
				ImportStateVerifyIgnore: []string{
					"description",
				},
			},
		},
	})
}

func testReportProfile_weekly(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "%[1]s"
}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = "%[2]s"
  category          = "weekly"
  name              = "%[1]s"
  topic_urn         = huaweicloud_smn_topic.test.topic_urn
  send_period       = "1"
  send_week_day     = "1"
  subscription_type = "0"
  status            = "1"
  description       = "test description"
}
`, name, acceptance.HW_CFW_INSTANCE_ID)
}

func testReportProfile_weekly_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test_another" {
  name         = "%[1]s_another"
  display_name = "%[1]s_another"
}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = "%[2]s"
  category          = "weekly"
  name              = "%[1]s_update"
  topic_urn         = huaweicloud_smn_topic.test_another.topic_urn
  send_period       = "3"
  send_week_day     = "2"
  subscription_type = "1"
  status            = "0"
  description       = "test description update"
}
`, name, acceptance.HW_CFW_INSTANCE_ID)
}

func TestAccReportProfile_custom(t *testing.T) {
	var (
		obj              interface{}
		name             = acceptance.RandomAccResourceName()
		rName            = "huaweicloud_cfw_report_profile.test"
		twoWeeksAgoStart = getDaysAgoStart(14)
		oneWeekAgoStart  = getDaysAgoStart(7)
		oneWeekAgoEnd    = getDaysAgoEnd(7)
		threeDaysAgoEnd  = getDaysAgoEnd(3)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getReportProfileResourceFunc,
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
				Config: testReportProfile_custom(name, twoWeeksAgoStart, oneWeekAgoEnd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "custom"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "subscription_type", "0"),
					resource.TestCheckResourceAttr(rName, "statistic_period.0.start_time", fmt.Sprintf("%d", twoWeeksAgoStart)),
					resource.TestCheckResourceAttr(rName, "statistic_period.0.end_time", fmt.Sprintf("%d", oneWeekAgoEnd)),
					resource.TestCheckResourceAttrPair(rName, "topic_urn", "huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "topic_name"),
				),
			},
			{
				Config: testReportProfile_custom_update(name, oneWeekAgoStart, threeDaysAgoEnd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "custom"),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "status", "0"),
					resource.TestCheckResourceAttr(rName, "subscription_type", "1"),
					resource.TestCheckResourceAttr(rName, "statistic_period.0.start_time", fmt.Sprintf("%d", oneWeekAgoStart)),
					resource.TestCheckResourceAttr(rName, "statistic_period.0.end_time", fmt.Sprintf("%d", threeDaysAgoEnd)),
					resource.TestCheckResourceAttrPair(rName, "topic_urn", "huaweicloud_smn_topic.test_another", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "topic_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testReportProfileImportState(rName),
				ImportStateVerifyIgnore: []string{
					"description",
				},
			},
		},
	})
}

func testReportProfile_custom(name string, startTime, endTime int64) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "%[1]s"
}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = "%[2]s"
  category          = "custom"
  name              = "%[1]s"
  topic_urn         = huaweicloud_smn_topic.test.topic_urn
  subscription_type = "0"
  status            = "1"
  description       = "test description"

  statistic_period {
    start_time = %[3]d
    end_time   = %[4]d
  }
}
`, name, acceptance.HW_CFW_INSTANCE_ID, startTime, endTime)
}

func testReportProfile_custom_update(name string, startTime, endTime int64) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test_another" {
  name         = "%[1]s_another"
  display_name = "%[1]s_another"
}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = "%[2]s"
  category          = "custom"
  name              = "%[1]s_update"
  topic_urn         = huaweicloud_smn_topic.test_another.topic_urn
  subscription_type = "1"
  status            = "0"
  description       = "test description update"

  statistic_period {
    start_time = %[3]d
    end_time   = %[4]d
  }
}
`, name, acceptance.HW_CFW_INSTANCE_ID, startTime, endTime)
}

func testReportProfileImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		fwInstanceId := rs.Primary.Attributes["fw_instance_id"]
		if fwInstanceId == "" {
			return "", fmt.Errorf("attribute (fw_instance_id) of Resource (%s) not found", name)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of Resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", fwInstanceId, rs.Primary.ID), nil
	}
}

func getDaysAgoStart(days int) int64 {
	now := time.Now()
	targetDay := now.AddDate(0, 0, -days)
	start := time.Date(
		targetDay.Year(),
		targetDay.Month(),
		targetDay.Day(),
		0, 0, 0, 0,
		now.Location(),
	)
	return start.UnixMilli()
}

func getDaysAgoEnd(days int) int64 {
	now := time.Now()
	targetDay := now.AddDate(0, 0, -days)
	end := time.Date(
		targetDay.Year(),
		targetDay.Month(),
		targetDay.Day(),
		23, 59, 59, 999999999,
		now.Location(),
	)
	return end.UnixMilli()
}
