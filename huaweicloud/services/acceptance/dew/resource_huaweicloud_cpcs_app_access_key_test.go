package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
)

func getCpcsAppAccessKeyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("kms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DEW client: %s", err)
	}

	appId := state.Primary.Attributes["app_id"]
	keyName := state.Primary.Attributes["key_name"]
	return dew.QueryCpcsAppAccessKeyByAppIdAndKeyName(client, appId, keyName)
}

// Currently, this resource is valid only in cn-north-9 region.
func TestAccCpcsAppAccessKey_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_cpcs_app_access_key.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCpcsAppAccessKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCpcsAppAccessKey_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "key_name", name),
					resource.TestCheckResourceAttr(rName, "status", "disable"),
					resource.TestCheckResourceAttrSet(rName, "app_id"),
					resource.TestCheckResourceAttrSet(rName, "access_key"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "download_time"),
					resource.TestCheckResourceAttrSet(rName, "is_downloaded"),
					resource.TestCheckResourceAttrSet(rName, "is_imported"),
				),
			},
			{
				Config: testCpcsAppAccessKey_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "enable"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCpcsAppAccessKeyImportState(rName),
			},
		},
	})
}

func testCpcsAppAccessKey_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cpcs_app_access_key" "test" {
  app_id   = huaweicloud_cpcs_app.test.id
  key_name = "%[2]s"
  status   = "disable"
}
`, testCpcsApp_basic(name), name)
}

func testCpcsAppAccessKey_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cpcs_app_access_key" "test" {
  app_id   = huaweicloud_cpcs_app.test.id
  key_name = "%[2]s"
  status   = "enable"
}
`, testCpcsApp_basic(name), name)
}

func testCpcsAppAccessKeyImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		if rs.Primary.Attributes["app_id"] == "" {
			return "", fmt.Errorf("attribute (app_id) of resource (%s) not found", name)
		}
		if rs.Primary.Attributes["key_name"] == "" {
			return "", fmt.Errorf("attribute (key_name) of resource (%s) not found", name)
		}

		appId := rs.Primary.Attributes["app_id"]
		keyName := rs.Primary.Attributes["key_name"]
		return fmt.Sprintf("%s/%s", appId, keyName), nil
	}
}
