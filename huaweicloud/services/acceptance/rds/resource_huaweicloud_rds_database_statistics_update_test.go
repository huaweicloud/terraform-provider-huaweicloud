package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsDatabaseStatisticsUpdate_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsDatabaseStatisticsUpdate_basic(rName),
			},
		},
	})
}

func testAccRdsDatabaseStatisticsUpdate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_sqlserver_database" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_rds_database_statistics_update" "test" {
  depends_on = [huaweicloud_rds_sqlserver_database.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_sqlserver_database.test.name
}
`, testAccRdsInstance_sqlserver(name), name)
}
