---
subcategory: "CodeArts Build"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_build_template"
description: |-
  Manages a CodeArts Build template resource within HuaweiCloud.
---

# huaweicloud_codearts_build_template

Manages a CodeArts Build template resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_codearts_build_template" "test" {
  name        = "test-api"
  description = "demo"

  steps {
    enable     = true
    module_id  = "devcloud2018.codeci_action_20057.action"
    name       = "update OBS"
    properties = {
      objectKey          = jsonencode("./")
      backetName         = jsonencode("test")
      uploadDirectory    = jsonencode(true)
      artifactSourcePath = jsonencode("bin/*")
      authorizationUser  = jsonencode({
        "displayName": "current user",
        "value": "build" 
      })
      obsHeaders = jsonencode([
        {
          "headerKey": "test",
          "headerValue": "test"
        }
      ])
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the build template.

* `steps` - (Required, List, NonUpdatable) Specifies the build execution steps.
  The [steps](#block--steps) structure is documented below.

* `description` - (Optional, String, NonUpdatable) Specifies the template description.

* `parameters` - (Optional, List, NonUpdatable) Specifies the build execution parameter list.
  The [parameters](#block--parameters) structure is documented below.

* `tool_type` - (Optional, String, NonUpdatable) Specifies the tool type.

<a name="block--steps"></a>
The `steps` block supports:

* `module_id` - (Required, String, NonUpdatable) Specifies the build step module ID.

* `name` - (Required, String, NonUpdatable) Specifies the build step name.

* `enable` - (Optional, Bool, NonUpdatable) Specifies whether to enable the step. Defaults to **false**.

* `properties` - (Optional, Map, NonUpdatable) Specifies the build step properties. Value is JSON format string.

* `version` - (Optional, String, NonUpdatable) Specifies the build step version.

<a name="block--parameters"></a>
The `parameters` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the parameter definition name.

* `params` - (Optional, List, NonUpdatable) Specifies the build execution sub-parameters.
  The [params](#block--parameters--params) structure is documented below.

<a name="block--parameters--params"></a>
The `params` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the parameter field name.

* `value` - (Optional, String, NonUpdatable) Specifies the parameter field value.

* `limits` - (Optional, List, NonUpdatable) Specifies the enumeration parameter restrictions.
  The [limits](#block--parameters--params--limits) structure is documented below.

<a name="block--parameters--params--limits"></a>
The `limits` block supports:

* `disable` - (Optional, String, NonUpdatable) Specifies whether it is effective.

* `display_name` - (Optional, String, NonUpdatable) Specifies the displayed name of the parameter.

* `name` - (Optional, String, NonUpdatable) Specifies the parameter name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - Indicates the template creation time.

* `favorite` - Indicates whether the template is favorite.

* `nick_name` - Indicates the nick name.

* `public` - Indicates whether the template is public.

* `scope` - Indicates the scope.

* `steps` - Indicates the build execution steps.
  The [steps](#attrblock--steps) structure is documented below.

* `template_id` - Indicates ID in database.

* `type` - Indicates the template type.

* `weight` - Indicates the weight of the template.

<a name="attrblock--steps"></a>
The `steps` block supports:

* `properties_all` - Indicates the build step properties.

## Import

The template can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_build_template.test <id>
```
