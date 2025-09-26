package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataTags_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dms_kafka_tags.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "tags.0.key"),
					resource.TestMatchResourceAttr(all, "tags.0.values.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("tags_validation", "true"),
				),
			},
		},
	})
}

func testAccDataTags_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type = "single"
}

locals {
  flavor = try(data.huaweicloud_dms_kafka_flavors.test.flavors[0], {})
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = "%[2]s"
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  flavor_id          = local.flavor.id
  storage_spec_code  = try(local.flavor.ios[0].storage_spec_code, null)
  engine_version     = "3.x"
  broker_num         = 1
  storage_space      = try(local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node, null)
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  tags = {
    owner = "terraform"
    foo   = "bar"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccDataTags_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_kafka_tags" "test" {
  depends_on = [huaweicloud_dms_kafka_instance.test]
}

output "tags_validation" {
  value = length([for v in data.huaweicloud_dms_kafka_tags.test.tags : v.key == "owner" &&
  alltrue([for k, v in huaweicloud_dms_kafka_instance.test[*].tags : contains(v.values, v) if k == "owner"])]) > 0
}
`, testAccDataTags_base())
}
