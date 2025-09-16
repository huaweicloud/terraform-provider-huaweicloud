package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrLongTermCredentials_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_long_term_credentials.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrLongTermCredentials_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "auth_tokens.#"),
					resource.TestCheckResourceAttrSet(dataSource, "auth_tokens.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "auth_tokens.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "auth_tokens.0.enable"),
					resource.TestCheckResourceAttrSet(dataSource, "auth_tokens.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "auth_tokens.0.user_profile"),
					resource.TestCheckResourceAttrSet(dataSource, "auth_tokens.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "auth_tokens.0.expire_date"),
				),
			},
		},
	})
}

func testDataSourceSwrLongTermCredentials_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_long_term_credentials" "test" {
  depends_on = [huaweicloud_swr_enterprise_long_term_credential.test]

  instance_id = huaweicloud_swr_enterprise_instance.test.id
}
`, testAccSwrEnterpriseLongTermCredential_update(name))
}
