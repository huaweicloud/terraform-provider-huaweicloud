---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_patch_compliance_report_items"
description: |-
  Use this data source to get the list of COC patch compliance report items.
---

# huaweicloud_coc_patch_compliance_report_items

Use this data source to get the list of COC patch compliance report items.

## Example Usage

```hcl
variable "instance_compliant_id" {}

data "huaweicloud_coc_patch_compliance_report_items" "test" {
  instance_compliant_id = var.instance_compliant_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_compliant_id` - (Required, String) Specifies the compliance report ID.

* `title` - (Optional, String) Specifies the patch name.

* `sort_dir` - (Optional, String) Specifies the sorting order.
  Values can be as follows:
  + **asc**: The query results are displayed in ascending order.
  + **desc**: The query results are displayed in the descending order.

* `sort_key` - (Optional, String) Specifies the sorting field.
  Values can be **installed_time**.

* `patch_status` - (Optional, String) Specifies the patch status.
  Values can be as follows:
  + **INSTALLED**: A patch has been installed.
  + **INSTALLED_OTHER**: Other patches have been installed.
  + **MISSING**: A patch is missing.
  + **REJECT**: A patch is rejected.
  + **FAILED**: A patch fails to be installed.
  + **PENDING_REBOOT**: A patch has been installed and is waiting to be restarted.

* `classification` - (Optional, String) Specifies the category.

* `severity_level` - (Optional, String) Specifies the severity level.

* `compliance_level` - (Optional, String) Specifies the compliance level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `compliance_items` - Indicates the patch compliance information.

  The [compliance_items](#compliance_items_struct) structure is documented below.

<a name="compliance_items_struct"></a>
The `compliance_items` block supports:

* `instance_id` - Indicates the node ID.

* `title` - Indicates the patch name.

* `classification` - Indicates the category.

* `severity_level` - Indicates the severity level.

* `compliance_level` - Indicates the compliance level.

* `patch_detail` - Indicates the patch details.

  The [patch_detail](#compliance_items_patch_detail_struct) structure is documented below.

<a name="compliance_items_patch_detail_struct"></a>
The `patch_detail` block supports:

* `installed_time` - Indicates the installation time.

* `patch_baseline_id` - Indicates the patch baseline ID.

* `patch_baseline_name` - Indicates the patch baseline name.

* `patch_status` - Indicates the patch status.
