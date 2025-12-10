package smn

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTopicSubscriber(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/subscriptions"
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{topic_urn}", state.Primary.Attributes["topic_urn"])

	getDwsAlarmSubsResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, golangsdk.ErrDefault404{}
	}

	getRespJson, err := json.Marshal(getDwsAlarmSubsResp)
	if err != nil {
		return nil, err
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}

	subscription := utils.PathSearch(fmt.Sprintf("subscriptions[?subscription_urn=='%s']|[0]", state.Primary.ID),
		getRespBody, nil)
	if subscription == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccTopicSubscriber_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_smn_topic_subscriber.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getTopicSubscriber,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSmnSubscribeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTopicSubscriber_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "topic_urn",
						"huaweicloud_smn_topic.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "owner"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "remark"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"subscribe_id"},
			},
		},
	})
}

func testAccTopicSubscriber_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "The display name"
}

resource "huaweicloud_smn_topic_subscriber" "test" {
  topic_urn    = huaweicloud_smn_topic.test.id
  subscribe_id = "%[2]s"
}
`, rName, acceptance.HW_SMN_SUBSCRIBE_ID)
}
