package elb

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getELBCertificatePrivateKeyEchoResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/certificates/settings/private-key-echo"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

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

	privateKeyEcho := utils.PathSearch("private_key_echo", getRespBody, false).(bool)
	if !privateKeyEcho {
		return nil, golangsdk.ErrDefault404{}
	}
	return getRespBody, nil
}

func TestAccElbCertificatePrivateKeyEcho_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_elb_certificate_private_key_echo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBCertificatePrivateKeyEchoResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbCertificatePrivateKeyEcho_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

func testAccElbCertificatePrivateKeyEcho_basic() string {
	return `resource "huaweicloud_elb_certificate_private_key_echo" "test" {}`
}
