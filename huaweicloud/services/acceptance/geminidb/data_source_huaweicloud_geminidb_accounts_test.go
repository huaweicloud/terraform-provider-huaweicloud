package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminidbAccountsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_geminidb_accounts.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminidbAccountsDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.privilege"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.databases.#"),
				),
			},
		},
	})
}

func testAccGeminidbAccountsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_geminidb_accounts" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
}
`, testAccGeminiDbInstance_basic(rName))
}
