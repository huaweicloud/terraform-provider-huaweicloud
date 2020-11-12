variable "iden_user_name_default" {
  type    = string
  default = "iden_user_default"
}

variable "iden_user_password_default" {
  type    = string
  default = "hw123456default"
}

variable "iden_users" {
    description = "List of identity users."
    type = list(object({
        name        = string
        description = string
    }))
    default = [
        {name = "iden_user_A", description = "This is a identity user.", password = "hw123456"},
        {name = "iden_user_B", description = "This is a identity user.", password = "hw654321"}
    ]
}