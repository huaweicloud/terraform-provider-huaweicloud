package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/iotda"
)

func getBatchTaskFileResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, iotda.WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	resp, err := client.ListBatchTaskFiles(&model.ListBatchTaskFilesRequest{})
	if err != nil {
		return nil, fmt.Errorf("error querying IoTDA batch task files")
	}

	for _, respFile := range *resp.Files {
		if *respFile.FileId == state.Primary.ID {
			return respFile, nil
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccBatchTaskFile_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_iotda_batchtask_file.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBatchTaskFileResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckIOTDABatchTaskFilePath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBatchTaskFile_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content"},
			},
		},
	})
}

func testBatchTaskFile_basic() string {
	return fmt.Sprintf(`

resource "huaweicloud_iotda_batchtask_file" "test" {
  content = "%s"
}

`, acceptance.HW_IOTDA_BATCHTASK_FILE_PATH)
}
