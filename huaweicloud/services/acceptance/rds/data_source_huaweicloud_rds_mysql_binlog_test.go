package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccMysqlBinlogDataSource_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_mysql_binlog.test"
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5),
		acctest.RandStringFromCharSet(2, "!#%^*"), acctest.RandIntRange(10, 99))
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBinlogDataSource_basic(name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "binlog_retention_hours"),
					resource.TestCheckResourceAttrPair(rName, "binlog_retention_hours",
						"huaweicloud_rds_mysql_binlog.test", "binlog_retention_hours"),
				),
			},
		},
	})
}

func testAccMysqlBinlogDataSource_basic(name string, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_mysql_binlog" "test" {
  depends_on  = [huaweicloud_rds_mysql_binlog.test]
  instance_id = huaweicloud_rds_mysql_binlog.test.instance_id
}

`, testMysqlBinlog_basic(name, dbPwd))
}
