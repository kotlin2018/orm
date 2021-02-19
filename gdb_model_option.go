// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package orm

// Option adds extra operation option for the model.
func (m *Model) Option(option int) *Model {
	model := m.getModel()
	model.option = model.option | option
	return model
}

// OmitEmpty sets OPTION_OMITEMPTY option for the model, which automatically filers
// the data and where attributes for empty values.
func (m *Model) OmitEmpty() *Model {
	return m.Option(OPTION_OMITEMPTY)
}
