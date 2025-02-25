package kafka

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDmsKafkaUserClientQuotaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getKafkaUserClientQuota: query DMS kafka user client quota
	var (
		getKafkaUserClientQuotaHttpUrl = "v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota"
		getKafkaUserClientQuotaProduct = "dms"
	)
	getKafkaUserClientQuotaClient, err := cfg.NewServiceClient(getKafkaUserClientQuotaProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS Client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<user>/<user_default>/<client>/<client_default>")
	}
	instanceID := parts[0]

	getKafkaUserClientQuotaPath := getKafkaUserClientQuotaClient.Endpoint + getKafkaUserClientQuotaHttpUrl
	getKafkaUserClientQuotaPath = strings.ReplaceAll(getKafkaUserClientQuotaPath, "{project_id}",
		getKafkaUserClientQuotaClient.ProjectID)
	getKafkaUserClientQuotaPath = strings.ReplaceAll(getKafkaUserClientQuotaPath, "{instance_id}", instanceID)

	getKafkaUserClientQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getKafkaUserClientQuotaResp, err := getKafkaUserClientQuotaClient.Request("GET", getKafkaUserClientQuotaPath,
		&getKafkaUserClientQuotaOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving the quota: %s", err)
	}

	getKafkaUserClientQuotaRespBody, err := utils.FlattenResponse(getKafkaUserClientQuotaResp)
	if err != nil {
		return nil, err
	}
	quota := filterUserClientQuota(parts, getKafkaUserClientQuotaRespBody)
	if quota == nil {
		return nil, fmt.Errorf("can not find the quota")
	}

	return quota, nil
}

func filterUserClientQuota(parts []string, resp interface{}) interface{} {
	quotaJson := utils.PathSearch("quotas", resp, make([]interface{}, 0))
	quotaArray := quotaJson.([]interface{})
	if len(quotaArray) < 1 || len(parts) != 5 {
		return nil
	}

	rawUserDefault, _ := strconv.ParseBool(parts[2])
	rawClientDefault, _ := strconv.ParseBool(parts[4])

	for _, quota := range quotaArray {
		user := utils.PathSearch("user", quota, nil)
		userDefault := utils.PathSearch(`"user-default"`, quota, false).(bool)
		client := utils.PathSearch("client", quota, nil)
		clientDefault := utils.PathSearch(`"client-default"`, quota, false).(bool)
		if parts[1] != user {
			continue
		}
		if rawUserDefault != userDefault {
			continue
		}
		if parts[3] != client {
			continue
		}
		if rawClientDefault != clientDefault {
			continue
		}
		return quota
	}
	return nil
}

func TestAccDmsKafkaUserClientQuota_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_user_client_quota.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsKafkaUserClientQuotaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkaUserClientQuota_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user", rName),
					resource.TestCheckResourceAttr(resourceName, "user_default", "false"),
					resource.TestCheckResourceAttr(resourceName, "client", rName),
					resource.TestCheckResourceAttr(resourceName, "client_default", "false"),
					resource.TestCheckResourceAttr(resourceName, "producer_byte_rate", "1024"),
					resource.TestCheckResourceAttr(resourceName, "consumer_byte_rate", "0"),
				),
			},
			{
				Config: testDmsKafkaUserClientQuota_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "producer_byte_rate", "0"),
					resource.TestCheckResourceAttr(resourceName, "consumer_byte_rate", "2048"),
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

func testDmsKafkaUserClientQuota_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_user_client_quota" "test" {
  instance_id        = huaweicloud_dms_kafka_instance.test.id
  user               = "%[2]s"
  client             = "%[2]s"
  producer_byte_rate = 1024
  consumer_byte_rate = 0
}
`, testAccKafkaInstance_newFormat(rName), rName)
}

func testDmsKafkaUserClientQuota_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_user_client_quota" "test" {
  instance_id        = huaweicloud_dms_kafka_instance.test.id
  user               = "%[2]s"
  client             = "%[2]s"
  producer_byte_rate = 0
  consumer_byte_rate = 2048
}
`, testAccKafkaInstance_newFormat(rName), rName)
}
