package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataNetworkAvailableSubnets_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_modelarts_network_available_subnets.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsNetworkAvailableSubnets(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNetworkAvailableSubnets_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "network_id"),
					resource.TestMatchResourceAttr(dataSource, "subnets.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSource, "subnets.0.cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "subnets.0.ip_version"),
					resource.TestCheckResourceAttrSet(dataSource, "subnets.0.used_ips"),
					resource.TestCheckResourceAttrSet(dataSource, "subnets.0.total_ips"),
				),
			},
		},
	})
}

func testAccDataNetworkAvailableSubnets_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_modelarts_network_available_subnets" "test" {
  network_name = "%[1]s"
  subnet_id    = "%[2]s"
}
`, acceptance.HW_MODELARTS_NETWORK_NAME, acceptance.HW_MODELARTS_SUBNET_ID)
}
