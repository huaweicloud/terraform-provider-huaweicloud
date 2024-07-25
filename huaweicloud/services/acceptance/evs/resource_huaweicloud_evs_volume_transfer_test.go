package evs

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getVolumeTransferFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region                   = acceptance.HW_REGION_NAME
		id                       = state.Primary.ID
		getVolumeTransferHttpUrl = "v2/{project_id}/os-volume-transfer/{transfer_id}"
		product                  = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	getVolumeTransferPath := client.Endpoint + getVolumeTransferHttpUrl
	getVolumeTransferPath = strings.ReplaceAll(getVolumeTransferPath, "{project_id}", client.ProjectID)
	getVolumeTransferPath = strings.ReplaceAll(getVolumeTransferPath, "{transfer_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResponse, err := client.Request("GET", getVolumeTransferPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving EVS volume transfer: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResponse)
	if err != nil {
		return nil, err
	}

	transfer := utils.PathSearch("transfer", getRespBody, nil)
	return transfer, nil
}
func TestAccVolumeTransfer_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_evs_volume_transfer.test"
		name  = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getVolumeTransferFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeTransfer_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "volume_id", "huaweicloud_evs_volume.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "auth_key"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccVolumeTransfer_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  description       = "Created by acc test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}
`, name)
}

func testAccVolumeTransfer_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume_transfer" "test" {
  volume_id = huaweicloud_evs_volume.test.id
  name      = "%[2]s"
}
`, testAccVolumeTransfer_base(name), name)
}
