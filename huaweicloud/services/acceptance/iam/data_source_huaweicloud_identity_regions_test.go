package iam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataRegions_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identity_regions.test"

		dc = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRegions,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "regions.#", "26"),
				),
			},
			{
				Config: testAccDataRegionsWithRegionId1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "regions.#", "1"),
					resource.TestCheckResourceAttr(dcName, "regions.0.id", "cn-north-4"),
					resource.TestCheckResourceAttr(dcName, "regions.0.type", "public"),
					resource.TestCheckResourceAttr(dcName, "regions.0.description", ""),
					resource.TestMatchResourceAttr(dcName, "regions.0.link",
						regexp.MustCompile("https://iam.*.myhuaweicloud.com/v3/regions/cn-north-4")),
					resource.TestMatchResourceAttr(dcName, "regions.0.locales.%", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.zh-cn", "华北-北京四"),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.en-us", "CN North-Beijing4"),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.pt-br", ""),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.es-us", ""),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.es-es", ""),
				),
			},
			{
				Config: testAccDataRegionsWithRegionId2,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "regions.#", "1"),
					resource.TestCheckResourceAttr(dcName, "regions.0.id", "la-south-2"),
					resource.TestCheckResourceAttr(dcName, "regions.0.type", "public"),
					resource.TestCheckResourceAttr(dcName, "regions.0.description", ""),
					resource.TestMatchResourceAttr(dcName, "regions.0.link",
						regexp.MustCompile("https://iam.*.myhuaweicloud.com/v3/regions/la-south-2")),
					resource.TestMatchResourceAttr(dcName, "regions.0.locales.%", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.en-us", "LA-Santiago"),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.pt-br", "AL-Santiago"),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.es-us", "AL-Santiago de Chile1"),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.es-es", "LA-Santiago"),
					resource.TestCheckResourceAttr(dcName, "regions.0.locales.zh-cn", "拉美-圣地亚哥"),
				),
			},
		},
	})
}

const testAccDataRegions = `
data "huaweicloud_identity_regions" "test" {}
`

const testAccDataRegionsWithRegionId1 = `
data "huaweicloud_identity_regions" "test" {
  region_id = "cn-north-4"
}
`

const testAccDataRegionsWithRegionId2 = `
data "huaweicloud_identity_regions" "test" {
  region_id = "la-south-2"
}
`
