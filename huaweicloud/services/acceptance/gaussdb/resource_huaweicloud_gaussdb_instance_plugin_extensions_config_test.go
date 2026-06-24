package gaussdb

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

func getInstancePluginExtensionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/plugin-extensions"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	pluginName := state.Primary.Attributes["plugin_name"]
	dbName := state.Primary.Attributes["db_name"]
	extensionName := state.Primary.Attributes["extension_name"]

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getPath += fmt.Sprintf("?plugin_name=%s&db_name=%s", pluginName, dbName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	extension := utils.PathSearch(
		fmt.Sprintf("[?extension_name=='%s'] | [0]", extensionName), respBody, nil)
	if extension == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	status := utils.PathSearch("status", extension, "").(string)
	if status == "" || status == "off" {
		return nil, golangsdk.ErrDefault404{}
	}

	return extension, nil
}

func TestAccGaussDbInstancePluginExtensionsConfig_basic(t *testing.T) {
	var (
		obj           interface{}
		resourceName  = "huaweicloud_gaussdb_instance_plugin_extensions_config.test"
		rc            = acceptance.InitResourceCheck(resourceName, &obj, getInstancePluginExtensionFunc)
		instanceID    = acceptance.HW_GAUSSDB_INSTANCE_ID
		dbName        = "gauss-6ead"
		pluginName    = "postgis"
		extensionName = "postgis-raster"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstancePluginExtensionsConfig_basic(instanceID, dbName, pluginName,
					extensionName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "db_name", dbName),
					resource.TestCheckResourceAttr(resourceName, "plugin_name", pluginName),
					resource.TestCheckResourceAttr(resourceName, "extension_name", extensionName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
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

func testAccGaussDbInstancePluginExtensionsConfig_basic(instanceID, dbName, pluginName, extensionName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_instance_plugin_extensions_config" "test" {
  instance_id      = "%[1]s"
  db_name          = "%[2]s"
  plugin_name      = "%[3]s"
  extension_name   = "%[4]s"
}
`, instanceID, dbName, pluginName, extensionName)
}
