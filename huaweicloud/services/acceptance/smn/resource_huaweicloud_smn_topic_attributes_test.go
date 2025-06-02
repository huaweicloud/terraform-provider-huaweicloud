package smn

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/smn"
)

func getTopicAttributesFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SmnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	return smn.GetTopicAttributes(client, state.Primary.Attributes["topic_urn"],
		state.Primary.Attributes["name"])
}

func TestAccTopicAttributes_basic(t *testing.T) {
	var (
		obj interface{}

		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_smn_topic_attributes.test"

		rc = acceptance.InitResourceCheck(resourceName, &obj, getTopicAttributesFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTopicAttributes_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "access_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "value"),
					resource.TestCheckResourceAttrSet(resourceName, "topic_urn"),
				),
			},
			{
				Config: testTopicAttributes_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "access_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"value",
				},
				ImportStateIdFunc: testTopicAttributesImportState(resourceName),
			},
		},
	})
}

func testTopic_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%s"
}
`, name)
}

func testTopicAttributes_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic_attributes" "test" {
  topic_urn = huaweicloud_smn_topic.test.topic_urn
  name      = "access_policy"
  value     = jsonencode({
    "Version": "2016-09-07",
    "Id": "__default_policy_ID",
    "Statement": [
      {
        "Sid": "__org_path_pub_0",
        "Effect": "Allow",
        "Principal": {
          "OrgPath": [
            "o-xxx/r-xxx/ou-xxx"
          ]
        },
        "Action": [
          "SMN:Publish",
          "SMN:QueryTopicDetail"
        ],
        "Resource": huaweicloud_smn_topic.test.topic_urn
    }]
  })
}
`, testTopic_base(name))
}

func testTopicAttributes_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic_attributes" "test" {
  topic_urn = huaweicloud_smn_topic.test.topic_urn
  name      = "access_policy"
  value     = jsonencode({
    "Version": "2016-09-07",
    "Id": "__default_policy_ID",
    "Statement": [
      {
        "Sid": "__org_path_pub_0",
        "Effect": "Allow",
        "Principal": {
          "OrgPath": [
            "o-xxx/r-xxx/ou-yyy"
          ]
        },
        "Action": [
          "SMN:Publish",
          "SMN:QueryTopicDetail"
        ],
        "Resource": huaweicloud_smn_topic.test.topic_urn
    }]
  })
}
`, testTopic_base(name))
}

func testTopicAttributesImportState(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		topicUrn := rs.Primary.Attributes["topic_urn"]
		name := rs.Primary.Attributes["name"]
		if topicUrn == "" || name == "" {
			return "", errors.New("topic_urn or name is missing")
		}
		return fmt.Sprintf("%s/%s", topicUrn, name), nil
	}
}
