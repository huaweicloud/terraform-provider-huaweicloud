---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_security_report"
description: |-
  Manages a SecMaster security report resource within HuaweiCloud.
---

# huaweicloud_secmaster_security_report

Manages a SecMaster security report resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "report_name" {}
variable "layout_id" {}
variable "binding_wizard" {}

resource "huaweicloud_secmaster_security_report" "test" {
  workspace_id   = var.workspace_id
  report_name    = var.report_name
  report_period  = "weekly"
  language       = "zh-cn"
  layout_id      = var.layout_id
  binding_wizard = var.binding_wizard

  report_range {
    start = "1662900406226"
    end   = "1662900406226"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.  

* `report_name` - (Required, String, NonUpdatable) Specifies the report name.

* `report_period` - (Required, String, NonUpdatable) Specifies the report period.
  The valid values are **weekly**, **daily**, **annual**, and **monthly**.

* `report_range` - (Required, List, NonUpdatable) Specifies the data range.

  The [report_range](#secmaster_security_report_report_range) structure is documented below.  

* `language` - (Required, String, NonUpdatable) Specifies the language.  

* `layout_id` - (Required, String, NonUpdatable) Specifies the layout ID.

* `binding_wizard` - (Required, String, NonUpdatable) Specifies the report page content.

* `status` - (Optional, String) Specifies the report status.  
  The valid values are **enable** and **disable**.

<a name="secmaster_security_report_report_range"></a>
The `report_range` block supports:

* `start` - (Required, String, NonUpdatable) Specifies the start time of the data range.

* `end` - (Required, String, NonUpdatable) Specifies the end time of the data range.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (the report ID).

## Import

The security report can be imported using the `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_security_report.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `report_range` and `language`.
It is generally recommended running `terraform plan` after importing a security report.
You can then decide if changes should be applied to the security report, or the resource definition should be updated to
align with the security report. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_security_report" "test" {
  ...

  lifecycle {
    ignore_changes = [
      report_range,
      language,
    ]
  }
}
```
