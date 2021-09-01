package huaweicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/maas/v1/task"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccMaasTask_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMaas(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMaasTaskV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaasTaskV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMaasTaskV1Exists("huaweicloud_oms_task.task_1"),
					resource.TestCheckResourceAttr("huaweicloud_oms_task.task_1", "description", "migration task"),
				),
			},
		},
	})
}

func testAccCheckMaasTaskV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	maasClient, err := config.MaasV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud maas client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_oms_task" {
			continue
		}

		_, err := task.Get(maasClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Maas task still exists")
		}
	}

	return nil
}

func testAccCheckMaasTaskV1Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		maasClient, err := config.MaasV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud maas client: %s", err)
		}

		found, err := task.Get(maasClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if strconv.FormatInt(found.ID, 10) != rs.Primary.ID {
			return fmtp.Errorf("Task not found")
		}

		return nil
	}
}

var testAccMaasTaskV1_basic = fmt.Sprintf(`
resource "huaweicloud_oms_task" "task_1" {
  description = "migration task"
  enable_kms = false
  thread_num = 1
  src_node {
    region = "cn-beijing"
	ak = "%s"
	sk = "%s"
    object_key = "123.txt"
    bucket = "oms-bucket"
  }
  dst_node {
    region = "%s"
	ak = "%s"
	sk = "%s"
    object_key = "oms"
    bucket = "oms-test"
  }
}
`, HW_SRC_ACCESS_KEY, HW_SRC_SECRET_KEY, HW_REGION_NAME, HW_ACCESS_KEY, HW_SECRET_KEY)
