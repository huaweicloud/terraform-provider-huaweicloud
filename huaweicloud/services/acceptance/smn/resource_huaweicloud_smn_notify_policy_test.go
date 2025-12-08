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

func getNotifyPolicy(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy"
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{topic_urn}", state.Primary.Attributes["topic_urn"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	id := utils.PathSearch("id", getRespBody, "").(string)
	if id == "" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccNotifyPolicy_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_smn_notify_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getNotifyPolicy,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSmnSubscribedTopicUrn(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNotifyPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "topic_urn",
						"data.huaweicloud_smn_subscriptions.test", "subscriptions.0.topic_urn"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "callnotify"),
					resource.TestCheckResourceAttr(resourceName, "polling.0.order", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "polling.0.subscription_urns.0",
						"data.huaweicloud_smn_subscriptions.test", "subscriptions.0.subscription_urn"),
					resource.TestCheckResourceAttrSet(resourceName, "polling.0.subscriptions.#"),
					resource.TestCheckResourceAttrSet(resourceName, "polling.0.subscriptions.0.subscription_urn"),
					resource.TestCheckResourceAttrSet(resourceName, "polling.0.subscriptions.0.endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "polling.0.subscriptions.0.remark"),
					resource.TestCheckResourceAttrSet(resourceName, "polling.0.subscriptions.0.status"),
				),
			},
			{
				Config: testAccNotifyPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "topic_urn",
						"data.huaweicloud_smn_subscriptions.test", "subscriptions.0.topic_urn"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "callnotify"),
					resource.TestCheckResourceAttr(resourceName, "polling.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "polling.0.order"),
					resource.TestCheckResourceAttrSet(resourceName, "polling.1.order"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testNotifyPolicyResourceImportState(resourceName),
			},
		},
	})
}

func testAccNotifyPolicy_basic() string {
	return `
data "huaweicloud_smn_subscriptions" "test" {
  protocol = "callnotify"
  status   = "1"
}

resource "huaweicloud_smn_notify_policy" "test" {
  topic_urn = data.huaweicloud_smn_subscriptions.test.subscriptions[0].topic_urn
  protocol  = "callnotify"
  polling {
    order             = 1
    subscription_urns = [data.huaweicloud_smn_subscriptions.test.subscriptions[0].subscription_urn]
  }
}
`
}

func testAccNotifyPolicy_update() string {
	return `
data "huaweicloud_smn_subscriptions" "test" {
  protocol = "callnotify"
  status   = "1"
}

resource "huaweicloud_smn_notify_policy" "test" {
  topic_urn = data.huaweicloud_smn_subscriptions.test.subscriptions[0].topic_urn
  protocol  = "callnotify"
  polling {
    order             = 1
    subscription_urns = [data.huaweicloud_smn_subscriptions.test.subscriptions[0].subscription_urn]
  }

  polling {
    order             = 2
    subscription_urns = [data.huaweicloud_smn_subscriptions.test.subscriptions[1].subscription_urn]
  }
}
`
}

func testNotifyPolicyResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		topicUrn := rs.Primary.Attributes["topic_urn"]
		return topicUrn, nil
	}
}
