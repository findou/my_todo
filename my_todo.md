# 待办项目

支持创建待办，删除待办，定时任务生成今日任务，查看今日任务，完成今日任务，查看周期内任务统计数据，查看或导出指定时间范围的统计数据，导入任务数据

## 代码目录介绍

```
- data # 包括了根据json文件随机生成历史待办数据的处理逻辑和数据导入的逻辑
- dataOracles # 对应 dataOracles 的http handler 方法
- todos # 对应 todos 的http handler 方法
- tasks # 对应 tasks 的http handler 方法
- statisViews  # 对应 statisViews 的http handler 方法
- handletodo # 对应待办任务相关处理逻辑，选用工厂+策略模式，不同的周期类型执行不同的策略
- model # 对应 数据库设计和一些返回接收结构体字段的定义
- utils # 一些工具函数，包括了日期处理，id生成，时间转化 方法
- schedule # 定时任务，用于在每天1点执行生成今日的待办数据，使用线程池控制并发量
- test # test文件，data_test.go 调用 data目录的 genHistoryData.go 里的代码根据这个目录里的data.json 文件生成历史待办，用于测试数据统计
- doc # 一些文件
 -- http # 存放http文件用于接口测试
 -- json # json文件，存放导出的json文件，或者导入文件

```



## 接口功能介绍

doc/http 目录的 http文件可以用来测试接口

### todo.http

这个文件包含了待办相关的接口，对应代码处理逻辑在待办todos包的todo.go 文件中，包括了创建待办，和删除待办两个接口

#### 1.创建待办

创建待办的时候设置待办名称，待办优先级，待办周期类型（1日，2周，3月，4年），待办开始时间结束时间（开始时间结束时间包括整数个周期时间），

cycleArrays：创建待办的周期值

比如，周期类型是2（周），cycleArrays的值是一周的天数的数组如（周1，周2，到周7）

周期类型和周期数组对应关系的说明如表

| 周期类型 | 周期类型数组取值范围                     | 示例              |
| -------- | ---------------------------------------- | ----------------- |
| 1 DAY    | 一天的任意时间 如8:30 ,用数组json表示    | "[\"8:30\"]"      |
| 2 WEEK   | 一周的 天数 ，取值范围为 周1到周7        | "[6,7]"           |
| 3 MONTH  | 一个月的天数，如 1号，到这个月最大的号数 | "[2,10,15,23,28]" |
| 4 YEAR   | 用 月-日 表示，表示某年的几月几号        | "[\"03-08\"]"     |

创建待办的格式

```
{
    "name": "51劳动节",
    "priority": 10,
    "startTime": "2022-01-01 00:00:00",
    "endTime": "2023-12-31 23:59:59",
    "cycleType": 4,
    "cycleArrays": "[\"05-01\"]"
}
```

接口地址

```http
### 创建待办
POST  http://127.0.0.1:8888/todo/
Content-Type: application/json

{
    "name": "51劳动节",
    "priority": 10,
    "startTime": "2022-01-01 00:00:00",
    "endTime": "2023-12-31 23:59:59",
    "cycleType": 4,
    "cycleArrays": "[\"05-01\"]"
}
```



#### 2.删除待办

更具待办的id删除已经创建的待办，实际没有删除，将待办的状态设置为了false，定时任务在生成每日任务的时候查询到待办的状态为false就不会生成这条待办的任务，待办删除后不可恢复，需要可以创建新的待办

接口地址

```http
### 删除待办
#DELETE http://127.0.0.1:8888/todo/<todo_id>
DELETE http://127.0.0.1:8888/todo/3967dac2-3a7a-4cd0-bb2f-f107f12bae51
```

### task.http

这个文件包含了查看任务相关的接口，代码对应 tasks目录的 task.go 文件，包括了查看今日任务列表，和修改任务完成状态为已完成两个接口

#### 1查看今日任务列表

使用get请求 返回今日任务列表

接口地址

```http
### 查看今日任务列表
GET http://127.0.0.1:8888/task/
```

#### 2 修改任务完成状态为已完成

