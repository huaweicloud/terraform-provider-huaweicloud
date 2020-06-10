## 1.15.0 (June 10, 2020)

ENHANCEMENTS:

* Mark sensitive for password parameters ([#314](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/314))
* Add tags support for VPC and Subnet resources ([#315](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/315))
* Make `auth_url` optional for provider configuration ([#328](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/328))
* Use `region` as tenant_name if not set ([#330](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/330))
* Add some validations for parameters of provider configuration ([#335](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/335))
* Set external logs according to TF_LOG instead of OS_DEBUG ([#339](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/339))

BUG FIXES:

* resource/huaweicloud_cdn_domain_v1: Fix resource not found issue ([#319](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/319))
* Ignore errors when fetching tags failed in ReadFunc ([#332](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/332))

## 1.14.0 (April 24, 2020)

FEATURES:

* **New Data Source:** `huaweicloud_dds_flavors_v3` ([#305](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/305))
* **New Resource:** `huaweicloud_evs_snapshot` ([#289](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/289))
* **New Resource:** `huaweicloud_cci_network_v1` ([#294](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/294))
* **New Resource:** `huaweicloud_dds_instance_v3` ([#305](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/305))

ENHANCEMENTS:

* resource/huaweicloud_cce_cluster_v3: Add authenticatingProxy.ca support ([#279](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/279))
* resource/huaweicloud_cce_node_v3: Add subnet_id parameter support ([#280](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/280))
* resource/huaweicloud_vpnaas_service_v2: Set admin_state_up default to true ([#293](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/293))
* resource/huaweicloud_compute_instance_v2: Make compute_instance_v2 importable ([#301](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/301))
* resource/huaweicloud_as_group_v1: Add tags support ([#306](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/306))
* resource/huaweicloud_lb_listener_v2: Add http2_enable parameter support ([#307](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/307))
* resource/huaweicloud_vbs_backup_policy_v2: Add week_frequency and rentention_day support ([#309](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/309))

BUG FIXES:

* resource/huaweicloud_nat_snat_rule_v2: Suppress diffs of floating_ip_id ([#274](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/274))
* resource/huaweicloud_fw_rule_v2: Fix removing FW rule assigned to FW policy ([#275](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/275))
* resource/huaweicloud_ecs_instance_v1: Fix DELETED status issue ([#276](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/276))
* resource/huaweicloud_cce_node_v3: Update docs for password parameter ([#282](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/282))
* resource/huaweicloud_nat_snat_rule_v2: Fix attribute type issue ([#291](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/291))
* resource/huaweicloud_obs_bucket: Fix region issue if not specified ([#292](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/292))
* resource/huaweicloud_cce_cluster_v3: Catch client creating exception ([#299](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/299))
* resource/huaweicloud_ecs_instance_v1: Fix PrePaid ECS instance issue ([#304](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/304))

## 1.13.0 (March 10, 2020)

FEATURES:

* **New Resource:** `huaweicloud_lb_whitelist_v2` ([#261](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/261))
* **New Resource:** `huaweicloud_nat_dnat_rule_v2` ([#265](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/265))
* **New Resource:** `huaweicloud_obs_bucket` ([#268](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/268))
* **New Resource:** `huaweicloud_api_gateway_group` ([#270](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/270))
* **New Resource:** `huaweicloud_api_gateway_api` ([#270](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/270))
* **New Resource:** `huaweicloud_fgs_function_v2` ([#271](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/271))

ENHANCEMENTS:

* resource/huaweicloud_identity_user_v3: Add description parameter support ([#266](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/266))
* resource/huaweicloud_s3_bucket: Add tags support ([#267](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/267))
* resource/huaweicloud_cce_node_v3: Add preinstall/postinstall script support ([#269](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/269))

## 1.12.0 (January 14, 2020)

ENHANCEMENTS:

* resource/huaweicloud_compute_volume_attach_v2: Add pci_address attribute support ([#251](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/251))
* resource/huaweicloud_compute_instance_v2: Add support for specifying deh host ([#253](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/253))
* resource/huaweicloud_ecs_instance_v1: Add port_id attribute to nics of ecs_instance ([#258](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/258))
* resource/huaweicloud_ecs_instance_v1: Add op_svc_userid support to ecs_instance ([#259](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/259))

BUG FIXES:

* resource/huaweicloud_as_group_v1: Fix desire/min/max_instance_number argument issue ([#250](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/250))
* resource/huaweicloud_as_group_v1: Fix usage docs issue ([#254](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/254))

## 1.11.0 (December 06, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_vpc_ids_v1` ([#233](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/233))
* **New Data Source:** `huaweicloud_compute_availability_zones_v2` ([#240](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/240))
* **New Data Source:** `huaweicloud_rds_flavors_v3` ([#248](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/248))
* **New Resource:** `huaweicloud_rds_instance_v3` ([#248](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/248))
* **New Resource:** `huaweicloud_rds_parametergroup_v3` ([#248](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/248))

ENHANCEMENTS:

* resource/huaweicloud_as_group_v1: Add lbaas_listeners to as_group_v1 ([#238](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/238))
* resource/huaweicloud_as_configuration_v1: Add kms_id to as_configuration_v1 ([#243](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/243))

BUG FIXES:

* resource/huaweicloud_ecs_instance_v1: Fix ecs instance prepaid issue ([#231](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/231))
* resource/huaweicloud_kms_key_v1: Fix kms client issue ([#234](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/234))

## 1.10.0 (November 13, 2019)

FEATURES:

* **New Resource:** `huaweicloud_cdn_domain_v1` ([#223](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/223))

ENHANCEMENTS:

* resource/huaweicloud_compute_instance_v2: Add volume_attached attribute support ([#214](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/214))
* resource/huaweicloud_cce_cluster_v3: Add eip parameter support ([#219](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/219))

BUG FIXES:

* resource/huaweicloud_compute_volume_attach_v2: Fix example issue for attaching volume ([#221](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/221))
* resource/huaweicloud_compute_instance_v2: Log fault message when build compute instance failed ([#225](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/225))
* resource/huaweicloud_ecs_instance_v1: Fix PrePaid ECS instance issue ([#226](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/226))

## 1.9.0 (September 30, 2019)

FEATURES:

* **New Resource:** `huaweicloud_dns_ptrrecord_v2` ([#200](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/200))
* **New Resource:** `huaweicloud_vpc_bandwidth_v2` ([#203](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/203))
* **New Resource:** `huaweicloud_lb_certificate_v2` ([#211](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/211))
* **New Resource:** `huaweicloud_networking_vip_v2` ([#212](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/212))
* **New Resource:** `huaweicloud_networking_vip_associate_v2` ([#212](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/212))

ENHANCEMENTS:

* resource/huaweicloud_vpc_eip_v1: Add shared bandwidth support ([#208](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/208))

BUG FIXES:

* resource/huaweicloud_ecs_instance_v1: Make ECS instance prePaid auto pay ([#202](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/202))
* Fix ELB resources job issue ([#207](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/207))

## 1.8.0 (August 28, 2019)

FEATURES:

* **New Resource:** `huaweicloud_ecs_instance_v1` ([#179](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/179))
* **New Resource:** `huaweicloud_compute_interface_attach_v2` ([#189](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/189))

ENHANCEMENTS:

* Add detailed error message for 404 ([#183](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/183))
* resource/huaweicloud_cce_node_v3: Add server_id to CCE node ([#185](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/185))
* resource/huaweicloud_cce_cluster_v3: Add certificates to CCE cluster ([#192](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/192))
* resource/huaweicloud_cce_node_v3: Add password support to CCE node ([#193](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/193))
* resource/huaweicloud_cce_cluster_v3: Add multi-az support to CCE cluster ([#194](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/194))

BUG FIXES:

* Fix OBS endpoint issue for new region ([#175](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/175))
* resource/huaweicloud_blockstorage_volume_v2: Add volume extending support ([#176](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/176))
* Update CCE client for new region ([#181](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/181))
* resource/huaweicloud_cce_node_v3: Fix data_volumes type of cce node ([#182](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/182))
* resource/huaweicloud_vpc_subnet_v1: Fix dns_list type issue ([#191](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/191))

## 1.7.0 (July 29, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_networking_port_v2` ([#152](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/152))
* **New Resource:** `huaweicloud_cs_cluster_v1` ([#153](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/153))
* **New Resource:** `huaweicloud_cs_peering_connect_v1` ([#154](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/154))
* **New Resource:** `huaweicloud_vpnaas_service_v2` ([#162](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/162))
* **New Resource:** `huaweicloud_vpnaas_endpoint_group_v2` ([#163](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/163))
* **New Resource:** `huaweicloud_vpnaas_ipsec_policy_v2` ([#164](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/164))
* **New Resource:** `huaweicloud_vpnaas_ike_policy_v2` ([#165](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/165))
* **New Resource:** `huaweicloud_vpnaas_site_connection_v2` ([#166](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/166))
* **New Resource:** `huaweicloud_dli_queue_v1` ([#170](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/170))
* **New Resource:** `huaweicloud_cs_route_v1` ([#171](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/171))

ENHANCEMENTS:

* resource/huaweicloud_networking_floatingip_v2: Add default value for floating_ip pool ([#160](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/160))
Make username/password authentication prior to ak/sk when they both provided ([#167](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/167))

BUG FIXES:

* Replace d.Set("id") with d.SetId to be compatible with terraform 0.12 ([#155](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/155))
* resource/huaweicloud_sfs_file_system_v2: Set availability_zone to Computed ([#156](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/156))
* resource/huaweicloud_compute_instance_v2: Remove personality from compute_instance_v2 ([#169](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/169))

## 1.6.0 (June 13, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_cdm_flavors_v1` ([#128](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/128))
* **New Data Source:** `huaweicloud_dis_partition_v2` ([#134](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/134))
* **New Resource:** `huaweicloud_cdm_cluster_v1` ([#128](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/128))
* **New Resource:** `huaweicloud_ges_graph_v1` ([#131](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/131))
* **New Resource:** `huaweicloud_css_cluster_v1` ([#132](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/132))
* **New Resource:** `huaweicloud_cloudtable_cluster_v2` ([#133](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/133))
* **New Resource:** `huaweicloud_dis_partition_v2` ([#134](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/134))

ENHANCEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

BUG FIXES:

* resource/huaweicloud_identity_role_assignment_v3: Fix role assignment issue ([#136](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/136))
* resource/huaweicloud_cce_node_v3: Fix cce node os option issue ([#145](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/145))
* resource/huaweicloud_vpc_subnet_v1: Fix vpc subnet delete issue ([#148](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/148))

## 1.5.0 (May 17, 2019)

ENHANCEMENTS:

* resource/huaweicloud_cce_cluster_v3: Add authentication mode option support for CCE cluster ([#98](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/98))
* resource/huaweicloud_cce_node_v3: Add os option support for CCE node ([#100](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/100))
* resource/huaweicloud_cce_node_v3: Add private/public IP attributes to CCE node ([#127](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/127))

BUG FIXES:

* resource/huaweicloud_cce_node_v3: Remove Abnormal from CCE node creating target state ([#112](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/112))
* resource/huaweicloud_cce_node_v3: Fix CCE node eip_count issue ([#115](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/115))
* resource/huaweicloud_s3_bucket: Fix OBS bucket domain name ([#124](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/124))
* resource/huaweicloud_cce_cluster_v3: Fix CCE cluster wait state error ([#125](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/125))

## 1.4.0 (March 21, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_identity_role_v3` ([#81](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_project_v3` ([#81](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_role_assignment_v3` ([#81](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_user_v3` ([#81](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_group_v3` ([#81](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_identity_group_membership_v3` ([#81](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/81))
* **New Resource:** `huaweicloud_lb_l7policy_v2` ([#82](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/82))
* **New Resource:** `huaweicloud_lb_l7rule_v2` ([#82](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/82))

ENHANCEMENTS:

* provider: Support authorized by token + agency ([#78](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/78))
* resource/huaweicloud_dns_zone_v2: Add multi router support for dns zone ([#80](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/80))
* resource/huaweicloud_networking_port_v2: Add DHCP opts to port resource ([#83](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/83))
* resource/huaweicloud_cce_cluster_v3: Add detailed options for cce cluster `flavor_id` and `container_network_type` ([#89](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/89))

BUG FIXES:

* resource/huaweicloud_dcs_instance_v1: Fix dcs instance update error ([#79](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/79))
* resource/huaweicloud_compute_instance_v2: Fix default security group error ([#86](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/86))
* resource/huaweicloud_dns_recordset_v2: Fix dns records update error ([#87](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/87))

## 1.3.0 (January 08, 2019)

FEATURES:

* **New Data Source:** `huaweicloud_dms_az_v1` ([#41](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/41))
* **New Data Source:** `huaweicloud_dms_product_v1` ([#41](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/41))
* **New Data Source:** `huaweicloud_dms_maintainwindow_v1` ([#41](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/41))
* **New Data Source:** `huaweicloud_vbs_backup_policy_v2` ([#44](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/44))
* **New Data Source:** `huaweicloud_vbs_backup_v2` ([#44](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/44))
* **New Data Source:** `huaweicloud_cce_cluster_v3` ([#19](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/19))
* **New Data Source:** `huaweicloud_cce_node_v3` ([#19](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/19))
* **New Data Source:** `huaweicloud_cts_tracker_v1` ([#46](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/46))
* **New Data Source:** `huaweicloud_csbs_backup_v1` ([#42](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/42))
* **New Data Source:** `huaweicloud_csbs_backup_policy_v1` ([#42](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/42))
* **New Data Source:** `huaweicloud_antiddos_v1` ([#47](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/47))
* **New Data Source:** `huaweicloud_dcs_az_v1` ([#55](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/55))
* **New Data Source:** `huaweicloud_dcs_maintainwindow_v1` ([#55](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/55))
* **New Data Source:** `huaweicloud_dcs_product_v1` ([#55](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/55))
* **New Resource:** `huaweicloud_dms_queue_v1` ([#41](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/41))
* **New Resource:** `huaweicloud_dms_group_v1` ([#41](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/41))
* **New Resource:** `huaweicloud_dms_instance_v1` ([#41](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/41))
* **New Resource:** `huaweicloud_vbs_backup_policy_v2` ([#44](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/44))
* **New Resource:** `huaweicloud_vbs_backup_v2` ([#44](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/44))
* **New Resource:** `huaweicloud_cce_cluster_v3` ([#19](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/19))
* **New Resource:** `huaweicloud_cce_node_v3` ([#19](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/19))
* **New Resource:** `huaweicloud_cts_tracker_v1` ([#46](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/46))
* **New Resource:** `huaweicloud_csbs_backup_v1` ([#42](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/42))
* **New Resource:** `huaweicloud_csbs_backup_policy_v1` ([#42](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/42))
* **New Resource:** `huaweicloud_mrs_cluster_v1` ([#56](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/56))
* **New Resource:** `huaweicloud_mrs_job_v1` ([#56](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/56))
* **New Resource:** `huaweicloud_dcs_instance_v1` ([#55](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/55))
* **New Resource:** `huaweicloud_maas_task_v1` ([#65](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/65))
* **New Resource:** `huaweicloud_networking_floatingip_associate_v2` ([#68](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/68))
* **New Resource:** `huaweicloud_dws_cluster` ([#69](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/69))
* **New Resource:** `huaweicloud_mls_instance` ([#69](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/69))

BUG FIXES:

* `resource/huaweicloud_elb_listener`: Fix certificate_id check ([#45](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/45))
* `resource/huaweicloud_smn_topic_v2`: Fix smn topic update error ([#48](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/48))
* `resource/huaweicloud_kms_key_v1`: Add default value of pending_days ([#62](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/62))
* `all resources`: Expose real error message of BadRequest error ([#63](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/63))
* `resource/huaweicloud_sfs_file_system_v2`: Suppress sfs system metadata ([#64](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/64))

## 1.2.0 (September 21, 2018)

FEATURES:

* **New Data Source:** `huaweicloud_vpc_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_peering_connection_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_route_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_route_ids_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_subnet_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_subnet_ids_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_rts_software_config_v1` ([#20](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/20))
* **New Data Source:** `huaweicloud_images_image_v2` ([#25](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/25))
* **New Resource:** `huaweicloud_vpc_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_peering_connection_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_peering_connection_accepter_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_route_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_subnet_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_rts_software_config_v1` ([#20](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/20))
* **New Resource:** `huaweicloud_images_image_v2` ([#25](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/25))
* **New Resource:** `huaweicloud_ces_alarmrule` ([#27](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/27))
* **New Resource:** `huaweicloud_as_configuration_v1` ([#29](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/29))
* **New Resource:** `huaweicloud_as_group_v1` ([#30](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/30))
* **New Resource:** `huaweicloud_as_policy_v1` ([#31](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/31))
* **New Resource:** `huaweicloud_cce_cluster_v3` ([#19](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/19))
* **New Resource:** `huaweicloud_cce_node_v3` ([#19](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/19))

ENHANCEMENTS:

* provider: Add AK/SK authentication support ([#33](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/33))

## 1.1.0 (July 20, 2018)

FEATURES:

* **New Data Source:** `huaweicloud_sfs_file_system_v2` ([#9](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/9))
* **New Data Source:** `huaweicloud_rts_stack_v1` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/10))
* **New Data Source:** `huaweicloud_rts_stack_resource_v1` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/10))
* **New Resource:** `huaweicloud_iam_agency_v3` ([#7](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/7))
* **New Resource:** `huaweicloud_sfs_file_system_v2` ([#9](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/9))
* **New Resource:** `huaweicloud_rts_stack_v1` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/10))
* **New Resource:** `huaweicloud_iam_agency_v3` ([#16](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/16))

ENHANCEMENTS:

* resource/huaweicloud_dns_recordset_v2: Add `PTR` type ([#12](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/12))

BUG FIXES:

* provider: Create only one token ([#5](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/5))

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
