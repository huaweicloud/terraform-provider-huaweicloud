# CHANGELOG

## 1.80.3 (November 14, 2025)

* **New Resource Source:**
  + `huaweicloud_aom_event_report` [GH-8237]
  + `huaweicloud_aom_uniagent_batch_install` [GH-8212]
  + `huaweicloud_aom_uniagent_batch_upgrade` [GH-8300]
  + `huaweicloud_apig_instance_ingress_port` [GH-8292]
  + `huaweicloud_cbh_delete_fault_instance` [GH-8301]
  + `huaweicloud_cbh_rollback_instance` [GH-8295]
  + `huaweicloud_cbh_upgrade_instance` [GH-8304]
  + `huaweicloud_dms_kafka_topic_message_batch_delete` [GH-8055]
  + `huaweicloud_dms_kafka_volume_auto_expand_configuration` [GH-8282]
  + `huaweicloud_hss_asset_assign_task` [GH-8071]
  + `huaweicloud_rds_notify_replace_node` [GH-8312]
  + `huaweicloud_swr_enterprise_replication_policy_execute` [GH-8315]
  + `huaweicloud_swr_enterprise_replication_policy_execution_stop` [GH-8315]
  + `huaweicloud_swr_enterprise_replication_policy` [GH-8315]

* **New Data Source:**
  + `huaweicloud_apig_instance_ingress_ports` [GH-8292]
  + `huaweicloud_cbh_instance_om_url` [GH-8307]
  + `huaweicloud_cbh_instances_by_tags` [GH-8287]
  + `huaweicloud_dms_kafka_consumer_group_message_offsets` [GH-8044]
  + `huaweicloud_dms_kafka_volume_auto_expand_configuration` [GH-8282]
  + `huaweicloud_fgs_feature` [GH-8285]
  + `huaweicloud_fgs_service_trusted_agencies` [GH-8255]
  + `huaweicloud_hss_antivirus_handle_history` [GH-8305]
  + `huaweicloud_hss_asset_web_app_service_hosts` [GH-8278]
  + `huaweicloud_hss_asset_web_framework_hosts` [GH-8290]
  + `huaweicloud_hss_asset_web_framework_statistics` [GH-8267]
  + `huaweicloud_hss_baseline_scan_status` [GH-8286]
  + `huaweicloud_hss_container_clusters_policy_template` [GH-8280]
  + `huaweicloud_hss_container_clusters_policy_templates` [GH-8280]
  + `huaweicloud_hss_container_kubernetes_clusters_risks` [GH-8296]
  + `huaweicloud_hss_container_network_cluster` [GH-8303]
  + `huaweicloud_hss_kubernetes_cronjobs` [GH-8306]
  + `huaweicloud_hss_kubernetes_daemonsets` [GH-8317]
  + `huaweicloud_hss_kubernetes_jobs` [GH-8319]
  + `huaweicloud_rds_configurable_distributor_instances` [GH-8274]
  + `huaweicloud_rds_configurable_subscriber_instances` [GH-8283]
  + `huaweicloud_rds_distribution` [GH-8274]
  + `huaweicloud_rds_publications` [GH-8291]
  + `huaweicloud_rds_top_sqls` [GH-8302]
  + `huaweicloud_swr_enterprise_replication_policies` [GH-8289]

## 1.80.2 (November 10, 2025)

* **New Resource Source:**
  + `huaweicloud_aom_alarm_inhibit_rule` [GH-8232]
  + `huaweicloud_cciv2_observability_configuration` [GH-8261]
  + `huaweicloud_cpcs_app_download_access_key` [GH-8231]
  + `huaweicloud_cpcs_cluster_authorize_access_key` [GH-8214]
  + `huaweicloud_cpcs_instance_status_action` [GH-8239]
  + `huaweicloud_csms_download_secret_backup` [GH-8256]
  + `huaweicloud_csms_restore_secret` [GH-8265]
  + `huaweicloud_csms_secret_rotate` [GH-8251]
  + `huaweicloud_cts_configuration` [GH-8223]
  + `huaweicloud_geminidb_instance` [GH-8253]
  + `huaweicloud_hss_cluster_protect_switch_mode` [GH-8257]
  + `huaweicloud_kms_dedicated_keystore_action` [GH-8260]
  + `huaweicloud_kms_key_replicate` [GH-8224]
  + `huaweicloud_kms_retire_grant` [GH-8220]
  + `huaweicloud_kps_export_private_key` [GH-8243]
  + `huaweicloud_swr_enterprise_instance_registry` [GH-8221]

* **New Data Source:**
  + `huaweicloud_aad_user_quotas` [GH-8211]
  + `huaweicloud_aadv2_instances` [GH-8210]
  + `huaweicloud_compute_attachable_nics` [GH-8222]
  + `huaweicloud_compute_flavor_capacity` [GH-8219]
  + `huaweicloud_cpcs_cluster_ports` [GH-8236]
  + `huaweicloud_cpcs_cluster_url` [GH-8236]
  + `huaweicloud_cpcs_instances` [GH-8198]
  + `huaweicloud_cpcs_vm_monitor` [GH-8226]
  + `huaweicloud_csms_notification_records` [GH-8213]
  + `huaweicloud_fgs_runtime_types` [GH-8234]
  + `huaweicloud_fgs_trigger_types` [GH-8218]
  + `huaweicloud_hss_asset_web_app_service_statistics` [GH-8249]
  + `huaweicloud_hss_asset_website_statistics` [GH-8254]
  + `huaweicloud_hss_cluster_protect_default_policies` [GH-8233]
  + `huaweicloud_hss_cluster_protect_info` [GH-8225]
  + `huaweicloud_hss_cluster_protect_overview` [GH-8228]
  + `huaweicloud_hss_cluster_protect_protection_items` [GH-8248]
  + `huaweicloud_hss_image_asset_statistics` [GH-8269]
  + `huaweicloud_hss_rasp_events` [GH-8263]
  + `huaweicloud_kms_dedicated_keystores` [GH-8230]
  + `huaweicloud_kms_retirable_grants` [GH-8216]
  + `huaweicloud_rds_instance_no_index_tables` [GH-8268]
  + `huaweicloud_rms_policy_states_statistics` [GH-8227]
  + `huaweicloud_swr_enterprise_instance_registries` [GH-8235]
  + `huaweicloud_workspace_app_publishable_applications` [GH-8229]

## 1.80.1 (October 31, 2025)

* **New Resource Source:**
  + `huaweicloud_cpcs_app_cluster_association` [GH-8181]
  + `huaweicloud_identity_provider_protocol` [GH-8200]
  + `huaweicloud_identity_token_with_id_token` [GH-8200]
  + `huaweicloud_identity_unscoped_token_saml` [GH-8200]
  + `huaweicloud_identity_unscoped_token_with_id_token` [GH-8200]
  + `huaweicloud_identity_user_info` [GH-8200]
  + `huaweicloud_identity_user_password` [GH-8200]
  + `huaweicloud_identitycenter_application_assignment`
  + `huaweicloud_identitycenter_application_certificate`
  + `huaweicloud_identitycenter_bearer_token`
  + `huaweicloud_identitycenter_identity_provider_certificate`
  + `huaweicloud_identitycenter_identity_provider`
  + `huaweicloud_identitycenter_instance` [GH-6047]
  + `huaweicloud_identitycenter_registered_region`
  + `huaweicloud_identitycenter_service_provider_certificate`
  + `huaweicloud_identitycenter_tenant`
  + `huaweicloud_identityv5_login_profile` [GH-8200]
  + `huaweicloud_identityv5_user` [GH-8200]
  + `huaweicloud_rgc_account_enroll` [GH-8208]
  + `huaweicloud_rgc_best_practice` [GH-8208]
  + `huaweicloud_rgc_control` [GH-8208]
  + `huaweicloud_rgc_landing_zone` [GH-8208]
  + `huaweicloud_rgc_organizational_unit_register` [GH-8208]
  + `huaweicloud_rgc_organizational_unit` [GH-8208]
  + `huaweicloud_rgc_template` [GH-8208]

* **New Data Source:**
  + `huaweicloud_elb_loadbalancers_by_tags` [GH-8196]
  + `huaweicloud_hss_antivirus_result` [GH-8178]
  + `huaweicloud_hss_rasp_statistics` [GH-8190]
  + `huaweicloud_identity_access_key` [GH-8200]
  + `huaweicloud_identity_auth_projects` [GH-8200]
  + `huaweicloud_identity_catalog` [GH-8200]
  + `huaweicloud_identity_check_agency_role_assignment` [GH-8200]
  + `huaweicloud_identity_check_group_membership` [GH-8200]
  + `huaweicloud_identity_check_group_role_assignment` [GH-8200]
  + `huaweicloud_identity_domain_quota` [GH-8200]
  + `huaweicloud_identity_endpoints` [GH-8200]
  + `huaweicloud_identity_federation_domains` [GH-8200]
  + `huaweicloud_identity_federation_projects` [GH-8200]
  + `huaweicloud_identity_keystone_metadata_file` [GH-8200]
  + `huaweicloud_identity_login_protects` [GH-8200]
  + `huaweicloud_identity_project_quota` [GH-8200]
  + `huaweicloud_identity_provider_protocols` [GH-8200]
  + `huaweicloud_identity_regions` [GH-8200]
  + `huaweicloud_identity_role_assignments` [GH-8200]
  + `huaweicloud_identity_security_compliance` [GH-8200]
  + `huaweicloud_identity_services` [GH-8200]
  + `huaweicloud_identity_user_projects` [GH-8200]
  + `huaweicloud_identity_user_token_info` [GH-8200]
  + `huaweicloud_identitycenter_account_assignments`
  + `huaweicloud_identitycenter_batch_query_groups`
  + `huaweicloud_identitycenter_batch_query_users`
  + `huaweicloud_identitycenter_service_provider_configuration`
  + `huaweicloud_nat_gateway_specs` [GH-8201]
  + `huaweicloud_nat_private_gateway_specs` [GH-8201]
  + `huaweicloud_nat_private_transit_ips_by_tags` [GH-8201]
  + `huaweicloud_organizations_delegated_services` [GH-8201]
  + `huaweicloud_rgc_blueprint` [GH-8208]
  + `huaweicloud_rgc_landing_zone_available_updates` [GH-8208]
  + `huaweicloud_rgc_landing_zone_configuration` [GH-8208]
  + `huaweicloud_rgc_landing_zone_identity_center` [GH-8208]
  + `huaweicloud_rgc_operation` [GH-8208]
  + `huaweicloud_rgc_organizational_unit_accounts` [GH-8208]
  + `huaweicloud_rgc_organizational_units` [GH-8208]
  + `huaweicloud_rgc_pre_launch_check` [GH-8208]
  + `huaweicloud_tms_resource_tags`
  + `huaweicloud_vpc_network_interface_tags`
  + `huaweicloud_vpc_network_interfaces_by_tags`
  + `huaweicloud_vpc_subnet_tags`
  + `huaweicloud_vpcep_resources_by_tags`

## 1.80.0 (October 30, 2025)

* **New Resource Source:**
  + `huaweicloud_coc_cloud_vendor_account` [GH-8158]
  + `huaweicloud_coc_cloud_vendor_user_resources_sync` [GH-8163]
  + `huaweicloud_csms_agency` [GH-8139]
  + `huaweicloud_dms_kafka_consumer_group_topic_batch_delete` [GH-8047]
  + `huaweicloud_dns_zone_authorization_verify` [GH-8183]
  + `huaweicloud_dns_zone_authorization` [GH-8183]
  + `huaweicloud_evs_recycle_bin_volume_revert` [GH-8119]
  + `huaweicloud_hss_antivirus_create_pay_per_scan_task` [GH-8159]
  + `huaweicloud_hss_antivirus_pay_per_scan_switch_status` [GH-8071]
  + `huaweicloud_hss_modify_webtamper_protection_policy` [GH-8137]
  + `huaweicloud_kms_ec_datakey_pair` [GH-8165]
  + `huaweicloud_kms_generate_mac` [GH-8174]
  + `huaweicloud_kms_key_update_primary_region` [GH-8174]
  + `huaweicloud_kms_rsa_datakey_pair` [GH-8174]
  + `huaweicloud_kms_verify_mac` [GH-8174]
  + `huaweicloud_kps_batch_export_private_key` [GH-8149]
  + `huaweicloud_kps_batch_import_keypair` [GH-8155]
  + `huaweicloud_rds_sql_statistics_view_reset` [GH-8152]
  + `huaweicloud_swr_enterprise_domain_name` [GH-8122]

* **New Data Source:**
  + `huaweicloud_coc_cloud_vendor_accounts` [GH-8158]
  + `huaweicloud_compute_appendable_volume_quota` [GH-8143]
  + `huaweicloud_cpcs_app_access_keys` [GH-8173]
  + `huaweicloud_cpcs_associations` [GH-8180]
  + `huaweicloud_cpcs_availability_zones` [GH-8173]
  + `huaweicloud_cpcs_cluster_access_keys` [GH-8173]
  + `huaweicloud_cpcs_clusters` [GH-8172]
  + `huaweicloud_cpcs_images` [GH-8173]
  + `huaweicloud_cpcs_resource_infos` [GH-8180]
  + `huaweicloud_csms_agencies` [GH-8142]
  + `huaweicloud_csms_function_templates` [GH-8140]
  + `huaweicloud_dms_kafka_consumer_group_members` [GH-8032]
  + `huaweicloud_dns_zone_nameservers` [GH-8167]
  + `huaweicloud_hss_antivirus_pay_per_scan_hosts` [GH-8175]
  + `huaweicloud_hss_asset_app_change_history` [GH-8123]
  + `huaweicloud_hss_asset_website_hosts` [GH-8106]
  + `huaweicloud_hss_cluster_protect_policies` [GH-8151]
  + `huaweicloud_hss_container_cluster_statistics` [GH-8134]
  + `huaweicloud_hss_container_iac_file_risks` [GH-8147]
  + `huaweicloud_hss_container_network_info` [GH-8169]
  + `huaweicloud_rds_sql_statistics` [GH-8152]
  + `huaweicloud_swr_enterprise_image_signature_policy_execution_record_sub_tasks` [GH-8164]
  + `huaweicloud_swr_enterprise_image_signature_policy_execution_record_tasks` [GH-8157]
  + `huaweicloud_swr_enterprise_instance_audit_logs` [GH-8171]
  + `huaweicloud_swr_enterprise_namespace_repositories` [GH-8170]
  + `huaweicloud_swr_enterprise_repositories` [GH-8166]

## 1.79.2 (October 22, 2025)

* **New Resource Source:**
  + `huaweicloud_cdn_certificate_associate_domains` [GH-8061]
  + `huaweicloud_cdn_domain_batch_copy` [GH-8076]
  + `huaweicloud_compute_instance_redeploy` [GH-8100]
  + `huaweicloud_compute_recycle_bin_server_delete` [GH-8080]
  + `huaweicloud_compute_recycle_bin_server_recover` [GH-8080]
  + `huaweicloud_compute_scheduled_event_accept` [GH-8103]
  + `huaweicloud_compute_scheduled_event_update` [GH-8108]
  + `huaweicloud_cpcs_app_access_key` [GH-8124]
  + `huaweicloud_cpcs_app` [GH-8115]
  + `huaweicloud_evs_recycle_bin_policy` [GH-8116]
  + `huaweicloud_evs_recycle_bin_volume_delete` [GH-8120]
  + `huaweicloud_hss_app_whitelist_policy_process` [GH-8105]
  + `huaweicloud_hss_container_network_policy_sync` [GH-8099]
  + `huaweicloud_swr_enterprise_immutable_tag_rule` [GH-8107]
  + `huaweicloud_swr_enterprise_retention_policy_execute` [GH-8083]
  + `huaweicloud_swr_enterprise_retention_policy` [GH-8091]

* **New Data Source:**
  + `huaweicloud_cdn_top_referrer_statistics` [GH-8056]
  + `huaweicloud_cdn_top_url_statistics` [GH-8056]
  + `huaweicloud_compute_auto_launch_group_instances` [GH-8117]
  + `huaweicloud_compute_auto_launch_groups` [GH-8117]
  + `huaweicloud_compute_availability_zones` [GH-8117]
  + `huaweicloud_compute_recycle_bin_servers` [GH-8117]
  + `huaweicloud_compute_resize_flavors` [GH-8117]
  + `huaweicloud_compute_scheduled_events` [GH-8117]
  + `huaweicloud_cpcs_apps` [GH-8115]
  + `huaweicloud_evs_recycle_bin_policy` [GH-8116]
  + `huaweicloud_evs_recycle_bin_volume_detail` [GH-8118]
  + `huaweicloud_evs_recycle_bin_volumes_detail` [GH-8116]
  + `huaweicloud_hss_antivirus_pay_per_scan_free_quotas` [GH-8085]
  + `huaweicloud_hss_app_whitelist_associate_hosts` [GH-8089]
  + `huaweicloud_hss_app_whitelist_optional_hosts` [GH-8089]
  + `huaweicloud_hss_app_whitelist_policies` [GH-8073]
  + `huaweicloud_hss_app_whitelist_policy_process_extend` [GH-8105]
  + `huaweicloud_hss_asset_kernel_module_hosts` [GH-8064]
  + `huaweicloud_hss_asset_kernel_module_statistics` [GH-8067]
  + `huaweicloud_hss_asset_overview_status_container_protection` [GH-8071]
  + `huaweicloud_hss_asset_process_detail` [GH-8095]
  + `huaweicloud_hss_asset_process_statistics` [GH-8082]
  + `huaweicloud_hss_baseline_white_lists` [GH-8071]
  + `huaweicloud_hss_container_iac_files` [GH-8071]
  + `huaweicloud_hss_files_change_hosts` [GH-8101]
  + `huaweicloud_hss_kubernetes_pod_detail` [GH-8096]
  + `huaweicloud_hss_kubernetes_services` [GH-8094]
  + `huaweicloud_hss_webtamper_host_management_hosts` [GH-8084]
  + `huaweicloud_swr_enterprise_domain_names` [GH-8122]
  + `huaweicloud_swr_enterprise_feature_gates` [GH-8112]
  + `huaweicloud_swr_enterprise_immutable_tag_rules` [GH-8107]
  + `huaweicloud_swr_enterprise_retention_policies` [GH-8091]
  + `huaweicloud_swr_enterprise_retention_policy_execution_records` [GH-8091]

## 1.79.1 (October 14, 2025)

* **New Resource Source:**
  + `huaweicloud_compute_password_delete` [GH-8011]
  + `huaweicloud_dms_kafka_topic_quota` [GH-8021]
  + `huaweicloud_hss_modify_webtamper_rasp_path` [GH-8049]
  + `huaweicloud_hss_setting_two_factor_login_config` [GH-8014]
  + `huaweicloud_hss_vulnerability_history_export_task` [GH-8063]
  + `huaweicloud_hss_vulnerability_task_user_trace` [GH-8010]

