package codearts

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDeployGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/host-groups/{group_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{group_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CodeArts deploy group: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	errorCode := utils.PathSearch("error_code", getRespBody, "")
	if errorCode == "Deploy.00021104" {
		// 'Deploy.00021104' means the group is not exist
		return nil, golangsdk.ErrDefault404{}
	}

	if errorCode != "" {
		errorMsg := utils.PathSearch("error_msg", getRespBody, "")
		return nil, fmt.Errorf("error retrieving CodeArts deploy group: error code: %s, error message: %s",
			errorCode, errorMsg)
	}

	return getRespBody, nil
}

func TestAccDeployGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "os_type", "linux"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "is_proxy_mode", "1"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "host_count"),
					resource.TestCheckResourceAttrSet(rName, "created_by.#"),
					resource.TestCheckResourceAttrSet(rName, "updated_by.#"),
					resource.TestCheckResourceAttrSet(rName, "permission.#"),
				),
			},
			{
				Config: testDeployGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"is_proxy_mode",
				},
			},
		},
	})
}

func TestAccDeployGroup_resourcePoolId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCodeArtsDeployResourcePoolID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployGroup_resourcePoolId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "os_type", "windows"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "is_proxy_mode", "0"),
					resource.TestCheckResourceAttr(rName, "resource_pool_id", acceptance.HW_CODEARTS_RESOURCE_POOL_ID),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "host_count"),
					resource.TestCheckResourceAttrSet(rName, "created_by.#"),
					resource.TestCheckResourceAttrSet(rName, "updated_by.#"),
					resource.TestCheckResourceAttrSet(rName, "permission.#"),
				),
			},
			{
				Config: testDeployGroup_resourcePoolId_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "resource_pool_id", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"is_proxy_mode",
				},
			},
		},
	})
}

func TestAccDeployGroup_errorCheck(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testDeployGroup_errorCheck(name),
				ExpectError: regexp.MustCompile(`error creating CodeArts deploy group: error code:`),
			},
		},
	})
}

func testDeployGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%s"
  os_type     = "linux"
  description = "test description"
}
`, testProject_basic(name), name)
}

func testDeployGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%s_update"
  os_type     = "linux"
  description = "test description update"
}
`, testProject_basic(name), name)
}

func testDeployGroup_resourcePoolId(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id       = huaweicloud_codearts_project.test.id
  name             = "%s"
  os_type          = "windows"
  description      = "test description"
  is_proxy_mode    = 0
  resource_pool_id = "%s"
}
`, testProject_basic(name), name, acceptance.HW_CODEARTS_RESOURCE_POOL_ID)
}

func testDeployGroup_resourcePoolId_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id    = huaweicloud_codearts_project.test.id
  name          = "%s_update"
  os_type       = "windows"
  description   = "test description update"
  is_proxy_mode = 0
}
`, testProject_basic(name), name)
}

func testDeployGroup_errorCheck(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_group" "test" {
  project_id    = huaweicloud_codearts_project.test.id
  name          = "%s"
  os_type       = "error_type"
  description   = "test description"
  is_proxy_mode = 0
}
`, testProject_basic(name), name)
}
