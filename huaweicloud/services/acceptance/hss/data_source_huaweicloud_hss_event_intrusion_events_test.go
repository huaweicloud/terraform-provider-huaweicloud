package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventIntrusionEvents_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_event_intrusion_events.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Because there is no available data for testing, the test case is only
			// used to verify that the API can be invoked.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventIntrusionEvents_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceEventIntrusionEvents_basic string = `
data "huaweicloud_hss_event_intrusion_events" "test" {
  category              = "host"
  enterprise_project_id = "all_granted_eps"
}
`
