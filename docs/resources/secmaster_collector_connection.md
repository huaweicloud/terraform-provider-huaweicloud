---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_connection"
description: |-
  Manages a collector connection resource within HuaweiCloud.
---

# huaweicloud_secmaster_collector_connection

Manages a collector connection resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "connection_name" {}
variable "template_id" {}

resource "huaweicloud_secmaster_collector_connection" "test" {
  workspace_id = var.workspace_id
  title        = "test_connection"
  name         = var.connection_name
  template_id  = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `title` - (Required, String) Specifies the title of the collector connection.

* `name` - (Required, String) Specifies the name of the collector connection.

* `template_id` - (Required, String) Specifies the template UUID of the collector connection.

* `description` - (Optional, String) Specifies the description of the collector connection.

* `fields` - (Optional, List) Specifies the list of field configurations.

  The [fields](#fields_struct) structure is documented below.

<a name="fields_struct"></a>
The `fields` block supports:

* `title` - (Required, String) Specifies the field title. It is an input-only field that the API does not return.

* `other` - (Required, String) Specifies other supplementary information. It is an input-only field that the API
  does not return.

* `name` - (Optional, String) Specifies the field name.

* `type` - (Optional, String) Specifies the field type.

* `value` - (Optional, String) Specifies the field value.

* `template_field_id` - (Optional, String) Specifies the template field UUID.

* `field_id` - (Optional, String) Specifies the unique UUID of the field. Only valid when updating an existing resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the connection UUID returned by the API.

* `fields_attribute` - The list of field configurations returned by the API.

* `module_id` - The module UUID of the collector connection.

## Import

The collector connection can be imported using the `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_collector_connection.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to the `fields` not being
imported from the API. It is generally recommended running `terraform plan` after importing this resource.
You can then decide if changes should be applied to the resource, or the definition should be updated to align with the
resource.
