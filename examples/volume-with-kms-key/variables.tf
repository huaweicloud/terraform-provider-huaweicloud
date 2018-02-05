variable "user_name" {
    default = "test_user_name"
    description = "The Username to login with."
}

variable "domain_name" {
    default = "test_domain_name"
    description = "The Name of the Domain to scope to (Identity v3)."
}

variable "tenant_name" {
    default = "cn-north-1"
    description = "The Name of the Tenant (Identity v2) or Project (Identity v3) to login with."
}

variable "region" {
    default = "cn-north-1"
    description = "The region of the HuaweiCloud cloud to use."
}

variable "password" {
    default = "********"
    description = "The Password to login with."
}

variable "auth_url" {
    default = "https://iam.cn-north-1.myhwclouds.com"
    description = "The Identity authentication URL."
}
