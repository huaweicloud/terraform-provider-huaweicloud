package huaweicloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/mutexkv"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/as"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/bms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ces"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cloudtable"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cpts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dis"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/drs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ecs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eps"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/gaussdb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/scm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/swr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/tms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
)

const defaultCloud string = "myhuaweicloud.com"

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

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
				DefaultFunc: schema.EnvDefaultFunc("HW_CLOUD", defaultCloud),
			},

			"endpoints": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["endpoints"],
				Elem:        &schema.Schema{Type: schema.TypeString},
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
			"huaweicloud_antiddos":           dataSourceAntiDdosV1(),
			"huaweicloud_availability_zones": DataSourceAvailabilityZones(),

			"huaweicloud_bms_flavors": bms.DataSourceBmsFlavors(),
			"huaweicloud_cbr_vaults":  cbr.DataSourceCbrVaultsV3(),

			"huaweicloud_cce_addon_template": DataSourceCCEAddonTemplateV3(),
			"huaweicloud_cce_cluster":        DataSourceCCEClusterV3(),
			"huaweicloud_cce_clusters":       cce.DataSourceCCEClusters(),
			"huaweicloud_cce_node":           DataSourceCCENodeV3(),
			"huaweicloud_cce_nodes":          cce.DataSourceCCENodes(),
			"huaweicloud_cce_node_pool":      DataSourceCCENodePoolV3(),
			"huaweicloud_cci_namespaces":     cci.DataSourceCciNamespaces(),

			"huaweicloud_cdm_flavors": DataSourceCdmFlavorV1(),

			"huaweicloud_compute_flavors":   DataSourceEcsFlavors(),
			"huaweicloud_compute_instance":  DataSourceComputeInstance(),
			"huaweicloud_compute_instances": DataSourceComputeInstances(),

			"huaweicloud_csbs_backup":        dataSourceCSBSBackupV1(),
			"huaweicloud_csbs_backup_policy": dataSourceCSBSBackupPolicyV1(),

			"huaweicloud_csms_secret_version": dew.DataSourceDewCsmsSecret(),
			"huaweicloud_css_flavors":         css.DataSourceCssFlavors(),

			"huaweicloud_dcs_az":             deprecated.DataSourceDcsAZV1(),
			"huaweicloud_dcs_maintainwindow": dcs.DataSourceDcsMaintainWindow(),
			"huaweicloud_dcs_product":        deprecated.DataSourceDcsProductV1(),

			"huaweicloud_dds_flavors":        dds.DataSourceDDSFlavorV3(),
			"huaweicloud_dms_az":             deprecated.DataSourceDmsAZ(),
			"huaweicloud_dms_product":        dms.DataSourceDmsProduct(),
			"huaweicloud_dms_maintainwindow": dms.DataSourceDmsMaintainWindow(),

			"huaweicloud_elb_flavors":        dataSourceElbFlavorsV3(),
			"huaweicloud_enterprise_project": eps.DataSourceEnterpriseProject(),
			"huaweicloud_evs_volumes":        evs.DataSourceEvsVolumesV2(),
			"huaweicloud_fgs_dependencies":   fgs.DataSourceFunctionGraphDependencies(),

			"huaweicloud_gaussdb_cassandra_dedicated_resource": dataSourceGeminiDBDehResource(),
			"huaweicloud_gaussdb_cassandra_flavors":            gaussdb.DataSourceCassandraFlavors(),
			"huaweicloud_gaussdb_nosql_flavors":                gaussdb.DataSourceGaussDBNoSQLFlavors(),
			"huaweicloud_gaussdb_cassandra_instance":           dataSourceGeminiDBInstance(),
			"huaweicloud_gaussdb_cassandra_instances":          dataSourceGeminiDBInstances(),
			"huaweicloud_gaussdb_opengauss_instance":           dataSourceOpenGaussInstance(),
			"huaweicloud_gaussdb_opengauss_instances":          gaussdb.DataSourceOpenGaussInstances(),
			"huaweicloud_gaussdb_mysql_configuration":          dataSourceGaussdbMysqlConfigurations(),
			"huaweicloud_gaussdb_mysql_dedicated_resource":     dataSourceGaussDBMysqlDehResource(),
			"huaweicloud_gaussdb_mysql_flavors":                dataSourceGaussdbMysqlFlavors(),
			"huaweicloud_gaussdb_mysql_instance":               dataSourceGaussDBMysqlInstance(),
			"huaweicloud_gaussdb_mysql_instances":              dataSourceGaussDBMysqlInstances(),
			"huaweicloud_gaussdb_redis_instance":               dataSourceGaussRedisInstance(),

			"huaweicloud_identity_role":        iam.DataSourceIdentityRoleV3(),
			"huaweicloud_identity_custom_role": iam.DataSourceIdentityCustomRole(),
			"huaweicloud_identity_group":       iam.DataSourceIdentityGroup(),

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

			"huaweicloud_images_image": ims.DataSourceImagesImageV2(),

			"huaweicloud_kms_key":      DataSourceKmsKeyV1(),
			"huaweicloud_kms_data_key": DataSourceKmsDataKeyV1(),

			"huaweicloud_lb_loadbalancer": lb.DataSourceELBV2Loadbalancer(),
			"huaweicloud_lb_certificate":  lb.DataSourceLBCertificateV2(),
			"huaweicloud_elb_certificate": elb.DataSourceELBCertificateV3(),

			"huaweicloud_nat_gateway":          DataSourceNatGatewayV2(),
			"huaweicloud_networking_port":      DataSourceNetworkingPortV2(),
			"huaweicloud_networking_secgroup":  DataSourceNetworkingSecGroup(),
			"huaweicloud_networking_secgroups": vpc.DataSourceNetworkingSecGroups(),

			"huaweicloud_modelarts_datasets":        modelarts.DataSourceDatasets(),
			"huaweicloud_modelarts_notebook_images": modelarts.DataSourceNotebookImages(),

			"huaweicloud_obs_buckets":       obs.DataSourceObsBuckets(),
			"huaweicloud_obs_bucket_object": DataSourceObsBucketObject(),

			"huaweicloud_rds_flavors":         rds.DataSourceRdsFlavor(),
			"huaweicloud_rds_engine_versions": rds.DataSourceRdsEngineVersionsV3(),
			"huaweicloud_rds_instances":       rds.DataSourceRdsInstances(),

			"huaweicloud_sfs_file_system":   DataSourceSFSFileSystemV2(),
			"huaweicloud_vbs_backup_policy": dataSourceVBSBackupPolicyV2(),
			"huaweicloud_vbs_backup":        dataSourceVBSBackupV2(),

			"huaweicloud_vpc_bandwidth":          eip.DataSourceBandWidth(),
			"huaweicloud_vpc_eip":                eip.DataSourceVpcEip(),
			"huaweicloud_vpc_eips":               eip.DataSourceVpcEips(),
			"huaweicloud_vpc":                    vpc.DataSourceVpcV1(),
			"huaweicloud_vpcs":                   vpc.DataSourceVpcs(),
			"huaweicloud_vpc_ids":                vpc.DataSourceVpcIdsV1(),
			"huaweicloud_vpc_peering_connection": vpc.DataSourceVpcPeeringConnectionV2(),
			"huaweicloud_vpc_route_table":        vpc.DataSourceVPCRouteTable(),
			"huaweicloud_vpc_subnet":             vpc.DataSourceVpcSubnetV1(),
			"huaweicloud_vpc_subnets":            vpc.DataSourceVpcSubnets(),
			"huaweicloud_vpc_subnet_ids":         vpc.DataSourceVpcSubnetIdsV1(),
			"huaweicloud_vpcep_public_services":  DataSourceVPCEPPublicServices(),

			"huaweicloud_waf_certificate":         waf.DataSourceWafCertificateV1(),
			"huaweicloud_waf_policies":            waf.DataSourceWafPoliciesV1(),
			"huaweicloud_waf_dedicated_instances": waf.DataSourceWafDedicatedInstancesV1(),
			"huaweicloud_waf_reference_tables":    waf.DataSourceWafReferenceTablesV1(),
			"huaweicloud_waf_instance_groups":     waf.DataSourceWafInstanceGroups(),
			"huaweicloud_dws_flavors":             dws.DataSourceDwsFlavlors(),

			// Legacy
			"huaweicloud_images_image_v2":        ims.DataSourceImagesImageV2(),
			"huaweicloud_networking_port_v2":     DataSourceNetworkingPortV2(),
			"huaweicloud_networking_secgroup_v2": DataSourceNetworkingSecGroup(),

			"huaweicloud_kms_key_v1":      DataSourceKmsKeyV1(),
			"huaweicloud_kms_data_key_v1": DataSourceKmsDataKeyV1(),

			"huaweicloud_rds_flavors_v3":     rds.DataSourceRdsFlavor(),
			"huaweicloud_sfs_file_system_v2": DataSourceSFSFileSystemV2(),

			"huaweicloud_vpc_v1":                    vpc.DataSourceVpcV1(),
			"huaweicloud_vpc_ids_v1":                vpc.DataSourceVpcIdsV1(),
			"huaweicloud_vpc_peering_connection_v2": vpc.DataSourceVpcPeeringConnectionV2(),
			"huaweicloud_vpc_subnet_v1":             vpc.DataSourceVpcSubnetV1(),
			"huaweicloud_vpc_subnet_ids_v1":         vpc.DataSourceVpcSubnetIdsV1(),

			"huaweicloud_cce_cluster_v3": DataSourceCCEClusterV3(),
			"huaweicloud_cce_node_v3":    DataSourceCCENodeV3(),

			"huaweicloud_csbs_backup_v1":        dataSourceCSBSBackupV1(),
			"huaweicloud_csbs_backup_policy_v1": dataSourceCSBSBackupPolicyV1(),

			"huaweicloud_dms_az_v1":             deprecated.DataSourceDmsAZ(),
			"huaweicloud_dms_product_v1":        dms.DataSourceDmsProduct(),
			"huaweicloud_dms_maintainwindow_v1": dms.DataSourceDmsMaintainWindow(),

			"huaweicloud_vbs_backup_policy_v2": dataSourceVBSBackupPolicyV2(),
			"huaweicloud_vbs_backup_v2":        dataSourceVBSBackupV2(),

			"huaweicloud_antiddos_v1": dataSourceAntiDdosV1(),

			"huaweicloud_dcs_az_v1":             deprecated.DataSourceDcsAZV1(),
			"huaweicloud_dcs_maintainwindow_v1": dcs.DataSourceDcsMaintainWindow(),
			"huaweicloud_dcs_product_v1":        deprecated.DataSourceDcsProductV1(),
			"huaweicloud_dcs_flavors":           dcs.DataSourceDcsFlavorsV2(),

			"huaweicloud_dds_flavors_v3":   dds.DataSourceDDSFlavorV3(),
			"huaweicloud_identity_role_v3": iam.DataSourceIdentityRoleV3(),
			"huaweicloud_cdm_flavors_v1":   DataSourceCdmFlavorV1(),

			// Deprecated ongoing (without DeprecationMessage), used by other providers
			"huaweicloud_vpc_route":        vpc.DataSourceVpcRouteV2(),
			"huaweicloud_vpc_route_ids":    vpc.DataSourceVpcRouteIdsV2(),
			"huaweicloud_vpc_route_v2":     vpc.DataSourceVpcRouteV2(),
			"huaweicloud_vpc_route_ids_v2": vpc.DataSourceVpcRouteIdsV2(),

			// Deprecated
			"huaweicloud_compute_availability_zones_v2": dataSourceComputeAvailabilityZonesV2(),
			"huaweicloud_networking_network_v2":         dataSourceNetworkingNetworkV2(),
			"huaweicloud_networking_subnet_v2":          dataSourceNetworkingSubnetV2(),
			"huaweicloud_cts_tracker":                   deprecated.DataSourceCTSTrackerV1(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"huaweicloud_api_gateway_api":        ResourceAPIGatewayAPI(),
			"huaweicloud_api_gateway_group":      ResourceAPIGatewayGroup(),
			"huaweicloud_apig_api":               apig.ResourceApigAPIV2(),
			"huaweicloud_apig_api_publishment":   apig.ResourceApigApiPublishment(),
			"huaweicloud_apig_instance":          apig.ResourceApigInstanceV2(),
			"huaweicloud_apig_application":       apig.ResourceApigApplicationV2(),
			"huaweicloud_apig_custom_authorizer": apig.ResourceApigCustomAuthorizerV2(),
			"huaweicloud_apig_environment":       apig.ResourceApigEnvironmentV2(),
			"huaweicloud_apig_group":             apig.ResourceApigGroupV2(),
			"huaweicloud_apig_response":          apig.ResourceApigResponseV2(),
			"huaweicloud_apig_throttling_policy": apig.ResourceApigThrottlingPolicyV2(),
			"huaweicloud_apig_vpc_channel":       apig.ResourceApigVpcChannelV2(),

			"huaweicloud_as_configuration":  as.ResourceASConfiguration(),
			"huaweicloud_as_group":          as.ResourceASGroup(),
			"huaweicloud_as_lifecycle_hook": as.ResourceASLifecycleHook(),
			"huaweicloud_as_policy":         as.ResourceASPolicy(),

			"huaweicloud_bms_instance": bms.ResourceBmsInstance(),
			"huaweicloud_bcs_instance": resourceBCSInstanceV2(),

			"huaweicloud_cbr_policy": cbr.ResourceCBRPolicyV3(),
			"huaweicloud_cbr_vault":  cbr.ResourceCBRVaultV3(),

			"huaweicloud_cce_cluster":     ResourceCCEClusterV3(),
			"huaweicloud_cce_node":        ResourceCCENodeV3(),
			"huaweicloud_cce_node_attach": ResourceCCENodeAttachV3(),
			"huaweicloud_cce_addon":       ResourceCCEAddonV3(),
			"huaweicloud_cce_node_pool":   ResourceCCENodePool(),
			"huaweicloud_cce_namespace":   cce.ResourceCCENamespaceV1(),
			"huaweicloud_cce_pvc":         cce.ResourceCcePersistentVolumeClaimsV1(),

			"huaweicloud_cts_tracker":      cts.ResourceCTSTracker(),
			"huaweicloud_cts_data_tracker": cts.ResourceCTSDataTracker(),
			"huaweicloud_cci_namespace":    cci.ResourceCciNamespace(),
			"huaweicloud_cci_network":      cci.ResourceCciNetworkV1(),
			"huaweicloud_cci_pvc":          ResourceCCIPersistentVolumeClaimV1(),

			"huaweicloud_cdm_cluster": cdm.ResourceCdmCluster(),
			"huaweicloud_cdm_job":     cdm.ResourceCdmJob(),
			"huaweicloud_cdm_link":    cdm.ResourceCdmLink(),

			"huaweicloud_cdn_domain":         resourceCdnDomainV1(),
			"huaweicloud_ces_alarmrule":      ces.ResourceAlarmRule(),
			"huaweicloud_cloudtable_cluster": cloudtable.ResourceCloudTableCluster(),

			"huaweicloud_compute_instance":         ResourceComputeInstanceV2(),
			"huaweicloud_compute_interface_attach": ResourceComputeInterfaceAttachV2(),
			"huaweicloud_compute_keypair":          ResourceComputeKeypairV2(),
			"huaweicloud_compute_servergroup":      ResourceComputeServerGroupV2(),
			"huaweicloud_compute_eip_associate":    ResourceComputeFloatingIPAssociateV2(),
			"huaweicloud_compute_volume_attach":    ecs.ResourceComputeVolumeAttach(),

			"huaweicloud_cs_cluster":         deprecated.ResourceCsClusterV1(),
			"huaweicloud_cs_peering_connect": deprecated.ResourceCsPeeringConnectV1(),
			"huaweicloud_cs_route":           deprecated.ResourceCsRouteV1(),

			"huaweicloud_csbs_backup":        resourceCSBSBackupV1(),
			"huaweicloud_csbs_backup_policy": resourceCSBSBackupPolicyV1(),

			"huaweicloud_csms_secret": dew.ResourceCsmsSecret(),

			"huaweicloud_css_cluster":   css.ResourceCssCluster(),
			"huaweicloud_css_snapshot":  css.ResourceCssSnapshot(),
			"huaweicloud_css_thesaurus": css.ResourceCssthesaurus(),

			"huaweicloud_dcs_instance": dcs.ResourceDcsInstance(),
			"huaweicloud_dds_instance": dds.ResourceDdsInstanceV3(),
			"huaweicloud_dis_stream":   dis.ResourceDisStream(),

			"huaweicloud_dli_database":     dli.ResourceDliSqlDatabaseV1(),
			"huaweicloud_dli_package":      dli.ResourceDliPackageV2(),
			"huaweicloud_dli_queue":        dli.ResourceDliQueue(),
			"huaweicloud_dli_spark_job":    dli.ResourceDliSparkJobV2(),
			"huaweicloud_dli_sql_job":      dli.ResourceSqlJob(),
			"huaweicloud_dli_table":        dli.ResourceDliTable(),
			"huaweicloud_dli_flinksql_job": dli.ResourceFlinkSqlJob(),
			"huaweicloud_dli_flinkjar_job": dli.ResourceFlinkJarJob(),
			"huaweicloud_dli_permission":   dli.ResourceDliPermission(),

			"huaweicloud_dms_instance":          deprecated.ResourceDmsInstancesV1(),
			"huaweicloud_dms_kafka_instance":    dms.ResourceDmsKafkaInstance(),
			"huaweicloud_dms_kafka_topic":       dms.ResourceDmsKafkaTopic(),
			"huaweicloud_dms_rabbitmq_instance": dms.ResourceDmsRabbitmqInstance(),

			"huaweicloud_dns_ptrrecord": ResourceDNSPtrRecordV2(),
			"huaweicloud_dns_recordset": ResourceDNSRecordSetV2(),
			"huaweicloud_dns_zone":      ResourceDNSZoneV2(),

			"huaweicloud_drs_job":     drs.ResourceDrsJob(),
			"huaweicloud_dws_cluster": dws.ResourceDwsCluster(),

			"huaweicloud_elb_certificate":  ResourceCertificateV3(),
			"huaweicloud_elb_l7policy":     ResourceL7PolicyV3(),
			"huaweicloud_elb_l7rule":       ResourceL7RuleV3(),
			"huaweicloud_elb_listener":     elb.ResourceListenerV3(),
			"huaweicloud_elb_loadbalancer": elb.ResourceLoadBalancerV3(),
			"huaweicloud_elb_monitor":      ResourceMonitorV3(),
			"huaweicloud_elb_ipgroup":      ResourceIpGroupV3(),
			"huaweicloud_elb_pool":         ResourcePoolV3(),
			"huaweicloud_elb_member":       ResourceMemberV3(),

			"huaweicloud_enterprise_project": eps.ResourceEnterpriseProject(),

			"huaweicloud_evs_snapshot": ResourceEvsSnapshotV2(),
			"huaweicloud_evs_volume":   evs.ResourceEvsVolume(),

			"huaweicloud_fgs_dependency": fgs.ResourceFgsDependency(),
			"huaweicloud_fgs_function":   fgs.ResourceFgsFunctionV2(),
			"huaweicloud_fgs_trigger":    fgs.ResourceFunctionGraphTrigger(),

			"huaweicloud_gaussdb_cassandra_instance": resourceGeminiDBInstanceV3(),
			"huaweicloud_gaussdb_mysql_instance":     resourceGaussDBInstance(),
			"huaweicloud_gaussdb_mysql_proxy":        gaussdb.ResourceGaussDBProxy(),
			"huaweicloud_gaussdb_opengauss_instance": resourceOpenGaussInstance(),
			"huaweicloud_gaussdb_redis_instance":     resourceGaussRedisInstanceV3(),

			"huaweicloud_ges_graph": ResourceGesGraphV1(),

			"huaweicloud_identity_access_key":       iam.ResourceIdentityKey(),
			"huaweicloud_identity_acl":              iam.ResourceIdentityACL(),
			"huaweicloud_identity_agency":           iam.ResourceIAMAgencyV3(),
			"huaweicloud_identity_group":            iam.ResourceIdentityGroupV3(),
			"huaweicloud_identity_group_membership": iam.ResourceIdentityGroupMembershipV3(),
			"huaweicloud_identity_project":          iam.ResourceIdentityProjectV3(),
			"huaweicloud_identity_role":             iam.ResourceIdentityRole(),
			"huaweicloud_identity_role_assignment":  iam.ResourceIdentityRoleAssignmentV3(),
			"huaweicloud_identity_user":             iam.ResourceIdentityUserV3(),
			"huaweicloud_identity_provider":         iam.ResourceIdentityProvider(),

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

			"huaweicloud_images_image": ResourceImsImage(),

			"huaweicloud_kms_key":     ResourceKmsKeyV1(),
			"huaweicloud_kps_keypair": dew.ResourceKeypair(),

			"huaweicloud_lb_certificate":  lb.ResourceCertificateV2(),
			"huaweicloud_lb_l7policy":     lb.ResourceL7PolicyV2(),
			"huaweicloud_lb_l7rule":       lb.ResourceL7RuleV2(),
			"huaweicloud_lb_listener":     lb.ResourceListenerV2(),
			"huaweicloud_lb_loadbalancer": lb.ResourceLoadBalancerV2(),
			"huaweicloud_lb_member":       lb.ResourceMemberV2(),
			"huaweicloud_lb_monitor":      lb.ResourceMonitorV2(),
			"huaweicloud_lb_pool":         lb.ResourcePoolV2(),
			"huaweicloud_lb_whitelist":    lb.ResourceWhitelistV2(),

			"huaweicloud_lts_group":  resourceLTSGroupV2(),
			"huaweicloud_lts_stream": resourceLTSStreamV2(),

			"huaweicloud_mls_instance": resourceMlsInstance(),

			"huaweicloud_mapreduce_cluster": mrs.ResourceMRSClusterV2(),
			"huaweicloud_mapreduce_job":     mrs.ResourceMRSJobV2(),

			"huaweicloud_modelarts_dataset":                modelarts.ResourceDataset(),
			"huaweicloud_modelarts_notebook":               modelarts.ResourceNotebook(),
			"huaweicloud_modelarts_notebook_mount_storage": modelarts.ResourceNotebookMountStorage(),

			"huaweicloud_mrs_cluster": ResourceMRSClusterV1(),
			"huaweicloud_mrs_job":     ResourceMRSJobV1(),

			"huaweicloud_nat_dnat_rule": ResourceNatDnatRuleV2(),
			"huaweicloud_nat_gateway":   ResourceNatGatewayV2(),
			"huaweicloud_nat_snat_rule": ResourceNatSnatRuleV2(),

			"huaweicloud_network_acl":              ResourceNetworkACL(),
			"huaweicloud_network_acl_rule":         ResourceNetworkACLRule(),
			"huaweicloud_networking_port":          ResourceNetworkingPortV2(),
			"huaweicloud_networking_secgroup":      ResourceNetworkingSecGroup(),
			"huaweicloud_networking_secgroup_rule": ResourceNetworkingSecGroupRule(),
			"huaweicloud_networking_vip":           vpc.ResourceNetworkingVip(),
			"huaweicloud_networking_vip_associate": resourceNetworkingVIPAssociateV2(),

			"huaweicloud_obs_bucket":        ResourceObsBucket(),
			"huaweicloud_obs_bucket_object": ResourceObsBucketObject(),
			"huaweicloud_obs_bucket_policy": ResourceObsBucketPolicy(),

			"huaweicloud_oms_task": resourceMaasTaskV1(),

			"huaweicloud_rds_account":               rds.ResourceRdsAccount(),
			"huaweicloud_rds_instance":              ResourceRdsInstanceV3(),
			"huaweicloud_rds_parametergroup":        ResourceRdsConfigurationV3(),
			"huaweicloud_rds_read_replica_instance": ResourceRdsReadReplicaInstance(),

			"huaweicloud_servicestage_repo_token_authorization":    servicestage.ResourceRepoTokenAuth(),
			"huaweicloud_servicestage_repo_password_authorization": servicestage.ResourceRepoPwdAuth(),

			"huaweicloud_sfs_access_rule": ResourceSFSAccessRuleV2(),
			"huaweicloud_sfs_file_system": ResourceSFSFileSystemV2(),
			"huaweicloud_sfs_turbo":       ResourceSFSTurbo(),

			"huaweicloud_smn_topic":        ResourceTopic(),
			"huaweicloud_smn_subscription": ResourceSubscription(),

			"huaweicloud_swr_organization":             swr.ResourceSWROrganization(),
			"huaweicloud_swr_organization_permissions": swr.ResourceSWROrganizationPermissions(),
			"huaweicloud_swr_repository":               swr.ResourceSWRRepository(),
			"huaweicloud_swr_repository_sharing":       swr.ResourceSWRRepositorySharing(),

			"huaweicloud_tms_tags": tms.ResourceTmsTag(),

			"huaweicloud_vbs_backup":        resourceVBSBackupV2(),
			"huaweicloud_vbs_backup_policy": resourceVBSBackupPolicyV2(),

			"huaweicloud_vpc_bandwidth":                   eip.ResourceVpcBandWidthV2(),
			"huaweicloud_vpc_eip":                         eip.ResourceVpcEIPV1(),
			"huaweicloud_vpc_eip_associate":               eip.ResourceEIPAssociate(),
			"huaweicloud_vpc_peering_connection":          vpc.ResourceVpcPeeringConnectionV2(),
			"huaweicloud_vpc_peering_connection_accepter": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"huaweicloud_vpc_route_table":                 vpc.ResourceVPCRouteTable(),
			"huaweicloud_vpc":                             vpc.ResourceVirtualPrivateCloudV1(),
			"huaweicloud_vpc_route":                       vpc.ResourceVPCRouteTableRoute(),
			"huaweicloud_vpc_subnet":                      vpc.ResourceVpcSubnetV1(),
			"huaweicloud_vpc_address_group":               vpc.ResourceVpcAddressGroup(),

			"huaweicloud_vpcep_approval": ResourceVPCEndpointApproval(),
			"huaweicloud_vpcep_endpoint": ResourceVPCEndpoint(),
			"huaweicloud_vpcep_service":  ResourceVPCEndpointService(),

			"huaweicloud_vpnaas_endpoint_group":  deprecated.ResourceVpnEndpointGroupV2(),
			"huaweicloud_vpnaas_ike_policy":      deprecated.ResourceVpnIKEPolicyV2(),
			"huaweicloud_vpnaas_ipsec_policy":    deprecated.ResourceVpnIPSecPolicyV2(),
			"huaweicloud_vpnaas_service":         deprecated.ResourceVpnServiceV2(),
			"huaweicloud_vpnaas_site_connection": deprecated.ResourceVpnSiteConnectionV2(),

			"huaweicloud_scm_certificate": scm.ResourceScmCertificate(),

			"huaweicloud_waf_certificate":                waf.ResourceWafCertificateV1(),
			"huaweicloud_waf_domain":                     waf.ResourceWafDomainV1(),
			"huaweicloud_waf_policy":                     waf.ResourceWafPolicyV1(),
			"huaweicloud_waf_rule_blacklist":             waf.ResourceWafRuleBlackListV1(),
			"huaweicloud_waf_rule_data_masking":          waf.ResourceWafRuleDataMaskingV1(),
			"huaweicloud_waf_rule_web_tamper_protection": waf.ResourceWafRuleWebTamperProtectionV1(),
			"huaweicloud_waf_dedicated_instance":         waf.ResourceWafDedicatedInstance(),
			"huaweicloud_waf_dedicated_domain":           waf.ResourceWafDedicatedDomainV1(),
			"huaweicloud_waf_instance_group":             waf.ResourceWafInstanceGroup(),
			"huaweicloud_waf_instance_group_associate":   waf.ResourceWafInstGroupAssociate(),
			"huaweicloud_waf_reference_table":            waf.ResourceWafReferenceTableV1(),

			"huaweicloud_cpts_project": cpts.ResourceProject(),

			// Legacy
			"huaweicloud_networking_eip_associate": eip.ResourceEIPAssociate(),

			"huaweicloud_compute_instance_v2":         ResourceComputeInstanceV2(),
			"huaweicloud_compute_interface_attach_v2": ResourceComputeInterfaceAttachV2(),
			"huaweicloud_compute_keypair_v2":          ResourceComputeKeypairV2(),
			"huaweicloud_compute_servergroup_v2":      ResourceComputeServerGroupV2(),
			"huaweicloud_compute_volume_attach_v2":    ecs.ResourceComputeVolumeAttach(),

			"huaweicloud_dns_ptrrecord_v2": ResourceDNSPtrRecordV2(),
			"huaweicloud_dns_recordset_v2": ResourceDNSRecordSetV2(),
			"huaweicloud_dns_zone_v2":      ResourceDNSZoneV2(),

			"huaweicloud_dcs_instance_v1": dcs.ResourceDcsInstance(),
			"huaweicloud_dds_instance_v3": dds.ResourceDdsInstanceV3(),

			"huaweicloud_fw_firewall_group_v2": resourceFWFirewallGroupV2(),
			"huaweicloud_fw_policy_v2":         resourceFWPolicyV2(),
			"huaweicloud_fw_rule_v2":           resourceFWRuleV2(),

			"huaweicloud_kms_key_v1": ResourceKmsKeyV1(),

			"huaweicloud_lb_certificate_v2":  lb.ResourceCertificateV2(),
			"huaweicloud_lb_loadbalancer_v2": lb.ResourceLoadBalancerV2(),
			"huaweicloud_lb_listener_v2":     lb.ResourceListenerV2(),
			"huaweicloud_lb_pool_v2":         lb.ResourcePoolV2(),
			"huaweicloud_lb_member_v2":       lb.ResourceMemberV2(),
			"huaweicloud_lb_monitor_v2":      lb.ResourceMonitorV2(),
			"huaweicloud_lb_l7policy_v2":     lb.ResourceL7PolicyV2(),
			"huaweicloud_lb_l7rule_v2":       lb.ResourceL7RuleV2(),
			"huaweicloud_lb_whitelist_v2":    lb.ResourceWhitelistV2(),

			"huaweicloud_mrs_cluster_v1": ResourceMRSClusterV1(),
			"huaweicloud_mrs_job_v1":     ResourceMRSJobV1(),

			"huaweicloud_networking_port_v2":          ResourceNetworkingPortV2(),
			"huaweicloud_networking_secgroup_v2":      ResourceNetworkingSecGroup(),
			"huaweicloud_networking_secgroup_rule_v2": ResourceNetworkingSecGroupRule(),

			"huaweicloud_smn_topic_v2":        ResourceTopic(),
			"huaweicloud_smn_subscription_v2": ResourceSubscription(),

			"huaweicloud_rds_instance_v3":       ResourceRdsInstanceV3(),
			"huaweicloud_rds_parametergroup_v3": ResourceRdsConfigurationV3(),

			"huaweicloud_nat_gateway_v2":   ResourceNatGatewayV2(),
			"huaweicloud_nat_snat_rule_v2": ResourceNatSnatRuleV2(),
			"huaweicloud_nat_dnat_rule_v2": ResourceNatDnatRuleV2(),

			"huaweicloud_sfs_access_rule_v2": ResourceSFSAccessRuleV2(),
			"huaweicloud_sfs_file_system_v2": ResourceSFSFileSystemV2(),

			"huaweicloud_iam_agency":    iam.ResourceIAMAgencyV3(),
			"huaweicloud_iam_agency_v3": iam.ResourceIAMAgencyV3(),

			"huaweicloud_vpc_bandwidth_v2":                   eip.ResourceVpcBandWidthV2(),
			"huaweicloud_vpc_eip_v1":                         eip.ResourceVpcEIPV1(),
			"huaweicloud_vpc_peering_connection_v2":          vpc.ResourceVpcPeeringConnectionV2(),
			"huaweicloud_vpc_peering_connection_accepter_v2": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"huaweicloud_vpc_route_v2":                       vpc.ResourceVPCRouteV2(),
			"huaweicloud_vpc_v1":                             vpc.ResourceVirtualPrivateCloudV1(),
			"huaweicloud_vpc_subnet_v1":                      vpc.ResourceVpcSubnetV1(),

			"huaweicloud_cce_cluster_v3": ResourceCCEClusterV3(),
			"huaweicloud_cce_node_v3":    ResourceCCENodeV3(),

			"huaweicloud_as_configuration_v1": as.ResourceASConfiguration(),
			"huaweicloud_as_group_v1":         as.ResourceASGroup(),
			"huaweicloud_as_policy_v1":        as.ResourceASPolicy(),

			"huaweicloud_csbs_backup_v1":        resourceCSBSBackupV1(),
			"huaweicloud_csbs_backup_policy_v1": resourceCSBSBackupPolicyV1(),
			"huaweicloud_vbs_backup_policy_v2":  resourceVBSBackupPolicyV2(),
			"huaweicloud_vbs_backup_v2":         resourceVBSBackupV2(),

			"huaweicloud_maas_task":    resourceMaasTaskV1(),
			"huaweicloud_maas_task_v1": resourceMaasTaskV1(),

			"huaweicloud_identity_project_v3":          iam.ResourceIdentityProjectV3(),
			"huaweicloud_identity_role_assignment_v3":  iam.ResourceIdentityRoleAssignmentV3(),
			"huaweicloud_identity_user_v3":             iam.ResourceIdentityUserV3(),
			"huaweicloud_identity_group_v3":            iam.ResourceIdentityGroupV3(),
			"huaweicloud_identity_group_membership_v3": iam.ResourceIdentityGroupMembershipV3(),
			"huaweicloud_identity_provider_conversion": iam.ResourceIAMProviderConversion(),

			"huaweicloud_cdm_cluster_v1": cdm.ResourceCdmCluster(),
			"huaweicloud_ges_graph_v1":   ResourceGesGraphV1(),
			"huaweicloud_css_cluster_v1": css.ResourceCssCluster(),
			"huaweicloud_dis_stream_v2":  dis.ResourceDisStream(),

			"huaweicloud_cs_cluster_v1":         deprecated.ResourceCsClusterV1(),
			"huaweicloud_cs_peering_connect_v1": deprecated.ResourceCsPeeringConnectV1(),

			"huaweicloud_vpnaas_ipsec_policy_v2":    deprecated.ResourceVpnIPSecPolicyV2(),
			"huaweicloud_vpnaas_service_v2":         deprecated.ResourceVpnServiceV2(),
			"huaweicloud_vpnaas_ike_policy_v2":      deprecated.ResourceVpnIKEPolicyV2(),
			"huaweicloud_vpnaas_endpoint_group_v2":  deprecated.ResourceVpnEndpointGroupV2(),
			"huaweicloud_vpnaas_site_connection_v2": deprecated.ResourceVpnSiteConnectionV2(),

			"huaweicloud_dli_queue_v1":                dli.ResourceDliQueue(),
			"huaweicloud_cs_route_v1":                 deprecated.ResourceCsRouteV1(),
			"huaweicloud_networking_vip_v2":           vpc.ResourceNetworkingVip(),
			"huaweicloud_networking_vip_associate_v2": resourceNetworkingVIPAssociateV2(),
			"huaweicloud_fgs_function_v2":             fgs.ResourceFgsFunctionV2(),
			"huaweicloud_cdn_domain_v1":               resourceCdnDomainV1(),

			// Deprecated
			"huaweicloud_blockstorage_volume_v2":          resourceBlockStorageVolumeV2(),
			"huaweicloud_networking_network_v2":           resourceNetworkingNetworkV2(),
			"huaweicloud_networking_subnet_v2":            resourceNetworkingSubnetV2(),
			"huaweicloud_networking_floatingip_v2":        resourceNetworkingFloatingIPV2(),
			"huaweicloud_networking_router_v2":            resourceNetworkingRouterV2(),
			"huaweicloud_networking_router_interface_v2":  resourceNetworkingRouterInterfaceV2(),
			"huaweicloud_networking_router_route_v2":      resourceNetworkingRouterRouteV2(),
			"huaweicloud_ecs_instance_v1":                 resourceEcsInstanceV1(),
			"huaweicloud_compute_secgroup_v2":             ResourceComputeSecGroupV2(),
			"huaweicloud_compute_floatingip_v2":           ResourceComputeFloatingIPV2(),
			"huaweicloud_compute_floatingip_associate_v2": ResourceComputeFloatingIPAssociateV2(),
			"huaweicloud_images_image_v2":                 deprecated.ResourceImagesImageV2(),
			"huaweicloud_dms_instance_v1":                 deprecated.ResourceDmsInstancesV1(),
			"huaweicloud_dms_group_v1":                    deprecated.ResourceDmsGroups(),
			"huaweicloud_dms_group":                       deprecated.ResourceDmsGroups(),
			"huaweicloud_dms_queue":                       deprecated.ResourceDmsQueues(),
			"huaweicloud_dms_queue_v1":                    deprecated.ResourceDmsQueues(),
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

		"shared_config_file": "The path to the shared config file. If not set, the default is ~/.hcloud/config.json.",

		"profile": "The profile name as set in the shared config file.",

		"max_retries": "How many times HTTP connection should be retried until giving up.",

		"enterprise_project_id": "enterprise project id",
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData, terraformVersion string) (interface{},
	diag.Diagnostics) {
	var tenantName, tenantID, delegated_project, identityEndpoint string
	region := d.Get("region").(string)
	cloud := d.Get("cloud").(string)

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
		delegated_project = v.(string)
	} else {
		delegated_project = region
	}

	// use auth_url as identityEndpoint if provided
	if v, ok := d.GetOk("auth_url"); ok {
		identityEndpoint = v.(string)
	} else {
		// use cloud as basis for identityEndpoint
		if cloud == defaultCloud {
			identityEndpoint = fmt.Sprintf("https://iam.%s:443/v3", cloud)
		} else {
			identityEndpoint = fmt.Sprintf("https://iam.%s.%s:443/v3", region, cloud)
		}
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
		DelegatedProject:    delegated_project,
		Cloud:               cloud,
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
	if len(assumeRoleList) == 1 {
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
