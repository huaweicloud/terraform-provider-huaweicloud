package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventOperations_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_event_operations.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEventOperations_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "operate_accept_list.#"),
				),
			},
		},
	})
}

const testDataSourceEventOperations_basic string = `
data "huaweicloud_hss_event_operations" "test" {
  event_type            = "1001"
  enterprise_project_id = "0"
}
`
