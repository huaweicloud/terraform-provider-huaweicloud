package swrenterprise

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

func getResourceSwrEnterpriseRepository(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	namespaceName := state.Primary.Attributes["namespace_name"]
	repositoryName := state.Primary.Attributes["repository_name"]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{namespace_name}", namespaceName)
	getPath = strings.ReplaceAll(getPath, "{repository_name}", repositoryName)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func TestAccSwrEnterpriseRepositoryUpdate_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_swr_enterprise_repository_update.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseRepository,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseRepositoryUpdate_basic("test"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
				),
			},
			{
				Config: testAccSwrEnterpriseRepositoryUpdate_basic(""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func testAccSwrEnterpriseRepositoryUpdate_basic(desc string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

resource "huaweicloud_swr_enterprise_repository_update" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = try(data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name, "")
  repository_name = try(data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name, "")
  description     = "%[1]s"

  lifecycle {
    ignore_changes = [
      namespace_name, repository_name,
    ]
  }
}`, desc)
}
