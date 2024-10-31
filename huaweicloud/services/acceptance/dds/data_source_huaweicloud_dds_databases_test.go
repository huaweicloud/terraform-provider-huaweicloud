package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_databases.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsDatabases_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.data_size"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.storage_size"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.collection_num"),
				),
			},
		},
	})
}

func testDataSourceDdsDatabases_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dds_databases" "test" {
  instance_id = huaweicloud_dds_instance.instance.id
}
`, testAccDDSInstanceReplicaSetBasic(name))
}
