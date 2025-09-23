---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_incident"
description: |-
  Manages a SecMaster incident resource within HuaweiCloud.
---

# huaweicloud_secmaster_incident

Manages a SecMaster incident resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

resource "huaweicloud_secmaster_incident" "test" {
  workspace_id = var.workspace_id
  name         = "test-incident"
  description  = "test incident"

  type {
    category      = "DDoS"
    incident_type = "ACK Flood"
  }

  level  = "Tips"
  status = "Open"
  owner  = "test-user"

  data_source {
    product_feature = "hss"
    product_name    = "hss"
    source_type     = 1
  }

  first_occurrence_time = "2023-04-18T13:00:00.000+08:00"
  last_occurrence_time  = "2023-04-19T14:00:00.000+08:00"
  planned_closure_time  = "2023-04-20T15:00:00.000+08:00"
  verification_status   = "Unknown"
  stage                 = "Preparation"
  debugging_data        = false
  labels                = "test1,test2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the incident belongs.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the incident name.

* `description` - (Required, String) Specifies the incident description.
  The value contains a maximum of 1024 characters.

* `type` - (Required, List, ForceNew) Specifies the incident type configuration.
  The [IncidentType](#Incident_IncidentType) structure is documented below.

  Changing this parameter will create a new resource.

* `level` - (Required, String) Specifies the incident level.
  The value can be **Tips**, **Low**, **Medium**, **High** and **Fatal**.

* `status` - (Required, String) Specifies the incident status.
  The value can be **Open**, **Block** and **Closed**.

* `data_source` - (Required, List, ForceNew) Specifies the data source configuration.
  The [IncidentDataSource](#IncidentDataSource) structure is documented below.

  Changing this parameter will create a new resource.

* `first_occurrence_time` - (Required, String) Specifies the first occurrence time of the incident.
  For example: 2023-04-18T13:00:00.000+08:00

* `owner` - (Optional, String) Specifies the user name of the owner.

* `last_occurrence_time` - (Optional, String) Specifies the last occurrence time of the incident.
  For example: 2023-04-18T13:00:00.000+08:00

* `planned_closure_time` - (Optional, String) Specifies the planned closure time of the incident.
  For example: 2023-04-18T13:00:00.000+08:00

* `verification_status` - (Optional, String) Specifies the verification status.
  The value can be **Unknown**, **True_Positive**, **False_Positive**, defaults to **Unknown**.

* `stage` - (Optional, String) Specifies the stage of the incident.
  The value can be **Preparation**, **Detection and Analysis**, **Containm,Eradication& Recovery** and **Post-Incident-Activity**.
  Defaults to **Preparation**.

* `debugging_data` - (Optional, Bool) Specifies whether it's a debugging data.

* `labels` - (Optional, String) Specifies the labels, separated by comma(,).

* `close_reason` - (Optional, String) Specifies the close reason.
  The value can be **False detection**, **Resolved**, **Repeated** and **Other**.

* `close_comment` - (Optional, String) Specifies the close comment.

<a name="Incident_IncidentType"></a>
The `IncidentType` block supports:

* `category` - (Required, String, ForceNew) Specifies the category.

  Changing this parameter will create a new resource.

* `incident_type` - (Required, String, ForceNew) Specifies the incident type.

  Changing this parameter will create a new resource.

<a name="IncidentDataSource"></a>
The `IncidentDataSource` block supports:

* `product_feature` - (Required, String, ForceNew) Specifies the product feature.

  Changing this parameter will create a new resource.

* `product_name` - (Required, String, ForceNew) Specifies the product name.

  Changing this parameter will create a new resource.

* `source_type` - (Required, Int, ForceNew) Specifies the source type.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `creator` - The creator name.

* `created_at` - The created time.

* `updated_at` - The updated time.

## Import

The incident can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_secmaster_incident.test 40b57838-2019-443a-bb07-30a7a50a4780
```
