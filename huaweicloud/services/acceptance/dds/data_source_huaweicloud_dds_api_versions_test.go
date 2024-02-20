package dds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsAPIVersions_basic(t *testing.T) {
	rName := "data.huaweicloud_dds_api_versions.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdsAPIVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "versions.#"),
					resource.TestCheckResourceAttrSet(rName, "versions.0.id"),
					resource.TestCheckResourceAttrSet(rName, "versions.0.status"),
					resource.TestCheckResourceAttrSet(rName, "versions.0.updated"),
				),
			},
		},
	})
}

func testAccDatasourceDdsAPIVersions_basic() string {
	return `

data "huaweicloud_dds_api_versions" "test" {}
`
}
