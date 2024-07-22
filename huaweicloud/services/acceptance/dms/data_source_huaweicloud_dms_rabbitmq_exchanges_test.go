package dms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRabbitmqExchanges_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rabbitmq_exchanges.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsRabbitmqExchanges_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "exchanges.#"),
					resource.TestCheckResourceAttrSet(dataSource, "exchanges.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "exchanges.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "exchanges.0.auto_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "exchanges.0.durable"),
					resource.TestCheckResourceAttrSet(dataSource, "exchanges.0.internal"),
					resource.TestCheckResourceAttrSet(dataSource, "exchanges.0.default"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsRabbitmqExchanges_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rabbitmq_exchanges" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_exchange.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.test.name
}
`, testRabbitmqExchange_basic(name))
}
