package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTemporaryAccessKey_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_identity_temporary_access_key.test"
		rc           = acceptance.InitDataSourceCheck(resourceName)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTemporaryAccessKey_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
					resource.TestCheckResourceAttrSet(resourceName, "access"),
					resource.TestCheckResourceAttrSet(resourceName, "secret"),
					resource.TestCheckResourceAttrSet(resourceName, "securitytoken"),
				),
			},
			{
				Config: testAccTemporaryAccessKey_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
					resource.TestCheckResourceAttrSet(resourceName, "access"),
					resource.TestCheckResourceAttrSet(resourceName, "secret"),
					resource.TestCheckResourceAttrSet(resourceName, "securitytoken"),
				),
			},
		},
	})
}

func testAccTemporaryAccessKey_basic_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

data "huaweicloud_identity_group" "test" {
  name = "admin"
}

resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = random_string.test.result
  enabled     = true
  description = "Created by terraform script"
}

resource "huaweicloud_identity_group_membership" "test" {
  group = data.huaweicloud_identity_group.test.id
  users = [huaweicloud_identity_user.test.id]
}

resource "huaweicloud_identity_user_token" "test" {
  depends_on = [huaweicloud_identity_group_membership.test]

  account_name = "%[2]s"
  user_name    = huaweicloud_identity_user.test.name
  password     = random_string.test.result
}
`, name, acceptance.HW_DOMAIN_NAME)
}

func testAccTemporaryAccessKey_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_agency" "test" {
  name                  = "%[2]s"
  description           = "Created by terraform script"
  delegated_domain_name = "%[3]s"
}

resource "huaweicloud_identity_temporary_access_key" "test" {
  token       = huaweicloud_identity_user_token.test.token
  methods     = "assume_role"
  agency_name = huaweicloud_identity_agency.test.name
  domain_name = huaweicloud_identity_agency.test.delegated_domain_name
  policy      = jsonencode({
	"Version": "1.1",
	"Statement": [
	  {
	    "Effect": "Allow",
	    "Action": [
	      "obs:object:GetObject"
	    ],
	    "Resource": [
	      "OBS:*:*:object:*"
	    ]
	  }
	]
  })
}
`, testAccTemporaryAccessKey_basic_base(name), name, acceptance.HW_DOMAIN_NAME)
}

func testAccTemporaryAccessKey_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_temporary_access_key" "test" {
  token            = huaweicloud_identity_user_token.test.token
  methods          = "token"
  duration_seconds = 3600
}
`, testAccTemporaryAccessKey_basic_base(name))
}
