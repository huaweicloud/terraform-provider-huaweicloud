package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocComponents_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_components.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocComponents_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.application_id"),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocComponents_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_components" "test" {
  application_id = huaweicloud_coc_application.test.id

  depends_on = [huaweicloud_coc_component.test]
}

output "application_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_components.test.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_components.test.data[*].application_id : v == huaweicloud_coc_application.test.id]
  )
}
`, testAccComponent_basic(name))
}
