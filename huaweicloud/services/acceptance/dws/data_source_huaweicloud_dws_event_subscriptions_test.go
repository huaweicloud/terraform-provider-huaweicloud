package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEventSubscriptionsDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_dws_event_subscriptions.name_filter"
	dc := acceptance.InitDataSourceCheck(resourceName)
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionsDataSourceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "event_subscriptions.#"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("source_type_filter_is_useful", "true"),
					resource.TestCheckOutput("category_filter_is_useful", "true"),
					resource.TestCheckOutput("notificate_filter_is_useful", "true"),
					resource.TestCheckOutput("enable_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccEventSubscriptionsDataSourceBasic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "The display name of topic"
}

resource "huaweicloud_dws_event_subscription" "test" {
  name                     = "%[1]s"
  enable                   = 0
  notification_target      = huaweicloud_smn_topic.test.id
  notification_target_type = "SMN"
  notification_target_name = huaweicloud_smn_topic.test.name
  category                 = "management,security"
  severity                 = "normal,warning"
  source_type              = "cluster,disaster-recovery"
  time_zone                = "GMT+09:00"
}

data "huaweicloud_dws_event_subscriptions" "name_filter" {
  name = huaweicloud_dws_event_subscription.test.name
}

data "huaweicloud_dws_event_subscriptions" "notificate_filter" {
  notification_target_name = huaweicloud_dws_event_subscription.test.notification_target_name
}

data "huaweicloud_dws_event_subscriptions" "enable_filter" {
  enable = huaweicloud_dws_event_subscription.test.enable
}

data "huaweicloud_dws_event_subscriptions" "source_type_filter" {
  source_type = "cluster,disaster-recovery"

  depends_on = [huaweicloud_dws_event_subscription.test]
}

data "huaweicloud_dws_event_subscriptions" "category_filter" {
  category = "security,management"

  depends_on = [huaweicloud_dws_event_subscription.test]
}

locals {
  name_filter = [for v in data.huaweicloud_dws_event_subscriptions.name_filter.event_subscriptions[*].name :
  v == huaweicloud_dws_event_subscription.test.name]

  notificate_filter = [for v in data.huaweicloud_dws_event_subscriptions.name_filter.event_subscriptions[*] :
    v.notification_target_name == huaweicloud_dws_event_subscription.test.notification_target_name]

  enable_filter = [for v in data.huaweicloud_dws_event_subscriptions.name_filter.event_subscriptions[*] :
    v.enable == huaweicloud_dws_event_subscription.test.enable]

  source_type_filter = [for v in data.huaweicloud_dws_event_subscriptions.source_type_filter.event_subscriptions[*] :
    length(regexall(replace(v.source_type, ",", "|"), huaweicloud_dws_event_subscription.test.source_type)) ==
  length(split(",", huaweicloud_dws_event_subscription.test.source_type))]

  category_filter = [for v in data.huaweicloud_dws_event_subscriptions.category_filter.event_subscriptions[*] :
    length(regexall(replace(v.category, ",", "|"), huaweicloud_dws_event_subscription.test.category)) ==
  length(split(",", huaweicloud_dws_event_subscription.test.category))]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter) && length(local.name_filter) > 0
}

output "notificate_filter_is_useful" {
  value = alltrue(local.notificate_filter) && length(local.notificate_filter) > 0
}
output "enable_filter_is_useful" {
  value = alltrue(local.enable_filter) && length(local.enable_filter) > 0
}

output "source_type_filter_is_useful" {
  value = alltrue(local.source_type_filter) && length(local.source_type_filter) > 0
}

output "category_filter_is_useful" {
  value = alltrue(local.category_filter) && length(local.category_filter) > 0
}
`, name)
}
