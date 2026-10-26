> **出品人**：ice7  
> **文档用途**：任何设备的 AI Agent（包括手机端微信/拍照云端工作模式）读到本文档后，即可作为系统的智能输入端，操作用户的错题本系统。  
> **核心定位**：AI Agent 是本系统最重要的**智能输入源**与**学情分析师**，负责承接「卷子拍照 → 智能批改 → 错题结构化清洗 → 标准 API 写入 → 统计洞察」。

## 一、服务连接与鉴权（AI Agent 必读）

- **服务根地址 (Base URL)**：`http://your-server-host:5173/cuoti`
- **强制鉴权 Header**：`X-API-Key: your-secret-api-key`
- **请求内容类型**：`Content-Type: application/json; charset=utf-8`
- **管理工作台**：`http://your-server-host:5173/cuoti/`（桌面端查看、详情抽屉与修改）
- **打印组卷页**：`http://your-server-host:5173/cuoti/print`（A4 纸质复习卷，带扫码二维码，答案不上纸）
- **学生免密扫码答题页**：`http://your-server-host:5173/cuoti/m/{题目ID}`（免鉴权，用于手机端自主练习）

> [!IMPORTANT]
> **路径注意**：服务端部署在 `/cuoti` 子路径下，所有 API 接口均必须包含 `/cuoti` 前缀，例如：
> `http://your-server-host:5173/cuoti/api/v1/questions`。若省略 `/cuoti` 会导致 404！

---

## 二、AI Agent 核心调用接口一览

| 请求方式 | 完整路径 | 说明 |
|---|---|---|
| **POST** | `/cuoti/api/v1/images` | **第一步：上传原题/试卷截图**（multipart/form-data） |
| **POST** | `/cuoti/api/v1/questions` | **第二步：错题结构化录入**（支持单题对象或多题数组批量写入） |
| **GET** | `/cuoti/api/v1/questions` | **错题检索**（支持 `subject`, `status`, `kp`, `q`, `limit`, `offset`） |
| **GET** | `/cuoti/api/v1/questions/:id` | **查询单题详情**（含该题历史重做记录） |
| **PUT** | `/cuoti/api/v1/questions/:id` | **修改/纠正错题**（全量更新题目字段） |
| **DELETE** | `/cuoti/api/v1/questions/:id` | **删除错题** |
| **POST** | `/cuoti/api/v1/questions/:id/reviews` | **记录重做结果**（`{"result":"对"}` 或 `{"result":"错"}`） |
| **GET** | `/cuoti/api/v1/stats` | **学情看板与薄弱点统计**（知识点排行、掌握率、学科分布） |

---

## 三、AI Agent 写输入标准两步流（重点）

当用户发出「收集到错题本」、「帮我录入这几道错题」、「批改这张卷子并存错题」时，严格按照以下两步执行：

### 步骤 1：上传原题图片（若有图片）

若错题包含几何图形、函数图象、长阅读材料或手写解题过程，先将图片持久化：

- **URL**: `POST http://your-server-host:5173/cuoti/api/v1/images`
- **Header**: `X-API-Key: your-secret-api-key`
- **Form-Data 参数**：
  - `file`: 图片二进制文件（支持 png / jpg / jpeg / webp）
  - `folder`: 分类归档目录，格式为 `学生-学科-学期`（如 `ice7-数学-八上`；未指明学生时用 `通用-数学`）
  - `name`: 文件命名，格式为 `YYYYMMDD-卷名`（如 `20261026-期中模拟卷`；同名系统会自动追加 `-2`、`-3`）
- **返回响应**：
  ```json
  {
    "url": "images/ice7-数学-八上/20261026-期中模拟卷.png"
  }
  ```
  > 将返回的 `url` 赋值给步骤 2 错题入库的 `image_url` 字段。

---

### 步骤 2：错题结构化入库 JSON 契约

- **URL**: `POST http://your-server-host:5173/cuoti/api/v1/questions`
- **Header**:
  - `X-API-Key: your-secret-api-key`
  - `Content-Type: application/json`
- **Body 规范**（支持传单个 JSON 对象，或多道错题组成的 JSON 数组 `[...]`）：

```json
[
  {
    "subject": "数学",
    "source": "2026年秋八年级期中模拟卷一 第16题",
    "exam_date": "2026-10-26",
    "stem": "如图，在平面直角坐标系中，反比例函数 y = k/x (k > 0) 的图象经过菱形 OABC 的顶点 A 和对角线交点 P，若菱形面积为 12，求反比例函数的解析式。",
    "options": [],
    "answer": "y = 3/x",
    "my_answer": "y = 6/x",
    "analysis": "【核心思路】连接 OB、AC 交于点 P。根据菱形的性质，对角线互相垂直平分，点 P 为 AC 和 OB 的中点。设 P 点坐标为 (a, b)，则 k = ab。根据菱形面积公式 S = 1/2 * OB * AC = 12，结合几何坐标关系可推导出 ab = 3，故解析式为 y = 3/x。【易错点】误将菱形面积直接当成矩形面积计算，导致求出 k = 6。",
    "knowledge_points": ["反比例函数", "菱形的性质", "反比例函数与几何综合"],
    "difficulty": 4,
    "error_reason": "概念不清",
    "image_url": "images/ice7-数学-八上/20261026-期中模拟卷.png",
    "status": "待重做"
  },
  {
    "subject": "英语",
    "source": "Unit 4 单元语法限时练习 第3题",
    "exam_date": "2026-10-26",
    "stem": "Neither of the two books ______ interesting, so I didn't finish reading them.",
    "options": ["A. is", "B. are", "C. were", "D. was"],
    "answer": "D",
    "my_answer": "B",
    "analysis": "【考点精析】1. neither of + 复数名词作主语时，谓语动词在正式文体中用单数（即 is 或 was）；2. 由后半句 didn't finish 可知句子时态为一般过去时，因此需选用过去式 was，故选 D。【陷阱提醒】注意时态呼应以及主谓一致原则。",
    "knowledge_points": ["主谓一致", "不定代词", "一般过去时"],
    "difficulty": 3,
    "error_reason": "审题错误",
    "status": "待重做"
  }
]
```

