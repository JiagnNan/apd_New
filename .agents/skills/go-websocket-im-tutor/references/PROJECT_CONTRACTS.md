# 实验与业务协议参考（需结合代码确认）

此文件是**教学目标的建议契约**，不是声明仓库已经实现。允许按现有业务规范调整字段名称，但不要改变可靠性和权限语义。首次进入各阶段时读取对应部分。

## 真实项目基础（2026-11-05 读取 GitHub master 后记录；开课必须再核实）

- Go Workspace `go.work` 管理：`api-gateway`、`server`、`commonlib`、`grpcclient`。
- Gateway 入口 `api-gateway/main.go`；`app.api` import `docs/main/auth.api`；goctl 生成 `internal/handler/routes.go`。
- `api-gateway/internal/svc/servicecontext.go` 历史内容只有 Config；`api-gateway/internal/logic/auth/*.go` 历史仍有 TODO。
- `server/rpc/auth/` 有 auth RPC 骨架；`goprobuf/auth/auth.proto` 来源于单独 Proto 文件；生成 Go 包实际在 `server/rpc/auth/authpb`。
- `server/apd.sql` 有 user_info/register_info/auth_code/auth_login 的示例建表 SQL；不意味着真实数据库已有这些表。
- IM RPC、WS Gateway、MySQL IM 数据表、正式认证链路目前无完成证据。不把建议架构误认为现状。

## 目标单机架构

客户端通过现有 HTTP API 登录，取得正式授权凭证 → 经同一 `api-gateway` 进程 `/ws` 建立连接 → Hub 管理在线 WebSocket → 业务操作调用 `server/rpc/im` → RPC 对会话成员及内容验证后写 MySQL → 返回 messageId → 网关返回 `stored` ACK 并向本地 Hub 投递。实时投递失败不抹掉已存数据，接收端通过历史补拉恢复。

初始 L01–L07 纯 Demo：只实现临时转发，不写 DB，因此只有 `accepted` 或 `delivery_queued` 之类的临时含义，绝不能写 `stored`。

## 建议消息 JSON 协议（先教材，L07/L10 各自演进）

客户端请求：

```json
{"type":"chat.send","requestId":"req-001","payload":{"conversationId":12,"clientMsgId":"uuid-001","content":"你好"}}
```

数据库持久化成功后，服务端可回：

```json
{"type":"chat.ack","requestId":"req-001","payload":{"clientMsgId":"uuid-001","messageId":501,"status":"stored"}}
```

接收方实时推送：

```json
{"type":"chat.message","payload":{"messageId":501,"conversationId":12,"senderId":1001,"content":"你好","createdAt":1783300000}}
```

错误响应：

```json
{"type":"error","requestId":"req-001","payload":{"code":"PERMISSION_DENIED","message":"not a conversation member"}}
```

必守约束：

- 服务端 `senderId` 必须从被验证的连接身份获取，不能采用客户端自报。
- `requestId` 是一次请求的关联 ID；`clientMsgId` 是跨网络重试稳定的幂等键；两者不能混淆。
- 保存成功不等于实时投递成功，也不等于阅读成功。额外 delivered/read ACK 需接收端协议定义与状态持久化，不是默认直接承诺。
- 防止高频恶意输入：设置 WebSocket 帧大小、单用户发送速率、内容长度、队列容量；不将明文 Token 或密钥输出日志。

## 建议数据库核心结构（L10 具体落地）

- `im_conversations(id, type, created_at, ...)`：会话；单聊唯一会话的判定规则需明确。
- `im_conversation_members(conversation_id, user_id, created_at, ...)`：会话成员，唯一 `(conversation_id,user_id)`，为 user_id 的反向查询按需索引。
- `im_messages(id, conversation_id, sender_id, client_msg_id, type, content, created_at, ...)`：消息。唯一 `(sender_id,client_msg_id)`，普通索引 `(conversation_id,id)`；对话授权按查 DB 的真实成员进行，不相信请求方传的目标身份。
- `clientMsgId` 应由客户端为每个逻辑消息生成并在重试时复用；若同一幂等键提交不同内容，应有拒绝或确定性处理规则，不可默默覆盖。
- `(conversation_id,id)` 用于稳定的历史记录排序；**不等于强保证并发事务提交顺序**。如产品要求严格的会话逻辑序号，另行设计每会话 seq 分配。
- 需要定义存储保留期、内容脱敏以及删除规则；第一版不默认保存敏感个人数据到调试日志。

## 连接与并发原则

- 单连接最多一位常规 reader + 一位常规 writer；业务逻辑通过 `send` Channel 排队写，Channel 必须有界。
- Hub 用 Mutex/RWMutex 或 Hub 单 owner 事件循环保护连接 Map；不在持锁期间执行网络 IO；连接删除要核对 connID，避免旧连接断开把新设备误删。
- 关闭时定义谁关闭 send Channel、何时从 Hub 移除、连接关闭的幂等性；不要从多个 Goroutine 对同一 Channel `close`。
- Ping/Pong 需要客户端正确响应，ReadDeadline/PongHandler 与实际读取循环配套，写超时独立设置，服务退出要停止 Ticker。
- 浏览器的 WebSocket 原生 API 不支持任意 Authorization 自定义 Header。正式验证用短期一次性 ticket 或合理的安全 Cookie + Origin；不要用长期 JWT URL 查询参数当生产默认。

## 当前架构不做什么

- 不先引入 MQ、Kafka、Nchan、OpenResty、Kubernetes、分布式在线路由。
- 不先拆 message.rpc、conversation.rpc、presence.rpc 三个服务；一开始一个 `im.rpc` 足够。
- 不把 Redis Pub/Sub 视为可靠队列：它不能取代消息持久化和断线补拉。
- 不因一个功能需要被动启动全部 RPC；跟当前调用链相关的服务按需运行。