* **New Data Source:**
  + `huaweicloud_cdn_domain_tags` [GH-8027]
  + `huaweicloud_cdn_ip_information` [GH-8020]
  + `huaweicloud_cdn_quotas` [GH-8036]
  + `huaweicloud_coc_incidents` [GH-8026]
  + `huaweicloud_compute_password` [GH-8009]
  + `huaweicloud_compute_supply_recommendations` [GH-8058]
  + `huaweicloud_dms_kafka_consumer_group_topics` [GH-8030]
  + `huaweicloud_dms_kafka_instance_coordinators` [GH-8024]
  + `huaweicloud_dms_kafka_maintainwindows` [GH-8021]
  + `huaweicloud_dms_kafka_topic_broker_disk_usages` [GH-8021]
  + `huaweicloud_dms_kafka_topic_quotas` [GH-8021]
  + `huaweicloud_hss_app_events` [GH-8046]
  + `huaweicloud_hss_asset_assign_task` [GH-8052]
  + `huaweicloud_hss_asset_overview_status_agent` [GH-8050]
  + `huaweicloud_hss_asset_overview_status_host_protection` [GH-8019]
  + `huaweicloud_hss_asset_overview_status_os` [GH-8048]
  + `huaweicloud_hss_baseline_check_rule_hab` [GH-8022]
  + `huaweicloud_hss_baseline_security_checks_directories` [GH-8034]
  + `huaweicloud_hss_change_files` [GH-8039]
  + `huaweicloud_hss_cicd_configurations` [GH-8019]
  + `huaweicloud_hss_common_task_statistics` [GH-8035]
  + `huaweicloud_hss_common_tasks` [GH-8019]
  + `huaweicloud_hss_configs` [GH-8045]
  + `huaweicloud_hss_container_cluster_risk_affected_resources` [GH-8019]
  + `huaweicloud_hss_container_cluster_risks` [GH-8019]
  + `huaweicloud_hss_container_network_policies` [GH-8038]
  + `huaweicloud_hss_files_statistic` [GH-8028]
  + `huaweicloud_hss_page_notices` [GH-8019]
  + `huaweicloud_hss_setting_login_white_ips` [GH-8019]
  + `huaweicloud_hss_setting_two_factor_login_hosts` [GH-8019]
  + `huaweicloud_hss_webtamper_policy` [GH-8059]
  + `huaweicloud_hss_webtamper_rasp_path` [GH-8042]
  + `huaweicloud_swr_enterprise_sub_resources_filter` [GH-8013]
  + `huaweicloud_waf_all_web_antitamper_rules` [GH-7989]

## 1.79.0 (September 29, 2025)

* **New Resource Source:**
  + `huaweicloud_apig_channel_member_batch_action` [GH-7977]
  + `huaweicloud_apig_global_certificate_batch_domains_associate` [GH-7923]
  + `huaweicloud_ces_resource_group_alarm_template_async_associate` [GH-7978]
  + `huaweicloud_coc_alarm_action` [GH-7938]
  + `huaweicloud_coc_alarm_clear` [GH-7951]
  + `huaweicloud_coc_alarm_linked_incident` [GH-7941]
  + `huaweicloud_coc_ticket_action` [GH-7921]
  + `huaweicloud_compute_kernel_dump_trigger` [GH-7954]
  + `huaweicloud_compute_os_change` [GH-7974]
  + `huaweicloud_compute_os_reinstall` [GH-7969]
  + `huaweicloud_compute_recycle_policy` [GH-7954]
  + `huaweicloud_dms_kafka_instance_batch_action` [GH-7961]
  + `huaweicloud_dms_kafka_instance_rebalance_log` [GH-7966]
  + `huaweicloud_dms_kafka_smart_connector_validate` [GH-7949]
  + `huaweicloud_dms_kafka_user_password_reset` [GH-7961]
  + `huaweicloud_hss_cicd_configuration` [GH-7992]
  + `huaweicloud_modelartsv2_node_batch_reboot` [GH-7962]
  + `huaweicloud_sfs_turbo_ldap_config` [GH-7940]
  + `huaweicloud_swr_enterprise_image_signature_policy_execute` [GH-7976]
  + `huaweicloud_swr_enterprise_image_signature_policy` [GH-7987]
  + `huaweicloud_swr_enterprise_job_delete` [GH-7950]
  + `huaweicloud_waf_alarm_notification` [GH-6866]
  + `huaweicloud_waf_cc_protection_rule_batch_delete` [GH-7956]
  + `huaweicloud_waf_domain_route_update` [GH-7948]
  + `huaweicloud_waf_geo_ip_rule_batch_update` [GH-7980]
  + `huaweicloud_waf_ip_intelligence_rule` [GH-7933]
  + `huaweicloud_waf_policies_batch_delete` [GH-7944]
  + `huaweicloud_waf_policy_copy` [GH-7986]

* **New Data Source:**
  + `huaweicloud_aom_event_statistic` [GH-7971]
  + `huaweicloud_aom_events` [GH-7971]
  + `huaweicloud_apig_application_authorize_statistic` [GH-7882]
  + `huaweicloud_apig_channel_members` [GH-7968]
  + `huaweicloud_apig_instance_associated_certificates` [GH-7985]
  + `huaweicloud_apig_instance_restriction` [GH-7990]
  + `huaweicloud_ces_host_configurations` [GH-7965]
  + `huaweicloud_coc_alarm_action_histories` [GH-7938]
  + `huaweicloud_compute_flavor_sales_policies` [GH-7952]
  + `huaweicloud_compute_quotas` [GH-7993]
  + `huaweicloud_compute_volume_attachments` [GH-7997]
  + `huaweicloud_dms_kafka_tags` [GH-7955]
  + `huaweicloud_sfs_turbo_quotas` [GH-7927]
  + `huaweicloud_swr_enterprise_image_signature_policies` [GH-7964]
  + `huaweicloud_swr_enterprise_image_signature_policy_execution_records` [GH-7987]
  + `huaweicloud_swr_enterprise_instance_tags` [GH-7975]
  + `huaweicloud_swr_enterprise_jobs` [GH-7939]
  + `huaweicloud_swr_enterprise_namespace_tags` [GH-7979]
  + `huaweicloud_swr_enterprise_resources_filter` [GH-7995]
  + `huaweicloud_swr_enterprise_triggers` [GH-7937]
  + `huaweicloud_waf_all_antileakage_rules` [GH-7981]
  + `huaweicloud_waf_all_data_masking_rules` [GH-7947]
  + `huaweicloud_waf_all_geo_ip_policy_rules` [GH-7967]
  + `huaweicloud_waf_all_global_whitelist_rules` [GH-7972]
  + `huaweicloud_waf_all_ip_reputation_policy_rules` [GH-7982]
  + `huaweicloud_waf_all_policy_cc_rules` [GH-7970]
  + `huaweicloud_waf_all_precise_protection_rules` [GH-7963]
  + `huaweicloud_waf_all_whiteblackip_rules` [GH-7943]
  + `huaweicloud_waf_tag_ip_reputation_map` [GH-7909]

## 1.78.5 (September 22, 2025)

* **New Resource Source:**
  + `huaweicloud_apig_api_batch_plugins_associate` [GH-7873]
  + `huaweicloud_apig_certificate_batch_domains_associate` [GH-7873]
  + `huaweicloud_apig_channel_member` [GH-7808]
  + `huaweicloud_apig_domain_certificate_associate` [GH-7839]
  + `huaweicloud_coc_change_delete` [GH-7917]
  + `huaweicloud_coc_change_update` [GH-7916]
  + `huaweicloud_coc_incident_action` [GH-7900]
  + `huaweicloud_coc_issue` [GH-7890]
  + `huaweicloud_coc_ticket_add` [GH-7914]
  + `huaweicloud_compute_template` [GH-7908]
  + `huaweicloud_swr_enterprise_long_term_credential` [GH-7905]
  + `huaweicloud_swr_enterprise_private_network_access_control` [GH-7905]
  + `huaweicloud_swr_enterprise_temporary_credential` [GH-7899]
  + `huaweicloud_swr_enterprise_trigger` [GH-7925]

* **New Data Source:**
  + `huaweicloud_apig_application_associated_quota` [GH-7889]
  + `huaweicloud_apig_certificate_associated_domains` [GH-7896]
  + `huaweicloud_apig_channel_member_groups` [GH-7808]
  + `huaweicloud_cbr_feature` [GH-7894]
  + `huaweicloud_cbr_features` [GH-7894]
  + `huaweicloud_coc_change_sub_tickets` [GH-7920]
  + `huaweicloud_coc_incident_action_histories` [GH-7900]
  + `huaweicloud_coc_scheduled_task_histories` [GH-7868]
  + `huaweicloud_coc_ticket_operation_histories` [GH-7861]
  + `huaweicloud_compute_template_versions` [GH-7911]
  + `huaweicloud_compute_templates` [GH-7908]
  + `huaweicloud_swr_enterprise_long_term_credentials` [GH-7905]
  + `huaweicloud_swr_enterprise_private_network_access_controls` [GH-7905]
  + `huaweicloud_waf_alarm_optional_event_types` [GH-7907]
  + `huaweicloud_waf_overviews_attack_ip` [GH-7897]
  + `huaweicloud_waf_overviews_attack_url` [GH-7895]
  + `huaweicloud_waf_rules_threat_intelligence` [GH-7880]
  + `huaweicloud_waf_tag_antileakage_map` [GH-7903]

## 1.78.4 (September 16, 2025)

* **New Resource Source:**
  + `huaweicloud_apig_channel_member_group` [GH-7801]
  + `huaweicloud_apig_plugin_batch_apis_associate` [GH-7870]
  + `huaweicloud_coc_diagnosis_task_retry` [GH-7824]
  + `huaweicloud_coc_public_script_execute` [GH-7741]
  + `huaweicloud_identitycenter_application_instance` [GH-7827]
  + `huaweicloud_secmaster_asset` [GH-7815]
  + `huaweicloud_swr_enterprise_namespace` [GH-7867]
  + `huaweicloud_vpn_gateway_job_delete` [GH-7720]
  + `huaweicloud_vpn_p2c_gateway_job_delete` [GH-7858]
  + `huaweicloud_vpn_p2c_gateway_upgrade` [GH-7858]

* **New Data Source:**
  + `huaweicloud_apig_api_associable_plugins` [GH-7838]
  + `huaweicloud_apig_plugin_associable_apis` [GH-7811]
  + `huaweicloud_coc_diagnosis_task_node_detail` [GH-7834]
  + `huaweicloud_coc_diagnosis_task_summary` [GH-7832]
  + `huaweicloud_coc_diagnosis_tasks` [GH-7828]
  + `huaweicloud_coc_scheduled_tasks` [GH-7853]
  + `huaweicloud_identitycenter_application_templates` [GH-7827]
  + `huaweicloud_identitycenter_catalog_applications` [GH-7829]
  + `huaweicloud_kms_key_regions` [GH-7810]
  + `huaweicloud_nat_gateways_by_tags` [GH-7822]
  + `huaweicloud_nat_private_gateways_by_tags` [GH-7822]
  + `huaweicloud_rgc_home_region` [GH-7813]
  + `huaweicloud_swr_enterprise_instances` [GH-7869]
  + `huaweicloud_swr_enterprise_namespaces` [GH-7867]
  + `huaweicloud_vpcep_tags` [GH-7833]
  + `huaweicloud_vpn_gateway_jobs` [GH-7720]
  + `huaweicloud_vpn_p2c_gateway_jobs` [GH-7858]

## 1.78.3 (September 11, 2025)

* **New Resource:**
  + `huaweicloud_apig_api_version_unpublish` [GH-7768]
  + `huaweicloud_cciv2_pool_binding` [GH-7794]
  + `huaweicloud_coc_diagnosis_task_cancel` [GH-7823]
  + `huaweicloud_coc_diagnosis_task` [GH-7791]
  + `huaweicloud_coc_scheduled_task` [GH-7750]
  + `huaweicloud_kms_cancel_key_deletion` [GH-7753]
  + `huaweicloud_kms_verify_sign` [GH-7765]
  + `huaweicloud_secmaster_collector_channel_group` [GH-7789]
  + `huaweicloud_swr_enterprise_instance` [GH-7816]
  + `huaweicloud_swr_temporary_login_command` [GH-7788]

* **New Data Source:**
  + `huaweicloud_apig_api_history_versions` [GH-7771]
  + `huaweicloud_apig_instance_tags` [GH-7756]
  + `huaweicloud_cc_bandwidth_package_tags` [GH-7802]
  + `huaweicloud_cc_central_networks_by_tags` [GH-7796]
  + `huaweicloud_cc_global_connection_bandwidth_tags` [GH-7799]
  + `huaweicloud_cc_site_network_quotas` [GH-7803]
  + `huaweicloud_coc_document_execution_step_instances` [GH-7779]
  + `huaweicloud_coc_incident_tasks` [GH-7821]
  + `huaweicloud_coc_script_tags` [GH-7779]
  + `huaweicloud_dms_rocketmq_brokers` [GH-7775]
  + `huaweicloud_lts_context_logs` [GH-7748]
  + `huaweicloud_secmaster_collector_channel_groups` [GH-7789]
  + `huaweicloud_secmaster_collector_channel_instances` [GH-7782]
  + `huaweicloud_secmaster_configuration_dictionaries` [GH-7762]
  + `huaweicloud_secmaster_playbook_instance` [GH-7680]
  + `huaweicloud_swr_image_auto_sync_jobs` [GH-7797]
  + `huaweicloud_swrv3_repositories` [GH-7773]
  + `huaweicloud_swrv3_shared_repositories` [GH-7777]
  + `huaweicloud_vpc_subnet_cidr_reservations` [GH-7769]

## 1.78.2 (September 8, 2025)

BUG FIXES:

* data/huaweicloud_dcs_flavors: fix configurations set issue [GH-7774]

## 1.78.1 (September 6, 2025)

* **New Resource:**
  + `huaweicloud_apig_api_action` [GH-7710]
  + `huaweicloud_apig_api_batch_action` [GH-7723]
  + `huaweicloud_apig_api_debug` [GH-7746]
  + `huaweicloud_coc_document_execution_operation` [GH-7737]
  + `huaweicloud_coc_enterprise_project_collection` [GH-7732]
  + `huaweicloud_coc_group_resource_relation` [GH-7732]
  + `huaweicloud_coc_other_resource_uniagent_sync` [GH-7685]
  + `huaweicloud_coc_resource_uniagent_sync` [GH-7675]
  + `huaweicloud_coc_script_approval` [GH-7654]
  + `huaweicloud_codearts_pipeline_plugin_version` [GH-7668]
  + `huaweicloud_dms_rocketmq_instance_diagnosis` [GH-7700]
  + `huaweicloud_kms_sign` [GH-7727]
  + `huaweicloud_secmaster_delete_policies` [GH-7709]
  + `huaweicloud_secmaster_workflow_version` [GH-7680]
  + `huaweicloud_vpn_connection_reset` [GH-7714]
  + `huaweicloud_vpn_gateway_upgrade` [GH-7697]

* **New Data Source:**
  + `huaweicloud_apig_instance_api_tags` [GH-7754]
  + `huaweicloud_apig_quotas` [GH-7751]
  + `huaweicloud_cc_across_area_bandwidth_package_flavors` [GH-7745]
  + `huaweicloud_cc_across_regions_bandwidth_package_flavors` [GH-7744]
  + `huaweicloud_cc_bandwidth_package_classes` [GH-7725]
  + `huaweicloud_cc_bandwidth_package_lines` [GH-7729]
  + `huaweicloud_cc_bandwidth_package_sites` [GH-7731]
  + `huaweicloud_cc_cloud_connection_capabilities` [GH-7722]
  + `huaweicloud_cc_cloud_connection_quotas` [GH-7759]
  + `huaweicloud_cc_global_connection_bandwidth_configs` [GH-7742]
  + `huaweicloud_cc_site_network_capabilities` [GH-7716]
  + `huaweicloud_cc_supported_areas` [GH-7763]
  + `huaweicloud_cc_supported_regions` [GH-7763]
  + `huaweicloud_coc_application_capacities` [GH-7732]
  + `huaweicloud_coc_application_capacity_orders` [GH-7732]
  + `huaweicloud_coc_document_execution_steps` [GH-7733]
  + `huaweicloud_coc_document_executions` [GH-7732]
  + `huaweicloud_coc_enterprise_project_collections` [GH-7732]
  + `huaweicloud_coc_group_resource_relations` [GH-7732]
  + `huaweicloud_coc_instance_batches` [GH-7732]
  + `huaweicloud_csms_secret_tags` [GH-7704]
  + `huaweicloud_dms_rocketmq_consumer_group_topics` [GH-7703]
  + `huaweicloud_dms_rocketmq_instance_diagnoses` [GH-7712]
  + `huaweicloud_eg_event_subscriptions` [GH-7752]
  + `huaweicloud_eg_traced_events` [GH-7620]
  + `huaweicloud_fgs_async_invocations` [GH-7713]
  + `huaweicloud_lts_logs` [GH-7738]
  + `huaweicloud_lts_member_group_streams` [GH-7739]
  + `huaweicloud_lts_stream_charts` [GH-7734]
  + `huaweicloud_secmaster_alert_rule_template_detail` [GH-7680]
  + `huaweicloud_secmaster_collector_logstash_parsers` [GH-7702]
  + `huaweicloud_secmaster_component_running_nodes` [GH-7736]
  + `huaweicloud_secmaster_components` [GH-7736]
  + `huaweicloud_secmaster_mappings_functions` [GH-7717]
  + `huaweicloud_secmaster_operation_connections` [GH-7717]
  + `huaweicloud_secmaster_reports_emails` [GH-7743]
  + `huaweicloud_secmaster_retrieve_scripts` [GH-7705]
  + `huaweicloud_secmaster_siem_directories` [GH-7747]
  + `huaweicloud_secmaster_table_consumption` [GH-7706]
  + `huaweicloud_secmaster_table_histograms` [GH-7715]
  + `huaweicloud_secmaster_tables` [GH-7705]
  + `huaweicloud_secmaster_vulnerabilities` [GH-7740]
  + `huaweicloud_secmasterv2_alert_rule_template_detail` [GH-7680]

## 1.78.0 (Aug 30, 2025)

* **New Resource:**
  + `huaweicloud_metastudio_instance` [GH-7691]

## 1.77.7 (Aug 29, 2025)

ENHANCEMENTS:

* add assume role with oidc support [GH-7683]

## 1.77.6 (Aug 27, 2025)

ENHANCEMENTS:

* resource/huaweicloud_apig_application_authorization: support partial management [GH-7645]

## 1.77.5 (Aug 22, 2025)

ENHANCEMENTS:

* refactor to use ErrDefault409 err check [GH-7595]

## 1.77.4 (Aug 19, 2025)

ENHANCEMENTS:

* refactor GetRawConfigTags to allow empty value [GH-7560]

## 1.77.3 (Aug 8, 2025)

