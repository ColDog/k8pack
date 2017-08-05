
syslog-server:
	@docker build -t coldog/test-syslog-server ./tests/syslog-server
	@docker run -p 24224:24224/udp coldog/test-syslog-server:latest
.PHONY: syslog-server
