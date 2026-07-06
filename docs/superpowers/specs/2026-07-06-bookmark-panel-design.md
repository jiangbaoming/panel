# 个人书签导航面板 - 前端设计文档

## 概述
基于 Go 后端 API 的个人书签管理前端，采用 Vue 3 + Vite (JS) + Element Plus 技术栈。设计参考 sun-panel 风格，大图标卡片+网格布局，聚焦书签浏览效率。

## 技术栈
- Vue 3 (Composition API, JS)
- Vite 构建
- Vue Router 4 (SPA 路由)
- Pinia 状态管理
- Axios HTTP 请求
- Element Plus UI 组件库
- vuedraggable@next (拖拽排序)

## 页面路由

| 路由 | 页面 | 说明 |
|------|------|------|
| `/login` | 登录 | 居中卡片式登录表单 |
| `/` | 主面板 | 书签首页 + 全部功能入口 |

所有二级功能（设置、图库、用户管理、新建/编辑分组和书签）均以对话框 (Dialog) 形式呈现，无独立页面。

## 主面板布局

### 顶栏
- 左侧空（或应用名）
- 右侧：用户名+头像 / 设置入口 / 图库入口 / 登出

### 内容区（主面板核心）
- **常驻书签行**：`pinned=1` 的书签，横向滚动卡片，带星标区分
- **分组列表**：每个分组为一个区块
  - 分组标题（Emoji 图标 + 分组名称）+ 「添加书签」按钮
  - 下方为书签卡片网格（自动换行）
- **新建分组按钮**：底部
- **悬浮搜索按钮**：右下角 FAB 按钮，点击或按 `Ctrl+K` 弹出搜索遮罩

### 书签卡片
- 尺寸约 120×110px，圆角 16px
- 白色背景，支持 `bg_color` / `icon_bg` 个性化颜色
- 中央大号图标（Emoji 或自定义上传图标）
- 下方小字书签名称
- 点击新窗口打开 URL
- Hover 出现编辑/删除操作按钮
- 支持拖拽排序

## 数据流

```
页面加载 → GET /api/settings/:userId → 设置背景/标题
        → GET /api/groups → 分组+书签列表
        → GET /api/pinned → 常驻书签
```

## 认证
- 登录获取 token，存入 localStorage
- Axios 拦截器自动携带 `Authorization: Bearer <token>`
- 401 响应自动跳转登录页
- 路由守卫保护所有页面（除 `/login`）

## 组件树

```
App.vue
├── LoginPage.vue (路由 /login)
├── MainPage.vue (路由 /)
│   ├── TopBar.vue (顶栏：头像+操作入口)
│   ├── PinnedBar.vue (常驻书签行)
│   ├── GroupSection.vue (分组区块)
│   │   └── BookmarkCard.vue (书签卡片)
│   ├── SearchOverlay.vue (搜索弹窗)
│   ├── SettingsDialog.vue (设置)
│   ├── ImageManageDialog.vue (图库)
│   ├── UserManageDialog.vue (用户管理)
│   ├── GroupFormDialog.vue (新建/编辑分组)
│   └── BookmarkFormDialog.vue (新建/编辑书签)
```

## 交互细节
- 卡片 Hover 显示编辑/删除图标
- 分组和书签支持拖拽排序（调用 PATCH reorder API）
- 搜索：防抖 300ms，从 `/api/search?q=` 获取结果
- 新建/编辑使用 Element Plus 对话框 (Dialog)
- 删除前确认弹窗
