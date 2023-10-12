---
subcategory: "Application Operations Management (AOM)"
---

# huaweicloud_aom_application

Manages an AOM application resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_aom_application" "test" {
  name                  = "app_demo"
  description           = "application description"
  enterprise_project_id = "0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name or *Unique Identifier* of the application. The value can contain
  2 to 64 characters. Only letters, digits, underscores (_), hyphens (-), and periods (.) are allowed.

* `display_name` - (Optional, String) Specifies the **display** name of the application. The value can contain
  2 to 64 characters. Only letters, digits, underscores (_), hyphens (-), and periods (.) are allowed.
  If not specified, it equals the value of `name` during creation, while it should be explicitly specified during modification.

* `description` - (Optional, String) Specifies the description about the application.
  The description can contain a maximum of 255 characters.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the application.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
* `register_type` - The register type of the application.
* `create_time` - The creation time.
* `modified_time` - The last modification time.
* `creator` - The creator of the application.
* `modifier` - The last modifier of the application.

## Import

The AOM application can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_application.test d61ef1ddb07f40e381ee37a000512caa
```
