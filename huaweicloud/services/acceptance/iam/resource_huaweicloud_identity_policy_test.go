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

func getIdentityPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccIdentityPolicy_basic(t *testing.T) {
	var object interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getIdentityPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckIAMV5(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test for terraform"),
					resource.TestCheckResourceAttr(resourceName, "policy_type", "custom"),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "v1"),
					resource.TestCheckResourceAttr(resourceName, "version_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment_count"),
					resource.TestCheckResourceAttrSet(resourceName, "version_ids"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccIdentityPolicy_update(rName),
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

func testAccIdentityPolicy_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_policy" "test" {
  name            = "%s"
  description     = "test for terraform"
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
`, rName)
}

func testAccIdentityPolicy_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_policy" "test" {
  name            = "%s"
  description     = "test for terraform"
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
`, rName)
}
