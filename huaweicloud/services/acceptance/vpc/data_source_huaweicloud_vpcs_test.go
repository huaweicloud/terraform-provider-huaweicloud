package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.huaweicloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_basic(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_base(rName, cidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}
`, rName, cidr)
}

func testAccDataSourceVpcs_basic(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcs" "test" {
  id = huaweicloud_vpc.test.id
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.huaweicloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byCidr(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byCidr(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcs" "test" {
  cidr = huaweicloud_vpc.test.cidr

  depends_on = [
    huaweicloud_vpc.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.huaweicloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byName(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byName(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcs" "test" {
  name = huaweicloud_vpc.test.name

  depends_on = [
    huaweicloud_vpc.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byAll(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.huaweicloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byAll(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${huaweicloud_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byAll(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcs" "test" {
  id                    = huaweicloud_vpc.test.id
  name                  = huaweicloud_vpc.test.name
  cidr                  = huaweicloud_vpc.test.cidr
  status                = "OK"

  depends_on = [
    huaweicloud_vpc.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_tags(t *testing.T) {
	randName1 := acceptance.RandomAccResourceName()
	randName2 := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.huaweicloud_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_tags(randName1, randName2, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "tags.foo", randName1),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName1),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_tags(rName1, rName2, randCidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test1" {
  name = "%[1]s"
  cidr = "%[3]s"
  tags = {
    foo = "%[1]s"
  }
}

resource "huaweicloud_vpc" "test2" {
  name = "%[2]s"
  cidr = "%[3]s"
  tags = {
    foo = "%[2]s"
  }
}

data "huaweicloud_vpcs" "test" {
  cidr = "%[3]s"
  tags = {
    foo = "%[1]s"
  }
  depends_on = [
    huaweicloud_vpc.test1,
    huaweicloud_vpc.test2,
  ]
}
`, rName1, rName2, randCidr)
}
