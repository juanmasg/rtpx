package rtp

import (
    "rtpx/rtp"
    "testing"
)

func BenchmarkRTPParseFull(b *testing.B){
    payl := make([]byte, 12)
    for n := 0; n< b.N; n++{
        rtp.RTPPacket(payl)
    }
}

func BenchmarkRTPParseQuick(b *testing.B){
    payl := make([]byte, 12)
    for n := 0; n< b.N; n++{
        rtp.RTPQuickPayload(payl)
    }
}
