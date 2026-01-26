package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam_no_version", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	getPolicyHttpUrl := "v5/policies/{policy_id}"
	getPolicyPath := client.Endpoint + getPolicyHttpUrl
	getPolicyPath = strings.ReplaceAll(getPolicyPath, "{policy_id}", state.Primary.ID)
	getPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPolicyResp, err := client.Request("GET", getPolicyPath, &getPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM identity policy: %s", err)
	}
	return utils.FlattenResponse(getPolicyResp)
}

func TestAccPolicy_basic(t *testing.T) {
	var (
		object interface{}

		resourceName = "huaweicloud_identity_policy.test"
		rc           = acceptance.InitResourceCheck(resourceName, &object, getPolicyResourceFunc)

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
				Config: testAccPolicy_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "policy_type", "custom"),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v1"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment_count"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccPolicy_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v2"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPolicy_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
  description     = "Created by terraform script"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Allow"
        },
      ]
      Version = "5.0"
    }
  )
}
`, name)
}

func testAccPolicy_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
  description     = "Created by terraform script"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Deny"
        },
      ]
      Version = "5.0"
    }
  )
}
`, name)
}
