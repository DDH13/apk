# --------------------------------------------------------------------
# Copyright (c) 2023, WSO2 LLC. (http://wso2.com) All Rights Reserved.
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
# -----------------------------------------------------------------------
FROM envoyproxy/envoy:v1.33.0
LABEL maintainer="WSO2 Docker Maintainers <wso2.com>"

RUN apt-get update && apt-get upgrade -y && apt-get install -y curl

ENV LANG=C.UTF-8

ARG MG_USER=wso2
ARG MG_USER_ID=10001
ARG MG_USER_GROUP=wso2
ARG MG_USER_GROUP_ID=10001
ARG MG_USER_HOME=/home/${MG_USER}

ARG MOTD="\n\
 Welcome to WSO2 Docker Resources \n\
 --------------------------------- \n\
 This Docker container comprises of a WSO2 product, running with its latest GA release \n\
 which is under the Apache License, Version 2.0. \n\
 Read more about Apache License, Version 2.0 here @ http://www.apache.org/licenses/LICENSE-2.0.\n"

RUN \
    groupadd --system -g ${MG_USER_GROUP_ID} ${MG_USER_GROUP} \
    && useradd --system --create-home --home-dir ${MG_USER_HOME} --no-log-init -g ${MG_USER_GROUP_ID} -u ${MG_USER_ID} ${MG_USER} \
    && chgrp -R 0 ${MG_USER_HOME} \
    && chmod -R g=u ${MG_USER_HOME} \
    && echo '[ ! -z "${TERM}" -a -r /etc/motd ] && cat /etc/motd' >> /etc/bash.bashrc; echo "${MOTD}" > /etc/motd

ENV ROUTER_ADMIN_HOST=0.0.0.0
ENV ROUTER_ADMIN_PORT=9000

ENV ROUTER_CLUSTER=default_cluster
ENV ROUTER_LABEL="default"
ENV ROUTER_PRIVATE_KEY_PATH=/home/wso2/security/keystore/mg.key
ENV ROUTER_PUBLIC_CERT_PATH=/home/wso2/security/keystore/mg.pem

ENV ADAPTER_HOST=adapter
ENV ADAPTER_PORT=18000
ENV ADAPTER_CA_CERT_PATH=/home/wso2/security/truststore/mg.pem

ENV ENFORCER_HOST=enforcer
ENV ENFORCER_PORT=8081
ENV ENFORCER_ANALYTICS_RECEIVER_PORT=18090
# Enforcer Analytics host
ENV ENFORCER_CA_CERT_PATH=/home/wso2/security/truststore/mg.pem
ENV ENFORCER_ANALYTICS_HOST=$ENFORCER_HOST

ENV CONCURRENCY=2
ENV TRAILING_ARGS=

USER ${MG_USER_ID}
ENV ENVOY_UID=0
ENV ENVOY_GID=0

