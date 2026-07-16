package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceRouteTables_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_er_route_tables.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byRouteTableId   = "data.huaweicloud_er_route_tables.filter_by_route_table_id"
		dcByRouteTableId = acceptance.InitDataSourceCheck(byRouteTableId)

		byName   = "data.huaweicloud_er_route_tables.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byTags   = "data.huaweicloud_er_route_tables.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteTables_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "route_tables.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Check whether filter parameter 'route_table_id' is effective.
					dcByRouteTableId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.id",
						"huaweicloud_er_route_table.test", "id"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.name",
						"huaweicloud_er_route_table.test", "name"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.description",
						"huaweicloud_er_route_table.test", "description"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.associations.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.associations.0.id",
						"huaweicloud_er_association.test", "id"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.associations.0.attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.associations.0.attachment_type", "vpc"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.propagations.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.propagations.0.id",
						"huaweicloud_er_propagation.test", "id"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.propagations.0.attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.propagations.0.attachment_type", "vpc"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.routes.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.routes.0.id",
						"huaweicloud_er_static_route.test", "id"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.routes.0.destination",
						"huaweicloud_er_static_route.test", "destination"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.routes.0.is_blackhole", "false"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.routes.0.attachments.#", "1"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.routes.0.attachments.0.attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.routes.0.attachments.0.attachment_type", "vpc"),
					resource.TestCheckResourceAttrPair(byRouteTableId, "route_tables.0.routes.0.attachments.0.resource_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.is_default_association", "false"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.is_default_propagation", "false"),
					resource.TestCheckResourceAttrSet(byRouteTableId, "route_tables.0.routes.0.status"),
					resource.TestCheckResourceAttrSet(byRouteTableId, "route_tables.0.status"),
					resource.TestMatchResourceAttr(byRouteTableId, "route_tables.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byRouteTableId, "route_tables.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.tags.%", "2"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(byRouteTableId, "route_tables.0.tags.key", "value"),
					resource.TestCheckOutput("is_route_table_id_filter_useful", "true"),
					// Check whether filter parameter 'name' is effective.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Check whether filter parameter 'tags' is effective.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRouteTables_base() string {
	var (
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

%[1]s

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name = "%[2]s"
  asn  = %[3]d
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[2]s"
  description = "Created by script"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  name        = "%[2]s"
}

resource "huaweicloud_er_static_route" "test" {
  route_table_id = huaweicloud_er_route_table.test.id
  destination    = huaweicloud_vpc.test.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}

resource "huaweicloud_er_association" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}

resource "huaweicloud_er_propagation" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}
`, common.TestVpc(name), name, bgpAsNum)
}

func testAccDataSourceRouteTables_basic_step1() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_route_tables" "test" {
  depends_on = [
    huaweicloud_er_static_route.test,
    huaweicloud_er_association.test,
    huaweicloud_er_propagation.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
}

# Filter by route table ID
locals {
  route_table_id = huaweicloud_er_route_table.test.id
}

data "huaweicloud_er_route_tables" "filter_by_route_table_id" {
  depends_on = [
    huaweicloud_er_static_route.test,
    huaweicloud_er_association.test,
    huaweicloud_er_propagation.test,
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = local.route_table_id
}

locals {
  route_table_id_filter_result = [
    for v in data.huaweicloud_er_route_tables.filter_by_route_table_id.route_tables[*].id : v == local.route_table_id
  ]
}

output "is_route_table_id_filter_useful" {
  value = length(local.route_table_id_filter_result) > 0 && alltrue(local.route_table_id_filter_result)
}

# Filter by name
locals {
  route_table_name = huaweicloud_er_route_table.test.name
}

data "huaweicloud_er_route_tables" "filter_by_name" {
  depends_on = [
    huaweicloud_er_static_route.test,
    huaweicloud_er_association.test,
    huaweicloud_er_propagation.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
  name        = local.route_table_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_er_route_tables.filter_by_route_table_id.route_tables[*].name : v == local.route_table_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by tags
locals {
  route_table_tags = huaweicloud_er_route_table.test.tags
}

data "huaweicloud_er_route_tables" "filter_by_tags" {
  depends_on = [
    huaweicloud_er_static_route.test,
    huaweicloud_er_association.test,
    huaweicloud_er_propagation.test,
  ]

  instance_id = huaweicloud_er_instance.test.id
  tags        = local.route_table_tags
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_er_route_tables.filter_by_tags.route_tables[*].tags : length(v) == length(local.route_table_tags) &&
    length(v) == length(merge(v, local.route_table_tags))
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testAccDataSourceRouteTables_base())
}
