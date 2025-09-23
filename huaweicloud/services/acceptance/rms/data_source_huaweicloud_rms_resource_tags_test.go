package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()
	baseConfig := testDataSourceResourceTags_base(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceTags_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.values.#"),
				),
			},
			{
				Config: testDataSourceResourceTags_specifiedResource(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.values.#"),
				),
			},
		},
	})
}

func testDataSourceResourceTags_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_tags" "test" {
  resource_type = "config:policyAssignments"

  depends_on = [huaweicloud_rms_policy_assignment.test]
}
`, baseConfig)
}

func testDataSourceResourceTags_specifiedResource(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_tags" "test" {
  resource_type = "config:policyAssignments"
  resource_id   = huaweicloud_rms_policy_assignment.test.id
}
`, baseConfig)
}

func testDataSourceResourceTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_rms_policy_definitions" "test" {
  name = "regular-matching-of-names"
}
	
resource "huaweicloud_rms_policy_assignment" "test" {
  name                 = "%[1]s"
  description          = "A resource name that does not match the regular expression is considered 'non-compliant'."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")
  status               = "Disabled"
	
  policy_filter {
    region            = "%[2]s"
    resource_provider = "vpc"
    resource_type     = "vpcs"
    tag_key           = "name"
    tag_value         = "bo"
  }
	
  parameters = {
    regularExpression = jsonencode("bo-form_")
  }

  tags = {
    foo = "bar"
  }
}
`, name, acceptance.HW_REGION_NAME)
}
