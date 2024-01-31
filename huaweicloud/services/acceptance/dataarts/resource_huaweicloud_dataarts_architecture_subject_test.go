package dataarts

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

func getArchitectureSubjectResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getSubjectClient, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	workspaceID := state.Primary.Attributes["workspace_id"]
	getSubjectHttpUrl := "v3/{project_id}/design/subjects"

	getSubjectPath := getSubjectClient.Endpoint + getSubjectHttpUrl
	getSubjectPath = strings.ReplaceAll(getSubjectPath, "{project_id}", getSubjectClient.ProjectID)

	getSubjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}

	getSubjectPath += "?limit=10"

	// using name to reduce the search results
	if val, ok := state.Primary.Attributes["name"]; ok && val != "" {
		getSubjectPath += fmt.Sprintf("&name=%s", val)
	}

	currentTotal := 0
	for {
		path := fmt.Sprintf("%s&offset=%v", getSubjectPath, currentTotal)
		getSubjectResp, err := getSubjectClient.Request("GET", path, &getSubjectOpt)
		if err != nil {
			return nil, err
		}
		getSubjectRespBody, err := utils.FlattenResponse(getSubjectResp)
		if err != nil {
			return nil, err
		}
		subjects := utils.PathSearch("data.value.records", getSubjectRespBody, make([]interface{}, 0)).([]interface{})
		total := utils.PathSearch("data.value.total", getSubjectRespBody, 0)
		for _, subject := range subjects {
			// using path to filter result for import, because ID can not be got from console
			// format of path in results using `/` to split, format of path from import using `.` to split
			id := utils.PathSearch("id", subject, "")
			path := strings.ReplaceAll(utils.PathSearch("path", subject, "").(string), "/", ".")
			if val, ok := state.Primary.Attributes["path"]; ok {
				if path != val {
					continue
				}
			} else if id != state.Primary.ID {
				continue
			}
			return subject, nil
		}
		currentTotal += len(subjects)
		// type of `total` is float64
		if float64(currentTotal) == total {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccResourceArchitectureSubject_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_architecture_subject.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getArchitectureSubjectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureSubject_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "code", rName),
					resource.TestCheckResourceAttr(resourceName, "owner", rName),
					resource.TestCheckResourceAttr(resourceName, "level", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "create"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccArchitectureSubject_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "code", rName),
					resource.TestCheckResourceAttr(resourceName, "owner", rName),
					resource.TestCheckResourceAttr(resourceName, "level", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "update"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceArchitectureSubjectImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccResourceArchitectureSubjectImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		path := rs.Primary.Attributes["path"]
		if workspaceID == "" || path == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<path>, but got '%s/%s'",
				workspaceID, path)
		}
		return fmt.Sprintf("%s/%s", workspaceID, path), nil
	}
}

func testAccArchitectureSubject_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_subject" "test" {
  workspace_id = "%s"
  name         = "%s"
  code         = "%s"
  owner        = "%s"
  level        = 1
  description  = "create"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, name, name)
}

func testAccArchitectureSubject_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_subject" "test" {
  workspace_id = "%s"
  name         = "%s"
  code         = "%s"
  owner        = "%s"
  level        = 1
  description  = "update"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, name, name)
}
