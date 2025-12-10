package swr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/swr/v2/repositories"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceRepository(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	repositoryName := strings.ReplaceAll(state.Primary.ID, "/", "$")

	return repositories.Get(swrClient, state.Primary.Attributes["organization"], repositoryName).Extract()
}

func TestAccSWRRepository_basic(t *testing.T) {
	var repo repositories.ImageRepository
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_swr_repository.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&repo,
		getResourceRepository,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWRRepository_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "category", "linux"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
				),
			},
			{
				Config: testAccSWRRepository_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "category", "windows"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSWRRepositoryImportStateIdFunc(),
			},
		},
	})
}

func testAccSWRRepositoryImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var organization string
		var repositoryID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_swr_organization" {
				organization = rs.Primary.Attributes["name"]
			} else if rs.Type == "huaweicloud_swr_repository" {
				repositoryID = rs.Primary.ID
			}
		}
		if organization == "" || repositoryID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", organization, repositoryID)
		}
		return fmt.Sprintf("%s/%s", organization, repositoryID), nil
	}
}

func testAccSWRRepository_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_repository" "test" {
  organization = huaweicloud_swr_organization.test.name
  name         = "%s"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}
`, testAccSWROrganization_basic(rName), rName)
}

func testAccSWRRepository_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_repository" "test" {
  organization = huaweicloud_swr_organization.test.name
  name         = "%s"
  description  = "Test repository"
  category     = "windows"
  is_public    = true
}
`, testAccSWROrganization_basic(rName), rName)
}

func TestAccSWRRepository_withcomma(t *testing.T) {
	var repo repositories.ImageRepository
	rName := acceptance.RandomAccResourceName()
	slashName := strings.ReplaceAll(rName, "_", "/")
	resourceName := "huaweicloud_swr_repository.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&repo,
		getResourceRepository,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWRRepository_withcomma(rName, slashName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "name", slashName),
					resource.TestCheckResourceAttr(resourceName, "category", "linux"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
				),
			},
			{
				Config: testAccSWRRepository_withcomma_update(rName, slashName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "name", slashName),
					resource.TestCheckResourceAttr(resourceName, "category", "windows"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSWRRepositoryImportStateIdFuncWithComma(),
			},
		},
	})
}

func testAccSWRRepository_withcomma(rName, slashName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_repository" "test" {
  organization = huaweicloud_swr_organization.test.name
  name         = "%s"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}
`, testAccSWROrganization_basic(rName), slashName)
}

func testAccSWRRepository_withcomma_update(rName, slashName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_repository" "test" {
  organization = huaweicloud_swr_organization.test.name
  name         = "%s"
  description  = "Test repository"
  category     = "windows"
  is_public    = true
}
`, testAccSWROrganization_basic(rName), slashName)
}

func testAccSWRRepositoryImportStateIdFuncWithComma() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var organization string
		var repositoryID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_swr_organization" {
				organization = rs.Primary.Attributes["name"]
			} else if rs.Type == "huaweicloud_swr_repository" {
				repositoryID = rs.Primary.ID
			}
		}
		if organization == "" || repositoryID == "" {
			return "", fmt.Errorf("resource not found: %s,%s", organization, repositoryID)
		}
		return fmt.Sprintf("%s,%s", organization, repositoryID), nil
	}
}

func TestAccSWRRepository_noSlashButWithcomma(t *testing.T) {
	var repo repositories.ImageRepository
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_swr_repository.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&repo,
		getResourceRepository,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWRRepository_noSlashButWithcomma(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "category", "linux"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
				),
			},
			{
				Config: testAccSWRRepository_noSlashButWithcomma_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "category", "windows"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSWRRepositoryImportStateIdFuncWithComma(),
			},
		},
	})
}

func testAccSWRRepository_noSlashButWithcomma(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_repository" "test" {
  organization = huaweicloud_swr_organization.test.name
  name         = "%s"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}
`, testAccSWROrganization_basic(rName), rName)
}

func testAccSWRRepository_noSlashButWithcomma_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_repository" "test" {
  organization = huaweicloud_swr_organization.test.name
  name         = "%s"
  description  = "Test repository"
  category     = "windows"
  is_public    = true
}
`, testAccSWROrganization_basic(rName), rName)
}