ENHANCEMENTS:

* resource/huaweicloud_vpcep_endpoint: add policy_document support [GH-7466]

## 1.77.2 (Aug 1, 2025)

ENHANCEMENTS:

* resource/huaweicloud_cce_node_pool: update min/max node count [GH-7409]

## 1.77.1 (July 31, 2025)

* **New Resource:**
  + `huaweicloud_cce_nodes_remove` [GH-7388]

ENHANCEMENTS:

* add endpoints key and description on documentation [GH-7342]
* add default_tags support on provider [GH-7334]

## 1.76.5 (July 15, 2025)

ENHANCEMENTS:

* resource/huaweicloud_cce_cluster: add encryption_config support [GH-7257]

## 1.76.4 (July 11, 2025)

ENHANCEMENTS:

* resource/huaweicloud_vpcep_endpoint: add condition support [GH-7237]

## 1.76.3 (July 10, 2025)

ENHANCEMENTS:

* resource/huaweicloud_modelarts_resource_pool: add data_volumes sorting logic [GH-7222]

## 1.76.2 (July 7, 2025)

ENHANCEMENTS:

* resource/huaweicloud_modelartsv2_node_batch_unsubscribe: use v2 resource name [GH-7116]

## 1.76.1 (July 2, 2025)

ENHANCEMENTS:

* update docs for temporary security credentials and backend [GH-7170]

## 1.76.0 (June 30, 2025)

* **New Resource:**
  + `huaweicloud_cbr_replicate_backup` [GH-7111]
  + `huaweicloud_sdrs_protected_instance_delete_nic` [GH-7096]
  + `huaweicloud_rds_restore_read_replica_database` [GH-7095]
  + `huaweicloud_sdrs_protected_instance_add_nic` [GH-7091]

* **New Data Source:**
  + `huaweicloud_rds_backup_database` [GH-7110]
  + `huaweicloud_cbr_backup_metadata` [GH-7088]

ENHANCEMENTS:

* resource/huaweicloud_workspace_desktop: suppress diff for name case changing [GH-7115]

## 1.75.5 (June 20, 2025)

ENHANCEMENTS:

* resource/huaweicloud_vpc_route_table: update changing rule logic [GH-7057]

## 1.75.4 (June 13, 2025)

ENHANCEMENTS:

* resource/huaweicloud_compute_interface_attach: add check deleted support [GH-7014]

## 1.75.3 (June 12, 2025)

ENHANCEMENTS:

* resource/huaweicloud_vpc_subnet: enhance dhcp_domain_name update [GH-6939]
* resource/huaweicloud_dms_kafkav2_smart_connect_task: add agency_name support [GH-6945]

## 1.75.2 (June 4, 2025)

ENHANCEMENTS:

* resource/huaweicloud_dms_kafka_instance: add kms_encrypted_password support [GH-6929]

## 1.75.1 (June 4, 2025)

ENHANCEMENTS:

* resource/huaweicloud_as_group: add pagination support [GH-6918]

## 1.75.0 (May 29, 2025)

* **New Data Source:**
  + `huaweicloud_rds_diagnosis` [GH-6882]
  + `huaweicloud_sms_migration_projects` [GH-6863]

ENHANCEMENTS:

* resource/huaweicloud_coc_script_execute: add is_sync parameter support [GH-6873]

## 1.74.1 (May 15, 2025)

* **New Source:**
  + `huaweicloud_cbc_resources_unsubscribe` [GH-6790]

ENHANCEMENTS:

* resource/huaweicloud_vpc_eip: report the error if the whole bandwidth is sold out [GH-6771]
* resource/huaweicloud_modelarts_resource_pool: disable the sorting logic for resources parameter [GH-6776]

## 1.74.0 (April 30, 2025)

* **New Data Source:**
  + `huaweicloud_deh_instances` [GH-6728]

## 1.73.9 (April 24, 2025)

ENHANCEMENTS:

* resource/huaweicloud_modelarts_resource_pool: add volume attaching support [GH-6718]

## 1.73.8 (April 14, 2025)

BUG FIXES:

* resource/huaweicloud_cce_node_attach: fix attributes setting issue [GH-6659]

## 1.73.7 (April 11, 2025)

ENHANCEMENTS:

* add waiting for logic for secmaster post_paid_order and workspace [GH-6647]

## 1.73.6 (April 7, 2025)

ENHANCEMENTS:

* make v5 assume role duration configurable [GH-6616]

## 1.73.5 (April 3, 2025)

* **New Data Source:**
  + `huaweicloud_networking_secgroup_tags` [GH-6611]

BUG FIXES:

* resource/huaweicloud_cfw_black_white_list: fix address issue [GH-6608]

## 1.73.4 (April 2, 2025)

* **New Resource:**
  + `huaweicloud_secmaster_workflow_action` [GH-6595]

## 1.73.3 (March 24, 2025)

* **New Data Source:**
  + `huaweicloud_ims_tags` [GH-6559]

## 1.73.2 (March 21, 2025)

* **New Data Source:**
  + `huaweicloud_cce_addons` [GH-6540]

## 1.73.1 (March 7, 2025)

ENHANCEMENTS:

* resource/huaweicloud_css_cluster: remove engine_type validate to support opensearch cluster [GH-6474]

## 1.73.0 (February 28, 2025)

* **New Resource:**
  + `huaweicloud_ces_metric_data_add` [GH-6290]
  + `huaweicloud_css_snapshot_restore` [GH-6403]

* **New Data Source:**
  + `huaweicloud_aom_access_codes` [GH-6411]

ENHANCEMENTS:

* resource/huaweicloud_compute_instance: add uuid and fixed_ip_v4 update support [GH-6410]

## 1.72.3 (February 20, 2025)

ENHANCEMENTS:

* resource/huaweicloud_cce_cluster: add enable_distribute_management support [GH-6375]

## 1.72.2 (February 19, 2025)

* **New Resource:**
  + `huaweicloud_cdn_domain_rule` [GH-6303]
  + `huaweicloud_ces_metric_data_add` [GH-6290]

* **New Data Source:**
  + `huaweicloud_ces_metric_data` [GH-6308]

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_opengauss_instance: add advance_features support [GH-6300]

## 1.72.1 (January 15, 2025)

* **New Resource:**
  + `huaweicloud_cae_vpc_egress` [GH-6175]
  + `huaweicloud_cae_environment` [GH-6134]
  + `huaweicloud_cae_notification_rule` [GH-6132]

* **New Data Source:**
  + `huaweicloud_cae_notification_rules` [GH-6168]

ENHANCEMENTS:

* resource/huaweicloud_cce_autopilot_cluster: add version support [GH-6191]

## 1.72.0 (December 31, 2024)

* **New Resource:**
  + `huaweicloud_live_stream_delay` [GH-6103]
  + `huaweicloud_modelarts_devserver_action` [GH-6098]
  + `huaweicloud_gaussdb_opengauss_restore` [GH-6093]
  + `huaweicloud_modelarts_devserver` [GH-6088]
  + `huaweicloud_gaussdb_opengauss_parameter_template` [GH-6059]

* **New Data Source:**
  + `huaweicloud_live_channels` [GH-6090]
  + `huaweicloud_cce_autopilot_cluster_certificate` [GH-6074]
  + `huaweicloud_cts_resources` [GH-6068]
  + `huaweicloud_gaussdb_opengauss_backups` [GH-6063]
  + `huaweicloud_cts_operations` [GH-6058]
  + `huaweicloud_cce_autopilot_addon_templates` [GH-6056]

ENHANCEMENTS:

* add default value for shared config file [GH-6091]
* resource/huaweicloud_opengauss_instance: add configuration id update support [GH-6110]
* resource/huaweicloud_identitycenter_permission_set: add tags support [GH-6099]

## 1.71.2 (December 13, 2024)

ENHANCEMENTS:

* resource/huaweicloud_dms_kafka_instance: add port protocols support [GH-6016]

## 1.71.1 (December 9, 2024)

ENHANCEMENTS:

* resource/huaweicloud_cce_cluster_openid_jwks: add header for jwks [GH-5988]

## 1.71.0 (November 30, 2024)

* **New Resource:**
  + `huaweicloud_cce_autopilot_addon` [GH-5949]
  + `huaweicloud_hss_cce_protection` [GH-5939]
  + `huaweicloud_rds_mysql_proxy_restart` [GH-5914]
  + `huaweicloud_workspace_app_image_server` [GH-5908]
  + `huaweicloud_workspace_app_personal_folders` [GH-5906]
  + `huaweicloud_rds_mysql_proxy` [GH-5893]

* **New Data Source:**
  + `huaweicloud_live_domains` [GH-5950]
  + `huaweicloud_smn_topic_subscriptions` [GH-5947]
  + `huaweicloud_css_snapshots` [GH-5944]
  + `huaweicloud_smn_logtanks` [GH-5912]
  + `huaweicloud_smn_subscriptions` [GH-5901]
  + `huaweicloud_rds_mysql_proxies` [GH-5900]
  + `huaweicloud_cph_server_bandwidths` [GH-5883]

ENHANCEMENTS:

* resource/huaweicloud_cce_node_pool: change data_volumes to be optional [GH-5934]
* resource/huaweicloud_cce_node_pool: add subnet_list support [GH-5934]
* resource/huaweicloud_cce_autopilot_cluster: add cluster update support [GH-5872]

## 1.70.3 (November 15, 2024)

ENHANCEMENTS:

* resource/huaweicloud_eg_event_subscription: add source.detail support [GH-5843]

## 1.70.2 (November 8, 2024)

BUG FIXES:

* resource/huaweicloud_ddm_instance: fix pathSearch issue [GH-5832]

## 1.70.1 (November 7, 2024)

ENHANCEMENTS:

* resource/huaweicloud_coc_script: handle more NotFound errors [GH-5825]

## 1.70.0 (October 30, 2024)

* **New Resource:**
  + `huaweicloud_cph_phone_stop` [GH-5769]
  + `huaweicloud_dds_primary_standby_switch` [GH-5768]
  + `huaweicloud_gaussdb_mysql_lts_log` [GH-5756]
  + `huaweicloud_dds_recycle_policy` [GH-5755]
  + `huaweicloud_vpn_access_policy` [GH-5751]
  + `huaweicloud_workspace_app_group` [GH-5747]

* **New Data Source:**
  + `huaweicloud_gaussdb_mysql_slow_logs` [GH-5777]
  + `huaweicloud_dds_error_logs` [GH-5772]
  + `huaweicloud_vpn_access_policies` [GH-5770]
  + `huaweicloud_ram_resource_shares` [GH-5749]
  + `huaweicloud_gaussdb_redis_flavors` [GH-5746]
  + `huaweicloud_iotda_batchtasks` [GH-5745]

ENHANCEMENTS:

* resource/huaweicloud_vpc_address_group: add ip_extra_set support [GH-5767]
* resource/huaweicloud_identity_policy: add policy_document updating support [GH-5757]

## 1.69.1 (October 16, 2024)

ENHANCEMENTS:

* resource/huaweicloud_cce_node_pool: add updating os support [GH-5681]

## 1.69.0 (September 29, 2024)

* **New Resource:**
  + `huaweicloud_servicestagev3_environment_associate` [GH-5634]
  + `huaweicloud_ces_microservice_engine_configuration` [GH-5633]
  + `huaweicloud_servicestagev3_application` [GH-5570]
  + `huaweicloud_dws_snapshot_copy` [GH-5596]
  + `huaweicloud_rgc_account` [GH-5571]
  + `huaweicloud_ces_one_click_alarm` [GH-5613]
  + `huaweicloud_dws_cluster_restart` [GH-5568]
  + `huaweicloud_gaussdb_mysql_recycling_policy` [GH-5590]
  + `huaweicloud_antiddos_default_protection_policy` [GH-5592]

* **New Data Source:**
  + `huaweicloud_vpn_p2c_gateways` [GH-5625]
  + `huaweicloud_ces_one_click_alarm_rules` [GH-5617]
  + `huaweicloud_gaussdb_influx_instances` [GH-5616]
  + `huaweicloud_cse_microservice_engines` [GH-5619]
  + `huaweicloud_gaussdb_mysql_incremental_backups` [GH-5600]
  + `huaweicloud_dbss_audit_sql_injection_rules` [GH-5615]
  + `huaweicloud_dbss_audit_data_masking_rules` [GH-5611]
  + `huaweicloud_dbss_audit_risk_rules` [GH-5610]
  + `huaweicloud_dbss_audit_rule_scopes` [GH-5603]
  + `huaweicloud_secmaster_playbook_statistics` [GH-5579]
  + `huaweicloud_dws_cluster_topo_rings` [GH-5604]
  + `huaweicloud_dws_statistics` [GH-5555]
  + `huaweicloud_gaussdb_mysql_auto_scaling_records` [GH-5601]

ENHANCEMENTS:

* resource/huaweicloud_cdn_domain: add user_agent_filter support [GH-5631]
* datasource/huaweicloud_compute_flavors: add storage_type support [GH-5621]
* resource/huaweicloud_cce_cluster: add upgrade version support [GH-5608]
* resource/huaweicloud_vpc_eip: add update enterprise project id support [GH-5587]

## 1.68.0 (August 30, 2024)

* **New Resource:**
  + `huaweicloud_ccm_private_ca_restore` [GH-5466]
  + `huaweicloud_ces_dashboard_widget` [GH-5447]
  + `huaweicloud_ims_cbr_whole_image` [GH-5440]
  + `huaweicloud_gaussdb_mysql_instance_restart` [GH-5435]
  + `huaweicloud_ims_ecs_whole_image` [GH-5430]
  + `huaweicloud_ccm_private_certificate_revoke` [GH-5429]
  + `huaweicloud_ccm_private_ca_revoke` [GH-5425]
  + `huaweicloud_ims_ecs_system_image` [GH-5411]

* **New Data Source:**
  + `huaweicloud_ces_alarm_templates` [GH-5461]
  + `huaweicloud_rms_advanced_query` [GH-5458]
  + `huaweicloud_ces_dashboard_widgets` [GH-5456]
  + `huaweicloud_rms_resource_aggregator_advanced_query` [GH-5453]
  + `huaweicloud_secmaster_alert_rule_templates` [GH-5452]
  + `huaweicloud_csms_secret_version_state` [GH-5450]
  + `huaweicloud_rms_resource_aggregator_policy_assignments` [GH-5448]
  + `huaweicloud_secmaster_alert_rules` [GH-5444]
  + `huaweicloud_rms_resource_aggregator_policy_states` [GH-5441]
  + `huaweicloud_rms_resource_aggregator_discovered_resources` [GH-5439]
  + `huaweicloud_secmaster_indicators` [GH-5437]

ENHANCEMENTS:

* resource/huaweicloud_elb_loadbalancer: add gateway type support [GH-5463]
* resource/huaweicloud_csms_secret: add epsId update support [GH-5434]
* resource/huaweicloud_kms_key: add epsId update support [GH-5432]

## 1.67.1 (August 9, 2024)

BUG FIXES:

* resource/huaweicloud_compute_instance: fix user_data problem while updating hostname [GH-5373]

## 1.67.0 (July 30, 2024)

* **New Resource:**
  + `huaweicloud_dataarts_dataservice_api_debug` [GH-5306]
  + `huaweicloud_dataarts_architecture_batch_publishment` [GH-5302]
  + `huaweicloud_rds_pg_account_privileges` [GH-5293]
  + `huaweicloud_evs_volume_transfer_accepter` [GH-5290]
  + `huaweicloud_rds_pg_database_privilege` [GH-5284]
  + `huaweicloud_evs_volume_transfer` [GH-5277]

* **New Data Source:**
  + `huaweicloud_dataarts_dataservice_instances` [GH-5297]
  + `huaweicloud_organizations_trusted_services` [GH-5295]
  + `huaweicloud_organizations_services` [GH-5294]
  + `huaweicloud_kms_keys` [GH-5289]
  + `huaweicloud_ccm_private_ca_export` [GH-5285]
  + `huaweicloud_cfw_capture_task_results` [GH-5280]

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_mysql_proxy: add access control support [GH-5270]

## 1.66.3 (July 19, 2024)

ENHANCEMENTS:

* resource/huaweicloud_cce_node: add spot support [GH-5245]

## 1.66.2 (July 17, 2024)

BUG FIXES:

* resource/huaweicloud_cts_data_tracker: add nil check before set agency_name [GH-5219]

## 1.66.1 (July 15, 2024)

ENHANCEMENTS:

* resource/huaweicloud_vpc_network_acl: add roll back for updating rules [GH-5201]

## 1.66.0 (June 29, 2024)

* **New Resource:**
  + `huaweicloud_css_cluster_az_migrate` [GH-5078]
  + `huaweicloud_access_analyzer_archive_rule` [GH-5093]
  + `huaweicloud_access_analyzer` [GH-5029]
  + `huaweicloud_ddm_instance_restart` [GH-5062]
  + `huaweicloud_ccm_certificate` [GH-5058]
  + `huaweicloud_asm_mesh` [GH-5041]
  + `huaweicloud_enterprise_project_authority` [GH-5094]

* **New Data Source:**
  + `huaweicloud_cdn_logs` [GH-5074]
  + `huaweicloud_er_available_routes` [GH-5073]
  + `huaweicloud_lts_aom_accesses` [GH-5075]
  + `huaweicloud_dms_kafkav2_smart_connect_tasks` [GH-5065]
  + `huaweicloud_lts_host_groups` [GH-5045]
  + `huaweicloud_lts_streams` [GH-5042]
  + `huaweicloud_hss_quotas` [GH-5028]

ENHANCEMENTS:

* resource/huaweicloud_vpc_network_acl: add tags support [GH-5076]
* resource/huaweicloud_css_cluster: add shrinking node support [GH-5044]
* resource/huaweicloud_cts_tracker: add deletion support [GH-5026]

## 1.65.2 (June 17, 2024)

BUG FIXES:

* data/huaweicloud_as_configurations: fix configurations set issue [GH-5014]

## 1.65.1 (June 13, 2024)

ENHANCEMENTS:

* authentication: add sso config support [GH-4998]

## 1.65.0 (June 3, 2024)

* **New Resource:**
  + `huaweicloud_apig_application_quota` [GH-4922]
  + `huaweicloud_apig_endpoint_connection_management` [GH-4921]
  + `huaweicloud_css_cluster_restart` [GH-4916]
  + `huaweicloud_hss_quota` [GH-4881]

* **New Data Source:**
  + `huaweicloud_apig_channels` [GH-4892]
  + `huaweicloud_rms_resource_aggregation_pending_requests` [GH-4893]
  + `huaweicloud_dms_rabbitmq_extend_flavors` [GH-4888]
  + `huaweicloud_cts_trackers` [GH-4880]
  + `huaweicloud_rms_resources` [GH-4879]
  + `huaweicloud_dds_database_roles` [GH-4876]

ENHANCEMENTS:

