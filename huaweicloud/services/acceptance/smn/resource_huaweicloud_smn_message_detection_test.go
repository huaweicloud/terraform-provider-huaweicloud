package smn

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

func getMessageDetectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("smn", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	getMessageDetectionHttpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/detection/{task_id}"
	getMessageDetectionPath := client.Endpoint + getMessageDetectionHttpUrl
	getMessageDetectionPath = strings.ReplaceAll(getMessageDetectionPath, "{project_id}", client.ProjectID)
	getMessageDetectionPath = strings.ReplaceAll(getMessageDetectionPath, "{topic_urn}", state.Primary.Attributes["topic_urn"])
	getMessageDetectionPath = strings.ReplaceAll(getMessageDetectionPath, "{task_id}", state.Primary.ID)

	getMessageDetectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getMessageDetectionResp, err := client.Request("GET", getMessageDetectionPath, &getMessageDetectionOpt)
	if err != nil {
		return nil, err
	}

	getMessageDetectionRespBody, err := utils.FlattenResponse(getMessageDetectionResp)
	if err != nil {
		return nil, err
	}

	return getMessageDetectionRespBody, nil
}

func TestAccMessageDetection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_smn_message_detection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getMessageDetectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSmnMessageDetection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "result"),
				),
			},
		},
	})
}

func testSmnMessageDetection_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_message_detection" "test" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  protocol  = "https"
  endpoint  = "https://example.com/notification/action"
}
`, testAccSMNV2TopicConfig_basic(name))
}
