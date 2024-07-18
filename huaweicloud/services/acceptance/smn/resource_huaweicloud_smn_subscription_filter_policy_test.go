package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/smn"
)

func getResourceSubscriptionFilterPolicy(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.SmnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	return smn.GetSubscriptionFilterPolicies(client, state.Primary.ID)
}

func TestAccSubscriptionFilterPolicy_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_smn_subscription_filter_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceSubscriptionFilterPolicy,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSubscriptionFilterPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "filter_policies.#", "1"),
					resource.TestCheckResourceAttr(rName, "filter_policies.0.name", "alarm"),
					resource.TestCheckResourceAttr(rName, "filter_policies.0.string_equals.#", "1"),
					resource.TestCheckResourceAttr(rName, "filter_policies.0.string_equals.0", "os"),
				),
			},
			{
				Config: testAccSubscriptionFilterPolicy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "filter_policies.#", "2"),
					resource.TestCheckResourceAttr(rName, "filter_policies.0.name", "alarm"),
					resource.TestCheckResourceAttr(rName, "filter_policies.0.string_equals.#", "1"),
					resource.TestCheckResourceAttr(rName, "filter_policies.0.string_equals.0", "os_update"),
					resource.TestCheckResourceAttr(rName, "filter_policies.1.name", "service"),
					resource.TestCheckResourceAttr(rName, "filter_policies.1.string_equals.#", "2"),
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

func testAccSubscriptionFilterPolicy_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_subscription_filter_policy" "test" {
  subscription_urn = huaweicloud_smn_subscription.test.id

  filter_policies {
    name          = "alarm"
    string_equals = ["os"]
  }
}
`, testAccSubscription_build(rName))
}

func testAccSubscriptionFilterPolicy_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_subscription_filter_policy" "test" {
  subscription_urn = huaweicloud_smn_subscription.test.id

  filter_policies {
    name          = "alarm"
    string_equals = ["os_update"]
  }

  filter_policies {
    name          = "service"
    string_equals = ["test", "db"]
  }
}
`, testAccSubscription_build(rName))
}

func testAccSubscription_build(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%s"
  display_name = "The display name of topic_test"
}

resource "huaweicloud_smn_subscription" "test" {
  topic_urn = huaweicloud_smn_topic.test.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "O&M"
}
`, rName)
}