* resource/huaweicloud_cdn_domain: add access_area_filter support [GH-4939]
* resource/huaweicloud_dms_kafka_instance: add arch_type support [GH-4919]
* resource/huaweicloud_organizations_account: add delete support [GH-4889]

## 1.64.4 (May 26, 2024)

ENHANCEMENTS:

* resource/huaweicloud_drs_job: make it possible to delete job if no order id returned [GH-4894]

## 1.64.3 (May 25, 2024)

BUG FIXES:

* resource/huaweicloud_drs_job: raise error if orderID is empty [GH-4887]

## 1.64.2 (May 17, 2024)

BUG FIXES:

* resource/huaweicloud_swr_image_trigger: Make cluster name/id optional for cci trigger [GH-4832]
* resource/huaweicloud_vpc_network_acl: Fix the problem of updating acl rules [GH-4821]

## 1.64.1 (May 6, 2024)

ENHANCEMENTS:

* resource/huaweicloud_vpn_gateway: add er_attachment_id support [GH-4731]

## 1.64.0 (April 29, 2024)

* **New Resource:**
  + `huaweicloud_ram_resource_share_accepter` [GH-4710]
  + `huaweicloud_dds_instance_restart` [GH-4691]
  + `huaweicloud_dataarts_security_permission_set_privilege` [GH-4685]
  + `huaweicloud_identity_login_policy` [GH-4681]
  + `huaweicloud_rds_recycling_policy` [GH-4678]
  + `huaweicloud_dataarts_security_permission_set_member` [GH-4663]
  + `huaweicloud_dataarts_dataservice_catalog` [GH-4651]
  + `huaweicloud_cc_global_connection_bandwidth_spec_codes` [GH-4649]
  + `huaweicloud_lts_log_converge` [GH-4633]
  + `huaweicloud_lts_log_converge_switch` [GH-4631]
  + `huaweicloud_cae_component_deployment` [GH-4609]

* **New Data Source:**
  + `huaweicloud_ram_resource_share_invatations` [GH-4712]
  + `huaweicloud_dataarts_security_data_secrecy_level` [GH-4692]
  + `huaweicloud_cfw_attack_logs` [GH-4683]
  + `huaweicloud_rds_cross_region_backup_instances` [GH-4682]
  + `huaweicloud_cc_permissions` [GH-4675]
  + `huaweicloud_waf_rules_blacklist` [GH-4673]
  + `huaweicloud_dcs_accounts` [GH-4667]
  + `huaweicloud_waf_rules_geolocation_access_control` [GH-4662]
  + `huaweicloud_cc_connection_tags` [GH-4659]
  + `huaweicloud_dli_sql_templates` [GH-4658]
  + `huaweicloud_identity_providers` [GH-4656]
  + `huaweicloud_waf_rules_precise_protection` [GH-4655]
  + `huaweicloud_cc_global_connection_bandwidth_sites` [GH-4653]
  + `huaweicloud_rds_cross_region_backups` [GH-4647]
  + `huaweicloud_rms_resource_aggregators` [GH-4645]
  + `huaweicloud_dli_spark_templates` [GH-4641]
  + `huaweicloud_cc_connection_routes` [GH-4640]
  + `huaweicloud_dli_flink_templates` [GH-4639]
  + `huaweicloud_hss_hosts` [GH-4638]
  + `huaweicloud_cc_central_network_policies_change_set` [GH-4635]
  + `huaweicloud_cc_global_connection_bandwidth_line_levels` [GH-4634]
  + `huaweicloud_cc_connection_routes` [GH-4623]
  + `huaweicloud_cc_central_network_capabilities` [GH-4620]
  + `huaweicloud_dns_nameservers` [GH-4617]
  + `huaweicloud_rds_sql_audit_operations` [GH-4615]
  + `huaweicloud_cc_bandwidth_packages` [GH-4613]
  + `huaweicloud_dli_spark_templates` [GH-4610]
  + `huaweicloud_hss_host_groups` [GH-4606]
  + `huaweicloud_hss_hosts` [GH-4605]
  + `huaweicloud_cc_network_instances` [GH-4600]
  + `huaweicloud_rms_advanced_queries` [GH-4596]
  + `huaweicloud_cc_inter_region_bandwidths` [GH-4595]
  + `huaweicloud_rms_organizational_assignment_packages` [GH-4592]
  + `huaweicloud_gaussdb_mysql_restore_time_ranges` [GH-4587]

ENHANCEMENTS:

* resource/huaweicloud_cfw_firewall: Add east-west firewall create afterwards support [GH-4703]
* resource/huaweicloud_cfw_firewall: Add tags update support [GH-4693]
* resource/huaweicloud_cfw_firewall: Add attachment_id import support [GH-4689]
* resource/huaweicloud_rds_instance: Add read write permission support [GH-4679]
* resource/huaweicloud_cbh_instance: Add security group update support [GH-4672]
* resource/huaweicloud_rds_read_replica_instance: Add maintain window support [GH-4660]

BUG FIXES:

* resource/huaweicloud_vpc: Fix return error of vpc v3 API [GH-4725]

## 1.63.2 (April 25, 2024)

ENHANCEMENTS:

* authentication: add v5 agency assume support [GH-4690]

## 1.63.1 (April 12, 2024)

* **New Resource:**
  + `huaweicloud_gaussdb_mysql_restore` [GH-4580]

ENHANCEMENTS:

* resource/huaweicloud_cdn_domain: Add epsId update support [GH-4597]
* resource/huaweicloud_cdn_domain: Add certificate_type and ocsp_stapling_status support [GH-4576]

BUG FIXES:

* resource/huaweicloud_cfw_address_group_member: fix the issue of fetching member with ip range [GH-4582]

## 1.63.0 (March 30, 2024)

* **New Resource:**
  + `huaweicloud_cfw_dns_resolution` [GH-4485]
  + `huaweicloud_identity_service_agency` [GH-4497]
  + `huaweicloud_lts_cce_access` [GH-4481]
  + `huaweicloud_ddm_instance_read_strategy` [GH-4450]
  + `huaweicloud_hss_host_protection` [GH-4474]
  + `huaweicloud_cc_global_connection_bandwidth_associate` [GH-4438]
  + `huaweicloud_css_logstash_custom_certificate` [GH-4416]
  + `huaweicloud_lts_cross_account_access` [GH-4423]
  + `huaweicloud_css_logstash_pipeline` [GH-4409]
  + `huaweicloud_dli_datasource_connection_privilege` [GH-4397]
  + `huaweicloud_workspace_desktop_name_rule` [GH-4381]
  + `huaweicloud_dcs_diagnosis_task` [GH-4385]

* **New Data Source:**
  + `huaweicloud_cfw_address_groups` [GH-4476]
  + `huaweicloud_iotda_device_certificates` [GH-4523]
  + `huaweicloud_iotda_devices` [GH-4526]
  + `huaweicloud_rms_services` [GH-4517]
  + `huaweicloud_dli_quotas` [GH-4521]
  + `huaweicloud_rms_advanced_query_schemas` [GH-4520]
  + `huaweicloud_rms_assignment_packages` [GH-4509]
  + `huaweicloud_er_propagations` [GH-4483]
  + `huaweicloud_huaweicloud_er_associations` [GH-4483]
  + `huaweicloud_rms_policy_assignments` [GH-4479]
  + `huaweicloud_iotda_amqps` [GH-4464]
  + `huaweicloud_dcs_bigkey_analyses` [GH-4454]
  + `huaweicloud_iotda_dataforwarding_rules` [GH-4480]
  + `huaweicloud_rms_regions` [GH-4473]
  + `huaweicloud_iotda_spaces` [GH-4436]
  + `huaweicloud_iotda_products` [GH-4433]
  + `huaweicloud_dws_workload_queues` [GH-4353]
  + `huaweicloud_dc_virtual_interfaces` [GH-4390]

ENHANCEMENTS:

* resource/huaweicloud_compute_instance: add extend_param parameter to bandwidth [GH-4534]
* resource/huaweicloud_drs_job: support policy config in job [GH-4530]
* resource/huaweicloud_drs_job: support charging mode in job [GH-4458]
* resource/huaweicloud_workspace_desktop: allow independent system disk updates [GH-4522]
* resource/huaweicloud_rds_instance: rds instance support msdtc hosts [GH-4494]
* resource/huaweicloud_cce_node: add new param details in the data_source of cce nodes [GH-4451]
* resource/huaweicloud_cfw_protection_rule: add explanation for protection rule type [GH-4440]
* resource/huaweicloud_mrs_cluster: cluster support prepaid charge mode [GH-4383]

BUG FIXES:

* resource/huaweicloud_lts_stream: fix the problem that eps cannot be set [GH-4527]

## 1.62.1 (March 6, 2024)

* **New Resource:**
  + `huaweicloud_css_logstash_configuration` [GH-4300]
  + `huaweicloud_fgs_function_trigger` [GH-4327]

* **New Data Source:**
  + `huaweicloud_csms_events` [GH-4292]
  + `huaweicloud_ga_endpoint_groups` [GH-4306]
  + `huaweicloud_ram_shared_principals` [GH-4313]
  + `huaweicloud_csms_secrets` [GH-4318]
  + `huaweicloud_dws_disaster_recovery_tasks` [GH-4325]

ENHANCEMENTS:

* resource/huaweicloud_fgs_function: Support config concurrent requests number [GH-4293]
* resource/huaweicloud_cc_connection: Add tags support [GH-4340]
* resource/huaweicloud_nat_snat_rule: Support `global_eip_id` parameter [GH-4345]
* resource/huaweicloud_drs_job: Support updating `tags` [GH-4352]

BUG FIXES:

* resource/huaweicloud_cce_node: Fix can't get resource ID issue when create a BMS node [GH-4323]

## 1.62.0 (February 29, 2024)

* **New Resource:**
  + `huaweicloud_dds_lts_log` [GH-4193]
  + `huaweicloud_identity_user_token` [GH-4237]
  + `huaweicloud_vpc_internet_gateway` [GH-4239]
  + `huaweicloud_dms_rocketmq_migration_task` [GH-4244]
  + `huaweicloud_rms_organizational_policy_assignment` [GH-4249]
  + `huaweicloud_cc_global_connection_bandwidth` [GH-4267]
  + `huaweicloud_css_logstash_cluster` [GH-4268]
  + `huaweicloud_global_eip_associate` [GH-4278]
  + `huaweicloud_dcs_account` [GH-4283]
  + `huaweicloud_dws_workload_plan_execution` [GH-4238]
  + `huaweicloud_dws_disaster_recovery_task` [GH-4262]
  + `huaweicloud_iotda_batchtask_file` [GH-4261]
  + `huaweicloud_iotda_upgrade_package` [GH-4282]

* **New Data Source:**
  + `huaweicloud_cdn_domain_certificates` [GH-4128]
  + `huaweicloud_cdn_domains` [GH-4221]
  + `huaweicloud_nat_snat_rules` [GH-4174]
  + `huaweicloud_nat_dnat_rules` [GH-4214]
  + `huaweicloud_dds_storage_types` [GH-4240]
  + `huaweicloud_fgs_application_templates` [GH-4242]
  + `huaweicloud_vpc_internet_gateways` [GH-4252]
  + `huaweicloud_bms_instances` [GH-4260]
  + `huaweicloud_ga_accelerators` [GH-4265]
  + `huaweicloud_ga_listeners` [GH-4275]
  + `huaweicloud_ga_address_groups` [GH-4285]
  + `huaweicloud_global_internet_bandwidths` [GH-4286]
  + `huaweicloud_global_eips` [GH-4299]

ENHANCEMENTS:

* resource/huaweicloud_identity_user: Support the verification method of user login protect [GH-4247]
* resource/huaweicloud_dc_virtual_interface: Support `resource_tenant_id` parameter [GH-4253]
* resource/huaweicloud_dli_queue: Support associate with a elastic resource pool [GH-4254]
* resource/huaweicloud_dms_kafka_instance: Support `parameters` block [GH-4263]
* resource/huaweicloud_evs_snapshot: Support `metadata` parameter [GH-4266]
* resource/huaweicloud_fgs_function: Support manage the reserved instance policies [GH-4272]
* resource/huaweicloud_dws_cluster: Support multiple availability zones [GH-4273]
* resource/huaweicloud_cce_cluster: Support resizing cluster with `flavor_id` [GH-4280]
* resource/huaweicloud_modelarts_resource_pool: Support pre-paid charging mode and lite resource pool [GH-4284]
* resource/huaweicloud_dws_cluster: Support to enable or disable LTS [GH-4287]
* resource/huaweicloud_dcs_instance: Support to enable or disable SSL [GH-4291]

## 1.61.1 (February 7, 2024)

* **New Resource:**
  + `huaweicloud_coc_script` [GH-3682]
  + `huaweicloud_coc_script_execute` [GH-3757]
  + `huaweicloud_dcs_bigkey_analysis` [GH-4030]
  + `huaweicloud_dcs_hotkey_analysis` [GH-4070]
  + `huaweicloud_fgs_function_event` [GH-4080]
  + `huaweicloud_vpc_network_acl` [GH-4111]
  + `huaweicloud_dws_workload_plan` [GH-4135]
  + `huaweicloud_dws_workload_plan_stage` [GH-4164]
  + `huaweicloud_dli_elastic_resource_pool` [GH-4186]

* **New Data Source:**
  + `huaweicloud_as_policy_execute_logs` [GH-4094]
  + `huaweicloud_compute_instance_remote_console` [GH-4123]
  + `huaweicloud_dli_datasource_connections` [GH-4130]
  + `huaweicloud_dli_datasource_auths` [GH-4146]
  + `huaweicloud_ram_shared_resources` [GH-4151]
  + `huaweicloud_dc_virtual_gateways` [GH-4184]

  + `huaweicloud_rds_pg_databases` [GH-4085]
  + `huaweicloud_rds_mysql_database_privileges` [GH-4092]
  + `huaweicloud_rds_sqlserver_databases` [GH-4117]
  + `huaweicloud_rds_sqlserver_database_privileges` [GH-4150]
  + `huaweicloud_rds_sqlserver_accounts` [GH-4158]

ENHANCEMENTS:

* resource/huaweicloud_bms_instance: Support updating `nics` block [GH-4073]
* resource/huaweicloud_bms_instance: Support `metadata` parameter [GH-4100]
* resource/huaweicloud_rds_instance: Support updating volume size in prepaid mode [GH-4078]
* resource/huaweicloud_vpc_bandwidth: Support changing billing mode to prePaid [GH-4113]
* resource/huaweicloud_vpc_eip: Support changing billing mode to prePaid [GH-4122]
* resource/huaweicloud_obs_bucket: Support agency configuration in `logging` block [GH-4138]
* resource/huaweicloud_identity_project: Support `status` parameter [GH-4143]
* resource/huaweicloud_vpn_gateway: Add tags support [GH-4149]
* resource/huaweicloud_vpn_customer_gateway: Add tags support [GH-4160]
* resource/huaweicloud_vpn_connection: Add tags support [GH-4161]

BUG FIXES:

* resource/huaweicloud_css_cluster: Fix the issue when creating cluster in postPaid billing mode [GH-4125]

## 1.61.0 (January 29, 2024)

* **New Resource:**
  + `huaweicloud_cdm_cluster_action` [GH-3408]
  + `huaweicloud_cc_central_network_attachment` [GH-3925]
  + `huaweicloud_cc_authorization` [GH-3973]
  + `huaweicloud_elb_active_standby_pool` [GH-3934]
  + `huaweicloud_identity_virtual_mfa_device` [GH-3948]
  + `huaweicloud_dataarts_architecture_code_table_values` [GH-3954]
  + `huaweicloud_compute_auto_launch_group` [GH-3980]
  + `huaweicloud_dms_rabbitmq_plugin` [GH-3985]
  + `huaweicloud_cfw_domain_name_group` [GH-3987]
  + `huaweicloud_dws_workload_queue` [GH-3994]
  + `huaweicloud_dws_logical_cluster` [GH-4024]
  + `huaweicloud_fgs_dependency_version` [GH-3999]
  + `huaweicloud_aom_prom_instance` [GH-4011]
  + `huaweicloud_dc_hosted_connect` [GH-4013]
  + `huaweicloud_fgs_application` [GH-4043]
  + `huaweicloud_oms_migration_sync_task` [GH-4049]
  + `huaweicloud_rds_pg_hba` [GH-4053]
  + `huaweicloud_ga_address_group` [GH-4059]
  + `huaweicloud_global_internet_bandwidth` [GH-4079]
  + `huaweicloud_global_eip` [GH-4086]
  + `huaweicloud_css_scan_task` [GH-4081]

  + `huaweicloud_vpc_traffic_mirror_filter` [GH-3944]
  + `huaweicloud_vpc_traffic_mirror_filter_rule` [GH-3968]
  + `huaweicloud_vpc_network_interface` [GH-3990]
  + `huaweicloud_vpc_sub_network_interface` [GH-3965]

* **New Data Source:**
  + `huaweicloud_networking_secgroup_rules` [GH-3930]
  + `huaweicloud_swr_image_tags` [GH-3943]
  + `huaweicloud_as_policies` [GH-3964]
  + `huaweicloud_waf_domains` [GH-3966]
  + `huaweicloud_dws_logical_cluster_rings` [GH-3995]
  + `huaweicloud_dms_rabbitmq_plugins` [GH-4000]
  + `huaweicloud_as_lifecycle_hooks` [GH-4002]
  + `huaweicloud_elb_active_standby_pools` [GH-4003]
  + `huaweicloud_elb_monitors` [GH-4016]
  + `huaweicloud_cts_notifications` [GH-4019]
  + `huaweicloud_rds_mysql_binlog` [GH-4027]
  + `huaweicloud_workspace_desktops` [GH-4032]
  + `huaweicloud_vpc_bandwidths` [GH-4034]
  + `huaweicloud_dms_rabbitmq_instances` [GH-4035]
  + `huaweicloud_nat_private_snat_rules` [GH-4036]
  + `huaweicloud_nat_private_dnat_rules` [GH-4095]
  + `huaweicloud_dms_rocketmq_flavors` [GH-4057]
  + `huaweicloud_cbh_flavors` [GH-4061]
  + `huaweicloud_rds_pg_accounts` [GH-4068]
  + `huaweicloud_global_eip_pools` [GH-4069]
  + `huaweicloud_global_eip_access_sites` [GH-4082]

  + `huaweicloud_vpcep_service_connections` [GH-3947]
  + `huaweicloud_vpcep_service_permissions` [GH-3955]
  + `huaweicloud_vpcep_endpoints` [GH-4025]

ENHANCEMENTS:

