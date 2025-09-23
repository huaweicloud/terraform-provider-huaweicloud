package sdrs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDeleteAllGroupFailureJobs_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceDeleteAllGroupFailureJobs_basic,
			},
		},
	})
}

const testResourceDeleteAllGroupFailureJobs_basic = `resource "huaweicloud_sdrs_delete_all_group_failure_jobs" "test" {}`
