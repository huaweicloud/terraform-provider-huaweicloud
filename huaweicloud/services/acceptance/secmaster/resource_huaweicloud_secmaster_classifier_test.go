package secmaster

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

func getClassifierResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		product     = "secmaster"
		region      = acceptance.HW_REGION_NAME
		workspaceId = state.Primary.Attributes["workspace_id"]
		id          = state.Primary.ID
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/classifiers/{classifier_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{classifier_id}", id)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SecMaster classifier: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	mappingInfo := utils.PathSearch("data.mapping_info", respBody, nil)
	if mappingInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return mappingInfo, nil
}

func TestAccClassifier_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_secmaster_classifier.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getClassifierResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testClassifier_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "dataclass_id",
						"data.huaweicloud_secmaster_data_classes.test", "data_classes.0.id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "data_source", "CFW"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "classifier.0.direct_classifier", "false"),
					resource.TestCheckResourceAttrSet(rName, "mapping_id"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testClassifier_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "classifier.0.direct_classifier", "true"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testClassifierImportState(rName),
				ImportStateVerifyIgnore: []string{"data_source"},
			},
		},
	})
}

func testClassifier_base() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = "%s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testClassifier_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_secmaster_classifier" "test" {
  workspace_id = "%s"
  name         = "%s"
  dataclass_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
  data_source  = "CFW"
  description  = "test description"

  classifier {
    direct_classifier = "false"
  }
}
`, testClassifier_base(), acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testClassifier_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_secmaster_classifier" "test" {
  workspace_id = "%s"
  name         = "%s_update"
  dataclass_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
  data_source  = "CFW"
  description  = "test description update"

  classifier {
    direct_classifier = "true"
  }
}
`, testClassifier_base(), acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testClassifierImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		id := rs.Primary.ID
		if workspaceId == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>',"+
				" but got '%s/%s'", workspaceId, id)
		}

		return fmt.Sprintf("%s/%s", workspaceId, id), nil
	}
}
