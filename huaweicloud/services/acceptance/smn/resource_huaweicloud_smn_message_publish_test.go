package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccMessagePublish_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSmnMessagePublish_basic(name),
			},
		},
	})
}

func testSmnMessagePublish_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_message_publish" "test1" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  subject   = "test1"
  message   = "test created by terraform"

  message_attributes {
    name  = "test"
    type  = "STRING"
    value = "aaa"
  }
}

resource "huaweicloud_smn_message_publish" "test2" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  subject   = "test2"

  message_structure = jsonencode({
    default       = "test"
    sms           = "test"
    email         = "test"
    http          = "test"
    functiongraph = "test"
    https         = "test"
  })

  message_attributes {
    name   = "aaa"
    type   = "STRING_ARRAY"
    values = ["aaa", "aaaa"]
  }
}

resource "huaweicloud_smn_message_publish" "test3" {
  topic_urn             = huaweicloud_smn_topic.topic_1.id
  subject               = "test1"
  message_template_name = "test"

  tags = {
    key = "value"
  }

  message_attributes {
    name   = "smn_protocol"
    type   = "PROTOCOL"
    values = ["https", "http", "email", "sms"]
  }
}
`, testAccSMNV2TopicConfig_basic(name))
}
