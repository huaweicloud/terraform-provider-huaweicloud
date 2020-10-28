package huaweicloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	OS_AVAILABILITY_ZONE      = os.Getenv("OS_AVAILABILITY_ZONE")
	OS_DEPRECATED_ENVIRONMENT = os.Getenv("OS_DEPRECATED_ENVIRONMENT")
	OS_DNS_ENVIRONMENT        = os.Getenv("OS_DNS_ENVIRONMENT")
	OS_EXTGW_ID               = os.Getenv("OS_EXTGW_ID")
	OS_FLAVOR_ID              = os.Getenv("OS_FLAVOR_ID")
	OS_FLAVOR_NAME            = os.Getenv("OS_FLAVOR_NAME")
	OS_IMAGE_ID               = os.Getenv("OS_IMAGE_ID")
	OS_IMAGE_NAME             = os.Getenv("OS_IMAGE_NAME")
	OS_NETWORK_ID             = os.Getenv("OS_NETWORK_ID")
	OS_SUBNET_ID              = os.Getenv("OS_SUBNET_ID")
	OS_POOL_NAME              = os.Getenv("OS_POOL_NAME")
	OS_REGION_NAME            = os.Getenv("OS_REGION_NAME")
	OS_CUSTOM_REGION_NAME     = os.Getenv("OS_CUSTOM_REGION_NAME")
	OS_ACCESS_KEY             = os.Getenv("OS_ACCESS_KEY")
	OS_SECRET_KEY             = os.Getenv("OS_SECRET_KEY")
	OS_SRC_ACCESS_KEY         = os.Getenv("OS_SRC_ACCESS_KEY")
	OS_SRC_SECRET_KEY         = os.Getenv("OS_SRC_SECRET_KEY")
	OS_VPC_ID                 = os.Getenv("OS_VPC_ID")
	OS_TENANT_ID              = os.Getenv("OS_TENANT_ID")
	OS_DOMAIN_ID              = os.Getenv("OS_DOMAIN_ID")
	OS_DWS_ENVIRONMENT        = os.Getenv("OS_DWS_ENVIRONMENT")
	OS_MRS_ENVIRONMENT        = os.Getenv("OS_MRS_ENVIRONMENT")
	OS_DMS_ENVIRONMENT        = os.Getenv("OS_DMS_ENVIRONMENT")
	OS_NAT_ENVIRONMENT        = os.Getenv("OS_NAT_ENVIRONMENT")
	OS_KMS_ENVIRONMENT        = os.Getenv("OS_KMS_ENVIRONMENT")
	OS_CCI_ENVIRONMENT        = os.Getenv("OS_CCI_ENVIRONMENT")
	OS_CDN_DOMAIN_NAME        = os.Getenv("OS_CDN_DOMAIN_NAME")
	OS_ADMIN                  = os.Getenv("OS_ADMIN")
	OS_ENTERPRISE_PROJECT_ID  = os.Getenv("OS_ENTERPRISE_PROJECT_ID")
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"huaweicloud": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	// Do not run the test if this is a deprecated testing environment.
	if OS_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}
}

func testAccPrecheckCustomRegion(t *testing.T) {
	if OS_CUSTOM_REGION_NAME == "" {
		t.Skip("This environment does not support custom region tests")
	}
}

func testAccPreCheckDeprecated(t *testing.T) {
	if OS_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}
}

func testAccPreCheckAdminOnly(t *testing.T) {
	if OS_ADMIN == "" {
		t.Skip("Skipping test because it requires the admin user group")
	}
}

func testAccPreCheckDNS(t *testing.T) {
	if OS_DNS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DNS tests")
	}
}

func testAccPreCheckULB(t *testing.T) {
	if OS_SUBNET_ID == "" {
		t.Skip("OS_SUBNET must be set for LB acceptance tests")
	}
}

func testAccPreCheckELB(t *testing.T) {
	if OS_TENANT_ID == "" {
		t.Skip("This environment does not support ELB tests")
	}
}

