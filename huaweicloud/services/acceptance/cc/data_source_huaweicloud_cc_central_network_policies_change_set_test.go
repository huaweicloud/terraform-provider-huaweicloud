package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCentralNetworkPoliciesChangeSet_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_central_network_policies_change_set.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcCentralNetworkPoliciesChangeSet_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_policy_change_set.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_network_policy_change_set.0.change_content.%"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcCentralNetworkPoliciesChangeSet_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cc_central_network_policies_change_set" "test" {
  depends_on = [
    huaweicloud_cc_central_network_policy.test2,
    huaweicloud_cc_central_network_policy_apply.test,
  ]
	
  central_network_id = huaweicloud_cc_central_network.test.id
  policy_id          = huaweicloud_cc_central_network_policy.test2.id
}
`, testCentralNetworkPolicies_dataBasic(name))
}
