package apig

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getPublishmentResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud APIG v2 client: %s", err)
	}
	return apig.GetVersionHistories(c, state.Primary.Attributes["instance_id"], state.Primary.Attributes["env_id"],
		state.Primary.Attributes["api_id"])
}

func TestAccApigApiPublishmentV2_basic(t *testing.T) {
	var histories []apis.ApiVersionInfo

	// The dedicated instance name only allow letters, digits and underscores (_).
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_apig_api_publishment.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&histories,
		getPublishmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApigApiPublishment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "instance_id",
						"${huaweicloud_apig_instance.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "env_id",
						"${huaweicloud_apig_environment.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "api_id",
						"${huaweicloud_apig_api.test.id}"),
					resource.TestCheckResourceAttrSet(resourceName, "env_name"),
					resource.TestCheckResourceAttrSet(resourceName, "publish_time"),
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

func testAccApigApiPublishment_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_environment" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}

resource "huaweicloud_apig_api_publishment" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_id      = huaweicloud_apig_api.test.id
}
`, testAccApigAPI_basic(rName), rName)
}
