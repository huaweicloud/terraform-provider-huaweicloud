package rabbitmq

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

func getRabbitmqPluginResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getRabbitmqPluginClient, err := cfg.NewServiceClient("dms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceID := parts[0]
	name := parts[1]
	getRabbitmqPluginHttpUrl := "v2/{project_id}/instances/{instance_id}/rabbitmq/plugins"
	getRabbitmqPluginPath := getRabbitmqPluginClient.Endpoint + getRabbitmqPluginHttpUrl
	getRabbitmqPluginPath = strings.ReplaceAll(getRabbitmqPluginPath, "{project_id}", getRabbitmqPluginClient.ProjectID)
	getRabbitmqPluginPath = strings.ReplaceAll(getRabbitmqPluginPath, "{instance_id}", instanceID)

	resp, err := getRabbitmqPluginClient.Request("GET", getRabbitmqPluginPath, &golangsdk.RequestOpts{
		KeepResponseBody: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error retrieving the plugin : %s", err)
	}

	body, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing the plugin : %s", err)
	}

	enable := utils.PathSearch(fmt.Sprintf("plugins|[?name=='%s']|[0].enable", name), body, false).(bool)
	if !enable {
		return nil, fmt.Errorf("the plugin %s is disabled", name)
	}

	return body, nil
}

func TestAccRabbitmqPlugin_basic(t *testing.T) {
	var obj interface{}
	randName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_plugin.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqPluginResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqPlugin_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "rabbitmq_consistent_hash_exchange"),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "running"),
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

func testRabbitmqPlugin_basic(randName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_plugin" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  name        = "rabbitmq_consistent_hash_exchange"
}
`, testAccDmsRabbitmqInstance_newFormat_single(randName))
}
