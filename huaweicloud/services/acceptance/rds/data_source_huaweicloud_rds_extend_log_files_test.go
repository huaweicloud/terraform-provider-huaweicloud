package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsExtendLogFiles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_extend_log_files.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsExtendLogFiles_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "files.#"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.file_size"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsExtendLogFiles_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_extend_log_files" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testAccRdsInstance_sqlserver(name))
}
