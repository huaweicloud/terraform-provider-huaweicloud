package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRabbitmqExchanges_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_dms_rabbitmq_exchanges.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		withArguments   = "data.huaweicloud_dms_rabbitmq_exchanges.with_arguments"
		dcWithArguments = acceptance.InitDataSourceCheck(withArguments)
	)

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
					dcWithArguments.CheckResourceExists(),
					resource.TestCheckOutput("arguments_is_set_and_valid", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceExchanges_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_rabbitmq_flavors" "with_arguments" {
  type = "single.professional"
}

locals {
  flavor_with_arguments = data.huaweicloud_dms_rabbitmq_flavors.with_arguments.flavors[0]
}

resource "huaweicloud_dms_rabbitmq_instance" "with_arguments" {
  name              = "%[2]s_with_arguments"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor_with_arguments.id
  engine_version    = "AMQP-0-9-1"
  storage_space     = local.flavor_with_arguments.properties[0].min_storage_per_node
  storage_spec_code = local.flavor_with_arguments.ios[0].storage_spec_code
}

resource "huaweicloud_dms_rabbitmq_vhost" "with_arguments" {
  instance_id = huaweicloud_dms_rabbitmq_instance.with_arguments.id
  name        = "%[2]s_with_arguments"
}

resource "huaweicloud_dms_rabbitmq_exchange" "with_arguments" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.with_arguments]

  instance_id = huaweicloud_dms_rabbitmq_instance.with_arguments.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.with_arguments.name
  name        = "%[2]s_with_arguments"
  type        = "x-delayed-message"
  auto_delete = false

  arguments   = jsonencode({
    "x-delayed-type" = "header"
  })
}
`, testRabbitmqExchange_basic(name), name)
}

func testDataSourceDataSourceDmsRabbitmqExchanges_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_rabbitmq_exchanges" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_exchange.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.test.name
}

data "huaweicloud_dms_rabbitmq_exchanges" "with_arguments" {
  depends_on = [huaweicloud_dms_rabbitmq_exchange.with_arguments]

  instance_id = huaweicloud_dms_rabbitmq_instance.with_arguments.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.with_arguments.name
}

output "arguments_is_set_and_valid" {
  value = length([for v in data.huaweicloud_dms_rabbitmq_exchanges.with_arguments.exchanges : v
  if v.arguments == huaweicloud_dms_rabbitmq_exchange.with_arguments.arguments]) > 0
}
`, testDataSourceDataSourceExchanges_base(name))
}
