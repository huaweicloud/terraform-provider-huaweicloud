package swr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceRepositoryTag(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	organization := state.Primary.Attributes["organization"]
	repository := strings.ReplaceAll(state.Primary.Attributes["repository"], "/", "$")
	tag := state.Primary.Attributes["tag"]

	getHttpUrl := "v2/manage/namespaces/{namespace}/repos/{repository}/tags/{tag}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{namespace}", organization)
	getPath = strings.ReplaceAll(getPath, "{repository}", repository)
	getPath = strings.ReplaceAll(getPath, "{tag}", tag)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SWR repository tag: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func TestAccSWRRepositoryTag_basic(t *testing.T) {
	var v interface{}
	resourceName := "huaweicloud_swr_repository_tag.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&v,
		getResourceRepositoryTag,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrRepository(t)
			acceptance.TestAccPreCheckSwrOrigination(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWRRepositoryTag_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "organization", acceptance.HW_SWR_ORGANIZATION),
					resource.TestCheckResourceAttr(resourceName, "repository", acceptance.HW_SWR_REPOSITORY),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_tag", "destination_tag", "override"},
			},
		},
	})
}

func testAccSWRRepositoryTag_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_swr_image_tags" "test" {
  organization = "%[1]s"
  repository   = "%[2]s"
}

resource "huaweicloud_swr_repository_tag" "test" {
  organization    = "%[1]s"
  repository      = "%[2]s"
  source_tag      = data.huaweicloud_swr_image_tags.test.image_tags[0].name
  destination_tag = "new"
  override        = true

  lifecycle {
    ignore_changes = [
      source_tag,
    ]
  }
}
`, acceptance.HW_SWR_ORGANIZATION, acceptance.HW_SWR_REPOSITORY)
}
