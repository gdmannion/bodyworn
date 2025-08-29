// events.js
import { generateUUID, updateStatus, showErrorModal, logEventToDOM, startWebRTC } from './webrtc_video.js';

const eventContainer = document.getElementById("events");
const statusEl = document.getElementById("status");
const errorModal = document.getElementById("errorModal");
const errorText = document.getElementById("errorText");

let wsEvents = null;
let sessionId = null;
let contextId = null;

export function getSessionContext() {
  return { sessionId, contextId };
}

export function startEventFeed() {
  const wsProtocol = location.protocol === "https:" ? "wss" : "ws";
  wsEvents = new WebSocket(`${wsProtocol}://${location.host}/events`);

  wsEvents.onopen = () => {
    console.log("Events WebSocket connected");
    updateStatus("Connected to events");
  };

  wsEvents.onerror = (err) => {
    console.error("Event feed WebSocket error", err);
    updateStatus("Event feed error");
    showErrorModal("Event WebSocket connection failed. Check backend server and network.");
  };

  wsEvents.onclose = () => {
    console.warn("Event feed WebSocket closed.");
  };

  wsEvents.onmessage = (event) => {
    const evt = JSON.parse(event.data);
    logEventToDOM(evt, "Event Received");

    if (evt.type === "com.axis.bodyworn.stream.started") {
      sessionId = evt.sessionId;
      contextId = generateUUID();
      console.log("Stream Start Detected, session:", sessionId);
      updateStatus("Starting WebRTC");
      startWebRTC(sessionId, contextId);
    }
  };
}
