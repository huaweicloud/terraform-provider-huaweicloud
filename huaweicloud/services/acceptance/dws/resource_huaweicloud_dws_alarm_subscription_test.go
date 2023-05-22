package dws

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDwsAlarmSubsResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDwsAlarmSubs: Query the DWS alarm subscription.
	var (
		getDwsAlarmSubsHttpUrl = "v2/{project_id}/alarm-subs"
		getDwsAlarmSubsProduct = "dws"
	)
	getDwsAlarmSubsClient, err := cfg.NewServiceClient(getDwsAlarmSubsProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	getDwsAlarmSubsPath := getDwsAlarmSubsClient.Endpoint + getDwsAlarmSubsHttpUrl
	getDwsAlarmSubsPath = strings.ReplaceAll(getDwsAlarmSubsPath, "{project_id}", getDwsAlarmSubsClient.ProjectID)

	getDwsAlarmSubsResp, err := pagination.ListAllItems(
		getDwsAlarmSubsClient,
		"offset",
		getDwsAlarmSubsPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsAlarmSubs: %s", err)
	}

	getDwsAlarmSubsRespJson, err := json.Marshal(getDwsAlarmSubsResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsEventSubs: %s", err)
	}
	var getDwsAlarmSubsRespBody interface{}
	err = json.Unmarshal(getDwsAlarmSubsRespJson, &getDwsAlarmSubsRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsEventSubs: %s", err)
	}

	jsonPath := fmt.Sprintf("alarm_subscriptions[?id=='%s']|[0]", state.Primary.ID)
	rawData := utils.PathSearch(jsonPath, getDwsAlarmSubsRespBody, nil)
	if rawData == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rawData, nil
}

func TestAccDwsAlarmSubs_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dws_alarm_subscription.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDwsAlarmSubsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDwsAlarmSubs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enable", "1"),
					resource.TestCheckResourceAttrPair(rName, "notification_target", "huaweicloud_smn_topic.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_name", "huaweicloud_smn_topic.test", "name"),
					resource.TestCheckResourceAttr(rName, "notification_target_type", "SMN"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "urgent,important"),
					resource.TestCheckResourceAttr(rName, "time_zone", "GMT+09:00"),
				),
			},
			{
				Config: testDwsAlarmSubs_basic_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "enable", "0"),
					resource.TestCheckResourceAttrPair(rName, "notification_target", "huaweicloud_smn_topic.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_name", "huaweicloud_smn_topic.test", "name"),
					resource.TestCheckResourceAttr(rName, "notification_target_type", "SMN"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "urgent,important,minor"),
					resource.TestCheckResourceAttr(rName, "time_zone", "GMT+09:00"),
				),
			},
			{
				Config: testDwsAlarmSubs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enable", "1"),
					resource.TestCheckResourceAttrPair(rName, "notification_target", "huaweicloud_smn_topic.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_name", "huaweicloud_smn_topic.test", "name"),
					resource.TestCheckResourceAttr(rName, "notification_target_type", "SMN"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "urgent,important"),
					resource.TestCheckResourceAttr(rName, "time_zone", "GMT+09:00"),
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

func testDwsAlarmSubs_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%s"
  display_name = "The display name of topic"
}

resource "huaweicloud_dws_alarm_subscription" "test" {
  name                     = "%s"
  enable                   = "1"
  notification_target      = huaweicloud_smn_topic.test.id
  notification_target_type = "SMN"
  notification_target_name = huaweicloud_smn_topic.test.name
  time_zone                = "GMT+09:00"
  alarm_level              = "urgent,important"
}
`, name, name)
}

func testDwsAlarmSubs_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%s"
  display_name = "The display name of topic"
}

resource "huaweicloud_dws_alarm_subscription" "test" {
  name                     = "%s"
  enable                   = "0"
  notification_target      = huaweicloud_smn_topic.test.id
  notification_target_type = "SMN"
  notification_target_name = huaweicloud_smn_topic.test.name
  time_zone                = "GMT+09:00"
  alarm_level              = "urgent,important,minor"
}
`, name, name)
}
