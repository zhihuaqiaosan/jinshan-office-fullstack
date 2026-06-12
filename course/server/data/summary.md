# 前端开发学习总结

### 一、HTML/CSS 基础

#### 1. HTML5 核心

| 分类       | 知识点                                                       |
| :--------- | :----------------------------------------------------------- |
| 语义化标签 | `<header>`、`<nav>`、`<main>`、`<article>`、`<section>`、`<footer>` |
| 表单增强   | `input` 新类型（email、tel、number、date）、`datalist`、表单验证属性 |
| 多媒体     | `<video>`、`<audio>`、`<canvas>`、`<svg>`                    |
| 存储       | localStorage、sessionStorage、IndexedDB                      |
| 其他       | Web Worker、WebSocket、Geolocation、拖拽API                  |

#### 2. CSS3 核心

| 分类     | 知识点                                                       |
| :------- | :----------------------------------------------------------- |
| 选择器   | 属性选择器、伪类（`:nth-child`、`:not`）、伪元素（`::before`、`::after`） |
| 盒模型   | `box-sizing`（border-box / content-box）、margin 塌陷与合并  |
| 布局     | Flex、Grid、浮动、定位（relative/absolute/fixed/sticky）、多列布局 |
| 响应式   | 媒体查询（`@media`）、视口单位（vw/vh）、`rem/em` 适配       |
| 视觉效果 | 过渡（transition）、动画（`@keyframes` + animation）、变换（transform） |
| 预处理器 | SCSS/Sass（变量、嵌套、混合、继承）、Less、PostCSS           |

------

### 二、JavaScript 核心

#### 1. 基础语法

| 分类        | 知识点                                                       |
| :---------- | :----------------------------------------------------------- |
| 数据类型    | 原始类型（string、number、boolean、null、undefined、symbol、bigint）、引用类型（object） |
| 作用域      | 全局作用域、函数作用域、块级作用域（let/const）              |
| 闭包        | 定义、用途（模块化、数据私有化）、内存泄漏风险               |
| 原型与继承  | 原型链、`prototype`、`__proto__`、`class` 语法糖             |
| 异步编程    | 回调函数、Promise（all/race/any/allSettled）、async/await、事件循环（宏任务/微任务） |
| ES6+ 新特性 | 解构赋值、展开运算符、箭头函数、模板字符串、可选链（`?.`）、空值合并（`??`） |

#### 2. 核心概念进阶

- **this 指向**：默认绑定、隐式绑定、显式绑定（call/apply/bind）、箭头函数 this
- **深拷贝与浅拷贝**：`Object.assign()`、扩展运算符、`JSON.parse(JSON.stringify())` 的缺陷、递归实现深拷贝
- **事件机制**：事件冒泡与捕获、事件委托、阻止默认行为与冒泡
- **模块化**：CommonJS、ES Module（import/export）、AMD/CMD

------

### 三、前端框架

#### 1. Vue.js

| 分类      | 知识点                                                       |
| :-------- | :----------------------------------------------------------- |
| 核心特性  | 响应式原理（Object.defineProperty / Proxy）、虚拟DOM、Diff算法 |
| 指令      | v-bind、v-model、v-for、v-if/v-show、v-on、v-html            |
| 组件通信  | props / emit、provide / inject、ref / $parent、Vuex / Pinia  |
| 生命周期  | beforeCreate → created → beforeMount → mounted → beforeUpdate → updated → beforeUnmount → unmounted |
| 组合式API | ref、reactive、computed、watch、watchEffect、setup 语法糖    |
| 路由      | Vue Router（history/hash模式、路由守卫、动态路由、嵌套路由） |
| 生态      | Pinia（状态管理）、Vite（构建工具）、Nuxt.js（SSR）          |

#### 2. React

| 分类     | 知识点                                                       |
| :------- | :----------------------------------------------------------- |
| 核心特性 | JSX、虚拟DOM、Fiber架构                                      |
| 组件     | 函数组件 vs 类组件、受控组件 vs 非受控组件                   |
| Hooks    | useState、useEffect、useContext、useReducer、useCallback、useMemo、自定义Hook |
| 状态管理 | Redux（action/reducer/store）、Zustand、MobX                 |
| 路由     | React Router（useParams、useNavigate、路由守卫实现）         |
| 生态     | Next.js（SSR/SSG）、React Native（移动端）                   |

#### 3. 框架对比

| 维度     | Vue                 | React      |
| :------- | :------------------ | :--------- |
| 学习曲线 | 较平缓              | 中等       |
| 数据流   | 双向绑定（v-model） | 单向数据流 |
| 模板语法 | 模板 + 指令         | JSX        |
| 生态     | 官方维护较好        | 社区更庞大 |

------

### 四、工程化与构建工具

