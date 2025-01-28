# Use Ubuntu as the base image
FROM ubuntu:latest

# Install system dependencies
RUN apt-get update && apt-get install -y build-essential openjdk-17-jdk python3 python3-pip g++ gcc golang 
RUN apt-get install -y  nodejs npm
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*
# RUN useradd -u 1000 -m myuser
# USER myuser
# RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh
# RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
# ENV NVM_DIR /usr/local/nvm
# ENV NODE_VERSION 20
# # RUN source ~/.bashrc
# RUN nvm install 20
# # Set up a working directory
WORKDIR /app

# Install Python packages (if needed)
# RUN pip3 install --no-cache-dir numpy pandas                     # Example Python package

# Install Node.js packages (if needed)
RUN npm install -g yarn

# Set up Go environment (if needed)
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
# Copy any necessary scripts or files
COPY . .

# Default command (can be overridden)
CMD ["bash"]