package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUserProjects_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_user_projects.test"
		dc  = acceptance.InitDataSourceCheck(all)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserProjects_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "projects.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "projects.0.id"),
					resource.TestCheckResourceAttrSet(all, "projects.0.name"),
					resource.TestCheckResourceAttrSet(all, "projects.0.enabled"),
				),
			},
		},
	})
}

func testAccIdentityUserProjects_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_user" "test" {
  name     = "%[1]s"
  password = random_string.test.result
}

resource "huaweicloud_identity_group" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identity_role" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"
  type        = "AX"
  policy      = <<EOT
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ]
    }
  ]
}
EOT
}

data "huaweicloud_identity_projects" "test" {
  name = "%[2]s"
}

resource "huaweicloud_identity_group_role_assignment" "test" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = huaweicloud_identity_role.test.id
  project_id = try(data.huaweicloud_identity_projects.test.projects[0].id, "NOT_FOUND")
}

resource "huaweicloud_identity_group_membership" "test" {
  depends_on = [huaweicloud_identity_group_role_assignment.test]

  group = huaweicloud_identity_group.test.id
  users = [huaweicloud_identity_user.test.id]
}
`, name, acceptance.HW_REGION_NAME)
}

func testAccIdentityUserProjects_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_user_projects" "without_user_id" {
  depends_on = [huaweicloud_identity_group_membership.test]
}

data "huaweicloud_identity_user_projects" "with_user_id" {
  depends_on = [huaweicloud_identity_group_membership.test]

  user_id = huaweicloud_identity_user.test.id
}
`, testAccIdentityUserProjects_base(name))
}
