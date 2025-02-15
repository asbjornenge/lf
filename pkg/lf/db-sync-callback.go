/*
 * LF: Global Fully Replicated Key/Value Store
 * Copyright (C) 2018-2019  ZeroTier, Inc.  https://www.zerotier.com/
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 *
 * --
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial closed-source software that incorporates or links
 * directly against ZeroTier software without disclosing the source code
 * of your own application.
 */

package lf

// Callbacks called from C have to be in a separate file due to cgo linking weirdness.

//#include <stdint.h>
import "C"
import "unsafe"

//export ztlfSyncCCallback
func ztlfSyncCCallback(dbp unsafe.Pointer, hash unsafe.Pointer, doff C.uint64_t, dlen C.uint, reputation C.int, arg unsafe.Pointer) {
	globalSyncCallbacksLock.RLock()
	defer func() {
		_ = recover() // should not happen since this is caught elsewhere, but make non-fatal
		globalSyncCallbacksLock.RUnlock()
	}()
	idx := int(uintptr(arg) & 0x7fffffff)
	if idx < len(globalSyncCallbacks) && globalSyncCallbacks[idx] != nil {
		var hash2 [32]byte
		copy(hash2[:], ((*[32]byte)(hash))[:]) // make copy since this buffer changes between callbacks
		globalSyncCallbacks[idx](uint64(doff), uint(dlen), int(reputation), &hash2)
	}
}
