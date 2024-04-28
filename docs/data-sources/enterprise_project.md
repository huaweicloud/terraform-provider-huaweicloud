---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project"
description: ""
---

# huaweicloud_enterprise_project

Use this data source to get an enterprise project from HuaweiCloud

## Example Usage

```hcl
data "huaweicloud_enterprise_project" "test" {
  name = "test"
}
```

## Resources Supported Currently

<!-- markdownlint-disable MD033 -->
Service Name | Resource Name
---- | ---
AOM | huaweicloud_aom_cmdb_application
APIG | huaweicloud_apig_instance
AS  | huaweicloud_as_group
BCS | huaweicloud_bcs_instance
BMS | huaweicloud_bms_instance
CBR | huaweicloud_cbr_vault
CC  | huaweicloud_cc_bandwidth_package<br>huaweicloud_cc_connection<br>huaweicloud_central_network
CCE | huaweicloud_cce_cluster
CCI | huaweicloud_cci_namespace
CCM | huaweicloud_ccm_private_certificate<br>huaweicloud_ccm_private_ca
CDM | huaweicloud_cdm_cluster
CDN | huaweicloud_cdn_domain
CES | huaweicloud_ces_alarmrule<br>huaweicloud_ces_resource_group
CodeArts | huaweicloud_codearts_project
CPH | huaweicloud_cph_server
CSE | huaweicloud_cse__microservice_engine
CSS | huaweicloud_css_cluster
DataArts | huaweicloud_dataarts_studio_instance
DBSS | huaweicloud_dbss_instance
DC  | huaweicloud_dc_virtual_gateway<br>huaweicloud_dc_virtual_interface
DCS | huaweicloud_dcs_instance
DDM | huaweicloud_ddm_instance
DDS | huaweicloud_dds_instance
DEW | huaweicloud_kms_key
DIS | huaweicloud_dis_stream
DLI | huaweicloud_dli_database<br>huaweicloud_dli_queue
DMS | huaweicloud_dms_kafka_instance<br>huaweicloud_dms_rabbitmq_instance<br>huaweicloud_dms_rocketmq_instance
DNS | huaweicloud_dns_ptrrecord<br>huaweicloud_dns_zone
DRS | huaweicloud_drs_job
DWS | huaweicloud_dws_cluster
ECS | huaweicloud_compute_instance
EG  | huaweicloud_eg_custom_event_channel
EIP | huaweicloud_vpc_eip<br>huaweicloud_vpc_bandwidth
ELB | huaweicloud_elb_loadbalancer<br>huaweicloud_elb_certificate<br>huaweicloud_elb_ipgroup<br>huaweicloud_elb_security_policy
ER  | huaweicloud_er_instance
EVS | huaweicloud_evs_volume
FGS | huaweicloud_fgs_function
GA  | huaweicloud_ga_accelerator
GaussDB | huaweicloud_gaussdb_cassandra_instance<br>huaweicloud_gaussdb_influx_instance<br>huaweicloud_gaussdb_mongo_instance<br>huaweicloud_gaussdb_mysql_instance<br>huaweicloud_gaussdb_opengauss_instance<br>huaweicloud_gaussdb_redis_instance
GES | huaweicloud_ges_graph
HSS | huaweicloud_hss_host_group
IAM | huaweicloud_identity_group_role_assignment<br>huaweicloud_identity_user_role_assignment
IMS | huaweicloud_images_image<br>huaweicloud_images_image_copy
LB  | huaweicloud_lb_certificate<br>huaweicloud_lb_loadbalancer
LTS | huaweicloud_lts_stream<br>huaweicloud_lts_waf_access<br>huaweicloud_lts_search_criteria
ModelArts | huaweicloud_modelarts_workspace
MRS | huaweicloud_mapreduce_cluster
NAT | huaweicloud_nat_gateway<br>huaweicloud_nat_private_gateway<br>huaweicloud_nat_private_transit_ip
OBS | huaweicloud_obs_bucket
RDS | huaweicloud_rds_instance<br>huaweicloud_rds_read_replica_instance
ServicesTage | huaweicloud_servicestage_application<br>huaweicloud_servicestage_environment
SFS | huaweicloud_sfs_file_system<br>huaweicloud_sfs_turbo
SMN | huaweicloud_smn_topic
VPC | huaweicloud_vpc<br>huaweicloud_networking_secgroup<br>huaweicloud_vpc_address_group
VPN | huaweicloud_vpn_gateway
WAF | huaweicloud_waf_address_group<br>huaweicloud_waf_certificate<br>huaweicloud_waf_cloud_instance<br>huaweicloud_waf_dedicated_domain<br>huaweicloud_waf_dedicated_instance<br>huaweicloud_waf_domain<br>huaweicloud_waf_policy<br>huaweicloud_waf_reference_table<br> huaweicloud_waf_rule_anti_crawler<br>huaweicloud_waf_rule_blacklist<br>huaweicloud_waf_rule_cc_protection<br>huaweicloud_waf_rule_data_masking<br>huaweicloud_waf_rule_geolocation_access_control<br>huaweicloud_waf_rule_global_protection_whitelist<br>huaweicloud_waf_rule_information_leakage_prevention<br>huaweicloud_waf_rule_known_attack_source<br>huaweicloud_waf_rule_precise_protection<br>huaweicloud_waf_rule_web_tamper_protection
WorkSpace | huaweicloud_workspace_desktop<br>
<!-- markdownlint-enable MD033 -->

## Argument Reference

* `name` - (Optional, String) Specifies the enterprise project name. Fuzzy search is supported.

* `id` - (Optional, String) Specifies the ID of an enterprise project. The value 0 indicates enterprise project default.

* `status` - (Optional, Int) Specifies the status of an enterprise project.
    + 1 indicates Enabled.
    + 2 indicates Disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Provides supplementary information about the enterprise project.

* `created_at` - Specifies the time (UTC) when the enterprise project was created. Example: 2018-05-18T06:49:06Z

* `updated_at` - Specifies the time (UTC) when the enterprise project was modified. Example: 2018-05-28T02:21:36Z
