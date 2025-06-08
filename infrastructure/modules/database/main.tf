terraform {
    required_providers {
        aws = {
            source  = "hashicorp/aws"
            version = "5.54.1"
        }
    }
}
resource "aws_docdb_cluster" "database" {
    cluster_identifier = "${var.project_name}-${var.environment}-docdb-cluster"
    engine             = "docdb"
    engine_version     = "5.0.0"
    master_username    = var.master_username
    master_password    = var.master_password
    skip_final_snapshot = true

    vpc_security_group_ids = var.app_security_group_ids
    db_subnet_group_name   = aws_docdb_subnet_group.mongo.name

    tags = {
        Name        = "${var.project_name}-${var.environment}-docdb-cluster"
        Environment = var.environment
    }
}

resource "aws_docdb_cluster_instance" "instance" {
    cluster_identifier = aws_docdb_cluster.database.cluster_identifier
    identifier = "${var.project_name}-${var.environment}-docdb-instance-${count.index}"

    instance_class     = var.instance_type
    count = var.instance_count
}

resource "aws_docdb_subnet_group" "mongo" {
    name       = "${var.project_name}-${var.environment}-docdb-subnet-group"
    subnet_ids = var.subnet_ids

    tags = {
        Name        = "${var.project_name}-${var.environment}-docdb-subnet-group"
        Environment = var.environment
    }
}
