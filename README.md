# scheduler-framework-demo
基于Kubernetes V1.17.16版本Scheduler Framework机制扩展Scheduler功能, 相关监控指标数据依赖于dynamic-annotator项目, 已在本地生产集群验证.

## 已支持扩展特性
- 优选打分阶段增加基于节点cpu使用率扩展
  - 修复当节点负载较高时node-exporter上报指标失败，导致dynamic-annotator获取到的CPU使用率为0然后该节点得分过高问题
- 预选阶段增加基于CPU使用率、内存使用率、load1、load5等指标过滤节点
