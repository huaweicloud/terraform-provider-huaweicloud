package organizations

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	region := acceptance.HW_REGION_NAME
	// getPolicy: Query Organizations policy
	var (
		getPolicyHttpUrl = "v1/organizations/policies/{policy_id}"
		getPolicyProduct = "organizations"
	)
	getPolicyClient, err := cfg.NewServiceClient(getPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	getPolicyPath := getPolicyClient.Endpoint + getPolicyHttpUrl
	getPolicyPath = strings.ReplaceAll(getPolicyPath, "{policy_id}", state.Primary.ID)

	getPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPolicyResp, err := getPolicyClient.Request("GET", getPolicyPath, &getPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations policy: %s", err)
	}

	getPolicyRespBody, err := utils.FlattenResponse(getPolicyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations policy: %s", err)
	}

	return getPolicyRespBody, nil
}

func TestAccPolicy_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_organizations_policy.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getPolicyResourceFunc)
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
				Config: testAccPolicy_basic_step1(name),
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
				),
			},
			{
				Config: testAccPolicy_basic_step2(updateName),
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

func checkPolicyContent(targetContent string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		var targetJson, valueJson interface{}
		if err := json.Unmarshal([]byte(targetContent), &targetJson); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value), &valueJson); err != nil {
			return err
		}
		if reflect.DeepEqual(targetJson, valueJson) {
			return nil
		}
		return fmt.Errorf("%#v is not equal target %#v", value, targetContent)
	}
}

func testAccPolicy_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_policy" "test" {
  name        = "%s"
  type        = "service_control_policy"
  description = "Created by terraform script"
  content     = jsonencode({
    Version : "5.0",
    Statement : [
      {
        Effect : "Deny",
        Action : []
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

func testAccPolicy_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_policy" "test" {
  name    = "%s"
  type    = "service_control_policy"
  content = jsonencode({
    Version : "5.0",
    Statement : [
      {
        Sid : "Statement1",
        Effect : "Deny",
        Action : ["vpc:subnets:delete"]
      }
    ]
  })

  tags = {
    foo1 = "bar_update"
  }
}
`, name)
}
