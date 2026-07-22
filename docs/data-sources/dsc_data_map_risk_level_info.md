---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_data_map_risk_level_info"
description: |-
  Use this data source to get the asset statistics of each risk level for the left menu bar within HuaweiCloud.
---

# huaweicloud_dsc_data_map_risk_level_info

Use this data source to get the asset statistics of each risk level for the left menu bar within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_data_map_risk_level_info" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `security_group_ids` - (Optional, List) Specifies the security group IDs used to filter assets.

* `security_level_ids` - (Optional, List) Specifies the sensitive level IDs used to filter assets.

* `type` - (Optional, String) Specifies the asset type used to filter assets.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_level` - The statistics data list of each risk level.

  The [data_level](#data_level_struct) structure is documented below.

<a name="data_level_struct"></a>
The `data_level` block supports:

* `level_color_number` - The front-end display color number of the level.

* `level_id` - The unique ID of the sensitive risk level.

* `level_name` - The name of the level.

* `total` - The total number of assets under the current level.

* `risk_list` - The risk asset groups under the current level.

  The [risk_list](#risk_list_struct) structure is documented below.

<a name="risk_list_struct"></a>
The `risk_list` block supports:

* `total` - The number of assets under the current risk type.

* `type` - The risk classification identifier.

* `detail_list` - The asset detail list in the group.

  The [detail_list](#detail_list_struct) structure is documented below.

<a name="detail_list_struct"></a>
The `detail_list` block supports:

* `asset_name` - The name of the asset.

* `asset_type` - The major category of the asset.

* `create_time` - The creation timestamp of the asset.

* `id` - The unique ID of the asset.

* `ins_type` - The subdivided instance type of the asset.

* `is_authorized` - Whether the asset has completed authorization.

* `is_scaned` - Whether the sensitive scan is executed.

* `vpc_id` - The VPC identifier to which the asset belongs.

* `scan_detail` - The sensitive scan result detail.

  The [scan_detail](#scan_detail_struct) structure is documented below.

* `security_strategy` - The security configuration strategy of the asset.

  The [security_strategy](#security_strategy_struct) structure is documented below.

* `threat_analysis` - The threat risk analysis of the asset.

  The [threat_analysis](#threat_analysis_struct) structure is documented below.

<a name="scan_detail_struct"></a>
The `scan_detail` block supports:

* `job_id` - The ID of the latest scan job.

* `last_scan_time` - The timestamp of the last scan.

* `object_num` - The total number of scan objects.

* `scan_risk` - The comprehensive scan risk level.

* `scan_template_id` - The ID of the scan template.

* `scan_template_name` - The name of the scan template.

* `security_level_color` - The color number corresponding to the level.

* `security_level_id` - The ID of the sensitive level.

* `security_level_name` - The name of the sensitive level.

* `sensitive_object_num` - The number of identified sensitive objects.

<a name="security_strategy_struct"></a>
The `security_strategy` block supports:

* `ssl_enabled` - The SSL switch status.

* `access_strategy` - The access control strategy.

* `access_strategy_level` - The risk level of the access strategy.

* `authority_enable` - The status of the fine-grained permission control.

* `authority_level` - The permission control level.

* `backup_and_restore` - The type of the backup and restore scheme.

* `backup_enable` - The automatic backup switch.

* `backup_level` - The backup security level.

* `data_volume_encrypt_enable` - The data volume encryption switch.

* `data_volume_encrypt_level` - The data volume encryption level.

* `dbss_audit_security_level` - The DBSS database audit level.

* `dbss_audit_status` - The DBSS audit enabled status.

* `disk_encrypted` - The overall disk encryption status.

* `disk_encrypted_enable` - The disk encryption switch.

* `encrypt_level` - The overall data encryption level.

* `https_enable` - The forced HTTPS switch.

* `https_level` - The HTTPS security level.

* `ik_enable` - The IK key management enabled status.

* `is_encrypt` - Whether the asset storage is encrypted.

* `obs_audit_level` - The OBS object storage audit level.

* `obs_audit_status` - The OBS audit switch.

* `public_network_access` - The public network access control strategy.

* `public_network_enable` - Whether the public network access is opened.

* `security_group_level` - The security group strategy level.

* `ssl_status` - The SSL certificate enabled status.

* `total_risk` - The comprehensive risk level of the security strategy.

* `working_mode` - The running mode of the asset.

* `working_type` - The business purpose type of the asset.

<a name="threat_analysis_struct"></a>
The `threat_analysis` block supports:

* `abnormal_login_level` - The abnormal login risk level.

* `risky_operation_level` - The high-risk operation risk level.

* `sql_inject_level` - The SQL injection attack risk level.
