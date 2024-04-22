package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dcs"
)

func getRestoreRecordResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getAccountProduct = "dcs"
	)
	getRestoreRecordClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	restoreRecord, err := dcs.GetRestoreRecord(state.Primary.Attributes["instance_id"], state.Primary.ID, getRestoreRecordClient)
	if err != nil {
		return nil, err
	}

	if restoreRecord == nil {
		return nil, fmt.Errorf("the restoration record (%s) is not found", state.Primary.ID)
	}

	return restoreRecord, nil
}

func TestAccDcsRestore_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance_restore.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRestoreRecordResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsRestore_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_dcs_instance.instance_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "backup_id", "huaweicloud_dcs_backup.test", "backup_id"),
					resource.TestCheckResourceAttr(rName, "description", "test DCS restoration"),
					resource.TestCheckResourceAttrSet(rName, "restore_name"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
		},
	})
}

func testDcsRestore_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_backup" "test" {
  instance_id   = huaweicloud_dcs_instance.instance_1.id
  description   = "test DCS backup"
  backup_format = "rdb"
}

resource "huaweicloud_dcs_instance_restore" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  backup_id   = huaweicloud_dcs_backup.test.backup_id
  description = "test DCS restoration"
}
`, testAccDcsV1Instance_basic(name))
}
