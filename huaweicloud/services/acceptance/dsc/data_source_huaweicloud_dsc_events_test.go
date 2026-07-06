package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscEvents_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_events.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscEvents_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "events.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.event_name"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.event_level"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.event_status"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.event_type"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.occur_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.scheduled_close_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.source_module"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.verification_status"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.affected_asset.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.responsible_person.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.source_instance.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.related_alarm_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceDscEvents_basic = `
data "huaweicloud_dsc_events" "test" {}
`
