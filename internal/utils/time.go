/*
 *  Copyright (c) 2022 NetEase Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
 * Project: DingoCli
 * Created Date: 2022-08-29
 * Author: chengyi (Cyber-SiKu)
 */

package cobrautil

import (
	"syscall"

	cmderror "github.com/dingodb/dingofs-tools/internal/error"
)

 func GetTimeofDayUs() (uint64, *cmderror.CmdError) {
	var now syscall.Timeval
	if e := syscall.Gettimeofday(&now); e != nil {
		retErr := cmderror.ErrGettimeofday()
		retErr.Format(e.Error())
		return uint64(0), retErr
	}
	return uint64(now.Sec) * 1000000 + uint64(now.Usec), cmderror.Success()
 }
 