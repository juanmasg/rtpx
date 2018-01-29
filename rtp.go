package main

import (
    "encoding/binary"
    "io"
)

type RTP struct{
    Version             uint8
    Padding             bool
    Extension           bool
    CC                  uint8
    Marker              bool
    PayloadType         uint8
    SequenceNumber      uint16
    Timestamp           uint32
    SSRC                uint32
    CSRC                []uint32
    ExtensionHeader     []*RTPExtension
    ExtensionHeaderID   uint16
    ExtensionHeaderLen  uint16

    Payload         []byte
}

type RTPExtension struct{
    Id                  uint8
    Len                 uint8
    Payload             []byte
}

func RTPPacket(b []byte) *RTP{
    rtp := &RTP{}
    rtp.Version = b[0] & 0xc0 >> 6
    rtp.Padding = b[0] & 0x20 != 0
    rtp.Extension = b[0] & 0x10 != 0
    rtp.CC = b[0] & 0xf
    rtp.Marker = b[1] & 0x80 != 0
    rtp.PayloadType = b[1] & 0x7f
    rtp.SequenceNumber = binary.BigEndian.Uint16(b[2:4])
    rtp.Timestamp = binary.BigEndian.Uint32(b[4:8])
    rtp.SSRC = binary.BigEndian.Uint32(b[8:12])

    offset := 12

    for len(rtp.CSRC) < int(rtp.CC){
        rtp.CSRC = append(rtp.CSRC, binary.BigEndian.Uint32(b[offset:offset+4]))
        offset += 4
    }

    if rtp.Extension{ //RFC8285
        rtp.ExtensionHeaderID = binary.BigEndian.Uint16(b[offset:offset+2])
        rtp.ExtensionHeaderLen = binary.BigEndian.Uint16(b[offset+2:offset+4])

        offset += 4

        if rtp.ExtensionHeaderID == 0xbede{
            for len(rtp.ExtensionHeader) < int(rtp.ExtensionHeaderLen){
                x := &RTPExtension{}
                x.Id = b[offset] & 0xf0 >> 4
                x.Len = (b[offset] & 0x0f) + 1
                offset += 1

                x.Payload = b[offset:offset + int(x.Len)]
                rtp.ExtensionHeader = append(rtp.ExtensionHeader, x)
                offset += 3 //padding to 32bit
            }
        }
    }

    rtp.Payload = b[offset:]

    return rtp
}


func ReadRTP(r io.Reader) *RTP{
    b := make([]byte, 1500)
	n, err := r.Read(b); if err != nil{
        return nil
	}
	rtp := RTPPacket(b[:n])
    return rtp
}

