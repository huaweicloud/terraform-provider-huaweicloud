package identitycenter

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterBatchQueryUsers_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_identitycenter_batch_query_users.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterBatchQueryUsers_basic(name),
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
					resource.TestCheckResourceAttrSet(rName, "users.0.created_by"),
					resource.TestCheckResourceAttrSet(rName, "users.0.updated_by"),
					resource.TestCheckResourceAttrSet(rName, "users.0.email_verified"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enabled"),
				),
			},
		},
	})
}

func testAccDatasourceIdentityCenterBatchQueryUsers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identitycenter_users" "test"{
  depends_on        = [huaweicloud_identitycenter_user.test]
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}

data "huaweicloud_identitycenter_batch_query_users" "test"{
  depends_on        = [huaweicloud_identitycenter_user.test]
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  user_ids          = [for user in data.huaweicloud_identitycenter_users.test.users: user.id]
}
`, testIdentityCenterUser_basic(name))
}
