package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterSubscriptionResource_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_subscription_resource.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterSubscriptionResource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "sku_attribute"),
					resource.TestCheckResourceAttrSet(dataSource, "upper_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "unit"),
					resource.TestCheckResourceAttrSet(dataSource, "step"),
					resource.TestCheckResourceAttrSet(dataSource, "used_amount"),
					resource.TestCheckResourceAttrSet(dataSource, "unused_amount"),
					resource.TestCheckResourceAttrSet(dataSource, "version"),
					resource.TestCheckResourceAttrSet(dataSource, "index_storage_upper_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "index_shards_upper_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "index_shards_unused"),
					resource.TestCheckResourceAttrSet(dataSource, "partitions_unused"),
					resource.TestCheckResourceAttrSet(dataSource, "partition_upper_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "update_time"),
				),
			},
		},
	})
}

func testDataSourceSecmasterSubscriptionResource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_subscription_resource" "test" {
  workspace_id = "%s"
  sku          = "FLOW_DATA_BANDWIDTH"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
