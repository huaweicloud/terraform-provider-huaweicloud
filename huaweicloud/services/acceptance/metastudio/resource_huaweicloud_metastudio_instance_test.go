package metastudio

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/metastudio"
)

func getMetaStudioResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("metastudio", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating MetaStudio client: %s", err)
	}
	return metastudio.GetResourceDetail(client, state.Primary.ID)
}

func TestAccResourceMetaStudio_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_metastudio_instance.test"
		instance     interface{}
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&instance,
			getMetaStudioResourceFunc,
		)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceMetaStudio_base(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "order_id"),
				),
			},
		},
	})
}

func testResourceMetaStudio_base() string {
	return `
resource "huaweicloud_metastudio_instance" "test" {
  period_type        = 2
  period_num         = 1
  is_auto_renew      = 0
  resource_spec_code = "hws.resource.type.metastudio.modeling.avatarlive.channel"
}
`
}
