# helm-maker
create helm chart by data

# note
It’s not completed yet, do not use it

# 均衡

- 易用性
- 适用范围
- 复杂度

# 对比

| 方式对                                 | 易用性 |  适用度 | 规范性|
|:------------------------------------- |:------:| :-----:|:------:|
| 将应用所有的组件的部署yaml都拿过来按格式打包成chart       | 高 | 高 | 低 |
| 参数和模板完全都是json,直接转成yaml并按chart存放文件         | 高 | 高 | 低 |
| 自定义一些模板,由传入的json来填充,json部分支持满足自定义需求  | 中 | 中 | 高 |

