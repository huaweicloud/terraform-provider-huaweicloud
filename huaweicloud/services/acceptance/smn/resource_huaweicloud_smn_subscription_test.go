package smn

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/smn/v2/subscriptions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceSMNSubscription(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	smnClient, err := conf.SmnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating smn: %s", err)
	}

	var subscription *subscriptions.SubscriptionGet
	foundList, err := subscriptions.List(smnClient).Extract()
	if err != nil {
		return nil, err
	}
	for _, subObject := range foundList {
		if subObject.SubscriptionUrn == state.Primary.ID {
			subscription = &subObject
		}
	}

	return subscription, nil
}

func TestAccSMNV2Subscription_basic(t *testing.T) {
	var subscription1 subscriptions.SubscriptionGet
	var subscription2 subscriptions.SubscriptionGet
	var subscription3 subscriptions.SubscriptionGet
	var subscription4 subscriptions.SubscriptionGet
	resourceName1 := "huaweicloud_smn_subscription.subscription_1"
	resourceName2 := "huaweicloud_smn_subscription.subscription_2"
	resourceName3 := "huaweicloud_smn_subscription.subscription_3"
	resourceName4 := "huaweicloud_smn_subscription.subscription_4"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc1 := acceptance.InitResourceCheck(
		resourceName1,
		&subscription1,
		getResourceSMNSubscription,
	)

	rc2 := acceptance.InitResourceCheck(
		resourceName2,
		&subscription2,
		getResourceSMNSubscription,
	)

	rc3 := acceptance.InitResourceCheck(
		resourceName3,
		&subscription3,
		getResourceSMNSubscription,
	)

	rc4 := acceptance.InitResourceCheck(
		resourceName4,
		&subscription4,
		getResourceSMNSubscription,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSMNSubscriptionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2SubscriptionConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					rc2.CheckResourceExists(),
					rc3.CheckResourceExists(),
					rc4.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName1, "endpoint", "mailtest@gmail.com"),
					resource.TestCheckResourceAttr(resourceName2, "endpoint", "13600000000"),
					resource.TestCheckResourceAttr(resourceName3, "endpoint", "https://test.com"),
					resource.TestCheckResourceAttrPair(
						resourceName4, "endpoint", "huaweicloud_fgs_function.test", "urn"),
				),
			},
			{
				ResourceName:      resourceName1,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSMNSubscriptionV2Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	smnClient, err := config.SmnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating SMN: %s", err)
	}

	foundList, err := subscriptions.List(smnClient).Extract()
	if err != nil {
		return err
	}

	var subscription *subscriptions.SubscriptionGet
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_smn_subscription" {
			continue
		}

		for _, subObject := range foundList {
			if subObject.SubscriptionUrn == rs.Primary.ID {
				subscription = &subObject
			}
		}
		if subscription != nil {
			return fmt.Errorf("subscription still exists")
		}
	}

	return nil
}

func testAccSMNV2SubscriptionConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}

resource "huaweicloud_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The display name of topic_1"
}

resource "huaweicloud_smn_subscription" "subscription_1" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "O&M"
}

resource "huaweicloud_smn_subscription" "subscription_2" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "13600000000"
  protocol  = "sms"
  remark    = "O&M"
}

resource "huaweicloud_smn_subscription" "subscription_3" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "https://test.com"
  protocol  = "https"
  remark    = "O&M"
}

resource "huaweicloud_smn_subscription" "subscription_4" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = huaweicloud_fgs_function.test.urn
  protocol  = "functionstage"
  remark    = "O&M"
}
`, rName, rName)
}
