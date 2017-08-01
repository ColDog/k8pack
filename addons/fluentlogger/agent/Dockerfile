FROM fluent/fluent-bit:0.11.15
USER root

ADD fluent-bit.conf /fluent-bit/etc/
ADD parsers.conf /fluent-bit/etc/

CMD ["/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/etc/fluent-bit.conf"]
