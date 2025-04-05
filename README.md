# Gin Blog 项目学习简介

## 1 项目来源

本项目是对已有项目[gin-vue-blog](https://github.com/szluyu99/gin-vue-blog) [https://github.com/szluyu99/gin-vue-blog] 的学习记录，通过学习该项目，从 0 到 1 完整搭建了一个 go-web 博客项目，并将搭建的路径完整记录了下来，方便后续调整。

项目目标：

1. 给出一个从 0 到 1 的完整的搭建记录，方便从零入手
2. 搭建一个属于自己的博客项目，方便后续调整模块、项目部署以及博客记录。
3. 前端学习的同学，可以直接部署现有后端，根据文档记录搭建前端并学习。
4. 后端学习的同学，可以直接部署现有前端，根据文章记录搭建后端并学习。
5. 在上述基础上，根据文档了解项目各部分功能，进而方便实现自己的自定义功能。

项目文档记录：

1. [架构与初始化](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/1%20%E6%9E%B6%E6%9E%84%E4%B8%8E%E5%88%9D%E5%A7%8B%E5%8C%96.md)
2. [gin-blog-front 首页搭建](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/2%20gin-blog-front%20%E9%A6%96%E9%A1%B5%E6%90%AD%E5%BB%BA.md)
3. [gin-blog-front BackToTopGlobalModal](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/3%20gin-blog-front%20BackToTopGlobalModal.md)
4. [gin-blog-front 归档-分类-标签](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/4%20gin-blog-front%20%E5%BD%92%E6%A1%A3-%E5%88%86%E7%B1%BB-%E6%A0%87%E7%AD%BE.md)
5. [gin-blog-front 文章-相册-友情链接-关于我](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/5%20gin-blog-front%20%E6%96%87%E7%AB%A0-%E7%9B%B8%E5%86%8C-%E5%8F%8B%E6%83%85%E9%93%BE%E6%8E%A5-%E5%85%B3%E4%BA%8E%E6%88%91.md)
6. [gin-blog-front 留言个人-中心-无法匹配404](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/6%20gin-blog-front%20%E7%95%99%E8%A8%80%E4%B8%AA%E4%BA%BA-%E4%B8%AD%E5%BF%83-%E6%97%A0%E6%B3%95%E5%8C%B9%E9%85%8D404.md)
7. [gin-blog-admin 首页搭建](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/7%20gin-blog-admin%20%E9%A6%96%E9%A1%B5%E6%90%AD%E5%BB%BA.md)
8. [gin-blog-admin layout](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/8%20gin-blog-admin%20layout.md)
9. [gin-blog-admin login-home-article-auth](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/9%20gin-blog-admin%20login-home-article-auth.md)
10. [gin-blog-admin errorpage-log-message-profile](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/10%20gin-blog-admin%20errorpage-log-message-profile.md)
11. [gin-blog-admin setting-test-user](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/11%20gin-blog-admin%20setting-test-user.md)
12. [gin-blog-server 后端初始化](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/12%20gin-blog-server%20%E5%90%8E%E7%AB%AF%E5%88%9D%E5%A7%8B%E5%8C%96.md)
13. [gin-blog-server 注册-邮箱验证-登陆](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/13%20gin-blog-server%20%E6%B3%A8%E5%86%8C-%E9%82%AE%E7%AE%B1%E9%AA%8C%E8%AF%81-%E7%99%BB%E9%99%86.md)
14. [gin-blog-server 配置相关以及上报信息](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/14%20gin-blog-server%20%E9%85%8D%E7%BD%AE%E7%9B%B8%E5%85%B3%E4%BB%A5%E5%8F%8A%E4%B8%8A%E6%8A%A5%E4%BF%A1%E6%81%AF.md)
15. [gin-blog-server INFO-关于我-前台首页-前台页面](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/15%20gin-blog-server%20INFO-%E5%85%B3%E4%BA%8E%E6%88%91-%E5%89%8D%E5%8F%B0%E9%A6%96%E9%A1%B5-%E5%89%8D%E5%8F%B0%E9%A1%B5%E9%9D%A2.md)
16. [gin-blog-server 菜单模块-用户模块](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/16%20gin-blog-server%20%E8%8F%9C%E5%8D%95%E6%A8%A1%E5%9D%97-%E7%94%A8%E6%88%B7%E6%A8%A1%E5%9D%97.md)
17. [gin-blog-server 分类模块-标签模块-文章模块](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/17%20gin-blog-server%20%E5%88%86%E7%B1%BB%E6%A8%A1%E5%9D%97-%E6%A0%87%E7%AD%BE%E6%A8%A1%E5%9D%97-%E6%96%87%E7%AB%A0%E6%A8%A1%E5%9D%97.md)
18. [gin-blog-server 评论-留言-友链-资源模块](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/18%20gin-blog-server%20%E8%AF%84%E8%AE%BA-%E7%95%99%E8%A8%80-%E5%8F%8B%E9%93%BE-%E8%B5%84%E6%BA%90%E6%A8%A1%E5%9D%97.md)
19. [gin-blog-server 角色模块-操作日志模块-页面模块-文件上传模块](https://github.com/Tjyy-1223/Gin-Vue-Admin/blob/main/19%20gin-blog-server%20%E8%A7%92%E8%89%B2%E6%A8%A1%E5%9D%97-%E6%93%8D%E4%BD%9C%E6%97%A5%E5%BF%97%E6%A8%A1%E5%9D%97-%E9%A1%B5%E9%9D%A2%E6%A8%A1%E5%9D%97-%E6%96%87%E4%BB%B6%E4%B8%8A%E4%BC%A0%E6%A8%A1%E5%9D%97.md)



## 2 功能点介绍

其主要功能如下：

本项目在以**博客**这个业务为主的前提下，提供一个完整的全栈项目代码（前台前端 + 后台前端 + 后端），技术点基本都是最新 + 最火的技术，代码轻量级，注释完善，适合学习。

### 1.1 前台功能

- 前台界面设计参考 Hexo 的 Butterfly 设计，美观简洁
- 响应式布局，适配了移动端
- 实现点赞，统计用户等功能 (Redis)
- 评论 + 回复评论功能
- 留言采用弹幕墙，效果炫酷
- 文章详情页有文章目录、推荐文章等功能，优化用户体验

### 1.2 后台功能

- 采用 Restful 风格的 API
- 前后端分离部署，前端使用 Nginx，后端使用 Docker
- 鉴权使用 JWT
- 权限管理使用 CASBIN，实现基于 RBAC 的权限管理
- 支持动态权限修改，前端菜单由后端生成（动态路由）
- 文章编辑使用 Markdown 编辑器
- 常规后台功能齐全：侧边栏、面包屑、标签栏等
- 实现记录操作日志功能（GET 不记录）
- 实现监听在线用户、强制下线功能
- 文件上传支持七牛云、本地（后续计划支持更多）
- 对 CRUD 操作封装了通用 Hook

### 1.3 补充功能

- ~~完善图片上传功能, 目前文件上传还没怎么处理~~ 🆗
- 后台首页重新设计（目前没放什么内容）
- ~~前台首页搜索文章（目前使用数据库模糊搜索）~~ 🆗
- ~~博客文章导入导出 (.md 文件)~~ 🆗 TODO：目前导出功能主要依赖前端，可以修改为后端主导
- ~~权限管理中菜单编辑时选择图标（现在只能输入图标字符串）~~ 🆗
- 后端日志切割
- ~~后台修改背景图片，博客配置等~~ 🆗
- ~~后端的 IP 地址检测 BUG 待修复~~ 🆗
- ~~博客前台适配移动端~~ 🆗
- ~~文章详情, 目录锚点跟随~~ 🆗
- ~~邮箱注册 + 邮件发送验证码~~ 🆗
- 修改测试环境的数据库为 SQLite3，方便运行

### 1.4 后续功能想法

- 黑夜模式
- 前台收缩侧边信息功能
- 说说
- 相册
- 音乐播放器
- 鼠标左击特效
- 看板娘
- 第三方登录: QQ、微信、Github ...
- 评论时支持选择表情，参考 Valine
- 单独部署：前后端 + 环境
- 重写单元测试，目前的单元测试是早期版本，项目架构更改后，无法跑通
- 前台首页搜索集成 ElasticSearch
- 国际化?

