---
layout: "huaweicloud"
page_title: "Provider: HuaweiCloud"
sidebar_current: "docs-huaweicloud-index"
description: |-
  The HuaweiCloud provider is used to interact with the many resources supported by HuaweiCloud. The provider needs to be configured with the proper credentials before it can be used.
---

# HuaweiCloud Provider

The HuaweiCloud provider is used to interact with the
many resources supported by HuaweiCloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the HuaweiCloud Provider
provider "huaweicloud" {
  user_name   = "${var.user_name}"
  password    = "${var.password}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  region      = "cn-north-1"
  auth_url    = "https://iam.cn-north-1.myhwclouds.com:443/v3"
}

# Create a web server
resource "huaweicloud_compute_instance_v2" "test-server" {
  # ...
}
```

## Authentication

This provider offers 4 means for authentication.

- User name + Password
- AKSK
- Token
- Assume Role

### User name + Password

```hcl
provider "huaweicloud" {
  user_name   = "${var.user_name}"
  password    = "${var.password}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  auth_url    = "https://iam.myhwclouds.com:443/v3"
  region      = "RegionOne"
}
```

### AKSK

```hcl
provider "huaweicloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  auth_url    = "https://iam.myhwclouds.com:443/v3"
  region      = "RegionOne"
}
```

### Token

```hcl
provider "huaweicloud" {
  token       = "${var.token}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  auth_url    = "https://iam.myhwclouds.com:443/v3"
  region      = "RegionOne"
}
```

### Assume Role

#### User name + Password

```hcl
provider "huaweicloud" {
  agency_name        = "${var.agency_name}"
  agency_domain_name = "${var.agency_domain_name}"
  delegated_project  = "${var.delegated_project}"
  user_name          = "${var.user_name}"
  password           = "${var.password}"
  domain_name        = "${var.domain_name}"
  auth_url           = "https://iam.myhwclouds.com:443/v3"
  region             = "RegionOne"
}
```

#### AKSK

```hcl
provider "huaweicloud" {
  agency_name        = "${var.agency_name}"
  agency_domain_name = "${var.agency_domain_name}"
  delegated_project  = "${var.delegated_project}"
  access_key         = "${var.access_key}"
  secret_key         = "${var.secret_key}"
  domain_name        = "${var.domain_name}"
  auth_url           = "https://iam.myhwclouds.com:443/v3"
  region             = "RegionOne"
}
```

#### Token

```hcl
provider "huaweicloud" {
  agency_name        = "${var.agency_name}"
  agency_domain_name = "${var.agency_domain_name}"
  delegated_project  = "${var.delegated_project}"
  token              = "${var.token}"
  auth_url           = "https://iam.myhwclouds.com:443/v3"
  region             = "RegionOne"
}
```
```token``` specified is not the normal token, but must have the authority of 'Agent Operator'

## Configuration Reference

The following arguments are supported:

* `access_key` - (Optional) The access key of the HuaweiCloud to use.
  If omitted, the `OS_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the HuaweiCloud to use.
  If omitted, the `OS_SECRET_KEY` environment variable is used.

* `auth_url` - (Required) The Identity authentication URL. If omitted, the
  `OS_AUTH_URL` environment variable is used.

* `region` - (Optional) The region of the HuaweiCloud to use. If omitted,
  the `OS_REGION_NAME` environment variable is used. If `OS_REGION_NAME` is
  not set, then no region will be used. It should be possible to omit the
  region in single-region HuaweiCloud environments, but this behavior may vary
  depending on the HuaweiCloud environment being used.

* `user_name` - (Optional) The Username to login with. If omitted, the
  `OS_USERNAME` environment variable is used.

* `user_id` - (Optional) The User ID to login with. If omitted, the
  `OS_USER_ID` environment variable is used.

* `tenant_id` - (Optional) The ID of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_ID` or
  `OS_PROJECT_ID` environment variables are used.

* `tenant_name` - (Optional) The Name of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_NAME` or
  `OS_PROJECT_NAME` environment variable are used.

* `password` - (Optional) The Password to login with. If omitted, the
  `OS_PASSWORD` environment variable is used.

* `token` - (Optional; Required if not using `user_name` and `password`)
  A token is an expiring, temporary means of access issued via the Keystone
  service. By specifying a token, you do not have to specify a username/password
  combination, since the token was already created by a username/password out of
  band of Terraform. If omitted, the `OS_AUTH_TOKEN` environment variable is used.

* `domain_id` - (Optional) The ID of the Domain to scope to (Identity v3). If
  If omitted, the following environment variables are checked (in this order):
  `OS_USER_DOMAIN_ID`, `OS_PROJECT_DOMAIN_ID`, `OS_DOMAIN_ID`.

* `domain_name` - (Optional) The Name of the Domain to scope to (Identity v3).
  If omitted, the following environment variables are checked (in this order):
  `OS_USER_DOMAIN_NAME`, `OS_PROJECT_DOMAIN_NAME`, `OS_DOMAIN_NAME`,
  `DEFAULT_DOMAIN`.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `OS_INSECURE` environment variable is used.

* `cacert_file` - (Optional) Specify a custom CA certificate when communicating
  over SSL. You can specify either a path to the file or the contents of the
  certificate. If omitted, the `OS_CACERT` environment variable is used.

* `cert` - (Optional) Specify client certificate file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the certificate. If omitted the `OS_CERT` environment variable is used.

* `key` - (Optional) Specify client private key file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the key. If omitted the `OS_KEY` environment variable is used.

* `endpoint_type` - (Optional) Specify which type of endpoint to use from the
  service catalog. It can be set using the OS_ENDPOINT_TYPE environment
  variable. If not set, public endpoints is used.

* `swauth` - (Optional) Set to `true` to authenticate against Swauth, a
  Swift-native authentication system. If omitted, the `OS_SWAUTH` environment
  variable is used. You must also set `username` to the Swauth/Swift username
  such as `username:project`. Set the `password` to the Swauth/Swift key.
  Finally, set `auth_url` as the location of the Swift service. Note that this
  will only work when used with the HuaweiCloud Object Storage resources.

* `use_octavia` - (Optional) If set to `true`, API requests will go the Load Balancer
  service (Octavia) instead of the Networking service (Neutron).

* `agency_name` - (Optional) if authorized by assume role, it must be set. The
  name of agency.

* `agency_domain_name` - (Optional) if authorized by assume role, it must be set.
  The name of domain who created the agency (Identity v3).

* `delegated_project` - (Optional) The name of delegated project (Identity v3).

## Additional Logging

This provider has the ability to log all HTTP requests and responses between
Terraform and the HuaweiCloud cloud which is useful for troubleshooting and
debugging.

To enable these logs, set the `OS_DEBUG` environment variable to `1` along
with the usual `TF_LOG=DEBUG` environment variable:

```shell
$ OS_DEBUG=1 TF_LOG=DEBUG terraform apply
```

If you submit these logs with a bug report, please ensure any sensitive
information has been scrubbed first!

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OS_REGION_NAME` - The region in which to create the server instance.

* `OS_IMAGE_ID` or `OS_IMAGE_NAME` - a UUID or name of an existing image in
    Glance.

* `OS_FLAVOR_ID` or `OS_FLAVOR_NAME` - an ID or name of an existing flavor.

* `OS_POOL_NAME` - The name of a Floating IP pool.

* `OS_NETWORK_ID` - The UUID of a network in your test environment.

* `OS_EXTGW_ID` - The UUID of the external gateway.

You should be able to use any HuaweiCloud environment to develop on as long as the
above environment variables are set.
