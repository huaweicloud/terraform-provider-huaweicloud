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

func getResourceSwrEnterpriseDomainName(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/domainname?uid={domain_name_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{domain_name_id}", state.Primary.Attributes["domain_name_id"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	domainName := utils.PathSearch("domain_name_infos[0]", getRespBody, nil)
	if domainName == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return domainName, nil
}

func TestAccSwrEnterpriseDomainName_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_domain_name.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseDomainName,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckScmCertificateId(t)
			acceptance.TestAccPreCheckScmCertificateDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseDomainName_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain_name", acceptance.HW_SCM_CERTIFICATE_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "certificate_id", acceptance.HW_SCM_CERTIFICATE_ID),
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

func testAccSwrEnterpriseDomainName_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_domain_name" "test" {
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  domain_name    = "%[2]s"
  certificate_id = "%[3]s"
}
`, testAccSwrEnterpriseInstance_update(rName), acceptance.HW_SCM_CERTIFICATE_DOMAIN_NAME, acceptance.HW_SCM_CERTIFICATE_ID)
}
