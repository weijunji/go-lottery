# go-lottery
[![Go Report Card](https://goreportcard.com/badge/github.com/weijunji/go-lottery)](https://goreportcard.com/report/github.com/weijunji/go-lottery)
[![Build Status](https://www.travis-ci.com/weijunji/go-lottery.svg?branch=main)](https://www.travis-ci.com/weijunji/go-lottery)
[![codecov](https://codecov.io/gh/weijunji/go-lottery/branch/main/graph/badge.svg?token=wLuLssUnbF)](https://codecov.io/gh/weijunji/go-lottery)


点击上方的徽章查看代码CI结果以及测试覆盖率情况。

## 文件规范
* `cmd/`保存程序的main包，只使用最少的代码，每个程序新建一个文件夹
* `internal/`保存程序的所有私有代码，即不会被其他程序使用的代码，每个程序新建一个文件夹
* `pkgs/`保存公共代码，即会被其他程序使用到的代码
* `web/`保存前端代码

## 项目规范
* 不得直接提交到main分支上，所有代码都新建一个分支来写，如`auth-dev`分支。
* 通过github的pull request来合并代码到main分支，合并之前需要通过CI的测试以及组员的review。
* 能写单元测试的函数都要写单元测试，尽量提高测试覆盖率。
