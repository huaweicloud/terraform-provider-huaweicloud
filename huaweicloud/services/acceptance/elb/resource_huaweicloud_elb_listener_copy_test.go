package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
)

func getResourceListenerCopyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("elb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	return elb.GetListenerCopy(client, state.Primary.ID)
}

func TestAccResourceListenerCopy_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_elb_listener_copy.test"
		name         = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			resourceName,
			&object,
			getResourceListenerCopyFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbListenerId(t)
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccListenerCopy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "listener_id", acceptance.HW_ELB_LISTENER_ID),
					resource.TestCheckResourceAttr(resourceName, "loadbalancer_id", acceptance.HW_ELB_LOADBALANCER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "8008"),
					resource.TestCheckResourceAttr(resourceName, "reuse_pool", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "http2_enable"),
					resource.TestCheckResourceAttrSet(resourceName, "default_pool_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"listener_id",
					"reuse_pool",
					"force_delete",
				},
			},
		},
	})
}

func testAccListenerCopy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_listener_copy" "test" {
  listener_id     = "%[1]s" 
  loadbalancer_id = "%[2]s"
  name            = "%[3]s"
  protocol_port   = "8008"
  reuse_pool      = false     

  force_delete = true
}
`, acceptance.HW_ELB_LISTENER_ID, acceptance.HW_ELB_LOADBALANCER_ID, name)
}
