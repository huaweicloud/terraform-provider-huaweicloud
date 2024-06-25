package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsSlowLogLink_basic(t *testing.T) {
	rName := "data.huaweicloud_rds_slow_log_link.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testRdsSlowLogLink_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_RDS_INSTANCE_ID),
					resource.TestCheckResourceAttrSet(rName, "file_size"),
					resource.TestCheckResourceAttrSet(rName, "file_link"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
		},
	})
}

func testRdsSlowLogLink_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_slow_log_link" "test" {
  instance_id = "%[2]s"
  file_name   = data.huaweicloud_rds_slow_log_files.test.files[0].file_name
}
`, testDataSourceDataSourceRdsSlowLogFiles_basic(), acceptance.HW_RDS_INSTANCE_ID)
}
