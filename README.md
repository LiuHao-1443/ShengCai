# ShengCai

如有使用上的问题，请发 Issue 或直接联系我，谢谢，欢迎 Start~

## 项目目的

本项目旨在实现一个系统，自动从飞书表格中采集文章数据，并将其展示在一个新的网页上。网页会持续更新，确保信息的及时性和准确性。同时，项目将处理加密的飞书表格，并支持自动更新功能。  

## 项目结构设计

1. 前端 (React)
* 负责展示从飞书表格中采集的信息。
* 页面包含：飞书表格编号、更新日志、文章摘要、文章链接、文章关键字（标签）。
* 摘要和标签通过 OpenAI API 获取。
* 支持自动更新功能，确保网页内容与原表格数据同步。

2. 后端 (Go)
* 使用nunu脚手架开发。
* 提供 API 接口供前端获取数据。
* 实现数据采集、处理和存储。
* 支持轮询任务，定期从飞书表格中采集数据。

3. 存储 (MySQL)
* 存储飞书表格数据、文章摘要、文章链接和关键字等信息。
* 数据表设计包括：sheet_info 表和 article_data 表。
* 数据采集 (RPA + Go)

4. 数据采集 (RPA + Go)
* 通过 RPA 脚本自动从飞书表格中提取数据。
* 使用 Go 编写轮询任务，定期更新数据库中的数据。

## Use Guide

1. 获取飞书 APP_ID 以及 APP_SECRET，获取方式参考 https://open.feishu.cn/document/server-docs/api-call-guide/calling-process/get-access-token  
2. 申请必需的飞书 API 权限，必需权限包括查看新版文档，查看云空间中文件元数据，查看、评论、编辑和管理电子表格，查看、评论和导出电子表格，查看知识空间节点信息，查看、编辑和管理知识库，查看知识库，申请方式参考 https://open.feishu.cn/document/server-docs/application-scope/introduction  
3. 新建一个个人电子表格，并记录 sheet_id，用于 rpa 采集加密的电子表格，红色方框框起来的文字即为 sheet_id
   ![image](https://github.com/user-attachments/assets/e3e4ca09-6ad7-4de1-824c-b3c42670b37f)
4. 导入 ./backend/sql 下的两个数据表到 MySQL
5. 修改 ./backend/config/local.yaml 配置文件，包括 MySQL 连接等，具体如图所示
  ![image](https://github.com/user-attachments/assets/32a5dadf-40f7-4bba-a2c2-c4b8dbbb0a69)
  ![image](https://github.com/user-attachments/assets/b782fc4a-e486-459c-b51b-51a3ab850527)
6. 大模型使用为 DeepSeek，API Key 请参考官网申请，https://www.deepseek.com/
7. 对于阶段一，参数配置好后可直接启动服务，Go 后端使用 nunu 脚手架，需要安装 nunu cli，安装方式请参考 https://github.com/go-nunu/nunu?tab=readme-ov-file#nunu-cli
8. Go 后端启动方式，启动好后自动开始进行数据采集
   ```bash
   cd ./backend
   go install github.com/go-nunu/nunu@latest
   go mod tidy
   nunu run ./cmd/server/main.go
   ```
9. 前端启动方式，node 版本 v20.15.1 pnpm 版本 9.6.0，如果需要修改代理，请修改 ./frontend/vite.config.ts 中的配置
   ```bash
   cd ./frontend/
   pnpm install
   pnpm dev
   ```
10. 访问地址，查看效果 http://127.0.0.1:5173/?sheet_id=N5Wts8V9Wh3gXJtyxPvcDbMZnJc
  ![image](https://github.com/user-attachments/assets/08e1dcc2-6d3b-4a89-96dd-16614f6bdf2f)
11. 对于阶段二，需要结构影刀 RPA 和个人电子表格进行采集，影刀 RPA 配置链接
12. 启动影刀 RPA，会自动将加密的电子表格内容拷贝到个人电子表格中，以方便采集，前置工作需要在浏览器中登录飞书以及输入加密电子表格的密码
13. 重复 5-9 的步骤即可完成采集任务

## TODO
1. 容器化部署
2. 影刀 RPA自动判断登录及输入密码
