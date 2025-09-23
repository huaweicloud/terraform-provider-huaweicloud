package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v2/batches"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getSparkJobResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.DliV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI v2 client: %s", err)
	}
	return batches.Get(c, state.Primary.ID)
}

func TestAccDliSparkJobV2_basic(t *testing.T) {
	var job batches.CreateResp

	rName := acceptance.RandomAccResourceName()
	dashName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dli_spark_job.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&job,
		getSparkJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliGenaralQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSparkJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDliSparkJob_basic(rName, dashName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "queue_name", acceptance.HW_DLI_GENERAL_QUEUE_NAME),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckDliSparkJobDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.DliV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Dli v2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dli_spark_job" {
			continue
		}

		resp, err := batches.GetState(client, rs.Primary.ID)
		// If the status of the spark job is "dead" or "success", it means that the life cycle of the job has ended.
		if err == nil && resp != nil && (resp.State != batches.StateDead && resp.State != batches.StateSuccess) {
			return fmt.Errorf("spark job (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccDliSparkJob_basic(name, dashName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_spark_job" "test" {
  queue_name = "%s"
  name       = "%s"
  app_name   = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  
 depends_on = [
    huaweicloud_obs_bucket.test,
    huaweicloud_obs_bucket_object.test,
  ]
}
`, testAccDliPackage_basic(dashName), acceptance.HW_DLI_GENERAL_QUEUE_NAME, name)
}
