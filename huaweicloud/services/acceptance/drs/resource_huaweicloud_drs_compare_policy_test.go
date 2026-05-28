package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/drs"
)

func getComparePolicyFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("drs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DRS client: %s", err)
	}

	return drs.GetComparePolicy(client, state.Primary.ID)
}

func TestAccResourceComparePolicy_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_drs_compare_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getComparePolicyFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComparePolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(resourceName, "period", "* * 1,3,5"),
					resource.TestCheckResourceAttr(resourceName, "begin_time", "00:00:00"),
					resource.TestCheckResourceAttr(resourceName, "end_time", "04:00:00"),
					resource.TestCheckResourceAttr(resourceName, "compare_type.0", "lines"),
					resource.TestCheckResourceAttr(resourceName, "compare_policy", "normal"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "next_compare_time"),
				),
			},
			{
				Config: testAccComparePolicy_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(resourceName, "period", "* * 1,2,3,4,5,6,7"),
					resource.TestCheckResourceAttr(resourceName, "begin_time", "01:00:00"),
					resource.TestCheckResourceAttr(resourceName, "end_time", "05:00:00"),
					resource.TestCheckResourceAttr(resourceName, "compare_type.0", "account"),
					resource.TestCheckResourceAttr(resourceName, "compare_policy", "normal"),
					resource.TestCheckResourceAttr(resourceName, "interval_hour", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "next_compare_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComparePolicy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_compare_policy" "test" {
  job_id         = "%s"
  period         = "* * 1,3,5"
  begin_time     = "00:00:00"
  end_time       = "04:00:00"
  compare_type   = ["lines"]
  compare_policy = "normal"
}
`, acceptance.HW_DRS_JOB_ID)
}

func testAccComparePolicy_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_compare_policy" "test" {
  job_id         = "%s"
  period         = "* * 1,2,3,4,5,6,7"
  begin_time     = "01:00:00"
  end_time       = "05:00:00"
  compare_type   = ["account"]
  compare_policy = "normal"
  interval_hour  = 2
}
`, acceptance.HW_DRS_JOB_ID)
}
