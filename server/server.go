// Copyright 2016 Jacob Taylor jacob.taylor@gmail.com
// License: Apache2
package main

import (
    "../utils"
    "bufio"
    "bytes"
    "encoding/binary"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net"
    "os"
)

const nbd_folder = "/sample_disks/"

var characters_per_line = 100
var newline = 0
var line_number = 0

func send_export_list_item(output *bufio.Writer, optionSentByClient []byte, export_name string) {
    data := make([]byte, 1024)
    length := len(export_name)
    offset := 0

    // length of export name
    binary.BigEndian.PutUint32(data[offset:], uint32(length)) // length of string
    offset += 4

    // export name
    copy(data[offset:], export_name)
    offset += length

    reply_type := uint32(2) // reply_type: NBD_REP_SERVER
    send_message(output, optionSentByClient, reply_type, uint32(offset), data)
}

func send_ack(output *bufio.Writer, optionSentByClient []byte) {
    send_message(output, optionSentByClient, utils.NBD_COMMAND_ACK, 0, nil)
}

func export_name(output *bufio.Writer, conn net.Conn, payload []byte, pad_with_zeros bool) {
    fmt.Printf("have request to bind to: %s\n", string(payload))

    defer conn.Close()

    // process an export
    var filename bytes.Buffer
    current_directory, err := os.Getwd()
    utils.ErrorCheck(err)
    filename.WriteString(current_directory)
    filename.WriteString(nbd_folder)
    filename.Write(payload)

    fmt.Printf("Opening file: %s\n", filename.String())

    file, err := os.OpenFile(filename.String(), os.O_RDWR, 0644)

    utils.ErrorCheck(err)
    if err != nil {
        return
    }

    buffer := make([]byte, 256)
    //offset := 0

    fs, err := file.Stat()
    file_size := uint64(fs.Size())

    //conn.Write(utils.GetBytes(utils.NBD_REPLY_MAGIC))
    //conn.Write(utils.GetBytes(0))
    //conn.Write(utils.GetBytes(0))
    fmt.Printf("file_size: %d\n", file_size)
    binary.BigEndian.PutUint64(buffer, file_size)
    _, err = output.Write(buffer[:8])
    if err != nil {
        fmt.Printf("Error encountered: %v\n", err)
    }
    binary.BigEndian.PutUint16(buffer, 1)
    _, err = output.Write(buffer[:2])
    if err != nil {
        fmt.Printf("Error encountered: %v\n", err)
    }
    if pad_with_zeros {
        fmt.Printf("padding with zeros\n")
        _, err = output.Write(make([]byte, 124))
        if err != nil {
            fmt.Printf("Error encountered: %v\n", err)
        }
    }
    output.Flush()
    //fileData, err := ioutil.ReadFile(filename.String())
    //conn.Write(fileData)
    // END:  process an export

    // process request message
    conn_reader := bufio.NewReader(conn)
    for {
        clientMagic := make([]byte, utils.NBD_REQUEST_MAGIC_LEN)
        _, err = conn_reader.Read(clientMagic)
        fmt.Printf("Error encountered: %v\n", err)
        if err == io.EOF {
            break
        }
        clientCommandFlags := make([]byte, utils.NBD_CLIENT_COMMAND_FLAGS_LEN)
        conn_reader.Read(clientCommandFlags)
        clientCommandType := make([]byte, utils.NBD_CLIENT_COMMAND_TYPE_LEN)
        conn_reader.Read(clientCommandType)
        clientHandle := make([]byte, utils.NBD_CLIENT_HANDLE_LEN)
        conn_reader.Read(clientHandle)
        clientOffset := make([]byte, utils.NBD_CLIENT_OFFSET_LEN)
        conn_reader.Read(clientOffset)
        clientLength := make([]byte, utils.NBD_CLIENT_LENGTH_LEN)
        conn_reader.Read(clientLength)

        //if binary.BigEndian.Uint32(buffer[4:8])
        fmt.Printf("clientMagic: %d, clientCommandType: %d, command: %d, offset: %d, length: %d\n",
            clientMagic, clientCommandType, clientCommandType, clientOffset, clientLength)
        switch binary.BigEndian.Uint16(clientCommandType) {
        case utils.NBD_COMMAND_READ:
            fmt.Printf(".")
            binary.BigEndian.PutUint32(buffer, utils.NBD_REPLY_MAGIC)
            conn.Write(buffer[:4])
            binary.BigEndian.PutUint32(buffer, 0)
            conn.Write(buffer[:4])
            buffer := make([]byte, binary.BigEndian.Uint32(clientLength))
            file.ReadAt(buffer, int64(binary.BigEndian.Uint64(clientOffset)))
            conn.Write(buffer)
        case utils.NBD_COMMAND_DISCONNECT:
            fmt.Printf("D")
            file.Sync()
            return
        }
    }
    // END: process request message

    //binary.BigEndian.PutUint64(buffer[offset:], file_size) // size
    //offset += 8
    //
    //binary.BigEndian.PutUint16(buffer[offset:], 1) // flags
    //offset += 2
    //
    //if pad_with_zeros {
    //    offset += 124 // pad with 124 zeroes
    //}
    //
    //_, err = output.Write(buffer[:offset])
    ////data_out, err := output.Write(buffer[:offset])
    //
    //output.Flush()
    //utils.ErrorCheck(err)
    ////fmt.Printf("Wrote %d chars: %v\n", data_out, buffer[:offset])
    //
    //buffer = make([]byte, 512*1024)
    //conn_reader := bufio.NewReader(conn)
    //commandBytes := make([]byte, 4)
    //
    //for {
    //    conn_reader.Read(commandBytes)
    //    utils.LogData("commandBytes.  ", len(commandBytes), commandBytes)
    //    command := binary.BigEndian.Uint32(commandBytes)
    //    switch command {
    //    case utils.NBD_COMMAND_READ:
    //        //fmt.Printf("We have a request to read. handle: %v, from: %v, length: %v\n", handle, from, length)
    //        //fmt.Printf("Read Resquest    Offset:%x length: %v     Handle %X\n", from, length, handle)
    //        fmt.Printf(".")
    //
    //        _, err = file.ReadAt(buffer[16:16 + length], int64(from))
    //        utils.ErrorCheck(err)
    //
    //        binary.BigEndian.PutUint32(buffer[:4], utils.NBD_REPLY_MAGIC)
    //        binary.BigEndian.PutUint32(buffer[4:8], 0) // error bits
    //
    //        //utils.LogData("About to reply with", int(16+length), buffer)
    //
    //        conn.Write(buffer[:16 + length])
    //
    //        continue
    //    }
    //}

    //abort := false
    //for {
    //    offset := 0
    //    waiting_for := 28 // wait for at least the minimum payload size
    //    //fmt.Printf("Sitting at top of export loop.  offset: %d, waiting_for: %d\n", offset, waiting_for)
    //    // Duplicate
    //    for offset < waiting_for {
    //        fmt.Printf("offset %d, waiting_for %d", offset, waiting_for)
    //        length, err := conn_reader.Read(buffer[offset:waiting_for])
    //        offset += length
    //        utils.ErrorCheck(err)
    //        if err == io.EOF {
    //            abort = true
    //            break
    //        }
    //        //utils.LogData("Reading instruction\n", offset, buffer)
    //        if offset < waiting_for {
    //            time.Sleep(5 * time.Millisecond)
    //        }
    //    }
    //    // Duplicate
    //    if abort {
    //        fmt.Printf("Abort detected, escaping processing loop\n")
    //        break
    //    }
    //
    //    //fmt.Printf("We read the buffer %v\n", buffer[:waiting_for])
    //
    //    //magic := binary.BigEndian.Uint32(buffer)
    //    command := binary.BigEndian.Uint32(buffer[4:8])
    //    //handle := binary.BigEndian.Uint64(buffer[8:16])
    //    from := binary.BigEndian.Uint64(buffer[16:24])
    //    length := binary.BigEndian.Uint32(buffer[24:28])
    //
    //    newline += 1
    //    if newline%characters_per_line == 0 {
    //        line_number++
    //        fmt.Printf("\n%3d: ", line_number)
    //        newline -= characters_per_line
    //
    //    }
    //
    //    switch command {
    //    case utils.NBD_COMMAND_READ:
    //        //fmt.Printf("We have a request to read. handle: %v, from: %v, length: %v\n", handle, from, length)
    //        //fmt.Printf("Read Resquest    Offset:%x length: %v     Handle %X\n", from, length, handle)
    //        fmt.Printf(".")
    //
    //        _, err = file.ReadAt(buffer[16:16+length], int64(from))
    //        utils.ErrorCheck(err)
    //
    //        binary.BigEndian.PutUint32(buffer[:4], utils.NBD_REPLY_MAGIC)
    //        binary.BigEndian.PutUint32(buffer[4:8], 0) // error bits
    //
    //        //utils.LogData("About to reply with", int(16+length), buffer)
    //
    //        conn.Write(buffer[:16+length])
    //
    //        continue
    //    case utils.NBD_COMMAND_WRITE:
    //        //fmt.Printf("We have a request to write. handle: %v, from: %v, length: %v\n", handle, from, length)
    //        fmt.Printf("W")
    //
    //        waiting_for += int(length) // wait for the additional payload
    //
    //        // Duplicate
    //        for offset < waiting_for {
    //            length, err := conn_reader.Read(buffer[offset:waiting_for])
    //            offset += length
    //            utils.ErrorCheck(err)
    //            if err == io.EOF {
    //                abort = true
    //                break
    //            }
    //            //utils.LogData("Reading write data\n", offset, buffer)
    //            if offset < waiting_for {
    //                time.Sleep(5 * time.Millisecond)
    //            }
    //        }
    //        // Duplicate
    //
    //        _, err = file.WriteAt(buffer[28:28+length], int64(from))
    //        utils.ErrorCheck(err)
    //
    //        file.Sync()
    //
    //        // let them know we are done
    //        binary.BigEndian.PutUint32(buffer[:4], utils.NBD_REPLY_MAGIC)
    //        binary.BigEndian.PutUint32(buffer[4:8], 0) // error bits
    //
    //        //utils.LogData("About to reply with", int(16), buffer)
    //        conn.Write(buffer[:16])
    //
    //        continue
    //
    //    case utils.NBD_COMMAND_DISCONNECT:
    //        fmt.Printf("D")
    //
    //        //fmt.Printf("We have received a request to disconnect\n%d\n", data_out)
    //        // close the file and return
    //
    //        file.Sync()
    //        return
    //    }
    //}
}

