package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalInternetBandwidthsByTags_basic(t *testing.T) {
	//test
	dataSource := "data.huaweicloud_global_internet_bandwidths_by_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGlobalInternetBandwidthsByTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "request_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(dataSource, "resources.0.tags.0.key", "tag1"),
					resource.TestCheckResourceAttr(dataSource, "resources.0.tags.0.value", "tag1"),
				),
			},
		},
	})
}

func TestAccDataSourceGlobalInternetBandwidthsByTags_withoutTags(t *testing.T) {
	dataSource := "data.huaweicloud_global_internet_bandwidths_by_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGlobalInternetBandwidthsByTags_withoutTags,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "request_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
				),
			},
		},
	})
}

func TestAccDataSourceGlobalInternetBandwidthsByTags_withMultipleTags(t *testing.T) {
	dataSource := "data.huaweicloud_global_internet_bandwidths_by_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGlobalInternetBandwidthsByTags_withMultipleTags(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "request_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttr(dataSource, "resources.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceGlobalInternetBandwidthsByTags_withMultipleTags(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_global_internet_bandwidths_by_tags" "test" {
  depends_on = [huaweicloud_global_internet_bandwidth.test]

  tags {
    key   = "foo"
    value = "bar"
  }

  tags {
    key   = "env"
    value = "test"
  }
}
`, testAccInternetBandwidth_basic(name))
}

const testAccDataSourceGlobalInternetBandwidthsByTags_withoutTags = `
data "huaweicloud_global_internet_bandwidths_by_tags" "test" {
}
`

func testAccDataSourceGlobalInternetBandwidthsByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_global_internet_bandwidths_by_tags" "test" {
  depends_on = [huaweicloud_global_internet_bandwidth.test]
  tags {
    key   = "tag1"
    value = "tag1"
  }
}
`, testAccInternetBandwidth_basic(name))
}
