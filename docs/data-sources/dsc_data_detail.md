---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_data_detail"
description: |-
  Use this data source to get the detail of a single asset in the asset map within HuaweiCloud.
---

# huaweicloud_dsc_data_detail

Use this data source to get the detail of a single asset in the asset map within HuaweiCloud.

## Example Usage

```hcl
variable "ins_id" {}

data "huaweicloud_dsc_data_detail" "test" {
  ins_id = var.ins_id
  type   = "DATABASE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `ins_id` - (Required, String) Specifies the asset ID used to identify a specific asset.

* `type` - (Required, String) Specifies the asset type used to specify the type of the detail to return.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `asset_name` - The name of the asset.

* `asset_type` - The type of the asset.

* `create_time` - The creation time of the asset.

* `ins_type` - The instance type of the asset.

* `is_authorized` - Whether the asset is authorized.

* `is_scaned` - Whether the asset is scanned.

* `scan_detail` - The sensitive scan result detail.

  The [scan_detail](#scan_detail_struct) structure is documented below.

* `security_strategy` - The security configuration strategy of the asset.

  The [security_strategy](#security_strategy_struct) structure is documented below.

* `threat_analysis` - The threat risk analysis of the asset.

  The [threat_analysis](#threat_analysis_struct) structure is documented below.

* `vpc_id` - The VPC identifier to which the asset belongs.

<a name="scan_detail_struct"></a>
The `scan_detail` block supports:

* `job_id` - The ID of the latest scan job.

* `last_scan_time` - The latest scan time.

* `object_num` - The number of scanned objects.

* `scan_risk` - The scan risk level.

* `scan_template_id` - The ID of the scan template.

* `scan_template_name` - The name of the scan template.

* `security_level_color` - The front-end display color of the security level.

* `security_level_id` - The ID of the security level.

* `security_level_name` - The name of the security level.

* `sensitive_object_num` - The number of identified sensitive objects.

<a name="security_strategy_struct"></a>
The `security_strategy` block supports:

* `ssl_enabled` - The SSL switch status.

* `access_strategy` - The access strategy.

* `access_strategy_level` - The access strategy level.

* `authority_enable` - Whether the authority control is enabled.

* `authority_level` - The authority level.

* `backup_and_restore` - The backup and restore strategy.

* `backup_enable` - Whether the backup is enabled.

* `backup_level` - The backup level.

* `data_volume_encrypt_enable` - Whether the data volume encryption is enabled.

* `data_volume_encrypt_level` - The data volume encryption level.

* `dbss_audit_security_level` - The DBSS audit security level.

* `dbss_audit_status` - The DBSS audit status.

* `disk_encrypted` - The disk encryption status.

* `disk_encrypted_enable` - Whether the disk encryption is enabled.

* `encrypt_level` - The encryption level.

* `https_enable` - Whether the HTTPS is enabled.

* `https_level` - The HTTPS level.

* `ik_enable` - The IK switch status.

* `is_encrypt` - Whether the asset is encrypted.

* `obs_audit_level` - The OBS audit level.

* `obs_audit_status` - The OBS audit status.

* `public_network_access` - The public network access strategy.

* `public_network_enable` - Whether the public network access is enabled.

* `security_group_level` - The security group level.

* `ssl_status` - The SSL status.

* `total_risk` - The overall risk level.

* `working_mode` - The working mode.

* `working_type` - The business purpose type of the asset.

<a name="threat_analysis_struct"></a>
The `threat_analysis` block supports:

* `abnormal_login_level` - The abnormal login risk level.

* `risky_operation_level` - The high risk operation level.

* `sql_inject_level` - The SQL injection risk level.
