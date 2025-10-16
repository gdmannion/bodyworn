import { startEventFeed } from './events.js';
import { fetchToken } from './webrtc_video.js';

document.addEventListener("DOMContentLoaded", async () => {
  await fetchToken();  // 🔁 Get fresh token on page load
  startEventFeed();
});
