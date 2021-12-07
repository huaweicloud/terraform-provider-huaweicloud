package dis

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/dis/v2/streams"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDisStreamsResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DisV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating DIS V2 client, err=%s", err)
	}

	return streams.Get(client, state.Primary.ID, streams.GetOpts{})
}

func TestAccResourceDisStream_basic(t *testing.T) {
	var streamInstance streams.CreateOpts
	resourceName := "huaweicloud_dis_stream.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&streamInstance,
		getDisStreamsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDisStream_basic(name, 1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", name),
					resource.TestCheckResourceAttr(resourceName, "partitions.#", "1"),
				),
			},
			{
				Config: testAccDisStream_basic(name, 4),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", name),
					resource.TestCheckResourceAttr(resourceName, "partitions.#", "4"),
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

func testAccDisStream_basic(rName string, partitionCount int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dis_stream" "test" {
  stream_name     = "%s"
  partition_count = %d
}
`, rName, partitionCount)
}

func TestAccResourceDisStream_all(t *testing.T) {
	var streamInstance streams.CreateOpts
	resourceName := "huaweicloud_dis_stream.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&streamInstance,
		getDisStreamsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDisStream_All(name, 1, 2, "bar"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", name),
					resource.TestCheckResourceAttr(resourceName, "partition_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale_min_partition_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale_max_partition_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "partitions.#", "1"),
				),
			},
			{
				Config: testAccDisStream_All(name, 2, 3, "bar2"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", name),
					resource.TestCheckResourceAttr(resourceName, "partition_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale_min_partition_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "auto_scale_max_partition_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "partitions.#", "2"),
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

func testAccDisStream_All(rName string, partitionCount int, scaleCount int, tagValue string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dis_stream" "test" {
  stream_name                    = "%s"
  partition_count                = %d
  stream_type                    = "COMMON"
  retention_period               = 24
  auto_scale_min_partition_count = %d
  auto_scale_max_partition_count = %d
  compression_format             = "zip"

  data_type     = "CSV"
  csv_delimiter = ";"
  data_schema   = "{\"type\":\"record\",\"name\":\"RecordName\",\"fields\":[{\"type\":\"string\",\"name\":\"name\"}]}"

  tags = {
    foo = "%s"
    key = "value"
  }
}
`, rName, partitionCount, scaleCount, scaleCount+1, tagValue)
}
