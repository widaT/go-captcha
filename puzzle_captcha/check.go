package captcha

import "errors"

const slipOffset = 5.0

var (
	ErrPostionErr = errors.New("postion error")
)

// Check 验证位置是否正确
func Check(paramInPoint *Point, cachedPoint *Point) error {
	if cachedPoint.X-slipOffset > paramInPoint.X ||
		paramInPoint.X > cachedPoint.X+slipOffset {
		return ErrPostionErr
	}
	return nil
}
