# prometheusSendToDingTlak
网上也有Prometheus告警到钉钉的插件，但公司的prometheus监控需要先连接vpn，使用网上的钉钉插件时vpn断开后也会发送告警，故网上的插件不能满足需求，当时心血来潮自己写了一个
# 说明
将prometheus配置告警到prometheus-Dingtalk的服务上，当prometheus发送告警后可以转发到钉钉；可选的<kbd>-i</kbd>参数可以对VPN主机进行检测，如果vpn断开可以通过可选的<kbd>-m</kbd>参数来向钉钉发送消息
## Node.js版本
master分支
## Go版本
1.1.0分支
# 使用方法
main [OPTION] URL1 URL2 ...


  -i 检查这个ip是否能ping通
  
  
  -s 连接成功时发送的消息，支持@
  
  
  -m 参数-i中的ip连接不通时，向钉钉发送的消息，支持@
  
  
  -p 指定服务启动的端口，默认5001
