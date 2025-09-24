package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrEnterpriseJobDelete_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseJobDelete_basic(rName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccSwrEnterpriseJobDelete_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_jobs" "test" {
  depends_on = [huaweicloud_swr_enterprise_instance.test]
}

resource "huaweicloud_swr_enterprise_job_delete" "test" {
  job_id = try(data.huaweicloud_swr_enterprise_jobs.test.jobs[0].id, "")
}
`, testAccSwrEnterpriseInstance_update(rName))
}
