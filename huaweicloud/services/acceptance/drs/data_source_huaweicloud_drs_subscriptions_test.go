package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsSubscriptions_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_subscriptions.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsSubscriptions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "subscriptions.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subscriptions.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subscriptions.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subscriptions.0.created_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subscriptions.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subscriptions.0.job_action.0.available_actions.#"),
				),
			},
		},
	})
}

const testAccDataSourceDrsSubscriptions_basic = `
data "huaweicloud_drs_subscriptions" "test" {
  job_type = "subscription"
}
`
