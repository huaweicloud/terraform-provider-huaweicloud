package dli

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v1/queues"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	act "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDliQueueResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating Dli v1 client, err=%s", err)
	}

	result := queues.Get(client, state.Primary.ID)
	return result.Body, result.Err
}

func TestAccDliQueue_basic(t *testing.T) {
	rName := act.RandomAccResourceName()
	resourceName := "huaweicloud_dli_queue.test"

	var obj queues.CreateOpts
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { act.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueue_basic(rName, dli.CU_16),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QUEUE_TYPE_SQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", fmt.Sprintf("%d", dli.CU_16)),
					resource.TestCheckResourceAttrSet(resourceName, "resource_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
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
  name     = "%s"
  cu_count = %d

  tags = {
    foo = "bar"
  }
}
`, rName, cuCount)
}

func TestAccDliQueue_cidr(t *testing.T) {
	rName := act.RandomAccResourceName()
	resourceName := "huaweicloud_dli_queue.test"

	var obj queues.CreateOpts
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { act.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueue_cidr(rName, "172.16.0.0/21"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QUEUE_TYPE_SQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", "16"),
					resource.TestCheckResourceAttr(resourceName, "resource_mode", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_cidr", "172.16.0.0/21"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{

				Config: testAccDliQueue_cidr(rName, "172.16.0.0/18"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QUEUE_TYPE_SQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", "16"),
					resource.TestCheckResourceAttr(resourceName, "resource_mode", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_cidr", "172.16.0.0/18"),
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

func testAccDliQueue_cidr(rName string, cidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_queue" "test" {
  name          = "%s"
  cu_count      = 16
  resource_mode = 1
  vpc_cidr      = "%s"

  tags = {
    foo = "bar"
  }
}`, rName, cidr)
}
