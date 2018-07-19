###############################################################################
# Copyright 2018 Samsung Electronics All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
###############################################################################
# Docker image for "pharos-anchor"
FROM alpine:3.6

# environment variables
ENV APP_DIR=/pharos
ENV APP_PORT=48099
ENV APP=pharos-anchor

# install MongoDB
RUN apk add --no-cache mongodb && \
    rm -rf /var/cache/apk/*

# copy files
COPY $APP run.sh $APP_DIR/

# expose notifications port
EXPOSE $APP_PORT

# set the working directory
WORKDIR $APP_DIR

# make mogodb volume
RUN mkdir /data && mkdir /data/db

# kick off the pharos-anchor container
CMD ["sh", "run.sh"]
