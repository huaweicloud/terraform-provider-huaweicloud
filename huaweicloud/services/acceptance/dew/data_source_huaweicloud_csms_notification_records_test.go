package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNotificationRecords_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_csms_notification_records.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceNotificationRecords_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.event_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.trigger_event_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.secret_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.secret_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.notification_target_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.notification_target_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.notification_content"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.notification_status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.create_time"),
				),
			},
		},
	})
}

const testDataSourceNotificationRecords_basic = `data "huaweicloud_csms_notification_records" "test" {}`
