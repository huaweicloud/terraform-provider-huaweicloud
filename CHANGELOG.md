# CHANGELOG

## 1.28.1 (September 3, 2021)

ENHANCEMENTS:

* resource/huaweicloud_cdn_domain: Add `service_area` argument support [GH-1466]

BUG FIXES:

* resource/huaweicloud_cce_node: fix an issue when unsubscribing a prePaid node [GH-1464]

## 1.28.0 (August 31, 2021)

FEATURES:

* **New Resurce:**
  + `huaweicloud_dms_kafka_topic` [GH-1379]
  + `huaweicloud_gaussdb_redis_instance` [GH-1399]
  + `huaweicloud_waf_dedicated_instance` [GH-1407]
  + `huaweicloud_waf_dedicated_domain` [GH-1409]
  + `huaweicloud_waf_reference_table` [GH-1426]

* **New Data Source:**
  + `huaweicloud_fgs_dependencies` [GH-1419]
  + `huaweicloud_gaussdb_redis_instance` [GH-1399]
  + `huaweicloud_gaussdb_cassandra_instances` [GH-1406]
  + `huaweicloud_waf_policies` [GH-1374]
  + `huaweicloud_waf_dedicated_instances` [GH-1420]
  + `huaweicloud_waf_reference_tables` [GH-1435]

ENHANCEMENTS:

* provider: Upgrade to terraform-plugin-sdk v2 [GH-1381]
* config: Add validation of domain name [GH-1429]
* resource/huaweicloud_cce_node_pool: Support to update `labels` and `taints` field [GH-1385]
* resource/huaweicloud_apig_application: Add app key and app secret attributes [GH-1401]
* resource/huaweicloud_mapreduce_cluster: Support custom type and `host_ips` field [GH-1436]

BUG FIXES:

* resource/huaweicloud_compute_instance: update sdk to fix destroying error [GH-1397]
* resource/huaweicloud_compute_servergroup: validate the server ID before removing [GH-1380]
* resource/huaweicloud_mapreduce_cluster: get subnet_name and assemble it into creation option [GH-1396]
* resource/huaweicloud_vpc_eip: set the eip status to *BOUND* after binding to a port [GH-1398]
* resource/huaweicloud_vpc_eip: use v1 API to allocate EIP in per-use mod [GH-1434]

DEPRECATE:

* resource/huaweicloud_cs_cluster [GH-1428]
* resource/huaweicloud_cs_route [GH-1428]
* resource/huaweicloud_cs_peering_connect [GH-1428]

## 1.27.2 (August 26, 2021)

FEATURES:

* **New Data Source:**
  + `huaweicloud_gaussdb_cassandra_dedicated_resource` [GH-1412]
  + `huaweicloud_gaussdb_mysql_dedicated_resource` [GH-1415]

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_cassandra_instance: Add dedicated_resource_id support [GH-1414]
* resource/huaweicloud_gaussdb_mysql_instance: Add dedicated_resource_id support [GH-1416]

## 1.27.1 (August 13, 2021)

FEATURES:

* **New Resurce:**
  + `huaweicloud_mapreduce_job` [GH-1324]
  + `huaweicloud_apig_api` [GH-1360]

ENHANCEMENTS:

* config: Introduce the retry func (honor 429 http code) [GH-1351]
* resource/huaweicloud_scm_certificate: Support import function [GH-1342]
* resource/huaweicloud_cce_cluster: Support hibernate/awake action [GH-1344]
* resource/huaweicloud_dli_queue: Make `cu_count` updatable [GH-1347]
* resource/huaweicloud_rds_instance: Make `db.0.port` and `security_group_id` updatable [GH-1317]

BUG FIXES:

* resource/huaweicloud_fgs_function: Mark runtimue parameter be forcenew [GH-1361]

## 1.27.0 (July 31, 2021)

FEATURES:

