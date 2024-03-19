package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getSpaceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	return client.ShowApplication(&model.ShowApplicationRequest{AppId: state.Primary.ID})
}

func TestAccSpace_basic(t *testing.T) {
	var obj model.ShowApplicationRequest

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_space.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSpaceResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSpace_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "is_default", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpace_derived(t *testing.T) {
	var obj model.ShowApplicationRequest

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_space.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSpaceResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSpace_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "is_default", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testSpace_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iotda_space" "test" {
  name = "%s"
}
`, buildIoTDAEndpoint(), name)
}
