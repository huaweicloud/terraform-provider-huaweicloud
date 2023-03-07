package as

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

func getAsNotificationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getASNotification: Query the AS notification.
	var (
		getASNotificationHttpUrl = "autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}"
		getASNotificationProduct = "autoscaling"
	)
	getASNotificationClient, err := cfg.NewServiceClient(getASNotificationProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AutoScaling Client: %s", err)
	}

	getASNotificationPath := getASNotificationClient.Endpoint + getASNotificationHttpUrl
	getASNotificationPath = strings.ReplaceAll(getASNotificationPath, "{project_id}",
		getASNotificationClient.ProjectID)
	getASNotificationPath = strings.ReplaceAll(getASNotificationPath, "{scaling_group_id}",
		fmt.Sprintf("%v", state.Primary.Attributes["scaling_group_id"]))

	getASNotificationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getASNotificationResp, err := getASNotificationClient.Request("GET", getASNotificationPath,
		&getASNotificationOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AS notification: %s", err)
	}

	getASNotificationRespBody, err := utils.FlattenResponse(getASNotificationResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten response: %s", err)
	}
	return filterTargetASNotificationByTopicUrn(getASNotificationRespBody, state.Primary.ID)
}

func filterTargetASNotificationByTopicUrn(resp interface{}, topicUrn string) (interface{}, error) {
	if resp == nil {
		return nil, fmt.Errorf("resp cannot be nil")
	}

	curJson := utils.PathSearch("topics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		urn := utils.PathSearch("topic_urn", v, "")
		if topicUrn == urn.(string) {
			return v, nil
		}
	}
	return nil, fmt.Errorf("the target AS notification: %s not found", topicUrn)
}

func TestAccAsNotification_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_as_notification.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAsNotificationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAsNotification_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "events.0", "SCALING_UP"),
					resource.TestCheckResourceAttrSet(rName, "scaling_group_id"),
					resource.TestCheckResourceAttrSet(rName, "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "topic_name"),
				),
			},
			{
				Config: testAsNotification_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "events.0", "SCALING_DOWN_FAIL"),
					resource.TestCheckResourceAttr(rName, "events.1", "SCALING_GROUP_ABNORMAL"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAsNotificationImportState(rName),
			},
		},
	})
}

func testAsNotification_base(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The display name of %s"
}

`, testASGroup_basic(name), name, name)
}

func testAsNotification_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_notification" "test" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  topic_urn        = huaweicloud_smn_topic.topic_1.id
  events           = ["SCALING_UP"]
}
`, testAsNotification_base(name))
}

func testAsNotification_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_notification" "test" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  topic_urn        = huaweicloud_smn_topic.topic_1.id
  events           = ["SCALING_DOWN_FAIL", "SCALING_GROUP_ABNORMAL"]
}
`, testAsNotification_base(name))
}

func testAsNotificationImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		var scalingGroupID, topicUrn string
		if scalingGroupID = rs.Primary.Attributes["scaling_group_id"]; scalingGroupID == "" {
			return "", fmt.Errorf("attribute (scaling_group_id) of Resource (%s) not found: %s", name, rs)
		}
		if topicUrn = rs.Primary.Attributes["topic_urn"]; topicUrn == "" {
			return "", fmt.Errorf("attribute (topic_urn) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", scalingGroupID, topicUrn), nil
	}
}
