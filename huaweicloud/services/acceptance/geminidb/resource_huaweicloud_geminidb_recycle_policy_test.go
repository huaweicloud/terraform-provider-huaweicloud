package geminidb

import (
	"errors"
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

func getRecyclePolicyResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geminidb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/recycle-policy"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GeminiDB recycle policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	retentionPeriod := utils.PathSearch("recycle_policy.retention_period_in_days", respBody, nil)
	if retentionPeriod == nil {
		return nil, errors.New("error retrieving GeminiDB recycling policy, retention_period_in_days is not found")
	}

	return respBody, nil
}

func TestAccGeminiDBRecyclePolicy_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_geminidb_recycle_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRecyclePolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBRecyclePolicy_basic(7),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "retention_period_in_days", "7"),
				),
			},
			{
				Config: testAccGeminiDBRecyclePolicy_basic(3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "retention_period_in_days", "3"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGeminiDBRecyclePolicy_basic(retentionPeriodInDays int) string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_recycle_policy" "test" {
  retention_period_in_days = %#v
}
`, retentionPeriodInDays)
}
