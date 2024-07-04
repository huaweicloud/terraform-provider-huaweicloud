package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceStorageType_basic(t *testing.T) {
	rName := "data.huaweicloud_rds_storage_types.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStoragetype_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.name"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.az_status.%"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.support_compute_group_type.#"),
				),
			},
		},
	})
}

func testAccDatasourceStoragetype_basic() string {
	return `
data "huaweicloud_rds_storage_types" "test" {
  db_type    = "MySQL"
  db_version = "8.0"
}`
}
