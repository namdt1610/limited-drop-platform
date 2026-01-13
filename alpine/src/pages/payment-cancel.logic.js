import { initRevealObserver } from "../lib/dom-utils.js";

export function PaymentCancelLogic() {
  return {
    orderCode: "",
    dropId: "",
    code: "",
    isCancelled: true,

    init() {
      this.initRevealObserver();
      const urlParams = new URLSearchParams(window.location.search);
      this.dropId = urlParams.get("drop_id") || "";
      this.code = urlParams.get("code") || "";
      this.orderCode = urlParams.get("orderCode") || "";
    },

    initRevealObserver() {
      initRevealObserver();
    },

    goHome() {
      window.location.href = "/";
    },
    retryPayment() {
      if (this.dropId)
        window.location.href = "/drop?id=" + encodeURIComponent(this.dropId);
    },
  };
}
