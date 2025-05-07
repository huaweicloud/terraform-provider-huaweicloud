package kafka

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDmsKafkaSmartConnectResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getKafkaSmartConnect: query DMS kafka smart connect
	var (
		getKafkaSmartConnectHttpUrl = "v2/{project_id}/instances/{instance_id}"
		getKafkaSmartConnectProduct = "dms"
	)
	getKafkaSmartConnectClient, err := cfg.NewServiceClient(getKafkaSmartConnectProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	getKafkaSmartConnectPath := getKafkaSmartConnectClient.Endpoint + getKafkaSmartConnectHttpUrl
	getKafkaSmartConnectPath = strings.ReplaceAll(getKafkaSmartConnectPath, "{project_id}",
		getKafkaSmartConnectClient.ProjectID)
	getKafkaSmartConnectPath = strings.ReplaceAll(getKafkaSmartConnectPath, "{instance_id}", instanceID)

	getKafkaSmartConnectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getKafkaSmartConnectResp, err := getKafkaSmartConnectClient.Request("GET", getKafkaSmartConnectPath,
		&getKafkaSmartConnectOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DMS kafka smart connect: %s", err)
	}

	getKafkaSmartConnectRespBody, err := utils.FlattenResponse(getKafkaSmartConnectResp)
	if err != nil {
		return nil, err
	}
	connectorId := utils.PathSearch("connector_id", getKafkaSmartConnectRespBody, nil)
	if connectorId != state.Primary.ID {
		return nil, fmt.Errorf("error retrieving DMS kafka smart connect: the connector can not be found")
	}
	return getKafkaSmartConnectRespBody, nil
}

func TestAccDmsKafkaSmartConnect_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_smart_connect.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsKafkaSmartConnectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkaSmartConnect_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "node_count", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_spec_code", "bandwidth", "node_count"},
				ImportStateIdFunc:       testKafkaSmartConnectResourceImportState(resourceName),
			},
		},
	})
}

func testDmsKafkaSmartConnect_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_smart_connect" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  node_count  = 2
}
`, testAccKafkaInstance_newFormat(name))
}

// testKafkaSmartConnectResourceImportState is used to return an import id with format <instance_id>/<id>
func testKafkaSmartConnectResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instance_id := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instance_id, rs.Primary.ID), nil
	}
}