func send_export_list(output *bufio.Writer, optionSentByClient []byte, pad_with_zeros bool) {
    current_directory, err := os.Getwd()
    files, err := ioutil.ReadDir(current_directory + nbd_folder)
    if err != nil {
        log.Fatal(err)
    }
    for _, file := range files {
        send_export_list_item(output, optionSentByClient, file.Name())
    }

    if pad_with_zeros {
        buffer := make([]byte, 256)
        _, err = output.Write(buffer[:124])
    }

    send_ack(output, optionSentByClient)
}

func send_message(output *bufio.Writer, optionSentByClient []byte, reply_type uint32, length uint32, data []byte) {
    endian := binary.BigEndian
    buffer := make([]byte, 1024)
    offset := 0

    endian.PutUint64(buffer[offset:], utils.NBD_SERVER_SEND_REPLY_MAGIC)
    offset += 8

    //endian.PutUint32(buffer[offset:], optionSentByClient) // option as sent by client
    copy(buffer[offset:offset+len(optionSentByClient)], optionSentByClient)
    offset += len(optionSentByClient)

    endian.PutUint32(buffer[offset:], reply_type) // reply_type: NBD_REP_SERVER
    offset += 4

    endian.PutUint32(buffer[offset:], length) // length of package
    offset += 4

    if data != nil {
        copy(buffer[offset:], data[0:length])
        offset += int(length)
    }

    data_to_send := buffer[:offset]
    output.Write(data_to_send)
    output.Flush()

    utils.LogData("Just sent:", offset, data_to_send)
}

