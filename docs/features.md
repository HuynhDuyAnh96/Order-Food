# Tính năng mở rộng

---

## 1. Đơn mang về (Takeaway)

### Mô tả
Khách đến quán nhưng muốn mua mang về. Nhân viên lên đơn bình thường, chọn loại "mang về" thay vì gắn số bàn.

### Khác biệt so với đơn ăn tại quán

| | Ăn tại quán (`dine_in`) | Mang về (`takeaway`) |
|---|---|---|
| `table_number` | 1–20 | 0 (không cần) |
| Flow trạng thái | `pending → preparing → ready → completed → paid` | `pending → preparing → ready → paid` |
| Thanh toán | Phải qua `completed` trước | Thu tiền bất kỳ lúc nào |
| Số đơn cùng lúc | Giới hạn theo bàn | Không giới hạn |

### Flow trạng thái

```
MANG VỀ:
pending → preparing → ready → paid ✅

ĂN TẠI QUÁN:
pending → preparing → ready → completed → paid ✅
```

### API

**Tạo đơn mang về:**
```http
POST /api/orders
Content-Type: application/json

{
  "order_type": "takeaway",
  "items": [
    {
      "id": "16",
      "title": "Ốc Hương Cháy Tỏi",
      "price": 50000,
      "quantity": 2
    }
  ],
  "total": 100000
}
```

**Tạo đơn ăn tại quán:**
```http
POST /api/orders
Content-Type: application/json

{
  "order_type": "dine_in",
  "table_number": 3,
  "items": [...],
  "total": 100000
}
```

**Thu tiền đơn mang về (không cần qua completed):**
```http
POST /api/orders/:orderId/pay
```

### Hiển thị trên KDS (bếp)

Response từ `GET /api/kitchen/board`:
```json
{
  "order_id": "order_xxx",
  "order_type": "takeaway",
  "table_number": 0,
  "status": "preparing",
  "items": [...]
}
```

- `order_type: "takeaway"` → frontend bếp hiển thị label **MANG VỀ** thay vì số bàn
- WebSocket events (`dish_cooking`, `dish_ready`, `order_completed`) đều có `order_type` để staff app phân biệt

---

## 2. Món không có trên menu (Món lậu)

### Mô tả
Khách gọi món đặc biệt không có trong menu chính. Nhân viên tự nhập tên và giá trực tiếp vào đơn, không cần thêm vào menu.

### Cách hoạt động

- Thêm `"is_custom": true` vào item
- `id` để trống (`""`) vì không có trong DB
- Tính tiền bình thường như các món khác
- Hiển thị rõ trên KDS để bếp nhận biết
- Không ảnh hưởng đến menu chính

### API

**Đơn có cả món menu + món lậu:**
```http
POST /api/orders
Content-Type: application/json

{
  "order_type": "dine_in",
  "table_number": 5,
  "items": [
    {
      "id": "16",
      "title": "Ốc Hương Cháy Tỏi",
      "price": 50000,
      "quantity": 2,
      "is_custom": false
    },
    {
      "id": "",
      "title": "Ghẹ hấp bia",
      "price": 120000,
      "quantity": 1,
      "is_custom": true,
      "note": "khách tự mang ghẹ vào"
    }
  ],
  "total": 220000
}
```

### Tính tổng hóa đơn

```
Tổng = Σ (item.price × item.quantity)  ← tất cả món, không phân biệt menu hay lậu
```

Ví dụ trên:
```
50,000 × 2 = 100,000  (Ốc Hương Cháy Tỏi)
120,000 × 1 = 120,000  (Ghẹ hấp bia - món lậu)
─────────────────────
Tổng: 220,000 đ
```

### Cấu trúc lưu trong DB

```json
{
  "item_id": "order_xxx_item_1",
  "dish_id": "",
  "title": "Ghẹ hấp bia",
  "price": 120000,
  "quantity": 1,
  "is_custom": true,
  "note": "khách tự mang ghẹ vào",
  "item_status": "pending"
}
```

### Kết hợp: Mang về + Món lậu

Hai tính năng hoạt động độc lập, có thể dùng cùng nhau:

```http
POST /api/orders

{
  "order_type": "takeaway",
  "items": [
    {
      "id": "22",
      "title": "Càng Ghẹ Cháy Tỏi",
      "price": 50000,
      "quantity": 1,
      "is_custom": false
    },
    {
      "id": "",
      "title": "Tôm sú hấp",
      "price": 150000,
      "quantity": 2,
      "is_custom": true,
      "note": "khách đặt riêng"
    }
  ],
  "total": 350000
}
```

---

## Tóm tắt thay đổi code

| File | Thay đổi |
|---|---|
| `domain/audit.go` | Thêm `OrderType`, field `order_type` vào `Order`/`KDSOrder`/request; thêm `is_custom`, `note` vào `OrderItem` |
| `service/order_service.go` | Validate table chỉ khi `dine_in`; PayOrder takeaway không cần qua `completed` |
| `service/kitchen_service.go` | Gán `order_type` vào `KDSOrder` |
| `http/kitchen_handler.go` | Thêm `order_type` vào tất cả WebSocket events |
