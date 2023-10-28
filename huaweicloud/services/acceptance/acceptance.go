//nolint:revive
package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
)

var (
	HW_REGION_NAME                        = os.Getenv("HW_REGION_NAME")
	HW_CUSTOM_REGION_NAME                 = os.Getenv("HW_CUSTOM_REGION_NAME")
	HW_AVAILABILITY_ZONE                  = os.Getenv("HW_AVAILABILITY_ZONE")
	HW_ACCESS_KEY                         = os.Getenv("HW_ACCESS_KEY")
	HW_SECRET_KEY                         = os.Getenv("HW_SECRET_KEY")
	HW_USER_ID                            = os.Getenv("HW_USER_ID")
	HW_USER_NAME                          = os.Getenv("HW_USER_NAME")
	HW_PROJECT_ID                         = os.Getenv("HW_PROJECT_ID")
	HW_DOMAIN_ID                          = os.Getenv("HW_DOMAIN_ID")
	HW_DOMAIN_NAME                        = os.Getenv("HW_DOMAIN_NAME")
	HW_ENTERPRISE_PROJECT_ID_TEST         = os.Getenv("HW_ENTERPRISE_PROJECT_ID_TEST")
	HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST = os.Getenv("HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST")

	HW_FLAVOR_ID             = os.Getenv("HW_FLAVOR_ID")
	HW_FLAVOR_NAME           = os.Getenv("HW_FLAVOR_NAME")
	HW_IMAGE_ID              = os.Getenv("HW_IMAGE_ID")
	HW_IMAGE_NAME            = os.Getenv("HW_IMAGE_NAME")
	HW_VPC_ID                = os.Getenv("HW_VPC_ID")
	HW_NETWORK_ID            = os.Getenv("HW_NETWORK_ID")
	HW_SUBNET_ID             = os.Getenv("HW_SUBNET_ID")
	HW_ENTERPRISE_PROJECT_ID = os.Getenv("HW_ENTERPRISE_PROJECT_ID")
	HW_ADMIN                 = os.Getenv("HW_ADMIN")

	HW_MAPREDUCE_CUSTOM           = os.Getenv("HW_MAPREDUCE_CUSTOM")
	HW_MAPREDUCE_BOOTSTRAP_SCRIPT = os.Getenv("HW_MAPREDUCE_BOOTSTRAP_SCRIPT")

	HW_CNAD_ENABLE_FLAG       = os.Getenv("HW_CNAD_ENABLE_FLAG")
	HW_CNAD_PROJECT_OBJECT_ID = os.Getenv("HW_CNAD_PROJECT_OBJECT_ID")

	HW_OBS_BUCKET_NAME        = os.Getenv("HW_OBS_BUCKET_NAME")
	HW_OBS_DESTINATION_BUCKET = os.Getenv("HW_OBS_DESTINATION_BUCKET")
	HW_OBS_USER_DOMAIN_NAME1  = os.Getenv("HW_OBS_USER_DOMAIN_NAME1")
	HW_OBS_USER_DOMAIN_NAME2  = os.Getenv("HW_OBS_USER_DOMAIN_NAME2")

	HW_OMS_ENABLE_FLAG = os.Getenv("HW_OMS_ENABLE_FLAG")

	HW_DEPRECATED_ENVIRONMENT = os.Getenv("HW_DEPRECATED_ENVIRONMENT")
	HW_INTERNAL_USED          = os.Getenv("HW_INTERNAL_USED")

	HW_WAF_ENABLE_FLAG = os.Getenv("HW_WAF_ENABLE_FLAG")

	HW_DEST_REGION          = os.Getenv("HW_DEST_REGION")
	HW_DEST_PROJECT_ID      = os.Getenv("HW_DEST_PROJECT_ID")
	HW_DEST_PROJECT_ID_TEST = os.Getenv("HW_DEST_PROJECT_ID_TEST")
	HW_CHARGING_MODE        = os.Getenv("HW_CHARGING_MODE")
	HW_HIGH_COST_ALLOW      = os.Getenv("HW_HIGH_COST_ALLOW")
	HW_SWR_SHARING_ACCOUNT  = os.Getenv("HW_SWR_SHARING_ACCOUNT")

	HW_RAM_SHARE_ACCOUNT_ID          = os.Getenv("HW_RAM_SHARE_ACCOUNT_ID")
	HW_RAM_SHARE_RESOURCE_URN        = os.Getenv("HW_RAM_SHARE_RESOURCE_URN")
	HW_RAM_SHARE_UPDATE_ACCOUNT_ID   = os.Getenv("HW_RAM_SHARE_UPDATE_ACCOUNT_ID")
	HW_RAM_SHARE_UPDATE_RESOURCE_URN = os.Getenv("HW_RAM_SHARE_UPDATE_RESOURCE_URN")

	HW_CDN_DOMAIN_NAME              = os.Getenv("HW_CDN_DOMAIN_NAME")
	HW_CDN_CERT_PATH                = os.Getenv("HW_CDN_CERT_PATH")
	HW_CDN_PRIVATE_KEY_PATH         = os.Getenv("HW_CDN_PRIVATE_KEY_PATH")
	HW_CERTIFICATE_KEY_PATH         = os.Getenv("HW_CERTIFICATE_KEY_PATH")
	HW_CERTIFICATE_CHAIN_PATH       = os.Getenv("HW_CERTIFICATE_CHAIN_PATH")
	HW_CERTIFICATE_PRIVATE_KEY_PATH = os.Getenv("HW_CERTIFICATE_PRIVATE_KEY_PATH")
	HW_CERTIFICATE_SERVICE          = os.Getenv("HW_CERTIFICATE_SERVICE")
	HW_CERTIFICATE_PROJECT          = os.Getenv("HW_CERTIFICATE_PROJECT")
	HW_CERTIFICATE_PROJECT_UPDATED  = os.Getenv("HW_CERTIFICATE_PROJECT_UPDATED")
	HW_CERTIFICATE_NAME             = os.Getenv("HW_CERTIFICATE_NAME")
	HW_DMS_ENVIRONMENT              = os.Getenv("HW_DMS_ENVIRONMENT")
	HW_SMS_SOURCE_SERVER            = os.Getenv("HW_SMS_SOURCE_SERVER")

	HW_DLI_FLINK_JAR_OBS_PATH           = os.Getenv("HW_DLI_FLINK_JAR_OBS_PATH")
	HW_DLI_DS_AUTH_CSS_OBS_PATH         = os.Getenv("HW_DLI_DS_AUTH_CSS_OBS_PATH")
	HW_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH = os.Getenv("HW_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH")
	HW_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH   = os.Getenv("HW_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH")
	HW_DLI_DS_AUTH_KRB_CONF_OBS_PATH    = os.Getenv("HW_DLI_DS_AUTH_KRB_CONF_OBS_PATH")
	HW_DLI_DS_AUTH_KRB_TAB_OBS_PATH     = os.Getenv("HW_DLI_DS_AUTH_KRB_TAB_OBS_PATH")
	HW_DLI_AGENCY_FLAG                  = os.Getenv("HW_DLI_AGENCY_FLAG")

	HW_GITHUB_REPO_HOST        = os.Getenv("HW_GITHUB_REPO_HOST")        // Repository host (Github, Gitlab, Gitee)
	HW_GITHUB_PERSONAL_TOKEN   = os.Getenv("HW_GITHUB_PERSONAL_TOKEN")   // Personal access token (Github, Gitlab, Gitee)
	HW_GITHUB_REPO_PWD         = os.Getenv("HW_GITHUB_REPO_PWD")         // Repository password (DevCloud, BitBucket)
	HW_GITHUB_REPO_URL         = os.Getenv("HW_GITHUB_REPO_URL")         // Repository URL (Github, Gitlab, Gitee)
	HW_OBS_STORAGE_URL         = os.Getenv("HW_OBS_STORAGE_URL")         // OBS storage URL where ZIP file is located
	HW_BUILD_IMAGE_URL         = os.Getenv("HW_BUILD_IMAGE_URL")         // SWR Image URL for component deployment
	HW_BUILD_IMAGE_URL_UPDATED = os.Getenv("HW_BUILD_IMAGE_URL_UPDATED") // SWR Image URL for component deployment update

	HW_VOD_WATERMARK_FILE   = os.Getenv("HW_VOD_WATERMARK_FILE")
	HW_VOD_MEDIA_ASSET_FILE = os.Getenv("HW_VOD_MEDIA_ASSET_FILE")

	HW_LTS_ENABLE_FLAG                 = os.Getenv("HW_LTS_ENABLE_FLAG")
	HW_LTS_STRUCT_CONFIG_TEMPLATE_ID   = os.Getenv("HW_LTS_STRUCT_CONFIG_TEMPLATE_ID")
	HW_LTS_STRUCT_CONFIG_TEMPLATE_NAME = os.Getenv("HW_LTS_STRUCT_CONFIG_TEMPLATE_NAME")

	HW_CHAIR_EMAIL              = os.Getenv("HW_CHAIR_EMAIL")
	HW_GUEST_EMAIL              = os.Getenv("HW_GUEST_EMAIL")
	HW_MEETING_ACCOUNT_NAME     = os.Getenv("HW_MEETING_ACCOUNT_NAME")
	HW_MEETING_ACCOUNT_PASSWORD = os.Getenv("HW_MEETING_ACCOUNT_PASSWORD")
	HW_MEETING_APP_ID           = os.Getenv("HW_MEETING_APP_ID")
	HW_MEETING_APP_KEY          = os.Getenv("HW_MEETING_APP_KEY")
	HW_MEETING_USER_ID          = os.Getenv("HW_MEETING_USER_ID")
	HW_MEETING_ROOM_ID          = os.Getenv("HW_MEETING_ROOM_ID")

	HW_AAD_INSTANCE_ID = os.Getenv("HW_AAD_INSTANCE_ID")
	HW_AAD_IP_ADDRESS  = os.Getenv("HW_AAD_IP_ADDRESS")

	HW_WORKSPACE_AD_DOMAIN_NAME = os.Getenv("HW_WORKSPACE_AD_DOMAIN_NAME") // Domain name, e.g. "example.com".
	HW_WORKSPACE_AD_SERVER_PWD  = os.Getenv("HW_WORKSPACE_AD_SERVER_PWD")  // The password of AD server.
	HW_WORKSPACE_AD_DOMAIN_IP   = os.Getenv("HW_WORKSPACE_AD_DOMAIN_IP")   // Active domain IP, e.g. "192.168.196.3".
	HW_WORKSPACE_AD_VPC_ID      = os.Getenv("HW_WORKSPACE_AD_VPC_ID")      // The VPC ID to which the AD server and desktops belongs.
	HW_WORKSPACE_AD_NETWORK_ID  = os.Getenv("HW_WORKSPACE_AD_NETWORK_ID")  // The network ID to which the AD server belongs.

	HW_FGS_TRIGGER_LTS_AGENCY = os.Getenv("HW_FGS_TRIGGER_LTS_AGENCY")

	HW_KMS_ENVIRONMENT = os.Getenv("HW_KMS_ENVIRONMENT")

	HW_MULTI_ACCOUNT_ENVIRONMENT            = os.Getenv("HW_MULTI_ACCOUNT_ENVIRONMENT")
	HW_ORGANIZATIONS_OPEN                   = os.Getenv("HW_ORGANIZATIONS_OPEN")
	HW_ORGANIZATIONS_ACCOUNT_NAME           = os.Getenv("HW_ORGANIZATIONS_ACCOUNT_NAME")
	HW_ORGANIZATIONS_INVITE_ACCOUNT_ID      = os.Getenv("HW_ORGANIZATIONS_INVITE_ACCOUNT_ID")
	HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID = os.Getenv("HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID")
	HW_ORGANIZATIONS_INVITATION_ID          = os.Getenv("HW_ORGANIZATIONS_INVITATION_ID")

	HW_IDENTITY_CENTER_ACCOUNT_ID = os.Getenv("HW_IDENTITY_CENTER_ACCOUNT_ID")

	HW_ER_TEST_ON = os.Getenv("HW_ER_TEST_ON") // Whether to run the ER related tests.

	// The OBS address where the HCL/JSON template archive (No variables) is located.
	HW_RF_TEMPLATE_ARCHIVE_NO_VARS_URI = os.Getenv("HW_RF_TEMPLATE_ARCHIVE_NO_VARS_URI")
	// The OBS address where the HCL/JSON template archive is located.
	HW_RF_TEMPLATE_ARCHIVE_URI = os.Getenv("HW_RF_TEMPLATE_ARCHIVE_URI")
	// The OBS address where the variable archive corresponding to the HCL/JSON template is located.
	HW_RF_VARIABLES_ARCHIVE_URI = os.Getenv("HW_RF_VARIABLES_ARCHIVE_URI")

	// The direct connection ID (provider does not support direct connection resource).
	HW_DC_DIRECT_CONNECT_ID = os.Getenv("HW_DC_DIRECT_CONNECT_ID")

	// The CFW instance ID
	HW_CFW_INSTANCE_ID = os.Getenv("HW_CFW_INSTANCE_ID")

	// The cluster ID of the CCE
	HW_CCE_CLUSTER_ID = os.Getenv("HW_CCE_CLUSTER_ID")
	// The cluster name of the CCE
	HW_CCE_CLUSTER_NAME = os.Getenv("HW_CCE_CLUSTER_NAME")
	// The cluster ID of the CCE
	HW_CCE_CLUSTER_ID_ANOTHER = os.Getenv("HW_CCE_CLUSTER_ID_ANOTHER")
	// The cluster name of the CCE
	HW_CCE_CLUSTER_NAME_ANOTHER = os.Getenv("HW_CCE_CLUSTER_NAME_ANOTHER")
	// The partition az of the CCE
	HW_CCE_PARTITION_AZ = os.Getenv("HW_CCE_PARTITION_AZ")
	// The namespace of the workload is located
	HW_WORKLOAD_NAMESPACE = os.Getenv("HW_WORKLOAD_NAMESPACE")
	// The workload type deployed in CCE/CCI
	HW_WORKLOAD_TYPE = os.Getenv("HW_WORKLOAD_TYPE")
	// The workload name deployed in CCE/CCI
	HW_WORKLOAD_NAME = os.Getenv("HW_WORKLOAD_NAME")
	// The target region of SWR image auto sync
	HW_SWR_TARGET_REGION = os.Getenv("HW_SWR_TARGET_REGION")
	// The target organization of SWR image auto sync
	HW_SWR_TARGET_ORGANIZATION = os.Getenv("HW_SWR_TARGET_ORGANIZATION")

	// The ID of the CBR backup
	HW_IMS_BACKUP_ID = os.Getenv("HW_IMS_BACKUP_ID")

	// The SecMaster workspace ID
	HW_SECMASTER_WORKSPACE_ID = os.Getenv("HW_SECMASTER_WORKSPACE_ID")

	HW_MODELARTS_HAS_SUBSCRIBE_MODEL = os.Getenv("HW_MODELARTS_HAS_SUBSCRIBE_MODEL")

	// The CMDB sub-application ID of AOM service
	HW_AOM_SUB_APPLICATION_ID = os.Getenv("HW_AOM_SUB_APPLICATION_ID")

	// Deprecated
	HW_SRC_ACCESS_KEY = os.Getenv("HW_SRC_ACCESS_KEY")
	HW_SRC_SECRET_KEY = os.Getenv("HW_SRC_SECRET_KEY")
	HW_EXTGW_ID       = os.Getenv("HW_EXTGW_ID")
	HW_POOL_NAME      = os.Getenv("HW_POOL_NAME")

	HW_IMAGE_SHARE_SOURCE_IMAGE_ID = os.Getenv("HW_IMAGE_SHARE_SOURCE_IMAGE_ID")

	HW_CERTIFICATE_CONTENT         = os.Getenv("HW_CERTIFICATE_CONTENT")
	HW_CERTIFICATE_PRIVATE_KEY     = os.Getenv("HW_CERTIFICATE_PRIVATE_KEY")
	HW_CERTIFICATE_ROOT_CA         = os.Getenv("HW_CERTIFICATE_ROOT_CA")
	HW_NEW_CERTIFICATE_CONTENT     = os.Getenv("HW_NEW_CERTIFICATE_CONTENT")
	HW_NEW_CERTIFICATE_PRIVATE_KEY = os.Getenv("HW_NEW_CERTIFICATE_PRIVATE_KEY")
	HW_NEW_CERTIFICATE_ROOT_CA     = os.Getenv("HW_NEW_CERTIFICATE_ROOT_CA")

	HW_CODEARTS_RESOURCE_POOL_ID = os.Getenv("HW_CODEARTS_RESOURCE_POOL_ID")
	HW_CODEARTS_ENABLE_FLAG      = os.Getenv("HW_CODEARTS_ENABLE_FLAG")

	HW_EG_CHANNEL_ID = os.Getenv("HW_EG_CHANNEL_ID")

	HW_KOOGALLERY_ASSET = os.Getenv("HW_KOOGALLERY_ASSET")

	HW_CCI_NAMESPACE = os.Getenv("HW_CCI_NAMESPACE")
)

