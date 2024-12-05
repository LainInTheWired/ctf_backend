package repository

import (
	"encoding/json"
	"testing"
)

// func TestGetUserWithRolesAndPermissions(t *testing.T) {
// 	r, err := NewDBClient()
// 	if err != nil {
// 		t.Errorf("%v", err)
// 	}
// 	ur := NewUserRepository(r)
// 	a, err := ur.GetUserWithRolesAndPermissions(1)
// 	if err != nil {
// 		t.Errorf("%v", err)
// 	}
// 	t.Logf("%+v", a)
// }

func TestGetPermissions(t *testing.T) {
	r, err := NewDBClient()
	if err != nil {
		t.Errorf("%v", err)
	}
	ur := NewUserRepository(r)
	a, err := ur.GetRolePermissions()
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("%+v", a)
}

func TestGetUserWithRole(t *testing.T) {
	r, err := NewDBClient()
	if err != nil {
		t.Errorf("%v", err)
	}
	ur := NewUserRepository(r)
	a, err := ur.GetUserWithRoles(1)
	if err != nil {
		t.Errorf("%v", err)
	}
	ja, err := json.Marshal(a)
	t.Logf("%s", ja)
}
