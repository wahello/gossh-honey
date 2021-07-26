# gossh-honey

使用 Go 语言编写的 ssh 服务端程序，可以只通过一个二进制文件运行来实现 ssh 服务，一个小型轻量级的 ssh 服务只需要一个文件即可完成，并且还拥有高度的可自定义化，方便不同场景使用。

**已经具有的功能**

- [x] 作为ssh服务端
- [x] 捕获用户输入的用户名和密码，日志输出

**需要进行完善的功能**

- [ ] 可以将用户名和密码分开输出到指定文件中，共统计分析使用，哪些是常用的爆破名
- [ ] 可以写一个脚本文件进行循环连接ssh服务端，方便记录
- [ ] 当我使用远程连接工具，输入正确的用户名和密码时却显示` Waiting for the pending transfer to complete...  `用户没有进入linux服务器进行操作
- [ ] ...

## 使用方法

将文件下载来，使用go build 进行构建生成二进制文件，直接运行监听端口即可。

![image-20210721171133023](https://raw.githubusercontent.com/zmk-c/blogImages/master/img/20210721171133.png)

## 配置文件
在 `config/config.json` 中配置了正确的用户名和密码等，可以按照你自己的需求进行更改。
```json
{
    "name": "root",         登录的用户名
    "password": "Zmk970309",     登录的密码
    "command": "/bin/bash", 验证后使用者需要运行程序
    "port": 2222            监听的端口
}
```

