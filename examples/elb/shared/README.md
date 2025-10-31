# Create a Basic ELB Load Balancer

This example provides best practice code for using Terraform to create a complete ELB (Elastic Load Balancer)
environment in HuaweiCloud. The example provides how to configure a shared ELB, listener, backend server group, health checks,
and EIP association for external access.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `loadbalancer_name` - The name of the load balancer
* `listener_name` - The name of the listener
* `instance_name` - The name of the ECS instance

#### Optional Variables

* `availability_zone` - The availability zone to which the ECS instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor (default: 4)
* `instance_image_id` - The image ID of the ECS instance (default: "")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The OS of the ECS instance image (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC (default: "172.16.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `subnet_dns_list` - The DNS list of the subnet (default: null)
* `is_associate_eip` - Whether to associate an EIP with the load balancer (default: false)
* `eip_address` - The address of the EIP (default: "")  
  Required when associating a new EIP with the load balancer and `bandwidth_name` is not provided.
* `bandwidth_name` - The name of the EIP bandwidth (default: "")  
  Required when associating a new EIP with the load balancer and `eip_address` is not provided.
* `bandwidth_size` - The bandwidth size of the EIP in Mbps (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `listener_protocol` - The protocol of the listener (default: "UDP")  
  Valid values: **HTTP**, **HTTPS**,  **TERMINATED_HTTPS**, **TCP**, **UDP**
* `listener_default_tls_container_ref` - The ID of the server certificate (default: "")
* `listener_server_certificate_name` - The name of the server certificate (default: "")  
  Required when `listener_protocol` is **TERMINATED_HTTPS** and `listener_default_tls_container_ref` is not provided
* `listener_server_certificate_private_key` - The private key of the server certificate (default: "")  
  Required when `listener_protocol` is **TERMINATED_HTTPS** and `listener_server_certificate_name` is provided
* `listener_server_certificate_certificate` - The content of the server certificate (default: "")  
  Required when `listener_protocol` is **TERMINATED_HTTPS** and `listener_server_certificate_name` is provided
* `listener_port` - The port of the listener (default: 80)
* `listener_description` - The description of the listener (default: "")
* `listener_tags` - The tags of the listener (default: {})
* `listener_http2_enable` - The HTTP/2 enable of the listener, only valid when `listener_protocol` is
  **TERMINATED_HTTPS** (default: false)
* `listener_client_ca_tls_container_ref` - The ID of the CA certificate of the listener, only valid when
  `listener_protocol` is **TERMINATED_HTTPS** (default: null)
* `listener_sni_container_refs` - The ID list of the SNI certificates of the listener, only valid when
  `listener_protocol` is **TERMINATED_HTTPS** (default: null)
* `listener_tls_ciphers_policy` - The TLS ciphers policy of the listener, only valid when
  `listener_protocol` is **TERMINATED_HTTPS** (default: null)
* `listener_insert_headers` - The insert headers of the listener, only valid when
  `listener_protocol` is **TERMINATED_HTTPS** (default: {})
* `pool_name` - The name of the backend server group (default: "")
* `pool_protocol` - The protocol of the backend server group (default: "UDP")
* `pool_method` - The load balancing algorithm of the backend server group (default: "ROUND_ROBIN")
* `pool_description` - The description of the backend server group (default: "")
* `pool_persistence` - The persistence of the backend server group (default: null)
* `member_protocol_port` - The protocol port of the backend server (default: 80)
* `member_weight` - The weight of the backend server (default: 1)
* `health_check_name` - The name of the health check (default: "health_check")
* `health_check_type` - The type of the health check (default: "UDP_CONNECT")
* `health_check_delay` - The delay between health checks in seconds (default: 10)
* `health_check_timeout` - The timeout for health checks in seconds (default: 5)
* `health_check_max_retries` - The maximum number of retries for health checks (default: 3)
* `health_check_port` - The port for health checks (default: null)
* `health_check_url_path` - The URL path for the health check (default: null)
* `health_check_http_method` - The HTTP method for the health check (default: null)
* `health_check_expected_codes` - The expected HTTP status codes for the health check (default: null)
* `health_check_domain_name` - The domain name for the health check (default: null)
* `security_group_rule_remote_ip_prefix` - The remote IP prefix of the security group rule (default: "100.125.0.0/16")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "tf_test_vpc"
  subnet_name         = "tf_test_subnet"
  security_group_name = "tf_test_security_group"
  loadbalancer_name   = "tf_test_loadbalancer"
  listener_name       = "tf_test_listener"
  instance_name       = "tf_test_ecs_instance"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Bind EIP to Load Balancer

To associate an Elastic IP (EIP) with the load balancer for external access, set the following variables:

### Option 1: Use Existing EIP

```hcl
is_associate_eip = true
eip_address      = "your_existing_eip_id"
```

### Option 2: Create New EIP and Bind to Load Balancer

```hcl
is_associate_eip = true
bandwidth_name   = "your_eip_bandwidth_name"
```

## Other Usage

This example creates a UDP listener by default. The following sections show how to configure other listener types.

### For TCP listener, configure the following required variables

```hcl
listener_protocol = "TCP"
pool_protocol     = "TCP"
health_check_type = "TCP"
```

### For HTTP listener, configure the following required variables

```hcl
listener_protocol           = "HTTP"
pool_protocol               = "HTTP"
health_check_type           = "HTTP"
health_check_expected_codes = "your_health_check_expected_codes"
health_check_url_path       = "your_health_check_url_path"
health_check_http_method    = "your_health_check_http_method"
```

### For TERMINATED_HTTPS listener, configuration one-way authentication listener

#### Option 1: Use existing server certificate

```hcl
listener_protocol                  = "TERMINATED_HTTPS"
listener_default_tls_container_ref = "your_existing_tls_container_ref_id"
pool_protocol                      = "HTTP"
health_check_type                  = "HTTP"
health_check_expected_codes        = "your_health_check_expected_codes"
health_check_url_path              = "your_health_check_url_path"
health_check_http_method           = "your_health_check_http_method"
```

#### Option 2: Create new server certificate

```hcl
listener_protocol                = "TERMINATED_HTTPS"
listener_server_certificate_name = "your_server_certificate_name"

listener_server_certificate_private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
[Your private key content here]
-----END RSA PRIVATE KEY-----
EOT
listener_server_certificate_certificate = <<EOT
-----BEGIN CERTIFICATE-----
    [Your certificate content here]
-----END CERTIFICATE-----
EOT

pool_protocol               = "HTTP"
health_check_type           = "HTTP"
health_check_expected_codes = "your_health_check_expected_codes"
health_check_url_path       = "your_health_check_url_path"
health_check_http_method    = "your_health_check_http_method"
```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The ELB is dependent on the VPC, subnet, and ECS instances
* Shared ELB backend server groups can only be used by one listener
* When health check is enabled, backend server security groups must allow ELB health check traffic.
  Please refer to [Configure Backend Server Security Groups](https://support.huaweicloud.com/usermanual-elb/elb_ug_hd_0002.html#elb_ug_hd_0002__zh-cn_topic_0000001434422737_section0207122815610)

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.73.4 |
