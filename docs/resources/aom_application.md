---
subcategory: "Application Operations Management (AOM)"
---

# huaweicloud_aom_application

Manages an AOM application resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_aom_application" "test" {
  description            = "application description"
  display_name           = "application_display"
  name                   = "application_demo"
  enterprise_project_id  = "0"
  register_type          = "CONSOLE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of application. The value can contain 2 to 64 characters. Only letters,
  digits, underscores (_), hyphens (-), and periods (.) are allowed.

* `aom_id` - (Optional, String) AOM ID. If you leave this parameter empty, it will not be displayed.

* `description` - (Optional, String) Description which can contain up to 255 characters.

* `display_name` - (Optional, String, ForceNew) Display name, which can contain 2 to 64 characters. Only letters, digits, 
  underscores (_), hyphens (-), and periods (.) are allowed.
  Changing this creates a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the application. This 
  parameter is mandatory for enterprise users.
  Changing this creates a new resource.

* `register_type` - (Optional, String, ForceNew) During frontend invocation, the default value is CONSOLE and no
  parameter needs to be transferred. During REST API invocation, the default value is API. You can change the value to 
  SERVICE_DISCOVERY when necessary.
  Changing this creates a new resource.
  Enumeration values: **API**, **CONSOLE**, **SERVICE_DISCOVERY**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Application ID.

* `create_time` - Creation time.

* `creator` - Creator.

* `modified_time` - Modification time.

* `modifier` - User who makes the modification.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The application operations management can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_application.test <id>
```
