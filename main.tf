provider "aws" {
  region = "us-east-1"
}

# Define the security group resource
resource "aws_security_group" "web_server_sg" {
  name        = "web_server_security_group"
  description = "Allow SSH and HTTP inbound traffic"
  # Optional: specify a non-default VPC ID if necessary
  # vpc_id      = "vpc-12345678" 

  # Ingress rule for SSH (port 22)
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # WARNING: 0.0.0.0/0 for SSH is insecure, use a specific IP range
    description = "Allow SSH from anywhere"
  }

  # Ingress rule for HTTP (port 80)
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow HTTP from anywhere"
  }
}
resource "aws_instance" "web" {
  ami           = "ami-0c02fb55956c7d316"
  instance_type = "t2.micro"

  user_data = <<-EOF
              #!/bin/bash
              sudo apt update -y
              sudo amazon-linux-extras install nginx1
              sudo systemctl start nginx
              sudo systemctl enable nginx
              EOF

  tags = {
    Name = "Ollama 0.2"
  }
}