package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPgRoles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_pg_roles.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsPgRoles_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "roles.#"),
					resource.TestCheckResourceAttr(dataSource, "roles.#", "1"),
					resource.TestCheckResourceAttrPair(dataSource, "roles.0",
						"huaweicloud_rds_pg_account.test", "name"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsPgRoles_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "Test@12345678"
}

data "huaweicloud_rds_pg_roles" "test" {
  depends_on = [huaweicloud_rds_pg_account.test]

  instance_id = huaweicloud_rds_instance.test.id
  account     = "root"
}
`, testAccRdsInstance_basic(name), name)
}
