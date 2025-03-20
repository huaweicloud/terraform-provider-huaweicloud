package lts

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

func getSearchCriteriaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSearchCriteria: Query the LTS search criteria detail
	var (
		getSearchCriteriaHttpUrl = "v1.0/{project_id}/groups/{group_id}/topics/{topic_id}/search-criterias"
		getSearchCriteriaProduct = "lts"
	)
	getSearchCriteriaClient, err := cfg.NewServiceClient(getSearchCriteriaProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS Client: %s", err)
	}

	groupID := state.Primary.Attributes["log_group_id"]
	streamID := state.Primary.Attributes["log_stream_id"]

	getSearchCriteriaPath := getSearchCriteriaClient.Endpoint + getSearchCriteriaHttpUrl
	getSearchCriteriaPath = strings.ReplaceAll(getSearchCriteriaPath, "{project_id}", getSearchCriteriaClient.ProjectID)
	getSearchCriteriaPath = strings.ReplaceAll(getSearchCriteriaPath, "{group_id}", groupID)
	getSearchCriteriaPath = strings.ReplaceAll(getSearchCriteriaPath, "{topic_id}", streamID)
	getSearchCriteriaPath = fmt.Sprintf("%s?search_type=%s", getSearchCriteriaPath, state.Primary.Attributes["type"])

	getSearchCriteriaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getSearchCriteriaResp, err := getSearchCriteriaClient.Request("GET", getSearchCriteriaPath, &getSearchCriteriaOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving search criteria: %s", err)
	}

	getSearchCriteriaRespBody, err := utils.FlattenResponse(getSearchCriteriaResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving search criteria: %s", err)
	}

	jsonPath := fmt.Sprintf("search_criterias[?id =='%s']|[0]", state.Primary.ID)
	getSearchCriteriaRespBody = utils.PathSearch(jsonPath, getSearchCriteriaRespBody, nil)
	if getSearchCriteriaRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getSearchCriteriaRespBody, nil
}

func TestAccSearchCriteria_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		searchCriteria      interface{}
		rName               = "huaweicloud_lts_search_criteria.test"
		withVisualization   = "huaweicloud_lts_search_criteria.visualization_log"
		rc                  = acceptance.InitResourceCheck(rName, &searchCriteria, getSearchCriteriaResourceFunc)
		withVisualizationRc = acceptance.InitResourceCheck(withVisualization, &searchCriteria, getSearchCriteriaResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			withVisualizationRc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testSearchCriteria_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "criteria", "context:test"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "ORIGINALLOG"),
					withVisualizationRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(withVisualization, "type", "VISUALIZATION"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceSearchCriteriaImportState(rName),
			},
			{
				ResourceName:      withVisualization,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceSearchCriteriaImportState(rName),
			},
		},
	})
}

func testSearchCriteria_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_lts_search_criteria" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  
  criteria = "context:test"
  name     = "%[1]s"
  type 	   = "ORIGINALLOG"
}

resource "huaweicloud_lts_search_criteria" "visualization_log" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  criteria      = "context:test"
  name          = "%[1]s_visual"
  type 	        = "VISUALIZATION"
}
`, name)
}

func resourceSearchCriteriaImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		searchCriteriaID := rs.Primary.ID
		groupID := rs.Primary.Attributes["log_group_id"]
		streamID := rs.Primary.Attributes["log_stream_id"]

		return fmt.Sprintf("%s/%s/%s", groupID, streamID, searchCriteriaID), nil
	}
}
