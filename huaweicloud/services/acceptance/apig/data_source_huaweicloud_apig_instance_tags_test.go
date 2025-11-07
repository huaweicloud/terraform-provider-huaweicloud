package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataInstanceTags_basic(t *testing.T) {
	var (
		dcAll                = "data.huaweicloud_apig_instance_tags.all"
		dcFilterByInstanceId = "data.huaweicloud_apig_instance_tags.test"
		byInstanceId         = acceptance.InitDataSourceCheck(dcFilterByInstanceId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.0.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					byInstanceId.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcFilterByInstanceId, "tags.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dcFilterByInstanceId, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dcFilterByInstanceId, "tags.0.value"),
					resource.TestMatchResourceAttr(dcAll, "tags.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckOutput("is_all_tags_include_test_tags", "true"),
				),
			},
		},
	})
}

func testAccDataSourceInstanceTags_basic() string {
	return fmt.Sprintf(`
resource "random_string" "access_key" {
  length  = 5
  special = false
}

data "huaweicloud_identity_projects" "test" {
  name = "%[1]s"
}

resource "huaweicloud_tms_resource_tags" "test" {
  project_id = try([for v in data.huaweicloud_identity_projects.test.projects: v.id if v.name == "%[1]s"][0], null)

  resources {
    resource_type = "apig" # Dedicated instance
    resource_id   = "%[2]s"
  }

  tags = {
    key   = "key${random_string.access_key.result}"
    value = "value${random_string.access_key.result}"
  }
}

data "huaweicloud_apig_instance_tags" "test" {
  depends_on = [
    huaweicloud_tms_resource_tags.test
  ]

  instance_id = "%[2]s"
}

data "huaweicloud_apig_instance_tags" "all" {
  depends_on = [
    huaweicloud_tms_resource_tags.test
  ]
}

output "is_all_tags_include_test_tags" {
  value = alltrue([for v in data.huaweicloud_apig_instance_tags.test.tags:
    length([for vv in data.huaweicloud_apig_instance_tags.all.tags:
      vv.key == v.key && vv.value == v.value
    ]) > 0
  ])
}
`, acceptance.HW_REGION_NAME, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
