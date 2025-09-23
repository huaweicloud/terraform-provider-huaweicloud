package deprecated

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/iotda"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getBatchTaskFileResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "iotda"
		httpUrl = "v5/iot/{project_id}/batchtask-files"
		ID      = state.Primary.ID
	)

	isDerived := iotda.WithDerivedAuth()
	client, err := conf.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying IoTDA batch task files: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	taskFile := utils.PathSearch(fmt.Sprintf("files[?file_id == '%s']|[0]", ID), respBody, nil)
	if taskFile == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return taskFile, nil
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

// When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify the IoTDA service
// endpoint using environment filed `HW_IOTDA_ACCESS_ADDRESS`.
func buildIoTDAEndpoint() string {
	endpoint := acceptance.HW_IOTDA_ACCESS_ADDRESS
	if endpoint == "" {
		return ""
	}

	// lintignore:AT004
	return fmt.Sprintf(`
provider "huaweicloud" {
  endpoints = {
    iotda = "%s"
  }
}
`, endpoint)
}

func testBatchTaskFile_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iotda_batchtask_file" "test" {
  content = "%s"
}

`, buildIoTDAEndpoint(), acceptance.HW_IOTDA_BATCHTASK_FILE_PATH)
}
