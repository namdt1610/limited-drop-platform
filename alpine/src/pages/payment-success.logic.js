import { initRevealObserver } from "../lib/dom-utils.js";

export function PaymentSuccessLogic() {
  return {
    orderCode: "",
    status: "",
    dropId: "",
    code: "",
    isPaid: false,

    init() {
      this.initRevealObserver();
      const urlParams = new URLSearchParams(window.location.search);
      this.dropId = urlParams.get("drop_id") || "";
      this.code = urlParams.get("code") || "";
      this.status = urlParams.get("status") || "";
      const cancel = urlParams.get("cancel") || "false";
      this.orderCode = urlParams.get("orderCode") || "";
      this.isPaid = this.status === "PAID" && cancel === "false";
    },

    initRevealObserver() {
      initRevealObserver();
    },

    goHome() {
      window.location.href = "/";
    },
    contactSupport() {
      window.location.href = "mailto:support@donaldvibe.xyz";
    },
  };
}
