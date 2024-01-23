package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRouteTablesDataSource_basic(t *testing.T) {
	var (
		dName    = "data.huaweicloud_er_route_tables.test"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesDataSource_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dName, "route_tables.#", "2"),
				),
			},
		},
	})
}

func TestAccRouteTablesDataSource_byName(t *testing.T) {
	var (
		dName    = "data.huaweicloud_er_route_tables.test"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesDataSource_byName(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("route_tables_count", "1"),
				),
			},
		},
	})
}

func TestAccRouteTablesDataSource_byId(t *testing.T) {
	var (
		dName    = "data.huaweicloud_er_route_tables.test"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesDataSource_byId(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("route_tables_count", "1"),
				),
			},
		},
	})
}

func TestAccRouteTablesDataSource_byTags(t *testing.T) {
	var (
		dName    = "data.huaweicloud_er_route_tables.test"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesDataSource_byTags(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("route_tables_count", "2"),
				),
			},
		},
	})
}

func testAccRouteTablesDataSource_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  name = "%[1]s"
  asn  = %[2]d
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[1]s"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}

resource "huaweicloud_er_route_table" "another" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[1]s_another"

  tags = {
    owner = "terraform"
  }
}
`, name, bgpAsNum)
}

func testAccRouteTablesDataSource_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_route_tables" "test" {
  depends_on = [
    huaweicloud_er_route_table.test
  ]

  instance_id = huaweicloud_er_instance.test.id
}
`, testAccRouteTablesDataSource_base(name, bgpAsNum), name)
}

func testAccRouteTablesDataSource_byName(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_route_tables" "test" {
  depends_on = [
    huaweicloud_er_route_table.test
  ]

  instance_id = huaweicloud_er_instance.test.id
  name        = "%[2]s"
}

output "route_tables_count" {
  value = length(data.huaweicloud_er_route_tables.test.route_tables)
}
`, testAccRouteTablesDataSource_base(name, bgpAsNum), name)
}

func testAccRouteTablesDataSource_byId(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_route_tables" "test" {
  depends_on = [
    huaweicloud_er_route_table.test
  ]

  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
}

output "route_tables_count" {
  value = length(data.huaweicloud_er_route_tables.test.route_tables)
}
`, testAccRouteTablesDataSource_base(name, bgpAsNum), name)
}

func testAccRouteTablesDataSource_byTags(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_route_tables" "test" {
  depends_on = [
    huaweicloud_er_route_table.test
  ]

  instance_id = huaweicloud_er_instance.test.id

  tags = {
    owner = "terraform"
  }
}

output "route_tables_count" {
  value = length(data.huaweicloud_er_route_tables.test.route_tables)
}
`, testAccRouteTablesDataSource_base(name, bgpAsNum), name)
}
