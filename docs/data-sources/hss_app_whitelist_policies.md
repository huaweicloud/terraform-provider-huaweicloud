---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_app_whitelist_policies"
description: |-
  Use this data source to get the list of process whitelist policies.
---

# huaweicloud_hss_app_whitelist_policies

Use this data source to get the list of process whitelist policies.

## Example Usage

```hcl
data "huaweicloud_hss_app_whitelist_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_name` - (Optional, String) Specifies the policy name.

* `policy_type` - (Optional, String) Specifies the policy type.
  The value can be **block** (Indicates daily operation mode).

* `learning_status` - (Optional, String) Specifies the policy learning status.
  The valid values are as follows:
  + **effecting**: Learning completed, strategy activated.
  + **learned**: Learning completed, pending confirmation.
  + **learning**: Learning in progress.
  + **pause**: Paused.
  + **abnormal**: Learning anomaly.

* `intercept` - (Optional, String) Specifies whether to enable blocking.
  The valid values are as follows:
  + **true**
  + **false**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of process whitelist policies.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `policy_id` - The policy ID.

* `policy_name` - The policy name.

* `policy_type` - The policy type.

* `learning_status` - The policy learning status.

* `learning_days` - The policy learning days.

* `specified_dir` - Whether to specify a learning directory.

* `dir_list` - The list of monitored directories.

* `file_extension_list` - The list of monitored file extensions.

* `intercept` - Whether to enable blocking.

* `auto_detect` - Whether to automatically enable detection.

* `not_effect_host_num` - The number of hosts where the learning completion strategy is not effective.

* `effect_host_num` - The number of hosts where the learning completion strategy has take effect.

* `trust_num` - The number of trusted processes identified.

* `suspicious_num` - The number of suspicious processes identified.

* `malicious_num` - The number of malicious processes identified.

* `unknown_num` - The number of unknown processes identified.

* `abnormal_info_list` - The list of reasons for learning abnormalities.

  The [abnormal_info_list](#abnormal_info_list_struct) structure is documented below.

* `auto_confirm` - Whether the learning results are automatically confirmed.

* `default_policy` - Whether to default the process whitelist policy.

* `host_id_list` - The list of host IDs.

<a name="abnormal_info_list_struct"></a>
The `abnormal_info_list` block supports:

* `abnormal_type` - The exception type.

* `abnormal_description` - The exception description.
