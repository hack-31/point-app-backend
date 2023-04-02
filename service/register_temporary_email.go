package service

type RegisterTemporaryEmail struct{}

func NewRegisterTemporaryEmail() *RegisterTemporaryEmail {
	return &RegisterTemporaryEmail{}
}

// メール仮登録
//
//
func (r *RegisterTemporaryEmail) RegisterTemporaryEmail() {
	println("testService")
}
