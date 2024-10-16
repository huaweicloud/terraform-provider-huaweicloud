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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getServiceGroupMemberResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getServiceGroupMember: Query the CFW service group member detail
	var (
		getServiceGroupMemberHttpUrl = "v1/{project_id}/service-items"
		getServiceGroupMemberProduct = "cfw"
	)
	getServiceGroupMemberClient, err := cfg.NewServiceClient(getServiceGroupMemberProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	getServiceGroupMemberPath := getServiceGroupMemberClient.Endpoint + getServiceGroupMemberHttpUrl
	getServiceGroupMemberPath = strings.ReplaceAll(getServiceGroupMemberPath, "{project_id}", getServiceGroupMemberClient.ProjectID)

	getServiceGroupMemberqueryParams := buildGetServiceGroupMemberQueryParams(state)
	getServiceGroupMemberPath += getServiceGroupMemberqueryParams

	getServiceGroupMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getServiceGroupMemberResp, err := getServiceGroupMemberClient.Request("GET", getServiceGroupMemberPath, &getServiceGroupMemberOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ServiceGroupMember: %s", err)
	}

	getServiceGroupMemberRespBody, err := utils.FlattenResponse(getServiceGroupMemberResp)
	if err != nil {
		return nil, err
	}

	members := utils.PathSearch("data.records", getServiceGroupMemberRespBody, make([]interface{}, 0)).([]interface{})
	return cfw.FilterServiceGroupMembers(members, state.Primary.ID)
}

func TestAccServiceGroupMember_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_service_group_member.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getServiceGroupMemberResourceFunc,
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
				Config: testServiceGroupMember_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_cfw_service_group.test", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "6"),
					resource.TestCheckResourceAttr(rName, "source_port", "80"),
					resource.TestCheckResourceAttr(rName, "dest_port", "22"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testServiceGroupMemberImportState(rName),
			},
		},
	})
}

func testServiceGroupMember_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_service_group_member" "test" {
  group_id    = huaweicloud_cfw_service_group.test.id
  protocol    = 6
  source_port = "80"
  dest_port   = "22"
}
`, testServiceGroup_basic(name))
}

func buildGetServiceGroupMemberQueryParams(state *terraform.ResourceState) string {
	res := "?offset=0&limit=1024"
	res = fmt.Sprintf("%s&set_id=%v", res, state.Primary.Attributes["group_id"])

	return res
}

func testServiceGroupMemberImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["group_id"] == "" {
			return "", fmt.Errorf("attribute (group_id) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["group_id"] + "/" + rs.Primary.ID, nil
	}
}
