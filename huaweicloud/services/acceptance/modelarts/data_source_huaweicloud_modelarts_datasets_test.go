package modelarts

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDatasets_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_modelarts_datasets.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	name := acceptance.RandomAccResourceName()
	obsName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatasets_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSourceName, "datasets.0.id",
						"huaweicloud_modelarts_dataset.test", "id"),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.type", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "datasets.0.description",
						"huaweicloud_modelarts_dataset.test", "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "datasets.0.output_path",
						"huaweicloud_modelarts_dataset.test", "output_path"),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.data_source.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.labels.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.schemas.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceDatasets_basic(rName, obsName string) string {
	datasets := testAccDateset_basic(rName, obsName)
	return fmt.Sprintf(`
%s

data "huaweicloud_modelarts_datasets" "test" {
  name = "%s"
  type = 1

  depends_on = [
    huaweicloud_modelarts_dataset.test
  ]
}
`, datasets, rName)
}
