package dli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/dli/v1/queues"
	act "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDliQueue_basic(t *testing.T) {
	rName := act.RandomAccResourceName()
	resourceName := "huaweicloud_dli_queue.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { act.TestAccPreCheck(t) },
		Providers:    act.TestAccProviders,
		CheckDestroy: testAccCheckDliQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueue_basic(rName, dli.CU_16),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDliQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QUEUE_TYPE_SQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", fmt.Sprintf("%d", dli.CU_16)),
					resource.TestCheckResourceAttrSet(resourceName, "resource_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			//test scale_out
			{
				Config: testAccDliQueue_basic(rName, 2*dli.CU_16),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDliQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QUEUE_TYPE_SQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", fmt.Sprintf("%d", 2*dli.CU_16)),
				),
			},
			//test scale_in
			{
				Config: testAccDliQueue_basic(rName, dli.CU_16),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDliQueueExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QUEUE_TYPE_SQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", fmt.Sprintf("%d", dli.CU_16)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func testAccDliQueue_basic(rName string, cuCount int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_queue" "test" {
  name          = "%s"
  cu_count      = %d
  
  tags = {
    k1 = "1"
  }
}`, rName, cuCount)
}

func testAccCheckDliQueueDestroy(s *terraform.State) error {
	config := act.TestAccProvider.Meta().(*config.Config)
	client, err := config.DliV1Client(act.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("error creating Dli client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dli_queue" {
			continue
		}

		res, err := fetchDliQueueByQueueNameOnTest(rs.Primary.ID, client)
		if err == nil && res != nil {
			return fmtp.Errorf("huaweicloud_dli_queue still exists:%s,%+v,%+v", rs.Primary.ID, err, res)
		}
	}

	return nil
}

func testAccCheckDliQueueExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := act.TestAccProvider.Meta().(*config.Config)
		client, err := config.DliV1Client(act.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("error creating Dli client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmtp.Errorf("Error checking huaweicloud_dli_queue.queue exist, err=not found this resource")
		}
		_, err = fetchDliQueueByQueueNameOnTest(rs.Primary.ID, client)
		if err != nil {
			if strings.Contains(err.Error(), "Error finding the resource by list api") {
				return fmtp.Errorf("huaweicloud_dli_queue is not exist")
			}
			return fmtp.Errorf("Error checking huaweicloud_dli_queue.queue exist, err=%s", err)
		}
		return nil
	}
}

func fetchDliQueueByQueueNameOnTest(primaryID string,
	client *golangsdk.ServiceClient) (interface{}, error) {
	result := queues.Get(client, primaryID)
	return result.Body, result.Err
}
