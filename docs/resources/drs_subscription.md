---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_subscription"
description: |-
  Use this resource to create a DRS subscription task within HuaweiCloud.
---

# huaweicloud_drs_subscription

Use this resource to create a DRS subscription task within HuaweiCloud.

-> This resource is a one-time action resource for creating a DRS subscription task. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "source_instance_id" {}

resource "huaweicloud_drs_subscription" "test" {
  name                  = "DRS-123"
  enterprise_project_id = "0"

  source_endpoint_info {
    id = var.source_instance_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the subscription task.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `name` - (Required, String, NonUpdatable) Specifies the name of the subscription task.  
  The name must be between 4 and 50 characters, and can contain letters, digits, hyphens or underscores.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the subscription task.

* `instance_type` - (Optional, String, NonUpdatable) Specifies the instance type of the subscription task.  
  Defaults to **rds**.

* `tags` - (Optional, List, NonUpdatable) Specifies the tags of the subscription task.  
  The [tags](#drs_subscription_tags) structure is documented below.

* `source_endpoint_info` - (Required, List, NonUpdatable) Specifies the source database information of the
  subscription task.  
  The [source_endpoint_info](#drs_subscription_source_endpoint_info) structure is documented below.

* `is_grant_new_agency` - (Optional, Bool, NonUpdatable) Specifies whether to create a new agency.

<a name="drs_subscription_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.

* `value` - (Required, String) Specifies the tag value.

<a name="drs_subscription_source_endpoint_info"></a>
The `source_endpoint_info` block supports:

* `id` - (Required, String) Specifies the database instance ID.

* `type` - (Optional, String) Specifies the database instance type.  
  Defaults to **mysql**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the job ID returned by the API.
