package cfw

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

func getBlackWhiteListResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBlackWhiteList: Query the CFW black white list detail
	var (
		getBlackWhiteListHttpUrl = "v1/{project_id}/black-white-lists"
		getBlackWhiteListProduct = "cfw"
	)
	getBlackWhiteListClient, err := cfg.NewServiceClient(getBlackWhiteListProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW Client: %s", err)
	}

	getBlackWhiteListPath := getBlackWhiteListClient.Endpoint + getBlackWhiteListHttpUrl
	getBlackWhiteListPath = strings.ReplaceAll(getBlackWhiteListPath, "{project_id}", getBlackWhiteListClient.ProjectID)

	getBlackWhiteListqueryParams := buildGetBlackWhiteListQueryParams(state)
	getBlackWhiteListPath += getBlackWhiteListqueryParams

	getBlackWhiteListOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getBlackWhiteListResp, err := getBlackWhiteListClient.Request("GET", getBlackWhiteListPath, &getBlackWhiteListOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving BlackWhiteList: %s", err)
	}

	getBlackWhiteListRespBody, err := utils.FlattenResponse(getBlackWhiteListResp)
	if err != nil {
		return nil, err
	}

	list := utils.PathSearch("data.records", getBlackWhiteListRespBody, make([]interface{}, 0)).([]interface{})
	if len(list) != 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return list[0], nil
}

func TestAccBlackWhiteList_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_cfw_black_white_list.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBlackWhiteListResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBlackWhiteList_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "list_type", "4"),
					resource.TestCheckResourceAttr(rName, "direction", "0"),
					resource.TestCheckResourceAttr(rName, "protocol", "6"),
					resource.TestCheckResourceAttr(rName, "port", "22"),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
					resource.TestCheckResourceAttr(rName, "address", "1.1.1.1"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
				),
			},
			{
				Config: testBlackWhiteList_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "list_type", "4"),
					resource.TestCheckResourceAttr(rName, "direction", "1"),
					resource.TestCheckResourceAttr(rName, "protocol", "-1"),
					resource.TestCheckResourceAttr(rName, "port", ""),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
					resource.TestCheckResourceAttr(rName, "address", "2.2.2.0/24"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testBlackWhiteListImportState(rName),
			},
		},
	})
}

func testBlackWhiteList_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_black_white_list" "test" {
  object_id    = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type    = 4
  direction    = 0
  protocol     = 6
  port         = "22"
  address_type = 0
  address      = "1.1.1.1"
  description  = "terraform test"
}
`, testAccDatasourceFirewalls_basic())
}

func testBlackWhiteList_basic_update() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_black_white_list" "test" {
  object_id    = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type    = 4
  direction    = 1
  protocol     = -1
  address_type = 0
  address      = "2.2.2.0/24"
  description  = ""
}
`, testAccDatasourceFirewalls_basic())
}

func buildGetBlackWhiteListQueryParams(state *terraform.ResourceState) string {
	res := "?offset=0&limit=10"
	res = fmt.Sprintf("%s&object_id=%v", res, state.Primary.Attributes["object_id"])
	res = fmt.Sprintf("%s&list_type=%v", res, state.Primary.Attributes["list_type"])
	res = fmt.Sprintf("%s&address=%v", res, strings.Split(state.Primary.Attributes["address"], "/")[0])

	return res
}

func testBlackWhiteListImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["object_id"] == "" {
			return "", fmt.Errorf("attribute (object_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["list_type"] == "" {
			return "", fmt.Errorf("attribute (list_type) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["address"] == "" {
			return "", fmt.Errorf("attribute (address) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["object_id"] + "/" +
			rs.Primary.Attributes["list_type"] + "/" + rs.Primary.Attributes["address"], nil
	}
}
