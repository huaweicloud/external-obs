# Based on centos
FROM centos:7.4.1708
LABEL maintainers="Kubernetes Authors"
LABEL description="OBS Flex Volume"
LABEL mountpath="none"

# Environments
ENV OBSAccessKey $OBSAccessKey
ENV OBSSecretKey $OBSSecretKey

# Install libraries
RUN yum update -y
RUN command yum install automake fuse-devel gcc-c++ git libcurl-devel libxml2-devel make openssl-devel -y

# Compile fuse
WORKDIR /
RUN git clone https://github.com/s3fs-fuse/s3fs-fuse.git
WORKDIR /s3fs-fuse
RUN ./autogen.sh  
RUN ./configure 
RUN make 
RUN make install

RUN mkdir -p /mnt/mountpoint
RUN echo $'#!/bin/sh\n \
           echo $OBSAccessKey:$OBSSecretKey > /etc/passwd-s3fs\n \
           chmod 600 /etc/passwd-s3fs\n \
           exec s3fs "$@"' > /obs-flexvolume.sh 
RUN chmod +x /obs-flexvolume.sh

# Define default command
ENTRYPOINT ["/obs-flexvolume.sh"]
