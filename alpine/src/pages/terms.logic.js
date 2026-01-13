import { initRevealObserver } from "../lib/dom-utils.js";

export function TermsLogic() {
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