* resource/huaweicloud_waf_domain: Support management protection status of domain [GH-3932]
* resource/huaweicloud_as_policy: Add `instance_percentage` parameter [GH-3942]
* resource/huaweicloud_rds_instance: Support updating replication mode [GH-3951]
* resource/huaweicloud_rds_instance: Support database switchover policy [GH-3958]
* resource/huaweicloud_dms_rabbitmq_instance: Support changing password [GH-3956]
* resource/huaweicloud_dc_virtual_gateway: Add tags support [GH-3963]
* resource/huaweicloud_dws_cluster: Support DWS cluster bind ELB [GH-3977]
* resource/huaweicloud_elb_listener: Add `port_ranges` and `gzip_enable` parameters [GH-3991]
* resource/huaweicloud_compute_instance: Support auto terminate time [GH-4015]
* resource/huaweicloud_cts_tracker: Add tags support [GH-4054]
* resource/huaweicloud_compute_instance: Support updating `user_data` parameter [GH-4075]

BUG FIXES:

* resource/huaweicloud_dns_recordset: Fix the issue when creating DNS record on **la-north-2** region [GH-3396]

## 1.60.1 (January 3, 2024)

* **New Resource:**
  + `huaweicloud_sfs_turbo_perm_rule` [GH-3903]
  + `huaweicloud_organizations_delegated_administrator` [GH-3940]

* **New Data Source:**
  + `huaweicloud_cce_cluster_certificate` [GH-3905]
  + `huaweicloud_evs_snapshots` [GH-3907]
  + `huaweicloud_swr_image_triggers` [GH-3909]
  + `huaweicloud_er_flow_logs` [GH-3921]

ENHANCEMENTS:

* resource/huaweicloud_rds_instance: Support maintain window parameter [GH-3892]
* resource/huaweicloud_antiddos_basic: Support alarm configuration [GH-3917]
* resource/huaweicloud_mapreduce_cluster: Support alarm configuration [GH-3931]

## 1.60.0 (December 29, 2023)

* **New Resource:**
  + `huaweicloud_dataarts_studio_data_connection` [GH-3834]
  + `huaweicloud_dataarts_architecture_process` [GH-3832]
  + `huaweicloud_dataarts_architecture_model` [GH-3836]
  + `huaweicloud_dataarts_architecture_data_standard` [GH-3850]
  + `huaweicloud_dataarts_architecture_code_table` [GH-3890]
  + `huaweicloud_dataarts_architecture_table_model` [GH-3894]
  + `huaweicloud_dataarts_architecture_data_standard_template` [GH-3900]
  + `huaweicloud_dataarts_architecture_reviewer` [GH-3902]
  + `huaweicloud_dataarts_factory_job` [GH-3837]
  + `huaweicloud_dataarts_factory_script` [GH-3893]
  + `huaweicloud_dms_kafka_user_client_quota` [GH-3882]
  + `huaweicloud_cce_chart` [GH-3884]
  + `huaweicloud_workspace_eip_associate` [GH-3885]
  + `huaweicloud_cc_central_network_policy_apply` [GH-3886]
  + `huaweicloud_cbr_backup_share_accepter` [GH-3888]
  + `huaweicloud_er_flow_log` [GH-3908]

* **New Data Source:**
  + `huaweicloud_dcs_backups` [GH-3841]
  + `huaweicloud_swr_repositories` [GH-3863]
  + `huaweicloud_dms_kafka_smart_connect_tasks` [GH-3864]
  + `huaweicloud_vpcep_services` [GH-3866]
  + `huaweicloud_nat_private_transit_ips` [GH-3879]
  + `huaweicloud_as_activity_logs` [GH-3891]
  + `huaweicloud_dataarts_architecture_ds_template_optionals` [GH-3901]

ENHANCEMENTS:

* resource/huaweicloud_workspace_desktop: Support migrating enterprise_project_id [GH-3835]
* resource/huaweicloud_vpn_gateway: Support migrating enterprise_project_id [GH-3847]
* resource/huaweicloud_vpn_gateway: Support creating a GM VPN gateway with certificate [GH-3865]
* resource/huaweicloud_vpc: Add `secondary_cidrs` to add more secondary cidrs [GH-3883]

## 1.59.1 (December 20, 2023)

* **New Resource:**
  + `huaweicloud_dms_kafka_smart_connect_task` [GH-3812]
  + `huaweicloud_dataarts_security_data_recognition_rule` [GH-3831]
  + `huaweicloud_rms_organizational_assignment_package` [GH-3839]

* **New Data Source:**
  + `huaweicloud_nat_gateways` [GH-3820]
  + `huaweicloud_swr_organizations` [GH-3849]

ENHANCEMENTS:

* resource/huaweicloud_elb_loadbalancer: Support updating `availability_zone` parameter [GH-3814]
* resource/huaweicloud_workspace_service: Support the OTP auxiliary authentication [GH-3856]
* resource/huaweicloud_networking_secgroup_rule: Allow creating security group rule without remote params [GH-3872]

BUG FIXES:

* resource/huaweicloud_dds_audit_log_policy: Waiting for instance audit log policy to be success [GH-3861]
* resource/huaweicloud_tms_resource: Ignore 404 error while querying resource tags [GH-3874]

## 1.59.0 (December 15, 2023)

* **New Resource:**
  + `huaweicloud_rds_mysql_binlog` [GH-3775]
  + `huaweicloud_dms_kafka_smart_connect` [GH-3784]
  + `huaweicloud_workspace_user_group` [GH-3789]
  + `huaweicloud_cfw_firewall` [GH-3830]

  + `huaweicloud_dataarts_dataservice_app` [GH-3793]
  + `huaweicloud_dataarts_architecture_directory` [GH-3791]
  + `huaweicloud_dataarts_security_permission_set` [GH-3801]
  + `huaweicloud_dataarts_factory_resource` [GH-3802]
  + `huaweicloud_dataarts_architecture_subject` [GH-3811]
  + `huaweicloud_dataarts_architecture_business_metric` [GH-3819]

* **New Data Source:**
  + `huaweicloud_rds_parametergroups` [GH-3783]
  + `huaweicloud_vpn_connection_health_checks` [GH-3790]
  + `huaweicloud_nat_private_gateways` [GH-3798]
  + `huaweicloud_rds_mysql_databases` [GH-3800]
  + `huaweicloud_waf_address_groups` [GH-3806]
  + `huaweicloud_rds_mysql_accounts` [GH-3816]
  + `huaweicloud_waf_dedicated_domains` [GH-3821]
  + `huaweicloud_dataarts_studio_workspaces` [GH-3828]

ENHANCEMENTS:

* resource/huaweicloud_waf_domain: Add `description`, `lb_algorithom`, `website_name` and `forward_header_map` params [GH-3781]
* data/huaweicloud_ccm_private_certificate_export: Support exporting IIS and TOMCAT type certificate [GH-3803]
* resource/huaweicloud_vpn_customer_gateway: Support `certificate_content` parameter [GH-3808]
* resource/huaweicloud_vpc: Support importing `secondary_cidr` attribute [GH-3815]
* resource/huaweicloud_obs_bucket: Support SSE-OBS encryption mode [GH-3825]
* resource/huaweicloud_dms_kafka_user: Support `description` parameter [GH-3826]

BUG FIXES:

* resource/huaweicloud_obs_bucket: Fix the issue when using default KMS key [GH-3804]

## 1.58.0 (December 1, 2023)

* **New Resource:**
  + `huaweicloud_as_planned_task` [GH-3650]
  + `huaweicloud_apig_endpoint_whitelist` [GH-3608]
  + `huaweicloud_cbr_backup_share` [GH-3664]
  + `huaweicloud_cc_central_network` [GH-3766]
  + `huaweicloud_cc_central_network_policy` [GH-3766]
  + `huaweicloud_ccm_private_ca` [GH-3674]
  + `huaweicloud_ccm_private_certificate` [GH-3665]
  + `huaweicloud_ccm_certificate_push` [GH-3769]
  + `huaweicloud_csms_event` [GH-3652]
  + `huaweicloud_dms_kafka_consumer_group` [GH-3756]
  + `huaweicloud_dns_endpoint` [GH-3689]
  + `huaweicloud_dns_resolver_rule` [GH-3714]
  + `huaweicloud_dns_resolver_rule_associate` [GH-3718]
  + `huaweicloud_dns_line_group` [GH-3731]
  + `huaweicloud_kms_dedicated_keystore` [Gh-3661]
  + `huaweicloud_rds_pg_account` [GH-3428]
  + `huaweicloud_rds_pg_database` [GH-3603]
  + `huaweicloud_rds_sqlserver_account` [GH-3598]
  + `huaweicloud_rds_sqlserver_database` [GH-3597]
  + `huaweicloud_rds_sqlserver_database_privilege` [GH-3631]
  + `huaweicloud_rds_pg_plugin` [GH-3695]
  + `huaweicloud_secmaster_indicator` [GH-3575]
  + `huaweicloud_secmaster_alert` [GH-3624]
  + `huaweicloud_secmaster_alert_rule` [GH-3678]
  + `huaweicloud_secmaster_playbook` [GH-3719]
  + `huaweicloud_secmaster_playbook_version` [GH-3728]
  + `huaweicloud_secmaster_playbook_rule` [GH-3749]
  + `huaweicloud_secmaster_playbook_action` [GH-3762]
  + `huaweicloud_sfs_turbo_dir` [GH-3744]
  + `huaweicloud_sfs_turbo_dir_quota` [GH-3760]
  + `huaweicloud_codearts_inspector_website` [GH-3730]
  + `huaweicloud_codearts_inspector_website_scan` [GH-3753]
  + `huaweicloud_workspace_access_policy` [GH-3772]
  + `huaweicloud_workspace_terminal_bindings` [GH-3773]
  + `huaweicloud_workspace_policy_group` [GH-3774]

* **New Data Source:**
  + `huaweicloud_cbr_policies` [GH-3681]
  + `huaweicloud_ccm_private_certificate_export` [GH-3724]
  + `huaweicloud_dbss_flavors` [GH-3638]
  + `huaweicloud_dms_rocketmq_topics` [GH-3667]
  + `huaweicloud_dms_rocketmq_users` [GH-3677]
  + `huaweicloud_dms_rocketmq_consumer_groups` [GH-3708]
  + `huaweicloud_enterprise_projects` [GH-3669]
  + `huaweicloud_er_availability_zones` [GH-3628]
  + `huaweicloud_identity_agencies` [GH-3596]
  + `huaweicloud_elb_listeners` [GH-3586]
  + `huaweicloud_elb_members` [GH-3619]
  + `huaweicloud_elb_security_policies` [GH-3679]
  + `huaweicloud_rds_pg_plugins` [GH-3733]
  + `huaweicloud_vpn_customer_gateways` [GH-3609]
  + `huaweicloud_vpn_connections` [GH-3715]
  + `huaweicloud_workspace_flavors` [GH-3736]

ENHANCEMENTS:

* resource/huaweicloud_cbr_vault: Support multiple AZ for backing [GH-3577]
* resource/huaweicloud_rds_instance: Support **MariaDB** engine type [GH-3673]
* resource/huaweicloud_sfs_turbo: Support updating the sfs turbo name [GH-3771]

## 1.57.0 (October 31, 2023)

* **New Resource:**
  + `huaweicloud_mapreduce_scaling_policy` [GH-3472]
  + `huaweicloud_lts_structuring_custom_configuration` [GH-3480]
  + `huaweicloud_organizations_policy` [GH-3491]
  + `huaweicloud_tms_resource_tags` [GH-3494]
  + `huaweicloud_lts_waf_access` [GH-3517]
  + `huaweicloud_css_configuration` [GH-3529]
  + `huaweicloud_organizations_policy_attach` [GH-3523]
  + `huaweicloud_cbr_checkpoint` [GH-3536]
  + `huaweicloud_smn_logtank` [GH-3567]
  + `huaweicloud_dcs_custom_template` [GH-3583]

  + `huaweicloud_aom_cmdb_application` [GH-3508]
  + `huaweicloud_aom_cmdb_component` [GH-3522]
  + `huaweicloud_aom_cmdb_environment` [GH-3537]

* **New Data Source:**
  + `huaweicloud_mapreduce_clusters` [GH-3488]
  + `huaweicloud_mapreduce_versions` [GH-3509]
  + `huaweicloud_lts_structuring_custom_templates` [GH-3496]
  + `huaweicloud_tms_resource_types` [GH-3499]
  + `huaweicloud_smn_message_templates` [GH-3514]
  + `huaweicloud_organizations_policies` [GH-3525]
  + `huaweicloud_apig_groups` [GH-3534]
  + `huaweicloud_vpn_gateway_availability_zones` [GH-3535]
  + `huaweicloud_vpn_gateways` [GH-3538]
  + `huaweicloud_fgs_functions` [GH-3540]
  + `huaweicloud_rds_sqlserver_collations` [GH-3542]
  + `huaweicloud_elb_loadbalancers` [GH-3544]
  + `huaweicloud_dcs_templates` [GH-3585]
  + `huaweicloud_dcs_template_detail` [GH-3585]

ENHANCEMENTS:

* resource/huaweicloud_compute_instance: Support updating `hostname` [GH-3476]
* resource/huaweicloud_organizations_organization: Support enable service conrtol policy [GH-3481]
* resource/huaweicloud_obs_bucket: Support user domain names [GH-3501]
* resource/huaweicloud_apig_channel: Add `label_key` and `label_value` params [GH-3510]
* resource/huaweicloud_compute_instance: Support creating system disk and data disks in DSS pool  [GH-3546]
* resource/huaweicloud_dcs_instance: Add `parameters` param [GH-3533]
* resource/huaweicloud_waf_certificate: Support updating certificate and private key [GH-3553]
* resource/huaweicloud_vpcep_service: Add `enable_policy` param [GH-3568]
* resource/huaweicloud_cts_tracker: Support the cross account feature [GH-3588]
* resource/huaweicloud_vpc: Support migrating enterprise_project_id [GH-3573]

BUG FIXES:

* resource/huaweicloud_antiddos_basic: Waiting a while for the immediately created EIP [GH-3576]

## 1.56.1 (October 13, 2023)

ENHANCEMENTS:

* resource/huaweicloud_cce_cluster: Support updating `tags` param [GH-3479]
* resource/huaweicloud_smn_subscription: Support more protocols and add `extension` block [GH-3497]

BUG FIXES:

* resource/huaweicloud_cce_addon: Fix a panic issue [GH-3399]

## 1.56.0 (September 28, 2023)

