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

func getHostCrossAccountAccessResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	ltsClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	listHostAccessConfigHttpUrl := "v3/{project_id}/lts/access-config-list"
	listHostAccessConfigPath := ltsClient.Endpoint + listHostAccessConfigHttpUrl
	listHostAccessConfigPath = strings.ReplaceAll(listHostAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	listHostAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	name := state.Primary.Attributes["name"]
	listHostAccessConfigOpt.JSONBody = map[string]interface{}{
		"access_config_name_list": []string{name},
	}

	listHostAccessConfigResp, err := ltsClient.Request("POST", listHostAccessConfigPath, &listHostAccessConfigOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving LTS cross account access: %s", err)
	}

	listHostAccessConfigRespBody, err := utils.FlattenResponse(listHostAccessConfigResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten LTS cross account access response: %s", err)
	}

	jsonPath := fmt.Sprintf("result[?access_config_name=='%s']|[0]", name)
	listHostAccessConfigRespBody = utils.PathSearch(jsonPath, listHostAccessConfigRespBody, nil)
	if listHostAccessConfigRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return listHostAccessConfigRespBody, nil
}

func TestAccCrossAccountAccess_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_cross_account_access.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHostCrossAccountAccessResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSCrossAccountAccess(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCrossAccuntAccessBasic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "access_config_type"),
				),
			},
			{
				Config: testCrossAccuntAccessBasicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
				),
			},
		},
	})
}

func testCrossAccuntAccessBasic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_cross_account_access" "test" {
  name               = "%s"
  agency_project_id  = "%s"
  agency_domain_name = "%s"
  agency_name        = "%s"

  log_agencystream_name = "%s"
  log_agencystream_id   = "%s"
  log_agencygroup_name  = "%s"
  log_agencygroup_id    = "%s"

  log_stream_name = "%s"
  log_stream_id   = "%s"
  log_group_name  = "%s"
  log_group_id    = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_LTS_AGENCY_PROJECT_ID, acceptance.HW_LTS_AGENCY_DOMAIN_NAME, acceptance.HW_LTS_AGENCY_NAME,
		acceptance.HW_LTS_AGENCY_STREAM_NAME, acceptance.HW_LTS_AGENCY_STREAM_ID, acceptance.HW_LTS_AGENCY_GROUP_NAME,
		acceptance.HW_LTS_AGENCY_GROUP_ID, acceptance.HW_LTS_LOG_STREAM_NAME, acceptance.HW_LTS_LOG_STREAM_ID,
		acceptance.HW_LTS_LOG_GROUP_NAME, acceptance.HW_LTS_LOG_GROUP_ID)
}

func testCrossAccuntAccessBasicUpdate(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_cross_account_access" "test" {
  name               = "%s"
  agency_project_id  = "%s"
  agency_domain_name = "%s"
  agency_name        = "%s"

  log_agencystream_name = "%s"
  log_agencystream_id   = "%s"
  log_agencygroup_name  = "%s"
  log_agencygroup_id    = "%s"

  log_stream_name = "%s"
  log_stream_id   = "%s"
  log_group_name  = "%s"
  log_group_id    = "%s"

  tags = {
    foo  = "bar1"
    key1 = "value1"
  }
}
`, name, acceptance.HW_LTS_AGENCY_PROJECT_ID, acceptance.HW_LTS_AGENCY_DOMAIN_NAME, acceptance.HW_LTS_AGENCY_NAME,
		acceptance.HW_LTS_AGENCY_STREAM_NAME, acceptance.HW_LTS_AGENCY_STREAM_ID, acceptance.HW_LTS_AGENCY_GROUP_NAME,
		acceptance.HW_LTS_AGENCY_GROUP_ID, acceptance.HW_LTS_LOG_STREAM_NAME, acceptance.HW_LTS_LOG_STREAM_ID,
		acceptance.HW_LTS_LOG_GROUP_NAME, acceptance.HW_LTS_LOG_GROUP_ID)
}
