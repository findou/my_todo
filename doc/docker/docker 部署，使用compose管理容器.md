# docker 部署，使用compose管理容器



## 1 编写 dockerfile 构建镜像

```dockerfile
# 基础镜像
FROM golang:1.19-alpine3.15 AS builder
# 设置go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \ 
    GOPROXY="https://goproxy.cn"

# 切换工作目录
WORKDIR "/build"
# 拷贝代码到容器中
COPY . .
# go构建可执行文件
RUN go build -o my_todo ./main.go

### 接下来创建一个小镜像
FROM scratch
# 从builder镜像中把 /dist/app 拷贝到当前目录
COPY --from=builder /build/my_todo /

# 暴露端口
EXPOSE 8888
# 启动容器时候运行的命令
ENTRYPOINT ["./my_todo"]
```

编写后可在当前目录执行 dockerfile，生成镜像，也可以在docker compose 执行

```
docker build -t my_todo:1.0.0 .
```

## 2 编写 dockercompose ，编排容器

```yaml
version: '3.8'
services:
  my_todo:
    restart: always #docker重启时候，容器也重启
    build: #构建docker镜像
      context: ./ #Dockerfile 文件的目录
      dockerfile: Dockerfile #Dockerfile 文件的名称
    image: my_todo:1.0.0 #镜像名称和版本号
    container_name: my_todo1.0.1 #容器名称
    ports: # 宿主机和容器之间映射端口
      - 8888:8888

  mysql:
    image: mysql
    restart: always
    container_name: todo_mysql
    ports:
      - 3306:3306
    environment:
      - MYSQL_ROOT_PASSWORD=123456
```

在当前目录执行命令，启动容器

```
docker compose up
```



## 3 本地链接MySQL

进入容器

Docker exec -it [CONTINER ID] bash

登录mysql

Mysql -u root -p

Alter user ‘root’@’localhost’ IDENTIFIED BY ‘密码’

Use mysql

查看plugin设置

Select host,user, plugin from user;

![img](E:\学习文档\docker\2.jpg) 

可以看到root的plugin是caching_sha2_password，我们希望改成mysql_native_password

而这里总共有两个root，一个代表的是远程连接，一个是本地连接，所以根据自己需求修改。

修改认证方式

![img](E:\学习文档\docker\3.jpg) 

添加远程登录用户

CREATE USER ‘用户名’@’%’ IDENTIFIED WITH mysql_native_password By ‘密码’

GRANT ALL PRIVILEGES ON *.* TO ‘用户名’@’%’

ALL代表所有权限，也可以选择增删改查等权限，小项目嫌麻烦，所以选择ALL

mydb是数据库的名称，您可更改为你自己数据库名称

myuser是数据库用户，您可更改为自己的数据库用户

%代表所有IP都能访问，您可更改为某个IP才能访问

```
GRANT ALL PRIVILEGES ON mydb.* TO ‘myuser’@’%’;
GRANT ALL PRIVILEGES ON mysql.* TO 'root'@'%';
```

]DBeaver连接mysql时Public Key Retrieval is not allowed错误解决附图片

在新建连接的时候，驱动属性里设置 allowPublicKeyRetrieval 的值为 true。

进入MySQL容器创建数据库 my_todo