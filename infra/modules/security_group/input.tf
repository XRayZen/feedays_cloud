
variable "security_group_name" {
    type = string
    description = "Name of the security group"
}

variable "security_group_description" {
    type = string
    description = "Description of the security group"
}

variable "security_group_vpc_id" {
    type = string
    description = "VPC ID"
}

variable "security_group_ingress_cidr_blocks" {
    type = list(string)
    description = "List of CIDR blocks to use on ingress rules"
}

