package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/identitycenter"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAccessControlAttributeConfigurationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center client: %s", err)
	}

	resp, err := identitycenter.GetAccessControlAttributeConfiguration(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", resp, "").(string)
	if status != "ENABLED" {
		return nil, golangsdk.ErrDefault404{}
	}

	return resp, nil
}

func TestAccAccessControlAttributeConfiguration_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_access_control_attribute_configuration.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccessControlAttributeConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccessControlAttributeConfiguration_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "data.huaweicloud_identitycenter_instance.system", "id"),
					resource.TestCheckResourceAttr(rName, "access_control_attributes.#", "1"),
					resource.TestCheckResourceAttr(rName, "access_control_attributes.0.key", name+"_1"),
					resource.TestCheckResourceAttr(rName, "access_control_attributes.0.value.0", "${user:email}"),
				),
			},
			{
				Config: testAccessControlAttributeConfiguration_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "instance_id", "data.huaweicloud_identitycenter_instance.system", "id"),
					resource.TestCheckResourceAttr(rName, "access_control_attributes.#", "2"),
					resource.TestCheckResourceAttr(rName, "access_control_attributes.0.key", name+"_1"),
					resource.TestCheckResourceAttr(rName, "access_control_attributes.0.value.0", "${user:email}"),
					resource.TestCheckResourceAttr(rName, "access_control_attributes.1.key", name+"_2"),
					resource.TestCheckResourceAttr(rName, "access_control_attributes.1.value.0", "${user:familyName}"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func testAccessControlAttributeConfiguration_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_access_control_attribute_configuration" "test" {
  instance_id = data.huaweicloud_identitycenter_instance.system.id

  access_control_attributes {
    key   = "%[1]s_1"
    value = ["$${user:email}"]
  }
}
`, name)
}

func testAccessControlAttributeConfiguration_update(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_access_control_attribute_configuration" "test" {
  instance_id = data.huaweicloud_identitycenter_instance.system.id

  access_control_attributes {
    key   = "%[1]s_1"
    value = ["$${user:email}"]
  }

  access_control_attributes {
    key   = "%[1]s_2"
    value = ["$${user:familyName}"]
  }
}
`, name)
}
