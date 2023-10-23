package aom

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

func getComponentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	getComponentHttpUrl := "v1/components/{id}"
	getComponentPath := client.Endpoint + getComponentHttpUrl
	getComponentPath = strings.ReplaceAll(getComponentPath, "{id}", state.Primary.ID)

	getComponentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getComponentResp, err := client.Request("GET", getComponentPath, &getComponentOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CMDB component: %s", err)
	}

	getComponentRespBody, err := utils.FlattenResponse(getComponentResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CMDB component: %s", err)
	}

	return getComponentRespBody, nil
}

func TestAccCmdbComponent_basic(t *testing.T) {
	var obj interface{}

	appName := acceptance.RandomAccResourceName()
	comName := acceptance.RandomAccResourceName()
	newName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_cmdb_component.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getComponentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCmdbComponent_basic(appName, comName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", comName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance"),
					resource.TestCheckResourceAttr(resourceName, "model_type", "APPLICATION"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "API"),
					resource.TestCheckResourceAttr(resourceName, "sub_app_id", ""),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrPair(resourceName, "model_id", "huaweicloud_aom_cmdb_application.test", "id"),
				),
			},
			{
				Config: testCmdbComponent_basic(appName, newName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", newName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance"),
					resource.TestCheckResourceAttr(resourceName, "model_type", "APPLICATION"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testCmdbComponent_update(appName, comName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", comName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "model_type", "APPLICATION"),
				),
			},
		},
	})
}

func TestAccCmdbComponent_sub_application(t *testing.T) {
	var obj interface{}

	subAppID := acceptance.HW_AOM_SUB_APPLICATION_ID
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_cmdb_component.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getComponentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAomSubApplicationId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCmdbComponent_sub_application(name, subAppID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "model_id", subAppID),
					resource.TestCheckResourceAttr(resourceName, "model_type", "SUB_APPLICATION"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "API"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCmdbComponent_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_cmdb_application" "test" {
  name                  = "app-%s"
  description           = "created by acceptance"
  enterprise_project_id = "0"
}`, name)
}

func testCmdbComponent_basic(appName, comName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_cmdb_component" "test" {
  name        = "%s"
  model_id    = huaweicloud_aom_cmdb_application.test.id
  model_type  = "APPLICATION"
  description = "created by acceptance"
}
`, testCmdbComponent_base(appName), comName)
}

// testCmdbComponent_update will clear description field
func testCmdbComponent_update(appName, comName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_cmdb_component" "test" {
  name       = "%s"
  model_id   = huaweicloud_aom_cmdb_application.test.id
  model_type = "APPLICATION"
}
`, testCmdbComponent_base(appName), comName)
}

func testCmdbComponent_sub_application(name, subAppID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_cmdb_component" "test" {
  name        = "%s"
  model_id    = "%s"
  model_type  = "SUB_APPLICATION"
  description = "created by acceptance"
}
`, name, subAppID)
}
