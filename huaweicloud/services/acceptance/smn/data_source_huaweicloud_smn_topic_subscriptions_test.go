package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmnTopicSubscriptions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_smn_topic_subscriptions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmnTopicSubscriptions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.endpoint"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmnTopicSubscriptions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_smn_topic_subscriptions" "test" {
  depends_on = [ 
    huaweicloud_smn_subscription.subscription_1,
    huaweicloud_smn_subscription.subscription_2,
    huaweicloud_smn_subscription.subscription_3,
    huaweicloud_smn_subscription.subscription_4,
  ]

  topic_urn = huaweicloud_smn_topic.topic_1.id
}
`, testAccSMNV2SubscriptionConfig_basic(name))
}