| 工具              | 核心知识点                                                   |
| :---------------- | :----------------------------------------------------------- |
| Webpack           | entry/output、loader（babel-loader、css-loader）、plugin（HtmlWebpackPlugin）、代码分割（SplitChunks）、热更新（HMR） |
| Vite              | 基于ESM的开发服务器、预构建依赖、HMR、Rollup 打包            |
| Babel             | 语法转换（@babel/preset-env）、Polyfill（core-js）、插件开发 |
| ESLint + Prettier | 代码规范与格式化、配置文件（`.eslintrc.js`）、常用规则（airbnb、standard） |
| Git               | 常用命令、分支策略（Git Flow）、Husky + lint-staged（提交前检查） |

------

### 五、性能优化

#### 1. 加载性能

- 资源压缩（HTML/CSS/JS 压缩、图片压缩）
- 路由懒加载 + 组件异步加载
- 图片优化（WebP、懒加载、响应式图片 `srcset`）
- CDN 加速
- 预加载（`<link rel="preload">`）与预连接（`dns-prefetch`）

#### 2. 运行时性能

- 减少重排与重绘（批量修改DOM、使用 `transform` 代替 `top/left`）
- 虚拟滚动（长列表优化）
- 防抖与节流（Debounce / Throttle）
- 事件委托
- Web Worker（复杂计算）

#### 3. 打包优化

- Tree Shaking（消除无用代码）
- 代码分割（按路由/按组件）
- 依赖体积分析（Webpack Bundle Analyzer）
- Gzip / Brotli 压缩

#### 4. 缓存策略

- HTTP 缓存（强缓存 `Cache-Control`、协商缓存 `ETag`）
- Service Worker（PWA）
- 浏览器本地存储缓存数据

------

### 六、网络与安全

#### 1. HTTP 知识

| 分类     | 知识点                                                       |
| :------- | :----------------------------------------------------------- |
| 请求方法 | GET、POST、PUT、DELETE、PATCH                                |
| 状态码   | 2xx（成功）、3xx（重定向）、4xx（客户端错误）、5xx（服务端错误） |
| 版本演进 | HTTP/1.1（队头阻塞）、HTTP/2（多路复用、头部压缩）、HTTP/3（QUIC） |
| 跨域     | CORS（简单请求/预检请求）、JSONP、代理转发                   |

#### 2. 前端安全

| 攻击类型             | 原理                   | 防护                          |
| :------------------- | :--------------------- | :---------------------------- |
| XSS（跨站脚本）      | 注入恶意脚本           | 输出转义、CSP（内容安全策略） |
| CSRF（跨站请求伪造） | 利用用户登录态伪造请求 | CSRF Token、SameSite Cookie   |
| SQL 注入             | 拼接SQL语句            | 参数化查询、输入校验          |
| 点击劫持             | 透明 iframe 诱导点击   | `X-Frame-Options: DENY`       |

------

### 七、浏览器原理

| 知识点          | 说明                                                         |
| :-------------- | :----------------------------------------------------------- |
| 渲染流程        | HTML解析 → CSS解析 → 布局树 → 绘制 → 合成（GPU加速）         |
| 重排（Reflow）  | 布局几何属性变化 → 重新计算布局                              |
| 重绘（Repaint） | 外观属性变化 → 重新绘制                                      |
| 事件循环        | 调用栈 → 宏任务（setTimeout）→ 微任务（Promise）→ 渲染       |
| 垃圾回收        | 标记清除、引用计数、内存泄漏排查                             |
| 存储机制        | Cookie（4KB）、localStorage（5-10MB）、IndexedDB（大量结构化数据） |

------

### 八、TypeScript（进阶）

| 知识点   | 示例                                                         |
| :------- | :----------------------------------------------------------- |
| 基础类型 | `string`、`number`、`boolean`、`array`、`tuple`、`enum`      |
| 高级类型 | `interface`、`type`、联合类型（`|`）、交叉类型（`&`）、泛型（`<T>`） |
| 类型工具 | `Partial`、`Required`、`Pick`、`Omit`、`Record`、`ReturnType` |
| 类型守卫 | `typeof`、`instanceof`、`in`、自定义守卫（`is`）             |
| 装饰器   | 类装饰器、方法装饰器（需开启 `experimentalDecorators`）      |

------

### 九、Node.js 基础

- 事件驱动、非阻塞 I/O
- 内置模块：`fs`、`path`、`http`、`crypto`、`child_process`
- 包管理：npm / yarn / pnpm（依赖管理、版本锁定）
- 常用框架：Express（中间件、路由）、Koa（async/await）

------

### 十、学习路线建议

```
基础阶段
├── HTML5 + CSS3（语义化、Flex/Grid、响应式）
├── JavaScript（ES6+、异步、DOM/BOM操作）
├── Git 版本控制
└── 网络基础（HTTP/HTTPS）

进阶阶段
├── 框架（Vue 或 React 选其一深入）
├── 工程化（Webpack/Vite、ESLint、模块化）
├── TypeScript
└── 小程序 / Uni-app（可选）

高阶阶段
├── 性能优化与监控
├── 浏览器原理与安全
├── Node.js 全栈
├── 微前端（qiankun、Module Federation）
└── 源码阅读（框架核心实现原理）
```