使用put请求，传入任务id，就会修改任务的状态为已完成

接口地址

```http
### 修改任务完成状态为已完成
# PUT http://127.0.0.1:8888/task/<task_id>
PUT http://127.0.0.1:8888/task/d912f753-5547-4dfc-ac01-fc71abc12f94
```

### statisView.http

这个文件包含了，根据任务周期，查询指定时间的周期内的任务纵览，返回的数据包括周期内任务的总量，任务的完成量，待办任务的总量与待办任务的完成量等，如果没有传入<cycle_type>默认查询一天以内的数据，如果没有传入<data> 默认查询系统时间的周期数据，如果传入的<cycle_type>大于4，将<cycle_type>设置为1

逻辑代码对应statisViews目录里的statisViews.go 文件

接口示例

```http
### 待办任务周期总览
# 根据日期查询这个日期某个周期类型的待办事项，返回结果优先级高的在前 ，
# 周期类型有cycle_type 1日，2周，3月，4年
# 没传指定日期，默认当前时间的某个周期
# GET http://127.0.0.1:8888/statisView/<cycle_type>/<date>

### 1 查询 2023-04-19 这一天的 待办任务总览
GET http://127.0.0.1:8888/statisView/1/2023-04-19

### 2 查询 2023-03-19 这一周的 待办任务总览
GET http://127.0.0.1:8888/statisView/2/2023-03-19

### 3 查询 当前调用时间 这一月的 待办任务总览
GET http://127.0.0.1:8888/statisView/3

### 4 查询 当前调用时间 这一年的 待办任务总览
GET http://127.0.0.1:8888/statisView/4
```

### dataOracle.http

包括了待办任务数据查看，待办任务数据导出，待办任务导入 功能，逻辑代码对应dataOracles目录的dataOracles.go 文件

#### 数据导出查询

实现了根据用户传入的开始结束时间 ，查询任务在开始结束时间之内的待办任务数据，如果用户指定了fileType=json，就会导出json格式的数据文件，如果没有指定返回json格式的数据到body

接口示例

```http
### 数据导出，不带fileType=json 不会生成json文件，返回数据到body
# GET http://127.0.0.1:8888/statisView?startTime=<start_time>&&endTime=<end_time>
# 在浏览器输入地址下载json文件
# GET http://127.0.0.1:8888/statisView?startTime=<start_time>&&endTime=<end_time>&&fileType=json

### 导出指定时间范围内的数据，格式为json，将地址放到浏览器执行会下载文件
GET http://127.0.0.1:8888/dataOracle?startTime=2023-03-19 00:00:00&&endTime=2023-03-19 23:59:59&&fileType=json

### 2 查看今日任务与完成情况，开始时间结束时间设置在今日开始结束
GET http://127.0.0.1:8888/dataOracle?startTime=2023-03-19 00:00:00&&endTime=2023-03-19 23:59:59


### 3 查看今日任务与完成情况，结束时间设置在晚上8点，比 上 2 少一条数据（喝牛奶）因为喝牛奶时间在晚上8点半
GET http://127.0.0.1:8888/dataOracle?startTime=2023-03-19 00:00:00&&endTime=2023-03-19 20:00:00
```

#### 数据导入

导入json格式的文件，解析保存到数据库，格式为导出的格式

接口示例

```
### 数据导入,上传json文件导入数据，json文件格式 见 doc/json 目录 的 uploadFile.json 文件
POST http://127.0.0.1:8888/dataOracle
```

## 随机生成待办数据（用于数据测试）

这部分代码对应 test目录里的data_test.go 文件

有两个方法 第一个更具json文件生成待办，第二个更具数据库里的待办，生成任务，并随机选择一些任务设置状态为已完成

```go

func TestGenHistoryTodoData(t *testing.T) {
	// 根据json文件生成历史待办
	data.GenHistoryTodoData()
}

func TestGenHistoryTaskData(t *testing.T) {
	// 根据历史待办生成历史任务，并随机选择一些任务更改状态为已完成
	data.GenHistoryTaskData()
}
```

