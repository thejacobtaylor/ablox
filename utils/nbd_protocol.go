package utils


type TransmissionRequestMessage struct {
	//C: 32 bits, 0x25609513, magic (NBD_REQUEST_MAGIC)
	//C: 16 bits, command flags
	//C: 16 bits, type
	//C: 64 bits, handle
	//C: 64 bits, offset (unsigned)
	//C: 32 bits, length (unsigned)
	//C: (length bytes of data if the request is of type NBD_CMD_WRITE)
	magic		[]byte
	commandFlags	[]byte
	typeFlags	[]byte
	handle		[]byte
	offset		[]byte
	length		[]byte
	data		[]byte
}

//func (trm *TransmissionRequestMessage) parseFromReader(reader *bufio.Reader) {
//	trm := TransmissionRequestMessage{nil, nil, nil, nil, nil, nil, nil}
//	data := make([]byte, 1024)
//	reader.Read(data[:NBD_MAGIC_NUMBER_LEN])
//
//
//}