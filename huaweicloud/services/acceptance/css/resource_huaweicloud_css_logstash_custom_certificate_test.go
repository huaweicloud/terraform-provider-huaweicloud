package css

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getLogstashCertificateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	cssV1Client, err := cfg.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	getCertificateHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/certs/{cert_id}"
	getCertificatePath := cssV1Client.Endpoint + getCertificateHttpUrl
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{project_id}", cssV1Client.ProjectID)
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{cluster_id}", state.Primary.Attributes["cluster_id"])
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{cert_id}", state.Primary.ID)

	getCertificatePathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getCertificateResp, err := cssV1Client.Request("GET", getCertificatePath, &getCertificatePathOpt)
	if err != nil {
		return nil, err
	}

	getCertificateRespBody, err := utils.FlattenResponse(getCertificateResp)
	if err != nil {
		return nil, fmt.Errorf("erorr retrieving CSS logstash cluster custom certificate: %s", err)
	}

	return getCertificateRespBody, nil
}

func TestAccLogstashCertificate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_css_logstash_custom_certificate.test"

	tmpFile, err := os.CreateTemp("", "tf-css-cert-test.cer")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write some data to the tempfile
	err = os.WriteFile(tmpFile.Name(), []byte("initial only test"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLogstashCertificateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogstashCertificate_basic(name, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_logstash_cluster.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "path"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket_name", "cert_object"},
				ImportStateIdFunc:       testAccResourceLogstashCertificateImportStateIDFunc(rName),
			},
		},
	})
}

func testAccResourceLogstashCertificateImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		certID := rs.Primary.ID
		if clusterID == "" || certID == "" {
			return "", fmt.Errorf("invalid format specified for import ID, " +
				"must be '<cluster_id>/<id>'")
		}
		return fmt.Sprintf("%s/%s", clusterID, certID), nil
	}
}

func testLogstashCertificate_basic(name, source string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_css_logstash_custom_certificate" "test" {
  cluster_id  = huaweicloud_css_logstash_cluster.test.id
  bucket_name = huaweicloud_obs_bucket.object_bucket.bucket
  cert_object = huaweicloud_obs_bucket_object.object.key
}
`, testAccLogstashCluster_basic(name, 1, "bar"), testCertificateBucket_basic(source))
}

func testCertificateBucket_basic(source string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
  bucket = "tf-css-cert-test-bucket"
}

resource "huaweicloud_obs_bucket_object" "object" {
  bucket       = huaweicloud_obs_bucket.object_bucket.bucket
  key          = "testCssCert.cer"
  source       = "%s"
  content_type = "binary/octet-stream"
}
`, source)
}
