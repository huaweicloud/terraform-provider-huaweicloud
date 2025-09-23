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

func getResourceSwrEnterpriseNamespace(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{namespace_name}", state.Primary.Attributes["name"])
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

func TestAccSwrEnterpriseNamespace_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_namespace.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseNamespace,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseNamespace_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.public", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccSwrEnterpriseNamespace_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.public", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
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

func testAccSwrEnterpriseNamespace_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_enterprise_namespace" "test" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  name        = "%s"

  metadata {
    public = "true"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccSwrEnterpriseInstance_update(rName), rName)
}

func testAccSwrEnterpriseNamespace_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_enterprise_namespace" "test" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  name        = "%s"

  metadata {
    public = "false"
  }

  tags = {
    key = "value1"
    foo = "bar1"
  }
}
`, testAccSwrEnterpriseInstance_update(rName), rName)
}
