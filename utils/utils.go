// Copyright 2016 Jacob Taylor jacob.taylor@gmail.com
// License: Apache2
package utils

import (
    "fmt"
    "encoding/binary"
    "os"
)

const (
    NBD_REQUEST_MAGIC_LEN =             4
    NBD_REQUEST_MAGIC =                 uint32(0x25609513)
    NBD_REPLY_MAGIC =                   uint32(0x67446698)
    NBD_SERVER_SEND_REPLY_MAGIC =       uint64(0x3e889045565a9)
    NBD_COMMAND_ACK =                   uint32(1)

    NBD_COMMAND_READ =                  uint16(0)
    NBD_COMMAND_WRITE =                 uint16(1)
    NBD_COMMAND_DISCONNECT =            uint16(2)
    NBD_COMMAND_FLUSH =                 uint16(3)
    NBD_COMMAND_TRIM =                  uint16(4)

    NBD_COMMAND_EXPORT_NAME =           uint32(1)
    NBD_COMMAND_LIST =                  uint32(3)

    NBD_FLAG_FIXED_NEW_STYLE =          uint32(1)
    NBD_FLAG_NO_ZEROES =                uint32(2)

    NBD_HANDSHAKE_SERVER_FLAG_LEN =     2
    NBD_HANDSHAKE_CLIENT_FLAG_LEN =     4
    NBD_MAGIC_NUMBER_LEN =              8
    NBD_CLIENT_OPTION_LEN =             4
    NBD_CLIENT_LENGTH_OF_OPTION_DATA_LEN = 4

    NBD_CLIENT_COMMAND_FLAGS_LEN =      2
    NBD_CLIENT_COMMAND_TYPE_LEN =       2
    NBD_CLIENT_HANDLE_LEN =             8
    NBD_CLIENT_OFFSET_LEN =             8
    NBD_CLIENT_LENGTH_LEN =             4
)

func ErrorCheck(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encountered: %v\n", err)
	}
}

func LogData(msg string, count int, data []byte) {
    fmt.Printf("%5s (count %3d) Data: '%s' (%v)\n", msg, count, string(data[0:count]), data[0:count])
}

func EncodeInt(val int) []byte {
    data := make([]byte, 4)
    binary.BigEndian.PutUint32(data, uint32(val))
    return data
}

