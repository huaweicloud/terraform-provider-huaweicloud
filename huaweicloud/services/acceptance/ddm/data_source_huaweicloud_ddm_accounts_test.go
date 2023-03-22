package ddm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmAccounts_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	instanceName := strings.ReplaceAll(name, "_", "-")
	rName := "data.huaweicloud_ddm_accounts.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmAccounts_basic(instanceName, name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "accounts.#", "1"),
					resource.TestCheckResourceAttr(rName, "accounts.0.name", name),
					resource.TestCheckResourceAttr(rName, "accounts.0.status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "accounts.0.permissions.#", "1"),
					resource.TestCheckResourceAttr(rName, "accounts.0.permissions.0", "SELECT"),
					resource.TestCheckResourceAttr(rName, "accounts.0.schemas.#", "0"),
				),
			},
		},
	})
}

func testAccDatasourceDdmAccounts_basic(instanceName, name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ddm_accounts" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
  name        = huaweicloud_ddm_account.test.name
}
`, testDdmAccount_basic(instanceName, name))
}
