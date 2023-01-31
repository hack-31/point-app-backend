package constant

const (
	// 確認コードの有効期限
	ConfirmationCodeExpiration_m int = 60
	// デフォルトの送信可能ポイント数
	DefaultSendingPoint int = 1000
	// 確認コードの桁数
	ConfirmCodeLength int = 4
	// キャッシュ側のトークンの有効期限
	TokenExpiration_m int = 60
	// JWT側の最大有効期限（連続操作この指定時間経てば期限が切れる）
	MaxTokenExpiration_m int = 3600
	// ランダムパスワードの桁数
	RandomPasswordLength int = 12
)