COPY resources/LICENSE.txt .
COPY resources/router_check_health.sh .
COPY resources/wasm/websocket/mgw-websocket.wasm /home/wso2/wasm/websocket/mgw-websocket.wasm
COPY resources/security/truststore/ca-certificates.crt /etc/ssl/certs
COPY resources/interceptor /home/wso2/interceptor
COPY resources/envoy.yaml /etc/envoy/envoy.yaml
CMD /usr/local/bin/envoy -c /etc/envoy/envoy.yaml --config-yaml "{admin: {address: {socket_address: {address: '${ROUTER_ADMIN_HOST}', port_value: '${ROUTER_ADMIN_PORT}'}}}, dynamic_resources: {ads_config: {api_type: GRPC, transport_api_version: V3, grpc_services: [{envoy_grpc: {cluster_name: xds_cluster}}]}, cds_config: {ads: {}, resource_api_version: V3}, lds_config: {ads: {}, resource_api_version: V3}}, node: {cluster: '${ROUTER_CLUSTER}', id: '${ROUTER_LABEL}', metadata: {instanceIdentifier : ${HOSTNAME}}}, static_resources: {clusters: [{name: xds_cluster, type: STRICT_DNS, connect_timeout: 1s, upstream_connection_options: {tcp_keepalive: {keepalive_probes: 3, keepalive_time: 300, keepalive_interval: 30}}, load_assignment: {cluster_name: xds_cluster, endpoints: [{lb_endpoints: [{endpoint: {address: {socket_address: {address: '${ADAPTER_HOST}', port_value: '${ADAPTER_PORT}'}}}}]}]}, typed_extension_protocol_options: {envoy.extensions.upstreams.http.v3.HttpProtocolOptions: {'@type': 'type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions', explicit_http_config: {http2_protocol_options: {}}}}, transport_socket: {name: envoy.transport_sockets.tls, typed_config: {'@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext, common_tls_context: {tls_params: {tls_minimum_protocol_version: TLSv1_2, tls_maximum_protocol_version: TLSv1_2}, tls_certificates: {private_key: {filename: '${ROUTER_PRIVATE_KEY_PATH}'}, certificate_chain: {filename: '${ROUTER_PUBLIC_CERT_PATH}'}}, validation_context: {trusted_ca: {filename: '${ADAPTER_CA_CERT_PATH}'}}}}}}, {name: ext-authz, type: STRICT_DNS, connect_timeout: 20s, typed_extension_protocol_options: {envoy.extensions.upstreams.http.v3.HttpProtocolOptions: {'@type': 'type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions', explicit_http_config: {http2_protocol_options: {}}}}, transport_socket: {name: envoy.transport_sockets.tls, typed_config: {'@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext, common_tls_context: {tls_params: {tls_minimum_protocol_version: TLSv1_2, tls_maximum_protocol_version: TLSv1_2}, tls_certificates: {private_key: {filename: '${ROUTER_PRIVATE_KEY_PATH}'}, certificate_chain: {filename: '${ROUTER_PUBLIC_CERT_PATH}'}}, validation_context: {trusted_ca: {filename: '${ENFORCER_CA_CERT_PATH}'}}}}}, load_assignment: {cluster_name: ext-authz, endpoints: [{lb_endpoints: [{endpoint: {address: {socket_address: {address: '${ENFORCER_HOST}', port_value: '${ENFORCER_PORT}'}}}}]}]}}, {name: access-logger, type: STRICT_DNS, connect_timeout: 200s, typed_extension_protocol_options: {envoy.extensions.upstreams.http.v3.HttpProtocolOptions: {'@type': 'type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions', explicit_http_config: {http2_protocol_options: {}}}}, transport_socket: {name: envoy.transport_sockets.tls, typed_config: {'@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext, common_tls_context: {tls_params: {tls_minimum_protocol_version: TLSv1_2, tls_maximum_protocol_version: TLSv1_2}, tls_certificates: {private_key: {filename: '${ROUTER_PRIVATE_KEY_PATH}'}, certificate_chain: {filename: '${ROUTER_PUBLIC_CERT_PATH}'}}, validation_context: {trusted_ca: {filename: '${ENFORCER_CA_CERT_PATH}'}}}}}, load_assignment: {cluster_name: access-logger, endpoints: [{lb_endpoints: [{endpoint: {address: {socket_address: {address: '${ENFORCER_ANALYTICS_HOST}', port_value: '${ENFORCER_ANALYTICS_RECEIVER_PORT}'}}}}]}]}}, {name: token_cluster, type: STRICT_DNS, connect_timeout: 20s, transport_socket: {name: envoy.transport_sockets.tls, typed_config: {'@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext, common_tls_context: {tls_params: {tls_minimum_protocol_version: TLSv1_2, tls_maximum_protocol_version: TLSv1_2}, tls_certificates: {private_key: {filename: '${ROUTER_PRIVATE_KEY_PATH}'}, certificate_chain: {filename: '${ROUTER_PUBLIC_CERT_PATH}'}}, validation_context: {trusted_ca: {filename: '${ENFORCER_CA_CERT_PATH}'}}}}}, load_assignment: {cluster_name: token_cluster, endpoints: [{lb_endpoints: [{endpoint: {address: {socket_address: {address: '${ENFORCER_HOST}', port_value: 8082}}}}]}]}},{name: jwks_cluster, type: STRICT_DNS, connect_timeout: 20s, transport_socket: {name: envoy.transport_sockets.tls, typed_config: {'@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext, common_tls_context: {tls_params: {tls_minimum_protocol_version: TLSv1_2, tls_maximum_protocol_version: TLSv1_2}, tls_certificates: {private_key: {filename: '${ROUTER_PRIVATE_KEY_PATH}'}, certificate_chain: {filename: '${ROUTER_PUBLIC_CERT_PATH}'}}, validation_context: {trusted_ca: {filename: '${ENFORCER_CA_CERT_PATH}'}}}}}, load_assignment: {cluster_name: jwks_cluster, endpoints: [{lb_endpoints: [{endpoint: {address: {socket_address: {address: '${ENFORCER_HOST}', port_value: 9092}}}}]}]}},{name: api_definition_cluster, type: STRICT_DNS, connect_timeout: 20s, transport_socket: {name: envoy.transport_sockets.tls, typed_config: {'@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext, common_tls_context: {tls_params: {tls_minimum_protocol_version: TLSv1_2, tls_maximum_protocol_version: TLSv1_2}, tls_certificates: {private_key: {filename: '${ROUTER_PRIVATE_KEY_PATH}'}, certificate_chain: {filename: '${ROUTER_PUBLIC_CERT_PATH}'}}, validation_context: {trusted_ca: {filename: '${ENFORCER_CA_CERT_PATH}'}}}}}, load_assignment: {cluster_name: api_definition_cluster, endpoints: [{lb_endpoints: [{endpoint: {address: {socket_address: {address: '${ENFORCER_HOST}', port_value: 8084}}}}]}]}}]} }" --concurrency "${CONCURRENCY}"  $TRAILING_ARGS
