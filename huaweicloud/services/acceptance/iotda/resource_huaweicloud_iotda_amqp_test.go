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

func getAmqpResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	detail, err := client.ShowQueue(&model.ShowQueueRequest{QueueId: state.Primary.ID})
	// When the queue does not exist, it still returns a empty struct
	if detail == nil || detail.QueueId == nil {
		return nil, fmt.Errorf("error retrieving IoTDA AMQP queue")
	}

	return detail, err
}

func TestAccAmqp_basic(t *testing.T) {
	var obj model.ShowQueueResponse

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_amqp.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAmqpResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAmqp_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
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

func TestAccAmqp_derived(t *testing.T) {
	var obj model.ShowQueueResponse

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_amqp.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAmqpResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAmqp_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
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

func testAmqp_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_amqp" "test" {
  name = "%[2]s"
}
`, buildIoTDAEndpoint(), name)
}
