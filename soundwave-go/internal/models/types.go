package models

// Role 用户角色类型
type Role string

const (
	RoleAdmin  Role = "admin"  // 管理员
	RoleUser   Role = "user"   // 普通用户
	RoleTester Role = "tester" // 测试用户
)

// Permission 权限类型
type Permission string

const (
	PermissionViewServices Permission = "view_services" // 查看服务列表
	PermissionViewStats    Permission = "view_stats"    // 查看服务统计
	PermissionManageSystem Permission = "manage_system" // 管理系统设置
	PermissionManageUsers  Permission = "manage_users"  // 管理用户
)
