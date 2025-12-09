package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/smn/v2/subscriptions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceSMNSubscription(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	smnClient, err := conf.SmnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	foundList, err := subscriptions.List(smnClient).Extract()
	if err != nil {
		return nil, err
	}

	var subscription *subscriptions.SubscriptionGet
	urn := state.Primary.ID
	for i := range foundList {
		if foundList[i].SubscriptionUrn == urn {
			subscription = &foundList[i]
		}
	}

	if subscription == nil {
		return nil, fmt.Errorf("the subscription does not exist")
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
					resource.TestCheckResourceAttr(resourceName1, "remark", "rest remark"),
					resource.TestCheckResourceAttr(resourceName2, "endpoint", "13600000000"),
					resource.TestCheckResourceAttr(resourceName2, "remark", "rest remark"),
					resource.TestCheckResourceAttr(resourceName3, "endpoint", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName3, "remark", "rest remark"),
					resource.TestCheckResourceAttr(resourceName3, "extension.#", "1"),
					resource.TestCheckResourceAttr(resourceName3, "extension.0.header.%", "1"),
					resource.TestCheckResourceAttr(resourceName3, "extension.0.header.X-Custom-Header", "test"),
					resource.TestCheckResourceAttr(resourceName4, "remark", "rest remark"),
					resource.TestCheckResourceAttrPair(
						resourceName4, "endpoint", "huaweicloud_fgs_function.test", "urn"),
				),
			},
			{
				Config: testAccSMNV2SubscriptionConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					rc2.CheckResourceExists(),
					rc3.CheckResourceExists(),
					rc4.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName1, "endpoint", "mailtest@gmail.com"),
					resource.TestCheckResourceAttr(resourceName1, "remark", "rest remark update"),
					resource.TestCheckResourceAttr(resourceName2, "endpoint", "13600000000"),
					resource.TestCheckResourceAttr(resourceName2, "remark", "rest remark update"),
					resource.TestCheckResourceAttr(resourceName3, "endpoint", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName3, "remark", ""),
					resource.TestCheckResourceAttr(resourceName3, "extension.#", "1"),
					resource.TestCheckResourceAttr(resourceName3, "extension.0.header.%", "1"),
					resource.TestCheckResourceAttr(resourceName3, "extension.0.header.X-Custom-Header", "test"),
					resource.TestCheckResourceAttr(resourceName4, "remark", ""),
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
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	smnClient, err := cfg.SmnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating SMN client: %s", err)
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

		urn := rs.Primary.ID
		for i := range foundList {
			if foundList[i].SubscriptionUrn == urn {
				subscription = &foundList[i]
			}
		}
		if subscription != nil {
			return fmt.Errorf("subscription still exists")
		}
	}

	return nil
}

func testAccSMNV2SubscriptionConfig_basic(rName string) string {
	funcCode := "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnIC" +
		"sganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "%[2]s"
}

resource "huaweicloud_smn_topic" "topic_1" {
  name         = "%[1]s"
  display_name = "The display name of topic_1"
}

resource "huaweicloud_smn_subscription" "subscription_1" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "rest remark"
}

resource "huaweicloud_smn_subscription" "subscription_2" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "13600000000"
  protocol  = "sms"
  remark    = "rest remark"
}

resource "huaweicloud_smn_subscription" "subscription_3" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "https://test.com"
  protocol  = "https"
  remark    = "rest remark"

  extension {
    header = {
      "X-Custom-Header" = "test"
    }
  }
}

resource "huaweicloud_smn_subscription" "subscription_4" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = huaweicloud_fgs_function.test.urn
  protocol  = "functionstage"
  remark    = "rest remark"
}
`, rName, funcCode)
}

func testAccSMNV2SubscriptionConfig_update(rName string) string {
	funcCode := "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnIC" +
		"sganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "%[2]s"
}

resource "huaweicloud_smn_topic" "topic_1" {
  name         = "%[1]s"
  display_name = "The display name of topic_1"
}

resource "huaweicloud_smn_subscription" "subscription_1" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "rest remark update"
}

resource "huaweicloud_smn_subscription" "subscription_2" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "13600000000"
  protocol  = "sms"
  remark    = "rest remark update"
}

resource "huaweicloud_smn_subscription" "subscription_3" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = "https://test.com"
  protocol  = "https"

  extension {
    header = {
      "X-Custom-Header" = "test"
    }
  }
}

resource "huaweicloud_smn_subscription" "subscription_4" {
  topic_urn = huaweicloud_smn_topic.topic_1.id
  endpoint  = huaweicloud_fgs_function.test.urn
  protocol  = "functionstage"
}
`, rName, funcCode)
}
