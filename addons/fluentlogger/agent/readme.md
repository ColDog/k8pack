# `fluentlogger-agent`

This is a build of `logspout` that includes the great `segmentio/logspout-fluentd` plugin. For Kubernetes, this setup sends the following payload to the fluentd forward handler.

```json
{
   "docker.hostname":"fluentlogger-agent-mhtw7",
   "docker.id":"4f3058d680642a355303b09e8083277b2cf71715a8b5dc28241ccf70b7c17b37",
   "docker.image":"coldog/fluentlogger-agent@sha256:3bba88afc4bcc5347c29f37a66309c3fdaddf2a33b58d1728a36c332725ebcdc",
   "docker.label.annotation.io.kubernetes.container.hash":"f0439694",
   "docker.label.annotation.io.kubernetes.container.restartCount":"0",
   "docker.label.annotation.io.kubernetes.container.terminationMessagePath":"/dev/termination-log",
   "docker.label.annotation.io.kubernetes.container.terminationMessagePolicy":"File",
   "docker.label.annotation.io.kubernetes.pod.terminationGracePeriod":"10",
   "docker.label.io.kubernetes.container.logpath":"/var/log/pods/9fc0c314-7545-11e7-b6fd-0221cffd1ba2/fluentlogger-agent_0.log",
   "docker.label.io.kubernetes.container.name":"fluentlogger-agent",
   "docker.label.io.kubernetes.docker.type":"container",
   "docker.label.io.kubernetes.pod.name":"fluentlogger-agent-mhtw7",
   "docker.label.io.kubernetes.pod.namespace":"kube-system",
   "docker.label.io.kubernetes.pod.uid":"9fc0c314-7545-11e7-b6fd-0221cffd1ba2",
   "docker.label.io.kubernetes.sandbox.id":"e145507bf4e3a7016e3563cffe4df224b5f648843a02b558466d14f74ee5343c",
   "docker.name":"/k8s_fluentlogger-agent_fluentlogger-agent-mhtw7_kube-system_9fc0c314-7545-11e7-b6fd-0221cffd1ba2_0",
   "message":"#   fluentd-tcp\t10.3.0.50:24224\t\t\t\tmap[]"
}
```
