/*
 * Project: Application Error Library
 * Filename: /uid.go
 * Created Date: Sunday September 3rd 2023 18:25:10 +0800
 * Author: Sallehuddin Abdul Latif (sallehuddin@berrypay.com)
 * Company: BerryPay (M) Sdn. Bhd.
 * --------------------------------------
 * Last Modified: Sunday September 3rd 2023 18:25:38 +0800
 * Modified By: Sallehuddin Abdul Latif (sallehuddin@berrypay.com)
 * --------------------------------------
 * Copyright (c) 2023 BerryPay (M) Sdn. Bhd.
 */

package apputil

import (
	"time"

	"github.com/muyo/sno"
)

var meta byte

func init() {
	meta = 88
}

func SetIDMeta(b byte) {
	meta = b
}

func GetIDMeta() byte {
	return meta
}

func GenerateRequestIDByte() []byte {
	return sno.NewWithTime(meta, time.Now()).Bytes()
}

func GenerateRequestID() sno.ID {
	return sno.NewWithTime(meta, time.Now())
}

func GenerateRequestIDString() string {
	return sno.NewWithTime(meta, time.Now()).String()
}
