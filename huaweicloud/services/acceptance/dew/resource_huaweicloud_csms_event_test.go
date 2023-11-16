package dew

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCsmsEventResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/csms/events/{event_name}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{event_name}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CSMS event: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func TestAccCsmsEvent_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_csms_event.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCsmsEventResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCsmsEvent_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "event_types.0", "SECRET_VERSION_CREATED"),
					resource.TestCheckResourceAttr(rName, "event_types.1", "SECRET_ROTATED"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLED"),
					resource.TestCheckResourceAttr(rName, "notification_target_type", "SMN"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_id",
						"huaweicloud_smn_topic.testA", "id"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_name",
						"huaweicloud_smn_topic.testA", "name"),
					resource.TestCheckResourceAttrSet(rName, "event_id"),
				),
			},
			{
				Config: testCsmsEvent_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "event_types.0", "SECRET_VERSION_EXPIRED"),
					resource.TestCheckResourceAttr(rName, "status", "DISABLED"),
					resource.TestCheckResourceAttr(rName, "notification_target_type", "SMN"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_id",
						"huaweicloud_smn_topic.testB", "id"),
					resource.TestCheckResourceAttrPair(rName, "notification_target_name",
						"huaweicloud_smn_topic.testB", "name"),
					resource.TestCheckResourceAttrSet(rName, "event_id"),
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

func testCsmsEvent_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "testA" {
  name         = "%[1]s_a"
  display_name = "The display nameA"
}

resource "huaweicloud_smn_topic" "testB" {
  name         = "%[1]s_b"
  display_name = "The display nameB"
}
`, name)
}

func testCsmsEvent_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_csms_event" "test" {
  name                     = "%s"
  event_types              = ["SECRET_VERSION_CREATED", "SECRET_ROTATED"]
  status                   = "ENABLED"
  notification_target_type = "SMN"
  notification_target_id   = huaweicloud_smn_topic.testA.id
  notification_target_name = huaweicloud_smn_topic.testA.name
}
`, testCsmsEvent_base(name), name)
}

func testCsmsEvent_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_csms_event" "test" {
  name                     = "%s"
  event_types              = ["SECRET_VERSION_EXPIRED"]
  status                   = "DISABLED"
  notification_target_type = "SMN"
  notification_target_id   = huaweicloud_smn_topic.testB.id
  notification_target_name = huaweicloud_smn_topic.testB.name
}
`, testCsmsEvent_base(name), name)
}
