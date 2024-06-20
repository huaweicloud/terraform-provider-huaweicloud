package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsExtendLogLinks_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_rds_extend_log_links.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsExtendLogLinks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "links.#"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.file_link"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceRdsExtendLogLinks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_extend_log_links" "test" {
  depends_on = [huaweicloud_rds_extend_log_link.test]

  instance_id = huaweicloud_rds_instance.test.id
  file_name   = data.huaweicloud_rds_extend_log_files.test.files[0].file_name
}
`, testRdsExtendLogLink_basic(name))
}
