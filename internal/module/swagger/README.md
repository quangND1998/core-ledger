# Swagger API Documentation

Module này cung cấp Swagger UI để xem và test API documentation.

## Cách truy cập

Sau khi chạy server, bạn có thể truy cập Swagger UI tại:

### Swagger UI (Mặc định)
```
http://localhost:8000/api/v1/swagger/
```

### ReDoc (Alternative UI)
```
http://localhost:8000/api/v1/swagger/redoc
```

### Swagger JSON file
```
http://localhost:8000/api/v1/swagger/swagger.json
```

## Tính năng

- **Swagger UI**: Giao diện tương tác để xem và test API
- **ReDoc**: Giao diện documentation đẹp hơn, dễ đọc hơn
- **Swagger JSON**: File OpenAPI specification dạng JSON

## Cập nhật documentation

Để cập nhật documentation:

1. Chỉnh sửa file `swagger.json` ở root directory
2. Copy file vào `internal/module/swagger/swagger.json`:
   ```bash
   copy swagger.json internal\module\swagger\swagger.json
   ```
3. Restart server

## Lưu ý

- Swagger UI sử dụng CDN từ unpkg.com, cần có kết nối internet
- File swagger.json được embed vào binary khi build, không cần file riêng khi deploy

