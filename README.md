## 简介

本示例来自于: https://github.com/github/platform-samples/tree/master/api/golang/basics-of-authentication

修改了配置读取的方式改为从config.toml中读取。下方是翻译为中文的原仓库README说明:

这是根据 developer.github.com 上的"认证基础指南"构建的示例项目，已移植到 Go 语言。

由于 Go 标准库没有内置的网页会话处理功能，因此只移植了简单示例。该示例还展示了如何使用 GitHub golang SDK。

## 配置

首先，你需要按照 GitHub OAuth 开发者指南中的步骤注册一个 OAuth 应用，可以复制下方地址去新建一个。

```
https://github.com/settings/applications/new
```

网站根地址`Homepage URL`设为:

```
http://localhost:4567
```

回调地址`Authorization callback URL`设为 :

```
http://localhost:4567/callback
```

创建好后，复制你新创建应用的客户端 ID 和密钥，并在`config.toml`文件将它们设置为环境变量：

## 启动

确保你已经安装了Go；然后通过运行以下命令获取 go-github 客户端库所需的模块：

```
go mod tidy
```

最后，启动项目，访问`http://localhost:4567`

```
go run main.go
```

在授权 GitHub OAuth 应用后，可以看到你的 GitHub 邮箱地址和用户名信息。

如果在重定向到 GitHub 时遇到任何错误，请检查你的环境变量以及注册 OAuth 应用时设置的回调 URL。