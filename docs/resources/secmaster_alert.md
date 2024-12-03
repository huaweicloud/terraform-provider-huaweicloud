---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert"
description: |-
  Manages a SecMaster alert resource within HuaweiCloud.
---

# huaweicloud_secmaster_alert

Manages a SecMaster alert resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

resource "huaweicloud_secmaster_alert" "test" {
  workspace_id = var.workspace_id
  name         = "test"
  description  = "created by terraform"

  type {
    category   = "Abnormal network behavior"
    alert_type = "Abnormal access frequency of IP address"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  first_occurrence_time = "2023-10-26T09:33:55.000+08:00"

  severity            = "Tips"
  status              = "Open"
  verification_status = "Unknown"
  stage               = "Preparation"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the alert belongs.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the alert name.

* `description` - (Required, String) Specifies the description of the alert.

* `type` - (Required, List) Specifies the alert type configuration.
  The [type](#Alert_AlertType) structure is documented below.

* `data_source` - (Required, List, ForceNew) Specifies the data source configuration.
  The [data_source](#Alert_DataSource) structure is documented below.

  Changing this parameter will create a new resource.

* `severity` - (Required, String) Specifies the alert severity.
  The value can be: **Tips**, **Low**, **Medium**, **High** and **Fatal**.

* `status` - (Required, String) Specifies the alert status.
  The value can be: **Open**, **Block** and **Closed**.

* `stage` - (Required, String) Specifies the alert stage.
  The value can be **Preparation**, **Detection and Analysis**, **Containm,Eradication& Recovery**
  and **Post-Incident-Activity**.

* `verification_status` - (Required, String) Specifies the alert verification status.
  The value can be: **Unknown**, **Positive** and **False positive**.

* `first_occurrence_time` - (Required, String) Specifies the first occurrence time of the indicator.

* `last_occurrence_time` - (Optional, String) Specifies the last occurrence time of the indicator.

* `owner` - (Optional, String) Specifies the owner name of the alert.

* `debugging_data` - (Optional, Bool, ForceNew) Specifies whether it's a debugging data.

  Changing this parameter will create a new resource.

* `labels` - (Optional, String) Specifies the labels of the alert in comma-separated string.

* `close_reason` - (Optional, String) Specifies the close reason.
  The value can be **False detection**, **Resolved**, **Repeated** and **Other**.

* `close_comment` - (Optional, String) Specifies the close comment.

<a name="Alert_AlertType"></a>
The `type` block supports:

* `category` - (Required, String) Specifies the category.

* `alert_type` - (Required, String) Specifies the alert type.

<a name="Alert_DataSource"></a>
The `data_source` block supports:

* `product_feature` - (Required, String, ForceNew) Specifies the product feature.

  Changing this parameter will create a new resource.

* `product_name` - (Required, String, ForceNew) Specifies the product name.

  Changing this parameter will create a new resource.

* `source_type` - (Required, Int, ForceNew) Specifies the source type.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The created time.

* `updated_at` - The updated time.

## Import

The indicator can be imported using the workspace ID and the alert ID, e.g.

```bash
$ terraform import huaweicloud_secmaster_alert.test <workspace_id>/<id>
```
