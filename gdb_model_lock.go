// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package orm

// LockUpdate 在当前操作上为 Model 设置"FOR UPDATE"锁。
//
// 避免选择行被其它共享锁修改或删除。
func (m *Model) LockUpdate() *Model {
	model := m.getModel()
	model.lockInfo = "FOR UPDATE"
	return model
}

// LockShared  在当前操作上为 Model 设置"共享锁"。
//
// 共享锁可以避免被选择的行被修改直到事务提交。
func (m *Model) LockShared() *Model {
	model := m.getModel()
	model.lockInfo = "LOCK IN SHARE MODE"
	return model
}
