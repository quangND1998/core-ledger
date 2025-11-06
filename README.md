# 🧭 Database Migration Guide (Go + PostgreSQL)

Hướng dẫn sử dụng **golang-migrate** để quản lý database migration trong dự án Golang.

---

## 📦 Cài đặt công cụ migrate

Cài đặt binary `migrate` (CLI):

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

> 📝 Lưu ý:  
> Nếu lệnh `migrate` chưa được nhận trong terminal, thêm dòng sau vào `~/.bashrc` hoặc `~/.zshrc`:
> ```bash
> export PATH=$PATH:$(go env GOPATH)/bin
> ```

---

## 📁 Cấu trúc thư mục

```
project/
│
├── migrations/
│   ├── 20251106163542_create_users_table.up.sql
│   └── 20251106163542_create_users_table.down.sql
│
└── README.md
```

---

## 🧱 Tạo migration mới

Chạy lệnh sau để tạo file migration:

```bash
migrate create -ext sql -dir migrations create_users_table
```

Lệnh này sẽ tự động tạo 2 file:

```
migrations/
  20251106163542_create_users_table.up.sql
  20251106163542_create_users_table.down.sql
```

---

## ✍️ Viết nội dung migration

**File `.up.sql`** — chứa các lệnh khi migrate **lên (apply)**:
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT,
  email TEXT UNIQUE,
  created_at TIMESTAMP DEFAULT NOW()
);
```

**File `.down.sql`** — chứa các lệnh khi migrate **xuống (rollback)**:
```sql
DROP TABLE users;
```

---

## 🚀 Chạy migration

### ✅ Chạy migration lên
```bash
migrate -path migrations -database "postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=disable" up
```

Ví dụ:
```bash
migrate -path migrations -database "postgres://postgres:123@localhost:5432/testdb?sslmode=disable" up
```

### 🔄 Rollback 1 bước
```bash
migrate -path migrations -database "postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=disable" down 1
```

### 🧹 Rollback toàn bộ
```bash
migrate -path migrations -database "postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=disable" down
```

---

## ⚙️ Các lệnh hữu ích khác

| Lệnh | Mô tả |
|------|--------|
| `migrate version` | Xem version hiện tại của database |
| `migrate force <version>` | Cập nhật thủ công version nếu bị lỗi |
| `migrate goto <version>` | Chuyển database tới version cụ thể |
| `migrate up` | Chạy tất cả migration chưa chạy |
| `migrate down` | Rollback tất cả migration đã chạy |

---

## 🧩 Gợi ý: dùng Makefile để tiện hơn

Tạo file `Makefile`:

```makefile
MIGRATIONS_PATH=migrations
DB_URL=postgres://postgres:123@localhost:5432/testdb?sslmode=disable

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) $(name)
```

Chạy:
```bash
make migrate-new name=create_users_table
make migrate-up
make migrate-down
```

---

## ✅ Kết luận

- **Tạo file migration:**  
  `migrate create -ext sql -dir migrations <tên_migration>`
- **Chạy lên:**  
  `migrate up`
- **Rollback:**  
  `migrate down`

> 🧠 Mẹo: nên commit cả file migration vào Git để đồng bộ schema giữa các môi trường (dev, staging, prod).

---
