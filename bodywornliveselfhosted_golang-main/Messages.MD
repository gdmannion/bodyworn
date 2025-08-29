here all the messages with complete payload sent and received starting from the begining to the end.

Event Received
{
  "bearerId": "c4ca449d-d348-4b8a-84f1-5d846364a3b8",
  "bearerName": "demo",
  "recordingId": "20250506_142153_FB97_B8A44FBB01B6",
  "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
  "subject": "B8A44FBB01B6",
  "time": "2025-05-06T14:22:04Z",
  "type": "com.axis.bodyworn.stream.started"
}
Sent Hello
{
  "type": "hello",
  "id": "c9954960-69d3-4d63-b602-69a840d08382",
  "correlationId": "6f11fac8-feac-4158-b4d8-4159005b5d1c"
}
Sent initSession
{
  "type": "initSession",
  "targetId": "B8A44FBB01B6",
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY1NDE1NzcsIm5iZiI6MTc0NjU0MTI3MiwiaWF0IjoxNzQ2NTQxMjc3fQ._8PCFHU-K5WuyiEZAdpk2fY4PYaq56ZN8kwB-5SUgJI",
  "correlationId": "43cefc4b-1cac-472f-a7ac-2b796c20148a",
  "data": {
    "apiVersion": "1.0",
    "type": "request",
    "method": "initSession",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "context": "26ef95b1-dbc8-47ac-ba58-ba3f3d9d5cf2",
    "params": {
      "type": "live",
      "videoReceive": {},
      "audioReceive": {}
    }
  }
}
Signaling Received
{
  "type": "hello",
  "id": "DDAvhkMlFDtduSvzXMrQbTygJcBN6wMiNrDeAg3k7B19hri2pQ2FaEQqzTsdMqzc::auto_assigned",
  "correlationId": "6f11fac8-feac-4158-b4d8-4159005b5d1c"
}
Signaling Received
{
  "type": "initSession",
  "targetId": "B8A44FBB01B6",
  "correlationId": "43cefc4b-1cac-472f-a7ac-2b796c20148a",
  "turnServers": [
    {
      "urls": [
        "turn://192.168.1.2:3478?transport=tcp",
        "turn://192.168.1.2:3478?transport=udp"
      ],
      "username": "1746627677:axis",
      "password": "HR6KX3j6S1oAjNUu0Y4K+blrqnA="
    }
  ],
  "stunServers": [
    {
      "urls": [
        "stun://192.168.1.2:3478"
      ]
    }
  ],
  "sourcePortRange": {
    "max": 65535,
    "min": 49152
  },
  "data": {
    "apiVersion": "1.0",
    "type": "request",
    "method": "initSession",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "context": "26ef95b1-dbc8-47ac-ba58-ba3f3d9d5cf2",
    "params": {
      "type": "live",
      "videoReceive": {},
      "audioReceive": {}
    }
  }
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "method": "initSession",
    "type": "response",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "context": "26ef95b1-dbc8-47ac-ba58-ba3f3d9d5cf2",
    "id": "1",
    "data": {}
  },
  "correlationId": "43cefc4b-1cac-472f-a7ac-2b796c20148a"
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "type": "request",
    "method": "setSdpOffer",
    "params": {
      "type": "offer",
      "sdp": "v=0\r\no=- 946424643130620735 0 IN IP4 0.0.0.0\r\ns=-\r\nt=0 0\r\na=ice-options:trickle\r\na=group:BUNDLE video0 audio1 application2\r\nm=video 9 UDP/TLS/RTP/SAVPF 96 98\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=ice-ufrag:48U7AfY877a1FDmvEuMhA/y69ZVYbpsm\r\na=ice-pwd:x1ZCATDcMIdniijdlTlqlCBOt4X8Nxfn\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=sendrecv\r\na=rtpmap:96 H264/90000\r\na=rtcp-fb:96 nack\r\na=rtcp-fb:96 nack pli\r\na=rtcp-fb:96 ccm fir\r\na=rtcp-fb:96 transport-cc\r\na=framerate:30\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=fmtp:96 packetization-mode=1;sprop-parameter-sets=Z2QAH6wbGoBQBb/wFuAgICgAAB9AAAdTBM11,aO44MA==;profile-level-id=42E01F;level-asymmetry-allowed=1\r\na=rtpmap:98 rtx/90000\r\na=fmtp:98 apt=96\r\na=ssrc-group:FID 3335252313 104965220\r\na=ssrc:3335252313 msid:user3091403972@host-cdece377 webrtctransceiver0\r\na=ssrc:3335252313 cname:user3091403972@host-cdece377\r\na=ssrc:104965220 msid:user3091403972@host-cdece377 webrtctransceiver0\r\na=ssrc:104965220 cname:user3091403972@host-cdece377\r\na=mid:video0\r\na=fingerprint:sha-256 F9:AC:5A:6A:76:60:7C:8C:2F:6D:9D:4D:1F:F6:AB:9D:AB:E3:1E:5A:F6:AE:41:12:EF:D2:12:49:D4:0C:2A:AC\r\na=rtcp-mux-only\r\nm=audio 0 UDP/TLS/RTP/SAVPF 97\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=ice-ufrag:48U7AfY877a1FDmvEuMhA/y69ZVYbpsm\r\na=ice-pwd:x1ZCATDcMIdniijdlTlqlCBOt4X8Nxfn\r\na=bundle-only\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=sendrecv\r\na=rtpmap:97 OPUS/48000/2\r\na=rtcp-fb:97 transport-cc\r\na=fmtp:97 sprop-stereo=0;sprop-maxcapturerate=48000\r\na=ssrc:1537521888 msid:user3091403972@host-cdece377 webrtctransceiver1\r\na=ssrc:1537521888 cname:user3091403972@host-cdece377\r\na=mid:audio1\r\na=fingerprint:sha-256 F9:AC:5A:6A:76:60:7C:8C:2F:6D:9D:4D:1F:F6:AB:9D:AB:E3:1E:5A:F6:AE:41:12:EF:D2:12:49:D4:0C:2A:AC\r\na=rtcp-mux-only\r\nm=application 0 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:actpass\r\na=ice-ufrag:48U7AfY877a1FDmvEuMhA/y69ZVYbpsm\r\na=ice-pwd:x1ZCATDcMIdniijdlTlqlCBOt4X8Nxfn\r\na=bundle-only\r\na=mid:application2\r\na=sctp-port:5000\r\na=fingerprint:sha-256 F9:AC:5A:6A:76:60:7C:8C:2F:6D:9D:4D:1F:F6:AB:9D:AB:E3:1E:5A:F6:AE:41:12:EF:D2:12:49:D4:0C:2A:AC\r\n"
    }
  },
  "correlationId": ""
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "method": "addIceCandidate",
    "type": "request",
    "params": {
      "candidate": "candidate:1 1 UDP 2015363327 192.168.1.12 46017 typ host",
      "sdpMLineIndex": 0
    },
    "context": "3016533917"
  },
  "correlationId": ""
}
Sent setSdpAnswer
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY1NDE1NzcsIm5iZiI6MTc0NjU0MTI3MiwiaWF0IjoxNzQ2NTQxMjc3fQ._8PCFHU-K5WuyiEZAdpk2fY4PYaq56ZN8kwB-5SUgJI",
  "correlationId": "f846f3e8-2d7a-45b7-95ea-5dcea1043538",
  "data": {
    "apiVersion": "1.0",
    "type": "request",
    "method": "setSdpAnswer",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "context": "26ef95b1-dbc8-47ac-ba58-ba3f3d9d5cf2",
    "params": {
      "type": "answer",
      "sdp": "v=0\r\no=- 7409636142890355322 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE video0 audio1 application2\r\na=msid-semantic: WMS\r\nm=video 9 UDP/TLS/RTP/SAVPF 96 98\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:xePs\r\na=ice-pwd:tJgPapK0AJFQ2q94ckIRMRZw\r\na=ice-options:trickle\r\na=fingerprint:sha-256 FF:18:29:25:42:DE:5F:4D:20:8A:FC:EA:73:98:65:B5:FB:0E:6B:07:D1:D2:E6:5F:46:9E:A0:E7:4D:1F:82:2D\r\na=setup:active\r\na=mid:video0\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=recvonly\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:96 H264/90000\r\na=rtcp-fb:96 transport-cc\r\na=rtcp-fb:96 ccm fir\r\na=rtcp-fb:96 nack\r\na=rtcp-fb:96 nack pli\r\na=fmtp:96 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f\r\na=rtpmap:98 rtx/90000\r\na=fmtp:98 apt=96\r\nm=audio 9 UDP/TLS/RTP/SAVPF 97\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:xePs\r\na=ice-pwd:tJgPapK0AJFQ2q94ckIRMRZw\r\na=ice-options:trickle\r\na=fingerprint:sha-256 FF:18:29:25:42:DE:5F:4D:20:8A:FC:EA:73:98:65:B5:FB:0E:6B:07:D1:D2:E6:5F:46:9E:A0:E7:4D:1F:82:2D\r\na=setup:active\r\na=mid:audio1\r\na=recvonly\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:97 OPUS/48000/2\r\na=rtcp-fb:97 transport-cc\r\na=fmtp:97 minptime=10;useinbandfec=1\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=ice-ufrag:xePs\r\na=ice-pwd:tJgPapK0AJFQ2q94ckIRMRZw\r\na=ice-options:trickle\r\na=fingerprint:sha-256 FF:18:29:25:42:DE:5F:4D:20:8A:FC:EA:73:98:65:B5:FB:0E:6B:07:D1:D2:E6:5F:46:9E:A0:E7:4D:1F:82:2D\r\na=setup:active\r\na=mid:application2\r\na=sctp-port:5000\r\n"
    }
  }
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "method": "addIceCandidate",
    "type": "request",
    "params": {
      "candidate": "candidate:2 1 TCP 1015021823 192.168.1.12 9 typ host tcptype active",
      "sdpMLineIndex": 0
    },
    "context": "3626776687"
  },
  "correlationId": ""
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "method": "addIceCandidate",
    "type": "request",
    "params": {
      "candidate": "candidate:3 1 TCP 1010827519 192.168.1.12 44551 typ host tcptype passive",
      "sdpMLineIndex": 0
    },
    "context": "4202913316"
  },
  "correlationId": ""
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "method": "addIceCandidate",
    "type": "request",
    "params": {
      "candidate": "candidate:4 1 UDP 2015363583 fe80::baa4:4fff:febb:1b6 41973 typ host",
      "sdpMLineIndex": 0
    },
    "context": "3897341148"
  },
  "correlationId": ""
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "method": "addIceCandidate",
    "type": "request",
    "params": {
      "candidate": "candidate:5 1 TCP 1015022079 fe80::baa4:4fff:febb:1b6 9 typ host tcptype active",
      "sdpMLineIndex": 0
    },
    "context": "1514263272"
  },
  "correlationId": ""
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "method": "addIceCandidate",
    "type": "request",
    "params": {
      "candidate": "candidate:6 1 TCP 1010827775 fe80::baa4:4fff:febb:1b6 54537 typ host tcptype passive",
      "sdpMLineIndex": 0
    },
    "context": "659864674"
  },
  "correlationId": ""
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "method": "addIceCandidate",
    "type": "request",
    "params": {
      "candidate": "candidate:10 1 UDP 505413887 192.168.1.2 58059 typ relay raddr 192.168.1.12 rport 46017",
      "sdpMLineIndex": 0
    },
    "context": "4289183868"
  },
  "correlationId": ""
}
Sent addIceCandidate
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY1NDE1NzcsIm5iZiI6MTc0NjU0MTI3MiwiaWF0IjoxNzQ2NTQxMjc3fQ._8PCFHU-K5WuyiEZAdpk2fY4PYaq56ZN8kwB-5SUgJI",
  "correlationId": "0dd65fb7-ab0e-40b8-9b9d-a6d4eaaf0b10",
  "data": {
    "apiVersion": "1.0",
    "type": "request",
    "method": "addIceCandidate",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "context": "26ef95b1-dbc8-47ac-ba58-ba3f3d9d5cf2",
    "params": {
      "candidate": "candidate:743055164 1 udp 2113937151 de0832d0-2ecc-4fc8-b06f-ca210c6aee17.local 60820 typ host generation 0 ufrag xePs network-cost 999",
      "sdpMid": "video0",
      "sdpMLineIndex": 0,
      "usernameFragment": "xePs"
    }
  }
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "method": "setSdpAnswer",
    "type": "response",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "context": "26ef95b1-dbc8-47ac-ba58-ba3f3d9d5cf2",
    "id": "2",
    "data": {}
  },
  "correlationId": "f846f3e8-2d7a-45b7-95ea-5dcea1043538"
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "method": "addIceCandidate",
    "type": "response",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "context": "26ef95b1-dbc8-47ac-ba58-ba3f3d9d5cf2",
    "id": "3",
    "data": {}
  },
  "correlationId": "0dd65fb7-ab0e-40b8-9b9d-a6d4eaaf0b10"
}
Signaling Received
{
  "type": "signaling",
  "targetId": "B8A44FBB01B6",
  "orgId": null,
  "data": {
    "apiVersion": "1.0",
    "sessionId": "7f788353-3b46-4828-a7c6-d85af2c993d8",
    "method": "addIceCandidate",
    "type": "request",
    "params": {
      "candidate": "",
      "sdpMLineIndex": 0
    },
    "context": "5381"
  },
  "correlationId": ""
}