package ucs

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getPolicy: Query the UCS Policy detail
	var (
		region           = acceptance.HW_REGION_NAME
		getPolicyHttpUrl = "v1/permissions/rules"
		getPolicyProduct = "ucs"
	)
	getPolicyClient, err := cfg.NewServiceClient(getPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating UCS Client: %s", err)
	}

	getPolicyPath := getPolicyClient.Endpoint + getPolicyHttpUrl

	getPolicyResp, err := pagination.ListAllItems(
		getPolicyClient,
		"offset",
		getPolicyPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Policy: %s", err)
	}

	getPolicyRespJson, err := json.Marshal(getPolicyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Policy: %s", err)
	}
	var getPolicyRespBody interface{}
	err = json.Unmarshal(getPolicyRespJson, &getPolicyRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Policy: %s", err)
	}

	jsonPath := fmt.Sprintf("items[?metadata.uid=='%s']|[0]", state.Primary.ID)
	getPolicyRespBody = utils.PathSearch(jsonPath, getPolicyRespBody, nil)
	if getPolicyRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getPolicyRespBody, nil
}

func TestAccPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_ucs_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "admin"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckOutput("is_iam_user_ids_different", "false"),
				),
			},
			{
				Config: testPolicy_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "custom"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(rName, "details.0.operations.0", "*"),
					resource.TestCheckResourceAttr(rName, "details.0.resources.0", "*"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckOutput("is_iam_user_ids_different", "false"),
				),
			},
			{
				Config:            testPolicy_basic_import(name),
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testPolicy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  count = 1

  name        = "%[1]s-${count.index}"
  password    = "Test@12345678"
}
    
resource "huaweicloud_ucs_policy" "test" {
  name         = "%[1]s"
  iam_user_ids = huaweicloud_identity_user.test[*].id
  type         = "admin"
  description  = "created by terraform"
}

output "is_iam_user_ids_different" {
  value = length(setsubtract(huaweicloud_ucs_policy.test.iam_user_ids,
    huaweicloud_identity_user.test[*].id)) != 0
}
`, name)
}

func testPolicy_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  count = 2

  name        = "%[1]s-${count.index}"
  password    = "Test@12345678"
}

resource "huaweicloud_ucs_policy" "test" {
  name         = "%[1]s"
  iam_user_ids = huaweicloud_identity_user.test[*].id
  type         = "custom"
  description  = "created by terraform update"
  details {
    operations = ["*"]
    resources  = ["*"]
  }
}

output "is_iam_user_ids_different" {
  value = length(setsubtract(huaweicloud_ucs_policy.test.iam_user_ids,
    huaweicloud_identity_user.test[*].id)) != 0
}
`, name)
}

func testPolicy_basic_import(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  count = 2

  name        = "%[1]s-${count.index}"
  password    = "Test@12345678"
}

resource "huaweicloud_ucs_policy" "test" {
	name         = "%[1]s"
	iam_user_ids = huaweicloud_identity_user.test[*].id
	type         = "custom"
	description  = "created by terraform update"
	details {
	  operations = ["*"]
	  resources  = ["*"]
	}
  }
`, name)
}
