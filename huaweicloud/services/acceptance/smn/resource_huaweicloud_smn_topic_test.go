package smn

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/smn/v2/topics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceSMNTopic(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	smnClient, err := conf.SmnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating smn: %s", err)
	}

	return topics.Get(smnClient, state.Primary.ID).Extract()
}

func TestAccSMNV2Topic_basic(t *testing.T) {
	var topic topics.TopicGet
	resourceName := "huaweicloud_smn_topic.topic_1"
	rName := acceptance.RandomAccResourceNameWithDash()
	displayName := fmt.Sprintf("The display name of %s", rName)
	update_displayName := fmt.Sprintf("The update display name of %s", rName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&topic,
		getResourceSMNTopic,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccSMNV2TopicConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", update_displayName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSMNV2Topic_withEpsID(t *testing.T) {
	var topic topics.TopicGet
	resourceName := "huaweicloud_smn_topic.topic_1"
	rName := acceptance.RandomAccResourceNameWithDash()
	displayName := fmt.Sprintf("The display name of %s", rName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&topic,
		getResourceSMNTopic,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_withEpsID(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(
						resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccSMNV2TopicConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The display name of %s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, rName)
}

func testAccSMNV2TopicConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The update display name of %s"

  tags = {
    foo        = "bar_update"
    key_update = "value"
  }
}
`, rName, rName)
}

func testAccSMNV2TopicConfig_withEpsID(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name                  = "%s"
  display_name          = "The display name of %s"
  enterprise_project_id = "%s"
}
`, rName, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
