---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_events"
description: |-
  Use this data source to get the list of CCI events within HuaweiCloud.
---

# huaweicloud_cciv2_events

Use this data source to get the list of CCI events within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_events" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace.

* `name` - (Optional, String) Specifies the name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `events` - The events.
  The [events](#events) structure is documented below.

<a name="events"></a>
The `events` block supports:

* `action` - The action.

* `api_version` - The API version.

* `count` - The count.

* `event_time` - The event time.

* `first_timestamp` - The first time.

* `involved_object` - The involved object.
  The [involved_object](#involved_object) structure is documented below.

* `kind` - The kind.

* `last_timestamp` - The last time.

* `message` - The message.

* `metadata` - The metadata.
  The [metadata](#metadata) structure is documented below.

* `reason` - The reason.

* `reporting_component` - The reporting component.

* `reporting_instance` - The reporting instance.

* `type` - The type.

* `source` - The source.
  The [source](#source) structure is documented below.

<a name="involved_object"></a>
The `involved_object` block supports:

* `field_path` - The field path of the involved object.

* `kind` - The kind of the involved object.

* `name` - The name of the involved object.

* `namespace` - The namespace of the involved object.

* `resource_version` - The resource version of the involved object.

* `uid` - The uid of the involved object.

<a name="metadata"></a>
The `metadata` block supports:

* `annotations` - The annotations of the metadata.

* `creation_timestamp` - The creation timestamp of the metadata.

* `name` - The name of the metadata.

* `namespace` - The namespace of the metadata.

* `resource_version` - The resource version of the metadata.

* `uid` - The uid of the metadata.

<a name="source"></a>
The `source` block supports:

* `component` - The component of the source.

* `host` - The host of the source.
