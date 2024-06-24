package ai

// SystemRoles 表名:roles，存放系统角色
type SystemRoles struct {
	Role string `gorm:"column:role"`
	Desc string `gorm:"column:desc"`
}
