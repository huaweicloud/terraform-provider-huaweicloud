package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccV5PolicyDefaultVersion_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_identity_policy.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV5PolicyDefaultVersion_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "default_version_id", "v1"),
					resource.TestCheckResourceAttr(rName, "version_ids.#", "1"),
				),
			},
			{
				Config: testAccV5PolicyDefaultVersion_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "default_version_id", "v2"),
					resource.TestCheckResourceAttr(rName, "version_ids.#", "2"),
				),
			},
			{
				Config: testAccV5PolicyDefaultVersion_step3(name),
				Check: resource.ComposeTestCheckFunc(
					// After setting the default version, the resource was not refreshed, so 'default_version_id' still has the old value.
					resource.TestCheckResourceAttr(rName, "default_version_id", "v2"),
					resource.TestCheckResourceAttr(rName, "version_ids.#", "2"),
				),
			},
			{
				// Verify the default version modification was successful.
				Config: testAccV5PolicyDefaultVersion_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "default_version_id", "v1"),
					resource.TestCheckResourceAttr(rName, "version_ids.#", "2"),
				),
			},
		},
	})
}

func testAccV5PolicyDefaultVersion_step1(name string) string {
	return fmt.Sprintf(`
# Default version is v1.
resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
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
`, name)
}

func testAccV5PolicyDefaultVersion_step2(name string) string {
	return fmt.Sprintf(`
# After updating the policy, the default version is v2.
resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Deny"
        }
      ]
      Version = "5.0"
    }
  )
}
`, name)
}

func testAccV5PolicyDefaultVersion_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Deny"
        }
      ]
      Version = "5.0"
    }
  )

  # After the policy version changes, the 'policy_document' will also be updated to the corresponding version.
  lifecycle {
    ignore_changes = [policy_document]
  }
}

# Set the default version to v1.
resource "huaweicloud_identityv5_policy_default_version" "test" {
  policy_id  = huaweicloud_identity_policy.test.id
  version_id = "v1"
}
`, name)
}
