# my_todo

#### 介绍
待办任务，支持创建待办，删除待办，定时任务生成今日任务，查看今日任务，完成今日任务，查看周期内任务统计数据，查看或导出指定时间范围的统计数据，导入任务数据

#### 软件架构
使用go http 处理请求，

gorm orm工具完成数据的存储查询

go timer 实现定时任务生成今日待办

策略+工厂 实现不同周期类型数据处理


#### 安装教程

1.  安装MySQL，设置MySQL用户名命名
2.  启动项目
3.  请求接口

#### 使用说明

1.  可使用 test 包里的 data_test.go 生成一些历史数据，使得调用接口时有数据
2.  可参照 doc/http 目录里的 http文件执行 接口http请求
3.  这两天不吃不喝完成的项目（开玩笑 ~_~），已经实现了目前能想到的所有功能，后续我有好的想法再持续改进
4.  项目可部署在docker中，部署文档即部署说明 见 doc/docker 目录

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
