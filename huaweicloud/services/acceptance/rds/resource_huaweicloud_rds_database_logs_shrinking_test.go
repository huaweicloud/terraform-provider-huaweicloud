package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsDbLogsShrinking_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsDbLogsShrinking_basic(name),
			},
		},
	})
}

func testAccRdsDbLogsShrinking_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_database_logs_shrinking" "test" {
  depends_on = [huaweicloud_rds_sqlserver_database.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_sqlserver_database.test.name
}
`, testSQLServerDatabase_basic(name))
}
