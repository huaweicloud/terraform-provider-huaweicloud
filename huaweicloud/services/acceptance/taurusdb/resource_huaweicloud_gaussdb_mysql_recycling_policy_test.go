package taurusdb

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

func getGaussDBRecyclingPolicyResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/recycle-policy"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB recycling policy: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	retentionPeriod := utils.PathSearch("recycle_policy.retention_period_in_days", getRespBody, nil)
	if retentionPeriod == nil {
		return nil, fmt.Errorf("error retrieving GaussDB recycling policy, retention_period_in_days is not found")
	}
	return getRespBody, nil
}

func TestAccGaussDBRecyclingPolicy_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_gaussdb_mysql_recycling_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBRecyclingPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testRecyclingPolicy_basic("3"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "retention_period_in_days", "3"),
				),
			},
			{
				Config: testRecyclingPolicy_basic("6"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "retention_period_in_days", "6"),
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

func testRecyclingPolicy_basic(retentionPeriodInDays string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_recycling_policy" "test" {
  retention_period_in_days = "%s"
}
`, retentionPeriodInDays)
}
