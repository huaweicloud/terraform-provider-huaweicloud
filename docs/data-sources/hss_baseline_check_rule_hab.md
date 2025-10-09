---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_check_rule_hab"
description: |-
  Use this data source to query the scope of operations for HSS baseline check rules.
---

# huaweicloud_hss_baseline_check_rule_hab

Use this data source to query the scope of operations for HSS baseline check rules.

## Example Usage

```hcl
data "huaweicloud_hss_baseline_check_rule_hab" "test" {
  action        = "ignore"
  handle_status = "unhandled"

  check_rule_list {
    check_name = "example_check"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

* `action` - (Required, String) Specifies the operation to be performed on the baseline check rule.
  Valid values are:
  + **add_to_whitelist**: Add to whitelist.
  + **ignore**: Ignore the check item.
  + **unignore**: Cancel ignoring the check item.
  + **fix**: Fix the check item.
  + **verify**: Verify the check item.

* `handle_status` - (Required, String) Specifies the current status of the baseline check rule.
  Valid values are:
  + **unhandled**: The check item is not processed.
  + **fix-failed**: The check item fails to be fixed.
  + **fixing**: The check item is being fixed.
  + **verifying**: The check item is being verified.
  + **ignored**: The check item is ignored.
  + **safe**: The check item is secure.

* `check_rule_list` - (Required, List) Specifies the list of baseline check rules to be handled.
  The [check_rule_list](#check_rule_list_Struct) structure is documented below.

* `host_id` - (Optional, String) Specifies the ID of the host to be queried.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project. The value **0** indicates
  the default enterprise project. To query all enterprise projects, set this parameter to **all_granted_eps**.
  This field is valid only for enterprise project users.

<a name="check_rule_list_Struct"></a>
The `check_rule_list` block supports:

* `check_name` - (Optional, String) Specifies the name of the baseline check.
  The value can contain up to `256` characters.

* `check_rule_id` - (Optional, String) Specifies the ID of the baseline check rule.
  The value can contain up to `256` characters.

* `standard` - (Optional, String) Specifies the standard type of the baseline check.
  Valid values are:
  + **cn_standard**: Compliance standard.
  + **hw_standard**: Cloud security practice standard.
  + **cis_standard**: General security standard.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_rule_num` - The total number of affected items.

* `rule_num` - The number of check items affected by the operation.

* `host_num` - The number of hosts affected by the operation.

* `data_list` - The detailed information about the operation impact scope.
  The [data_list](#data_list_Struct) structure is documented below.

<a name="data_list_Struct"></a>
The `data_list` block supports:

* `host_id` - The ID of the affected host.

* `host_name` - The name of the affected host.

* `public_ip` - The public IP address of the host.

* `private_ip` - The private IP address of the host.

* `asset_value` - The asset value of the host. Valid values are:
  + **important**: Important asset.
  + **common**: Common asset.
  + **test**: Test asset.

* `check_type` - The name of the baseline check.

* `standard` - The standard type of the host. Valid values are:
  + **cn_standard**: Compliance standard.
  + **hw_standard**: Cloud security practice standard.
  + **cis_standard**: General security standard.

* `tag` - The check type of check items in the baseline check.

* `check_rule_name` - The name of check items in the baseline check.
