---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_timer_rule"
description: |-
  Use this resource to manage a timer rule for starting and stopping components within HuaweiCloud.
---

# huaweicloud_cae_timer_rule

Use this resource to manage a timer rule for starting and stopping components within HuaweiCloud.

## Example Usage

```hcl
variable "environment_id" {}
variable "rule_name" {}
variable "cron" {}
variable "component_configurations" {
  type = list(object({
    id   = string
    name = string
  }))
}

resource "huaweicloud_cae_timer_rule" "test" {
  environment_id   = var.environment_id
  name             = var.rule_name
  type             = "start"
  status           = "on"
  effective_range  = "component"
  effective_policy = "onetime"
  cron             =  var.cron

  dynamic "components" {
    for_each = var.component_configurations
    content {
      id   = components.value["id"]
      name = components.value["name"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the CAE environment.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the timer rule.  
  The maximum length of the name is `64` characters, only letters, digits, underscores (_) and hyphens (-)
  are allowed.  
  The name must start and end with a letter or a digit.

* `type` - (Required, String) Specifies the type of the timer rule.  
  The valid values are as follows:
  + **stop**: The components will be started in batches. The components that have been started are not affected.
  + **start**: The components will be stopped in batches. The components that have been stopped are not affected.

* `status` - (Required, String) Specifies the status of the timer rule.  
  The valid values are as follows:
  + **on**
  + **off**

* `effective_range` - (Required, String) Specifies the effective range of the timer rule.  
  The valid values are as follows:
  + **environment**: The rule takes effect for all components in the environment.
  + **application**: The rule takes effect for all components in the application.
  + **component**: The rule takes effect for the specified components.

* `effective_policy` - (Required, String) Specifies the effective policy of the timer rule.  
  The valid values are as follows:
  + **onetime**: The rule is executed only once.
  + **periodic**: The rule is executed periodically.

* `cron` - (Required, String) Specifies the cron expression of the timer rule.  
  The triggered time of the rule must be at least two minutes later than the current time.  
  When `effective_policy` is set to **periodic**, the rule can only be executed by week of day.

* `applications` - (Optional, List) Specifies the list of the applications in which the timer rule takes effect.  
  The [applications](#timer_rule_applications) structure is documented below.  
  This parameter is required and available only when the `effective_range` parameter is set to **application**.

* `components` - (Optional, List) Specifies the list of the components in which the timer rule takes effect.  
  The [components](#timer_rule_components) structure is documented below.  
  This parameter is required and available only when the `effective_range` parameter is set to **component**.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  timer rule belongs.  
  Changing this creates a new resource.

-> If the `environment_id` belongs to the non-default enterprise project, this parameter is required and is
   only valid for enterprise users.

<a name="timer_rule_applications"></a>
The `applications` block supports:

* `id` - (Required, String) Specifies the ID of the application.

* `name` - (Optional, String) Specifies the name of the application.

<a name="timer_rule_components"></a>
The `components` block supports:

* `id` - (Required, String) Specifies the ID of the component.

* `name` - (Optional, String) Specifies the name of the component.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The resource can be imported using `environment_id` and `name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_cae_timer_rule.test <environment_id>/<name>
```

For the timer rule with the non-default enterprise project ID, its enterprise project ID need to be specified
additionanlly when importing. All fields are separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_cae_timer_rule.test <environment_id>/<name>/<enterprise_project_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `status`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cae_timer_rule" "test" {
  ...

  lifecycle {
    ignore_changes = [
      status,
    ]
  }
}
```
