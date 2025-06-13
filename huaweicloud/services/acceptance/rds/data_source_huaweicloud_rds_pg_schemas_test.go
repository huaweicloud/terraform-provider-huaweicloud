package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPgSchemas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_pg_schemas.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsPgSchemas_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "database_schemas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_schemas.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource, "database_schemas.0.owner"),
				),
			},
		},
	})
}

func testDataSourceRdsPgSchemas_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_pg_schemas" "test" {
  depends_on = [huaweicloud_rds_pg_schema.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = "%[2]s"
}
`, testPgSchema_basic(name), name)
}
