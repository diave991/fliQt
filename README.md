# fliQt

`fliQt` 是一個使用 Go 1.24 和 Gin 框架開發的簡易人資後端 API，整合 MySQL、Redis、GORM 自動遷移及種子資料功能，提供員工管理、請假管理與打卡功能。

## 功能特色

- **員工管理 (Employee)**  
  - 建立員工：新增員工資料  
  - 查詢員工：取得單筆或全部員工資料  
  - 更新員工：編輯員工資訊  
  - 刪除員工：軟刪除 (status = 0) 保留歷史

- **請假管理 (Leave)**  
  - 提交請假：新增請假記錄  
  - 查詢請假：可列出請假紀錄，並同時顯示員工姓名與職位

- **打卡功能 (Attendance)**  
  - 上下班打卡：新增打卡記錄 (IN/OUT)  
  - 查詢打卡：依員工與日期範圍查詢，並顯示員工姓名與職位

## 快速開始

1. **取得程式碼**  
   ```bash
   git clone <https://github.com/diave991/fliQt.git>
   cd fliQt
   ```

2. **設定環境參數**  
   複製範本並編輯 `.env`：  
   ```bash
   cp .env.example .env
   ```
   編輯 `.env` 內容：
   ```dotenv
   # MySQL
   DB_HOST=db
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=123456
   DB_NAME=fliqt

   # Redis
   REDIS_ADDR=redis:6379
   REDIS_PASSWORD=
   REDIS_DB=0

   # 服務
   SERVER_PORT=8080
   ```

3. **一鍵啟動 (Docker Compose)**  
   ```bash
   make run
   ```
   服務預設運行於 `http://localhost:8080`

4. **本地開發**  
   ```bash
   go mod tidy
   go test ./services ./repositories ./controllers -v
   go build -o fliQt
   ./fliQt
   ```

## API 路徑與 IN/OUT 範例

### 員工管理 (Employee)

#### 新增員工

- **路徑**：`POST /api/v1/employees`  
- **Request (IN)**  
  ```json
  {
    "name": "王小明",
    "position": "資深後端工程師",
    "contact": "xiaoming.wang@example.com",
    "salary": 1200000
  }
  ```
- **Response (OUT)**  
  ```json
  {
    "id": 1,
    "name": "王小明",
    "position": "資深後端工程師",
    "contact": "xiaoming.wang@example.com",
    "salary": 1200000,
    "status": 1,
    "created_at": "2025-05-09T12:00:00Z",
    "updated_at": "2025-05-09T12:00:00Z"
  }
  ```

#### 查詢所有員工

- **路徑**：`GET /api/v1/employees`  
- **Response (OUT)**  
  ```json
  [
    {
      "id": 1,
      "name": "王小明",
      "position": "資深後端工程師",
      "contact": "xiaoming.wang@example.com",
      "salary": 1200000,
      "status": 1,
      "created_at": "...",
      "updated_at": "..."
    },
    {
      "id": 2,
      "name": "李麗華",
      "position": "前端工程師",
      "contact": "li.lihua@example.com",
      "salary": 1000000,
      "status": 1,
      "created_at": "...",
      "updated_at": "..."
    }
  ]
  ```

#### 查詢單筆員工

- **路徑**：`GET /api/v1/employees/{id}`  
- **Response (OUT)**  
  ```json
  {
    "id": 1,
    "name": "王小明",
    "position": "資深後端工程師",
    "contact": "xiaoming.wang@example.com",
    "salary": 1200000,
    "status": 1,
    "created_at": "...",
    "updated_at": "..."
  }
  ```

#### 更新員工

- **路徑**：`PUT /api/v1/employees/{id}`  
- **Request (IN)**  
  ```json
  {
    "name": "王小明",
    "position": "技術主管",
    "contact": "xiaoming.wang@example.com",
    "salary": 1300000
  }
  ```
- **Response (OUT)**  
  ```json
  {
    "id": 1,
    "name": "王小明",
    "position": "技術主管",
    "contact": "xiaoming.wang@example.com",
    "salary": 1300000,
    "status": 1,
    "created_at": "...",
    "updated_at": "..."
  }
  ```

#### 刪除員工 (軟刪除)

- **路徑**：`DELETE /api/v1/employees/{id}`  
- **Response (OUT)**  
  無內容 (204 No Content)

### 請假管理 (Leave)

#### 新增請假

- **路徑**：`POST /api/v1/leaves`  
- **Request (IN)**  
  ```json
  {
    "employee_id": 1,
    "start_date": "2025-06-01T09:00:00Z",
    "end_date": "2025-06-03T18:00:00Z",
    "reason": "家用因素"
  }
  ```
