package hss

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

func getEventLoginWhiteListResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	var (
		epsId         = state.Primary.Attributes["enterprise_project_id"]
		privateIP     = state.Primary.Attributes["private_ip"]
		loginIP       = state.Primary.Attributes["login_ip"]
		loginUserName = state.Primary.Attributes["login_user_name"]
	)

	queryPath := client.Endpoint + "v5/{project_id}/event/white-list/login"
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = fmt.Sprintf("%s?limit=200&offset=0&private_ip=%s&login_ip=%s&login_user_name=%s",
		queryPath, privateIP, loginIP, loginUserName)

	if epsId != "" {
		queryPath = fmt.Sprintf("%s&enterprise_project_id=%s", queryPath, epsId)
	}

	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", queryPath, &queryOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS event login white list: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(dataList) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return dataList[0], nil
}

func TestAccEventLoginWhiteList_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_hss_event_login_white_list.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getEventLoginWhiteListResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEventLoginWhiteList_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "private_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "login_ip", "192.168.0.2"),
					resource.TestCheckResourceAttr(resourceName, "login_user_name", "user_test"),
					resource.TestCheckResourceAttr(resourceName, "remarks", "remarks_test"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_name"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
				),
			},
		},
	})
}

const testEventLoginWhiteList_basic string = `
resource "huaweicloud_hss_event_login_white_list" "test" {
  private_ip            = "192.168.0.1"
  login_ip              = "192.168.0.2"
  login_user_name       = "user_test"
  remarks               = "remarks_test"
  handle_event          = true
  enterprise_project_id = "0"
}
`
