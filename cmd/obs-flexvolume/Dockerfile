# Based on centos
FROM centos:7.4.1708
LABEL maintainers="Kubernetes Authors"
LABEL description="OBS Flex Volume"

# Install libraries
RUN yum update -y

# Copy from build directory
COPY obs-flexvolume /obsflexvolume
COPY obs-flexvolume.sh /obs-flexvolume.sh
RUN chmod +x /obs-flexvolume.sh

# Define default command
ENTRYPOINT ["/obs-flexvolume.sh"]
