# gossh-honey

使用 Go 语言编写的 ssh 服务端蜜罐，不限制用户的登录，模拟伪终端系统，监听用户操作。

**已经具有的功能**

- [x] 作为ssh服务端
- [x] 记录连接信息(src ip / dest ip)  
- [x] 捕获用户输入的用户名和密码，日志输出
- [x] 当用户进入伪终端后，记录用户的命令操作

**待完善的功能**
- [ ] 更完善的shell命令模拟
- [ ] 降权操作

## 使用方法

```
git clone https://github.com/zmk-c/gossh-honey.git
cd gossh-honey/
go build
./gossh-honey
```

