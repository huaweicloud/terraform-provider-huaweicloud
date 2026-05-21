package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDatasets_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_datasets.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_modelarts_datasets.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_modelarts_datasets.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatasets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "datasets.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byName, "datasets.0.id",
						"huaweicloud_modelarts_dataset.test", "id"),
					resource.TestCheckResourceAttr(byName, "datasets.0.name", name),
					resource.TestCheckResourceAttr(byName, "datasets.0.type", "1"),
					resource.TestCheckResourceAttrPair(byName, "datasets.0.description",
						"huaweicloud_modelarts_dataset.test", "description"),
					resource.TestCheckResourceAttrPair(byName, "datasets.0.output_path",
						"huaweicloud_modelarts_dataset.test", "output_path"),
					resource.TestCheckResourceAttr(byName, "datasets.0.data_source.#", "1"),
					resource.TestCheckResourceAttr(byName, "datasets.0.labels.#", "1"),
					resource.TestCheckResourceAttr(byName, "datasets.0.schemas.#", "0"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDatasets_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true

  lifecycle {
    ignore_changes = [
      cors_rule,
    ]
  }
}

resource "huaweicloud_obs_bucket_object" "input" {
  bucket  = huaweicloud_obs_bucket.bucket.bucket
  key     = "input/t1"
  content = "some_bucket_content"
}

resource "huaweicloud_obs_bucket_object" "output" {
  bucket  = huaweicloud_obs_bucket.bucket.bucket
  key     = "output/t2"
  content = "some_bucket_content"
}

resource "huaweicloud_modelarts_dataset" "test" {
  name        = "%[1]s"
  type        = 1
  output_path = "/${huaweicloud_obs_bucket.bucket.bucket}/output/"
  description = "Created by terraform script"
  data_source {
    path = "/${huaweicloud_obs_bucket.bucket.bucket}/input/"
  }

  labels {
    name = "%[1]s"
  }

  depends_on = [
    huaweicloud_obs_bucket_object.input,
    huaweicloud_obs_bucket_object.output
  ]
}
`, name)
}

func testAccDataSourceDatasets_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all datasets without any filter
data "huaweicloud_modelarts_datasets" "all" {
  depends_on = [
    huaweicloud_modelarts_dataset.test
  ]
}

# Filter by name
locals {
  dataset_name = huaweicloud_modelarts_dataset.test.name
}

data "huaweicloud_modelarts_datasets" "filter_by_name" {
  depends_on = [
    huaweicloud_modelarts_dataset.test
  ]

  name = local.dataset_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_modelarts_datasets.filter_by_name.datasets[*].name : v == local.dataset_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by type
locals {
  dataset_type = huaweicloud_modelarts_dataset.test.type
}

data "huaweicloud_modelarts_datasets" "filter_by_type" {
  depends_on = [
    huaweicloud_modelarts_dataset.test
  ]

  type = local.dataset_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_modelarts_datasets.filter_by_type.datasets[*].type : v == local.dataset_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`, testAccDataSourceDatasets_base(name), name)
}
