---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_group"
description: ""
---

# huaweicloud_apig_group

Manages an APIG (API) group resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}
variable "description" {}
variable "environment_id" {}

resource "huaweicloud_apig_group" "test" {
  instance_id = var.instance_id
  name        = var.group_name
  description = var.description

  environment {
    variable {
      name  = "TERRAFORM"
      value = "/stage/terraform"
    }
    environment_id = var.environment_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the APIG (API) group is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the group belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the group name.  
  The valid length is limited from `3` to `64`, only chinese and english letters, digits and hyphens (-) are
  allowed.  
  The name must start with a chinese or english letter, and the Chinese characters must be in **UTF-8** or **Unicode**
  format.

* `description` - (Optional, String) Specifies the group description.  
  The description contain a maximum of 255 characters and the angle brackets (< and >) are not allowed.  
  Chinese characters must be in **UTF-8** or **Unicode** format.

* `environment` - (Optional, List) Specifies an array of one or more environments of the associated group.  
  The [object](#group_environment) structure is documented below.

 -> The `environment` paramater is conflict with `huaweicloud_apig_environment_variable` resource.

* `url_domains` - (Optional, List) Specifies independent domain names of the associated with group.  
  The [url_domains](#group_url_domains) structure is documented below.

  -> Different groups under the same dedicated instance cannot be bound to the same independent domain name.
     Each API group can be associated with up to `5` domain names.

* `domain_access_enabled` - (Optional, Bool) Specifies whether to use the debugging domain name to access the APIs
  within the group. The default value is `false`.

* `force_destroy` - (Optional, Bool) Specifies whether to delete all sub-resources (for API) when deleting this group.  
  Defaults to **false**.

<a name="group_environment"></a>
The `environment` block supports:

* `variable` - (Required, List) Specifies an array of one or more environment variables.  
  The [object](#group_environment_variable) structure is documented below.

  -> The environment variables of different groups are isolated in the same environment.

* `environment_id` - (Required, String) Specifies the environment ID of the associated group.

<a name="group_environment_variable"></a>
The `variable` block supports:

* `name` - (Required, String) Specifies the variable name.  
  The valid length is limited from `3` to `32` characters.  
  Only letters, digits, hyphens (-), and underscores (_) are allowed, and must start with a letter.  
  In the definition of an API, `name` (case-sensitive) indicates a variable, such as #Name#.
  It is replaced by the actual value when the API is published in an environment.  
  The variable names are not allowed to be repeated for an API group.

* `value` - (Required, String) Specifies the variable value.  
  The valid length is limited from `1` to `255` characters.  
  Only letters, digits and special characters (_-/.:) are allowed.

  -> **NOTE:** The variable value will be displayed in plain text on the console.

<a name="group_url_domains"></a>
The `url_domains` block supports:

* `name` - (Required, String) Specifies the domain name. The valid must comply with the domian name specifications.

* `min_ssl_version` - (Optional, String) Specifies the minimum TLS version that can be used to access the domain name,
  the default value is `TLSv1.2`.
  The valid values are as follows:
  + **TLSv1.1**
  + **TLSv1.2**

  -> This parameter applies only to `HTTPS` and does not take effect for `HTTP` and other access modes.
     Configure `HTTPS` cipher suites using the `ssl_ciphers` parameter on the parameters tab,
     please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-apig/apig_03_0039.html).

* `is_http_redirect_to_https` - (Optional, Bool) Specifies whether to enable redirection from `HTTP` to `HTTPS`.
  The default value is `false`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The group ID.

* `created_at` - The creation time of the group, in RFC3339 format.

* `updated_at` - The latest update time of the group, in RFC3339 format.

* `environment` - The array of one or more environments of the associated group.  
  The [object](#group_environment_attr) structure is documented below.

<a name="group_environment_attr"></a>
The `environment` block supports:

* `variable` - The array of one or more environment variables.  
  The [object](#group_environment_variable_attr) structure is documented below.

<a name="group_environment_variable_attr"></a>
The `variable` block supports:

* `id` - The variable ID.

## Import

API groups can be imported using their `id` and the ID of the related dedicated instance, separated by a slash, e.g.

```shell
$ terraform import huaweicloud_apig_group.test <instance_id>/<id>
```
