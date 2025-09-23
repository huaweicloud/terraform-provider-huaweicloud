package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesMetrics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_metrics.filter_by_namespace"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesMetrics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "metrics.0.namespace", "SYS.ER"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.metric_name"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.dimensions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.unit"),
					resource.TestCheckOutput("is_namespace_filter_useful", "true"),
					resource.TestCheckOutput("is_metric_name_filter_useful", "true"),
					resource.TestCheckOutput("is_dim0_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesMetrics_basic() string {
	return fmt.Sprintf(`
%[1]s

locals {
  namespace   = "SYS.ER"
  metric_name = "instance_bits_rate_in"
  dim0        = format("er_instance_id,%%s", huaweicloud_er_instance.test.id)
}

data "huaweicloud_ces_metrics" "filter_by_namespace" {
  namespace = local.namespace

  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]
}

output "is_namespace_filter_useful" {
  value = length(data.huaweicloud_ces_metrics.filter_by_namespace.metrics) >= 1 && alltrue(
    [for v in data.huaweicloud_ces_metrics.filter_by_namespace.metrics[*].namespace : v == local.namespace]
  )
}

data "huaweicloud_ces_metrics" "filter_by_metric_name" {
  metric_name = local.metric_name

  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]
}

output "is_metric_name_filter_useful" {
  value = length(data.huaweicloud_ces_metrics.filter_by_metric_name.metrics) >= 1 && alltrue(
    [for v in data.huaweicloud_ces_metrics.filter_by_metric_name.metrics[*].metric_name : v == local.metric_name]
  )
}

data "huaweicloud_ces_metrics" "filter_by_dim0" {
  dim_0 = local.dim0

  depends_on = [
    huaweicloud_er_vpc_attachment.test,
  ]
}

output "is_dim0_filter_useful" {
  value = length(data.huaweicloud_ces_metrics.filter_by_dim0.metrics) >= 1 && alltrue([
    for v in data.huaweicloud_ces_metrics.filter_by_dim0.metrics[*].dimensions : 
      v[0].name == "er_instance_id" && v[0].value == huaweicloud_er_instance.test.id
  ])
}
`, testAccDataSourceCesMetrics_base())
}

func testAccDataSourceCesMetrics_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}
	
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name = "%[1]s"
  asn  = 64512
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name                   = "%[1]s"
  description            = "Create by acc test"
  auto_create_vpc_routes = true

  tags = {
    foo = "bar"
  }

  provisioner "local-exec" {
    command = "sleep 70"
  }
}
`, name)
}
