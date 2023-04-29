package service

type UpdateEmail struct{}

func NewUpdateEmail() *UpdateEmail {
	return &UpdateEmail{}
}

// メール本変更サービス
func (ue *UpdateEmail) UpdateEmail() {
	println("サービステスト")
}