* **New Resource:**
  + `huaweicloud_codearts_deploy_host` [GH-3368]
  + `huaweicloud_codearts_deploy_application` [GH-3420]
  + `huaweicloud_modelarts_network` [GH-3374]
  + `huaweicloud_modelarts_resource_pool` [GH-3395]
  + `huaweicloud_rds_cross_region_backup_strategy` [GH-3417]
  + `huaweicloud_rds_sql_audit` [GH-3419]
  + `huaweicloud_rms_advanced_query` [GH-3430]
  + `huaweicloud_rms_assignment_package` [GH-3446]

  + `huaweicloud_eg_endpoint` [GH-3404]
  + `huaweicloud_eg_connection` [GH-3433]
  + `huaweicloud_eg_custom_event_source` [GH-3411]
  + `huaweicloud_eg_custom_event_channel` [GH-3435]
  + `huaweicloud_eg_event_subscription` [GH-3460]`

  + `huaweicloud_lts_aom_access` [GH-3436]
  + `huaweicloud_lts_notification_template` [GH-3437]
  + `huaweicloud_lts_host_access` [GH-3443]
  + `huaweicloud_lts_search_criteria` [GH-3447]
  + `huaweicloud_lts_sql_alarm_rule` [GH-3448]
  + `huaweicloud_lts_keywords_alarm_rule` [GH-3454]
  + `huaweicloud_lts_structuring_configuration` [GH-3455]

* **New Data Source:**
  + `huaweicloud_organizations_accounts` [GH-3409]
  + `huaweicloud_organizations_organizational_units` [GH-3409]
  + `huaweicloud_cdm_clusters` [GH-3418]
  + `huaweicloud_eg_custom_event_sources` [GH-3421]
  + `huaweicloud_eg_custom_event_channels` [GH-3435]
  + `huaweicloud_rms_assignment_package_templates` [GH-3441]

ENHANCEMENTS:

* resource/huaweicloud_apig_api: Add `passthrough`, `enumeration` and `retry_count` parameters [GH-3369]
* resource/huaweicloud_apig_instance: Add tags support [GH-3394]
* resource/huaweicloud_cbr_vault: Support updating `consistent_level` param [GH-3413]
* resource/huaweicloud_compute_instance: Add `metadata` param [GH-3414]
* resource/huaweicloud_workspace_desktop: Add tags support [GH-3416]
* resource/huaweicloud_fgs_trigger: Support Kafka user_ame and password [GH-3434]
* resource/huaweicloud_fgs_function: Add log config support[GH-3458]

BUG FIXES:

* resource/huaweicloud_identity_group_membership: Add checkDeleted in read function [GH-3387]
* resource/huaweicloud_cci_namespace: Add checkDeleted in read function [GH-3424]

## 1.55.0 (August 31, 2023)

* **New Resource:**
  + `huaweicloud_lts_transfer` [GH-3288]
  + `huaweicloud_modelarts_authorization` [GH-3297]
  + `huaweicloud_apig_appcode` [GH-3318]
  + `huaweicloud_apig_application_authorization` [GH-3343]
  + `huaweicloud_codearts_deploy_group` [GH-3356]

* **New Data Source:**
  + `huaweicloud_account` [GH-3327]
  + `huaweicloud_dns_zones` [GH-3302]
  + `huaweicloud_dns_recordsets` [GH-3308]
  + `huaweicloud_modelarts_model_templates` [GH-3317]
  + `huaweicloud_modelarts_services` [GH-3339]
  + `huaweicloud_modelarts_workspaces` [GH-3347]
  + `huaweicloud_modelarts_resource_flavors` [GH-3359]

ENHANCEMENTS:

* resource/huaweicloud_mapreduce_cluster: Support bootstrap script [GH-3312]
* resource/huaweicloud_identity_provider: Support specify IAM user SSO type [GH-3319]
* resource/huaweicloud_identity_user: Support `external_identity_id` parameter [GH-3321]
* resource/huaweicloud_fgs_function: Support configuring DNS network [GH-3323]
* resource/huaweicloud_waf_policy: Support more options [GH-3326]
* resource/huaweicloud_rds_instance: Support `description` and `dss_pool_id` parameters [GH-3329]
* resource/huaweicloud_rds_instance: Support updating the private IP of RDS instance [GH-3366]
* resource/huaweicloud_rds_mysql_account: Support `description` and `hosts` parameters [GH-3342]
* resource/huaweicloud_lts_group: Add tags support [GH-3346]
* resource/huaweicloud_cce_cluster: Support more parameters [GH-3348]

BUG FIXES:

* resource/huaweicloud_sms_task: Fix panic in when specifying `target_server_disks` block [GH-3325]
* resource/huaweicloud_cbh_instance: Ignore the error when the unbinding eip does not exist [GH-3355]

## 1.54.1 (August 17, 2023)

ENHANCEMENTS:

* resource/huaweicloud_cce_node: Support extension nics param [GH-3303]
* resource/huaweicloud_cce_node: Support dedicated_host_id and initialized_conditions params [GH-3307]

BUG FIXES:

* resource/huaweicloud_obs_bucket_policy: Add checkDeleted in read function [GH-3306]
* resource/huaweicloud_lts_stream: Add checkDeleted in read function [GH-3311]

## 1.54.0 (August 15, 2023)

* **New Resource:**
  + `huaweicloud_modelarts_workspace` [GH-3241]
  + `huaweicloud_waf_rule_information_leakage_prevention` [GH-3278]
  + `huaweicloud_waf_rule_anti_crawler` [GH-3286]
  + `huaweicloud_cfw_eip_protection` [GH-3281]

* **New Data Source:**
  + `huaweicloud_modelarts_models` [GH-3253]
  + `huaweicloud_identity_permissions` [GH-3277]

ENHANCEMENTS:

* resource/huaweicloud_sfs_turbo: Support HPC and HPC cache share type [GH-3244]
* resource/huaweicloud_elb_loadbalancer: Add `backend_subnets`, `protection_status` and `protection_reason` params [GH-3254]
* resource/huaweicloud_elb_listener: Add `forward_port`, `forward_request_port` and `forward_host` params [GH-3259]
* resource/huaweicloud_elb_pool: Add `type`, `protection_status`, `protection_reason` and `slow_start_enabled` params [GH-3274]
* resource/huaweicloud_elb_monitor: Support checking healthy by HTTP or HTTPs status code [GH-3275]
* resource/huaweicloud_elb_l7rule: Add `condition` block param [GH-3279]
* resource/huaweicloud_vpc_bandwidth: Support pre-paid charging mode [GH-3265]
* resource/huaweicloud_compute_instance: Add iops and throughput to disk [GH-3282]
* resource/huaweicloud_cce_node: Support dss in node and node pool [GH-3296]

BUG FIXES:

* resource/huaweicloud_cts_tracker: Create system tracker if it does not exist [GH-3251]
* resource/huaweicloud_rds_instance: Fix panic when availability_zone is an empty list [GH-3258]

## 1.53.0 (July 31, 2023)

* **New Resource:**
  + `huaweicloud_api_gateway_environment` [GH-3186]
  + `huaweicloud_apig_certificate` [GH-3187]
  + `huaweicloud_cc_bandwidth_package` [GH-3192]
  + `huaweicloud_cc_inter_region_bandwidth` [GH-3203]
  + `huaweicloud_ucs_policy` [GH-3126]
  + `huaweicloud_ucs_fleet` [GH-3197]
  + `huaweicloud_ucs_cluster` [GH-3213]
  + `huaweicloud_mapreduce_data_connection` [GH-3215]
  + `huaweicloud_identity_protection_policy` [GH-3219]
  + `huaweicloud_waf_rule_geolocation_access_control` [GH-3226]
  + `huaweicloud_waf_rule_known_attack_source` [GH-3227]
  + `huaweicloud_lts_host_group` [GH-3229]

* **New Data Source:**
  + `huaweicloud_er_attachments` [GH-3208]
  + `huaweicloud_modelarts_service_flavors` [GH-3232]

ENHANCEMENTS:

* resource/huaweicloud_mapreduce_cluster: Support external datasource configurations [GH-3138]
* resource/huaweicloud_ram_resource_share: Support updating principals and resource_urns [GH-3207]
* resource/huaweicloud_cce_node: Add `extend_params` block parameter [GH-3210]
* resource/huaweicloud_elb_l7policy: Support redirect to URL and fixed response [GH-3212]
* resource/huaweicloud_fgs_function: Support version aliases management [GH-3216]
* resource/huaweicloud_evs_volume: Support GPSSD2 and ESSD2 volume type [GH-3217]
* resource/huaweicloud_rds_instance: Support restore to a new instance using backup_id [GH-3218]
* resource/huaweicloud_dms_kafka_instance: Support PLAIN and SCRAM-SHA-512 authentication mechanisms [GH-3234]

## 1.52.1 (July 18, 2023)

ENHANCEMENTS:

* resource/huaweicloud_dcs_instance: Support scaling down the instance [GH-3185]
* resource/huaweicloud_vpc_address_group: Add `enterprise_project_id` support [GH-3191]
* resource/huaweicloud_waf_cloud_instance: Support postpaid charging mode [GH-3195]
* resource/huaweicloud_vpc_subnet: Try to get the DNS list through API [GH-3198]
* resource/huaweicloud_vpn_gateway: Support attach to ER instance[GH-3199]

BUG FIXES:

* resource/huaweicloud_dcs_instance:: Add retry mechanism when restarting [GH-3196]
* resource/huaweicloud_identity_password_policy: Update the endpoint of IAM client with region [GH-3202]

## 1.52.0 (July 14, 2023)

* **New Resource:**
  + `huaweicloud_cnad_advanced_protected_object` [GH-3094]
  + `huaweicloud_cnad_advanced_policy` [GH-3158]
  + `huaweicloud_cnad_advanced_black_white_list` [GH-3160]
  + `huaweicloud_cnad_advanced_policy_associate` [GH-3164]
  + `huaweicloud_modelarts_service` [GH-3151]
  + `huaweicloud_modelarts_model` [GH-3171]
  + `huaweicloud_er_static_route` [GH-3159]

* **New Data Source:**
  + `huaweicloud_cnad_advanced_instances` [GH-3095]
  + `huaweicloud_cnad_advanced_available_objects` [GH-3096]
  + `huaweicloud_cnad_advanced_protected_objects` [GH-3173]
  + `huaweicloud_er_instances` [GH-3152]

ENHANCEMENTS:

* resource/huaweicloud_elb_l7policy: Support redirect the traffic to HTTP listener [GH-3173]
* resource/huaweicloud_dcs_instance: Add `template_id` and support update `password` and `rename_commands` [GH-3176]

BUG FIXES:

* resource/huaweicloud_rds_instance: Add retry mechanism to enable SSL [GH-3190]

## 1.51.0 (June 30, 2023)

* **New Resource:**
  + `huaweicloud_ges_backup` [GH-3112]
  + `huaweicloud_vpc_bandwidth_associate` [GH-3153]
  + `huaweicloud_sdrs_protection_group` [GH-3117]
  + `huaweicloud_sdrs_protected_instance` [GH-3119]
  + `huaweicloud_sdrs_replication_pair` [GH-3121]
  + `huaweicloud_sdrs_drill` [GH-3129]
  + `huaweicloud_sdrs_replication_attach` [GH-3130]
  + `huaweicloud_gaussdb_mysql_parameter_template` [GH-3113]
  + `huaweicloud_gaussdb_mysql_database` [GH-3114]
  + `huaweicloud_gaussdb_mysql_account` [GH-3115]
  + `huaweicloud_gaussdb_mysql_sql_control_rule` [GH-3118]
  + `huaweicloud_gaussdb_mysql_account_privilege` [GH-3122]
  + `huaweicloud_gaussdb_redis_eip_associate` [GH-3125]
  + `huaweicloud_organizations_trusted_service` [GH-3150]
  + `huaweicloud_organizations_account_invite` [GH-3155]
  + `huaweicloud_organizations_account_invite_accepter` [GH-3156]
  + `huaweicloud_organizations_account_associate` [GH-3163]
  + `huaweicloud_organizations_account` [GH-3163]

* **New Data Source:**
  + `huaweicloud_sdrs_domain` [GH-3106]

ENHANCEMENTS:

* resource/huaweicloud_rds_instance: Support specifying the backup cycle in `backup_strategy` block [GH-3059]
* resource/huaweicloud_gaussdb_redis_instance: Add `ssl` and `port` parameters [GH-3100]
* resource/huaweicloud_cce_cluster: Support specifying the default worker node security group ID [GH-3102]
* resource/huaweicloud_evs_volume: Add `dedicated_storage_id` parameter [GH-3107]
* resource/huaweicloud_sfs_turbo: Add `dedicated_storage_id` parameter [GH-3108]
* resource/huaweicloud_dms_rabbitmq_instance: Support updating broker num and storage space [GH-3109]
* resource/huaweicloud_gaussdb_mysql_instance: Support switching SQL filter [GH-3124]
* resource/huaweicloud_cbr_vault: Add `backup_name_prefix` parameter [GH-3139]
* resource/huaweicloud_cbr_policy: Add `enable_acceleration` and `full_back_interval` parameters [GH-3139]
* resource/huaweicloud_mapreduce_cluster: Support updating cluster name [GH-3140]
* resource/huaweicloud_fgs_function: Add tags support [GH-3148]
* resource/huaweicloud_er_instance: Add tags support [GH-3149]

BUG FIXES:

* resource/huaweicloud_identity_agency: Sleep 200ms to avoid API rate limiting [GH-3120]

## 1.50.0 (June 13, 2023)

* **New Resource:**
  + `huaweicloud_organizations_organization` [GH-3076]
  + `huaweicloud_organizations_organizational_unit` [GH-3076]
  + `huaweicloud_ram_resource_share` [GH-3081]
  + `huaweicloud_ges_graph` [GH-3079]
  + `huaweicloud_ges_metadata` [GH-3090]

* **New Data Source:**
  + `huaweicloud_ram_resource_permissions` [GH-3069]
  + `huaweicloud_organizations_organization` [GH-3076]

ENHANCEMENTS:

* resource/huaweicloud_compute_instance: support root volume encryption in some regions [GH-3073]
* resource/huaweicloud_cce_cluster: support setting and updating multi ENI subnet IDs [GH-3077]
* resource/huaweicloud_smn_topic: support policies in topic [GH-3080]
* resource/huaweicloud_cbr_vault: add `policy` parameter to bind all type of policies [GH-3082]
* resource/huaweicloud_images_image: support updating description, min_ram and max_ram fields [GH-3083]

BUG FIXES:

* resource/huaweicloud_identity_agency: fix append project_role issue when creating [GH-3072]
* resource/huaweicloud_images_image: set `backup_id` even if the image is created by instance_id [GH-3084]

## 1.49.0 (May 31, 2023)

* **New Resource:**
  + `huaweicloud_secmaster_incident` [GH-2943]
  + `huaweicloud_live_bucket_authorization` [GH-2992]
  + `huaweicloud_live_snapshot` [GH-2994]
  + `huaweicloud_apig_plugin` [GH-3008]
  + `huaweicloud_apig_plugin_associate` [GH-3013]
  + `huaweicloud_vpn_connection_health_check` [GH-3022]
  + `huaweicloud_obs_bucket_object_acl` [GH-3027]

  + `huaweicloud_aom_alarm_action_rule` [GH-3006]
  + `huaweicloud_aom_alarm_silence_rule` [GH-3011]
  + `huaweicloud_aom_event_alarm_rule` [GH-3016]

  + `huaweicloud_rms_resource_aggregator` [GH-3000]
  + `huaweicloud_rms_resource_aggregation_authorization` [GH-3000]
  + `huaweicloud_rms_resource_recorder` [GH-3029]

  + `huaweicloud_waf_rule_precise_protection` [GH-2984]
  + `huaweicloud_waf_rule_global_protection_whitelist` [GH-2996]
  + `huaweicloud_waf_rule_cc_protection` [GH-3001]

  + `huaweicloud_dws_snapshot` [GH-2981]
  + `huaweicloud_dws_snapshot_policy` [GH-2995]
  + `huaweicloud_dws_event_subscription` [GH-3007]
  + `huaweicloud_dws_alarm_subscription` [GH-3031]
  + `huaweicloud_dws_ext_data_source` [GH-3025]

* **New Data Source:**
  + `huaweicloud_cbr_backup` [GH-2988]

ENHANCEMENTS:

* data/huaweicloud_images_images: export `backup_id` attribute [GH-2979]
* resource/huaweicloud_vpc_peering_connection: support `description` parameter [GH-2987]
* resource/huaweicloud_cce_cluster: support multi container network CIDR [GH-2990]
* resource/huaweicloud_vpcep_service: support `description` parameter [GH-3012]
* resource/huaweicloud_as_group: support updating `agency_name` parameter [GH-3018]
* resource/huaweicloud_images_image: support building whole images [GH-3023]
* resource/huaweicloud_rds_instance: support automatic expansion [GH-3024]
* resource/huaweicloud_cce_addon: support updating version and values [GH-3034]
* resource/huaweicloud_rds_instance: support case-sensitive table names [GH-3037]
* resource/huaweicloud_vpcep_endpoint: support updating whitelist and enable_whitelist [GH-3045]
* resource/huaweicloud_css_cluster: support creating clusters with local dist flavor [GH-3056]
* resource/huaweicloud_identity_agency: support specific days for duration [GH-3058]
* resource/huaweicloud_vpc_subnet: support extra dhcp options [GH-3060]
* resource/huaweicloud_identity_agency: support granting roles to agency on all resources [GH-3062]
* resource/huaweicloud_dms_rabbitmq_instance: support prePaid charging mode [GH-3063]

BUG FIXES:

* resource/huaweicloud_vpn_connection: fix issue caused by `enable_nqa` field [GH-3009]
* data/huaweicloud_modelarts_notebook_images: fix the next page calculation of URL [GH-3026]
* data/huaweicloud_enterprise_project: support exact match by name [GH-3053]

## 1.48.0 (April 28, 2023)

* **New Resource:**
  + `huaweicloud_images_image_share` [GH-2885]
  + `huaweicloud_images_image_share_accepter` [GH-2885]
  + `huaweicloud_apig_acl_policy_associate` [GH-2892]
  + `huaweicloud_dli_template_sql` [GH-2898]
  + `huaweicloud_dds_parameter_template` [GH-2899]
  + `huaweicloud_dds_audit_log_policy` [GH-2900]
  + `huaweicloud_dcs_backup` [GH-2901]
  + `huaweicloud_ces_resource_group` [GH-2902]
  + `huaweicloud_waf_address_group` [GH-2903]
  + `huaweicloud_as_instance_attach` [GH-2904]
  + `huaweicloud_ces_alarm_template` [GH-2905]
  + `huaweicloud_smn_message_template` [GH-2906]
  + `huaweicloud_dli_template_flink` [GH-2909]
  + `huaweicloud_dli_template_spark` [GH-2920]
  + `huaweicloud_dli_global_variable` [GH-2924]
  + `huaweicloud_oms_migration_task_group` [GH-2928]
  + `huaweicloud_dli_agency` [GH-2931]
  + `huaweicloud_apig_signature` [GH-2932]
  + `huaweicloud_identity_group_role_assignment` [GH-2934]
  + `huaweicloud_identity_user_role_assignment` [GH-2934]
  + `huaweicloud_apig_signature_associate` [GH-2942]
  + `huaweicloud_obs_bucket_acl` [GH-2946]
  + `huaweicloud_obs_bucket_replication` [GH-2954]
  + `huaweicloud_swr_image_permissions` [GH-2948]
  + `huaweicloud_swr_image_trigger` [GH-2949]
  + `huaweicloud_swr_image_retention_policy` [GH-2950]
  + `huaweicloud_swr_image_auto_sync` [GH-2951]
  + `huaweicloud_apig_channel` [GH-2958]

ENHANCEMENTS:

* resource/WAF: Add enterprise_project_id support for all of the WAF resources [GH-2939,GH-2955,GH-2964,GH-2965,GH-2967]
* resource/huaweicloud_waf_rule_blacklist: Support `name` and `address_group_id` parameters [GH-2910]
* resource/huaweicloud_gaussdb_opengauss_instance: Support setting `replica_num` to 2 [GH-2936]
* resource/huaweicloud_dds_instance: Support configuration template [GH-2938]
* resource/huaweicloud_as_group: Support `agency_name` parameter [GH-2953]
* resource/huaweicloud_sfs_turbo: Support prePaid charging mode [GH-2961]
* resource/huaweicloud_rds_instance: Support configuration parameter [GH-2962]
* resource/huaweicloud_compute_instance: Support updating `agency_name`, `agent_list` and `enterprise_project_id` [GH-2968,GH-2969]

BUG FIXES:

* resource/huaweicloud_cce_node_attach: Fix the issue that `lvm_config` parameter doesn't work [GH-2811]
* resource/huaweicloud_servicestage_component_instance: Fix the built logic and update param behavior [GH-2915]
* resource/huaweicloud_dli_queue: Fix the wrong query URI when queue_type is general [GH-2941]

DEPRECATE:

* resource/huaweicloud_compute_keypair [GH-2918]
* resource/huaweicloud_identity_role_assignment [GH-2934]
* resource/huaweicloud_apig_vpc_channel` [GH-2958]

## 1.47.1 (April 7, 2023)

ENHANCEMENTS:

* resource/huaweicloud_ces_alarmrule: Support `ALL_INSTANCE` alarm type and add new `resources` param [GH-2896]
* resource/huaweicloud_mapreduce_cluster: Support setting component configurations during cluster creation [GH-2912]

## 1.47.0 (March 31, 2023)

* **New Resource:**
  + `huaweicloud_nat_private_dnat_rule` [GH-2817]
  + `huaweicloud_nat_private_snat_rule` [GH-2821]
  + `huaweicloud_dli_datasource_connection` [GH-2809]
  + `huaweicloud_dli_datasource_auth` [GH-2830]
  + `huaweicloud_dns_custom_line` [GH-2833]
  + `huaweicloud_images_image_copy` [GH-2836]
  + `huaweicloud_cfw_address_group` [GH-2832]
  + `huaweicloud_cfw_address_group_member` [GH-2832]
  + `huaweicloud_cfw_service_group` [GH-2841]
  + `huaweicloud_cfw_service_group_member` [GH-2841]
  + `huaweicloud_cfw_black_white_list` [GH-2846]
  + `huaweicloud_apig_acl_policy` [GH-2859]
  + `huaweicloud_fgs_async_invoke_configuration` [GH-2875]
  + `huaweicloud_apig_instance_routes` [GH-2879]

* **New Data Source:**
  + `huaweicloud_ddm_schemas` [GH-2793]
  + `huaweicloud_ddm_accounts` [GH-2823]
  + `huaweicloud_dms_rabbitmq_flavors` [GH-2834]
  + `huaweicloud_modelarts_notebook_flavors` [GH-2842]

ENHANCEMENTS:

* resource/huaweicloud_as_group: Add `source_dest_check` parameter [GH-2827]
* resource/huaweicloud_dns_recordset: Support multi-line record set configuration [GH-2825]
* resource/huaweicloud_cse_microservice_engine: Add `enterprise_project_id` field [GH-2835]
* resource/huaweicloud_ddm_instance: Support updating `auto_renew` parameter [GH-2853]
* data/huaweicloud_networking_secgroups: Support filter security groups by id [GH-2850]
* data/huaweicloud_cfw_firewalls: Support filter CFW instances by id [GH-2871]

## 1.46.0 (March 17, 2023)

* **New Resource:**
  + `huaweicloud_ddm_schema` [GH-2764]
  + `huaweicloud_as_notification` [GH-2771]
  + `huaweicloud_nat_private_gateway` [GH-2785]
  + `huaweicloud_nat_private_transit_ip` [GH-2814]
  + `huaweicloud_ddm_account` [GH-2794]
  + `huaweicloud_dds_backup` [GH-2798]

