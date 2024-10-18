package dli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	funccommon "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDatasourceConnectionPrivilegeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{connection_id}/privileges"
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connection_id}", state.Primary.Attributes["connection_id"])

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, funccommon.ConvertExpected400ErrInto404Err(err, "error_code", dli.ErrCodeConnNotFound)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	privilege := utils.PathSearch(fmt.Sprintf("privileges[?project_id=='%v']|[0]", state.Primary.Attributes["project_id"]), respBody, nil)
	if privilege == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return privilege, nil
}

func TestAccDatasourceConnectionPrivilege_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_datasource_connection_privilege.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDatasourceConnectionPrivilegeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDatasourceConnectionPrivilege_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "connection_id", "huaweicloud_dli_datasource_connection.test", "id"),
					resource.TestCheckResourceAttr(rName, "project_id", acceptance.HW_DEST_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "privileges.0", "BIND_QUEUE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDatasourceConnectionPrivilegeImportStateFunc(rName),
			},
		},
	})
}

func testAccDatasourceConnectionPrivilegeImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var connectionId, projectId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of datasource connection privilege is not found in the tfstate", rsName)
		}
		connectionId = rs.Primary.Attributes["connection_id"]
		projectId = rs.Primary.Attributes["project_id"]
		if connectionId == "" || projectId == "" {
			return "", fmt.Errorf("the project ID or related connection ID is missing")
		}
		return fmt.Sprintf("%s/%s", connectionId, projectId), nil
	}
}

func testDatasourceConnectionPrivilege_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_datasource_connection" "test" {
  name      = "%[2]s"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
}
`, common.TestVpc(name), name)
}

func testDatasourceConnectionPrivilege_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_datasource_connection_privilege" "test" {
  connection_id = huaweicloud_dli_datasource_connection.test.id
  project_id    = "%[2]s"
}
`, testDatasourceConnectionPrivilege_base(name), acceptance.HW_DEST_PROJECT_ID_TEST)
}
