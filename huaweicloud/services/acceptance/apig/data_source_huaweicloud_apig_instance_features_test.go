package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceFeatures_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_instance_features.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byName   = "data.huaweicloud_apig_instance_features.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_apig_instance_features.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceFeatures_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "features.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByName.CheckResourceExists(),
					resource.TestMatchResourceAttr(byName, "features.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceInstanceFeatures_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = local.instance_id
  name        = "ratelimit"
  enabled     = true

  config = jsonencode({
    api_limits = 200
  })
}

data "huaweicloud_apig_instance_features" "test" {
  depends_on = [huaweicloud_apig_instance_feature.test]

  instance_id = local.instance_id
}

# Filter by name
locals {
  feature_name = huaweicloud_apig_instance_feature.test.name
}

data "huaweicloud_apig_instance_features" "filter_by_name" {
  depends_on = [huaweicloud_apig_instance_feature.test]

  instance_id = local.instance_id
  name        = local.feature_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_instance_features.filter_by_name.features[*].name : v == local.feature_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_apig_instance_features" "filter_by_not_found_name" {
  depends_on = [huaweicloud_apig_instance_feature.test]

  instance_id = local.instance_id
  name        = local.not_found_name
}

locals {
  not_found_name_filter_result = [
    for v in data.huaweicloud_apig_instance_features.filter_by_not_found_name.features[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.not_found_name_filter_result) == 0
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