* **New Resurce:**
  + `huaweicloud_apig_throttling_policy` ([#1296](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1296))
  + `huaweicloud_apig_custom_authorizer` ([#1297](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1297))
  + `huaweicloud_mapreduce_cluster` ([#1324](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1324))
  + `huaweicloud_cce_node_attach` ([#1326](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1326))

ENHANCEMENTS:

* resource/huaweicloud_cce_node: Add ability to remove cce node by `keep_ecs` ([#1314](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1314))
* resource/huaweicloud_network_acl_rule: cancel the MaxItems limitation of inbound_rules and outbound_rules ([#1315](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1315))
* resource/huaweicloud_gaussdb_mysql_instance: Support to enlarge proxy node ([1258](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1258))
* enterprise_project_id support:
  + `huaweicloud_smn_topic`: ([1305](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1305))
  + `huaweicloud_css_cluster`: ([1307](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1307))
  + `huaweicloud_dis_stream`: ([1313](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1313))
  + `huaweicloud_dws_cluster`: ([1313](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1313))
  + `huaweicloud_dli_queue`: ([1321](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1321))

BUG FIXES:

* resource/huaweicloud_obs_bucket: Support to create parallel file system ([#1312](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1312))

## 1.26.1 (July 23, 2021)

FEATURES:

* **New Resurce:**
  + `huaweicloud_waf_certificate` ([#1255](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1255))
  + `huaweicloud_waf_domain` ([#1255](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1255))
  + `huaweicloud_waf_policy` ([#1257](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1257))
  + `huaweicloud_waf_rule_blacklist` ([#1283](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1283))
  + `huaweicloud_waf_rule_data_masking` ([#1295](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1295))
  + `huaweicloud_waf_rule_web_tamper_protection` ([#1298](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1298))
  + `huaweicloud_apig_environment` ([#1267](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1267))
  + `huaweicloud_apig_group` ([#1284](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1284))
  + `huaweicloud_apig_response` ([#1294](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1294))
  + `huaweicloud_apig_vpc_channel` ([#1273](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1273))

* **New Data Source:**
  + `huaweicloud_waf_certificate` ([#1279](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1279))
  + `huaweicloud_elb_certificate` ([#1301](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1301))
  + `huaweicloud_lb_certificate` ([#1303](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1303))

ENHANCEMENTS:

* provider: use `cloud` value as basis for `auth_url` ([#1285](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1285))
* resource/huaweicloud_vpc_eip: Add tags support ([#1262](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1262))
* resource/huaweicloud_dds_instance: Support to update `flavor` field ([#1286](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1286))

BUG FIXES:

* resource/huaweicloud_compute_instance: Fix power action error ([#1268](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1268))
* resource/huaweicloud_identity_role: Support the policy for cloud services and agencies ([#1289](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1289))

## 1.26.0 (June 30, 2021)

FEATURES:

* **New Resource:** `huaweicloud_apig_instance` ([#1221](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1221))
* **New Resource:** `huaweicloud_apig_application` ([#1198](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1198))

ENHANCEMENTS:

* resource/huaweicloud_obs_bucket: Support to enable multi-AZ mode ([#1190](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1190))
* resource/huaweicloud_vpc_eip: support prePaid charging mode ([#963](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/963))
* resource/huaweicloud_gaussdb_mysql_instance: Add proxy support for gaussdb mysql ([#1136](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1136))
* resource/huaweicloud_cce_*: expand PollInterval in WaitForState to avoid the API rate limits ([#1251](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1251))

BUG FIXES:

* resource/huaweicloud_rds_instance: Fix the exception of empty value conversion ([#1204](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1204))

## 1.25.1 (June 24, 2021)

FEATURES:

* **New Resource:** `huaweicloud_scm_certificate` ([#1218](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1218))

ENHANCEMENTS:

* resource/huaweicloud_fgs_function: Add urn and version support ([#1203](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1203))
* resource/huaweicloud_gaussdb_mysql_instance: Add volume_size support ([#1201](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1201))

## 1.25.0 (May 31, 2021)

FEATURES:

* **New Data Source:** `huaweicloud_cce_addon_template` ([#1039](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1039))
* **New Data Source:** `huaweicloud_iec_port` ([#1152](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1152))
* **New Data Source:** `huaweicloud_iec_vpc` ([#1152](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1152))
* **New Data Source:** `huaweicloud_iec_vpc_subnets` ([#1152](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1152))
* **New Data Source:** `huaweicloud_iec_network_acl` ([#1159](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1159))
* **New Data Source:** `huaweicloud_iec_security_group` ([#1159](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1159))
* **New Data Source:** `huaweicloud_iec_eips` ([#1164](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1164))
* **New Data Source:** `huaweicloud_iec_keypair` ([#1164](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1164))
* **New Data Source:** `huaweicloud_iec_server` ([#1169](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1169))
* **New Resource:** `huaweicloud_as_lifecycle_hook` ([#1069](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1069))
* **New Resource:** `huaweicloud_cci_pvc` ([#1081](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1081))
* **New Resource:** `huaweicloud_elb_listener` ([#1021](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1021))
* **New Resource:** `huaweicloud_elb_certificate` ([#1148](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1148))
* **New Resource:** `huaweicloud_elb_ipgroup` ([#1148](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1148))
* **New Resource:** `huaweicloud_elb_pool` ([#1150](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1150))
* **New Resource:** `huaweicloud_elb_member` ([#1150](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1150))
* **New Resource:** `huaweicloud_elb_l7policy` ([#1161](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1161))
* **New Resource:** `huaweicloud_elb_l7rule` ([#1161](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1161))
* **New Resource:** `huaweicloud_elb_monitor` ([#1163](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1163))
* **New Resource:** `huaweicloud_dms_kafka_instance` ([#1162](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1162))
* **New Resource:** `huaweicloud_dms_rabbitmq_instance` ([#1170](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1170))

ENHANCEMENTS:

* resource/huaweicloud_cce_node: Support import function ([#958](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/958))
* resource/huaweicloud_cce_node_pool: Support import function ([#1005](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1005))
* resource/huaweicloud_compute_instance: Support power action ([#914](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/914))
* resource/huaweicloud_fgs_function: Support to update parameters ([#1140](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1040))
* resource/huaweicloud_iec_vpc_subnet: Set the default DNS list if it was not specified ([#1157](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1157))
* enterprise_project_id support:
  + `huaweicloud_ces_alarmrule`: ([1137](https://github.com/huaweicloudterraform-provider-huaweicloud/pull/1137))
  + `huaweicloud_dds_instance`: ([#1145](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1145))

BUG FIXES:

* data/huaweicloud_vpc_route_ids: Make list instead of set for ids ([#1141](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1141))
* resource/huaweicloud_compute_instance: Support security_group_ids parameter ([#1128](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1128))

DEPRECATE:

* resource/huaweicloud_dms_instance: ([#1176](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1176))

## 1.24.2 (May 18, 2021)

BUG FIXES:

* resource/huaweicloud_sfs_turbo: Remove SFS turbo size validation ([#1140](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1140))

## 1.24.1 (May 12, 2021)

FEATURES:

* **New Data Source:** `huaweicloud_lb_loadbalancer` ([#1113](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1113))

ENHANCEMENTS:

* data/huaweicloud_identity_role: Support to filter system-defined IAM role by display_name ([#1105](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1105))
* resource/huaweicloud_ces_alarmrule: Support update function ([#1116](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1116))
* resource/huaweicloud_lb_loadbalancer: Support import function ([#1125](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1125))
* resource/huaweicloud_compute_instance: Support ESSD type for volume_type field ([#1126](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1126))

DEPRECATE:

* resource/huaweicloud_cts_tracker: ([#1102](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1102))
* data/huaweicloud_cts_tracker: ([#1102](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1102))

## 1.24.0 (April 30, 2021)

FEATURES:

* **New Data Source:** `huaweicloud_cce_node_pool` ([#1005](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1005))
* **New Resource:** `huaweicloud_swr_organization` ([#428](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/428))
* **New Resource:** `huaweicloud_bcs_instance` ([#1064](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1064))
* **New Resource:** `huaweicloud_bms_instance` ([#1024](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1024))
* **New Resource:** `huaweicloud_cbr_policy` ([#1025](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1025))
* **New Resource:** `huaweicloud_cbr_vault` ([#1025](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1025))

ENHANCEMENTS:

* Provider: Support `security_token` to authenticate with a temporary security credential([#1062](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1062))
* Support `enterprise_project_id` in AS group, SFS turbo, ELB loadbalancer, IMS and DNS ([#1019](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1019))
* data/huaweicloud_rds_flavors: support to filter rds flavors with replica mode ([#1070](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1070))
* resource/huaweicloud_lb_monitor: Add port option to lb_monitor ([#1059](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1059))
* resource/huaweicloud_rds_instance: support prePaid charging mode ([#1066](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1066))
* resource/huaweicloud_vpc_subnet: Try to set default DNS server if it was not specified ([#1074](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1074))
* resource/huaweicloud_nat_gateway: Unify network parameters ([#1087](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1087))
* resource/huaweicloud_ces_alarmrule: Add alarm_level parameter ([#1085](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1085))
* resource/huaweicloud_fgs_function:
  + Rename xrole and app parameters ([#1076](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1095))
  + `func_code` support both base64 and plain text format ([#1077](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1077))

BUG FIXES:

* resource/huaweicloud_rds_instance: Fix RDS deployment crash with v1.23.1 version ([#1054](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1054))
* resource/huaweicloud_dcs_instance: Fix DCS backup policy issue for single instance ([#1092](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1092))

## 1.23.1 (April 7, 2021)

BUG FIXES:

* resource/huaweicloud_dds_instance: Fix backup_strategy update issue ([#1041](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1041))
* resource/huaweicloud_cce_node: Unsubscribe eip as well in prePaid mode ([#1043](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1043))

## 1.23.0 (April 2, 2021)

FEATURES:

* **New Resource:** `huaweicloud_identity_acl` ([#982](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/982))

ENHANCEMENTS:

* resource/huaweicloud_vpc & huaweicloud_vpc_subent: Support IPv6 ([#989](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/989))

* resource/huaweicloud_cce_node_pool: Support `tags` field ([#980](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/980))

* resource/huaweicloud_cce_node:
  + Support prePaid charging mode ([#1001](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1001))
  + Add possibility to set `runtime` ([#1026](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1026))

* resource/huaweicloud_cce_cluster:
  + Support prePaid charging mode ([#1027](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1027))
  + Add `delete_*` parameters to delete associated resources ([#1007](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1007))

* resource/huaweicloud_api_gateway_api: Support `CORS` field ([#1015](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/1015))

* resource/huaweicloud_as_group: Support `enterprise_project_id` field ([#1028](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1028))

* resource/huaweicloud_sfs_turbo: Support `enterprise_project_id` field ([#1030](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1030))

BUG FIXES:

* data/huaweicloud_dcs_az: Filter avaliable zones by code ([#990](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/990))

* resource/huaweicloud_vpcep_approval: Make vpcep approval can work cross-project ([#1010](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1010))

Removed:

* data/huaweicloud_s3_bucket_object: ([#973](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/973))
* data/huaweicloud_rds_flavors_v1: ([#1032](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1032))
* resource/huaweicloud_rds_instance_v1: ([#1032](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/1032))
* resource/huaweicloud_s3_bucket: ([#973](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/973))
* resource/huaweicloud_s3_bucket_policy: ([#973](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/973))
* resource/huaweicloud_s3_bucket_object: ([#973](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/973))

## 1.22.3 (March 26, 2021)

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_mysql_instance: Add table_name_case_sensitivity support ([#998](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/998))

## 1.22.2 (March 19, 2021)

ENHANCEMENTS:

Do not fetch twice the first page in AllPages request ([#981](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/981))

## 1.22.1 (March 5, 2021)

BUG FIXES:

* resource/huaweicloud_obs_bucket: Fix wrong bucket domain name in customizing cloud scene ([#957](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/957))

* resource/huaweicloud_gaussdb_opengauss_instance: Set sharding_num and coordinator_num default to 3 ([#959](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/959))

* resource/huaweicloud_cce_node & huaweicloud_cce_node_pool: revoke `extend_param` and set to deprecated  ([#966](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/966))

* resource/huaweicloud_vpc & data/huaweicloud_vpc: revoke `shared` attribute  ([#967](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/967))

DEPRECATE:

* the `tenant_id` is marked as deprecated in resources ([#952][#954])

## 1.22.0 (February 27, 2021)

ENHANCEMENTS:

* resource/huaweicloud_networking_secgroup_rule: Support `description` field ([#905](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/905))

* resource/huaweicloud_compute_servergroup: Support attach ECS to server group ([#913](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/913))

* resource/huaweicloud_networking_vip: Support import of virtual IP ([#915](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/915))

* resource/huaweicloud_cce_cluster: Support eni network for turbo cluster ([#934](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/934))

* resource/huaweicloud_cce_node:
    + Support ECS group_id param ([#936](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/936))
    + Support `extend_param`, `fixed_ip` and `hw_passthrough` ([#947](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/947))

* resource/huaweicloud_identity_agency: Support `duration` param ([#946](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/946))

* resource/huaweicloud_identity_user: Support `email` and `phone` param (([#949](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/949)))

BUG FIXES:

* resource/huaweicloud_evs_volume: Fix missing fields when importing volume ([#916](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/916))

## 1.21.2 (February 19, 2021)

BUG FIXES:

* resource/huaweicloud_lb_pool: support UDP protocol and update docs ([#923](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/923))

* Get `enterprise_project_id` form the config when it was empty in resources and data sources ([#910](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/910))

## 1.21.1 (February 7, 2021)

ENHANCEMENTS:

* provider: Support to customize nat service endpoint by `nat` key ([#899](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/899))

* data/huaweicloud_nat_gateway: Support to query by enterprise_project_id ([#891](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/891))

* resource/huaweicloud_gaussdb_opengauss_instance: Support to update name and password ([#898](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/898))

* resource/huaweicloud_cce_cluster: Support `service_network_cidr` parm ([#901](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/901))

* resource/huaweicloud_cce_node: change the type of volume/extend_param to map ([#904](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/904))

BUG FIXES:

* resource/huaweicloud_cce_cluster: Fix validate bug when `masters` param is empty ([#892](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/892))

* resource/huaweicloud_sfs_file_system: Make access_type and access_level to be computed ([#902](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/902))

## 1.21.0 (February 1, 2021)

ENHANCEMENTS:

* data/huaweicloud_gaussdb_mysql_instances: Allow to get no instances with a search criteria ([#872](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/872))

* resource/huaweicloud_iec_vip: Support to associate ports ([#876](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/876))

* resource/huaweicloud_cce_cluster: Support `masters` parm ([#885](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/885))

* resource/huaweicloud_rds_instance: Support to update `name` ([#888](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/888))

BUG FIXES:

* data/huaweicloud_iec_flavors: Support to query iec flavors with name and site_ids ([#859](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/859))

* resource/huaweicloud_iec_eip: Fix the resource can't be destroyed when bind with port ([#857](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/857))

* reousrce/huaweicloud_iec_network_acl: Fix `networks` attribute can't be importted ([#871](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/871))

* reousrce/huaweicloud_compute_eip_associate: Fix API response code 202 ([#878](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/878))

* resource/huaweicloud_nat_dnat_rule: Fix `internal_service_port` and `internal_service_port` can't be 0 issue ([#880](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/880))

## 1.20.4 (January 26, 2021)

FEATURES:

* **New Data Source:** `huaweicloud_gaussdb_mysql_instances` ([#855](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/855))

ENHANCEMENTS:

* resource/huaweicloud_cce_node: Set subnet_id attribute ([#841](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/841))

## 1.20.3 (January 19, 2021)

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_opengauss_instance: Add backup_strategy update support ([#823](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/823))

## 1.20.2 (December 28, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_vpcep_public_services` ([#769](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/769))
* **New Data Source:** `huaweicloud_iec_flavors` ([#779](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/779))
* **New Data Source:** `huaweicloud_iec_images` ([#780](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/780))
* **New Data Source:** `huaweicloud_iec_sites` ([#782](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/782))
* **New Resource:** `huaweicloud_identity_role` ([#761](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/761))
* **New Resource:** `huaweicloud_vpcep_service` ([#766](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/766))
* **New Resource:** `huaweicloud_vpcep_endpoint` ([#772](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/772))
* **New Resource:** `huaweicloud_iec_vpc` ([#775](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/775))
* **New Resource:** `huaweicloud_vpcep_approval` ([#783](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/783))

ENHANCEMENTS:

* resource/huaweicloud_compute_instance: Add status and public_ip attributes ([#750](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/750))
* resource/huaweicloud_rds_instance: Add time_zone attribute ([#751](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/751))

BUG FIXES:

* resource/huaweicloud_cce_node_pool: Fix initial_node_count can't be 0 issue ([#757](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/757))

## 1.20.1 (December 16, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_vpc_eip` ([#743](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/743))
* **New Data Source:** `huaweicloud_compute_instance` ([#744](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/744))

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_mysql_instance: Add prePaid support ([#733](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/733))
* resource/huaweicloud_compute_instance: Add fault_domain support ([#735](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/735))
* resource/huaweicloud_gaussdb_cassandra_instance: Add prePaid support ([#740](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/740))
* Add custom endpoints support ([#741](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/741))

## 1.20.0 (November 30, 2020)

ENHANCEMENTS:

* Update resource Attributes Reference in docs ([#715](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/715))

## 1.19.3 (November 28, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_gaussdb_mysql_instance` ([#682](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/682))
* **New Data Source:** `huaweicloud_gaussdb_cassandra_instance` ([#690](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/690))
* **New Data Source:** `huaweicloud_gaussdb_opengauss_instance` ([#699](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/699))
* **New Resource:** `huaweicloud_images_image` ([#706](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/706))

ENHANCEMENTS:

* resource/huaweicloud_kms_key: Add enterprise_project_id suppport ([#693](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/693))

## 1.19.2 (November 19, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_vpc_bandwidth` ([#595](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/595))
* **New Data Source:** `huaweicloud_compute_flavors` ([#609](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/609))
* **New Data Source:** `huaweicloud_enterprise_project` ([#620](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/620))
* **New Resource:** `huaweicloud_css_snapshot` ([#603](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/603))

ENHANCEMENTS:

* resource/huaweicloud_css_cluster: Add security mode support ([#592](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/592))
* Add enterprise_project_id support to cce_cluster, rds, obs, sfs ([#593](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/593))
* resource/huaweicloud_cce_node: Add tags support ([#598](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/598))
* Add tags support to dns and vpn resources ([#599](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/599))
* Make name argument support Chinese character ([#600](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/600))
* Add enterprise_project_id support to dcs, nat_gateway ([#601](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/601))
* Add enterprise_project_id to secgroup ([#606](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/606))
* Add tags support to rds instance resource ([#607](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/607))
* Add tags support to dds and dcs instance resource ([#610](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/610))
* Add tags support to dms instance resource ([#611](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/611))
* Add tags support to elb resources ([#613](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/613))
* resource/huaweicloud_mrs_cluster: Set login mode default to keypair ([#614](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/614))
* Add resource-level region support ([#616](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/616))
* resource/huaweicloud_mrs_cluster: Add tags support ([#617](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/617))
* resource/huaweicloud_mrs_cluster: Add support to login cluster with password ([#628](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/628))
* resource/huaweicloud_networking_vip: Make subnet_id parameter optional ([#648](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/648))
* resource/huaweicloud_networking_vip_associate: Make port_ids updatable ([#650](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/650))
* resource/huaweicloud_gaussdb_mysql_instance: Add force_import support ([#654](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/654))
* resource/huaweicloud_gaussdb_cassandra_instance: Add force_import support ([#656](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/656))
* resource/huaweicloud_gaussdb_opengauss_instance: Add force_import support ([#658](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/658))

BUG FIXES:

* resource/huaweicloud_oms_task: Fix endpoint issue ([#651](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/651))
* resource/huaweicloud_cce_addon: Fix value type issue ([#657](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/657))

## 1.19.1 (October 15, 2020)

ENHANCEMENTS:

* resource/huaweicloud_cce_node_pool: Add type parameter support ([#554](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/554))
* Update max_retries default to 5 ([#577](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/577))
* resource/huaweicloud_obs_bucket: Add obs bucket quota support ([#579](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/579))
* Add enterprise_project_id to vpc, eip, and bandwidth resources ([#585](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/585))
* resource/huaweicloud_dcs_instance: Make whitelists optional for Redis 4.0 and 5.0 ([#588](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/588))
* resource/huaweicloud_dcs_instance: Update capacity from into to float ([#589](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/589))

## 1.19.0 (September 16, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_nat_gateway` ([#501](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/501))
* **New Data Source:** `huaweicloud_gaussdb_mysql_configuration` ([#529](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/529))
* **New Resource:** `huaweicloud_cce_node_pool` ([#511](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/511))

ENHANCEMENTS:

* resource/huaweicloud_dcs_instance: Add IP whitelists support for Redis 4.0 and 5.0 ([#510](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/510))
* resource/huaweicloud_cce_cluster: Add kube_config_raw support ([#512](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/512))
* data/huaweicloud_cce_cluster: Add TLS certificates support ([#516](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/516))
* resource/huaweicloud_gaussdb_cassandra_instance: Add configuration_id update support ([#522](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/522))
* resource/huaweicloud_evs_volume: Add evs volume extend support ([#524](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/524))
* resource/huaweicloud_compute_instance: Add system disk extend support ([#527](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/527))

## 1.18.1 (August 31, 2020)

* Add subcategories to frontmatter for Terrafrom Registry website

## 1.18.0 (August 29, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_obs_bucket` ([#482](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/482))
* **New Resource:** `huaweicloud_evs_volume` ([#429](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/429))
* **New Resource:** `huaweicloud_sfs_turbo` ([#433](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/433))
* **New Resource:** `huaweicloud_lts_group` ([#446](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/446))
* **New Resource:** `huaweicloud_lts_stream` ([#446](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/446))
* **New Resource:** `huaweicloud_sfs_access_rule` ([#451](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/451))
* **New Resource:** `huaweicloud_cce_addon` ([#484](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/484))
* **New Resource:** `huaweicloud_network_acl_rule` ([#495](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/495))
* **New Resource:** `huaweicloud_network_acl` ([#496](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/496))

ENHANCEMENTS:

* resource/huaweicloud_compute_instance: Add disk related parameters ([#440](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/440))
* resource/huaweicloud_gaussdb_opengauss: Set security_group_id to Optional ([#445](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/445))
* resource/huaweicloud_compute_instance: Add enterprise_project_id support ([#450](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/450))
* resource/huaweicloud_gaussdb_cassandra: Add extend-volume support ([#444](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/444))
* provider: Add max_retries support ([#463](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/463))

## 1.17.0 (July 31, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_availability_zones` ([#376](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/376))
* **New Resource:** `huaweicloud_obs_bucket_policy` ([#407](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/407))

ENHANCEMENTS:

* resource/huaweicloud_compute_instance_v2: Add Sensitive to admin_pass argument ([#370](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/370))
* resource/huaweicloud_vpc_eip_v1: Add address attribute ([#375](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/375))
* resource/huaweicloud_cce_node_v3: Add eip_id argument support ([#380](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/380))
* resource/huaweicloud_compute_eip_associate_v2: Add public_ip argument support ([#384](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/384))
* resource/huaweicloud_networking_eip_associate_v2: Add public_ip argument support ([#385](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/385))
* resource/huaweicloud_gaussdb_mysql_instance: Add az mode and configuration_id support ([#396](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/396))
* resource/huaweicloud_gaussdb_cassandra_instance: Add private_ips support ([#406](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/406))
* resource/huaweicloud_cce_cluster_v3: Add kube_proxy_mode support ([#424](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/424))
* resource/huaweicloud_cce_node_v3: Add taints support ([#424](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/424))

BUG FIXES:

* resource/huaweicloud_vpc_eip_v1: Ignore eip unbind error during deleting ([#368](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/368))
* resource/huaweicloud_cce_node_v3: Fix max_pods argument issue ([#369](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/369))

## 1.16.0 (July 03, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_gaussdb_mysql_flavors` ([#354](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/354))
* **New Resource:** `huaweicloud_gaussdb_mysql_instance` ([#350](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/350))
* **New Resource:** `huaweicloud_gaussdb_opengauss_instance` ([#353](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/353))
* **New Resource:** `huaweicloud_geminidb_instance` ([#347](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/347))

ENHANCEMENTS:

* Improvement on dds flavors data source ([#355](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/355))
* Support `port` and `nodes` attributes in dds instance ([#349](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/349))
* Support `ssl` parameter in dds instance resource ([#343](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/343))
* Support `routes` attribute in vpc resource ([#342](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/342))
* Support `status` and `current_instance_number` attributes in as group ([#344](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/344))
* Support `auto_renew` for ecs instance ([#359](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/359))

BUG FIXES:

* resource/huaweicloud_rds_instance_v3: fix document issue ([#348](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/348))

## 1.15.0 (June 10, 2020)

ENHANCEMENTS:

* Mark sensitive for password parameters ([#314](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/314))
* Add tags support for VPC and Subnet resources ([#315](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/315))
* Make `auth_url` optional for provider configuration ([#328](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/328))
* Use `region` as tenant_name if not set ([#330](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/330))
* Add some validations for parameters of provider configuration ([#335](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/335))
* Set external logs according to TF_LOG instead of OS_DEBUG ([#339](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/339))

BUG FIXES:

* resource/huaweicloud_cdn_domain_v1: Fix resource not found issue ([#319](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/319))
* Ignore errors when fetching tags failed in ReadFunc ([#332](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/332))

## 1.14.0 (April 24, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_dds_flavors_v3` ([#305](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/305))
* **New Resource:** `huaweicloud_evs_snapshot` ([#289](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/289))
* **New Resource:** `huaweicloud_cci_network_v1` ([#294](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/294))
* **New Resource:** `huaweicloud_dds_instance_v3` ([#305](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/305))

ENHANCEMENTS:

* resource/huaweicloud_cce_cluster_v3: Add authenticatingProxy.ca support ([#279](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/279))
* resource/huaweicloud_cce_node_v3: Add subnet_id parameter support ([#280](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/280))
* resource/huaweicloud_vpnaas_service_v2: Set admin_state_up default to true ([#293](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/293))
* resource/huaweicloud_compute_instance_v2: Make compute_instance_v2 importable ([#301](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/301))
* resource/huaweicloud_as_group_v1: Add tags support ([#306](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/306))
* resource/huaweicloud_lb_listener_v2: Add http2_enable parameter support ([#307](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/307))
* resource/huaweicloud_vbs_backup_policy_v2: Add week_frequency and rentention_day support ([#309](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/309))

BUG FIXES:

* resource/huaweicloud_nat_snat_rule_v2: Suppress diffs of floating_ip_id ([#274](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/274))
* resource/huaweicloud_fw_rule_v2: Fix removing FW rule assigned to FW policy ([#275](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/275))
* resource/huaweicloud_ecs_instance_v1: Fix DELETED status issue ([#276](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/276))
* resource/huaweicloud_cce_node_v3: Update docs for password parameter ([#282](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/282))
* resource/huaweicloud_nat_snat_rule_v2: Fix attribute type issue ([#291](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/291))
* resource/huaweicloud_obs_bucket: Fix region issue if not specified ([#292](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/292))
* resource/huaweicloud_cce_cluster_v3: Catch client creating exception ([#299](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/299))
* resource/huaweicloud_ecs_instance_v1: Fix PrePaid ECS instance issue ([#304](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/304))

## 1.13.0 (March 10, 2020)

FEATURES:

* **New Resource:** `huaweicloud_lb_whitelist_v2` ([#261](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/261))
* **New Resource:** `huaweicloud_nat_dnat_rule_v2` ([#265](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/265))
* **New Resource:** `huaweicloud_obs_bucket` ([#268](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/268))
* **New Resource:** `huaweicloud_api_gateway_group` ([#270](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/270))
* **New Resource:** `huaweicloud_api_gateway_api` ([#270](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/270))
* **New Resource:** `huaweicloud_fgs_function_v2` ([#271](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/271))

ENHANCEMENTS:

* resource/huaweicloud_identity_user_v3: Add description parameter support ([#266](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/266))
* resource/huaweicloud_s3_bucket: Add tags support ([#267](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/267))
* resource/huaweicloud_cce_node_v3: Add preinstall/postinstall script support ([#269](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/269))

## 1.12.0 (January 14, 2020)

ENHANCEMENTS:

* resource/huaweicloud_compute_volume_attach_v2: Add pci_address attribute support ([#251](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/251))
* resource/huaweicloud_compute_instance_v2: Add support for specifying deh host ([#253](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/253))
* resource/huaweicloud_ecs_instance_v1: Add port_id attribute to nics of ecs_instance ([#258](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/258))
* resource/huaweicloud_ecs_instance_v1: Add op_svc_userid support to ecs_instance ([#259](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/259))

BUG FIXES:

* resource/huaweicloud_as_group_v1: Fix desire/min/max_instance_number argument issue ([#250](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/250))
* resource/huaweicloud_as_group_v1: Fix usage docs issue ([#254](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/254))

## 1.11.0 (December 06, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_vpc_ids_v1` ([#233](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/233))
* **New Data Source:** `huaweicloud_compute_availability_zones_v2` ([#240](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/240))
* **New Data Source:** `huaweicloud_rds_flavors_v3` ([#248](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/248))
* **New Resource:** `huaweicloud_rds_instance_v3` ([#248](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/248))
* **New Resource:** `huaweicloud_rds_parametergroup_v3` ([#248](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/248))

ENHANCEMENTS:

* resource/huaweicloud_as_group_v1: Add lbaas_listeners to as_group_v1 ([#238](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/238))
* resource/huaweicloud_as_configuration_v1: Add kms_id to as_configuration_v1 ([#243](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/243))

BUG FIXES:

* resource/huaweicloud_ecs_instance_v1: Fix ecs instance prepaid issue ([#231](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/231))
* resource/huaweicloud_kms_key_v1: Fix kms client issue ([#234](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/234))

## 1.10.0 (November 13, 2019)

FEATURES:

* **New Resource:** `huaweicloud_cdn_domain_v1` ([#223](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/223))

ENHANCEMENTS:

* resource/huaweicloud_compute_instance_v2: Add volume_attached attribute support ([#214](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/214))
* resource/huaweicloud_cce_cluster_v3: Add eip parameter support ([#219](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/219))

BUG FIXES:

* resource/huaweicloud_compute_volume_attach_v2: Fix example issue for attaching volume ([#221](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/221))
* resource/huaweicloud_compute_instance_v2: Log fault message when build compute instance failed ([#225](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/225))
* resource/huaweicloud_ecs_instance_v1: Fix PrePaid ECS instance issue ([#226](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/226))

## 1.9.0 (September 30, 2019)

FEATURES:

* **New Resource:** `huaweicloud_dns_ptrrecord_v2` ([#200](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/200))
* **New Resource:** `huaweicloud_vpc_bandwidth_v2` ([#203](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/203))
* **New Resource:** `huaweicloud_lb_certificate_v2` ([#211](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/211))
* **New Resource:** `huaweicloud_networking_vip_v2` ([#212](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/212))
* **New Resource:** `huaweicloud_networking_vip_associate_v2` ([#212](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/212))

ENHANCEMENTS:

* resource/huaweicloud_vpc_eip_v1: Add shared bandwidth support ([#208](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/208))

BUG FIXES:

* resource/huaweicloud_ecs_instance_v1: Make ECS instance prePaid auto pay ([#202](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/202))
* Fix ELB resources job issue ([#207](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/207))

## 1.8.0 (August 28, 2019)

FEATURES:

* **New Resource:** `huaweicloud_ecs_instance_v1` ([#179](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/179))
* **New Resource:** `huaweicloud_compute_interface_attach_v2` ([#189](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/189))

ENHANCEMENTS:

* Add detailed error message for 404 ([#183](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/183))
* resource/huaweicloud_cce_node_v3: Add server_id to CCE node ([#185](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/185))
* resource/huaweicloud_cce_cluster_v3: Add certificates to CCE cluster ([#192](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/192))
* resource/huaweicloud_cce_node_v3: Add password support to CCE node ([#193](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/193))
* resource/huaweicloud_cce_cluster_v3: Add multi-az support to CCE cluster ([#194](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/194))

BUG FIXES:

* Fix OBS endpoint issue for new region ([#175](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/175))
* resource/huaweicloud_blockstorage_volume_v2: Add volume extending support ([#176](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/176))
* Update CCE client for new region ([#181](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/181))
* resource/huaweicloud_cce_node_v3: Fix data_volumes type of cce node ([#182](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/182))
* resource/huaweicloud_vpc_subnet_v1: Fix dns_list type issue ([#191](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/191))

## 1.7.0 (July 29, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_networking_port_v2` ([#152](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/152))
* **New Resource:** `huaweicloud_cs_cluster_v1` ([#153](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/153))
* **New Resource:** `huaweicloud_cs_peering_connect_v1` ([#154](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/154))
* **New Resource:** `huaweicloud_vpnaas_service_v2` ([#162](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/162))
* **New Resource:** `huaweicloud_vpnaas_endpoint_group_v2` ([#163](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/163))
* **New Resource:** `huaweicloud_vpnaas_ipsec_policy_v2` ([#164](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/164))
* **New Resource:** `huaweicloud_vpnaas_ike_policy_v2` ([#165](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/165))
* **New Resource:** `huaweicloud_vpnaas_site_connection_v2` ([#166](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/166))
* **New Resource:** `huaweicloud_dli_queue_v1` ([#170](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/170))
* **New Resource:** `huaweicloud_cs_route_v1` ([#171](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/171))

ENHANCEMENTS:

* resource/huaweicloud_networking_floatingip_v2: Add default value for floating_ip pool ([#160](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/160))
Make username/password authentication prior to ak/sk when they both provided ([#167](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/167))

BUG FIXES:

* Replace d.Set("id") with d.SetId to be compatible with terraform 0.12 ([#155](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/155))
* resource/huaweicloud_sfs_file_system_v2: Set availability_zone to Computed ([#156](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/156))
* resource/huaweicloud_compute_instance_v2: Remove personality from compute_instance_v2 ([#169](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/169))

## 1.6.0 (June 13, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_cdm_flavors_v1` ([#128](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/128))
* **New Data Source:** `huaweicloud_dis_partition_v2` ([#134](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/134))
* **New Resource:** `huaweicloud_cdm_cluster_v1` ([#128](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/128))
* **New Resource:** `huaweicloud_ges_graph_v1` ([#131](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/131))
* **New Resource:** `huaweicloud_css_cluster_v1` ([#132](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/132))
* **New Resource:** `huaweicloud_cloudtable_cluster_v2` ([#133](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/133))
* **New Resource:** `huaweicloud_dis_partition_v2` ([#134](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/134))

ENHANCEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

BUG FIXES:

* resource/huaweicloud_identity_role_assignment_v3: Fix role assignment issue ([#136](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/136))
* resource/huaweicloud_cce_node_v3: Fix cce node os option issue ([#145](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/145))
* resource/huaweicloud_vpc_subnet_v1: Fix vpc subnet delete issue ([#148](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/148))

## 1.5.0 (May 17, 2019)

ENHANCEMENTS:

* resource/huaweicloud_cce_cluster_v3: Add authentication mode option support for CCE cluster ([#98](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/98))
* resource/huaweicloud_cce_node_v3: Add os option support for CCE node ([#100](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/100))
* resource/huaweicloud_cce_node_v3: Add private/public IP attributes to CCE node ([#127](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/127))

BUG FIXES:

* resource/huaweicloud_cce_node_v3: Remove Abnormal from CCE node creating target state ([#112](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/112))
* resource/huaweicloud_cce_node_v3: Fix CCE node eip_count issue ([#115](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/115))
* resource/huaweicloud_s3_bucket: Fix OBS bucket domain name ([#124](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/124))
* resource/huaweicloud_cce_cluster_v3: Fix CCE cluster wait state error ([#125](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/125))

## 1.4.0 (March 21, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_identity_role_v3` ([#81](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_project_v3` ([#81](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_role_assignment_v3` ([#81](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_user_v3` ([#81](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_group_v3` ([#81](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_group_membership_v3` ([#81](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_lb_l7policy_v2` ([#82](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/82))
* **New Resource:** `huaweicloud_lb_l7rule_v2` ([#82](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/82))

ENHANCEMENTS:

* provider: Support authorized by token + agency ([#78](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/78))
* resource/huaweicloud_dns_zone_v2: Add multi router support for dns zone ([#80](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/80))
* resource/huaweicloud_networking_port_v2: Add DHCP opts to port resource ([#83](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/83))
* resource/huaweicloud_cce_cluster_v3: Add detailed options for cce cluster `flavor_id` and `container_network_type` ([#89](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/89))

BUG FIXES:

* resource/huaweicloud_dcs_instance_v1: Fix dcs instance update error ([#79](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/79))
* resource/huaweicloud_compute_instance_v2: Fix default security group error ([#86](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/86))
* resource/huaweicloud_dns_recordset_v2: Fix dns records update error ([#87](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/87))

## 1.3.0 (January 08, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_dms_az_v1` ([#41](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/41))
* **New Data Source:** `huaweicloud_dms_product_v1` ([#41](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/41))
* **New Data Source:** `huaweicloud_dms_maintainwindow_v1` ([#41](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/41))
* **New Data Source:** `huaweicloud_vbs_backup_policy_v2` ([#44](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/44))
* **New Data Source:** `huaweicloud_vbs_backup_v2` ([#44](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/44))
* **New Data Source:** `huaweicloud_cce_cluster_v3` ([#19](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/19))
* **New Data Source:** `huaweicloud_cce_node_v3` ([#19](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/19))
* **New Data Source:** `huaweicloud_cts_tracker_v1` ([#46](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/46))
* **New Data Source:** `huaweicloud_csbs_backup_v1` ([#42](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/42))
* **New Data Source:** `huaweicloud_csbs_backup_policy_v1` ([#42](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/42))
* **New Data Source:** `huaweicloud_antiddos_v1` ([#47](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/47))
* **New Data Source:** `huaweicloud_dcs_az_v1` ([#55](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/55))
* **New Data Source:** `huaweicloud_dcs_maintainwindow_v1` ([#55](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/55))
* **New Data Source:** `huaweicloud_dcs_product_v1` ([#55](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/55))
* **New Resource:** `huaweicloud_dms_queue_v1` ([#41](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/41))
* **New Resource:** `huaweicloud_dms_group_v1` ([#41](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/41))
* **New Resource:** `huaweicloud_dms_instance_v1` ([#41](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/41))
* **New Resource:** `huaweicloud_vbs_backup_policy_v2` ([#44](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/44))
* **New Resource:** `huaweicloud_vbs_backup_v2` ([#44](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/44))
* **New Resource:** `huaweicloud_cce_cluster_v3` ([#19](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/19))
* **New Resource:** `huaweicloud_cce_node_v3` ([#19](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/19))
* **New Resource:** `huaweicloud_cts_tracker_v1` ([#46](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/46))
* **New Resource:** `huaweicloud_csbs_backup_v1` ([#42](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/42))
* **New Resource:** `huaweicloud_csbs_backup_policy_v1` ([#42](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/42))
* **New Resource:** `huaweicloud_mrs_cluster_v1` ([#56](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/56))
* **New Resource:** `huaweicloud_mrs_job_v1` ([#56](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/56))
* **New Resource:** `huaweicloud_dcs_instance_v1` ([#55](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/55))
* **New Resource:** `huaweicloud_maas_task_v1` ([#65](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/65))
* **New Resource:** `huaweicloud_networking_floatingip_associate_v2` ([#68](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/68))
* **New Resource:** `huaweicloud_dws_cluster` ([#69](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/69))
* **New Resource:** `huaweicloud_mls_instance` ([#69](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/69))

BUG FIXES:

* `resource/huaweicloud_elb_listener`: Fix certificate_id check ([#45](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/45))
* `resource/huaweicloud_smn_topic_v2`: Fix smn topic update error ([#48](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/48))
* `resource/huaweicloud_kms_key_v1`: Add default value of pending_days ([#62](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/62))
* `all resources`: Expose real error message of BadRequest error ([#63](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/63))
* `resource/huaweicloud_sfs_file_system_v2`: Suppress sfs system metadata ([#64](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/64))

## 1.2.0 (September 21, 2018)

FEATURES:

* **New Data Source:** `huaweicloud_vpc_v1` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_peering_connection_v2` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_route_v2` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_route_ids_v2` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_subnet_v1` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_subnet_ids_v1` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_rts_software_config_v1` ([#20](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/20))
* **New Data Source:** `huaweicloud_images_image_v2` ([#25](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/25))
* **New Resource:** `huaweicloud_vpc_v1` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_peering_connection_v2` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_peering_connection_accepter_v2` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_route_v2` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_subnet_v1` ([#14](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_rts_software_config_v1` ([#20](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/20))
* **New Resource:** `huaweicloud_images_image_v2` ([#25](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/25))
* **New Resource:** `huaweicloud_ces_alarmrule` ([#27](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/27))
* **New Resource:** `huaweicloud_as_configuration_v1` ([#29](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/29))
* **New Resource:** `huaweicloud_as_group_v1` ([#30](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/30))
* **New Resource:** `huaweicloud_as_policy_v1` ([#31](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/31))
* **New Resource:** `huaweicloud_cce_cluster_v3` ([#19](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/19))
* **New Resource:** `huaweicloud_cce_node_v3` ([#19](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/19))

ENHANCEMENTS:

* provider: Add AK/SK authentication support ([#33](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/33))

## 1.1.0 (July 20, 2018)

FEATURES:

* **New Data Source:** `huaweicloud_sfs_file_system_v2` ([#9](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/9))
* **New Data Source:** `huaweicloud_rts_stack_v1` ([#10](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/10))
* **New Data Source:** `huaweicloud_rts_stack_resource_v1` ([#10](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/10))
* **New Resource:** `huaweicloud_iam_agency_v3` ([#7](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/7))
* **New Resource:** `huaweicloud_sfs_file_system_v2` ([#9](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/9))
* **New Resource:** `huaweicloud_rts_stack_v1` ([#10](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/10))
* **New Resource:** `huaweicloud_iam_agency_v3` ([#16](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/16))

ENHANCEMENTS:

* resource/huaweicloud_dns_recordset_v2: Add `PTR` type ([#12](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/12))

BUG FIXES:

* provider: Create only one token ([#5](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/5))

## 1.0.0 (July 13, 2018)

FEATURES:

* **New Data Source:** `huaweicloud_networking_network_v2`
* **New Data Source:** `huaweicloud_networking_subnet_v2`
* **New Data Source:** `huaweicloud_networking_secgroup_v2`
* **New Data Source:** `huaweicloud_s3_bucket_object`
* **New Data Source:** `huaweicloud_kms_key_v1`
* **New Data Source:** `huaweicloud_kms_data_key_v1`
* **New Data Source:** `huaweicloud_rds_flavors_v1`
* **New Resource:** `huaweicloud_blockstorage_volume_v2`
* **New Resource:** `huaweicloud_compute_instance_v2`
* **New Resource:** `huaweicloud_compute_keypair_v2`
* **New Resource:** `huaweicloud_compute_secgroup_v2`
* **New Resource:** `huaweicloud_compute_servergroup_v2`
* **New Resource:** `huaweicloud_compute_floatingip_v2`
* **New Resource:** `huaweicloud_compute_floatingip_associate_v2`
* **New Resource:** `huaweicloud_compute_volume_attach_v2`
* **New Resource:** `huaweicloud_dns_recordset_v2`
* **New Resource:** `huaweicloud_dns_zone_v2`
* **New Resource:** `huaweicloud_fw_firewall_group_v2`
* **New Resource:** `huaweicloud_fw_policy_v2`
* **New Resource:** `huaweicloud_fw_rule_v2`
* **New Resource:** `huaweicloud_kms_key_v1`
* **New Resource:** `huaweicloud_elb_loadbalancer`
* **New Resource:** `huaweicloud_elb_listener`
* **New Resource:** `huaweicloud_elb_healthcheck`
* **New Resource:** `huaweicloud_lb_loadbalancer_v2`
* **New Resource:** `huaweicloud_lb_listener_v2`
* **New Resource:** `huaweicloud_lb_pool_v2`
* **New Resource:** `huaweicloud_lb_member_v2`
* **New Resource:** `huaweicloud_lb_monitor_v2`
* **New Resource:** `huaweicloud_networking_network_v2`
* **New Resource:** `huaweicloud_networking_subnet_v2`
* **New Resource:** `huaweicloud_networking_floatingip_v2`
* **New Resource:** `huaweicloud_networking_port_v2`
* **New Resource:** `huaweicloud_networking_router_v2`
* **New Resource:** `huaweicloud_networking_router_interface_v2`
* **New Resource:** `huaweicloud_networking_router_route_v2`
* **New Resource:** `huaweicloud_networking_secgroup_v2`
* **New Resource:** `huaweicloud_networking_secgroup_rule_v2`
* **New Resource:** `huaweicloud_s3_bucket`
* **New Resource:** `huaweicloud_s3_bucket_policy`
* **New Resource:** `huaweicloud_s3_bucket_object`
* **New Resource:** `huaweicloud_smn_topic_v2`
* **New Resource:** `huaweicloud_smn_subscription_v2`
* **New Resource:** `huaweicloud_rds_instance_v1`
* **New Resource:** `huaweicloud_nat_gateway_v2`
* **New Resource:** `huaweicloud_nat_snat_rule_v2`
