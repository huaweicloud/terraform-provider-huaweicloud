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

func getInstancePluginFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/plugins?plugin_name={plugin_name}"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	pluginName := state.Primary.Attributes["plugin_name"]

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getPath = strings.ReplaceAll(getPath, "{plugin_name}", pluginName)

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

	plugin := utils.PathSearch(fmt.Sprintf("plugins[?plugin_name=='%s' && installed] | [0]", pluginName),
		respBody, nil)
	if plugin == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return plugin, nil
}

func TestAccGaussDbInstancePlugin_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_gaussdb_instance_plugin.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getInstancePluginFunc)
		instanceID   = acceptance.HW_GAUSSDB_INSTANCE_ID
		pluginName   = "postgis"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstancePlugin_basic(instanceID, pluginName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "plugin_name", pluginName),
					resource.TestCheckResourceAttr(resourceName, "url", "https://obs.bucket1"),
					resource.TestCheckResourceAttr(resourceName, "sha_256",
						"791a8d68064ca3208b52ac2584b3b1ab89e4945069baf48e2b14ed5a7151889b"),
					resource.TestCheckResourceAttrSet(resourceName, "installed"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "plugin_version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"url",
					"sha_256",
				},
			},
		},
	})
}

func testAccGaussDbInstancePlugin_basic(instanceID, pluginName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_instance_plugin" "test" {
  instance_id = "%[1]s"
  plugin_name = "%[2]s"
  url         = "https://obs.bucket1"
  sha_256     = "791a8d68064ca3208b52ac2584b3b1ab89e4945069baf48e2b14ed5a7151889b"
}
`, instanceID, pluginName)
}
