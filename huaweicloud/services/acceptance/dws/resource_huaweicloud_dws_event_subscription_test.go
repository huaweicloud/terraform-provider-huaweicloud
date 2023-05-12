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

func getDwsEventSubsResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDwsEventSubs: Query the DWS event subscription.
	var (
		getDwsEventSubsHttpUrl = "v2/{project_id}/event-subs"
		getDwsEventSubsProduct = "dws"
	)
	getDwsEventSubsClient, err := cfg.NewServiceClient(getDwsEventSubsProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	getDwsEventSubsPath := getDwsEventSubsClient.Endpoint + getDwsEventSubsHttpUrl
	getDwsEventSubsPath = strings.ReplaceAll(getDwsEventSubsPath, "{project_id}", getDwsEventSubsClient.ProjectID)

	getDwsEventSubsResp, err := pagination.ListAllItems(
		getDwsEventSubsClient,
		"offset",
		getDwsEventSubsPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsEventSubs: %s", err)
	}

	getDwsEventSubsRespJson, err := json.Marshal(getDwsEventSubsResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsEventSubs: %s", err)
	}
	var getDwsEventSubsRespBody interface{}
	err = json.Unmarshal(getDwsEventSubsRespJson, &getDwsEventSubsRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsEventSubs: %s", err)
	}

	jsonPath := fmt.Sprintf("event_subscriptions[?id=='%s']|[0]", state.Primary.ID)
	rawData := utils.PathSearch(jsonPath, getDwsEventSubsRespBody, nil)
	if rawData == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rawData, nil
}

func TestAccDwsEventSubs_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dws_event_subscription.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDwsEventSubsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDwsEventSubs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enable", "1"),
					resource.TestCheckResourceAttrPair(rName, "notification_target", "huaweicloud_smn_topic.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_name", "huaweicloud_smn_topic.test", "name"),
					resource.TestCheckResourceAttr(rName, "notification_target_type", "SMN"),
					resource.TestCheckResourceAttr(rName, "category", "management,monitor,security"),
					resource.TestCheckResourceAttr(rName, "severity", "normal,warning"),
					resource.TestCheckResourceAttr(rName, "source_type", "cluster,backup,disaster-recovery"),
					resource.TestCheckResourceAttr(rName, "time_zone", "GMT+09:00"),
				),
			},
			{
				Config: testDwsEventSubs_basic_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "enable", "0"),
					resource.TestCheckResourceAttrPair(rName, "notification_target", "huaweicloud_smn_topic.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_name", "huaweicloud_smn_topic.test", "name"),
					resource.TestCheckResourceAttr(rName, "notification_target_type", "SMN"),
					resource.TestCheckResourceAttr(rName, "category", "management,monitor"),
					resource.TestCheckResourceAttr(rName, "severity", "normal"),
					resource.TestCheckResourceAttr(rName, "source_type", "cluster,backup"),
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

func testDwsEventSubs_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%s"
  display_name = "The display name of topic"
}

resource "huaweicloud_dws_event_subscription" "test" {
  name                     = "%s"
  enable                   = 1
  notification_target      = huaweicloud_smn_topic.test.id
  notification_target_type = "SMN"
  notification_target_name = huaweicloud_smn_topic.test.name
  category                 = "management,monitor,security"
  severity                 = "normal,warning"
  source_type              = "cluster,backup,disaster-recovery"
  time_zone                = "GMT+09:00"
}
`, name, name)
}

func testDwsEventSubs_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%s"
  display_name = "The display name of topic"
}

resource "huaweicloud_dws_event_subscription" "test" {
  name                     = "%s"
  enable                   = 0
  notification_target      = huaweicloud_smn_topic.test.id
  notification_target_type = "SMN"
  notification_target_name = huaweicloud_smn_topic.test.name
  category                 = "management,monitor"
  severity                 = "normal"
  source_type              = "cluster,backup"
  time_zone                = "GMT+09:00"
}
`, name, name)
}
