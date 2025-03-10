package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceElbLoadbalancerFeatureConfigurations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_loadbalancer_feature_configurations.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbLoadbalancerFeatureConfigurations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "features.#"),
					resource.TestCheckResourceAttrSet(dataSource, "features.0.feature"),
					resource.TestCheckResourceAttrSet(dataSource, "features.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "features.0.value"),
				),
			},
		},
	})
}

func testDataSourceElbLoadbalancerFeatureConfigurations_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
 name            = "%[2]s"
 ipv4_subnet_id  = huaweicloud_vpc_subnet.test.ipv4_subnet_id

 availability_zone = [
   data.huaweicloud_availability_zones.test.names[0]
 ]
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceElbLoadbalancerFeatureConfigurations_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_elb_loadbalancer_feature_configurations" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}
`, testDataSourceElbLoadbalancerFeatureConfigurations_base(name))
}
