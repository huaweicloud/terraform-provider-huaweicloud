package kafka

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDmsKafkaConsumerGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getKafkaConsumerGroup: query DMS kafka consumer group
	var (
		getKafkaConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}"
		getKafkaConsumerGroupProduct = "dms"
	)
	getKafkaConsumerGroupClient, err := cfg.NewServiceClient(getKafkaConsumerGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS Client: %s", err)
	}

	// Split instance_id and group from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<consumerGroup>")
	}
	instanceID := parts[0]
	name := parts[1]
	getKafkaConsumerGroupPath := getKafkaConsumerGroupClient.Endpoint + getKafkaConsumerGroupHttpUrl
	getKafkaConsumerGroupPath = strings.ReplaceAll(getKafkaConsumerGroupPath, "{project_id}",
		getKafkaConsumerGroupClient.ProjectID)
	getKafkaConsumerGroupPath = strings.ReplaceAll(getKafkaConsumerGroupPath, "{instance_id}", instanceID)
	getKafkaConsumerGroupPath = strings.ReplaceAll(getKafkaConsumerGroupPath, "{group}", name)

	getKafkaConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getKafkaConsumerGroupResp, err := getKafkaConsumerGroupClient.Request("GET", getKafkaConsumerGroupPath,
		&getKafkaConsumerGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DMS Kafka consumer group: %s", err)
	}

	getKafkaConsumerGroupRespBody, err := utils.FlattenResponse(getKafkaConsumerGroupResp)
	if err != nil {
		return nil, err
	}
	groupJson := utils.PathSearch("group", getKafkaConsumerGroupRespBody, nil)
	groupState := utils.PathSearch("state", groupJson, nil)
	if groupState == "DEAD" {
		return nil, fmt.Errorf("error retrieving DMS Kafka consumer group")
	}
	return getKafkaConsumerGroupRespBody, nil
}

func TestAccDmsKafkaConsumerGroup_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_consumer_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsKafkaConsumerGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkaConsumerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "lag"),
					resource.TestCheckResourceAttrSet(resourceName, "coordinator_id"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "add"),
					resource.TestCheckResourceAttr(resourceName, "state", "EMPTY"),
					resource.TestCheckResourceAttr(resourceName, "lag", "0"),
				),
			},
			{
				Config: testDmsKafkaConsumerGroup_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func testDmsKafkaConsumerGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_consumer_group" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%s"
  description = "add"
}
`, testAccKafkaInstance_newFormat(name), name)
}

func testDmsKafkaConsumerGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_consumer_group" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id  
  name        = "%s"  
  description = ""
}
`, testAccKafkaInstance_newFormat(name), name)
}
