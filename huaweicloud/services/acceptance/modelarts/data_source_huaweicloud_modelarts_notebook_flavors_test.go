package modelarts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceNotebookFlavors_basic(t *testing.T) {
	rName := "data.huaweicloud_modelarts_notebook_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceNotebookFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.arch"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.category"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.description"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.feature"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.memory"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.free"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.sold_out"),

					resource.TestCheckOutput("category_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceNotebookFlavors_basic() string {
	return `
data "huaweicloud_modelarts_notebook_flavors" "test" {
}

data "huaweicloud_modelarts_notebook_flavors" "category_filter" {
  category = "CPU"
}
output "category_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_modelarts_notebook_flavors.category_filter.flavors[*].category : v == "CPU"])
}
`
}
