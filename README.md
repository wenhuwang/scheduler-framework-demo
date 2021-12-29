# scheduler-framework-demo
基于Kubernetes v1.19.16版本Scheduler Framework机制扩展Scheduler功能, 相关监控指标数据依赖于dynamic-annotator项目, 已在本地生产集群验证.

## 已支持扩展特性
- 优选阶段增加基于节点当前cpu和内存使用量进行打分
- 节点负载(开发中)
