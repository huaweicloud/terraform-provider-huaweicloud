package identitycenter

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterUsers_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_identitycenter_users.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterUsers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "users.0.id"),
					resource.TestCheckResourceAttrSet(rName, "users.0.user_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.family_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.given_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.display_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.email"),
					resource.TestCheckResourceAttrSet(rName, "users.0.phone_number"),
					resource.TestCheckResourceAttrSet(rName, "users.0.user_type"),
					resource.TestCheckResourceAttrSet(rName, "users.0.title"),
					resource.TestCheckResourceAttrSet(rName, "users.0.addresses.#"),
					resource.TestCheckResourceAttrSet(rName, "users.0.addresses.0.country"),
					resource.TestCheckResourceAttrSet(rName, "users.0.addresses.0.formatted"),
					resource.TestCheckResourceAttrSet(rName, "users.0.addresses.0.locality"),
					resource.TestCheckResourceAttrSet(rName, "users.0.addresses.0.postal_code"),
					resource.TestCheckResourceAttrSet(rName, "users.0.addresses.0.region"),
					resource.TestCheckResourceAttrSet(rName, "users.0.addresses.0.street_address"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enterprise.#"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enterprise.0.cost_center"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enterprise.0.department"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enterprise.0.division"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enterprise.0.employee_number"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enterprise.0.organization"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enterprise.0.manager"),
					resource.TestMatchResourceAttr(rName, "users.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "users.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "users.0.updated_by"),
					resource.TestCheckResourceAttrSet(rName, "users.0.email_verified"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enabled"),

					resource.TestCheckOutput("user_name_filter_is_useful", "true"),
					resource.TestCheckOutput("family_name_filter_is_useful", "true"),
					resource.TestCheckOutput("given_name_filter_is_useful", "true"),
					resource.TestCheckOutput("display_name_filter_is_useful", "true"),
					resource.TestCheckOutput("email_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceIdentityCenterUsers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identitycenter_users" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  user_name         = huaweicloud_identitycenter_user.test.user_name
}

data "huaweicloud_identitycenter_users" "user_name_filter" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  user_name         = huaweicloud_identitycenter_user.test.user_name
}

data "huaweicloud_identitycenter_users" "family_name_filter" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  family_name       = huaweicloud_identitycenter_user.test.family_name
}

data "huaweicloud_identitycenter_users" "given_name_filter" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  given_name        = huaweicloud_identitycenter_user.test.given_name
}

data "huaweicloud_identitycenter_users" "display_name_filter" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  display_name      = huaweicloud_identitycenter_user.test.display_name
}

data "huaweicloud_identitycenter_users" "email_filter" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  email             = huaweicloud_identitycenter_user.test.email
}

locals {
  user_name_filter_result = [for v in data.huaweicloud_identitycenter_users.user_name_filter.users[*].user_name:
  v == data.huaweicloud_identitycenter_users.test.users.0.user_name]
  family_name_filter_result = [for v in data.huaweicloud_identitycenter_users.family_name_filter.users[*].family_name:
  v == data.huaweicloud_identitycenter_users.test.users.0.family_name]
  given_name_filter_result = [for v in data.huaweicloud_identitycenter_users.given_name_filter.users[*].given_name:
  v == data.huaweicloud_identitycenter_users.test.users.0.given_name]
  display_name_filter_result = [for v in data.huaweicloud_identitycenter_users.display_name_filter.users[*].display_name:
  v == data.huaweicloud_identitycenter_users.test.users.0.display_name]
  email_filter_filter_result = [for v in data.huaweicloud_identitycenter_users.email_filter.users[*].email:
  v == data.huaweicloud_identitycenter_users.test.users.0.email]
}

output "user_name_filter_is_useful" {
  value = alltrue(local.user_name_filter_result) && length(local.user_name_filter_result) > 0
}

output "family_name_filter_is_useful" {
  value = alltrue(local.family_name_filter_result) && length(local.family_name_filter_result) > 0
}

output "given_name_filter_is_useful" {
  value = alltrue(local.given_name_filter_result) && length(local.given_name_filter_result) > 0
}

output "display_name_filter_is_useful" {
  value = alltrue(local.display_name_filter_result) && length(local.display_name_filter_result) > 0
}

output "email_filter_is_useful" {
  value = alltrue(local.email_filter_filter_result) && length(local.email_filter_filter_result) > 0
}
`, testIdentityCenterUser_basic(name))
}