---

## 四、入库字段约束与规范表（AI 批改打标核心准则）

| 字段名 | 类型 | 必填 | 取值范围与质量规范 |
|---|---|---|---|
| `subject` | String | **是** | **严格限定 9 门科目**：`数学`、`物理`、`化学`、`英语`、`语文`、`道法`、`历史`、`地理`、`生物` |
| `stem` | String | **是** | 题目题干完整描述（保留原始条件，支持 Markdown 格式与公式排版） |
| `answer` | String | **是** | 正确标准答案（选择题填选项字母如 `C`；填空或解答题写最终化简结果） |
| `options` | Array | 否 | 选择题选项数组（如 `["A. 2", "B. 3", "C. 4", "D. 5"]`）；非选择题传空数组 `[]` |
| `my_answer` | String | 否 | 学生的原始错误答案（用于错因对比，如 `A` 或算错的过程值） |
| `analysis` | String | 否 | **AI 详尽解析**：包含【解题思路】、【标准推导过程】与【易错陷阱/避免方法】 |
| `knowledge_points` | Array | **推荐** | **初中教材标准知识点**数组（如 `["反比例函数", "待定系数法"]`）；**必须规范统一用词**，严禁使用“不会解方程”等口语，统计看板直接依赖该字段聚合 |
| `error_reason` | String | 否 | **限定 5 大归因**：`概念不清`、`计算失误`、`审题错误`、`方法不会`、`粗心马虎` |
| `difficulty` | Integer| 否 | 难度评级 `1` 到 `5`（默认 `3`） |
| `source` | String | 否 | 试卷或练习来源（如 `2026学年八上数学期末测试 第18题`） |
| `exam_date` | String | 否 | 试卷或纠错日期，格式 `YYYY-MM-DD`（默认当天） |
| `image_url` | String | 否 | 原图相对路径（通过 `/api/v1/images` 上传获取） |
| `status` | String | 否 | 初始状态（默认留空或填 `待重做`；云端会自动根据做题结果流转：连续 2 次做对自动变 `已掌握`） |
| `ima_media_id` | String | 否 | 腾讯 ima 个人知识库关联的媒体 ID（用于跨平台搜索） |

---

## 五、AI Agent 标准交互场景与提示词示例

### 场景 1：拍照批改后自动录入错题
1. **AI 批改并汇总**：识别卷子后，先在对话中输出结构化清单，例如：
   > 发现 2 道错题，已完成诊断，准备录入错题本：
   > 1. 【数学】第 16 题（考点：反比例函数几何综合 · 概念不清）
   > 2. 【英语】第 3 题（考点：主谓一致、不定代词 · 审题错误）
   > 是否立即录入系统？
2. **用户下达确认指令**（如「录入」、「收集到错题本」）：
3. **AI 执行写入**：调用 `POST /cuoti/api/v1/questions`，写入成功后反馈：
   > ✅ 已成功录入 2 道错题（编号 #9、#10）。
   > 当前题库共有 10 道错题待巩固，可随时在网页端查看或打印复习卷：
   > 🔗 [错题本工作台](http://your-server-host:5173/cuoti/) · [组卷打印](http://your-server-host:5173/cuoti/print)

### 场景 2：学情诊断与出卷推荐
1. **用户询问**：「我哪部分最薄弱？」
2. **AI 执行查询**：调用 `GET http://your-server-host:5173/cuoti/api/v1/stats`
3. **AI 诊断口径**：
   - 提取 `kp_stats` 中 `pending_count`（待重做错题数）最高、且重做通过率最低的知识点；
   - 给出复习建议：「你在【反比例函数】和【一次函数】板块积压错题最多，建议点击进入 [打印组卷页](http://your-server-host:5173/cuoti/print) 筛选这 2 个知识点，抽取 5-10 题进行纸笔重做。」

---

## 六、错误处理与稳健性规范

1. **401 Unauthorized**：
   - 原因：请求头缺少 `X-API-Key` 或密钥不匹配。
   - 处置：确保请求 Header 中包含 `X-API-Key: your-secret-api-key`。
2. **404 Not Found**：
   - 原因：URL 未拼接 `/cuoti` 前缀。
   - 处置：检查请求 URL 是否为 `http://your-server-host:5173/cuoti/api/v1/...`。
3. **400 Bad Request**：
   - 原因：JSON 格式错误、必填项缺少（如未填 `subject` 或 `stem`）、或科目不在 9 大科目列表中。
   - 处置：严格按照第四节字段表校验字段后再发请求。
