package user

// 送信可能ポイントオブジェクト
type SendablePoint struct {
	value int
}

func NewSendablePoint(sendablePoint int) *SendablePoint {
	return &SendablePoint{value: sendablePoint}
}

// 送付可能かどうか判断する
//
// @params
// sendPlannedPoint 送付予定ポイント
//
// @returns
// true: 送付できる, false: 送付できない
func (sp *SendablePoint) CanSendPoint(sendPlannedPoint int) bool {
	return sp.value-sendPlannedPoint >= 0
}

// ポイント残高を計算
//
// @params
// sendPlannedPoint 送付予定ポイント
//
// @returns
// 残高
func (sp *SendablePoint) CalculatePointBalance(sendPlannedPoint int) int {
	return sp.value - sendPlannedPoint
}
