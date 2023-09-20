/*
 * Project: Application Utility Library
 * Filename: /json.go
 * Created Date: Wednesday September 20th 2023 15:20:52 +0800
 * Author: Sallehuddin Abdul Latif (sallehuddin@berrypay.com)
 * Company: BerryPay (M) Sdn. Bhd.
 * --------------------------------------
 * Last Modified: Wednesday September 20th 2023 15:22:03 +0800
 * Modified By: Sallehuddin Abdul Latif (sallehuddin@berrypay.com)
 * --------------------------------------
 * Copyright (c) 2023 BerryPay (M) Sdn. Bhd.
 */

package apputil

func IsEmptyJson(data interface{}) bool {
	switch val := data.(type) {
	case map[string]interface{}:
		return len(val) == 0
	case []interface{}:
		return len(val) == 0
	default:
		return true // Data is not a JSON object or array
	}
}
