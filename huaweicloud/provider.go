package huaweicloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aad"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/antiddos"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apigateway"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/as"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/bms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbh"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ces"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cloudtable"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cmdb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cnad"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codearts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cph"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cpts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dbss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ddm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dis"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/drs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ecs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eg"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eps"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/er"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ga"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/gaussdb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ges"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/identitycenter"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/koogallery"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/live"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/meeting"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/oms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/organizations"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ram"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rfs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/scm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sdrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sfs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/smn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/swr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/tms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ucs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vod"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpcep"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

const (
	defaultCloud       string = "myhuaweicloud.com"
	defaultEuropeCloud string = "myhuaweicloud.eu"
	prefixEuropeRegion string = "eu-west-1"
)

// Provider returns a schema.Provider for HuaweiCloud.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["region"],
				InputDefault: "cn-north-1",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_REGION_NAME",
					"OS_REGION_NAME",
				}, nil),
			},

			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["access_key"],
				RequiredWith: []string{"secret_key"},
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_ACCESS_KEY",
					"OS_ACCESS_KEY",
				}, nil),
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["secret_key"],
				RequiredWith: []string{"access_key"},
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_SECRET_KEY",
					"OS_SECRET_KEY",
				}, nil),
			},

			"security_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["security_token"],
				RequiredWith: []string{"access_key"},
				DefaultFunc:  schema.EnvDefaultFunc("HW_SECURITY_TOKEN", nil),
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_DOMAIN_ID",
					"OS_DOMAIN_ID",
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
				}, ""),
			},

			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
				}, ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_USER_NAME",
					"OS_USERNAME",
				}, ""),
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_USER_ID",
					"OS_USER_ID",
				}, ""),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: descriptions["password"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_USER_PASSWORD",
					"OS_PASSWORD",
				}, ""),
			},

			"assume_role": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_agency_name"],
							DefaultFunc: schema.EnvDefaultFunc("HW_ASSUME_ROLE_AGENCY_NAME", nil),
						},
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_domain_name"],
							DefaultFunc: schema.EnvDefaultFunc("HW_ASSUME_ROLE_DOMAIN_NAME", nil),
						},
					},
				},
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_PROJECT_ID",
					"OS_PROJECT_ID",
				}, nil),
			},

			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_PROJECT_NAME",
					"OS_PROJECT_NAME",
				}, nil),
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_id"],
				DefaultFunc: schema.EnvDefaultFunc("OS_TENANT_ID", ""),
			},

			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_name"],
				DefaultFunc: schema.EnvDefaultFunc("OS_TENANT_NAME", ""),
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["token"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_AUTH_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["insecure"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_INSECURE",
					"OS_INSECURE",
				}, false),
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"agency_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OS_AGENCY_NAME", nil),
				Description:  descriptions["agency_name"],
				RequiredWith: []string{"agency_domain_name"},
			},

			"agency_domain_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OS_AGENCY_DOMAIN_NAME", nil),
				Description:  descriptions["agency_domain_name"],
				RequiredWith: []string{"agency_name"},
			},

			"delegated_project": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELEGATED_PROJECT", ""),
				Description: descriptions["delegated_project"],
			},

			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["auth_url"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_AUTH_URL",
					"OS_AUTH_URL",
				}, nil),
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["cloud"],
				DefaultFunc: schema.EnvDefaultFunc("HW_CLOUD", ""),
			},

			"endpoints": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["endpoints"],
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"regional": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["regional"],
			},

			"shared_config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_config_file"],
				DefaultFunc: schema.EnvDefaultFunc("HW_SHARED_CONFIG_FILE", ""),
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("HW_PROFILE", ""),
			},

			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["enterprise_project_id"],
				DefaultFunc: schema.EnvDefaultFunc("HW_ENTERPRISE_PROJECT_ID", ""),
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: descriptions["max_retries"],
				DefaultFunc: schema.EnvDefaultFunc("HW_MAX_RETRIES", 5),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"huaweicloud_apig_environments": apig.DataSourceEnvironments(),
			"huaweicloud_apig_groups":       apig.DataSourceGroups(),

			"huaweicloud_as_configurations": as.DataSourceASConfigurations(),
			"huaweicloud_as_groups":         as.DataSourceASGroups(),

			"huaweicloud_account":            DataSourceAccount(),
			"huaweicloud_availability_zones": DataSourceAvailabilityZones(),

			"huaweicloud_bms_flavors": bms.DataSourceBmsFlavors(),

			"huaweicloud_cbr_backup": cbr.DataSourceBackup(),
			"huaweicloud_cbr_vaults": cbr.DataSourceVaults(),

			"huaweicloud_cbh_instances": cbh.DataSourceCbhInstances(),

			"huaweicloud_cce_addon_template": cce.DataSourceAddonTemplate(),
			"huaweicloud_cce_cluster":        cce.DataSourceCCEClusterV3(),
			"huaweicloud_cce_clusters":       cce.DataSourceCCEClusters(),
			"huaweicloud_cce_node":           cce.DataSourceNode(),
			"huaweicloud_cce_nodes":          cce.DataSourceNodes(),
			"huaweicloud_cce_node_pool":      cce.DataSourceCCENodePoolV3(),
			"huaweicloud_cci_namespaces":     cci.DataSourceCciNamespaces(),

			"huaweicloud_cdm_flavors": DataSourceCdmFlavorV1(),

			"huaweicloud_cdn_domain_statistics": cdn.DataSourceStatistics(),

			"huaweicloud_cfw_firewalls": cfw.DataSourceFirewalls(),

			"huaweicloud_cnad_advanced_instances":         cnad.DataSourceInstances(),
			"huaweicloud_cnad_advanced_available_objects": cnad.DataSourceAvailableProtectedObjects(),
			"huaweicloud_cnad_advanced_protected_objects": cnad.DataSourceProtectedObjects(),

			"huaweicloud_compute_flavors":      ecs.DataSourceEcsFlavors(),
			"huaweicloud_compute_instance":     ecs.DataSourceComputeInstance(),
			"huaweicloud_compute_instances":    ecs.DataSourceComputeInstances(),
			"huaweicloud_compute_servergroups": ecs.DataSourceComputeServerGroups(),

			"huaweicloud_cdm_clusters": cdm.DataSourceCdmClusters(),

			"huaweicloud_cph_server_flavors": cph.DataSourceServerFlavors(),
			"huaweicloud_cph_phone_flavors":  cph.DataSourcePhoneFlavors(),
			"huaweicloud_cph_phone_images":   cph.DataSourcePhoneImages(),

			"huaweicloud_csms_secret_version": dew.DataSourceDewCsmsSecret(),
			"huaweicloud_css_flavors":         css.DataSourceCssFlavors(),

			"huaweicloud_dcs_flavors":         dcs.DataSourceDcsFlavorsV2(),
			"huaweicloud_dcs_maintainwindow":  dcs.DataSourceDcsMaintainWindow(),
			"huaweicloud_dcs_instances":       dcs.DataSourceDcsInstance(),
			"huaweicloud_dcs_templates":       dcs.DataSourceTemplates(),
			"huaweicloud_dcs_template_detail": dcs.DataSourceTemplateDetail(),

			"huaweicloud_dds_flavors":   dds.DataSourceDDSFlavorV3(),
			"huaweicloud_dds_instances": dds.DataSourceDdsInstance(),

			"huaweicloud_dms_kafka_flavors":   dms.DataSourceKafkaFlavors(),
			"huaweicloud_dms_kafka_instances": dms.DataSourceDmsKafkaInstances(),
			"huaweicloud_dms_product":         dms.DataSourceDmsProduct(),
			"huaweicloud_dms_maintainwindow":  dms.DataSourceDmsMaintainWindow(),

			"huaweicloud_dms_rabbitmq_flavors": dms.DataSourceRabbitMQFlavors(),

			"huaweicloud_dms_rocketmq_broker":    dms.DataSourceDmsRocketMQBroker(),
			"huaweicloud_dms_rocketmq_instances": dms.DataSourceDmsRocketMQInstances(),

			"huaweicloud_dns_zones":      dns.DataSourceZones(),
			"huaweicloud_dns_recordsets": dns.DataSourceRecordsets(),

			"huaweicloud_eg_custom_event_channels": eg.DataSourceCustomEventChannels(),
			"huaweicloud_eg_custom_event_sources":  eg.DataSourceCustomEventSources(),

			"huaweicloud_enterprise_project": eps.DataSourceEnterpriseProject(),

			"huaweicloud_er_attachments":  er.DataSourceAttachments(),
			"huaweicloud_er_instances":    er.DataSourceInstances(),
			"huaweicloud_er_route_tables": er.DataSourceRouteTables(),

			"huaweicloud_evs_volumes":      evs.DataSourceEvsVolumesV2(),
			"huaweicloud_fgs_dependencies": fgs.DataSourceFunctionGraphDependencies(),
			"huaweicloud_fgs_functions":    fgs.DataSourceFunctionGraphFunctions(),

			"huaweicloud_gaussdb_cassandra_dedicated_resource": gaussdb.DataSourceGeminiDBDehResource(),
			"huaweicloud_gaussdb_cassandra_flavors":            gaussdb.DataSourceCassandraFlavors(),
			"huaweicloud_gaussdb_nosql_flavors":                gaussdb.DataSourceGaussDBNoSQLFlavors(),
			"huaweicloud_gaussdb_cassandra_instance":           gaussdb.DataSourceGeminiDBInstance(),
			"huaweicloud_gaussdb_cassandra_instances":          gaussdb.DataSourceGeminiDBInstances(),
			"huaweicloud_gaussdb_opengauss_instance":           gaussdb.DataSourceOpenGaussInstance(),
			"huaweicloud_gaussdb_opengauss_instances":          gaussdb.DataSourceOpenGaussInstances(),
			"huaweicloud_gaussdb_mysql_configuration":          gaussdb.DataSourceGaussdbMysqlConfigurations(),
			"huaweicloud_gaussdb_mysql_dedicated_resource":     gaussdb.DataSourceGaussDBMysqlDehResource(),
			"huaweicloud_gaussdb_mysql_flavors":                gaussdb.DataSourceGaussdbMysqlFlavors(),
			"huaweicloud_gaussdb_mysql_instance":               gaussdb.DataSourceGaussDBMysqlInstance(),
			"huaweicloud_gaussdb_mysql_instances":              gaussdb.DataSourceGaussDBMysqlInstances(),
			"huaweicloud_gaussdb_redis_instance":               gaussdb.DataSourceGaussRedisInstance(),

			"huaweicloud_identity_permissions": iam.DataSourceIdentityPermissions(),
			"huaweicloud_identity_role":        iam.DataSourceIdentityRole(),
			"huaweicloud_identity_custom_role": iam.DataSourceIdentityCustomRole(),
			"huaweicloud_identity_group":       iam.DataSourceIdentityGroup(),
			"huaweicloud_identity_projects":    iam.DataSourceIdentityProjects(),
			"huaweicloud_identity_users":       iam.DataSourceIdentityUsers(),

			"huaweicloud_identitycenter_instance": identitycenter.DataSourceIdentityCenter(),
			"huaweicloud_identitycenter_groups":   identitycenter.DataSourceIdentityCenterGroups(),
			"huaweicloud_identitycenter_users":    identitycenter.DataSourceIdentityCenterUsers(),

			"huaweicloud_iec_bandwidths":     dataSourceIECBandWidths(),
			"huaweicloud_iec_eips":           dataSourceIECNetworkEips(),
			"huaweicloud_iec_flavors":        dataSourceIecFlavors(),
			"huaweicloud_iec_images":         dataSourceIecImages(),
			"huaweicloud_iec_keypair":        dataSourceIECKeypair(),
			"huaweicloud_iec_network_acl":    dataSourceIECNetworkACL(),
			"huaweicloud_iec_port":           DataSourceIECPort(),
			"huaweicloud_iec_security_group": dataSourceIECSecurityGroup(),
			"huaweicloud_iec_server":         dataSourceIECServer(),
			"huaweicloud_iec_sites":          dataSourceIecSites(),
			"huaweicloud_iec_vpc":            DataSourceIECVpc(),
			"huaweicloud_iec_vpc_subnets":    DataSourceIECVpcSubnets(),

			"huaweicloud_images_image":  ims.DataSourceImagesImageV2(),
			"huaweicloud_images_images": ims.DataSourceImagesImages(),

			"huaweicloud_kms_key":      dew.DataSourceKmsKey(),
			"huaweicloud_kms_data_key": dew.DataSourceKmsDataKeyV1(),
			"huaweicloud_kps_keypairs": dew.DataSourceKeypairs(),

			"huaweicloud_koogallery_assets": koogallery.DataSourceKooGalleryAssets(),

			"huaweicloud_lb_listeners":    lb.DataSourceListeners(),
			"huaweicloud_lb_loadbalancer": lb.DataSourceELBV2Loadbalancer(),
			"huaweicloud_lb_certificate":  lb.DataSourceLBCertificateV2(),
			"huaweicloud_lb_pools":        lb.DataSourcePools(),

			"huaweicloud_lts_structuring_custom_templates": lts.DataSourceCustomTemplates(),

			"huaweicloud_elb_certificate":   elb.DataSourceELBCertificateV3(),
			"huaweicloud_elb_flavors":       elb.DataSourceElbFlavorsV3(),
			"huaweicloud_elb_pools":         elb.DataSourcePools(),
			"huaweicloud_elb_loadbalancers": elb.DataSourceElbLoadbalances(),

			"huaweicloud_nat_gateway": nat.DataSourcePublicGateway(),

			"huaweicloud_networking_port":      vpc.DataSourceNetworkingPortV2(),
			"huaweicloud_networking_secgroup":  vpc.DataSourceNetworkingSecGroup(),
			"huaweicloud_networking_secgroups": vpc.DataSourceNetworkingSecGroups(),

			"huaweicloud_mapreduce_versions": mrs.DataSourceMrsVersions(),

			"huaweicloud_modelarts_datasets":         modelarts.DataSourceDatasets(),
			"huaweicloud_modelarts_dataset_versions": modelarts.DataSourceDatasetVerions(),
			"huaweicloud_modelarts_notebook_images":  modelarts.DataSourceNotebookImages(),
			"huaweicloud_modelarts_notebook_flavors": modelarts.DataSourceNotebookFlavors(),
			"huaweicloud_modelarts_service_flavors":  modelarts.DataSourceServiceFlavors(),
			"huaweicloud_modelarts_models":           modelarts.DataSourceModels(),
			"huaweicloud_modelarts_model_templates":  modelarts.DataSourceModelTemplates(),
			"huaweicloud_modelarts_workspaces":       modelarts.DataSourceWorkspaces(),
			"huaweicloud_modelarts_services":         modelarts.DataSourceServices(),
			"huaweicloud_modelarts_resource_flavors": modelarts.DataSourceResourceFlavors(),

			"huaweicloud_mapreduce_clusters": mrs.DataSourceMrsClusters(),

			"huaweicloud_obs_buckets":       obs.DataSourceObsBuckets(),
			"huaweicloud_obs_bucket_object": obs.DataSourceObsBucketObject(),

			"huaweicloud_ram_resource_permissions": ram.DataSourceRAMPermissions(),

			"huaweicloud_rds_flavors":              rds.DataSourceRdsFlavor(),
			"huaweicloud_rds_engine_versions":      rds.DataSourceRdsEngineVersionsV3(),
			"huaweicloud_rds_instances":            rds.DataSourceRdsInstances(),
			"huaweicloud_rds_backups":              rds.DataSourceBackup(),
			"huaweicloud_rds_storage_types":        rds.DataSourceStoragetype(),
			"huaweicloud_rds_sqlserver_collations": rds.DataSourceSQLServerCollations(),

			"huaweicloud_rms_policy_definitions":           rms.DataSourcePolicyDefinitions(),
			"huaweicloud_rms_assignment_package_templates": rms.DataSourceTemplates(),

			"huaweicloud_sdrs_domain": sdrs.DataSourceSDRSDomain(),

			"huaweicloud_servicestage_component_runtimes": servicestage.DataSourceComponentRuntimes(),

			"huaweicloud_smn_topics":            smn.DataSourceTopics(),
			"huaweicloud_smn_message_templates": smn.DataSourceSmnMessageTemplates(),

			"huaweicloud_sms_source_servers": sms.DataSourceServers(),

			"huaweicloud_scm_certificates": scm.DataSourceCertificates(),

			"huaweicloud_sfs_file_system": sfs.DataSourceSFSFileSystemV2(),
			"huaweicloud_sfs_turbos":      sfs.DataSourceTurbos(),

			"huaweicloud_tms_resource_types": tms.DataSourceResourceTypes(),

			"huaweicloud_vpc_bandwidth": eip.DataSourceBandWidth(),
			"huaweicloud_vpc_eip":       eip.DataSourceVpcEip(),
			"huaweicloud_vpc_eips":      eip.DataSourceVpcEips(),

			"huaweicloud_vpc":                    vpc.DataSourceVpcV1(),
			"huaweicloud_vpcs":                   vpc.DataSourceVpcs(),
			"huaweicloud_vpc_ids":                vpc.DataSourceVpcIdsV1(),
			"huaweicloud_vpc_peering_connection": vpc.DataSourceVpcPeeringConnectionV2(),
			"huaweicloud_vpc_route_table":        vpc.DataSourceVPCRouteTable(),
			"huaweicloud_vpc_subnet":             vpc.DataSourceVpcSubnetV1(),
			"huaweicloud_vpc_subnets":            vpc.DataSourceVpcSubnets(),
			"huaweicloud_vpc_subnet_ids":         vpc.DataSourceVpcSubnetIdsV1(),

			"huaweicloud_vpcep_public_services": vpcep.DataSourceVPCEPPublicServices(),

			"huaweicloud_vpn_gateway_availability_zones": vpn.DataSourceVpnGatewayAZs(),
			"huaweicloud_vpn_gateways":                   vpn.DataSourceGateways(),

			"huaweicloud_waf_certificate":         waf.DataSourceWafCertificateV1(),
			"huaweicloud_waf_policies":            waf.DataSourceWafPoliciesV1(),
			"huaweicloud_waf_dedicated_instances": waf.DataSourceWafDedicatedInstancesV1(),
			"huaweicloud_waf_reference_tables":    waf.DataSourceWafReferenceTablesV1(),
			"huaweicloud_waf_instance_groups":     waf.DataSourceWafInstanceGroups(),
			"huaweicloud_dws_flavors":             dws.DataSourceDwsFlavors(),

			// Legacy
			"huaweicloud_images_image_v2":        ims.DataSourceImagesImageV2(),
			"huaweicloud_networking_port_v2":     vpc.DataSourceNetworkingPortV2(),
			"huaweicloud_networking_secgroup_v2": vpc.DataSourceNetworkingSecGroup(),

			"huaweicloud_kms_key_v1":      dew.DataSourceKmsKey(),
			"huaweicloud_kms_data_key_v1": dew.DataSourceKmsDataKeyV1(),

			"huaweicloud_rds_flavors_v3":     rds.DataSourceRdsFlavor(),
			"huaweicloud_sfs_file_system_v2": sfs.DataSourceSFSFileSystemV2(),

			"huaweicloud_vpc_v1":                    vpc.DataSourceVpcV1(),
			"huaweicloud_vpc_ids_v1":                vpc.DataSourceVpcIdsV1(),
			"huaweicloud_vpc_peering_connection_v2": vpc.DataSourceVpcPeeringConnectionV2(),
			"huaweicloud_vpc_subnet_v1":             vpc.DataSourceVpcSubnetV1(),
			"huaweicloud_vpc_subnet_ids_v1":         vpc.DataSourceVpcSubnetIdsV1(),

			"huaweicloud_cce_cluster_v3": cce.DataSourceCCEClusterV3(),
			"huaweicloud_cce_node_v3":    cce.DataSourceNode(),

			"huaweicloud_dms_product_v1":        dms.DataSourceDmsProduct(),
			"huaweicloud_dms_maintainwindow_v1": dms.DataSourceDmsMaintainWindow(),

			"huaweicloud_dcs_maintainwindow_v1": dcs.DataSourceDcsMaintainWindow(),

			"huaweicloud_dds_flavors_v3":   dds.DataSourceDDSFlavorV3(),
			"huaweicloud_identity_role_v3": iam.DataSourceIdentityRole(),
			"huaweicloud_cdm_flavors_v1":   DataSourceCdmFlavorV1(),

			"huaweicloud_ddm_engines":        ddm.DataSourceDdmEngines(),
			"huaweicloud_ddm_flavors":        ddm.DataSourceDdmFlavors(),
			"huaweicloud_ddm_instance_nodes": ddm.DataSourceDdmInstanceNodes(),
			"huaweicloud_ddm_instances":      ddm.DataSourceDdmInstances(),
			"huaweicloud_ddm_schemas":        ddm.DataSourceDdmSchemas(),
			"huaweicloud_ddm_accounts":       ddm.DataSourceDdmAccounts(),

			"huaweicloud_organizations_organization":         organizations.DataSourceOrganization(),
			"huaweicloud_organizations_organizational_units": organizations.DataSourceOrganizationalUnits(),
			"huaweicloud_organizations_accounts":             organizations.DataSourceAccounts(),
			"huaweicloud_organizations_policies":             organizations.DataSourcePolicies(),

			// Deprecated ongoing (without DeprecationMessage), used by other providers
			"huaweicloud_vpc_route":        vpc.DataSourceVpcRouteV2(),
			"huaweicloud_vpc_route_ids":    vpc.DataSourceVpcRouteIdsV2(),
			"huaweicloud_vpc_route_v2":     vpc.DataSourceVpcRouteV2(),
			"huaweicloud_vpc_route_ids_v2": vpc.DataSourceVpcRouteIdsV2(),

			// Deprecated
			"huaweicloud_antiddos":                      deprecated.DataSourceAntiDdosV1(),
			"huaweicloud_antiddos_v1":                   deprecated.DataSourceAntiDdosV1(),
			"huaweicloud_compute_availability_zones_v2": deprecated.DataSourceComputeAvailabilityZonesV2(),
			"huaweicloud_csbs_backup":                   deprecated.DataSourceCSBSBackupV1(),
			"huaweicloud_csbs_backup_policy":            deprecated.DataSourceCSBSBackupPolicyV1(),
			"huaweicloud_csbs_backup_policy_v1":         deprecated.DataSourceCSBSBackupPolicyV1(),
			"huaweicloud_csbs_backup_v1":                deprecated.DataSourceCSBSBackupV1(),
			"huaweicloud_networking_network_v2":         deprecated.DataSourceNetworkingNetworkV2(),
			"huaweicloud_networking_subnet_v2":          deprecated.DataSourceNetworkingSubnetV2(),
			"huaweicloud_cts_tracker":                   deprecated.DataSourceCTSTrackerV1(),
			"huaweicloud_dcs_az":                        deprecated.DataSourceDcsAZV1(),
			"huaweicloud_dcs_az_v1":                     deprecated.DataSourceDcsAZV1(),
			"huaweicloud_dcs_product":                   deprecated.DataSourceDcsProductV1(),
			"huaweicloud_dcs_product_v1":                deprecated.DataSourceDcsProductV1(),
			"huaweicloud_dms_az":                        deprecated.DataSourceDmsAZ(),
			"huaweicloud_dms_az_v1":                     deprecated.DataSourceDmsAZ(),
			"huaweicloud_vbs_backup_policy":             deprecated.DataSourceVBSBackupPolicyV2(),
			"huaweicloud_vbs_backup":                    deprecated.DataSourceVBSBackupV2(),
			"huaweicloud_vbs_backup_policy_v2":          deprecated.DataSourceVBSBackupPolicyV2(),
			"huaweicloud_vbs_backup_v2":                 deprecated.DataSourceVBSBackupV2(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"huaweicloud_aad_forward_rule": aad.ResourceForwardRule(),

			"huaweicloud_antiddos_basic": antiddos.ResourceCloudNativeAntiDdos(),

			"huaweicloud_aom_alarm_rule":             aom.ResourceAlarmRule(),
			"huaweicloud_aom_event_alarm_rule":       aom.ResourceEventAlarmRule(),
			"huaweicloud_aom_service_discovery_rule": aom.ResourceServiceDiscoveryRule(),
			"huaweicloud_aom_alarm_action_rule":      aom.ResourceAlarmActionRule(),
			"huaweicloud_aom_alarm_silence_rule":     aom.ResourceAlarmSilenceRule(),

			"huaweicloud_aom_cmdb_application": aom.ResourceCmdbApplication(),
			"huaweicloud_aom_cmdb_component":   aom.ResourceCmdbComponent(),
			"huaweicloud_aom_cmdb_environment": aom.ResourceCmdbEnvironment(),

			"huaweicloud_rfs_stack": rfs.ResourceStack(),

			"huaweicloud_api_gateway_api":         ResourceAPIGatewayAPI(),
			"huaweicloud_api_gateway_environment": apigateway.ResourceEnvironment(),
			"huaweicloud_api_gateway_group":       ResourceAPIGatewayGroup(),

			"huaweicloud_apig_acl_policy":                  apig.ResourceAclPolicy(),
			"huaweicloud_apig_acl_policy_associate":        apig.ResourceAclPolicyAssociate(),
			"huaweicloud_apig_api":                         apig.ResourceApigAPIV2(),
			"huaweicloud_apig_api_publishment":             apig.ResourceApigApiPublishment(),
			"huaweicloud_apig_appcode":                     apig.ResourceAppcode(),
			"huaweicloud_apig_application":                 apig.ResourceApigApplicationV2(),
			"huaweicloud_apig_application_authorization":   apig.ResourceAppAuth(),
			"huaweicloud_apig_certificate":                 apig.ResourceCertificate(),
			"huaweicloud_apig_channel":                     apig.ResourceChannel(),
			"huaweicloud_apig_custom_authorizer":           apig.ResourceApigCustomAuthorizerV2(),
			"huaweicloud_apig_environment":                 apig.ResourceApigEnvironmentV2(),
			"huaweicloud_apig_group":                       apig.ResourceApigGroupV2(),
			"huaweicloud_apig_instance_routes":             apig.ResourceInstanceRoutes(),
			"huaweicloud_apig_instance":                    apig.ResourceApigInstanceV2(),
			"huaweicloud_apig_plugin_associate":            apig.ResourcePluginAssociate(),
			"huaweicloud_apig_plugin":                      apig.ResourcePlugin(),
			"huaweicloud_apig_response":                    apig.ResourceApigResponseV2(),
			"huaweicloud_apig_signature_associate":         apig.ResourceSignatureAssociate(),
			"huaweicloud_apig_signature":                   apig.ResourceSignature(),
			"huaweicloud_apig_throttling_policy_associate": apig.ResourceThrottlingPolicyAssociate(),
			"huaweicloud_apig_throttling_policy":           apig.ResourceApigThrottlingPolicyV2(),

			"huaweicloud_as_configuration":    as.ResourceASConfiguration(),
			"huaweicloud_as_group":            as.ResourceASGroup(),
			"huaweicloud_as_lifecycle_hook":   as.ResourceASLifecycleHook(),
			"huaweicloud_as_instance_attach":  as.ResourceASInstanceAttach(),
			"huaweicloud_as_notification":     as.ResourceAsNotification(),
			"huaweicloud_as_policy":           as.ResourceASPolicy(),
			"huaweicloud_as_bandwidth_policy": as.ResourceASBandWidthPolicy(),

			"huaweicloud_bms_instance": bms.ResourceBmsInstance(),
			"huaweicloud_bcs_instance": resourceBCSInstanceV2(),

			"huaweicloud_cbr_checkpoint": cbr.ResourceCheckpoint(),
			"huaweicloud_cbr_policy":     cbr.ResourcePolicy(),
			"huaweicloud_cbr_vault":      cbr.ResourceVault(),

			"huaweicloud_cbh_instance": cbh.ResourceCBHInstance(),

			"huaweicloud_cc_connection":             cc.ResourceCloudConnection(),
			"huaweicloud_cc_network_instance":       cc.ResourceNetworkInstance(),
			"huaweicloud_cc_bandwidth_package":      cc.ResourceBandwidthPackage(),
			"huaweicloud_cc_inter_region_bandwidth": cc.ResourceInterRegionBandwidth(),

			"huaweicloud_cce_cluster":     cce.ResourceCluster(),
			"huaweicloud_cce_node":        cce.ResourceNode(),
			"huaweicloud_cce_node_attach": cce.ResourceNodeAttach(),
			"huaweicloud_cce_addon":       cce.ResourceAddon(),
			"huaweicloud_cce_node_pool":   cce.ResourceNodePool(),
			"huaweicloud_cce_namespace":   cce.ResourceCCENamespaceV1(),
			"huaweicloud_cce_pvc":         cce.ResourceCcePersistentVolumeClaimsV1(),
			"huaweicloud_cce_partition":   cce.ResourcePartition(),

			"huaweicloud_cts_tracker":      cts.ResourceCTSTracker(),
			"huaweicloud_cts_data_tracker": cts.ResourceCTSDataTracker(),
			"huaweicloud_cts_notification": cts.ResourceCTSNotification(),
			"huaweicloud_cci_namespace":    cci.ResourceCciNamespace(),
			"huaweicloud_cci_network":      cci.ResourceCciNetworkV1(),
			"huaweicloud_cci_pvc":          cci.ResourcePersistentVolumeClaimV1(),

			"huaweicloud_cdm_cluster": cdm.ResourceCdmCluster(),
			"huaweicloud_cdm_job":     cdm.ResourceCdmJob(),
			"huaweicloud_cdm_link":    cdm.ResourceCdmLink(),

			"huaweicloud_cdn_domain":         cdn.ResourceCdnDomainV1(),
			"huaweicloud_ces_alarmrule":      ces.ResourceAlarmRule(),
			"huaweicloud_ces_resource_group": ces.ResourceResourceGroup(),
			"huaweicloud_ces_alarm_template": ces.ResourceCesAlarmTemplate(),

			"huaweicloud_cfw_address_group":        cfw.ResourceAddressGroup(),
			"huaweicloud_cfw_address_group_member": cfw.ResourceAddressGroupMember(),
			"huaweicloud_cfw_black_white_list":     cfw.ResourceBlackWhiteList(),
			"huaweicloud_cfw_eip_protection":       cfw.ResourceEipProtection(),
			"huaweicloud_cfw_protection_rule":      cfw.ResourceProtectionRule(),
			"huaweicloud_cfw_service_group":        cfw.ResourceServiceGroup(),
			"huaweicloud_cfw_service_group_member": cfw.ResourceServiceGroupMember(),

			"huaweicloud_cloudtable_cluster": cloudtable.ResourceCloudTableCluster(),

			"huaweicloud_cnad_advanced_black_white_list": cnad.ResourceBlackWhiteList(),
			"huaweicloud_cnad_advanced_policy":           cnad.ResourceCNADAdvancedPolicy(),
			"huaweicloud_cnad_advanced_policy_associate": cnad.ResourcePolicyAssociate(),
			"huaweicloud_cnad_advanced_protected_object": cnad.ResourceProtectedObject(),

			"huaweicloud_compute_instance":         ecs.ResourceComputeInstance(),
			"huaweicloud_compute_interface_attach": ecs.ResourceComputeInterfaceAttach(),
			"huaweicloud_compute_keypair":          ResourceComputeKeypairV2(),
			"huaweicloud_compute_servergroup":      ecs.ResourceComputeServerGroup(),
			"huaweicloud_compute_eip_associate":    ecs.ResourceComputeEIPAssociate(),
			"huaweicloud_compute_volume_attach":    ecs.ResourceComputeVolumeAttach(),

			"huaweicloud_cph_server": cph.ResourceCphServer(),

			"huaweicloud_cse_microservice":          cse.ResourceMicroservice(),
			"huaweicloud_cse_microservice_engine":   cse.ResourceMicroserviceEngine(),
			"huaweicloud_cse_microservice_instance": cse.ResourceMicroserviceInstance(),

			"huaweicloud_csms_secret": dew.ResourceCsmsSecret(),

			"huaweicloud_css_cluster":       css.ResourceCssCluster(),
			"huaweicloud_css_snapshot":      css.ResourceCssSnapshot(),
			"huaweicloud_css_thesaurus":     css.ResourceCssthesaurus(),
			"huaweicloud_css_configuration": css.ResourceCssConfiguration(),

			"huaweicloud_dbss_instance": dbss.ResourceInstance(),

			"huaweicloud_dc_virtual_gateway":   dc.ResourceVirtualGateway(),
			"huaweicloud_dc_virtual_interface": dc.ResourceVirtualInterface(),

			"huaweicloud_dcs_instance":        dcs.ResourceDcsInstance(),
			"huaweicloud_dcs_backup":          dcs.ResourceDcsBackup(),
			"huaweicloud_dcs_custom_template": dcs.ResourceCustomTemplate(),

			"huaweicloud_dds_database_role":      dds.ResourceDatabaseRole(),
			"huaweicloud_dds_database_user":      dds.ResourceDatabaseUser(),
			"huaweicloud_dds_instance":           dds.ResourceDdsInstanceV3(),
			"huaweicloud_dds_backup":             dds.ResourceDdsBackup(),
			"huaweicloud_dds_parameter_template": dds.ResourceDdsParameterTemplate(),
			"huaweicloud_dds_audit_log_policy":   dds.ResourceDdsAuditLogPolicy(),

			"huaweicloud_ddm_instance": ddm.ResourceDdmInstance(),
			"huaweicloud_ddm_schema":   ddm.ResourceDdmSchema(),
			"huaweicloud_ddm_account":  ddm.ResourceDdmAccount(),

			"huaweicloud_dis_stream": dis.ResourceDisStream(),

			"huaweicloud_dli_database":              dli.ResourceDliSqlDatabaseV1(),
			"huaweicloud_dli_package":               dli.ResourceDliPackageV2(),
			"huaweicloud_dli_queue":                 dli.ResourceDliQueue(),
			"huaweicloud_dli_spark_job":             dli.ResourceDliSparkJobV2(),
			"huaweicloud_dli_sql_job":               dli.ResourceSqlJob(),
			"huaweicloud_dli_table":                 dli.ResourceDliTable(),
			"huaweicloud_dli_flinksql_job":          dli.ResourceFlinkSqlJob(),
			"huaweicloud_dli_flinkjar_job":          dli.ResourceFlinkJarJob(),
			"huaweicloud_dli_permission":            dli.ResourceDliPermission(),
			"huaweicloud_dli_datasource_connection": dli.ResourceDatasourceConnection(),
			"huaweicloud_dli_datasource_auth":       dli.ResourceDatasourceAuth(),
			"huaweicloud_dli_template_sql":          dli.ResourceSQLTemplate(),
			"huaweicloud_dli_template_flink":        dli.ResourceFlinkTemplate(),
			"huaweicloud_dli_global_variable":       dli.ResourceGlobalVariable(),
			"huaweicloud_dli_template_spark":        dli.ResourceSparkTemplate(),
			"huaweicloud_dli_agency":                dli.ResourceDliAgency(),

			"huaweicloud_dms_kafka_user":        dms.ResourceDmsKafkaUser(),
			"huaweicloud_dms_kafka_permissions": dms.ResourceDmsKafkaPermissions(),
			"huaweicloud_dms_kafka_instance":    dms.ResourceDmsKafkaInstance(),
			"huaweicloud_dms_kafka_topic":       dms.ResourceDmsKafkaTopic(),
			"huaweicloud_dms_rabbitmq_instance": dms.ResourceDmsRabbitmqInstance(),

			"huaweicloud_dms_rocketmq_instance":       dms.ResourceDmsRocketMQInstance(),
			"huaweicloud_dms_rocketmq_consumer_group": dms.ResourceDmsRocketMQConsumerGroup(),
			"huaweicloud_dms_rocketmq_topic":          dms.ResourceDmsRocketMQTopic(),
			"huaweicloud_dms_rocketmq_user":           dms.ResourceDmsRocketMQUser(),

			"huaweicloud_dns_custom_line": dns.ResourceDNSCustomLine(),
			"huaweicloud_dns_ptrrecord":   dns.ResourceDNSPtrRecord(),
			"huaweicloud_dns_recordset":   dns.ResourceDNSRecordset(),
			"huaweicloud_dns_zone":        dns.ResourceDNSZone(),

			"huaweicloud_drs_job": drs.ResourceDrsJob(),

			"huaweicloud_dws_cluster":            dws.ResourceDwsCluster(),
			"huaweicloud_dws_event_subscription": dws.ResourceDwsEventSubs(),
			"huaweicloud_dws_alarm_subscription": dws.ResourceDwsAlarmSubs(),
			"huaweicloud_dws_snapshot":           dws.ResourceDwsSnapshot(),
			"huaweicloud_dws_snapshot_policy":    dws.ResourceDwsSnapshotPolicy(),
			"huaweicloud_dws_ext_data_source":    dws.ResourceDwsExtDataSource(),

			"huaweicloud_eg_connection":           eg.ResourceConnection(),
			"huaweicloud_eg_custom_event_channel": eg.ResourceCustomEventChannel(),
			"huaweicloud_eg_custom_event_source":  eg.ResourceCustomEventSource(),
			"huaweicloud_eg_endpoint":             eg.ResourceEndpoint(),
			"huaweicloud_eg_event_subscription":   eg.ResourceEventSubscription(),

			"huaweicloud_elb_certificate":     elb.ResourceCertificateV3(),
			"huaweicloud_elb_l7policy":        elb.ResourceL7PolicyV3(),
			"huaweicloud_elb_l7rule":          elb.ResourceL7RuleV3(),
			"huaweicloud_elb_listener":        elb.ResourceListenerV3(),
			"huaweicloud_elb_loadbalancer":    elb.ResourceLoadBalancerV3(),
			"huaweicloud_elb_monitor":         elb.ResourceMonitorV3(),
			"huaweicloud_elb_ipgroup":         elb.ResourceIpGroupV3(),
			"huaweicloud_elb_pool":            elb.ResourcePoolV3(),
			"huaweicloud_elb_member":          elb.ResourceMemberV3(),
			"huaweicloud_elb_logtank":         elb.ResourceLogTank(),
			"huaweicloud_elb_security_policy": elb.ResourceSecurityPolicy(),

			"huaweicloud_enterprise_project": eps.ResourceEnterpriseProject(),

			"huaweicloud_er_association":    er.ResourceAssociation(),
			"huaweicloud_er_instance":       er.ResourceInstance(),
			"huaweicloud_er_propagation":    er.ResourcePropagation(),
			"huaweicloud_er_route_table":    er.ResourceRouteTable(),
			"huaweicloud_er_static_route":   er.ResourceStaticRoute(),
			"huaweicloud_er_vpc_attachment": er.ResourceVpcAttachment(),

			"huaweicloud_evs_snapshot": evs.ResourceEvsSnapshotV2(),
			"huaweicloud_evs_volume":   evs.ResourceEvsVolume(),

			"huaweicloud_fgs_async_invoke_configuration": fgs.ResourceAsyncInvokeConfiguration(),
			"huaweicloud_fgs_dependency":                 fgs.ResourceFgsDependency(),
			"huaweicloud_fgs_function":                   fgs.ResourceFgsFunctionV2(),
			"huaweicloud_fgs_trigger":                    fgs.ResourceFunctionGraphTrigger(),

			"huaweicloud_ga_accelerator":    ga.ResourceAccelerator(),
			"huaweicloud_ga_listener":       ga.ResourceListener(),
			"huaweicloud_ga_endpoint_group": ga.ResourceEndpointGroup(),
			"huaweicloud_ga_endpoint":       ga.ResourceEndpoint(),
			"huaweicloud_ga_health_check":   ga.ResourceHealthCheck(),

			"huaweicloud_gaussdb_cassandra_instance": gaussdb.ResourceGeminiDBInstanceV3(),

			"huaweicloud_gaussdb_mysql_instance":           gaussdb.ResourceGaussDBInstance(),
			"huaweicloud_gaussdb_mysql_proxy":              gaussdb.ResourceGaussDBProxy(),
			"huaweicloud_gaussdb_mysql_database":           gaussdb.ResourceGaussDBDatabase(),
			"huaweicloud_gaussdb_mysql_account":            gaussdb.ResourceGaussDBAccount(),
			"huaweicloud_gaussdb_mysql_account_privilege":  gaussdb.ResourceGaussDBAccountPrivilege(),
			"huaweicloud_gaussdb_mysql_sql_control_rule":   gaussdb.ResourceGaussDBSqlControlRule(),
			"huaweicloud_gaussdb_mysql_parameter_template": gaussdb.ResourceGaussDBMysqlTemplate(),

			"huaweicloud_gaussdb_opengauss_instance": gaussdb.ResourceOpenGaussInstance(),

			"huaweicloud_gaussdb_redis_instance":      gaussdb.ResourceGaussRedisInstanceV3(),
			"huaweicloud_gaussdb_redis_eip_associate": gaussdb.ResourceGaussRedisEipAssociate(),

			"huaweicloud_gaussdb_influx_instance": gaussdb.ResourceGaussDBInfluxInstanceV3(),
			"huaweicloud_gaussdb_mongo_instance":  gaussdb.ResourceGaussDBMongoInstanceV3(),

			"huaweicloud_ges_graph":    ges.ResourceGesGraph(),
			"huaweicloud_ges_metadata": ges.ResourceGesMetadata(),
			"huaweicloud_ges_backup":   ges.ResourceGesBackup(),

			"huaweicloud_hss_host_group": hss.ResourceHostGroup(),

			"huaweicloud_identity_access_key":            iam.ResourceIdentityKey(),
			"huaweicloud_identity_acl":                   iam.ResourceIdentityACL(),
			"huaweicloud_identity_agency":                iam.ResourceIAMAgencyV3(),
			"huaweicloud_identity_group":                 iam.ResourceIdentityGroup(),
			"huaweicloud_identity_group_membership":      iam.ResourceIdentityGroupMembership(),
			"huaweicloud_identity_group_role_assignment": iam.ResourceIdentityGroupRoleAssignment(),
			"huaweicloud_identity_project":               iam.ResourceIdentityProject(),
			"huaweicloud_identity_role":                  iam.ResourceIdentityRole(),
			"huaweicloud_identity_role_assignment":       iam.ResourceIdentityGroupRoleAssignment(),
			"huaweicloud_identity_user":                  iam.ResourceIdentityUser(),
			"huaweicloud_identity_user_role_assignment":  iam.ResourceIdentityUserRoleAssignment(),
			"huaweicloud_identity_provider":              iam.ResourceIdentityProvider(),
			"huaweicloud_identity_password_policy":       iam.ResourceIdentityPasswordPolicy(),
			"huaweicloud_identity_protection_policy":     iam.ResourceIdentityProtectionPolicy(),

			"huaweicloud_identitycenter_user":                     identitycenter.ResourceIdentityCenterUser(),
			"huaweicloud_identitycenter_group":                    identitycenter.ResourceIdentityCenterGroup(),
			"huaweicloud_identitycenter_group_membership":         identitycenter.ResourceGroupMembership(),
			"huaweicloud_identitycenter_permission_set":           identitycenter.ResourcePermissionSet(),
			"huaweicloud_identitycenter_system_policy_attachment": identitycenter.ResourceSystemPolicyAttachment(),
			"huaweicloud_identitycenter_account_assignment":       identitycenter.ResourceIdentityCenterAccountAssignment(),
			"huaweicloud_identitycenter_custom_policy_attachment": identitycenter.ResourceCustomPolicyAttachment(),

			"huaweicloud_iec_eip":                 resourceIecNetworkEip(),
			"huaweicloud_iec_keypair":             resourceIecKeypair(),
			"huaweicloud_iec_network_acl":         resourceIecNetworkACL(),
			"huaweicloud_iec_network_acl_rule":    resourceIecNetworkACLRule(),
			"huaweicloud_iec_security_group":      resourceIecSecurityGroup(),
			"huaweicloud_iec_security_group_rule": resourceIecSecurityGroupRule(),
			"huaweicloud_iec_server":              resourceIecServer(),
			"huaweicloud_iec_vip":                 resourceIecVipV1(),
			"huaweicloud_iec_vpc":                 ResourceIecVpc(),
			"huaweicloud_iec_vpc_subnet":          resourceIecSubnet(),

			"huaweicloud_images_image":                ims.ResourceImsImage(),
			"huaweicloud_images_image_copy":           ims.ResourceImsImageCopy(),
			"huaweicloud_images_image_share":          ims.ResourceImsImageShare(),
			"huaweicloud_images_image_share_accepter": ims.ResourceImsImageShareAccepter(),

			"huaweicloud_iotda_space":               iotda.ResourceSpace(),
			"huaweicloud_iotda_product":             iotda.ResourceProduct(),
			"huaweicloud_iotda_device":              iotda.ResourceDevice(),
			"huaweicloud_iotda_device_group":        iotda.ResourceDeviceGroup(),
			"huaweicloud_iotda_dataforwarding_rule": iotda.ResourceDataForwardingRule(),
			"huaweicloud_iotda_amqp":                iotda.ResourceAmqp(),
			"huaweicloud_iotda_device_certificate":  iotda.ResourceDeviceCertificate(),
			"huaweicloud_iotda_device_linkage_rule": iotda.ResourceDeviceLinkageRule(),

			"huaweicloud_kms_key":     dew.ResourceKmsKey(),
			"huaweicloud_kps_keypair": dew.ResourceKeypair(),
			"huaweicloud_kms_grant":   dew.ResourceKmsGrant(),

			"huaweicloud_lb_certificate":  lb.ResourceCertificateV2(),
			"huaweicloud_lb_l7policy":     lb.ResourceL7PolicyV2(),
			"huaweicloud_lb_l7rule":       lb.ResourceL7RuleV2(),
			"huaweicloud_lb_loadbalancer": lb.ResourceLoadBalancer(),
			"huaweicloud_lb_listener":     lb.ResourceListener(),
			"huaweicloud_lb_member":       lb.ResourceMemberV2(),
			"huaweicloud_lb_monitor":      lb.ResourceMonitorV2(),
			"huaweicloud_lb_pool":         lb.ResourcePoolV2(),
			"huaweicloud_lb_whitelist":    lb.ResourceWhitelistV2(),

			"huaweicloud_live_domain":               live.ResourceDomain(),
			"huaweicloud_live_recording":            live.ResourceRecording(),
			"huaweicloud_live_record_callback":      live.ResourceRecordCallback(),
			"huaweicloud_live_transcoding":          live.ResourceTranscoding(),
			"huaweicloud_live_snapshot":             live.ResourceLiveSnapshot(),
			"huaweicloud_live_bucket_authorization": live.ResourceLiveBucketAuthorization(),

			"huaweicloud_lts_aom_access":                       lts.ResourceAOMAccess(),
			"huaweicloud_lts_group":                            lts.ResourceLTSGroup(),
			"huaweicloud_lts_host_group":                       lts.ResourceHostGroup(),
			"huaweicloud_lts_host_access":                      lts.ResourceHostAccessConfig(),
			"huaweicloud_lts_stream":                           lts.ResourceLTSStream(),
			"huaweicloud_lts_structuring_configuration":        lts.ResourceStructConfig(),
			"huaweicloud_lts_structuring_custom_configuration": lts.ResourceStructCustomConfig(),
			"huaweicloud_lts_transfer":                         lts.ResourceLtsTransfer(),
			"huaweicloud_lts_keywords_alarm_rule":              lts.ResourceKeywordsAlarmRule(),
			"huaweicloud_lts_sql_alarm_rule":                   lts.ResourceSQLAlarmRule(),
			"huaweicloud_lts_notification_template":            lts.ResourceNotificationTemplate(),
			"huaweicloud_lts_search_criteria":                  lts.ResourceSearchCriteria(),
			"huaweicloud_lts_waf_access":                       lts.ResourceWAFAccess(),

			"huaweicloud_mapreduce_cluster":         mrs.ResourceMRSClusterV2(),
			"huaweicloud_mapreduce_job":             mrs.ResourceMRSJobV2(),
			"huaweicloud_mapreduce_data_connection": mrs.ResourceDataConnection(),
			"huaweicloud_mapreduce_scaling_policy":  mrs.ResourceScalingPolicy(),

			"huaweicloud_meeting_admin_assignment": meeting.ResourceAdminAssignment(),
			"huaweicloud_meeting_conference":       meeting.ResourceConference(),
			"huaweicloud_meeting_user":             meeting.ResourceUser(),

			"huaweicloud_modelarts_dataset":                modelarts.ResourceDataset(),
			"huaweicloud_modelarts_dataset_version":        modelarts.ResourceDatasetVersion(),
			"huaweicloud_modelarts_notebook":               modelarts.ResourceNotebook(),
			"huaweicloud_modelarts_notebook_mount_storage": modelarts.ResourceNotebookMountStorage(),
			"huaweicloud_modelarts_model":                  modelarts.ResourceModelartsModel(),
			"huaweicloud_modelarts_service":                modelarts.ResourceModelartsService(),
			"huaweicloud_modelarts_workspace":              modelarts.ResourceModelartsWorkspace(),
			"huaweicloud_modelarts_authorization":          modelarts.ResourceModelArtsAuthorization(),
			"huaweicloud_modelarts_network":                modelarts.ResourceModelartsNetwork(),
			"huaweicloud_modelarts_resource_pool":          modelarts.ResourceModelartsResourcePool(),

			"huaweicloud_dataarts_studio_instance": dataarts.ResourceStudioInstance(),

			"huaweicloud_mpc_transcoding_template":       mpc.ResourceTranscodingTemplate(),
			"huaweicloud_mpc_transcoding_template_group": mpc.ResourceTranscodingTemplateGroup(),

			"huaweicloud_mrs_cluster": ResourceMRSClusterV1(),
			"huaweicloud_mrs_job":     ResourceMRSJobV1(),

			"huaweicloud_nat_dnat_rule": nat.ResourcePublicDnatRule(),
			"huaweicloud_nat_gateway":   nat.ResourcePublicGateway(),
			"huaweicloud_nat_snat_rule": nat.ResourcePublicSnatRule(),

			"huaweicloud_nat_private_dnat_rule":  nat.ResourcePrivateDnatRule(),
			"huaweicloud_nat_private_gateway":    nat.ResourcePrivateGateway(),
			"huaweicloud_nat_private_snat_rule":  nat.ResourcePrivateSnatRule(),
			"huaweicloud_nat_private_transit_ip": nat.ResourcePrivateTransitIp(),

			"huaweicloud_network_acl":              ResourceNetworkACL(),
			"huaweicloud_network_acl_rule":         ResourceNetworkACLRule(),
			"huaweicloud_networking_secgroup":      vpc.ResourceNetworkingSecGroup(),
			"huaweicloud_networking_secgroup_rule": vpc.ResourceNetworkingSecGroupRule(),
			"huaweicloud_networking_vip":           vpc.ResourceNetworkingVip(),
			"huaweicloud_networking_vip_associate": vpc.ResourceNetworkingVIPAssociateV2(),

			"huaweicloud_obs_bucket":             obs.ResourceObsBucket(),
			"huaweicloud_obs_bucket_acl":         obs.ResourceOBSBucketAcl(),
			"huaweicloud_obs_bucket_object":      obs.ResourceObsBucketObject(),
			"huaweicloud_obs_bucket_object_acl":  obs.ResourceOBSBucketObjectAcl(),
			"huaweicloud_obs_bucket_policy":      obs.ResourceObsBucketPolicy(),
			"huaweicloud_obs_bucket_replication": obs.ResourceObsBucketReplication(),

			"huaweicloud_oms_migration_task":       oms.ResourceMigrationTask(),
			"huaweicloud_oms_migration_task_group": oms.ResourceMigrationTaskGroup(),

			"huaweicloud_ram_resource_share": ram.ResourceRAMShare(),

			"huaweicloud_rds_mysql_account":                rds.ResourceMysqlAccount(),
			"huaweicloud_rds_mysql_database":               rds.ResourceMysqlDatabase(),
			"huaweicloud_rds_mysql_database_privilege":     rds.ResourceMysqlDatabasePrivilege(),
			"huaweicloud_rds_instance":                     rds.ResourceRdsInstance(),
			"huaweicloud_rds_parametergroup":               rds.ResourceRdsConfiguration(),
			"huaweicloud_rds_read_replica_instance":        rds.ResourceRdsReadReplicaInstance(),
			"huaweicloud_rds_backup":                       rds.ResourceBackup(),
			"huaweicloud_rds_cross_region_backup_strategy": rds.ResourceBackupStrategy(),
			"huaweicloud_rds_sql_audit":                    rds.ResourceSQLAudit(),

			"huaweicloud_rms_policy_assignment":                  rms.ResourcePolicyAssignment(),
			"huaweicloud_rms_resource_aggregator":                rms.ResourceAggregator(),
			"huaweicloud_rms_resource_aggregation_authorization": rms.ResourceAggregationAuthorization(),
			"huaweicloud_rms_resource_recorder":                  rms.ResourceRecorder(),
			"huaweicloud_rms_advanced_query":                     rms.ResourceAdvancedQuery(),
			"huaweicloud_rms_assignment_package":                 rms.ResourceAssignmentPackage(),

			"huaweicloud_sdrs_drill":              sdrs.ResourceDrill(),
			"huaweicloud_sdrs_replication_pair":   sdrs.ResourceReplicationPair(),
			"huaweicloud_sdrs_protection_group":   sdrs.ResourceProtectionGroup(),
			"huaweicloud_sdrs_protected_instance": sdrs.ResourceProtectedInstance(),
			"huaweicloud_sdrs_replication_attach": sdrs.ResourceReplicationAttach(),

			"huaweicloud_secmaster_incident": secmaster.ResourceIncident(),

			"huaweicloud_servicestage_application":                 servicestage.ResourceApplication(),
			"huaweicloud_servicestage_component_instance":          servicestage.ResourceComponentInstance(),
			"huaweicloud_servicestage_component":                   servicestage.ResourceComponent(),
			"huaweicloud_servicestage_environment":                 servicestage.ResourceEnvironment(),
			"huaweicloud_servicestage_repo_token_authorization":    servicestage.ResourceRepoTokenAuth(),
			"huaweicloud_servicestage_repo_password_authorization": servicestage.ResourceRepoPwdAuth(),

			"huaweicloud_sfs_access_rule": sfs.ResourceSFSAccessRuleV2(),
			"huaweicloud_sfs_file_system": sfs.ResourceSFSFileSystemV2(),
			"huaweicloud_sfs_turbo":       sfs.ResourceSFSTurbo(),

			"huaweicloud_smn_topic":            smn.ResourceTopic(),
			"huaweicloud_smn_subscription":     smn.ResourceSubscription(),
			"huaweicloud_smn_message_template": smn.ResourceSmnMessageTemplate(),
			"huaweicloud_smn_logtank":          smn.ResourceSmnLogtank(),

			"huaweicloud_sms_server_template": sms.ResourceServerTemplate(),
			"huaweicloud_sms_task":            sms.ResourceMigrateTask(),

			"huaweicloud_swr_organization":             swr.ResourceSWROrganization(),
			"huaweicloud_swr_organization_permissions": swr.ResourceSWROrganizationPermissions(),
			"huaweicloud_swr_repository":               swr.ResourceSWRRepository(),
			"huaweicloud_swr_repository_sharing":       swr.ResourceSWRRepositorySharing(),
			"huaweicloud_swr_image_permissions":        swr.ResourceSwrImagePermissions(),
			"huaweicloud_swr_image_trigger":            swr.ResourceSwrImageTrigger(),
			"huaweicloud_swr_image_retention_policy":   swr.ResourceSwrImageRetentionPolicy(),
			"huaweicloud_swr_image_auto_sync":          swr.ResourceSwrImageAutoSync(),

			"huaweicloud_tms_resource_tags": tms.ResourceResourceTags(),
			"huaweicloud_tms_tags":          tms.ResourceTmsTag(),

			"huaweicloud_ucs_fleet":   ucs.ResourceFleet(),
			"huaweicloud_ucs_cluster": ucs.ResourceCluster(),
			"huaweicloud_ucs_policy":  ucs.ResourcePolicy(),

			"huaweicloud_vod_media_asset":                vod.ResourceMediaAsset(),
			"huaweicloud_vod_media_category":             vod.ResourceMediaCategory(),
			"huaweicloud_vod_transcoding_template_group": vod.ResourceTranscodingTemplateGroup(),
			"huaweicloud_vod_watermark_template":         vod.ResourceWatermarkTemplate(),

			"huaweicloud_vpc_bandwidth":           eip.ResourceVpcBandWidthV2(),
			"huaweicloud_vpc_bandwidth_associate": eip.ResourceBandWidthAssociate(),
			"huaweicloud_vpc_eip":                 eip.ResourceVpcEIPV1(),
			"huaweicloud_vpc_eip_associate":       eip.ResourceEIPAssociate(),

			"huaweicloud_vpc_peering_connection":          vpc.ResourceVpcPeeringConnectionV2(),
			"huaweicloud_vpc_peering_connection_accepter": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"huaweicloud_vpc_route_table":                 vpc.ResourceVPCRouteTable(),
			"huaweicloud_vpc_route":                       vpc.ResourceVPCRouteTableRoute(),
			"huaweicloud_vpc":                             vpc.ResourceVirtualPrivateCloudV1(),
			"huaweicloud_vpc_subnet":                      vpc.ResourceVpcSubnetV1(),
			"huaweicloud_vpc_address_group":               vpc.ResourceVpcAddressGroup(),
			"huaweicloud_vpc_flow_log":                    vpc.ResourceVpcFlowLog(),

			"huaweicloud_vpcep_approval": vpcep.ResourceVPCEndpointApproval(),
			"huaweicloud_vpcep_endpoint": vpcep.ResourceVPCEndpoint(),
			"huaweicloud_vpcep_service":  vpcep.ResourceVPCEndpointService(),

			"huaweicloud_vpn_gateway":                 vpn.ResourceGateway(),
			"huaweicloud_vpn_customer_gateway":        vpn.ResourceCustomerGateway(),
			"huaweicloud_vpn_connection":              vpn.ResourceConnection(),
			"huaweicloud_vpn_connection_health_check": vpn.ResourceConnectionHealthCheck(),

			"huaweicloud_scm_certificate": scm.ResourceScmCertificate(),

			"huaweicloud_waf_address_group":                       waf.ResourceWafAddressGroup(),
			"huaweicloud_waf_certificate":                         waf.ResourceWafCertificateV1(),
			"huaweicloud_waf_cloud_instance":                      waf.ResourceCloudInstance(),
			"huaweicloud_waf_domain":                              waf.ResourceWafDomainV1(),
			"huaweicloud_waf_policy":                              waf.ResourceWafPolicyV1(),
			"huaweicloud_waf_rule_anti_crawler":                   waf.ResourceRuleAntiCrawler(),
			"huaweicloud_waf_rule_blacklist":                      waf.ResourceWafRuleBlackListV1(),
			"huaweicloud_waf_rule_cc_protection":                  waf.ResourceRuleCCProtection(),
			"huaweicloud_waf_rule_data_masking":                   waf.ResourceWafRuleDataMaskingV1(),
			"huaweicloud_waf_rule_global_protection_whitelist":    waf.ResourceRuleGlobalProtectionWhitelist(),
			"huaweicloud_waf_rule_known_attack_source":            waf.ResourceRuleKnownAttack(),
			"huaweicloud_waf_rule_precise_protection":             waf.ResourceRulePreciseProtection(),
			"huaweicloud_waf_rule_web_tamper_protection":          waf.ResourceWafRuleWebTamperProtectionV1(),
			"huaweicloud_waf_rule_geolocation_access_control":     waf.ResourceRuleGeolocation(),
			"huaweicloud_waf_rule_information_leakage_prevention": waf.ResourceRuleLeakagePrevention(),
			"huaweicloud_waf_dedicated_instance":                  waf.ResourceWafDedicatedInstance(),
			"huaweicloud_waf_dedicated_domain":                    waf.ResourceWafDedicatedDomainV1(),
			"huaweicloud_waf_instance_group":                      waf.ResourceWafInstanceGroup(),
			"huaweicloud_waf_instance_group_associate":            waf.ResourceWafInstGroupAssociate(),
			"huaweicloud_waf_reference_table":                     waf.ResourceWafReferenceTableV1(),

			"huaweicloud_workspace_desktop": workspace.ResourceDesktop(),
			"huaweicloud_workspace_service": workspace.ResourceService(),
			"huaweicloud_workspace_user":    workspace.ResourceUser(),

			"huaweicloud_cpts_project": cpts.ResourceProject(),
			"huaweicloud_cpts_task":    cpts.ResourceTask(),

			// CodeArts
			"huaweicloud_codearts_project":            codearts.ResourceProject(),
			"huaweicloud_codearts_repository":         codearts.ResourceRepository(),
			"huaweicloud_codearts_deploy_application": codearts.ResourceDeployApplication(),
			"huaweicloud_codearts_deploy_group":       codearts.ResourceDeployGroup(),
			"huaweicloud_codearts_deploy_host":        codearts.ResourceDeployHost(),

			"huaweicloud_dsc_instance":  dsc.ResourceDscInstance(),
			"huaweicloud_dsc_asset_obs": dsc.ResourceAssetObs(),

			// internal only
			"huaweicloud_apm_aksk":                apm.ResourceApmAkSk(),
			"huaweicloud_aom_alarm_policy":        aom.ResourceAlarmPolicy(),
			"huaweicloud_aom_prometheus_instance": aom.ResourcePrometheusInstance(),

			"huaweicloud_aom_application":                 cmdb.ResourceAomApplication(),
			"huaweicloud_aom_component":                   cmdb.ResourceAomComponent(),
			"huaweicloud_aom_cmdb_resource_relationships": cmdb.ResourceCiRelationships(),
			"huaweicloud_aom_environment":                 cmdb.ResourceAomEnvironment(),

			"huaweicloud_lts_access_rule":     lts.ResourceAomMappingRule(),
			"huaweicloud_lts_dashboard":       lts.ResourceLtsDashboard(),
			"huaweicloud_elb_log":             lts.ResourceLtsElb(),
			"huaweicloud_lts_struct_template": lts.ResourceLtsStruct(),

			// Legacy
			"huaweicloud_networking_eip_associate": eip.ResourceEIPAssociate(),

			"huaweicloud_projectman_project": codearts.ResourceProject(),
			"huaweicloud_codehub_repository": codearts.ResourceRepository(),

			"huaweicloud_compute_instance_v2":             ecs.ResourceComputeInstance(),
			"huaweicloud_compute_interface_attach_v2":     ecs.ResourceComputeInterfaceAttach(),
			"huaweicloud_compute_keypair_v2":              ResourceComputeKeypairV2(),
			"huaweicloud_compute_servergroup_v2":          ecs.ResourceComputeServerGroup(),
			"huaweicloud_compute_volume_attach_v2":        ecs.ResourceComputeVolumeAttach(),
			"huaweicloud_compute_floatingip_associate_v2": ecs.ResourceComputeEIPAssociate(),

			"huaweicloud_dns_ptrrecord_v2": dns.ResourceDNSPtrRecord(),
			"huaweicloud_dns_recordset_v2": dns.ResourceDNSRecordSetV2(),
			"huaweicloud_dns_zone_v2":      dns.ResourceDNSZone(),

			"huaweicloud_dcs_instance_v1": dcs.ResourceDcsInstance(),
			"huaweicloud_dds_instance_v3": dds.ResourceDdsInstanceV3(),

			"huaweicloud_kms_key_v1": dew.ResourceKmsKey(),

			"huaweicloud_lb_certificate_v2":  lb.ResourceCertificateV2(),
			"huaweicloud_lb_loadbalancer_v2": lb.ResourceLoadBalancer(),
			"huaweicloud_lb_listener_v2":     lb.ResourceListener(),
			"huaweicloud_lb_pool_v2":         lb.ResourcePoolV2(),
			"huaweicloud_lb_member_v2":       lb.ResourceMemberV2(),
			"huaweicloud_lb_monitor_v2":      lb.ResourceMonitorV2(),
			"huaweicloud_lb_l7policy_v2":     lb.ResourceL7PolicyV2(),
			"huaweicloud_lb_l7rule_v2":       lb.ResourceL7RuleV2(),
			"huaweicloud_lb_whitelist_v2":    lb.ResourceWhitelistV2(),

			"huaweicloud_mrs_cluster_v1": ResourceMRSClusterV1(),
			"huaweicloud_mrs_job_v1":     ResourceMRSJobV1(),

			"huaweicloud_networking_secgroup_v2":      vpc.ResourceNetworkingSecGroup(),
			"huaweicloud_networking_secgroup_rule_v2": vpc.ResourceNetworkingSecGroupRule(),

			"huaweicloud_smn_topic_v2":        smn.ResourceTopic(),
			"huaweicloud_smn_subscription_v2": smn.ResourceSubscription(),

			"huaweicloud_rds_account":            rds.ResourceMysqlAccount(),
			"huaweicloud_rds_database":           rds.ResourceMysqlDatabase(),
			"huaweicloud_rds_database_privilege": rds.ResourceMysqlDatabasePrivilege(),
			"huaweicloud_rds_instance_v3":        rds.ResourceRdsInstance(),
			"huaweicloud_rds_parametergroup_v3":  rds.ResourceRdsConfiguration(),

			"huaweicloud_rf_stack": rfs.ResourceStack(),

			"huaweicloud_nat_dnat_rule_v2": nat.ResourcePublicDnatRule(),
			"huaweicloud_nat_gateway_v2":   nat.ResourcePublicGateway(),
			"huaweicloud_nat_snat_rule_v2": nat.ResourcePublicSnatRule(),

			"huaweicloud_sfs_access_rule_v2": sfs.ResourceSFSAccessRuleV2(),
			"huaweicloud_sfs_file_system_v2": sfs.ResourceSFSFileSystemV2(),

			"huaweicloud_iam_agency":    iam.ResourceIAMAgencyV3(),
			"huaweicloud_iam_agency_v3": iam.ResourceIAMAgencyV3(),

			"huaweicloud_vpc_bandwidth_v2":                   eip.ResourceVpcBandWidthV2(),
			"huaweicloud_vpc_eip_v1":                         eip.ResourceVpcEIPV1(),
			"huaweicloud_vpc_peering_connection_v2":          vpc.ResourceVpcPeeringConnectionV2(),
			"huaweicloud_vpc_peering_connection_accepter_v2": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"huaweicloud_vpc_v1":                             vpc.ResourceVirtualPrivateCloudV1(),
			"huaweicloud_vpc_subnet_v1":                      vpc.ResourceVpcSubnetV1(),

			"huaweicloud_cce_cluster_v3": cce.ResourceCCEClusterV3(),
			"huaweicloud_cce_node_v3":    cce.ResourceNode(),

			"huaweicloud_as_configuration_v1": as.ResourceASConfiguration(),
			"huaweicloud_as_group_v1":         as.ResourceASGroup(),
			"huaweicloud_as_policy_v1":        as.ResourceASPolicy(),

			"huaweicloud_identity_project_v3":          iam.ResourceIdentityProject(),
			"huaweicloud_identity_role_assignment_v3":  iam.ResourceIdentityGroupRoleAssignment(),
			"huaweicloud_identity_user_v3":             iam.ResourceIdentityUser(),
			"huaweicloud_identity_group_v3":            iam.ResourceIdentityGroup(),
			"huaweicloud_identity_group_membership_v3": iam.ResourceIdentityGroupMembership(),
			"huaweicloud_identity_provider_conversion": iam.ResourceIAMProviderConversion(),

			"huaweicloud_cdm_cluster_v1": cdm.ResourceCdmCluster(),
			"huaweicloud_css_cluster_v1": css.ResourceCssCluster(),
			"huaweicloud_dis_stream_v2":  dis.ResourceDisStream(),

			"huaweicloud_organizations_organization":            organizations.ResourceOrganization(),
			"huaweicloud_organizations_organizational_unit":     organizations.ResourceOrganizationalUnit(),
			"huaweicloud_organizations_account":                 organizations.ResourceAccount(),
			"huaweicloud_organizations_account_associate":       organizations.ResourceAccountAssociate(),
			"huaweicloud_organizations_account_invite":          organizations.ResourceAccountInvite(),
			"huaweicloud_organizations_account_invite_accepter": organizations.ResourceAccountInviteAccepter(),
			"huaweicloud_organizations_trusted_service":         organizations.ResourceTrustedService(),
			"huaweicloud_organizations_policy":                  organizations.ResourcePolicy(),
			"huaweicloud_organizations_policy_attach":           organizations.ResourcePolicyAttach(),

			"huaweicloud_dli_queue_v1":                dli.ResourceDliQueue(),
			"huaweicloud_networking_vip_v2":           vpc.ResourceNetworkingVip(),
			"huaweicloud_networking_vip_associate_v2": vpc.ResourceNetworkingVIPAssociateV2(),
			"huaweicloud_fgs_function_v2":             fgs.ResourceFgsFunctionV2(),
			"huaweicloud_cdn_domain_v1":               cdn.ResourceCdnDomainV1(),

			// Deprecated
			"huaweicloud_apig_vpc_channel":               deprecated.ResourceApigVpcChannelV2(),
			"huaweicloud_blockstorage_volume_v2":         deprecated.ResourceBlockStorageVolumeV2(),
			"huaweicloud_csbs_backup":                    deprecated.ResourceCSBSBackupV1(),
			"huaweicloud_csbs_backup_policy":             deprecated.ResourceCSBSBackupPolicyV1(),
			"huaweicloud_csbs_backup_policy_v1":          deprecated.ResourceCSBSBackupPolicyV1(),
			"huaweicloud_csbs_backup_v1":                 deprecated.ResourceCSBSBackupV1(),
			"huaweicloud_networking_network_v2":          deprecated.ResourceNetworkingNetworkV2(),
			"huaweicloud_networking_subnet_v2":           deprecated.ResourceNetworkingSubnetV2(),
			"huaweicloud_networking_floatingip_v2":       deprecated.ResourceNetworkingFloatingIPV2(),
			"huaweicloud_networking_router_v2":           deprecated.ResourceNetworkingRouterV2(),
			"huaweicloud_networking_router_interface_v2": deprecated.ResourceNetworkingRouterInterfaceV2(),
			"huaweicloud_networking_router_route_v2":     deprecated.ResourceNetworkingRouterRouteV2(),
			"huaweicloud_networking_port":                deprecated.ResourceNetworkingPortV2(),
			"huaweicloud_networking_port_v2":             deprecated.ResourceNetworkingPortV2(),
			"huaweicloud_vpc_route_v2":                   deprecated.ResourceVPCRouteV2(),
			"huaweicloud_ecs_instance_v1":                deprecated.ResourceEcsInstanceV1(),
			"huaweicloud_compute_secgroup_v2":            deprecated.ResourceComputeSecGroupV2(),
			"huaweicloud_compute_floatingip_v2":          deprecated.ResourceComputeFloatingIPV2(),
			"huaweicloud_oms_task":                       deprecated.ResourceMaasTaskV1(),

			"huaweicloud_fw_firewall_group_v2": deprecated.ResourceFWFirewallGroupV2(),
			"huaweicloud_fw_policy_v2":         deprecated.ResourceFWPolicyV2(),
			"huaweicloud_fw_rule_v2":           deprecated.ResourceFWRuleV2(),

			"huaweicloud_images_image_v2": deprecated.ResourceImagesImageV2(),

			"huaweicloud_dms_instance":    deprecated.ResourceDmsInstancesV1(),
			"huaweicloud_dms_instance_v1": deprecated.ResourceDmsInstancesV1(),
			"huaweicloud_dms_group":       deprecated.ResourceDmsGroups(),
			"huaweicloud_dms_group_v1":    deprecated.ResourceDmsGroups(),
			"huaweicloud_dms_queue":       deprecated.ResourceDmsQueues(),
			"huaweicloud_dms_queue_v1":    deprecated.ResourceDmsQueues(),

			"huaweicloud_cs_cluster":            deprecated.ResourceCsClusterV1(),
			"huaweicloud_cs_cluster_v1":         deprecated.ResourceCsClusterV1(),
			"huaweicloud_cs_route":              deprecated.ResourceCsRouteV1(),
			"huaweicloud_cs_route_v1":           deprecated.ResourceCsRouteV1(),
			"huaweicloud_cs_peering_connect":    deprecated.ResourceCsPeeringConnectV1(),
			"huaweicloud_cs_peering_connect_v1": deprecated.ResourceCsPeeringConnectV1(),

			"huaweicloud_vbs_backup":           deprecated.ResourceVBSBackupV2(),
			"huaweicloud_vbs_backup_policy":    deprecated.ResourceVBSBackupPolicyV2(),
			"huaweicloud_vbs_backup_policy_v2": deprecated.ResourceVBSBackupPolicyV2(),
			"huaweicloud_vbs_backup_v2":        deprecated.ResourceVBSBackupV2(),

			"huaweicloud_vpnaas_ipsec_policy_v2":    deprecated.ResourceVpnIPSecPolicyV2(),
			"huaweicloud_vpnaas_service_v2":         deprecated.ResourceVpnServiceV2(),
			"huaweicloud_vpnaas_ike_policy_v2":      deprecated.ResourceVpnIKEPolicyV2(),
			"huaweicloud_vpnaas_endpoint_group_v2":  deprecated.ResourceVpnEndpointGroupV2(),
			"huaweicloud_vpnaas_site_connection_v2": deprecated.ResourceVpnSiteConnectionV2(),

			"huaweicloud_vpnaas_endpoint_group":  deprecated.ResourceVpnEndpointGroupV2(),
			"huaweicloud_vpnaas_ike_policy":      deprecated.ResourceVpnIKEPolicyV2(),
			"huaweicloud_vpnaas_ipsec_policy":    deprecated.ResourceVpnIPSecPolicyV2(),
			"huaweicloud_vpnaas_service":         deprecated.ResourceVpnServiceV2(),
			"huaweicloud_vpnaas_site_connection": deprecated.ResourceVpnSiteConnectionV2(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11 cc
			terraformVersion = "0.11+compatible"
		}

		return configureProvider(ctx, d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"auth_url": "The Identity authentication URL.",

		"region": "The HuaweiCloud region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"project_id": "The ID of the project to login with.",

		"project_name": "The name of the project to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"domain_id": "The ID of the Domain to scope to.",

		"domain_name": "The name of the Domain to scope to.",

		"access_key":     "The access key of the HuaweiCloud to use.",
		"secret_key":     "The secret key of the HuaweiCloud to use.",
		"security_token": "The security token to authenticate with a temporary security credential.",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"agency_name": "The name of agency",

		"agency_domain_name": "The name of domain who created the agency (Identity v3).",

		"delegated_project": "The name of delegated project (Identity v3).",

		"assume_role_agency_name": "The name of agency for assume role.",

		"assume_role_domain_name": "The name of domain for assume role.",

		"cloud": "The endpoint of cloud provider, defaults to myhuaweicloud.com",

		"endpoints": "The custom endpoints used to override the default endpoint URL.",

		"regional": "Whether the service endpoints are regional",

		"shared_config_file": "The path to the shared config file. If not set, the default is ~/.hcloud/config.json.",

		"profile": "The profile name as set in the shared config file.",

		"max_retries": "How many times HTTP connection should be retried until giving up.",

		"enterprise_project_id": "enterprise project id",
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData, terraformVersion string) (interface{},
	diag.Diagnostics) {
	var tenantName, tenantID, delegatedProject, identityEndpoint string
	region := d.Get("region").(string)
	cloud := getCloudDomain(d.Get("cloud").(string), region)

	isRegional := d.Get("regional").(bool)
	if strings.HasPrefix(region, prefixEuropeRegion) {
		// the default format of endpoints in Europe site is xxx.{{region}}.{{cloud}}
		isRegional = true
	}

	// project_name is prior to tenant_name
	// if neither of them was set, use region as the default project
	if v, ok := d.GetOk("project_name"); ok && v.(string) != "" {
		tenantName = v.(string)
	} else if v, ok := d.GetOk("tenant_name"); ok && v.(string) != "" {
		tenantName = v.(string)
	} else {
		tenantName = region
	}

	// project_id is prior to tenant_id
	if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
		tenantID = v.(string)
	} else {
		tenantID = d.Get("tenant_id").(string)
	}

	// Use region as delegated_project if it's not set
	if v, ok := d.GetOk("delegated_project"); ok && v.(string) != "" {
		delegatedProject = v.(string)
	} else {
		delegatedProject = region
	}

	// use auth_url as identityEndpoint if specified
	if v, ok := d.GetOk("auth_url"); ok {
		identityEndpoint = v.(string)
	} else {
		// use cloud as basis for identityEndpoint
		identityEndpoint = fmt.Sprintf("https://iam.%s.%s/v3", region, cloud)
	}

	config := config.Config{
		AccessKey:           d.Get("access_key").(string),
		SecretKey:           d.Get("secret_key").(string),
		CACertFile:          d.Get("cacert_file").(string),
		ClientCertFile:      d.Get("cert").(string),
		ClientKeyFile:       d.Get("key").(string),
		DomainID:            d.Get("domain_id").(string),
		DomainName:          d.Get("domain_name").(string),
		IdentityEndpoint:    identityEndpoint,
		Insecure:            d.Get("insecure").(bool),
		Password:            d.Get("password").(string),
		Token:               d.Get("token").(string),
		SecurityToken:       d.Get("security_token").(string),
		Region:              region,
		TenantID:            tenantID,
		TenantName:          tenantName,
		Username:            d.Get("user_name").(string),
		UserID:              d.Get("user_id").(string),
		AgencyName:          d.Get("agency_name").(string),
		AgencyDomainName:    d.Get("agency_domain_name").(string),
		DelegatedProject:    delegatedProject,
		Cloud:               cloud,
		RegionClient:        isRegional,
		MaxRetries:          d.Get("max_retries").(int),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		SharedConfigFile:    d.Get("shared_config_file").(string),
		Profile:             d.Get("profile").(string),
		TerraformVersion:    terraformVersion,
		RegionProjectIDMap:  make(map[string]string),
		RPLock:              new(sync.Mutex),
		SecurityKeyLock:     new(sync.Mutex),
	}

	// get assume role
	assumeRoleList := d.Get("assume_role").([]interface{})
	if len(assumeRoleList) == 0 {
		// without assume_role block in provider
		delegatedAgencyName := os.Getenv("HW_ASSUME_ROLE_AGENCY_NAME")
		delegatedDomianName := os.Getenv("HW_ASSUME_ROLE_DOMAIN_NAME")
		if delegatedAgencyName != "" && delegatedDomianName != "" {
			config.AssumeRoleAgency = delegatedAgencyName
			config.AssumeRoleDomain = delegatedDomianName
		}
	} else {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		config.AssumeRoleAgency = assumeRole["agency_name"].(string)
		config.AssumeRoleDomain = assumeRole["domain_name"].(string)
	}

	// get custom endpoints
	endpoints, err := flattenProviderEndpoints(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	config.Endpoints = endpoints

	if err := config.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &config, nil
}

func flattenProviderEndpoints(d *schema.ResourceData) (map[string]string, error) {
	endpoints := d.Get("endpoints").(map[string]interface{})
	epMap := make(map[string]string)

	for key, val := range endpoints {
		endpoint := strings.TrimSpace(val.(string))
		// check empty string
		if endpoint == "" {
			return nil, fmt.Errorf("the value of customer endpoint %s must be specified", key)
		}

		// add prefix "https://" and suffix "/"
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}
		if !strings.HasSuffix(endpoint, "/") {
			endpoint = fmt.Sprintf("%s/", endpoint)
		}
		epMap[key] = endpoint
	}

	// unify the endpoint which has multiple versions
	for key := range endpoints {
		ep, ok := epMap[key]
		if !ok {
			continue
		}

		multiKeys := config.GetServiceDerivedCatalogKeys(key)
		for _, k := range multiKeys {
			epMap[k] = ep
		}
	}

	log.Printf("[DEBUG] customer endpoints: %+v", epMap)
	return epMap, nil
}

func getCloudDomain(cloud, region string) string {
	// first, use the specified value
	if cloud != "" {
		return cloud
	}

	// then check whether the region(eu-west-1xx) is located in Europe
	if strings.HasPrefix(region, prefixEuropeRegion) {
		return defaultEuropeCloud
	}
	return defaultCloud
}
