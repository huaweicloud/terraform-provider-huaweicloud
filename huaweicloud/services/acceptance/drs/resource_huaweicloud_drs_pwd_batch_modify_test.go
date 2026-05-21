package drs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourcePwdBatchModify_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_pwd_batch_modify.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
			acceptance.TestAccPreCheckDrsDbPassword(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePwdBatchModify_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "id",
						regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestMatchResourceAttr(resourceName, "results.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.id"),
					resource.TestCheckResourceAttr(resourceName, "results.0.id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.end_point_type"),
					resource.TestCheckResourceAttr(resourceName, "results.0.end_point_type", "so"),
					resource.TestCheckResourceAttr(resourceName, "jobs.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "jobs.0.job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(resourceName, "jobs.0.db_password", acceptance.HW_DRS_DB_PASSWORD),
					resource.TestCheckResourceAttr(resourceName, "jobs.0.end_point_type", "so"),
					resource.TestCheckResourceAttr(resourceName, "jobs.0.kerberos.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "jobs.1.job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(resourceName, "jobs.1.db_password", acceptance.HW_DRS_DB_PASSWORD),
					resource.TestCheckResourceAttr(resourceName, "jobs.1.end_point_type", "ta"),
					resource.TestCheckResourceAttr(resourceName, "jobs.1.kerberos.#", "0"),
				),
			},
		},
	})
}

func testAccResourcePwdBatchModify_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_pwd_batch_modify" "test" {
  jobs {
    job_id         = "%[1]s"
    db_password    = "%[2]s"
    end_point_type = "so"
  }

  jobs {
    job_id         = "%[1]s"
    db_password    = "%[2]s"
    end_point_type = "ta"
  }
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_DB_PASSWORD)
}
