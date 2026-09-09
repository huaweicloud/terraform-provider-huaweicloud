package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV5PolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5PolicyById(client, state.Primary.ID)
}

func TestAccV5Policy_basic(t *testing.T) {
	var (
		object interface{}

		resourceName = "huaweicloud_identityv5_policy.test"
		rc           = acceptance.InitResourceCheck(resourceName, &object, getV5PolicyResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5Policy_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "policy_type", "custom"),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v1"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.0", "v1"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment_count"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccV5Policy_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v2"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.0", "v2"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.1", "v1"),
				),
			},
			{
				Config: testAccV5Policy_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v3"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.0", "v3"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.1", "v2"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.2", "v1"),
				),
			},
			{
				Config: testAccV5Policy_basic_step4(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v4"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.0", "v4"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.1", "v3"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.2", "v2"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.3", "v1"),
				),
			},
			{
				Config: testAccV5Policy_basic_step5(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v5"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.0", "v5"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.1", "v4"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.2", "v3"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.3", "v2"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.4", "v1"),
				),
			},
			{
				Config: testAccV5Policy_basic_step6(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v6"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.0", "v6"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.1", "v5"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.2", "v4"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.3", "v2"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.4", "v1"),
				),
			},
			{
				Config: testAccV5Policy_basic_step7(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v7"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.0", "v7"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.1", "v6"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.2", "v5"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.3", "v4"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.4", "v2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"version_to_delete",
				},
			},
		},
	})
}

func testAccV5Policy_basic_base(name, policyAction, deleteVersionId string, isPolicyAllow bool) string {
	return fmt.Sprintf(`
variable "is_policy_allow" {
  type    = bool
  default = %[3]v
}

variable "delete_version_id" {
  type    = string
  default = "%[4]s"
}

resource "huaweicloud_identityv5_policy" "test" {
  name            = "%[1]s"
  description     = "Created by terraform script"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["%[2]s"]
          Effect = var.is_policy_allow ? "Allow" : "Deny"
        },
      ]
      Version = "5.0"
    }
  )

  version_to_delete = var.delete_version_id != "" ? var.delete_version_id : null
}
`, name, policyAction, isPolicyAllow, deleteVersionId)
}

// Create a policy and with a default version.
func testAccV5Policy_basic_step1(name string) string {
	return testAccV5Policy_basic_base(name, "*", "", true)
}

// Update the policy document and supplement a new version (now we have two versions).
func testAccV5Policy_basic_step2(name string) string {
	return testAccV5Policy_basic_base(name, "eip:globalEips:listTags", "", false)
}

// Update the policy document and supplement a new version (now we have three versions).
func testAccV5Policy_basic_step3(name string) string {
	return testAccV5Policy_basic_base(name, "eip:globalEips:listJobs", "", true)
}

// Update the policy document and supplement a new version (now we have four versions).
func testAccV5Policy_basic_step4(name string) string {
	return testAccV5Policy_basic_base(name, "eip:internetBandwidths:listTags", "", false)
}

// Update the policy document and supplement a new version (now we have five versions).
func testAccV5Policy_basic_step5(name string) string {
	return testAccV5Policy_basic_base(name, "eip:globalEips:list", "", true)
}

// Update the policy document and replace a specific version (v3) (the number of version list is upper limit).
func testAccV5Policy_basic_step6(name string) string {
	return testAccV5Policy_basic_base(name, "eip:globalEips:listQuotas", "v3", false)
}

// Update the policy document and replace the earliest version (v1) (the number of version list is upper limit).
func testAccV5Policy_basic_step7(name string) string {
	return testAccV5Policy_basic_base(name, "eip:internetBandwidths:list", "", false)
}
