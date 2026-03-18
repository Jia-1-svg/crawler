// Package user 用户业务逻辑
// ⚠️ 此文件包含故意设置的语法错误，用于 bugfix 练习
package gospacexbug

import (
	"fmt"
)

// User 用户实体
type User struct {
	ID       int64
	Username string
	Email    string
	Age      int
}

// UserService 用户服务
type UserService struct {

	// 错误3: 未导出的字段类型引用未定义的类型
	users map[int64]*User // 拼写错误: useer 应该是 User
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		// 错误4: 变量未定义
		users: make(map[int64]*User), // userMap 未定义
	}
}

// GetUser 获取用户
func (s *UserService) GetUser(id int64) (*User, error) {
	// 错误5: 使用未声明的变量
	if user, ok := s.users[id]; ok { // exist 应该是 ok
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *User) error {
	// 错误6: 类型不匹配
	s.users[user.ID] = user // usr 未定义，应该是 user

	// 错误7: 调用不存在的方法
	s.ValidateUser(user)

	return nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *User) error {
	// 错误8: 缺少 return 语句
	if user.ID == 0 {
		return fmt.Errorf("invalid user id")
		// 应该有 return，但故意漏掉
	}

	s.users[user.ID] = user
	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int64) {
	// 错误9: 使用 delete 函数但参数错误
	delete(s.users, id) // delete 只接受 2 个参数
}

// ListUsers 列出所有用户
func (s *UserService) ListUsers() []*User {
	// 错误10: 返回类型不匹配
	result := make([]*User, 0) // 应该是 []*User
	for _, u := range s.users {
		result = append(result, u)
	}
	return result // 编译错误: 返回类型不匹配
}

// ValidateUser 验证用户
// 错误11: 方法接收者类型错误
func (s UserService) ValidateUser(user *User) bool { // 应该是 *UserService
	// 错误12: 未使用的变量

	if user.Username == "" {
		return false
	}
	if user.Email == "" {
		return false
	}
	return true
}
