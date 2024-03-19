package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/certificates"
	wafcertificates "github.com/chnsz/golangsdk/openstack/waf/v1/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

const (
	targetServiceWaf = "WAF"
	targetServiceElb = "ELB"
)

func getCertificatePushResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	projectName := state.Primary.Attributes["targets.1.project_name"]
	certId := state.Primary.Attributes["targets.1.cert_id"]
	targetService := state.Primary.Attributes["service"]
	if targetService == targetServiceElb {
		elbClient, err := conf.ElbV3Client(projectName)
		if err != nil {
			return nil, fmt.Errorf("error creating ELB client: %s", err)
		}

		return certificates.Get(elbClient, certId).Extract()
	}
	if targetService == targetServiceWaf {
		wafClient, err := conf.WafV1Client(projectName)
		if err != nil {
			return nil, fmt.Errorf("error creating WAF client: %s", err)
		}

		return wafcertificates.Get(wafClient, certId).Extract()
	}

	return nil, fmt.Errorf("find certificate is failed")
}
func TestAccCcmCertificatePush_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_ccm_certificate_push.test"
	project := "cn-north-4"
	project2 := "cn-east-3"
	changeProject := "cn-south-1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCertificatePushResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMPushCertificateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesCcmCertificatePush_basic(targetServiceElb, project, project2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.1.project_name", project2),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_name"),
				),
			},
			{
				Config: tesCcmCertificatePush_basic(targetServiceElb, project, changeProject),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.1.project_name", changeProject),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_name"),
				),
			},
		},
	})
}

func TestAccCcmCertificatePush_waf(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_ccm_certificate_push.test"
	project := "cn-north-4"
	project2 := "cn-east-3"
	changeProject := "cn-south-1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCertificatePushResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMPushCertificateID(t)
			acceptance.TestAccPreCheckCCMPushWAFInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesCcmCertificatePush_basic(targetServiceWaf, project, project2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.1.project_name", project2),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_name"),
				),
			},
			{
				Config: tesCcmCertificatePush_basic(targetServiceWaf, project, changeProject),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "targets.1.project_name", changeProject),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_id"),
					resource.TestCheckResourceAttrSet(resourceName, "targets.1.cert_name"),
				),
			},
		},
	})
}

func tesCcmCertificatePush_basic(service string, project string, project2 string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_push" "test"{
  region         = "cn-north-4"
  certificate_id = "%s"
  service        = "%s"
  targets {
    project_name = "%s"
 }
 targets {
    project_name = "%s"
 }
}`, acceptance.HW_CERT_BATCH_PUSH_ID, service, project, project2)
}
