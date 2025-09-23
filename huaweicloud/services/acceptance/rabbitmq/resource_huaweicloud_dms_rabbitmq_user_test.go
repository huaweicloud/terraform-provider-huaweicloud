package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rabbitmq"
)

func getRabbitmqUserResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return rabbitmq.GetRabbitmqUser(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["access_key"])
}

func TestAccRabbitmqUser_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_user.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRabbitmqUserResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRabbitmqUser_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "access_key", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vhosts.#"),
				),
			},
			{
				Config: testRabbitmqUser_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "access_key", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vhosts.#"),
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

func testRabbitmqUser_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_rabbitmq_user" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  access_key  = "%[2]s"
  secret_key  = "Terraform@123"

  vhosts {
    vhost = "default"
    conf  = ".*"
    write = ".*"
    read  = ".*"
  }
}
`, testAccDmsRabbitmqInstance_amqp_single(rName, true), rName)
}

func testRabbitmqUser_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_rabbitmq_user" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  access_key  = "%[2]s"
  secret_key  = "Terraform@1234"

  vhosts {
    vhost = "default"
    conf  = "a"
    write = "b"
    read  = "c"
  }
}
`, testAccDmsRabbitmqInstance_amqp_single(rName, true), rName)
}
