package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsErrorLogLink_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_rds_error_log_link.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testRdsErrorLogLink_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "file_link"),
					resource.TestCheckResourceAttrSet(dataSource, "created_at"),
				),
			},
		},
	})
}

func testRdsErrorLogLink_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_error_log_link" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testAccRdsInstance_basic(name))
}
