package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/organizations"
)

func getDryRunPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("organizations", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	return organizations.GetDryRunPolicyById(client, state.Primary.ID)
}

func TestAccDryRunPolicy_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_organizations_dry_run_policy.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDryRunPolicyResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDryRunPolicy_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrWith(rName, "content",
						checkPolicyContent("{\"Version\":\"5.0\",\"Statement\":[{\"Effect\":\"Deny\","+
							"\"Action\":[]}]}")),
					resource.TestCheckResourceAttr(rName, "type", "service_control_policy"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttr(rName, "is_builtin", "false"),
				),
			},
			{
				Config: testAccDryRunPolicy_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttrWith(rName, "content",
						checkPolicyContent("{\"Version\":\"5.0\",\"Statement\":[{\"Sid\":\"Statement1\","+
							"\"Effect\":\"Deny\",\"Action\":[\"vpc:subnets:delete\"]}]}")),
					resource.TestCheckResourceAttr(rName, "type", "service_control_policy"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "tags.%", "1"),
					resource.TestCheckResourceAttr(rName, "tags.foo1", "bar_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDryRunPolicy_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_dry_run_policy" "test" {
  name        = "%s"
  type        = "service_control_policy"
  description = "Created by terraform script"
  content     = jsonencode({
    Version = "5.0",

    Statement = [
      {
        Effect = "Deny",
        Action = []
      }
    ]
  })

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, name)
}

func testAccDryRunPolicy_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_dry_run_policy" "test" {
  name    = "%s"
  type    = "service_control_policy"
  content = jsonencode({
    Version = "5.0",

    Statement = [
      {
        Sid = "Statement1",
        Effect = "Deny",
        Action = ["vpc:subnets:delete"]
      }
    ]
  })

  tags = {
    foo1 = "bar_update"
  }
}
`, name)
}
