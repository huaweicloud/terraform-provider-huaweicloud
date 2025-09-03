package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocEnterpriseProjectCollections_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_enterprise_project_collections.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocEnterpriseProjectCollections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ep_id_list.#"),
					resource.TestCheckOutput("unique_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocEnterpriseProjectCollections_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_enterprise_project_collections" "test" {
  unique_id = huaweicloud_coc_enterprise_project_collection.test.id
}

output "unique_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_enterprise_project_collections.test.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_enterprise_project_collections.test.data[*].id :
      v == huaweicloud_coc_enterprise_project_collection.test.id]
  )
}
`, testCocEnterpriseProjectCollection_basic())
}