// TestAccProviders is a static map containing only the main provider instance.
//
// Deprecated: Terraform Plugin SDK version 2 uses TestCase.ProviderFactories
// but supports this value in TestCase.Providers for backwards compatibility.
// In the future Providers: TestAccProviders will be changed to
// ProviderFactories: TestAccProviderFactories
var TestAccProviders map[string]*schema.Provider

// TestAccProviderFactories is a static map containing only the main provider instance
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// TestAccProvider is the "main" provider instance
var TestAccProvider *schema.Provider

func init() {
	TestAccProvider = huaweicloud.Provider()

	TestAccProviders = map[string]*schema.Provider{
		"huaweicloud": TestAccProvider,
	}

	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"huaweicloud": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

func preCheckRequiredEnvVars(t *testing.T) {
	if HW_REGION_NAME == "" {
		t.Fatal("HW_REGION_NAME must be set for acceptance tests")
	}
}

// use this function to precheck langding zone services, such as Organizations and Identity Center
// lintignore:AT003
func TestAccPreCheckMultiAccount(t *testing.T) {
	if HW_MULTI_ACCOUNT_ENVIRONMENT == "" {
		t.Skip("This environment does not support multi-account tests")
	}
}

// when this variable is set, the Organizations service should be enabled, and the organization info
// can be get by the API
// lintignore:AT003
func TestAccPreCheckOrganizationsOpen(t *testing.T) {
	if HW_ORGANIZATIONS_OPEN == "" {
		t.Skip("HW_ORGANIZATIONS_OPEN must be set for the acceptance test")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsAccountName(t *testing.T) {
	if HW_ORGANIZATIONS_ACCOUNT_NAME == "" {
		t.Skip("HW_ORGANIZATIONS_ACCOUNT_NAME must be set for the acceptance test")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsInviteAccountId(t *testing.T) {
	if HW_ORGANIZATIONS_INVITE_ACCOUNT_ID == "" {
		t.Skip("HW_ORGANIZATIONS_INVITE_ACCOUNT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsInvitationId(t *testing.T) {
	if HW_ORGANIZATIONS_INVITATION_ID == "" {
		t.Skip("HW_ORGANIZATIONS_INVITATION_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOrganizationsOrganizationalUnitId(t *testing.T) {
	if HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID == "" {
		t.Skip("HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckIdentityCenterAccountId(t *testing.T) {
	if HW_IDENTITY_CENTER_ACCOUNT_ID == "" {
		t.Skip("HW_IDENTITY_CENTER_ACCOUNT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheck(t *testing.T) {
	// Do not run the test if this is a deprecated testing environment.
	if HW_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

// lintignore:AT003
func TestAccPrecheckDomainId(t *testing.T) {
	if HW_DOMAIN_ID == "" {
		t.Skip("HW_DOMAIN_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPrecheckDomainName(t *testing.T) {
	if HW_DOMAIN_NAME == "" {
		t.Skip("HW_DOMAIN_NAME must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPrecheckCustomRegion(t *testing.T) {
	if HW_CUSTOM_REGION_NAME == "" {
		t.Skip("HW_CUSTOM_REGION_NAME must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckDeprecated(t *testing.T) {
	if HW_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

// lintignore:AT003
func TestAccPreCheckInternal(t *testing.T) {
	if HW_INTERNAL_USED == "" {
		t.Skip("HW_INTERNAL_USED must be set for internal acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckEpsID(t *testing.T) {
	// The environment variables in tests take HW_ENTERPRISE_PROJECT_ID_TEST instead of HW_ENTERPRISE_PROJECT_ID to
	// ensure that other data-resources that support enterprise projects query the default project without being
	// affected by this variable.
	if HW_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("The environment variables does not support Enterprise Project ID for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckMigrateEpsID(t *testing.T) {
	if HW_ENTERPRISE_PROJECT_ID_TEST == "" || HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST == "" {
		t.Skip("The environment variables does not support Migrate Enterprise Project ID for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckUserId(t *testing.T) {
	if HW_USER_ID == "" {
		t.Skip("The environment variables does not support the user ID (HW_USER_ID) for acc tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSms(t *testing.T) {
	if HW_SMS_SOURCE_SERVER == "" {
		t.Skip("HW_SMS_SOURCE_SERVER must be set for SMS acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckMrsCustom(t *testing.T) {
	if HW_MAPREDUCE_CUSTOM == "" {
		t.Skip("HW_MAPREDUCE_CUSTOM must be set for acceptance tests:custom type cluster of map reduce")
	}
}

// lintignore:AT003
func TestAccPreCheckMrsBootstrapScript(t *testing.T) {
	if HW_MAPREDUCE_BOOTSTRAP_SCRIPT == "" {
		t.Skip("HW_MAPREDUCE_BOOTSTRAP_SCRIPT must be set for acceptance tests: cluster of map reduce with bootstrap")
	}
}

// lintignore:AT003
func TestAccPreCheckFgsTrigger(t *testing.T) {
	if HW_FGS_TRIGGER_LTS_AGENCY == "" {
		t.Skip("HW_FGS_TRIGGER_LTS_AGENCY must be set for FGS trigger acceptance tests")
	}
}

// Deprecated
// lintignore:AT003
func TestAccPreCheckMaas(t *testing.T) {
	if HW_ACCESS_KEY == "" || HW_SECRET_KEY == "" || HW_SRC_ACCESS_KEY == "" || HW_SRC_SECRET_KEY == "" {
		t.Skip("HW_ACCESS_KEY, HW_SECRET_KEY, HW_SRC_ACCESS_KEY, and HW_SRC_SECRET_KEY  must be set for MAAS acceptance tests")
	}
}

func RandomAccResourceName() string {
	return fmt.Sprintf("tf_test_%s", acctest.RandString(5))
}

func RandomAccResourceNameWithDash() string {
	return fmt.Sprintf("tf-test-%s", acctest.RandString(5))
}

func RandomCidr() string {
	return fmt.Sprintf("172.16.%d.0/24", acctest.RandIntRange(0, 255))
}

func RandomCidrAndGatewayIp() (string, string) {
	seed := acctest.RandIntRange(0, 255)
	return fmt.Sprintf("172.16.%d.0/24", seed), fmt.Sprintf("172.16.%d.1", seed)
}

func RandomPassword() string {
	return fmt.Sprintf("%s%s%s%d",
		acctest.RandStringFromCharSet(2, "ABCDEFGHIJKLMNOPQRSTUVWXZY"),
		acctest.RandStringFromCharSet(3, acctest.CharSetAlpha),
		acctest.RandStringFromCharSet(2, "~!@#%^*-_=+?"),
		acctest.RandIntRange(1000, 9999))
}

// lintignore:AT003
func TestAccPrecheckWafInstance(t *testing.T) {
	if HW_WAF_ENABLE_FLAG == "" {
		t.Skip("Skip the WAF acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckCNADInstance(t *testing.T) {
	if HW_CNAD_ENABLE_FLAG == "" {
		t.Skip("Skip the CNAD acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckCNADProtectedObject(t *testing.T) {
	if HW_CNAD_PROJECT_OBJECT_ID == "" {
		t.Skip("Skipping test because HW_CNAD_PROJECT_OBJECT_ID is required for this acceptance test.")
	}
}

// lintignore:AT003
func TestAccPreCheckOmsInstance(t *testing.T) {
	if HW_OMS_ENABLE_FLAG == "" {
		t.Skip("Skip the OMS acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckAdminOnly(t *testing.T) {
	if HW_ADMIN == "" {
		t.Skip("Skipping test because it requires the admin privileges")
	}
}

// lintignore:AT003
func TestAccPreCheckReplication(t *testing.T) {
	if HW_DEST_REGION == "" || HW_DEST_PROJECT_ID == "" {
		t.Skip("Skip the replication policy acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckProjectId(t *testing.T) {
	if HW_DEST_PROJECT_ID_TEST == "" {
		t.Skip("Skipping test because it requires the test project id.")
	}
}

// lintignore:AT003
func TestAccPreCheckProject(t *testing.T) {
	if HW_ENTERPRISE_PROJECT_ID_TEST != "" {
		t.Skip("This environment does not support project tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBS(t *testing.T) {
	if HW_ACCESS_KEY == "" || HW_SECRET_KEY == "" {
		t.Skip("HW_ACCESS_KEY and HW_SECRET_KEY must be set for OBS acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBSBucket(t *testing.T) {
	if HW_OBS_BUCKET_NAME == "" {
		t.Skip("HW_OBS_BUCKET_NAME must be set for OBS object acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBSDestinationBucket(t *testing.T) {
	if HW_OBS_DESTINATION_BUCKET == "" {
		t.Skip("HW_OBS_DESTINATION_BUCKET must be set for OBS destination tests")
	}
}

// lintignore:AT003
func TestAccPreCheckOBSUserDomainNames(t *testing.T) {
	if HW_OBS_USER_DOMAIN_NAME1 == "" || HW_OBS_USER_DOMAIN_NAME2 == "" {
		t.Skip("HW_OBS_USER_DOMAIN_NAME1 and HW_OBS_USER_DOMAIN_NAME2 must be set for OBS user domain name tests")
	}
}

// lintignore:AT003
func TestAccPreCheckChargingMode(t *testing.T) {
	if HW_CHARGING_MODE != "prePaid" {
		t.Skip("This environment does not support prepaid tests")
	}
}

// lintignore:AT003
func TestAccPreCheckHighCostAllow(t *testing.T) {
	if HW_HIGH_COST_ALLOW == "" {
		t.Skip("Do not allow expensive testing")
	}
}

// lintignore:AT003
func TestAccPreCheckScm(t *testing.T) {
	if HW_CERTIFICATE_KEY_PATH == "" || HW_CERTIFICATE_CHAIN_PATH == "" ||
		HW_CERTIFICATE_PRIVATE_KEY_PATH == "" || HW_CERTIFICATE_SERVICE == "" ||
		HW_CERTIFICATE_PROJECT == "" || HW_CERTIFICATE_PROJECT_UPDATED == "" {
		t.Skip("HW_CERTIFICATE_KEY_PATH, HW_CERTIFICATE_CHAIN_PATH, HW_CERTIFICATE_PRIVATE_KEY_PATH, " +
			"HW_CERTIFICATE_SERVICE, HW_CERTIFICATE_PROJECT and HW_CERTIFICATE_TARGET_UPDATED " +
			"can not be empty for SCM certificate tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSWRDomian(t *testing.T) {
	if HW_SWR_SHARING_ACCOUNT == "" {
		t.Skip("HW_SWR_SHARING_ACCOUNT must be set for swr domian tests, " +
			"the value of HW_SWR_SHARING_ACCOUNT should be another IAM user name")
	}
}

// lintignore:AT003
func TestAccPreCheckRAM(t *testing.T) {
	if HW_RAM_SHARE_ACCOUNT_ID == "" || HW_RAM_SHARE_RESOURCE_URN == "" {
		t.Skip("HW_RAM_SHARE_ACCOUNT_ID and HW_RAM_SHARE_RESOURCE_URN " +
			"must be set for create ram resource tests.")
	}

	if HW_RAM_SHARE_UPDATE_ACCOUNT_ID == "" || HW_RAM_SHARE_UPDATE_RESOURCE_URN == "" {
		t.Skip("HW_RAM_SHARE_UPDATE_ACCOUNT_ID and HW_RAM_SHARE_UPDATE_RESOURCE_URN" +
			" must be set for update ram resource tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDms(t *testing.T) {
	if HW_DMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DMS tests")
	}
}

// lintignore:AT003
func TestAccPreCheckDliJarPath(t *testing.T) {
	if HW_DLI_FLINK_JAR_OBS_PATH == "" {
		t.Skip("HW_DLI_FLINK_JAR_OBS_PATH must be set for DLI Flink Jar job acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthCss(t *testing.T) {
	if HW_DLI_DS_AUTH_CSS_OBS_PATH == "" {
		t.Skip("HW_DLI_DS_AUTH_CSS_OBS_PATH must be set for DLI datasource CSS Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthKafka(t *testing.T) {
	if HW_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH == "" || HW_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH == "" {
		t.Skip("HW_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH,HW_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH must be set for DLI datasource Kafka Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliDsAuthKrb(t *testing.T) {
	if HW_DLI_DS_AUTH_KRB_CONF_OBS_PATH == "" || HW_DLI_DS_AUTH_KRB_TAB_OBS_PATH == "" {
		t.Skip("HW_DLI_DS_AUTH_KRB_CONF_OBS_PATH,HW_DLI_DS_AUTH_KRB_TAB_OBS_PATH must be set for DLI datasource Kafka Auth acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckDliAgency(t *testing.T) {
	if HW_DLI_AGENCY_FLAG == "" {
		t.Skip("HW_DLI_AGENCY_FLAG must be set for DLI datasource DLI agency acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckRepoTokenAuth(t *testing.T) {
	if HW_GITHUB_REPO_HOST == "" || HW_GITHUB_PERSONAL_TOKEN == "" {
		t.Skip("Repository configurations are not completed for acceptance test of personal access token authorization.")
	}
}

// lintignore:AT003
func TestAccPreCheckRepoPwdAuth(t *testing.T) {
	if HW_DOMAIN_NAME == "" || HW_USER_NAME == "" || HW_GITHUB_REPO_PWD == "" {
		t.Skip("Repository configurations are not completed for acceptance test of password authorization.")
	}
}

// lintignore:AT003
func TestAccPreCheckComponent(t *testing.T) {
	if HW_DOMAIN_NAME == "" || HW_GITHUB_REPO_URL == "" || HW_OBS_STORAGE_URL == "" {
		t.Skip("Repository (package) configurations are not completed for acceptance test of component.")
	}
}

// lintignore:AT003
func TestAccPreCheckComponentDeployment(t *testing.T) {
	if HW_BUILD_IMAGE_URL == "" {
		t.Skip("SWR image URL configuration is not completed for acceptance test of component deployment.")
	}
}

// lintignore:AT003
func TestAccPreCheckImageUrlUpdated(t *testing.T) {
	if HW_BUILD_IMAGE_URL_UPDATED == "" {
		t.Skip("SWR image update URL configuration is not completed for acceptance test of component deployment.")
	}
}

// lintignore:AT003
func TestAccPreCheckVODWatermark(t *testing.T) {
	if HW_VOD_WATERMARK_FILE == "" {
		t.Skip("HW_VOD_WATERMARK_FILE must be set for VOD watermark template acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckVODMediaAsset(t *testing.T) {
	if HW_VOD_MEDIA_ASSET_FILE == "" {
		t.Skip("HW_VOD_MEDIA_ASSET_FILE must be set for VOD media asset acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckPwdAuth(t *testing.T) {
	if HW_MEETING_ACCOUNT_NAME == "" || HW_MEETING_ACCOUNT_PASSWORD == "" {
		t.Skip("The account name (HW_MEETING_ACCOUNT_NAME) or password (HW_MEETING_ACCOUNT_PASSWORD) is not " +
			"completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckAppAuth(t *testing.T) {
	if HW_MEETING_APP_ID == "" || HW_MEETING_APP_KEY == "" || HW_MEETING_USER_ID == "" {
		t.Skip("The app ID (HW_MEETING_APP_ID), app KEY (HW_MEETING_APP_KEY) or user ID (HW_MEETING_USER_ID) is not " +
			"completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckMeetingRoom(t *testing.T) {
	if HW_MEETING_ROOM_ID == "" {
		t.Skip("The vmr ID (HW_MEETING_ROOM_ID) is not completed for acceptance test of conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckParticipants(t *testing.T) {
	if HW_CHAIR_EMAIL == "" || HW_GUEST_EMAIL == "" {
		t.Skip("The chair (HW_CHAIR_EMAIL) or guest (HW_GUEST_EMAIL) mailbox is not completed for acceptance test of " +
			"conference.")
	}
}

// lintignore:AT003
func TestAccPreCheckAadForwardRule(t *testing.T) {
	if HW_AAD_INSTANCE_ID == "" || HW_AAD_IP_ADDRESS == "" {
		t.Skip("The instance information is not completed for AAD rule acceptance test.")
	}
}

// lintignore:AT003
func TestAccPreCheckScmCertificateName(t *testing.T) {
	if HW_CERTIFICATE_NAME == "" {
		t.Skip("HW_CERTIFICATE_NAME must be set for SCM acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckKms(t *testing.T) {
	if HW_KMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support KMS tests")
	}
}

// lintignore:AT003
func TestAccPreCheckProjectID(t *testing.T) {
	if HW_PROJECT_ID == "" {
		t.Skip("HW_PROJECT_ID must be set for acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkspaceAD(t *testing.T) {
	if HW_WORKSPACE_AD_DOMAIN_NAME == "" || HW_WORKSPACE_AD_SERVER_PWD == "" || HW_WORKSPACE_AD_DOMAIN_IP == "" ||
		HW_WORKSPACE_AD_VPC_ID == "" || HW_WORKSPACE_AD_NETWORK_ID == "" {
		t.Skip("The configuration of AD server is not completed for Workspace service acceptance test.")
	}
}

// lintignore:AT003
func TestAccPreCheckER(t *testing.T) {
	if HW_ER_TEST_ON == "" {
		t.Skip("Skip all ER acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckRfArchives(t *testing.T) {
	if HW_RF_TEMPLATE_ARCHIVE_NO_VARS_URI == "" || HW_RF_TEMPLATE_ARCHIVE_URI == "" ||
		HW_RF_VARIABLES_ARCHIVE_URI == "" {
		t.Skip("Skip the archive URI parameters acceptance test for RF resource stack.")
	}
}

// lintignore:AT003
func TestAccPreCheckDcDirectConnection(t *testing.T) {
	if HW_DC_DIRECT_CONNECT_ID == "" {
		t.Skip("Skip the interface acceptance test because of the direct connection ID is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckCfw(t *testing.T) {
	if HW_CFW_INSTANCE_ID == "" {
		t.Skip("HW_CFW_INSTANCE_ID must be set for CFW acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadType(t *testing.T) {
	if HW_WORKLOAD_TYPE == "" {
		t.Skip("HW_WORKLOAD_TYPE must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadName(t *testing.T) {
	if HW_WORKLOAD_NAME == "" {
		t.Skip("HW_WORKLOAD_NAME must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCceClusterId(t *testing.T) {
	if HW_CCE_CLUSTER_ID == "" {
		t.Skip("HW_CCE_CLUSTER_ID must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckWorkloadNameSpace(t *testing.T) {
	if HW_WORKLOAD_NAMESPACE == "" {
		t.Skip("HW_WORKLOAD_NAMESPACE must be set for SWR image trigger acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSwrTargetRegion(t *testing.T) {
	if HW_SWR_TARGET_REGION == "" {
		t.Skip("HW_SWR_TARGET_REGION must be set for SWR image auto sync tests")
	}
}

// lintignore:AT003
func TestAccPreCheckSwrTargetOrigination(t *testing.T) {
	if HW_SWR_TARGET_ORGANIZATION == "" {
		t.Skip("HW_SWR_TARGET_ORGANIZATION must be set for SWR image auto sync tests")
	}
}

// lintignore:AT003
func TestAccPreCheckImsBackupId(t *testing.T) {
	if HW_IMS_BACKUP_ID == "" {
		t.Skip("HW_IMS_BACKUP_ID must be set for IMS whole image with CBR backup id")
	}
}

// lintignore:AT003
func TestAccPreCheckSourceImage(t *testing.T) {
	if HW_IMAGE_SHARE_SOURCE_IMAGE_ID == "" {
		t.Skip("Skip the interface acceptance test because of the source image ID is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckSecMaster(t *testing.T) {
	if HW_SECMASTER_WORKSPACE_ID == "" {
		t.Skip("HW_SECMASTER_WORKSPACE_ID must be set for SecMaster acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCcePartitionAz(t *testing.T) {
	if HW_CCE_PARTITION_AZ == "" {
		t.Skip("Skip the interface acceptance test because of the cce partition az is missing.")
	}
}

// lintignore:AT003
func TestAccPreCheckCnEast3(t *testing.T) {
	if HW_REGION_NAME != "cn-east-3" {
		t.Skip("HW_REGION_NAME must be cn-east-3 for this test.")
	}
}

// lintignore:AT003
func TestAccPreCheckCertificateWithoutRootCA(t *testing.T) {
	if HW_CERTIFICATE_CONTENT == "" || HW_CERTIFICATE_PRIVATE_KEY == "" ||
		HW_NEW_CERTIFICATE_CONTENT == "" || HW_NEW_CERTIFICATE_PRIVATE_KEY == "" {
		t.Skip("HW_CERTIFICATE_CONTENT, HW_CERTIFICATE_PRIVATE_KEY, HW_NEW_CERTIFICATE_CONTENT and " +
			"HW_NEW_CERTIFICATE_PRIVATE_KEY must be set for simple acceptance tests of SSL certificate resource")
	}
}

// lintignore:AT003
func TestAccPreCheckCertificateFull(t *testing.T) {
	TestAccPreCheckCertificateWithoutRootCA(t)
	if HW_CERTIFICATE_ROOT_CA == "" || HW_NEW_CERTIFICATE_ROOT_CA == "" {
		t.Skip("HW_CERTIFICATE_ROOT_CA and HW_NEW_CERTIFICATE_ROOT_CA must be set for root CA validation")
	}
}

// lintignore:AT003
func TestAccPreCheckCodeArtsDeployResourcePoolID(t *testing.T) {
	if HW_CODEARTS_RESOURCE_POOL_ID == "" {
		t.Skip("HW_CODEARTS_RESOURCE_POOL_ID must be set for this acceptance test")
	}
}

// lintignore:AT003
func TestAccPreCheckCodeArtsEnableFlag(t *testing.T) {
	if HW_CODEARTS_ENABLE_FLAG == "" {
		t.Skip("Skip the CodeArts acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckModelArtsHasSubscribeModel(t *testing.T) {
	if HW_MODELARTS_HAS_SUBSCRIBE_MODEL == "" {
		t.Skip("Subscribe two free models from market and set HW_MODELARTS_HAS_SUBSCRIBE_MODEL" +
			" for modelarts service acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckEgChannelId(t *testing.T) {
	if HW_EG_CHANNEL_ID == "" {
		t.Skip("The sub-resource acceptance test of the EG channel must set 'HW_EG_CHANNEL_ID'")
	}
}

// lintignore:AT003
func TestAccPreCheckLtsAomAccess(t *testing.T) {
	if HW_CCE_CLUSTER_ID == "" || HW_CCE_CLUSTER_NAME == "" {
		t.Skip("HW_CCE_CLUSTER_ID and HW_CCE_CLUSTER_NAME must be set for LTS AOM access acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckLtsAomAccessUpdate(t *testing.T) {
	if HW_CCE_CLUSTER_ID_ANOTHER == "" || HW_CCE_CLUSTER_NAME_ANOTHER == "" {
		t.Skip("HW_CCE_CLUSTER_ID_ANOTHER and HW_CCE_CLUSTER_NAME_ANOTHER must be set for LTS AOM access" +
			" acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckAomSubApplicationId(t *testing.T) {
	if HW_AOM_SUB_APPLICATION_ID == "" {
		t.Skip("HW_AOM_SUB_APPLICATION_ID must be set for the acceptance test")
	}
}

// lintignore:AT003
func TestAccPrecheckKooGallery(t *testing.T) {
	if HW_KOOGALLERY_ASSET == "" {
		t.Skip("Skip the KooGallery acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckLtsStructConfigCustom(t *testing.T) {
	if HW_LTS_STRUCT_CONFIG_TEMPLATE_ID == "" || HW_LTS_STRUCT_CONFIG_TEMPLATE_NAME == "" {
		t.Skip("HW_LTS_STRUCT_CONFIG_TEMPLATE_ID and HW_LTS_STRUCT_CONFIG_TEMPLATE_NAME must be" +
			" set for LTS struct config custom acceptance tests")
	}
}

// lintignore:AT003
func TestAccPreCheckLtsEnableFlag(t *testing.T) {
	if HW_LTS_ENABLE_FLAG == "" {
		t.Skip("Skip the LTS acceptance tests.")
	}
}

// lintignore:AT003
func TestAccPreCheckCCINamespace(t *testing.T) {
	if HW_CCI_NAMESPACE == "" {
		t.Skip("This environment does not support CCI Namespace tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCDN(t *testing.T) {
	if HW_CDN_DOMAIN_NAME == "" {
		t.Skip("This environment does not support CDN tests")
	}
}

// lintignore:AT003
func TestAccPreCheckCERT(t *testing.T) {
	if HW_CDN_CERT_PATH == "" || HW_CDN_PRIVATE_KEY_PATH == "" {
		t.Skip("This environment does not support CDN certificate tests")
	}
}