func testAccPreCheckMaas(t *testing.T) {
	if OS_ACCESS_KEY == "" || OS_SECRET_KEY == "" || OS_SRC_ACCESS_KEY == "" || OS_SRC_SECRET_KEY == "" {
		t.Skip("OS_ACCESS_KEY, OS_SECRET_KEY, OS_SRC_ACCESS_KEY, and OS_SRC_SECRET_KEY  must be set for MAAS acceptance tests")
	}
}

func testAccPreCheckS3(t *testing.T) {
	if OS_ACCESS_KEY == "" || OS_SECRET_KEY == "" {
		t.Skip("OS_ACCESS_KEY and OS_SECRET_KEY  must be set for S3 acceptance tests")
	}
}

func testAccPreCheckImage(t *testing.T) {
	if OS_ACCESS_KEY != "" && OS_SECRET_KEY != "" {
		t.Skip("AK/SK authentication doesn't support images tests")
	}
}

func testAccPreCheckDws(t *testing.T) {
	if OS_DWS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DWS tests")
	}
}

func testAccPreCheckMrs(t *testing.T) {
	if OS_MRS_ENVIRONMENT == "" {
		t.Skip("This environment does not support MRS tests")
	}
}

func testAccPreCheckDms(t *testing.T) {
	if OS_DMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DMS tests")
	}
}

func testAccPreCheckNat(t *testing.T) {
	if OS_NAT_ENVIRONMENT == "" {
		t.Skip("This environment does not support NAT tests")
	}
}

func testAccPreCheckKms(t *testing.T) {
	if OS_KMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support KMS tests")
	}
}

func testAccPreCheckCDN(t *testing.T) {
	if OS_CDN_DOMAIN_NAME == "" {
		t.Skip("This environment does not support CDN tests")
	}
}

func testAccPreCheckCCI(t *testing.T) {
	if OS_CCI_ENVIRONMENT == "" {
		t.Skip("This environment does not support CCI tests")
	}
}

func testAccPreCheckEpsID(t *testing.T) {
	if OS_ENTERPRISE_PROJECT_ID == "" {
		t.Skip("This environment does not support EPS_ID tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

// Steps for configuring HuaweiCloud with SSL validation are here:
// https://github.com/hashicorp/terraform/pull/6279#issuecomment-219020144
func TestAccProvider_caCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloud SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping HuaweiCloud CA test.")
	}

	p := Provider()

	caFile, err := envVarFile("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caFile)

	raw := map[string]interface{}{
		"cacert_file": caFile,
	}

	err = p.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected err when specifying HuaweiCloud CA by file: %s", err)
	}
}

func TestAccProvider_caCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloud SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping HuaweiCloud CA test.")
	}

	p := Provider()

	caContents, err := envVarContents("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	raw := map[string]interface{}{
		"cacert_file": caContents,
	}

	err = p.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected err when specifying HuaweiCloud CA by string: %s", err)
	}
}

func TestAccProvider_clientCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloud SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping HuaweiCloud client SSL auth test.")
	}

	p := Provider()

	certFile, err := envVarFile("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(certFile)
	keyFile, err := envVarFile("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(keyFile)

	raw := map[string]interface{}{
		"cert": certFile,
		"key":  keyFile,
	}

	err = p.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected err when specifying HuaweiCloud Client keypair by file: %s", err)
	}
}

func TestAccProvider_clientCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping HuaweiCloud SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping HuaweiCloud client SSL auth test.")
	}

	p := Provider()

	certContents, err := envVarContents("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	keyContents, err := envVarContents("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"cert": certContents,
		"key":  keyContents,
	}

	err = p.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected err when specifying HuaweiCloud Client keypair by contents: %s", err)
	}
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmt.Errorf("Error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile("", varName)
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}

func testAccAsConfigPreCheck(t *testing.T) {
	if OS_FLAVOR_ID == "" {
		t.Skip("OS_FLAVOR_ID must be set for acceptance tests")
	}
}
