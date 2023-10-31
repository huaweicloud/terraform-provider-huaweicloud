package codearts

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

func getRepositoryResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getRepository: Query the resource detail of the CodeHub repository
	var (
		getRepositoryHttpUrl = "v2/repositories/{repository_uuid}"
		getRepositoryProduct = "codehub"
	)
	getRepositoryClient, err := cfg.NewServiceClient(getRepositoryProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating repository client: %s", err)
	}

	getRepositoryPath := getRepositoryClient.Endpoint + getRepositoryHttpUrl
	getRepositoryPath = strings.ReplaceAll(getRepositoryPath, "{repository_uuid}", state.Primary.ID)

	getRepositoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRepositoryResp, err := getRepositoryClient.Request("GET", getRepositoryPath, &getRepositoryOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CodeHub repository: %s", err)
	}
	return utils.FlattenResponse(getRepositoryResp)
}

func TestAccRepository_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_repository.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRepositoryResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRepository_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test",
						"id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform acc test"),
					resource.TestCheckResourceAttr(rName, "gitignore_id", "Go"),
					resource.TestCheckResourceAttr(rName, "enable_readme", "0"),
					resource.TestCheckResourceAttr(rName, "visibility_level", "20"),
					resource.TestCheckResourceAttr(rName, "license_id", "2"),
					resource.TestCheckResourceAttr(rName, "import_members", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name",
					"description",
					"gitignore_id",
					"enable_readme",
					"license_id",
					"import_members",
				},
			},
		},
	})
}

func TestAccRepository_default(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_repository.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRepositoryResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRepository_default(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test",
						"id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enable_readme", "1"),
					resource.TestCheckResourceAttr(rName, "visibility_level", "0"),
					resource.TestCheckResourceAttr(rName, "license_id", "1"),
					resource.TestCheckResourceAttr(rName, "import_members", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name",
					"description",
					"gitignore_id",
					"enable_readme",
					"license_id",
					"import_members",
				},
			},
		},
	})
}

func testRepository_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_project" "test" {
  name = "%[1]s"
  type = "scrum"
}

resource "huaweicloud_codearts_repository" "test" {
  project_id = huaweicloud_codearts_project.test.id

  name             = "%[1]s"
  description      = "Created by terraform acc test"
  gitignore_id     = "Go"
  enable_readme    = 0  // Do not generate README.md
  visibility_level = 20 // Public read-only
  license_id       = 2  // MIT License
  import_members   = 0  // Do not import members of the project into this repository when creating
}
`, name)
}

func testRepository_default(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_project" "test" {
  name = "%[1]s"
  type = "scrum"
}

resource "huaweicloud_codearts_repository" "test" {
  project_id = huaweicloud_codearts_project.test.id

  name = "%[1]s"
}
`, name)
}
