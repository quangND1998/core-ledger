# Hướng dẫn Test Backoff

## Cách kiểm tra backoff hoạt động đúng

### 1. Quan sát Logs

Khi job retry, bạn sẽ thấy logs như sau:

#### Với custom backoff:
```
🔄 [Backoff] Job data:process - Retry #1: Using custom backoff from array[0] = 2 seconds
📦 [Job] Handling DataProcessJob: Type=test, Action=fail | Attempt=2/4 | Backoff=[2 5 10] seconds
❌ [Job] Forcing failure for testing (will retry with backoff)
🔄 [Backoff] Job data:process - Retry #2: Using custom backoff from array[1] = 5 seconds
📦 [Job] Handling DataProcessJob: Type=test, Action=fail | Attempt=3/4 | Backoff=[2 5 10] seconds
```

#### Với default backoff (không set backoff):
```
🔄 [Backoff] Job data:process - Retry #1: Using default exponential backoff = 1s
🔄 [Backoff] Job data:process - Retry #2: Using default exponential backoff = 2s
🔄 [Backoff] Job data:process - Retry #3: Using default exponential backoff = 4s
```

### 2. Test bằng cách dispatch job với Action="fail"

```go
// Trong service hoặc handler nào đó
job := jobs.NewDataProcessJob("test", "fail", map[string]interface{}{
    "test": "backoff",
})
job.SetBackoff([]int{2, 5, 10}) // Custom backoff: 2s, 5s, 10s
job.SetRetry(3) // Cho phép retry 3 lần

dispatcher.Dispatch(job)
```

### 3. Đo thời gian thực tế

**Cách 1: Dùng timestamps trong logs**
- Ghi lại thời gian khi job fail lần đầu
- Ghi lại thời gian khi job retry
- Khoảng cách = thời gian retry - thời gian fail
- So sánh với giá trị backoff đã set

**Cách 2: Dùng Redis để xem task**
```bash
redis-cli
> KEYS asynq:*
> HGETALL asynq:task:<task_id>
> TTL asynq:task:<task_id>  # Xem thời gian còn lại
```

**Cách 3: Dùng asynqmon (Web UI)**
```bash
# Cài đặt
go install github.com/hibiken/asynq/tools/asynqmon@latest

# Chạy
asynqmon --redis-addr=localhost:6379
# Mở browser: http://localhost:8080
```

### 4. Test Cases

#### Test Case 1: Custom backoff đủ giá trị
```go
job.SetBackoff([]int{1, 2, 3})
job.SetRetry(3)
// Kỳ vọng: Retry 1 sau 1s, Retry 2 sau 2s, Retry 3 sau 3s
```

#### Test Case 2: Custom backoff thiếu giá trị
```go
job.SetBackoff([]int{1, 2})
job.SetRetry(5)
// Kỳ vọng: Retry 1 sau 1s, Retry 2 sau 2s, Retry 3-5 sau 2s (dùng giá trị cuối)
```

#### Test Case 3: Không có backoff (dùng default)
```go
// Không set backoff
job.SetRetry(5)
// Kỳ vọng: Exponential backoff: 1s, 2s, 4s, 8s, 16s, 30s (max)
```

### 5. Kiểm tra trong code

Bạn có thể thêm breakpoint hoặc log trong hàm `retryDelayFuncWithJobBackoff` để xem:
- Payload có được parse đúng không
- Backoff array có được đọc đúng không
- Index có đúng không
- Delay có được tính đúng không

### 6. Ví dụ test script

Tạo một endpoint test hoặc command để dispatch job test:

```go
// Test endpoint
func TestBackoffHandler(c *gin.Context) {
    job := jobs.NewDataProcessJob("backoff_test", "fail", map[string]interface{}{
        "test_id": time.Now().Unix(),
    })
    job.SetBackoff([]int{2, 5, 10})
    job.SetRetry(3)
    
    if err := dispatcher.Dispatch(job); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "message": "Test job dispatched",
        "backoff": []int{2, 5, 10},
        "retry": 3,
    })
}
```

Sau đó quan sát logs trong worker để xác nhận backoff hoạt động đúng.

