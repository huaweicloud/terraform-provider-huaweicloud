package iam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityRegions_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_regions.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTestDataSourceIdentityRegions,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "regions.#", "26"),
				),
			},
			{
				Config: testTestDataSourceIdentityRegionsWithRegionId1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "regions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.id", "cn-north-4"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.type", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.description", ""),
					resource.TestMatchResourceAttr(dataSourceName, "regions.0.link",
						regexp.MustCompile("https://iam.*.myhuaweicloud.com/v3/regions/cn-north-4")),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.%", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.zh-cn", "华北-北京四"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.en-us", "CN North-Beijing4"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.pt-br", ""),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.es-us", ""),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.es-es", ""),
				),
			},
			{
				Config: testTestDataSourceIdentityRegionsWithRegionId2,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "regions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.id", "la-south-2"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.type", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.description", ""),
					resource.TestMatchResourceAttr(dataSourceName, "regions.0.link",
						regexp.MustCompile("https://iam.*.myhuaweicloud.com/v3/regions/la-south-2")),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.%", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.en-us", "LA-Santiago"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.pt-br", "AL-Santiago"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.es-us", "AL-Santiago de Chile1"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.es-es", "LA-Santiago"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.locales.zh-cn", "拉美-圣地亚哥"),
				),
			},
		},
	})
}

const testTestDataSourceIdentityRegions = `
data "huaweicloud_identity_regions" "test" {}
`

const testTestDataSourceIdentityRegionsWithRegionId1 = `
data "huaweicloud_identity_regions" "test" {
  region_id = "cn-north-4"
}
`

const testTestDataSourceIdentityRegionsWithRegionId2 = `
data "huaweicloud_identity_regions" "test" {
  region_id = "la-south-2"
}
`
