package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocApplicationViews_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_application_views.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocApplicationViews_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.path"),
					resource.TestCheckOutput("name_like_filter_is_useful", "true"),
					resource.TestCheckOutput("code_list_filter_is_useful", "true"),
					resource.TestCheckOutput("is_collection_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocApplicationViews_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_coc_application_views" "test" {
  depends_on = [huaweicloud_coc_application.test]
}

data "huaweicloud_coc_application_views" "name_like_filter" {
  name_like = "%[2]s"

  depends_on = [huaweicloud_coc_application.test]
}

output "name_like_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_views.name_like_filter.data) > 0
}

data "huaweicloud_coc_application_views" "code_list_filter" {
  code_list = [huaweicloud_coc_application.test.code]
}

output "code_list_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_views.code_list_filter.data) > 0
}

data "huaweicloud_coc_application_views" "is_collection_filter" {
  is_collection = true

  depends_on = [huaweicloud_coc_application.test]
}

output "is_collection_filter_is_useful" {
  value = length(data.huaweicloud_coc_application_views.is_collection_filter.data) > 0
}
`, testAccApplication_basic(rName), rName)
}
