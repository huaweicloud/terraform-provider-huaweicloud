package cci

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cci/v1/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getV2SecretResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}

	getSecretHttpUrl := "apis/cci/v2/namespaces/{namespace}/secret/{name}"
	getSecretPath := client.Endpoint + getSecretHttpUrl
	getSecretPath = strings.ReplaceAll(getSecretPath, "{namespace}", state.Primary.Attributes["namespace"])
	getSecretPath = strings.ReplaceAll(getSecretPath, "{name}", state.Primary.Attributes["name"])
	getSecretOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getSecretResp, err := client.Request("GET", getSecretPath, &getSecretOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getSecretResp)
}

func TestAccV2Secret_basic(t *testing.T) {
	var ns namespaces.Namespace
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_secret.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ns,
		getV2SecretResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Secret_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "api_version", "cci/v2"),
					resource.TestCheckResourceAttr(resourceName, "kind", "Secret"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
					resource.TestCheckResourceAttrSet(resourceName, "data.%"),
					resource.TestCheckOutput("dockerconfigjson_verify", "true"),
				),
			},
			{
				Config: testAccV2Secret_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckOutput("dockerconfigjson_verify", "true"),
					resource.TestCheckOutput("expired_at_verify", "true"),
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

func testAccV2Secret_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cciv2_secret" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[2]s"

  data = {
    ".dockerconfigjson" = "%[3]s"
  }
}

output "dockerconfigjson_verify" {
  value = huaweicloud_cciv2_secret.test.data[".dockerconfigjson"] == "%[3]s"
}
`, testAccV2Namespace_basic(rName), rName, acceptance.HW_CCI_SECRET_DOCKERCONFIGJSON)
}

func testAccV2Secret_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cciv2_secret" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[2]s"

  data = {
    ".dockerconfigjson" = "%[3]s"
    "expired.at"        = "MjAyNS0wNC0xNlQwNTo1NzowMVo="
  }
}

output "dockerconfigjson_verify" {
  value = huaweicloud_cciv2_secret.test.data[".dockerconfigjson"] == "%[3]s"
}

output "expired_at_verify" {
  value = huaweicloud_cciv2_secret.test.data["expired.at"] == "MjAyNS0wNC0xNlQwNTo1NzowMVo="
}
`, testAccV2Namespace_basic(rName), rName, acceptance.HW_CCI_SECRET_DOCKERCONFIGJSON)
}
