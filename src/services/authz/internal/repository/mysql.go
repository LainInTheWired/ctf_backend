package repository

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u *userRepository) CreatePermission(permission *model.Permission) error {
	ins, err := u.DB.Prepare("INSERT INTO permissions (name,description)  VALUES(? ,?)")
	if err != nil {
		return errors.Wrap(err, "user insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(permission.Name, permission.Description)
	if err != nil {
		return errors.Wrap(err, "can't insert permissions")
	}
	return nil
}
func (u *userRepository) CreateRole(role *model.Role) (int, error) {
	ins, err := u.DB.Prepare("INSERT INTO roles (name,namespace)  VALUES(? ,?)")
	if err != nil {
		return 0, errors.Wrap(err, "user insert error")
	}
	defer ins.Close()

	r, err := ins.Exec(role.Name, role.Namespace)
	if err != nil {
		return 0, errors.Wrap(err, "can't insert roles")
	}
	lastInsertID, err := r.LastInsertId()
	if err != nil {
		// エラーハンドリング
		return 0, errors.Wrap(err, "failed to retrieve last insert ID")
	}

	return int(lastInsertID), nil
}
func (u *userRepository) BindRolePermissions(rid int, pid int) error {
	ins, err := u.DB.Prepare("INSERT INTO role_permissions ( role_id, permission_id)  VALUES(? ,?)")
	if err != nil {
		return errors.Wrap(err, "user insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(rid, pid)
	if err != nil {
		return errors.Wrap(err, "can't bind role_permissions")
	}
	return nil
}
func (u *userRepository) BindUserRoles(uid int, rid int) error {
	ins, err := u.DB.Prepare("INSERT INTO user_roles (user_id,role_id) VALUES(?,?)")
	if err != nil {
		return errors.Wrap(err, "can't bind user_roles")
	}
	defer ins.Close()

	_, err = ins.Exec(uid, rid)
	if err != nil {
		return errors.Wrap(err, "can't bind user_roles")
	}
	return nil
}

func (u *userRepository) GetUserWithRolesAndPermissions(id int) (*model.User, error) {
	query := `
    SELECT
        u.id 
        u.name 
        r.id 
        r.name 
        p.id 
        p.name 
    FROM
        users AS u
        JOIN user_roles AS ur ON u.id = ur.user_id
        JOIN roles AS r ON ur.role_id = r.id
        JOIN role_permissions AS rp ON r.id = rp.role_id
        JOIN permissions AS p ON rp.permission_id = p.id
    WHERE
        u.id = $1
    ORDER BY
        r.id, p.id;
    `
	if _, err := u.GetUserByID(id); err != nil {
		return nil, errors.Wrap(err, "can't get user by id")
	}
	user := &model.User{}
	// rows, err := u.DB.Query("SELECT r.id ,r.name,r.namespace FROM roles AS r INNER JOIN user_roles ON r.id = user_roles.role_id WHERE user_roles.user_id = ?", id)
	rows, err := u.DB.Query(query, id)
	if err != nil {
		return nil, errors.Wrap(err, "can't select user_roles")
	}
	defer rows.Close()

	for rows.Next() {
		role := model.Role{}
		if err := rows.Scan(&role.ID, &role.Name, &role.Namespace); err != nil {
			return nil, errors.Wrap(err, "can't scan role")
		}
		user.Role = append(user.Role, role)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "can't scan role")
	}

	return nil, nil
}
