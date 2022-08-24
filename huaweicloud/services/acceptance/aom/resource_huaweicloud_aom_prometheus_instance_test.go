package aom

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPrometheusInstanceResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, _ := httpclient_go.NewHttpClientGo(config)

	client.WithMethod(httpclient_go.MethodGet).WithUrlWithoutEndpoint(config, "aom",
		config.Region, "v1/"+state.Primary.Attributes["project_id"]+"/prometheus-instances?action=prom_for_cloud_service")

	resp, err := client.Do()

	mErr := &multierror.Error{}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	rlt := &entity.PrometheusInstanceParams{}

	err = json.Unmarshal(body, rlt)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	return nil, fmt.Errorf("error getting HuaweiCloud Resource")
}

func TestAccAOMPrometheusInstance_basic(t *testing.T) {
	var ar []entity.AddAlarmRuleParams
	resourceName := "huaweicloud_aom_prometheus_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getPrometheusInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckInternal(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMPrometheusInstance_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAOMPrometheusInstance_basic() string {
	return fmt.Sprintf(`
provider "huaweicloud" {
  region     = "cn-north-7"
  access_key = "KKILFTQC8DJYDINQQZXI"
  secret_key = "yM9DG7GaISH9Ob2n6zrq89IvdiZC68keqOad9oVu"
  auth_url = "https://iam.cn-north-7.myhuaweicloud.com"
  endpoints = {
    aom : "aom.cn-north-7.myhuaweicloud.com"
  }

  insecure = true
  domain_id = "40de487942a74a70b4666fa32d11ffa8"
  project_id = "2a473356cca5487f8373be891bffc1cf"
}


resource "huaweicloud_aom_prometheus_instance" "test" {
  project_id = "2a473356cca5487f8373be891bffc1cf"
  action         = "prom_for_cloud_service"
  prom_for_cloud_service  {
  ces_metric_namespaces =	["SYS.ELB","SYS.VPC","SYS.DMS","SYS.RDS"]
 }
}
`)
}
