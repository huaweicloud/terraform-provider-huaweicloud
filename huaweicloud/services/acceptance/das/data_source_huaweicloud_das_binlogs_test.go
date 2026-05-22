package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataBinlogs_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_binlogs.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBinlogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "binlogs.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "binlogs.0.file_name"),
					resource.TestCheckResourceAttrSet(all, "binlogs.0.file_size"),
				),
			},
		},
	})
}

func testAccDataBinlogs_base() string {
	return fmt.Sprintf(`
data "huaweicloud_das_database_users" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_RDS_INSTANCE_ID)
}

func testAccDataBinlogs_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_binlogs" "all" {
  user_id     = try(data.huaweicloud_das_database_users.test.users.0.id, "")
  binlog_type = "latest"
}
`, testAccDataBinlogs_base())
}
