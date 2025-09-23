package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getParameterConfigurationsFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	return dws.GetParameterConfigurations(client, state.Primary.ID)
}

// lintignore:AT001
func TestAccParameterConfigurations_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_dws_parameter_configurations.test"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getParameterConfigurationsFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testParameterConfigurations_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "configurations.0.name", "agg_max_mem"),
					resource.TestCheckResourceAttr(rName, "configurations.0.type", "cn"),
					resource.TestCheckResourceAttr(rName, "configurations.0.value", "2097151"),
				),
			},
			{
				Config: testParameterConfigurations_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "configurations.0.name", "ssl_ciphers"),
					resource.TestCheckResourceAttr(rName, "configurations.0.type", "dn"),
					resource.TestCheckResourceAttr(rName, "configurations.0.value", "TLS_CHACHA20_POLY1305_SHA256,TLS_AES_128_GCM_SHA256"),
					resource.TestCheckResourceAttr(rName, "configurations.1.value", "2097152"),
				),
			},
		},
	})
}

func testParameterConfigurations_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_parameter_configurations" "test" {
  cluster_id = "%s"

  configurations {
    name  = "agg_max_mem"
    type  = "cn"
    value = "2097151"
  }
}
`, acceptance.HW_DWS_CLUSTER_ID)
}

func testParameterConfigurations_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_parameter_configurations" "test" {
  cluster_id = "%[1]s"

  configurations {
    name  = "ssl_ciphers"
    type  = "dn"
    value = "TLS_CHACHA20_POLY1305_SHA256,TLS_AES_128_GCM_SHA256"
  }
  configurations {
    name  = "agg_max_mem"
    type  = "cn"
    value = "2097152"
  }
}

resource "huaweicloud_dws_cluster_restart" "test" {
  depends_on = [
    huaweicloud_dws_parameter_configurations.test
  ]

  cluster_id = "%[1]s"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
