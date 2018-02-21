package readers

import(
    "time"
    "net"
    "fmt"
    "../../rtpx/rtp"
)

type RTPReader struct{
    mcastr      MulticastReader
}

func NewRTPReader(ifacename string, addr string, port int) (r *RTPReader, err error){
    mcastr := NewMulticastReader(ifacename, addr, port)
    return &RTPReader{mcastr}, nil
}

func (rtpr *RTPReader) ReadFrom() (*rtp.RTP, string, error){
    b := make([]byte, 1500)
    n, cm, _, err := rtpr.mcastr.ReadFrom(b); if err != nil{
        return nil, "", err
    }

    rtp := rtp.RTPPacket(b[:n])
    return rtp, cm.Dst.String(), nil
}

func (rtpr *RTPReader) ReadPayloadFrom() (uint16, []byte, string, error){
    b := make([]byte, 1500)
    n, cm, _, err := rtpr.mcastr.ReadFrom(b); if err != nil{
        return 0, nil, "", err
    }

    seq, pay := rtp.RTPQuickPayload(b[:n])
    return seq, pay, cm.Dst.String(), nil
}

func (rtpr *RTPReader) String() string{
    return fmt.Sprintf("RTP/%s", rtpr.mcastr)
}

func (rtpr *RTPReader) Close() (err error){
    return rtpr.mcastr.Close()
}

func (rtpr *RTPReader) JoinGroup(group net.IP) error{
    return rtpr.mcastr.JoinGroup(group)
}

func (rtpr *RTPReader) LeaveGroup(group net.IP) error{
    return rtpr.mcastr.LeaveGroup(group)
}

func (rtpr *RTPReader) SetReadDeadline(t time.Time) error{
    return rtpr.mcastr.SetReadDeadline(t)
}
