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

func getV2ConfigMapResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}

	getConfigMapHttpUrl := "apis/cci/v2/namespaces/{namespace}/configmaps/{name}"
	getConfigMapPath := client.Endpoint + getConfigMapHttpUrl
	getConfigMapPath = strings.ReplaceAll(getConfigMapPath, "{namespaces}", state.Primary.Attributes["namespace"])
	getConfigMapPath = strings.ReplaceAll(getConfigMapPath, "{name}", state.Primary.Attributes["name"])
	getConfigMapOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getConfigMapResp, err := client.Request("GET", getConfigMapPath, &getConfigMapOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getConfigMapResp)
}

func TestAccV2ConfigMap_basic(t *testing.T) {
	var ns namespaces.Namespace
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_config_map.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ns,
		getV2ConfigMapResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2ConfigMap_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "api_version", "cci/v2"),
					resource.TestCheckResourceAttr(resourceName, "kind", "CongfigMap"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
					resource.TestCheckResourceAttrSet(resourceName, "data"),
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

func testAccV2ConfigMap_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cciv2_config_map" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = %[2]s

  data = {
	key   = "xxx"
	value = "xxx"
  }
}
`, testAccV2Namespace_basic(rName), rName)
}