* **New Data Source:**
  + `huaweicloud_cph_phone_flavors` [GH-2761]
  + `huaweicloud_cph_phone_images` [GH-2769]
  + `huaweicloud_ddm_instances` [GH-2765]
  + `huaweicloud_ddm_instance_nodes` [GH-2766]

ENHANCEMENTS:

* provider: Support `regional` parameter for ally clouds [GH-2801]
* resource/huaweicloud_lb_loadbalancer: Add `public_ip` attribute [GH-2786]
* resource/huaweicloud_cce_node_pool: Support `storage` configuration [GH-2795]
* resource/huaweicloud_cce_node_pool: Support custom security groups configuration [GH-2803]
* resource/huaweicloud_dms_rocketmq_instance: Support prePaid charging mode [GH-2802]
* resource/huaweicloud_sfs_turbo: Add tags support [GH-2804]
* resource/huaweicloud_dms_rocketmq_instance: Add tags support [GH-2810]

BUG FIXES:

* resource/huaweicloud_vpn_gateway: Support more flavor values [GH-2788]
* resource/huaweicloud_fgs_function: Unset conflicts between `custom_image` and `handler` [GH-2791]

## 1.45.1 (March 10, 2023)

ENHANCEMENTS:

* resource/huaweicloud_dms_kafka_instance: Support prePaid charging mode [GH-2777]
* resource/huaweicloud_ces_alarmrule: Support more `condition` blocks with v2 API [GH-2776]

BUG FIXES:

* data/huaweicloud_gaussdb_opengauss_instances: Fix integer divide by zero dnNum issue [GH-1814]
* resource/huaweicloud_fgs_function: Add CustomImage param when updating function created by mirror [GH-2770]

## 1.45.0 (March 3, 2023)

* **New Resource:**
  + `huaweicloud_kms_grant` [GH-2707]
  + `huaweicloud_vpc_flow_log` [GH-2714]
  + `huaweicloud_rms_policy_assignment` [GH-2721]
  + `huaweicloud_identity_password_policy` [GH-2734]
  + `huaweicloud_cph_server` [GH-2743]
  + `huaweicloud_ddm_instance` [GH-2751]

* **New Data Source:**
  + `huaweicloud_cph_server_flavors` [GH-2716]
  + `huaweicloud_rms_policy_definitions` [GH-2719]
  + `huaweicloud_ddm_engines` [GH-2723]
  + `huaweicloud_as_groups` [GH-2724]
  + `huaweicloud_as_configurations` [GH-2729]
  + `huaweicloud_ddm_flavors` [GH-2736]

ENHANCEMENTS:

* resource/huaweicloud_compute_instance: Support data_disk encryption [GH-2720]
* resource/huaweicloud_compute_instance: Add `description` field [GH-2722]
* resource/huaweicloud_nat_gateway: Add tags support [GH-2726]
* resource/huaweicloud_waf_domain: Support both prePaid and postPaid charging mode [GH-2732]
* resource/huaweicloud_compute_interface_attach: Support security groups parameter when attaching [GH-2739]
* resource/huaweicloud_nat_dnat_rule: Support internal_service_port_range and external_service_port_range [GH-2742]
* resource/huaweicloud_waf_dedicated_domain: Add pci_3ds and pci_dss support [GH-2752]

BUG FIXES:

* resource/huaweicloud_cce_node_pool: Fix prePaid charging mode issue [GH-2741]
* data/huaweicloud_css_flavors: Fix not being able to filter by `vpcs` or `memory` [GH-2750]

## 1.44.2 (February 17, 2023)

ENHANCEMENTS:

* data/huaweicloud_dws_flavors: Support `datastore_type` filter [GH-2652]
* resource/huaweicloud_compute_instance: Add created and updated attributes [GH-2694]
* resource/huaweicloud_waf_dedicated_instance: Add `enterprise_project_id` field [GH-2697]
* resource/huaweicloud_waf_dedicated_domain: Add TLS support [GH-2704]
* resource/huaweicloud_servicestage_component: Add obs config and builder args [GH-2706]

BUG FIXES:

* resource/huaweicloud_dms_rocketmq_instance: Fix an invalid address to set error [GH-2698]
* resource/huaweicloud_compute_servergroup: Add a lock by instance ID when resources are binding or unbinding [GH-2702]

## 1.44.1 (January 20, 2023)

* **New Resource:**
  + `huaweicloud_dms_rocketmq_user` [GH-2576]

* **New Data Source:**
  + `huaweicloud_dms_rocketmq_broker` [GH-2574]
  + `huaweicloud_dms_rocketmq_instances` [GH-2578]

ENHANCEMENTS:

* resource/huaweicloud_dms_rocketmq_instance: Add `enterprise_project_id` field [GH-2631]
* resource/huaweicloud_bucket: Support migrating enterprise_project_id [GH-2638]
* resource/huaweicloud_cce_node_pool: Support ECS group ID when creating [GH-2649]
* resource/huaweicloud_css_cluster: Support enable or disable HTTPS [GH-2664]
* resource/huaweicloud_cce_node: Support updating `key_pair` and `password` [GH-2667]

BUG FIXES:

* resource/huaweicloud_dds_instance: Wait for the instance be ready before update actions [GH-2651]
* resource/huaweicloud_dms_rocketmq_instance: Fix ForceNew issue caused by the order of `availability_zones` [GH-2668]
* resource/huaweicloud_dms_kafka_instance: Avoid overriding the accesses configuration when creating [GH-2669]

## 1.44.0 (December 30, 2022)

* **New Resource:**
  + `huaweicloud_ga_accelerator` [GH-2545]
  + `huaweicloud_ga_listener` [GH-2571]
  + `huaweicloud_ga_endpoint_group` [GH-2584]
  + `huaweicloud_ga_health_check` [GH-2596]
  + `huaweicloud_dms_rocketmq_consumer_group` [GH-2572]
  + `huaweicloud_dms_rocketmq_topic` [GH-2573]
  + `huaweicloud_dc_virtual_gateway` [GH-2575]
  + `huaweicloud_dc_virtual_interface` [GH-2583]
  + `huaweicloud_cfw_protection_rule` [GH-2585]
  + `huaweicloud_dsc_instance` [GH-2587]
  + `huaweicloud_dsc_asset_obs` [GH-2599]
  + `huaweicloud_cbh_instance` [GH-2588]
  + `huaweicloud_hss_host_group` [GH-2590]
  + `huaweicloud_dbss_instance` [GH-2593]
  + `huaweicloud_waf_cloud_instance` [GH-2601]

* **New Data Source:**
  + `huaweicloud_cfw_firewalls` [GH-2585]
  + `huaweicloud_cbh_instances` [GH-2588]
  + `huaweicloud_compute_servergroups` [GH-2591]
  + `huaweicloud_dds_instances` [GH-2594]
  + `huaweicloud_dcs_instances` [GH-2595]

ENHANCEMENTS:

* resource/huaweicloud_servicestage_component_instance: Add `secret_name` for the storage of the secret storge type [GH-2570]
* resource/huaweicloud_cce_node: Support encrypt the root volume with kms key [GH-2602]

## 1.43.0 (November 30, 2022)

* **New Resource:**
  + `huaweicloud_codehub_repository` [GH-2354]
  + `huaweicloud_rds_backup` [GH-2467]
  + `huaweicloud_er_instance` [GH-2471]
  + `huaweicloud_er_vpc_attachment` [GH-2525]
  + `huaweicloud_er_route_table` [GH-2529]
  + `huaweicloud_er_association` [GH-2532]
  + `huaweicloud_er_propagation` [GH-2533]
  + `huaweicloud_elb_logtank` [GH-2520]
  + `huaweicloud_elb_security_policy` [GH-2521]
  + `huaweicloud_vpn_gateway` [GH-2519]
  + `huaweicloud_vpn_customer_gateway` [GH-2528]
  + `huaweicloud_vpn_connection` [GH-2550]
  + `huaweicloud_rf_stack` [GH-2551]
  + `huaweicloud_dms_rocketmq_instance` [GH-2556]

* **New Data Source:**
  + `huaweicloud_rds_backups` [GH-2491]
  + `huaweicloud_rds_storage_types` [GH-2496]
  + `huaweicloud_er_route_tables` [GH-2529]
  + `huaweicloud_images_images` [GH-2536]

ENHANCEMENTS:

* resource/huaweicloud_rds_read_replica_instance: Add prepaid support [GH-2482]
* resource/huaweicloud_vpc_eip: Support updating `auto_renew` parameter [GH-2498]
* resource/huaweicloud_elb_loadbalancer: Support autoscaling option [GH-2493]
* resource/huaweicloud_compute_instance: Support spot price charging mode [GH-2514]
* resource/huaweicloud_rds_instance: Support resetting the password of root user [GH-2515]
* resource/huaweicloud_vpc_address_group: Support IPv6 addresses [GH-2524]
* resource/huaweicloud_cce_node_pool: Add prepaid support [GH-2526]
* resource/huaweicloud_as_configuration: Support more parameters and importing feature [GH-2537]
* resource/huaweicloud_as_group: Support more parameters and IPv6 feature [GH-2539]
* resource/huaweicloud_elb_listener: Support advanced_forwarding [GH-2542]
* resource/huaweicloud_elb_pool: Add `timeout` param of persistence [GH-2546]
* resource/huaweicloud_elb_member: Support cross vpc backend [GH-2548]
* data/huaweicloud_networking_port: Add port allowed IPs attribute [GH-2527]
* data/huaweicloud_compute_instances: Support filtering instances by ID [GH-2555]

BUG FIXES:

* resource/huaweicloud_identity_access_key: Downgrade to warning when saving secret key failed [GH-2461]
* resource/huaweicloud_cse_microservice_engine: Fix the issue of the jobs query fails [GH-2554]

## 1.42.0 (October 31, 2022)

* **New Resource:**
  + `huaweicloud_workspace_desktop` [GH-2420]
  + `huaweicloud_workspace_user` [GH-2426]
  + `huaweicloud_workspace_service` [GH-2453]
  + `huaweicloud_dataarts_studio_instance` [GH-2476]

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_opengauss_instance: Add prepaid support [GH-2406]
* resource/huaweicloud_cce_cluster: Support binding or unbinding EIP without rebuild [GH-2418]
* resource/huaweicloud_cdn_domain: Add tags and range-based interival support [GH-2421] [GH-2425]
* resource/huaweicloud_compute_instance: Support updating keypair without rebuild [GH-2431]
* resource/huaweicloud_servicestage_component: Support `DevCloud` repository [GH-2462]

BUG FIXES:

* resource/huaweicloud_cce_node_attach: Fix issue when attaching a prePaid instance [GH-2395]
* resource/huaweicloud_dms_rabbitmq_instance: Fix crash when creating Rabbitmq instance [GH-2456]

## 1.41.1 (October 15, 2022)

BUG FIXES:

* resource/huaweicloud_cce_node: fix storage spec issue [GH-2432]

## 1.41.0 (September 30, 2022)

* **New Resource:**
  + `huaweicloud_rds_database_privilege` [GH-2374]
  + `huaweicloud_as_bandwidth_policy` [GH-2375]
  + `huaweicloud_gaussdb_influx_instance` [GH-2393]
  + `huaweicloud_gaussdb_mongo_instance` [GH-2393]
  + `huaweicloud_cc_connection` [GH-2407]
  + `huaweicloud_cc_network_instance` [GH-2407]
  + `huaweicloud_dms_kafka_permissions` [GH-2408]
  + `huaweicloud_dms_kafka_user` [GH-2408]
  + `huaweicloud_lts_struct_template` [GH-2411]

* **New Data Source:**
  + `huaweicloud_cdn_domain_statistics` [GH-2352]
  + `huaweicloud_scm_certificates` [GH-2353]

ENHANCEMENTS:

* resource/huaweicloud_cdn_domain: support retrieval host configuration [GH-2355]
* resource/huaweicloud_css_cluster: support master-node and purchase a yearly/monthly cluster [GH-2383]
* resource/huaweicloud_gaussdb_opengauss_instance: support new ha mode for centralized type [GH-2387]
* resource/huaweicloud_gaussdb_opengauss_instance: refactor instance resource [GH-2398]
* resource/huaweicloud_gaussdb_opengauss_instance: support new payment for instance resource [GH-2406]
* resource/huaweicloud_rds_instance: add collation support [GH-2404]
* resource/huaweicloud_fgs_trigger: add LTS trigger type support [GH-2410]
* data/huaweicloud_obs_bucket_object:  support to display content of an object data source with `body` [GH-2368]

BUG FIXES:

* resource/huaweicloud_lts_stream: check the resource deleted when refreshing [GH-2397]

## 1.40.2 (September 16, 2022)

* **New Resource:**
  + `huaweicloud_projectman_project` [GH-2279]

ENHANCEMENTS:

* resource/huaweicloud_cdn_domain: Support IPv6 function [GH-2347]
* resource/huaweicloud_dms_kafka_instance: Support the cross-vpc modification [GH-2360]
* resource/huaweicloud_dds_instance: Add prepaid support [GH-2361]
* resource/huaweicloud_fgs_function: Support create functions by SWR image [GH-2379]

BUG FIXES:

* resource/huaweicloud_dcs_instance: ignore the error of ShowIpWhitelist API when reading [GH-2380]
* resource/huaweicloud_cce_node: fix node can't loggin issue [GH-2381]

## 1.40.1 (August 31, 2022)

* **New Data Source:**
  + `huaweicloud_elb_pools` [GH-2341]

ENHANCEMENTS:

* resource/huaweicloud_compute_instance: Make `auto-renew` updatable [GH-2340]
* resource/huaweicloud_dds_instance: Make `port` updatable [GH-2343]
* resource/huaweicloud_obs_bucket: Support `kms_key_project_id` field [GH-2344]

## 1.40.0 (August 27, 2022)

* **New Resource:**
  + `huaweicloud_aad_forward_rule` [GH-2312]

* **New Data Source:**
  + `huaweicloud_kps_keypairs` [GH-2311]
  + `huaweicloud_sfs_turbos` [GH-2326]
  + `huaweicloud_lb_pools` [GH-2334]
  + `huaweicloud_lb_listener` [GH-2337]

ENHANCEMENTS:

* resource/huaweicloud_servicestage_environment: Support `enterprise_project_id` field [GH-2306]
* resource/huaweicloud_lb_pool: Support `persistence.timeout` field [GH-2324]
* data/huaweicloud_rds_flavors: Support query flavors by `group_type` option [GH-2318]

BUG FIXES:

* resource/huaweicloud_servicestage_component_instance: Fix wrong param name and logic [GH-2303]
* resource/huaweicloud_compute_instance: Upgrade the API to reset password [GH-2317]
* resource/huaweicloud_apig_instance: Ignore the error if ingress address was not found [GH-2332]

## 1.39.0 (August 15, 2022)

* **New Resource:**
  + `huaweicloud_meeting_user` [GH-2249]
  + `huaweicloud_meeting_admin_assignment` [GH-2251]
  + `huaweicloud_dds_database_role` [GH-2270]
  + `huaweicloud_dds_database_user` [GH-2286]

* **New Data Source:**
  + `huaweicloud_identity_users` [GH-2267]

ENHANCEMENTS:

* log: Support setting log level by TF_LOG_PROVIDER environment variable  [GH-2269]
* resource/huaweicloud_cdn_domain: Support configs and cache_settings in cdn domain [GH-2274]
* resource/huaweicloud_servicestage_component_instance: Support log collection policy for host mounting
  and container mounting [GH-2276]
* resource/huaweicloud_networking_secgroup_rule: Support `action` and `priority` parameters [GH-2289]

BUG FIXES:

* resource/huaweicloud_rds_instance: Avoid BACKING_UP state issue [GH-2272]
* resource/huaweicloud_cce_node_pool: Convert the value of subfield in `extend_param` to string [GH-2288]

## 1.38.2 (July 30, 2022)

ENHANCEMENTS:

* resource/huaweicloud_cbr_vault: Add `auto_bind` support [GH-2252]

BUG FIXES:

* schema: fix the problem that may cause the plan error [GH-2248]
* resource/huaweicloud_compute_eip_associate: Waiting for the EIP is associated [GH-2256]
* resource/huaweicloud_dms_rabbitmq_instance: fix the problem that single instance query is empty [GH-2259]

## 1.38.1 (July 15, 2022)

* **New Data Source:**
  + `huaweicloud_smn_topics` [GH-2237]
  + `huaweicloud_dms_kafka_flavors` [GH-2238]

ENHANCEMENTS:

* resource/huaweicloud_dms_kafka_instance: Add `product_id` and `broker_num` support [GH-2238]

BUG FIXES:

* resource/huaweicloud_servicestage_component_instance: Fix index out of bound error [GH-2234]

## 1.38.0 (July 2, 2022)

* **New Resource:**
  + `huaweicloud_iotda_amqp` [GH-2216]
  + `huaweicloud_iotda_device` [GH-2182]
  + `huaweicloud_iotda_device_group` [GH-2206]
  + `huaweicloud_iotda_dataforwarding_rule` [GH-2211]
  + `huaweicloud_iotda_device_linkage_rule` [GH-2225]
  + `huaweicloud_iotda_device_certificate` [GH-2219]
  + `huaweicloud_vod_media_asset` [GH-2201]
  + `huaweicloud_oms_migration_task` [GH-2214]
  + `huaweicloud_meeting_conference` [GH-2218]
  + `huaweicloud_rds_database` [GH-2220]
  + `huaweicloud_apig_throttling_policy_associate` [GGH-2221]

* **New Data Source:**
  + `huaweicloud_apig_environments` [GH-2222]

ENHANCEMENTS:

* resource/huaweicloud_cbr_vault: Add prepaid support [GH-2183]
* resource/huaweicloud_kms_key: Support key rotation [GH-2210]

## 1.37.1 (June 10, 2022)

BUG FIXES:

* resource/huaweicloud_mapreduce_cluster: fix misusing SetNewComputed issue [GH-2186]
* resource/huaweicloud_sfs_file_system: does not delete multiple times during terraform destroy [GH-2189]
* resource/huaweicloud_rds_read_replica_instance: fix "auto_pay" conversion issue [GH-2191]
* resource/huaweicloud_compute_volume_attach: try to delete when status code is 400 [GH-2195]

## 1.37.0 (June 4, 2022)

