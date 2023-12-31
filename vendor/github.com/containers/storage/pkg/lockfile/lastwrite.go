package lockfile

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/binary"
	"os"
	"sync/atomic"
	"time"
)

// LastWrite is an opaque identifier of the last write to some *LockFile.
// It can be used by users of a *LockFile to determine if the lock indicates changes
// since the last check.
//
// Never construct a LastWrite manually; only accept it from *LockFile methods, and pass it back.
type LastWrite struct {
	// Never modify fields of a LastWrite object; it has value semantics.
	state []byte // Contents of the lock file.
}

var lastWriterIDCounter uint64 // Private state for newLastWriterID

const lastWriterIDSize = 64 // This must be the same as len(stringid.GenerateRandomID)
// newLastWrite returns a new "last write" ID.
// The value must be different on every call, and also differ from values
// generated by other processes.
func newLastWrite() LastWrite {
	// The ID is (PID, time, per-process counter, random)
	// PID + time represents both a unique process across reboots,
	// and a specific time within the process; the per-process counter
	// is an extra safeguard for in-process concurrency.
	// The random part disambiguates across process namespaces
	// (where PID values might collide), serves as a general-purpose
	// extra safety, _and_ is used to pad the output to lastWriterIDSize,
	// because other versions of this code exist and they don't work
	// efficiently if the size of the value changes.
	pid := os.Getpid()
	tm := time.Now().UnixNano()
	counter := atomic.AddUint64(&lastWriterIDCounter, 1)

	res := make([]byte, lastWriterIDSize)
	binary.LittleEndian.PutUint64(res[0:8], uint64(tm))
	binary.LittleEndian.PutUint64(res[8:16], counter)
	binary.LittleEndian.PutUint32(res[16:20], uint32(pid))
	if n, err := cryptorand.Read(res[20:lastWriterIDSize]); err != nil || n != lastWriterIDSize-20 {
		panic(err) // This shouldn't happen
	}

	return LastWrite{
		state: res,
	}
}

// serialize returns bytes to write to the lock file to represent the specified write.
func (lw LastWrite) serialize() []byte {
	if lw.state == nil {
		panic("LastWrite.serialize on an uninitialized object")
	}
	return lw.state
}

// Equals returns true if lw matches other
func (lw LastWrite) equals(other LastWrite) bool {
	if lw.state == nil {
		panic("LastWrite.equals on an uninitialized object")
	}
	if other.state == nil {
		panic("LastWrite.equals with an uninitialized counterparty")
	}
	return bytes.Equal(lw.state, other.state)
}

// newLastWriteFromData returns a LastWrite corresponding to data that came from a previous LastWrite.serialize
func newLastWriteFromData(serialized []byte) LastWrite {
	if serialized == nil {
		panic("newLastWriteFromData with nil data")
	}
	return LastWrite{
		state: serialized,
	}
}
