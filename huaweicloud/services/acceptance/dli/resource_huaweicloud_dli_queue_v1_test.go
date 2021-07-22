package dli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/dli/v1/queues"
	act "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDliQueueV1_basic(t *testing.T) {
	rName := fmt.Sprintf("tf_acc_test_dli_queue_%s", acctest.RandString(5))
	resourceName := "huaweicloud_dli_queue_v1.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { act.TestAccPreCheck(t) },
		Providers:    act.TestAccProviders,
		CheckDestroy: testAccCheckDliQueueV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueueV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDliQueueV1Exists(resourceName),
				),
			},
		},
	})
}

func testAccDliQueueV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_queue_v1" "test" {
  name          = "%s"
  cu_count      = 16
  resource_mode = 0
  
  tags = {
    k1 = "1"
  }
}`, rName)
}

func testAccCheckDliQueueV1Destroy(s *terraform.State) error {
	config := act.TestAccProvider.Meta().(*config.Config)
	client, err := config.DliV1Client(act.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("error creating Dli client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dli_queue_v1" {
			continue
		}

		res, err := fetchDliQueueV1ByQueueNameOnTest(rs.Primary.ID, client)
		if err == nil && res != nil {
			return fmtp.Errorf("huaweicloud_dli_queue_v1 still exists:%s,%+v,%+v", rs.Primary.ID, err, res)
		}
	}

	return nil
}

func testAccCheckDliQueueV1Exists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := act.TestAccProvider.Meta().(*config.Config)
		client, err := config.DliV1Client(act.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("error creating Dli client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmtp.Errorf("Error checking huaweicloud_dli_queue_v1.queue exist, err=not found this resource")
		}
		_, err = fetchDliQueueV1ByQueueNameOnTest(rs.Primary.ID, client)
		if err != nil {
			if strings.Contains(err.Error(), "Error finding the resource by list api") {
				return fmtp.Errorf("huaweicloud_dli_queue_v1 is not exist")
			}
			return fmtp.Errorf("Error checking huaweicloud_dli_queue_v1.queue exist, err=%s", err)
		}
		return nil
	}
}

func fetchDliQueueV1ByQueueNameOnTest(primaryID string,
	client *golangsdk.ServiceClient) (interface{}, error) {
	result := queues.Get(client, primaryID)
	return result.Body, result.Err
}
