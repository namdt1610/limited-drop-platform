import { verifySymbicode } from "../lib/api.js";
import { initRevealObserver } from "../lib/dom-utils.js";

export function VerifyLogic() {
  return {
    code: "",
    verifying: false,
    result: null,

    init() {
      this.initRevealObserver();
      this.setCodeFromParams();
    },

    initRevealObserver() {
      initRevealObserver();
    },

    setCodeFromParams() {
      const params = window.AlpineRouter?.getParams?.() || {};
      if (params.code) this.code = params.code.toUpperCase();
    },

    async handleVerify() {
      if (!this.code.trim()) {
        this.result = { success: false, message: "Vui lòng nhập mã SYMBICODE" };
        return;
      }
      this.verifying = true;
      this.result = null;
      try {
        const data = await verifySymbicode(this.code.trim());
        this.result = data;
      } catch (error) {
        this.result = {
          success: false,
          message: error.message || "Lỗi hệ thống. Vui lòng thử lại.",
        };
      } finally {
        this.verifying = false;
      }
    },

    // handleScan() {
    //   TODO: Implement QR Scanner
    // },
    reset() {
      this.code = "";
      this.result = null;
    },
  };
}
