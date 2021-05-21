package ps

import (
	"errors"
)

var ErrStackUnderFlow = errors.New("スタックが枯渇した状態でデータ呼び出し発生")

var ErrEOF = errors.New("スクリプト解析中にファイル末尾まで到達")

var ErrInvalidToken = errors.New("無効なトークン")
