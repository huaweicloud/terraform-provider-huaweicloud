package cce

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

func getAccessPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getAccessPolicyHttpUrl = "api/v3/access-policies/{policy_id}"
		getAccessPolicyProduct = "cce"
	)
	getAccessPolicyClient, err := cfg.NewServiceClient(getAccessPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE Client: %s", err)
	}

	getAccessPolicyPath := getAccessPolicyClient.Endpoint + getAccessPolicyHttpUrl
	getAccessPolicyPath = strings.ReplaceAll(getAccessPolicyPath, "{policy_id}", state.Primary.ID)

	getAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAccessPolicyResp, err := getAccessPolicyClient.Request("GET", getAccessPolicyPath, &getAccessPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE access policy: %s", err)
	}

	return utils.FlattenResponse(getAccessPolicyResp)
}

func TestAccAccessPolicy_basic(t *testing.T) {
	var (
		accessPolicy interface{}
		resourceName = "huaweicloud_cce_access_policy.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		updateName   = rName + "-update"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&accessPolicy,
			getAccessPolicyFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "clusters.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "access_scope.0.namespaces.0", "default"),
					resource.TestCheckResourceAttr(resourceName, "policy_type", "CCEClusterAdminPolicy"),
					resource.TestCheckResourceAttr(resourceName, "principal.0.type", "user"),
					resource.TestCheckResourceAttr(resourceName, "principal.0.ids.0", acceptance.HW_USER_ID),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessPolicy_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
				),
			},
		},
	})
}

func testAccAccessPolicy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_access_policy" "test" {
  name     = "%[1]s"
  clusters = ["*"]

  access_scope {
    namespaces = ["default"]
  }

  policy_type = "CCEClusterAdminPolicy"

  principal {
    type = "user"
    ids  = ["%[2]s"]
  }
}
`, name, acceptance.HW_USER_ID)
}
