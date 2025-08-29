// webrtc_video.js
export const remoteVideo = document.getElementById("remoteVideo");
export const statusEl = document.getElementById("status");
export const errorModal = document.getElementById("errorModal");
export const errorText = document.getElementById("errorText");
export const eventContainer = document.getElementById("events");

let token = null;
let targetId = null;
let ws = null;
let pc = null;
let pendingCandidates = [];

export function generateUUID() {
  if (window.crypto && crypto.randomUUID) return crypto.randomUUID();
  return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}

export function updateStatus(text) {
  statusEl.textContent = text;
}

export function showErrorModal(msg) {
  errorText.textContent = msg;
  errorModal.style.display = "block";
}

export function logEventToDOM(data, label = "") {
  const div = document.createElement("div");
  div.className = "event-card";
  div.textContent = `${label}\n${JSON.stringify(data, null, 2)}`;
  eventContainer.appendChild(div); 
}

async function fetchToken() {
  const res = await fetch("/token");
  const data = await res.json();
  token = data.token;
  targetId = data.targetId;
  console.log("Token Received:", token);
  console.log("Target ID:", targetId);
}

export async function startWebRTC(sessionId, contextId) {
  await fetchToken();

  const wsProtocol = location.protocol === "https:" ? "wss:" : "ws:";
  const wsUrl = `${wsProtocol}//${location.host}/ws-proxy?token=${encodeURIComponent(token)}`;
  ws = new WebSocket(wsUrl);

  pc = new RTCPeerConnection();
  pendingCandidates = [];

  ws.onopen = () => {
    const hello = {
      type: "hello",
      id: generateUUID(),
      correlationId: generateUUID()
    };
    ws.send(JSON.stringify(hello));
    logEventToDOM(hello, "Sent Hello");

    const initSession = {
      type: "initSession",
      targetId,
      accessToken: token,
      correlationId: generateUUID(),
      data: {
        apiVersion: "1.0",
        type: "request",
        method: "initSession",
        sessionId,
        context: contextId,
        params: {
          type: "live",
          videoReceive: {},
          audioReceive: {}
        }
      }
    };
    ws.send(JSON.stringify(initSession));
    logEventToDOM(initSession, "Sent initSession");
  };

  ws.onerror = (err) => {
    console.error("WebRTC signaling error:", err);
    updateStatus("Signaling error");
    showErrorModal("WebRTC signaling WebSocket error. Check device or backend.");
  };

  ws.onclose = () => {
    console.warn("WebRTC signaling closed.");
    updateStatus("Signaling closed");
  };

  ws.onmessage = async (event) => {
    const msg = JSON.parse(event.data);
    logEventToDOM(msg, "Signaling Received");
    const { method, params } = msg?.data || {};

    if (method === "setSdpOffer") {
      await pc.setRemoteDescription(new RTCSessionDescription(params));
      const answer = await pc.createAnswer();
      await pc.setLocalDescription(answer);

      const sdpAnswer = {
        type: "signaling",
        targetId,
        accessToken: token,
        correlationId: generateUUID(),
        data: {
          apiVersion: "1.0",
          type: "request",
          method: "setSdpAnswer",
          sessionId,
          context: contextId,
          params: {
            type: answer.type,
            sdp: answer.sdp
          }
        }
      };
      ws.send(JSON.stringify(sdpAnswer));
      logEventToDOM(sdpAnswer, "Sent setSdpAnswer");
    }

    if (method === "addIceCandidate") {
      if (!params || params.candidate === '') {
        console.log("ICE gathering complete (end-of-candidates)");
        return;
      }

      if (params.candidate && (params.sdpMid != null || params.sdpMLineIndex != null)) {
        try {
          const candidate = new RTCIceCandidate(params);
          if (pc.remoteDescription) {
            await pc.addIceCandidate(candidate);
          } else {
            pendingCandidates.push(candidate);
          }
        } catch (err) {
          console.warn("Failed to add ICE candidate:", err, params);
        }
      }
    }
  };

  pc.onicecandidate = (event) => {
    if (event.candidate) {
      const candidateMsg = {
        type: "signaling",
        targetId,
        accessToken: token,
        correlationId: generateUUID(),
        data: {
          apiVersion: "1.0",
          type: "request",
          method: "addIceCandidate",
          sessionId,
          context: contextId,
          params: event.candidate.toJSON()
        }
      };
      ws.send(JSON.stringify(candidateMsg));
      logEventToDOM(candidateMsg, "Sent addIceCandidate");
    }
  };

  pc.ontrack = (event) => {
    remoteVideo.srcObject = event.streams[0];
  };
}
