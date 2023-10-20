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

func getEnvironmentResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	getEnvironmentHttpUrl := "v1/environments/{id}"
	getEnvironmentPath := client.Endpoint + getEnvironmentHttpUrl
	getEnvironmentPath = strings.ReplaceAll(getEnvironmentPath, "{id}", state.Primary.ID)

	getEnvironmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getEnvironmentResp, err := client.Request("GET", getEnvironmentPath, &getEnvironmentOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CMDB environment: %s", err)
	}

	getEnvironmentRespBody, err := utils.FlattenResponse(getEnvironmentResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CMDB environment: %s", err)
	}

	return getEnvironmentRespBody, nil
}

func TestAccCmdbEnvironment_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_cmdb_environment.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getEnvironmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesCmdbEnvironment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("env-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "ONLINE"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "LINUX"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "API"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrPair(resourceName, "component_id", "huaweicloud_aom_cmdb_component.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: tesCmdbEnvironment_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("env-%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "DEV"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "LINUX"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "API"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
		},
	})
}

// enterprise_project_id is required for the testing account, we use the default value "0"
func tesCmdbEnvironment_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_cmdb_application" "test" {
  name                  = "app-%[1]s"
  description           = "created by acceptance"
  enterprise_project_id = "0"
}

resource "huaweicloud_aom_cmdb_component" "test" {
  name        = "com-%[1]s"
  model_id    = huaweicloud_aom_cmdb_application.test.id
  model_type  = "APPLICATION"
  description = "created by acceptance"
}`, name)
}

func tesCmdbEnvironment_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_cmdb_environment" "test" {
  component_id = huaweicloud_aom_cmdb_component.test.id
  name         = "env-%s"
  type         = "ONLINE"
  os_type      = "LINUX"
  description  = "created by acceptance"
}`, tesCmdbEnvironment_base(name), name)
}

// update `name` and `type`, clear up `description`
func tesCmdbEnvironment_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_cmdb_environment" "test" {
  component_id = huaweicloud_aom_cmdb_component.test.id
  name         = "env-%s-update"
  type         = "DEV"
  os_type      = "LINUX"
}`, tesCmdbEnvironment_base(name), name)
}
