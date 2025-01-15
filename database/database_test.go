package database_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"github.com/stretchr/testify/assert"
)

// 用户模型
type User struct {
	Name    string  `json:"name" gorm:"type:varchar(50);index"`         // 用户名
	Email   string  `json:"email" gorm:"type:varchar(100);uniqueIndex"` // 邮箱
	Phone   string  `json:"phone" gorm:"type:varchar(20)"`              // 电话
	Balance float64 `json:"balance"`                                    // 账户余额
	// Orders  []Order `json:"orders" gorm:"foreignKey:UserID;references:ID"` // 一对多关系
	Orders []Order `json:"orders"`

	model.Base
}

// 订单模型
type Order struct {
	OrderNo     string         `json:"order_no" gorm:"type:varchar(32);uniqueIndex"` // 订单编号
	UserID      string         `json:"user_id" gorm:"type:varchar(32);index"`        // 用户ID
	User        User           `json:"user" gorm:"foreignKey:UserID;references:ID"`  // 关联用户
	Status      string         `json:"status" gorm:"type:varchar(20);index"`         // 订单状态
	Amount      float64        `json:"amount"`                                       // 订单金额
	PaymentTime model.GormTime `json:"payment_time"`                                 // 支付时间

	// 只是用来 Join 查询的
	OrderJoined `gorm:"-"`
	model.Base
}

type OrderJoined struct {
	OrderNo     string         `json:"order_no"`
	UserName    string         `json:"user_name"`
	UserEmail   string         `json:"user_email"`
	PaymentTime model.GormTime `json:"payment_time"`
	OrderCount  int            `json:"order_count"`
}

var users = []*User{
	{Name: "user01", Email: "user01@gmail.com", Base: model.Base{ID: "user01"}},
	{Name: "user02", Email: "user02@gmail.com", Base: model.Base{ID: "user02"}},
	{Name: "user03", Email: "user03@gmail.com", Base: model.Base{ID: "user03"}},
	{Name: "user04", Email: "user04@gmail.com", Base: model.Base{ID: "user04"}},
	{Name: "user05", Email: "user05@gmail.com", Base: model.Base{ID: "user05"}},
	{Name: "user06", Email: "user06@gmail.com", Base: model.Base{ID: "user06"}},
	{Name: "user07", Email: "user07@gmail.com", Base: model.Base{ID: "user07"}},
	{Name: "user08", Email: "user08@gmail.com", Base: model.Base{ID: "user08"}},
	{Name: "user09", Email: "user09@gmail.com", Base: model.Base{ID: "user09"}},
	{Name: "user10", Email: "user10@gmail.com", Base: model.Base{ID: "user10"}},
}

var orders = []*Order{
	// user01 的订单
	{OrderNo: "ORD20240101001", UserID: "user01", Status: "completed", Amount: 199.99, PaymentTime: model.GormTime(time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local)), Base: model.Base{ID: "order01"}},
	{OrderNo: "ORD20240102001", UserID: "user01", Status: "pending", Amount: 299.99, Base: model.Base{ID: "order02"}},
	// user02 的订单
	{OrderNo: "ORD20240101002", UserID: "user02", Status: "completed", Amount: 159.99, PaymentTime: model.GormTime(time.Date(2024, 1, 1, 11, 0, 0, 0, time.Local)), Base: model.Base{ID: "order03"}},
	{OrderNo: "ORD20240102002", UserID: "user02", Status: "cancelled", Amount: 99.99, Base: model.Base{ID: "order04"}},
	// user03 的订单
	{OrderNo: "ORD20240101003", UserID: "user03", Status: "paid", Amount: 399.99, PaymentTime: model.GormTime(time.Date(2024, 1, 1, 14, 0, 0, 0, time.Local)), Base: model.Base{ID: "order05"}},
	{OrderNo: "ORD20240102003", UserID: "user03", Status: "shipped", Amount: 499.99, PaymentTime: model.GormTime(time.Date(2024, 1, 2, 9, 0, 0, 0, time.Local)), Base: model.Base{ID: "order06"}},

	// user04 的订单
	{OrderNo: "ORD20240101004", UserID: "user04", Status: "completed", Amount: 199.99, PaymentTime: model.GormTime(time.Date(2024, 1, 1, 16, 0, 0, 0, time.Local)), Base: model.Base{ID: "order07"}},

	// user05 的订单
	{OrderNo: "ORD20240101005", UserID: "user05", Status: "pending", Amount: 599.99, Base: model.Base{ID: "order08"}},
	{OrderNo: "ORD20240102005", UserID: "user05", Status: "cancelled", Amount: 299.99, Base: model.Base{ID: "order09"}},

	// user06 的订单 (大额订单)
	{OrderNo: "ORD20240101006", UserID: "user06", Status: "completed", Amount: 1299.99, PaymentTime: model.GormTime(time.Date(2024, 1, 1, 15, 0, 0, 0, time.Local)), Base: model.Base{ID: "order10"}},

	// user07 的订单
	{OrderNo: "ORD20240101007", UserID: "user07", Status: "shipped", Amount: 199.99, PaymentTime: model.GormTime(time.Date(2024, 1, 1, 10, 30, 0, 0, time.Local)), Base: model.Base{ID: "order11"}},

	// user08 的订单 (多个订单)
	{OrderNo: "ORD20240101008", UserID: "user08", Status: "completed", Amount: 299.99, PaymentTime: model.GormTime(time.Date(2024, 1, 1, 9, 0, 0, 0, time.Local)), Base: model.Base{ID: "order12"}},
	{OrderNo: "ORD20240102008", UserID: "user08", Status: "completed", Amount: 399.99, PaymentTime: model.GormTime(time.Date(2024, 1, 2, 10, 0, 0, 0, time.Local)), Base: model.Base{ID: "order13"}},
	{OrderNo: "ORD20240103008", UserID: "user08", Status: "pending", Amount: 199.99, Base: model.Base{ID: "order14"}},

	// user09 和 user10 暂时没有订单
}

