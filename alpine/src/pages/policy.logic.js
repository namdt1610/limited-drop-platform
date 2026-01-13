import { initRevealObserver } from "../lib/dom-utils.js";

export function PolicyLogic() {
  return {
    init() {
      this.initRevealObserver();
    },
    initRevealObserver() {
      initRevealObserver();
    },
    goBack() {
      window.history.back();
    },
  };
}
