package repo

//EmailExistError struct
type EmailExistError struct{}

func (*EmailExistError) Error() string {
	return "email exists"
}

//EmailNotExistsError error email not exsist
type EmailNotExistsError struct{}

func (*EmailNotExistsError) Error() string {
	return "email not exists"
}

type PasswordError struct{}

func (*PasswordError) Error() string {
	return "Password does not match"
}

type RecordNotFoundError struct{}

func (*RecordNotFoundError) Error() string {
	return "Record not found"
}
