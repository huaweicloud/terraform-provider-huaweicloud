package apm

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"testing"
)

var c *httpclient_go.HttpClientGo

func init() {
	conf := &config.Config{AccessKey: acceptance.HW_ACCESS_KEY, SecretKey: acceptance.HW_ACCESS_KEY, Region: acceptance.HW_CUSTOM_REGION_NAME}
	c, _ = httpclient_go.NewHttpClientGo(conf)
}

func getAppResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "apm", conf.Region, "v1/apm2/access-keys")
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.GetAkSkListVO{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}
	return rlt, nil
}

func TestAccApmAkSk_basic(t *testing.T) {
	var instance entity.GetAkSkListVO
	resourceName := "huaweicloud_apm_aksk.aksk_basic"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getAppResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			//acceptance.TestAccPreCheckInternal(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testApmAkSk_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "create aksk by terraform"),
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

func testApmAkSk_basic() string {
	return fmt.Sprintf(`
provider "huaweicloud"{
  region     ="cn-north-7"
  access_key ="DADYWPU8JMUV3UGPEI9"
  secret_key ="jUtvcc0oIIcGZGoAUvtlSi80Z6sZDFI2ZqFKBGUZ"
  auth_url   ="https://iam.cn-north-7.myhuaweicloud.com"
  endpoint = {
     aom : "aom.cn-north-7.myhuaweicloud.com"
     cce : "cce.cn-north-7.myhuaweicloud.com"
     rds : "rds.cn-north-7.myhuaweicloud.com"
     iam : "iam.cn-north-7.myhuaweicloud.com:31943"
     dcsv2 : "dcs.cn-north-7.myhuaweicloud.com"
     obs : "obs.cn-north-7.myhuaweicloud.com"
     vpc : "vpc.cn-north-7.myhuaweicloud.com"
     elb : "elb.cn-north-7.myhuaweicloud.com"
     apm : "apm2.cn-north-7.myhuaweicloud.com"
  }
  insecure = true
  domain_id = "40de487942a74a70b4666fa32d11ffa8"
}

resource "huaweicloud_apm_aksk" "aksk_basic" {
  description = "create aksk by terraform"
}`)
}
