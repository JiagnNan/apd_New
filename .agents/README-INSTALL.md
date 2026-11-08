# Codex WebSocket / IM 教学 Skill v2 安装说明

此版本基于原版 `go-websocket-im-tutor` 升级，保持同一个 Skill 名称与 14 课顺序，但加入每项具体实现任务、阶段 Gate、协议约束和持久进度。**解压覆盖原 `.agents/skills/go-websocket-im-tutor` 文件夹即可，不要同时保留同名 v1/v2 Skill。**

目标目录：

```text
apd_New/
└── .agents/
    └── skills/
        └── go-websocket-im-tutor/
            ├── SKILL.md
            ├── references/
            │   ├── curriculum.md
            │   ├── EXECUTION_PLAN.md
            │   └── PROJECT_CONTRACTS.md
            └── templates/
                └── PROGRESS.md
```

Codex 启动时执行：

```text
$go-websocket-im-tutor 请先审查完整执行计划及本仓库当前代码，初始化学习进度，告诉我 G0 要核对什么。不要写业务代码，不要跳到下一任务。
```

首次确认后开始：

```text
$go-websocket-im-tutor 从当前未完成任务开始授课。我自己编写代码；只教一个任务，给验收标准。
```

中途恢复：

```text
$go-websocket-im-tutor 读取 learning/websocket/PROGRESS.md，从第一个未通过任务继续。不要重复已通过课程。
```

验收代码：

```text
$go-websocket-im-tutor 审查我当前 Lxx-Tyy 的实现，指出错误和失败路径，执行或告诉我验收命令；未通过不要解锁下一任务。
```

核对计划：

```text
$go-websocket-im-tutor 只显示我尚未完成的任务 ID、完成标准、当前障碍，不要临时扩展课程。
```

提示：计划不以时间长度决定是否毕业。必须有代码、运行证据和设计理解。学习进度模板在 Skill 目录，正式进度默认在项目根目录 `learning/websocket/PROGRESS.md`。
