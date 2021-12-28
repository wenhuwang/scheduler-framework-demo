# scheduler-framework-demo
基于Kubernetes V1.17.16版本Scheduler Framework机制扩展Scheduler功能, 相关监控指标数据依赖于dynamic-annotator项目, 已在本地生产集群验证.

## 已支持扩展特性
- 优选打分阶段增加基于节点cpu使用率扩展
- 节点负载(开发中)