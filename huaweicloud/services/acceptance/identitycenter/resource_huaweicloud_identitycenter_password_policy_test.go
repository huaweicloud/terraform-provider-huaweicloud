package identitycenter

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

func getIdentityCenterPasswordPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/password-policy"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{identity_store_id}", state.Primary.Attributes["identity_store_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center password policy: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccIdentityCenterPasswordPolicy_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_password_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterPasswordPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPasswordPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "max_password_age", "10"),
					resource.TestCheckResourceAttr(rName, "minimum_password_length", "10"),
					resource.TestCheckResourceAttr(rName, "password_reuse_prevention", "1"),
					resource.TestCheckResourceAttr(rName, "require_lowercase_characters", "false"),
					resource.TestCheckResourceAttr(rName, "require_numbers", "false"),
					resource.TestCheckResourceAttr(rName, "require_symbols", "true"),
					resource.TestCheckResourceAttr(rName, "require_uppercase_characters", "true"),
				),
			},
			{
				Config: testPasswordPolicy_basic_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "max_password_age", "100"),
					resource.TestCheckResourceAttr(rName, "minimum_password_length", "100"),
					resource.TestCheckResourceAttr(rName, "password_reuse_prevention", "1"),
					resource.TestCheckResourceAttr(rName, "require_lowercase_characters", "true"),
					resource.TestCheckResourceAttr(rName, "require_numbers", "true"),
					resource.TestCheckResourceAttr(rName, "require_symbols", "false"),
					resource.TestCheckResourceAttr(rName, "require_uppercase_characters", "false"),
				),
			},
		},
	})
}

const testPasswordPolicy_basic = `
data "huaweicloud_identitycenter_instance" "test" {}

resource "huaweicloud_identitycenter_password_policy" "test" {
  identity_store_id            = data.huaweicloud_identitycenter_instance.test.identity_store_id
  max_password_age             = 10
  minimum_password_length      = 10
  password_reuse_prevention    = 1
  require_lowercase_characters = false
  require_numbers              = false
  require_symbols              = true
  require_uppercase_characters = true
}
`

const testPasswordPolicy_basic_update = `
data "huaweicloud_identitycenter_instance" "test" {}

resource "huaweicloud_identitycenter_password_policy" "test" {
  identity_store_id            = data.huaweicloud_identitycenter_instance.test.identity_store_id
  max_password_age             = 100
  minimum_password_length      = 100
  password_reuse_prevention    = 1
  require_lowercase_characters = true
  require_numbers              = true
  require_symbols              = false
  require_uppercase_characters = false
}
`