- **Response (OUT)**  
  ```json
  {
    "id": 1,
    "employee_id": 1,
    "start_date": "2025-06-01T09:00:00Z",
    "end_date": "2025-06-03T18:00:00Z",
    "reason": "家用因素",
    "created_at": "...",
    "updated_at": "..."
  }
  ```

#### 查詢所有請假 (含員工姓名/職位)

- **路徑**：`GET /api/v1/leaves-with-staff`  
- **Response (OUT)**  
  ```json
  [
    {
      "id": 1,
      "employee_id": 1,
      "employee_name": "王小明",
      "employee_position": "技術主管",
      "start_date": "2025-06-01T09:00:00Z",
      "end_date": "2025-06-03T18:00:00Z",
      "reason": "家用因素",
      "created_at": "...",
      "updated_at": "..."
    }
  ]
  ```

### 打卡功能 (Attendance)

#### 打卡 (IN/OUT)

- **路徑**：`POST /api/v1/attendance`  
- **Request (IN)**  
  ```json
  {
    "employee_id": 1,
    "type": "IN"
  }
  ```
- **Response (OUT)**  
  ```json
  {
    "id": 1,
    "employee_id": 1,
    "type": "IN",
    "timestamp": "2025-06-01T08:30:00Z",
    "created_at": "...",
    "updated_at": "..."
  }
  ```

#### 查詢員工當日打卡 (含員工姓名/職位)

- **路徑**：`GET /api/v1/attendance/by_employee?employee_id={id}&date={YYYY-MM-DD}`  
- **Response (OUT)**  
  ```json
  [
    {
      "id": 1,
      "employee_id": 1,
      "employee_name": "王小明",
      "employee_position": "技術主管",
      "type": "IN",
      "timestamp": "2025-06-01T08:30:00Z"
    },
    {
      "id": 2,
      "employee_id": 1,
      "employee_name": "王小明",
      "employee_position": "技術主管",
      "type": "OUT",
      "timestamp": "2025-06-01T17:45:00Z"
    }
  ]
  ```

### 出缺勤報表 (Reports)

#### 單一員工報表 (IN/OUT)

- **路徑**：`GET /api/v1/reports/{employee_id}`
  
- **Response (OUT)**
  ```json
  [
       { "date": "2025-05-05", "status": "absent" },
       { "date": "2025-05-06", "status": "present" },
       { "date": "2025-05-07", "status": "leave" },
       { "date": "2025-05-08", "status": "present" },
       { "date": "2025-05-09", "status": "present" },
       { "date": "2025-05-10", "status": "absent" },
       { "date": "2025-05-11", "status": "present" }
    ]
  ```

#### 查詢員工當日打卡 (含員工姓名/職位)

- **路徑**：`GET /api/v1/reports?page=1`
- **Response (OUT)**
  ```json
  {
  "page": 1,
  "data": [
    {
      "employee_id": 1,
      "employee_name": "王小明",
      "employee_position": "資深後端工程師",
      "report": [
        { "date": "2025-05-05", "status": "absent" },
        { "date": "2025-05-06", "status": "present" },
        { "date": "2025-05-07", "status": "leave" },
        { "date": "2025-05-08", "status": "present" },
        { "date": "2025-05-09", "status": "present" },
        { "date": "2025-05-10", "status": "absent" },
        { "date": "2025-05-11", "status": "present" }
      ]
    },
    {
      "employee_id": 2,
      "employee_name": "李麗華",
      "employee_position": "前端工程師",
      "report": [
        { "date": "2025-05-05", "status": "absent" },
        { "date": "2025-05-06", "status": "absent" },
        { "date": "2025-05-07", "status": "absent" },
        { "date": "2025-05-08", "status": "present" },
        { "date": "2025-05-09", "status": "present" },
        { "date": "2025-05-10", "status": "present" },
        { "date": "2025-05-11", "status": "absent" }
      ]
    }
    // ...更多筆資料...
   ]
  ```
- **狀態判斷**
- - **請假(leave)：當天有請假區間**
- - **出勤(present)：無請假但有打卡記錄**
- - **缺席(absent)：既無請假也無打卡**

## Makefile 指令

```bash
make run    # 啟動並同步執行 Docker Compose
make stop   # 停止所有容器
make down   # 移除所有容器與網路
make test   # 執行 services 目錄單元測試
```

© 2025 fliQt. MIT License.