var (
	configPath    = "/tmp/config.ini"
	configContent = `
[server]
mode = dev
port = 8002
db = "mysql"

[mysql]
database = test
password = qQk5zXWHfj4LD2Nxm9vF3YpBZt8a6JhUTdsS7RgyruGCAEebVP

[logger]
dir = "/tmp/golib/logs"
`
)

var RunOrDie = util.RunOrDie

func init() {
	// // Register table and table records that should automatically created in database.
	// model.Register(users...)
	// model.Register(orders...)
	model.Register[*User]()
	model.Register[*Order]()

	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		panic(err)
	}
	// Bootstrap all initializers.
	config.SetConfigFile(configPath)
	RunOrDie(bootstrap.Bootstrap)
}

// TestWithJoin, 测试 WithJoin 之前请开启 model.Register
func TestWithJoin(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"按时间范围查询订单及用户", testJoinOrdersByTimeRange},
		{"查询订单状态及用户信息", testJoinOrdersWithStatus},
		{"查询大额订单及用户", testJoinHighValueOrders},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}

func testJoinOrdersByTimeRange(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2024, 1, 1, 23, 59, 59, 0, time.Local)
	results := make([]*Order, 0)

	err := database.Database[*Order]().
		WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id").
		WithSelectRaw("orders.order_no, users.name as user_name, orders.payment_time").
		WithTimeRange("payment_time", startTime, endTime).
		WithOrder("payment_time ASC").
		List(&results)

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	// for i := range results {
	// 	fmt.Println(results[i].OrderNo, results[i].UserName, time.Time(results[i].PaymentTime).String())
	// }
}

func testJoinOrdersWithStatus(t *testing.T) {
	results := make([]*Order, 0)
	err := database.Database[*Order]().
		WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id").
		WithSelectRaw("orders.*, users.name as user_name, users.email as user_email").
		WithQuery(&Order{Status: "completed"}).
		List(&results)
	assert.NoError(t, err)
	// for _, o := range results {
	// 	fmt.Println(o)
	// }
}

// 查询大额订单及用户
func testJoinHighValueOrders(t *testing.T) {
	results := make([]*Order, 0)
	err := database.Database[*Order]().
		WithJoinRaw("LEFT JOIN users ON users.id = orders.user_id").
		WithSelectRaw("orders.*, users.name as user_name").
		WithQueryRaw("orders.amount > ?", 500).
		WithOrder("amount DESC").
		List(&results)

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	// 验证所有订单金额都大于 500
	for _, order := range results {
		assert.Greater(t, order.Amount, float64(500))
	}
	// for _, order := range results {
	// 	fmt.Println(order.ID, order.Amount, order.UserName)
	// }
}

func TestWithTransaction(t *testing.T) {
	count := 10
	users := make([]*User, 0)
	for i := 1; i <= count; i++ {
		users = append(users, &User{
			Base: model.Base{ID: strconv.Itoa(i)},
			Name: "user" + strconv.Itoa(i),
		})
	}

	var err error = database.Database[*User]().WithLimit(-1).Update(users...)
	assert.NoError(t, err)

	// var err error = mysql.Default.Transaction(func(tx *gorm.DB) error {
	// 	return database.Database[*User]().WithTransaction(tx).WithLock().WithLimit(-1).Update(users...)
	// })
	assert.NoError(t, err)
}
