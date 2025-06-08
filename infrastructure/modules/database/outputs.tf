output "cluster_endpoint" {
  description = "The endpoint of the DocumentDB cluster"
  value = aws_docdb_cluster.database.endpoint
}