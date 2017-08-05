FROM fluent/fluentd:stable

RUN apk add --update --virtual .build-deps sudo build-base ruby-dev \
    && sudo gem install \
        fluent-plugin-rewrite-tag-filter \
        fluent-plugin-cloudwatch-logs \
        fluent-plugin-remote-syslog \
    && sudo gem sources --clear-all \
    && apk del .build-deps \
    && rm -rf /var/cache/apk/* /home/fluent/.gem/ruby/2.3.0/cache/*.gem

ADD fluentd.conf /etc/fluentd/fluentd.conf

USER root
EXPOSE 24224
CMD ["fluentd", "-c", "/etc/fluentd/fluentd.conf"]
