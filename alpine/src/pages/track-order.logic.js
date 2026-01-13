import { trackOrder } from "../lib/api.js";
import { initRevealObserver } from "../lib/dom-utils.js";

export function TrackOrderLogic() {
  return {
    order: null,
    orderNumber: "",
    phone: "",
    email: "",
    errors: {},
    loading: false,
    error: null,

    init() {
      this.initRevealObserver();
    },

    initRevealObserver() {
      initRevealObserver();
    },

    validate() {
      this.errors = {};
      if (!this.orderNumber || !String(this.orderNumber).trim())
        this.errors.orderNumber = "Vui lòng nhập mã đơn hàng";
      if (!this.phone || !String(this.phone).trim())
        this.errors.phone = "Vui lòng nhập số điện thoại";
      return Object.keys(this.errors).length === 0;
    },

    async submitForm() {
      if (!this.validate()) return;
      this.loading = true;
      this.error = null;
      try {
        const res = await trackOrder({
          orderNumber: this.orderNumber.trim(),
          phone: this.phone.trim(),
          email: this.email?.trim() || "",
        });
        this.order = res.data || res.order || null;
      } catch (e) {
        this.error = e.message || "Không thể tra cứu đơn hàng";
      } finally {
        this.loading = false;
      }
    },

    getStatusLabel(status) {
      const map = {
        pending: "Chờ xử lý",
        shipped: "Đang vận chuyển",
        delivered: "Đã giao",
        cancelled: "Đã hủy",
      };
      return map[status] || status || "Unknown";
    },

    getStatusIcon(status) {
      if (status === "delivered") return "check-circle";
      if (status === "shipped") return "truck";
      if (status === "pending") return "clock";
      return "clock";
    },

    getStatusColor(status) {
      if (status === "delivered") return "text-success";
      if (status === "shipped") return "text-blue-400";
      return "text-muted-foreground";
    },

    getStatusSteps(status) {
      // Simplified steps
      const steps = [
        { label: "Đặt hàng", completed: true },
        {
          label: "Xử lý",
          completed: status !== "pending",
          active: status === "pending",
        },
        {
          label: "Vận chuyển",
          completed: status === "shipped" || status === "delivered",
          active: status === "shipped",
        },
        {
          label: "Giao hàng",
          completed: status === "delivered",
          active: status === "delivered",
        },
      ];
      return steps;
    },

    resetForm() {
      this.order = null;
      this.orderNumber = "";
      this.phone = "";
      this.email = "";
      this.errors = {};
      this.error = null;
    },
  };
}
