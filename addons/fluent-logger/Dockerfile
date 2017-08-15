FROM fluent/fluentd:v0.14.20-debian

RUN buildDeps="sudo make gcc g++ libc-dev ruby-dev libffi-dev" \
     && apt-get update \
     && apt-get install \
     -y --no-install-recommends \
     $buildDeps net-tools \
    && echo 'gem: --no-document' >> /etc/gemrc \
    && gem install fluent-plugin-secure-forward \
    && gem install fluent-plugin-record-reformer \
    && gem install fluent-plugin-cloudwatch-logs \
    && gem install fluent-plugin-kubernetes_metadata_filter \
    && gem install ffi \
    && gem install fluent-plugin-systemd \
    && SUDO_FORCE_REMOVE=yes \
    apt-get purge -y --auto-remove \
                  -o APT::AutoRemove::RecommendsImportant=false \
                  $buildDeps \
 && rm -rf /var/lib/apt/lists/* \
    && gem sources --clear-all \
    && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem

COPY plugins /fluentd/plugins/
COPY fluentd.conf /etc/fluentd/fluentd.conf

EXPOSE 24224

USER root
ENTRYPOINT exec fluentd -c /etc/fluentd/fluentd.conf -p /fluentd/plugin
