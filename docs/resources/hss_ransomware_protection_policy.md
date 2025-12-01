---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ransomware_protection_policy"
description: |-
  Manages a ransomware protection policy resource within HuaweiCloud HSS.
---

# huaweicloud_hss_ransomware_protection_policy

Manages a ransomware protection policy resource within HuaweiCloud HSS.

## Example Usage

```hcl
resource "huaweicloud_hss_ransomware_protection_policy" "test" {
  policy_name              = "test-policy"
  protection_mode          = "alarm_only"
  protection_directory     = "/root;/home"
  protection_type          = "rtf,doc,txt"
  operating_system         = "Linux"
  enterprise_project_id    = "all_granted_eps"
  deploy_mode              = "opened"
  exclude_directory        = "/boot;/lib;/roo"
  runtime_detection_status = "closed"
  ai_protection_status     = "opened"
  bait_protection_status   = "opened"

  process_whitelist {
    path = "/usr/bin/safe_process"
    hash = "a1b2c3d4e5f6"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the HSS ransomware protection policy belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_name` - (Required, String) Specifies the name of the ransomware protection policy.

* `protection_mode` - (Required, String) Specifies the protection mode of the policy.
  The valid values are as follows:
  + **alarm_and_isolation**: Report an alarm and isolate.
  + **alarm_only**: Only report alarms.

* `protection_directory` - (Required, String) Specifies the directory to be protected.
  Separate multiple directories with semicolons (;). You can configure up to `20` directories.

* `protection_type` - (Required, String) Specifies the protection type. For example, **rtf,doc,txt**.

* `operating_system` - (Required, String, NonUpdatable) Specifies the OS. Its value can be **Windows** or **Linux**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID to which the policy belongs.
  This parameter only needs to be configured after the Enterprise Project feature is enabled.
  For enterprise users, if omitted, default enterprise project will be used.
  Value **0** means default enterprise project.
  Value **all_granted_eps** means all enterprise projects to which the user has been granted access.

* `deploy_mode` - (Optional, String) Specifies whether to enable dynamic honeypots.
  The options are **opened** and **closed**. By default, dynamic honeypot protection is disabled.

* `exclude_directory` - (Optional, String) Specifies the directory to be excluded from protection.
  Separate multiple directories with semicolons (;). You can configure up to `20` directories.

* `process_whitelist` - (Optional, List) Specifies the process whitelist. The value of this field will not be filled back.
  Please check the attribute field `process_whitelist_attribute`.
  The [process_whitelist](#process_whitelist_object) structure is documented below.

* `agent_id_list` - (Optional, List) Specifies the list of agent IDs that have enabled this ransomware protection policy.
  The list can contain a maximum of `1000` items. The value of this field will not be filled back.

* `runtime_detection_status` - (Optional, String) Specifies whether to detect at runtime.
  The valid values are **opened** and **closed**. Defaults to **closed**.

* `ai_protection_status` - (Optional, String) Specifies whether to enable AI ransomware protection.
  Valid values are **opened** and **closed**.

* `bait_protection_status` - (Optional, String) Specifies whether to enable decoy protection.
  Valid value is **opened**.

<a name="process_whitelist_object"></a>
The `process_whitelist` block supports:

* `path` - (Optional, String) Specifies the path of the process to be added to the whitelist.
  The value can contain `0` to `128` characters.

* `hash` - (Optional, String) Specifies the hash value of the process to be added to the whitelist.
  The value can contain `0` to `128` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the policy ID).

* `process_whitelist_attribute` - The list of whitelisted processes with their attributes.
  The [process_whitelist_attribute](#process_whitelist_attribute_object) structure is documented below.

* `count_associated_server` - The number of servers associated with the policy.

* `default_policy` - Whether it is the default policy. Valid values are:
  + `0`: Not the default policy.
  + `1`: The default policy.

<a name="process_whitelist_attribute_object"></a>
The `process_whitelist_attribute` block supports:

* `path` - The path of the process to be added to the whitelist.

* `hash` - The hash value of the process to be added to the whitelist.

## Import

Ransomware protection policy can be imported using the `enterprise_project_id` and `id` separated by a slash, e.g.

### Import resource under the default enterprise project

```bash
$ terraform import huaweicloud_hss_ransomware_protection_policy.test 0/<id>
```

### Import resource from non default enterprise project

```bash
$ terraform import huaweicloud_hss_ransomware_protection_policy.test <enterprise_project_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `enterprise_project_id`,
`process_whitelist`, `agent_id_list`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_ransomware_protection_policy" "test" { 
  # ...

  lifecycle {
    ignore_changes = [
      enterprise_project_id, process_whitelist, agent_id_list,
    ]
  }
}
```