* **New Resource:**
  + `huaweicloud_live_domain` [GH-2151]
  + `huaweicloud_live_record_callback` [GH-2155]
  + `huaweicloud_live_transcoding` [GH-2158]
  + `huaweicloud_live_recording` [GH-2166]
  + `huaweicloud_vod_media_category` [GH-2164]
  + `huaweicloud_vod_transcoding_template_group` [GH-2173]
  + `huaweicloud_vod_watermark_template` [GH-2176]
  + `huaweicloud_mpc_transcoding_template` [GH-2159]
  + `huaweicloud_mpc_transcoding_template_group` [GH-2159]
  + `huaweicloud_servicestage_component_instance` [GH-2154]
  + `huaweicloud_cse_microservice_engine` [GH-2161]
  + `huaweicloud_cse_microservice` [GH-2169]
  + `huaweicloud_cse_microservice_instance` [GH-2169]
  + `huaweicloud_sms_server_template` [GH-2162]
  + `huaweicloud_sms_task` [GH-2177]
  + `huaweicloud_iotda_space` [GH-2170]
  + `huaweicloud_iotda_product` [GH-2174]

* **New Data Source:**
  + `huaweicloud_dms_kafka_instances` [GH-2144]
  + `huaweicloud_sms_source_servers` [GH-2168]
  + `huaweicloud_identity_projects` [GH-2180]

ENHANCEMENTS:

* resource/huaweicloud_dms_kafka_instance: Add `cross_vpc_accesses` attribute [GH-2141]
* resource/huaweicloud_compute_instance: Support to bind EIP when creating [GH-2153]
* resource/huaweicloud_cse_microservice_engine: Add entrypoint attributes [GH-2167]
* resource/huaweicloud_nat_snat_rule: Add `description` field and make `floating_ip_id` updatable [GH-2171]
* resource/huaweicloud_nat_dnat_rule: Add `description` field [GH-2172]
* resource/huaweicloud_networking_vip_associate: Support import function [GH-2165]
* Support import function for all of ELB resources [GH-2146]
* Support import function for all of dedicated ELB resources [GH-2148]

## 1.36.0 (May 13, 2022)

* **New Resource:**
  + `huaweicloud_servicestage_environment` [GH-2064]
  + `huaweicloud_servicestage_application` [GH-2099]
  + `huaweicloud_servicestage_component` [GH-2133]
  + `huaweicloud_cpts_task` [GH-2109]
  + `huaweicloud_aom_alarm_rule` [GH-2113]
  + `huaweicloud_aom_service_discovery_rule` [GH-2126]
  + `huaweicloud_cts_notification` [GH-2114]
  + `huaweicloud_antiddos_basic` [GH-2120]
  + `huaweicloud_modelarts_dataset_version` [GH-2121]

* **New Data Source:**
  + `huaweicloud_servicestage_component_runtimes` [GH-2115]
  + `huaweicloud_modelarts_dataset_versions` [GH-2129]

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_mysql_instance: Add tags support in special regions [GH-2123]
* resource/huaweicloud_vpc_subnet: Support description field [GH-2134]
* resource/huaweicloud_compute_instance: Support agent_list field [GH-2135]
* resource/huaweicloud_cce_pvc: Support import function [GH-2142]
* resource/huaweicloud_obs_bucket_object: Support import function [GH-2145]

## 1.35.2 (April 24, 2022)

* **New Resource:**
  + `huaweicloud_cpts_project` [GH-2090]
  + `huaweicloud_cts_data_tracker` [GH-2093]
  + `huaweicloud_rds_account` [GH-2095]
  + `huaweicloud_servicestage_repo_password_authorization` [GH-2084]
  + `huaweicloud_servicestage_repo_token_authorization` [GH-2084]

ENHANCEMENTS:

* resource/huaweicloud_as_group: support to forcibly delete an AS group [GH-2086]
* resource/huaweicloud_vpc: support to add a secondary CIDR into a VPC [GH-2088]
* resource/huaweicloud_elb_loadbalancer: support to update flavor_id in prePaid charging mode [GH-2102]

BUG FIXES:

* resource/huaweicloud_compute_volume_attach: add a mutex to make actions serial [GH-2106]
* resource/huaweicloud_compute_eip_associate: ignore errors if the ECS instance was deleted [GH-2111]
* resource/huaweicloud_compute_instance: fix a bug of changing source_dest_check when creating [GH-2116]
* resource/huaweicloud_cce_node: fix a bug of unsubscribing a prepaid cce node [GH-2117]

## 1.35.1 (April 12, 2022)

BUG FIXES:

* resource/huaweicloud_gaussdb_cassandra_instance: fix reducing nodes issue in prepaid mode [GH-2074]
* resource/huaweicloud_networking_secgroup: fix ForceNew issue [GH-2075]

## 1.35.0 (March 31, 2022)

* **New Resource:**
  + `huaweicloud_tms_tags` [GH-2007]
  + `huaweicloud_modelarts_dataset` [GH-2010]
  + `huaweicloud_kps_keypair` [GH-2032]
  + `huaweicloud_cts_tracker` [GH-2066]

* **New Data Source:**
  + `huaweicloud_cci_namespaces` [GH-2004]
  + `huaweicloud_modelarts_datasets` [GH-2017]

ENHANCEMENTS:

* Authentication: Add **shared config profile** support [GH-2025]
* Authentication: Add assume_role support [GH-2040]
* resource/huaweicloud_cce_cluster: Add tags support [GH-2011]
* resource/huaweicloud_elb_loadbalancer: Add prepaid support [GH-2019]
* resource/huaweicloud_fgs_function: Support to specify the functiongraph version [GH-2042]
* resource/huaweicloud_compute_instance: Add `auto_pay` support [GH-2046]
* resource/huaweicloud_vpc_eip: Add `auto_pay` support [GH-2046]
* resource/huaweicloud_vpc_eip_associate: Support to associate an EIP with network_id and fixed_ip [GH-2059]

BUG FIXES:

* resource/huaweicloud_cce_cluster: diff suppressing `cluster_version` [GH-1980]
* data/huaweicloud_dms_product: fix the problem of `vm_specification` parameter [GH-2005]

## 1.34.1 (March 4, 2022)

* **New Resource:**
  + `huaweicloud_vpc_address_group` [GH-1972]
  + `huaweicloud_fgs_dependency` [GH-1999]

* **New Data Source:**
  + `huaweicloud_networking_secgroups` [GH-1992]

ENHANCEMENTS:

* support EPS Authorization for the following resources or data sources:
  + huaweicloud_evs_volumes [GH-1995]
  + huaweicloud_compute_instance [GH-1996]
  + huaweicloud_as [GH-1983]

## 1.34.0 (February 28, 2022)

* **New Resource:**
  + `huaweicloud_drs_job` [GH-1978]
  + `huaweicloud_modelarts_notebook` [GH-1920]
  + `huaweicloud_modelarts_notebook_mount_storage` [GH-1941]

* **New Data Source:**
  + `huaweicloud_modelarts_notebook_images` [GH-1921]

ENHANCEMENTS:

* config: update all versions for a specified customizing service endpoint [GH-1967]
* resource/huaweicloud_fgs_trigger: support apig type for trigger [GH-1982]
* resource/huaweicloud_evs_volume: support prepaid resource creation [GH-1945]
* resource/huaweicloud_cbr_policy: support long-term retention settings [GH-1971]
* support EPS Authorization for the following resources or data sources:
  + huaweicloud_vpc [GH-1958]
  + huaweicloud_vpc_eip [GH-1968]
  + huaweicloud_vpc_bandwidth [GH-1968]
  + huaweicloud_networking_secgroup [GH-1913]
  + huaweicloud_networking_secgroup_rule [GH-1913]
  + huaweicloud_images_image [GH-1943]

BUG FIXES:

* dns: do not update other fields when only tags was changed [GH-1957]
* resource/huaweicloud_compute_instance: fix regexp error for IsIPv4Address [GH-1944]

## 1.33.0 (January 29, 2022)

* **New Resource:**
  + `huaweicloud_csms_secret` [GH-1889]

* **New Data Source:**
  + `huaweicloud_gaussdb_nosql_flavors` [GH-1862]
  + `huaweicloud_css_flavors` [GH-1895]
  + `huaweicloud_csms_secret_version` [GH-1898]

ENHANCEMENTS:

* Authentication: Add ECS metadata API authentication support [GH-1907]
* rename `huaweicloud_networking_eip_associate` to `huaweicloud_vpc_eip_associate` resource [GH-1908]
* resource/huaweicloud_vpc_bandwidth: add `publicips` attribute [GH-1888]
* resource/huaweicloud_vpc_eip: add `port_id` attribute [GH-1916]
* resource/huaweicloud_cce_node: support `password` as plain or salted format [GH-1918]
* resource/huaweicloud_compute_eip_associate: support to associate IPv6 address to a shared bandwidth [GH-1919]
* resource/huaweicloud_gaussdb_cassandra_instance: add `dedicated_resource_name` support [GH-1925]
* resource/huaweicloud_lb_*: update shared ELB resources with v2 API [GH-1900]

BUG FIXES:

* resource/huaweicloud_elb_loadbalancer: release eip when deleting the loadbalancer [GH-1893]
* resource/huaweicloud_nat_dnat_rule: add checkDeleted in read function [GH-1899]

DEPRECATE:

* resource/huaweicloud_dms_group [GH-1878]
* resource/huaweicloud_dms_queue [GH-1880]

## 1.32.1 (January 11, 2022)

ENHANCEMENTS:

* resource/huaweicloud_compute_eip_associate: add `port_id` attribute [GH-1864]
* resource/huaweicloud_lb_certificate: support `enterprise_project_id` field [GH-1865]

BUG FIXES:

* resource/huaweicloud_evs_volume: fix incorrect device_type configuration [GH-1852]
* resource/huaweicloud_vpc_eip: set `port_id` to computed and deprecated [GH-1856]
* resource/huaweicloud_dms_kafka_topic: fix an API issue when updating [GH-1874]
* fix a URL splicing error for customizing endpoint of IAM service [GH-1866]

## 1.32.0 (December 31, 2021)

* **New Resource:**
  + `huaweicloud_cdm_link` [GH-1819]
  + `huaweicloud_cdm_job` [GH-1840]

* **New Data Source:**
  + `huaweicloud_vpc_eips` [GH-1792]
  + `huaweicloud_evs_volumes` [GH-1794]
  + `huaweicloud_rds_instances` [GH-1826]

ENHANCEMENTS:

* resource/huaweicloud_elb_pool: make name and description can be updated to empty [GH-1816]
* resource/huaweicloud_networking_vip: support IPv6 function [GH-1818]
* resource/huaweicloud_vpc_eip: support IPv6 function [GH-1821]
* resource/huaweicloud_compute_instance: support IPv6 function [GH-1834]

BUG FIXES:

* resource/huaweicloud_compute_instance: fix checkdeleted issue [GH-1793]
* resource/huaweicloud_cbr_vaults: fix the resources cannot be removed [GH-1796]
* resource/huaweicloud_vpc_subnet: retry to delete when the error code was 403 [GH-1841]
* resource/huaweicloud_lb_*: don't resty to create resources when an error occurs [GH-1842]

DEPRECATE:

* data/huaweicloud_dms_az [GH-1828]

## 1.31.1 (December 10, 2021)

* **New Resource:**
  + `huaweicloud_dli_flinkjar_job` [GH-1666]
  + `huaweicloud_dli_permission` [GH-1695]
  + `huaweicloud_identity_provider` [GH-1625]
  + `huaweicloud_identity_provider_conversion` [GH-1737]
  + `huaweicloud_waf_instance_group_associate` [GH-1684]

* **New Data Source:**
  + `huaweicloud_cbr_vaults` [GH-1687]
  + `huaweicloud_obs_buckets` [GH-1691]
  + `huaweicloud_iec_bandwidths` [GH-1762]

ENHANCEMENTS:

* resource/huaweicloud_iec_eip: support multi line [GH-1755]
* resource/huaweicloud_fgs_function: add encrypted user data support [GH-1766]
* resource/huaweicloud_dis_stream: add `partitions` attribute [GH-1771]
* resource/huaweicloud_mapreduce_cluster: support `public_ip` parameter [GH-1765]
* resoure/huaweicloud_dms_kafka_instance: support to update storage capacity and bandwidth [GH-1776]
* data/huaweicloud_rds_flavors: add availability_zone filter [GH-1767]

BUG FIXES:

* resource/huaweicloud_rds_configuration: ignore case for `type` [GH-1756]

Removed:

* data/huaweicloud_dis_partition [GH-1768]

## 1.31.0 (November 30, 2021)

* **New Resource:**
  + `huaweicloud_cci_namespace` [GH-1648]
  + `huaweicloud_swr_repository_sharing` [GH-1671]
  + `huaweicloud_enterprise_project` [GH-1731]
  + `huaweicloud_cci_network` [GH-1726]

* **New Data Source:**
  + `huaweicloud_rds_engine_versions` [GH-1729]
  + `huaweicloud_vpcs` [GH-1694]
  + `huaweicloud_vpc_subnets` [GH-1707]

ENHANCEMENTS:

* resource/huaweicloud_gaussdb_mysql_instance: Add configuration_name and dedicated_resource_name [GH-1709]
* resource/huaweicloud_kms_key: Add key_algorithm parameter [GH-1712]
* data/huaweicloud_iec_sites: Add lines attribute [GH-1733]
* resource/huaweicloud_compute_instance: Add source_dest_check parameter [GH-1727]
* resource/huaweicloud_iec_vip: Add ip_address support [GH-1736]
* resource/huaweicloud_identity_user: Add pwd_reset support [GH-1725]

BUG FIXES:

* resource/huaweicloud_evs_volume: Fix update of tag [GH-1705]
* resource/huaweicloud_elb_pool: Add HTTPS and QUIC to protocol [GH-1715]
* resource/huaweicloud_obs_bucket: Ignore FsNotSupport error [GH-1723]

## 1.30.1 (November 27, 2021)

BUG FIXES:

* data/huaweicloud_gaussdb_mysql_instances: Update public_ips type [GH-1740]

## 1.30.0 (October 30, 2021)

* **New Resource:**
  + `huaweicloud_apig_api_publishment` [GH-1595]
  + `huaweicloud_cce_namespace` [GH-1650]
  + `huaweicloud_dli_database` [GH-1607]
  + `huaweicloud_dli_table` [GH-1621]
  + `huaweicloud_dli_package` [GH-1622]
  + `huaweicloud_dli_sql_job` [GH-1623]
  + `huaweicloud_dli_spark_job` [GH-1651]
  + `huaweicloud_dli_flinksql_job` [GH-1659]
  + `huaweicloud_gaussdb_mysql_proxy` [GH-1635]
  + `huaweicloud_swr_organization_permissions` [GH-1575]
  + `huaweicloud_swr_repository` [GH-1658]
  + `huaweicloud_waf_instance_group` [GH-1628]

* **New Data Source:**
  + `huaweicloud_compute_instances` [GH-1645]
  + `huaweicloud_dws_flavors` [GH-1593]
  + `huaweicloud_gaussdb_cassandra_flavors` [GH-1652]
  + `huaweicloud_waf_instance_groups` [GH-1637]

ENHANCEMENTS:

* cce: Support data volume encryption [GH-1616]
* vpc: Add `description` field and deprecate `routes` field [GH-1644]
* resource/huaweicloud_dds_instance: Support import function [GH-1613]
* resource/huaweicloud_gaussdb_redis_instance: Support to update `flavor` and `security_group_id` fields  [GH-1576]
* resource/huaweicloud_waf_dedicated_instance: Add `group_id` parameter [GH-1638]
* resource/huaweicloud_vpc_route_table: Support more than 5 routes when creating [GH-1624]
* resource/huaweicloud_vpc_route: Support `description` and more route types [GH-1619]
* data/huaweicloud_images_image: Add more querying options: name_regex, architecture, os, os_version and image_type [GH-1597]
* data/huaweicloud_gaussdb_mysql_flavors: Add `type` and `az_status` attributes [GH-1614]
* data/huaweicloud_waf_dedicated_instances: Add `group_id` attribute [GH-1640]
* data/huaweicloud_waf_reference_tables: Support to filter data by `name` [GH-1570]

BUG FIXES:

* cce: add missing ForceNew limits [GH-1582]
* data/huaweicloud_vpc_route_table: Retrieve the default route table if name was not specified [GH-1643]

DEPRECATE:

* resource/huaweicloud_images_image_v2 [GH-1604]
* data/huaweicloud_vpc_route_ids [GH-1641]
* data/huaweicloud_vpc_route [GH-1641]

## 1.29.2 (October 9, 2021)

BUG FIXES:

* resource/huaweicloud_dns_recordset: Update to use region client for private zone [GH-1579]

## 1.29.1 (October 6, 2021)

ENHANCEMENTS:

* resource/huaweicloud_cce_pvc: Update with kubernetes API [GH-1572]

## 1.29.0 (September 29, 2021)

* **New Resource:**
  + `huaweicloud_vpc_route_table` [GH-1359]
  + `huaweicloud_fgs_trigger` [GH-1372]
  + `huaweicloud_css_thesaurus` [GH-1513]
  + `huaweicloud_cce_pvc` [GH-1552]

* **New Data Source:**
  + `huaweicloud_gaussdb_opengauss_instances` [GH-1483]
  + `huaweicloud_identity_group` [GH-1485]
  + `huaweicloud_vpc_route_table` [GH-1493]
  + `huaweicloud_dcs_flavors` [GH-1498]
  + `huaweicloud_bms_flavors` [GH-1557]

ENHANCEMENTS:

* IAM: Make `domain_name` optional for IAM resources [GH-1480]
* data/huaweicloud_rds_flavors: Support vcpus and memory arguments [GH-1484]
* resource/huaweicloud_obs_bucket: Support default encryption for a bucket [GH-1518]
* resource/huaweicloud_compute_keypair: Add key file parameter of the private key creation path [GH-1524]
* resource/huaweicloud_dcs_instance: Upgrade DCS API version to V2 [GH-1547]
* resource/huaweicloud_cce_node_attach: Support reset operation [GH-1555]

BUG FIXES:

* resource/huaweicloud_cce_addon: Support json string in values block [GH-1479]
* resource/huaweicloud_vpc_eip: Check whether the eip exists before delete [GH-1522]
* resource/huaweicloud_mapreduce_cluster: fix the type error when handling assigned_roles of node [GH-1526]
* resource/huaweicloud_identity_group_membership: Check whether the user exists before remove [GH-1533]

DEPRECATE:

* data/huaweicloud_dcs_az [GH-1521]
* data/huaweicloud_dcs_product [GH-1521]

## 1.28.1 (September 3, 2021)

ENHANCEMENTS:

* resource/huaweicloud_cdn_domain: Add `service_area` argument support [GH-1466]

BUG FIXES:

* resource/huaweicloud_cce_node: fix an issue when unsubscribing a prePaid node [GH-1464]

## 1.28.0 (August 31, 2021)

FEATURES:

* **New Resource:**
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

* **New Resource:**
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

* **New Resource:**
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

* **New Resource:**
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

* data/huaweicloud_dcs_az: Filter available zones by code ([#990](https://github.com/huaweicloud/terraform-provider-huaweicloud/issues/990))

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

* the `tenant_id` is marked as deprecated in resources ([#952](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/952)[#954](https://github.com/huaweicloud/terraform-provider-huaweicloud/pull/954))

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
