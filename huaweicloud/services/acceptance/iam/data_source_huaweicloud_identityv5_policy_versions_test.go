package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV5PolicyVersions_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identityv5_policy_versions.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byVersionId   = "data.huaweicloud_identityv5_policy_versions.filter_by_version_id"
		dcByVersionId = acceptance.InitDataSourceCheck(byVersionId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5PolicyVersions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "versions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'version_id' parameter.
					dcByVersionId.CheckResourceExists(),
					resource.TestCheckOutput("is_version_id_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.version_id",
						"huaweicloud_identity_policy.test", "default_version_id"),
					resource.TestCheckResourceAttrPair(byVersionId, "versions.0.document",
						"huaweicloud_identity_policy.test", "policy_document"),
					resource.TestCheckResourceAttr(byVersionId, "versions.0.is_default", "true"),
					resource.TestCheckResourceAttrSet(byVersionId, "versions.0.created_at"),
				),
			},
		},
	})
}

func testAccDataV5PolicyVersions_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
  description     = "created by terraform script"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Allow"
        }
      ]
      Version = "5.0"
    }
  )
}

# Without any filter parameter.
data "huaweicloud_identityv5_policy_versions" "test" {
  policy_id = huaweicloud_identity_policy.test.id
}

# Filter by 'version_id' parameter.
locals {
  version_id = huaweicloud_identity_policy.test.default_version_id
}

data "huaweicloud_identityv5_policy_versions" "filter_by_version_id" {
  policy_id  = huaweicloud_identity_policy.test.id
  version_id = local.version_id
}

locals {
  version_id_filter_result = [for v in data.huaweicloud_identityv5_policy_versions.filter_by_version_id.versions[*].version_id :
  v == local.version_id]
}

output "is_version_id_filter_useful" {
  value = length(local.version_id_filter_result) > 0 && alltrue(local.version_id_filter_result)
}
`, name)
}
