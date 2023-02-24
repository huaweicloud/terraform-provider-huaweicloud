package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getPublishmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apig.GetVersionHistories(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["env_id"],
		state.Primary.Attributes["api_id"])
}

func TestAccApiPublishment_basic(t *testing.T) {
	var (
		histories []apis.ApiVersionInfo

		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_apig_api_publishment.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&histories,
		getPublishmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApiPublishment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "instance_id",
						"${huaweicloud_apig_instance.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "env_id",
						"${huaweicloud_apig_environment.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "api_id",
						"${huaweicloud_apig_api.test.id}"),
					resource.TestCheckResourceAttrSet(resourceName, "env_name"),
					resource.TestCheckResourceAttrSet(resourceName, "published_at"),
					resource.TestCheckResourceAttrSet(resourceName, "publish_id"),
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

func testAccApiPublishment_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_environment" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_apig_api_publishment" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_id      = huaweicloud_apig_api.test.id
}
`, testAccApi_basic(rName), rName)
}
