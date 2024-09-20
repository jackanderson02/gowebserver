package mock
import (
	"acme/model"
)

type MockRepository struct {
	MockGetUsers    func() ([]model.User, error)
	MockGetUser     func(id int) (model.User, error)
	MockAddUser     func(user model.User) (int, error)
	MockUpdateUser  func(id int, user model.User)  error
	MockDeleteUser  func(id int) error
	MockClose       func()
}

func (m *MockRepository) GetUsers() ([]model.User, error) {
	return m.MockGetUsers()
}
func (m *MockRepository) GetUser(id int) (model.User, error) {
	return m.MockGetUser(id)
}
func (m *MockRepository) AddUser(user model.User) (int, error) {
	return m.MockAddUser(user)
}
func (m *MockRepository) UpdateUser(id int, user model.User) error {
	return m.MockUpdateUser(id,user)
}
func (m *MockRepository) DeleteUser(id int) error {
	return m.MockDeleteUser(id)
}
func (m *MockRepository) Close() {
	m.MockClose()
}