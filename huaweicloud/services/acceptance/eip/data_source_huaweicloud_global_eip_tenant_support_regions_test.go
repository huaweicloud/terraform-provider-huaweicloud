package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipTenantSupportRegions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_global_eip_tenant_support_regions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGlobalEipTenantSupportRegions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.instance_type"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.access_site"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceGlobalEipTenantSupportRegions_base() string {
	return `
data "huaweicloud_global_eip_access_sites" "all" {}
`
}

func testDataSourceGlobalEipTenantSupportRegions_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_global_eip_tenant_support_regions" "test" {
  access_site = data.huaweicloud_global_eip_access_sites.all.access_sites[0].name

  fields = [
    "id",
    "instance_type",
    "region_id",
    "public_border_group",
    "access_site",
    "status",
    "created_at",
    "updated_at",
  ]
}
`, testDataSourceGlobalEipTenantSupportRegions_base())
}
