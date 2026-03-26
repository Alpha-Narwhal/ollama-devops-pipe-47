provider "aws" {
  region = "us-east-1"
}

# Define the security group resource
resource "aws_security_group" "ollama_server_sg" {
  name        = "ollama_server_security_group"
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

  #Ingress rule for Ollama (Port 11434)
  ingress {
    from_port = 11434 
    to_port = 11434 
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allows Ollama through"
  }

}
resource "aws_instance" "Ollama" {
  ami           = "ami-0ec10929233384c7f"
  instance_type = "t2.medium"

   root_block_device {
    volume_size = 20 # Sets the root volume size to 50 GB
    volume_type = "gp3" # Optional: specify volume type, e.g., gp2, gp3, io1, etc.
  }

  user_data = <<-EOF
              #!/bin/bash
              sudo apt update -y
              sudo amazon-linux-extras install nginx1
              sudo systemctl start nginx
              sudo systemctl enable nginx

              sudo yum install curl -y
              sudo yum install zstd -y
              curl -fsSL https://ollama.com/install.sh | sh
              EOF

  tags = {
    Name = "Ollama 0.4"
  }
}