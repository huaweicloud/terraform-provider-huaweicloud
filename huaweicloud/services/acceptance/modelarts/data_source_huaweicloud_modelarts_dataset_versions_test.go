package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDatasetVersions_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_modelarts_dataset_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	name := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatasetVersions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.id",
						"huaweicloud_modelarts_dataset_version.test", "version_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.description",
						"huaweicloud_modelarts_dataset_version.test", "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.name",
						"huaweicloud_modelarts_dataset_version.test", "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.split_ratio",
						"huaweicloud_modelarts_dataset_version.test", "split_ratio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.status",
						"huaweicloud_modelarts_dataset_version.test", "status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.name",
						"huaweicloud_modelarts_dataset_version.test", "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.files",
						"huaweicloud_modelarts_dataset_version.test", "files"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.storage_path",
						"huaweicloud_modelarts_dataset_version.test", "storage_path"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.is_current",
						"huaweicloud_modelarts_dataset_version.test", "is_current"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.created_at",
						"huaweicloud_modelarts_dataset_version.test", "created_at"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.updated_at",
						"huaweicloud_modelarts_dataset_version.test", "updated_at"),
				),
			},
			{
				Config: testAccDataSourceDatasetVersions_name(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "versions.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceDatasetVersions_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_dataset_version" "test" {
  name       = "%[2]s"
  dataset_id = huaweicloud_modelarts_dataset.test.id
}
`, testAccDatasetVersion_base(name), name)
}

func testAccDataSourceDatasetVersions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_modelarts_dataset_versions" "test" {
  dataset_id  = huaweicloud_modelarts_dataset.test.id
  split_ratio = "0,2.9"

  depends_on = [
    huaweicloud_modelarts_dataset_version.test
  ]
}
`, testAccDataSourceDatasetVersions_base(name))
}

func testAccDataSourceDatasetVersions_name(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_modelarts_dataset_versions" "test" {
  dataset_id  = huaweicloud_modelarts_dataset.test.id
  split_ratio = "0,2.9"
  name        = "wrong_name"

  depends_on = [
    huaweicloud_modelarts_dataset_version.test
  ]
}
`, testAccDataSourceDatasetVersions_base(name))
}
