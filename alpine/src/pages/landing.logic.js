import { getDrops, getDropStatus } from "../lib/api.js";
import {
  initRevealObserver,
  handleMagnetic,
  resetMagneticStyle,
  formatCountdown,
} from "../lib/dom-utils.js";

export function LandingLogic() {
  return {
    activeDropId: null,
    status: "CONNECTING...",
    magneticStyle: "",
    cursorX: 0,
    cursorY: 0,
    isHovering: false,
    countdown: { d: "00", h: "00", m: "00", s: "00" },
    nextDropTime: null,
    serverOffset: 0,
    isDropLive: false,
    dropStartTime: null,

    async init() {
      try {
        const drops = await getDrops();
        const now = Date.now();

        // Map backend response to expected format
        const mappedDrops = (drops || []).map((d) => ({
          drop_id: d.id,
          available: (d.total_stock || 0) - (d.sold || 0),
          next_drop_at: d.starts_at,
        }));

        const candidates = mappedDrops.filter((d) => {
          if (d.available > 0) return true;
          if (d.next_drop_at) {
            const t = new Date(d.next_drop_at).getTime();
            return t > now;
          }
          return false;
        });

        const chosen = (candidates[0] ?? mappedDrops[0]) || null;
        this.activeDropId = chosen?.drop_id ?? null;
        this.nextDropTime = chosen?.next_drop_at
          ? new Date(chosen.next_drop_at).getTime()
          : null;
        this.dropStartTime = this.nextDropTime;

        // Get server time offset for accurate countdown
        if (this.activeDropId) {
          try {
            const dropStatus = await getDropStatus(this.activeDropId);
            const serverNow = new Date(dropStatus.now).getTime();
            const clientNow = Date.now();
            this.serverOffset = serverNow - clientNow;

            // Check if drop is already live
            this.checkIfDropLive();
          } catch (e) {
            // Sync failed, will use client time
          }
        }

        this.status = "SYSTEM READY";

        if (this.nextDropTime) {
          this.startCountdown();
        }
      } catch (e) {
        this.status = "OFFLINE";
      }

      // Listen for hover on interactive elements
      document.addEventListener("mouseover", (e) => {
        if (
          e.target.closest('a, button, .group, [x-data="techSpecsSection"] div')
        ) {
          this.isHovering = true;
        } else {
          this.isHovering = false;
        }
      });

      // Initialize Intersection Observer for reveal animations
      this.initRevealObserver();
    },

    updateCursor(e) {
      this.cursorX = e.clientX;
      this.cursorY = e.clientY;
    },

    handleMagnetic(e) {
      const result = handleMagnetic(e.currentTarget, e);
      // Move button 30% towards cursor instead of 20%
      this.magneticStyle = `transform: translate(${
        result.transform.split("translate(")[1].split(")")[0]
      }, 0.3))`;
    },

    resetMagnetic() {
      this.magneticStyle =
        "transform: translate(0, 0); transition: transform 0.5s cubic-bezier(0.23, 1, 0.32, 1)";
    },

    checkIfDropLive() {
      if (!this.dropStartTime) return;
      const now = Date.now() + this.serverOffset;
      this.isDropLive = now >= this.dropStartTime;
    },

    startCountdown() {
      const update = () => {
        this.checkIfDropLive();
        const now = Date.now() + this.serverOffset;
        const diff = this.nextDropTime - now;
        this.countdown =
          diff <= 0
            ? { d: "00", h: "00", m: "00", s: "00" }
            : formatCountdown(diff);
        if (diff > 0) requestAnimationFrame(update);
      };
      update();
    },

    initRevealObserver() {
      initRevealObserver();
    },

    canEnter() {
      return !!this.activeDropId;
    },
  };
}
