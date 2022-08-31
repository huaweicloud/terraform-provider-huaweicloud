package cmdb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
)

func getDashboardResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf, "lts", acceptance.HW_REGION_NAME)
	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "lts", conf.Region, "v2/"+conf.GetProjectID(conf.Region)+
			"/dashboards?id="+state.Primary.ID).WithHeader(map[string]string{"content-type": "application/json;charset=UTF8"})
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.ReadDashBoardResp{}
	if response.StatusCode == 200 {
		err = json.Unmarshal(body, rlt)
		if err != nil {
			return nil, fmt.Errorf("error umarshal")
		}
		if len(rlt.Results) == 0 {
			return nil, fmt.Errorf("error resource donesn`t exist")
		}
		return rlt, nil
	}
	return nil, fmt.Errorf("error get resource %v,%v,%s", response.StatusCode, err, string(body))
}

func TestAccDashboard_basic(t *testing.T) {
	var instance entity.ReadDashBoardResp
	resourceName := "huaweicloud_lts_dashboard.dashboard_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDashboardResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckInternal(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesDashboard_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "log_group_id", "9ac33c09-7f00-4eed-b9a0-0ffaad7a64d1"),
					resource.TestCheckResourceAttr(resourceName, "log_group_name", "CTS"),
					resource.TestCheckResourceAttr(resourceName, "log_stream_id", "c3ab6968-a903-493d-a49a-5c45caaf32b4"),
					resource.TestCheckResourceAttr(resourceName, "log_stream_name", "test-zhb"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_title"},
				ImportStateIdFunc:       testAccDashboardImportStateIdFunc(),
			},
		},
	})
}

func testAccDashboardImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var logGroupId, logGroupName, logStreamId, logStreamName, id string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_lts_dashboard" {
				logGroupId = rs.Primary.Attributes["log_group_id"]
				logGroupName = rs.Primary.Attributes["log_group_name"]
				logStreamId = rs.Primary.Attributes["log_stream_id"]
				logStreamName = rs.Primary.Attributes["log_stream_name"]
				id = rs.Primary.ID
			}
		}
		if logGroupId == "" || logGroupName == "" || logStreamId == "" || logStreamName == "" || id == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s/%s/%s", id, logGroupId, logGroupName, logStreamId, logStreamName)
		}
		return fmt.Sprintf("%s/%s/%s/%s/%s", id, logGroupId, logGroupName, logStreamId, logStreamName), nil
	}
}

func tesDashboard_basic() string {
	return `
resource "huaweicloud_lts_dashboard" "dashboard_1" {
  log_group_id        = "9ac33c09-7f00-4eed-b9a0-0ffaad7a64d1"
  log_group_name      = "CTS"
  log_stream_id = "c3ab6968-a903-493d-a49a-5c45caaf32b4"
  log_stream_name = "test-zhb"
  is_delete_charts      = "true"
  template_title   = ["cfw-log-analysis"]
  
}`
}
