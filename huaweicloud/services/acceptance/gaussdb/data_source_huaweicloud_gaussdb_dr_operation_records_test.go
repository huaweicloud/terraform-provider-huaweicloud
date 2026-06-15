package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbDrOperationRecords_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_dr_operation_records.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbDrOperationRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "records.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.action"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.entity_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.entity_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.updated_at"),
				),
			},
		},
	})
}

func testAccGaussDbDrOperationRecords_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_dr_relationships" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_dr_operation_records" "test" {
  instance_id = "%[1]s"
  entity_id   = data.huaweicloud_gaussdb_dr_relationships.test.relations[0].synchronization_id
  entity_type = "dr"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