func padZerosToClient(data []byte) bool {
    optionsFromClient := binary.BigEndian.Uint32(data[:4])
    fmt.Printf("Options are: %v\n", optionsFromClient)
    if (optionsFromClient & utils.NBD_FLAG_FIXED_NEW_STYLE) == utils.NBD_FLAG_FIXED_NEW_STYLE {
        fmt.Printf("Fixed New Style option requested\n")
    }
    if (optionsFromClient & utils.NBD_FLAG_NO_ZEROES) == utils.NBD_FLAG_NO_ZEROES {
        fmt.Printf("No Zero Padding option requested\n")
        return false
    }
    return true
}

func main() {
    if len(os.Args) < 3 {
        panic("missing arguments:  (ipaddress) (portnumber)")
        return
    }

    listener, err := net.Listen("tcp", os.Args[1]+":"+os.Args[2])
    utils.ErrorCheck(err)

    fmt.Printf("Hello World, we have %v\n", listener)
    reply_magic := make([]byte, 4)
    binary.BigEndian.PutUint32(reply_magic, utils.NBD_REPLY_MAGIC)

    for {
        conn, err := listener.Accept()
        utils.ErrorCheck(err)

        fmt.Printf("We have a new connection from: %s\n", conn.RemoteAddr())
        output := bufio.NewWriter(conn)

        // initial handshake from Server to Client
        output.WriteString("NBDMAGIC") // init password
        output.WriteString("IHAVEOPT") // Magic
        output.Write([]byte{0, 3})     // Server side handshake flags, bit0 is NBD_FLAG_FIXED_NEWSTYLE,
        // bit1 is NBD_FLAG_NO_ZEROS
        output.Flush()

        // initial handshake from Client to Server
        data := make([]byte, 1024)
        conn.Read(data[:utils.NBD_HANDSHAKE_CLIENT_FLAG_LEN])
        utils.LogData("handshake from Client to Server. ", utils.NBD_HANDSHAKE_CLIENT_FLAG_LEN, data)
        // END: initial handshake

        // option haggling
        conn.Read(data[:utils.NBD_MAGIC_NUMBER_LEN])
        utils.LogData("option haggling:  client magic number. ", utils.NBD_MAGIC_NUMBER_LEN, data)

        optionSentByClient := make([]byte, utils.NBD_CLIENT_OPTION_LEN)
        conn.Read(optionSentByClient)
        utils.LogData("option haggling:  client option. ", utils.NBD_CLIENT_OPTION_LEN, optionSentByClient)

        padZerosToClient := padZerosToClient(optionSentByClient)

        conn.Read(data[:utils.NBD_CLIENT_LENGTH_OF_OPTION_DATA_LEN])
        lengthOfOptionDataLen := data[:utils.NBD_CLIENT_LENGTH_OF_OPTION_DATA_LEN]
        utils.LogData("option haggling:  lengthOfOptionDataLen. ", utils.NBD_CLIENT_LENGTH_OF_OPTION_DATA_LEN, lengthOfOptionDataLen)

        var optionData []byte
        uintLengthOfOptionDataLen := binary.BigEndian.Uint32(lengthOfOptionDataLen)
        if uintLengthOfOptionDataLen == 0 {
            fmt.Printf("option haggling:  no option data!\n")
        } else {
            fmt.Printf("option haggling:  reading option data!\n")
            optionData = make([]byte, int(uintLengthOfOptionDataLen))
            conn.Read(optionData)
        }
        // END: option haggling

        fmt.Printf("option haggling:  FINISHED %d", binary.BigEndian.Uint32(optionSentByClient))
        padZerosToClient = false
        switch binary.BigEndian.Uint32(optionSentByClient) {
        case utils.NBD_COMMAND_EXPORT_NAME:
            fmt.Printf("command is export name\n")
            go export_name(output, conn, optionData, padZerosToClient)
            break
        case utils.NBD_COMMAND_LIST:
            fmt.Printf("command is list\n")
            send_export_list(output, optionSentByClient, padZerosToClient)
            conn.Close()
            break
        }

        //// Fetch the data until we get the initial options
        //offset := 0
        //waiting_for := 16 // wait for at least the minimum payload size
        //
        //packet_count := 0
        //for offset < waiting_for {
        //    length, err := conn.Read(data[offset:])
        //    if length > 0 {
        //        packet_count += 1
        //    }
        //    offset += length
        //    utils.ErrorCheck(err)
        //    //utils.LogData("Reading instruction", offset, data)
        //    if offset < waiting_for {
        //        time.Sleep(5 * time.Millisecond)
        //    }
        //    // If we are requesting an export, make sure we have the length of the data for the export name.
        //    if offset > 15 && binary.BigEndian.Uint32(data[12:]) == utils.NBD_COMMAND_EXPORT_NAME {
        //        waiting_for = 20
        //    }
        //}
        //
        //fmt.Printf("%d packets processed to get %d bytes\n", packet_count, offset)
        //utils.LogData("Received from client", offset, data)
        //padZerosToClient := padZerosToClient(data)
        //
        //command := binary.BigEndian.Uint32(data[12:])
        //payload_size := int(binary.BigEndian.Uint32(data[16:]))
        //
        //fmt.Sprintf("command is: %d\npayload_size is: %d\n", command, payload_size)
        //waiting_for += int(payload_size)
        //for offset < waiting_for {
        //    length, err := conn.Read(data[offset:])
        //    offset += length
        //    utils.ErrorCheck(err)
        //    utils.LogData("Reading instruction", offset, data)
        //    if offset < waiting_for {
        //        time.Sleep(5 * time.Millisecond)
        //    }
        //}
        //payload := make([]byte, payload_size)
        //
        //if payload_size > 0 {
        //    copy(payload, data[20:])
        //}
        //
        //utils.LogData("Payload is:", payload_size, payload)
        //fmt.Printf("command is: %v\n", command)

        // At this point, we have the command, payload size, and payload.
        //switch command {
        //case utils.NBD_COMMAND_LIST:
        //    send_export_list(output, optionSentByClient, padZerosToClient)
        //    conn.Close()
        //    break
        //case utils.NBD_COMMAND_EXPORT_NAME:
        //    export_name(output, conn, payload_size, payload, padZerosToClient)
        //    break
        //}
    }

}
